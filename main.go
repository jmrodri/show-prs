package main

import (
	"bytes"
	"fmt"
	"net/smtp"
	"os"
	"strconv"
	"strings"

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
		// TODO: consider using a template
		msg += fmt.Sprintf("%s\tauthor: (%s)\n%s\n---\n", *pr.Title, *pr.User.Login, *pr.URL)
	}

	return buildMsg(heading, msg, footer)
}

func sendEmail(from_addr string, to_addr string, project string, msg string) {
	fmt.Println("Sending email to: ", to_addr)

	c, err := smtp.Dial("smtp.corp.redhat.com:25")
	if err != nil {
		fmt.Println("Could not send email", err)
		os.Exit(1)
	}
	defer c.Close()

	c.Mail(from_addr)
	c.Rcpt(to_addr)
	wc, err := c.Data()
	if err != nil {
		fmt.Println("Could not send email 2", err)
		os.Exit(1)
	}
	defer wc.Close()
	buf := bytes.NewBufferString("Subject: PR report for:" + project + "\n" + msg)
	if _, err = buf.WriteTo(wc); err != nil {
		fmt.Println("Could not write message", err)
	}
}

func main() {
	getopt.BoolLong("send", 's', "", "Send Email otherwise just prints to stdout")
	getopt.StringLong("from", 'f', "", "Sender email address")
	getopt.StringLong("to", 't', "", "Recipient email address")
	getopt.StringLong("org", 'o', "fusor", "org")
	getopt.ListLong("project", 'p', "Project")

	getopt.Parse()

	org := getopt.GetValue("org")
	to := getopt.GetValue("to")
	from := getopt.GetValue("from")
	projects := strings.Split(getopt.GetValue("project"), ",")
	send, _ := strconv.ParseBool(getopt.GetValue("send"))

	if send {
		// if we are sending email we need to know the from and to addresses
		if from == "" || to == "" {
			fmt.Println("You must specify a from and to email address with the send option")
			os.Exit(1)
		}
	}

	if len(projects) == 1 && projects[0] == "" {
		fmt.Println("no projects, defaulting to fusor")
		projects = []string{"fusor"}
	}

	client := github.NewClient(nil)

	for _, project := range projects {
		msg := processProject(client, org, project)
		if send {
			sendEmail(from, to, project, msg)
		} else {
			fmt.Println(msg)
		}
	}
}
