package github

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"
)

const (
	githubV3Endpoint = "https://api.github.com"
)

// Client represents a Github HTTP client.
type Client struct {
	username    string
	accessToken string
	client      *http.Client
	orgs        []string
}

// New creates a new Github client.
func New(username, accessToken string, orgs []string) *Client {
	return &Client{
		username:    username,
		accessToken: accessToken,
		client:      &http.Client{},
		orgs:        orgs,
	}
}

// Send sends an HTTP request.
func (c *Client) Send(method, url, path, body, authToken string) (*http.Response, error) {
	fullURL := url + path
	req, err := http.NewRequest(method, fullURL, strings.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "bearer "+authToken)
	return c.client.Do(req)
}

// Repositories returns a slice of strings representing owned repositories of
// the user and whitelisted organizations provided.
func (c *Client) Repositories() ([]string, error) {
	var r []string
	type repositoryResponse struct {
		FullName string `json:"full_name"`
	}

	// Get all repos that user owns.
	resp, err := c.Send(http.MethodGet, githubV3Endpoint, "/user/repos?type=owner", ``, c.accessToken)
	defer resp.Body.Close()
	if err != nil {
		return r, errors.Wrap(err, "failed to send request successfully")
	}
	var repos []repositoryResponse
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return r, errors.Wrap(err, "failed to read response body")
	}
	err = json.Unmarshal(body, &repos)
	if err != nil {
		return r, errors.Wrap(err, "failed to unmarshal response")
	}

	// Get all repos of whitelisted orgs.
	var orgRepos []repositoryResponse
	for _, org := range c.orgs {
		resp, err = c.Send(http.MethodGet, githubV3Endpoint, "/orgs/"+org+"/repos", "", c.accessToken)
		defer resp.Body.Close()
		if err != nil {
			return r, errors.Wrap(err, "failed to send request successfully")
		}

		body, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return r, errors.Wrap(err, "failed to read response body")
		}

		err = json.Unmarshal(body, &orgRepos)
		if err != nil {
			return r, errors.Wrap(err, "failed to unmarshal response")
		}
	}

	repos = append(repos, orgRepos...)

	for _, repo := range repos {
		r = append(r, repo.FullName)
	}
	return r, nil
}

// CommitCreatedToday checks the latest commit of each repo until it finds a commit that
// was created today.
func (c *Client) CommitCreatedToday(repos []string, location *time.Location) (*Commit, error) {
	var commits []Commit
	today := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, location)

	for _, repo := range repos {
		resp, err := c.Send(http.MethodGet, githubV3Endpoint, "/repos/"+repo+"/commits", "", c.accessToken)
		defer resp.Body.Close()
		if err != nil {
			return nil, errors.Wrap(err, "failed to send request successfully")
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, errors.Wrap(err, "failed to read response body")
		}
		err = json.Unmarshal(body, &commits)
		if err != nil {
			return nil, errors.Wrap(err, "failed to unmarshal response")
		}

		// Look at latest commit and check whether it was made today.
		if commits[0].Author.Login == c.username {
			localCommitTime := commits[0].Commit.Author.Date.In(location)
			if today.Before(localCommitTime) {
				return &(commits[0]), nil
			}
		}
	}
	return nil, nil
}
