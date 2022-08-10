package utils

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"time"
	"unsafe"
)

// 浮点数除法
func Div(n, b float64) float64 {
	return math.Log(n) / math.Log(b)
}

// PowInt is int type of math.Pow function.
func PowInt(x int, y int) int {
	if y <= 0 {
		return 1
	} else {
		if y%2 == 0 {
			sqrt := PowInt(x, y/2)
			return sqrt * sqrt
		} else {
			return PowInt(x, y-1) * x
		}
	}
}

func RandInt(start int, end int) int {
	// 范围检查
	if end < start {
		return 0
	}
	// 随机数生成器，加入时间戳保证每次生成的随机数不一样
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	// 生成随机数
	return r.Intn(end-start) + start
}

func RandInt64(start int64, end int64) int64 {
	// 范围检查
	if end < start {
		return 0
	}
	// 随机数生成器，加入时间戳保证每次生成的随机数不一样
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	// 生成随机数
	return r.Int63n(end-start) + start
}

// 生成六位随机码
func GenerateRandomCode() string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	code := fmt.Sprintf("%06v", rnd.Int31n(1000000))
	return code
}

// 指定长度随机中文字符(包含复杂字符)
func GenFixedLengthChineseChars(length int) string {
	var buf bytes.Buffer
	for i := 0; i < length; i++ {
		buf.WriteRune(rune(RandInt(19968, 40869)))
	}
	return buf.String()
}

// 指定范围随机中文字符
func GenRandomLengthChineseChars(start, end int) string {
	length := RandInt(start, end)
	return GenFixedLengthChineseChars(length)
}

// 随机英文小写字母
func RandStr(len int) string {
	rand.Seed(time.Now().UnixNano())
	data := make([]byte, len)
	for i := 0; i < len; i++ {
		data[i] = byte(rand.Intn(26) + 97)
	}
	return string(data)
}

/**
 * 生成指定长度的随机字母和数字字符串，包括0-9、a-z、A-Z的所有字符。
 *
 * @param length
 * @return
 */
func RandString(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return *(*string)(unsafe.Pointer(&b))
}

var src = rand.NewSource(time.Now().UnixNano())

const letterBytes = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)
