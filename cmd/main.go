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
	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error loading .env file: %s", err.Error())
	}

	logger.InitLogger()

	if err := initConfig(); err != nil {
		logger.Error("Error initializing configs: " + err.Error())
		os.Exit(1)
	}

	db, err := database.NewPostgresDB(database.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: os.Getenv("DB_PASSWORD"),
	})
	if err != nil {
		logger.Error("Failed to initialize database: " + err.Error())
		os.Exit(1)
	}
	defer db.Close()

	if err := migrations.Migrate(db); err != nil {
		logger.Error("Error running migrations: " + err.Error())
		os.Exit(1)
	}

	repo := repository.NewSongRepository(db)
	services := service.NewSongService(repo)
	handlers := handlers.NewHandler(services)

	srv := new(app.Server)
	go func() {
		if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
			logger.Error("Error occurred while running HTTP server: " + err.Error())
		}
	}()

	logger.Info("Server is running")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logger.Info("Server shutting down")
	if err := srv.Shutdown(context.Background()); err != nil {
		logger.Error("Error occurred during server shutdown: " + err.Error())
	}

	if err := db.Close(); err != nil {
		logger.Error("Error occurred while closing database connection: " + err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	return nil
}
