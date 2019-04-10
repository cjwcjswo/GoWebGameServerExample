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
	Address       string
	RootName      string
	Password      string
	PoolSize      int
	RequestBuffer int
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
	allConfig.MySqlConfig.RootName = section.Key("RootName").String()
	allConfig.MySqlConfig.Password = section.Key("Password").String()
	allConfig.MySqlConfig.PoolSize, _ = section.Key("PoolSize").Int()
	allConfig.MySqlConfig.RequestBuffer, _ = section.Key("RequestBuffer").Int()

	if checkConfig(&allConfig) == false {
		return AllConfig{}, false
	}

	writeInfoLog(&allConfig)
	return allConfig, true
}

func checkConfig(config *AllConfig) bool {
	// Check GameServer
	if config.GameServerConfig.Address == "" {
		log.LocalLogger.Error("GameSererConfig address is empty!")
		return false
	}

	// Check Redis
	if config.RedisConfig.Address == "" {
		log.LocalLogger.Error("RedisConfig address is empty!")
		return false
	}

	// Check Mysql
	if config.MySqlConfig.Address == "" {
		log.LocalLogger.Error("MysqlConfig address is empty!")
		return false
	}
	if config.MySqlConfig.RootName == "" {
		log.LocalLogger.Error("MysqlConfig root name is empty!")
		return false
	}
	if config.MySqlConfig.PoolSize < 1 {
		log.LocalLogger.Error("MysqlConfig PoolSize < 1!", zap.Int("PoolSize", config.MySqlConfig.PoolSize))
		return false
	}
	if config.MySqlConfig.RequestBuffer < 1 {
		log.LocalLogger.Error("MysqlConfig RequestBuffer < 1!", zap.Int("RequestBuffer", config.MySqlConfig.RequestBuffer))
		return false
	}

	return true
}

func writeInfoLog(config *AllConfig) {
	log.LocalLogger.Info("GameServerConfig", zap.String("Address", config.GameServerConfig.Address))
	log.LocalLogger.Info("RedisConfig", zap.String("Address", config.RedisConfig.Address))
	log.LocalLogger.Info("MysqlConfig",
		zap.String("Address", config.MySqlConfig.Address),
		zap.String("RootName", config.MySqlConfig.RootName),
		zap.String("Password", config.MySqlConfig.Password),
		zap.Int("PoolSize", config.MySqlConfig.PoolSize),
		zap.Int("RequestBuffer", config.MySqlConfig.RequestBuffer))
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
