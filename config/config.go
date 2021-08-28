package config

import (
	"flag"
	"net"
	"net/url"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

const SchemaVersion = 1

// Config consists of service params
type Config struct {
	Debug bool `default:"false"`

	Http HttpConfig
	TLS  TLSConfig

	Db DBConfig

	ApiVersion string `default:"v1"`
}

type HttpConfig struct {
	Host          string        `default:"0.0.0.0"`
	Port          int           `default:"8080"`
	ReadTimeout   time.Duration `default:"15s"`
	WriteTimeout  time.Duration `default:"15s"`
	SessionSecret string        `default:"secret"`
	SessionName   string        `default:"grpc"`
}

type DBConfig struct {
	SchemaVersion  uint   `default:"1"`
	MigrationPath  string `default:"migrations"`
	Debug          bool   `default:"false"`
	Dialect        string `default:"pgx"`
	Host           string `default:"localhost"`
	Port           string `default:"5432"`
	Name           string `default:"asscRegistration"`
	User           string `default:"assc"`
	Password       string `default:"blablabla"`
	ConnectionPool int    `default:"10"`
	ConnectionMax  int    `default:"100"`
}

func (cfg DBConfig) String() string {
	q := url.Values{}
	q.Add("sslmode", "disable")
	q.Add("connect_timeout", "5")
	q.Add("statement_cache_mode", "describe") // pgBouncer

	u := url.URL{
		Scheme:   "postgresql",
		User:     url.UserPassword(cfg.User, cfg.Password),
		Host:     net.JoinHostPort(cfg.Host, cfg.Port),
		Path:     cfg.Name,
		RawQuery: q.Encode(),
	}

	return u.String()
}

type TLSConfig struct {
	Enabled  bool
	CertPath string
	KeyPath  string
}

func Read() *Config {
	configPath := flag.String("config", "deploy/example.env", "config path")
	flag.Parse()

	_ = godotenv.Load(*configPath)

	var cfg Config

	envconfig.MustProcess("", &cfg)

	cfg.Db.SchemaVersion = SchemaVersion

	return &cfg
}
