package db

import (
	"GoWebGameServerExample/pkg/config"
	"GoWebGameServerExample/pkg/log"
	"github.com/go-xorm/xorm"
	"github.com/go-xorm/xorm-redis-cache"
	"go.uber.org/zap"
	"time"
)

type InterfaceStore interface {
	InitEngine(mysqlConfig *config.MySqlConfig, cacheConfig *config.CacheConfig, dataSource string) bool
	GetShardName() string
	GetShardId() int
	Close()
}

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
	ormEngine.ShowSQL(true)

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
