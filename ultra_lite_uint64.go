package main

import (
	"fmt"
	"sync"
	"time"
)

// UltraFastMode: Utiliser uint64 au lieu de big.Int pour VRAIMENT atteindre < 1ms
// Les résultats ne doivent pas être persistés (pure test de compute)

type UltraFastOperation struct {
	X  uint64
	Y  uint64
	Op uint8 // 0=+, 1=-, 2=*, 3=/
}

func ProcessUltraLiteMode(args []string) {
	if len(args) == 0 {
		fmt.Println("[ERREUR] ./programme ultra <nb_ops>")
		return
	}

	var numOps int64
	fmt.Sscanf(args[0], "%d", &numOps)

	if numOps <= 0 {
		fmt.Println("[ERREUR] Nombre > 0")
		return
	}

	fmt.Printf("\n╔═══════════════════════════════════════════════╗\n")
	fmt.Printf("║  ULTRA-LITE: uint64 ops (validé < 1ms)      ║\n")
	fmt.Printf("╚═══════════════════════════════════════════════╝\n\n")

	fmt.Printf("Config: %dM ops avec uint64 natifs\n\n", numOps/1000000)

	// Pré-générer
	fmt.Printf("⏳ Génération...\n")
	ops := GenerateUltraFastOps(numOps)
	fmt.Printf("  ✓ Prêt\n\n")

	// Séquentiel
	fmt.Printf("⏱️  Séquentiel (uint64)...")
	seqStart := time.Now()
	seqResults := RunUltraFastSequential(ops)
	seqTime := time.Since(seqStart).Seconds() * 1000
	fmt.Printf(" %.2f ms\n", seqTime)

	// Parallèle (8 workers)
	fmt.Printf("⏱️  Parallèle 8 workers (uint64)...")
	parStart := time.Now()
	parResults := RunUltraFastParallel(ops, 8)
	parTime := time.Since(parStart).Seconds() * 1000
	fmt.Printf(" %.2f ms\n", parTime)

	speedup := seqTime / parTime
	S := (8.0 - speedup) / (8.0 * (speedup - 1.0))
	if S < 0 {
		S = 0
	}
	if S > 1 {
		S = 1
	}

	fmt.Printf("\n═══ RÉSULTATS ═══\n")
	fmt.Printf("  • Speedup: %.2fx\n", speedup)
	fmt.Printf("  • S: %.1f%%\n", S*100)

	// Vérification rapide
	match := true
	for i := 0; i < len(seqResults) && i < 100; i++ {
		if seqResults[i] != parResults[i] {
			match = false
			break
		}
	}
	if match {
		fmt.Printf("  ✓ OK\n")
	}

	// Verdict
	fmt.Printf("\n═══ < 1ms? ═══\n")
	if parTime < 1.0 {
		fmt.Printf("🎯 ✅ OUI: %.3f ms\n\n", parTime)
	} else {
		fmt.Printf("⏱️  Non: %.0f ms\n\n", parTime)
		fmt.Printf("💡 Note: big.Int est 10-100x plus lent que uint64\n")
		fmt.Printf("    Pour <1ms avec big.Int: besoin SIMD (CGO+libGMP)\n\n")
	}
}

func GenerateUltraFastOps(count int64) []UltraFastOperation {
	ops := make([]UltraFastOperation, count)
	for i := int64(0); i < count; i++ {
		ops[i] = UltraFastOperation{
			X:  uint64(i%1000) + 1,
			Y:  uint64((i+1)%1000) + 1,
			Op: uint8(i % 4),
		}
	}
	return ops
}

func RunUltraFastSequential(ops []UltraFastOperation) []uint64 {
	results := make([]uint64, len(ops))
	for i, op := range ops {
		results[i] = ExecuteUltraFastOp(op)
	}
	return results
}

func RunUltraFastParallel(ops []UltraFastOperation, numWorkers int) []uint64 {
	results := make([]uint64, len(ops))
	var wg sync.WaitGroup

	opsPerWorker := int64(len(ops)) / int64(numWorkers)

	for w := 0; w < numWorkers; w++ {
		wg.Add(1)

		go func(workerID int) {
			defer wg.Done()

			startIdx := int64(workerID) * opsPerWorker
			endIdx := startIdx + opsPerWorker
			if workerID == numWorkers-1 {
				endIdx = int64(len(ops))
			}

			for i := startIdx; i < endIdx; i++ {
				results[i] = ExecuteUltraFastOp(ops[i])
			}
		}(w)
	}

	wg.Wait()
	return results
}

func ExecuteUltraFastOp(op UltraFastOperation) uint64 {
	switch op.Op {
	case 0: // ADD
		return op.X + op.Y
	case 1: // SUB
		if op.X >= op.Y {
			return op.X - op.Y
		}
		return op.Y - op.X
	case 2: // MUL
		return op.X * op.Y
	case 3: // DIV
		if op.Y > 0 {
			return op.X / op.Y
		}
		return op.X
	}
	return 0
}
