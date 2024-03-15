package slice_test

import (
	"fmt"
	"gitee.com/yanwenqing/backend-lib/slice"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

type MyInt int32

func TestSliceContains(t *testing.T) {
	tests := []struct {
		name   string
		arr    interface{}
		item   interface{}
		expect bool
	}{
		{
			name:   "int",
			arr:    []int{1, 2, 3},
			item:   1,
			expect: true,
		},
		{
			name:   "string",
			arr:    []string{"s1", "s2", "s3"},
			item:   "s1",
			expect: true,
		},
		{
			name:   "custom true",
			arr:    []MyInt{1, 2, 3},
			item:   MyInt(1),
			expect: true,
		},
		{
			name:   "custom false",
			arr:    []MyInt{1, 2, 3},
			item:   MyInt(4),
			expect: false,
		},
		{
			name:   "map true",
			arr:    map[string]int{"a": 1, "b": 2},
			item:   "a",
			expect: true,
		},
		{
			name:   "map false",
			arr:    map[string]int{"a": 1, "b": 2},
			item:   "d",
			expect: false,
		},
		{
			name:   "unsupported type",
			arr:    12,
			item:   "d",
			expect: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, slice.Contains(test.arr, test.item), test.expect)
		})
	}
}

func TestIntSliceToString(t *testing.T) {
	type args struct {
		arr []int64
		sep string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"normal", args{arr: []int64{1, 2, 3}, sep: ","}, "1,2,3"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := slice.IntSliceToString(tt.args.arr, tt.args.sep); got != tt.want {
				t.Errorf("IntSliceToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJoin(t *testing.T) {
	tests := []struct {
		name      string
		s1        []int64
		s2        []int64
		disExpect []int64
		norExpect []int64
	}{
		{
			name:      "join",
			s1:        []int64{1, 2, 3},
			s2:        []int64{2, 3, 4},
			disExpect: []int64{1, 2, 3, 4},
			norExpect: []int64{1, 2, 3, 2, 3, 4},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := slice.Concat(tt.s1, tt.s2, true); !reflect.DeepEqual(got, tt.disExpect) {
				t.Errorf("slice.Concat() = %v, want %v", got, tt.disExpect)
			}
			if got := slice.Concat(tt.s1, tt.s2, false); !reflect.DeepEqual(got, tt.norExpect) {
				t.Errorf("slice.Concat() = %v, want %v", got, tt.norExpect)
			}
		})
	}
}

func TestFilter(t *testing.T) {
	tests := []struct {
		name      string
		arr       []int64
		condition func(int64) bool
		expect    []int64
	}{
		{
			name: "normal1",
			arr:  []int64{1, 2, 3},
			condition: func(i int64) bool {
				return i != 2
			},
			expect: []int64{1, 3},
		},
		{
			name: "normal2",
			arr:  []int64{2},
			condition: func(i int64) bool {
				return i != 2
			},
			expect: []int64{},
		},
		{
			name: "not exist",
			arr:  []int64{1, 2, 3},
			condition: func(i int64) bool {
				return i == 4
			},
			expect: []int64{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := slice.Filter(tt.arr, tt.condition)
			if !reflect.DeepEqual(got, tt.expect) {
				t.Errorf("slice.Filter() = %v, want %v", got, tt.expect)
			}
		})
	}

	s1 := []int64{}
	s2 := []int64{1}
	fmt.Println(reflect.DeepEqual(s1, s2))
}

func TestSubtract(t *testing.T) {
	tests := []struct {
		name   string
		a      []int64
		b      []int64
		expect []int64
	}{
		{
			name:   "empty",
			a:      []int64{},
			b:      []int64{},
			expect: []int64{},
		},
		{
			name:   "a = nil",
			a:      nil,
			b:      []int64{},
			expect: []int64{},
		},
		{
			name:   "b = nil",
			a:      []int64{},
			b:      nil,
			expect: []int64{},
		},
		{
			name:   "both = nil",
			a:      nil,
			b:      nil,
			expect: []int64{},
		},
		{
			name:   "a > b",
			a:      []int64{1, 2, 3},
			b:      []int64{2},
			expect: []int64{1, 3},
		},
		{
			name:   "a < b",
			a:      []int64{1, 2},
			b:      []int64{1, 2, 3, 4},
			expect: []int64{},
		},
		{
			name:   "a = b",
			a:      []int64{1, 2},
			b:      []int64{1, 2},
			expect: []int64{},
		},
		{
			name:   "no intersect",
			a:      []int64{1, 2},
			b:      []int64{3, 4},
			expect: []int64{1, 2},
		},
		{
			name:   "cross",
			a:      []int64{1, 2, 3},
			b:      []int64{2, 3, 4},
			expect: []int64{1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := slice.Subtract(tt.a, tt.b)
			if !reflect.DeepEqual(got, tt.expect) {
				t.Errorf("slice.Subtract() = %v, want %v", got, tt.expect)
			}
		})
	}
}

func TestSameDistinctElements(t *testing.T) {
	tests := []struct {
		name   string
		a      []interface{}
		b      []interface{}
		expect bool
	}{
		{
			name:   "empty",
			a:      []interface{}{},
			b:      []interface{}{},
			expect: true,
		},
		{
			name:   "a = nil",
			a:      nil,
			b:      []interface{}{},
			expect: true,
		},
		{
			name:   "b = nil",
			a:      []interface{}{},
			b:      nil,
			expect: true,
		},
		{
			name:   "both = nil",
			a:      nil,
			b:      nil,
			expect: true,
		},
		{
			name:   "a > b",
			a:      []interface{}{1, 2, 3},
			b:      []interface{}{2},
			expect: false,
		},
		{
			name:   "a < b",
			a:      []interface{}{1, 2},
			b:      []interface{}{1, 2, 3, 4},
			expect: false,
		},
		{
			name:   "no intersect",
			a:      []interface{}{1, 2, 3},
			b:      []interface{}{4, 5, 6},
			expect: false,
		},
		{
			name:   "same element, different order",
			a:      []interface{}{1, 2, 3},
			b:      []interface{}{1, 3, 2},
			expect: true,
		},
		{
			name:   "same element, same order",
			a:      []interface{}{1, 2, 3},
			b:      []interface{}{1, 2, 3},
			expect: true,
		},
		{
			name:   "same element with duplicate",
			a:      []interface{}{1, 2, 3, 2, 2},
			b:      []interface{}{1, 2, 3},
			expect: true,
		},
		{
			name:   "same element with duplicate, reverse",
			a:      []interface{}{1, 2, 3, 2, 2},
			b:      []interface{}{1, 2, 3},
			expect: true,
		},
		{
			name:   "different element with duplicate",
			a:      []interface{}{1, 2, 3, 2, 2, 4},
			b:      []interface{}{1, 2, 3},
			expect: false,
		},
		{
			name:   "different element with duplicate, reverse",
			a:      []interface{}{1, 2, 3},
			b:      []interface{}{1, 2, 3, 2, 2, 4},
			expect: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := slice.SameDistinctElements(tt.a, tt.b)
			if tt.expect != got {
				t.Errorf("slice.SameElements() = %v, want %v", got, tt.expect)
			}
		})
	}
}
