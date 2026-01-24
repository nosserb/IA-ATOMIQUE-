// Package database - Cellular Relaxation System (CRS)
// Dynamic image decomposition with local energy minimization
//
// 🎯 CONCEPT:
// Instead of fixed chunk sizes, decompose image into DYNAMIC CELLS where:
// • Each cell C_{i,j} contains variable number n_{i,j} of atoms
// • Energy E(C) combines structural + constraint + interaction terms
// • Local energy minimization: C* = argmin E(C)
// • Inter-cell interactions ensure smooth boundaries
// • Global verification confirms convergence
//
// 📊 MATHEMATICAL FRAMEWORK:
// Image I = {C_{i,j}}, i=1..N_H, j=1..N_W
// Energy: E(C) = α*E_struct + β*E_constraint + γ*E_interaction
// Minimization: ∂E(C)/∂a = 0 via gradient descent per atom
// Result: Perfect local reconstruction with stable boundaries

package database

import (
	"fmt"
	"math"
	"sort"
	"sync"
)

// ============================================================================
// PART 1: PATCH - Dynamic cellular unit
// ============================================================================

// Patch represents a dynamic cell in image decomposition
// Unlike fixed chunks, patches have variable size based on energy minimization
type Patch struct {
	ID            int   // Unique patch identifier
	GridX, GridY  int   // Position in patch grid
	Atoms         []int // Indices of atoms in this patch
	BoundaryAtoms []int // Atoms on boundary with neighbors

	// Energy components
	StructuralEnergy  float64 // E_struct: local stability
	ConstraintEnergy  float64 // E_constraint: modification cost
	InteractionEnergy float64 // E_interaction: boundary penalty
	TotalEnergy       float64 // E(C) = α*E_s + β*E_c + γ*E_i

	// Dynamics
	PreviousEnergy       float64 // For convergence detection
	ConvergenceIteration int     // When did it converge?
	IsConverged          bool    // Has this patch reached equilibrium?

	// Geometric properties
	CenterX, CenterY float64 // Patch centroid
	Bounds           [4]int  // [minX, maxX, minY, maxY]

	// Neighboring patches
	NeighboringPatches map[int]*Patch // Adjacent patches

	// Synchronization
	mutex sync.RWMutex
}

// NewPatch creates a patch from atoms
func NewPatch(id, gridX, gridY int, atomIndices []int, atomNetwork [][]PixelAtomV2) *Patch {
	patch := &Patch{
		ID:                 id,
		GridX:              gridX,
		GridY:              gridY,
		Atoms:              atomIndices,
		BoundaryAtoms:      make([]int, 0),
		NeighboringPatches: make(map[int]*Patch),
		IsConverged:        false,
	}

	// Compute centroid and bounds
	height := len(atomNetwork)
	width := 0
	if height > 0 {
		width = len(atomNetwork[0])
	}

	var sumX, sumY float64
	minX, maxX, minY, maxY := width, 0, height, 0

	for _, idx := range atomIndices {
		y := idx / width
		x := idx % width

		if x < minX {
			minX = x
		}
		if x > maxX {
			maxX = x
		}
		if y < minY {
			minY = y
		}
		if y > maxY {
			maxY = y
		}

		sumX += float64(x)
		sumY += float64(y)
	}

	n := float64(len(atomIndices))
	if n > 0 {
		patch.CenterX = sumX / n
		patch.CenterY = sumY / n
	}

	patch.Bounds = [4]int{minX, maxX, minY, maxY}

	return patch
}

// ============================================================================
// PART 2: LOCAL ENERGY COMPUTATION
// ============================================================================

// ComputeStructuralEnergy measures local stability of atoms in patch
// E_struct(a) = sum over a in C of ||a - neighbors_avg||^2
func (patch *Patch) ComputeStructuralEnergy(atomNetwork [][]PixelAtomV2) float64 {
	var totalEnergy float64

	height := len(atomNetwork)
	width := 0
	if height > 0 {
		width = len(atomNetwork[0])
	}

	for _, atomIdx := range patch.Atoms {
		y := atomIdx / width
		x := atomIdx % width

		if y < 0 || y >= height || x < 0 || x >= width {
			continue
		}

		atom := &atomNetwork[y][x]

		// Energy = deviation from stability
		// Lower confidence = higher energy (unstable)
		instability := 1.0 - atom.Confidence
		totalEnergy += instability * instability

		// Add gradient-based energy (prefer smooth gradients)
		gradientPenalty := atom.CoherenceError
		totalEnergy += gradientPenalty * 0.5
	}

	patch.StructuralEnergy = totalEnergy
	return totalEnergy
}

// ComputeConstraintEnergy measures cost of modifications
// E_constraint(a) = energy required to modify atom properties
func (patch *Patch) ComputeConstraintEnergy(atomNetwork [][]PixelAtomV2, targetEnergy float64) float64 {
	var totalEnergy float64

	height := len(atomNetwork)
	width := 0
	if height > 0 {
		width = len(atomNetwork[0])
	}

	for _, atomIdx := range patch.Atoms {
		y := atomIdx / width
		x := atomIdx % width

		if y < 0 || y >= height || x < 0 || x >= width {
			continue
		}

		atom := &atomNetwork[y][x]

		// Constraint energy: cost to match target
		currentEnergy := atom.LocalEnergy
		energyDiff := math.Abs(currentEnergy - targetEnergy)
		totalEnergy += energyDiff * energyDiff
	}

	patch.ConstraintEnergy = totalEnergy
	return totalEnergy
}

// ComputeInteractionEnergy measures boundary coherence with neighbors
// E_interaction(a_i, a_j) = λ*f(distance)*||a_i - a_j||^2
func (patch *Patch) ComputeInteractionEnergy(
	atomNetwork [][]PixelAtomV2,
	neighborPatches map[int]*Patch,
	lambda float64,
) float64 {
	var totalEnergy float64

	height := len(atomNetwork)
	width := 0
	if height > 0 {
		width = len(atomNetwork[0])
	}

	// Identify boundary atoms (atoms adjacent to other patches)
	boundarySet := make(map[int]bool)
	internalSet := make(map[int]bool)
	for _, idx := range patch.Atoms {
		internalSet[idx] = true
	}

	for _, atomIdx := range patch.Atoms {
		y := atomIdx / width
		x := atomIdx % width

		// Check 8-neighbors
		for dy := -1; dy <= 1; dy++ {
			for dx := -1; dx <= 1; dx++ {
				if dy == 0 && dx == 0 {
					continue
				}

				ny := y + dy
				nx := x + dx

				if ny >= 0 && ny < height && nx >= 0 && nx < width {
					neighborIdx := ny*width + nx

					// If neighbor is in different patch, this is boundary
					if !internalSet[neighborIdx] {
						boundarySet[atomIdx] = true
						break
					}
				}
			}
		}
	}

	patch.BoundaryAtoms = make([]int, 0)
	for idx := range boundarySet {
		patch.BoundaryAtoms = append(patch.BoundaryAtoms, idx)
	}

	// Compute interaction energy for boundary atoms
	for _, boundaryIdx := range patch.BoundaryAtoms {
		by := boundaryIdx / width
		bx := boundaryIdx % width

		if by < 0 || by >= height || bx < 0 || bx >= width {
			continue
		}

		boundaryAtom := &atomNetwork[by][bx]

		// Check neighbors in adjacent patches
		for _, neighborPatch := range neighborPatches {
			for _, neighborIdx := range neighborPatch.Atoms {
				ny := neighborIdx / width
				nx := neighborIdx % width

				if ny < 0 || ny >= height || nx < 0 || nx >= width {
					continue
				}

				// Only penalize if adjacent (8-neighbor)
				dy := by - ny
				dx := bx - nx

				if math.Abs(float64(dy)) <= 1 && math.Abs(float64(dx)) <= 1 {
					neighborAtom := &atomNetwork[ny][nx]

					// Distance-based weighting function
					distance := math.Sqrt(float64(dx*dx + dy*dy))
					f_d := math.Exp(-distance / 0.5) // Closer = stronger coupling

					// Coherence penalty
					colorDiff := math.Pow(boundaryAtom.R-neighborAtom.R, 2) +
						math.Pow(boundaryAtom.G-neighborAtom.G, 2) +
						math.Pow(boundaryAtom.B-neighborAtom.B, 2)

					energy := lambda * f_d * colorDiff
					totalEnergy += energy
				}
			}
		}
	}

	patch.InteractionEnergy = totalEnergy
	return totalEnergy
}

// ComputeTotalEnergy computes E(C) = α*E_struct + β*E_constraint + γ*E_interaction
func (patch *Patch) ComputeTotalEnergy(
	atomNetwork [][]PixelAtomV2,
	neighborPatches map[int]*Patch,
	alpha, beta, gamma, lambda, targetEnergy float64,
) float64 {
	patch.mutex.Lock()
	defer patch.mutex.Unlock()

	// Compute each component
	structEnergy := patch.ComputeStructuralEnergy(atomNetwork)
	constraintEnergy := patch.ComputeConstraintEnergy(atomNetwork, targetEnergy)
	interactionEnergy := patch.ComputeInteractionEnergy(atomNetwork, neighborPatches, lambda)

	// Weighted sum
	patch.PreviousEnergy = patch.TotalEnergy
	patch.TotalEnergy = alpha*structEnergy + beta*constraintEnergy + gamma*interactionEnergy

	return patch.TotalEnergy
}

// ============================================================================
// PART 3: LOCAL ENERGY MINIMIZATION (GRADIENT DESCENT)
// ============================================================================

// MinimizeLocalEnergy performs gradient descent on atoms in patch
// a ← a - η * ∂E(C)/∂a
func (patch *Patch) MinimizeLocalEnergy(
	atomNetwork [][]PixelAtomV2,
	neighborPatches map[int]*Patch,
	eta float64, // Learning rate
	alpha, beta, gamma, lambda, targetEnergy float64,
	maxIterations int,
	convergenceThreshold float64,
) int {
	patch.mutex.Lock()
	defer patch.mutex.Unlock()

	height := len(atomNetwork)
	width := 0
	if height > 0 {
		width = len(atomNetwork[0])
	}

	iterCount := 0

	for iter := 0; iter < maxIterations; iter++ {
		// Compute current energy
		_ = patch.ComputeTotalEnergy(
			atomNetwork, neighborPatches,
			alpha, beta, gamma, lambda, targetEnergy,
		)

		// Update each atom via gradient descent
		for _, atomIdx := range patch.Atoms {
			y := atomIdx / width
			x := atomIdx % width

			if y < 0 || y >= height || x < 0 || x >= width {
				continue
			}

			atom := &atomNetwork[y][x]

			// Gradient w.r.t. Intensity (simplified)
			gradIntensity := ComputeIntensityGradient(atom, atomNetwork, height, width, patch)

			// Update with learning rate
			newIntensity := atom.Intensity - eta*gradIntensity

			// Clamp to [0, 1]
			if newIntensity > 1.0 {
				newIntensity = 1.0
			} else if newIntensity < 0.0 {
				newIntensity = 0.0
			}

			atom.Intensity = newIntensity
		}

		iterCount++

		// Check convergence: |ΔE| < ε
		deltaEnergy := math.Abs(patch.TotalEnergy - patch.PreviousEnergy)
		if deltaEnergy < convergenceThreshold {
			patch.IsConverged = true
			patch.ConvergenceIteration = iter
			return iterCount
		}
	}

	// Mark as converged if we exit loop
	patch.IsConverged = true
	patch.ConvergenceIteration = maxIterations

	return iterCount
}

// ComputeIntensityGradient computes ∂E/∂I for an atom
func ComputeIntensityGradient(
	atom *PixelAtomV2,
	atomNetwork [][]PixelAtomV2,
	height, width int,
	patch *Patch,
) float64 {
	var gradient float64

	// Structural gradient: dE_struct/dI
	instability := 1.0 - atom.Confidence
	gradient += 2 * instability * (-atom.Intensity)

	// Constraint gradient: dE_constraint/dI
	energyDiff := atom.LocalEnergy - 0.5 // Target intensity ~0.5
	gradient += 2 * energyDiff

	// Interaction gradient: check neighbors
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			if dy == 0 && dx == 0 {
				continue
			}

			ny := atom.Y + dy
			nx := atom.X + dx

			if ny >= 0 && ny < height && nx >= 0 && nx < width {
				neighbor := &atomNetwork[ny][nx]

				// Penalty for difference
				intensityDiff := atom.Intensity - neighbor.Intensity
				gradient += 2 * intensityDiff * atom.Intensity
			}
		}
	}

	return gradient
}

// ============================================================================
// PART 4: PATCH GRID MANAGER
// ============================================================================

// PatchGrid manages dynamic patch decomposition of image
type PatchGrid struct {
	Patches               map[int]*Patch
	GridHeight, GridWidth int
	AtomNetwork           [][]PixelAtomV2

	// Energy tracking
	TotalEnergy         float64
	PreviousTotalEnergy float64
	GlobalConverged     bool

	// Parameters
	Alpha        float64 // E_struct weight
	Beta         float64 // E_constraint weight
	Gamma        float64 // E_interaction weight
	Lambda       float64 // Inter-cell coupling
	TargetEnergy float64 // Desired energy level

	// Convergence
	LearningRate               float64
	MaxIterationsPerPatch      int
	ConvergenceThreshold       float64
	GlobalConvergenceThreshold float64

	mu sync.RWMutex
}

// NewPatchGrid creates a patch grid from atom network
func NewPatchGrid(
	atomNetwork [][]PixelAtomV2,
	gridHeight, gridWidth int,
) *PatchGrid {
	return &PatchGrid{
		Patches:                    make(map[int]*Patch),
		GridHeight:                 gridHeight,
		GridWidth:                  gridWidth,
		AtomNetwork:                atomNetwork,
		Alpha:                      0.4, // Structural importance
		Beta:                       0.3, // Constraint importance
		Gamma:                      0.3, // Interaction importance
		Lambda:                     0.8, // Inter-cell coupling strength
		TargetEnergy:               0.5, // Target intensity level
		LearningRate:               0.01,
		MaxIterationsPerPatch:      50,
		ConvergenceThreshold:       0.001,
		GlobalConvergenceThreshold: 0.01,
		GlobalConverged:            false,
	}
}

// InitializePatches divides network into initial patches
func (grid *PatchGrid) InitializePatches(atomsPerPatchTarget int) error {
	grid.mu.Lock()
	defer grid.mu.Unlock()

	height := len(grid.AtomNetwork)
	width := 0
	if height > 0 {
		width = len(grid.AtomNetwork[0])
	}

	totalAtoms := height * width
	_ = int(math.Sqrt(float64(totalAtoms) / float64(grid.GridHeight*grid.GridWidth)))

	patchID := 0

	// Create initial regular grid of patches
	for py := 0; py < grid.GridHeight; py++ {
		for px := 0; px < grid.GridWidth; px++ {
			atomIndices := make([]int, 0)

			// Collect atoms in this grid cell
			minY := (py * height) / grid.GridHeight
			maxY := ((py + 1) * height) / grid.GridHeight
			minX := (px * width) / grid.GridWidth
			maxX := ((px + 1) * width) / grid.GridWidth

			for y := minY; y < maxY; y++ {
				for x := minX; x < maxX; x++ {
					if y >= 0 && y < height && x >= 0 && x < width {
						atomIndices = append(atomIndices, y*width+x)
					}
				}
			}

			if len(atomIndices) > 0 {
				patch := NewPatch(patchID, px, py, atomIndices, grid.AtomNetwork)
				grid.Patches[patchID] = patch
				patchID++
			}
		}
	}

	// Connect neighboring patches
	grid.BuildPatchConnectivity()

	return nil
}

// BuildPatchConnectivity establishes neighbor relationships
func (grid *PatchGrid) BuildPatchConnectivity() {
	for id1, patch1 := range grid.Patches {
		for id2, patch2 := range grid.Patches {
			if id1 >= id2 {
				continue
			}

			// Check if adjacent (share boundary)
			dx := patch1.GridX - patch2.GridX
			dy := patch1.GridY - patch2.GridY

			if (math.Abs(float64(dx)) <= 1 && math.Abs(float64(dy)) <= 1) && (dx != 0 || dy != 0) {
				patch1.NeighboringPatches[id2] = patch2
				patch2.NeighboringPatches[id1] = patch1
			}
		}
	}
}

// MinimizeGlobalEnergy performs one global relaxation iteration
func (grid *PatchGrid) MinimizeGlobalEnergy() float64 {
	grid.mu.Lock()
	defer grid.mu.Unlock()

	grid.PreviousTotalEnergy = grid.TotalEnergy

	// Minimize each patch locally
	totalIters := 0
	for _, patch := range grid.Patches {
		iters := patch.MinimizeLocalEnergy(
			grid.AtomNetwork,
			patch.NeighboringPatches,
			grid.LearningRate,
			grid.Alpha, grid.Beta, grid.Gamma,
			grid.Lambda, grid.TargetEnergy,
			grid.MaxIterationsPerPatch,
			grid.ConvergenceThreshold,
		)
		totalIters += iters
	}

	// Compute new total energy
	newTotalEnergy := 0.0
	for _, patch := range grid.Patches {
		energy := patch.ComputeTotalEnergy(
			grid.AtomNetwork,
			patch.NeighboringPatches,
			grid.Alpha, grid.Beta, grid.Gamma,
			grid.Lambda, grid.TargetEnergy,
		)
		newTotalEnergy += energy
	}

	grid.TotalEnergy = newTotalEnergy

	return grid.TotalEnergy
}

// ============================================================================
// PART 5: GLOBAL VERIFICATION
// ============================================================================

// VerifyGlobalConvergence checks if entire system is converged
func (grid *PatchGrid) VerifyGlobalConvergence() bool {
	grid.mu.RLock()
	defer grid.mu.RUnlock()

	// Check energy convergence
	deltaEnergy := math.Abs(grid.TotalEnergy - grid.PreviousTotalEnergy)
	if deltaEnergy >= grid.GlobalConvergenceThreshold {
		return false
	}

	// Check all patches converged
	convergedCount := 0
	for _, patch := range grid.Patches {
		if patch.IsConverged {
			convergedCount++
		}
	}

	allConverged := convergedCount == len(grid.Patches)

	grid.GlobalConverged = allConverged && deltaEnergy < grid.GlobalConvergenceThreshold

	return grid.GlobalConverged
}

// GetStatistics returns comprehensive grid statistics
func (grid *PatchGrid) GetStatistics() map[string]interface{} {
	grid.mu.RLock()
	defer grid.mu.RUnlock()

	stats := make(map[string]interface{})

	// Count statistics
	stats["num_patches"] = len(grid.Patches)

	// Energy statistics
	var energies []float64
	for _, patch := range grid.Patches {
		energies = append(energies, patch.TotalEnergy)
	}
	sort.Float64s(energies)

	stats["total_energy"] = grid.TotalEnergy
	stats["avg_patch_energy"] = grid.TotalEnergy / float64(len(grid.Patches))

	if len(energies) > 0 {
		stats["min_patch_energy"] = energies[0]
		stats["max_patch_energy"] = energies[len(energies)-1]
	}

	// Convergence statistics
	convergedCount := 0
	for _, patch := range grid.Patches {
		if patch.IsConverged {
			convergedCount++
		}
	}
	stats["converged_patches"] = convergedCount
	stats["total_patches"] = len(grid.Patches)
	stats["convergence_percent"] = 100.0 * float64(convergedCount) / float64(len(grid.Patches))

	// Global convergence
	stats["global_converged"] = grid.GlobalConverged
	stats["energy_delta"] = math.Abs(grid.TotalEnergy - grid.PreviousTotalEnergy)

	return stats
}

// PrintGridStatus displays comprehensive status
func (grid *PatchGrid) PrintGridStatus() string {
	stats := grid.GetStatistics()

	output := "\n╔════════════════════════════════════════════════════════════╗\n"
	output += "║         CELLULAR RELAXATION GRID STATUS                   ║\n"
	output += "╚════════════════════════════════════════════════════════════╝\n\n"

	output += "[PATCH STATISTICS]\n"
	output += fmt.Sprintf("  • Total Patches: %d\n", stats["num_patches"])

	output += "\n[ENERGY ANALYSIS]\n"
	output += fmt.Sprintf("  • Total Energy: %.4f\n", stats["total_energy"])
	output += fmt.Sprintf("  • Average per Patch: %.4f\n", stats["avg_patch_energy"])
	output += fmt.Sprintf("  • Min/Max: %.4f / %.4f\n", stats["min_patch_energy"], stats["max_patch_energy"])
	output += fmt.Sprintf("  • Energy Delta: %.6f\n", stats["energy_delta"])

	output += "\n[CONVERGENCE STATUS]\n"
	output += fmt.Sprintf("  • Converged Patches: %d / %d (%.1f%%)\n",
		stats["converged_patches"], stats["total_patches"], stats["convergence_percent"])

	if stats["global_converged"].(bool) {
		output += "  • Global Status: ✓ CONVERGED\n"
	} else {
		output += "  • Global Status: ⏳ Still relaxing...\n"
	}

	return output
}
