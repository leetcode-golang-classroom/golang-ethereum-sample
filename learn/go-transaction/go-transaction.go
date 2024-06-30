package main

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/leetcode-golang-classroom/golang-ethereum-sample/internal/config"
)

func main() {
	client, err := ethclient.Dial(config.AppConfig.EthJsonRpcURL)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()
	sender := common.HexToAddress("44e8c6ee2eaD166A03b93d4e2bA6F9d299f2Cf44")
	receiver := common.HexToAddress("95fe3f0a532a4feb1bc92c45d447ac66ac54b4d5")
	_, err = GetBalances(client, sender, receiver)
	if err != nil {
		log.Fatal(err)
	}
	// 1 ether = 1000000000000000000 wei
	// 0.1 ETH
	err = Transfer(client, sender, receiver, 100000000000000000, config.AppConfig.SenderPriveKey)
	if err != nil {
		log.Fatal(err)
	}
	_, err = GetBalances(client, sender, receiver)
	if err != nil {
		log.Fatal(err)
	}
}

func GetBalances(client *ethclient.Client, sender, receiver common.Address) ([]*big.Int, error) {
	senderBalance, err := client.BalanceAt(context.Background(), sender, nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	receiverBalance, err := client.BalanceAt(context.Background(), receiver, nil)
	if err != nil {
		log.Println(err)
	}
	fmt.Println("senderBalance:", senderBalance)
	fmt.Println("receiverBalance:", receiverBalance)
	return []*big.Int{senderBalance, receiverBalance}, nil
}
func Transfer(client *ethclient.Client, sender common.Address,
	receiver common.Address, transferAmount int64, senderPrivKey string) error {
	nonce, err := client.PendingNonceAt(context.Background(), sender)
	if err != nil {
		log.Fatal(err)
	}
	amount := big.NewInt(transferAmount)
	suggestGasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	tx := types.NewTransaction(nonce, receiver, amount, 21000, suggestGasPrice, nil)
	chainId, err := client.NetworkID(context.Background())
	if err != nil {
		log.Println(err)
		return err
	}
	prvKey, err := crypto.HexToECDSA(senderPrivKey)
	if err != nil {
		log.Println(err)
		return err
	}
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainId), prvKey)
	if err != nil {
		log.Println(err)
		return err
	}
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Println(err)
		return err
	}
	fmt.Printf("tx send: %s\n", signedTx.Hash().Hex())
	return nil
}
