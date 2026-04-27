package main

import (
	"context"
	"fmt"
	"math/big"
	"runtime"
	"sync"
	"time"

	"zk-htlc/actors"
	"zk-htlc/contracts"
	"zk-htlc/zkp"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
)

func startTPSBenchmark(client *ethclient.Client, auth *bind.TransactOpts, instA *contracts.DataMigration, chainID *big.Int) {
	restoreAllOutput()
	fmt.Printf("🚀 Running TPS Benchmark (Nodes: %d)\n", *flagNodeCount)
	fmt.Printf("⚙️  Pacing: Lock=%dms | Unlock=%dms\n", *flagLockMs, *flagUnlockMs)
	fmt.Println("   (Detailed logs -> debug.log)")
	silenceAllOutput()

	var allMetrics []TPSMetrics
	warmup(client, auth, instA)
	drainPendingPool(client)

	for i := 0; i < *flagRepeat; i++ {
		runtime.GC()
		metrics := runSingleTPSPass(client, auth, instA, chainID)
		allMetrics = append(allMetrics, metrics)
		
		printResultForce("   Run %d/%d -> Lock: %.2f | Unlock: %.2f\n", i+1, *flagRepeat, metrics.TxLockTPS, metrics.TxUnlockTPS)

		if i < *flagRepeat-1 {
			time.Sleep(2 * time.Second)
			drainPendingPool(client)
			newNonce, _ := client.PendingNonceAt(context.Background(), auth.From)
			auth.Nonce = big.NewInt(int64(newNonce))
		}
	}

	smartResult := calculateSmartTPS(allMetrics)
	printResultForce("\n🏆 Final TPS: Lock: %.2f | Unlock: %.2f\n", smartResult.TxLockTPS, smartResult.TxUnlockTPS)
	
	filename := fmt.Sprintf("tps_result_nodes_%d.json", *flagNodeCount)
	ExportTPSMetrics(smartResult, filename)
	printResultForce("💾 Saved: %s\n", filename)
}

func runSingleTPSPass(client *ethclient.Client, auth *bind.TransactOpts, instA *contracts.DataMigration, chainID *big.Int) TPSMetrics {
	txCount := *flagNodeCount * 20; if txCount < 1000 { txCount = 1000 }
	highGasPrice := big.NewInt(200000000000) 
	
	zkpHandler, _ := zkp.NewZKPHandler()
	user, _ := actors.NewUser()
	dspaPkg, _ := user.DistributeInfo()
	dspa := actors.NewDSP("DSPA", dspaPkg.Pre, dspaPkg.Sn, dspaPkg.H, zkpHandler)
	assignA, _ := dspa.GenerateUnlockProof()
	realProof, _ := zkpHandler.Prove(assignA)
	solProof, solPub := formatProofForSolidity(realProof, []*big.Int{dspa.H, dspa.Sn})

	tpsMonitor := NewTPSMonitor(client)
	dataID := [32]byte{0x01}
	timeout := big.NewInt(3600)
	dummyHash := bigIntTo32Bytes(dspa.H)
	startNonce, _ := client.PendingNonceAt(context.Background(), auth.From)
	
	var wg sync.WaitGroup
	guard := make(chan struct{}, 200) 
	currentNonce := startNonce
	pacingLock := time.Duration(*flagLockMs) * time.Millisecond
	pacingUnlock := time.Duration(*flagUnlockMs) * time.Millisecond

	// Congestion Control
	checkCongestion := func() { 
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		p, _ := client.PendingTransactionCount(ctx)
		if p > 1000 { time.Sleep(1 * time.Second) } 
	}

	// Lock Phase
	logToFile("Phase 1: Lock")
	for i := 0; i < txCount; i++ {
		wg.Add(1); guard <- struct{}{}
		if i%50==0 { checkCongestion() }
		txAuth := &bind.TransactOpts{From: auth.From, Signer: auth.Signer, Context: context.Background(), GasLimit: 200000, GasPrice: highGasPrice, Nonce: big.NewInt(int64(currentNonce))}
		currentNonce++
		go func(idx int, opts *bind.TransactOpts) { 
			defer wg.Done(); defer func() { <-guard }()
			tx, err := instA.Lock(opts, dummyHash, dataID, timeout)
			if err == nil { tpsMonitor.RecordProtocolTx(tx.Hash(), fmt.Sprintf("TxLock-%d", idx), time.Now()); go watchTx(client, tx.Hash(), tpsMonitor) }
		}(i, txAuth)
		time.Sleep(pacingLock)
	}
	wg.Wait()
	
	waitTimeout := time.Duration(60 + (*flagNodeCount * 6)) * time.Second
	if !tpsMonitor.WaitForConfirmations(int(float64(txCount)*0.8), waitTimeout) { return tpsMonitor.Stop() }
	
	drainPendingPool(client)
	realNonce, _ := client.PendingNonceAt(context.Background(), auth.From); currentNonce = realNonce 

	// Unlock Phase
	logToFile("Phase 2: Unlock")
	for i := 0; i < txCount; i++ {
		wg.Add(1); guard <- struct{}{}
		if i%50==0 { checkCongestion() }
		txAuth := &bind.TransactOpts{From: auth.From, Signer: auth.Signer, Context: context.Background(), GasLimit: 500000, GasPrice: highGasPrice, Nonce: big.NewInt(int64(currentNonce))}
		currentNonce++
		go func(idx int, opts *bind.TransactOpts) { 
			defer wg.Done(); defer func() { <-guard }()
			tx, err := instA.Unlock(opts, solProof, solPub)
			if err == nil { tpsMonitor.RecordProtocolTx(tx.Hash(), fmt.Sprintf("TxUnlock-%d", idx), time.Now()); go watchTx(client, tx.Hash(), tpsMonitor) }
		}(i, txAuth)
		time.Sleep(pacingUnlock)
	}
	wg.Wait()
	tpsMonitor.WaitForConfirmations(int(float64(txCount)*0.8), waitTimeout)
	
	return tpsMonitor.Stop()
}