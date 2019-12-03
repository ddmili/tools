package command

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
	"sync"
	"tools/console"
	"tools/ssh"
)

type Command struct {
	Host      string
	User      string
	Stdout    io.Reader
	Stderr    io.Reader
	Server    *Server
	client    *ssh.Client
	outDone   map[int]chan bool
	outDoneID int
	lock      *sync.RWMutex
}

// Message The message used by channel to transport log line by line
type Message struct {
	Host    string
	Content string
}

// NewCommand Create a new command
func NewCommand(server *Server) *Command {
	cmd := &Command{
		Host:   server.Hostname,
		User:   server.User,
		Server: server,
		lock:   new(sync.RWMutex),
	}
	if !strings.Contains(cmd.Host, ":") {
		cmd.Host = cmd.Host + ":" + strconv.Itoa(server.Port)
	}
	cmd.connect()
	return cmd
}

func (cmd *Command) Close() {
	cmd.client.Close()
}

func (cmd *Command) Clear() {
	//s := ioutil.NopCloser(cmd.Stdout)
	//s.Close()
	//e := ioutil.NopCloser(cmd.Stderr)
	//e.Close()
	//fmt.Printf("%+v", cmd.outDone)
	for _, o := range cmd.outDone {
		o <- true
	}
	//初始化
	cmd.outDone = map[int]chan bool{}
}

// connect 连接
func (cmd *Command) connect() {
	cmd.client = &ssh.Client{
		Host:     cmd.Host,
		User:     cmd.User,
		Password: cmd.Server.Password,
	}

	if err := cmd.client.Connect(); err != nil {
		fmt.Printf("[%s] unable to connect: %s", cmd.Host, err)
		return
	}
}

// AddOutDone 记录开启的日志协程
func (cmd *Command) AddOutDone(o chan bool) int {
	cmd.lock.Lock()
	defer cmd.lock.Unlock()
	cmd.outDoneID++
	cmd.outDone[cmd.outDoneID] = o
	return cmd.outDoneID
}

// AddOutDone 记录开启的日志协程
func (cmd *Command) RemoveOutDone(l int) {
	cmd.lock.Lock()
	defer cmd.lock.Unlock()
	delete(cmd.outDone, l)
}

// bindOut 绑定输出
func (cmd *Command) bindOut(input *io.Reader, prefix string, color int) {
	o := make(chan bool, 1)
	l := cmd.AddOutDone(o)
	defer cmd.RemoveOutDone(l)
	reader := bufio.NewReader(*input)
	for {
		line, err := reader.ReadString('\n')
		if err != nil || io.EOF == err {
			if err != io.EOF {
				panic(fmt.Sprintf("[%s] faield to execute command: %s", cmd.Host, err))
			}
			break
		}
		//用于关闭协程
		select {
		case <-o:
			return
		default:
		}
		line = prefix + line
		if color != 0 {
			line = console.ColorfulText(color, line)
		}

		cmd.Server.output <- Message{
			Host:    cmd.Server.ServerName,
			Content: line,
		}
	}
}

// Execute 执行脚本命令
func (cmd *Command) Execute(script string) {

	//关闭上一条命令的输出
	cmd.Clear()
	session, err := cmd.client.NewSession()
	if err != nil {
		fmt.Printf("[%s] unable to create session: %s", cmd.Host, err)
		return
	}
	defer session.Close()
	if err = session.RequestPty("xterm", 80, 40, *ssh.CreateTerminalModes()); err != nil {
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
	go cmd.bindOut(&cmd.Stdout, "", 0)
	go cmd.bindOut(&cmd.Stderr, "Error:", console.TextRed)
	if err = session.Start(script); err != nil {
		fmt.Printf("[%s] failed to execute command: %s", cmd.Host, err)
		return
	}
	if err = session.Wait(); err != nil {
		fmt.Printf("[%s] failed to wait command: %s", cmd.Host, err)
		return
	}
}
