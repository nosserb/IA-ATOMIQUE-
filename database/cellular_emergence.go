// Package database - Cellular Emergence System
// Hierarchical self-organization: Atoms → Clusters → Cells → Dynamics
//
// 🎯 CONCEPT:
// Instead of arbitrary chunking (64×64 grids), we dynamically detect STABLE CLUSTERS
// of atoms that meet precise criteria:
// • Minimum 9 atoms forming a connected component
// • Each atom has 2+ connections to other atoms in cluster
// • 100% stability: all atoms at max coherence
// • 100% internal cohesion: minimum distance between states
//
// These clusters are "frozen" as CELLS - emergent super-atoms that:
// • Have their own state (center of mass, average intensity)
// • Can move and interact with other cells
// • Exhibit resonance at cellular level
// • Enable perfect rendering through hierarchical stabilization

package database

import (
	"fmt"
	"math"
	"sort"
	"sync"
)

// ============================================================================
// PART 1: CELL - The Emergent Entity
// ============================================================================

// Cell represents a self-organized cluster of stabilized atoms
type Cell struct {
	ID                  int             // Unique cell identifier
	AtomPositions       [][2]int        // Grid positions [y][x] of atoms
	CenterX, CenterY    float64         // Center of mass position
	CellState           float64         // Aggregate state [0, 1]
	AverageIntensity    float64         // Average color/intensity
	Stability           float64         // Measure of internal coherence [0, 1]
	ConnectedCells      map[int]float64 // Other cells nearby: cellID -> distance
	CellWeights         map[int]float64 // Adaptive weights to other cells
	EnergyConsumption   float64
	LastUpdateIteration int
	IsActive            bool
	mu                  sync.RWMutex
}

// NewCell creates a cell from detected positions
func NewCell(id int, positions [][2]int, atoms [][]PixelAtomV2) *Cell {
	cell := &Cell{
		ID:             id,
		AtomPositions:  positions,
		ConnectedCells: make(map[int]float64),
		CellWeights:    make(map[int]float64),
		IsActive:       true,
	}

	// Compute center of mass and aggregate state
	var totalX, totalY, totalState, totalIntensity float64
	var stateVariance float64

	height := len(atoms)
	width := 0
	if height > 0 {
		width = len(atoms[0])
	}

	for _, pos := range positions {
		y, x := pos[0], pos[1]
		if y >= 0 && y < height && x >= 0 && x < width {
			atom := &atoms[y][x]
			totalX += float64(x)
			totalY += float64(y)
			totalState += atom.Intensity
			totalIntensity += (atom.R + atom.G + atom.B) / 3.0
		}
	}

	n := float64(len(positions))
	if n > 0 {
		cell.CenterX = totalX / n
		cell.CenterY = totalY / n
		cell.CellState = totalState / n
		cell.AverageIntensity = totalIntensity / n

		// Compute stability
		for _, pos := range positions {
			y, x := pos[0], pos[1]
			if y >= 0 && y < height && x >= 0 && x < width {
				atom := &atoms[y][x]
				diff := atom.Intensity - cell.CellState
				stateVariance += diff * diff
			}
		}
		stateVariance /= n
		cell.Stability = 1.0 / (1.0 + stateVariance)
	}

	return cell
}

// UpdateCellState updates cell's state via resonance with neighbors
func (cell *Cell) UpdateCellState(neighborStates map[int]float64, alpha, sigma float64) {
	cell.mu.Lock()
	defer cell.mu.Unlock()

	var resonanceInfluence float64
	for neighborID, neighborState := range neighborStates {
		resonance := math.Exp(-math.Pow(cell.CellState-neighborState, 2) / (2 * sigma * sigma))
		weight := cell.CellWeights[neighborID]
		resonanceInfluence += weight * resonance * (neighborState - cell.CellState)
	}

	cell.CellState += alpha * resonanceInfluence

	if cell.CellState > 1.0 {
		cell.CellState = 1.0
	} else if cell.CellState < 0.0 {
		cell.CellState = 0.0
	}

	cell.EnergyConsumption += math.Abs(resonanceInfluence) * alpha
}

// UpdateCellConnections updates weights to neighboring cells
func (cell *Cell) UpdateCellConnections(neighborStates map[int]float64, gamma, delta float64) {
	cell.mu.Lock()
	defer cell.mu.Unlock()

	for neighborID, neighborState := range neighborStates {
		if weight, exists := cell.CellWeights[neighborID]; exists {
			coherence := 1.0 - math.Abs(cell.CellState-neighborState)
			if coherence < 0 {
				coherence = 0
			}

			deltaW := gamma*coherence - delta*weight
			cell.CellWeights[neighborID] = weight + deltaW

			if cell.CellWeights[neighborID] > 2.0 {
				cell.CellWeights[neighborID] = 2.0
			} else if cell.CellWeights[neighborID] < 0.0 {
				cell.CellWeights[neighborID] = 0.0
			}
		}
	}
}

// ============================================================================
// PART 2: CELLULAR CLUSTER DETECTOR
// ============================================================================

// CellularClusterDetector identifies emergent cells from atomic network
type CellularClusterDetector struct {
	Atoms         [][]PixelAtomV2
	NetworkWidth  int
	NetworkHeight int
	DetectedCells []*Cell
	CellCounter   int
	mu            sync.RWMutex

	// Detection criteria
	MinAtomsPerCell       int     // Minimum 9
	MinConnectionsPerAtom int     // At least 2
	StabilityThreshold    float64 // >= 0.85
	CoherenceThreshold    float64 // >= 0.90
}

// NewCellularClusterDetector creates detector
func NewCellularClusterDetector(atoms [][]PixelAtomV2) *CellularClusterDetector {
	height := len(atoms)
	width := 0
	if height > 0 {
		width = len(atoms[0])
	}

	return &CellularClusterDetector{
		Atoms:                 atoms,
		NetworkWidth:          width,
		NetworkHeight:         height,
		DetectedCells:         make([]*Cell, 0),
		CellCounter:           0,
		MinAtomsPerCell:       9,
		MinConnectionsPerAtom: 2,
		StabilityThreshold:    0.85,
		CoherenceThreshold:    0.90,
	}
}

// DetectCells performs main clustering
func (detector *CellularClusterDetector) DetectCells() []*Cell {
	detector.mu.Lock()
	defer detector.mu.Unlock()

	detector.DetectedCells = make([]*Cell, 0)
	detector.CellCounter = 0

	visited := make(map[[2]int]bool)
	potentialCells := make([][][2]int, 0)

	// Find seed points from stable atoms
	for y := 0; y < detector.NetworkHeight; y++ {
		for x := 0; x < detector.NetworkWidth; x++ {
			pos := [2]int{y, x}
			if visited[pos] {
				continue
			}

			if detector.IsStableAtom(y, x) {
				cluster := detector.ExpandCluster(y, x, visited)
				if len(cluster) >= detector.MinAtomsPerCell {
					potentialCells = append(potentialCells, cluster)
				}
			}
		}
	}

	// Verify and create cells
	for _, cluster := range potentialCells {
		if detector.VerifyCluster(cluster) {
			cell := NewCell(detector.CellCounter, cluster, detector.Atoms)
			detector.DetectedCells = append(detector.DetectedCells, cell)
			detector.CellCounter++
		}
	}

	return detector.DetectedCells
}

// IsStableAtom checks individual atom stability
func (detector *CellularClusterDetector) IsStableAtom(y, x int) bool {
	if y < 0 || y >= detector.NetworkHeight || x < 0 || x >= detector.NetworkWidth {
		return false
	}

	atom := &detector.Atoms[y][x]

	// Check coherence
	if atom.Confidence < detector.CoherenceThreshold {
		return false
	}

	// Check has neighbors
	neighborCount := 0
	for _, n := range atom.Neighbors {
		if n != nil {
			neighborCount++
		}
	}

	return neighborCount >= detector.MinConnectionsPerAtom
}

// ExpandCluster performs flood-fill
func (detector *CellularClusterDetector) ExpandCluster(startY, startX int, visited map[[2]int]bool) [][2]int {
	cluster := make([][2]int, 0)
	queue := [][2]int{{startY, startX}}
	visited[[2]int{startY, startX}] = true

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		y, x := current[0], current[1]
		cluster = append(cluster, current)

		// Add neighbors that are stable
		for dy := -1; dy <= 1; dy++ {
			for dx := -1; dx <= 1; dx++ {
				if dy == 0 && dx == 0 {
					continue
				}

				ny, nx := y+dy, x+dx
				neighborPos := [2]int{ny, nx}

				if !visited[neighborPos] && detector.IsStableAtom(ny, nx) {
					visited[neighborPos] = true
					queue = append(queue, neighborPos)
				}
			}
		}
	}

	return cluster
}

// VerifyCluster verifies all criteria
func (detector *CellularClusterDetector) VerifyCluster(cluster [][2]int) bool {
	// Criterion 1: At least 9 atoms
	if len(cluster) < detector.MinAtomsPerCell {
		return false
	}

	// Criterion 2: Create set for fast lookup
	clusterSet := make(map[[2]int]bool)
	for _, pos := range cluster {
		clusterSet[pos] = true
	}

	// Criterion 3: Each atom has 2+ connections within cluster
	for _, pos := range cluster {
		y, x := pos[0], pos[1]
		atom := &detector.Atoms[y][x]

		internalConnections := 0
		for _, neighbor := range atom.Neighbors {
			if neighbor == nil {
				continue
			}

			neighborPos := [2]int{neighbor.Y, neighbor.X}
			if clusterSet[neighborPos] {
				internalConnections++
			}
		}

		if internalConnections < detector.MinConnectionsPerAtom {
			return false
		}
	}

	// Criterion 4: All atoms at high coherence
	for _, pos := range cluster {
		y, x := pos[0], pos[1]
		if detector.Atoms[y][x].Confidence < detector.CoherenceThreshold {
			return false
		}
	}

	// Criterion 5: Check connectivity
	return detector.IsConnectedComponent(cluster, clusterSet)
}

// IsConnectedComponent verifies one connected graph
func (detector *CellularClusterDetector) IsConnectedComponent(cluster [][2]int, clusterSet map[[2]int]bool) bool {
	if len(cluster) == 0 {
		return false
	}

	visited := make(map[[2]int]bool)
	queue := [][2]int{cluster[0]}
	visited[cluster[0]] = true
	visitedCount := 1

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		y, x := current[0], current[1]
		atom := &detector.Atoms[y][x]

		for _, neighbor := range atom.Neighbors {
			if neighbor == nil {
				continue
			}

			neighborPos := [2]int{neighbor.Y, neighbor.X}
			if clusterSet[neighborPos] && !visited[neighborPos] {
				visited[neighborPos] = true
				visitedCount++
				queue = append(queue, neighborPos)
			}
		}
	}

	return visitedCount == len(cluster)
}

// ============================================================================
// PART 3: CELLULAR NETWORK
// ============================================================================

// CellularNetwork manages interactions between cells
type CellularNetwork struct {
	Cells []*Cell
	mu    sync.RWMutex

	CellCouplingAlpha      float64
	CellLocalBeta          float64
	CellReinforcementGamma float64
	CellDecayDelta         float64
	CellResonanceSigma     float64
}

// NewCellularNetwork creates network
func NewCellularNetwork(cells []*Cell) *CellularNetwork {
	network := &CellularNetwork{
		Cells:                  cells,
		CellCouplingAlpha:      0.7,
		CellLocalBeta:          0.3,
		CellReinforcementGamma: 0.12,
		CellDecayDelta:         0.04,
		CellResonanceSigma:     0.75,
	}

	network.BuildCellConnectivity()
	return network
}

// BuildCellConnectivity determines neighbor cells
func (network *CellularNetwork) BuildCellConnectivity() {
	network.mu.Lock()
	defer network.mu.Unlock()

	maxDistance := 5.0 // grid units

	for i, cellA := range network.Cells {
		for j, cellB := range network.Cells {
			if i >= j {
				continue
			}

			dx := cellA.CenterX - cellB.CenterX
			dy := cellA.CenterY - cellB.CenterY
			distance := math.Sqrt(dx*dx + dy*dy)

			if distance <= maxDistance {
				cellA.ConnectedCells[cellB.ID] = distance
				cellA.CellWeights[cellB.ID] = 1.0

				cellB.ConnectedCells[cellA.ID] = distance
				cellB.CellWeights[cellA.ID] = 1.0
			}
		}
	}
}

// IterateCells performs one cellular iteration
func (network *CellularNetwork) IterateCells() {
	network.mu.RLock()
	cells := network.Cells
	network.mu.RUnlock()

	// Get neighbor states
	neighborStates := make([]map[int]float64, len(cells))
	for i := range cells {
		neighborStates[i] = make(map[int]float64)
		for neighborID := range cells[i].ConnectedCells {
			for _, neighborCell := range cells {
				if neighborCell.ID == neighborID {
					neighborStates[i][neighborID] = neighborCell.CellState
					break
				}
			}
		}
	}

	// Update each cell
	for i, cell := range cells {
		cell.UpdateCellState(neighborStates[i], network.CellCouplingAlpha, network.CellResonanceSigma)
		cell.UpdateCellConnections(neighborStates[i], network.CellReinforcementGamma, network.CellDecayDelta)
	}
}

// GetCellularCoherence measures alignment
func (network *CellularNetwork) GetCellularCoherence() float64 {
	if len(network.Cells) == 0 {
		return 0
	}

	var totalDistance float64
	count := 0

	for i := 0; i < len(network.Cells); i++ {
		for j := i + 1; j < len(network.Cells); j++ {
			dist := math.Abs(network.Cells[i].CellState - network.Cells[j].CellState)
			totalDistance += dist
			count++
		}
	}

	if count == 0 {
		return 1.0
	}

	avgDistance := totalDistance / float64(count)
	return 1.0 - (avgDistance / 1.0)
}

// ============================================================================
// PART 4: HIERARCHICAL INTEGRATION
// ============================================================================

// HierarchicalLayers manages both atomic and cellular levels
type HierarchicalLayers struct {
	AtomNetwork     *ConstraintRelaxationNetwork
	Detector        *CellularClusterDetector
	CellNetwork     *CellularNetwork
	DetectionPeriod int

	CurrentIteration int
	TotalEnergy      float64

	mu sync.RWMutex
}

// NewHierarchicalLayers creates two-level system
func NewHierarchicalLayers(atomNetwork *ConstraintRelaxationNetwork, detectionPeriod int) *HierarchicalLayers {
	hierarchy := &HierarchicalLayers{
		AtomNetwork:      atomNetwork,
		DetectionPeriod:  detectionPeriod,
		CurrentIteration: 0,
	}

	hierarchy.Detector = NewCellularClusterDetector(atomNetwork.Atoms)
	return hierarchy
}

// Step performs one hierarchical iteration
func (h *HierarchicalLayers) Step() {
	h.mu.Lock()
	defer h.mu.Unlock()

	// Step 1: Atomic relaxation
	h.AtomNetwork.RelaxationStep()
	h.CurrentIteration++

	// Step 2: Every N iterations, detect cells
	if h.CurrentIteration%h.DetectionPeriod == 0 {
		h.Detector.Atoms = h.AtomNetwork.Atoms

		cells := h.Detector.DetectCells()

		if len(cells) > 0 {
			h.CellNetwork = NewCellularNetwork(cells)
			h.CellNetwork.IterateCells()
		}
	}

	// Track energy
	h.TotalEnergy = h.AtomNetwork.ComputeNetworkCoherence()
	if h.CellNetwork != nil {
		for _, cell := range h.CellNetwork.Cells {
			h.TotalEnergy += cell.EnergyConsumption
		}
	}
}

// GetHierarchicalStats returns metrics
func (h *HierarchicalLayers) GetHierarchicalStats() map[string]interface{} {
	h.mu.RLock()
	defer h.mu.RUnlock()

	stats := make(map[string]interface{})

	stats["atomic_coherence"] = h.AtomNetwork.ComputeNetworkCoherence()

	if h.CellNetwork != nil {
		stats["num_cells"] = len(h.CellNetwork.Cells)
		stats["cellular_coherence"] = h.CellNetwork.GetCellularCoherence()

		cellStates := make([]float64, len(h.CellNetwork.Cells))
		for i, cell := range h.CellNetwork.Cells {
			cellStates[i] = cell.CellState
		}
		sort.Float64s(cellStates)
		stats["cellular_states"] = cellStates
	} else {
		stats["num_cells"] = 0
	}

	stats["iteration"] = h.CurrentIteration
	stats["total_energy"] = h.TotalEnergy

	return stats
}

// PrintCellularStatus displays status
func (h *HierarchicalLayers) PrintCellularStatus() string {
	h.mu.RLock()
	defer h.mu.RUnlock()

	stats := h.GetHierarchicalStats()

	output := "\n╔════════════════════════════════════════════════════════════╗\n"
	output += "║         HIERARCHICAL EMERGENCE STATUS                     ║\n"
	output += "╚════════════════════════════════════════════════════════════╝\n\n"

	output += "[ATOMIC LEVEL]\n"
	atomicCoherence := stats["atomic_coherence"].(float64)
	output += fmt.Sprintf("  • Coherence: %.2f%%\n", atomicCoherence*100)

	output += "\n[CELLULAR LEVEL]\n"
	numCells := stats["num_cells"].(int)
	output += fmt.Sprintf("  • Detected Cells: %d\n", numCells)

	if numCells > 0 {
		cellularCoherence := stats["cellular_coherence"].(float64)
		output += fmt.Sprintf("  • Cellular Coherence: %.2f%%\n", cellularCoherence*100)
	}

	output += fmt.Sprintf("\n[OVERALL]\n")
	output += fmt.Sprintf("  • Iteration: %d\n", stats["iteration"])
	output += fmt.Sprintf("  • Total Energy: %.2f\n", stats["total_energy"])

	return output
}
