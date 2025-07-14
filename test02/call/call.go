package call

import (
	"context"
	"crypto/ecdsa"
	"dapp/test02/counter"
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
)

func InteractWithContract(contractAddress common.Address) {
	// 根据实际路径调整
	err := godotenv.Load("../.env")
	client, err := ethclient.Dial("https://eth-sepolia.g.alchemy.com/v2/" + os.Getenv("ALCHEMY_API_KEY"))
	if err != nil {
		log.Fatal(err)
	}

	privateKeyStr := os.Getenv("PRIVATE_KEY")
	privateKey, err := crypto.HexToECDSA(privateKeyStr)
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	chainId, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainId)
	if err != nil {
		log.Fatal(err)
	}

	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasLimit = uint64(300000)
	auth.GasPrice = gasPrice

	instance, err := counter.NewCounter(contractAddress, client)
	if err != nil {
		log.Fatal(err)
	}

	// 调用 increment 方法
	tx, err := instance.Increment(auth)
	if err != nil {
		log.Fatal("Failed to send transaction:", err)
	}
	fmt.Println("Increment transaction sent:", tx.Hash().Hex())

	// 等待交易确认
	receipt, err := bind.WaitMined(context.Background(), client, tx)
	if err != nil {
		log.Fatal("Transaction not mined:", err)
	}
	fmt.Println("Transaction mined in block:", receipt.BlockNumber)

	// 获取当前计数器值
	countValue, err := instance.GetCount(nil)
	if err != nil {
		log.Fatal("Failed to get count:", err)
	}

	fmt.Printf("Current count: %v\n", countValue)
}
