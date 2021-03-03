package main

import (
	"github.com/naaltunian/paca-agent/agent"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.Info("Starting agent...")
	agent.Start()

}
