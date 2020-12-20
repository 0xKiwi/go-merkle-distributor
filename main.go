package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"os"
	"path/filepath"

	"github.com/ethereum/go-ethereum/common"
)


var zeroAddr = common.HexToAddress("0x0000000000000000000000000000000000000000")
var safe1Address = common.HexToAddress("0x1Aa61c196E76805fcBe394eA00e4fFCEd24FC469")

var savedDataFileName = filepath.Join("output", "saved-balances-and-supply.json")
var addrToClaimFileName = filepath.Join("output", "addr-to-claim.json")

var (
	jsonFile = flag.String("json-file", "", "JSON file of addresses to balances in wei")
)

func main() {
	flag.Parse()
	
	fullPath, err := expandPath(*jsonFile)
	if err != nil {
		log.Fatalf("Could not expand path: %v", err)
	}
	jsonBytes, err := ioutil.ReadFile(fullPath)
	if err != nil {
		log.Fatalf("Could not read file: %v", err)
	}
	var stringJson map[string]string
	if err := json.Unmarshal(jsonBytes, &stringJson); err != nil {
		log.Fatalf("Could not unmarshal json: %v", err)
	}

	i := 0
	balArray := make([]*TokenHolder, len(stringJson))
	for addr, bal := range stringJson {
		bigInt, ok := big.NewInt(0).SetString(bal, 10)
		if !ok {
			log.Fatalf("could not cast %s to big int", bal)
		}
		balArray[i] = &TokenHolder{
			addr: common.HexToAddress(addr),
			balance: bigInt,
		}
		i++
	}
	_, addrToClaim, err := CreateDistributionTree(balArray)
	if err != nil {
		log.Fatalf("Could not create distribution tree: %v", err)
	}
	if _, err := createFile(addrToClaimFileName, addrToClaim); err != nil {
		log.Fatalf("Could not create file: %v", err)
	}
}

func unmarshalJSON(jsonMap map[string]string) (map[common.Address]*big.Int, error) {
	balMap := make(map[common.Address]*big.Int, len(jsonMap) - 1)
	for k, v := range jsonMap {
		bigInt, ok := big.NewInt(0).SetString(v, 10)
		if !ok {
			return nil, fmt.Errorf("could not cast %s to big int", v)
		}
		balMap[common.HexToAddress(k)] = bigInt
	}

	return balMap, nil
}

func createFile(filename string, contents interface{}) (*os.File, error) {
	file, err := os.Create(filename)
	if err != nil {
		return nil, fmt.Errorf("could not create file: %v", err)
	}
	totalBytes, err := json.Marshal(contents)
	if err != nil {
		return nil, fmt.Errorf("could not marshal data: %v", err)
	}
	if _, err := file.Write(totalBytes); err != nil {
		return nil, fmt.Errorf("could not write data: %v", err)
	}
	return file, nil
}
