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

func ListRepos(config *Config) {
	t := &oauth.Transport{
		Token: &oauth.Token{AccessToken: config.OauthToken},
	}

	client := github.NewClient(t.Client())

	// list all repositories for the authenticated user
	for _, repo := range config.Repositories {
		prs, _, _ := client.PullRequests.List(config.Organization, repo, nil)

		if len(prs) == 0 {
			continue
		}

		terminal.Stdout.Color("y")
		// Get color codes here: https://github.com/wsxiaoys/terminal/blob/master/color/color.go
		color.Printf("@{!kW}**** Pull requests for [%s]", repo)
		color.Println()
		for _, pr := range prs {
			color.Printf("@w%d: %s\n", *pr.Number, *pr.Title)
			if len(*pr.Body) == 0 {
				color.Println("@r<no body>")
			} else {
				color.Printf("@b%s\n", *pr.Body)
			}

			comments, _, _ := client.PullRequests.ListComments(config.Organization, repo, *pr.Number, nil)
			color.Printf("@{/}(%d comments)\n", len(comments))
			color.Printf("@g%s\n", *pr.HTMLURL)
			color.Println()
		}
	}
}

func main() {
	file, e := ioutil.ReadFile("./privy.cfg")
	if e != nil {
		log.Fatal("File error: ", e)
	}

	var config Config
	json.Unmarshal(file, &config)
	//fmt.Printf("Config: %v\n", config)

	ListRepos(&config)
}
