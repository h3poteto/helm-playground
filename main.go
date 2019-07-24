package main

import (
	log "github.com/sirupsen/logrus"
)

func main() {
	client, err := NewClient()
	if err != nil {
		log.Fatal(err)
	}
	log.Info(client)
}
