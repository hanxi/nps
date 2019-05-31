package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/ProtonMail/go-autostart"
	"github.com/cnlh/nps/client"
	"github.com/cnlh/nps/lib/common"
	"github.com/cnlh/nps/lib/daemon"
	"github.com/cnlh/nps/lib/version"
	"github.com/cnlh/nps/vender/github.com/astaxie/beego/logs"
	"github.com/getlantern/systray"
	"github.com/hanxi/nps/cmd/npc_windows/icon"
	"github.com/monochromegane/conflag"
	"github.com/skratchdot/open-golang/open"
)

var (
	serverAddr string
	verifyKey  string
	logType    string
	connType   string
	proxyURL   string
	logLevel   string
	logPath    string
)

var confFile = "npc.toml"
var flags *flag.FlagSet

func updateTips() {
	tips := fmt.Sprintf("server='%s'\nvkey='%s'", serverAddr, verifyKey)
	systray.SetTooltip(tips)
}

func onReady() {
	systray.SetIcon(icon.Data)
	systray.SetTitle("npc")
	updateTips()

	mChecked := systray.AddMenuItem("Auto Startup", "Auto Startup npc on boot")
	filename := os.Args[0] // get command line first parameter
	app := &autostart.App{
		Name:        "npc",
		DisplayName: "npc",
		Exec:        []string{filename},
	}
	if app.IsEnabled() {
		mChecked.Check()
	}

	mOpenConfig := systray.AddMenuItem("OpenConfig", "Open npc Config file")

	go func() {
		for {
			select {
			case <-mChecked.ClickedCh:
				if mChecked.Checked() {
					if err := app.Disable(); err != nil {
						logs.Error("Disable Autostart Failed.")
					} else {
						mChecked.Uncheck()
					}
				} else {
					if err := app.Enable(); err != nil {
						logs.Error("Enable Autostart Failed.")
					} else {
						mChecked.Check()
					}
				}
			case <-mOpenConfig.ClickedCh:
				{
					confPath, err := getConfPath()
					if err == nil {
						open.Run(filepath.Dir(confPath))
					}
				}
			}
		}
	}()

	mQuit := systray.AddMenuItem("Quit", "Quit npc")
	go func() {
		<-mQuit.ClickedCh
		systray.Quit()
	}()

	go start()
}

func onExit() {
	// clean up here
}

func main() {
	flags = getFlags()

	systray.Run(onReady, onExit)
}

func getConfPath() (string, error) {
	file, err := os.Open(confFile)
	defer func() {
		file.Close()
	}()

	if os.IsNotExist(err) {
		os.MkdirAll(filepath.Dir(confFile), os.ModePerm)
		s := "#server='xx.com:8080'\r\n#vkey='xxxx'\r\n"
		ioutil.WriteFile(confFile, []byte(s), os.ModePerm)
	}

	return confFile, nil
}

func getFlags() *flag.FlagSet {
	flags := flag.NewFlagSet("npc", flag.ContinueOnError)
	flags.StringVar(&serverAddr, "server", "", "Server addr (ip:port)")
	flags.StringVar(&verifyKey, "vkey", "", "Authentication key")
	flags.StringVar(&logType, "log", "stdout", "Log output mode（stdout|file）")
	flags.StringVar(&connType, "type", "tcp", "Connection type with the server（kcp|tcp）")
	flags.StringVar(&logLevel, "log_level", "7", "log level 0~7")
	flags.StringVar(&logPath, "log_path", "npc.log", "npc log path")
	return flags
}

func start() {
	confPath, err := getConfPath()
	if err == nil {
		if confArgs, err := conflag.ArgsFrom(confPath); err == nil {
			flags.Parse(confArgs)
		} else {
			logs.Info("parse error:%s", err.Error())
		}
	}

	if len(os.Args) > 1 {
		err = flags.Parse(os.Args[1:])
		if err != nil {
			logs.Error("args error.")
			return
		}
	}

	daemon.InitDaemon("npc", common.GetRunPath(), common.GetTmpPath())
	logs.EnableFuncCallDepth(true)
	logs.SetLogFuncCallDepth(3)
	if logType == "stdout" {
		logs.SetLogger(logs.AdapterConsole, `{"level":`+logLevel+`,"color":true}`)
	} else {
		logs.SetLogger(logs.AdapterFile, `{"level":`+logLevel+`,"filename":"`+logPath+`","daily":false,"maxlines":100000,"color":true}`)
	}

	logs.Info("the version of client is %s, the core version of client is %s", version.VERSION, version.GetVersion())
	go func() {
		for {
			updateTips()
			logs.Info("serverAddr:%s, verifyKey:%s", serverAddr, verifyKey)
			client.NewRPClient(serverAddr, verifyKey, connType, proxyURL, nil).Start()
			logs.Info("It will be reconnected in five seconds")
			time.Sleep(time.Second * 5)
		}
	}()
}
