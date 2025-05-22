package dto

// MarketOverviewResponse represents overall market data
type MarketOverviewResponse struct {
	TotalValueLocked    string `json:"totalValueLocked"`
	TotalBorrowed       string `json:"totalBorrowed"`
	ActiveUsers         int64  `json:"activeUsers"`
	ActivePositions     int64  `json:"activePositions"`
	AverageLendingAPY   string `json:"averageLendingAPY"`
	AverageBorrowingAPY string `json:"averageBorrowingAPY"`
}

// TokenMetadata represents information about a token
type TokenMetadata struct {
	Address  string `json:"address"`
	Symbol   string `json:"symbol"`
	Name     string `json:"name"`
	Decimals int    `json:"decimals"`
	PriceUSD string `json:"priceUSD"`
	LogoURI  string `json:"logoURI,omitempty"`
}

// TokenMarketData represents market data for a specific token
type TokenMarketData struct {
	Token              TokenMetadata `json:"token"`
	TotalSupply        string        `json:"totalSupply"`
	TotalDeposited     string        `json:"totalDeposited"`
	TotalBorrowed      string        `json:"totalBorrowed"`
	AvailableLiquidity string        `json:"availableLiquidity"`
	LendingAPY         string        `json:"lendingAPY"`
	BorrowingAPY       string        `json:"borrowingAPY"`
	CollateralFactor   string        `json:"collateralFactor"`
}

// TokensMarketResponse represents market data for multiple tokens
type TokensMarketResponse struct {
	Tokens []TokenMarketData `json:"tokens"`
}
