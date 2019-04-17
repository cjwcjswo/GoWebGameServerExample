package model

import "GoWebGameServerExample/pkg/protocol"

type Version struct {
	Seq                 int    `xorm:"pk" json:"seq"`
	AppVersion          string `json:"app_version"`
	LatestServerVersion string `json:"latest_server_version"`
	IsActiveVersion     int    `json:"is_active_version"`
	AppId               string `json:"app_id"`
}

func (Version) TableName() string {
	return "TB_VERSION"
}

func (Version) GetShardKey() int {
	return protocol.SHARD_ID_COMMON
}
