package account

import (
	"time"

	"github.com/alpacahq/alpaca-trade-api-go/alpaca"
	"github.com/alpacahq/alpaca-trade-api-go/common"
)

type Client struct {
	AlpacaClient *alpaca.Client
	Account      *alpaca.Account
	MarketOpen   bool
	NextOpen     time.Time
}

// InitializeClient initializes the client and checks if the market is open
func InitializeClient() (*Client, error) {
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

	Client := &Client{
		AlpacaClient: alpacaClient,
		Account:      acct,
		MarketOpen:   clock.IsOpen,
		NextOpen:     clock.NextOpen,
	}

	return Client, nil
}

func (c *Client) GetAccount() *alpaca.Account {
	acct := c.Account

	return acct
}

func (c *Client) GetEquity() (string, string) {
	equity := c.Account.Equity
	balanceChange := equity.Sub(c.Account.LastEquity)

	return equity.String(), balanceChange.String()
}
