package nocopy

import "unsafe"

// BytesToString convert []byte to string via unsafe.Pointer
func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// StringToBytes convert string to []byte via unsafe.Pointer,
//
// WARNING: if the string s is allocated in constant pool, do not change the returned []byte value
func StringToBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(&s))
}
