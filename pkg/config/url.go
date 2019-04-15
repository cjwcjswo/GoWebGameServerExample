package config

var UrlConfig InterfaceUrlConfig

type InterfaceUrlConfig interface {
	GetUpdateUrlGroup() map[string]string
	GetPatchFileDownloadUrl() string
}

const (
	PLATFORM_TYPE_ANDROID = "A"
	PLATFORM_TYPE_IOS     = "I"
)

type UrlConfigLocal struct{}

func (UrlConfigLocal) GetUpdateUrlGroup() map[string]string {
	return map[string]string{
		PLATFORM_TYPE_ANDROID: "https://rink.hockeyapp.net/apps/c99b93f9ccea47e997042c4b6f52adc2",
		PLATFORM_TYPE_IOS:     "https://rink.hockeyapp.net/apps/c99b93f9ccea47e997042c4b6f52adc2",
	}
}

func (UrlConfigLocal) GetPatchFileDownloadUrl() string {
	return "http://test-danceville-hub.com2us.net/cdn/patch/dev"
}

type UrlConfigDev struct{}

func (UrlConfigDev) GetUpdateUrlGroup() map[string]string {
	return map[string]string{
		PLATFORM_TYPE_ANDROID: "https://rink.hockeyapp.net/apps/c99b93f9ccea47e997042c4b6f52adc2",
		PLATFORM_TYPE_IOS:     "https://rink.hockeyapp.net/apps/c99b93f9ccea47e997042c4b6f52adc2",
	}
}

func (UrlConfigDev) GetPatchFileDownloadUrl() string {
	return "http://test-danceville-hub.com2us.net/cdn/patch/dev"
}

type UrlConfigQa struct{}

func (UrlConfigQa) GetUpdateUrlGroup() map[string]string {
	return map[string]string{
		PLATFORM_TYPE_ANDROID: "https://rink.hockeyapp.net/apps/c99b93f9ccea47e997042c4b6f52adc2",
		PLATFORM_TYPE_IOS:     "https://rink.hockeyapp.net/apps/c99b93f9ccea47e997042c4b6f52adc2",
	}
}

func (UrlConfigQa) GetPatchFileDownloadUrl() string {
	return "http://test-danceville-hub.com2us.net/cdn/patch/dev"
}

type UrlConfigLive struct{}

func (UrlConfigLive) GetUpdateUrlGroup() map[string]string {
	return map[string]string{
		PLATFORM_TYPE_ANDROID: "https://rink.hockeyapp.net/apps/c99b93f9ccea47e997042c4b6f52adc2",
		PLATFORM_TYPE_IOS:     "https://rink.hockeyapp.net/apps/c99b93f9ccea47e997042c4b6f52adc2",
	}
}

func (UrlConfigLive) GetPatchFileDownloadUrl() string {
	return "http://test-danceville-hub.com2us.net/cdn/patch/dev"
}
