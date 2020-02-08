package bot

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/goodbuns/nozerodays/github"
	"github.com/rs/zerolog"
)

// New creates a new bot config.
func New() *Config {
	// Get configuration details from command line flags.
	username := flag.String("username", "", "github username")
	accessToken := flag.String("accessToken", "", "github access token")
	webhookURL := flag.String("webhook", "", "slack webhook URL")
	orgs := flag.String("organizations", "", "list of organizations to whitelist, with spaces between")
	l := flag.String("location", "", "location name corresponding to a file in the IANA Time Zone database")
	flag.Parse()

	host, _ := os.Hostname()
	zl := zerolog.New(os.Stdout).With().Timestamp().Str("host", host)
	logger := zl.Logger()

	logger.Info().Msg(fmt.Sprint("set configuration:", *username, *orgs, *l))

	location, err := time.LoadLocation(*l)
	if err != nil {
		logger.Err(err).Msg("")
		panic(err)
	}

	return &Config{
		webhookURL: *webhookURL,
		github:     github.New(*username, *accessToken, strings.Split(*orgs, " ")),
		location:   location,
		logger:     logger,
	}
}

// Config holds all configuration necessary for the bot to run.
type Config struct {
	github     *github.Client
	webhookURL string
	location   *time.Location
	logger     zerolog.Logger
}

// Start starts the bot.
func (c *Config) Start() {
	for true {
		currentTime := time.Now().In(c.location)
		// check whether it's past 8pm, if not, sleep until 8pm.
		today8PM := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 20, 0, 0, 0, c.location)
		timeUntil8 := time.Until(today8PM)
		if timeUntil8 > 0 {
			c.logger.Info().Msg(fmt.Sprintf("it is not yet 8pm. sleeping for %v until %v", timeUntil8, today8PM))
			time.Sleep(timeUntil8)
		}

		c.logger.Info().Msg("start scan")
		c.logger.Info().Msg(fmt.Sprintf("current time: %v", currentTime))
		repos, err := c.github.Repositories()
		if err != nil {
			c.logger.Err(err).Msg("failed to find repositories associated with the user/orgs requested")
		}
		r := fmt.Sprintf("%v", repos)
		c.logger.Info().Msg(r)

		commit, err := c.github.CommitCreatedToday(repos, c.location, currentTime)
		if err != nil {
			c.logger.Err(err).Msg("failed to find commit created today successfully")
		}

		if commit != nil {
			c.logger.Info().Msg(fmt.Sprintf("found commit for today (%s) with commit URL (%s)", commit.Commit.Author.Date.String(), commit.URL))
			// Send slack message.
			msg := slackMsg{
				Text: "good work! commit was made today at " + commit.Commit.Author.Date.String() + ", commit link: " + commit.URL,
			}
			err = c.sendSlackMsg(msg)
			if err != nil {
				c.logger.Err(err).Msg(fmt.Sprintf("unable to send the following message to slack: %s", msg.Text))
			}

			// check again tomorrow at 8pm
			tomorrow := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day()+1, 20, 0, 0, 0, c.location)
			duration := time.Until(tomorrow)
			c.logger.Info().Msg(fmt.Sprintf("sleeping for %v until %v", duration, tomorrow))
			time.Sleep(duration)
		} else {
			// send a slack reminder
			msg := slackMsg{
				Text: "hey! you haven't made a commit today yet. i'll check again in an hour. remember, you want to code!",
			}
			err = c.sendSlackMsg(msg)
			if err != nil {
				c.logger.Err(err).Msg(fmt.Sprintf("unable to send the following message to slack: %s", msg.Text))
			}

			// then sleep for an hour and check again
			c.logger.Info().Msg("sleeping for an hour to check again")
			time.Sleep(time.Hour)
		}
	}
}

type slackMsg struct {
	Text string `json:"text"`
}

// sendSlackMsg sends a request to the provided slack webhook with the provided message.
func (c *Config) sendSlackMsg(msg slackMsg) error {
	msgBody, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	_, err = c.github.Send(http.MethodPost, c.webhookURL, "", string(msgBody), "")
	return err
}
