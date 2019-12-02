package main

import (
	"flag"
	"fmt"
	"os"
	"tools/command"
	"tools/console"
)

var mossSep = "------------------------------------------------ \n"

var welcomeMessage = getWelcomeMessage() + console.ColorfulText(console.TextMagenta, mossSep)

//var filePath = flag.String("file", "", "-file=\"/var/log/*.log\"")
//var hostStr = flag.String("hosts", "", "-hosts=root@192.168.1.101,root@192.168.1.102")
var configFile = flag.String("conf", "", "-conf=conf.yaml")

//var tailFlags = flag.String("tail-flags", "--retry --follow=name", "flags for tail command, you can use -f instead if your server does't support `--retry --follow=name` flags")
//var slient = flag.Bool("slient", false, "-slient=false")

var Version = ""
var GitCommit = ""

// usageAndExit 退出
func usageAndExit(message string) {

	if message != "" {
		fmt.Fprintln(os.Stderr, message)
	}

	flag.Usage()
	fmt.Fprint(os.Stderr, "\n")

	os.Exit(1)
}

func main() {

	flag.Usage = func() {
		fmt.Fprint(os.Stderr, welcomeMessage)
		fmt.Fprint(os.Stderr, "Options:\n\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if *configFile == "" {
		usageAndExit("配置文件为空")
	}
	fmt.Println(welcomeMessage)
	config := command.ParseConf(*configFile)
	for name, s := range config.Server {
		fmt.Printf("请选择服务器 %s \n", name)
		s.Init()
	}
	fmt.Printf("\n%s\n", console.ColorfulText(console.TextCyan, mossSep))

}

func getWelcomeMessage() string {
	return `
 _____ 				  _
|_   _|  ___    ___  | |
  | |   / _ \  / _ \ | |
  | |  | (_) || (_) || |
  |_|   \___/  \___/ |_|
Author: song
Homepage: github.com/disciple-song/tools
Version: ` + Version + "(" + GitCommit + ")" + `
`
}
