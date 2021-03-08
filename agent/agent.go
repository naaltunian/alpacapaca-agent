package agent

import (
	"os"
	"strconv"
	"time"

	"github.com/naaltunian/paca-agent/account"
	"github.com/naaltunian/paca-agent/mailer"
	"github.com/naaltunian/paca-agent/positions"

	log "github.com/sirupsen/logrus"
)

func Start() {
	// add microsoft, amazon, netflix, walmart, target, etc
	stockToWatch := []string{"UBER", "AAPL", "TSLA"}
	// TODO: move to db once tested. keeping track of stock changes in memory
	memPos := make(map[string]*positions.PositionTracking)
	memPos["UBER"] = &positions.PositionTracking{Name: "UBER", Owned: false}
	memPos["TSLA"] = &positions.PositionTracking{Name: "TSLA", Owned: false}
	memPos["AAPL"] = &positions.PositionTracking{Name: "AAPL", Owned: false}
	for {

		// Notifies user agent is down if a panic occurs.
		defer recovery()

		// Initialize account and get current account information/balance
		profile, err := account.InitializeClient()
		if err != nil {
			log.Error("Error initializing client: ", err)
			// email notifying agent is down
			mailer.Notify("Error", "Could not initialize client: "+err.Error())

			os.Exit(1)
		}

		// Get user's account information
		acct := profile.GetAccount()
		if acct.TradingBlocked || acct.AccountBlocked {
			log.Error("Account is blocked")
			// email notifying agent is down.
			mailer.Notify("Error", "Account is blocked. Trading Blocked: "+strconv.FormatBool(acct.TradingBlocked)+" Account Blocked: "+strconv.FormatBool(acct.AccountBlocked))

			os.Exit(1)
		}

		// Check if market is closed or closing in 15 min. If closed email current equity and balance change and sleep until the market reopens.
		// if profile.MarketClosing || !profile.MarketOpen {
		// 	log.Info("Market is closing. Selling all open positions")

		// 	totalEquity, balanceChange := profile.GetEquityAndBalanceChange()
		// 	mailer.Notify("Days End", "Current equity: "+totalEquity+"\n"+"Today's change: "+balanceChange)

		// 	sleep := profile.NextOpen.Sub(time.Now())
		// 	log.Info("Sleeping for ", sleep)
		// 	time.Sleep(sleep)

		// 	continue
		// }

		// Get all current positions
		positions, err := profile.AlpacaClient.ListPositions()
		if err != nil {
			mailer.Notify("Error", "Couldn't list positions to error: "+err.Error())
			continue
		}

		// memPos["UBER"].Owned = true

		// TODO: cleanup loop
		// loop through current positions to determine hold/sell
		for _, position := range positions {
			log.Info("Starting hold/sell run")
			name := position.Symbol
			currentPrice, _ := position.CurrentPrice.Float64()
			log.Info(name, currentPrice)
			memPos[position.Symbol].UpdatePosition(currentPrice)

			// get current price and do math. if total profit >= 1.5% sell all unless price rose >= 0.5% over past 5 mins
			if memPos[position.Symbol].CurrentPercentChange <= 0.30 {
				log.Info("SELL THIS POSITION")
				err := profile.AlpacaClient.ClosePosition(name)
				if err != nil {
					log.Error("Could not sell position ", err)
					continue
				}
			} else if memPos[position.Symbol].OverallPercentChange >= 1.5 {
				log.Info("SELL THIS POSITION")
				err := profile.AlpacaClient.ClosePosition(name)
				if err != nil {
					log.Error("Could not sell position ", err)
					continue
				}
			}
		}

		// TODO: cleanup loop. If holding position don't buy more
		// loop through positions to buy
		if profile.BuyingPower >= 5000 {
			for _, stock := range stockToWatch {
				// if holding stock continue. Don't buy more
				if memPos[stock].Owned {
					log.Info("Holding stock ", stock+". Continuing")
					continue
				}

				name := memPos[stock].Name

				// now := time.Now()
				// startTime := now.Add(time.Duration(-5) * time.Minute)
				// quotes := profile.AlpacaClient.GetQuotes(name, startTime, time.Now(), 5)
				quotes, err := profile.AlpacaClient.GetLastTrade(name)
				if err != nil {
					log.Error(err)
					continue
				}
				log.Info(quotes.Last.Price)
				memPos[stock].UpdatePosition(float64(quotes.Last.Price))

				// get lowest? mid?
				// for quote := range quotes {
				// 	log.Info(quote.Quote.AskPrice)
				// }
				// if err := stream.SubscribeTrades(account.TradeHandler, "AAPL"); err != nil {
				// 	panic(err)
				// }
				// quote, err := profile.AlpacaClient.ListAssets()
				// log.Info("QUOTE ", quote.Last.AskPrice)
				if err != nil {
					log.Error(err)
					continue
				}
				// pos.UpdatePosition(pos.CurrentPrice)

				// percentChange, err := profile.CheckPositionChange(stock)
				// if err != nil {
				// 	log.Error("Couldn't get balance change ", err)
				// 	continue
				// }
				if memPos[stock].CurrentPercentChange >= 0.2 {
					// buy and set sell limit stop loss -0.05%
					log.Info("BUYING", stock)
				}
			}
		}

		log.Infof("%+v", memPos["AAPL"])
		log.Infof("%+v", memPos["TSLA"])
		log.Infof("%+v", memPos["UBER"])

		log.Info("Sleeping for 5 minutes")
		time.Sleep(5 * time.Minute)
	}
}
