package cmd

import (
	"Templatest/internal/api"
	"Templatest/internal/config"
	customLogger "Templatest/internal/logger"
	"Templatest/internal/repo"
	"Templatest/internal/service"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	// Config
	if err := godotenv.Load(config.EnvPath); err != nil {
		log.Fatalf("failed to load config file %s: %v", config.EnvPath, err)
	}

	var cfg config.AppConfig
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatalf("failed to process config file %s: %v", config.EnvPath, err)
	}

	// Logger
	logger, err := customLogger.NewLogger(cfg.LogLevel)
	if err != nil {
		log.Fatal("failed to initialize logger: ", err)
	}

	// Repository
	repository, err := repo.NewRepository(cfg.Memory)
	if err != nil {
		log.Fatal("failed to initialize repo: ", err)
	}

	// Service initialization
	serviceInstance := service.NewService(repository, logger)

	// Routers initialization
	app := api.NewRouters(&api.Routers{Service: serviceInstance}, "token")

	// Listening and serving
	go func() {
		logger.Infof("Starting server on %s", cfg.Rest.Port)
		if err := app.Listen(cfg.Rest.Port); err != nil {
			log.Fatalf("failed to start server: %v", err)
		}
	}()

	// Fold operations
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	<-signalChan

	logger.Info("Shutting down gracefully...")
}
