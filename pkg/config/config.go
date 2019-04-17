package config

var Config InterfaceConfig

type InterfaceConfig interface {
	GetUpdateUrlGroup() map[string]string
	GetPatchFileDownloadUrl() string
}

const (
	PLATFORM_TYPE_ANDROID = "A"
	PLATFORM_TYPE_IOS     = "I"
)
