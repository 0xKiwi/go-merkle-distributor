package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/common"
)


var zeroAddr = common.HexToAddress("0x0000000000000000000000000000000000000000")
var safe1Address = common.HexToAddress("0x1Aa61c196E76805fcBe394eA00e4fFCEd24FC469")

var (
	jsonFile = flag.String("json-file", "", "JSON file of addresses to balances in wei")
	outputFile = flag.String("output-file", "addr-to-claim.json", "JSON file of addresses to claiming info")
)

func main() {
	flag.Parse()

	if *jsonFile == "" {
		log.Fatal("Expected --json-file to contain a file path")
	}
	
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
	balArray, err := ArrayFromAddrBalMap(stringJson)
	if err != nil {
		log.Fatalf("Could not get array from string map json")
	}
	_, addrToClaim, err := CreateDistributionTree(balArray)
	if err != nil {
		log.Fatalf("Could not create distribution tree: %v", err)
	}
	if _, err := createFile(*outputFile, addrToClaim); err != nil {
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
