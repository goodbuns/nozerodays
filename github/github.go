package github

import (
	"net/http"
	"strings"
)

type Client struct {
	Client *http.Client
}

func New() *Client {
	return &Client{Client: &http.Client{}}
}

func (c *Client) SendRequest(method, url, path, body, authToken string) (*http.Response, error) {
	fullURL := url + path
	req, err := http.NewRequest(method, fullURL, strings.NewReader(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "bearer "+authToken)

	return c.Client.Do(req)
}

// func (c *Client) ReadResponse(resp *http.Response, respBody *[]interface{}) error {
// 	defer resp.Body.Close()
// 	body, err := ioutil.ReadAll(resp.Body)
// 	c.check(err)

// 	err = json.Unmarshal(body, respBody)
// 	c.check(err)

// 	return nil
// }

// func (c *Client) check(err error) {
// 	if err != nil {
// 		fmt.Println(err)
// 		panic(err)
// 	}
// }
