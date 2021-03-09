package account

import (
	"math"
	"time"

	"github.com/alpacahq/alpaca-trade-api-go/alpaca"
	"github.com/alpacahq/alpaca-trade-api-go/common"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
)

type Profile struct {
	AlpacaClient  *alpaca.Client
	Account       *alpaca.Account
	MarketOpen    bool
	NextOpen      time.Time
	BuyingPower   float64
	MarketClosing bool
}

// InitializeClient initializes the client and checks if the market is open
func InitializeClient() (*Profile, error) {
	// paper-trading
	alpaca.SetBaseUrl("https://paper-api.alpaca.markets")
	// prod
	// alpaca.SetBaseUrl("https://api.alpaca.markets")

	// Set credentials as environment variables
	// APCA_API_KEY_ID=
	// APCA_API_SECRET_KEY=
	alpacaClient := alpaca.NewClient(common.Credentials())

	acct, err := alpacaClient.GetAccount()
	if err != nil {
		return nil, err
	}

	clock, err := alpacaClient.GetClock()
	if err != nil {
		return nil, err
	}

	nextClose := clock.NextClose
	isMarketClosing := checkMarketClosing(nextClose)

	buyingPower, _ := acct.BuyingPower.Float64()

	profile := &Profile{
		AlpacaClient:  alpacaClient,
		Account:       acct,
		MarketOpen:    clock.IsOpen,
		MarketClosing: isMarketClosing,
		NextOpen:      clock.NextOpen,
		BuyingPower:   buyingPower,
	}

	return profile, nil
}

// GetAccount returns the user's account details
func (p *Profile) GetAccount() *alpaca.Account {
	acct := p.Account

	return acct
}

// GetEquityAndBalanceChange returns the user's current equity and today's balance change
func (p *Profile) GetEquityAndBalanceChange() (string, string) {
	equity := p.Account.Equity
	balanceChange := equity.Sub(p.Account.LastEquity)

	return equity.String(), balanceChange.String()
}

func (p *Profile) CheckPositionChange(stock string) (float32, error) {

	barCount := 5
	bars, err := p.AlpacaClient.GetSymbolBars(stock, alpaca.ListBarParams{
		Timeframe: "day",
		Limit:     &barCount,
	})
	log.Info(bars[0])
	// asset, err := p.AlpacaClient.
	if err != nil {
		return 0.0, err
	}
	// log.Info(asset.)

	return 0.0, nil
}

// checkMarketClosing returns whether the market is closing within 15 minutes.
// Used to determine whether to sell all positions before market closes
func checkMarketClosing(timeToClose time.Time) bool {
	now := time.Now()
	sellTime := timeToClose.Add(time.Duration(-15) * time.Minute)

	if sellTime.Sub(now) <= 15*time.Minute {
		return true
	}
	return false
}

// func TradeHandler(trade stream.Trade) {
// 	log.Info("here")
// 	log.Info("trade", trade.Price)
// 	// return trade.Price
// }

func (a *Profile) PlaceOrder(symbol string, currentPrice float64) error {
	acct := a.AlpacaClient
	// buyingPower := a.BuyingPower

	qtyToBuy := math.Floor(700 / currentPrice)
	newBuyingPower := currentPrice * qtyToBuy
	log.Info("BUYING ", symbol, " ", qtyToBuy)
	// asset := symbol
	// log.Info(decimal.NewFromFloat(qtyToBuy))

	_, err := acct.PlaceOrder(alpaca.PlaceOrderRequest{
		AssetKey:    &symbol,
		Qty:         decimal.NewFromFloat(qtyToBuy),
		Type:        alpaca.Market,
		Side:        alpaca.Buy,
		TimeInForce: alpaca.IOC,
	})
	if err != nil {
		log.Error("Error buying ", &symbol)
		log.Error("ERROR: ", err)
		return err
	}

	a.BuyingPower = newBuyingPower
	// order.FilledQty
	return nil
	// to let the order go through. otherwise it happens too fast and the TrailingStop fails
	// time.Sleep(3 * time.Second)

	// stopLossPercent := 0.5
	// stopLossDecimal := decimal.NewFromFloat(stopLossPercent)
	// update buying power: will be approximated (depends on position buy in but this rough estimate will work until loop runs again and buying power is refreshed)
	// set sell limit here
	// _, err = alpaca.PlaceOrder(alpaca.PlaceOrderRequest{
	// 	AssetKey:     &symbol,
	// 	Qty:          order.Qty,
	// 	Side:         alpaca.Sell,
	// 	Type:         alpaca.TrailingStop,
	// 	TrailPercent: &stopLossDecimal,
	// 	TimeInForce:  alpaca.Day,
	// })
	// if err != nil {
	// 	log.Error("Error setting trailing stop ", err)
	// }

	// return nil
}

func SetNewStopTrailingPrice(name string, qty decimal.Decimal, stopLossId string) (string, error) {
	stopLossPercent := 0.5
	stopLossDecimal := decimal.NewFromFloat(stopLossPercent)

	if stopLossId != "" {
		err := alpaca.CancelOrder(stopLossId)
		if err != nil {
			log.Error("Couldn't cancel stop loss order with err: ", err)
			return "", err
		}
	}

	time.Sleep(3 * time.Second)

	order, err := alpaca.PlaceOrder(alpaca.PlaceOrderRequest{
		AssetKey:     &name,
		Qty:          qty,
		Side:         alpaca.Sell,
		Type:         alpaca.TrailingStop,
		TrailPercent: &stopLossDecimal,
		TimeInForce:  alpaca.Day,
	})
	if err != nil {
		log.Error("Error setting trailing stop ", err)
		return "", err
	}

	newStopLossId := order.ClientOrderID

	return newStopLossId, nil
}

// func setSellLimit()
