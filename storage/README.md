# 一些常用的存储结构

> * **bitmap** 并发安全的bitmap
> * NewBitMap(size uint64, lock_granularity uint64) *BitMap
>   > 返回一个容量为size的bitmap，lock_granularity的意义为多少个bit共用一把锁。lock_granularity为0表示不加锁，即非线程安全的位图。
> * (b *BitMap) Get(pos uint64) bool 获取pos位置的值
> * (b *BitMap) Set(pos uint64, value bool) 设置pos位置的值

``` go
package main

import (
	"fmt"

	"github.com/ACV-er/gotools/storage"
)

func main() {
	bm := storage.NewBitMap(10000, 1000)
	bm.Set(15, true)
	bm.Set(50, true)

	// 15 true 16 false
	fmt.Printf("15 %#v 16 %#v\n", bm.Get(15), bm.Get(16))
	bm.Set(50, false)

	// 50 false
	fmt.Printf("50 %#v\n", bm.Get(50))
}

```