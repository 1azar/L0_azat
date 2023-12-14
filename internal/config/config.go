package config

import (
	"flag"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

const configEnvName = "CONFIG_PATH"

type Config struct {
	Env           string               `yaml:"env" env-required:"true"`
	DbCredentials DbCredentials        `yaml:"dbCredentials" env-required:"true"`
	NatsCfg       NatsConfiguration    `yaml:"natsCfg" env-required:"true"`
	ServiceCfg    ServiceConfiguration `yaml:"serviceCfg" env-required:"true"`
	HttpCfg       HttpConfiguration    `yaml:"httpCfg" env-required:"true"`
}

type HttpConfiguration struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idleTimeout" env-default:"60s"`
}

type ServiceConfiguration struct {
	CacheSize int `yaml:"cacheSize" env-required:"true"`
}

type NatsConfiguration struct {
	ClusterIdEnv string `yaml:"clusterIdEnv" env-required:"true"`
	ClientIdENv  string `yaml:"clientIdENv" env-required:"true"`
	SubjectEnv   string `yaml:"subjectEnv" env-required:"true"`
}

// DbCredentials store env var names for db credentials.
// real values kept in env vars of a machine due to the critical nature of the data
type DbCredentials struct {
	UsernameEnv string `yaml:"usernameEnv" env-required:"true"`
	PasswordEnv string `yaml:"passwordEnv" env-required:"true"`
	AddressEnv  string `yaml:"addressEnv" env-required:"true"`
	PortEnv     string `yaml:"portEnv" env-required:"true"`
	DbNameEnv   string `yaml:"dbNameEnv" env-required:"true"`
}

func MustLoad() *Config {
	// fetches config path: flag > env > panic
	var path string

	flag.StringVar(&path, "config", "", "path to config file")
	flag.Parse()

	if path == "" { // no flag
		path = os.Getenv(configEnvName)
	}

	if path == "" { // nothing in env var
		panic(fmt.Sprintf("configuration file not specified.\n\tuse flag --config or set %s env variable", configEnvName))
	}

	if _, err := os.Stat(path); os.IsNotExist(err) { // no such file
		panic("config file does not exist: " + path)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("cannot read config file: " + err.Error())
	}

	return &cfg
}
