package config

import (
	"GoWebGameServerExample/pkg/log"
	"github.com/go-ini/ini"
	"go.uber.org/zap"
	"os"
	"path/filepath"
)

const (
	LOCAL = "LOCAL"
	DEV   = "DEV"
	QA    = "QA"
	LIVE  = "LIVE"
)

var serverVersion string

type InitConfig struct {
	GameServerConfig
	MySqlConfig
	CacheConfig
}

type GameServerConfig struct {
	Address string
}

type MySqlConfig struct {
	Address     string
	RootName    string
	Password    string
	MaxPoolSize int
}

type CacheConfig struct {
	Address  string
	Password string
}

func LoadConfig(version string) (InitConfig, bool) {
	if CheckServerVersion(version) == false {
		return InitConfig{}, false
	} else {
		serverVersion = version
	}
	applyServerVersion()

	log.LocalLogger.Info("LoadConfig!", zap.String("ServerVersion", version))
	currentDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return InitConfig{}, false
	}

	configDir := currentDir + "/../../configs"
	config, err := ini.Load(filepath.FromSlash(configDir + "/config" + version + ".ini"))
	if err != nil {
		return InitConfig{}, false
	}

	var allConfig InitConfig

	section := config.Section("GameServer")
	allConfig.GameServerConfig.Address = section.Key("Address").String()

	section = config.Section("MySQL")
	allConfig.MySqlConfig.Address = section.Key("Address").String()
	allConfig.MySqlConfig.RootName = section.Key("RootName").String()
	allConfig.MySqlConfig.Password = section.Key("Password").String()
	allConfig.MySqlConfig.MaxPoolSize, _ = section.Key("MaxPoolSize").Int()

	section = config.Section("Cache")
	allConfig.CacheConfig.Address = section.Key("Address").String()
	allConfig.CacheConfig.Password = section.Key("Password").String()

	if checkConfig(&allConfig) == false {
		return InitConfig{}, false
	}

	writeInfoLog(&allConfig)
	return allConfig, true
}

func checkConfig(config *InitConfig) bool {
	// Check GameServer
	if config.GameServerConfig.Address == "" {
		log.LocalLogger.Error("GameSererConfig address is empty!")
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
	if config.MySqlConfig.MaxPoolSize < 0 {
		log.LocalLogger.Error("MysqlConfig PoolSize < 0!", zap.Int("PoolSize", config.MySqlConfig.MaxPoolSize))
		return false
	}

	// Check Cache
	if config.CacheConfig.Address == "" {
		log.LocalLogger.Error("CacheConfig address is empty!")
		return false
	}
	return true
}

func writeInfoLog(config *InitConfig) {
	log.LocalLogger.Info("GameServerConfig", zap.String("Address", config.GameServerConfig.Address))
	log.LocalLogger.Info("MysqlConfig",
		zap.String("Address", config.MySqlConfig.Address),
		zap.String("RootName", config.MySqlConfig.RootName),
		zap.String("Password", config.MySqlConfig.Password),
		zap.Int("MaxPoolSize", config.MySqlConfig.MaxPoolSize))
	log.LocalLogger.Info("CacheConfig",
		zap.String("Address", config.CacheConfig.Address),
		zap.String("Password", config.CacheConfig.Password))
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

func applyServerVersion() {
	switch serverVersion {
	case LOCAL:
		Config = LocalConfig{}
	case DEV:
		Config = DevConfig{}
	case QA:
		Config = DevConfig{}
	case LIVE:
		Config = DevConfig{}
	}
}

func GetServerVersion() string {
	return serverVersion
}

func SetServerVersion(version string) {
	serverVersion = version
}
