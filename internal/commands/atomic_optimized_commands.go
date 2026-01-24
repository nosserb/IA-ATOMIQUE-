package commands

import (
	"IA-ATOMIQUE/database"
	"fmt"
	"strings"
	"time"
)

// ============================================================================
// OPTIMIZED ATOMIC SIMULATION COMMANDS
// ============================================================================

// ProcessOptimizedAtomicCommand gère les commandes de simulation optimisée
func ProcessOptimizedAtomicCommand(args []string) {
	if len(args) < 2 {
		fmt.Println("Usage: ./programme atomic-optimized <command> [options]")
		fmt.Println("\nCommands:")
		fmt.Println("  simulate <iterations> <atoms> <topK>    - Run optimized simulation")
		fmt.Println("  compare <iterations> <atoms>             - Compare old vs optimized")
		fmt.Println("  benchmark <iterations> <atoms>           - Performance benchmark")
		return
	}

	command := strings.ToLower(args[1])

	switch command {
	case "simulate":
		if len(args) < 5 {
			fmt.Println("Usage: ./programme atomic-optimized simulate <iterations> <atoms> <topK>")
			return
		}
		var iterations, numAtoms, topK int
		fmt.Sscanf(args[2], "%d", &iterations)
		fmt.Sscanf(args[3], "%d", &numAtoms)
		fmt.Sscanf(args[4], "%d", &topK)

		RunOptimizedSimulation(iterations, numAtoms, topK)

	case "compare":
		if len(args) < 4 {
			fmt.Println("Usage: ./programme atomic-optimized compare <iterations> <atoms>")
			return
		}
		var iterations, numAtoms int
		fmt.Sscanf(args[2], "%d", &iterations)
		fmt.Sscanf(args[3], "%d", &numAtoms)

		CompareOptimizations(iterations, numAtoms)

	case "benchmark":
		if len(args) < 4 {
			fmt.Println("Usage: ./programme atomic-optimized benchmark <iterations> <atoms>")
			return
		}
		var iterations, numAtoms int
		fmt.Sscanf(args[2], "%d", &iterations)
		fmt.Sscanf(args[3], "%d", &numAtoms)

		RunPerformanceBenchmark(iterations, numAtoms)

	default:
		fmt.Printf("Unknown command: %s\n", command)
	}
}

// ============================================================================
// OPTIMIZED SIMULATION
// ============================================================================

// RunOptimizedSimulation lance une simulation avec les optimisations
func RunOptimizedSimulation(iterations, numAtoms, topK int) {
	fmt.Println("\n╔════════════════════════════════════════════════════════════╗")
	fmt.Println("║     ATOMIC OPTIMIZED SIMULATION - MATHEMATICAL MODEL     ║")
	fmt.Println("╚════════════════════════════════════════════════════════════╝\n")

	// Paramètres du réseau
	network := database.NewAtomicNetwork(numAtoms)

	// Configuration des atomes (voisins aléatoires pour dense network)
	for i := 0; i < numAtoms; i++ {
		network.Atoms[i].Neighbors = make([]int, 0)
		for j := 0; j < numAtoms; j++ {
			if i != j && (j < i+8 || j == i+1) { // Chaque atome a 8 voisins
				network.Atoms[i].Neighbors = append(network.Atoms[i].Neighbors, j)
				network.Atoms[i].ConnectionWeights[j] = 0.5
			}
		}
	}

	fmt.Printf("📊 Configuration:\n")
	fmt.Printf("  • Total atoms: %d\n", numAtoms)
	fmt.Printf("  • Iterations: %d\n", iterations)
	fmt.Printf("  • Top-K neighbors: %d\n", topK)
	fmt.Printf("  • α (coupling): %.2f\n", network.CouplingCoefficient)
	fmt.Printf("  • β (local rules): %.2f\n", network.LocalRulesCoefficient)
	fmt.Printf("  • γ (reinforcement): %.2f\n", network.ReinforcementFactor)
	fmt.Printf("  • δ (decay): %.2f\n", network.DecayFactor)
	fmt.Printf("  • σ (resonance sensitivity): %.2f\n\n", network.ResonanceSensitivity)

	// Simulation
	startTime := time.Now()

	for iter := 0; iter < iterations; iter++ {
		// Mise à jour vectorisée avec workers
		numWorkers := 4 // Go routines parallel
		database.VectorizedStateUpdate(
			network,
			network.CouplingCoefficient,
			network.LocalRulesCoefficient,
			network.ResonanceSensitivity,
			topK,
			numWorkers,
		)

		if iter%100 == 0 && iter > 0 {
			fmt.Printf("  Iteration %d: Energy=%.4f, Active=%d\n",
				iter, network.TotalEnergy, numAtoms-network.FrozenAtomsCount)
		}
	}

	elapsed := time.Since(startTime)

	// Résultats
	metrics := database.AnalyzePerformance(network, topK)

	fmt.Printf("\n⏱️  Processing:\n")
	fmt.Printf("  • Total time: %.2f ms\n", elapsed.Seconds()*1000)
	fmt.Printf("  • Time per iteration: %.3f ms\n", elapsed.Seconds()*1000/float64(iterations))
	fmt.Printf("  • Time per atom per iteration: %.3f µs\n", elapsed.Microseconds()/int64(iterations*numAtoms))

	fmt.Printf("\n📈 Performance Metrics:\n")
	fmt.Printf("  • Average neighbors per atom: %.1f\n", metrics.AverageNeighbors)
	fmt.Printf("  • Top-K selected: %d\n", metrics.TopKUsed)
	fmt.Printf("  • Reduction factor: %.2f%% (O(|N|) → O(k))\n", metrics.ReductionFactor*100)
	fmt.Printf("  • Estimated speedup: %.2fx\n", metrics.EstimatedSpeedup)
	fmt.Printf("  • Energy per atom: %.6f\n", metrics.EnergyPerAtom)
	fmt.Printf("  • Total energy used: %.4f\n", metrics.TotalEnergyUsed)

	// Analyse des états finaux
	fmt.Printf("\n📊 Final State Distribution:\n")
	lowState := 0
	midState := 0
	highState := 0
	for i := range network.Atoms {
		s := network.Atoms[i].InternalState
		if s < 0.33 {
			lowState++
		} else if s < 0.66 {
			midState++
		} else {
			highState++
		}
	}
	fmt.Printf("  • Low (0-0.33):   %d atoms (%.1f%%)\n", lowState, float64(lowState)*100/float64(numAtoms))
	fmt.Printf("  • Mid (0.33-0.66): %d atoms (%.1f%%)\n", midState, float64(midState)*100/float64(numAtoms))
	fmt.Printf("  • High (0.66-1.0): %d atoms (%.1f%%)\n", highState, float64(highState)*100/float64(numAtoms))

	fmt.Printf("\n✅ Simulation Complete\n")
}

// ============================================================================
// COMPARISON: Old vs Optimized
// ============================================================================

// CompareOptimizations compare l'ancienne approche avec la nouvelle
func CompareOptimizations(iterations, numAtoms int) {
	fmt.Println("\n╔════════════════════════════════════════════════════════════╗")
	fmt.Println("║     COMPARISON: Original vs Optimized Atomic Model       ║")
	fmt.Println("╚════════════════════════════════════════════════════════════╝\n")

	// Configuration réseau
	network := database.NewAtomicNetwork(numAtoms)

	// Ajouter voisins
	for i := 0; i < numAtoms; i++ {
		network.Atoms[i].Neighbors = make([]int, 0)
		for j := 0; j < numAtoms; j++ {
			if i != j && j < i+12 { // Chaque atome a ~12 voisins
				network.Atoms[i].Neighbors = append(network.Atoms[i].Neighbors, j)
				network.Atoms[i].ConnectionWeights[j] = 0.5
			}
		}
	}

	// Configuration paramètres
	avgNeighbors := float64(0)
	for i := range network.Atoms {
		avgNeighbors += float64(len(network.Atoms[i].Neighbors))
	}
	avgNeighbors /= float64(numAtoms)

	topK := 5 // Top-K sélectionné

	fmt.Printf("📊 Configuration:\n")
	fmt.Printf("  • Atoms: %d\n", numAtoms)
	fmt.Printf("  • Iterations: %d\n", iterations)
	fmt.Printf("  • Average neighbors: %.1f\n", avgNeighbors)
	fmt.Printf("  • Top-K: %d\n\n", topK)

	// --- APPROACH 1: BRUTE FORCE (Original) ---
	fmt.Printf("═══ Approach 1: BRUTE FORCE (tous les voisins) ═══\n")
	start1 := time.Now()

	network1 := database.NewAtomicNetwork(numAtoms)
	for i := 0; i < numAtoms; i++ {
		network1.Atoms[i].Neighbors = make([]int, 0)
		for j := 0; j < numAtoms; j++ {
			if i != j && j < i+12 {
				network1.Atoms[i].Neighbors = append(network1.Atoms[i].Neighbors, j)
				network1.Atoms[i].ConnectionWeights[j] = 0.5
			}
		}
	}

	for iter := 0; iter < iterations; iter++ {
		database.VectorizedStateUpdate(network1, network1.CouplingCoefficient,
			network1.LocalRulesCoefficient, network1.ResonanceSensitivity,
			1000, 1) // 1000 = pas de limitation Top-K
	}

	elapsed1 := time.Since(start1)
	fmt.Printf("  ⏱️  Time: %.2f ms\n", elapsed1.Seconds()*1000)
	fmt.Printf("  📊 Operations per atom per iteration: %.0f (all neighbors)\n", avgNeighbors)

	// --- APPROACH 2: TOP-K OPTIMIZATION ---
	fmt.Printf("\n═══ Approach 2: TOP-K OPTIMIZED (k=%d) ═══\n", topK)
	start2 := time.Now()

	network2 := database.NewAtomicNetwork(numAtoms)
	for i := 0; i < numAtoms; i++ {
		network2.Atoms[i].Neighbors = make([]int, 0)
		for j := 0; j < numAtoms; j++ {
			if i != j && j < i+12 {
				network2.Atoms[i].Neighbors = append(network2.Atoms[i].Neighbors, j)
				network2.Atoms[i].ConnectionWeights[j] = 0.5
			}
		}
	}

	for iter := 0; iter < iterations; iter++ {
		database.VectorizedStateUpdate(network2, network2.CouplingCoefficient,
			network2.LocalRulesCoefficient, network2.ResonanceSensitivity,
			topK, 4) // topK avec parallelization
	}

	elapsed2 := time.Since(start2)
	fmt.Printf("  ⏱️  Time: %.2f ms\n", elapsed2.Seconds()*1000)
	fmt.Printf("  📊 Operations per atom per iteration: %.0f (Top-K selected)\n", float64(topK))

	// Résultats comparatifs
	speedup := elapsed1.Seconds() / elapsed2.Seconds()
	reduction := (1.0 - float64(topK)/avgNeighbors) * 100

	fmt.Printf("\n📈 Comparison Results:\n")
	fmt.Printf("  • Brute Force: %.2f ms\n", elapsed1.Seconds()*1000)
	fmt.Printf("  • Top-K Opt:   %.2f ms\n", elapsed2.Seconds()*1000)
	fmt.Printf("  • Speedup:     %.2fx\n", speedup)
	fmt.Printf("  • Reduction:   %.1f%% (from %.0f → %d neighbors)\n", reduction, avgNeighbors, topK)

	if speedup > 1.0 {
		fmt.Printf("\n✅ Top-K optimization is %.2fx FASTER!\n", speedup)
	} else {
		fmt.Printf("\n⚠️  Overhead detected (possibly due to sorting). Consider increasing atoms count.\n")
	}
}

// ============================================================================
// PERFORMANCE BENCHMARK
// ============================================================================

// RunPerformanceBenchmark lance un benchmark complet
func RunPerformanceBenchmark(iterations, numAtoms int) {
	fmt.Println("\n╔════════════════════════════════════════════════════════════╗")
	fmt.Println("║          PERFORMANCE BENCHMARK - SCALABILITY STUDY       ║")
	fmt.Println("╚════════════════════════════════════════════════════════════╝\n")

	// Test avec différentes valeurs de k
	kValues := []int{2, 5, 10, 15, 20}

	fmt.Printf("Testing with %d atoms, %d iterations\n\n", numAtoms, iterations)
	fmt.Printf("┌─────┬─────────────┬──────────┬──────────┬───────────┐\n")
	fmt.Printf("│ k   │ Time (ms)   │ Time/iter│ Energy   │ Speedup   │\n")
	fmt.Printf("├─────┼─────────────┼──────────┼──────────┼───────────┤\n")

	baselineTime := time.Duration(0)

	for _, k := range kValues {
		network := database.NewAtomicNetwork(numAtoms)

		// Configurer voisins
		for i := 0; i < numAtoms; i++ {
			network.Atoms[i].Neighbors = make([]int, 0)
			for j := 0; j < numAtoms; j++ {
				if i != j && j < i+12 {
					network.Atoms[i].Neighbors = append(network.Atoms[i].Neighbors, j)
					network.Atoms[i].ConnectionWeights[j] = 0.5
				}
			}
		}

		start := time.Now()

		for iter := 0; iter < iterations; iter++ {
			database.VectorizedStateUpdate(network, network.CouplingCoefficient,
				network.LocalRulesCoefficient, network.ResonanceSensitivity,
				k, 4)
		}

		elapsed := time.Since(start)

		if k == kValues[0] {
			baselineTime = elapsed
		}

		speedup := baselineTime.Seconds() / elapsed.Seconds()

		fmt.Printf("│ %2d  │ %11.2f │ %8.4f │ %8.4f │ %9.2fx │\n",
			k,
			elapsed.Seconds()*1000,
			elapsed.Seconds()*1000/float64(iterations),
			network.TotalEnergy,
			speedup,
		)
	}

	fmt.Printf("└─────┴─────────────┴──────────┴──────────┴───────────┘\n")

	fmt.Printf("\n✅ Benchmark Complete\n")
	fmt.Printf("Recommendation: Use k=5-10 for best balance of speed and accuracy\n")
}
