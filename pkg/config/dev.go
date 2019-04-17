package config

type DevConfig struct{}

func (DevConfig) GetUpdateUrlGroup() map[string]string {
	return map[string]string{
		PLATFORM_TYPE_ANDROID: "https://rink.hockeyapp.net/apps/c99b93f9ccea47e997042c4b6f52adc2",
		PLATFORM_TYPE_IOS:     "https://rink.hockeyapp.net/apps/c99b93f9ccea47e997042c4b6f52adc2",
	}
}

func (DevConfig) GetPatchFileDownloadUrl() string {
	return "http://test-danceville-hub.com2us.net/cdn/patch/dev"
}
