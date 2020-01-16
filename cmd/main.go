package main

import (
	"os"

	"github.com/goodbuns/nozerodays/bot"
)

func main() {
	// todo: parse in flags instead
	token := os.Getenv("GITHUB_ACCESS_TOKEN")
	username := os.Getenv("GITHUB_USERNAME")
	webhookURL := os.Getenv("WEBHOOK_URL")
	config := bot.New(username, token, webhookURL)
	config.Start()
	// r := config.Repositories()
	// c := config.Commits(r)
	// fmt.Println("found commit for today", *c)
}
