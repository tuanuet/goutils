package utils

import (
	"fmt"

	jsoniter "github.com/json-iterator/go"
	"github.com/tidwall/gjson"
)

func ArrToMapIdentify[T any, K comparable](arrInput []T, funcGetKey func(T) K) map[K]T {
	return ArrToMap(arrInput, funcGetKey, func(t T) T { return t })
}

func ArrToMap[T any, M any, K comparable](arrInput []T, funcGetKey func(T) K, funcGetValue func(T) M) map[K]M {
	mapRes := make(map[K]M, len(arrInput))

	for _, e := range arrInput {
		mapRes[funcGetKey(e)] = funcGetValue(e)
	}

	return mapRes
}

func InSlice[T comparable](element T, slice []T) bool {
	for _, v := range slice {
		if v == element {
			return true
		}
	}
	return false
}

func Unique[T comparable](s []T) []T {
	inResult := make(map[T]bool)
	var result []T
	for _, str := range s {
		if _, ok := inResult[str]; !ok {
			inResult[str] = true
			result = append(result, str)
		}
	}
	return result
}

func UniqueElement[T any, K comparable](input []T, f func(t T) K) []T {
	inResult := make(map[K]bool)
	var result []T
	for _, t := range input {
		key := f(t)
		if _, ok := inResult[key]; !ok {
			inResult[key] = true
			result = append(result, t)
		}
	}
	return result
}

func GroupArray[T any, K comparable](items []T, getKey func(t T) K) map[K][]T {
	result := make(map[K][]T)
	for _, item := range items {
		key := getKey(item)
		result[key] = append(result[key], item)
	}
	return result
}

// RemoveIndex ...
func RemoveIndex[T any](slice []T, pos int) []T {
	newArr := make([]T, len(slice)-1)
	k := 0
	for i := 0; i < (len(slice) - 1); {
		if i != pos {
			newArr[i] = slice[k]
			k++
		} else {
			k++
		}
		i++
	}

	return newArr
}

// GetArrayFromJsonString ... return array by key, and array must not be empty
func GetArrayFromJsonString[T any](jsonData, key string) ([]T, error) {
	// 1. Get key value
	cfg := gjson.Get(jsonData, key)
	if !cfg.Exists() {
		return nil, fmt.Errorf("missing variable %s, cfgVariables = %s\n", key, jsonData)
	}

	// 2. Is Array?
	if !cfg.IsArray() {
		return nil, fmt.Errorf("%s must be array", key)
	}

	// 3. Convert Array
	resArr := cfg.Array()
	var res []T
	for _, item := range resArr {
		var tmp T
		if err := jsoniter.Unmarshal([]byte(item.String()), &tmp); err == nil {
			res = append(res, tmp)
		}
	}

	// 4. Is empty?
	if len(res) == 0 {
		return nil, fmt.Errorf("%s must not be empty", key)
	}

	// 5. Return
	return res, nil
}

// Map ... A[] -> B[] by function
func Map[T, U any](ts []T, f func(T) U) []U {
	us := make([]U, len(ts))
	for i := range ts {
		us[i] = f(ts[i])
	}
	return us
}

// Filter ...
func Filter[T any](ts []T, f func(T) bool) []T {
	var res []T
	for i := range ts {
		if ok := f(ts[i]); ok {
			res = append(res, ts[i])
		}
	}
	return res
}

// FindOne ... only return one element, if you want more, please use Filter
func FindOne[T any](ts []T, f func(T) bool) *T {
	for i := range ts {
		if ok := f(ts[i]); ok {
			return &ts[i]
		}
	}
	return nil
}

// MapFilter ... A[] -> B[] by function
func MapFilter[T, U any](ts []T, f func(T) (U, bool)) []U {
	var us []U
	for i := range ts {
		if newData, ok := f(ts[i]); ok {
			us = append(us, newData)
		}
	}
	return us
}

// DifferentArrays find elements only exist in arr1 (comparable by common field in each arr)
func DifferentArrays[T, M any, K comparable](arr1 []T, arr2 []M, funcGetCompareKey1 func(T) K, funcGetCompareKey2 func(M) K) []T {
	appearedMap := make(map[K]bool)

	onlyInArr1 := make([]T, 0)

	for _, item := range arr2 {
		appearedMap[funcGetCompareKey2(item)] = true
	}

	for _, item := range arr1 {
		_, f := appearedMap[funcGetCompareKey1(item)]
		if !f {
			onlyInArr1 = append(onlyInArr1, item)
		}
	}

	return onlyInArr1
}

func Sum[T any, F int | int32 | int64](array []T, getValue func(T) F) F {
	var sum F = 0
	for _, element := range array {
		sum += getValue(element)
	}
	return sum
}
