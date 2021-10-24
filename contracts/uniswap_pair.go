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

// UniswapPair is a solidity contract
type UniswapPair struct {
	c *contract.Contract
}

// NewUniswapPair creates a new instance of the contract at a specific address
func NewUniswapPair(addr web3.Address, provider *jsonrpc.Client) *UniswapPair {
	return &UniswapPair{c: contract.NewContract(addr, UniswapPairAbi(), provider)}
}

// Contract returns the contract object
func (up *UniswapPair) Contract() *contract.Contract {
	return up.c
}

// calls
// Token0 calls the token0 method in the solidity contract
func (up *UniswapPair) Token0(block ...web3.BlockNumber) (retval0 web3.Address, err error) {
	var out map[string]interface{}
	var ok bool

	out, err = up.c.Call("token0", web3.EncodeBlock(block...))
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

// Token1 calls the token1 method in the solidity contract
func (up *UniswapPair) Token1(block ...web3.BlockNumber) (retval0 web3.Address, err error) {
	var out map[string]interface{}
	var ok bool

	out, err = up.c.Call("token1", web3.EncodeBlock(block...))
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

// Allowance calls the allowance method in the solidity contract
func (up *UniswapPair) Allowance(
	owner web3.Address, spender web3.Address, block ...web3.BlockNumber,
) (retval0 *big.Int, err error) {
	var out map[string]interface{}
	var ok bool

	out, err = up.c.Call("allowance", web3.EncodeBlock(block...), owner, spender)
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

// BalanceOf calls the balanceOf method in the solidity contract
func (up *UniswapPair) BalanceOf(owner web3.Address, block ...web3.BlockNumber) (retval0 *big.Int, err error) {
	var out map[string]interface{}
	var ok bool

	out, err = up.c.Call("balanceOf", web3.EncodeBlock(block...), owner)
	if err != nil {
		return
	}

	// decode outputs
	retval0, ok = out["balance"].(*big.Int)
	if !ok {
		err = fmt.Errorf("failed to encode output at index 0")
		return
	}

	return
}

func (up *UniswapPair) GetReserves(
	owner web3.Address, block ...web3.BlockNumber,
) (retval0 *big.Int, retval1 *big.Int, err error) {
	var out map[string]interface{}
	var ok bool

	out, err = up.c.Call("getReserves", web3.EncodeBlock(block...), owner)
	if err != nil {
		return
	}

	// decode outputs
	retval0, ok = out["_reserve0"].(*big.Int)
	retval1, ok = out["_reserve1"].(*big.Int)

	if !ok {
		err = fmt.Errorf("failed to encode output at index 0")
		return
	}

	return
}

// Decimals calls the decimals method in the solidity contract
func (up *UniswapPair) Decimals(block ...web3.BlockNumber) (retval0 uint8, err error) {
	var out map[string]interface{}
	var ok bool

	out, err = up.c.Call("decimals", web3.EncodeBlock(block...))
	if err != nil {
		return
	}

	// decode outputs
	retval0, ok = out["0"].(uint8)
	if !ok {
		err = fmt.Errorf("failed to encode output at index 0")
		return
	}

	return
}

// Name calls the name method in the solidity contract
func (up *UniswapPair) Name(block ...web3.BlockNumber) (retval0 string, err error) {
	var out map[string]interface{}
	var ok bool

	out, err = up.c.Call("name", web3.EncodeBlock(block...))
	if err != nil {
		return
	}

	// decode outputs
	retval0, ok = out["0"].(string)
	if !ok {
		err = fmt.Errorf("failed to encode output at index 0")
		return
	}

	return
}

// Symbol calls the symbol method in the solidity contract
func (up *UniswapPair) Symbol(block ...web3.BlockNumber) (retval0 string, err error) {
	var out map[string]interface{}
	var ok bool

	out, err = up.c.Call("symbol", web3.EncodeBlock(block...))
	if err != nil {
		return
	}

	// decode outputs
	retval0, ok = out["0"].(string)
	if !ok {
		err = fmt.Errorf("failed to encode output at index 0")
		return
	}

	return
}

// TotalSupply calls the totalSupply method in the solidity contract
func (up *UniswapPair) TotalSupply(block ...web3.BlockNumber) (retval0 *big.Int, err error) {
	var out map[string]interface{}
	var ok bool

	out, err = up.c.Call("totalSupply", web3.EncodeBlock(block...))
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

// txns

// Approve sends a approve transaction in the solidity contract
func (up *UniswapPair) Approve(spender web3.Address, value *big.Int) *contract.Txn {
	return up.c.Txn("approve", spender, value)
}

// Transfer sends a transfer transaction in the solidity contract
func (up *UniswapPair) Transfer(to web3.Address, value *big.Int) *contract.Txn {
	return up.c.Txn("transfer", to, value)
}

// TransferFrom sends a transferFrom transaction in the solidity contract
func (up *UniswapPair) TransferFrom(from web3.Address, to web3.Address, value *big.Int) *contract.Txn {
	return up.c.Txn("transferFrom", from, to, value)
}

// events

func (up *UniswapPair) ApprovalEventSig() web3.Hash {
	return up.c.ABI().Events["Approval"].ID()
}

func (up *UniswapPair) TransferEventSig() web3.Hash {
	return up.c.ABI().Events["Transfer"].ID()
}
