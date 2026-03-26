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
	token   string
}

func NewClient(baseUrl, token string) *Client {
	return &Client{
		http:    http.DefaultClient,
		baseUrl: baseUrl,
		token:   token,
	}
}

type Repo struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	LanguagesUrl string `json:"languages_url"`
	Languages    []string
}

func (c *Client) listRepos(user string) ([]Repo, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/users/%s/repos", c.baseUrl, user), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "github-explorer")
	req.Header.Set("Authorization", "token "+c.token)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("github status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var repos []Repo
	if err = json.Unmarshal(body, &repos); err != nil {
		return nil, err
	}

	for i := range repos {
		langReq, err := http.NewRequest(http.MethodGet, repos[i].LanguagesUrl, nil)
		if err != nil {
			return nil, err
		}
		langReq.Header.Set("Accept", "application/vnd.github+json")
		langReq.Header.Set("User-Agent", "github-explorer")
		langReq.Header.Set("Authorization", "token "+c.token)

		langResp, err := c.http.Do(langReq)
		if err != nil {
			return nil, err
		}
		defer langResp.Body.Close()

		langBody, err := io.ReadAll(langResp.Body)
		if err != nil {
			return nil, err
		}

		var languages map[string]int
		if err = json.Unmarshal(langBody, &languages); err != nil {
			return nil, err
		}

		repos[i].Languages = []string{}
		for language := range languages {
			repos[i].Languages = append(repos[i].Languages, language)
		}
	}

	return repos, nil
}

func main() {
	client := NewClient("https://api.github.com", os.Getenv("GITHUB_TOKEN"))
	repos, err := client.listRepos("alkalax")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, repo := range repos {
		fmt.Println(repo.Name, repo.Languages)
	}
}
