// Package database - Image Atomic Generation Module
// Implements Atomic Resonance Technology (T.R.A.) for image generation
// Each pixel/patch is treated as a computational atom that interacts locally
//
// Core principle: Intelligence emerges from local atomic interactions,
// guided by prompts translated into local constraints

package database

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"math/rand"
	"os"
	"strings"
	"sync"
)

// PixelAtom represents a pixel or patch as a computational atom
type PixelAtom struct {
	ID                 int                // Unique identifier
	X, Y               int                // Position in image
	Color              [3]float64         // RGB normalized [0, 1]
	Intensity          float64            // Brightness/intensity value
	Features           map[string]float64 // Feature vector (texture, orientation, etc.)
	State              float64            // Internal state of atom
	Neighbors          []int              // IDs of neighboring atoms
	NeighborStates     map[int]float64    // States of neighbors (for update)
	ConnectionWeights  map[int]float64    // Weights wij to neighbors
	ExternalConstraint float64            // Constraint from prompt (0-1)
	ConstraintType     string             // Type of constraint ("color", "brightness", "texture")
	EnergyConsumption  float64            // Energy used by this atom
	IsFrozen           bool               // Is atom in low-energy freeze state
	FreezeCount        int                // Iterations at low activity
	mutex              sync.Mutex         // Thread-safe operations
}

// AtomicImageNetwork represents the complete image generation network
type AtomicImageNetwork struct {
	Atoms                 [][]PixelAtom      // 2D grid of atoms
	Width, Height         int                // Image dimensions
	PatchSize             int                // Size of each patch (8, 16, etc.)
	CouplingCoefficient   float64            // α - Influence of neighbors
	LocalRulesCoefficient float64            // β - Impact of local constraints
	ReinforcementFactor   float64            // γ - Reinforcement strength
	DecayFactor           float64            // δ - Connection weight decay
	ResonanceSensitivity  float64            // σ - Resonance mechanism sensitivity
	FreezeThreshold       float64            // ϵ - Activity threshold for freeze
	FreezeIterations      int                // T - Iterations before freeze
	GlobalIteration       int                // Current iteration counter
	TotalEnergy           float64            // Total energy consumed
	FrozenAtomsCount      int                // Number of frozen atoms
	PromptEmbedding       map[string]float64 // Parsed prompt constraints
	StyleVector           [3]float64         // Global style influence (R, G, B)
	mutex                 sync.RWMutex       // Network-level locking
}

// PromptConstraint defines how prompt terms map to local atom modifications
type PromptConstraint struct {
	Term              string                 // Prompt term (e.g., "red", "night", "rough")
	ConstraintType    string                 // "color", "brightness", "texture"
	ColorModifier     [3]float64             // RGB modification [-1, 1]
	IntensityModifier float64                // Brightness modification [-1, 1]
	SpatialMask       func(x, y int) float64 // Spatial application (0-1)
	Strength          float64                // Overall strength [0, 1]
}

// NewPixelAtom creates a new pixel atom
func NewPixelAtom(id, x, y int, patchSize int) *PixelAtom {
	atom := &PixelAtom{
		ID:                 id,
		X:                  x,
		Y:                  y,
		Color:              [3]float64{rand.Float64(), rand.Float64(), rand.Float64()},
		Intensity:          rand.Float64() * 0.5,
		Features:           make(map[string]float64),
		State:              rand.Float64() * 0.3,
		Neighbors:          make([]int, 0, 8),
		NeighborStates:     make(map[int]float64),
		ConnectionWeights:  make(map[int]float64),
		ExternalConstraint: 0.0,
		ConstraintType:     "none",
		EnergyConsumption:  0.0,
		IsFrozen:           false,
		FreezeCount:        0,
	}

	// Initialize features
	atom.Features["texture"] = rand.Float64()
	atom.Features["orientation"] = rand.Float64() * 360
	atom.Features["coherence"] = 0.5

	return atom
}

// NewAtomicImageNetwork initializes a new atomic image generation network
func NewAtomicImageNetwork(width, height, patchSize int) *AtomicImageNetwork {
	network := &AtomicImageNetwork{
		Atoms:                 make([][]PixelAtom, height/patchSize),
		Width:                 width,
		Height:                height,
		PatchSize:             patchSize,
		CouplingCoefficient:   0.7,  // α - Neighbor influence
		LocalRulesCoefficient: 0.3,  // β - Local constraints
		ReinforcementFactor:   0.15, // γ - Connection reinforcement
		DecayFactor:           0.05, // δ - Weight decay
		ResonanceSensitivity:  0.1,  // σ - Resonance sensitivity
		FreezeThreshold:       0.3,  // ϵ - Freeze activation threshold
		FreezeIterations:      2,    // T - Iterations before freeze
		GlobalIteration:       0,
		TotalEnergy:           0.0,
		FrozenAtomsCount:      0,
		PromptEmbedding:       make(map[string]float64),
		StyleVector:           [3]float64{1.0, 1.0, 1.0}, // Neutral style initially
	}

	// Initialize 2D grid of atoms
	gridHeight := height / patchSize
	gridWidth := width / patchSize

	for y := 0; y < gridHeight; y++ {
		network.Atoms[y] = make([]PixelAtom, gridWidth)
		for x := 0; x < gridWidth; x++ {
			id := y*gridWidth + x
			network.Atoms[y][x] = *NewPixelAtom(id, x, y, patchSize)
		}
	}

	// Connect neighbors in 2D grid
	for y := 0; y < gridHeight; y++ {
		for x := 0; x < gridWidth; x++ {
			network.connectNeighbors(x, y, gridWidth, gridHeight)
		}
	}

	return network
}

// connectNeighbors establishes 8-neighborhood connections (or 4-neighborhood at edges)
func (net *AtomicImageNetwork) connectNeighbors(x, y, gridWidth, gridHeight int) {
	directions := [][2]int{
		{-1, -1}, {0, -1}, {1, -1},
		{-1, 0}, {1, 0},
		{-1, 1}, {0, 1}, {1, 1},
	}

	for _, dir := range directions {
		nx, ny := x+dir[0], y+dir[1]
		if nx >= 0 && nx < gridWidth && ny >= 0 && ny < gridHeight {
			neighborID := ny*gridWidth + nx

			net.Atoms[y][x].Neighbors = append(net.Atoms[y][x].Neighbors, neighborID)
			net.Atoms[y][x].ConnectionWeights[neighborID] = rand.Float64() * 0.5
		}
	}
}

// ParsePrompt translates a text prompt into atomic constraints
// Maps semantic meaning to local modifications
func (net *AtomicImageNetwork) ParsePrompt(prompt string) {
	promptLower := strings.ToLower(prompt)

	// Color mappings
	colorMappings := map[string][3]float64{
		"red":    {0.8, -0.4, -0.4},
		"blue":   {-0.4, -0.4, 0.8},
		"green":  {-0.4, 0.8, -0.4},
		"yellow": {0.8, 0.8, -0.4},
		"purple": {0.8, -0.4, 0.8},
		"orange": {0.8, 0.3, -0.4},
		"pink":   {0.8, -0.2, 0.4},
		"cyan":   {-0.4, 0.8, 0.8},
	}

	// Brightness mappings
	brightnessTerms := map[string]float64{
		"dark":    -0.6,
		"night":   -0.6,
		"bright":  0.6,
		"light":   0.6,
		"sunny":   0.7,
		"dim":     -0.4,
		"glowing": 0.8,
	}

	// Texture/style mappings
	textureMappings := map[string]float64{
		"rough":    0.8,
		"smooth":   -0.7,
		"detailed": 0.6,
		"blurry":   -0.5,
		"sharp":    0.7,
		"noisy":    0.5,
		"clean":    -0.4,
	}

	// Extract colors
	for color, modifier := range colorMappings {
		if strings.Contains(promptLower, color) {
			net.StyleVector[0] += modifier[0] / 3
			net.StyleVector[1] += modifier[1] / 3
			net.StyleVector[2] += modifier[2] / 3
		}
	}

	// Extract brightness modifiers
	for term, intensity := range brightnessTerms {
		if strings.Contains(promptLower, term) {
			net.PromptEmbedding[term] = intensity
		}
	}

	// Extract texture modifiers
	for term, strength := range textureMappings {
		if strings.Contains(promptLower, term) {
			net.PromptEmbedding[term] = strength
		}
	}

	// Clamp style vector to [-1, 1]
	for i := 0; i < 3; i++ {
		if net.StyleVector[i] > 1.0 {
			net.StyleVector[i] = 1.0
		} else if net.StyleVector[i] < -1.0 {
			net.StyleVector[i] = -1.0
		}
	}
}

// ComputeExternalConstraint calculates prompt influence on an atom at (x, y)
func (net *AtomicImageNetwork) ComputeExternalConstraint(x, y int) float64 {
	// Simple spatial influence: stronger at edges for prompts like "frame" or "border"
	gridWidth := net.Width / net.PatchSize
	gridHeight := net.Height / net.PatchSize

	// Distance to center (0-1, 0 at center, 1 at corners)
	centerX := float64(gridWidth) / 2.0
	centerY := float64(gridHeight) / 2.0
	distToCenter := math.Sqrt(
		math.Pow(float64(x)-centerX, 2)+math.Pow(float64(y)-centerY, 2),
	) / math.Sqrt(centerX*centerX+centerY*centerY)

	// For now, return normalized constraint
	// This can be extended with spatial masks
	return math.Min(1.0, math.Max(0.0, 0.5+distToCenter*0.5))
}

// UpdateAtomState updates a single atom's state based on neighbors and constraints
// Implements: a_i(t+1) = a_i(t) + α*Σ(w_ij*f(a_i,a_j,c)) + β*g(constraints)
func (net *AtomicImageNetwork) UpdateAtomState(atom *PixelAtom, alpha, beta float64) {
	atom.mutex.Lock()
	defer atom.mutex.Unlock()

	// If frozen, minimal update for energy savings
	if atom.IsFrozen {
		// Very slow decay
		atom.State *= 0.98
		atom.Intensity *= 0.99
		atom.EnergyConsumption += 0.01
		return
	}

	// Phase 1: Gather neighbor influences
	gridWidth := net.Width / net.PatchSize
	gridHeight := net.Height / net.PatchSize
	neighborInfluence := 0.0
	neighborCount := 0

	for _, neighborID := range atom.Neighbors {
		ny := neighborID / gridWidth
		nx := neighborID % gridWidth

		if ny >= 0 && ny < gridHeight && nx >= 0 && nx < gridWidth {
			neighbor := &net.Atoms[ny][nx]
			neighborState := neighbor.State
			resonance := atom.computePixelResonance(neighborState, net.ResonanceSensitivity)
			weight := atom.ConnectionWeights[neighborID]

			neighborInfluence += weight * resonance * (neighborState - atom.State)
			neighborCount++
		}
	}

	if neighborCount > 0 {
		neighborInfluence /= float64(neighborCount)
	}

	// Phase 2: Apply external constraints from prompt
	constraintInfluence := 0.0
	if atom.ExternalConstraint > 0.0 {
		// Constraint pulls atom state toward the specified value
		constraintInfluence = atom.ExternalConstraint * (0.5 - atom.State)
	}

	// Phase 3: Apply color modifications based on style
	colorChange := 0.0
	if intensity, ok := net.PromptEmbedding["bright"]; ok {
		colorChange += intensity * 0.1
	}
	if intensity, ok := net.PromptEmbedding["dark"]; ok {
		colorChange += intensity * 0.1
	}

	// Phase 4: Update state and color
	oldState := atom.State
	atom.State += alpha*neighborInfluence + beta*constraintInfluence
	atom.Intensity += colorChange

	// Clamp to valid ranges [0, 1]
	if atom.State > 1.0 {
		atom.State = 1.0
	} else if atom.State < 0.0 {
		atom.State = 0.0
	}

	if atom.Intensity > 1.0 {
		atom.Intensity = 1.0
	} else if atom.Intensity < 0.0 {
		atom.Intensity = 0.0
	}

	// Calculate energy consumption
	stateChange := math.Abs(atom.State - oldState)
	atom.EnergyConsumption += stateChange * 0.5

	// Update connection weights via: dw_ij/dt = γ*coherence - δ*w_ij
	coherence := 1.0 - stateChange
	for neighborID := range atom.ConnectionWeights {
		atom.ConnectionWeights[neighborID] =
			atom.ConnectionWeights[neighborID] +
				net.ReinforcementFactor*coherence -
				net.DecayFactor*atom.ConnectionWeights[neighborID]

		// Clamp weights
		if atom.ConnectionWeights[neighborID] > 1.0 {
			atom.ConnectionWeights[neighborID] = 1.0
		} else if atom.ConnectionWeights[neighborID] < 0.0 {
			atom.ConnectionWeights[neighborID] = 0.0
		}
	}

	// Check freeze condition
	if stateChange < net.FreezeThreshold {
		atom.FreezeCount++
		if atom.FreezeCount >= net.FreezeIterations {
			atom.IsFrozen = true
			net.FrozenAtomsCount++
		}
	} else {
		atom.FreezeCount = 0
		if atom.IsFrozen {
			atom.IsFrozen = false
			net.FrozenAtomsCount--
		}
	}
}

// computePixelResonance calculates resonance between pixels: R(s_i, s_j) = exp(-||s_i-s_j||²/2σ²)
func (atom *PixelAtom) computePixelResonance(neighborState, sigma float64) float64 {
	diff := atom.State - neighborState
	distSquared := diff * diff
	exponent := -distSquared / (2 * sigma * sigma)
	resonance := math.Exp(exponent)

	if resonance > 1.0 {
		resonance = 1.0
	} else if resonance < 0.0 {
		resonance = 0.0
	}

	return resonance
}

// IterateGeneration performs one generation iteration over all atoms
func (net *AtomicImageNetwork) IterateGeneration() {
	net.mutex.Lock()
	net.GlobalIteration++
	net.mutex.Unlock()

	gridHeight := net.Height / net.PatchSize
	gridWidth := net.Width / net.PatchSize

	// Update all atoms asynchronously
	var wg sync.WaitGroup
	for y := 0; y < gridHeight; y++ {
		for x := 0; x < gridWidth; x++ {
			wg.Add(1)
			go func(atomX, atomY int) {
				defer wg.Done()
				net.UpdateAtomState(
					&net.Atoms[atomY][atomX],
					net.CouplingCoefficient,
					net.LocalRulesCoefficient,
				)
			}(x, y)
		}
	}
	wg.Wait()

	// Accumulate total energy
	totalEnergyDelta := 0.0
	for y := 0; y < gridHeight; y++ {
		for x := 0; x < gridWidth; x++ {
			totalEnergyDelta += net.Atoms[y][x].EnergyConsumption
		}
	}
	net.TotalEnergy += totalEnergyDelta
}

// GenerateImage runs the complete generation pipeline
func (net *AtomicImageNetwork) GenerateImage(iterations int) image.Image {
	// Main generation loop
	for i := 0; i < iterations; i++ {
		net.IterateGeneration()

		// Optional: Adaptive learning rate based on convergence
		if i%100 == 0 && i > 0 {
			energyDensity := net.TotalEnergy / float64(net.Width*net.Height)
			if energyDensity < 0.01 {
				// System is very stable, can reduce coupling to preserve state
				net.CouplingCoefficient *= 0.98
			}
		}
	}

	return net.RenderImage()
}

// RenderImage converts the atomic network state to a PNG image
func (net *AtomicImageNetwork) RenderImage() image.Image {
	img := image.NewRGBA(image.Rect(0, 0, net.Width, net.Height))

	gridHeight := net.Height / net.PatchSize
	gridWidth := net.Width / net.PatchSize

	// Fill image with atom colors, expanded to patch size
	for y := 0; y < gridHeight; y++ {
		for x := 0; x < gridWidth; x++ {
			atom := &net.Atoms[y][x]

			// Convert atom state and color to RGB
			r := uint8(math.Min(255, (atom.Color[0]+atom.State)*128))
			g := uint8(math.Min(255, (atom.Color[1]+atom.State)*128))
			b := uint8(math.Min(255, (atom.Color[2]+atom.State)*128))

			// Apply intensity as brightness
			brightness := uint8(atom.Intensity * 255)
			r = (r + brightness) / 2
			g = (g + brightness) / 2
			b = (b + brightness) / 2

			// Fill patch area
			for py := 0; py < net.PatchSize; py++ {
				for px := 0; px < net.PatchSize; px++ {
					imgX := x*net.PatchSize + px
					imgY := y*net.PatchSize + py
					if imgX < net.Width && imgY < net.Height {
						img.Set(imgX, imgY, color.RGBA{r, g, b, 255})
					}
				}
			}
		}
	}

	return img
}

// SaveImage saves the generated image to file
func (net *AtomicImageNetwork) SaveImage(filename string) error {
	img := net.RenderImage()

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	return png.Encode(file, img)
}

// LocalSmoothing applies post-processing local smoothing to the image
func (net *AtomicImageNetwork) LocalSmoothing(radius int) {
	gridHeight := net.Height / net.PatchSize
	gridWidth := net.Width / net.PatchSize

	// Create temporary copy
	tempAtoms := make([][]PixelAtom, gridHeight)
	for y := 0; y < gridHeight; y++ {
		tempAtoms[y] = make([]PixelAtom, gridWidth)
		copy(tempAtoms[y], net.Atoms[y])
	}

	// Apply local averaging
	for y := 0; y < gridHeight; y++ {
		for x := 0; x < gridWidth; x++ {
			var sumR, sumG, sumB, sumIntensity float64
			count := 0

			for dy := -radius; dy <= radius; dy++ {
				for dx := -radius; dx <= radius; dx++ {
					ny, nx := y+dy, x+dx
					if nx >= 0 && nx < gridWidth && ny >= 0 && ny < gridHeight {
						sumR += tempAtoms[ny][nx].Color[0]
						sumG += tempAtoms[ny][nx].Color[1]
						sumB += tempAtoms[ny][nx].Color[2]
						sumIntensity += tempAtoms[ny][nx].Intensity
						count++
					}
				}
			}

			if count > 0 {
				net.Atoms[y][x].Color[0] = sumR / float64(count)
				net.Atoms[y][x].Color[1] = sumG / float64(count)
				net.Atoms[y][x].Color[2] = sumB / float64(count)
				net.Atoms[y][x].Intensity = sumIntensity / float64(count)
			}
		}
	}
}

// EdgeEnhancement applies local edge sharpening
func (net *AtomicImageNetwork) EdgeEnhancement(strength float64) {
	gridHeight := net.Height / net.PatchSize
	gridWidth := net.Width / net.PatchSize

	for y := 1; y < gridHeight-1; y++ {
		for x := 1; x < gridWidth-1; x++ {
			atom := &net.Atoms[y][x]

			// Compute local gradient
			var gradX, gradY float64
			for i := 0; i < 3; i++ {
				gradX += (net.Atoms[y][x+1].Color[i] - net.Atoms[y][x-1].Color[i]) / 2.0
				gradY += (net.Atoms[y+1][x].Color[i] - net.Atoms[y-1][x].Color[i]) / 2.0
			}

			magnitude := math.Sqrt(gradX*gradX + gradY*gradY)

			// Enhance edges
			if magnitude > 0.3 {
				for i := 0; i < 3; i++ {
					atom.Color[i] += strength * (0.5 - atom.Color[i]) * magnitude
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

// GetNetworkStats returns statistics about the current network state
func (net *AtomicImageNetwork) GetNetworkStats() map[string]float64 {
	gridHeight := net.Height / net.PatchSize
	gridWidth := net.Width / net.PatchSize

	stats := make(map[string]float64)

	var totalState, totalIntensity, activeAtoms float64
	minState, maxState := 1.0, 0.0

	for y := 0; y < gridHeight; y++ {
		for x := 0; x < gridWidth; x++ {
			atom := &net.Atoms[y][x]
			totalState += atom.State
			totalIntensity += atom.Intensity

			if atom.State > 0.1 {
				activeAtoms++
			}

			if atom.State < minState {
				minState = atom.State
			}
			if atom.State > maxState {
				maxState = atom.State
			}
		}
	}

	totalAtoms := float64(gridWidth * gridHeight)
	stats["avg_state"] = totalState / totalAtoms
	stats["avg_intensity"] = totalIntensity / totalAtoms
	stats["active_atoms"] = activeAtoms
	stats["active_percent"] = (activeAtoms / totalAtoms) * 100
	stats["min_state"] = minState
	stats["max_state"] = maxState
	stats["frozen_atoms"] = float64(net.FrozenAtomsCount)
	stats["total_energy"] = net.TotalEnergy
	stats["energy_density"] = net.TotalEnergy / (float64(net.Width) * float64(net.Height))
	stats["iteration"] = float64(net.GlobalIteration)

	return stats
}
