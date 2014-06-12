package main

import (
	"fmt"

	"code.google.com/p/goauth2/oauth"
	"github.com/google/go-github/github"
)

func main() {
	t := &oauth.Transport{
		Token: &oauth.Token{AccessToken: ""},
	}

	client := github.NewClient(t.Client())

	// list all repositories for the authenticated user
	repos, _, err := client.Repositories.ListByOrg("", nil)

	if err != nil {
		fmt.Println(err)
	}

	//fmt.Println(github.Stringify(repos))
	fmt.Println("Repositories")
	for _, repo := range repos {
		fmt.Printf("  %s\n", *repo.Name)
	}
}
