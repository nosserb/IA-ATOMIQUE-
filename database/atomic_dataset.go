// Package database - Atomic Dataset Management
// Loads reference images and creates target atomic states for training
// Extracts micro-patterns (gradients, edges) as state vectors

package database

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"math"
	"os"
	"path/filepath"
	"sync"
)

// TargetAtomicState represents the goal state for a pixel/block
type TargetAtomicState struct {
	PixelX, PixelY int        // Position in image
	ColorTarget    [3]float64 // Target RGB [0, 1]
	GradientX      float64    // X-direction gradient
	GradientY      float64    // Y-direction gradient
	EdgeStrength   float64    // Edge magnitude
	Curvature      float64    // Local curvature
	TextureEnergy  float64    // Texture complexity
	Orientation    float64    // Dominant orientation [0, 2π]
}

// AtomicDataset represents a collection of reference images
type AtomicDataset struct {
	Images       []image.Image         // Loaded images
	TargetStates [][]TargetAtomicState // Target states per image
	ImagePaths   []string              // Source file paths
	DatasetName  string
	PatchSize    int // Patch size for blocks
	TotalPixels  int
	Statistics   DatasetStatistics
	mutex        sync.RWMutex
}

// DatasetStatistics contains aggregate statistics
type DatasetStatistics struct {
	NumImages            int
	AvgImageWidth        float64
	AvgImageHeight       float64
	AvgGradientMagnitude float64
	AvgEdgeStrength      float64
	AvgTextureEnergy     float64
	ColorHistogram       map[string]int // Distribution of colors
}

// NewAtomicDataset creates an empty dataset
func NewAtomicDataset(name string, patchSize int) *AtomicDataset {
	return &AtomicDataset{
		Images:       make([]image.Image, 0),
		TargetStates: make([][]TargetAtomicState, 0),
		ImagePaths:   make([]string, 0),
		DatasetName:  name,
		PatchSize:    patchSize,
		TotalPixels:  0,
		Statistics:   DatasetStatistics{ColorHistogram: make(map[string]int)},
	}
}

// LoadImage loads a single image and extracts atomic target states
func (ds *AtomicDataset) LoadImage(filePath string) error {
	ds.mutex.Lock()
	defer ds.mutex.Unlock()

	// Open image file
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open image: %v", err)
	}
	defer file.Close()

	// Decode image
	var img image.Image
	ext := filepath.Ext(filePath)

	switch ext {
	case ".png":
		img, err = png.Decode(file)
	case ".jpg", ".jpeg":
		img, err = jpeg.Decode(file)
	default:
		return fmt.Errorf("unsupported image format: %s", ext)
	}

	if err != nil {
		return fmt.Errorf("failed to decode image: %v", err)
	}

	// Extract atomic states from image
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	states := ds.extractAtomicStates(img, width, height)

	// Store
	ds.Images = append(ds.Images, img)
	ds.TargetStates = append(ds.TargetStates, states)
	ds.ImagePaths = append(ds.ImagePaths, filePath)
	ds.TotalPixels += width * height

	return nil
}

// LoadDirectory loads all images from a directory
func (ds *AtomicDataset) LoadDirectory(dirPath string) (int, error) {
	ds.mutex.Lock()
	ds.mutex.Unlock()

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return 0, err
	}

	loadedCount := 0
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		ext := filepath.Ext(entry.Name())
		if ext != ".png" && ext != ".jpg" && ext != ".jpeg" {
			continue
		}

		fullPath := filepath.Join(dirPath, entry.Name())
		if err := ds.LoadImage(fullPath); err != nil {
			fmt.Printf("⚠️  Failed to load %s: %v\n", entry.Name(), err)
			continue
		}
		loadedCount++
	}

	ds.mutex.Lock()
	defer ds.mutex.Unlock()
	ds.Statistics.NumImages = len(ds.Images)

	return loadedCount, nil
}

// extractAtomicStates converts an image to atomic target states
func (ds *AtomicDataset) extractAtomicStates(
	img image.Image,
	width, height int,
) []TargetAtomicState {
	states := make([]TargetAtomicState, width*height)
	bounds := img.Bounds()

	// Precompute gradients for entire image
	gradientX := make([][]float64, height)
	gradientY := make([][]float64, height)
	edgeStrength := make([][]float64, height)
	curvature := make([][]float64, height)

	for y := 0; y < height; y++ {
		gradientX[y] = make([]float64, width)
		gradientY[y] = make([]float64, width)
		edgeStrength[y] = make([]float64, width)
		curvature[y] = make([]float64, width)
	}

	// Compute gradients
	for y := 1; y < height-1; y++ {
		for x := 1; x < width-1; x++ {
			// Get neighboring pixels
			r0, g0, b0, _ := img.At(bounds.Min.X+x-1, bounds.Min.Y+y).RGBA()
			r1, g1, b1, _ := img.At(bounds.Min.X+x+1, bounds.Min.Y+y).RGBA()
			r2, g2, b2, _ := img.At(bounds.Min.X+x, bounds.Min.Y+y-1).RGBA()
			r3, g3, b3, _ := img.At(bounds.Min.X+x, bounds.Min.Y+y+1).RGBA()

			// Convert to [0,1] range and compute gradients
			normalizeColor := func(r, g, b uint32) float64 {
				return float64(r+g+b) / (3 * 65536)
			}

			// X gradient
			c0 := normalizeColor(r0, g0, b0)
			c1 := normalizeColor(r1, g1, b1)
			gradientX[y][x] = (c1 - c0) / 2.0

			// Y gradient
			c2 := normalizeColor(r2, g2, b2)
			c3 := normalizeColor(r3, g3, b3)
			gradientY[y][x] = (c3 - c2) / 2.0

			// Edge strength (magnitude of gradient)
			edgeStrength[y][x] = math.Sqrt(
				gradientX[y][x]*gradientX[y][x] +
					gradientY[y][x]*gradientY[y][x],
			)

			// Curvature (Laplacian)
			r, g, b, _ := img.At(bounds.Min.X+x, bounds.Min.Y+y).RGBA()
			center := normalizeColor(r, g, b)
			laplacian := (c0 + c1 + c2 + c3 - 4*center) / 4.0
			curvature[y][x] = math.Abs(laplacian)
		}
	}

	// Create state for each pixel
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			idx := y*width + x

			// Get pixel color
			r, g, b, _ := img.At(bounds.Min.X+x, bounds.Min.Y+y).RGBA()
			colorR := float64(r) / 65536
			colorG := float64(g) / 65536
			colorB := float64(b) / 65536

			// Clamp to [0, 1]
			if colorR > 1.0 {
				colorR = 1.0
			}
			if colorG > 1.0 {
				colorG = 1.0
			}
			if colorB > 1.0 {
				colorB = 1.0
			}

			state := TargetAtomicState{
				PixelX:        x,
				PixelY:        y,
				ColorTarget:   [3]float64{colorR, colorG, colorB},
				GradientX:     0.0,
				GradientY:     0.0,
				EdgeStrength:  0.0,
				Curvature:     0.0,
				TextureEnergy: 0.0,
				Orientation:   0.0,
			}

			// Add gradient info if in valid range
			if y > 0 && y < height-1 && x > 0 && x < width-1 {
				state.GradientX = gradientX[y][x]
				state.GradientY = gradientY[y][x]
				state.EdgeStrength = edgeStrength[y][x]
				state.Curvature = curvature[y][x]

				// Orientation: atan2(gradY, gradX)
				if state.GradientX != 0.0 || state.GradientY != 0.0 {
					state.Orientation = math.Atan2(state.GradientY, state.GradientX)
					if state.Orientation < 0 {
						state.Orientation += 2 * math.Pi
					}
				}

				// Texture energy: combination of edge and curvature
				state.TextureEnergy = state.EdgeStrength*0.7 + state.Curvature*0.3
			}

			states[idx] = state
		}
	}

	return states
}

// GetImage returns an image by index
func (ds *AtomicDataset) GetImage(idx int) (image.Image, error) {
	ds.mutex.RLock()
	defer ds.mutex.RUnlock()

	if idx < 0 || idx >= len(ds.Images) {
		return nil, fmt.Errorf("image index out of range: %d", idx)
	}

	return ds.Images[idx], nil
}

// GetTargetStates returns target states for an image
func (ds *AtomicDataset) GetTargetStates(idx int) ([]TargetAtomicState, error) {
	ds.mutex.RLock()
	defer ds.mutex.RUnlock()

	if idx < 0 || idx >= len(ds.TargetStates) {
		return nil, fmt.Errorf("states index out of range: %d", idx)
	}

	return ds.TargetStates[idx], nil
}

// ComputeStatistics analyzes the dataset
func (ds *AtomicDataset) ComputeStatistics() {
	ds.mutex.Lock()
	defer ds.mutex.Unlock()

	stats := DatasetStatistics{
		NumImages:      len(ds.Images),
		ColorHistogram: make(map[string]int),
	}

	if stats.NumImages == 0 {
		ds.Statistics = stats
		return
	}

	totalGradient := 0.0
	totalEdge := 0.0
	totalTexture := 0.0
	stateCount := 0

	totalWidth := 0.0
	totalHeight := 0.0

	for i, img := range ds.Images {
		bounds := img.Bounds()
		width := float64(bounds.Dx())
		height := float64(bounds.Dy())

		totalWidth += width
		totalHeight += height

		states := ds.TargetStates[i]
		for _, state := range states {
			totalGradient += math.Sqrt(
				state.GradientX*state.GradientX +
					state.GradientY*state.GradientY,
			)
			totalEdge += state.EdgeStrength
			totalTexture += state.TextureEnergy
			stateCount++

			// Color histogram
			colorKey := fmt.Sprintf("%.1f_%.1f_%.1f",
				state.ColorTarget[0],
				state.ColorTarget[1],
				state.ColorTarget[2],
			)
			stats.ColorHistogram[colorKey]++
		}
	}

	if stateCount > 0 {
		stats.AvgGradientMagnitude = totalGradient / float64(stateCount)
		stats.AvgEdgeStrength = totalEdge / float64(stateCount)
		stats.AvgTextureEnergy = totalTexture / float64(stateCount)
	}

	if stats.NumImages > 0 {
		stats.AvgImageWidth = totalWidth / float64(stats.NumImages)
		stats.AvgImageHeight = totalHeight / float64(stats.NumImages)
	}

	ds.Statistics = stats
}

// PrintStatistics displays dataset information
func (ds *AtomicDataset) PrintStatistics() {
	ds.mutex.RLock()
	defer ds.mutex.RUnlock()

	fmt.Printf("\n📊 Dataset Statistics: %s\n", ds.DatasetName)
	fmt.Printf("═══════════════════════════════════════\n")
	fmt.Printf("Number of images: %d\n", ds.Statistics.NumImages)
	fmt.Printf("Average size: %.0f × %.0f\n", ds.Statistics.AvgImageWidth, ds.Statistics.AvgImageHeight)
	fmt.Printf("Total pixels: %d\n", ds.TotalPixels)
	fmt.Printf("\n🎨 Image Statistics:\n")
	fmt.Printf("  Avg gradient magnitude: %.4f\n", ds.Statistics.AvgGradientMagnitude)
	fmt.Printf("  Avg edge strength: %.4f\n", ds.Statistics.AvgEdgeStrength)
	fmt.Printf("  Avg texture energy: %.4f\n", ds.Statistics.AvgTextureEnergy)
	fmt.Printf("\n📈 Color Distribution: %d unique colors\n", len(ds.Statistics.ColorHistogram))
	fmt.Printf("═══════════════════════════════════════\n")
}

// ExportTargetImage creates a visual representation of target states
func (ds *AtomicDataset) ExportTargetImage(idx int, outFile string, mode string) error {
	ds.mutex.RLock()
	defer ds.mutex.RUnlock()

	if idx < 0 || idx >= len(ds.TargetStates) {
		return fmt.Errorf("invalid image index")
	}

	img, err := ds.GetImage(idx)
	if err != nil {
		return err
	}

	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	states := ds.TargetStates[idx]

	// Create output image
	outImg := image.NewRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			idx := y*width + x
			if idx >= len(states) {
				continue
			}

			state := states[idx]
			var r, g, b uint8

			switch mode {
			case "color":
				// Original color
				r = uint8(state.ColorTarget[0] * 255)
				g = uint8(state.ColorTarget[1] * 255)
				b = uint8(state.ColorTarget[2] * 255)

			case "gradient":
				// Gradient magnitude
				mag := math.Sqrt(
					state.GradientX*state.GradientX +
						state.GradientY*state.GradientY,
				)
				intensity := uint8(math.Min(255, mag*1000))
				r, g, b = intensity, intensity, intensity

			case "edge":
				// Edge strength
				intensity := uint8(math.Min(255, state.EdgeStrength*500))
				r, g, b = intensity, intensity, intensity

			case "texture":
				// Texture energy
				intensity := uint8(math.Min(255, state.TextureEnergy*500))
				r, g, b = intensity, intensity, intensity

			case "orientation":
				// Orientation as hue
				hue := state.Orientation / (2 * math.Pi)
				r = uint8(math.Max(0, math.Min(255, hue*255)))
				g = uint8(math.Max(0, math.Min(255, (1-hue)*255)))
				b = 128

			default:
				// Color by default
				r = uint8(state.ColorTarget[0] * 255)
				g = uint8(state.ColorTarget[1] * 255)
				b = uint8(state.ColorTarget[2] * 255)
			}

			outImg.Set(x, y, color.RGBA{r, g, b, 255})
		}
	}

	// Save
	outFile2, err := os.Create(outFile)
	if err != nil {
		return err
	}
	defer outFile2.Close()

	return png.Encode(outFile2, outImg)
}

// SampleSubset returns a subset of images for validation
func (ds *AtomicDataset) SampleSubset(percentage float64) *AtomicDataset {
	ds.mutex.RLock()
	defer ds.mutex.RUnlock()

	subset := NewAtomicDataset(ds.DatasetName+"_subset", ds.PatchSize)

	sampleSize := int(float64(len(ds.Images)) * percentage / 100.0)
	if sampleSize < 1 {
		sampleSize = 1
	}

	for i := 0; i < sampleSize && i < len(ds.Images); i++ {
		subset.Images = append(subset.Images, ds.Images[i])
		subset.TargetStates = append(subset.TargetStates, ds.TargetStates[i])
		subset.ImagePaths = append(subset.ImagePaths, ds.ImagePaths[i])
	}

	return subset
}

// GetDatasetSize returns number of images
func (ds *AtomicDataset) GetDatasetSize() int {
	ds.mutex.RLock()
	defer ds.mutex.RUnlock()
	return len(ds.Images)
}
