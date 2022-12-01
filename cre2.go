package cre2

//#cgo LDFLAGS: -lre2 -lstdc++ -O3
//#cgo CFLAGS: -I${SRCDIR}/. -O3
//#include <stdlib.h>
//#include "cre2.h"
//#include "cre2_cgo.h"
import "C"
import (
	"fmt"
	"unsafe"

	"github.com/gensliu/nocopy"
)

type unsafeptr = unsafe.Pointer

// Regexp is the representation of a compiled regular expression.
// A Regexp is safe for concurrent use by multiple goroutines,
// except for configuration methods, such as Longest.
type Regexp struct {
	opt    unsafeptr //	*C.cre2_options_t
	rex    unsafeptr // *C.cre2_regexp_t
	nGroup int       // num of capturing groups
}

// Compile parses a regular expression and returns, if successful,
// a Regexp object that can be used to match against text.
func Compile(s string) (*Regexp, error) {
	pattern := *(*C.cre2_string_t)(unsafeptr(&s))

	opt := C.cre2_opt_new()
	rex := C.cre2_new(pattern.data, C.int(pattern.length), opt)

	if errCode := C.cre2_error_code(rex); errCode != C.CRE2_NO_ERROR {
		errMsg := C.GoString(C.cre2_error_string(rex))
		return nil, fmt.Errorf("cre2: Compile(`%s`): error parsing regexp: %s", s, errMsg)
	}

	return &Regexp{opt, rex, int(C.cre2_num_capturing_groups(rex))}, nil
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

// Match reports whether the byte slice b
// contains any match of the regular expression r.
func (r *Regexp) Match(b []byte) bool {
	cstr := (*C.cre2_string_t)(unsafe.Pointer(&b))
	return bool(C.match(r.rex, cstr.data, cstr.length))
}

// MatchString reports whether the string s
// contains any match of the regular expression r.
func (r *Regexp) MatchString(s string) bool {
	// cstr := (*C.cre2_string_t)(unsafe.Pointer(&s))
	// return bool(C.match(r.rex, cstr.data, cstr.length))
	return r.Match(nocopy.StringToBytes(s))
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
	if n == 0 {
		return nil
	}

	cstr := (*C.cre2_string_t)(unsafe.Pointer(&s))
	if n < 0 {
		n = int(cstr.length) + 1
	}

	matched := make([]string, n)
	len := C.all_matches(
		/* regexp    */ r.rex,
		/* textaddr  */ cstr.data,
		/* textlen   */ cstr.length,
		/* match     */ (*C.cre2_string_t)(unsafe.Pointer(&matched[0])),
		/* nmatch    */ C.int(n),
		/* nsubmatch */ 1,
	)
	if len == 0 {
		return nil
	}

	return matched[:len]
}

// FindAllStringSubmatch it returns a slice of all successive matches of the expression,
// as defined by the 'All' description in the package comment.
// A return value of nil indicates no match.
func (r *Regexp) FindAllStringSubmatch(s string, n int) [][]string {
	// debug.SetGCPercent(-1)
	if n == 0 {
		return nil
	}

	cstr := (*C.cre2_string_t)(unsafe.Pointer(&s))
	if n < 0 {
		n = int(cstr.length) + 1
	}

	// NOTE: the following code will cause a gc panic, because the memory of [][]string is not continuous.
	// 	match := make([][]string, n)
	// 	for i := range match {
	// 		match[i] = make([]string, re.nGroup+1)
	// 	}
	//	C.find_all_string_submatch(..., (**C.cre2_string_t)(unsafe.Pointer(&match[0][0])),...)
	rawMatch := make([]string, n*(r.nGroup+1))
	len := C.all_matches(
		/* regexp    */ r.rex,
		/* textaddr  */ cstr.data,
		/* textlen   */ cstr.length,
		/* match     */ (*C.cre2_string_t)(unsafe.Pointer(&rawMatch[0])),
		/* nmatch    */ C.int(n),
		/* nsubmatch */ C.int(r.nGroup+1),
	)

	if len == 0 {
		return nil
	}

	rawMatch = rawMatch[:int(len)*(r.nGroup+1)]
	match := make([][]string, len)
	for i := 0; i < int(len); i++ {
		match[i] = rawMatch[i*(r.nGroup+1) : (i+1)*(r.nGroup+1)]
	}

	return match
}

// FindAllStringIndex is the 'All' version of FindStringIndex; it returns a
// slice of all successive matches of the expression, as defined by the 'All'
// description in the package comment.
// A return value of nil indicates no match.
func (r *Regexp) FindAllStringIndex(s string, n int) [][]int {
	if n == 0 {
		return nil
	}

	cstr := (*C.cre2_string_t)(unsafe.Pointer(&s))
	if n < 0 {
		n = int(cstr.length) + 1
	}

	rawMatch := make([]C.int, n*2)
	len := C.find_all_string_index(
		/* regexp   */ r.rex,
		/* textaddr */ cstr.data,
		/* textlen  */ cstr.length,
		/* match    */ (**C.int)(unsafe.Pointer(&rawMatch[0])),
		/* nmatch   */ C.int(n),
	)

	if len == 0 {
		return nil
	}

	rawMatch = rawMatch[:len*2]
	match := make([][]int, len)
	for i := 0; i < int(len); i++ {
		// match[i] = rawMatch[i*2 : i*2+1]
		match[i] = []int{int(rawMatch[i*2]), int(rawMatch[i*2+1])}
	}

	return match
}
