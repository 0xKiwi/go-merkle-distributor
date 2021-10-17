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

type TokenMetadata struct {
	Id *big.Int
	Metadata *big.Int
}

type ClaimInfo struct {
	Id  uint64
	Metadata string
	Proof  []string
}

var contractAddr = common.HexToAddress("0xF75140376D246D8B1E5B8a48E3f00772468b3c0c")

// CreateDistributionTree uses sol-merkle-tree-go to construct a tree of balance items and returns a 
// human readable mapping of address to claim info such as the proof, amount, and index. 
func CreateDistributionTree(holderArray []*TokenMetadata) (*solTree.MerkleTree, map[string]ClaimInfo, error) {
	sort.Slice(holderArray,  func(i, j int) bool {
		return holderArray[i].Id.Cmp(holderArray[j].Id) > 0
	})
	nodes := make([][]byte, len(holderArray))
	for i, user := range holderArray {
		packed := append(
			contractAddr.Bytes(),
			append(
				common.LeftPadBytes(user.Id.Bytes(), 32),
				common.LeftPadBytes(user.Metadata.Bytes(), 32)...,
			)...,
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
		addrToProof[holder.Id.String()] = ClaimInfo{
			Id:  uint64(i),
			Metadata: holder.Metadata.String(),
			Proof:  stringArrayFrom2DBytes(proof),
		}
	}
	distributionRoot := tree.Root()
	addrToProof["root"] = ClaimInfo{Proof: []string{fmt.Sprintf("%#x", distributionRoot)}}
	return tree, addrToProof, nil
}

func MetadataFromJSON(stringMap map[string]string) ([]*TokenMetadata, error) {
	i := 0
	dataArray := make([]*TokenMetadata, len(stringMap))
	for tokenId, metadata := range stringMap {
		bytes := common.Hex2Bytes(metadata[2:])
		tokenBig, ok := big.NewInt(0).SetString(tokenId, 10)
		if !ok {
			return nil, fmt.Errorf("could not cast %s to big int", tokenBig)
		}
		metadataBig := big.NewInt(0).SetBytes(bytes)
		dataArray[i] = &TokenMetadata{
			Id: tokenBig,
			Metadata: metadataBig,
		}
		i++
	}
	return dataArray, nil
}
