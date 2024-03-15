package slice

import (
	"fmt"
	"reflect"
	"strings"
)

// Contains for now, only item of same type is supported.
func Contains(arr interface{}, item interface{}) bool {
	arrValue := reflect.ValueOf(arr)
	switch reflect.TypeOf(arr).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < arrValue.Len(); i++ {
			if arrValue.Index(i).Interface() == item {
				return true
			}
		}
	case reflect.Map:
		if arrValue.MapIndex(reflect.ValueOf(item)).IsValid() {
			return true
		}
	}
	return false
}

func IntSliceToString(arr []int64, sep string) string {
	return strings.Replace(strings.Trim(fmt.Sprint(arr), "[]"), " ", sep, -1)
}

func Concat(a, b []int64, distinct bool) []int64 {
	if !distinct {
		return append(a, b...)
	}
	m := make(map[int64]bool)
	for _, v := range a {
		m[v] = true
	}
	for _, v := range b {
		if _, exist := m[v]; !exist {
			a = append(a, v)
		}
	}
	return a
}

func Filter(arr []int64, meet func(int64) bool) []int64 {
	res := []int64{}
	for j := 0; j < len(arr); j++ {
		if meet(arr[j]) {
			res = append(res, arr[j])
		}
	}
	return res
}

func Subtract(a, b []int64) []int64 {
	bm := make(map[int64]bool, len(b))
	for _, v := range b {
		bm[v] = true
	}
	res := make([]int64, 0)
	for _, v := range a {
		if _, ok := bm[v]; ok {
			continue
		}
		res = append(res, v)
	}
	return res
}

func SameDistinctElements(a, b []interface{}) bool {
	m := make(map[interface{}]int)
	for _, v := range a {
		if _, ok := m[v]; !ok {
			m[v] = 0
		}
	}

	for _, v := range b {
		if _, ok := m[v]; !ok {
			return false
		}
		delete(m, v)
	}

	return len(m) == 0
}
