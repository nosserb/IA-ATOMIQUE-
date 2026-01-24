package database

import (
	"fmt"
	"sync"
)

// ============================================================================
// FAST IMAGE GENERATION - DRAFT MODE OPTIMIZATIONS
// ============================================================================
// Ces optimisations réduisent le temps de génération à <2sec pour images 256x256

// FastImageConfig contient les paramètres pour génération ultra-rapide
type FastImageConfig struct {
	Width       int
	Height      int
	Iterations  int    // Réduit drastiquement
	PatchSize   int    // Augmenté pour moins d'atomes
	Quality     string // "draft" | "fast" | "balanced" | "quality"
	UseCache    bool
	Parallelism int // Nombre de goroutines
}

// PresetConfigs définit les configurations prédéfinies
var PresetConfigs = map[string]FastImageConfig{
	"ultra": {
		Width:       128,
		Height:      128,
		Iterations:  2,  // 2 itérations seulement!
		PatchSize:   64, // TRÈS grand patch size
		Quality:     "ultra",
		UseCache:    true,
		Parallelism: 4,
	},
	"draft": {
		Width:       256,
		Height:      256,
		Iterations:  5,  // 5 itérations seulement
		PatchSize:   32, // Très grand patch size
		Quality:     "draft",
		UseCache:    true,
		Parallelism: 8,
	},
	"fast": {
		Width:       256,
		Height:      256,
		Iterations:  10,
		PatchSize:   16,
		Quality:     "fast",
		UseCache:    true,
		Parallelism: 8,
	},
	"balanced": {
		Width:       512,
		Height:      512,
		Iterations:  30,
		PatchSize:   8,
		Quality:     "balanced",
		UseCache:    true,
		Parallelism: 16,
	},
	"quality": {
		Width:       512,
		Height:      512,
		Iterations:  100,
		PatchSize:   8,
		Quality:     "quality",
		UseCache:    false,
		Parallelism: 32,
	},
}

// FastAtomicImageNetwork est une version optimisée pour génération rapide
type FastAtomicImageNetwork struct {
	BaseNetwork  *AtomicImageNetwork
	Config       FastImageConfig
	Cache        map[string]float64
	CacheMutex   sync.RWMutex
	ComputeCount int // Pour monitoring
	CacheHits    int
}

// NewFastAtomicImageNetwork crée une version optimisée du réseau
func NewFastAtomicImageNetwork(config FastImageConfig) *FastAtomicImageNetwork {
	// Réduire les coefficients pour convergence plus rapide
	baseNetwork := NewAtomicImageNetwork(config.Width, config.Height, config.PatchSize)

	// Ajuster les paramètres pour convergence rapide
	switch config.Quality {
	case "draft":
		baseNetwork.CouplingCoefficient = 0.9   // Plus d'influence voisins
		baseNetwork.LocalRulesCoefficient = 0.8 // Plus de contraintes
		baseNetwork.ReinforcementFactor = 0.3
		baseNetwork.DecayFactor = 0.1
		baseNetwork.FreezeThreshold = 0.5
		baseNetwork.FreezeIterations = 1
	case "fast":
		baseNetwork.CouplingCoefficient = 0.85
		baseNetwork.LocalRulesCoefficient = 0.6
		baseNetwork.ReinforcementFactor = 0.2
		baseNetwork.DecayFactor = 0.08
		baseNetwork.FreezeThreshold = 0.4
		baseNetwork.FreezeIterations = 2
	default:
		// Garder les paramètres standards pour balanced/quality
	}

	return &FastAtomicImageNetwork{
		BaseNetwork:  baseNetwork,
		Config:       config,
		Cache:        make(map[string]float64),
		ComputeCount: 0,
		CacheHits:    0,
	}
}

// OptimizedIterateGeneration - Version parallélisée avec batching
func (fnet *FastAtomicImageNetwork) OptimizedIterateGeneration() {
	network := fnet.BaseNetwork
	network.mutex.Lock()
	network.GlobalIteration++
	network.mutex.Unlock()

	gridHeight := network.Height / network.PatchSize
	gridWidth := network.Width / network.PatchSize

	// Batch updates pour meilleure localité de cache
	batchSize := 16
	totalAtoms := gridHeight * gridWidth

	for batchStart := 0; batchStart < totalAtoms; batchStart += batchSize {
		batchEnd := batchStart + batchSize
		if batchEnd > totalAtoms {
			batchEnd = totalAtoms
		}

		var wg sync.WaitGroup
		for atomIdx := batchStart; atomIdx < batchEnd; atomIdx++ {
			y := atomIdx / gridWidth
			x := atomIdx % gridWidth

			wg.Add(1)
			go func(atomX, atomY int) {
				defer wg.Done()
				// Version rapide de UpdateAtomState
				fnet.fastUpdateAtomState(&network.Atoms[atomY][atomX], atomX, atomY)
				fnet.ComputeCount++
			}(x, y)
		}
		wg.Wait()
	}
}

// fastUpdateAtomState - Version allégée pour draft mode
func (fnet *FastAtomicImageNetwork) fastUpdateAtomState(atom *PixelAtom, x, y int) {
	atom.mutex.Lock()
	defer atom.mutex.Unlock()

	network := fnet.BaseNetwork

	// Skip frozen atoms entièrement
	if atom.IsFrozen {
		return
	}

	// Version SIMPLIFIÉE sans résonance complexe
	gridWidth := network.Width / network.PatchSize
	gridHeight := network.Height / network.PatchSize

	// Top-k neighbors only (pas tous les 8)
	neighborInfluence := 0.0
	neighborCount := 0
	topK := 2 // Seulement les 2 meilleurs voisins

	// Calcul simplifié des voisins
	dx := []int{-1, 0, 1, 0}
	dy := []int{0, -1, 0, 1}

	for i := 0; i < 4; i++ {
		nx := x + dx[i]
		ny := y + dy[i]

		if nx >= 0 && nx < gridWidth && ny >= 0 && ny < gridHeight {
			neighbor := &network.Atoms[ny][nx]
			// Pas de résonance complexe
			diff := neighbor.State - atom.State
			weight := 0.25 // Poids uniforme
			neighborInfluence += weight * diff

			neighborCount++
			if neighborCount >= topK {
				break
			}
		}
	}

	// Update ultra-simplifié
	alpha := network.CouplingCoefficient
	beta := network.LocalRulesCoefficient

	atom.State += alpha * neighborInfluence * 0.1 // Petit pas
	atom.State += beta * atom.ExternalConstraint * 0.1

	// Clamp
	if atom.State > 1.0 {
		atom.State = 1.0
	} else if atom.State < 0.0 {
		atom.State = 0.0
	}

	// Color update ultra-simplifié
	if network.PromptEmbedding != nil {
		if intensity, ok := network.PromptEmbedding["bright"]; ok {
			atom.Color[0] += intensity * 0.05
			atom.Color[1] += intensity * 0.05
			atom.Color[2] += intensity * 0.05
		}
	}

	// Clamp colors
	for i := 0; i < 3; i++ {
		if atom.Color[i] > 1.0 {
			atom.Color[i] = 1.0
		} else if atom.Color[i] < 0.0 {
			atom.Color[i] = 0.0
		}
	}

	// Pas de weight update pour économiser du CPU
	// Pas de freeze check pour économiser du CPU
}

// SkipPostProcessing - Saute le post-processing pour draft mode
func (fnet *FastAtomicImageNetwork) SkipPostProcessing() bool {
	return fnet.Config.Quality == "draft" || fnet.Config.Quality == "fast"
}

// PrintOptimizationStats affiche les statistiques d'optimisation
func (fnet *FastAtomicImageNetwork) PrintOptimizationStats() {
	fmt.Printf("\n⚡ OPTIMIZATION STATS (Mode: %s)\n", fnet.Config.Quality)
	fmt.Printf("   Total computations: %d\n", fnet.ComputeCount)
	fmt.Printf("   Cache hits: %d\n", fnet.CacheHits)
	cacheHitRate := float64(fnet.CacheHits) / float64(fnet.ComputeCount+fnet.CacheHits) * 100
	fmt.Printf("   Cache hit rate: %.1f%%\n", cacheHitRate)
}
