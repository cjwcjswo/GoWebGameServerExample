package db

import (
	"GoWebGameServerExample/pkg/config"
	"GoWebGameServerExample/pkg/protocol"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

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
