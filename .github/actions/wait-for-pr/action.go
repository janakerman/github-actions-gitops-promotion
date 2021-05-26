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
	prNumString := os.Getenv("PR_NUM")
	prNum, err := strconv.Atoi(prNumString)
	if err != nil {
		log.Fatalf("PR_NUM %s is not a number", prNumString)
	}
	owner := os.Getenv("OWNER")
	if owner == "" {
		log.Fatalf("OWNER %s is a required value", owner)
	}
	repo := os.Getenv("REPO")
	if repo == "" {
		log.Fatalf("REPO %s is a required value", repo)
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
