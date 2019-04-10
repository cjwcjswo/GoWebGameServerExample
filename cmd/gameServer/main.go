package main

import (
	"GoWebGameServerExample/pkg/config"
	"GoWebGameServerExample/pkg/game"
	"GoWebGameServerExample/pkg/log"
	"flag"
	"go.uber.org/zap"
)

func main() {
	// Input Command Line Argument
	version := flag.String("version", config.LOCAL, "Server Version")
	flag.Parse()

	// Local Logger Setting
	if log.InitLocalLog() == false {
		log.LocalLogger.Error("Init local logger fail!")
		return
	}

	// Load Config
	loadConfig, result := config.LoadConfig(*version)
	if result == false {
		log.LocalLogger.Error("Config Load Fail!: ", zap.String("Version", *version))
		return
	}

	// Start GameServer
	server, result := game.NewGameServer(loadConfig)
	if result == false {
		log.LocalLogger.Error("NewGameServer Fail!: ", zap.String("Version", *version))
		return
	}
	server.StartServer()
}
