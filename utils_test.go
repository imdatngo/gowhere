package gowhere

import (
	"reflect"
	"testing"
	"time"
)

func TestToString(t *testing.T) {
	testStr := "I'm String!"
	resultStr := "I'm String!"
	testNum := 123.456
	resultNum := "123.456"
	testTime := time.Now()
	resultTime := testTime.Format("2006-01-02 15:04:05")

	tests := []struct {
		name  string
		input interface{}
		want  string
	}{
		{
			name:  "string",
			input: testStr,
			want:  resultStr,
		},
		{
			name:  "number",
			input: testNum,
			want:  resultNum,
		},
		{
			name:  "string pointer",
			input: &testStr,
			want:  resultStr,
		},
		{
			name:  "number pointer",
			input: &testNum,
			want:  resultNum,
		},
		{
			name:  "byte",
			input: byte('a'),
			want:  "a",
		},
		{
			name:  "[]byte",
			input: []byte(testStr),
			want:  resultStr,
		},
		{
			name:  "rune",
			input: 'b',
			want:  "b",
		},
		{
			name:  "nil",
			input: nil,
			want:  "",
		},
		{
			name:  "time",
			input: testTime,
			want:  resultTime,
		},
		{
			name:  "time pointer",
			input: &testTime,
			want:  resultTime,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Utils.ToString(tt.input); got != tt.want {
				t.Errorf("Utils.ToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToSlice(t *testing.T) {
	tests := []struct {
		name  string
		input interface{}
		want  interface{}
	}{
		{
			name:  "a slice",
			input: []string{"a", "b", "c"},
			want:  []string{"a", "b", "c"},
		},
		{
			name:  "an array",
			input: [3]int{1, 2, 3},
			want:  [3]int{1, 2, 3},
		},
		{
			name:  "single value",
			input: "a",
			want:  []interface{}{"a"},
		},
		{
			name:  "nil",
			input: nil,
			want:  []interface{}{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Utils.ToSlice(tt.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Utils.ToSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToSQLVar(t *testing.T) {
	testTime := time.Now()
	resultTime := testTime.Format("2006-01-02 15:04:05")

	tests := []struct {
		name  string
		input interface{}
		want  interface{}
	}{
		{
			name:  "nil",
			input: nil,
			want:  "",
		},
		{
			name:  "time",
			input: testTime,
			want:  resultTime,
		},
		{
			name:  "*time",
			input: &testTime,
			want:  resultTime,
		},
		{
			name:  "string",
			input: "Hello World!",
			want:  "Hello World!",
		},
		{
			name:  "number",
			input: 123,
			want:  123,
		},
		{
			name:  "slice",
			input: []interface{}{"empty"},
			want:  []interface{}{"empty"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Utils.ToSQLVar(tt.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Utils.ToSQLVar() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToDateTime(t *testing.T) {
	testTime := time.Now()
	resultTime := testTime.Format("2006-01-02 15:04:05")
	testDate := time.Date(2019, time.April, 19, 0, 0, 0, 0, time.Local)
	resultDate := testDate.Format("2006-01-02")

	tests := []struct {
		name  string
		input interface{}
		want  interface{}
	}{
		{
			name:  "nil",
			input: nil,
			want:  "",
		},
		{
			name:  "time",
			input: testTime,
			want:  resultTime,
		},
		{
			name:  "date + pointer",
			input: &testDate,
			want:  resultDate,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Utils.ToDateTime(tt.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Utils.ToDateTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToDate(t *testing.T) {
	testTime := time.Now()
	resultTime := testTime.Format("2006-01-02")
	testDate := time.Date(2019, time.April, 19, 0, 0, 0, 0, time.Local)
	resultDate := testDate.Format("2006-01-02")

	tests := []struct {
		name  string
		input interface{}
		want  interface{}
	}{
		{
			name:  "nil",
			input: nil,
			want:  "",
		},
		{
			name:  "time",
			input: testTime,
			want:  resultTime,
		},
		{
			name:  "date + pointer",
			input: &testDate,
			want:  resultDate,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Utils.ToDate(tt.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Utils.ToDate() = %v, want %v", got, tt.want)
			}
		})
	}
}
