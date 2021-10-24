package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/nikolalosic/dex_pairs/dex"
	"github.com/umbracle/go-web3"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"sync"
	"time"
)

var uniswapV1FactoryAddress = web3.HexToAddress("0xc0a47dfe034b400b47bdad5fecda2621de6c4d95")
var uniswapV2FactoryAddress = web3.HexToAddress("0x5c69bee701ef814a2b6a3edd4b1652cb9cc5aa6f")
var uniswapV3FactoryAddress = web3.HexToAddress("0x1f98431c8ad98523631ae4a59f267346ea31f984")
var pancakeSwapV1FactoryAddress = web3.HexToAddress("0xbcfccbde45ce874adcb698cc183debcf17952812")
var pancakeSwapV2FactoryAddress = web3.HexToAddress("0xca143ce32fe78f1f7019d7d551a6402fc5350c73")

var factoryContracts = map[string]map[int]map[int]web3.Address{
	"pancakeswap": {
		56: {
			1: pancakeSwapV1FactoryAddress,
			2: pancakeSwapV2FactoryAddress,
		},
	},
	"uniswap": {
		1: {
			1: uniswapV1FactoryAddress,
			2: uniswapV2FactoryAddress,
			3: uniswapV3FactoryAddress,
		},
	},
}

type job struct {
	start       int
	end         int
	dexExchange string
}
type result struct {
	pairs []dex.Pair
	err   error
}

type fileTemplate struct {
	Name      string     `json:"name"`
	Timestamp time.Time  `json:"timestamp"`
	Version   version    `json:"version"`
	Keywords  []string   `json:"keywords"`
	Tokens    []dex.Pair `json:"tokens"`
}

type version struct {
	Major int `json:"major"`
	Minor int `json:"minor"`
	Patch int `json:"patch"`
}

var counter = 0
var m = sync.RWMutex{}

func getDataFromFile(fileName string) (*fileTemplate, error) {
	log.Printf("Reading data from file %s", fileName)
	ft := fileTemplate{}
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		return &ft, nil
	}
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Printf("Error opening file %s", fileName)
		return nil, err
	}
	err = json.Unmarshal(content, &ft)
	if err != nil {
		log.Printf("Error unmarshaling json %s", fileName)
		return nil, err
	}
	return &ft, nil
}

func getPairs(start int, end int, dexExchange string, dexVersion int, chainId int) []dex.Pair {
	log.Printf("Getting dex pairs. start=%d, end=%d", start, end)
	d, err := getDex(dexExchange, dexVersion, chainId)
	if err != nil {
		panic("Cannot Get DEX")
	}

	var res []dex.Pair
	for i := start; i < end; i++ {
		pair, err := d.GetPair(int64(i))
		m.Lock()
		counter++
		m.Unlock()
		if err != nil {
			log.Printf("Error getting pair n=%d. Error=%s", i, err.Error())
			continue
		}
		m.RLock()
		log.Printf("Fetched pair %d, total fetched=%d", i, counter)
		m.RUnlock()
		res = append(res, *pair)
	}
	return res
}

func saveToFile(fileData *fileTemplate, fileName string) error {
	log.Printf("Saving data to file %s", fileName)
	data, err := json.Marshal(fileData)
	if err != nil {
		log.Printf("Error marshalling file template")
		return err
	}
	err = ioutil.WriteFile(fileName, data, 0755)
	if err != nil {
		log.Printf("Error writing to file %s", fileName)
		return err
	}
	return nil
}

func getDex(dexId string, dexVersion int, chainId int) (dex.DexExchange, error) {
	if dexId == "uniswap" {
		return dex.NewUniswap(factoryContracts[dexId][chainId][dexVersion], chainId, os.Getenv("NODE_URL"))
	} else if dexId == "pancakeswap" {
		return dex.NewPancakeSwap(factoryContracts[dexId][chainId][dexVersion], chainId, os.Getenv("NODE_URL"))
	}
	return nil, errors.New("no appropriate DEX found")
}

func main() {

	// flags declaration using flag package
	var inputFile, outputFile, dexExchange string
	var cores, chainId, dexVersion int
	flag.StringVar(&inputFile, "input-file", "dex_pairs.json", "Specify input file.")
	flag.StringVar(&outputFile, "output-file", "dex_pairs.json", "Specify output file.")
	flag.IntVar(&cores, "cores", runtime.NumCPU()/2, "Specify number of cores to use. Default is runtime.NumCPU()/2.")
	flag.StringVar(&dexExchange, "dex-exchange", "uniswap", "Specify from which DEX exchange to get pairs.")
	flag.IntVar(&chainId, "chain-id", 1, "Specify chain id.")
	flag.IntVar(&dexVersion, "dex-version", 2, "Specify from which DEX exchange version to get pairs.")
	flag.Parse()

	exchange, err := getDex(dexExchange, dexVersion, chainId)
	if err != nil {
		panic("Cannot Get DEX")
	}
	pn, err := exchange.GetPairNumber()
	if err != nil {
		log.Printf("Error getting all pairs length")
		return
	}
	pairCount := int(pn.Int64())
	data, err := getDataFromFile(inputFile)
	if err != nil {
		log.Printf("Error reading data from input file")
	}
	n := len(data.Tokens) - 1
	if n < 0 {
		n = 0
	}
	step := 300
	if step > pairCount {
		step = pairCount
	}
	jobCount := (pairCount - n) / step
	fmt.Println(fmt.Sprintf("pairCount=%d, n=%d, step=%d", pairCount, n, step))
	if (pairCount-n)%step != 0 {
		jobCount++
	}
	if jobCount == 0 {
		return
	}
	if jobCount < cores {
		cores = jobCount
	}

	jobs := make(chan job, cores)
	results := make(chan result)
	runtime.GOMAXPROCS(cores)

	for i := 0; i < cores; i++ {
		go func(js <-chan job, rs chan<- result) {
			for j := range js {
				rs <- result{err: nil, pairs: getPairs(j.start, j.end, dexExchange, dexVersion, chainId)}
			}
		}(jobs, results)
	}
	for i := n; i < pairCount; i += step {
		if i+step > pairCount {
			jobs <- job{start: i, end: pairCount, dexExchange: dexExchange}
		} else {
			jobs <- job{start: i, end: i + step, dexExchange: dexExchange}
		}
	}

	close(jobs)
	for i := 0; i < jobCount; i++ {
		r := <-results
		data.Tokens = append(data.Tokens, r.pairs...)
		if r.err != nil {
			log.Printf(r.err.Error())
		}
	}

	log.Printf("Completed getting pairs")
	err = saveToFile(data, outputFile)
	if err != nil {
		return
	}

}
