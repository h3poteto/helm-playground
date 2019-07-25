package main

import (
	"github.com/h3poteto/helm-playground/deploy"
	log "github.com/sirupsen/logrus"
)

func main() {
	client, err := deploy.New("", "")
	if err != nil {
		log.Fatal(err)
	}
	version, err := client.Version()
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("helm version: %s", version)
	res, err := client.NewRelease("../../lapras-inc/charts/scout")
	if err != nil {
		log.Fatal(err)
	}
	log.Info(res)
}
