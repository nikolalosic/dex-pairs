package contracts

import (
	"fmt"
	"math/big"

	"github.com/umbracle/go-web3"
	"github.com/umbracle/go-web3/contract"
	"github.com/umbracle/go-web3/jsonrpc"
)

var (
	_ = big.NewInt
)

var (
	_ = big.NewInt
)

// PancakeFactory is a solidity contract
type PancakeFactory struct {
	c *contract.Contract
}

// NewPancakeFactory creates a new instance of the contract at a specific address
func NewPancakeFactory(addr web3.Address, provider *jsonrpc.Client) *PancakeFactory {
	return &PancakeFactory{c: contract.NewContract(addr, PancakeFactoryAbi(), provider)}
}

// Contract returns the contract object
func (pf *PancakeFactory) Contract() *contract.Contract {
	return pf.c
}

// calls

func (pf *PancakeFactory) AllPairs(n int64, block ...web3.BlockNumber) (retval0 web3.Address, err error) {
	var out map[string]interface{}
	var ok bool

	out, err = pf.c.Call("allPairs", web3.EncodeBlock(block...), n)
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

func (pf *PancakeFactory) AllPairsLength(block ...web3.BlockNumber) (retval0 *big.Int, err error) {
	var out map[string]interface{}
	var ok bool

	out, err = pf.c.Call("allPairsLength", web3.EncodeBlock(block...))
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
