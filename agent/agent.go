package agent

import (
	"os"
	"strconv"
	"time"

	"github.com/naaltunian/paca-agent/account"
	"github.com/naaltunian/paca-agent/mailer"
	log "github.com/sirupsen/logrus"
)

func Start() {
	for {
		// Initialize account and get current account information/balance
		profile, err := account.InitializeClient()
		if err != nil {
			log.Error("Error initializing client: ", err)
			// email notifying agent is down
			mailer.Notify("Could not initialize client: " + err.Error())

			continue
		}

		// Get user's account information
		acct := profile.GetAccount()
		if acct.TradingBlocked || acct.AccountBlocked {
			log.Error("Account is blocked")
			// email notifying agent is down.
			mailer.Notify("Account is blocked. Trading Blocked: " + strconv.FormatBool(acct.TradingBlocked) + " Account Blocked: " + strconv.FormatBool(acct.AccountBlocked))

			os.Exit(0)
		}

		// Check if market is open. If closed email current equity and sleep until the market reopens.
		if !profile.MarketOpen {
			log.Info("Market is closed")

			totalEquity, balanceChange := profile.GetEquityAndBalanceChange()
			mailer.Notify("Current equity: " + totalEquity + "\n" + "Today's change: " + balanceChange)

			sleep := profile.NextOpen.Sub(time.Now())
			log.Info("Sleeping for ", sleep)
			time.Sleep(sleep)

			continue
		}

		// do trades here

		log.Info("Sleeping for 5 minutes")
		time.Sleep(5 * time.Minute)
	}
}
