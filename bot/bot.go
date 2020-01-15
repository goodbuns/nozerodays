package bot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/goodbuns/nozerodays/github"
)

const (
	githubGraphQLEndpoint = "https://api.github.com/graphql"
	githubV3Endpoint      = "https://api.github.com"
)

type Config struct {
	Username    string
	AccessToken string
	Client      *github.Client
	// Logger      *log.Logger
}

func New(username, accessToken string) *Config {
	return &Config{
		Username:    username,
		AccessToken: accessToken,
		Client:      github.New(),
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

	// todo: have orgs be some sort of whitelist of organizations, passed in via flags
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

// features:
// 1. remind 8pm everyday on slack if commit hasn't been made by user
// 2. keep track of metrics for how many days i did it

// get list of all organizations
// check whether i'm an admin of the organization
// create webhooks for all repos in organization that do not yet have expected webhook
// create new webhook on new repos in my account or goodbuns
// receive webhooks
// create API for webhook
// when webhook fires, update database w latest time commit
// worker/server or cron/batch job?
