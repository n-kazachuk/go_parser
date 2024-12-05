package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env    string       `yaml:"env" env-default:"local"`
	Pgsql  PgsqlConfig  `yaml:"pgsql"`
	Kafka  KafkaConfig  `yaml:"kafka"`
	Parser ParserConfig `yaml:"parser"`
}

type PgsqlConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DbName   string `yaml:"dbname"`
}

type KafkaConfig struct {
	Interval int      `yaml:"interval"`
	Brokers  []string `yaml:"brokers"`
}

type ParserConfig struct {
	Timeout  time.Duration      `yaml:"timeout"`
	Interval time.Duration      `yaml:"interval"`
	Worker   ParserWorkerConfig `yaml:"worker"`
}

type ParserWorkerConfig struct {
	Count    int           `yaml:"count"`
	Interval time.Duration `yaml:"interval"`
}

func MustLoad() *Config {
	configPath := fetchConfigPath()
	if configPath == "" {
		panic("config path is empty")
	}

	return MustLoadPath(configPath)
}

func MustLoadPath(configPath string) *Config {
	// check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file does not exist: " + configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("cannot read config: " + err.Error())
	}

	return &cfg
}

// fetchConfigPath fetches config path from command line flag or environment variable.
// Priority: flag > env > default.
// Default value is empty string.
func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
