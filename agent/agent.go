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

		// Check if market is open. If closed email current equity and sleep until the market reopens.
		if !client.MarketOpen {
			log.Info("Market is closed")

			totalEquity, TodaysEquity := client.GetEquity()
			mailer.Notify("Current equity: " + totalEquity + "\n" + "Today's change: " + TodaysEquity)

			sleep := client.NextOpen.Sub(time.Now())
			log.Info("Sleeping for ", sleep)
			time.Sleep(sleep)
			continue
		}

		// do trades here

		log.Info("Sleeping for 5 minutes")
		time.Sleep(5 * time.Minute)
	}
}
