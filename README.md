# golang-ethereum-sample

This repository is to demo how to interact with ethereum jsonrpc server by golang

## install dependency

```shell
go get github.com/ethereum/go-ethereum
```

## connect to ETH JSONRPC

```golang
package main

import (
	"context"
	"fmt"
	"log"

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
	fmt.Println(block.Number())
}
```

## setup local ethereum node with ganache

1. install ganache-cli
  
```shell
pnpm add -g ganache-cli
```

2. start ganache node

```shell
ganache-cli
```

![ganache-node-info](ganache-node-info.png)

## setup ENV for local node

```yaml
ETH_JSON_RPC_URL=http://localhost:8545
```

## add lookup balance logic

```golang
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
```

## wallet concept

![wallet-address](wallet-address.png)

## write logic to generate private key, public key and address

```golang
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
```

## wallet (keystore) concept

![wallet-keystore-concept](wallet-keystore-concept.png)

1. Generate keystore with password
```golang
func GenerateKeyStore() {
	// generate keystore
	key := keystore.NewKeyStore("./wallet", keystore.StandardScryptN, keystore.StandardScryptP)
	account, err := key.NewAccount(config.AppConfig.KeyStorePassword)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(account.Address)
}
```
2. Read key from keystore with password
```golang
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
```

## make transaction concept

![transaction-concept](transaction-concept.png)

## implementation

```golang
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
```

## use remix.ethereum.org as IDE

https://remix.ethereum.org/#lang=en&optimize=false&runs=200&evmVersion=null&version=soljson-v0.8.26+commit.8a97fa7a.js

create solidity

```solidity
pragma solidity >=0.8.2 <0.9.0;

contract Todo {
    Task[] tasks;
    struct Task {
        string content;
        bool status;
    }
    constructor() {

    }
    function add(string memory _content) public {
      tasks.push(Task(_content, false));
    }
    function get(uint _id) public view returns (Task memory) {
        return tasks[_id];
    }
    function list() public  view returns (Task[] memory) {
        return tasks;
    }
    function update(uint _id, string memory _content) public {
        tasks[_id].content = _content;
    }
    function remove(uint _id) public {
        delete tasks[_id];
    }
}
```

test with remix vm

![deploy-contract-to-vm](./deploy-contract-to-vm.png)

## interact contract with public function

![interact-with-public-function](./interact-with-public-function.png)

## check execution log of vm on console

![execution-console-log](./execution-console-log.png)

## add permission for access check

```solidity
pragma solidity >=0.8.2 <0.9.0;

contract Todo {
    address public owner;
    Task[] tasks;
    struct Task {
        string content;
        bool status;
    }
    constructor() {
        owner = msg.sender;
    }
    modifier isOwner() {
        require(owner == msg.sender);
        _;  
    }
    function add(string memory _content) public isOwner {
        tasks.push(Task(_content, false));
    }
    function get(uint _id) public isOwner view returns (Task memory) {
        return tasks[_id];
    }
    function list() public  isOwner view returns (Task[] memory) {
        return tasks;
    }
    function update(uint _id, string memory _content) public isOwner{
        tasks[_id].content = _content;
    }
    function remove(uint _id) public isOwner {
        for (uint i = _id; i < tasks.length -1; i++ ) {
            tasks[i] = tasks[i+1];
        }
        tasks.pop();
    }
}
```
## setup owner storage
```solidity
	address public owner;
	constructor() {
		owner = msg.sender;
	}
```
the msg.sender is special key word refer contract creator

owner will set to the creator of the contract
## implement Modifier isOwner for check only allow owner execution

```solidity
	modifier isOwner() {
			require(owner == msg.sender);
			_;  
	}
```
use the required keyword for assertion on check owner is msg.sender

_ the underscope is for rest of the other code

## success execution on modifier
![success-execution-modifier](success-execution-modifier.png)

## failed execution on modifier

![fail-execution-modifier](fail-execution-modifier.png)

[modifier-contract](https://remix.ethereum.org/#lang=en&optimize=false&runs=200&evmVersion=null&version=soljson-v0.8.26+commit.8a97fa7a.js)