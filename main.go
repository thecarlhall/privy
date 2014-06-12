package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"code.google.com/p/goauth2/oauth"
	"github.com/google/go-github/github"
)

func ListRepos(config *Config) {
	t := &oauth.Transport{
		Token: &oauth.Token{AccessToken: config.OauthToken},
	}

	client := github.NewClient(t.Client())

	// list all repositories for the authenticated user
	repos, _, err := client.Repositories.ListByOrg(config.Organization, nil)

	if err != nil {
		fmt.Println(err)
	}

	//fmt.Println(github.Stringify(repos))
	fmt.Println("Repositories")
	for _, repo := range repos {
		fmt.Printf("  %s\n", *repo.Name)
	}
}

func main() {
	file, e := ioutil.ReadFile("./privy.cfg")
	if e != nil {
		fmt.Printf("File error: %v\n", e)
		os.Exit(1)
	}

	var config Config
	json.Unmarshal(file, &config)
	fmt.Printf("Config: %v\n", config)

	ListRepos(&config)
}
