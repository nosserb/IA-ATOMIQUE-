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

// PatternPixel represents a single pixel as an atomic unit
// P_ij = (R, G, B) with internal state and connections
type PatternPixel struct {
	Color    [3]float64 // RGB values [0, 1]
	Gradient [3]float64 // Local gradient information
	Velocity [3]float64 // Velocity for momentum-based updates
}

// SeedPoint represents a known/anchored pixel that guides the pattern
// Seeds "anchor" the waves and prevent them from drifting
type SeedPoint struct {
	X, Y   int
	Color  [3]float64 // True RGB value
	Weight float64    // How strongly this seed influences neighbors
}

// PixelConnection represents weights between adjacent pixels
// W_ij;kl = weight from pixel (k,l) to pixel (i,j)
type PixelConnection struct {
	TargetX, TargetY int
	SourceX, SourceY int
	Weight           float64
}

// PatternEmergenceEngine implements pixel-level diffusion and pattern learning
// Transforms abstract waves into recognizable structures through local interactions
type PatternEmergenceEngine struct {
	Width, Height int
	Pixels        [][]PatternPixel
	Connections   map[string][]PixelConnection // "x,y" -> connections from that pixel
	Seeds         []SeedPoint

	// Parameters matching mathematical formulation
	DiffusionAlpha     float64 // α: influence of neighbors (0.0-1.0)
	ReinforcementGamma float64 // γ: weight update rate
	SeedWeight         float64 // How strongly seeds anchor the pattern
	VelocityDamping    float64 // β: momentum damping (0-1)

	// Tracking
	IterationCount int
	AverageLoss    float64
	ConnectedCount int
	mu             sync.RWMutex
}

// NewPatternEmergenceEngine creates a new pattern emergence system
func NewPatternEmergenceEngine(width, height int) *PatternEmergenceEngine {
	pee := &PatternEmergenceEngine{
		Width:              width,
		Height:             height,
		Pixels:             make([][]PatternPixel, height),
		Connections:        make(map[string][]PixelConnection),
		Seeds:              make([]SeedPoint, 0),
		DiffusionAlpha:     0.15, // α: moderate neighbor influence
		ReinforcementGamma: 0.05, // γ: gradual weight strengthening
		SeedWeight:         0.8,  // Seeds strongly anchor
		VelocityDamping:    0.95, // β: smooth momentum
		IterationCount:     0,
		AverageLoss:        0.0,
	}

	// Initialize pixel grid with random values
	for i := 0; i < height; i++ {
		pee.Pixels[i] = make([]PatternPixel, width)
		for j := 0; j < width; j++ {
			pee.Pixels[i][j] = PatternPixel{
				Color:    [3]float64{0.5, 0.5, 0.5}, // Start neutral gray
				Gradient: [3]float64{0, 0, 0},
				Velocity: [3]float64{0, 0, 0},
			}
		}
	}

	// Initialize uniform connections (8-neighborhood)
	pee.initializeConnections()

	return pee
}

// initializeConnections sets up uniform weights between neighboring pixels
// Initially W_ij;kl = 1.0 for all neighbors
func (pee *PatternEmergenceEngine) initializeConnections() {
	for i := 0; i < pee.Height; i++ {
		for j := 0; j < pee.Width; j++ {
			key := fmt.Sprintf("%d,%d", j, i)
			connections := make([]PixelConnection, 0)

			// 8-neighborhood
			for di := -1; di <= 1; di++ {
				for dj := -1; dj <= 1; dj++ {
					if di == 0 && dj == 0 {
						continue // Skip self
					}

					ni, nj := i+di, j+dj
					if ni >= 0 && ni < pee.Height && nj >= 0 && nj < pee.Width {
						conn := PixelConnection{
							TargetX: j,
							TargetY: i,
							SourceX: nj,
							SourceY: ni,
							Weight:  1.0, // Uniform initially
						}
						connections = append(connections, conn)
					}
				}
			}

			pee.Connections[key] = connections
			pee.ConnectedCount += len(connections)
		}
	}
}

// AddSeedPoint anchors a pixel to a known value
// P_ij = P_real if (i,j) ∈ seeds, else compute via diffusion
func (pee *PatternEmergenceEngine) AddSeedPoint(x, y int, r, g, b float64) error {
	if x < 0 || x >= pee.Width || y < 0 || y >= pee.Height {
		return fmt.Errorf("seed point out of bounds: (%d,%d)", x, y)
	}

	pee.mu.Lock()
	defer pee.mu.Unlock()

	// Normalize RGB to [0,1]
	r = math.Max(0, math.Min(1, r/255.0))
	g = math.Max(0, math.Min(1, g/255.0))
	b = math.Max(0, math.Min(1, b/255.0))

	seed := SeedPoint{
		X:      x,
		Y:      y,
		Color:  [3]float64{r, g, b},
		Weight: pee.SeedWeight,
	}
	pee.Seeds = append(pee.Seeds, seed)

	// Immediately set pixel to seed value
	pee.Pixels[y][x].Color = seed.Color

	return nil
}

// AddSeedsFromImage extracts seeds from reference image at sample density
// Higher density = more anchoring points, more structure
func (pee *PatternEmergenceEngine) AddSeedsFromImage(img image.Image, sampleDensity float64) error {
	bounds := img.Bounds()
	scaleX := float64(pee.Width) / float64(bounds.Dx())
	scaleY := float64(pee.Height) / float64(bounds.Dy())

	pee.mu.Lock()
	defer pee.mu.Unlock()

	// Determine sampling step (1 = every pixel, 2 = every other, etc.)
	sampleStep := int(math.Max(1, 1.0/sampleDensity))

	count := 0
	for y := bounds.Min.Y; y < bounds.Max.Y; y += sampleStep {
		for x := bounds.Min.X; x < bounds.Max.X; x += sampleStep {
			// Map to target image coordinates
			tx := int(float64(x-bounds.Min.X) * scaleX)
			ty := int(float64(y-bounds.Min.Y) * scaleY)

			if tx >= 0 && tx < pee.Width && ty >= 0 && ty < pee.Height {
				r, g, b, _ := img.At(x, y).RGBA()
				// Convert from 16-bit to [0,1]
				seed := SeedPoint{
					X:      tx,
					Y:      ty,
					Color:  [3]float64{float64(r) / 65535.0, float64(g) / 65535.0, float64(b) / 65535.0},
					Weight: pee.SeedWeight,
				}
				pee.Seeds = append(pee.Seeds, seed)
				pee.Pixels[ty][tx].Color = seed.Color
				count++
			}
		}
	}

	fmt.Printf("  📌 Added %d seed points at density %.1f\n", count, sampleDensity)
	return nil
}

// DiffuseStep performs one iteration of pixel diffusion
// P_ij(t+1) = P_ij(t) + α·Σ W_ij;kl·(P_kl(t) - P_ij(t))
func (pee *PatternEmergenceEngine) DiffuseStep() {
	pee.mu.Lock()
	defer pee.mu.Unlock()

	newPixels := make([][]PatternPixel, pee.Height)
	totalLoss := 0.0

	for i := 0; i < pee.Height; i++ {
		newPixels[i] = make([]PatternPixel, pee.Width)
		copy(newPixels[i], pee.Pixels[i])
	}

	// Apply diffusion to each pixel
	for i := 0; i < pee.Height; i++ {
		for j := 0; j < pee.Width; j++ {
			key := fmt.Sprintf("%d,%d", j, i)
			neighbors := pee.Connections[key]

			if len(neighbors) == 0 {
				continue
			}

			// Compute weighted sum of neighbor influences
			influence := [3]float64{0, 0, 0}
			for _, conn := range neighbors {
				neighborColor := pee.Pixels[conn.SourceY][conn.SourceX].Color
				diff := [3]float64{
					neighborColor[0] - pee.Pixels[i][j].Color[0],
					neighborColor[1] - pee.Pixels[i][j].Color[1],
					neighborColor[2] - pee.Pixels[i][j].Color[2],
				}

				// W_ij;kl · (P_kl - P_ij)
				for c := 0; c < 3; c++ {
					influence[c] += conn.Weight * diff[c]
				}
			}

			// P_ij(t+1) = P_ij(t) + α·influence + β·velocity
			for c := 0; c < 3; c++ {
				newPixels[i][j].Velocity[c] = pee.VelocityDamping * newPixels[i][j].Velocity[c]
				newPixels[i][j].Color[c] += pee.DiffusionAlpha*influence[c] + newPixels[i][j].Velocity[c]

				// Store new velocity for next iteration
				newPixels[i][j].Velocity[c] = pee.DiffusionAlpha * influence[c]

				// Clamp to [0, 1]
				newPixels[i][j].Color[c] = math.Max(0, math.Min(1, newPixels[i][j].Color[c]))
			}

			// Track loss for monitoring
			for _, seed := range pee.Seeds {
				if seed.X == j && seed.Y == i {
					for c := 0; c < 3; c++ {
						diff := newPixels[i][j].Color[c] - seed.Color[c]
						totalLoss += diff * diff
					}
					break
				}
			}
		}
	}

	// Apply seed constraints: P_ij = P_real if (i,j) ∈ seeds
	for _, seed := range pee.Seeds {
		newPixels[seed.Y][seed.X].Color = seed.Color
	}

	pee.Pixels = newPixels
	pee.IterationCount++

	if len(pee.Seeds) > 0 {
		pee.AverageLoss = totalLoss / float64(len(pee.Seeds))
	}
}

// ReinforceConnections strengthens weights between similar-colored neighboring pixels
// W_ij;kl(t+1) = W_ij;kl(t) + γ·exp(-||P_ij - P_kl||²)
// This amplifies correct patterns
func (pee *PatternEmergenceEngine) ReinforceConnections() {
	pee.mu.Lock()
	defer pee.mu.Unlock()

	for i := 0; i < pee.Height; i++ {
		for j := 0; j < pee.Width; j++ {
			key := fmt.Sprintf("%d,%d", j, i)
			neighbors := pee.Connections[key]

			for idx, conn := range neighbors {
				neighborColor := pee.Pixels[conn.SourceY][conn.SourceX].Color
				pixelColor := pee.Pixels[i][j].Color

				// Compute color difference
				colorDiff := 0.0
				for c := 0; c < 3; c++ {
					diff := neighborColor[c] - pixelColor[c]
					colorDiff += diff * diff
				}

				// γ·exp(-||P_ij - P_kl||²)
				reinforcement := pee.ReinforcementGamma * math.Exp(-colorDiff)
				pee.Connections[key][idx].Weight += reinforcement

				// Prevent weight explosion
				if pee.Connections[key][idx].Weight > 10.0 {
					pee.Connections[key][idx].Weight = 10.0
				}
			}
		}
	}
}

// IteratePattern runs a full pattern emergence cycle
// Steps: diffuse → reinforce → apply seeds
func (pee *PatternEmergenceEngine) IteratePattern(steps int) {
	for step := 0; step < steps; step++ {
		pee.DiffuseStep()

		// Reinforce every 5 diffusion steps to avoid instability
		if step%5 == 0 {
			pee.ReinforceConnections()
		}

		if step%10 == 0 {
			fmt.Printf("  ✓ Step %d/%d | Loss: %.6f | Iter: %d\n",
				step+1, steps, pee.AverageLoss, pee.IterationCount)
		}
	}
}

// ComputeGradient calculates local gradients for each pixel
// Used for edge detection and texture analysis
func (pee *PatternEmergenceEngine) ComputeGradient() {
	pee.mu.Lock()
	defer pee.mu.Unlock()

	for i := 1; i < pee.Height-1; i++ {
		for j := 1; j < pee.Width-1; j++ {
			// Sobel-like operators
			gx := 0.0
			gy := 0.0

			for c := 0; c < 3; c++ {
				// X gradient
				gx += pee.Pixels[i-1][j-1].Color[c]*-1 + pee.Pixels[i-1][j+1].Color[c]*1 +
					pee.Pixels[i][j-1].Color[c]*-2 + pee.Pixels[i][j+1].Color[c]*2 +
					pee.Pixels[i+1][j-1].Color[c]*-1 + pee.Pixels[i+1][j+1].Color[c]*1

				// Y gradient
				gy += pee.Pixels[i-1][j-1].Color[c]*-1 + pee.Pixels[i+1][j-1].Color[c]*1 +
					pee.Pixels[i-1][j].Color[c]*-2 + pee.Pixels[i+1][j].Color[c]*2 +
					pee.Pixels[i-1][j+1].Color[c]*-1 + pee.Pixels[i+1][j+1].Color[c]*1
			}

			magnitude := math.Sqrt(gx*gx+gy*gy) / 8.0
			pee.Pixels[i][j].Gradient[0] = math.Max(0, math.Min(1, gx/3.0))
			pee.Pixels[i][j].Gradient[1] = math.Max(0, math.Min(1, gy/3.0))
			pee.Pixels[i][j].Gradient[2] = math.Max(0, math.Min(1, magnitude))
		}
	}
}

// SaveImage exports current pattern state as PNG
func (pee *PatternEmergenceEngine) SaveImage(filename string) error {
	pee.mu.RLock()
	defer pee.mu.RUnlock()

	img := image.NewRGBA(image.Rect(0, 0, pee.Width, pee.Height))

	for i := 0; i < pee.Height; i++ {
		for j := 0; j < pee.Width; j++ {
			r := uint8(pee.Pixels[i][j].Color[0] * 255)
			g := uint8(pee.Pixels[i][j].Color[1] * 255)
			b := uint8(pee.Pixels[i][j].Color[2] * 255)
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

// GetPixel safely returns a pixel value
func (pee *PatternEmergenceEngine) GetPixel(x, y int) [3]float64 {
	pee.mu.RLock()
	defer pee.mu.RUnlock()

	if x >= 0 && x < pee.Width && y >= 0 && y < pee.Height {
		return pee.Pixels[y][x].Color
	}
	return [3]float64{0.5, 0.5, 0.5}
}

// GetStats returns pattern emergence statistics
func (pee *PatternEmergenceEngine) GetStats() map[string]interface{} {
	pee.mu.RLock()
	defer pee.mu.RUnlock()

	return map[string]interface{}{
		"iterations":   pee.IterationCount,
		"average_loss": pee.AverageLoss,
		"seed_count":   len(pee.Seeds),
		"dimensions":   fmt.Sprintf("%dx%d", pee.Width, pee.Height),
		"alpha":        pee.DiffusionAlpha,
		"gamma":        pee.ReinforcementGamma,
		"seed_weight":  pee.SeedWeight,
	}
}

// PrintStatistics outputs pattern emergence metrics
func (pee *PatternEmergenceEngine) PrintStatistics() {
	stats := pee.GetStats()
	fmt.Printf("\n📊 Pattern Emergence Statistics\n")
	fmt.Printf("═══════════════════════════════════════\n")
	fmt.Printf("Iterations:     %v\n", stats["iterations"])
	fmt.Printf("Average Loss:   %.8f\n", stats["average_loss"])
	fmt.Printf("Seed Points:    %v\n", stats["seed_count"])
	fmt.Printf("Dimensions:     %v\n", stats["dimensions"])
	fmt.Printf("Diffusion α:    %.4f\n", stats["alpha"])
	fmt.Printf("Reinforcement γ: %.4f\n", stats["gamma"])
	fmt.Printf("Seed Weight:    %.4f\n", stats["seed_weight"])
	fmt.Printf("═══════════════════════════════════════\n\n")
}
