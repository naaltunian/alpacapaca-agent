package agent

import (
	"time"

	"github.com/naaltunian/paca-agent/account"
	log "github.com/sirupsen/logrus"
)

func Start() {
	for {
		Account, err := account.InitializeClient()
		if err != nil {
			log.Error(err)
			// email notifying agent is down
		}
		act := Account.GetAccount()
		if act.TradingBlocked || act.AccountBlocked {
			log.Error("Account is blocked")
			// email notifying agent is down. send act.TradingBlocked and AccountBlocked values
		}

		// do trades here

		log.Info("Sleeping for 5 minutes")
		time.Sleep(5 * time.Minute)
	}
}
