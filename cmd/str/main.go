package main

import "strings"

func concat(str []string) string {
	result := ""

	for _, v := range str { // loop over a copy
		result += v //  new string allocation
	}

	return result
}

func concatStringsBuilder(str []string) string {
	var b strings.Builder

	var n int
	for i := range str {
		n += len(str[i])
	}

	b.Grow(n)

	for i := range str {
		b.WriteString(str[i])
	}

	return b.String()
}

func concatByteArr(str []string) string {
	var n int
	for i := range str {
		n += len(str[i])
	}

	arr := make([]byte, 0, n)

	n = 0
	for i := range str {
		arr = append(arr, str[i]...)
	}

	return string(arr)
}

func concatByteArrCopy(str []string) string {
	var n int
	for i := range str {
		n += len(str[i])
	}

	arr := make([]byte, n)

	n = 0
	for i := range str {
		n += copy(arr[n:], str[i])
	}

	return string(arr)
}
