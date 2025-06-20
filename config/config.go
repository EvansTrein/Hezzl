package config

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

const (
	ModeLocal      = "local"
	ModeDev        = "dev"
	ModeProd       = "prod"
	defaultLogPath = "./logs/out.log"
)

var conf config

type config struct {
	Env         string `env:"ENV" env-required:"true"`
	StoragePath string `env:"POSTGRES_PATH" env-required:"true"`
	LogOutput   string `env:"LOG_OUTPUT"`
	LogLevel    string `env:"LOG_LEVEL"`
	Redis       `env-prefix:"REDIS_"`
	HttpServer  `env-prefix:"HTTP_"`
	Clickhouse  `env-prefix:"CLICKHOUSE_"`
	Nats        `env-prefix:"NATS_"`
}

type HttpServer struct {
	Host string `env:"API_HOST" env-required:"true"`
	Port string `env:"API_PORT" env-required:"true"`
}

type Redis struct {
	Host     string        `env:"HOST" env-required:"true"`
	Port     string        `env:"PORT" env-required:"true"`
	Password string        `env:"PASSWORD" env-required:"true"`
	TTLKeys  time.Duration `env:"TTL" env-required:"true"`
	NumberDB int           `env:"DB_NUMBER"` // default == 0
}

type Clickhouse struct {
	Host     string `env:"HOST" env-required:"true"`
	Port     string `env:"TCP_PORT" env-required:"true"`
	DB       string `env:"DB" env-required:"true"`
	User     string `env:"USER" env-required:"true"`
	Password string `env:"PASSWORD" env-required:"true"`
}

type Nats struct {
	Host     string `env:"HOST" env-required:"true"`
	Port     string `env:"PORT" env-required:"true"`
	NameMess string `env:"NAME_MESSAGES" env-required:"true"`
}

func MustLoad() {
	var filePath string

	flag.StringVar(&filePath, "config", "", "path to config file")
	flag.Parse()

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		log.Fatalf("env file does not exist: %s", filePath)
	}

	if err := cleanenv.ReadConfig(filePath, &conf); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	log.Println("configuration file successfully loaded")
}

func GetConfig() *config {
	return &conf
}
