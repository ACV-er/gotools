package lib

import (
	"reflect"
	"testing"
)

type test_stringer struct {
	s string
}

func (t test_stringer) String() string {
	return t.s
}

type test_not_stringer struct {
	S string `json:"s"`
}

func TestAnyToString(t *testing.T) {
	type args struct {
		src interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantRet string
	}{
		{
			name: "string",
			args: args{
				src: "aadd",
			},
			wantRet: "aadd",
		},
		{
			name: "int",
			args: args{
				src: 16,
			},
			wantRet: "16",
		},
		{
			name: "int64",
			args: args{
				src: int64(16),
			},
			wantRet: "16",
		},
		{
			name: "float",
			args: args{
				src: 16.0,
			},
			wantRet: "16",
		},
		{
			name: "bool",
			args: args{
				src: true,
			},
			wantRet: "true",
		},
		{
			name: "nil",
			args: args{
				src: nil,
			},
			wantRet: "",
		},
		{
			name: "map",
			args: args{
				src: map[string]string{
					"a": "a",
					"b": "b",
				},
			},
			wantRet: `{"a":"a","b":"b"}`,
		},
		{
			name: "have string",
			args: args{
				src: test_stringer{
					s: "aadd",
				},
			},
			wantRet: "aadd",
		},
		{
			name: "struct指针 嵌套",
			args: args{
				src: &struct {
					Query string `json:"q"`
					Data  *struct {
						A string `json:"a"`
					} `json:"data"`
				}{
					Query: "aadd",
					Data: &struct {
						A string `json:"a"`
					}{
						A: "a",
					},
				},
			},
			wantRet: `{"q":"aadd","data":{"a":"a"}}`,
		},
		{
			name: "not have string",
			args: args{
				src: test_not_stringer{
					S: "aadd",
				},
			},
			wantRet: `{"s":"aadd"}`,
		},
		{
			name: "have string ptr",
			args: args{
				src: &test_stringer{
					s: "aadd",
				},
			},
			wantRet: "aadd",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRet := AnyToString(tt.args.src); gotRet != tt.wantRet {
				t.Errorf("AnyToString() = %v, want %v", gotRet, tt.wantRet)
			}
		})
	}
}

func TestGenMapFromObject(t *testing.T) {
	type args struct {
		src interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantRet map[string]interface{}
		wantErr bool
	}{
		{
			name: "结构体转换成map",
			args: args{
				src: &struct {
					A string `json:"a"`
					B int    `json:"b"`
				}{
					A: "a",
					B: 5,
				},
			},
			wantRet: map[string]interface{}{
				"a": "a",
				"b": 5,
			},
			wantErr: false,
		},
		{
			name: "map转换成map",
			args: args{
				src: map[string]string{
					"a": "a",
					"b": "b",
				},
			},
			wantRet: map[string]interface{}{
				"a": "a",
				"b": "b",
			},
			wantErr: false,
		},
		{
			name: "struct指针转换成map",
			args: args{
				src: &struct {
					A string `json:"a"`
					B string `json:"b"`
					C float64
				}{
					A: "a",
					B: "b",
					C: 3.14,
				},
			},
			wantRet: map[string]interface{}{
				"a": "a",
				"b": "b",
				"C": 3.14,
			},
			wantErr: false,
		},
		{
			name: "struct指针嵌套",
			args: args{
				src: &struct {
					A string `json:"a"`
					B *struct {
						C string `json:"c"`
					} `json:"b"`
				}{
					A: "a",
					B: &struct {
						C string `json:"c"`
					}{
						C: "c",
					},
				},
			},
			wantRet: map[string]interface{}{
				"a": "a",
				"b": struct {
					C string `json:"c"`
				}{
					C: "c",
				},
			},
			wantErr: false,
		},
		{
			name: "can not convert",
			args: args{
				src: 123,
			},
			wantRet: map[string]interface{}{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRet, err := GenMapFromObject(tt.args.src)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenMapFromObject() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRet, tt.wantRet) {
				t.Errorf("GenMapFromObject() = %v, want %v", gotRet, tt.wantRet)
			}
		})
	}
}

func TestGenArrFromObjectOrderByFieldNameAsc(t *testing.T) {
	type args struct {
		src interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantRet SortMap
		wantErr bool
	}{
		{
			name: "普通结构体",
			args: args{
				src: &struct {
					B int    `json:"b"`
					X string `json:"x"`
					Z int    `json:"z"`
					H string `json:"h"`
					A string `json:"a"`
				}{
					B: 5,
					X: "x",
					Z: 6,
					H: "h",
					A: "a",
				},
			},
			wantRet: SortMap{
				{"a", "a"},
				{"b", 5},
				{"h", "h"},
				{"x", "x"},
				{"z", 6},
			},
			wantErr: false,
		},
		{
			name: "报错",
			args: args{
				src: 123,
			},
			wantRet: SortMap{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRet, err := GenArrFromObjectOrderByFieldNameAsc(tt.args.src)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenArrFromObjectOrderByFieldNameAsc() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRet, tt.wantRet) {
				t.Errorf("GenArrFromObjectOrderByFieldNameAsc() = %v, want %v", gotRet, tt.wantRet)
			}
		})
	}
}
