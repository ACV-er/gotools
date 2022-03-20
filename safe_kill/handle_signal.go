package safe_kill

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func init() {
	go handleSignal()
}

var execOnceBeforeExitFuncs sync.Once

func execBeforeExitFuncs() chan bool {
	// 等待注册的退出函数结束并设置超时时间
	wg := sync.WaitGroup{}

	done := make(chan bool)

	execOnceBeforeExitFuncs.Do(func() {
		// 执行注册的退出函数
		wg.Add(len(beforeExitFuncs))
		for _, f := range beforeExitFuncs {
			tmp_func := f
			go func() {
				defer wg.Done()
				defer func() {
					if err := recover(); err != nil {
						fmt.Println(err)
					}

				}()

				tmp_func()
			}()
		}

		// 所有注册的退出函数执行完毕后设置done
		go func() {
			wg.Wait()
			close(done)
		}()
	})

	return done
}

func exit(code int) {
	done := execBeforeExitFuncs()

	// 等待所有注册的退出函数执行完毕或超时
	select {
	case <-done:
		// 执行完成
	case <-time.After(timeOut):
		// 超时退出
		os.Exit(1)
	}
	os.Exit(code)
}

// 监听系统关闭信号
func handleSignal() {
	// 监听系统关闭信号
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	for {
		// 阻塞等待系统关闭信号
		<-signalChan
		// 关闭系统
		if !nonexit {
			exit(0)
		}
		nonexit = false
	}
}

func GraceExit(code int) {
	exit(code)
}
