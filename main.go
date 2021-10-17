package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"strconv"
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
	log.Printf("Generating claim info for %s\n", fullPath)
	jsonBytes, err := ioutil.ReadFile(fullPath)
	if err != nil {
		log.Fatalf("Could not read file: %v", err)
	}
	var stringJson []map[string]string
	if err := json.Unmarshal(jsonBytes, &stringJson); err != nil {
		log.Fatalf("Could not unmarshal json: %v", err)
	}
	allMetadataUints := allMetadataToUint256(stringJson)
	allMetadata, err := MetadataFromJSON(allMetadataUints)
	if err != nil {
		log.Fatal(err)
	}
	_, addrToClaim, err := CreateDistributionTree(allMetadata)
	if err != nil {
		log.Fatalf("Could not create distribution tree: %v", err)
	}
	log.Printf("Root: %s\n", addrToClaim["root"].Proof[0])
	if _, err := createFile(*outputFile, addrToClaim); err != nil {
		log.Fatalf("Could not create file: %v", err)
	}
	log.Printf("Created claim info file at %s\n", *outputFile)
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

var background = []string{
	"Blush",
	"Gray",
	"Hot Pink",
	"Lavender",
	"Marsh",
	"Ocean",
	"Orange",
	"Pink",
	"Rose",
	"Ruby",
	"Sky",
}

var skin = []string{
	"Blushing",
	"Freckled",
	"Bronze",
	"Sepia",
}

var outfits = []string{
	"White Tanktop",
	"Baka Tee",
	"Green Tanktop",
	"White Dress",
	"Petite Gym",
	"Zipper Top",
	"Busty Gym",
	"Black Dress",
	"Bored Tee",
	"Hoodie",
	"Orange Sweater",
	"Jersey",
	"Designer",
	"YFI Tee",
	"Occult Dress",
	"Dapper",
	"Qipao",
	"Kimono",
	"Military Coat",
	"NFTX Hoodie",
	"Biker",
	"Wagie",
	"Pink Overalls",
	"Baelien",
	"Stealth Dress",
	"Racer",
	"Admiral",
	"Queen",
	"Bunny",
	"Black Suit",
	"Mechsuit",
	"Armor",
}

var rightaccessory = []string{
	"No Right Accessory",
	"Boba",
	"Sushi",
	"Pizza",
	"Basketball",
	"Pancakes",
	"Katana",
	"Blue Fire",
	"Giant Sword",
	"Gold Scythe",
}

var leftaccessory = []string{
	"No Left Accessory",
	"Peace Hand",
	"Lollipop",
	"DEFI NOTE",
	"Doggo",
	"Skulls",
	"Fomo Doll",
	"Fren",
	"Cyborg Arm",
	"Gold Claws",
	"Brown Doggo",
}

var hair = []string{
	"Blonde Ponytail",
	"Rose Ponytail",
	"Blonde Long",
	"Brown Short Bangs",
	"Black Short Bangs",
	"Black Pixie",
	"Black Braids",
	"Blonde Twintails",
	"Purple Messy",
	"Orange Long Bangs",
	"Purple Braids",
	"Pink Messy",
	"Silver Pixie",
	"Blue Short Bangs",
	"Pink Twintails",
	"Turquoise Twintails",
	"White Long Bangs",
	"Rainbow Braids",
	"Black Long Bangs",
}

var hat = []string{
	"No Hat",
	"Black Bow",
	"Pink Bow",
	"Hibiscus Flower",
	"Elf Ears",
	"Bone Pins",
	"Horns",
	"Flower Crown",
	"Alien Antenna",
	"Demon Horns",
	"Mr. Whiskers",
	"Fox Mask",
	"Floor Cap",
	"Occult Hat",
	"Ornate Flower",
	"Wagie Cap",
	"Bunny Ears",
	"Admiral's Hat",
	"Queen's Crown",
	"Hair Bands",
}

var eyes = []string{
	"Round Glasses",
	"Blue Neutral",
	"Green Wink",
	"Eyepatch",
	"Green Determined",
	"Relaxed",
	"Gold Determined",
	"Purple Neutral",
	"Closed",
	"Pink Wink",
	"Hypnotized",
	"Dual Color Neutral",
	"Lovey Dovey",
	"Pink Glasses",
	"Heterochromia Determined",
	"Blindfold",
	"Heterochromia Neutral",
	"Shades",
	"Gold Shades",
}

var mouths = []string{
	"Content Smile",
	"uwu",
	"Grin",
	"Playful",
	"Pout",
	"Hmph",
	"Tongue Out",
	"Lipstick",
	"Open Mouth",
	"Yelling",
	"Bubblegum",
	"Black Mask",
}

var attributes = map[string][]string{
	"Background": background,
	"Skin": skin,
	"Outfit": outfits,
	"Right Accessory": rightaccessory,
	"Left Accessory": leftaccessory,
	"Hair": hair,
	"Hat": hat,
	"Eyes": eyes,
	"Mouth": mouths,
}

func allMetadataToUint256(allMetadata []map[string]string) map[string]string {
	allMetadataUint := make(map[string]string, len(allMetadata))
	for i, metadatum := range allMetadata {
		totalData := make([]byte, 0, len(metadatum))
		for key, value := range metadatum {
			dataId := findIndexIn(attributes[key], value)
			fmt.Printf("%s: %s into key %d: %#x\n", key, value, dataId, totalData)
			if dataId == -1 {
				log.Fatal(dataId)
			}
			totalData = append(totalData, uint8(dataId))
		}
		iStr := strconv.Itoa(i)
		allMetadataUint[iStr] = fmt.Sprintf("%#x", totalData)
	}
	return allMetadataUint
}

func findIndexIn(strings []string, item string) int {
	for i, s := range strings {
		if s == item {
			return i
		}
	}
	return -1
}
