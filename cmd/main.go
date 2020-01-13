package main

import (
	"os"
)

func main() {
	// token := os.Getenv("GITHUB_ACCESS_TOKEN")
	username := os.Getenv("GITHUB_USERNAME")
	// bot.LatestContributionDate(username, token)
	bot.ScrapeContributions(username)
}
