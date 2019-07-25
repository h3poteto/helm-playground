package main

import (
	log "github.com/sirupsen/logrus"
)

func main() {
	client, err := NewClient()
	if err != nil {
		log.Fatal(err)
	}
	version, err := client.GetVersion()
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("helm version: %s", version.Version.SemVer)
}
