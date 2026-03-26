package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type Client struct {
	http    *http.Client
	baseUrl string
}

func NewClient(baseUrl string) *Client {
	return &Client{
		http:    http.DefaultClient,
		baseUrl: baseUrl,
	}
}

type Repo struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	LanguagesUrl string `json:"languages_url"`
}

func (c *Client) listRepos(user string) ([]Repo, error) {
	resp, err := c.http.Get(fmt.Sprintf("%s/users/%s/repos", c.baseUrl, user))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var repos []Repo
	if err = json.Unmarshal(body, &repos); err != nil {
		return nil, err
	}

	return repos, nil
}

func main() {
	client := NewClient("https://api.github.com")
	repos, err := client.listRepos("alkalax")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, repo := range repos {
		fmt.Println(repo.Name)
	}
}
