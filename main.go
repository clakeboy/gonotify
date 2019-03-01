package main

import (
	"fmt"
	"github.com/asdine/storm"
	"github.com/clakeboy/golib/utils"
	"gonotify/command"
	"gonotify/common"
	"gonotify/queue"
	"gonotify/service"
	"os"
)

//go:generate go-bindata -o=asset/asset.go -pkg=asset html/...

var _VERSION_ = "0.0.1"
var server *service.HttpServer
var out chan os.Signal

func main() {
	go utils.ExitApp(out, func(s os.Signal) {
		os.Remove(command.CmdPidName)
	})
	queue.Notify.Start()
	server.Start()
}

func init() {
	var err error
	command.InitCommand()

	common.Conf = common.NewYamlConfig(command.CmdConfFile)

	common.BDB, err = storm.Open(common.Conf.BDB.Path)

	if err != nil {
		fmt.Println("open database error:", err)
	}

	utils.WritePid(command.CmdPidName)
	out = make(chan os.Signal, 1)
	server = service.NewHttpServer(common.Conf.System.Ip+":"+common.Conf.System.Port, command.CmdDebug, command.CmdCross, command.CmdPProf)
	queue.Notify = queue.NewNotifyServer()
}
