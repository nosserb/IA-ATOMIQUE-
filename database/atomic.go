// Package database implements the Atomic Resonance Technology (T.R.A.)
// IA-ATOMIQUE: An Asynchronous Inference Engine Based on Atomic Resonance Technology
//
// This file implements the core atomic computational units and their local interactions
// as described in the academic paper "IA atomique: un moteur d'inférence asynchrone
// fondé sur la Technologie de Résonance Atomique (T.R.A.)"

package database

import (
	"math"
	"math/rand"
	"sync"
)

// ComputationalAtom represents an elementary unit in the Atomic Resonance Network
// Each atom is autonomous and interacts only with its direct neighbors through resonance
type ComputationalAtom struct {
	ID                    int                // Unique identifier
	InternalState         float64            // si - Internal state of the atom
	LocalRules            map[string]float64 // Ri - Local rules governing behavior
	Perceptions           map[int]float64    // pi - Perceptions of local environment
	Neighbors             []int              // List of neighboring atom IDs
	ConnectionWeights     map[int]float64    // wij - Weights of connections to neighbors
	LastUpdateTime        int64              // Timestamp of last update (asynchronous)
	EnergyConsumption     float64            // Energy used by this atom
	IsFrozen              bool               // Is this atom in freeze state
	FreezeIterationCount  int                // Count of low-activity iterations
	LastActivationLevel   float64            // Previous activation level (for Δa)
	LastLocalCoherence    float64            // Previous local coherence (for Δc)
	LastEnergyConsumption float64            // Previous energy consumption (for ΔE)
	ActivityWeights       map[string]float64 // Weights for activity calculation
	mutex                 sync.Mutex         // For thread-safe operations
}

// AtomicNetwork represents the entire distributed network of computational atoms
type AtomicNetwork struct {
	Atoms                 []ComputationalAtom
	CouplingCoefficient   float64      // α - Influence of neighbors on atom state
	LocalRulesCoefficient float64      // β - Impact of local rules and perceptions
	ReinforcementFactor   float64      // γ - Strength of connection reinforcement
	DecayFactor           float64      // δ - Decay of weak connections
	ResonanceSensitivity  float64      // σ - Sensitivity of resonance mechanism
	FreezeThreshold       float64      // ϵ - Activity threshold for freeze
	FreezeIterations      int          // T - Consecutive low-activity iterations before freeze
	WakeThreshold         float64      // σ_wake - Resonance threshold to wake up
	GlobalIteration       int          // Current global iteration counter
	TotalEnergy           float64      // Total energy consumed by network
	FrozenAtomsCount      int          // Number of currently frozen atoms
	mutex                 sync.RWMutex // For thread-safe network updates
}

// NewComputationalAtom creates a new autonomous atom unit
func NewComputationalAtom(id int, numNeighbors int) *ComputationalAtom {
	atom := &ComputationalAtom{
		ID:                    id,
		InternalState:         rand.Float64() * 0.5, // Initialize with low state
		LocalRules:            make(map[string]float64),
		Perceptions:           make(map[int]float64),
		Neighbors:             make([]int, 0, numNeighbors),
		ConnectionWeights:     make(map[int]float64),
		EnergyConsumption:     0.0,
		IsFrozen:              false,
		FreezeIterationCount:  0,
		LastActivationLevel:   0.25, // Initialize to half of max state
		LastLocalCoherence:    0.5,  // Start at middle value
		LastEnergyConsumption: 0.0,
		ActivityWeights: map[string]float64{
			"activation": 0.35,
			"coherence":  0.35,
			"resonance":  0.05, // Minimal weight for resonance
			"energy":     0.25,
		}}
	return atom
}

// NewAtomicNetwork initializes a new distributed atomic resonance network
func NewAtomicNetwork(numAtoms int) *AtomicNetwork {
	network := &AtomicNetwork{
		Atoms:                 make([]ComputationalAtom, numAtoms),
		CouplingCoefficient:   0.7,  // α - High influence from neighbors
		LocalRulesCoefficient: 0.3,  // β - Lower weight for local rules
		ReinforcementFactor:   0.15, // γ - Moderate reinforcement
		DecayFactor:           0.05, // δ - Gradual decay of weak connections
		ResonanceSensitivity:  0.1,  // σ - ULTRA-REDUCED for freeze activation (was 0.8, then 0.3)
		FreezeThreshold:       0.3,  // ϵ - Resonance threshold for freeze (ADJUSTED for new sigma)
		FreezeIterations:      2,    // T - Iterations before freeze (minimal)
		WakeThreshold:         0.40, // σ_wake - Resonance to wake up
		GlobalIteration:       0,
		TotalEnergy:           0.0,
		FrozenAtomsCount:      0,
	}

	// Initialize all atoms
	for i := 0; i < numAtoms; i++ {
		network.Atoms[i] = *NewComputationalAtom(i, 8)
	}

	return network
}

// AddNeighbor adds a neighboring atom connection (bidirectional)
func (atom *ComputationalAtom) AddNeighbor(neighborID int) {
	atom.mutex.Lock()
	defer atom.mutex.Unlock()

	// Check if already a neighbor
	for _, n := range atom.Neighbors {
		if n == neighborID {
			return
		}
	}

	atom.Neighbors = append(atom.Neighbors, neighborID)
	// Initialize connection weight randomly
	atom.ConnectionWeights[neighborID] = rand.Float64() * 0.5
}

// ComputeResonance calculates the resonance between this atom and a neighbor
// Based on equation: R(si, sj) = exp(-||si - sj||^2 / 2σ^2)
// This measures alignment/compatibility between two atomic states
func (atom *ComputationalAtom) ComputeResonance(neighborState float64, sigma float64) float64 {
	if sigma <= 0 {
		sigma = 0.1
	}

	// Calculate Euclidean distance between states
	diff := atom.InternalState - neighborState
	distanceSquared := diff * diff

	// Apply exponential decay based on distance
	exponent := -distanceSquared / (2 * sigma * sigma)
	resonance := math.Exp(exponent)

	// Clamp to [0, 1]
	if resonance > 1.0 {
		resonance = 1.0
	} else if resonance < 0.0 {
		resonance = 0.0
	}

	return resonance
}

// UpdateState performs asynchronous state update based on local resonance and rules
// Implements: si(t+1) = si(t) + α * Σ(wij * Rij) + β * (Ri + pi)
// Si frozen: n'update que faiblement pour économiser l'énergie
func (atom *ComputationalAtom) UpdateState(neighbors map[int]float64, alpha float64, beta float64, sigma float64) {
	atom.mutex.Lock()
	defer atom.mutex.Unlock()

	// Si l'atome est gelé, mettre à jour très faiblement seulement
	if atom.IsFrozen {
		// Decay très lent de l'état pour écouter les voisins actifs
		atom.InternalState *= 0.98
		atom.EnergyConsumption += 0.01 // Consommation minimale
		atom.LastActivationLevel = atom.InternalState
		atom.LastEnergyConsumption = 0.01
		return
	}

	// Phase 1: Resonance-based alignment with neighbors
	resonanceInfluence := 0.0
	neighborCount := 0

	for neighborID, neighborState := range neighbors {
		resonance := atom.ComputeResonance(neighborState, sigma)
		weight := atom.ConnectionWeights[neighborID]

		// Accumulate weighted resonance influence
		resonanceInfluence += weight * resonance * (neighborState - atom.InternalState)
		neighborCount++
	}

	// Average resonance influence
	if neighborCount > 0 {
		resonanceInfluence /= float64(neighborCount)
	}

	// Phase 2: Apply local rules and perceptions
	localInfluence := 0.0
	if rules, ok := atom.LocalRules["activation"]; ok {
		localInfluence += rules
	}
	if perception, ok := atom.Perceptions[1]; ok {
		localInfluence += perception * 0.5
	}

	// Phase 3: Update internal state
	oldState := atom.InternalState
	atom.InternalState += alpha*resonanceInfluence + beta*localInfluence

	// Clamp state to valid range [0, 1]
	if atom.InternalState > 1.0 {
		atom.InternalState = 1.0
	} else if atom.InternalState < 0.0 {
		atom.InternalState = 0.0
	}

	// Energy consumption is proportional to state change
	atom.EnergyConsumption += math.Abs(atom.InternalState - oldState)

	// Save values for next iteration's activity calculation
	atom.LastActivationLevel = atom.InternalState
	atom.LastEnergyConsumption = math.Abs(atom.InternalState - oldState)
}

// UpdateConnections implements adaptive weight dynamics
// Based on equation: dwij/dt = γ * cohesion(si,sj) - δ * wij
// This reinforces coherent interactions and weakens unstable ones
// When frozen, learning is disabled (gamma = 0)
func (atom *ComputationalAtom) UpdateConnections(neighbors map[int]float64, gamma float64, delta float64) {
	atom.mutex.Lock()
	defer atom.mutex.Unlock()

	// When frozen, disable learning (gamma becomes 0)
	effectiveGamma := gamma
	if atom.IsFrozen {
		effectiveGamma = 0.0
	}

	for neighborID, neighborState := range neighbors {
		if weight, exists := atom.ConnectionWeights[neighborID]; exists {
			// Measure coherence between states
			coherence := 1.0 - math.Abs(atom.InternalState-neighborState)
			if coherence < 0 {
				coherence = 0
			}

			// Update weight: reinforce coherent connections, decay weak ones
			// When frozen, only decay occurs (learning disabled)
			deltaW := effectiveGamma*coherence - delta*weight
			atom.ConnectionWeights[neighborID] = weight + deltaW

			// Keep weights in reasonable bounds [0, 2]
			if atom.ConnectionWeights[neighborID] > 2.0 {
				atom.ConnectionWeights[neighborID] = 2.0
			} else if atom.ConnectionWeights[neighborID] < 0.0 {
				atom.ConnectionWeights[neighborID] = 0.0
			}
		}
	}

	// Calculate and save average local coherence for next iteration
	var totalCoherence float64
	if len(neighbors) > 0 {
		for _, neighborState := range neighbors {
			coherence := 1.0 - math.Abs(atom.InternalState-neighborState)
			if coherence < 0 {
				coherence = 0
			}
			totalCoherence += coherence
		}
		atom.LastLocalCoherence = totalCoherence / float64(len(neighbors))
	} else {
		atom.LastLocalCoherence = 0.0
	}
}

// PerceiveEnvironment updates the atom's perception of its local environment
func (atom *ComputationalAtom) PerceiveEnvironment(signal float64, sourceID int) {
	atom.mutex.Lock()
	defer atom.mutex.Unlock()

	atom.Perceptions[sourceID] = signal
}

// GetState safely reads the current internal state
func (atom *ComputationalAtom) GetState() float64 {
	atom.mutex.Lock()
	defer atom.mutex.Unlock()

	return atom.InternalState
}

// ComputeLocalActivity calcule l'indice d'activité locale d'un atome
// Alocal = wa·|Δa| + wc·|Δc| + wR·R + wE·|ΔE|
func (atom *ComputationalAtom) ComputeLocalActivity(neighborStates map[int]float64, receivedResonance float64) float64 {
	atom.mutex.Lock()
	defer atom.mutex.Unlock()

	// Variation d'activation: |Δa|
	deltaActivation := math.Abs(atom.InternalState - atom.LastActivationLevel)

	// Variation de cohérence locale: moyenne des différences avec voisins
	var totalCoherenceDiff float64
	neighborCount := len(neighborStates)
	if neighborCount > 0 {
		for _, neighborState := range neighborStates {
			totalCoherenceDiff += math.Abs(atom.InternalState - neighborState)
		}
		totalCoherenceDiff /= float64(neighborCount)
	}
	deltaCoherence := math.Abs(totalCoherenceDiff - atom.LastLocalCoherence)

	// Résonance perçue: R (moyenne des signaux reçus)
	resonancePerceived := receivedResonance

	// Variation d'énergie: |ΔE|
	deltaEnergy := math.Abs(atom.EnergyConsumption - atom.LastEnergyConsumption)

	// Normaliser toutes les valeurs
	maxVariation := 1.0
	deltA := deltaActivation / maxVariation
	deltC := deltaCoherence / maxVariation
	deltE := deltaEnergy / maxVariation

	// Calculer Alocal avec poids
	wa := atom.ActivityWeights["activation"]
	wc := atom.ActivityWeights["coherence"]
	wR := atom.ActivityWeights["resonance"]
	wE := atom.ActivityWeights["energy"]

	aLocal := wa*deltA + wc*deltC + wR*resonancePerceived + wE*deltE

	// Clamp to [0, 1]
	if aLocal > 1.0 {
		aLocal = 1.0
	} else if aLocal < 0.0 {
		aLocal = 0.0
	}

	return aLocal
}

// UpdateFreezeState gère l'entrée/sortie de l'état freeze
// Freeze basé uniquement sur la résonance reçue (isolation du réseau)
func (atom *ComputationalAtom) UpdateFreezeState(localActivity float64, freezeThreshold float64, freezeIterations int, wakeThreshold float64, receivedResonance float64) {
	atom.mutex.Lock()
	defer atom.mutex.Unlock()

	if atom.IsFrozen {
		// Condition de réveil: résonance reçue > seuil de réveil
		if receivedResonance > wakeThreshold {
			atom.IsFrozen = false
			atom.FreezeIterationCount = 0
		}
	} else {
		// Comptage des itérations sans résonance significative
		// Freeze si l'atome n'est pas en contact significatif avec ses voisins
		if receivedResonance < freezeThreshold {
			atom.FreezeIterationCount++
			// Freeze après T itérations consecutives
			if atom.FreezeIterationCount >= freezeIterations {
				atom.IsFrozen = true
				atom.FreezeIterationCount = 0
			}
		} else {
			// Réinitialiser le compteur si résonance détectée
			atom.FreezeIterationCount = 0
		}
	}
}

// IterateNetwork performs one asynchronous iteration of the entire network
// Each atom operates independently, updating its state based on neighbor states
func (network *AtomicNetwork) IterateNetwork() {
	network.mutex.Lock()
	network.GlobalIteration++
	network.mutex.Unlock()

	// Collect neighbor states for all atoms (read phase)
	neighborStates := make([]map[int]float64, len(network.Atoms))
	receivedResonances := make([]float64, len(network.Atoms)) // NEW: Track resonance received

	for i := range network.Atoms {
		neighborStates[i] = make(map[int]float64)
		for _, neighborID := range network.Atoms[i].Neighbors {
			if neighborID >= 0 && neighborID < len(network.Atoms) {
				neighborStates[i][neighborID] = network.Atoms[neighborID].GetState()
			}
		}
		receivedResonances[i] = 0.0 // Initialize
	}

	// Calculate resonances that each atom perceives (freeze detection)
	for i := range network.Atoms {
		totalResonance := 0.0
		for _, neighborState := range neighborStates[i] {
			resonance := network.Atoms[i].ComputeResonance(neighborState, network.ResonanceSensitivity)
			totalResonance += resonance
		}
		// Use average resonance instead of max
		if len(neighborStates[i]) > 0 {
			receivedResonances[i] = totalResonance / float64(len(neighborStates[i]))
		} else {
			receivedResonances[i] = 0.0
		}
	}

	// Update all atoms asynchronously (write phase)
	// Each atom acts independently without waiting for others
	frozenCount := 0 // Track frozen atoms for statistics
	for i := range network.Atoms {
		// Calculate local activity for freeze system
		localActivity := network.Atoms[i].ComputeLocalActivity(
			neighborStates[i],
			receivedResonances[i],
		)

		// Update freeze state
		network.Atoms[i].UpdateFreezeState(
			localActivity,
			network.FreezeThreshold,
			network.FreezeIterations,
			network.WakeThreshold,
			receivedResonances[i],
		)

		// Count frozen atoms
		if network.Atoms[i].IsFrozen {
			frozenCount++
		}

		// Atom i updates based on its current perception of neighbors
		network.Atoms[i].UpdateState(
			neighborStates[i],
			network.CouplingCoefficient,
			network.LocalRulesCoefficient,
			network.ResonanceSensitivity,
		)

		// Update connection weights based on coherence
		network.Atoms[i].UpdateConnections(
			neighborStates[i],
			network.ReinforcementFactor,
			network.DecayFactor,
		)
	}

	// Update frozen atom count
	network.FrozenAtomsCount = frozenCount

	// Update global network metrics
	network.UpdateNetworkMetrics()
}

// UpdateNetworkMetrics calculates global network statistics
func (network *AtomicNetwork) UpdateNetworkMetrics() {
	network.mutex.Lock()
	defer network.mutex.Unlock()

	totalEnergy := 0.0
	for _, atom := range network.Atoms {
		totalEnergy += atom.EnergyConsumption
	}
	network.TotalEnergy = totalEnergy
}

// GetNetworkCoherence measures global coherence of the atomic network
// Returns value between 0 (completely incoherent) and 1 (perfectly coherent)
func (network *AtomicNetwork) GetNetworkCoherence() float64 {
	network.mutex.RLock()
	defer network.mutex.RUnlock()

	if len(network.Atoms) == 0 {
		return 0.0
	}

	totalDistance := 0.0
	maxDistance := 1.0 // Max distance between states [0, 1]

	// Calculate average distance between all atom states
	for i := 0; i < len(network.Atoms); i++ {
		for j := i + 1; j < len(network.Atoms); j++ {
			si := network.Atoms[i].GetState()
			sj := network.Atoms[j].GetState()
			totalDistance += math.Abs(si - sj)
		}
	}

	// Calculate average normalized distance
	numPairs := len(network.Atoms) * (len(network.Atoms) - 1) / 2
	if numPairs == 0 {
		return 1.0
	}

	avgDistance := totalDistance / float64(numPairs)
	coherence := 1.0 - (avgDistance / maxDistance)

	if coherence < 0 {
		coherence = 0
	}
	return coherence
}

// GetAverageActivation returns the mean activation level across all atoms
func (network *AtomicNetwork) GetAverageActivation() float64 {
	network.mutex.RLock()
	defer network.mutex.RUnlock()

	if len(network.Atoms) == 0 {
		return 0.0
	}

	totalActivation := 0.0
	for _, atom := range network.Atoms {
		totalActivation += atom.GetState()
	}

	return totalActivation / float64(len(network.Atoms))
}

// ExtractEmergentBehavior identifies stable global patterns from local interactions
// Returns a map of detected behavioral clusters
func (network *AtomicNetwork) ExtractEmergentBehavior() map[string]interface{} {
	network.mutex.RLock()
	defer network.mutex.RUnlock()

	behavior := make(map[string]interface{})

	// Identify high-activation clusters

	behavior["coherence"] = network.GetNetworkCoherence()
	behavior["average_activation"] = network.GetAverageActivation()
	behavior["iteration"] = network.GlobalIteration

	return behavior
}
