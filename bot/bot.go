package bot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const (
	githubGraphQLEndpoint = "https://api.github.com/graphql"
)

// LatestContributionDate returns the date of the latest contribution made by the user.
func LatestContributionDate(username, accessToken string) time.Time {
	query := "{user(login:\"" + username + "\"){contributionsCollection{ latestRestrictedContributionDate}}}"

	client := http.Client{}

	request, err := http.NewRequest(http.MethodPost, githubGraphQLEndpoint, strings.NewReader(query))
	check(err)
	request.Header.Set("Authorization", accessToken)

	resp, err := client.Do(request)
	check(err)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	var respBody interface{}
	err = json.Unmarshal(body, &respBody)
	check(err)

	fmt.Println(respBody)
	return time.Now()
}

func check(err error) {
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
}
