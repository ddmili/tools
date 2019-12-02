package command

import (
	"fmt"
	"strings"
	"tools/console"
)

type Server struct {
	ServerName string `yaml:"server_name"`
	Hostname   string `yaml:"hostname"`
	Port       int    `yaml:"port"`
	User       string `yaml:"user"`
	Password   string `yaml:"password"`
	output     chan Message
	cmd        *Command
}

func (s *Server) Init() {
	s.cmd = NewCommand(s)
	s.output = make(chan Message, 255)
}

func (s *Server) Tail(file string) {
	script := fmt.Sprintf("tail -f %s", file)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Printf(console.ColorfulText(console.TextRed, "Error: %s\n"), err)
			}
		}()
		s.cmd.Execute(script)
	}()

	go func() {
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
	}()
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf(console.ColorfulText(console.TextRed, "Error: %s\n"), err)
		}
	}()
	cmd := NewCommand(s)

	cmd.Execute(fmt.Sprintf("tail -f %s", file))
}
