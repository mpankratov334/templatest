package config

import "time"

const EnvPath = "local.env"

type AppConfig struct {
	LogLevel string
	Rest     Rest
	Memory   InMemory
}

type Rest struct {
	Port           string        `envconfig:"PORT" default:"8080"`
	RequestTimeout time.Duration `envconfig:"REQUEST_TIMEOUT" default:"30"`
}

type InMemory struct {
	Capacity    int `envconfig:"MAX_ITEMS" default:"10000"`
	MaxItemSize int `envconfig:"MAX_ITEM_SIZE" default:"1024"`
}
