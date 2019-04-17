package db

import (
	"GoWebGameServerExample/pkg/config"
	"GoWebGameServerExample/pkg/log"
	"GoWebGameServerExample/pkg/protocol"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	xormrediscache "github.com/go-xorm/xorm-redis-cache"
	"go.uber.org/zap"
	"time"
)

type InterfaceStore interface {
	InitEngine(mysqlConfig *config.MySqlConfig, cacheConfig *config.CacheConfig, dataSource string) bool
	GetShardName() string
	GetShardId() int
	Close()
}

var storeList []InterfaceStore

func InitRdbEngine(mysqlConfig config.MySqlConfig, cacheConfig config.CacheConfig) bool {
	storeList = []InterfaceStore{
		&MainDataStore{shardName: "DANCEVILLE_" + config.GetServerVersion() + "_SHARD_00", shardId: protocol.SHARD_ID_MAIN},
		&DataStore{shardName: "DANCEVILLE_" + config.GetServerVersion() + "_SHARD_01", shardId: protocol.SHARD_ID_DATA_1},
		&DataStore{shardName: "DANCEVILLE_" + config.GetServerVersion() + "_SHARD_02", shardId: protocol.SHARD_ID_DATA_2},
		&DataStore{shardName: "DANCEVILLE_" + config.GetServerVersion() + "_SHARD_03", shardId: protocol.SHARD_ID_DATA_3},
		&CommonDataStore{shardName: "DANCEVILLE_" + config.GetServerVersion() + "_SHARD_50", shardId: protocol.SHARD_ID_COMMON},
	}
	storeListLength := len(storeList)

	// Engine Setting
	for i := 0; i < storeListLength; i++ {
		dataSource := parsingInfoToDataSourceName(mysqlConfig.RootName, mysqlConfig.Password, mysqlConfig.Address, storeList[i].GetShardName())

		if storeList[i].InitEngine(&mysqlConfig, &cacheConfig, dataSource) == false {
			return false
		}
	}
	return true
}

func CloseRdbEngine() {
	for _, rdbEngine := range storeList {
		rdbEngine.Close()
	}
	storeList = nil
}

func FindEngineByShardId(shardId int) InterfaceStore {
	for _, store := range storeList {
		if store.GetShardId() == shardId {
			return store
		}
	}
	return nil
}

func parsingInfoToDataSourceName(rootName string, password string, address string, dbName string) string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s", rootName, password, address, dbName)
}

func newEngine(mysqlConfig *config.MySqlConfig, cacheConfig *config.CacheConfig, dataSource string) *xorm.Engine {
	var err error

	// Engine Create
	ormEngine, err := xorm.NewEngine("mysql", dataSource)
	if err != nil {
		log.LocalLogger.Error("Init Rdb Fail! xorm.NewEngine", zap.String("Error", err.Error()), zap.String("DataSource", dataSource))
		return nil
	}

	// Ping Check
	if err = ormEngine.Ping(); err != nil {
		log.LocalLogger.Error("Init Rdb Fail! xorm.NewEngine", zap.String("Error", err.Error()), zap.String("DataSource", dataSource))
		return nil
	}
	log.LocalLogger.Info("Init Max Pool Size", zap.String("DataSource", dataSource), zap.Int("MaxPoolSize", mysqlConfig.MaxPoolSize))
	ormEngine.SetMaxOpenConns(mysqlConfig.MaxPoolSize)
	//ormEngine.ShowSQL(true)

	// Cache Setting
	cache := xormrediscache.NewRedisCacher(cacheConfig.Address, cacheConfig.Password, time.Duration(24*time.Hour), nil)
	if cache == nil {
		log.LocalLogger.Error("Cache Init Error!", zap.String("Address", cacheConfig.Address), zap.String("Password", cacheConfig.Password))
		return nil
	}
	pool, err := cache.GetPool()
	if err != nil {
		log.LocalLogger.Error("Cache Get Pool Error!", zap.String("Address", cacheConfig.Address), zap.String("Password", cacheConfig.Password))
		return nil
	}
	_, err = pool.Dial()
	if err != nil {
		log.LocalLogger.Error("Cache Ping!", zap.String("Address", cacheConfig.Address), zap.String("Password", cacheConfig.Password))
		return nil
	}
	ormEngine.SetDefaultCacher(cache)

	log.LocalLogger.Info("Engine Setting Finish", zap.String("DataSource", dataSource))
	return ormEngine
}
