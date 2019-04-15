package dao

import (
	"GoWebGameServerExample/pkg/db"
	"GoWebGameServerExample/pkg/log"
	"GoWebGameServerExample/pkg/model"
	"GoWebGameServerExample/pkg/protocol"
	"encoding/gob"
)

var versionDao *VersionDao

type VersionDao struct {
	store *db.CommonDataStore
}

func NewVersionDao() *VersionDao {
	versionDao := new(VersionDao)

	gob.Register(new(model.Version))
	storeInterface := db.FindEngineByShardId(protocol.SHARD_ID_COMMON)
	if storeInterface == nil {
		log.LocalLogger.Error("VersionDao Init Fail - no store interface")
		return nil
	}
	versionDao.store = storeInterface.(*db.CommonDataStore)
	return versionDao
}

func GetVersionDao() *VersionDao {
	if versionDao == nil {
		versionDao = NewVersionDao()
		return versionDao
	}
	return versionDao
}

func (dao *VersionDao) SelectAll() ([]model.Version, error) {
	var result []model.Version
	err := dao.store.SelectAll(&result)
	return result, err
}
