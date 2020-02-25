package main

import (
	"fmt"
	"log"
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
	vkNotifier := notify.NewVKNotifier(cfg.VKAccessToken, logrus.StandardLogger())

	// tg notifier
	tgNotifier, err := notify.NewTGNotifier(
		cfg.TGAccessToken,
		logrus.StandardLogger(),
	)
	if err != nil {
		logrus.Fatal(err)
	}

	// options
	opts := &app.Options{
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
	}

	// joke vk service
	vkService := app.NewJokeService(
		vkNotifier,
		jokeRepo,
		logrus.StandardLogger(),
		opts,
	)

	// joke tg service
	tgService := app.NewJokeService(
		tgNotifier,
		jokeRepo,
		logrus.StandardLogger(),
		opts,
	)

	// tg listener
	tgListener, err := callback.NewTGListener(
		cfg.TGAccessToken,
		cfg.TGWebhook,
		tgService,
		logrus.StandardLogger(),
	)
	if err != nil {
		log.Fatal(err)
	}
	err = tgListener.Listen("/tg")
	if err != nil {
		log.Fatal(err)
	}

	// vk callback
	vkListener := callback.NewVKListener(
		cfg.VKConfirmToken,
		cfg.VKSecretKey,
		vkService,
		logrus.StandardLogger(),
	)
	vkListener.Listen("/")

	logrus.Info("application server start...")

	addr := fmt.Sprintf(`%s:%d`, cfg.Host, cfg.Port)
	logrus.Fatal(http.ListenAndServe(addr, nil))
}

func loadConfig() (config, error) {
	var cfg config
	err := envconfig.Process("", &cfg)
	return cfg, err
}
