package main

import (
	"context"
	"crypto/ecdsa"
	"flag"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

const RPCEndpoint = "http://127.0.0.1:8545"

func main() {
	flagNodes := flag.Int("nodes", 20, "节点数量")
	flagPriv := flag.String("priv", "", "用于发送交易的私钥 (Hex)")
	flag.Parse()

	if *flagPriv == "" {
		log.Fatal("❌ 必须提供私钥 (-priv)")
	}

	targetPending := int64(*flagNodes * 5)
	if targetPending < 20 {
		targetPending = 20
	}

	targetTPS := 35.0 - 0.3*float64(*flagNodes)
	if targetTPS < 5 {
		targetTPS = 5
	}
	sendInterval := time.Duration(1000.0/targetTPS) * time.Millisecond

	fmt.Printf("   - 节点数: %d\n", *flagNodes)
	fmt.Printf("   - 目标Pending水位: %d\n", targetPending)
	fmt.Printf("   - 发送限速: %v/tx (~%.1f TPS)\n", sendInterval, targetTPS)

	client, err := ethclient.Dial(RPCEndpoint)
	if err != nil {
		log.Fatal("RPC 连接失败:", err)
	}

	privKey, err := crypto.HexToECDSA(*flagPriv)
	if err != nil {
		log.Fatal("私钥格式错误:", err)
	}

	publicKey := privKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("公钥生成失败")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal("获取 NetworkID 失败:", err)
	}
	signer := types.NewEIP155Signer(chainID)

	startNonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal("获取 Nonce 失败:", err)
	}
	localNonce := startNonce

	to := common.HexToAddress("0x0000000000000000000000000000000000000000")
	payload := make([]byte, 100)

	// 检查间隔不要太快，减轻 RPC 压力
	checkInterval := 200 * time.Millisecond

	for {
		globalPendingUint, err := client.PendingTransactionCount(context.Background())
		if err != nil {
			time.Sleep(1 * time.Second)
			continue
		}
		globalPending := int64(globalPendingUint)
		gap := targetPending - globalPending

		if gap > 0 {

			batchSize := gap
			if batchSize > 10 {
				batchSize = 10
			}
			gasPrice := big.NewInt(2000000000)

			for i := int64(0); i < batchSize; i++ {
				tx := types.NewTransaction(localNonce, to, big.NewInt(1), 50000, gasPrice, payload)
				signedTx, _ := types.SignTx(tx, signer, privKey)

				err := client.SendTransaction(context.Background(), signedTx)
				if err != nil {
					// 只有 Nonce 问题才重置，其他错误稍微冷却
					newNonce, _ := client.PendingNonceAt(context.Background(), fromAddress)
					if newNonce > localNonce {
						localNonce = newNonce
					} else {
						time.Sleep(500 * time.Millisecond)
					}
					break
				}
				localNonce++

				time.Sleep(sendInterval)
			}
		} else {

			time.Sleep(checkInterval)
		}
	}
}
