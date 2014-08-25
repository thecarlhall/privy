package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os/user"

	"code.google.com/p/goauth2/oauth"
	"github.com/google/go-github/github"
)

func homeDir() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	return usr.HomeDir
}

func main() {
	file, e := ioutil.ReadFile(fmt.Sprintf("%s/.privy.cfg", homeDir()))
	if e != nil {
		log.Fatal("File error: ", e)
	}

	var config Config
	json.Unmarshal(file, &config)

	if config.Debug {
		log.Println("Config", config)
	}

	t := &oauth.Transport{
		Token: &oauth.Token{AccessToken: config.OauthToken},
	}

	if config.Debug {
		log.Println("Creating GitHub client")
	}
	client := github.NewClient(t.Client())

	pr := NewPullRequestor(&config, client)
	pr.ListRepos()
}
