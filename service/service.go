// 服务命令
package service

import (
	"fmt"
	"strings"
)

type CmdInterface interface {
	Run()
}

// CmdService 命令服务
type CmdService struct {
	Name  string
	Param []string
}

// AddCmd 新增命令参数
func (s *CmdService) AddCmd(val string) {
	s.Param = append(s.Param, val)
}

// Run 执行命令服务
func (s *CmdService) Run() {
	err := Run(s)
	if err != nil {
		fmt.Println(err.Error())
	}
}

// analysisCmd 分析命令
func AnalysisCmd(cmdString string) (CmdService, error) {
	var cmd CmdService
	cmdString = strings.TrimSuffix(cmdString, "\n")
	cmdString = strings.TrimSuffix(cmdString, "\r")
	cmdString = strings.Trim(cmdString, " ")
	if cmdString == "" {
		return cmd, fmt.Errorf("请输入正确的指令~")
	}
	cmdArr := strings.Split(cmdString, " ")
	//fmt.Printf("打印命令参数 ", cmdArr, len(cmdArr))
	for _, v := range cmdArr {
		if v == "" {
			continue
		}
		if cmd.Name == "" {
			if !CheckService(v) {
				return cmd, fmt.Errorf("请输入正确的指令~")
			}
			cmd.Name = v
		} else {
			cmd.AddCmd(v)
		}
	}
	return cmd, nil

}
