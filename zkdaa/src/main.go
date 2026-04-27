package main

import (
	"flag"
	"fmt"
	"log"
	"math/big"

	//"os"
	"time"
)

// 全局变量定义

var (
	flagNodeCount     = flag.Int("nodes", 25, "Nodes count")
	flagContractAddr  = flag.String("addr", "", "Contract Address (deprecated)")
	flagContractAddrA = flag.String("addrA", "", "Contract Address A")
	flagContractAddrB = flag.String("addrB", "", "Contract Address B")

	// 🔥 新增：用于 Single 模式的 flags
	flagSingle    = flag.Bool("single", false, "Run single test mode")
	flagFileSize  = flag.Int("filesize", 1024*1024, "File size in bytes")
	flagChunkSize = flag.Int("chunksize", 1024, "Chunk size in bytes")

	// 原有 flags
	flagRunLatency = flag.Bool("latency", false, "Run latency mode")
	flagRunTPS     = flag.Bool("tps", false, "Run TPS mode")
	flagRunCircuit = flag.Bool("circuit", false, "Run circuit benchmark")
	flagSimulate   = flag.Bool("simulate", false, "Simulate physics")
	flagBatch      = flag.Bool("batch", false, "Run batch tests (deprecated)")
	flagSilent     = flag.Bool("silent", false, "Silent mode")
	flagLogFile    = flag.String("log", "debug.log", "Log file path")
	flagBenchmark  = flag.Bool("benchmark", false, "Run benchmark mode")
	flagRounds     = flag.Int("rounds", 50, "Benchmark iterations")

	flagRepeat   = flag.Int("repeat", 5, "TPS Repetitions")
	flagLockMs   = flag.Int("lock-ms", 12, "Lock interval (ms)")
	flagUnlockMs = flag.Int("unlock-ms", 18, "Unlock interval (ms)")
	privacyTest  = flag.Bool("privacy", false, "Run privacy attack experiment")

	// 🔥 关键修复：恢复全局变量供 latency.go 使用
	fileSize   = 10 * 1024
	chunkSize  = 1024
	outputJSON = flag.String("output", "", "Output JSON path")
	cpuCores   = 0
	skipAudit  = flag.Bool("skipAudit", false, "Skip audit")

	bn254FieldModulus = new(big.Int)
)

func init() {
	bn254FieldModulus.SetString("21888242871839275222246405745257275088548364400416034343698204186575808495617", 10)
}

func runCircuitBenchmark() {
	fmt.Println("Circuit benchmark not implemented yet")
}

// 占位，因为 batch.go 里删除了 runBatchTests
func runBatchTests(chunkSize int, silent bool) {
	fmt.Println("Batch mode is deprecated. Use -single via run_batch_tests.sh instead.")
}

func main() {
	flag.Parse()
	setupFileLogging()
	setupEnv()

	// 1. 不需要链的模式
	if *flagRunCircuit {
		runCircuitBenchmark()
		return
	}

	if *flagSimulate {
		runSimulatedBenchmark(*flagNodeCount * 10)
		return
	}

	if *privacyTest {
		RunPrivacyAttackExperiment()
		return
	}

	if *flagBenchmark {
		if *flagContractAddrA == "" || *flagContractAddrB == "" {
			log.Fatal("❌ -benchmark requires -addrA and -addrB")
		}

		runBenchmarkMode(*flagFileSize, *flagChunkSize,
			*flagContractAddrA, *flagContractAddrB,
			*flagRounds, *outputJSON)
		return
	}

	// 2. 🔥 Single 模式 (Shell 脚本调用的核心入口)
	if *flagSingle {
		if *flagContractAddrA == "" || *flagContractAddrB == "" {
			log.Fatal("❌ Error: -single mode requires -addrA and -addrB")
		}

		fmt.Printf("🔧 Mode: Single Test | Size: %d | Chunk: %d\n", *flagFileSize, *flagChunkSize)
		fmt.Printf("   AddrA: %s\n   AddrB: %s\n", *flagContractAddrA, *flagContractAddrB)

		// 调用 batch.go 中的 runSingleTest
		metrics := runSingleTest(*flagFileSize, *flagChunkSize, *flagContractAddrA, *flagContractAddrB)

		if metrics != nil && *outputJSON != "" {
			exportMetricsToJSON(metrics, *outputJSON)
		}
		return
	}

	// 3. Batch 模式 (旧)
	if *flagBatch {
		runBatchTests(*flagChunkSize, *flagSilent)
		return
	}

	// 4. Latency / TPS 模式
	if *flagRunLatency || *flagRunTPS {
		privKey := "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
		client, auth, chainID, _ := setupBlockchain("http://127.0.0.1:8545", privKey)
		defer client.Close()

		time.Sleep(1 * time.Second)

		addrA := *flagContractAddrA
		addrB := *flagContractAddrB
		if addrA == "" {
			addrA = *flagContractAddr
		}
		if addrB == "" {
			addrB = *flagContractAddr
		}

		if addrA == "" || addrB == "" {
			log.Fatal("Error: Latency/TPS mode requires -addrA and -addrB flags")
		}

		instA := setupContract(addrA, client)
		instB := setupContract(addrB, client)

		if *flagRunLatency {
			startLatencyBenchmark(client, auth, instA, instB)
		} else {
			startTPSBenchmark(client, auth, instA, chainID)
		}
	} else {
		fmt.Println("Usage:")
		fmt.Println("  -single   : Run single benchmark (used by script)")
		fmt.Println("  -latency  : Run latency test")
		fmt.Println("  -tps      : Run TPS test")
	}
}
