# Go Merkle Distributor

This repo is a golang implementation (using [sol-merkle-tree-go](https://github.com/0xKiwi/sol-merkle-tree-go)) that generates the required data for a [MerkleDrop](https://github.com/Uniswap/merkle-distributor/blob/master/contracts/MerkleDistributor.sol). It takes in a JSON of address->balances (both in strings) and writes a file containing the data neccessary for a cliam transaction (index, address, amount, and proof).

## Usage
```
go-merkle-distributor --json-file=/path/to/input.json --output-file=/path/to/output.json
```

### Helper functions
```golang
// CreateDistributionTree uses sol-merkle-tree-go to construct a tree of balance items and returns a 
// human readable mapping of address to claim info such as the proof, amount, and index. 
func CreateDistributionTree(holderArray []*TokenHolder) (*solTree.MerkleTree, map[string]ClaimInfo, error)

// ArrayFromAddrBalMap takes a string mapping of address -> balance and converts it into an array for sorting and easier usage. 
func ArrayFromAddrBalMap(stringMap map[string]string) ([]*TokenHolder, error)

// BalancesAndSupplyAtBlock takes in a go-ethereum client, token address and traces 
// through transfer logs from the startBlock, to the endBlock.
func BalancesAndSupplyAtBlock(client *ethclient.Client, tokenAddress common.Address, startBlock int64, endBlock int64) (map[common.Address]*big.Int, *big.Int, error)
```



## Example input
```json
{
    "0xffd342920aD20AE56769b9fC3ba4DD0E5AC75688":"1000000000000000000",
    "0x00dF49020aD20AE56769b9fC341bDD0E5AC75688":"15000000000000000000",
    "0x12dF43220D650AE56769b9fC3b4dDD0E5AC75688":"100000000000000000000",
    "0xf5dF46020aD20AE56769b9fC3da04D0E5AC75688":"1500000000000000000"
}
```


## Example Output
```json
{
    "0x00dF49020AD20Ae56769B9fc341bdd0E5AC75688": {
        "Index": 1,
        "Amount": "15000000000000000000",
        "Proof": [
            "0x21bc7a62dae272690dd4c8b808490d31b99bcb26628732f2ed499c632f04c97a",
            "0x7962ac1a7b4dc9e2befa1d1b063ea0dbb2d29519325718c7688bac06c34505bd"
        ]
    },
    "0x12DF43220d650ae56769b9FC3b4ddD0e5ac75688": {
        "Index": 0,
        "Amount": "100000000000000000000",
        "Proof": [
            "0x9055b91635dde0567458347be8d078524a72dcb754613a122bc5a8b2eb1676d4",
            "0x57e7cdb442d4d4498bbedfe5fbb29278e3b919fe42d3af5913983437f1c2a370"
        ]
    },
    "0xF5DF46020AD20Ae56769B9fC3Da04D0e5Ac75688": {
        "Index": 2,
        "Amount": "1500000000000000000",
        "Proof": [
            "0x17936c9a0951dcca9726377608802c03040c3c5901e0659cf6465a17b80b04dd",
            "0x7962ac1a7b4dc9e2befa1d1b063ea0dbb2d29519325718c7688bac06c34505bd"
        ]
    },
    "0xffd342920aD20aE56769B9FC3BA4DD0E5AC75688": {
        "Index": 3,
        "Amount": "1000000000000000000",
        "Proof": [
            "0x3811a5ba731e1016403500e277995d096d5f0a6e8344070c7b468f8e38e75e7e",
            "0x57e7cdb442d4d4498bbedfe5fbb29278e3b919fe42d3af5913983437f1c2a370"
        ]
    },
    "root": {
        "Index": 0,
        "Amount": "",
        "Proof": [
            "0x00bfd4c68462ab5f09d52c2e9680b52c2a55db9a7cea62d55164d726eeb36f8c"
        ]
    }
}
```
