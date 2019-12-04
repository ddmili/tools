package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"tools/command"
	"tools/console"
	"tools/service"
)

var mossSep = "------------------------------------------------ \n"

var welcomeMessage = getWelcomeMessage() + console.ColorfulText(console.TextMagenta, mossSep)

//var configFile = flag.String("conf", "", "-conf=conf.yaml")
var configFile = "server.yaml"

var Version = "v0.0.1"
var GitCommit = "v0.0.1"

// usageAndExit 退出
func usageAndExit(message string) {

	if message != "" {
		fmt.Fprintln(os.Stderr, message)
	}

	flag.Usage()
	fmt.Fprint(os.Stderr, "\n")

	os.Exit(1)
}

// currentServer 当前服务
var currentServer *command.Server

var config *command.Yaml

func main() {

	flag.Usage = func() {
		fmt.Fprint(os.Stderr, welcomeMessage)
		fmt.Fprint(os.Stderr, "Options:\n\n")
		flag.PrintDefaults()
	}

	//flag.Parse()

	//if *configFile == "" {
	//	usageAndExit("配置文件为空")
	//}
	fmt.Println(welcomeMessage)
	config = command.ParseConf(configFile)
	fmt.Println("请选择服务器:")
	for _, v := range config.Server {
		fmt.Printf("   %s \n", v.ServerName)
	}

	registerHandler()

	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("\n>>")
	for {

		cmdString, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}

		cmd, err := service.AnalysisCmd(cmdString)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		cmd.Run()
	}
	fmt.Printf("\n%s\n", console.ColorfulText(console.TextCyan, mossSep))

}

// registerHandler 注册命令服务
func registerHandler() {
	service.RegisterService("change", func(param []string) error {
		if len(param) != 1 {
			return fmt.Errorf("请选中服务器名称")
		}
		s, ok := config.Server[param[0]]
		if !ok {
			return fmt.Errorf("服务器名称错误")
		}
		currentServer = &s
		currentServer.Init()
		currentServer.Ping()
		fmt.Println("切换服务器-" + currentServer.ServerName)
		return nil
	})
	service.RegisterService("pm2", func(param []string) error {
		if currentServer == nil {
			return fmt.Errorf("请选择服务器")
		}
		currentServer.Execute("pm2 " + strings.Join(param, " "))
		return nil
	})
}

func getWelcomeMessage() string {
	return `
 _____ 		      _    
|_   _|  ___    ___  | |
  | |   / _ \  / _ \ | |
  | |  | (_) || (_) || | 
  |_|   \___/  \___/ |_|   
Author: song
Homepage: github.com/disciple-song/tools
Version: ` + Version + "(" + GitCommit + ")" + `
`
}
