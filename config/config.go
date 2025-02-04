package config

type Settings map[string]string

type Config struct {
	settings Settings
}

// Comment
func Init() *Config {
	return &Config{
		settings: make(Settings),
	}
}

// Comment
func (ctx *Config) Set(key string, value string) *Config {
	ctx.settings[key] = value

	return ctx
}

// Comment
func (ctx *Config) Get(key string) string {
	setting, ok := ctx.settings[key]

	if !ok {
		return ""
	}

	return setting
}
