package main

type config struct {
	Host           string `envconfig:"HOST"`
	Port           int    `envconfig:"PORT"`
	VKAccessToken  string `envconfig:"VK_ACCESS_TOKEN" default:"default"`
	VKConfirmToken string `envconfig:"VK_CONFIRM_TOKEN" default:"default"`
	VKSecretKey    string `envconfig:"VK_SECRET_KEY" default:"default"`
	DBDriver       string `envconfig:"DB_DRIVER" default:"sqlite3"`
	DBAddr         string `envconfig:"DB_ADDR" default:"../../database/jokebot.db"`
	DBDebug        bool   `envconfig:"DB_DEBUG" default:"false"`
}
