package main

import (
	"bytes"
	"fmt"
	"log"
	"strings"

	"github.com/google/go-github/github"

	"github.com/wsxiaoys/terminal"
	"github.com/wsxiaoys/terminal/color"
)

func NewPullRequestor(config *Config, client *github.Client) *PullRequestor {
	return &PullRequestor{config, client}
}

type PullRequestor struct {
	config *Config
	client *github.Client
}

func (self *PullRequestor) writePullRequest(organization string, project string, pr github.PullRequest, out chan<- string) {
	if self.config.Debug {
		log.Println("Listing comments")
	}
	comments, _, err := self.client.PullRequests.ListComments(organization, project, *pr.Number, nil)

	if err != nil {
		log.Fatal(err)
	}

	var buffer bytes.Buffer

	buffer.WriteString(color.Sprintf("@w[%d] %s\n", *pr.Number, *pr.Title))

	buffer.WriteString(color.Sprintf("@{/}(Created on %s :: %d comments)\n", *pr.CreatedAt, len(comments)))

	if pr.Body == nil || len(*pr.Body) == 0 {
		buffer.WriteString(color.Sprint("@r<no body>\n"))
	} else {
		bodyLen := len(*pr.Body)
		if bodyLen > 120 {
			bodyLen = 120
		}

		buffer.WriteString(color.Sprintf("@b%s\n", (*pr.Body)[:bodyLen]))
	}

	buffer.WriteString(color.Sprintf("@g%s\n\n", (*pr.HTMLURL)))

	if self.config.Debug {
		log.Println("Done printing pull requests")
	}

	out <- buffer.String()
}

func (self *PullRequestor) PrintPullRequests(repo Repository, done chan<- bool) {
	for _, project := range repo.Projects {
		var buffer bytes.Buffer

		if self.config.Debug {
			log.Println("Getting pull requests for", project)
		}
		prs, _, err := self.client.PullRequests.List(repo.Organization, project, nil)

		if err != nil {
			log.Fatal(err)
		}

		if len(prs) == 0 {
			continue
		}

		// use a channel to collect all of the writes to a buffer then write it all at once
		printer := make(chan string)
		defer close(printer)
		for _, pr := range prs {
			go self.writePullRequest(repo.Organization, project, pr, printer)
		}

		// write out the printer to the terminal after everything has been collected
		terminal.Stdout.Color("y")

		// Get color codes here: https://github.com/wsxiaoys/terminal/blob/master/color/color.go
		title := fmt.Sprintf(" [ %s ] ", strings.ToUpper(project))
		paddingWidth := (80 - len(title)) / 2

		buffer.WriteString(color.Sprintf("%s\n", strings.Repeat("=", 80)))
		buffer.WriteString(color.Sprintf("@{!m}%s%s%s\n", strings.Repeat("-", paddingWidth), title, strings.Repeat("-", paddingWidth)))

		for i := 0; i < len(prs); i++ {
			buffer.WriteString(<-printer)
		}
		color.Print(buffer.String())
	}
	done <- true
}

func (self *PullRequestor) ListRepos() {
	// list all repositories for the authenticated user
	done := make(chan bool)
	defer close(done)
	for _, repo := range self.config.Repositories {
		go self.PrintPullRequests(repo, done)
	}

	for i := 0; i < len(self.config.Repositories); i++ {
		<-done
	}
}
