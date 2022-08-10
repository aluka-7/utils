package utils

import "reflect"

// Contains返回数组中val的索引位置,当为-1时代表不包含
func Contains(array interface{}, val interface{}) (index int) {
	index = -1
	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		{
			s := reflect.ValueOf(array)
			for i := 0; i < s.Len(); i++ {
				if reflect.DeepEqual(val, s.Index(i).Interface()) {
					index = i
					return
				}
			}
		}
	}
	return
}

// ContainsString返回string类型'val'在数组中的索引位置
func ContainsString(array []string, val string) (index int) {
	index = -1
	for i := 0; i < len(array); i++ {
		if array[i] == val {
			index = i
			return
		}
	}
	return
}

// ContainsInt返回int64类型'val'在数组中的索引位置
func ContainsInt(array []int64, val int64) (index int) {
	index = -1
	for i := 0; i < len(array); i++ {
		if array[i] == val {
			index = i
			return
		}
	}
	return
}

// ContainsUint返回uint64类型'val'在数组中的索引位置
func ContainsUint(array []uint64, val uint64) (index int) {
	index = -1
	for i := 0; i < len(array); i++ {
		if array[i] == val {
			index = i
			return
		}
	}
	return
}

// ContainsBool返回bool类型'val'在数组中的索引位置
func ContainsBool(array []bool, val bool) (index int) {
	index = -1
	for i := 0; i < len(array); i++ {
		if array[i] == val {
			index = i
			return
		}
	}
	return
}

// ContainsFloat返回float64类型'val'在数组中的索引位置
func ContainsFloat(array []float64, val float64) (index int) {
	index = -1
	for i := 0; i < len(array); i++ {
		if array[i] == val {
			index = i
			return
		}
	}
	return
}

// ContainsComplex返回complex128类型'val'在数组中的索引位置
func ContainsComplex(array []complex128, val complex128) (index int) {
	index = -1
	for i := 0; i < len(array); i++ {
		if array[i] == val {
			index = i
			return
		}
	}
	return
}
