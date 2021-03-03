package agent

import (
	"strconv"
	"time"

	"github.com/naaltunian/paca-agent/account"
	"github.com/naaltunian/paca-agent/mailer"
	log "github.com/sirupsen/logrus"
)

func Start() {
	for {
		Account, err := account.InitializeClient()
		if err != nil {
			log.Error("Error initializing client: ", err)
			// email notifying agent is down
			mailer.Notify("Could not initialize client: " + err.Error())
		}
		act := Account.GetAccount()
		if act.TradingBlocked || act.AccountBlocked {
			log.Error("Account is blocked")
			// email notifying agent is down.
			mailer.Notify("Account is blocked. Trading Blocked: " + strconv.FormatBool(act.TradingBlocked) + " Account Blocked: " + strconv.FormatBool(act.AccountBlocked))
		}

		// do trades here

		log.Info("Sleeping for 5 minutes")
		time.Sleep(5 * time.Minute)
	}
}
