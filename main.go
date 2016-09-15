package main

import (
	"fmt"
	"os"

	"github.com/google/go-github/github"
	"github.com/pborman/getopt"
)

func buildMsg(header string, msg string, footer string) string {
	return fmt.Sprintf("%s\n\n%s\n%s\n", header, msg, footer)
}

func processProject(client *github.Client, org string, project string) string {
	prs, _, err := client.PullRequests.List(org, project, nil)
	if err != nil {
		fmt.Println("error listing prs", err)
		os.Exit(1)
	}

	var msg string
	heading := "We still have OPEN PRs to be reviewed:"
	footer := ""

	for _, pr := range prs {
		// Title, User, URL
		// string, string, string
		msg += fmt.Sprintf("%s\tauthor: (%s)\n%s\n---\n", *pr.Title, *pr.User.Login, *pr.URL)
	}

	return buildMsg(heading, msg, footer)
}

func main() {
	send := *getopt.BoolLong("send", 's', "", "Send Email")
	email := *getopt.StringLong("email", 'e', "", "email address")
	org := *getopt.StringLong("org", 'o', "fusor", "org")
	projects := *getopt.ListLong("project", 'p', "Project")

	getopt.Parse()

	fmt.Println(projects)
	if len(projects) <= 0 {
		fmt.Println("no projects, defaulting to fusor")
		projects = []string{"fusor"}
	}

	client := github.NewClient(nil)

	for _, project := range projects {
		msg := processProject(client, org, project)
		if send {
			fmt.Println(email)
			//fmt.Println(msg)
		} else {
			fmt.Println(msg)
		}
	}

}
