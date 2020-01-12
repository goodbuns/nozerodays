package main

import (
	"os"

	"github.com/goodbuns/nozerodays/bot"
)

func main() {
	token := os.Getenv("GITHUB_ACCESS_TOKEN")
	username := os.Getenv("GITHUB_USERNAME")
	bot.LatestContributionDate(username, token)
}
