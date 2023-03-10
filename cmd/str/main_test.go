package main

import "testing"

var testStrings = []string{
	"foo", "bar", "baz", "foo", "bar", "baz", "foo", "bar", "baz", "foo", "bar", "baz", "foo", "bar", "baz", "foo", "bar", "baz", "foo", "bar", "baz", "foo", "bar", "baz", "foo", "bar", "baz", "foo", "bar", "baz", "foo", "bar", "baz", "foo", "bar", "baz", "foo", "bar", "baz", "foo", "bar", "baz", "foo", "bar", "baz", "foo", "bar", "baz", "foo", "bar", "baz", "foo", "bar", "baz", "foo", "bar", "baz", "foo", "bar", "baz", "foo", "bar", "baz", "foo", "bar", "baz", "foo", "bar", "baz", "foo", "bar", "baz",
}

func BenchmarkConcat(b *testing.B) {
	for i := 0; i < b.N; i++ {
		concat(testStrings)
	}
}

func BenchmarkConcatNu(b *testing.B) {
	for i := 0; i < b.N; i++ {
		concatNu(testStrings)
	}
}

func BenchmarkConcatNuNu(b *testing.B) {
	for i := 0; i < b.N; i++ {
		concatNuNu(testStrings)
	}
}
