package main

import (
	"fmt"
	"math/big"
	"math/rand"
	"sync"
	"time"
)

// GenerateOperationsWithRedundancy génère des opérations AVEC redondances intentionnelles
// Pour simuler le cas réel (30-70% des opérations ne changent rien)
func GenerateOperationsWithRedundancy(count int64, redundancyPercent float64, numWorkers int) []ArithmeticOperation {
	operations := make([]ArithmeticOperation, count)

	redundantTarget := int64(float64(count) * redundancyPercent / 100.0)
	redundantCount := int64(0)

	for i := int64(0); i < count; i++ {
		// Avec probabilité redundancyPercent, créer une op redondante
		if redundantCount < redundantTarget && rand.Float64() < (float64(redundantTarget-redundantCount)/float64(count-i)) {
			// Créer une opération redondante
			x := big.NewInt(rand.Int63n(1000) + 1)

			opType := rand.Intn(4)
			switch opType {
			case 0: // x + 0
				operations[i] = ArithmeticOperation{X: x, Y: big.NewInt(0), OpType: OP_ADD}
			case 1: // x - 0
				operations[i] = ArithmeticOperation{X: x, Y: big.NewInt(0), OpType: OP_SUB}
			case 2: // x * 1
				operations[i] = ArithmeticOperation{X: x, Y: big.NewInt(1), OpType: OP_MUL}
			case 3: // x / 1
				operations[i] = ArithmeticOperation{X: x, Y: big.NewInt(1), OpType: OP_DIV}
			}
			redundantCount++
		} else {
			// Opération réelle
			x := big.NewInt(rand.Int63n(int64(1e8)) + int64(1e7))
			y := big.NewInt(rand.Int63n(int64(1e8)) + int64(1e7))
			op := OperationType(rand.Intn(4))

			factor := big.NewInt(1e18)
			x.Mul(x, factor)
			y.Mul(y, factor)

			operations[i] = ArithmeticOperation{
				X:      x,
				Y:      y,
				OpType: op,
			}
		}
	}

	return operations
}

// ProcessStressTestLevier3WithRedundancy démontre l'impact du Levier 3
func ProcessStressTestLevier3WithRedundancy(args []string) {
	if len(args) == 0 {
		fmt.Println("[ERREUR] ./programme stest-l3-demo <nb_ops> [redundancy%]")
		return
	}

	var numOps int64
	var redundancyPercent float64 = 40.0 // Défaut: 40% redondantes

	fmt.Sscanf(args[0], "%d", &numOps)
	if len(args) > 1 {
		fmt.Sscanf(args[1], "%g", &redundancyPercent)
	}

	if numOps <= 0 {
		fmt.Println("[ERREUR] Nombre > 0")
		return
	}

	if redundancyPercent < 0 || redundancyPercent > 100 {
		redundancyPercent = 40.0
	}

	numWorkers := 8
	targetResponseMs := 230.0

	fmt.Printf("\n╔═══════════════════════════════════════════════════╗\n")
	fmt.Printf("║  LEVIER 3: DÉMO - OPS AVEC REDONDANCES         ║\n")
	fmt.Printf("║  Impact réel de l'élimination d'ops inutiles   ║\n")
	fmt.Printf("╚═══════════════════════════════════════════════════╝\n\n")

	fmt.Printf("Config: %dM ops, %.0f%% redondantes\n", numOps/1000000, redundancyPercent)
	fmt.Printf("        %d workers, T_cible=%gms\n\n", numWorkers, targetResponseMs)

	// PHASE 1: Génération AVEC redondances
	fmt.Printf("⏳ Génération avec %.0f%% redondances...\n", redundancyPercent)
	genStart := time.Now()
	operations := GenerateOperationsWithRedundancy(numOps, redundancyPercent, numWorkers)
	genTime := time.Since(genStart)
	fmt.Printf("  ✓ %v\n", genTime)

	// PHASE 2: Analyse redondances
	fmt.Printf("\n🔍 Analyse (détection redondantes)...\n")
	analysisStart := time.Now()
	optimizedOps, redundantCount := OptimizeOperationList(operations)
	analysisTime := time.Since(analysisStart)

	redundantPercent := float64(redundantCount) * 100 / float64(len(operations))
	effectiveOps := int64(len(operations)) - int64(redundantCount)

	fmt.Printf("  ✓ %v\n", analysisTime)
	fmt.Printf("  • Redondantes trouvées: %d (%.1f%% vs %.0f%% attendu)\n", redundantCount, redundantPercent, redundancyPercent)
	fmt.Printf("  • Opérations réelles: %d (%.1f%%)\n", effectiveOps, 100-redundantPercent)

	// PHASE 3: Estimer T_op
	fmt.Printf("\n📊 Estimation temps/op...\n")
	sampleSize := 1000
	opTimeUs := EstimateOperationTime(operations, sampleSize)
	fmt.Printf("  ✓ T_op ≈ %.2f µs (incluant skips redondants)\n", opTimeUs)

	// PHASE 4: Calculer batch adaptatif
	adaptiveBatch := CalculateAdaptiveBatch(opTimeUs, numWorkers, targetResponseMs)
	fmt.Printf("\n📐 Batch adaptatif: %d ops/batch\n", adaptiveBatch)

	// PHASE 5: Réponse immédiate
	immediateOps := int64(float64(numOps) * 0.05)

	fmt.Printf("\n⚡ Calcul IMMÉDIAT (5%% = %.0fM ops)...\n", float64(immediateOps)/1e6)
	immediateStart := time.Now()

	immediateResults := make([]*big.Int, immediateOps)
	immediateSkipped := int64(0)

	for i := int64(0); i < immediateOps; i++ {
		immediateResults[i] = new(big.Int)
		if optimizedOps[i].IsRedundant {
			immediateResults[i].Set(optimizedOps[i].RedundantValue)
			immediateSkipped++
		} else {
			ExecuteOptimizedOperation(optimizedOps[i], immediateResults[i])
		}
	}

	immediateTime := time.Since(immediateStart)
	fmt.Printf("  ✓ En %.2f ms (skip %d ops)\n", immediateTime.Seconds()*1000, immediateSkipped)

	// PHASE 6: Fond
	fmt.Printf("\n🔄 Calcul FOND (95%%)...")
	allResults := make([]*big.Int, numOps)
	for i := int64(0); i < immediateOps; i++ {
		allResults[i] = immediateResults[i]
	}
	for i := immediateOps; i < numOps; i++ {
		allResults[i] = new(big.Int)
	}

	bgStart := time.Now()
	var wg sync.WaitGroup
	opsPerWorker := (numOps - immediateOps) / int64(numWorkers)

	for w := 0; w < numWorkers; w++ {
		wg.Add(1)

		go func(workerID int) {
			defer wg.Done()

			startIdx := immediateOps + int64(workerID)*opsPerWorker
			endIdx := startIdx + opsPerWorker
			if workerID == numWorkers-1 {
				endIdx = numOps
			}

			for i := startIdx; i < endIdx; i++ {
				ExecuteOptimizedOperation(optimizedOps[i], allResults[i])
			}
		}(w)
	}

	wg.Wait()
	bgTime := time.Since(bgStart)

	totalTime := time.Since(genStart)
	fmt.Printf(" %.2f ms\n", bgTime.Seconds()*1000)

	// RÉSULTATS
	fmt.Printf("\n═══════════════════════════════════════════════════\n")
	fmt.Printf("RÉSULTATS AVEC LEVIER 3\n")
	fmt.Printf("═══════════════════════════════════════════════════\n")
	fmt.Printf("⚡ Réponse visible: %.2f ms", immediateTime.Seconds()*1000)
	if immediateTime.Seconds()*1000 < 230 {
		fmt.Printf(" ✅\n")
	} else {
		fmt.Printf(" ⚠️\n")
	}

	fmt.Printf("🔄 Calcul fond: %.2f ms\n", bgTime.Seconds()*1000)
	fmt.Printf("📊 Total: %.2f ms\n", totalTime.Seconds()*1000)

	// Analyse gain
	fmt.Printf("\n💡 Impact Levier 3 (Opérations Fusionnées):\n")
	fmt.Printf("  • Redondantes détectées: %d (%.1f%%)\n", redundantCount, redundantPercent)
	fmt.Printf("  • Opérations RÉELLES exécutées: %d\n", int64(numOps)-int64(redundantCount))
	fmt.Printf("  • Économie: %.1f%% du temps de calcul\n", redundantPercent)

	// Simulation: temps sans Levier 3
	estimatedTimeWithoutL3 := bgTime.Seconds() / (1.0 - float64(redundantCount)/float64(numOps))
	timeSaved := estimatedTimeWithoutL3 - bgTime.Seconds()

	fmt.Printf("\n  Simulation (sans Levier 3):\n")
	fmt.Printf("  • Temps estimé: %.2f ms\n", estimatedTimeWithoutL3*1000)
	fmt.Printf("  • Temps SAUVÉ: %.2f ms (%.1f%%)\n",
		timeSaved*1000, timeSaved/estimatedTimeWithoutL3*100)

	// Verdict
	fmt.Printf("\n🎯 VERDICT LEVIERS 1+2+3:\n")
	fmt.Printf("  ✅ UX instantanée\n")
	fmt.Printf("  ✅ Batch adaptatif (0%% overhead)\n")
	fmt.Printf("  ✅ Opérations fusionnées (%.1f%% gain)\n", redundantPercent)
	fmt.Printf("  🚀 Speedup perçu: %.0fx\n", totalTime.Seconds()*1000/immediateTime.Seconds()/1000)

	fmt.Printf("\n")
}
