package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/leetcode-golang-classroom/golang-ethereum-sample/internal/config"
)

func main() {
	// setup client connect to jsonprc node
	client, err := ethclient.DialContext(context.Background(), config.AppConfig.EthJsonRpcURL)
	if err != nil {
		log.Fatalf("Error to create a ether client: %v", err)
	}
	defer client.Close()
	block, err := client.BlockByNumber(context.Background(), nil)
	if err != nil {
		log.Fatalf("Error to get a block: %v", err)
	}
	// get latest block number
	fmt.Println("The block number:", block.Number())

	// find specific address balance
	addr := config.AppConfig.EthAddress
	address := common.HexToAddress(addr)

	balance, err := client.BalanceAt(context.Background(), address, nil)
	if err != nil {
		log.Fatalf("Error to get the balance:%v", err)
	}
	fmt.Println("The balance:", balance)
	// 1 ether = 10^18 wei
	fBalance := new(big.Float)
	fBalance.SetString(balance.String())
	balanceEther := new(big.Float).Quo(fBalance, big.NewFloat(math.Pow10(18)))
	fmt.Println("address:", config.AppConfig.EthAddress, "has", balanceEther, "ether")
}
