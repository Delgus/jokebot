package main

import (
	"net/http"

	api "github.com/delgus/go-vk/callback-api"
	"github.com/delgus/go-vk/client"
	"github.com/delgus/jokebot/internal/inrastructure/store/sql"
	"github.com/delgus/jokebot/internal/vkhooks"
	"github.com/sirupsen/logrus"
)

func main() {
	cfg, err := loadConfig()
	if err != nil {
		logrus.Fatalf("vk-server: can't load config: %v", err)
	}

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
