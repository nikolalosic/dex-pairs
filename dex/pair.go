package dex

type Pair struct {
	Token0   string `json:"token0"`
	Token1   string `json:"token1"`
	Name     string `json:"name"`
	Address  string `json:"address"`
	Symbol   string `json:"symbol"`
	Decimals int    `json:"decimals"`
	ChainId  int    `json:"chainId"`
}
