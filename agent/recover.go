package agent

import (
	"fmt"

	"github.com/naaltunian/paca-agent/mailer"
	log "github.com/sirupsen/logrus"
)

// Recover will alert the user via email that the agent is down. TODO: sell all positions??
func recovery() {
	if r := recover(); r != nil {
		log.Error(r)
		errStr := fmt.Sprintf("%v", r)
		mailer.Notify("PANIC", "Agent panic: "+errStr)
	}
}
