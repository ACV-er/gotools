# 优雅关闭

* 在程序退出前执行传入的方法，用于回收资源

> 因为该模块只作为退出前执行资源回收和收尾工作，对注册的函数，均只会执行一次，用了一个sync.Once进行限制，并且不关注错误、panic等。

* 使用

``` go
package main

import (
	"fmt"
	"time"

	"github.com/ACV-er/gotools/safe_kill"
)

// 文件为测试用
func test() {
	fmt.Println(123)
	// 释放资源
}

func loneTime() {
	// 阻塞
	time.Sleep(time.Second * 5)
}

func main() {
	fmt.Println("进入main")
	// 注册方法，会在程序退出前执行
	// 限于kill ctrl+c方式关闭，kill -9或异常退出无法清理资源
	// 后置处理函数会并发执行，请保证任务完成后退出，函数退出视为处理完毕，避免有残留协程等
	safe_kill.Register(test)

	// 可以有多个后置处理函数，这里模拟超时
	// safe_kill.Register(loneTime)

	// 后置处理函数超时时间，超时直接退出，返回值为1，默认20秒
	safe_kill.SetTimeOut(time.Second * 3)

	// 忽略下一次退出信号(多次调用无效)，不会进行清理，在后置处理函数中无效
	safe_kill.KeepAliveOnce()

	// 模拟阻塞
	loneTime()

	// 终止程序，并在此之前执行后置处理函数，可在main函数自然退出前使用
	safe_kill.GraceExit(0)
}

```