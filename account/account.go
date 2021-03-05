package account

import (
	"time"

	"github.com/alpacahq/alpaca-trade-api-go/alpaca"
	"github.com/alpacahq/alpaca-trade-api-go/common"
	"github.com/shopspring/decimal"
)

type Profile struct {
	AlpacaClient *alpaca.Client
	Account      *alpaca.Account
	MarketOpen   bool
	NextOpen     time.Time
	BuyingPower  decimal.Decimal
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

	profile := &Profile{
		AlpacaClient: alpacaClient,
		Account:      acct,
		MarketOpen:   clock.IsOpen,
		NextOpen:     clock.NextOpen,
		BuyingPower:  acct.BuyingPower,
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
