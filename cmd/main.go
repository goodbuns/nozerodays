package main

import (
	"fmt"
	"os"

	"github.com/goodbuns/nozerodays/bot"
)

func main() {
	token := os.Getenv("GITHUB_ACCESS_TOKEN")
	username := os.Getenv("GITHUB_USERNAME")
	config := bot.New(username, token)
	r := config.Repositories()
	fmt.Println(r)
	// bot.LatestContributionDate(username, token)
	// bot.ScrapeContributions(username)
}
