package main

import (
	"dapp/test02/call"
	"dapp/test02/deploy"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/joho/godotenv"
)

func main() {
	// 部署合约
	deploy.DeployCount()
	// 根据实际路径调整
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal(err)
	}
	comAddress := os.Getenv("CONTRACT_ADDRESS")
	contractAddress := common.HexToAddress(comAddress) // 替换为实际合约地址
	call.InteractWithContract(contractAddress)
}
