package bot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/gocolly/colly"
)

const (
	githubGraphQLEndpoint = "https://api.github.com/graphql"
)

// LatestContributionDate returns the date of the latest contribution made by
// the user.
func LatestContributionDate(username, accessToken string) time.Time {
	queryString := `query { user(login: "` + username + `") { contributionsCollection { totalCommitContributions restrictedContributionsCount hasAnyContributions startedAt endedAt latestRestrictedContributionDate }}}`

	type reqBody struct {
		Query string `json:"query"`
	}
	r, err := json.Marshal(reqBody{Query: queryString})
	check(err)

	req, err := http.NewRequest(http.MethodPost, githubGraphQLEndpoint, bytes.NewReader(r))
	check(err)

	// Set Authorization token.
	req.Header.Set("Authorization", "bearer "+accessToken)

	dump, err := httputil.DumpRequestOut(req, true)
	check(err)
	fmt.Printf("%q", dump)

	client := http.Client{}
	resp, err := client.Do(req)
	check(err)

	dump, err = httputil.DumpResponse(resp, true)
	check(err)
	fmt.Printf("%q", dump)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	var respBody interface{}
	err = json.Unmarshal(body, &respBody)
	check(err)

	fmt.Println("\n\n\n", respBody)
	return time.Now()
}

func check(err error) {
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
}

// get list of all organizations
// check whether i'm an admin of the organization
// create webhooks for all repos in organization that do not yet have expected webhook
// create new webhook on new repos in my account or goodbuns
// receive webhooks
