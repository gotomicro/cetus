package kc

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"reflect"
)

// Union 求并集
func Union(slice1, slice2 []string) []string {
	m := make(map[string]int)
	for _, v := range slice1 {
		m[v]++
	}
	for _, v := range slice2 {
		times := m[v]
		if times == 0 {
			slice1 = append(slice1, v)
		}
	}
	return slice1
}

// Intersect 求交集
func Intersect(slice1, slice2 []string) []string {
	m := make(map[string]int)
	nn := make([]string, 0)
	for _, v := range slice1 {
		m[v]++
	}
	for _, v := range slice2 {
		times := m[v]
		if times == 1 {
			nn = append(nn, v)
		}
	}
	return nn
}

// Difference 求差集 slice1-并集
func Difference(slice1, slice2 []string) []string {
	m := make(map[string]int)
	nn := make([]string, 0)
	inter := Intersect(slice1, slice2)
	for _, v := range inter {
		m[v]++
	}
	for _, value := range slice1 {
		times := m[value]
		if times == 0 {
			nn = append(nn, value)
		}
	}
	return nn
}

// DeepCopy 深度拷贝
func DeepCopy(dst, src interface{}) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(src); err != nil {
		return err
	}
	return gob.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(dst)
}

// IsSliceEqual ..
func IsSliceEqual(a, b interface{}) bool {
	if (a == nil) && (b == nil) {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	if reflect.TypeOf(a).Kind() != reflect.Array && reflect.TypeOf(a).Kind() != reflect.Slice {
		return false
	}

	lenA := reflect.ValueOf(a).Len()
	lenB := reflect.ValueOf(b).Len()

	if lenA != lenB {
		return false
	}
	valA := reflect.ValueOf(a)
	valB := reflect.ValueOf(b)

	mapA := make(map[interface{}]int)
	for i := 0; i < lenA; i++ {
		av := valA.Index(i).Interface()
		if _, exist := mapA[av]; !exist {
			mapA[av] = 0
		} else {
			mapA[av]++
		}
	}

	for i := 0; i < lenB; i++ {
		bv := valB.Index(i).Interface()
		if _, exist := mapA[bv]; !exist {
			return false
		}
		mapA[bv]--
		if mapA[bv] < 0 {
			delete(mapA, bv)
		}
	}
	return len(mapA) == 0

}

func Ints2String(ints []int) string {
	var res = ""
	for k, v := range ints {
		if k == 0 {
			res = fmt.Sprintf("%d", v)
		} else {
			res = fmt.Sprintf("%s,%d", res, v)
		}
	}
	return res
}
