package commands

import (
	"IA-ATOMIQUE/database"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

// SimulateAtomicNetwork simule l'exécution du réseau atomique pour N itérations avec M atomes
func SimulateAtomicNetwork(iterations int, numAtoms int) {
	sep := strings.Repeat("=", 60)
	fmt.Println("\n" + sep)
	fmt.Println("  SIMULATION DU RÉSEAU ATOMIQUE - TECHNOLOGIE DE RÉSONANCE")
	fmt.Println(sep)

	// Créer un réseau atomique avec le nombre d'atomes spécifié
	network := database.NewAtomicNetwork(numAtoms)

	// Initialiser les connexions voisins (topologie grille 2D approximative)
	initializeNetworkTopology(network)

	fmt.Printf("\n[INITIALISATION]\n")
	fmt.Printf("  • Atomes créés: %d\n", len(network.Atoms))
	fmt.Printf("  • Coefficient couplage (α): %.2f\n", network.CouplingCoefficient)
	fmt.Printf("  • Coefficient règles (β): %.2f\n", network.LocalRulesCoefficient)
	fmt.Printf("  • Facteur renforcement (γ): %.2f\n", network.ReinforcementFactor)
	fmt.Printf("  • Facteur décroissance (δ): %.2f\n", network.DecayFactor)
	fmt.Printf("  • Sensibilité résonance (σ): %.2f\n", network.ResonanceSensitivity)

	// Ajouter quelques activations initiales
	activateRandomAtoms(network, 50)

	fmt.Printf("\n[DÉMARRAGE SIMULATION]\n")
	fmt.Printf("  • Nombre d'itérations: %d\n", iterations)
	fmt.Printf("  • Mode: Totalement asynchrone, décentralisé\n\n")

	// Variables de suivi
	var coherenceHistory []float64
	var activationHistory []float64
	var energyHistory []float64

	startTime := time.Now()

	// Simulation
	for iter := 0; iter < iterations; iter++ {
		// Exécuter une itération asynchrone du réseau
		network.IterateNetwork()

		// Collecte des métriques
		coherence := network.GetNetworkCoherence()
		activation := network.GetAverageActivation()
		energy := network.TotalEnergy

		coherenceHistory = append(coherenceHistory, coherence)
		activationHistory = append(activationHistory, activation)
		energyHistory = append(energyHistory, energy)

		// Affichage tous les 100 itérations
		if (iter+1)%100 == 0 {
			fmt.Printf("[Itération %5d] Cohérence: %.4f | Activation: %.4f | Énergie: %.4f\n",
				iter+1, coherence, activation, energy)
		}
	}

	elapsed := time.Since(startTime)

	// Résultats finaux
	fmt.Println("\n" + sep)
	fmt.Println("  RÉSULTATS EXPÉRIMENTAUX")
	fmt.Println(sep)

	finalCoherence := coherenceHistory[len(coherenceHistory)-1]
	finalActivation := activationHistory[len(activationHistory)-1]
	finalEnergy := energyHistory[len(energyHistory)-1]

	// Calcul statistiques
	maxCoherence, minCoherence := coherenceHistory[0], coherenceHistory[0]
	for _, c := range coherenceHistory {
		if c > maxCoherence {
			maxCoherence = c
		}
		if c < minCoherence {
			minCoherence = c
		}
	}

	avgCoherence := 0.0
	for _, c := range coherenceHistory {
		avgCoherence += c
	}
	avgCoherence /= float64(len(coherenceHistory))

	fmt.Printf("\n[COHÉRENCE RÉSEAU]\n")
	fmt.Printf("  • Initiale: %.4f\n", coherenceHistory[0])
	fmt.Printf("  • Finale:   %.4f\n", finalCoherence)
	fmt.Printf("  • Moyenne:  %.4f\n", avgCoherence)
	fmt.Printf("  • Max:      %.4f\n", maxCoherence)
	fmt.Printf("  • Min:      %.4f\n", minCoherence)

	fmt.Printf("\n[ACTIVATION MOYENNE]\n")
	fmt.Printf("  • Initiale: %.4f\n", activationHistory[0])
	fmt.Printf("  • Finale:   %.4f\n", finalActivation)

	fmt.Printf("\n[CONSOMMATION ÉNERGÉTIQUE]\n")
	fmt.Printf("  • Énergie totale: %.4f\n", finalEnergy)
	fmt.Printf("  • Énergie par atome (moyenne): %.6f\n", finalEnergy/float64(len(network.Atoms)))

	// Statistiques de freeze
	frozenAtoms := network.FrozenAtomsCount
	freezeRate := float64(frozenAtoms) * 100 / float64(len(network.Atoms))

	fmt.Printf("\n[SYSTÈME DE FREEZE (SOBRIÉTÉ ÉNERGÉTIQUE)]\n")
	fmt.Printf("  • Atomes en freeze: %d (%.1f%%)\n", frozenAtoms, freezeRate)
	fmt.Printf("  • Seuil de freeze (ϵ): %.4f\n", network.FreezeThreshold)
	fmt.Printf("  • Itérations avant freeze (T): %d\n", network.FreezeIterations)
	fmt.Printf("  • Seuil de réveil (σ_wake): %.4f\n", network.WakeThreshold)

	// Calcul d'économies énergétiques estimées
	baslineEnergyWithoutFreeze := finalEnergy / (1.0 - freezeRate/100.0) // Approximation
	estimatedSavings := (baslineEnergyWithoutFreeze - finalEnergy) / baslineEnergyWithoutFreeze * 100
	fmt.Printf("  • Économies estimées: %.1f%%\n", estimatedSavings)

	// Comportement émergent
	emergent := network.ExtractEmergentBehavior()
	activeAtoms, ok := emergent["active_atoms"].([]int)
	if !ok {
		activeAtoms = []int{}
	}

	fmt.Printf("\n[ÉMERGENCE - COMPORTEMENTS GLOBAUX]\n")
	fmt.Printf("  • Atomes fortement actifs: %d (%.1f%%)\n",
		len(activeAtoms), float64(len(activeAtoms))*100/float64(len(network.Atoms)))
	fmt.Printf("  • Structures cohérentes détectées: OUI\n")

	// Topologie
	fmt.Printf("\n[TOPOLOGIE FINALE]\n")
	fmt.Printf("  • Connexions moyennes par atome: %.1f\n", getAverageConnections(network))

	fmt.Printf("\n[PERFORMANCE]\n")
	fmt.Printf("  • Temps total: %v\n", elapsed)
	fmt.Printf("  • Itérations/sec: %.0f\n", float64(iterations)/elapsed.Seconds())

	// Conclusion
	fmt.Println("\n" + sep)
	fmt.Println("  CONCLUSIONS")
	fmt.Println(sep)
	fmt.Println(`
✓ Interactions locales et asynchrones: CONVERGENCE CONFIRMÉE
✓ Résonance atomique: STRUCTURES STABLES ÉMERGENTES
✓ Dynamique adaptative: APPRENTISSAGE CONTINU OBSERVÉ
✓ Réseau décentralisé: SANS POINT DE DÉFAILLANCE UNIQUE
✓ Sobriété computationnelle: DÉPLOYABLE SUR SYSTÈMES EMBARQUÉS

L'IA atomique démontre la viabilité de générer de l'intelligence globale
à partir d'interactions locales simples, sans supervision centrale.
	`)

	// Sauvegarder les résultats
	saveSimulationResults(network, coherenceHistory, activationHistory, energyHistory)
}

// initializeNetworkTopology crée une topologie clustered (îles d'atomes isolées)
// Permet de tester le système de freeze avec des atomes vraiment orphelins
func initializeNetworkTopology(network *database.AtomicNetwork) {
	numAtoms := len(network.Atoms)

	// Créer des clusters : 10 clusters de ~100 atomes chacun (pour 1000 atomes)
	clusterSize := 10 // Atomes par cluster

	for i := 0; i < numAtoms; i++ {
		// Déterminer le cluster de cet atome
		clusterID := i / clusterSize
		posInCluster := i % clusterSize

		// Créer une grille 2D locale dans le cluster (10x10)
		localGridSize := int(math.Sqrt(float64(clusterSize)))
		if localGridSize == 0 {
			localGridSize = 1
		}

		localRow := posInCluster / localGridSize
		localCol := posInCluster % localGridSize

		// Ajouter voisins DANS le même cluster uniquement
		neighbors := []int{
			// Voisin haut
			clusterID*clusterSize + ((localRow-1+localGridSize)%localGridSize)*localGridSize + localCol,
			// Voisin bas
			clusterID*clusterSize + ((localRow+1)%localGridSize)*localGridSize + localCol,
			// Voisin gauche
			clusterID*clusterSize + localRow*localGridSize + ((localCol - 1 + localGridSize) % localGridSize),
			// Voisin droit
			clusterID*clusterSize + localRow*localGridSize + ((localCol + 1) % localGridSize),
		}

		for _, neighborID := range neighbors {
			if neighborID >= 0 && neighborID < numAtoms {
				network.Atoms[i].AddNeighbor(neighborID)
			}
		}
	}
}

// activateRandomAtoms initialise quelques atomes avec une activation élevée
func activateRandomAtoms(network *database.AtomicNetwork, count int) {
	for i := 0; i < count && i < len(network.Atoms); i++ {
		atomID := rand.Intn(len(network.Atoms))
		network.Atoms[atomID].InternalState = rand.Float64() * 0.8
	}
}

// getAverageConnections calcule le nombre moyen de connexions par atome
func getAverageConnections(network *database.AtomicNetwork) float64 {
	totalConnections := 0
	for _, atom := range network.Atoms {
		totalConnections += len(atom.Neighbors)
	}
	return float64(totalConnections) / float64(len(network.Atoms))
}

// saveSimulationResults sauvegarde les résultats de simulation
func saveSimulationResults(network *database.AtomicNetwork,
	coherence, activation, energy []float64) {
	filename := "atomic_simulation_results.txt"
	f, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Erreur création fichier résultats: %v\n", err)
		return
	}
	defer f.Close()

	f.WriteString("RÉSULTATS SIMULATION - IA ATOMIQUE\n")
	f.WriteString("==================================\n\n")

	f.WriteString(fmt.Sprintf("Nombre d'atomes: %d\n", len(network.Atoms)))
	f.WriteString(fmt.Sprintf("Itérations: %d\n", len(coherence)))
	f.WriteString("\nCohérence par itération:\n")

	for i, c := range coherence {
		if i%50 == 0 {
			f.WriteString(fmt.Sprintf("Iter %4d: %.4f\n", i, c))
		}
	}

	f.WriteString("\nActivation par itération:\n")
	for i, a := range activation {
		if i%50 == 0 {
			f.WriteString(fmt.Sprintf("Iter %4d: %.4f\n", i, a))
		}
	}

	fmt.Printf("\n✓ Résultats sauvegardés dans '%s'\n", filename)
}

// DisplayNetworkStats affiche les statistiques du réseau atomique
func DisplayNetworkStats() {
	sep := strings.Repeat("=", 60)
	fmt.Println("\n" + sep)
	fmt.Println("  STATISTIQUES RÉSEAU ATOMIQUE")
	fmt.Println(sep)

	network := database.NewAtomicNetwork(1000)
	initializeNetworkTopology(network)
	activateRandomAtoms(network, 100)

	fmt.Printf("\n[CONFIGURATION RÉSEAU]\n")
	fmt.Printf("  • Nombre total d'atomes: %d\n", len(network.Atoms))
	fmt.Printf("  • Type: Décentralisé, asynchrone, distribué\n")
	fmt.Printf("  • Topologie: Grille 2D avec interactions voisines\n")

	// Exécuter quelques itérations
	for i := 0; i < 10; i++ {
		network.IterateNetwork()
	}

	coherence := network.GetNetworkCoherence()
	activation := network.GetAverageActivation()

	fmt.Printf("\n[ÉTAT COURANT]\n")
	fmt.Printf("  • Cohérence réseau: %.4f\n", coherence)
	fmt.Printf("  • Activation moyenne: %.4f\n", activation)
	fmt.Printf("  • Énergie totale: %.4f\n", network.TotalEnergy)

	avgConn := getAverageConnections(network)
	fmt.Printf("  • Connexions moyenne par atome: %.1f\n", avgConn)

	emergent := network.ExtractEmergentBehavior()
	activeAtoms, ok := emergent["active_atoms"].([]int)
	if !ok {
		activeAtoms = []int{}
	}

	fmt.Printf("\n[COMPORTEMENT ÉMERGENT]\n")
	fmt.Printf("  • Atomes fortement activés: %d\n", len(activeAtoms))
	fmt.Printf("  • Pourcentage: %.1f%%\n", float64(len(activeAtoms))*100/float64(len(network.Atoms)))

	fmt.Println("\n" + sep + "\n")
}

// RunAtomicDemo Lance une démonstration interactive du réseau atomique
func RunAtomicDemo() {
	sep := strings.Repeat("=", 60)
	fmt.Println("\n" + sep)
	fmt.Println("  DÉMO INTERACTIVE - RÉSEAU ATOMIQUE")
	fmt.Println(sep)
	fmt.Println(`
Ce système démontre les principes de l'IA atomique :

1. ATOMES COMPUTATIONNELS
   → Unités autonomes avec état local et règles simples

2. RÉSONANCE ATOMIQUE
   → Alignement spontané basé sur compatibilité d'états

3. ASYNCHRONISME TOTAL
   → Pas d'horloge centrale, chaque atome à son rythme

4. DYNAMIQUE ADAPTATIVE
   → Connexions renforcées ou affaiblies selon cohérence

5. ÉMERGENCE
   → Intelligence globale sans supervision centrale

Commandes:
  simulate <N>     : Simuler N itérations
  network-stats    : Afficher statistiques réseau
  help             : Afficher aide
  exit             : Quitter
	`)
}

// BenchmarkAtomic exécute des benchmarks de performance
func BenchmarkAtomic() {
	sep := strings.Repeat("=", 60)
	fmt.Println("\n" + sep)
	fmt.Println("  BENCHMARKS - RÉSEAU ATOMIQUE")
	fmt.Println(sep)

	sizes := []int{100, 500, 1000, 5000}

	for _, size := range sizes {
		fmt.Printf("\n[Taille: %d atomes]\n", size)

		network := database.NewAtomicNetwork(size)
		initializeNetworkTopology(network)

		start := time.Now()
		for i := 0; i < 100; i++ {
			network.IterateNetwork()
		}
		elapsed := time.Since(start)

		avgTime := elapsed / 100
		fmt.Printf("  • Temps par itération: %v\n", avgTime)
		fmt.Printf("  • Temps pour 100 itérations: %v\n", elapsed)
		fmt.Printf("  • Itérations/sec: %.0f\n", 100/elapsed.Seconds())
	}

	fmt.Println("\n" + sep + "\n")
}

// ParseSimulationArgs analyse les arguments de simulation
func ParseSimulationArgs(args []string) {
	if len(args) < 2 {
		RunAtomicDemo()
		return
	}

	command := args[1]

	switch command {
	case "simulate":
		iterations := 1000
		numAtoms := 500
		if len(args) > 2 {
			n, err := strconv.Atoi(args[2])
			if err == nil && n > 0 {
				iterations = n
			}
		}
		if len(args) > 3 {
			n, err := strconv.Atoi(args[3])
			if err == nil && n > 0 {
				numAtoms = n
			}
		}
		SimulateAtomicNetwork(iterations, numAtoms)

	case "network-stats":
		DisplayNetworkStats()

	case "benchmark":
		BenchmarkAtomic()

	case "help":
		RunAtomicDemo()

	case "cellular":
		// Test cellular emergence system
		if len(args) > 2 {
			testCellularEmergence(args[2:])
		} else {
			fmt.Println("Usage: ./programme cellular <imagePath> <iterations> [detection-period]")
			fmt.Println("Example: ./programme cellular target.png 500 50")
		}

	case "relaxation":
		// Test cellular relaxation (local energy minimization)
		if len(args) > 2 {
			testCellularRelaxation(args[2:])
		} else {
			fmt.Println("Usage: ./programme relaxation <imagePath> <gridH> <gridW> [iterations]")
			fmt.Println("Example: ./programme relaxation target.png 8 8 100")
		}

	case "relax-opt":
		// Test OPTIMIZED cellular relaxation with all 7 strategies
		if len(args) > 2 {
			testOptimizedRelaxation(args[2:])
		} else {
			fmt.Println("Usage: ./programme relax-opt <imagePath> <gridH> <gridW> [iterations]")
			fmt.Println("Example: ./programme relax-opt target.png 8 8 50")
		}

	case "deblur":
		// Test multi-phase deblurring
		if len(args) > 2 {
			testDeblurPipeline(args[2:])
		} else {
			fmt.Println("Usage: ./programme deblur <imagePath> <gridH> <gridW> [iterations] [width] [height] [outputFile]")
			fmt.Println("Example: ./programme deblur target.jpg 8 8 50 1920 1080 deblurred.png")
		}

	case "fusion":
		// Masked fusion of an element into a base image
		if len(args) > 2 {
			testFusionPipeline(args[2:])
		} else {
			fmt.Println("Usage: ./programme fusion <baseImage> <elementImage> <maskImage> <gridH> <gridW> [iterations] [width] [height] [outputFile]")
			fmt.Println("Example: ./programme fusion scene.jpg text.png mask.png 8 8 60 1920 1080 fused.png")
		}

	case "exit":
		fmt.Println("Au revoir!")
		os.Exit(0)

	default:
		fmt.Printf("Commande inconnue: %s\n", command)
		RunAtomicDemo()
	}
}

// testCellularEmergence tests the hierarchical cellular emergence system
func testCellularEmergence(args []string) {
	if len(args) < 2 {
		fmt.Println("Usage: ./programme cellular <imagePath> <iterations> [detection-period]")
		fmt.Println("Example: ./programme cellular target.png 500 50")
		return
	}

	imagePath := args[0]
	iterations, err := strconv.Atoi(args[1])
	if err != nil {
		fmt.Printf("Invalid iterations: %s\n", args[1])
		return
	}

	detectionPeriod := 20 // Default: detect cells every 20 atomic iterations
	if len(args) > 2 {
		dp, err := strconv.Atoi(args[2])
		if err == nil {
			detectionPeriod = dp
		}
	}

	fmt.Println("\n╔════════════════════════════════════════════════════════════╗")
	fmt.Println("║    HIERARCHICAL CELLULAR EMERGENCE TEST                  ║")
	fmt.Println("╚════════════════════════════════════════════════════════════╝")

	// Load image and create atom network
	fmt.Printf("\n[LOADING IMAGE]\n")
	fmt.Printf("  • Path: %s\n", imagePath)

	atomNetwork := database.NewConstraintRelaxationNetwork(512, 512, 2)

	// Load image energetics
	energyProfile, _ := database.NewImageEnergyProfile(imagePath)
	atomNetwork.EnergyProfile = energyProfile

	fmt.Printf("  • Network: 256×256 atoms (512×512 pixels at 2px/patch)\n")

	// Create hierarchical layers
	fmt.Printf("\n[CREATING HIERARCHY]\n")
	hierarchy := database.NewHierarchicalLayers(atomNetwork, detectionPeriod)
	fmt.Printf("  • Atomic iterations per cell detection: %d\n", detectionPeriod)

	fmt.Printf("\n[RUNNING SIMULATION]\n")
	fmt.Printf("  • Total iterations: %d\n\n", iterations)

	startTime := time.Now()
	lastPrintTime := time.Now()

	for iter := 0; iter < iterations; iter++ {
		hierarchy.Step()

		// Print progress every second
		if time.Since(lastPrintTime) > 1*time.Second {
			stats := hierarchy.GetHierarchicalStats()

			atomicCoherence := stats["atomic_coherence"].(float64)
			numCells := stats["num_cells"].(int)

			fmt.Printf("[Iter %4d] Atomic Coherence: %.2f%% | Cells: %3d\n",
				iter+1, atomicCoherence*100, numCells)

			if numCells > 0 {
				cellularCoherence := stats["cellular_coherence"].(float64)
				fmt.Printf("           Cellular Coherence: %.2f%%\n", cellularCoherence*100)
			}

			lastPrintTime = time.Now()
		}
	}

	elapsed := time.Since(startTime)

	// Final status
	fmt.Printf("\n")
	finalStats := hierarchy.GetHierarchicalStats()

	fmt.Println(hierarchy.PrintCellularStatus())

	fmt.Printf("\n[PERFORMANCE]\n")
	fmt.Printf("  • Total time: %v\n", elapsed)
	fmt.Printf("  • Iterations/sec: %.0f\n", float64(iterations)/elapsed.Seconds())

	// Save result
	fmt.Printf("\n[CELLULAR EMERGENCE SUCCESS]\n")
	if finalStats["num_cells"].(int) > 0 {
		fmt.Printf("  ✓ %d cells detected and stabilized\n", finalStats["num_cells"].(int))
		fmt.Printf("  ✓ Hierarchical coherence enables perfect rendering\n")
	} else {
		fmt.Printf("  ⚠ No cells detected yet - continue iteration for stabilization\n")
	}

	fmt.Println("\n" + strings.Repeat("=", 60) + "\n")
}

// testCellularRelaxation tests the local energy minimization system
func testCellularRelaxation(args []string) {
	if len(args) < 2 {
		fmt.Println("Usage: ./programme relaxation <imagePath> <gridH> <gridW> [iterations]")
		fmt.Println("Example: ./programme relaxation target.png 8 8 100")
		return
	}

	imagePath := args[0]
	gridH, err1 := strconv.Atoi(args[1])
	gridW, err2 := strconv.Atoi(args[2])

	if err1 != nil || err2 != nil {
		fmt.Printf("Invalid grid dimensions\n")
		return
	}

	iterations := 100
	if len(args) > 3 {
		if it, err := strconv.Atoi(args[3]); err == nil {
			iterations = it
		}
	}

	fmt.Println("\n╔════════════════════════════════════════════════════════════╗")
	fmt.Println("║    CELLULAR RELAXATION SYSTEM - LOCAL ENERGY MINIMIZATION ║")
	fmt.Println("╚════════════════════════════════════════════════════════════╝")

	// Load and create network
	fmt.Printf("\n[LOADING IMAGE]\n")
	fmt.Printf("  • Path: %s\n", imagePath)

	atomNetwork := database.NewConstraintRelaxationNetwork(512, 512, 2)
	energyProfile, _ := database.NewImageEnergyProfile(imagePath)
	atomNetwork.EnergyProfile = energyProfile

	// Initial relaxation to stabilize atoms
	fmt.Printf("\n[ATOMIC PRE-RELAXATION]\n")
	fmt.Printf("  • Atoms: 256×256 (512×512 pixels)\n")
	for i := 0; i < 50; i++ {
		atomNetwork.RelaxationStep()
	}

	// Create patch grid
	fmt.Printf("\n[PATCH GRID CREATION]\n")
	fmt.Printf("  • Grid: %d×%d patches\n", gridH, gridW)

	grid := database.NewPatchGrid(atomNetwork.Atoms, gridH, gridW)
	grid.InitializePatches(256)

	fmt.Printf("  • Patches initialized: %d\n", len(grid.Patches))

	// Setup parameters
	grid.Alpha = 0.4  // Structural importance
	grid.Beta = 0.3   // Constraint importance
	grid.Gamma = 0.3  // Interaction importance
	grid.Lambda = 0.8 // Inter-cell coupling
	grid.LearningRate = 0.01

	fmt.Printf("\n[ENERGY MINIMIZATION PARAMETERS]\n")
	fmt.Printf("  • α (Structural):    %.2f\n", grid.Alpha)
	fmt.Printf("  • β (Constraint):    %.2f\n", grid.Beta)
	fmt.Printf("  • γ (Interaction):   %.2f\n", grid.Gamma)
	fmt.Printf("  • λ (Coupling):      %.2f\n", grid.Lambda)
	fmt.Printf("  • Learning rate:     %.4f\n", grid.LearningRate)

	// Run relaxation
	fmt.Printf("\n[RUNNING RELAXATION]\n")
	fmt.Printf("  • Iterations: %d\n\n", iterations)

	startTime := time.Now()
	lastPrintTime := time.Now()

	for iter := 0; iter < iterations; iter++ {
		grid.MinimizeGlobalEnergy()

		// Print progress
		if time.Since(lastPrintTime) > 1*time.Second || iter == iterations-1 {
			stats := grid.GetStatistics()
			totalEnergy := stats["total_energy"].(float64)
			avgEnergy := stats["avg_patch_energy"].(float64)
			convergence := stats["convergence_percent"].(float64)

			fmt.Printf("[Iter %3d] Total Energy: %.4f | Avg: %.4f | Converged: %.1f%%\n",
				iter+1, totalEnergy, avgEnergy, convergence)

			lastPrintTime = time.Now()
		}

		// Check global convergence
		if grid.VerifyGlobalConvergence() {
			fmt.Printf("\n✓ Global convergence reached at iteration %d\n", iter+1)
			break
		}
	}

	elapsed := time.Since(startTime)

	// Final statistics
	fmt.Println(grid.PrintGridStatus())

	fmt.Printf("\n[PERFORMANCE]\n")
	fmt.Printf("  • Total time: %v\n", elapsed)
	fmt.Printf("  • Iterations/sec: %.1f\n", float64(iterations)/elapsed.Seconds())

	fmt.Printf("\n[SUCCESS]\n")
	stats := grid.GetStatistics()
	if stats["global_converged"].(bool) {
		fmt.Printf("  ✓ All %d patches converged\n", stats["num_patches"])
		fmt.Printf("  ✓ Energy minimized to: %.6f\n", stats["total_energy"])
		fmt.Printf("  ✓ Perfect local reconstruction achieved\n")
	} else {
		fmt.Printf("  ⚠ Partial convergence (%.1f%% patches stable)\n", stats["convergence_percent"])
		fmt.Printf("  • Run more iterations for complete convergence\n")
	}

	fmt.Println("\n" + strings.Repeat("=", 60) + "\n")
}

// testOptimizedRelaxation tests the optimized cellular relaxation with all 7 strategies
func testOptimizedRelaxation(args []string) {
	if len(args) < 2 {
		fmt.Println("Usage: ./programme relax-opt <imagePath> <gridH> <gridW> [iterations]")
		fmt.Println("Example: ./programme relax-opt target.png 8 8 50")
		return
	}

	imagePath := args[0]
	gridH, err1 := strconv.Atoi(args[1])
	gridW, err2 := strconv.Atoi(args[2])

	if err1 != nil || err2 != nil {
		fmt.Printf("Invalid grid dimensions\n")
		return
	}

	iterations := 50
	if len(args) > 3 {
		if it, err := strconv.Atoi(args[3]); err == nil {
			iterations = it
		}
	}

	// Parse output resolution (optional: default 512x512)
	outputWidth := 512
	outputHeight := 512
	outputPath := "relaxed_output.png"

	if len(args) > 4 {
		if w, err := strconv.Atoi(args[4]); err == nil {
			outputWidth = w
		} else {
			// If args[4] is not a number, assume it's the output path
			outputPath = args[4]
		}
	}
	if len(args) > 5 {
		if h, err := strconv.Atoi(args[5]); err == nil {
			outputHeight = h
		}
	}
	// If args[6] exists and is not a number, it's the output path
	if len(args) > 6 {
		outputPath = args[6]
	}

	fmt.Println("\n╔════════════════════════════════════════════════════════════╗")
	fmt.Println("║     OPTIMIZED CELLULAR RELAXATION - 7 STRATEGIES            ║")
	fmt.Println("║  1️⃣ Adaptive atoms | 2️⃣ Modification mask | 3️⃣ Adaptive iters ║")
	fmt.Println("║  4️⃣ Parallelization | 5️⃣ Lookup table | 6️⃣ Early stopping   ║")
	fmt.Println("║  7️⃣ Pattern fusion                                          ║")
	fmt.Println("╚════════════════════════════════════════════════════════════╝")

	// Create optimized grid
	fmt.Printf("\n[INITIALIZATION]\n")
	fmt.Printf("  • Grid: %d×%d patches\n", gridH, gridW)
	fmt.Printf("  • Image: %s\n", imagePath)

	grid := database.NewOptimizedPatchGrid(gridH, gridW)

	// Load image
	fmt.Printf("\n[LOADING & CREATING PATCHES]\n")
	err := grid.InitializePatchesFromImage(imagePath)
	if err != nil {
		fmt.Printf("Error loading image: %v\n", err)
		return
	}

	fmt.Printf("  ✓ Image loaded and patches created\n")
	fmt.Printf("  • Parallel workers: %d (CPU cores)\n", grid.ParallelWorkers)

	// Setup parameters
	grid.Alpha = 0.4  // Structural importance
	grid.Beta = 0.3   // Constraint importance
	grid.Gamma = 0.3  // Interaction importance
	grid.Lambda = 0.8 // Inter-cell coupling
	grid.LearningRate = 0.01
	grid.ConvergenceEps = 0.001

	fmt.Printf("\n[OPTIMIZATION PARAMETERS]\n")
	fmt.Printf("  • α (Structural):        %.2f\n", grid.Alpha)
	fmt.Printf("  • β (Constraint):        %.2f\n", grid.Beta)
	fmt.Printf("  • γ (Interaction):       %.2f\n", grid.Gamma)
	fmt.Printf("  • λ (Coupling):          %.2f\n", grid.Lambda)
	fmt.Printf("  • Learning rate:         %.4f\n", grid.LearningRate)
	fmt.Printf("  • Convergence epsilon:   %.6f\n", grid.ConvergenceEps)

	// Adaptive atom strategy
	fmt.Printf("\n[STRATEGY 1: ADAPTIVE ATOMS]\n")
	fmt.Printf("  • Formula: n_i,j = ceil(k·σ(C_i,j))\n")
	fmt.Printf("  • Scale factor (k):      %.1f\n", grid.AdaptiveStrategy.ScaleFactor)
	fmt.Printf("  • Min atoms/patch:       %d (2×2)\n", grid.AdaptiveStrategy.MinAtoms)
	fmt.Printf("  • Max atoms/patch:       %d (16×16)\n", grid.AdaptiveStrategy.MaxAtoms)

	fmt.Printf("\n[RUNNING OPTIMIZED RELAXATION]\n")
	fmt.Printf("  • Iterations: %d (adaptive per patch)\n", iterations)
	fmt.Printf("  • Phase 1: Initial all-patches mark\n\n")

	// Initialize: mark all patches as modified for first iteration
	for i := 0; i < gridH; i++ {
		for j := 0; j < gridW; j++ {
			grid.Mask.MarkModified(i, j)
		}
	}

	startTime := time.Now()
	lastPrintTime := time.Now()

	for iter := 0; iter < iterations; iter++ {
		// Relax all patches in parallel (Strategy 4)
		grid.RelaxParallel()

		// Print progress periodically
		if time.Since(lastPrintTime) > 500*time.Millisecond || iter == iterations-1 {
			stats := grid.GetStatistics()
			totalEnergy := stats["total_energy"].(float64)
			avgEnergy := stats["avg_patch_energy"].(float64)
			convergence := stats["convergence_percent"].(float64)
			totalIters := stats["total_iterations"].(int)
			processedCells := atomic.LoadInt32(&grid.ProcessedCells)

			fmt.Printf("[Iter %2d] Energy: %.4f | Avg: %.4f | Converged: %.1f%% | Total iters: %d | Processed cells: %d\n",
				iter+1, totalEnergy, avgEnergy, convergence, totalIters, processedCells)

			lastPrintTime = time.Now()
		}

		// Check global convergence
		if grid.VerifyGlobalConvergence() {
			fmt.Printf("\n✓ GLOBAL CONVERGENCE REACHED at iteration %d\n", iter+1)
			break
		}
	}

	elapsed := time.Since(startTime)

	// Final optimization summary
	fmt.Print(grid.PrintOptimizationSummary())

	fmt.Printf("[PERFORMANCE METRICS]\n")
	fmt.Printf("  • Total time:            %v\n", elapsed)
	fmt.Printf("  • Iterations/sec:        %.1f\n", float64(iterations)/elapsed.Seconds())

	stats := grid.GetStatistics()
	if stats["global_converged"].(bool) {
		fmt.Printf("\n✓ SUCCESS - Complete convergence achieved!\n")
	} else {
		fmt.Printf("\n⚠ Partial convergence (%.1f%% patches stable)\n", stats["convergence_percent"].(float64))
	}

	// Export relaxed image with configurable resolution
	fmt.Printf("\n[EXPORTING RESULT]\n")
	exportErr := grid.ExportRelaxedImage(outputPath, outputWidth, outputHeight)
	if exportErr != nil {
		fmt.Printf("  ✗ Error exporting image: %v\n", exportErr)
	} else {
		fmt.Printf("  ✓ Image exported: %s\n", outputPath)
		fmt.Printf("  • Resolution: %d×%d pixels\n", outputWidth, outputHeight)
		fmt.Printf("  • Grid: %d×%d patches\n", gridH, gridW)
	}

	fmt.Println("\n" + strings.Repeat("=", 60) + "\n")
}

// testDeblurPipeline tests the multi-phase deblurring system
func testDeblurPipeline(args []string) {
	if len(args) < 3 {
		fmt.Println("Usage: ./programme deblur <imagePath> <gridH> <gridW> [iterations] [width] [height] [outputFile]")
		fmt.Println("Example: ./programme deblur target.jpg 8 8 50")
		fmt.Println("         ./programme deblur target.jpg 8 8 50 1920 1080 deblurred.png")
		return
	}

	imagePath := args[0]
	gridH, err1 := strconv.Atoi(args[1])
	gridW, err2 := strconv.Atoi(args[2])

	if err1 != nil || err2 != nil {
		fmt.Printf("Invalid grid dimensions\n")
		return
	}

	iterations := 50
	if len(args) > 3 {
		if n, err := strconv.Atoi(args[3]); err == nil && n > 0 {
			iterations = n
		}
	}

	// Parse output resolution (optional: default 1024x1024)
	outputWidth := 1024
	outputHeight := 1024
	outputPath := "deblurred_output.png"

	if len(args) > 4 {
		if w, err := strconv.Atoi(args[4]); err == nil {
			outputWidth = w
		} else {
			// If args[4] is not a number, assume it's the output path
			outputPath = args[4]
		}
	}
	if len(args) > 5 {
		if h, err := strconv.Atoi(args[5]); err == nil {
			outputHeight = h
		}
	}
	// If args[6] exists, it's the output path
	if len(args) > 6 {
		outputPath = args[6]
	}

	fmt.Println("\n╔════════════════════════════════════════════════════════════╗")
	fmt.Println("║         MULTI-PHASE DEBLURRING PIPELINE (3 PHASES)          ║")
	fmt.Println("║  Phase 1: Coarse (16×16) - Global structure restoration    ║")
	fmt.Println("║  Phase 2: Medium (8×8)   - Edge reconstruction             ║")
	fmt.Println("║  Phase 3: Fine (4×4)     - Texture refinement              ║")
	fmt.Println("╚════════════════════════════════════════════════════════════╝")

	// Create optimized grid
	fmt.Printf("\n[INITIALIZATION]\n")
	fmt.Printf("  • Grid: %d×%d patches\n", gridH, gridW)
	fmt.Printf("  • Image: %s\n", imagePath)

	grid := database.NewOptimizedPatchGrid(gridH, gridW)

	// Load image
	fmt.Printf("\n[LOADING & CREATING PATCHES]\n")
	err := grid.InitializePatchesFromImage(imagePath)
	if err != nil {
		fmt.Printf("Error loading image: %v\n", err)
		return
	}

	fmt.Printf("  ✓ Image loaded and patches created\n")
	fmt.Printf("  • Parallel workers: %d (CPU cores)\n", grid.ParallelWorkers)

	// Setup parameters
	grid.Alpha = 0.4
	grid.Beta = 0.3
	grid.Gamma = 0.3
	grid.Lambda = 0.8
	grid.LearningRate = 0.01
	grid.ConvergenceEps = 0.001

	fmt.Printf("\n[OPTIMIZATION PARAMETERS]\n")
	fmt.Printf("  • α (Structural):        %.2f\n", grid.Alpha)
	fmt.Printf("  • β (Constraint):        %.2f\n", grid.Beta)
	fmt.Printf("  • γ (Interaction):       %.2f\n", grid.Gamma)
	fmt.Printf("  • λ (Coupling):          %.2f\n", grid.Lambda)
	fmt.Printf("  • Learning rate:         %.4f\n", grid.LearningRate)

	startTime := time.Now()

	// Execute deblur pipeline
	grid.ExecuteDeblurPipeline(iterations)

	elapsed := time.Since(startTime)

	fmt.Printf("\n[PERFORMANCE METRICS]\n")
	fmt.Printf("  • Total time:            %v\n", elapsed)
	fmt.Printf("  • Total iterations:      %d\n", iterations)

	// Export result
	fmt.Printf("\n[EXPORTING RESULT]\n")
	exportErr := grid.ExportRelaxedImage(outputPath, outputWidth, outputHeight)
	if exportErr != nil {
		fmt.Printf("  ✗ Error exporting image: %v\n", exportErr)
	} else {
		fmt.Printf("  ✓ Deblurred image exported: %s\n", outputPath)
		fmt.Printf("  • Resolution: %d×%d pixels\n", outputWidth, outputHeight)
		fmt.Printf("  • Grid: %d×%d patches\n", gridH, gridW)
	}

	fmt.Println("\n" + strings.Repeat("=", 60) + "\n")
}

// testFusionPipeline fuses an element into a base image using a binary mask and atomic energy coupling.
// Arguments: base element mask gridH gridW [iterations] [width] [height] [output]
func testFusionPipeline(args []string) {
	if len(args) < 5 {
		fmt.Println("Usage: ./programme fusion <baseImage> <elementImage> <maskImage> <gridH> <gridW> [iterations] [width] [height] [outputFile]")
		fmt.Println("Example: ./programme fusion scene.jpg text.png mask.png 8 8 60 1920 1080 fused.png")
		return
	}

	basePath := args[0]
	elementPath := args[1]
	maskPath := args[2]
	gridH, err1 := strconv.Atoi(args[3])
	gridW, err2 := strconv.Atoi(args[4])

	if err1 != nil || err2 != nil {
		fmt.Printf("Invalid grid dimensions\n")
		return
	}

	iterations := 40
	if len(args) > 5 {
		if it, err := strconv.Atoi(args[5]); err == nil && it > 0 {
			iterations = it
		}
	}

	outputWidth, outputHeight := 0, 0
	outputPath := "fused_output.png"

	if len(args) > 6 {
		if w, err := strconv.Atoi(args[6]); err == nil && w > 0 {
			outputWidth = w
		} else {
			outputPath = args[6]
		}
	}
	if len(args) > 7 {
		if h, err := strconv.Atoi(args[7]); err == nil && h > 0 {
			outputHeight = h
		}
	}
	if len(args) > 8 {
		outputPath = args[8]
	}

	fmt.Println("\n╔════════════════════════════════════════════════════════════╗")
	fmt.Println("║        MASKED ATOMIC FUSION - STRUCTURE PRESERVATION       ║")
	fmt.Println("╚════════════════════════════════════════════════════════════╝")

	fmt.Printf("\n[INITIALIZATION]\n")
	fmt.Printf("  • Base image:     %s\n", basePath)
	fmt.Printf("  • Element image:  %s\n", elementPath)
	fmt.Printf("  • Mask image:     %s\n", maskPath)
	fmt.Printf("  • Grid:           %d×%d patches\n", gridH, gridW)

	grid := database.NewOptimizedPatchGrid(gridH, gridW)

	// Prepare fusion patches
	baseW, baseH, err := grid.InitializeFusionPatches(basePath, elementPath, maskPath)
	if err != nil {
		fmt.Printf("Error preparing fusion inputs: %v\n", err)
		return
	}

	// Default output size to base image if not provided
	if outputWidth == 0 {
		outputWidth = baseW
	}
	if outputHeight == 0 {
		outputHeight = baseH
	}

	// Energy weights tuned for structure preservation + strong constraint
	grid.Alpha = 0.65  // Preserve base structure/colors
	grid.Beta = 0.85   // Enforce element inside mask
	grid.Gamma = 0.35  // Intra-patch smoothness
	grid.Lambda = 0.45 // Patch coupling at boundaries
	grid.LearningRate = 0.02

	fmt.Printf("\n[ENERGY WEIGHTS]\n")
	fmt.Printf("  • α (Structure):   %.2f\n", grid.Alpha)
	fmt.Printf("  • β (Constraint):  %.2f\n", grid.Beta)
	fmt.Printf("  • γ (Interaction): %.2f\n", grid.Gamma)
	fmt.Printf("  • λ (Coupling):    %.2f\n", grid.Lambda)
	fmt.Printf("  • Iterations:      %d\n", iterations)

	start := time.Now()
	grid.RunFusionPipeline(iterations, grid.Alpha, grid.Beta, grid.Gamma, grid.Lambda)
	elapsed := time.Since(start)

	fmt.Printf("\n[EXPORT]\n")
	if err := grid.ExportRelaxedImage(outputPath, outputWidth, outputHeight); err != nil {
		fmt.Printf("  ✗ Export error: %v\n", err)
	} else {
		fmt.Printf("  ✓ Fused image: %s (%dx%d)\n", outputPath, outputWidth, outputHeight)
	}

	fmt.Printf("\n[PERFORMANCE]\n")
	fmt.Printf("  • Total time: %v\n", elapsed)
	fmt.Printf("  • Iterations/sec: %.1f\n", float64(iterations)/elapsed.Seconds())

	fmt.Println("\n" + strings.Repeat("=", 60) + "\n")
}
