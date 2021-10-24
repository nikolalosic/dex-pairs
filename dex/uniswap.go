package dex

import (
	"fmt"
	"github.com/nikolalosic/dex_pairs/contracts"
	"github.com/umbracle/go-web3"
	"github.com/umbracle/go-web3/contract/builtin/erc20"
	"github.com/umbracle/go-web3/jsonrpc"
	"log"
	"math/big"
	"strings"
)

type Uniswap struct {
	factory *contracts.UniswapFactory
	client  *jsonrpc.Client
	chainId int
}

func NewUniswap(factoryAddress web3.Address, chainId int, nodeUrl string) (*Uniswap, error) {
	client, err := jsonrpc.NewClient(nodeUrl)
	if err != nil {
		log.Printf("Error openning rpc client")
		return nil, err
	}
	return &Uniswap{
		factory: contracts.NewUniswapFactory(factoryAddress, client),
		client:  client,
		chainId: chainId,
	}, nil
}

func (us *Uniswap) GetPair(n int64) (*Pair, error) {
	log.Printf("Getting Uniswap pair. n=%d", n)
	pairAddress, err := us.factory.AllPairs(n, web3.Latest)
	if err != nil {
		return nil, err
	}
	pairContract := contracts.NewUniswapPair(pairAddress, us.client)
	pairSymbol, _ := pairContract.Symbol(web3.Latest)
	pairName, _ := pairContract.Name(web3.Latest)
	pairDecimals, _ := pairContract.Decimals(web3.Latest)

	token0, _ := pairContract.Token0(web3.Latest)
	token1, _ := pairContract.Token1(web3.Latest)

	token0Contract := erc20.NewERC20(token0, us.client)
	token1Contract := erc20.NewERC20(token1, us.client)

	token0Symbol := "UNK"
	token1Symbol := "UNK"
	if token0Contract.Contract().Addr() != zeroAddress {
		token0Symbol, err = token0Contract.Symbol(web3.Latest)
		if err != nil || !allowedRegex.MatchString(token0Symbol) {
			token0Symbol = "UNK"
		}
	}
	if token1Contract.Contract().Addr() != zeroAddress {
		token1Symbol, err = token1Contract.Symbol(web3.Latest)
		if err != nil || !allowedRegex.MatchString(token1Symbol) {
			token1Symbol = "UNK"
		}
	}
	if len(token0Symbol) > 13 {
		token0Symbol = token0Symbol[:13]
	}
	if len(token1Symbol) > 13 {
		token0Symbol = token1Symbol[:13]
	}

	pair := Pair{
		Token0:   strings.ToLower(token0.String()),
		Token1:   strings.ToLower(token1.String()),
		Name:     fmt.Sprintf("%s - %s/%s", pairName, token0Symbol, token1Symbol),
		Address:  strings.ToLower(pairAddress.String()),
		Symbol:   pairSymbol,
		Decimals: int(pairDecimals),
		ChainId:  us.chainId,
	}
	return &pair, nil
}

func (us *Uniswap) GetPairNumber() (*big.Int, error) {
	return us.factory.AllPairsLength(web3.Latest)
}
