package database

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"sync"
)

// AtomicCell represents a single atom in the generation grid
// Each atom has: color, internal state, and connections to neighbors
type AtomicCell struct {
	// Color representation
	Color [3]float64 // RGB [0,1]

	// Internal state vector
	// s_ij contains: [intensity, orientation, phase, frequency, ...]
	State [5]float64 // intensity, orientation, phase, frequency, coherence

	// Connection weights to neighbors
	Weights [8]float64 // w_ij;kl for 8 neighbors

	// Pattern information
	Pattern [3]float64 // Pattern input (from waves/diffusion)

	// Velocity for momentum
	Velocity [5]float64
}

// GenerationGrid represents the complete atomic grid
type GenerationGrid struct {
	Width, Height int
	Cells         [][]AtomicCell

	// Parameters
	ResonanceAlpha  float64 // α: neighbor influence (0.1-0.5)
	PatternBeta     float64 // β: pattern importance (0.3-0.7)
	SmoothingGamma  float64 // γ: color smoothing (0.1-0.3)
	FeedbackWeight  float64 // δ: feedback strength (0.0-1.0)
	VelocityDamping float64 // ε: momentum damping

	// Target pattern for feedback
	TargetImage image.Image
	TargetData  [][]PatternPixel // Pre-extracted target patterns

	// Iteration tracking
	IterationCount int
	AverageLoss    float64
	Convergence    float64

	mu sync.RWMutex
}

// NewGenerationGrid creates a new atomic generation grid
func NewGenerationGrid(width, height int) *GenerationGrid {
	gg := &GenerationGrid{
		Width:           width,
		Height:          height,
		Cells:           make([][]AtomicCell, height),
		ResonanceAlpha:  0.3, // α: moderate neighbor influence
		PatternBeta:     0.5, // β: balanced pattern weight
		SmoothingGamma:  0.2, // γ: subtle smoothing
		FeedbackWeight:  0.2, // δ: gentle feedback
		VelocityDamping: 0.9, // ε: smooth momentum
		IterationCount:  0,
		AverageLoss:     0.0,
		Convergence:     0.0,
	}

	// Initialize grid
	for i := 0; i < height; i++ {
		gg.Cells[i] = make([]AtomicCell, width)
		for j := 0; j < width; j++ {
			// Initialize with neutral gray
			gg.Cells[i][j] = AtomicCell{
				Color:    [3]float64{0.5, 0.5, 0.5},
				State:    [5]float64{0.5, 0.0, 0.0, 1.0, 0.5},
				Weights:  [8]float64{0.125, 0.125, 0.125, 0.125, 0.125, 0.125, 0.125, 0.125},
				Pattern:  [3]float64{0.5, 0.5, 0.5},
				Velocity: [5]float64{0, 0, 0, 0, 0},
			}
		}
	}

	return gg
}

// SetPattern injects a pattern (from wave diffusion) into the grid
// This guides the generation process
func (gg *GenerationGrid) SetPattern(pattern *PatternEmergenceEngine) {
	gg.mu.Lock()
	defer gg.mu.Unlock()

	for i := 0; i < gg.Height && i < pattern.Height; i++ {
		for j := 0; j < gg.Width && j < pattern.Width; j++ {
			pixelColor := pattern.GetPixel(j, i)
			gg.Cells[i][j].Pattern = pixelColor

			// Initialize color from pattern
			gg.Cells[i][j].Color = pixelColor

			// Extract state information from pattern
			// Intensity: average of RGB
			intensity := (pixelColor[0] + pixelColor[1] + pixelColor[2]) / 3.0
			gg.Cells[i][j].State[0] = intensity

			// Orientation: based on R-G difference
			orientation := (pixelColor[0] - pixelColor[1]) / 2.0
			gg.Cells[i][j].State[1] = orientation

			// Phase: based on G-B difference
			phase := (pixelColor[1] - pixelColor[2]) / 2.0
			gg.Cells[i][j].State[2] = phase

			// Frequency: fixed for now
			gg.Cells[i][j].State[3] = 1.0

			// Coherence: color saturation
			maxC := math.Max(pixelColor[0], math.Max(pixelColor[1], pixelColor[2]))
			minC := math.Min(pixelColor[0], math.Min(pixelColor[1], pixelColor[2]))
			if maxC > 0 {
				coherence := (maxC - minC) / maxC
				gg.Cells[i][j].State[4] = coherence
			}
		}
	}
}

// PropagateLocal performs one iteration of local state propagation
// s_ij(t+1) = α·(1/|N|)·Σ s_kl(t) + β·P_ij
func (gg *GenerationGrid) PropagateLocal() {
	gg.mu.Lock()
	defer gg.mu.Unlock()

	newCells := make([][]AtomicCell, gg.Height)
	totalLoss := 0.0

	for i := 0; i < gg.Height; i++ {
		newCells[i] = make([]AtomicCell, gg.Width)
		copy(newCells[i], gg.Cells[i])
	}

	// Propagate state for each cell
	for i := 0; i < gg.Height; i++ {
		for j := 0; j < gg.Width; j++ {
			// Compute average neighbor state
			neighborSum := [5]float64{0, 0, 0, 0, 0}
			neighborCount := 0

			for di := -1; di <= 1; di++ {
				for dj := -1; dj <= 1; dj++ {
					if di == 0 && dj == 0 {
						continue
					}

					ni, nj := i+di, j+dj
					if ni >= 0 && ni < gg.Height && nj >= 0 && nj < gg.Width {
						neighbor := gg.Cells[ni][nj]
						for s := 0; s < 5; s++ {
							neighborSum[s] += neighbor.State[s]
						}
						neighborCount++
					}
				}
			}

			// s_ij(t+1) = α·(neighbor_avg) + β·P_ij + ε·V
			for s := 0; s < 5; s++ {
				neighborAvg := neighborSum[s] / float64(neighborCount)

				// Combine neighbor influence + pattern + momentum
				newCells[i][j].Velocity[s] *= gg.VelocityDamping
				// Use only first 3 state components for pattern (color)
				patternIdx := s
				if s > 2 {
					patternIdx = 2
				}
				newCells[i][j].State[s] = gg.ResonanceAlpha*neighborAvg +
					gg.PatternBeta*gg.Cells[i][j].Pattern[patternIdx] +
					newCells[i][j].Velocity[s]

				// Store velocity for next iteration
				newCells[i][j].Velocity[s] = gg.ResonanceAlpha * (neighborAvg - gg.Cells[i][j].State[s])

				// Clamp state to [0, 1]
				newCells[i][j].State[s] = math.Max(0, math.Min(1, newCells[i][j].State[s]))
			}

			// Track loss (divergence from pattern)
			for c := 0; c < 3; c++ {
				diff := newCells[i][j].State[c] - gg.Cells[i][j].Pattern[c]
				totalLoss += diff * diff
			}
		}
	}

	gg.Cells = newCells
	gg.IterationCount++
	gg.AverageLoss = totalLoss / float64(gg.Width*gg.Height)
}

// StateToColor converts internal state to RGB color
// f(s_ij) = RGB representation of the state
func (gg *GenerationGrid) StateToColor(i, j int) [3]float64 {
	cell := gg.Cells[i][j]

	// Extract state components
	intensity := cell.State[0]   // Average brightness
	orientation := cell.State[1] // R-G axis
	_ = cell.State[2]            // G-B axis (phase)
	_ = cell.State[3]            // Frequency component
	coherence := cell.State[4]   // Color saturation

	// Convert to HSL-like representation
	hue := orientation // 0-1 maps to hue-like rotation
	saturation := coherence
	lightness := intensity

	// HSL to RGB conversion
	c := (1 - math.Abs(2*lightness-1)) * saturation
	hPrime := hue * 6.0
	x := c * (1 - math.Abs(math.Mod(hPrime, 2)-1))

	var r1, g1, b1 float64
	if hPrime >= 0 && hPrime < 1 {
		r1, g1, b1 = c, x, 0
	} else if hPrime >= 1 && hPrime < 2 {
		r1, g1, b1 = x, c, 0
	} else if hPrime >= 2 && hPrime < 3 {
		r1, g1, b1 = 0, c, x
	} else if hPrime >= 3 && hPrime < 4 {
		r1, g1, b1 = 0, x, c
	} else if hPrime >= 4 && hPrime < 5 {
		r1, g1, b1 = x, 0, c
	} else {
		r1, g1, b1 = c, 0, x
	}

	m := lightness - c/2
	return [3]float64{
		r1 + m,
		g1 + m,
		b1 + m,
	}
}

// SmoothColors applies local smoothing to avoid "cube" artifacts
// C_ij' = (1/|N|)·Σ C_kl
func (gg *GenerationGrid) SmoothColors() {
	gg.mu.Lock()
	defer gg.mu.Unlock()

	smoothed := make([][]AtomicCell, gg.Height)
	for i := 0; i < gg.Height; i++ {
		smoothed[i] = make([]AtomicCell, gg.Width)
		copy(smoothed[i], gg.Cells[i])
	}

	for i := 0; i < gg.Height; i++ {
		for j := 0; j < gg.Width; j++ {
			colorSum := [3]float64{0, 0, 0}
			neighborCount := 0

			for di := -1; di <= 1; di++ {
				for dj := -1; dj <= 1; dj++ {
					ni, nj := i+di, j+dj
					if ni >= 0 && ni < gg.Height && nj >= 0 && nj < gg.Width {
						neighbor := gg.Cells[ni][nj]
						for c := 0; c < 3; c++ {
							colorSum[c] += neighbor.Color[c]
						}
						neighborCount++
					}
				}
			}

			// Smooth: blend own color with neighbor average
			for c := 0; c < 3; c++ {
				neighborAvg := colorSum[c] / float64(neighborCount)
				smoothed[i][j].Color[c] = (1-gg.SmoothingGamma)*gg.Cells[i][j].Color[c] +
					gg.SmoothingGamma*neighborAvg
			}
		}
	}

	gg.Cells = smoothed
}

// GenerateColors converts all internal states to final colors
// This runs after state propagation for each iteration
func (gg *GenerationGrid) GenerateColors() {
	gg.mu.Lock()
	defer gg.mu.Unlock()

	for i := 0; i < gg.Height; i++ {
		for j := 0; j < gg.Width; j++ {
			color := gg.StateToColor(i, j)
			gg.Cells[i][j].Color = color
		}
	}
}

// ApplyFeedback adjusts state based on target image similarity
// Provides guidance toward target pattern
func (gg *GenerationGrid) ApplyFeedback() error {
	if gg.TargetImage == nil {
		return nil // No feedback if no target
	}

	gg.mu.Lock()
	defer gg.mu.Unlock()

	bounds := gg.TargetImage.Bounds()
	scaleX := float64(gg.Width) / float64(bounds.Dx())
	scaleY := float64(gg.Height) / float64(bounds.Dy())

	for i := 0; i < gg.Height; i++ {
		for j := 0; j < gg.Width; j++ {
			// Map grid coordinates to target image
			srcX := int(float64(j) / scaleX)
			srcY := int(float64(i) / scaleY)

			if srcX >= bounds.Min.X && srcX < bounds.Max.X &&
				srcY >= bounds.Min.Y && srcY < bounds.Max.Y {

				r, g, b, _ := gg.TargetImage.At(srcX, srcY).RGBA()
				targetColor := [3]float64{
					float64(r) / 65535.0,
					float64(g) / 65535.0,
					float64(b) / 65535.0,
				}

				// Compute error signal
				for c := 0; c < 3; c++ {
					error := targetColor[c] - gg.Cells[i][j].Color[c]
					// Apply feedback to internal state
					gg.Cells[i][j].State[c%5] += gg.FeedbackWeight * error
					gg.Cells[i][j].State[c%5] = math.Max(0, math.Min(1, gg.Cells[i][j].State[c%5]))
				}
			}
		}
	}

	return nil
}

// GenerateStep performs one complete generation iteration:
// 1. Propagate local state
// 2. Generate colors from state
// 3. Smooth colors locally
// 4. Apply feedback (if target available)
func (gg *GenerationGrid) GenerateStep() {
	gg.PropagateLocal()
	gg.GenerateColors()
	gg.SmoothColors()
	gg.ApplyFeedback()
}

// Generate runs full generation pipeline for N iterations
func (gg *GenerationGrid) Generate(iterations int) {
	step := int(math.Max(1, float64(iterations/10)))
	for i := 0; i < iterations; i++ {
		gg.GenerateStep()

		if i%step == 0 {
			fmt.Printf("  ✓ Generation step %d: Loss %.8f | Convergence: %.4f\n",
				i, gg.AverageLoss, gg.Convergence)
		}
	}
}

// SaveImage exports the current grid as PNG
func (gg *GenerationGrid) SaveImage(filename string) error {
	gg.mu.RLock()
	defer gg.mu.RUnlock()

	img := image.NewRGBA(image.Rect(0, 0, gg.Width, gg.Height))

	for i := 0; i < gg.Height; i++ {
		for j := 0; j < gg.Width; j++ {
			r := uint8(gg.Cells[i][j].Color[0] * 255)
			g := uint8(gg.Cells[i][j].Color[1] * 255)
			b := uint8(gg.Cells[i][j].Color[2] * 255)
			img.Set(j, i, color.RGBA{r, g, b, 255})
		}
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	return png.Encode(file, img)
}

// GetStatistics returns generation metrics
func (gg *GenerationGrid) GetStatistics() map[string]interface{} {
	gg.mu.RLock()
	defer gg.mu.RUnlock()

	return map[string]interface{}{
		"iterations":       gg.IterationCount,
		"average_loss":     gg.AverageLoss,
		"convergence":      gg.Convergence,
		"dimensions":       fmt.Sprintf("%dx%d", gg.Width, gg.Height),
		"resonance_alpha":  gg.ResonanceAlpha,
		"pattern_beta":     gg.PatternBeta,
		"smoothing_gamma":  gg.SmoothingGamma,
		"feedback_weight":  gg.FeedbackWeight,
		"velocity_damping": gg.VelocityDamping,
	}
}

// PrintStatistics outputs generation metrics
func (gg *GenerationGrid) PrintStatistics() {
	stats := gg.GetStatistics()
	fmt.Printf("\n📊 Atomic Generation Statistics\n")
	fmt.Printf("═══════════════════════════════════════\n")
	fmt.Printf("Iterations:       %v\n", stats["iterations"])
	fmt.Printf("Average Loss:     %.8f\n", stats["average_loss"])
	fmt.Printf("Convergence:      %.6f\n", stats["convergence"])
	fmt.Printf("Dimensions:       %v\n", stats["dimensions"])
	fmt.Printf("Resonance α:      %.4f\n", stats["resonance_alpha"])
	fmt.Printf("Pattern β:        %.4f\n", stats["pattern_beta"])
	fmt.Printf("Smoothing γ:      %.4f\n", stats["smoothing_gamma"])
	fmt.Printf("Feedback δ:       %.4f\n", stats["feedback_weight"])
	fmt.Printf("Velocity Damp:    %.4f\n", stats["velocity_damping"])
	fmt.Printf("═══════════════════════════════════════\n\n")
}
