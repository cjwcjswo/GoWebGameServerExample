package handler

import (
	"GoWebGameServerExample/pkg/config"
	"GoWebGameServerExample/pkg/dao"
	"GoWebGameServerExample/pkg/log"
	"GoWebGameServerExample/pkg/model"
	"GoWebGameServerExample/pkg/protocol"
	"backup/goTcpLib"
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
)

const API_NAME_CHECK_UPDATE = "CheckUpdate"

const (
	UPDATE_TYPE_NORMAL = 0
	UPDATE_TYPE_FORCE  = 1
	UPDATE_TYPE_SERVER = 2
	UPDATE_TYPE_CLIENT = 3
)

type checkUpdateParams struct {
	AppVersion          string `json:"app_version"`
	Platform            string `json:"platform"`
	ClientServerVersion string `json:"server_version"`
}

type checkUpdateResponse struct {
	ResourceMetaDto     string `json:"resource_meta_dto"`
	Path                string `json:"path"`
	UpdateState         int    `json:"update_state"`
	UpdateUrl           string `json:"update_url"`
	LatestServerVersion string `json:"latest_server_version"`
}

func (params *checkUpdateParams) Decoding(data *json.RawMessage) protocol.ErrorCode {
	if err := json.Unmarshal(*data, &params); err != nil {
		goTcpLib.Logger.Error("checkUpdateParams - Decoding Error", zap.String("Error", err.Error()))
		return protocol.ERROR_MYSQL_FAIL
	}
	return protocol.ERROR_SUCCESS
}

type checkUpdateHandler struct{}

func NewCheckUpdateHandler() *checkUpdateHandler {
	return new(checkUpdateHandler)
}

func (handler *checkUpdateHandler) Handle(writer http.ResponseWriter, request *http.Request, gatewayRequest *gatewayRequest) (interface{}, protocol.ServerError) {
	var params checkUpdateParams
	if errorCode := params.Decoding(gatewayRequest.Params); errorCode != protocol.ERROR_SUCCESS {
		return nil, protocol.ServerError{ErrorCode: errorCode}
	}

	updateUrl := ""
	updateState := UPDATE_TYPE_NORMAL
	forceUpdate := false
	versionList, errorCode := checkAppUpdate(params.AppVersion)
	if errorCode != protocol.ERROR_SUCCESS {
		// Fail Fetch Version
		result, errorCode := getUpdateUrl(params.Platform)
		if errorCode != protocol.ERROR_SUCCESS {
			return nil, protocol.ServerError{ErrorCode: errorCode}
		}
		updateUrl = result
		updateState = UPDATE_TYPE_FORCE
		forceUpdate = true
	}

	latestServerVersion, errorCode := getServerVersion(versionList, params.AppVersion)
	if errorCode != protocol.ERROR_SUCCESS {
		return nil, protocol.ServerError{ErrorCode: errorCode}
	}

	if !forceUpdate && params.ClientServerVersion < latestServerVersion.AppVersion {
		updateState = UPDATE_TYPE_SERVER
	}
	if !forceUpdate {
		for _, version := range versionList {
			if version.AppVersion > params.AppVersion {
				updateState = UPDATE_TYPE_CLIENT
				break
			}
		}
	}

	// Skip...
	resourceMetaDto := "Test"

	return checkUpdateResponse{
		ResourceMetaDto:     resourceMetaDto,
		Path:                config.UrlConfig.GetPatchFileDownloadUrl(),
		UpdateState:         updateState,
		UpdateUrl:           updateUrl,
		LatestServerVersion: latestServerVersion.LatestServerVersion,
	}, protocol.ServerError{ErrorCode: protocol.ERROR_SUCCESS}
}

func (handler *checkUpdateHandler) GetApiName() string {
	return API_NAME_CHECK_UPDATE
}

func checkAppUpdate(appVersion string) ([]model.Version, protocol.ErrorCode) {
	versionList, err := dao.GetVersionDao().SelectAll()
	if err != nil {
		log.LocalLogger.Error("SelectAll Error", zap.String("Error", err.Error()))
		return nil, protocol.ERROR_MYSQL_FAIL
	}

	isOk := false
	for _, version := range versionList {
		if appVersion == version.AppVersion {
			isOk = version.IsActiveVersion == 1
			break
		}
	}

	if !isOk {
		return nil, protocol.ERROR_YOU_NEED_TO_UPDATE_APP_VERSION
	}
	return versionList, protocol.ERROR_SUCCESS
}

func getUpdateUrl(platform string) (string, protocol.ErrorCode) {
	updateUrlGroup := config.UrlConfig.GetUpdateUrlGroup()
	result, isExist := updateUrlGroup[platform]
	if !isExist {
		return "", protocol.ERROR_WRONG_PLATFORM
	}
	return result, protocol.ERROR_SUCCESS
}

func getServerVersion(versionList []model.Version, appVersion string) (model.Version, protocol.ErrorCode) {
	for _, version := range versionList {
		if version.AppVersion == appVersion {
			return version, protocol.ERROR_SUCCESS
		}
	}
	return model.Version{}, protocol.ERROR_APP_VERSION_DATA_NOT_FOUND
}
