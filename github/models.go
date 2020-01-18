package github

import "time"

// Commit represents a commit on Github.
type Commit struct {
	URL    string `json:"html_url"`
	Commit commit
	Author githubAuthor
}

type commit struct {
	Author  commitAuthor
	Message string
}

type commitAuthor struct {
	Name string
	Date time.Time
}

type githubAuthor struct {
	Login string
}
