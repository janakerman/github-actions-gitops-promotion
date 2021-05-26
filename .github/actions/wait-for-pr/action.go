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
	prNumString := os.Getenv("INPUT_PRNUMBER")
	prNum, err := strconv.Atoi(prNumString)
	if err != nil {
		log.Fatalf("INPUT_PRNUMBER %s is not a number", prNumString)
	}
	owner := os.Getenv("INPUT_OWNER")
	if owner == "" {
		log.Fatalf("INPUT_OWNER %s is a required value", owner)
	}
	repo := os.Getenv("INPUT_REPO")
	if repo == "" {
		log.Fatalf("INPUT_REPO %s is a required value", repo)
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
