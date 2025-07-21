package transaction

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/blocto/solana-go-sdk/client"
	"github.com/blocto/solana-go-sdk/common"
	"github.com/blocto/solana-go-sdk/program/system"
	"github.com/blocto/solana-go-sdk/types"
	"github.com/joho/godotenv"
)

func Transaction() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	// 1. 连接到 Solana Devnet
	c := client.NewClient("https://solana-devnet.g.alchemy.com/v2/" + os.Getenv("ALCHEMY_API_KEY"))
	fmt.Println("✅ 已连接到 Solana Devnet")

	// 2. 定义发送方和接收方账户
	// 注意：请替换为您的实际私钥
	senderPrivateKey := os.Getenv("AC1_PRIVATE_KEY")  // 发送方私钥
	receiverPublicKey := os.Getenv("AC2_PRIVATE_KEY") // 接收方公钥

	// 3. 导入发送方账户
	senderAccount, err := importAccount(senderPrivateKey)
	if err != nil {
		log.Fatalf("❌ 导入发送方账户失败: %v", err)
	}
	fmt.Printf("👤 发送方账户导入成功:\n公钥: %s\n", senderAccount.PublicKey.ToBase58())

	// 4. 验证接收方地址
	receiverPubkey := common.PublicKeyFromString(receiverPublicKey)
	if err != nil {
		log.Fatalf("❌ 接收方地址无效: %v", err)
	}
	fmt.Printf("👤 接收方账户: %s\n", receiverPubkey.ToBase58())

	// 5. 检查发送方余额
	fmt.Println("💰 检查发送方余额...")
	balance, err := c.GetBalance(context.Background(), senderAccount.PublicKey.ToBase58())
	if err != nil {
		log.Fatalf("❌ 获取余额失败: %v", err)
	}

	balanceSOL := float64(balance) / 1e9
	fmt.Printf("发送方余额: %.9f SOL\n", balanceSOL)

	// 如果余额不足，建议手动获取测试 SOL
	if balance < 0.01*1e9 { // 小于 0.01 SOL
		fmt.Println("⚠️ 余额不足，请先获取测试 SOL")
		fmt.Println("手动获取方法: https://faucet.solana.com/")
		return
	}

	// 6. 设置转账金额 (0.01 SOL)
	amountLamports := uint64(0.01 * 1e9)
	fmt.Printf("\n🔄 准备转账 %.9f SOL\n", float64(amountLamports)/1e9)

	// 7. 创建转账交易
	fmt.Println("🔄 创建转账交易...")
	recentBlockhash, err := c.GetLatestBlockhash(context.Background())
	if err != nil {
		log.Fatalf("❌ 获取区块哈希失败: %v", err)
	}

	// 创建转账指令（修复点）
	instruction := system.Transfer(
		system.TransferParam{
			From:   senderAccount.PublicKey,
			To:     receiverPubkey,
			Amount: amountLamports,
		},
	)

	// 构建交易消息
	message := types.NewMessage(types.NewMessageParam{
		FeePayer:        senderAccount.PublicKey,
		RecentBlockhash: recentBlockhash.Blockhash,
		Instructions:    []types.Instruction{instruction},
	})

	// 构建交易
	tx, err := types.NewTransaction(types.NewTransactionParam{
		Message: message,
		Signers: []types.Account{senderAccount},
	})
	if err != nil {
		log.Fatalf("❌ 创建交易失败: %v", err)
	}
	// 8. 发送交易
	fmt.Println("📤 发送交易...")
	txhash, err := c.SendTransaction(context.Background(), tx)
	if err != nil {
		log.Fatalf("❌ 发送交易失败: %v", err)
	}
	fmt.Printf("⏳ 交易已发送，哈希: %s\n等待确认...\n", txhash)

	maxRetries := 10
	retryInterval := 3 * time.Second

	for i := 0; i < maxRetries; i++ {
		txStatus, err := c.GetTransaction(context.Background(), txhash)
		if err != nil {
			log.Printf("⚠️ 获取交易状态失败 (重试 %d/%d): %v", i+1, maxRetries, err)
			time.Sleep(retryInterval)
			continue
		}

		if txStatus != nil {
			if txStatus.Meta != nil && txStatus.Meta.Err != nil {
				log.Fatalf("❌ 交易执行失败: %v", txStatus.Meta.Err)
			}
			fmt.Println("✅ 交易已确认!")
			break
		}

		fmt.Printf("⏳ 交易未确认 (重试 %d/%d)...\n", i+1, maxRetries)
		time.Sleep(retryInterval)

		if i == maxRetries-1 {
			log.Fatalf("❌ 交易确认超时，请检查交易: %s", txhash)
		}
	} // 修复结束

	// 10. 验证结果
	fmt.Println("\n✅ 转账成功!")
	fmt.Printf("发送方: %s\n接收方: %s\n金额: %.9f SOL\n",
		senderAccount.PublicKey.ToBase58(), receiverPubkey.ToBase58(), float64(amountLamports)/1e9)
	fmt.Printf("🌐 在浏览器查看交易: https://explorer.solana.com/tx/%s?cluster=devnet\n", txhash)

	// 	✅ 已连接到 Solana Devnet
	// 👤 发送方账户导入成功:
	// 公钥: D4UH9vLwgweWJXAoC59fciWCukYFuAeUD63e39dt9j68
	// 👤 接收方账户: 6FC7k5MpSGemsvoYTL4UFokRmz74euZXV7q1UaehXKMY
	// 💰 检查发送方余额...
	// 发送方余额: 2.000000000 SOL

	// 🔄 准备转账 0.010000000 SOL
	// 🔄 创建转账交易...
	// 📤 发送交易...
	// ⏳ 交易已发送，哈希: 2q8eupP64Suq65GPfoTRJ4EGqwdoStRidBqwgcNFvXAZ5voeUFoMexmKVetGHqsWzprb6YUe98hwEXfu7JwcZZnk
	// 等待确认...
	// ⏳ 交易未确认 (重试 1/10)...
	// ⏳ 交易未确认 (重试 2/10)...
	// ⏳ 交易未确认 (重试 3/10)...
	// ⏳ 交易未确认 (重试 4/10)...
	// ✅ 交易已确认!

	// ✅ 转账成功!
	// 发送方: D4UH9vLwgweWJXAoC59fciWCukYFuAeUD63e39dt9j68
	// 接收方: 6FC7k5MpSGemsvoYTL4UFokRmz74euZXV7q1UaehXKMY
	// 金额: 0.010000000 SOL
	// 🌐 在浏览器查看交易: https://explorer.solana.com/tx/2q8eupP64Suq65GPfoTRJ4EGqwdoStRidBqwgcNFvXAZ5voeUFoMexmKVetGHqsWzprb6YUe98hwEXfu7JwcZZnk?cluster=devnet
}

// 从私钥导入账户
func importAccount(privateKey string) (types.Account, error) {
	return types.AccountFromBase58(privateKey)
}
