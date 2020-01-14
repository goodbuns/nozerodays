package github

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type Client struct {
	Client *http.Client
}

func (c *Client) New() *Client {
	return &Client{Client: &http.Client{}}
}

func (c *Client) SendRequest(method, url, body, authToken string) (*http.Response, error) {
	req, err := http.NewRequest(method, url, strings.NewReader(body))
	c.check(err)

	req.Header.Set("Authorization", "bearer "+authToken)

	return c.Client.Do(req)
}

func (c *Client) ReadResponse(resp *http.Response, respBody *interface{}) error {
	body, err := ioutil.ReadAll(resp.Body)
	c.check(err)

	err = json.Unmarshal(body, respBody)
	c.check(err)

	return nil
}

func (c *Client) check(err error) {
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}
