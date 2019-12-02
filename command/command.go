package command

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
	"tools/console"
	"tools/ssh"
)

type Command struct {
	Host   string
	User   string
	Stdout io.Reader
	Stderr io.Reader
	Server *Server
	output chan Message
}

// Message The message used by channel to transport log line by line
type Message struct {
	Host    string
	Content string
}

// NewCommand Create a new command
func NewCommand(server *Server) (cmd *Command) {
	cmd = &Command{
		Host:   server.Hostname,
		User:   server.User,
		Server: server,
		output: server.output,
	}
	if !strings.Contains(cmd.Host, ":") {
		cmd.Host = cmd.Host + ":" + strconv.Itoa(server.Port)
	}
	return
}

// Execute 执行脚本命令
func (cmd *Command) Execute(script string) {

	client := &ssh.Client{
		Host:     cmd.Host,
		User:     cmd.User,
		Password: cmd.Server.Password,
	}

	if err := client.Connect(); err != nil {
		fmt.Printf("[%s] unable to connect: %s", cmd.Host, err)
		return
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		fmt.Printf("[%s] unable to create session: %s", cmd.Host, err)
		return
	}
	defer session.Close()

	if err := session.RequestPty("xterm", 80, 40, *ssh.CreateTerminalModes()); err != nil {
		fmt.Printf("[%s] unable to create pty: %v", cmd.Host, err)
		return
	}

	cmd.Stdout, err = session.StdoutPipe()
	if err != nil {
		fmt.Printf("[%s] redirect stdout failed: %s", cmd.Host, err)
		return
	}

	cmd.Stderr, err = session.StderrPipe()
	if err != nil {
		fmt.Printf("[%s] redirect stderr failed: %s", cmd.Host, err)
		return
	}

	go bindOutput(cmd.Host, cmd.output, &cmd.Stdout, "", 0)
	go bindOutput(cmd.Host, cmd.output, &cmd.Stderr, "Error:", console.TextRed)

	if err = session.Start(script); err != nil {
		fmt.Printf("[%s] failed to execute command: %s", cmd.Host, err)
		return
	}

	if err = session.Wait(); err != nil {
		fmt.Printf("[%s] failed to wait command: %s", cmd.Host, err)
		return
	}
}

// bindOutput 绑定输出信息
func bindOutput(host string, output chan Message, input *io.Reader, prefix string, color int) {
	reader := bufio.NewReader(*input)
	for {
		line, err := reader.ReadString('\n')
		if err != nil || io.EOF == err {
			if err != io.EOF {
				panic(fmt.Sprintf("[%s] faield to execute command: %s", host, err))
			}
			break
		}

		line = prefix + line
		if color != 0 {
			line = console.ColorfulText(color, line)
		}

		output <- Message{
			Host:    host,
			Content: line,
		}
	}
}
