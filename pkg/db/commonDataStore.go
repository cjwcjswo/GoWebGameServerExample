package db

import (
	"GoWebGameServerExample/pkg/config"
	"github.com/go-xorm/xorm"
)

type CommonDataStore struct {
	ormEngine *xorm.Engine
	shardName string
	shardId   int
}

func (store *CommonDataStore) InitEngine(mysqlConfig *config.MySqlConfig, cacheConfig *config.CacheConfig, dataSource string) bool {
	if store.ormEngine = newEngine(mysqlConfig, cacheConfig, dataSource); store.ormEngine == nil {
		return false
	}
	return true
}

func (store *CommonDataStore) GetShardId() int {
	return store.shardId
}

func (store *CommonDataStore) GetShardName() string {
	return store.shardName
}

func (store *CommonDataStore) Close() {
	if store.ormEngine != nil {
		_ = store.ormEngine.Close()
	}
}

// Common Data Store Query
func (store CommonDataStore) SelectAll(beanList interface{}) error {
	return store.ormEngine.Find(beanList)
}

func (store CommonDataStore) Query(query string) ([]map[string]interface{}, error) {
	result, err := store.ormEngine.QueryInterface(query)
	if err != nil {
		return nil, err
	}
	return result, nil
}
