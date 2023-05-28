package main

import (
	"github.com/aplr/pubsub-emulator/cmd"
	log "github.com/sirupsen/logrus"
)

var Version string
var Buildtime string

func main() {
	version := "local"
	if Version != "" {
		version = Version
	}

	cmd.Execute(version)
}

func init() {
	log.SetFormatter(&log.JSONFormatter{})
}