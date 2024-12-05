package data

var (
	GConfig *Config
)

func InitGlobal() {
	GConfig = &Config{}
	GConfig.Load()
}
