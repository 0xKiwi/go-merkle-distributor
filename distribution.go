package main

import (
	"encoding/binary"
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

type userInfo struct {
	index uint64
	account common.Address
	amount *big.Int
}

type ClaimInfo struct {
	Index  uint64
	Amount string
	Proof  []string
}

func CreateDistributionTree(holderArray []*TokenHolder) (*solTree.MerkleTree, map[string]ClaimInfo, error) {
	sort.Slice(holderArray,  func(i, j int) bool {
		return holderArray[i].balance.Cmp(holderArray[j].balance) > 0
	})
	elements := make([]*userInfo, len(holderArray))
	for i, holder := range holderArray {
		elements[i] = &userInfo{
			index: uint64(i),
			account: holder.addr,
			amount: holder.balance,
		}
	}

	nodes := make([][]byte, len(elements))
	for i, user := range elements {
		// TODO: this needs to be padded to 256 bits.
		indexBytes := uint64ToBytesLittleEndian(user.index)
		accountBytes := user.account.Bytes()
		amountBytes := user.amount.Bytes()
		nodes[i] = crypto.Keccak256(indexBytes, accountBytes, amountBytes)
	}

	tree, err := solTree.GenerateTreeFromItems(nodes)
	if err != nil {
		return nil, nil, fmt.Errorf("could not generate trie: %v", err)
	}
	//distributionRoot := tree.Root()

	addrToProof := make(map[string]ClaimInfo, len(holderArray))
	for i, holder := range holderArray {
		proof, err := tree.MerkleProof(nodes[i])
		if err != nil {
			return nil, nil, fmt.Errorf("could not generate proof: %v", err)
		}
		addrToProof[holder.addr.String()] = ClaimInfo{
			Index:  uint64(i),
			Amount: holder.balance.String(),
			Proof:  stringArrayFrom2DBytes(proof),
		}
	}
	return tree, addrToProof, nil
}

func stringArrayFrom2DBytes(bytes2d [][]byte) []string {
	stringArray := make([]string, len(bytes2d))
	for i, bytes := range bytes2d {
		stringArray[i] = fmt.Sprintf("%#x", bytes)
	}
	return stringArray
}

func bytes2DFromStringArray(strings []string) [][]byte {
	bytesArray := make([][]byte, len(strings))
	for i, ss := range strings {
		bytesArray[i] = common.HexToHash(ss).Bytes()
	}
	return bytesArray
}

func uint64ToBytesLittleEndian(i uint64) []byte {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, i)
	return buf
}

func uint64To256BytesLittleEndian(i uint64) []byte {
	buf := make([]byte, 32)
	PutUint64Padded(buf, i)
	return buf
}

func PutUint64Padded(b []byte, v uint64) {
	_ = b[31] // early bounds check to guarantee safety of writes below
	b[0] = byte(v)
	b[1] = byte(v >> 8)
	b[2] = byte(v >> 16)
	b[3] = byte(v >> 24)
	b[4] = byte(v >> 32)
	b[5] = byte(v >> 40)
	b[6] = byte(v >> 48)
	b[7] = byte(v >> 56)
	for i := 8; i < len(b); i++ {
		b[i] = byte(0)
	}
}


