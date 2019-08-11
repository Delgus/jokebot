package main

import (
	"net/http"
	"os"

	api "github.com/delgus/go-vk/callback-api"
	"github.com/delgus/go-vk/client"
	"github.com/delgus/jokebot/internal/inrastructure/store/sql"
	"github.com/delgus/jokebot/internal/vkhooks"
	"github.com/sirupsen/logrus"
	"github.com/subosito/gotenv"
)

func main() {
	//load config
	if err := gotenv.Load(".env"); err != nil {
		logrus.Fatalf("vk-server: can't load environments: %v", err)
	}
	cfg, err := loadConfig()
	if err != nil {
		logrus.Fatalf("vk-server: can't load config: %v", err)
	}
	//file for logs
	file, err := os.OpenFile(cfg.VKLogFile, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		logrus.Fatalf("vk-server: can't open log file: %v", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			logrus.Fatalf("vk-server: can't close log file: %v", err)
		}
	}()
	logrus.SetOutput(file)

	//db connection
	db, err := sql.NewConnection(sql.ConnectionOptions{
		Addr:   cfg.DBAddr,
		Driver: cfg.DBDriver,
		Debug:  cfg.DBDebug})
	if err != nil {
		logrus.Fatalf("vk-server: can't connect database: %v", err)
	}

	//service
	jokeRepo := sql.NewJokeRepo(db)

	//client
	vk := new(client.VKClient)
	vk.SetAccessToken(cfg.VKAccessToken)

	//callback api server
	apiServer := api.NewApiServer()
	apiServer.SetConfirmToken(cfg.VKConfirmToken)
	apiServer.SetSecret(cfg.VKSecretKey)

	apiServer.OnMessageNew = vkhooks.OnMessageNew(jokeRepo, vk)

	logrus.Println("vk server start...")
	if err := http.ListenAndServe(":9000", apiServer); err != nil {
		panic(err)
	}

}
