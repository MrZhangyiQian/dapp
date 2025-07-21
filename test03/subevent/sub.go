package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/gagliardetto/solana-go/rpc/ws"
)

func main() {
	log.Println("🔍 启动 Solana 交易监控...")
	defer log.Println("🛑 监控服务停止")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 使用更可靠的 RPC 节点
	rpcEndpoints := []string{
		rpc.MainNetBeta_WS, // 主网
		rpc.TestNet_WS,
		"wss://api.mainnet-beta.solana.com",          // 官方主网
		"wss://solana-api.projectserum.com",          // Serum 提供的节点
		"wss://solana-mainnet.g.alchemy.com/v2/demo", // Alchemy 提供的节点
	}

	// 尝试连接多个节点
	for i, endpoint := range rpcEndpoints {
		log.Printf("🔌 尝试连接节点 [%d/%d]: %s", i+1, len(rpcEndpoints), endpoint)

		client, err := ws.Connect(ctx, endpoint)
		if err != nil {
			log.Printf("⚠️ 连接失败: %v", err)
			continue
		}

		log.Printf("✅ 成功连接到节点: %s", endpoint)

		// 订阅代币转账事件
		sub, err := client.LogsSubscribeMentions(
			solana.TokenProgramID,
			rpc.CommitmentConfirmed,
		)
		if err != nil {
			log.Printf("⚠️ 订阅失败: %v", err)
			client.Close()
			continue
		}

		log.Println("✅ 监控已启动 (按 Ctrl+C 退出)")

		// 设置优雅退出
		signalCh := make(chan os.Signal, 1)
		signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

		// 处理接收事件
		go func() {
			for {
				select {
				case <-ctx.Done():
					return
				default:
					got, err := sub.Recv(ctx)
					if err != nil {
						log.Printf("⚠️ 接收错误: %v", err)
						return
					}

					sig := got.Value.Signature
					if !sig.IsZero() {
						log.Printf("📥 捕获交易: %s...", sig.String()[:8])
					}
				}
			}
		}()

		// 等待退出信号
		<-signalCh
		log.Println("🛑 收到停止信号，关闭连接")
		sub.Unsubscribe()
		client.Close()
		return
	}

	log.Println("❌ 所有节点连接尝试失败，请检查网络连接")
}
