package utils

import (
	"reflect"
	"testing"
)

type testStruct struct {
	MinCoverage float32 `json:"min_coverage"`
	IndexValue  float32 `json:"index_value"`
}

func TestUnique(t *testing.T) {
	type testCase[T any] struct {
		name string
		args []T
		want []T
	}

	tests := []testCase[int32]{
		{name: "When slice number", args: []int32{1, 2, 4, 2}, want: []int32{1, 2, 4}},
		//{name: "When slice string", args: []string{"1", "2", "4", "2"}, want: []string{"1", "2", "4"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Unique(tt.args)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getArrayValueFromCfgVariables() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_GetArrayFromJsonString(t *testing.T) {
	type args struct {
		cfgVariables string
	}
	type testCase[T any] struct {
		name    string
		args    args
		want    []T
		wantErr bool
	}
	tests := []testCase[testStruct]{
		{"test coverageLadder case empty -> return error", args{""}, nil, true},
		{"test coverageLadder case json wrong format -> return error", args{`{abc: ede}`}, nil, true},
		{"test coverageLadder case not contains key -> return error", args{`{"cov":[]}`}, nil, true},
		{"test coverageLadder case it is not array -> return error", args{`{"coverage_ladders":123}`}, nil, true},
		{"test coverageLadder case empty array -> return error", args{`{"coverage_ladders":[]}`}, nil, true},

		{"test coverageLadder case happy -> return correct coverageLadders", args{coverageLaddersConfigMock},
			[]testStruct{{20.0, 0.1}, {70.0, 0.9}, {80.0, 0.95}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetArrayFromJsonString[testStruct](tt.args.cfgVariables, "coverage_ladders")
			if (err != nil) != tt.wantErr {
				t.Errorf("getArrayValueFromCfgVariables() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getArrayValueFromCfgVariables() got = %v, want %v", got, tt.want)
			}
		})
	}
}

const coverageLaddersConfigMock = `
		{
		  "coverage_ladders": [
			{"min_coverage": 20, "index_value": 0.1},
			{"min_coverage": 70, "index_value": 0.9},
			{"min_coverage": 80, "index_value": 0.95}
		  ]
		}
	`

func TestDifferentArrays(t *testing.T) {
	type MyStruct1 struct {
		ID   int
		Name string
	}

	type MyStruct2 struct {
		Age  int
		Name string
	}

	funcGetCompareKey1 := func(item MyStruct1) string {
		return item.Name
	}

	funcGetCompareKey2 := func(item MyStruct2) string {
		return item.Name
	}

	type input struct {
		arr1 []MyStruct1
		arr2 []MyStruct2
	}

	type expect struct {
		expectedOnlyInArr1 []MyStruct1
	}

	type testCase struct {
		name   string
		input  input
		expect expect
	}
	tests := []testCase{
		{
			name: "happy case",
			input: input{
				arr1: []MyStruct1{
					{ID: 1, Name: "A"},
					{ID: 2, Name: "B"},
					{ID: 3, Name: "C"},
				},
				arr2: []MyStruct2{
					{Age: 1, Name: "B"},
					{Age: 2, Name: "C"},
					{Age: 3, Name: "D"},
				},
			},

			expect: expect{
				expectedOnlyInArr1: []MyStruct1{
					{ID: 1, Name: "A"},
				},
			},
		},
		{
			name: "when empty",
			input: input{
				arr1: []MyStruct1{},
				arr2: []MyStruct2{
					{Age: 1, Name: "B"},
				},
			},

			expect: expect{
				expectedOnlyInArr1: []MyStruct1{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			onlyInArr1 := DifferentArrays(tt.input.arr1, tt.input.arr2, funcGetCompareKey1, funcGetCompareKey2)

			if !reflect.DeepEqual(onlyInArr1, tt.expect.expectedOnlyInArr1) {
				t.Errorf("onlyInArr1 not equal to expectedOnlyInArr1: %v", onlyInArr1)
			}
		})
	}
}

func TestSum(t *testing.T) {
	type MyStruct struct {
		f1 int
		f2 int32
	}

	type args[T any, F interface{ int | int32 | int64 }] struct {
		array    []T
		getValue func(T) F
	}
	type testCase[T any, F interface{ int | int32 | int64 }] struct {
		name string
		args args[T, F]
		want F
	}
	funcGetF1 := func(s MyStruct) int { return s.f1 }
	tests := []testCase[MyStruct, int]{
		{"test case empty", args[MyStruct, int]{[]MyStruct{}, funcGetF1}, 0},
		{"test case 1 element", args[MyStruct, int]{[]MyStruct{{f1: 1, f2: 2}}, funcGetF1}, 1},
		{"test case n element", args[MyStruct, int]{[]MyStruct{{f1: 1}, {f1: 2}, {f1: 3}, {f1: 4}}, funcGetF1}, 10},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Sum(tt.args.array, tt.args.getValue); got != tt.want {
				t.Errorf("Sum() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_GroupArray(t *testing.T) {
	type Item struct {
		Key  string
		Data int
	}

	type args[T any, K comparable] struct {
		items  []T
		getKey func(t T) K
	}
	type testCase[T any, K comparable] struct {
		name string
		args args[T, K]
		want map[K][]T
	}
	getKey := func(i Item) string { return i.Key }
	tests := []testCase[Item, string]{
		{"test case nil", args[Item, string]{}, map[string][]Item{}},
		{"test case empty", args[Item, string]{[]Item{}, getKey}, map[string][]Item{}},
		{"test case 1 item", args[Item, string]{[]Item{{"key1", 1}}, getKey}, map[string][]Item{"key1": {{"key1", 1}}}},
		{"test case 2 item, 2 group", args[Item, string]{[]Item{{"key1", 1}, {"key2", 2}}, getKey}, map[string][]Item{"key1": {{"key1", 1}}, "key2": {{"key2", 2}}}},
		{"test case 2 item, 1 group", args[Item, string]{[]Item{{"key1", 1}, {"key1", 2}}, getKey}, map[string][]Item{"key1": {{"key1", 1}, {"key1", 2}}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GroupArray(tt.args.items, tt.args.getKey); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GroupArray() = %v, want %v", got, tt.want)
			}
		})
	}
}
