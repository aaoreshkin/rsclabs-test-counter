package common

import (
	"strconv"
	"testing"
)

func BenchmarkIntConversion_Atoi(b *testing.B) {
	s := "12345"

	for b.Loop() {
		strconv.Atoi(s)
	}
}

func BenchmarkIntConversion_Manual(b *testing.B) {
	s := "12345"

	for b.Loop() {
		result := 0
		for j := 0; j < len(s); j++ {
			result = result*10 + int(s[j]-'0')
		}
	}
}
