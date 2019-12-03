package command

import (
	"fmt"
	"strings"
	"tools/console"
)

type Server struct {
	ServerName  string `yaml:"server_name"`
	Hostname    string `yaml:"hostname"`
	Port        int    `yaml:"port"`
	User        string `yaml:"user"`
	Password    string `yaml:"password"`
	output      chan Message
	outputState int //标示是否有输出（有命令在执行）
	cmd         *Command
}

func (s *Server) Init() {
	s.output = make(chan Message, 255)
	s.cmd = NewCommand(s)
	go s.BindOutput()
}

// Use 选中当前服务器
func (s *Server) Use() {
	fmt.Println("当前服务器", s.ServerName)
	s.Pwd()
}

// Pwd 查看当前目录
func (s *Server) Ping() {
	s.Execute("pwd")
}

// Pwd 查看当前目录
func (s *Server) Pwd() {
	s.Execute("pwd")
}

// Pm2Ls 查看项目
func (s *Server) Pm2Ls() {
	s.Execute("pm2 ls")
}

// Tail 查看日志
func (s *Server) Tail(file string) {
	s.Execute(fmt.Sprintf("tail -f %s", file))
}

// BindOutput 绑定输出
func (s *Server) BindOutput() {
	for output := range s.output {
		content := strings.Trim(output.Content, "\r\n")
		// 去掉文件名称输出
		if content == "" || (strings.HasPrefix(content, "==>") && strings.HasSuffix(content, "<==")) {
			continue
		}

		fmt.Printf(
			"%s %s %s\n",
			console.ColorfulText(console.TextGreen, output.Host),
			console.ColorfulText(console.TextYellow, "->"),
			content,
		)
	}
}

func (s *Server) Execute(script string) {

	for i := 0; i < s.outputState; i++ {
		s.output <- Message{Host: s.cmd.Host, Content: "close"}
	}
	s.outputState = 0
	script = fmt.Sprintf(`
cd /data/html
%s
`, script)
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf(console.ColorfulText(console.TextRed, "Error: %s\n"), err)
		}
	}()
	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Printf(console.ColorfulText(console.TextRed, "Error: %s\n"), err)
			}
		}()
		s.cmd.Execute(script)
	}()

}
