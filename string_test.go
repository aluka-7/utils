package utils

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

var sMap *sync.Map

// TestTokenizeToStringArray
func TestTokenizeToStringArray(t *testing.T) {
	tokens := TokenizeToStringArray("/bla /**/**/bla/", "/", false, false)
	for _, item := range tokens {
		fmt.Println(*item)
	}
}

func TestIsBlank(t *testing.T) {
	t.Log(IsBlank(""))
	t.Log(IsBlank(" "))
	t.Log(IsBlank("		"))
	t.Log(IsBlank(" t"))
	t.Log(IsBlank("t "))
	t.Log(IsBlank("t t"))
	t.Log(IsBlank("tt"))

}

func init() {
	sMap = new(sync.Map)
}

// 测试sync.Map的性能
// TestSyncMap 62050000 466883000
func TestSyncMap(t *testing.T) {
	t1 := time.Now().Nanosecond()
	for i := 0; i < 65536; i++ {
		sMap.Store(i, i)
	}
	timeSpan := time.Now().Nanosecond() - t1
	t.Log(timeSpan)
}

func BenchmarkAppendStr(b *testing.B) {
	s := []string{"a"}
	for i := 0; i < b.N; i++ {
		s = AppendStr(s, fmt.Sprint(b.N%3))
	}
}
