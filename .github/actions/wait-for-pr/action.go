package main

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/google/go-github/v35/github"
)

func main() {
	if len(os.Args) != 1 {
		log.Fatal("Unexpected number of arguments")
	}
	prNum, err := strconv.Atoi(os.Getenv("INPUT_PRNUMBER"))
	if err != nil {
		log.Fatalf("INPUT_PRNUMBER is not a number")
	}
	owner := os.Getenv("INPUT_OWNER")
	if owner == "" {
		log.Fatalf("INPUT_OWNER is a required value")
	}
	repo := os.Getenv("INPUT_REPO")
	if repo == "" {
		log.Fatalf("INPUT_REPO is a required value")
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
