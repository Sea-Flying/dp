package utils

import "testing"

func TestMergeStringSlice(t *testing.T) {
	a := []string{"Tom", "Jack"}
	b := []string{"Jack", "Tom", "Rose"}
	a = MergeStringSlice(a, b)
	for _, i := range a {
		println(i)
	}
}
