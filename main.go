package main

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"code.google.com/p/goauth2/oauth"
	"github.com/google/go-github/github"

	"github.com/wsxiaoys/terminal"
	"github.com/wsxiaoys/terminal/color"
)

func printPullRequest(client *github.Client, config *Config, repo *string, pr *github.PullRequest) {
	color.Printf("@w[%d] %s\n", *pr.Number, *pr.Title)
	if len(*pr.Body) == 0 {
		color.Println("@r<no body>")
	} else {
		color.Printf("@b%s\n", (*pr.Body)[:120])
	}

	comments, _, err := client.PullRequests.ListComments(config.Organization, *repo, *pr.Number, nil)

	if err != nil {
		log.Fatal(err)
	}
	color.Printf("@{/}(%d comments)\n", len(comments))
	color.Printf("@g%s\n", *pr.HTMLURL)
	color.Println()
}

func printPullRequests(client *github.Client, config *Config, repo *string) {
	prs, _, err := client.PullRequests.List(config.Organization, *repo, nil)

	if err != nil {
		log.Fatal(err)
	}

	if len(prs) == 0 {
		return
	}

	terminal.Stdout.Color("y")
	// Get color codes here: https://github.com/wsxiaoys/terminal/blob/master/color/color.go
	color.Println("================================================================================")
	color.Printf("@{!m}**** Pull requests for [%s]", *repo)
	color.Println()
	for _, pr := range prs {
		printPullRequest(client, config, repo, &pr)
	}
}

func listRepos(config *Config) {
	t := &oauth.Transport{
		Token: &oauth.Token{AccessToken: config.OauthToken},
	}

	client := github.NewClient(t.Client())

	// list all repositories for the authenticated user
	for _, repo := range config.Repositories {
		printPullRequests(client, config, &repo)
	}
}

func main() {
	//fmt.Println("Reading up config...")
	file, e := ioutil.ReadFile("./privy.cfg")
	if e != nil {
		log.Fatal("File error: ", e)
	}

	var config Config
	json.Unmarshal(file, &config)
	//fmt.Printf("Config: %v\n", config)

	listRepos(&config)
}
