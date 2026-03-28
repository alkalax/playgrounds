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
	Name        string `json:"name"`
	Description string `json:"description"`
	Owner       struct {
		Login string `json:"login"`
	} `json:"owner"`
	LanguagesUrl string `json:"languages_url"`
	Languages    []string
}

func (c *Client) get(path string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "github-explorer")
	req.Header.Set("Authorization", "token "+c.token)

	return c.http.Do(req)
}

func (c *Client) searchRepos(query string) ([]Repo, error) {
	resp, err := c.get(fmt.Sprintf("%s/search/repositories?q=%s&per_page=10", c.baseUrl, query))
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

	var respObj struct {
		Repos []Repo `json:"items"`
	}
	if err = json.Unmarshal(body, &respObj); err != nil {
		return nil, err
	}

	return respObj.Repos, nil
}

func (c *Client) listRepos(user string) ([]Repo, error) {
	resp, err := c.get(fmt.Sprintf("%s/users/%s/repos", c.baseUrl, user))
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
		langResp, err := c.get(repos[i].LanguagesUrl)
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
	//repos, err := client.listRepos("alkalax")
	repos, err := client.searchRepos("language:go")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for i, repo := range repos {
		//fmt.Println(repo.Name, repo.Languages)
		if i > 0 {
			fmt.Println("==============================")
		}
		fmt.Println("Name:", repo.Name)
		fmt.Println("Owner:", repo.Owner.Login)
		fmt.Println(repo.Description)
	}
}
