package cre2

//#cgo CXXFLAGS: -std=c++17 -O3 -g
//#cgo pkg-config: re2
//#include <stdlib.h>
//#include "cre2.h"
//#include "cre2_cgo.h"
import "C"
import (
	"fmt"
	"sync"
	"unsafe"
)

const (
	bufferSize = 100000
	sliceSize  = bufferSize / 2
)

type cBuffer struct {
	b    []C.int
	lock sync.Mutex
}

func newBuffer() cBuffer {
	return cBuffer{
		b:    make([]C.int, bufferSize),
		lock: sync.Mutex{},
	}
}

type unsafeptr = unsafe.Pointer

// Regexp is the representation of a compiled regular expression.
// A Regexp is safe for concurrent use by multiple goroutines,
// except for configuration methods, such as Longest.
type Regexp struct {
	expr   string
	opt    unsafeptr //	*C.cre2_options_t
	rex    unsafeptr // *C.cre2_regexp_t
	nGroup int       // num of capturing groups
	buffer cBuffer
}

// Compile parses a regular expression and returns, if successful,
// a Regexp object that can be used to match against text.
func Compile(s string) (*Regexp, error) {
	pattern := *(*C.cre2_string_t)(unsafeptr(&s))

	opt := C.cre2_opt_new()
	C.cre2_opt_set_max_mem(opt, 50<<20) //	50MB, default is 8MB
	rex := C.cre2_new(pattern.data, C.int(pattern.length), opt)

	if errCode := C.cre2_error_code(rex); errCode != C.CRE2_NO_ERROR {
		errMsg := C.GoString(C.cre2_error_string(rex))
		return nil, fmt.Errorf("cre2: Compile(`%s`): error parsing regexp: %s", s, errMsg)
	}

	return &Regexp{s, opt, rex, int(C.cre2_num_capturing_groups(rex)), newBuffer()}, nil
}

// MustCompile is like Compile but panics if the expression cannot be parsed.
// It simplifies safe initialization of global variables holding compiled regular
// expressions.
func MustCompile(s string) *Regexp {
	r, err := Compile(s)
	if err != nil {
		panic(err)
	}
	return r
}

func (r *Regexp) String() string {
	return r.expr
}

// Close will free the members of the Regexp, which was allocated in C heap.
func (r *Regexp) Close() {
	if r.rex != nil {
		C.cre2_delete(r.rex)
		r.rex = nil
	}
	if r.opt != nil {
		C.cre2_opt_delete(r.opt)
		r.opt = nil
	}
}

// MatchString reports whether the string s
// contains any match of the regular expression r.
func (r *Regexp) MatchString(s string) bool {
	cstr := (*C.cre2_string_t)(unsafeptr(&s))
	return bool(C.match(r.rex, cstr.data, cstr.length))
}

// FindString returns a string holding the text of the leftmost match in s of the regular
// expression. If there is no match, the return value is an empty string,
// but it will also be empty if the regular expression successfully matches
// an empty string.
func (r *Regexp) FindString(s string) string {
	match := r.FindAllString(s, 1)
	if len(match) == 0 {
		return ""
	}
	return match[0]
}

// FindAllString returns a slice of all successive matches of the expression,
// as defined by the 'All' description in the package comment.
// A return value of nil indicates no match.
func (r *Regexp) FindAllString(s string, n int) []string {
	groups := r.FindAllStringIndex(s, n)
	if groups == nil {
		return nil
	}

	matched := make([]string, 0, len(groups))
	for _, group := range groups {
		matched = append(matched, s[group[0]:group[1]])
	}
	return matched
}

func (r *Regexp) FindStringSubmatch(s string) []string {
	match := r.FindAllStringSubmatch(s, 1)
	if len(match) == 0 {
		return nil
	}
	return match[0]
}

// FindAllStringSubmatch it returns a slice of all successive matches of the expression,
// as defined by the 'All' description in the package comment.
// A return value of nil indicates no match.
func (r *Regexp) FindAllStringSubmatch(s string, n int) [][]string {
	groups := r.FindAllStringSubmatchIndex(s, n)
	if groups == nil {
		return nil
	}

	nGroup := r.nGroup + 1

	matched := make([][]string, 0, len(groups))
	for _, group := range groups {
		tmp := make([]string, 0, nGroup)
		for i := 0; i < nGroup; i++ {
			tmp = append(tmp, s[group[i*2]:group[i*2+1]])
		}
		matched = append(matched, tmp)
	}
	return matched
}

// FindAllStringIndex is the 'All' version of FindStringIndex; it returns a
// slice of all successive matches of the expression, as defined by the 'All'
// description in the package comment.
// A return value of nil indicates no match.
func (r *Regexp) FindAllStringIndex(s string, n int) [][]int {
	if n == 0 {
		return nil
	}

	cstr := (*C.cre2_string_t)(unsafeptr(&s))
	if n < 0 {
		n = int(cstr.length) + 1
	}
	if n > sliceSize {
		n = sliceSize
	}

	r.buffer.lock.Lock()
	defer r.buffer.lock.Unlock()

	len := C.find_all_string_index(
		/* regexp   */ r.rex,
		/* textaddr */ cstr.data,
		/* textlen  */ cstr.length,
		/* match    */ (**C.int)(unsafeptr(&r.buffer.b[0])),
		/* nmatch   */ C.int(n),
	)

	if len == 0 {
		return nil
	}

	matched := make([][]int, 0, len)
	for i := 0; i < int(len); i++ {
		matched = append(matched, []int{int(r.buffer.b[i*2]), int(r.buffer.b[i*2+1])})
	}
	return matched
}

func (r *Regexp) FindAllStringSubmatchIndex(s string, n int) [][]int {
	if n == 0 {
		return nil
	}

	cstr := (*C.cre2_string_t)(unsafe.Pointer(&s))
	nGroup := r.nGroup + 1
	maxN := sliceSize / nGroup
	if n < 0 {
		n = (int(cstr.length) + 1) / nGroup
	}
	if n > maxN {
		n = maxN
	}

	r.buffer.lock.Lock()
	defer r.buffer.lock.Unlock()

	len := C.find_all_string_submatch_index(
		/* regexp   */ r.rex,
		/* textaddr */ cstr.data,
		/* textlen  */ cstr.length,
		/* match    */ (**C.int)(unsafe.Pointer(&r.buffer.b[0])),
		/* nmatch   */ C.int(n),
		/* nsubmatch   */ C.int(nGroup),
	)

	if len == 0 {
		return nil
	}

	matched := make([][]int, 0, len)
	for i := 0; i < int(len); i++ {
		tmp := make([]int, 0, nGroup)
		for j := 0; j < nGroup; j++ {
			tmp = append(
				tmp,
				int(r.buffer.b[i*nGroup*2+j*2]),
				int(r.buffer.b[i*nGroup*2+j*2+1]),
			)
		}
		matched = append(matched, tmp)
	}

	return matched
}
