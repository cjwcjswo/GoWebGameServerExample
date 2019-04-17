package service

import (
	"GoWebGameServerExample/pkg/dao"
	"GoWebGameServerExample/pkg/log"
	"GoWebGameServerExample/pkg/model"
	"go.uber.org/zap"
)

var patchGameData *PatchGameData

type PatchGameData struct {
	patchGameDataGroup            []model.PatchGameData
	patchGameDataGroupByKey       map[string][]model.PatchGameData
	minPatchGameDataServerVersion string
}

func NewPatchGameData() *PatchGameData {
	result := new(PatchGameData)

	return result
}

func (service *PatchGameData) init() {
	if service.patchGameDataGroupByKey != nil {
		return
	}

	service.patchGameDataGroupByKey = make(map[string][]model.PatchGameData)
	if service.patchGameDataGroup == nil {
		service.patchGameDataGroup = make([]model.PatchGameData, 0, 60)
	}
	tableName := model.PatchGameData{}.TableName()
	query := "SELECT * FROM " + tableName + " ORDER BY design_data_name, server_version DESC;"

	patchGameDao := dao.GetPatchGameData()
	patchGameDataGroup, err := patchGameDao.GetPatchGameDataGroup(query)
	if err != nil {
		log.LocalLogger.Error("Init PatchGameData Error", zap.String("Error", err.Error()))
	}
	if patchGameDataGroup == nil || len(patchGameDataGroup) == 0 {
		return
	}

	service.minPatchGameDataServerVersion = patchGameDataGroup[0].ServerVersion
	for _, patchGameData := range patchGameDataGroup {
		dataMap, isExist := service.patchGameDataGroupByKey[patchGameData.DesignDataName]
		if !isExist {
			dataMap = make([]model.PatchGameData, 0, 5)
		}
		service.patchGameDataGroupByKey[patchGameData.DesignDataName] = append(dataMap, patchGameData)
		service.patchGameDataGroup = append(service.patchGameDataGroup, patchGameData)

		if service.minPatchGameDataServerVersion <= patchGameData.ServerVersion {
			continue
		}
		service.minPatchGameDataServerVersion = patchGameData.ServerVersion
	}
}

func GetPatchGameData() *PatchGameData {
	if patchGameData == nil {
		patchGameData := NewPatchGameData()
		patchGameData.init()
		return patchGameData
	}
	return patchGameData
}

func (service PatchGameData) GetPatchGameDataGroup() []model.PatchGameData {
	return service.patchGameDataGroup
}

func (service PatchGameData) GetMinPatchGameDataServerVersion() string {
	return service.minPatchGameDataServerVersion
}
