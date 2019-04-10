package game

import (
	"GoWebGameServerExample/pkg/config"
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

	// Handler 세팅
	log.LocalLogger.Info("Handler Init Start!")
	handler.InitGameHandler()
	log.LocalLogger.Info("Handler Init Finish!")

	return server, true
}

func (server gameServer) StartServer() {
	log.LocalLogger.Info("Start Server!")

	gatewayHandler := handler.GetGateWayHandler()
	http.HandleFunc("/", gatewayHandler.Handle)
	err := http.ListenAndServe(server.GameServerConfig.Address, nil)
	if err != nil {
		log.LocalLogger.Error("Server Start Fail!", zap.String("Error", err.Error()))
	}
}
