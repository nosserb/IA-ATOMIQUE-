package database

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"os"
)

// loadRGBImage decodes an image into a [height][width][3] float matrix in [0,1].
func loadRGBImage(path string) ([][][3]float64, int, int, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("cannot open image %s: %w", path, err)
	}
	defer file.Close()

	// Register decoders by importing image/jpeg and image/png above
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("cannot decode image %s: %w", path, err)
	}

	bounds := img.Bounds()
	w := bounds.Dx()
	h := bounds.Dy()

	data := make([][][3]float64, h)
	for y := 0; y < h; y++ {
		data[y] = make([][3]float64, w)
		for x := 0; x < w; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			data[y][x] = [3]float64{float64(r) / 65535.0, float64(g) / 65535.0, float64(b) / 65535.0}
		}
	}

	return data, w, h, nil
}

// loadMaskImage decodes a mask to [height][width] float values in [0,1].
func loadMaskImage(path string) ([][]float64, int, int, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("cannot open mask %s: %w", path, err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("cannot decode mask %s: %w", path, err)
	}

	bounds := img.Bounds()
	w := bounds.Dx()
	h := bounds.Dy()

	data := make([][]float64, h)
	for y := 0; y < h; y++ {
		data[y] = make([]float64, w)
		for x := 0; x < w; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			// Convert to grayscale intensity
			data[y][x] = (float64(r) + float64(g) + float64(b)) / (3.0 * 65535.0)
		}
	}

	return data, w, h, nil
}

func resizeRGBNearest(src [][][3]float64, newH, newW int) [][][3]float64 {
	if len(src) == 0 || newH <= 0 || newW <= 0 {
		return nil
	}

	out := make([][][3]float64, newH)
	h := len(src)
	w := len(src[0])
	for y := 0; y < newH; y++ {
		out[y] = make([][3]float64, newW)
		srcY := int(float64(y) / float64(newH) * float64(h))
		if srcY >= h {
			srcY = h - 1
		}
		for x := 0; x < newW; x++ {
			srcX := int(float64(x) / float64(newW) * float64(w))
			if srcX >= w {
				srcX = w - 1
			}
			out[y][x] = src[srcY][srcX]
		}
	}

	return out
}

func resizeMaskNearest(src [][]float64, newH, newW int) [][]float64 {
	if len(src) == 0 || newH <= 0 || newW <= 0 {
		return nil
	}

	out := make([][]float64, newH)
	h := len(src)
	w := len(src[0])
	for y := 0; y < newH; y++ {
		out[y] = make([]float64, newW)
		srcY := int(float64(y) / float64(newH) * float64(h))
		if srcY >= h {
			srcY = h - 1
		}
		for x := 0; x < newW; x++ {
			srcX := int(float64(x) / float64(newW) * float64(w))
			if srcX >= w {
				srcX = w - 1
			}
			out[y][x] = src[srcY][srcX]
		}
	}

	return out
}

func sampleRegionRGBToGrid(src [][][3]float64, x0, x1, y0, y1, grid int) [][][3]float64 {
	regionH := y1 - y0
	regionW := x1 - x0
	out := make([][][3]float64, grid)

	for gy := 0; gy < grid; gy++ {
		out[gy] = make([][3]float64, grid)
		srcY := y0 + int(float64(gy)/float64(grid)*float64(regionH))
		if srcY >= len(src) {
			srcY = len(src) - 1
		}
		for gx := 0; gx < grid; gx++ {
			srcX := x0 + int(float64(gx)/float64(grid)*float64(regionW))
			if srcX >= len(src[0]) {
				srcX = len(src[0]) - 1
			}
			out[gy][gx] = src[srcY][srcX]
		}
	}

	return out
}

func sampleRegionMaskToGrid(src [][]float64, x0, x1, y0, y1, grid int) [][]float64 {
	regionH := y1 - y0
	regionW := x1 - x0
	out := make([][]float64, grid)

	for gy := 0; gy < grid; gy++ {
		out[gy] = make([]float64, grid)
		srcY := y0 + int(float64(gy)/float64(grid)*float64(regionH))
		if srcY >= len(src) {
			srcY = len(src) - 1
		}
		for gx := 0; gx < grid; gx++ {
			srcX := x0 + int(float64(gx)/float64(grid)*float64(regionW))
			if srcX >= len(src[0]) {
				srcX = len(src[0]) - 1
			}
			out[gy][gx] = src[srcY][srcX]
		}
	}

	return out
}

// InitializeFusionPatches prepares the grid to fuse a masked element into the base image.
// Returns the base image width/height for downstream export.
func (grid *OptimizedPatchGrid) InitializeFusionPatches(basePath, elementPath, maskPath string) (int, int, error) {
	baseRGB, baseW, baseH, err := loadRGBImage(basePath)
	if err != nil {
		return 0, 0, err
	}

	elementRGB, elemW, elemH, err := loadRGBImage(elementPath)
	if err != nil {
		return 0, 0, err
	}
	if elemW != baseW || elemH != baseH {
		elementRGB = resizeRGBNearest(elementRGB, baseH, baseW)
	}

	maskData, maskW, maskH, err := loadMaskImage(maskPath)
	if err != nil {
		return 0, 0, err
	}
	if maskW != baseW || maskH != baseH {
		maskData = resizeMaskNearest(maskData, baseH, baseW)
	}

	gridH := len(grid.Patches)
	gridW := len(grid.Patches[0])

	for i := 0; i < gridH; i++ {
		for j := 0; j < gridW; j++ {
			y0 := (i * baseH) / gridH
			y1 := ((i + 1) * baseH) / gridH
			x0 := (j * baseW) / gridW
			x1 := ((j + 1) * baseW) / gridW

			// Extract region references
			patchPixels := make([][][3]float64, y1-y0)
			for yy := y0; yy < y1; yy++ {
				patchPixels[yy-y0] = baseRGB[yy][x0:x1]
			}

			patch := CreateOptimizedPatchRGB(patchPixels, grid.AdaptiveStrategy)
			gridSize := len(patch.Atoms)

			baseSample := sampleRegionRGBToGrid(baseRGB, x0, x1, y0, y1, gridSize)
			constraintSample := sampleRegionRGBToGrid(elementRGB, x0, x1, y0, y1, gridSize)
			maskSample := sampleRegionMaskToGrid(maskData, x0, x1, y0, y1, gridSize)

			// Fill fusion buffers and atoms with blended start state
			var maskSum float64
			patch.BaseRGB = make([][][3]float64, gridSize)
			patch.ConstraintRGB = make([][][3]float64, gridSize)
			patch.Mask = make([][]float64, gridSize)

			for yy := 0; yy < gridSize; yy++ {
				patch.BaseRGB[yy] = make([][3]float64, gridSize)
				patch.ConstraintRGB[yy] = make([][3]float64, gridSize)
				patch.Mask[yy] = make([]float64, gridSize)

				for xx := 0; xx < gridSize; xx++ {
					patch.BaseRGB[yy][xx] = baseSample[yy][xx]
					patch.ConstraintRGB[yy][xx] = constraintSample[yy][xx]
					m := math.Max(0, math.Min(1, maskSample[yy][xx]))
					patch.Mask[yy][xx] = m
					maskSum += m

					fusedR := baseSample[yy][xx][0]*(1-m) + constraintSample[yy][xx][0]*m
					fusedG := baseSample[yy][xx][1]*(1-m) + constraintSample[yy][xx][1]*m
					fusedB := baseSample[yy][xx][2]*(1-m) + constraintSample[yy][xx][2]*m

					patch.Atoms[yy][xx].R = fusedR
					patch.Atoms[yy][xx].G = fusedG
					patch.Atoms[yy][xx].B = fusedB
					patch.Atoms[yy][xx].Intensity = (fusedR + fusedG + fusedB) / 3.0
				}
			}

			patch.HasMask = maskSum > 1e-3
			patch.OriginalGradient = ComputeGradientFieldRGB(patch.Atoms)

			grid.Patches[i][j] = *patch

			if patch.HasMask {
				grid.Mask.MarkModified(i, j)
			} else {
				grid.Mask.MarkConverged(i, j)
			}
		}
	}

	return baseW, baseH, nil
}

// RunFusionPipeline executes masked fusion using the energy weights provided.
func (grid *OptimizedPatchGrid) RunFusionPipeline(iterations int, alpha, beta, gamma, lambda float64) {
	gridH := len(grid.Patches)
	gridW := len(grid.Patches[0])

	for iter := 0; iter < iterations; iter++ {
		for i := 0; i < gridH; i++ {
			for j := 0; j < gridW; j++ {
				patch := &grid.Patches[i][j]
				if !patch.HasMask {
					continue
				}

				neighbors := map[string]*OptimizedPatch{}
				if i > 0 {
					neighbors["up"] = &grid.Patches[i-1][j]
				}
				if i < gridH-1 {
					neighbors["down"] = &grid.Patches[i+1][j]
				}
				if j > 0 {
					neighbors["left"] = &grid.Patches[i][j-1]
				}
				if j < gridW-1 {
					neighbors["right"] = &grid.Patches[i][j+1]
				}

				patch.RelaxFusionStep(alpha, beta, gamma, lambda, grid.LearningRate, neighbors)
			}
		}
	}
}
