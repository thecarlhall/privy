package main

import (
	"log"
	"sync"

	"github.com/google/go-github/github"

	"github.com/wsxiaoys/terminal"
	"github.com/wsxiaoys/terminal/color"
)

var (
	mut sync.Mutex
)

func NewPullRequestor(config *Config, client *github.Client) *PullRequestor {
	return &PullRequestor{config, client}
}

type PullRequestor struct {
	config *Config
	client *github.Client
}

func (self *PullRequestor) printPullRequest(repo string, pr github.PullRequest, wg *sync.WaitGroup) {
	comments, _, err := self.client.PullRequests.ListComments(self.config.Organization, repo, *pr.Number, nil)

	if err != nil {
		log.Fatal(err)
	}

	// Keep the mutex down here so we can make the list requests in parallel, but
	// not clobber stdout while printing stuff. A better approach would be to
	// throw this in a channel.
	mut.Lock()

	color.Printf("@w[%d] %s\n", *pr.Number, *pr.Title)

	if pr.Body == nil || len(*pr.Body) == 0 {
		color.Println("@r<no body>")
	} else {
		color.Printf("@b%s\n", (*pr.Body)[:120])
	}

	color.Printf("@{/}(%d comments)\n", len(comments))
	color.Printf("@g%s\n", (*pr.HTMLURL))
	color.Println()
	mut.Unlock()

	// Signals we're done printing this one element.
	wg.Done()
}

func (self *PullRequestor) PrintPullRequests(repo string) {
	prs, _, err := self.client.PullRequests.List(self.config.Organization, repo, nil)

	if err != nil {
		log.Fatal(err)
	}

	if len(prs) == 0 {
		return
	}

	terminal.Stdout.Color("y")

	// Get color codes here: https://github.com/wsxiaoys/terminal/blob/master/color/color.go
	color.Println("================================================================================")
	color.Printf("@{!m}**** Pull requests for [%s]", repo)
	color.Println()

	var wg sync.WaitGroup

	for _, pr := range prs {
		wg.Add(1)
		go self.printPullRequest(repo, pr, &wg)
	}

	// Wait for the threads to finish their stuff.
	wg.Wait()
}

func (self *PullRequestor) ListRepos() {
	// list all repositories for the authenticated user
	for _, repo := range self.config.Repositories {
		self.PrintPullRequests(repo)
	}
}