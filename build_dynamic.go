//go:build dynamic

package cre2

//#cgo pkg-config: re2
import "C"
import "log"

func init() {
	log.Println("[cre2-go] load dynamic")
}
