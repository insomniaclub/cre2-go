package nocopy

import "testing"

func TestBytesToString(t *testing.T) {
	target := "abcdefg"
	b := []byte(target)
	s := BytesToString(b)
	if s != target {
		t.Errorf("expected %s, but got %s", target, s)
		return
	}
	t.Log(s)
}

func TestStringToBytes(t *testing.T) {
	target := "abcdefg"
	b := StringToBytes(target)
	s := string(b)
	if s != target {
		t.Errorf("expected %s, but got %s", target, s)
		return
	}
	t.Log(s)
}

func BenchmarkBytesToStringCopy(b *testing.B) {
	target := make([]byte, 10*1024)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = string(target)
	}
}

func BenchmarkBytesToStringNoCopy(b *testing.B) {
	target := make([]byte, 10*1024)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = BytesToString(target)
	}
}

func BenchmarkStringToBytesCopy(b *testing.B) {
	target := string(make([]byte, 10*1024))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = []byte(target)
	}
}

func BenchmarkStringToBytesNoCopy(b *testing.B) {
	target := string(make([]byte, 10*1024))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = StringToBytes(target)
	}
}
