package config

import (
	"GoWebGameServerExample/pkg/log"
	"github.com/go-ini/ini"
	"go.uber.org/zap"
	"os"
	"path/filepath"
)

const (
	LOCAL = "Local"
	DEV   = "Dev"
	QA    = "Qa"
	LIVE  = "Live"
)

var serverVersion string

type AllConfig struct {
	GameServerConfig
	MySqlConfig
	RedisConfig
}

type GameServerConfig struct {
	Address string
}

type MySqlConfig struct {
	Address string
}

type RedisConfig struct {
	Address string
}

func LoadConfig(version string) (AllConfig, bool) {
	if CheckServerVersion(version) == false {
		return AllConfig{}, false
	} else {
		serverVersion = version
	}

	currentDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	configDir := currentDir + "/../../configs"

	if err != nil {
		return AllConfig{}, false
	}

	config, err := ini.Load(filepath.FromSlash(configDir + "/config" + version + ".ini"))
	if err != nil {
		return AllConfig{}, false
	}

	var allConfig AllConfig

	section := config.Section("GameServer")
	allConfig.GameServerConfig.Address = section.Key("Address").String()

	section = config.Section("Redis")
	allConfig.RedisConfig.Address = section.Key("Address").String()

	section = config.Section("MySQL")
	allConfig.MySqlConfig.Address = section.Key("Address").String()

	if checkConfig(&allConfig) == false {
		return AllConfig{}, false
	}

	writeInfoLog(&allConfig)
	return allConfig, true
}

func checkConfig(config *AllConfig) bool {
	if config.GameServerConfig.Address == "" {
		log.LocalLogger.Error("GameSererConfig address is empty!")
		return false
	}
	if config.RedisConfig.Address == "" {
		log.LocalLogger.Error("RedisConfig address is empty!")
		return false
	}
	if config.MySqlConfig.Address == "" {
		log.LocalLogger.Error("MysqlConfig address is empty!")
		return false
	}
	return true
}

func writeInfoLog(config *AllConfig) {
	log.LocalLogger.Info("GameServerConfig", zap.String("Address", config.GameServerConfig.Address))
	log.LocalLogger.Info("RedisConfig", zap.String("Address", config.RedisConfig.Address))
	log.LocalLogger.Info("MysqlConfig", zap.String("Address", config.MySqlConfig.Address))
}

func CheckServerVersion(version string) bool {
	switch version {
	case LOCAL, DEV, QA, LIVE:
		{
			return true
		}
	default:
		{
			return false
		}
	}
}

func GetServerVersion() string {
	return serverVersion
}
