package main

import (
	"fmt"
	"net/http"

	"github.com/delgus/jokebot/internal/app"
	"github.com/delgus/jokebot/internal/inrastructure/callback"
	"github.com/delgus/jokebot/internal/inrastructure/notify"
	"github.com/delgus/jokebot/internal/inrastructure/store/sql"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

func main() {
	cfg, err := loadConfig()
	if err != nil {
		logrus.Fatalf("vk-server: can't load config: %v", err)
	}

	// db connection
	db, err := sql.NewConnection(sql.ConnectionOptions{
		Addr:   cfg.DBAddr,
		Driver: cfg.DBDriver,
		Debug:  cfg.DBDebug,
		Logger: logrus.StandardLogger(),
	})
	if err != nil {
		logrus.Fatalf("vk-server: can't connect database: %v", err)
	}

	// repository
	jokeRepo := sql.NewJokeRepo(db)

	// vk notifier
	notifier := notify.NewVKNotifier(cfg.VKAccessToken, logrus.StandardLogger())
	notifier.Keyboard(`{
		"buttons": [
		  [
			{
			  "action": {
				"type": "text",
				"label": "Анекдот",
				"payload": "{\"command\":\"joke\"}"
			  },
			  "color": "positive"
			}
		  ],
		  [
			{
			  "action": {
				"type": "text",
				"label": "Категории анекдотов",
				"payload": "{\"command\":\"list\"}"
			  },
			  "color": "negative"
			}
		  ],
		  [
			{
			  "action": {
				"type": "text",
				"label": "Помощь",
				"payload": "{\"command\":\"help\"}"
			  },
			  "color": "primary"
			}
		  ]
		]
	  }`)

	// joke service
	service := app.NewJokeService(
		notifier,
		jokeRepo,
		logrus.StandardLogger(),
		&app.Options{
			JokeCommand:            "joke",
			ListCommand:            "list",
			HelpCommand:            "help",
			JokesAreOverText:       "К сожалению шутки закончились",
			TryAnotherCategoryText: "Попробуйте другую категорию!",
			InternalErrorText:      "К сожалению произошла ошибка. Попробуйте получить шутку позднее",
			HelpMessageText: `
			Команды для бота:
			
			list - список категорий анекдотов
			
			Чтобы получить анекдот из категории,отправьте номер категории
			
			joke - возвращает анекдот из любой категории
			
			help - помощь
			`,
			NotCorrectCommandText: "Неверная команда! \n",
		})

	// vk callback
	cb := callback.NewVKCallback(
		cfg.VKConfirmToken,
		cfg.VKSecretKey,
		service,
		logrus.StandardLogger())
	http.HandleFunc("/", cb.HandleFunc)

	logrus.Info("application server start...")
	addr := fmt.Sprintf(`%s:%d`, cfg.Host, cfg.Port)
	logrus.Fatal(http.ListenAndServe(addr, nil))
}

func loadConfig() (config, error) {
	var cfg config
	err := envconfig.Process("", &cfg)
	return cfg, err
}
