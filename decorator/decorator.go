package decorator

import "reflect"

// go 通用装饰器
func AddDecorator(decorator func(func()), target interface{}) func(...interface{}) []interface{} {
	// 判断是否是函数
	if reflect.TypeOf(target).Kind() != reflect.Func {
		panic("target is not a function")
	}

	args_nums := reflect.TypeOf(target).NumIn()
	ret_nums := reflect.TypeOf(target).NumOut()
	target_func_reflect := reflect.ValueOf(target)

	// 反射函数
	ret := func(args ...interface{}) (interface_ret []interface{}) {
		tmp_func := func() {
			func_args := make([]reflect.Value, args_nums)
			for idx, arg := range args {
				func_args[idx] = reflect.ValueOf(arg)
			}
			func_ret := target_func_reflect.Call(func_args)
			interface_ret = make([]interface{}, ret_nums)
			for idx, ret := range func_ret {
				interface_ret[idx] = ret.Interface()
			}
		}

		decorator(tmp_func)

		return interface_ret
	}

	return ret
}
