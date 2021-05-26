package main

import (
	"context"
	"github.com/google/go-github/v35/github"
    "fmt"
	"golang.org/x/oauth2"
	"log"
	"os"
	"strconv"
	"time"
)

func main() {
	if len(os.Args) != 1 {
		log.Fatal("Unexpected number of arguments")
	}
	prNum, err := strconv.Atoi(os.Getenv("PR_NUM"))
	if err != nil {
		log.Fatalf("Argument is not a number")
	}
	owner := os.Getenv("OWNER")
	if owner == "" {
		log.Fatalf("OWNER is a required value")
	}
	repo := os.Getenv("REPO")
	if repo == "" {
		log.Fatalf("REPO is a required value")
	}

	client := github.NewClient(nil)

	t := time.Tick(10 * time.Second)
	for {
		select {
		case <-t:
			pr, _, err := client.PullRequests.Get(context.Background(), owner, repo, prNum)
			if err != nil {
				log.Fatalf("Unable to fetch state for %s/%s#%d: %v", owner, repo, prNum, err)
			}
			if pr.GetState() == "closed" {
				if pr.GetMerged() {
					log.Printf("PR has been merged.")
					break
				} else {
					log.Fatalf("PR closed without merging!")
				}
			} else {
				log.Printf("Waiting for PR to merge..")
			}
		}
	}
}
