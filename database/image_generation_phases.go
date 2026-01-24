// Package database - Multi-Scale Image Generation with Atomic Resonance
// Implements 5 phases of image generation based on local atomic interactions
// Phase 1: Multi-scale structuration (pixels → coherent patterns)
// Phase 2: Shape emergence (capsules of resonance)
// Phase 3: Prompt conditioning (user guidance)
// Phase 4: Iterative refinement (detail and texture)
// Phase 5: Atomic coherence verification (quality control)

package database

import (
	"image"
	"image/color"
	"math"
	"sync"
)

// ========== PHASE 1: MULTI-SCALE STRUCTURATION ==========

// MultiScaleResonanceLayer represents one scale level (micro, meso, macro)
type MultiScaleResonanceLayer struct {
	Scale          int             // Patch size for this layer (1, 4, 8, 16...)
	Atoms          [][]PixelAtom   // Atoms at this scale
	Width, Height  int             // Dimensions of grid
	Alpha          float64         // Coupling coefficient α
	Iteration      int             // Current iteration at this scale
	ResonanceCache map[int]float64 // Cached resonance values for efficiency
	mutex          sync.RWMutex
}

// NewMultiScaleResonanceLayer creates a layer at a specific scale
func NewMultiScaleResonanceLayer(width, height, scale int) *MultiScaleResonanceLayer {
	layer := &MultiScaleResonanceLayer{
		Scale:          scale,
		Atoms:          make([][]PixelAtom, height/scale),
		Width:          width / scale,
		Height:         height / scale,
		Alpha:          0.7, // Tunable coupling coefficient
		Iteration:      0,
		ResonanceCache: make(map[int]float64),
	}

	// Initialize atoms at this scale
	for y := 0; y < height/scale; y++ {
		layer.Atoms[y] = make([]PixelAtom, width/scale)
		for x := 0; x < width/scale; x++ {
			id := y*(width/scale) + x
			atom := NewPixelAtom(id, x, y, scale)
			// Add some initial coherence
			atom.Features["coherence"] = 0.3
			layer.Atoms[y][x] = *atom
		}
	}

	return layer
}

// PhaseOne_StructurationMultiEchelle performs multi-scale pixel alignment
// Goal: Transform isolated pixels into coherent patterns through local resonance
func (net *AtomicImageNetwork) PhaseOne_StructurationMultiEchelle(iterations int) {
	// Create multiple scale layers
	layers := make([]*MultiScaleResonanceLayer, 0)
	scales := []int{1, 4, 8} // Micro, meso, macro

	for _, scale := range scales {
		layer := NewMultiScaleResonanceLayer(net.Width, net.Height, scale)
		layers = append(layers, layer)
	}

	// Iterate with local resonance at each scale
	for iter := 0; iter < iterations; iter++ {
		// Update from micro → macro (bottom-up emergence)
		for layerIdx, layer := range layers {
			net.IterateScale(layer)

			// Cross-scale influence: propagate patterns from smaller to larger scales
			if layerIdx > 0 {
				net.PropagatePatterns(layers[layerIdx-1], layer)
			}
		}
	}

	// Copy converged patterns back to main network
	layer := layers[len(layers)-1] // Use macro scale
	for y := 0; y < net.Height/net.PatchSize; y++ {
		for x := 0; x < net.Width/net.PatchSize; x++ {
			if y < len(layer.Atoms) && x < len(layer.Atoms[y]) {
				net.Atoms[y][x].Color = layer.Atoms[y][x].Color
				net.Atoms[y][x].Intensity = layer.Atoms[y][x].Intensity
				net.Atoms[y][x].State = layer.Atoms[y][x].State
			}
		}
	}
}

// IterateScale performs one iteration at a given scale
// Implements: c_ij(t+1) = c_ij(t) + α * Σ(w_ijk * (c_k(t) - c_ij(t)))
func (net *AtomicImageNetwork) IterateScale(layer *MultiScaleResonanceLayer) {
	var wg sync.WaitGroup
	height := len(layer.Atoms)
	width := 0
	if height > 0 {
		width = len(layer.Atoms[0])
	}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			wg.Add(1)
			go func(atomX, atomY int) {
				defer wg.Done()
				atom := &layer.Atoms[atomY][atomX]

				// Phase 1a: Compute local resonance with neighbors
				neighborInfluence := [3]float64{0, 0, 0} // RGB influence
				neighborCount := 0

				directions := [][2]int{
					{-1, -1}, {0, -1}, {1, -1},
					{-1, 0}, {1, 0},
					{-1, 1}, {0, 1}, {1, 1},
				}

				for _, dir := range directions {
					ny, nx := atomY+dir[0], atomX+dir[1]
					if nx >= 0 && nx < width && ny >= 0 && ny < height {
						neighbor := &layer.Atoms[ny][nx]

						// Resonance R(si, sj) = exp(-||si - sj||²/2σ²)
						colorDiff := 0.0
						for i := 0; i < 3; i++ {
							colorDiff += (atom.Color[i] - neighbor.Color[i]) *
								(atom.Color[i] - neighbor.Color[i])
						}
						sigma := net.ResonanceSensitivity
						resonance := math.Exp(-colorDiff / (2 * sigma * sigma))

						weight := atom.ConnectionWeights[neighbor.ID]
						for i := 0; i < 3; i++ {
							neighborInfluence[i] += weight * resonance *
								(neighbor.Color[i] - atom.Color[i])
						}
						neighborCount++
					}
				}

				// Phase 1b: Apply coupling
				if neighborCount > 0 {
					alpha := layer.Alpha
					for i := 0; i < 3; i++ {
						neighborInfluence[i] /= float64(neighborCount)
						atom.Color[i] += alpha * neighborInfluence[i]

						// Clamp to [0, 1]
						if atom.Color[i] > 1.0 {
							atom.Color[i] = 1.0
						} else if atom.Color[i] < 0.0 {
							atom.Color[i] = 0.0
						}
					}
				}

				// Update coherence feature
				if neighborCount > 0 {
					maxDiff := 0.0
					for _, dir := range directions {
						ny, nx := atomY+dir[0], atomX+dir[1]
						if nx >= 0 && nx < width && ny >= 0 && ny < height {
							neighbor := &layer.Atoms[ny][nx]
							for i := 0; i < 3; i++ {
								diff := math.Abs(atom.Color[i] - neighbor.Color[i])
								if diff > maxDiff {
									maxDiff = diff
								}
							}
						}
					}
					atom.Features["coherence"] = 1.0 - maxDiff
				}
			}(x, y)
		}
	}
	wg.Wait()

	layer.Iteration++
}

// PropagatePatterns transfers pattern information between scales
func (net *AtomicImageNetwork) PropagatePatterns(sourceLayer, targetLayer *MultiScaleResonanceLayer) {
	// Upscale patterns from finer resolution
	sourceHeight := len(sourceLayer.Atoms)
	sourceWidth := 0
	if sourceHeight > 0 {
		sourceWidth = len(sourceLayer.Atoms[0])
	}

	targetHeight := len(targetLayer.Atoms)
	targetWidth := 0
	if targetHeight > 0 {
		targetWidth = len(targetLayer.Atoms[0])
	}

	scale := targetLayer.Scale / sourceLayer.Scale

	for sy := 0; sy < sourceHeight; sy++ {
		for sx := 0; sx < sourceWidth; sx++ {
			ty := sy / scale
			tx := sx / scale

			if ty < targetHeight && tx < targetWidth {
				sourceAtom := &sourceLayer.Atoms[sy][sx]
				targetAtom := &targetLayer.Atoms[ty][tx]

				// Mix patterns: target learns from finer scale
				blendFactor := 0.3 // Don't overwrite completely
				for i := 0; i < 3; i++ {
					targetAtom.Color[i] = targetAtom.Color[i]*(1-blendFactor) +
						sourceAtom.Color[i]*blendFactor
				}
				targetAtom.Intensity = targetAtom.Intensity*(1-blendFactor) +
					sourceAtom.Intensity*blendFactor
			}
		}
	}
}

// ========== PHASE 2: SHAPE EMERGENCE ==========

// ResonanceCapsule represents a motif/shape with state vector
type ResonanceCapsule struct {
	ID            int
	BlockID       int             // Which block does this capsule represent
	State         [6]float64      // Texture, intensity, orientation, coherence, entropy, stability
	NeighborCaps  []int           // Neighboring capsule IDs
	CompatWeights map[int]float64 // Compatibility with neighbor capsules
	Energy        float64
	Recognition   float64 // How well recognized (0-1)
}

// ShapeEmergenceEngine manages capsule resonance
type ShapeEmergenceEngine struct {
	Capsules       map[int]*ResonanceCapsule
	Grid           [][]int // Grid of capsule IDs
	Width, Height  int
	Gamma          float64 // Reinforcement factor
	Iterations     int
	ConvergedMasks map[int][]uint8 // Binary masks of recognized shapes
	mutex          sync.RWMutex
}

// NewShapeEmergenceEngine creates a new capsule resonance engine
func NewShapeEmergenceEngine(width, height, blockSize int) *ShapeEmergenceEngine {
	engine := &ShapeEmergenceEngine{
		Capsules:       make(map[int]*ResonanceCapsule),
		Grid:           make([][]int, height/blockSize),
		Width:          width / blockSize,
		Height:         height / blockSize,
		Gamma:          0.15, // Reinforcement factor γ
		Iterations:     0,
		ConvergedMasks: make(map[int][]uint8),
	}

	// Initialize capsules
	for y := 0; y < engine.Height; y++ {
		engine.Grid[y] = make([]int, engine.Width)
		for x := 0; x < engine.Width; x++ {
			id := y*engine.Width + x
			capsule := &ResonanceCapsule{
				ID:            id,
				BlockID:       id,
				State:         [6]float64{0.3, 0.4, 0, 0.2, 0.1, 0.5}, // Initial state
				NeighborCaps:  make([]int, 0, 8),
				CompatWeights: make(map[int]float64),
				Energy:        0.0,
				Recognition:   0.0,
			}

			// Set up neighborhood
			for dy := -1; dy <= 1; dy++ {
				for dx := -1; dx <= 1; dx++ {
					if dx == 0 && dy == 0 {
						continue
					}
					ny, nx := y+dy, x+dx
					if nx >= 0 && nx < engine.Width && ny >= 0 && ny < engine.Height {
						neighborID := ny*engine.Width + nx
						capsule.NeighborCaps = append(capsule.NeighborCaps, neighborID)
						capsule.CompatWeights[neighborID] = 0.5 + (0.5 * (1.0 - float64(dy*dy+dx*dx)/9.0))
					}
				}
			}

			engine.Capsules[id] = capsule
			engine.Grid[y][x] = id
		}
	}

	return engine
}

// PhaseTwo_ShapeEmergence makes primitive shapes appear through capsule resonance
// Implements: s_m(t+1) = s_m(t) + γ * Σ(R(s_n(t), s_m(t)))
func (net *AtomicImageNetwork) PhaseTwo_ShapeEmergence(engine *ShapeEmergenceEngine, iterations int) {
	for iter := 0; iter < iterations; iter++ {
		net.IterateCapsuleResonance(engine)

		// Periodically check for shape emergence
		if iter%10 == 0 {
			net.DetectPrimitiveShapes(engine)
		}
	}

	// Stabilize recognized shapes
	net.StabilizeShapes(engine)
}

// IterateCapsuleResonance performs capsule state updates
func (net *AtomicImageNetwork) IterateCapsuleResonance(engine *ShapeEmergenceEngine) {
	var wg sync.WaitGroup

	for capID, capsule := range engine.Capsules {
		wg.Add(1)
		go func(cid int, cap *ResonanceCapsule) {
			defer wg.Done()

			// Phase 2a: Compute compatibility with neighbors
			neighborCompatibility := [6]float64{}
			for _, neighborID := range cap.NeighborCaps {
				if neighbor, ok := engine.Capsules[neighborID]; ok {
					// Compatibility function R(a, b)
					compat := net.ComputeMotifCompatibility(cap.State, neighbor.State)
					weight := cap.CompatWeights[neighborID]

					for i := 0; i < 6; i++ {
						neighborCompatibility[i] += weight * compat *
							(neighbor.State[i] - cap.State[i])
					}
				}
			}

			// Phase 2b: Update capsule state
			gamma := engine.Gamma
			for i := 0; i < 6; i++ {
				cap.State[i] += gamma * neighborCompatibility[i]
				// Clamp to [0, 1]
				if cap.State[i] > 1.0 {
					cap.State[i] = 1.0
				} else if cap.State[i] < 0.0 {
					cap.State[i] = 0.0
				}
			}

			// Update energy and recognition
			cap.Energy += cap.State[3] * 0.1 // Coherence drives energy
			cap.Recognition = math.Min(1.0, cap.Recognition+cap.State[3]*0.05)
		}(capID, capsule)
	}
	wg.Wait()

	engine.Iterations++
}

// ComputeMotifCompatibility calculates R(a, b) - similarity between motif states
func (net *AtomicImageNetwork) ComputeMotifCompatibility(stateA, stateB [6]float64) float64 {
	// Multi-dimensional Gaussian
	distance := 0.0
	for i := 0; i < 6; i++ {
		distance += (stateA[i] - stateB[i]) * (stateA[i] - stateB[i])
	}
	sigma := 0.15
	compat := math.Exp(-distance / (2 * sigma * sigma))
	return compat
}

// DetectPrimitiveShapes identifies contours, lines, curves from capsule patterns
func (net *AtomicImageNetwork) DetectPrimitiveShapes(engine *ShapeEmergenceEngine) {
	for y := 0; y < engine.Height; y++ {
		for x := 0; x < engine.Width; x++ {
			capID := engine.Grid[y][x]
			capsule := engine.Capsules[capID]

			// Detect edges: high intensity + orientation changes with neighbors
			orientationVariance := 0.0
			for _, neighborID := range capsule.NeighborCaps {
				if neighbor, ok := engine.Capsules[neighborID]; ok {
					orientationVariance +=
						math.Abs(capsule.State[2] - neighbor.State[2]) // orientation
				}
			}

			if capsule.State[1] > 0.6 && orientationVariance > 0.3 {
				capsule.Recognition = math.Min(1.0, capsule.Recognition+0.1)
			}
		}
	}
}

// StabilizeShapes reinforces recognized shapes
func (net *AtomicImageNetwork) StabilizeShapes(engine *ShapeEmergenceEngine) {
	for _, capsule := range engine.Capsules {
		if capsule.Recognition > 0.5 {
			// Stabilize coherence of recognized shapes
			capsule.State[3] = math.Min(1.0, capsule.State[3]+0.2) // coherence
			capsule.State[5] = math.Min(1.0, capsule.State[5]+0.1) // stability
		}
	}
}

// ========== PHASE 3: PROMPT CONDITIONING ==========

// PromptGuideVector encodes prompt instructions as local modifications
type PromptGuideVector struct {
	DirectionX, DirectionY float64 // Spatial guidance
	ColorTarget            [3]float64
	Intensity              float64
	FeatureWeights         map[string]float64
	Strength               float64 // β coefficient
	SpatialDecay           float64 // How strength decays from focus points
}

// PhaseThree_PromptConditioning applies user guidance to pixels/blocks
// Implements: c_ij(t+1) = c_ij(t) + α*...+ β*G_ij(P)
func (net *AtomicImageNetwork) PhaseThree_PromptConditioning(
	promptGuide PromptGuideVector,
	iterations int,
) {
	beta := promptGuide.Strength

	for iter := 0; iter < iterations; iter++ {
		var wg sync.WaitGroup
		gridHeight := len(net.Atoms)
		gridWidth := 0
		if gridHeight > 0 {
			gridWidth = len(net.Atoms[0])
		}

		for y := 0; y < gridHeight; y++ {
			for x := 0; x < gridWidth; x++ {
				wg.Add(1)
				go func(atomX, atomY int) {
					defer wg.Done()
					atom := &net.Atoms[atomY][atomX]

					// Compute spatial guidance strength (stronger near directional focus)
					spatialStrength := beta * math.Exp(
						-(math.Pow(float64(atomX)-promptGuide.DirectionX, 2)+
							math.Pow(float64(atomY)-promptGuide.DirectionY, 2))/
							(2*promptGuide.SpatialDecay*promptGuide.SpatialDecay),
					)

					// Apply color pull
					for i := 0; i < 3; i++ {
						atom.Color[i] += spatialStrength *
							(promptGuide.ColorTarget[i] - atom.Color[i])

						if atom.Color[i] > 1.0 {
							atom.Color[i] = 1.0
						} else if atom.Color[i] < 0.0 {
							atom.Color[i] = 0.0
						}
					}

					// Apply intensity guidance
					atom.Intensity += spatialStrength *
						(promptGuide.Intensity - atom.Intensity)
					if atom.Intensity > 1.0 {
						atom.Intensity = 1.0
					} else if atom.Intensity < 0.0 {
						atom.Intensity = 0.0
					}
				}(x, y)
			}
		}
		wg.Wait()
	}
}

// ========== PHASE 4: ITERATIVE REFINEMENT ==========

// PhaseF our_IterativeRefinement adds fine details and realistic texture
// Implements: c_ij(n+1) = c_ij(n) + δ*Laplacian(c_ij(n)) + ϵ*NoiseAdjust(c_ij(n))
func (net *AtomicImageNetwork) PhaseFour_IterativeRefinement(
	laplacianStrength float64,
	noiseStrength float64,
	iterations int,
) {
	gridHeight := len(net.Atoms)
	gridWidth := 0
	if gridHeight > 0 {
		gridWidth = len(net.Atoms[0])
	}

	for iter := 0; iter < iterations; iter++ {
		// Compute Laplacian for smoothing
		laplacianMap := make([][]float64, gridHeight)
		for y := 0; y < gridHeight; y++ {
			laplacianMap[y] = make([]float64, gridWidth)
		}

		// Compute Laplacian: ∇²c = c_xx + c_yy
		for y := 1; y < gridHeight-1; y++ {
			for x := 1; x < gridWidth-1; x++ {
				laplacian := [3]float64{}
				for i := 0; i < 3; i++ {
					laplacian[i] = (net.Atoms[y-1][x].Color[i] +
						net.Atoms[y+1][x].Color[i] +
						net.Atoms[y][x-1].Color[i] +
						net.Atoms[y][x+1].Color[i] -
						4*net.Atoms[y][x].Color[i])
				}
				// Average across RGB
				laplacianMap[y][x] = (laplacian[0] + laplacian[1] + laplacian[2]) / 3.0
			}
		}

		// Apply Laplacian and noise adjustment
		for y := 0; y < gridHeight; y++ {
			for x := 0; x < gridWidth; x++ {
				atom := &net.Atoms[y][x]

				// Phase 4a: Smooth transitions
				if y > 0 && y < gridHeight-1 && x > 0 && x < gridWidth-1 {
					for i := 0; i < 3; i++ {
						atom.Color[i] += laplacianStrength * laplacianMap[y][x]
					}
				}

				// Phase 4b: Add texture variation
				noiseValue := (math.Sin(float64(iter)*0.3+float64(x)*0.7) *
					math.Cos(float64(y)*0.5)) * noiseStrength
				for i := 0; i < 3; i++ {
					atom.Color[i] += noiseValue * 0.1

					// Clamp
					if atom.Color[i] > 1.0 {
						atom.Color[i] = 1.0
					} else if atom.Color[i] < 0.0 {
						atom.Color[i] = 0.0
					}
				}
			}
		}
	}
}

// ========== PHASE 5: COHERENCE VERIFICATION ==========

// CoherenceReport contains verification results
type CoherenceReport struct {
	GlobalCoherence    float64
	CoherenceMap       [][]float64
	FaultyAtoms        []CoherenceFault
	RepairCount        int
	OverallHealthScore float64
}

// CoherenceFault represents a pixel with low coherence
type CoherenceFault struct {
	X, Y      int
	Coherence float64
	DiffScore float64
}

// PhaseFive_CoherenceVerification checks and corrects coherence issues
// Implements: coherence_ij = 1 - Σ||c_ij - c_k||/max_possible_diff
func (net *AtomicImageNetwork) PhaseFive_CoherenceVerification() *CoherenceReport {
	report := &CoherenceReport{
		CoherenceMap: make([][]float64, len(net.Atoms)),
		FaultyAtoms:  make([]CoherenceFault, 0),
	}

	gridHeight := len(net.Atoms)
	gridWidth := 0
	if gridHeight > 0 {
		gridWidth = len(net.Atoms[0])
	}

	totalCoherence := 0.0
	faultCount := 0

	for y := 0; y < gridHeight; y++ {
		report.CoherenceMap[y] = make([]float64, gridWidth)

		for x := 0; x < gridWidth; x++ {
			atom := &net.Atoms[y][x]

			// Compute local coherence
			diffSum := 0.0
			maxPossibleDiff := 3.0 * math.Sqrt(3) // Max RGB distance

			directions := [][2]int{
				{-1, -1}, {0, -1}, {1, -1},
				{-1, 0}, {1, 0},
				{-1, 1}, {0, 1}, {1, 1},
			}

			count := 0
			for _, dir := range directions {
				ny, nx := y+dir[0], x+dir[1]
				if nx >= 0 && nx < gridWidth && ny >= 0 && ny < gridHeight {
					neighbor := &net.Atoms[ny][nx]
					colorDiff := 0.0
					for i := 0; i < 3; i++ {
						colorDiff += math.Abs(atom.Color[i] - neighbor.Color[i])
					}
					diffSum += colorDiff
					count++
				}
			}

			if count > 0 {
				diffSum /= float64(count)
			}

			coherence := 1.0 - (diffSum / maxPossibleDiff)
			if coherence < 0.0 {
				coherence = 0.0
			}

			report.CoherenceMap[y][x] = coherence
			totalCoherence += coherence

			// Mark low-coherence atoms
			if coherence < 0.3 {
				faultCount++
				report.FaultyAtoms = append(report.FaultyAtoms, CoherenceFault{
					X:         x,
					Y:         y,
					Coherence: coherence,
					DiffScore: diffSum,
				})
			}
		}
	}

	report.GlobalCoherence = totalCoherence / float64(gridHeight*gridWidth)
	report.OverallHealthScore = 1.0 - float64(faultCount)/float64(gridHeight*gridWidth)

	// Repair low-coherence atoms
	report.RepairCount = net.RepairCoherentAtoms(report.FaultyAtoms)

	return report
}

// RepairCoherentAtoms reapplies resonance to fix low-coherence regions
func (net *AtomicImageNetwork) RepairCoherentAtoms(faults []CoherenceFault) int {
	repaired := 0

	for _, fault := range faults {
		atom := &net.Atoms[fault.Y][fault.X]

		// Reapply local resonance (Phase 1 technique)
		directions := [][2]int{
			{-1, -1}, {0, -1}, {1, -1},
			{-1, 0}, {1, 0},
			{-1, 1}, {0, 1}, {1, 1},
		}

		avgColor := [3]float64{}
		count := 0

		for _, dir := range directions {
			ny, nx := fault.Y+dir[0], fault.X+dir[1]
			if nx >= 0 && nx < len(net.Atoms[0]) && ny >= 0 && ny < len(net.Atoms) {
				neighbor := &net.Atoms[ny][nx]
				for i := 0; i < 3; i++ {
					avgColor[i] += neighbor.Color[i]
				}
				count++
			}
		}

		if count > 0 {
			// Blend toward neighbor average
			blendFactor := 0.4
			for i := 0; i < 3; i++ {
				avgColor[i] /= float64(count)
				atom.Color[i] = atom.Color[i]*(1-blendFactor) + avgColor[i]*blendFactor
			}
			repaired++
		}
	}

	return repaired
}

// ========== COMPLETE PIPELINE ==========

// FullImageGenerationPipeline runs all 5 phases sequentially
func (net *AtomicImageNetwork) FullImageGenerationPipeline(
	prompt string,
	totalIterations int,
) *CoherenceReport {
	// Phase 1: Multi-scale structuration
	net.PhaseOne_StructurationMultiEchelle(totalIterations / 5)

	// Phase 2: Shape emergence
	engine := NewShapeEmergenceEngine(net.Width, net.Height, net.PatchSize)
	net.PhaseTwo_ShapeEmergence(engine, totalIterations/5)

	// Phase 3: Prompt conditioning
	guide := net.ParsePromptToGuide(prompt)
	net.PhaseThree_PromptConditioning(guide, totalIterations/5)

	// Phase 4: Iterative refinement
	net.PhaseFour_IterativeRefinement(0.1, 0.08, totalIterations/5)

	// Phase 5: Coherence verification
	report := net.PhaseFive_CoherenceVerification()

	return report
}

// ParsePromptToGuide converts a text prompt to a PromptGuideVector
func (net *AtomicImageNetwork) ParsePromptToGuide(prompt string) PromptGuideVector {
	guide := PromptGuideVector{
		DirectionX:     float64(net.Width/net.PatchSize) / 2.0,
		DirectionY:     float64(net.Height/net.PatchSize) / 2.0,
		ColorTarget:    [3]float64{0.5, 0.5, 0.5},
		Intensity:      0.5,
		Strength:       0.3,
		SpatialDecay:   10.0,
		FeatureWeights: make(map[string]float64),
	}

	// Parse prompt (using existing ParsePrompt logic)
	net.ParsePrompt(prompt)

	// Apply style vector as color target
	guide.ColorTarget = net.StyleVector

	// Extract brightness
	if intensity, ok := net.PromptEmbedding["bright"]; ok {
		guide.Intensity = 0.7 + intensity*0.1
	} else if intensity, ok := net.PromptEmbedding["dark"]; ok {
		guide.Intensity = 0.3 + intensity*0.1
	}

	return guide
}

// ========== VISUALIZATION & DEBUGGING ==========

// RenderCoherenceMap creates an image showing coherence levels
func (report *CoherenceReport) RenderCoherenceMap(width, height int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < len(report.CoherenceMap); y++ {
		for x := 0; x < len(report.CoherenceMap[y]); x++ {
			coherence := report.CoherenceMap[y][x]

			// Map coherence to color: red (low) → green (high)
			r := uint8(math.Max(0, (1.0-coherence)*255))
			g := uint8(math.Max(0, coherence*255))
			b := uint8(100)

			// Fill block
			img.Set(x, y, color.RGBA{r, g, b, 255})
		}
	}

	return img
}
