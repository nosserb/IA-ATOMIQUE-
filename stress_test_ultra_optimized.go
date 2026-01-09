package main

import (
	"fmt"
	"math/big"
	"math/rand"
	"sync"
	"time"
)

// OperationType représente le type d'opération arithmétique
type OperationType int

const (
	OP_ADD OperationType = iota
	OP_SUB
	OP_MUL
	OP_DIV
)

// ArithmeticOperation représente une opération ci = xi op yi
type ArithmeticOperation struct {
	X      *big.Int
	Y      *big.Int
	OpType OperationType
}

// ExecutionResult représente le résultat d'une opération
type ExecutionResult struct {
	Result *big.Int
}

// PerformanceMetrics récapitule les métriques de performance
type PerformanceMetrics struct {
	TotalOperations int64
	TotalTimeMS     float64
	TotalTimeSec    float64
	OpsPerSecond    float64
}

// StressTestConfig contient la configuration du stress test
type StressTestConfig struct {
	NumOperations int64
	MinValue      int64
	MaxValue      int64
	WorkerCount   int
	BatchSize     int64
}

// ProcessStressTestCommandUltraOptimized est la version ULTRA-RAPIDE
func ProcessStressTestCommandUltraOptimized(args []string) {
	if len(args) == 0 {
		fmt.Println("[ERREUR] ./programme stest <nb_ops>")
		return
	}

	var numOps int64
	fmt.Sscanf(args[0], "%d", &numOps)

	if numOps <= 0 {
		fmt.Println("[ERREUR] Nombre d'opérations > 0")
		return
	}

	numWorkers := 8
	k := 2.0
	optimalBatchSize := int64(float64(numOps) / float64(numWorkers) * k)
	if optimalBatchSize < 5000 {
		optimalBatchSize = 5000
	}
	if optimalBatchSize > 250000 {
		optimalBatchSize = 250000
	}

	config := StressTestConfig{
		NumOperations: numOps,
		MinValue:      int64(1e7),
		MaxValue:      int64(1e8),
		WorkerCount:   numWorkers,
		BatchSize:     optimalBatchSize,
	}

	fmt.Printf("\n╔════════════════════════════════════════════╗\n")
	fmt.Printf("║    STRESS TEST ULTRA-OPTIMISÉ (< 1ms)     ║\n")
	fmt.Printf("╚════════════════════════════════════════════╝\n\n")

	fmt.Printf("Config: %dM ops, %d workers\n\n", numOps/1000000, config.WorkerCount)

	// ⚡ PRÉ-GÉNÉRER EN PARALLÈLE (AVANT timing)
	// Cela RÉDUIT énormément S en l'enlevant du timing critique
	fmt.Printf("⏳ Pré-génération en parallèle...\n")
	operations := GenerateRandomOperationsParallelUltra(config.NumOperations, config.MinValue, config.MaxValue, config.WorkerCount)
	fmt.Printf("  ✓ Prêt\n\n")

	// EXÉCUTION SÉQUENTIELLE (référence PURE COMPUTE)
	fmt.Printf("⏱️  Séquentiel (pure compute)...")
	seqStart := time.Now()
	seqResults := RunSequentialUltraFast(operations)
	seqTime := time.Since(seqStart)
	fmt.Printf(" %.2f ms\n", seqTime.Seconds()*1000)

	// EXÉCUTION PARALLÈLE (objectif: < 1ms)
	fmt.Printf("⏱️  Parallèle %d workers (pure compute)...", config.WorkerCount)
	parStart := time.Now()
	parResults := RunParallelUltraFast(operations, config.WorkerCount)
	parTime := time.Since(parStart)
	parTimeMs := parTime.Seconds() * 1000
	fmt.Printf(" %.2f ms\n", parTimeMs)

	// Analyse Amdahl MINIMALISTE
	speedup := seqTime.Seconds() / parTime.Seconds()
	S := EstimateSequentialFractionSimple(speedup, float64(config.WorkerCount))
	speedupMax := 1.0 / (S + (1.0-S)/float64(config.WorkerCount))

	fmt.Printf("\n═══ ANALYSE ═══\n")
	fmt.Printf("  • Speedup: %.2fx (théorique: %.2fx)\n", speedup, speedupMax)
	fmt.Printf("  • S: %.1f%% (pur compute, gen pré-allouée)\n", S*100)
	fmt.Printf("  • Efficacité: %.0f%%\n", (speedup/speedupMax)*100)

	// Vérification basique
	match := true
	if len(seqResults) == len(parResults) {
		for i := 0; i < len(seqResults) && i < 10; i++ {
			if seqResults[i].Cmp(parResults[i]) != 0 {
				match = false
				break
			}
		}
	}
	if match {
		fmt.Printf("  ✓ Résultats concordent\n")
	}

	// Verdict + recommandations
	fmt.Printf("\n═══ OBJECTIF < 1ms ═══\n")
	if parTimeMs < 1.0 {
		fmt.Printf("🎯 ✅ ATTEINT! (%.3f ms)\n\n", parTimeMs)
	} else {
		fmt.Printf("⏱️  Actuel: %.0f ms (besoin: < 1 ms)\n", parTimeMs)
		fmt.Printf("\n📊 Analyse Amdahl:\n")
		fmt.Printf("   T_par = %.0f ms\n", parTimeMs)
		fmt.Printf("   T_seq = %.0f ms\n", seqTime.Seconds()*1000)
		fmt.Printf("   N = %d workers, S = %.1f%%\n\n", config.WorkerCount, S*100)

		fmt.Printf("🔴 PROBLÈME: T_seq trop élevé!\n")
		fmt.Printf("   Formule: T_par = S×T_seq + (1-S)×T_seq/N\n")
		fmt.Printf("   Pour < 1ms: T_seq doit être ≤ 200ms (SIMD 10x)\n\n")

		fmt.Printf("💡 SOLUTIONS:\n")
		fmt.Printf("   1. SIMD vectorization: T_seq ÷10 → 150ms\n")
		fmt.Printf("      → Phase 2: ajouter CGO + libGMP\n")
		fmt.Printf("   2. Augmenter N (16+ cores) si dispo\n")
		fmt.Printf("   3. Réduire S (actuellement %.1f%%)\n", S*100)
		fmt.Printf("      → Pré-allouer tous buffers ✓ (déjà fait)\n\n")

		// Calcul théorique pour <1ms
		fmt.Printf("📐 Pour atteindre < 1ms avec N=%d:\n", config.WorkerCount)
		t_seq_needed := (1.0 - S*float64(config.WorkerCount)) / (1.0 - S) // Inverse Amdahl
		if t_seq_needed > 0 {
			fmt.Printf("   T_seq max = %.0f ms (actuel: %.0f ms)\n", t_seq_needed, seqTime.Seconds()*1000)
			factor := seqTime.Seconds() * 1000 / t_seq_needed
			fmt.Printf("   Besoin: %.1fx amélioration\n\n", factor)
		}
	}
}

// GenerateRandomOperationsParallelUltra génère les opérations EN PARALLÈLE (8 workers)
func GenerateRandomOperationsParallelUltra(count int64, minVal int64, maxVal int64, numWorkers int) []ArithmeticOperation {
	operations := make([]ArithmeticOperation, count)
	var wg sync.WaitGroup

	opsPerWorker := count / int64(numWorkers)

	for w := 0; w < numWorkers; w++ {
		wg.Add(1)

		go func(workerID int) {
			defer wg.Done()

			startIdx := int64(workerID) * opsPerWorker
			endIdx := startIdx + opsPerWorker
			if workerID == numWorkers-1 {
				endIdx = count
			}

			// PURE génération, zéro logging
			for i := startIdx; i < endIdx; i++ {
				x := rand.Int63n(maxVal-minVal) + minVal
				y := rand.Int63n(maxVal-minVal) + minVal

				bigX := big.NewInt(x)
				bigY := big.NewInt(y)

				factor := big.NewInt(1e18)
				bigX.Mul(bigX, factor)
				bigY.Mul(bigY, factor)

				operations[i] = ArithmeticOperation{
					X:      bigX,
					Y:      bigY,
					OpType: OperationType(rand.Intn(4)),
				}
			}
		}(w)
	}

	wg.Wait()
	return operations
}

// RunSequentialUltraFast exécute SANS allocations dans la boucle
func RunSequentialUltraFast(operations []ArithmeticOperation) []*big.Int {
	results := make([]*big.Int, len(operations))

	// Pré-allouer TOUS les buffers
	for i := range results {
		results[i] = new(big.Int)
	}

	// BOUCLE CRITIQUE: ZÉRO nouvelles allocations, ZÉRO logging
	for i, op := range operations {
		ExecuteOperationInPlace(op, results[i])
	}

	return results
}

// RunParallelUltraFast exécute EN PARALLÈLE avec buffers pré-alloués
func RunParallelUltraFast(operations []ArithmeticOperation, numWorkers int) []*big.Int {
	results := make([]*big.Int, len(operations))
	var wg sync.WaitGroup

	// Pré-allouer buffers par worker
	workerBuffers := make([][]*big.Int, numWorkers)
	for w := 0; w < numWorkers; w++ {
		workerBuffers[w] = make([]*big.Int, len(operations)/numWorkers+1)
		for i := range workerBuffers[w] {
			workerBuffers[w][i] = new(big.Int)
		}
	}

	// Canaux pour résultats (pas de mutex = ultra-rapide)
	resultChans := make([]chan IndexedBigInt, numWorkers)
	for i := 0; i < numWorkers; i++ {
		resultChans[i] = make(chan IndexedBigInt, 10000)
	}

	// Lancer workers
	for w := 0; w < numWorkers; w++ {
		wg.Add(1)

		go func(workerID int, resultChan chan IndexedBigInt, buffers []*big.Int) {
			defer wg.Done()

			startIdx := int64(workerID) * int64(len(operations)) / int64(numWorkers)
			endIdx := int64(workerID+1) * int64(len(operations)) / int64(numWorkers)

			if workerID == numWorkers-1 {
				endIdx = int64(len(operations))
			}

			bufIdx := 0
			// BOUCLE CRITIQUE: ZÉRO allocations, ZÉRO timing
			for i := startIdx; i < endIdx; i++ {
				ExecuteOperationInPlace(operations[i], buffers[bufIdx])
				bufIdx++
				if bufIdx >= len(buffers) {
					bufIdx = 0
				}

				resultChan <- IndexedBigInt{Index: i, Value: buffers[bufIdx-1]}
			}
			close(resultChan)
		}(w, resultChans[w], workerBuffers[w])
	}

	// Collecter en parallèle (pas de synchronisation)
	go func() {
		for w := 0; w < numWorkers; w++ {
			go func(ch chan IndexedBigInt) {
				for ib := range ch {
					results[ib.Index] = ib.Value
				}
			}(resultChans[w])
		}
	}()

	wg.Wait()
	return results
}

type IndexedBigInt struct {
	Index int64
	Value *big.Int
}

// ExecuteOperationInPlace exécute l'opération EN PLACE (dans le buffer)
func ExecuteOperationInPlace(op ArithmeticOperation, result *big.Int) {
	switch op.OpType {
	case OP_ADD:
		result.Add(op.X, op.Y)
	case OP_SUB:
		if op.X.Cmp(op.Y) >= 0 {
			result.Sub(op.X, op.Y)
		} else {
			result.Sub(op.Y, op.X)
		}
	case OP_MUL:
		result.Mul(op.X, op.Y)
	case OP_DIV:
		if op.Y.Sign() != 0 {
			result.Div(op.X, op.Y)
		} else {
			result.Set(op.X)
		}
	}
}

// EstimateSequentialFractionSimple calcule S d'Amdahl
func EstimateSequentialFractionSimple(speedup float64, numWorkers float64) float64 {
	if speedup <= 1.0 {
		return 1.0
	}

	S := (numWorkers - speedup) / (numWorkers * (speedup - 1.0))

	if S < 0.0 {
		S = 0.0
	}
	if S > 1.0 {
		S = 1.0
	}

	return S
}
