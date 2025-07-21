package query

import (
	"context"
	"fmt"

	"github.com/blocto/solana-go-sdk/client"
)

func QueryBlock() {
	// 创建RPC客户端
	c := client.NewClient("https://solana-devnet.g.alchemy.com/v2/M0quawxer68-8mXsa8UBw1y4Dtmp-qx3")

	// 获取最新 slot
	latestSlot, err := c.GetSlot(context.Background())
	if err != nil {
		panic("获取最新 slot 失败: " + err.Error())
	}

	// 获取最新区块
	recentBlock, err := c.GetBlock(context.Background(), latestSlot)
	if err != nil {
		panic("查询失败: " + err.Error())
	}

	fmt.Printf("区块高度: %d\n", recentBlock.BlockHeight)
	fmt.Printf("交易数量: %d\n", len(recentBlock.Transactions))

	// 查询的账户地址
	accountPubkey := "D4UH9vLwgweWJXAoC59fciWCukYFuAeUD63e39dt9j68"
	accountInfo, err := c.GetAccountInfo(context.Background(), accountPubkey)
	if err != nil {
		panic("查询账户信息失败: " + err.Error())
	}
	balanceInSol := float64(accountInfo.Lamports)
	fmt.Printf("账户地址： %s\n", accountPubkey)
	fmt.Printf("账户余额： %d lamports\n", accountInfo.Lamports)
	fmt.Printf("账户余额： %.9f SOL\n", balanceInSol)
}
