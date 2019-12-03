package service

import "fmt"

var serviceHandler = make(map[string]func(param []string) error)

// RegisterService 注册命令服务
func RegisterService(name string, fun func(param []string) error) {
	serviceHandler[name] = fun
}

// CheckService 检查是否有该命令
func CheckService(name string) (ok bool) {
	_, ok = serviceHandler[name]
	return
}

// Run 执行命令
func Run(cmd *CmdService) error {
	h, ok := serviceHandler[cmd.Name]
	if ok {
		return h(cmd.Param)
	}
	return fmt.Errorf("命令句柄不存在")
}
