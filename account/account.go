package account

import (
	"time"

	"github.com/alpacahq/alpaca-trade-api-go/alpaca"
	"github.com/alpacahq/alpaca-trade-api-go/common"
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
	log.Info(bars)
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
