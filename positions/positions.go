package positions

type PositionTracking struct {
	Name                 string  // don't update
	Owned                bool    // owned
	CurrentPercentChange float64 // update when checking price
	OverallPercentChange float64 // update when checking price
	PriceBought          float64 // don't update
	LastPrice            float64 // update current price to lastPrice for future loop iterations
	// CurrentPrice         float64 // delete?
}

// UpdatePosition updates the position with current information used for tracking buy/sell/hold
func (p *PositionTracking) UpdatePosition(currentPrice float64) {
	var overallPercentchange float64
	lastPrice := p.LastPrice
	entryPrice := p.PriceBought

	// to account for market opening
	if lastPrice == 0 {
		lastPrice = currentPrice
	}

	currentPercentChange := getPositionPercentChange(currentPrice, lastPrice)
	if p.Owned {
		overallPercentchange = getPositionPercentChange(entryPrice, lastPrice)
	}
	p.LastPrice = currentPrice // set last price to current price for the next loop iteration
	p.CurrentPercentChange = currentPercentChange
	p.OverallPercentChange = overallPercentchange
}

// getPositionPercentChange tracks the percent change between each loop
func getPositionPercentChange(currentPrice, lastPrice float64) float64 {
	change := currentPrice - lastPrice
	percentChange := (change / lastPrice) * 100

	return percentChange
}

// DELETE: Can use to test in agent loop
// loop through current positions to determine hold/sell
// for _, position := range positions {
// log.Info(position)
// name := "UBER"
// entryPrice := 55.00
// currentPrice := 55.16 // position.CurrentPrice

// pos := memPos[name]
// pos.PriceBought = entryPrice
// pos.LastPrice = entryPrice
// pos.UpdatePosition(currentPrice)
// log.Info("SHOULD BE 55.16 ---- ", pos.LastPrice)
// log.Info("PERCENT CHANGE", pos.CurrentPercentChange)

// pos.UpdatePosition(56)
// log.Info("SHOULD BE 56 ----", pos.LastPrice)
// log.Info("PERCENT CHANGE", pos.CurrentPercentChange)

// pos.UpdatePosition(55.67)
// log.Info("SHOULD BE 56 ----", pos.LastPrice)
// log.Info("PERCENT CHANGE", pos.CurrentPercentChange)
// get current price and do math. if total profit >= 1.5% sell all unless price rose >= 0.5% over past 5 mins
// }
