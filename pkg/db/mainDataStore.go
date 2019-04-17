package db

import (
	"GoWebGameServerExample/pkg/config"
	"github.com/go-xorm/xorm"
)

type MainDataStore struct {
	ormEngine *xorm.Engine
	shardName string
	shardId   int
}

func (store *MainDataStore) InitEngine(mysqlConfig *config.MySqlConfig, cacheConfig *config.CacheConfig, dataSource string) bool {
	if store.ormEngine = newEngine(mysqlConfig, cacheConfig, dataSource); store.ormEngine == nil {
		return false
	}
	return true
}

func (store *MainDataStore) GetShardId() int {
	return store.shardId
}

func (store *MainDataStore) GetShardName() string {
	return store.shardName
}

func (store *MainDataStore) Close() {
	if store.ormEngine != nil {
		_ = store.ormEngine.Close()
	}
}
