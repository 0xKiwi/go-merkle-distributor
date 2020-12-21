package web3

import (
	"bytes"
	"context"
	"fmt"
	"math/big"
	"os"
	"path/filepath"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

var totalSupplyKey = "totalSupply"
var erc20ABIFileName = filepath.Join("..", "abi", "IERC20.json")

var zeroAddr = common.HexToAddress("0x0000000000000000000000000000000000000000")
var deadAddr = common.HexToAddress("0x000000000000000000000000000000000000dead")

// BalancesAndSupplyAtBlock takes in a go-ethereum client, token address and traces 
// through transfer logs from the startBlock, to the endBlock.
func BalancesAndSupplyAtBlock(client *ethclient.Client, tokenAddress common.Address, startBlock int64, endBlock int64) (map[common.Address]*big.Int, *big.Int, error) {
	var totalLogs []types.Log
	for i := startBlock; i < endBlock; i+=2000 {
		epochEnd := i+1999
		if epochEnd > endBlock {
			epochEnd = endBlock
		}
		fmt.Printf("Reading blocks %d - %d\n", i, endBlock)
		query := ethereum.FilterQuery{
			FromBlock: big.NewInt(i),
			ToBlock:   big.NewInt(epochEnd),
			Addresses: []common.Address{
				tokenAddress,
			},
		}
		logs, err := client.FilterLogs(context.Background(), query)
		if err != nil {
			return nil, nil, err
		}
		totalLogs = append(totalLogs, logs...)
	}
	abiFile, err := os.Open(erc20ABIFileName)
	if err != nil {
		return nil, nil, err
	}
	contractAbi, err := abi.JSON(abiFile)
	if err != nil {
		return nil, nil, err
	}

	type logTransfer struct {
		From   common.Address
		To     common.Address
		Tokens *big.Int
	}

	balances := make(map[common.Address]*big.Int)
	logTransferSig := []byte("Transfer(address,address,uint256)")
	logTransferSigHash := crypto.Keccak256Hash(logTransferSig)
	totalSupply := big.NewInt(0)
	for _, vLog := range totalLogs {
		switch vLog.Topics[0].Hex() {
		case logTransferSigHash.Hex():
			var transferEvent logTransfer
			result, err := contractAbi.Unpack("Transfer", vLog.Data)
			if err != nil {
				return nil, nil, err
			}
			transferEvent.From = common.HexToAddress(vLog.Topics[1].Hex())
			transferEvent.To = common.HexToAddress(vLog.Topics[2].Hex())
			transferEvent.Tokens = result[0].(*big.Int)

			fromBal := balances[transferEvent.From]
			toBal := balances[transferEvent.To]
			if fromBal == nil {
				fromBal = big.NewInt(0)
			}
			if toBal == nil {
				toBal = big.NewInt(0)
			}

			if bytes.Equal(transferEvent.From.Bytes(), zeroAddr.Bytes()) {
				// Mint if from 0 addr.
				totalSupply.Add(totalSupply, transferEvent.Tokens)
			} else {
				// Add to from if not.
				fromBal = fromBal.Sub(fromBal, transferEvent.Tokens)
			}

			if bytes.Equal(transferEvent.To.Bytes(), deadAddr.Bytes()) {
				// Burn transfers to burn addr.
				totalSupply.Sub(totalSupply, transferEvent.Tokens)
			} else {
				// Add balance if to user.
				toBal = toBal.Add(toBal, transferEvent.Tokens)
			}

			balances[transferEvent.From] = fromBal
			balances[transferEvent.To] = toBal

			if fromBal.Cmp(big.NewInt(0)) == 0 {
				delete(balances, transferEvent.From)
			}
		}
	}
	return balances, totalSupply, nil
}

func marshalJSON(balances map[common.Address]*big.Int, totalSupply *big.Int) map[string]string {
	stringMap := make(map[string]string, len(balances) + 1)
	for k, v := range balances {
		stringMap[k.String()] = v.String()
	}
	return stringMap
}

func unmarshalJSON(jsonMap map[string]string) (map[common.Address]*big.Int, error) {
	balMap := make(map[common.Address]*big.Int, len(jsonMap) - 1)
	for k, v := range jsonMap {
		if k == totalSupplyKey {
			continue
		}
		bigInt, ok := big.NewInt(0).SetString(v, 10)
		if !ok {
			return nil, fmt.Errorf("could not cast %s to big int", v)
		}
		balMap[common.HexToAddress(k)] = bigInt
	}

	return balMap, nil
}
