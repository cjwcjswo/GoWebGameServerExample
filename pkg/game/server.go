package game

import (
	"GoWebGameServerExample/pkg/config"
	"GoWebGameServerExample/pkg/db"
	"GoWebGameServerExample/pkg/handler"
	"GoWebGameServerExample/pkg/log"
	"go.uber.org/zap"
	"net/http"
)

type gameServer struct {
	config.AllConfig
}

func NewGameServer(config config.AllConfig) (*gameServer, bool) {
	server := new(gameServer)
	server.AllConfig = config

	// RDB 세팅
	log.LocalLogger.Info("Rdb Init Start!")
	if db.InitRdb(config.MySqlConfig) == false {
		return nil, false
	}
	log.LocalLogger.Info("Rdb Init Finish!")

	// Handler 세팅
	log.LocalLogger.Info("Handler Init Start!")
	handler.InitGameHandler()
	log.LocalLogger.Info("Handler Init Finish!")

	return server, true
}

func (server gameServer) StartServer() {
	defer server.Close()
	log.LocalLogger.Info("Start Server!")

	gatewayHandler := handler.GetGateWayHandler()
	http.HandleFunc("/", gatewayHandler.Handle)
	err := http.ListenAndServe(server.GameServerConfig.Address, nil)
	if err != nil {
		log.LocalLogger.Error("Server Start Fail!", zap.String("Error", err.Error()))
	}
}

func (server gameServer) Close() {
	log.LocalLogger.Info("End Server!")
	db.CloseRdb()
}
