package concat

import (
	_ "embed"
	"strings"
	"testing"
)

//go:embed testdata/data
var testArr []byte
var testStr = strings.Fields(string(testArr))

func BenchmarkConcat(b *testing.B) {
	for i := 0; i < b.N; i++ {
		concat(testStr)
	}

}

func BenchmarkStringsBuilder(b *testing.B) {
	for i := 0; i < b.N; i++ {
		concatStringsBuilder(testStr)
	}
}

func BenchmarkByteArr(b *testing.B) {
	for i := 0; i < b.N; i++ {
		concatByteArr(testStr)
	}
}

func BenchmarkByteArrCopy(b *testing.B) {
	for i := 0; i < b.N; i++ {
		concatByteArrCopy(testStr)
	}
}
