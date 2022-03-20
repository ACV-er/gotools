# 压缩数据

* 一共两个方法，压缩和解压缩

* 示例

```
package main

import (
	"fmt"

	"github.com/ACV-er/gotools/compress"
)

// 重复原字符串
func repeatStr(str string, times int) string {
	var ret string
	for i := 0; i < times; i++ {
		ret += str
	}
	return ret
}

// 文件为测试用

func main() {
	// 原始数据
	src := `{"code":"0","msg":"success","data":{"email":"", "message":"原始数据，一个字符串，一个待压缩的字符串，一个待压缩的数组，一个待压缩的map"}}`

	fmt.Println("原始数据: ", src)

	ret, _ := compress.Compress(compress.COMPRESS_GZIP, []byte(src))

	// base64后的压缩数据
	fmt.Println("压缩后数据: ", ret)

	// 解压缩
	ret_uncompress, _ := compress.UnCompress(compress.COMPRESS_GZIP, ret)

	fmt.Println("还原后数据: ", string(ret_uncompress))
	fmt.Printf("还原后数据等于原始数据: %t \n", string(ret_uncompress) == src)

	// 重复数据压缩率比较高
	src = repeatStr(src, 1000)
	fmt.Println("原始数据长度: ", len(src))

	ret, _ = compress.Compress(compress.COMPRESS_GZIP, []byte(src))
	fmt.Println("压缩后数据长度: ", len(ret))

	ret_uncompress, _ = compress.UnCompress(compress.COMPRESS_GZIP, ret)
	fmt.Printf("还原后数据等于原始数据: %t \n", string(ret_uncompress) == src)
}

```