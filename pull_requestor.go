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
		log.Println("Listing comments for", project)
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
		log.Println("Done printing pull requests for", project)
	}

	out <- buffer.String()
}

func (self *PullRequestor) PrintPullRequests(repo Repository, done chan<- struct{}) {
	// write out the printer to the terminal after everything has been collected
	terminal.Stdout.Color("y")

	outs := make([]chan string, len(repo.Projects))

	for i, project := range repo.Projects {
		outs[i] = make(chan string)
		defer close(outs[i])

		go func(project string, out chan<- string) {
			if self.config.Debug {
				log.Println("Getting pull requests for", project)
			}
			prs, _, err := self.client.PullRequests.List(repo.Organization, project, nil)

			if err != nil {
				log.Fatal(err)
			}

			if len(prs) == 0 {
				out <- ""
				return
			}

			// use a channel to collect all of the writes to a buffer then write it all at once
			printer := make(chan string)
			defer close(printer)

			for _, pr := range prs {
				go self.writePullRequest(repo.Organization, project, pr, printer)
			}

			// Get color codes here: https://github.com/wsxiaoys/terminal/blob/master/color/color.go
			width := 80
			title := fmt.Sprintf(" [ %s (%d) ] ", strings.ToUpper(project), len(prs))
			paddingWidth := (width - len(title)) / 2

			header := strings.Repeat("=", width)
			padding := strings.Repeat("-", paddingWidth)

			var buffer bytes.Buffer
			buffer.WriteString(color.Sprintf("@y%s\n", header))
			buffer.WriteString(color.Sprintf("@{!m}%s%s%s\n", padding, title, padding))

			for i := 0; i < len(prs); i++ {
				buffer.WriteString(<-printer)
			}

			out <- buffer.String()
		}(project, outs[i])
	}

	for _, out := range outs {
		color.Print(<-out)
	}

	done <- struct{}{}
}

func (self *PullRequestor) ListRepos() {
	// list all repositories for the authenticated user
	done := make(chan struct{})
	defer close(done)

	for _, repo := range self.config.Repositories {
		go self.PrintPullRequests(repo, done)
	}

	for i := 0; i < len(self.config.Repositories); i++ {
		<-done
	}
}
