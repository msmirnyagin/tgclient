package main

import (
	"log"
	"tgclient/config"
	"tgclient/internal/botcode"
	"tgclient/internal/webhook"
)

func main() {

	cfg := config.ReadEnv()

	code, err := botcode.GetCode(cfg.Bot)
	if err != nil {
		log.Panic(err)
	}
	err = webhook.GetPost(cfg.Url, code)
	if err != nil {
		log.Panic(err)
	}
}
