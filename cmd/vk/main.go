package main

import (
	"fmt"
	"net/http"

	"github.com/SevereCloud/vksdk/api"
	"github.com/SevereCloud/vksdk/callback"
	"github.com/delgus/jokebot/internal/inrastructure/store/sql"
	"github.com/delgus/jokebot/internal/vkhooks"
	"github.com/kelseyhightower/envconfig"
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

	// vk client
	vk := api.Init(cfg.VKAccessToken)

	// vk callback
	var cb callback.Callback
	cb.ConfirmationKey = cfg.VKConfirmToken
	cb.SecretKey = cfg.VKSecretKey

	cb.MessageNew(vkhooks.OnMessageNew(jokeRepo, vk))

	http.HandleFunc("/", cb.HandleFunc)

	logrus.Println("vk server start...")
	addr := fmt.Sprintf(`%s:%d`, cfg.Host, cfg.Port)
	logrus.Fatal(http.ListenAndServe(addr, nil))
}

func loadConfig() (config, error) {
	var cfg config
	err := envconfig.Process("", &cfg)
	return cfg, err
}
