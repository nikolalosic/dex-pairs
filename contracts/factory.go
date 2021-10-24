package contracts

import (
	"github.com/umbracle/go-web3"
	"math/big"
)

type Factory interface {
	AllPairs(n big.Int, block ...web3.BlockNumber) (retval0 *string, err error)
	AllPairsLength(block ...web3.BlockNumber) (retval0 *big.Int, err error)
}
