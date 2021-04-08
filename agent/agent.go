package agent

import (
	"os"
	"strconv"
	"time"

	"github.com/alpacahq/alpaca-trade-api-go/alpaca"
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
		if profile.MarketClosing || !profile.MarketOpen {
			log.Info("Market is closing. Selling all open positions")

			profile.AlpacaClient.CloseAllPositions()

			// give the account 3 seconds to settle before re-retrieving the account info
			time.Sleep(3 * time.Second)

			profile, err := account.InitializeClient()
			// continue until client connection reestablished
			if err != nil {
				log.Error("reininitizling client: ", err)
				continue
			}

			totalEquity, balanceChange := profile.GetEquityAndBalanceChange()
			mailer.Notify("Days End", "Current equity: "+totalEquity+"\n"+"Today's change: "+balanceChange)

			sleep := profile.NextOpen.Sub(time.Now())
			log.Info("Sleeping for ", sleep)
			time.Sleep(sleep)

			continue
		}

		// Get all current positions
		positions, err := profile.AlpacaClient.ListPositions()
		if err != nil {
			mailer.Notify("Error", "Couldn't list positions to error: "+err.Error())
			continue
		}

		// deprecate once in db
		// used keep track of what is held in memory and mark everything as not owned. Mark as true once validated after getting current positions. This helps keep track of stop orders being filled in between loops
		for _, pos := range memPos {
			pos.Owned = false
		}

		log.Info("before looping through positions")

		// TODO: cleanup loop
		// loop through current positions to determine hold/sell
		for _, position := range positions {
			memPos[position.Symbol].Owned = true
			name := position.Symbol
			// qty := position.Qty
			currentPrice, _ := position.CurrentPrice.Float64()
			// entryPrice, _ := position.EntryPrice.Float64()
			// stopLossId := memPos[position.Symbol].StopLossOrderId

			// if currentPrice > entryPrice {
			// 	stopLossId, err := account.SetNewStopTrailingPrice(name, qty, stopLossId)
			// 	if err != nil {
			// 		log.Error("error setting trailingprice ", err)
			// 		continue
			// 	}
			// 	memPos[position.Symbol].StopLossOrderId = stopLossId
			// }

			memPos[position.Symbol].UpdatePosition(currentPrice)
			// delete below most likely
			// get current price and do math. if total profit >= 1.5% sell all unless price rose >= 0.5% over past 5 mins
			if memPos[position.Symbol].CurrentPercentChange <= 0.30 {
				log.Info("SELL THIS POSITION STOP LOSSES ", position.Symbol)
				// err := profile.AlpacaClient.ClosePosition(name)
				// if err != nil {
				// 	log.Error("Could not sell position ", err)
				// 	continue
				// }
			} else if memPos[position.Symbol].OverallPercentChange >= 1.5 {
				log.Info("SELL THIS POSITION MAKE PROFIT ", position.Symbol)
				err := profile.AlpacaClient.ClosePosition(name)
				if err != nil {
					log.Error("Could not sell position ", err)
					continue
				}
				memPos[position.Symbol].Owned = false
			}
		}

		log.Info("before looping through stock to buy")

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

				quotes, err := profile.AlpacaClient.GetLastTrade(name)
				if err != nil {
					log.Error(err)
					continue
				}

				log.Info("after getting last trade")

				memPos[stock].UpdatePosition(float64(quotes.Last.Price))

				log.Info("before placing order")
				if memPos[stock].CurrentPercentChange >= 0.2 {
					// buy and set sell limit stop loss -0.05%
					log.Info("BUYING ", stock)
					err := profile.PlaceOrder(stock, float64(quotes.Last.Price))
					if err != nil {
						log.Error("Error buying stock, ", stock, " with error: ", err)
					}

					time.Sleep(3 * time.Second)

					// get accurate position data
					position, err := alpaca.GetPosition(stock)
					if err != nil {
						log.Error("Error getting position ", err)
						continue
					}

					stopId, err := account.SetNewStopTrailingPrice(stock, position.Qty, "")
					if err != nil {
						log.Error("Error setting new stop loss: ", err)
						continue
					}
					memPos[stock].Owned = true
					memPos[stock].StopLossOrderId = stopId
				}
			}
		}

		log.Info("after buy loop")

		log.Infof("%+v", memPos["AAPL"])
		log.Infof("%+v", memPos["TSLA"])
		log.Infof("%+v", memPos["UBER"])
		log.Info("BUYING POWER ", profile.BuyingPower)

		log.Info("Sleeping for 5 minutes")
		log.Info("--------------------------------------------------------")
		time.Sleep(5 * time.Minute)
	}
}
