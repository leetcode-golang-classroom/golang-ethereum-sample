package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/leetcode-golang-classroom/golang-ethereum-sample/internal/config"
)

func main() {
	ReadKeyFromKeyStore()
}

func GenerateKeyStore() {
	// generate keystore
	key := keystore.NewKeyStore("./wallet", keystore.StandardScryptN, keystore.StandardScryptP)
	account, err := key.NewAccount(config.AppConfig.KeyStorePassword)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(account.Address)
}

func ReadKeyFromKeyStore() {
	b, err := os.ReadFile(fmt.Sprintf("%s/%s", "./wallet", config.AppConfig.KeyStorefile))
	if err != nil {
		log.Fatal(err)
	}
	key, err := keystore.DecryptKey(b, config.AppConfig.KeyStorePassword)
	if err != nil {
		log.Fatal(err)
	}
	pData := crypto.FromECDSA(key.PrivateKey)
	fmt.Println("private key:", hexutil.Encode(pData))
	pubData := crypto.FromECDSAPub(&key.PrivateKey.PublicKey)
	fmt.Println("public key:", hexutil.Encode(pubData))
	address := crypto.PubkeyToAddress(key.PrivateKey.PublicKey).Hex()
	fmt.Println("address:", address)
}
