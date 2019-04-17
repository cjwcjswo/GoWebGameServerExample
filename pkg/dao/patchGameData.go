package dao

import (
	"GoWebGameServerExample/pkg/db"
	"GoWebGameServerExample/pkg/log"
	"GoWebGameServerExample/pkg/model"
	"GoWebGameServerExample/pkg/protocol"
	"encoding/gob"
	"strconv"
)

var patchGameData *PatchGameData

type PatchGameData struct {
	store *db.CommonDataStore
}

func NewPatchGameData() *PatchGameData {
	patchGameData := new(PatchGameData)

	gob.Register(new(model.PatchGameData))
	storeInterface := db.FindEngineByShardId(protocol.SHARD_ID_COMMON)
	if storeInterface == nil {
		log.LocalLogger.Error("PatchGameData Init Fail - no store interface")
		return nil
	}
	patchGameData.store = storeInterface.(*db.CommonDataStore)
	return patchGameData
}

func GetPatchGameData() *PatchGameData {
	if patchGameData == nil {
		patchGameData = NewPatchGameData()
		return patchGameData
	}
	return patchGameData
}

func (dao *PatchGameData) GetPatchGameDataGroup(query string) ([]model.PatchGameData, error) {
	result, err := dao.store.Query(query)
	if err != nil {
		return nil, err
	}

	length := len(result)
	list := make([]model.PatchGameData, length)

	for i := 0; i < length; i++ {
		data := result[i]
		seqByte := data["seq"].([]byte)
		sequence, _ := strconv.Atoi(string(seqByte))
		list[i] = model.PatchGameData{
			Seq:            sequence,
			ServerVersion:  string(data["server_version"].([]byte)),
			DesignDataName: string(data["design_data_name"].([]byte)),
		}
	}
	return list, nil
}
