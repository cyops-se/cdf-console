package main

import (
	"flag"
	"log"
	"server/db"
	"server/devices"
	"server/flows"
	"server/logger"
	"server/routes"
	"server/types"
	"server/web"
)

var ctx types.Context
var GitVersion string
var GitCommit string

func main() {
	flag.StringVar(&ctx.Cmd, "cmd", "debug", "Windows service command (try 'usage' for more info)")
	flag.StringVar(&ctx.Wdir, "workdir", ".", "Sets the working directory for the process")
	flag.BoolVar(&ctx.Trace, "trace", false, "Prints traces of application engine to the console")
	flag.BoolVar(&ctx.Version, "v", false, "Prints the commit hash and exits")
	flag.Parse()

	routes.SysInfo.GitVersion = GitVersion
	routes.SysInfo.GitCommit = GitCommit

	logger.InitLogger(ctx)
	logger.Trace("SysInfo", "GetVersion: %s, GetCommit: %s", routes.SysInfo.GitVersion, routes.SysInfo.GitCommit)
	runEngine()
	log.Printf("Exiting ...")
}

func runEngine() {
	db.ConnectDatabase(ctx)
	go flows.Syslogd()
	go devices.CheckDefaultPresentWorker()
	go devices.CheckAvailabilityWorker()
	web.RunWeb()
}
