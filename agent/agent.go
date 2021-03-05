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

		// Notifies user agent is down if a panic occurs.
		defer recovery()

		// Initialize account and get current account information/balance
		profile, err := account.InitializeClient()
		if err != nil {
			log.Error("Error initializing client: ", err)
			// email notifying agent is down
			mailer.Notify("Could not initialize client: " + err.Error())

			os.Exit(1)
		}

		// Get user's account information
		acct := profile.GetAccount()
		if acct.TradingBlocked || acct.AccountBlocked {
			log.Error("Account is blocked")
			// email notifying agent is down.
			mailer.Notify("Account is blocked. Trading Blocked: " + strconv.FormatBool(acct.TradingBlocked) + " Account Blocked: " + strconv.FormatBool(acct.AccountBlocked))

			os.Exit(1)
		}

		// Check if market is open. If closed email current equity and balance change and sleep until the market reopens.
		if !profile.MarketOpen {
			log.Info("Market is closed")

			totalEquity, balanceChange := profile.GetEquityAndBalanceChange()
			mailer.Notify("Current equity: " + totalEquity + "\n" + "Today's change: " + balanceChange)

			sleep := profile.NextOpen.Sub(time.Now())
			log.Info("Sleeping for ", sleep)
			time.Sleep(sleep)

			continue
		}

		// Get all current positions
		positions, err := profile.AlpacaClient.ListPositions()
		if err != nil {
			mailer.Notify("Couldn't list positions to error: " + err.Error())
			continue
		}

		for _, position := range positions {
			log.Info(position)
		}

		log.Info("Sleeping for 5 minutes")
		time.Sleep(5 * time.Minute)
	}
}
