package db

import (
	"GoWebGameServerExample/pkg/config"
	"GoWebGameServerExample/pkg/log"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"go.uber.org/zap"
	"strings"
)

type rdbEngine struct {
	shardInfo
	ormEngine *xorm.Engine
}

func (engine *rdbEngine) init(dataSource string, info shardInfo) bool {
	var err error
	engine.ormEngine, err = xorm.NewEngine("mysql", dataSource)
	if err != nil {
		log.LocalLogger.Error("Init Rdb Fail! xorm.NewEngine", zap.String("Error", err.Error()), zap.String("DataSource", dataSource))
		return false
	}
	engine.shardInfo = info
	log.LocalLogger.Info("Engine Setting", zap.String("DataSource", dataSource), zap.String("ShardName", engine.shardName))
	return true
}

type shardInfo struct {
	shardName string
	shardNum  int
}

var engineList []*rdbEngine
var rdbExitChan chan struct{}

func InitRdb(config config.MySqlConfig) bool {
	shardInfoList := getShardNameList()
	engineLength := len(shardInfoList)
	engineList = make([]*rdbEngine, engineLength)
	rdbExitChan = make(chan struct{}, 1)

	// Engine Setting
	for i := 0; i < engineLength; i++ {
		engineList[i] = new(rdbEngine)
		dataSource := parsingInfoToDataSourceName(config.RootName, config.Password, config.Address, shardInfoList[i].shardName)
		engineList[i].init(dataSource, shardInfoList[i])
	}

	// Worker Setting
	for i := 0; i < config.PoolSize; i++ {
		rdbWorker := new(rdbWorker)
		rdbWorker.initRdbWorker(config.RequestBuffer, rdbExitChan)
		rdbWorker.start()
	}

	return true
}

func CloseRdb() {
	if engineList != nil {
		for _, engine := range engineList {
			if err := engine.ormEngine.Close(); err != nil {
				log.LocalLogger.Error("CloseRdb - Close Fail", zap.String("Error", err.Error()))
			}
		}
	}
	close(rdbExitChan)
}

func parsingInfoToDataSourceName(rootName string, password string, address string, dbName string) string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s", rootName, password, address, dbName)
}

func getShardNameList() []shardInfo {
	return []shardInfo{
		{shardName: "DANCEVILLE_" + strings.ToUpper(config.GetServerVersion()) + "_SHARD_00", shardNum: 0},
		{shardName: "DANCEVILLE_" + strings.ToUpper(config.GetServerVersion()) + "_SHARD_01", shardNum: 1},
		{shardName: "DANCEVILLE_" + strings.ToUpper(config.GetServerVersion()) + "_SHARD_02", shardNum: 2},
		{shardName: "DANCEVILLE_" + strings.ToUpper(config.GetServerVersion()) + "_SHARD_03", shardNum: 3},
		{shardName: "DANCEVILLE_" + strings.ToUpper(config.GetServerVersion()) + "_SHARD_50", shardNum: 50},
	}
}

type rdbRequest struct {
	command string
	object  []interface{}
}

type rdbWorker struct {
	requestChan chan rdbRequest
	exitChan    chan struct{}
}

func (worker *rdbWorker) initRdbWorker(requestBuffer int, exitChan chan struct{}) {
	worker.requestChan = make(chan rdbRequest, requestBuffer)
	worker.exitChan = exitChan
}

func (worker *rdbWorker) start() {
	go worker.startGoroutine()
}

func (worker *rdbWorker) startGoroutine() {
	for {
		if worker.startGoroutineImpl() {
			log.LocalLogger.Info("RdbWorkerGoroutine End....")
			break
		}
	}
}

func (worker *rdbWorker) startGoroutineImpl() bool {
	log.LocalLogger.Info("RdbWorkerGoroutine Start....")
	isWantedTermination := false
	for {
		if isWantedTermination {
			break
		}
		select {
		case <-worker.exitChan:
			{
				isWantedTermination = true
				return isWantedTermination
			}
		}

	}
	return isWantedTermination
}
