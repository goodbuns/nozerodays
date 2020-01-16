package bot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/goodbuns/nozerodays/github"
)

const (
	githubGraphQLEndpoint = "https://api.github.com/graphql"
	githubV3Endpoint      = "https://api.github.com"
)

type Config struct {
	Username    string
	AccessToken string
	WebhookURL  string
	Client      *github.Client
	// Logger      *log.Logger
}

type GithubCommit struct {
	URL    string `json:"html_url"`
	Commit Commit
	Author GithubAuthor
}

type Commit struct {
	Author  CommitAuthor
	Message string
}

type CommitAuthor struct {
	Name string
	Date time.Time
}

type GithubAuthor struct {
	Login string
}

func New(username, accessToken, webhookURL string) *Config {
	return &Config{
		Username:    username,
		AccessToken: accessToken,
		Client:      github.New(),
		WebhookURL:  webhookURL,
		// Logger:      log.New(),
	}
}

// todo: figure out logging instead of panicking
func check(err error) {
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
}

func (c *Config) Repositories() []string {
	// get all repositories that user owns
	resp, err := c.Client.SendRequest(http.MethodGet, githubV3Endpoint, "/user/repos?type=owner", ``, c.AccessToken)

	defer resp.Body.Close()
	check(err)

	type RepositoryResponse struct {
		FullName string `json:"full_name"`
	}
	var repos []RepositoryResponse
	body, err := ioutil.ReadAll(resp.Body)
	check(err)

	err = json.Unmarshal(body, &repos)
	check(err)

	// todo: have orgs be some sort of whitelist of organizations, passed in via
	// flags
	orgs := []string{"goodbuns"}
	// get all repos within whitelisted orgs
	var orgRepos []RepositoryResponse
	for _, org := range orgs {
		resp, err = c.Client.SendRequest(http.MethodGet, githubV3Endpoint, "/orgs/"+org+"/repos", "", c.AccessToken)
		defer resp.Body.Close()
		check(err)

		body, err = ioutil.ReadAll(resp.Body)
		check(err)

		err = json.Unmarshal(body, &orgRepos)
		check(err)
	}

	repos = append(repos, orgRepos...)

	var r []string
	for _, repo := range repos {
		r = append(r, repo.FullName)
	}
	return r
}

func (c *Config) Commits(repos []string) *GithubCommit {
	// go through commit history for each repo until (1) the date is before
	// today's date or (2) i am the author and the date is today if i am the
	// author, return if not, continue if we go through all commits for the day
	// in all repos and i am not the author of any of those commits, i have not yet
	// made a commit for this day

	var commits []GithubCommit
	today := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.UTC)

	for _, repo := range repos {
		resp, err := c.Client.SendRequest(http.MethodGet, githubV3Endpoint, "/repos/"+repo+"/commits", "", c.AccessToken)
		defer resp.Body.Close()
		check(err)

		body, err := ioutil.ReadAll(resp.Body)
		check(err)

		err = json.Unmarshal(body, &commits)
		check(err)

		// really only need to look at the latest commit
		if commits[0].Author.Login == c.Username {
			if today.Before(commits[0].Commit.Author.Date) {
				return &(commits[0])
			}
		}
	}
	return nil
}

type SlackMessage struct {
	Text string `json:"text"`
}

func (c *Config) Start() {
	for true {
		r := c.Repositories()
		commit := c.Commits(r)
		if c != nil {
			fmt.Println("found commit for today", *c)

			// send slack message
			message := SlackMessage{
				Text: "good work! commit was made today at " + commit.Commit.Author.Date.String() + ", commit link: " + commit.URL,
			}
			c.SlackMessage(message)

			// check again tomorrow at 8pm
			tomorrow := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day()+1, 20, 0, 0, 0, time.Local)
			duration := time.Until(tomorrow.Local())
			fmt.Println("sleeping for ", duration, "until", tomorrow.Local(), "now", time.Now())
			time.Sleep(duration)
		} else {
			// do a reminder
			message := SlackMessage{
				Text: "hey! it's 8pm and you haven't made a commit today yet. i'll check again in an hour. remember, you want to code!",
			}
			c.SlackMessage(message)
			// then sleep for an hour and check again
			time.Sleep(time.Hour * 4)
		}
	}

}

func (c *Config) SlackMessage(message SlackMessage) {
	msgBody, err := json.Marshal(message)
	check(err)

	_, err = c.Client.SendRequest(http.MethodPost, c.WebhookURL, "", string(msgBody), "")
	check(err)
}

// features: 1. remind 8pm everyday on slack if commit hasn't been made by user
// 2. keep track of metrics for how many days i did it

// get list of all organizations check whether i'm an admin of the organization
// create webhooks for all repos in organization that do not yet have expected
// webhook create new webhook on new repos in my account or goodbuns receive
// webhooks create API for webhook when webhook fires, update database w latest
// time commit worker/server or cron/batch job?
