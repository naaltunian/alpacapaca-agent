package agent

import (
	"fmt"

	"github.com/naaltunian/paca-agent/mailer"
)

// Recover will alert the user via email that the agent is down. TODO: sell all positions??
func recovery() {
	if r := recover(); r != nil {
		errStr := fmt.Sprintf("%v", r)
		mailer.Notify("PANIC", "Agent panic: "+errStr)
	}
}
