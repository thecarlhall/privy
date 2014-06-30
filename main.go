package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"code.google.com/p/goauth2/oauth"
	"github.com/google/go-github/github"
)

const (
	Debug = false
)

func main() {
	//fmt.Println("Reading up config...")
	file, e := ioutil.ReadFile("./privy.cfg")
	if e != nil {
		log.Fatal("File error: ", e)
	}

	var config Config
	json.Unmarshal(file, &config)

	if Debug {
		fmt.Printf("Config: %v\n", config)
	}

	t := &oauth.Transport{
		Token: &oauth.Token{AccessToken: config.OauthToken},
	}

	client := github.NewClient(t.Client())

	pr := NewPullRequestor(&config, client)
	pr.ListRepos()
}
