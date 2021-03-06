package main

import (
	"fmt"
	"net/http"

	easybot "github.com/delgus/easy-bot"
	"github.com/delgus/jokebot/internal/jokebot"
	"github.com/delgus/jokebot/internal/jokebot/store/sql"
	tghook "github.com/delgus/tg-logrus-hook"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

type config struct {
	Host             string `envconfig:"HOST"`
	Port             int    `envconfig:"PORT"`
	VKAccessToken    string `envconfig:"VK_ACCESS_TOKEN" default:"default"`
	VKConfirmToken   string `envconfig:"VK_CONFIRM_TOKEN" default:"default"`
	VKSecretKey      string `envconfig:"VK_SECRET_KEY" default:"default"`
	DBDriver         string `envconfig:"DB_DRIVER" default:"sqlite3"`
	DBAddr           string `envconfig:"DB_ADDR" default:"../../database/jokebot.db"`
	DBDebug          bool   `envconfig:"DB_DEBUG" default:"false"`
	TGAccessToken    string `envconfig:"TG_ACCESS_TOKEN"`
	TGWebhook        string `envconfig:"TG_WEBHOOK"`
	LogTGChatID      int64  `envconfig:"LOG_TG_CHAT_ID"`
	LogTGAccessToken string `envconfig:"LOG_TG_ACCESS_TOKEN"`
}

func main() {
	// configurations
	cfg, err := loadConfig()
	if err != nil {
		logrus.Fatalf("can't load config: %v", err)
	}

	hook, err := tghook.NewHook(cfg.LogTGAccessToken, cfg.LogTGChatID, logrus.AllLevels)
	if err != nil {
		logrus.Errorf(`can not create tg hook for logging: %v`, err)
	} else {
		logrus.Info(`create tg hook for logging`)
		logrus.StandardLogger().AddHook(hook)
		logrus.Info(`add hook!!!`)
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

	// jokebot
	jokeBot := jokebot.NewBot(sql.NewJokeRepo(db))

	// common options for apps
	opts := easybot.Options{
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

	// vk app with bot
	vkBotApp := vkApp(cfg, jokeBot, opts, logrus.StandardLogger())
	go func() {
		if err := vkBotApp.Run("/vk"); err != nil {
			logrus.Fatal("can not start vk bot app")
		}
	}()

	// tg app with bot
	tgBotApp, err := tgApp(cfg, jokeBot, opts, logrus.StandardLogger())
	if err != nil {
		logrus.Fatal(err)
	}
	go func() {
		if err := tgBotApp.Run("/tg"); err != nil {
			logrus.Fatal("can not start tg bot app")
		}
	}()

	logrus.Info("start server for tg and app")
	addr := fmt.Sprintf(`%s:%d`, cfg.Host, cfg.Port)
	logrus.Fatal(http.ListenAndServe(addr, nil))
}

func loadConfig() (*config, error) {
	var cfg config
	err := envconfig.Process("", &cfg)
	return &cfg, err
}
