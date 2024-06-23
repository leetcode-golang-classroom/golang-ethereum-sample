package main

import (
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

func main() {
	pvk, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}
	// hex value
	pData := crypto.FromECDSA(pvk)
	// private key
	fmt.Println(hexutil.Encode(pData))

	pubData := crypto.FromECDSAPub(&pvk.PublicKey)
	// public key
	fmt.Println(hexutil.Encode(pubData))
	// address
	fmt.Println(crypto.PubkeyToAddress(pvk.PublicKey).Hex())
}
