package model

import "GoWebGameServerExample/pkg/protocol"

type PatchGameData struct {
	Seq            int `xorm:"pk"`
	ServerVersion  string
	DesignDataName string
}

func (PatchGameData) TableName() string {
	return "TB_PATCH_GAME_DATA"
}

func (PatchGameData) GetShardKey() int {
	return protocol.SHARD_ID_COMMON
}
