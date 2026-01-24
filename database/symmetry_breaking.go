// Package database - Symmetry Breaking for Image Generation
// Implements the 4 critical solutions to transform "thermodynamically correct noise"
// into structured, recognizable images without violating energy physics.
//
// 🎯 THE CORE PROBLEM SOLVED:
// Same energy ≠ Same image
// Multiple micro-states can have identical energy (degeneracy)
// We need to BREAK THE SYMMETRY to guide which specific micro-state to realize
//
// 📊 PHYSICS ANALOGY:
// ✅ We copy: temperature, pressure, density
// ❌ We DON'T copy: crystal structure
// Result: Thermodynamically correct but geometrically random
//
// 🛠️ 4 SOLUTIONS IMPLEMENTED:
// 1️⃣ Phase Continuity Term - Force coherent gradient propagation
// 2️⃣ Weak Directional Field - Low-res guidance (blur source image to 16×16)
// 3️⃣ Topological Constraint - Match edge topology, not intensity
// 4️⃣ Multi-Scale Pipeline - CRUCIAL: 32×32 → 64×64 → 128×128 → 256×256

package database

import (
	"image/png"
	"math"
	"os"
)

// ============================================================================
// SOLUTION 2: WEAK DIRECTIONAL FIELD EXTRACTION
// ============================================================================

// ExtractLowResGuidanceField creates a weak orientation guide from source image
// This is NOT a mask, NOT a model - just weak structural hints
// Physical interpretation: External weak field that orients but doesn't dictate
func ExtractLowResGuidanceField(sourcePath string, gridSize int) (*GlobalCoherenceField, error) {
	file, err := os.Open(sourcePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, err := png.Decode(file)
	if err != nil {
		return nil, err
	}

	field := NewGlobalCoherenceField()
	field.CreateDirectionalFieldFromImage(img, gridSize)

	// Extract additional weak constraints
	bounds := img.Bounds()
	width := bounds.Max.X - bounds.Min.X
	height := bounds.Max.Y - bounds.Min.Y

	var avgBrightness float64
	var totalPixels int

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			brightness := float64(r+g+b) / (3 * 65535)
			avgBrightness += brightness
			totalPixels++
		}
	}

	if totalPixels > 0 {
		field.AverageBrightness = avgBrightness / float64(totalPixels)
	}

	// CRUCIAL: This field has LOW weight (5-15%)
	field.FieldStrength = 0.10 // Only 10% influence - WEAK guidance
	field.InfluenceWeight = 0.08

	return field, nil
}

// ============================================================================
// SOLUTION 3: TOPOLOGICAL CONSTRAINT (Edge Rank Matching)
// ============================================================================

// EdgeTopologyMap extracts the topology of edges, NOT their values
// Instead of: I(x) ≈ I_ref(x)
// We enforce: rank(|∇I(x)|) ≈ rank(|∇I_ref(x)|)
// Result: Edges appear in the right PLACES, with flexible VALUES
type EdgeTopologyMap struct {
	Width, Height int
	EdgeRanks     [][]float64 // Normalized rank of edge strength [0, 1]
	EdgeMask      [][]bool    // True where edges exist
	EdgeDensity   float64     // Fraction of pixels that are edges
}

// ComputeEdgeTopology analyzes source image and extracts edge structure
func ComputeEdgeTopology(sourcePath string) (*EdgeTopologyMap, error) {
	file, err := os.Open(sourcePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, err := png.Decode(file)
	if err != nil {
		return nil, err
	}

	bounds := img.Bounds()
	width := bounds.Max.X - bounds.Min.X
	height := bounds.Max.Y - bounds.Min.Y

	topo := &EdgeTopologyMap{
		Width:     width,
		Height:    height,
		EdgeRanks: make([][]float64, height),
		EdgeMask:  make([][]bool, height),
	}

	// Compute gradient magnitude at each pixel (Sobel)
	gradMagnitudes := make([][]float64, height)
	var allGradients []float64

	for y := 0; y < height; y++ {
		gradMagnitudes[y] = make([]float64, width)
		topo.EdgeRanks[y] = make([]float64, width)
		topo.EdgeMask[y] = make([]bool, width)

		for x := 0; x < width; x++ {
			// Sobel operator
			var gx, gy float64

			for dy := -1; dy <= 1; dy++ {
				for dx := -1; dx <= 1; dx++ {
					px := x + dx
					py := y + dy

					if px >= 0 && px < width && py >= 0 && py < height {
						r, g, b, _ := img.At(px, py).RGBA()
						intensity := float64(r+g+b) / (3 * 65535)

						// Sobel kernel
						if dy == -1 && dx == -1 {
							gx += -1 * intensity
							gy += -1 * intensity
						} else if dy == -1 && dx == 0 {
							gy += -2 * intensity
						} else if dy == -1 && dx == 1 {
							gx += 1 * intensity
							gy += -1 * intensity
						} else if dy == 0 && dx == -1 {
							gx += -2 * intensity
						} else if dy == 0 && dx == 1 {
							gx += 2 * intensity
						} else if dy == 1 && dx == -1 {
							gx += -1 * intensity
							gy += 1 * intensity
						} else if dy == 1 && dx == 0 {
							gy += 2 * intensity
						} else if dy == 1 && dx == 1 {
							gx += 1 * intensity
							gy += 1 * intensity
						}
					}
				}
			}

			gradMag := math.Sqrt(gx*gx + gy*gy)
			gradMagnitudes[y][x] = gradMag
			allGradients = append(allGradients, gradMag)
		}
	}

	// Compute ranks (percentile positions)
	// Sort gradients to get rank
	sortedGrads := make([]float64, len(allGradients))
	copy(sortedGrads, allGradients)
	quickSort(sortedGrads)

	// Map each gradient to its rank [0, 1]
	edgeCount := 0
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			grad := gradMagnitudes[y][x]

			// Find rank via binary search
			rank := float64(binarySearchRank(sortedGrads, grad)) / float64(len(sortedGrads))
			topo.EdgeRanks[y][x] = rank

			// Mark as edge if rank > 0.7 (top 30% of gradients)
			if rank > 0.7 {
				topo.EdgeMask[y][x] = true
				edgeCount++
			}
		}
	}

	topo.EdgeDensity = float64(edgeCount) / float64(width*height)

	return topo, nil
}

// ApplyTopologicalConstraint adds energy penalty if edges don't match topology
func (topo *EdgeTopologyMap) ApplyTopologicalConstraint(atom *PixelAtomV2, atomX, atomY, patchSize int) float64 {
	// Map atom position to topology map
	imgX := atomX * patchSize
	imgY := atomY * patchSize

	if imgX >= topo.Width || imgY >= topo.Height {
		return 0 // Out of bounds
	}

	// Check if this location SHOULD be an edge
	shouldBeEdge := topo.EdgeMask[imgY][imgX]
	expectedRank := topo.EdgeRanks[imgY][imgX]

	// Compute current edge strength of this atom
	var maxColorDiff float64
	for _, neighbor := range atom.Neighbors {
		if neighbor != nil {
			colorDiff := math.Sqrt(
				math.Pow(atom.R-neighbor.R, 2) +
					math.Pow(atom.G-neighbor.G, 2) +
					math.Pow(atom.B-neighbor.B, 2),
			)
			if colorDiff > maxColorDiff {
				maxColorDiff = colorDiff
			}
		}
	}

	currentEdgeStrength := maxColorDiff // [0, ~1.7]
	currentRank := currentEdgeStrength / 1.7

	// Penalty = difference in ranks
	// If expectedRank = 0.9 (strong edge expected) but currentRank = 0.1 (smooth) → high penalty
	rankError := math.Abs(expectedRank - currentRank)

	// Only penalize if topology says there should/shouldn't be an edge
	if shouldBeEdge && currentRank < 0.5 {
		// Expected edge but too smooth
		return rankError * 0.3
	} else if !shouldBeEdge && currentRank > 0.7 {
		// Expected smooth but too sharp
		return rankError * 0.2
	}

	return 0
}

// ============================================================================
// SOLUTION 4: MULTI-SCALE PIPELINE (CRITICAL!)
// ============================================================================

// MultiScalePipeline implements the ESSENTIAL coarse-to-fine approach
// WITHOUT this, energy distributes randomly at local scale → chaos
//
// Pipeline:
// 32×32   → relaxation (establish global structure)
// 64×64   → relaxation (refine structure)
// 128×128 → relaxation (add detail)
// 256×256 → relaxation (final detail)
//
// Each scale inherits from previous, providing CONSTRAINTS for next level
type MultiScalePipeline struct {
	Scales             []int // [32, 64, 128, 256]
	RelaxationPerScale []int // Iterations per scale [100, 200, 300, 500]
	CurrentScale       int
}

// NewMultiScalePipeline creates the pipeline
func NewMultiScalePipeline() *MultiScalePipeline {
	return &MultiScalePipeline{
		Scales:             []int{32, 64, 128, 256},
		RelaxationPerScale: []int{150, 200, 250, 400},
		CurrentScale:       0,
	}
}

// NewAdaptiveMultiScalePipeline creates a pipeline adapted to target size
func NewAdaptiveMultiScalePipeline(targetWidth, targetHeight int) *MultiScalePipeline {
	// Determine scales based on target size
	maxDim := targetWidth
	if targetHeight > maxDim {
		maxDim = targetHeight
	}

	scales := []int{}
	iterations := []int{}

	// Start with coarse scale (32×32 grid minimum)
	currentScale := 32
	for currentScale <= maxDim {
		scales = append(scales, currentScale)
		// More iterations for finer scales
		iterCount := 100 + (len(scales) * 50)
		iterations = append(iterations, iterCount)
		currentScale *= 2
	}

	// Ensure we have at least 2 scales
	if len(scales) < 2 {
		scales = []int{32, maxDim}
		iterations = []int{150, 300}
	}

	return &MultiScalePipeline{
		Scales:             scales,
		RelaxationPerScale: iterations,
		CurrentScale:       0,
	}
}

// RunMultiScalePipeline executes the full coarse-to-fine generation
func (pipeline *MultiScalePipeline) RunMultiScalePipeline(
	targetWidth, targetHeight int,
	energyProfile *ImageEnergyProfile,
	guidanceField *GlobalCoherenceField,
	edgeTopology *EdgeTopologyMap,
) *ConstraintRelaxationNetwork {

	var network *ConstraintRelaxationNetwork

	for i, scale := range pipeline.Scales {
		patchSize := targetWidth / scale // Scale determines patch size (adapt to target size)
		if patchSize < 1 {
			patchSize = 1
		}

		// Create or resize network
		if i == 0 {
			// Start at coarsest scale
			network = NewConstraintRelaxationNetwork(targetWidth, targetHeight, patchSize)
		} else {
			// Upscale from previous network
			network = UpscaleNetwork(network, patchSize)
		}

		// Apply energy profile and constraints
		network.EnergyProfile = energyProfile
		network.GlobalField = guidanceField

		// Relax at this scale
		iterations := pipeline.RelaxationPerScale[i]
		for iter := 0; iter < iterations; iter++ {
			network.RelaxationStep()

			// Apply topological constraint every 10 iterations
			if edgeTopology != nil && iter%10 == 0 {
				network.ApplyTopologyConstraint(edgeTopology)
			}

			// Progress feedback
			if iter%50 == 0 && iter > 0 {
				coherence := network.ComputeNetworkCoherence()
				_ = coherence // Use or print as needed
			}
		}

		// Capture state before moving to next scale
		network.CapturePhaseState()
		pipeline.CurrentScale++
	}

	return network
}

// UpscaleNetwork doubles the resolution while preserving structure
func UpscaleNetwork(source *ConstraintRelaxationNetwork, newPatchSize int) *ConstraintRelaxationNetwork {
	// Create new network at double resolution
	newNet := NewConstraintRelaxationNetwork(source.Width, source.Height, newPatchSize)

	// Copy global field and energy profile
	newNet.GlobalField = source.GlobalField
	newNet.EnergyProfile = source.EnergyProfile

	// Inherit state from coarser network (interpolation)
	srcHeight := len(source.Atoms)
	srcWidth := 0
	if srcHeight > 0 {
		srcWidth = len(source.Atoms[0])
	}

	newHeight := len(newNet.Atoms)
	newWidth := 0
	if newHeight > 0 {
		newWidth = len(newNet.Atoms[0])
	}

	// Bilinear interpolation from coarse to fine
	for y := 0; y < newHeight; y++ {
		for x := 0; x < newWidth; x++ {
			// Map to source coordinates
			srcX := float64(x) * float64(srcWidth) / float64(newWidth)
			srcY := float64(y) * float64(srcHeight) / float64(newHeight)

			// Get integer part
			x0 := int(srcX)
			y0 := int(srcY)
			x1 := x0 + 1
			y1 := y0 + 1

			// Clamp
			if x1 >= srcWidth {
				x1 = srcWidth - 1
			}
			if y1 >= srcHeight {
				y1 = srcHeight - 1
			}

			// Fractional part
			fx := srcX - float64(x0)
			fy := srcY - float64(y0)

			// Bilinear interpolation (use pointers to avoid mutex copy)
			c00 := &source.Atoms[y0][x0]
			c10 := &source.Atoms[y0][x1]
			c01 := &source.Atoms[y1][x0]
			c11 := &source.Atoms[y1][x1]

			// Interpolate color
			newNet.Atoms[y][x].R = (1-fx)*(1-fy)*c00.R + fx*(1-fy)*c10.R +
				(1-fx)*fy*c01.R + fx*fy*c11.R
			newNet.Atoms[y][x].G = (1-fx)*(1-fy)*c00.G + fx*(1-fy)*c10.G +
				(1-fx)*fy*c01.G + fx*fy*c11.G
			newNet.Atoms[y][x].B = (1-fx)*(1-fy)*c00.B + fx*(1-fy)*c10.B +
				(1-fx)*fy*c01.B + fx*fy*c11.B

			// Interpolate other properties
			newNet.Atoms[y][x].Intensity = (1-fx)*(1-fy)*c00.Intensity +
				fx*(1-fy)*c10.Intensity + (1-fx)*fy*c01.Intensity + fx*fy*c11.Intensity
			newNet.Atoms[y][x].Orientation = c00.Orientation // Inherit orientation
			newNet.Atoms[y][x].Confidence = 0.5              // Reset confidence for refinement
		}
	}

	return newNet
}

// ApplyTopologyConstraint applies topological constraint to all atoms
func (net *ConstraintRelaxationNetwork) ApplyTopologyConstraint(topo *EdgeTopologyMap) {
	gridH := len(net.Atoms)
	gridW := 0
	if gridH > 0 {
		gridW = len(net.Atoms[0])
	}

	for y := 0; y < gridH; y++ {
		for x := 0; x < gridW; x++ {
			atom := &net.Atoms[y][x]
			penalty := topo.ApplyTopologicalConstraint(atom, x, y, net.PatchSize)

			// Add penalty to atom's energy
			atom.LocalEnergy += penalty
		}
	}
}

// ComputeNetworkCoherence returns overall network coherence [0, 1]
func (net *ConstraintRelaxationNetwork) ComputeNetworkCoherence() float64 {
	gridH := len(net.Atoms)
	gridW := 0
	if gridH > 0 {
		gridW = len(net.Atoms[0])
	}

	var totalCoherence float64
	count := 0

	for y := 0; y < gridH; y++ {
		for x := 0; x < gridW; x++ {
			totalCoherence += net.Atoms[y][x].Confidence
			count++
		}
	}

	if count > 0 {
		return totalCoherence / float64(count)
	}
	return 0
}

// Helper functions
func quickSort(arr []float64) {
	if len(arr) < 2 {
		return
	}
	quickSortHelper(arr, 0, len(arr)-1)
}

func quickSortHelper(arr []float64, low, high int) {
	if low < high {
		pi := partition(arr, low, high)
		quickSortHelper(arr, low, pi-1)
		quickSortHelper(arr, pi+1, high)
	}
}

func partition(arr []float64, low, high int) int {
	pivot := arr[high]
	i := low - 1
	for j := low; j < high; j++ {
		if arr[j] < pivot {
			i++
			arr[i], arr[j] = arr[j], arr[i]
		}
	}
	arr[i+1], arr[high] = arr[high], arr[i+1]
	return i + 1
}

func binarySearchRank(sortedArr []float64, value float64) int {
	left, right := 0, len(sortedArr)-1
	for left <= right {
		mid := (left + right) / 2
		if sortedArr[mid] < value {
			left = mid + 1
		} else if sortedArr[mid] > value {
			right = mid - 1
		} else {
			return mid
		}
	}
	return left
}
