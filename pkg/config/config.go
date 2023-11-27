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

	// Port is the port where to run the RESTful service.
	Port int

	// APIDocURL is the url to the API documentation.
	APIDocURL string

	// InitClasses limits embedded import to certain taxa.
	InitClasses []string

	// InitTaxa limitd embedded import to main taxons of items.
	InitTaxa []string

	// ScoreThreshold filters out results with too low score.
	ScoreThreshold float64

	// MaxResultsNum limits the maximum number of results returned
	// in an answer.
	MaxResultsNum int

	// WithoutConfirm when true, remves confirmation dialogs.
	WithoutConfirm bool

	// WithRebuildDb flag is true if the bhlquest database needs to
	// be rebuilt from scratch.
	WithRebuildDb bool

	// WithText is a flag that control appearance of matched texts with
	// the results. If the flag is false, the text does not appear.
	WithText bool
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

func OptPort(i int) Option {
	return func(cfg *Config) {
		cfg.Port = i
	}
}

func OptInitClasses(ss []string) Option {
	return func(cfg *Config) {
		cfg.InitClasses = ss
	}
}

func OptInitTaxa(ss []string) Option {
	return func(cfg *Config) {
		cfg.InitTaxa = ss
	}
}

func OptScoreThreshold(f float64) Option {
	return func(cfg *Config) {
		cfg.ScoreThreshold = f
	}
}

func OptMaxResultsNum(i int) Option {
	return func(cfg *Config) {
		cfg.MaxResultsNum = i
	}
}

func OptWithoutConfirm(b bool) Option {
	return func(cfg *Config) {
		cfg.WithoutConfirm = b
	}
}

func OptWithRebuildDb(b bool) Option {
	return func(cfg *Config) {
		cfg.WithRebuildDb = b
	}
}

func OptWithText(b bool) Option {
	return func(cfg *Config) {
		cfg.WithText = b
	}
}

func New(opts ...Option) Config {
	res := Config{
		BHLDir:         "/opt/bhl/",
		LlmUtilURL:     "http://0.0.0.0:8000/api/v1/",
		DbHost:         "0.0.0.0",
		DbUser:         "postgres",
		DbPass:         "postgres",
		DbBHLQuest:     "bhlquest",
		DbBHLNames:     "bhlnames",
		Port:           8555,
		APIDocURL:      "https://apidoc.globalnames.org/bhlquest",
		ScoreThreshold: 0.4,
		MaxResultsNum:  5,
	}
	for _, opt := range opts {
		opt(&res)
	}
	return res
}
