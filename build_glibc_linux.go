//go:build !dynamic && !musl

package cre2

//#cgo CFLAGS: -I${SRCDIR}/include/. -O3
//#cgo LDFLAGS: -L${SRCDIR}/lib/libre2_glibc_linux.a -lre2 -lstdc++ -lpthread -O3
import "C"
import "log"

func init() {
	log.Println("[cre2-go] load glibc_linux")
}
