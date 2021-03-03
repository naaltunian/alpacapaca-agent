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
		client, err := account.InitializeClient()
		if err != nil {
			log.Error("Error initializing client: ", err)
			// email notifying agent is down
			mailer.Notify("Could not initialize client: " + err.Error())
		}

		acct := client.GetAccount()
		if acct.TradingBlocked || acct.AccountBlocked {
			log.Error("Account is blocked")
			// email notifying agent is down.
			mailer.Notify("Account is blocked. Trading Blocked: " + strconv.FormatBool(acct.TradingBlocked) + " Account Blocked: " + strconv.FormatBool(acct.AccountBlocked))
		}

		// Check if market is open. If closed sleep until open.
		if !client.MarketOpen {
			log.Info("Market is closed")
			// log.Info("clock ", clock.NextOpen.Sub(time.Now()))
			sleep := client.NextOpen.Sub(time.Now())
			time.Sleep(sleep)
			continue
		}

		log.Info("Account equity ", acct.Equity.Floor())
		log.Info("Account buying power ", acct.BuyingPower)

		// do trades here

		log.Info("Sleeping for 5 minutes")
		time.Sleep(5 * time.Minute)
	}
}
