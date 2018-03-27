package config

import (
	"github.com/BurntSushi/toml"
)

type Config struct {
	Database Database `toml:"database"`
	Routes   Routes   `toml:"routes"`
}

func New(f string) (*Config, error) {
	c := Config{}
	_, err := toml.DecodeFile(f, &c)

	return &c, err
}

type Database struct {
	Postgres struct {
		DSN string `toml:"dsn"`
	} `toml:"postgres"`

	Redis struct {
		Addrs []string `toml:"addrs"`
	} `toml:"redis"`
}

type Routes struct {
	TemplatesDir    string `toml:"templates_dir"`
	RecaptchaSecret string `toml:"recaptcha_secret"`
	Certs           string `toml:"certs"`
}
