package http_tool

import (
	"net/url"
	"strings"

	"github.com/ACV-er/gotools/lib"
)

// 将类、map生成&连接的http请求参数
func GenFromParam(src_param interface{}) (param string, err error) {
	param_map, err := lib.GenMapFromObject(src_param)
	if err != nil {
		return
	}

	paramArr := make([]string, 0, len(param_map))
	for k, v := range param_map {
		paramArr = append(paramArr, url.QueryEscape(k)+"="+url.QueryEscape(lib.AnyToString(v)))
	}
	param = strings.Join(paramArr, "&")
	return
}

func GenGetUrl(baseUrl string, src_param interface{}) (retUrl string, err error) {
	retUrl = baseUrl
	param, err := GenFromParam(src_param)
	if err != nil {
		return
	}

	retUrl = baseUrl + "?" + param
	return
}
