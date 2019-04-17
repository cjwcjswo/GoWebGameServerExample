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
	config.InitConfig

	gatewayHandler *handler.GatewayHandler
}

func NewGameServer(config config.InitConfig) (*gameServer, bool) {
	server := new(gameServer)
	server.InitConfig = config

	// RDB Setting
	log.LocalLogger.Info("Rdb Init Start!")
	if db.InitRdbEngine(config.MySqlConfig, config.CacheConfig) == false {
		return nil, false
	}
	log.LocalLogger.Info("Rdb Init Finish!")

	// Gateway Setting
	log.LocalLogger.Info("Gateway Init Start!")
	server.gatewayHandler = handler.NewGatewayHandler()
	log.LocalLogger.Info("Gateway Init Finish!")

	return server, true
}

func (server gameServer) StartServer() {
	defer server.Close()
	log.LocalLogger.Info("Start Server!")

	// Enroll Gateway Handler
	http.HandleFunc("/", server.gatewayHandler.Handle)

	// Start Listen And Serve
	err := http.ListenAndServe(server.GameServerConfig.Address, nil)
	if err != nil {
		log.LocalLogger.Error("Server Start Fail!", zap.String("Error", err.Error()))
	}
}

func (server gameServer) Close() {
	log.LocalLogger.Info("End Server!")
	db.CloseRdbEngine()
}
