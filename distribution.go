package main

import (
	"fmt"
	"math/big"
	"sort"

	solTree "github.com/0xKiwi/sol-merkle-tree-go"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type TokenHolder struct {
	addr common.Address
	balance *big.Int
}

type ClaimInfo struct {
	Amount string
	Proof  []string
}

// CreateDistributionTree uses sol-merkle-tree-go to construct a tree of balance items and returns a 
// human readable mapping of address to claim info such as the proof, amount, and index. 
func CreateDistributionTree(holderArray []*TokenHolder) (*solTree.MerkleTree, map[string]ClaimInfo, error) {
	sort.Slice(holderArray,  func(i, j int) bool {
		return holderArray[i].balance.Cmp(holderArray[j].balance) > 0
	})
	nodes := make([][]byte, len(holderArray))
	for i, user := range holderArray {
		packed := append(
			user.addr.Bytes(),
			common.LeftPadBytes(user.balance.Bytes(), 32)...,
		)
		nodes[i] = crypto.Keccak256(packed)
	}

	tree, err := solTree.GenerateTreeFromHashedItems(nodes)
	if err != nil {
		return nil, nil, fmt.Errorf("could not generate trie: %v", err)
	}

	addrToProof := make(map[string]ClaimInfo, len(holderArray))
	for i, holder := range holderArray {
		proof, err := tree.MerkleProof(nodes[i])
		if err != nil {
			return nil, nil, fmt.Errorf("could not generate proof: %v", err)
		}
		addrToProof[holder.addr.String()] = ClaimInfo{
			Amount: holder.balance.String(),
			Proof:  stringArrayFrom2DBytes(proof),
		}
	}
	distributionRoot := tree.Root()
	addrToProof["root"] = ClaimInfo{Proof: []string{fmt.Sprintf("%#x", distributionRoot)}}
	return tree, addrToProof, nil
}

// ArrayFromAddrBalMap takes a string mapping of address -> balance and converts it into an array for sorting and easier usage. 
func ArrayFromAddrBalMap(stringMap map[string]string) ([]*TokenHolder, error) {
	i := 0
	balArray := make([]*TokenHolder, len(stringMap))
	for addr, bal := range stringMap {
		bigInt, ok := big.NewInt(0).SetString(bal, 10)
		if !ok {
			return nil, fmt.Errorf("could not cast %s to big int", bal)
		}
		balArray[i] = &TokenHolder{
			addr: common.HexToAddress(addr),
			balance: bigInt,
		}
		i++
	}
	return balArray, nil
}
