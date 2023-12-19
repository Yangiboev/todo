package main

import (
	"log"
	"os"

	"github.com/Yangiboev/todo/config"
	"github.com/Yangiboev/todo/internal/server"
	"github.com/Yangiboev/todo/pkg/db/postgres"
	"github.com/Yangiboev/todo/pkg/logger"
	"github.com/Yangiboev/todo/pkg/utils"
)

// @title Go ToDo app
// @version 1.0
// @description Golang TODO app
// @contact.name Dilmurod Yangiboev
// @contact.url https://github.com/Yangiboev
// @contact.email dilmurod.yangiboev@gmail.com
// @BasePath /api/v1
func main() {
	log.Println("Starting api server")

	configPath := utils.GetConfigPath(os.Getenv("config"))

	cfgFile, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("LoadConfig: %v", err)
	}

	cfg, err := config.ParseConfig(cfgFile)
	if err != nil {
		log.Fatalf("ParseConfig: %v", err)
	}

	appLogger := logger.NewApiLogger(cfg)

	appLogger.InitLogger()
	appLogger.Infof("AppVersion: %s, LogLevel: %s, Mode: %s", cfg.Server.AppVersion, cfg.Logger.Level, cfg.Server.Mode)

	psqlDB, err := postgres.NewPsqlDB(cfg)
	if err != nil {
		appLogger.Fatalf("Postgresql init: %s", err)
	} else {
		appLogger.Infof("Postgres connected, Status: %#v", psqlDB.Stats())
	}
	defer psqlDB.Close()

	s := server.NewServer(cfg, psqlDB, appLogger)
	if err = s.Run(); err != nil {
		log.Fatal(err)
	}
}
