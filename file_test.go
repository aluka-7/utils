package utils

import (
	"testing"
)

func TestIsFile(t *testing.T) {
	if !IsFile("file.go") {
		t.Errorf("IsExist:\n Expect => %v\n Got => %v\n", true, false)
	}

	if IsFile("testdata") {
		t.Errorf("IsExist:\n Expect => %v\n Got => %v\n", false, true)
	}

	if IsFile("files.go") {
		t.Errorf("IsExist:\n Expect => %v\n Got => %v\n", false, true)
	}
}

func BenchmarkIsFile(b *testing.B) {
	for i := 0; i < b.N; i++ {
		IsFile("file.go")
	}
}

func BenchmarkIsExist(b *testing.B) {
	for i := 0; i < b.N; i++ {
		IsExist("file.go")
	}
}

func BenchmarkIsDir(b *testing.B) {
	for i := 0; i < b.N; i++ {
		IsDir("file.go")
	}
}
