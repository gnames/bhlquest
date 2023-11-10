package config

type Config struct {
	// BHLDir is the path to BHL items and OCRed texts.
	BHLDir string

	// LlmUtilURL is the URL to the llmutil service.
	LlmUtilURL string

	// DbHost is the database host.
	DbHost string

	// DbUser is the username used to access the database.
	DbUser string

	// DbPass is the DbUser's password.
	DbPass string

	// DbBHLQuest is the database name where BHLquest keeps its data.
	DbBHLQuest string

	// DbBHLNames is the database where BHLnames keeps its data.
	DbBHLNames string
}

type Option func(*Config)

func OptBHLDir(s string) Option {
	return func(cfg *Config) {
		cfg.BHLDir = s
	}
}

func OptLlmUtilURL(s string) Option {
	return func(cfg *Config) {
		cfg.LlmUtilURL = s
	}
}

func OptDbHost(s string) Option {
	return func(cfg *Config) {
		cfg.DbHost = s
	}
}

func OptDbUser(s string) Option {
	return func(cfg *Config) {
		cfg.DbUser = s
	}
}

func OptDbPass(s string) Option {
	return func(cfg *Config) {
		cfg.DbPass = s
	}
}

func OptDbBHLQuest(s string) Option {
	return func(cfg *Config) {
		cfg.DbBHLQuest = s
	}
}

func OptDbBHLNames(s string) Option {
	return func(cfg *Config) {
		cfg.DbBHLNames = s
	}
}

func New(opts ...Option) Config {
	res := Config{}
	for _, opt := range opts {
		opt(&res)
	}
	return res
}
