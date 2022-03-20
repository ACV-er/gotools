package http_tool

import (
	"strings"
	"testing"
)

type test_stringer struct {
	S string
}

func (t test_stringer) String() string {
	return t.S
}

type test_not_stringer struct {
	S string `json:"s"`
}

// 判断数组内是否存在某个值,辅助函数
func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// 判断两个数组是否有相同元素，不管顺序，辅助函数
func equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for _, v := range a {
		if !contains(b, v) {
			return false
		}
	}
	return true
}

func TestGenFromParam(t *testing.T) {
	type args struct {
		src_param interface{}
	}
	tests := []struct {
		name      string
		args      args
		wantParam []string
		wantErr   bool
	}{
		{
			name: "结构体转换成url参数-单参数",
			args: args{
				src_param: &struct {
					Query string `json:"q"`
				}{
					Query: "aadd",
				},
			},
			wantParam: []string{"q=aadd"},
			wantErr:   false,
		},
		{
			name: "结构体转换成url参数-带数字",
			args: args{
				src_param: &struct {
					Query int `json:"q"`
				}{
					Query: 16,
				},
			},
			wantParam: []string{"q=16"},
			wantErr:   false,
		},
		{
			name: "结构体转换成url参数-多参数",
			args: args{
				src_param: &struct {
					Query string `json:"q"`
					Limit string `json:"limit"`
				}{
					Query: "aadd",
					Limit: "20",
				},
			},
			wantParam: []string{"q=aadd", "limit=20"},
			wantErr:   false,
		},
		{
			name: "结构体转换成url参数-多参数-url编码",
			args: args{
				src_param: &struct {
					Query string `json:"q"`
					Limit string `json:"limit"`
				}{
					Query: "+=&/测试",
					Limit: "20",
				},
			},
			wantParam: []string{"q=%2B%3D%26%2F%E6%B5%8B%E8%AF%95", "limit=20"},
			wantErr:   false,
		},
		{
			name: "结构体转换成url参数-多参数-各种类型-url编码",
			args: args{
				src_param: &struct {
					Query    string             `json:"q"`
					Limit    string             `json:"limit"`
					Page     int                `json:"page"`
					Score    float64            `json:"score"`
					Time     int64              `json:"time"`
					Is       bool               `json:"is"`
					SObject  *test_stringer     `json:"sobject"`
					SObject2 *test_not_stringer `json:"sobject2"`
				}{
					Query: "+=&/测试",
					Limit: "20",
					Page:  1,
					Score: 0.5,
					Time:  123456789,
					Is:    true,
					SObject: &test_stringer{
						S: "测试",
					},
					SObject2: &test_not_stringer{
						S: "测试",
					},
				},
			},
			wantParam: []string{"q=%2B%3D%26%2F%E6%B5%8B%E8%AF%95", "limit=20", "page=1", "score=0.5", "time=123456789", "is=true", "sobject=%E6%B5%8B%E8%AF%95", "sobject2=%7B%22s%22%3A%22%E6%B5%8B%E8%AF%95%22%7D"},
			wantErr:   false,
		},
		{
			name: "map转换成url参数-单参数",
			args: args{
				src_param: map[string]string{
					"q": "aadd",
				},
			},
			wantParam: []string{"q=aadd"},
			wantErr:   false,
		},
		{
			name: "map转换成url参数-多参数",
			args: args{
				src_param: map[string]string{
					"q":     "aadd",
					"limit": "20",
				},
			},
			wantParam: []string{"q=aadd", "limit=20"},
			wantErr:   false,
		},
		{
			name: "map转换成url参数-多参数-url编码",
			args: args{
				src_param: map[string]string{
					"q":     "+=&/测试",
					"limit": "20",
				},
			},
			wantParam: []string{"q=%2B%3D%26%2F%E6%B5%8B%E8%AF%95", "limit=20"},
			wantErr:   false,
		},
		{
			name: "验证无法生成的参数报错",
			args: args{
				src_param: 18,
			},
			wantParam: []string{""},
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotParam, err := GenFromParam(tt.args.src_param)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenFromParam() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !equal(tt.wantParam, strings.Split(gotParam, "&")) {
				t.Errorf("GenFromParam() = %v, want in %v", gotParam, tt.wantParam)
			}
		})
	}
}

func TestGenGetUrl(t *testing.T) {
	type args struct {
		baseUrl   string
		src_param interface{}
	}
	tests := []struct {
		name       string
		args       args
		wantRetUrl []string
		wantErr    bool
	}{
		{
			name: "正常流程",
			args: args{
				baseUrl: "http://www.baidu.com",
				src_param: map[string]string{
					"q":     "+=&/测试",
					"limit": "20",
				},
			},
			wantRetUrl: []string{"http://www.baidu.com?q=%2B%3D%26%2F%E6%B5%8B%E8%AF%95&limit=20", "http://www.baidu.com?limit=20&q=%2B%3D%26%2F%E6%B5%8B%E8%AF%95"},
			wantErr:    false,
		},
		{
			name: "验证无法生成的参数报错",
			args: args{
				baseUrl:   "http://baidu.com",
				src_param: 18,
			},
			wantRetUrl: []string{"http://baidu.com"},
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRetUrl, err := GenGetUrl(tt.args.baseUrl, tt.args.src_param)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenGetUrl() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !contains(tt.wantRetUrl, gotRetUrl) {
				t.Errorf("GenGetUrl() = %v, want in %v", gotRetUrl, tt.wantRetUrl)
			}
		})
	}
}
