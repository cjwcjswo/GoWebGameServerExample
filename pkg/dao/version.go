package dao

import (
	"GoWebGameServerExample/pkg/db"
	"GoWebGameServerExample/pkg/log"
	"GoWebGameServerExample/pkg/model"
	"GoWebGameServerExample/pkg/protocol"
	"encoding/gob"
)

var version *Version

type Version struct {
	store *db.CommonDataStore
}

func NewVersion() *Version {
	version := new(Version)

	gob.Register(new(model.Version))
	storeInterface := db.FindEngineByShardId(protocol.SHARD_ID_COMMON)
	if storeInterface == nil {
		log.LocalLogger.Error("Version Init Fail - no store interface")
		return nil
	}
	version.store = storeInterface.(*db.CommonDataStore)
	return version
}

func GetVersion() *Version {
	if version == nil {
		version = NewVersion()
		return version
	}
	return version
}

func (dao *Version) SelectAll() ([]model.Version, error) {
	var result []model.Version
	err := dao.store.SelectAll(&result)
	return result, err
}
