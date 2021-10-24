package dex

import (
	"github.com/umbracle/go-web3"
	"math/big"
	"regexp"
)

var allowedRegex = regexp.MustCompile("^[ \\w.'+\\-$%/]+$")
var zeroAddress = web3.HexToAddress("0x0000000000000000000000000000000000000000")

type DexExchange interface {
	GetPair(n int64) (*Pair, error)
	GetPairNumber() (*big.Int, error)
}
