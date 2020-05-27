package gerror

import (
	"testing"
)

func Benchmark_Error(B *testing.B) {
	B.ReportAllocs()
	B.ResetTimer()
	for i := 0; i < B.N; i++ {
		Error(-1002, "ddfdfdfdfd", "haodddfd")
	}
}

func Benchmark_New(B *testing.B) {
	B.ReportAllocs()
	B.ResetTimer()
	for i := 0; i < B.N; i++ {
		New(-1002, "ddfdfdfdfd", "haodddfd")
	}

}
func Benchmark_String(B *testing.B) {
	B.ReportAllocs()
	B.ResetTimer()
	var a string
	for i := 0; i < B.N; i++ {
		test_r(-1002, "ddfdfdfdfd", "haodddfd")
	}
	if a == "" {
	}
}
