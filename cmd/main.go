package cmd

import (
	"TempletefullDDDCRUD/internal/config"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"log"
)

func main() {
	if err := godotenv.Load(config.EnvPath); err != nil {
		log.Fatalf("failed to load config file %s: %v", config.EnvPath, err)
	}

	var cfg config.AppConfig
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatalf("failed to process config file %s: %v", config.EnvPath, err)
	}

}
