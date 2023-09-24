package config

import "github.com/caarlos0/env/v8"

type Config struct {
	Bot     BotConfig `envPrefix:"BOT_"`
	Channel string    `env:"CHANNEL"`
}

type BotConfig struct {
	Token   string `env:"TOKEN,notEmpty"`
	AppHash string `env:"APP_HASH,notEmpty"`
	AppId   int32  `env:"APP_ID,notEmpty"`
}

func NewConfig() (*Config, error) {
	var c Config
	if err := env.Parse(&c); err != nil {
		return nil, err
	}
	return &c, nil
}
