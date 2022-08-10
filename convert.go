package utils

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"strconv"
)

// 转换字符串以指定类型。
type StrTo string

func (f StrTo) Exist() bool {
	return string(f) != string(0x1E)
}

func (f StrTo) Uint8() (uint8, error) {
	v, err := strconv.ParseUint(f.String(), 10, 8)
	return uint8(v), err
}

func (f StrTo) Int() (int, error) {
	v, err := strconv.ParseInt(f.String(), 10, 0)
	return int(v), err
}

func (f StrTo) Int64() (int64, error) {
	v, err := strconv.ParseInt(f.String(), 10, 64)
	return int64(v), err
}

func (f StrTo) MustUint8() uint8 {
	v, _ := f.Uint8()
	return v
}

func (f StrTo) MustInt() int {
	v, _ := f.Int()
	return v
}

func (f StrTo) MustInt64() int64 {
	v, _ := f.Int64()
	return v
}
func (f StrTo) Float64() float64 {
	v, _ := strconv.ParseFloat(f.String(), 10)
	return v
}
func (f StrTo) String() string {
	if f.Exist() {
		return string(f)
	}
	return ""
}
func (f StrTo) Hash() string {
	if len(f) > 0 {
		h := fnv.New32a()
		h.Write([]byte(f))
		return strconv.FormatUint(uint64(h.Sum32()), 10)
	} else {
		return ""
	}
}

// 将任何类型转换为字符串。
func ToStr(value interface{}, args ...int) (s string) {
	switch v := value.(type) {
	case bool:
		s = strconv.FormatBool(v)
	case float32:
		s = strconv.FormatFloat(float64(v), 'f', ArgInt(args).Get(0, -1), ArgInt(args).Get(1, 32))
	case float64:
		s = strconv.FormatFloat(v, 'f', ArgInt(args).Get(0, -1), ArgInt(args).Get(1, 64))
	case int:
		s = strconv.FormatInt(int64(v), ArgInt(args).Get(0, 10))
	case int8:
		s = strconv.FormatInt(int64(v), ArgInt(args).Get(0, 10))
	case int16:
		s = strconv.FormatInt(int64(v), ArgInt(args).Get(0, 10))
	case int32:
		s = strconv.FormatInt(int64(v), ArgInt(args).Get(0, 10))
	case int64:
		s = strconv.FormatInt(v, ArgInt(args).Get(0, 10))
	case uint:
		s = strconv.FormatUint(uint64(v), ArgInt(args).Get(0, 10))
	case uint8:
		s = strconv.FormatUint(uint64(v), ArgInt(args).Get(0, 10))
	case uint16:
		s = strconv.FormatUint(uint64(v), ArgInt(args).Get(0, 10))
	case uint32:
		s = strconv.FormatUint(uint64(v), ArgInt(args).Get(0, 10))
	case uint64:
		s = strconv.FormatUint(v, ArgInt(args).Get(0, 10))
	case string:
		s = v
	case []byte:
		s = string(v)
	default:
		s = fmt.Sprintf("%v", v)
	}
	return s
}

type ArgInt []int

func (a ArgInt) Get(i int, args ...int) (r int) {
	if i >= 0 && i < len(a) {
		r = a[i]
	} else if len(args) > 0 {
		r = args[0]
	}
	return
}

type ArgString []string

func (a ArgString) Get(i int, args ...string) (r string) {
	if i >= 0 && i < len(a) {
		r = a[i]
	} else if len(args) > 0 {
		r = args[0]
	}
	return
}

type ArgAny []interface{}

func (a ArgAny) Get(i int, args ...interface{}) (r interface{}) {
	if i >= 0 && i < len(a) {
		r = a[i]
	}
	if len(args) > 0 {
		r = args[0]
	}
	return
}

// 将十六进制格式字符串转换为十进制数.
func HexStr2int(hexStr string) (int, error) {
	num := 0
	length := len(hexStr)
	for i := 0; i < length; i++ {
		char := hexStr[length-i-1]
		factor := -1

		switch {
		case char >= '0' && char <= '9':
			factor = int(char) - '0'
		case char >= 'a' && char <= 'f':
			factor = int(char) - 'a' + 10
		default:
			return -1, fmt.Errorf("invalid hex: %s", string(char))
		}

		num += factor * PowInt(16, i)
	}
	return num, nil
}

// 将十进制数格式转换为十六进制字符串.
func Int2HexStr(num int) (hex string) {
	if num == 0 {
		return "0"
	}

	for num > 0 {
		r := num % 16

		c := "?"
		if r >= 0 && r <= 9 {
			c = string(r + '0')
		} else {
			c = string(r + 'a' - 10)
		}
		hex = c + hex
		num = num / 16
	}
	return hex
}

// 整形数组转换成字符串
func ArrayToString(A []int, denim string) string {
	var buffer bytes.Buffer
	for i := 0; i < len(A); i++ {
		buffer.WriteString(strconv.Itoa(A[i]))
		if i != len(A)-1 {
			buffer.WriteString(denim)
		}
	}

	return buffer.String()
}
