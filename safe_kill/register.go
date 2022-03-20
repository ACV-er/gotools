package safe_kill

import "time"

var (
	beforeExitFuncs []func()
	nonexit         bool
	timeOut         = time.Second * 20
)

// 注册进程退出时的处理函数
func Register(f func()) {
	beforeExitFuncs = append(beforeExitFuncs, f)
}

// 调用后下一次关闭信号不退出系统
func KeepAliveOnce() {
	nonexit = true
}

// 设置超时时间
func SetTimeOut(t time.Duration) {
	timeOut = t
}
