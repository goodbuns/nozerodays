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

func ScrapeContributions(username string) {
	c := colly.NewCollector()

	c.OnHTML("#day", func(e *colly.HTMLElement) {
		e.ForEach(".data-date", func(_ int, el *colly.HTMLElement) {
			// name := strings.Split(el.ChildText(".card-header"), " ")
			// g := Lottery{
			// 	Name:      strings.Join(name[:len(name)-2], " "),
			// 	CashValue: el.ChildText(".draw-cards--cash-value"),
			// 	DrawDate:  el.ChildText(".draw-cards--next-draw-date"),
			// }

			// lotteryValue := strings.Split(el.ChildText(".draw-cards--lottery-amount"), " ")
			// g.Value = strings.Join(lotteryValue, " ")
			// if lotteryValue[len(lotteryValue)-1] == "MILLION*" {
			// 	g.Millions, err = strconv.Atoi(lotteryValue[0][1:])
			// 	if err != nil {
			// 		fmt.Println("issues converting from string to int")
			// 	}
			// }
			// lotteries = append(lotteries, g)
		})
	})

	// page := scraper.Scrape("http://.com")

}
