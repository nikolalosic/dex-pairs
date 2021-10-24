package contracts

import (
	"fmt"
	"math/big"

	"github.com/umbracle/go-web3"
	"github.com/umbracle/go-web3/contract"
	"github.com/umbracle/go-web3/jsonrpc"
)

// UniswapFactory is a solidity contract
type UniswapFactory struct {
	c *contract.Contract
}

// NewUniswapFactory creates a new instance of the contract at a specific address
func NewUniswapFactory(addr web3.Address, provider *jsonrpc.Client) *UniswapFactory {
	return &UniswapFactory{c: contract.NewContract(addr, UniswapFactoryAbi(), provider)}
}

// Contract returns the contract object
func (usf *UniswapFactory) Contract() *contract.Contract {
	return usf.c
}

// calls

func (usf *UniswapFactory) AllPairs(n int64, block ...web3.BlockNumber) (retval0 web3.Address, err error) {
	var out map[string]interface{}
	var ok bool

	out, err = usf.c.Call("allPairs", web3.EncodeBlock(block...), n)
	if err != nil {
		return
	}

	// decode outputs
	retval0, ok = out["0"].(web3.Address)
	if !ok {
		err = fmt.Errorf("failed to encode output at index 0")
		return
	}
	return
}

func (usf *UniswapFactory) AllPairsLength(block ...web3.BlockNumber) (retval0 *big.Int, err error) {
	var out map[string]interface{}
	var ok bool

	out, err = usf.c.Call("allPairsLength", web3.EncodeBlock(block...))
	if err != nil {
		return
	}

	// decode outputs
	retval0, ok = out["0"].(*big.Int)
	if !ok {
		err = fmt.Errorf("failed to encode output at index 0")
		return
	}
	return
}
