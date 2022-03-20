# 工具包

> 放置一些不好归类的工具

* http生成按key排序的查询串示例

```go
package main

import (
	"bytes"
	"fmt"

	"github.com/ACV-er/gotools/lib"
)

// 文件为测试用

type GetRoleRequest struct {
	RoleId int    `json:"role_id"`
	Time   int64  `json:"time"`
	Onece  string `json:"onece"`
	Data   struct {
		Name  string `json:"name"`
		Age   int    `json:"age"`
		Sorec struct {
			Name string `json:"name"`
			Age  int    `json:"age"`
		} `json:"sorec"`
	} `json:"data"`
	Sgin string `json:"sgin"`
}

func main() {
	req := &GetRoleRequest{
		RoleId: 1,
		Time:   245,
		Onece:  "dwa4dsf12",
		Data: struct {
			Name  string `json:"name"`
			Age   int    `json:"age"`
			Sorec struct {
				Name string `json:"name"`
				Age  int    `json:"age"`
			} `json:"sorec"`
		}{
			Name: "dwa4dsf12",
			Age:  12,
			Sorec: struct {
				Name string `json:"name"`
				Age  int    `json:"age"`
			}{
				Name: "dwa4dsf12",
				Age:  12,
			},
		},
		Sgin: "sgin",
	}

	str_buf := bytes.NewBufferString("")

	ret, _ := lib.GenArrFromObjectOrderByFieldNameAsc(req)

	for _, v := range ret {
		str_buf.WriteString(v.Key + "=" + lib.AnyToString(v.Value) + "&")
	}

	// 去掉末尾的&
	if str_buf.Len() > 0 {
		str_buf.Truncate(str_buf.Len() - 1)
	}

	str := str_buf.String()
	fmt.Println(str)
}

```