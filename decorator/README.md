# 通用装饰器

* 反射实现的一个通用装饰器，性能较定制化的装饰器较差

``` go
package main

import (
	"fmt"
	"time"

	"github.com/ACV-er/gotools/decorator"
)

func print(s string) int {
	fmt.Println(s)

	sum := 0
	for i := 0; i < 1000; i++ {
		sum += i
	}

	return sum
}

func count_cost(target func()) {
	start_time := time.Now().UnixNano()

	target()

	end_time := time.Now().UnixNano()

	fmt.Printf("耗时: %d纳秒\n", (end_time - start_time))
}

func _count_cost(target func(string) int) func(string) int {
	return func(s string) int {
		start_time := time.Now().UnixNano()

		ret := target(s)

		end_time := time.Now().UnixNano()

		fmt.Printf("耗时: %d纳秒\n", (end_time - start_time))

		return ret
	}
}

func main() {
	fmt.Println("---------反射装饰器---------")
	print_with_count := decorator.AddDecorator(count_cost, print)
	ret := print_with_count("这是一个字符串\n")
	fmt.Printf("%#v\n", ret)

	fmt.Println("---------闭包装饰器---------")
	func_decorator := _count_cost(print)
	iret := func_decorator("这是一个字符串\n")
	fmt.Printf("%#v\n", iret)
}

```