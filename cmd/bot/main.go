package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/delgus/jokebot/internal/app"
	"github.com/delgus/jokebot/internal/bots/jokebot"
	"github.com/delgus/jokebot/internal/bots/jokebot/store/sql"
	"github.com/delgus/jokebot/internal/tg"
	"github.com/delgus/jokebot/internal/vk"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

type config struct {
	Host           string `envconfig:"HOST"`
	Port           int    `envconfig:"PORT"`
	VKAccessToken  string `envconfig:"VK_ACCESS_TOKEN" default:"default"`
	VKConfirmToken string `envconfig:"VK_CONFIRM_TOKEN" default:"default"`
	VKSecretKey    string `envconfig:"VK_SECRET_KEY" default:"default"`
	DBDriver       string `envconfig:"DB_DRIVER" default:"sqlite3"`
	DBAddr         string `envconfig:"DB_ADDR" default:"../../database/jokebot.db"`
	DBDebug        bool   `envconfig:"DB_DEBUG" default:"false"`
	TGAccessToken  string `envconfig:"TG_ACCESS_TOKEN"`
	TGWebhook      string `envconfig:"TG_WEBHOOK"`
}

func main() {
	// configurations
	cfg, err := loadConfig()
	if err != nil {
		logrus.Fatalf("can't load config: %v", err)
	}

	// db connection for jokebot
	db, err := sql.NewConnection(sql.ConnectionOptions{
		Addr:   cfg.DBAddr,
		Driver: cfg.DBDriver,
		Debug:  cfg.DBDebug,
		Logger: logrus.StandardLogger(),
	})
	if err != nil {
		logrus.Fatalf("can't connect database: %v", err)
	}

	// repository for jokebot
	jokeRepo := sql.NewJokeRepo(db)

	jokeBot := jokebot.NewBot(jokeRepo)

	// options for notifiers
	opts := app.Options{
		EOFText:           `Шутки кончились господа!`,
		InternalErrorText: `Внутренняя ошибка. На нашей стороне ошибка, попробуйте позднее`,
		HelpText: `
	Команды для бота:
	
	list - список категорий анекдотов
	
	Чтобы получить анекдот из категории,отправьте номер категории
	
	joke - возвращает анекдот из любой категории
	
	help - помощь
	`,

		NotCorrectCommandText: "Неверная команда! \n",
	}

	// vk bot
	vkBotApp := &app.App{
		Notifier: vk.NewNotifier(cfg.VKAccessToken),
		Bot:      jokeBot,
		Listener: vk.NewListener(cfg.VKConfirmToken, cfg.VKSecretKey),
		Logger:   logrus.StandardLogger(),
		Options:  opts,
	}

	// tg notifier
	tgNotifier, err := tg.NewNotifier(cfg.TGAccessToken)
	if err != nil {
		logrus.Fatal(err)
	}

	// tg listener
	tgListener, err := tg.NewListener(cfg.TGAccessToken, cfg.TGWebhook)
	if err != nil {
		log.Fatal(err)
	}

	// tg bot
	tgBotApp := &app.App{
		Notifier: tgNotifier,
		Bot:      jokeBot,
		Listener: tgListener,
		Logger:   logrus.StandardLogger(),
		Options:  opts,
	}

	go vkBotApp.Run("/vk")

	go tgBotApp.Run("/tg")

	logrus.Info("application server start...")

	addr := fmt.Sprintf(`%s:%d`, cfg.Host, cfg.Port)
	logrus.Fatal(http.ListenAndServe(addr, nil))
}

func loadConfig() (config, error) {
	var cfg config
	err := envconfig.Process("", &cfg)
	return cfg, err
}
