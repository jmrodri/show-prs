/*
 * main.go (show-prs) - list open pull requests from your github project
 *
 * Copyright (C) 2016 Jesus M. Rodriguez
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 2 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful, but
 * WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, write to the Free Software
 * Foundation, Inc., 59 Temple Place - Suite 330, Boston, MA 02111-1307, USA.
 */
package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/smtp"
	"os"
	"strconv"
	"strings"

	"golang.org/x/oauth2"

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

	if len(prs) <= 0 {
		// nothing to send
		return ""
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
	buf := bytes.NewBufferString("Subject: PR report for: " + project + "\n" + msg)
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
	getopt.StringLong("token", 'a', "", "Token")

	getopt.Parse()

	org := getopt.GetValue("org")
	to := getopt.GetValue("to")
	from := getopt.GetValue("from")
	token := getopt.GetValue("token")
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

	var tc *http.Client

	if token != "" {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		)
		tc = oauth2.NewClient(oauth2.NoContext, ts)
	}

	client := github.NewClient(tc)

	for _, project := range projects {
		msg := processProject(client, org, project)
		if send {
			// we could optimize this to be if send && msg but
			// then that triggers the else block printing an empty string.
			//
			// only send email if we actually have a message to send
			if msg != "" {
				sendEmail(from, to, project, msg)
			}
		} else {
			fmt.Println(msg)
		}
	}
}
