package main

import "strings"

func concat(str []string) string {
	var result string // smells less than result := ""

	for _, v := range str {
		result += v // new string allocating on the heap
	}

	return result
}

func concatNu(str []string) string {
	return strings.Join(str, "")
}

func concatNuNu(str []string) string {
	var b strings.Builder

	for _, s := range str {
		b.WriteString(s)
	}

	return b.String()
}
