package main

import (
	"context"
	"log"
	"time"
	"yc-backend/common"
	"yc-backend/engine"
	"yc-backend/internals"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	config, err := common.LoadConfiguration(common.ConfEnvSetting{YamlFilePath: []string{"./dev.yml", "./dev.example.yml"}})
	if err != nil {
		log.Fatal(err)
	}
	logger := internals.GetLogger()
	serverCtx, serverStopCtx := context.WithCancel(context.Background())
	defer serverStopCtx()
	if config.LogLevel == "debug" {
		logger.Infof("[Gin-Debug] SET gin.forceConsoleLog")
		gin.ForceConsoleColor()
	}

	clientOpts := options.Client().
		ApplyURI(config.MongoDB.DBUri).
		SetMaxPoolSize(20).
		SetMinPoolSize(5)
	client := setupDatabase(clientOpts)
	app := &engine.Application{
		DB:      client,
		Config:  config,
		Logger:  logger,
		Context: serverCtx,
	}

	app.Setup().
		RegisterRoute().
		GracefulShutdown()

	if err := app.ListenAndServe(); err != nil {
		log.Panic(err)
	}
}

func setupDatabase(clientOpts *options.ClientOptions) *mongo.Client {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancelFunc()

	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancelFunc = context.WithTimeout(context.Background(), 1*time.Second)
	defer cancelFunc()
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	return client
}
