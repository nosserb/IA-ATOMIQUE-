// Package database - Energy-Based Image Generation (Constraint Relaxation)
// Instead of "drawing", we relax a system to equilibrium under constraints
// Physics-inspired: each atom minimizes local tension, global coherence emerges

package database

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"math/rand"
	"os"
	"sync"
)

// ============================================================================
// LEVEL 1: PIXEL ATOMS (Elementary units with local state)
// ============================================================================

type PixelAtomV2 struct {
	// Position
	X, Y int

	// Color state [0, 1]
	R, G, B float64

	// Local properties
	Intensity    float64 // Overall brightness
	Orientation  float64 // Local gradient direction (0-2π)
	Confidence   float64 // How stable is this pixel? (0-1)
	TextureIndex float64 // Type of texture (0-1)

	// Energy tracking
	LocalEnergy     float64 // Dissonance with neighbors
	CoherenceError  float64 // Deviation from global pattern
	StabilityScore  float64 // Positive=stable, Negative=oscillating
	LastStateEnergy float64 // Previous energy (for trend detection)
	EnergyTrend     float64 // Is energy rising or falling?

	// Neighbors cache
	Neighbors [8]*PixelAtomV2 // Direct 8-neighbor pointers
	IsEdge    bool            // Is on image boundary?

	// Threading
	mutex sync.Mutex
}

// ============================================================================
// LEVEL 2: PATTERN DETECTION (Local structures)
// ============================================================================

type PatternRegion struct {
	// Bounding box
	X1, Y1, X2, Y2 int

	// Detected properties
	EdgeStrength      float64    // How sharp are boundaries?
	GradientDirection float64    // Average gradient direction
	TextureType       string     // "smooth", "rough", "detailed"
	DominantColor     [3]float64 // Average RGB
	LocalSymmetry     float64    // Symmetry score (0-1)
	Coherence         float64    // Internal consistency

	// Atoms in this region
	Atoms []*PixelAtomV2
}

// ============================================================================
// LEVEL 3: GLOBAL COHERENCE FIELD (Weak influence on all atoms)
// ============================================================================

type GlobalCoherenceField struct {
	// This field doesn't command directly, just influences
	// Think of it as "pressure" that atoms respond to

	// Global properties
	AverageBrightness    float64
	DominantGradientDir  float64
	GlobalSymmetryTarget float64 // Target symmetry level
	TextureConsistency   float64 // How uniform should textures be?
	ShadowDirection      float64 // Where should shadows go?
	EdgeCohesion         float64 // How sharp should edges be?

	// Influence strength (weak, not dominant)
	InfluenceWeight float64 // Typically 0.05-0.15

	// Symmetry breaking: directional field (low resolution guide)
	DirectionalField [][]float64 // Low-res guide: orientation field
	FieldStrength    float64     // How much to follow directional field (0.05-0.15)

	// Violation penalties
	SymmetryViolation float64
	TextureViolation  float64
	LightingViolation float64
	ContourViolation  float64
}

// ============================================================================
// IMAGE ENERGY SIGNATURE EXTRACTION
// ============================================================================

type ImageEnergyProfile struct {
	// Extracted energy coefficients from target image
	// These define the "physical signature" to match

	// Weights for energy terms
	LambdaGradient  float64 // Gradient energy: ||∇I||²
	LambdaLocal     float64 // Local coherence: ||I(x) - mean(neighbors)||²
	LambdaTexture   float64 // Texture/frequency: FFT or variance
	LambdaScale     float64 // Multi-scale distribution
	LambdaSmoothing float64 // Overall smoothness penalty

	// Computed statistics
	AverageGradient       float64 // Mean gradient magnitude
	LocalCoherence        float64 // Mean local coherence
	TextureLevel          float64 // Texture amount (0=smooth, 1=rough)
	SharpnessRatio        float64 // Ratio of sharp to smooth regions
	FrequencyDistribution float64 // Average frequency (0=low, 1=high)

	// Distribution info
	HistogramShape          string  // "uniform", "gaussian", "bimodal", "sparse"
	EdgeDensity             float64 // Fraction of edges (0-1)
	FlatRegionsFraction     float64 // Fraction of low-variance areas
	TexturedRegionsFraction float64 // Fraction of high-variance areas

	// Source info
	SourceImage   string // Path to analyzed image
	Width, Height int
}

// NewImageEnergyProfile analyzes a PNG image and extracts its energy signature
func NewImageEnergyProfile(imagePath string) (*ImageEnergyProfile, error) {
	file, err := os.Open(imagePath)
	if err != nil {
		return nil, fmt.Errorf("cannot open image: %v", err)
	}
	defer file.Close()

	img, err := png.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("cannot decode PNG: %v", err)
	}

	bounds := img.Bounds()
	width := bounds.Max.X - bounds.Min.X
	height := bounds.Max.Y - bounds.Min.Y

	profile := &ImageEnergyProfile{
		SourceImage: imagePath,
		Width:       width,
		Height:      height,
	}

	// Analyze the image
	profile.analyzeEnergySignature(img)

	return profile, nil
}

func (profile *ImageEnergyProfile) analyzeEnergySignature(img image.Image) {
	bounds := img.Bounds()
	width := bounds.Max.X - bounds.Min.X
	height := bounds.Max.Y - bounds.Min.Y

	// Convert to RGBA for processing
	var gradients []float64
	var localCoherences []float64
	var textureValues []float64
	var edgeCount int
	var flatCount int

	// Analyze each pixel
	for y := 1; y < height-1; y++ {
		for x := 1; x < width-1; x++ {
			// Get current and neighbor pixels
			curr := rgbaToFloat(img.At(x, y))

			// COMPUTE GRADIENT ENERGY
			// E_grad = ||∇I||² (Sobel-like)
			left := rgbaToFloat(img.At(x-1, y))
			right := rgbaToFloat(img.At(x+1, y))
			top := rgbaToFloat(img.At(x, y-1))
			bottom := rgbaToFloat(img.At(x, y+1))

			gradX := (right[0] + right[1] + right[2]) - (left[0] + left[1] + left[2])
			gradY := (bottom[0] + bottom[1] + bottom[2]) - (top[0] + top[1] + top[2])
			gradMag := math.Sqrt(gradX*gradX+gradY*gradY) / 3.0
			gradients = append(gradients, gradMag)

			// COMPUTE LOCAL COHERENCE
			// E_local = ||I(x) - mean(neighbors)||²
			var neighborSum [3]float64
			for dy := -1; dy <= 1; dy++ {
				for dx := -1; dx <= 1; dx++ {
					if dx == 0 && dy == 0 {
						continue
					}
					n := rgbaToFloat(img.At(x+dx, y+dy))
					neighborSum[0] += n[0]
					neighborSum[1] += n[1]
					neighborSum[2] += n[2]
				}
			}
			neighborMean := [3]float64{
				neighborSum[0] / 8,
				neighborSum[1] / 8,
				neighborSum[2] / 8,
			}

			localCoh := math.Sqrt(
				math.Pow(curr[0]-neighborMean[0], 2) +
					math.Pow(curr[1]-neighborMean[1], 2) +
					math.Pow(curr[2]-neighborMean[2], 2),
			)
			localCoherences = append(localCoherences, localCoh)

			// TEXTURE VARIANCE (local)
			var pixelVariance float64
			for dy := -1; dy <= 1; dy++ {
				for dx := -1; dx <= 1; dx++ {
					if dx == 0 && dy == 0 {
						continue
					}
					n := rgbaToFloat(img.At(x+dx, y+dy))
					diff := math.Pow(curr[0]-n[0], 2) +
						math.Pow(curr[1]-n[1], 2) +
						math.Pow(curr[2]-n[2], 2)
					pixelVariance += diff
				}
			}
			pixelVariance /= 8.0
			textureValues = append(textureValues, pixelVariance)

			// Classification
			if gradMag > 0.15 {
				edgeCount++
			} else if pixelVariance < 0.01 {
				flatCount++
			}
		}
	}

	// Compute statistics
	if len(gradients) > 0 {
		profile.AverageGradient = mean(gradients)
		profile.SharpnessRatio = float64(edgeCount) / float64(len(gradients))
	}

	if len(localCoherences) > 0 {
		profile.LocalCoherence = mean(localCoherences)
	}

	if len(textureValues) > 0 {
		profile.TextureLevel = mean(textureValues)
	}

	// Determine histogram shape
	if profile.SharpnessRatio > 0.3 {
		profile.HistogramShape = "bimodal" // Mix of edges and flat
	} else if profile.SharpnessRatio < 0.05 {
		profile.HistogramShape = "gaussian" // Smooth gradient
	} else {
		profile.HistogramShape = "uniform"
	}

	profile.EdgeDensity = float64(edgeCount) / float64(len(gradients))
	profile.FlatRegionsFraction = float64(flatCount) / float64(len(gradients))
	profile.TexturedRegionsFraction = 1.0 - profile.FlatRegionsFraction

	// Set energy lambdas based on analysis
	// Strongly weighted by what we found
	profile.LambdaGradient = profile.AverageGradient * 0.3
	profile.LambdaLocal = profile.LocalCoherence * 0.25
	profile.LambdaTexture = profile.TextureLevel * 0.2
	profile.LambdaScale = profile.SharpnessRatio * 0.25
	profile.LambdaSmoothing = (1.0 - profile.SharpnessRatio) * 0.1

	// Normalize
	total := profile.LambdaGradient + profile.LambdaLocal + profile.LambdaTexture +
		profile.LambdaScale + profile.LambdaSmoothing
	if total > 0 {
		profile.LambdaGradient /= total
		profile.LambdaLocal /= total
		profile.LambdaTexture /= total
		profile.LambdaScale /= total
		profile.LambdaSmoothing /= total
	}
}

// Helper functions
func rgbaToFloat(c color.Color) [3]float64 {
	r, g, b, _ := c.RGBA()
	return [3]float64{
		float64(r) / 65535.0,
		float64(g) / 65535.0,
		float64(b) / 65535.0,
	}
}

func mean(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}

// ============================================================================
// MULTI-PHASE INHERITANCE TRACKING
// ============================================================================

type PhaseMemory struct {
	// Héritage inter-phase: garde les tensions stables de la phase précédente
	PreviousEnergies  map[string]float64 // Format: "x,y" -> energy
	StableStructures  []string           // Indices of stable atoms
	CanBreakThreshold float64            // Énergie minimale pour "casser" une structure stable
	InheritanceWeight float64            // Poids du souvenir (0.0-1.0)
}

func NewPhaseMemory() *PhaseMemory {
	return &PhaseMemory{
		PreviousEnergies:  make(map[string]float64),
		StableStructures:  make([]string, 0),
		CanBreakThreshold: 0.15, // Doit réduire 15% d'énergie pour casser
		InheritanceWeight: 0.4,  // 40% de poids au souvenir
	}
}

// ============================================================================
// THE CONSTRAINT-RELAXATION NETWORK
// ============================================================================

type ConstraintRelaxationNetwork struct {
	// Grid of atoms
	Atoms [][]PixelAtomV2

	Width, Height int
	PatchSize     int

	// Multi-level structure
	GlobalField    *GlobalCoherenceField
	PatternRegions []*PatternRegion
	EnergyProfile  *ImageEnergyProfile // Target energy signature (if analyzing)

	// Energy tracking
	TotalEnergy         float64
	PreviousTotalEnergy float64
	AverageLocalEnergy  float64
	SystemStability     float64 // -1 (oscillating) to +1 (stable)
	ConvergencePhase    int     // Track which phase we're in

	// Iteration tracking
	Iteration           int
	IterationsAtPlateau int // How many iterations at same energy level?

	// Phase inheritance
	PhaseMemory *PhaseMemory // Héritage de phase précédente

	// Threading
	mutex sync.RWMutex
}

// NewConstraintRelaxationNetwork creates a new energy-based image generator
func NewConstraintRelaxationNetwork(width, height, patchSize int) *ConstraintRelaxationNetwork {
	net := &ConstraintRelaxationNetwork{
		Atoms:               make([][]PixelAtomV2, height/patchSize),
		Width:               width,
		Height:              height,
		PatchSize:           patchSize,
		GlobalField:         NewGlobalCoherenceField(),
		PatternRegions:      make([]*PatternRegion, 0),
		EnergyProfile:       nil, // Will be set if analyzing an image
		TotalEnergy:         0,
		PreviousTotalEnergy: 0,
		AverageLocalEnergy:  0,
		SystemStability:     0,
		ConvergencePhase:    0,
		Iteration:           0,
		IterationsAtPlateau: 0,
		PhaseMemory:         NewPhaseMemory(),
	}

	// Initialize atoms
	gridH := height / patchSize
	gridW := width / patchSize

	for y := 0; y < gridH; y++ {
		net.Atoms[y] = make([]PixelAtomV2, gridW)
		for x := 0; x < gridW; x++ {
			atom := &PixelAtomV2{
				X:               x,
				Y:               y,
				R:               rand.Float64()*0.3 + 0.35, // Random color
				G:               rand.Float64()*0.3 + 0.35,
				B:               rand.Float64()*0.3 + 0.35,
				Intensity:       rand.Float64()*0.4 + 0.3, // Random brightness
				Orientation:     rand.Float64() * 2 * math.Pi,
				Confidence:      0.1, // Start low
				TextureIndex:    rand.Float64(),
				LocalEnergy:     0,
				CoherenceError:  0,
				StabilityScore:  0,
				LastStateEnergy: 0,
				EnergyTrend:     0,
				Neighbors:       [8]*PixelAtomV2{},
				IsEdge:          false,
			}
			net.Atoms[y][x] = *atom
		}
	}

	// Connect neighbors (8-connectivity)
	for y := 0; y < gridH; y++ {
		for x := 0; x < gridW; x++ {
			net.connectAtomNeighbors(x, y, gridW, gridH)
		}
	}

	return net
}

func NewGlobalCoherenceField() *GlobalCoherenceField {
	return &GlobalCoherenceField{
		AverageBrightness:    0.5,
		DominantGradientDir:  0,   // Arbitrary starting direction
		GlobalSymmetryTarget: 0.3, // Slight symmetry
		TextureConsistency:   0.4,
		ShadowDirection:      math.Pi / 4, // Top-left
		EdgeCohesion:         0.5,
		InfluenceWeight:      0.08, // Weak influence
		DirectionalField:     nil,  // Will be set if using symmetry breaking
		FieldStrength:        0.10, // Weak guidance (5-15%)
		SymmetryViolation:    0,
		TextureViolation:     0,
		LightingViolation:    0,
		ContourViolation:     0,
	}
}

// CreateDirectionalFieldFromImage extracts low-res orientation guide from image
func (field *GlobalCoherenceField) CreateDirectionalFieldFromImage(sourceImage image.Image, gridSize int) {
	bounds := sourceImage.Bounds()
	width := bounds.Max.X - bounds.Min.X
	height := bounds.Max.Y - bounds.Min.Y

	// Create low-resolution field (gridSize × gridSize)
	field.DirectionalField = make([][]float64, gridSize)
	for i := range field.DirectionalField {
		field.DirectionalField[i] = make([]float64, gridSize)
	}

	// Sample and blur to extract orientation
	for fy := 0; fy < gridSize; fy++ {
		for fx := 0; fx < gridSize; fx++ {
			// Map field coordinates to image
			imgY := (fy * height) / gridSize
			imgX := (fx * width) / gridSize

			// Compute local gradient at this point (via Sobel)
			var gradX, gradY float64

			// Sample 3×3 neighborhood
			for dy := -1; dy <= 1; dy++ {
				for dx := -1; dx <= 1; dx++ {
					y := imgY + dy
					x := imgX + dx

					if x >= 0 && x < width && y >= 0 && y < height {
						c := sourceImage.At(x, y)
						r, g, b, _ := c.RGBA()
						brightness := float64(r+g+b) / (3 * 65535)

						// Sobel coefficients
						if dy == -1 && dx == -1 {
							gradY -= brightness
							gradX -= brightness
						} else if dy == -1 && dx == 0 {
							gradY -= 2 * brightness
						} else if dy == -1 && dx == 1 {
							gradY -= brightness
							gradX += brightness
						} else if dy == 0 && dx == -1 {
							gradX -= 2 * brightness
						} else if dy == 0 && dx == 1 {
							gradX += 2 * brightness
						} else if dy == 1 && dx == -1 {
							gradY += brightness
							gradX -= brightness
						} else if dy == 1 && dx == 0 {
							gradY += 2 * brightness
						} else if dy == 1 && dx == 1 {
							gradY += brightness
							gradX += brightness
						}
					}
				}
			}

			// Compute orientation from gradient
			orientation := math.Atan2(gradY, gradX)
			field.DirectionalField[fy][fx] = orientation
		}
	}
}

func (net *ConstraintRelaxationNetwork) connectAtomNeighbors(x, y, gridW, gridH int) {
	directions := [8][2]int{
		{-1, -1}, {0, -1}, {1, -1},
		{-1, 0}, {1, 0},
		{-1, 1}, {0, 1}, {1, 1},
	}

	for i, dir := range directions {
		nx, ny := x+dir[0], y+dir[1]
		if nx >= 0 && nx < gridW && ny >= 0 && ny < gridH {
			net.Atoms[y][x].Neighbors[i] = &net.Atoms[ny][nx]
		}
	}

	// Mark edges
	if x == 0 || x == gridW-1 || y == 0 || y == gridH-1 {
		net.Atoms[y][x].IsEdge = true
	}
}

// ============================================================================
// ENERGY COMPUTATION (The core physics)
// ============================================================================

// ComputeLocalEnergy calculates the "tension" an atom experiences
// Lower energy = more stable, in equilibrium with neighbors
func (atom *PixelAtomV2) ComputeLocalEnergy(globalField *GlobalCoherenceField) float64 {
	var energy float64

	// TERM 1: Continuity with neighbors
	// Penalty if color differs too much from neighbors
	for _, neighbor := range atom.Neighbors {
		if neighbor == nil {
			continue
		}

		colorDiff := math.Sqrt(
			math.Pow(atom.R-neighbor.R, 2) +
				math.Pow(atom.G-neighbor.G, 2) +
				math.Pow(atom.B-neighbor.B, 2),
		)

		// Smooth transitions preferred, but sharp edges are OK if consistent
		edgeDetection := atom.detectEdgeSharpness(neighbor)
		if edgeDetection < 0.7 { // Not a real edge
			energy += colorDiff * (1 - edgeDetection) * 0.3
		}

		// Intensity continuity
		intensityDiff := math.Abs(atom.Intensity - neighbor.Intensity)
		energy += intensityDiff * 0.2
	}

	// 🔥 SOLUTION 1: TERME DE CONTINUITÉ DE PHASE (Phase Continuity Term)
	// Force les gradients cohérents à se propager (E_phase = Σ ||∇I(x) - ∇I(x+1)||²)
	// Cela BRISE LA SYMÉTRIE en favorisant les structures étendues plutôt que le bruit
	phaseContinuityEnergy := atom.computePhaseContinuity()
	energy += phaseContinuityEnergy * 0.15 // Poids réduit pour éviter oscillations (était 0.4)

	// TERM 2: Orientation/Gradient consistency
	// Atoms want their gradients to align with their neighbors
	for _, neighbor := range atom.Neighbors {
		if neighbor == nil {
			continue
		}

		orientationDiff := math.Abs(atom.Orientation - neighbor.Orientation)
		// Normalize to [0, π]
		if orientationDiff > math.Pi {
			orientationDiff = 2*math.Pi - orientationDiff
		}
		energy += orientationDiff * 0.15
	}

	// TERM 3: Global field influence (weak)
	// The atom doesn't strongly obey global rules, just slightly attracted
	brightnessDeviation := math.Abs((atom.R+atom.G+atom.B)/3 - globalField.AverageBrightness)
	energy += brightnessDeviation * globalField.InfluenceWeight * 0.5

	// Slight attraction to global shadow direction
	shadowAttractionError := math.Abs(atom.Orientation - globalField.ShadowDirection)
	if shadowAttractionError > math.Pi {
		shadowAttractionError = 2*math.Pi - shadowAttractionError
	}
	energy += shadowAttractionError * globalField.InfluenceWeight * 0.3

	// TERM 4: Texture consistency in local region + TERM 6: RÉGION COHERENCE ENERGY
	// Combined: Texture variance AND structure rewards

	var structureEnergy float64
	var colorVariance float64
	var gradientVariance float64
	var localTextureVariance float64
	neighborCount := 0

	for _, neighbor := range atom.Neighbors {
		if neighbor == nil {
			continue
		}
		neighborCount++

		// Variance de couleur locale
		colorDiff := math.Sqrt(
			math.Pow(atom.R-neighbor.R, 2) +
				math.Pow(atom.G-neighbor.G, 2) +
				math.Pow(atom.B-neighbor.B, 2),
		)
		colorVariance += colorDiff

		// Variance de gradient
		gradDiff := math.Abs(atom.Orientation - neighbor.Orientation)
		if gradDiff > math.Pi {
			gradDiff = 2*math.Pi - gradDiff
		}
		gradientVariance += gradDiff

		// Texture variance
		localTextureVariance += math.Abs(atom.TextureIndex - neighbor.TextureIndex)
	}

	if neighborCount > 0 {
		colorVariance /= float64(neighborCount)
		gradientVariance /= float64(neighborCount)
		localTextureVariance /= float64(neighborCount)
	}

	// Texture penalty
	energy += localTextureVariance * 0.1

	// RÉCOMPENSE: variation bien structurée (0.2-0.8)
	// PÉNALITÉ: soit trop uniforme (< 0.1), soit trop chaotique (> 0.9)
	idealVariance := 0.5 // Variance idéale
	colorVariationPenalty := math.Abs(colorVariance-idealVariance) * 0.15
	gradientVariationPenalty := math.Abs(gradientVariance-idealVariance*0.5) * 0.1

	structureEnergy = colorVariationPenalty + gradientVariationPenalty
	energy += structureEnergy

	// TERM 5: Confidence penalty (prefer stable states)
	// Low confidence = higher energy cost
	confidencePenalty := (1 - atom.Confidence) * 0.2
	energy += confidencePenalty

	// ⭐ TERM 7: AESTHETIC PENALTY (NEW)
	// Pénalité douce pour absence de structure
	// Si tout est complètement uniforme = mort visuelle
	totalUniformity := 1.0 - (colorVariance + 0.01) // Petit biais pour éviter div par zéro
	if totalUniformity > 0.95 {
		// C'est trop uniforme: appliquer pénalité esthétique
		estheticPenalty := (totalUniformity - 0.95) * 0.1
		energy += estheticPenalty
	}

	// ⭐ TERM 8: PHASE CONTINUITY (SYMMETRY BREAKING!)
	// Force gradients à se propager, pas apparaître au hasard
	// E_phase = ||∇I(x) - ∇I(x+1)||²
	// Ceci crée des contours continus, pas du bruit aléatoire
	var gradientContinuity float64
	for _, neighbor := range atom.Neighbors {
		if neighbor == nil {
			continue
		}
		// Penalty if my gradient is very different from neighbor's
		myGradMag := atom.Intensity // Proxy for gradient magnitude
		neighborGradMag := neighbor.Intensity
		gradContinuityError := math.Abs(myGradMag - neighborGradMag)
		gradientContinuity += gradContinuityError
	}
	if neighborCount > 0 {
		gradientContinuity /= float64(neighborCount)
	}
	energy += gradientContinuity * 0.12 // Moderate weight

	// ⭐ TERM 9: DIRECTIONAL FIELD GUIDANCE (SYMMETRY BREAKING!)
	// Si un champ directeur existe (low-res guide), l'utiliser faiblement
	if len(globalField.DirectionalField) > 0 {
		// Map atom position to low-res field
		fieldH := len(globalField.DirectionalField)
		fieldW := len(globalField.DirectionalField[0])

		// Scale atom position to field coordinates
		fieldY := (atom.Y * fieldH) / 1000 // Normalized position
		fieldX := (atom.X * fieldW) / 1000

		if fieldY >= 0 && fieldY < fieldH && fieldX >= 0 && fieldX < fieldW {
			targetOrientation := globalField.DirectionalField[fieldY][fieldX]
			orientationError := math.Abs(atom.Orientation - targetOrientation)
			if orientationError > math.Pi {
				orientationError = 2*math.Pi - orientationError
			}
			energy += orientationError * globalField.FieldStrength * 0.08
		}
	}

	return energy
}

func (atom *PixelAtomV2) detectEdgeSharpness(neighbor *PixelAtomV2) float64 {
	// Returns 0-1, where 1 = sharp edge, 0 = smooth
	colorDiff := math.Sqrt(
		math.Pow(atom.R-neighbor.R, 2) +
			math.Pow(atom.G-neighbor.G, 2) +
			math.Pow(atom.B-neighbor.B, 2),
	)

	// Sharp edges have high color differences
	// Use sigmoid to get 0-1 range
	return 2 / (1 + math.Exp(-10*(colorDiff-0.3)))
}

// 🔥 SOLUTION 1: CONTINUITÉ DE PHASE - BRISURE DE SYMÉTRIE
// Cette méthode pénalise les gradients incohérents entre voisins
// E_phase = Σ ||∇I(x) - ∇I(x+1)||²
// Résultat : Les bords SE PROPAGENT au lieu d'apparaître au hasard
func (atom *PixelAtomV2) computePhaseContinuity() float64 {
	var phaseContinuityError float64
	validNeighbors := 0

	// Calculer le gradient local de cet atome
	myGradient := atom.computeLocalGradient()

	// Comparer avec le gradient de chaque voisin
	for _, neighbor := range atom.Neighbors {
		if neighbor == nil {
			continue
		}

		neighborGradient := neighbor.computeLocalGradient()

		// Distance entre les deux gradients (en norme L2)
		gradDiffX := myGradient[0] - neighborGradient[0]
		gradDiffY := myGradient[1] - neighborGradient[1]
		gradientDiscontinuity := math.Sqrt(gradDiffX*gradDiffX + gradDiffY*gradDiffY)

		phaseContinuityError += gradientDiscontinuity
		validNeighbors++
	}

	if validNeighbors > 0 {
		phaseContinuityError /= float64(validNeighbors)
	}

	return phaseContinuityError
}

// Calculer le gradient local (approximation finite difference)
// Retourne [gradX, gradY] basé sur les voisins immédiats
func (atom *PixelAtomV2) computeLocalGradient() [2]float64 {
	// Gradient = différence de couleur avec voisins horizontaux/verticaux
	var gradX, gradY float64
	var countX, countY int

	// Horizontal gradient (left vs right)
	if atom.Neighbors[3] != nil { // Left
		leftIntensity := (atom.Neighbors[3].R + atom.Neighbors[3].G + atom.Neighbors[3].B) / 3
		myIntensity := (atom.R + atom.G + atom.B) / 3
		gradX -= (myIntensity - leftIntensity)
		countX++
	}
	if atom.Neighbors[4] != nil { // Right
		rightIntensity := (atom.Neighbors[4].R + atom.Neighbors[4].G + atom.Neighbors[4].B) / 3
		myIntensity := (atom.R + atom.G + atom.B) / 3
		gradX += (rightIntensity - myIntensity)
		countX++
	}

	// Vertical gradient (top vs bottom)
	if atom.Neighbors[1] != nil { // Top
		topIntensity := (atom.Neighbors[1].R + atom.Neighbors[1].G + atom.Neighbors[1].B) / 3
		myIntensity := (atom.R + atom.G + atom.B) / 3
		gradY -= (myIntensity - topIntensity)
		countY++
	}
	if atom.Neighbors[6] != nil { // Bottom
		bottomIntensity := (atom.Neighbors[6].R + atom.Neighbors[6].G + atom.Neighbors[6].B) / 3
		myIntensity := (atom.R + atom.G + atom.B) / 3
		gradY += (bottomIntensity - myIntensity)
		countY++
	}

	if countX > 0 {
		gradX /= float64(countX)
	}
	if countY > 0 {
		gradY /= float64(countY)
	}

	return [2]float64{gradX, gradY}
}

// ============================================================================
// RELAXATION STEP (Atoms minimize energy)
// ============================================================================

// RelaxationStep performs one iteration of energy minimization
// Each atom adjusts to reduce its local energy
func (net *ConstraintRelaxationNetwork) RelaxationStep() {
	net.mutex.Lock()
	net.Iteration++
	net.mutex.Unlock()

	gridH := net.Height / net.PatchSize
	gridW := net.Width / net.PatchSize

	// Phase 1: Compute all local energies
	var totalEnergy float64
	energies := make([][]float64, gridH)

	for y := 0; y < gridH; y++ {
		energies[y] = make([]float64, gridW)
		for x := 0; x < gridW; x++ {
			energy := net.Atoms[y][x].ComputeLocalEnergy(net.GlobalField)
			energies[y][x] = energy
			totalEnergy += energy
		}
	}

	// Phase 2: Update atoms to reduce energy
	for y := 0; y < gridH; y++ {
		for x := 0; x < gridW; x++ {
			net.minimizeAtomEnergy(&net.Atoms[y][x], net.GlobalField)
		}
	}

	// Phase 3: Track convergence
	net.PreviousTotalEnergy = net.TotalEnergy
	net.TotalEnergy = totalEnergy
	net.AverageLocalEnergy = totalEnergy / float64(gridW*gridH)

	// Detect if we're at a plateau (convergence)
	energyChange := math.Abs(net.TotalEnergy - net.PreviousTotalEnergy)
	if energyChange < 0.001 {
		net.IterationsAtPlateau++
		net.SystemStability = math.Min(1.0, net.SystemStability+0.02)
	} else {
		net.IterationsAtPlateau = 0
		net.SystemStability = math.Max(-1.0, net.SystemStability-0.01)
	}

	// Phase 4: Auto-réévaluation - compare this state with previous
	net.evaluateAndStabilize()
}

// minimizeAtomEnergy adjusts an atom's state to reduce energy
func (net *ConstraintRelaxationNetwork) minimizeAtomEnergy(atom *PixelAtomV2, globalField *GlobalCoherenceField) {
	atom.mutex.Lock()
	defer atom.mutex.Unlock()

	// Save old state for comparison
	oldR, oldG, oldB := atom.R, atom.G, atom.B
	oldEnergy := atom.ComputeLocalEnergy(globalField)

	// Adaptive step size with damping based on confidence AND iteration
	// High confidence = small steps (stable), Low confidence = larger steps (exploring)
	// Later iterations = smaller steps for convergence
	baseStepSize := 0.03                     // Réduit de 0.05 à 0.03 pour meilleure stabilité
	damping := 0.4 + (atom.Confidence * 0.4) // damping ∈ [0.4, 0.8]
	stepSize := baseStepSize * damping

	// Direction 1: Move color toward neighbor average (if not at edge)
	if !atom.IsEdge {
		var avgR, avgG, avgB float64
		neighborCount := 0
		for _, neighbor := range atom.Neighbors {
			if neighbor != nil {
				avgR += neighbor.R
				avgG += neighbor.G
				avgB += neighbor.B
				neighborCount++
			}
		}

		if neighborCount > 0 {
			avgR /= float64(neighborCount)
			avgG /= float64(neighborCount)
			avgB /= float64(neighborCount)

			// Move partially toward average (smoothing)
			atom.R = atom.R*(1-stepSize*0.3) + avgR*stepSize*0.3
			atom.G = atom.G*(1-stepSize*0.3) + avgG*stepSize*0.3
			atom.B = atom.B*(1-stepSize*0.3) + avgB*stepSize*0.3
		}
	}

	// Direction 2: Move intensity toward global average
	targetIntensity := globalField.AverageBrightness
	atom.Intensity = atom.Intensity*(1-stepSize*0.1) + targetIntensity*stepSize*0.1

	// Direction 3: Align orientation with global shadow direction
	orientationDiff := globalField.ShadowDirection - atom.Orientation
	// Normalize to [-π, π]
	for orientationDiff > math.Pi {
		orientationDiff -= 2 * math.Pi
	}
	for orientationDiff < -math.Pi {
		orientationDiff += 2 * math.Pi
	}
	atom.Orientation += orientationDiff * stepSize * 0.05

	// Clamp values
	atom.clampState()

	// Check if energy improved
	newEnergy := atom.ComputeLocalEnergy(globalField)
	if newEnergy > oldEnergy {
		// Energy got worse, revert
		atom.R = oldR
		atom.G = oldG
		atom.B = oldB
		atom.StabilityScore -= 0.1
	} else {
		// Energy improved
		atom.StabilityScore += 0.05
		atom.LastStateEnergy = oldEnergy
		atom.EnergyTrend = oldEnergy - newEnergy // Positive = improving
	}

	// Update confidence based on stability
	if atom.StabilityScore > 0.5 {
		atom.Confidence = math.Min(1.0, atom.Confidence+0.02)
	} else {
		atom.Confidence = math.Max(0.1, atom.Confidence-0.01)
	}
}

func (atom *PixelAtomV2) clampState() {
	if atom.R > 1.0 {
		atom.R = 1.0
	} else if atom.R < 0 {
		atom.R = 0
	}

	if atom.G > 1.0 {
		atom.G = 1.0
	} else if atom.G < 0 {
		atom.G = 0
	}

	if atom.B > 1.0 {
		atom.B = 1.0
	} else if atom.B < 0 {
		atom.B = 0
	}

	if atom.Intensity > 1.0 {
		atom.Intensity = 1.0
	} else if atom.Intensity < 0 {
		atom.Intensity = 0
	}
}

// ============================================================================
// AUTO-RÉÉVALUATION (The key to quality improvement)
// ============================================================================

// evaluateAndStabilize implements self-assessment
// The network compares state n vs n-1 and penalizes oscillations
func (net *ConstraintRelaxationNetwork) evaluateAndStabilize() {
	gridH := net.Height / net.PatchSize
	gridW := net.Width / net.PatchSize

	// Detect oscillations - atoms bouncing back and forth
	oscillatingAtoms := 0
	for y := 0; y < gridH; y++ {
		for x := 0; x < gridW; x++ {
			atom := &net.Atoms[y][x]
			// If trend is small or changing sign, we're oscillating
			if math.Abs(atom.EnergyTrend) < 0.005 && net.Iteration > 10 {
				oscillatingAtoms++
				// Suppress oscillation by increasing "damping"
				atom.Confidence = math.Max(0.7, atom.Confidence+0.05)
			}
		}
	}

	// If too many atoms are oscillating, reduce global field influence
	oscillationRatio := float64(oscillatingAtoms) / float64(gridW*gridH)
	if oscillationRatio > 0.3 {
		net.GlobalField.InfluenceWeight *= 0.95 // Reduce influence
	} else if oscillationRatio < 0.1 {
		net.GlobalField.InfluenceWeight = math.Min(0.2, net.GlobalField.InfluenceWeight*1.02)
	}

	// Update global statistics
	var totalBrightness float64
	for y := 0; y < gridH; y++ {
		for x := 0; x < gridW; x++ {
			totalBrightness += (net.Atoms[y][x].R + net.Atoms[y][x].G + net.Atoms[y][x].B) / 3
		}
	}
	net.GlobalField.AverageBrightness = totalBrightness / float64(gridW*gridH)
}

// ============================================================================
// PATTERN DETECTION (Level 2)
// ============================================================================

// DetectPatterns identifies coherent regions in the current state
func (net *ConstraintRelaxationNetwork) DetectPatterns() []*PatternRegion {
	gridH := net.Height / net.PatchSize
	gridW := net.Width / net.PatchSize

	regions := make([]*PatternRegion, 0)

	// Simple region growing: find connected components with similar properties
	visited := make([][]bool, gridH)
	for y := 0; y < gridH; y++ {
		visited[y] = make([]bool, gridW)
	}

	for y := 0; y < gridH; y++ {
		for x := 0; x < gridW; x++ {
			if !visited[y][x] {
				region := net.growRegion(x, y, visited)
				if len(region.Atoms) > 4 { // Only keep significant regions
					regions = append(regions, region)
				}
			}
		}
	}

	net.PatternRegions = regions
	return regions
}

func (net *ConstraintRelaxationNetwork) growRegion(startX, startY int, visited [][]bool) *PatternRegion {
	region := &PatternRegion{
		X1:    startX,
		Y1:    startY,
		X2:    startX,
		Y2:    startY,
		Atoms: make([]*PixelAtomV2, 0),
	}

	gridH := len(visited)
	gridW := len(visited[0])

	// BFS to grow the region
	queue := [][2]int{{startX, startY}}
	visited[startY][startX] = true

	startColor := [3]float64{
		net.Atoms[startY][startX].R,
		net.Atoms[startY][startX].G,
		net.Atoms[startY][startX].B,
	}

	for len(queue) > 0 {
		pos := queue[0]
		queue = queue[1:]
		x, y := pos[0], pos[1]

		atom := &net.Atoms[y][x]
		region.Atoms = append(region.Atoms, atom)

		// Update bounds
		if x < region.X1 {
			region.X1 = x
		}
		if x > region.X2 {
			region.X2 = x
		}
		if y < region.Y1 {
			region.Y1 = y
		}
		if y > region.Y2 {
			region.Y2 = y
		}

		// Check neighbors
		for dy := -1; dy <= 1; dy++ {
			for dx := -1; dx <= 1; dx++ {
				if dx == 0 && dy == 0 {
					continue
				}
				nx, ny := x+dx, y+dy
				if nx >= 0 && nx < gridW && ny >= 0 && ny < gridH && !visited[ny][nx] {
					neighbor := &net.Atoms[ny][nx]
					colorDist := math.Sqrt(
						math.Pow(neighbor.R-startColor[0], 2) +
							math.Pow(neighbor.G-startColor[1], 2) +
							math.Pow(neighbor.B-startColor[2], 2),
					)
					if colorDist < 0.2 { // Similar color threshold
						visited[ny][nx] = true
						queue = append(queue, [2]int{nx, ny})
					}
				}
			}
		}
	}

	// Calculate region properties
	region.computeProperties()
	return region
}

func (region *PatternRegion) computeProperties() {
	if len(region.Atoms) == 0 {
		return
	}

	// Average color
	var sumR, sumG, sumB float64
	var sumOrient float64
	var maxEdge float64

	for _, atom := range region.Atoms {
		sumR += atom.R
		sumG += atom.G
		sumB += atom.B
		sumOrient += atom.Orientation
		for _, neighbor := range atom.Neighbors {
			if neighbor != nil {
				edge := atom.detectEdgeSharpness(neighbor)
				if edge > maxEdge {
					maxEdge = edge
				}
			}
		}
	}

	n := float64(len(region.Atoms))
	region.DominantColor = [3]float64{sumR / n, sumG / n, sumB / n}
	region.GradientDirection = sumOrient / n

	region.EdgeStrength = maxEdge

	// Determine texture type
	if region.EdgeStrength > 0.7 {
		region.TextureType = "detailed"
	} else if region.EdgeStrength < 0.3 {
		region.TextureType = "smooth"
	} else {
		region.TextureType = "rough"
	}

	// ⭐ COMPUTE COHERENCE (NEW)
	// Mesure la cohérence interne: comment les couleurs/gradients s'alignent-ils?
	// Basse variance = cohérence haute
	if len(region.Atoms) > 1 {
		var colorVariance float64
		var orientationVariance float64

		avgColor := region.DominantColor

		for _, atom := range region.Atoms {
			// Variance de couleur par rapport à la moyenne
			colorDiff := math.Sqrt(
				math.Pow(atom.R-avgColor[0], 2) +
					math.Pow(atom.G-avgColor[1], 2) +
					math.Pow(atom.B-avgColor[2], 2),
			)
			colorVariance += colorDiff

			// Variance d'orientation
			orientDiff := math.Abs(atom.Orientation - region.GradientDirection)
			if orientDiff > math.Pi {
				orientDiff = 2*math.Pi - orientDiff
			}
			orientationVariance += orientDiff
		}

		colorVariance /= n
		orientationVariance /= n

		// Cohérence = 1 - variance normalisée (0-1)
		// Plus la variance est basse, plus la cohérence est haute
		maxColorVar := 1.0  // Variance de couleur maximale attendue
		maxOrientVar := 1.0 // Variance d'orientation maximale

		colorCoherence := 1.0 - math.Min(1.0, colorVariance/maxColorVar)
		orientCoherence := 1.0 - math.Min(1.0, orientationVariance/maxOrientVar)

		// Cohérence globale = moyenne pondérée
		region.Coherence = (colorCoherence*0.7 + orientCoherence*0.3)

		// Clamp to [0, 1]
		region.Coherence = math.Max(0, math.Min(1, region.Coherence))
	} else {
		region.Coherence = 0.0
	}
}

// ============================================================================
// GENERATION & RENDERING
// ============================================================================

// Generate performs constraint relaxation until convergence
func (net *ConstraintRelaxationNetwork) Generate(maxIterations int) {
	for i := 0; i < maxIterations; i++ {
		net.RelaxationStep()

		// Print progress
		if (i+1)%50 == 0 {
			fmt.Printf("[%4d/%4d] Energy: %.4f | Stability: %.2f | Plateau: %d\n",
				i+1, maxIterations,
				net.AverageLocalEnergy,
				net.SystemStability,
				net.IterationsAtPlateau,
			)
		}

		// Early stopping if converged
		if net.IterationsAtPlateau > 200 && net.SystemStability > 0.8 {
			fmt.Printf("✓ Converged at iteration %d\n", i+1)
			break
		}
	}

	// Final pattern detection
	net.DetectPatterns()
}

// RenderToImage converts the relaxed state to an image
func (net *ConstraintRelaxationNetwork) RenderToImage() image.Image {
	img := image.NewRGBA(image.Rect(0, 0, net.Width, net.Height))

	gridH := net.Height / net.PatchSize
	gridW := net.Width / net.PatchSize

	for y := 0; y < gridH; y++ {
		for x := 0; x < gridW; x++ {
			atom := &net.Atoms[y][x]

			// Render atom with confidence-based contrast
			r := uint8(atom.R * 255 * atom.Confidence)
			g := uint8(atom.G * 255 * atom.Confidence)
			b := uint8(atom.B * 255 * atom.Confidence)

			c := color.RGBA{r, g, b, 255}

			// Fill patch
			for py := 0; py < net.PatchSize; py++ {
				for px := 0; px < net.PatchSize; px++ {
					imgX := x*net.PatchSize + px
					imgY := y*net.PatchSize + py
					if imgX < net.Width && imgY < net.Height {
						img.Set(imgX, imgY, c)
					}
				}
			}
		}
	}

	return img
}

// ============================================================================
// MULTI-PHASE INHERITANCE
// ============================================================================

// CapturePhaseState sauvegarde l'état stable de la phase actuelle
// pour l'héritage dans la phase suivante
func (net *ConstraintRelaxationNetwork) CapturePhaseState() {
	if net.PhaseMemory == nil {
		net.PhaseMemory = NewPhaseMemory()
	}

	// Sauvegarder les énergies de tous les atomes
	for y := range net.Atoms {
		for x := range net.Atoms[y] {
			key := fmt.Sprintf("%d,%d", x, y)
			net.PhaseMemory.PreviousEnergies[key] = net.Atoms[y][x].LocalEnergy

			// Marquer comme stable si confiance > 0.6
			if net.Atoms[y][x].Confidence > 0.6 {
				net.PhaseMemory.StableStructures = append(net.PhaseMemory.StableStructures, key)
			}
		}
	}

	fmt.Printf("📸 Capture de phase: %d structures stables identifiées\n",
		len(net.PhaseMemory.StableStructures))
}

// InheritPhaseState transfère les tensions stables de la phase précédente
// Les atomes héritent des positions/énergies, mais peuvent s'améliorer
func (net *ConstraintRelaxationNetwork) InheritPhaseState(previousNetwork *ConstraintRelaxationNetwork) {
	if previousNetwork.PhaseMemory == nil || len(previousNetwork.PhaseMemory.PreviousEnergies) == 0 {
		fmt.Println("⚠️  Pas d'état précédent à hériter")
		return
	}

	if net.PhaseMemory == nil {
		net.PhaseMemory = NewPhaseMemory()
	}

	// Copier les énergies précédentes
	net.PhaseMemory.PreviousEnergies = make(map[string]float64)
	for k, v := range previousNetwork.PhaseMemory.PreviousEnergies {
		net.PhaseMemory.PreviousEnergies[k] = v
	}

	// Copier les structures stables
	net.PhaseMemory.StableStructures = make([]string, len(previousNetwork.PhaseMemory.StableStructures))
	copy(net.PhaseMemory.StableStructures, previousNetwork.PhaseMemory.StableStructures)

	// Propager les couleurs/états des atomes (interpolation)
	minH := len(previousNetwork.Atoms)
	if len(net.Atoms) < minH {
		minH = len(net.Atoms)
	}
	minW := 0
	if len(previousNetwork.Atoms) > 0 {
		minW = len(previousNetwork.Atoms[0])
		if len(net.Atoms) > 0 && len(net.Atoms[0]) < minW {
			minW = len(net.Atoms[0])
		}
	}

	for y := 0; y < minH; y++ {
		for x := 0; x < minW; x++ {
			prev := &previousNetwork.Atoms[y][x]
			curr := &net.Atoms[y][x]

			// Hériter les états (avec petit décalage pour variation)
			curr.R = prev.R + (rand.Float64()-0.5)*0.05
			curr.G = prev.G + (rand.Float64()-0.5)*0.05
			curr.B = prev.B + (rand.Float64()-0.5)*0.05
			curr.Intensity = prev.Intensity + (rand.Float64()-0.5)*0.05
			curr.Orientation = prev.Orientation
			curr.Confidence = prev.Confidence * 0.8 // Légèrement réinitialiser
			curr.TextureIndex = prev.TextureIndex

			// Clamp à [0, 1]
			curr.R = math.Max(0, math.Min(1, curr.R))
			curr.G = math.Max(0, math.Min(1, curr.G))
			curr.B = math.Max(0, math.Min(1, curr.B))
			curr.Intensity = math.Max(0, math.Min(1, curr.Intensity))
		}
	}

	fmt.Printf("🔗 État de phase hérité: %d atomes, poids d'héritage=%.1f%%\n",
		minH*minW, net.PhaseMemory.InheritanceWeight*100)
}
