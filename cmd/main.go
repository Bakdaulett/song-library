package main

import (
	"os"
	"os/signal"
	"song-library/internal/app"
	"song-library/internal/handlers"
	"song-library/internal/repository"
	"song-library/internal/service"
	"song-library/pkg/database"
	"song-library/pkg/logger"
	"song-library/pkg/migrations"
	"syscall"

	"context"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error loading .env file: %s", err.Error())
	}

	// Initialize logger based on config.yml
	logger.InitLogger()

	// Initialize configurations from config.yml
	if err := initConfig(); err != nil {
		logger.Error("Error initializing configs: " + err.Error())
		os.Exit(1)
	}

	// Get database configuration from viper (loaded from config.yml)
	db, err := database.NewPostgresDB(database.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: os.Getenv("DB_PASSWORD"), // Sensitive data from .env file
	})
	if err != nil {
		logger.Error("Failed to initialize database: " + err.Error())
		os.Exit(1)
	}
	defer db.Close()

	// Apply migrations
	if err := migrations.Migrate(db); err != nil {
		logger.Error("Error running migrations: " + err.Error())
		os.Exit(1)
	}

	repo := repository.NewSongRepository(db)
	services := service.NewSongService(repo)
	handlers := handlers.NewHandler(services)

	// Start HTTP server
	srv := new(app.Server)
	go func() {
		if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
			logger.Error("Error occurred while running HTTP server: " + err.Error())
		}
	}()

	logger.Info("Server is running")

	// Graceful shutdown handling
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	// Log shutdown process
	logger.Info("Server shutting down")
	if err := srv.Shutdown(context.Background()); err != nil {
		logger.Error("Error occurred during server shutdown: " + err.Error())
	}

	if err := db.Close(); err != nil {
		logger.Error("Error occurred while closing database connection: " + err.Error())
	}
}

func initConfig() error {
	// Load configuration from config.yml
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	return nil
}
