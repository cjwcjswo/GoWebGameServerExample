package db

import (
	"GoWebGameServerExample/pkg/config"
	"github.com/go-xorm/xorm"
)

type DataStore struct {
	ormEngine *xorm.Engine
	shardName string
	shardId   int
}

func (store *DataStore) InitEngine(mysqlConfig *config.MySqlConfig, cacheConfig *config.CacheConfig, dataSource string) bool {
	if store.ormEngine = newEngine(mysqlConfig, cacheConfig, dataSource); store.ormEngine == nil {
		return false
	}
	return true
}

func (store *DataStore) GetShardId() int {
	return store.shardId
}

func (store *DataStore) GetShardName() string {
	return store.shardName
}

func (store *DataStore) Close() {
	if store.ormEngine != nil {
		_ = store.ormEngine.Close()
	}
}
