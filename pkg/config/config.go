package config

type Config struct {
	// OpenaiAPIKey a key to OpenAI API
	OpenaiAPIKey string

	// CohereAPIKey is a key to Cohere service
	CohereAPIKey string

	// BHLDir is the path to BHL items and OCRed texts.
	BHLDir string

	// LlmUtilURL is the URL to the llmutil service.
	LlmUtilURL string

	// QdrantHost is the host for creation qdrant grpc connection.
	QdrantHost string

	// QdrantSegmentsNum sets the number of segmens for Qdrant.
	QdrantSegmentsNum uint64

	// VectorSize sets the number of dimentions for embeding
	// vectors.
	VectorSize uint64

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

	// WithSummary adds LLM-generated summary generated from
	// received BHL data.
	WithSummary bool

	// WithCrossEmbed flag controls use of Cross-Embed comparison
	// of a question with results.
	WithCrossEmbed bool
}

type Option func(*Config)

func OptOpenaiAPIKey(s string) Option {
	return func(cfg *Config) {
		cfg.OpenaiAPIKey = s
	}
}

func OptCohereAPIKey(s string) Option {
	return func(cfg *Config) {
		cfg.CohereAPIKey = s
	}
}

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

func OptQdrantHost(s string) Option {
	return func(cfg *Config) {
		cfg.QdrantHost = s
	}
}

func OptVectorSize(i uint64) Option {
	return func(cfg *Config) {
		cfg.VectorSize = i
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

func New(opts ...Option) Config {
	res := Config{
		BHLDir:            "/opt/bhl/",
		LlmUtilURL:        "http://0.0.0.0:8000/api/v1/",
		QdrantHost:        "0.0.0.0:6334",
		QdrantSegmentsNum: 2,
		VectorSize:        768,
		DbHost:            "0.0.0.0",
		DbUser:            "postgres",
		DbPass:            "postgres",
		DbBHLQuest:        "bhlquest",
		DbBHLNames:        "bhlnames",
		Port:              8555,
		APIDocURL:         "https://apidoc.globalnames.org/bhlquest",
		ScoreThreshold:    0.4,
		MaxResultsNum:     5,
		WithSummary:       true,
		WithCrossEmbed:    false,
	}
	for _, opt := range opts {
		opt(&res)
	}
	return res
}
