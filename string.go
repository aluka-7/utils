package utils

import (
	"reflect"
	"strings"
	"unicode"
	"unicode/utf8"
	"unsafe"
)

// EmptySpace
const EmptySpace = " "

// EmptyString
const EmptyString = ""

// 根据分隔符进行分割处理,形成包路径数组.默认分割符为:",; \t\n"
func TokenizeToStringArray(str, delimiters string, trimTokens, ignoreEmptyTokens bool) []*string {
	if str == EmptyString {
		return nil
	}
	tokens := make([]*string, 0)
	for _, token := range strings.Split(str, delimiters) {
		if trimTokens {
			token = strings.Trim(token, EmptySpace)
		}
		if !ignoreEmptyTokens || token != EmptyString {
			var item = token
			tokens = append(tokens, &item)
		}
	}
	return tokens
}

// TokenizeToStringArray1
func TokenizeToStringArray1(str, delimiters string) []*string {
	return TokenizeToStringArray(str, delimiters, true, true)
}

// Str2Bytes 字符串转[]byte
func Str2Bytes(s string) []byte {
	x := *(*reflect.StringHeader)(unsafe.Pointer(&s))
	h := *(*[]byte)(unsafe.Pointer(&x))
	return h
}

// Bytes2Str []byte转字符串
func Bytes2Str(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// StartsWith
func StartsWith(str, prefix string, offset int) bool {
	ta := Str2Bytes(str)
	to := offset
	pa := Str2Bytes(prefix)
	po := 0
	pc := utf8.RuneCountInString(prefix)
	// 注意:偏移量可能接近 -1>>>1.
	if (offset < 0) || (offset > utf8.RuneCountInString(str)-pc) {
		return false
	}
	for {
		if pc--; pc >= 0 {
			if ta[to] != pa[po] {
				to++
				po++
				return false
			}
		} else {
			break
		}
	}
	return true
}

// IsBlank 判断是否存在空格
func IsBlank(source string) bool {
	if strings.EqualFold(EmptyString, source) {
		return true
	}
	for i := len(source); i > 0; {
		r, size := utf8.DecodeLastRuneInString(source[0:i])
		i -= size
		if !unicode.IsSpace(r) {
			return false
		}
	}
	return true
}

// HasText 判断是否有值
func HasText(source string) bool {
	return !IsBlank(source)
}

// AppendUniqueStr 将字符串追加到数组中,且去重
func AppendUniqueStr(strs []string, str string) []string {
	for _, s := range strs {
		if s == str {
			return strs
		}
	}
	return append(strs, str)
}

// Reverse 字符翻转
func Reverse(s []rune) []rune {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}
