package benchmark

import (
	"strings"
	"testing"
)

func BenchmarkPlusOperator(b *testing.B) {
	for i := 0; i < b.N; i++ { // <= 를 < 로 수정 권장
		str := ""
		for j := 0; j < 100; j++ {
			str += "hello"
		}
	}
}

func BenchmarkStringBuilder(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var builder strings.Builder
		for j := 0; j < 100; j++ {
			builder.WriteString("hello")
		}
		_ = builder.String()
	}
}
