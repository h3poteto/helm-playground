package main

import (
	"os"
	"regexp"

	"github.com/h3poteto/helm-playground/deploy"
	"github.com/h3poteto/helm-playground/repo"
	log "github.com/sirupsen/logrus"
)

func main() {
	branch := os.Getenv("BRANCH")
	github := repo.New(os.Getenv("GITHUB_TOKEN"))
	revision, err := github.GetRevision(os.Getenv("OWNER"), os.Getenv("REPOSITORY"), branch)
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("Target SHA1: %s", revision)

	stack := transformBranchName(branch)
	log.Infof("Target stack name: %s", stack)

	client, err := deploy.New(stack, revision, "", "")
	if err != nil {
		log.Fatal(err)
	}
	version, err := client.Version()
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("helm version: %s", version)
	res, err := client.NewRelease(os.Getenv("CHART_PATH"))
	if err != nil {
		log.Fatal(err)
	}
	err = client.PrintRelease(res)
	if err != nil {
		log.Fatal(err)
	}
}

func transformBranchName(name string) string {
	reg := regexp.MustCompile(`[/.*@:;~_]`)
	return reg.ReplaceAllString(name, "-")
}
