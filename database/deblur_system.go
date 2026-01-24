// Package database - Deblur System
// Advanced deblurring with gradient amplification and deconvolution
//
// 1️⃣ Gradient extraction: ∇I_blur from original blurry image
// 2️⃣ Energy minimization: E_sharpen = Σ ||∇I_recon - k·∇I_blur||² where k > 1
//    k > 1 → amplifies details (adaptive sharpening)
// 3️⃣ Gaussian deconvolution: I_sharp ≈ Deconvolve(I_blur, G_σ)
// 4️⃣ Multi-phase: Coarse (16×16) → Medium (8×8) → Fine (4×4)

package database

import (
	"fmt"
	"math"
)

// ============================================================================
// DECONVOLUTION & SHARPENING PARAMETERS
// ============================================================================

// DeconvolutionParams controls the deblurring intensity
type DeconvolutionParams struct {
	GradientAmplification float64 // k > 1: amplify gradients for detail hallucination
	GaussianSigma         float64 // σ for Gaussian blur kernel (if known)
	SharpeningStrength    float64 // 0-1: how much to enhance details
	NoiseReduction        float64 // 0-1: prevent noise amplification
	EdgeEnhancementLambda float64 // λ > 0: rewards high-gradient regions (edges)
}

// NewDefaultDeconvolutionParams creates sensible defaults
func NewDefaultDeconvolutionParams() *DeconvolutionParams {
	return &DeconvolutionParams{
		GradientAmplification: 2.2,  // 120% gradient boost (aggressive detail hallucination)
		GaussianSigma:         1.5,  // Assume moderate blur
		SharpeningStrength:    0.85, // 85% sharpening
		NoiseReduction:        0.2,  // 20% noise suppression (allow more texture)
		EdgeEnhancementLambda: 0.4,  // λ = 0.4: moderate edge reward (0.0-1.0)
	}
}

// NewAdaptiveDeconvolutionParams creates params based on blur level
func NewAdaptiveDeconvolutionParams(blurSigma float64) *DeconvolutionParams {
	// More blur → more aggressive amplification
	k := 1.5 + (blurSigma * 0.5) // k ∈ [1.5, 3.0]
	if k > 3.0 {
		k = 3.0
	}
	// More blur → more edge enhancement
	lambda := 0.3 + (blurSigma * 0.15) // λ ∈ [0.3, 1.0]
	if lambda > 1.0 {
		lambda = 1.0
	}
	return &DeconvolutionParams{
		GradientAmplification: k,
		GaussianSigma:         blurSigma,
		SharpeningStrength:    0.8,
		NoiseReduction:        0.15, // Less suppression for texture
		EdgeEnhancementLambda: lambda,
	}
}

// ============================================================================
// GAUSSIAN DECONVOLUTION (SIMPLIFIED WIENER FILTER)
// ============================================================================

// ApplyGaussianDeconvolution applies simplified Wiener-like deconvolution
// This is a spatial-domain approximation (not full FFT Wiener)
func ApplyGaussianDeconvolution(atoms [][]PixelAtomV2, sigma float64, strength float64) {
	h := len(atoms)
	w := len(atoms[0])
	if h < 3 || w < 3 {
		return
	}

	// Create temporary buffer for deconvolved values
	tempR := make([][]float64, h)
	tempG := make([][]float64, h)
	tempB := make([][]float64, h)
	for i := range tempR {
		tempR[i] = make([]float64, w)
		tempG[i] = make([]float64, w)
		tempB[i] = make([]float64, w)
	}

	// Simplified Gaussian blur inversion (unsharp mask + amplification)
	// I_sharp ≈ I_blur + k·(I_blur - Gaussian(I_blur))
	kernelSize := int(math.Ceil(3 * sigma))
	if kernelSize%2 == 0 {
		kernelSize++
	}
	halfKernel := kernelSize / 2

	// Generate Gaussian kernel (Point Spread Function)
	kernel := make([][]float64, kernelSize)
	kernelSum := 0.0
	for i := 0; i < kernelSize; i++ {
		kernel[i] = make([]float64, kernelSize)
		for j := 0; j < kernelSize; j++ {
			x := float64(i - halfKernel)
			y := float64(j - halfKernel)
			kernel[i][j] = math.Exp(-(x*x + y*y) / (2 * sigma * sigma))
			kernelSum += kernel[i][j]
		}
	}
	// Normalize kernel (PSF normalization)
	for i := 0; i < kernelSize; i++ {
		for j := 0; j < kernelSize; j++ {
			kernel[i][j] /= kernelSum
		}
	}

	// Apply Gaussian blur to get low-frequency component
	for i := halfKernel; i < h-halfKernel; i++ {
		for j := halfKernel; j < w-halfKernel; j++ {
			var blurR, blurG, blurB float64

			for ki := 0; ki < kernelSize; ki++ {
				for kj := 0; kj < kernelSize; kj++ {
					ii := i + ki - halfKernel
					jj := j + kj - halfKernel
					weight := kernel[ki][kj]

					blurR += atoms[ii][jj].R * weight
					blurG += atoms[ii][jj].G * weight
					blurB += atoms[ii][jj].B * weight
				}
			}

			// Unsharp mask: I_sharp = I + k·(I - Blur(I))
			// Higher k → more sharpening
			detail_amplification := 1.5 * strength

			tempR[i][j] = atoms[i][j].R + detail_amplification*(atoms[i][j].R-blurR)
			tempG[i][j] = atoms[i][j].G + detail_amplification*(atoms[i][j].G-blurG)
			tempB[i][j] = atoms[i][j].B + detail_amplification*(atoms[i][j].B-blurB)

			// Clamp to [0, 1]
			tempR[i][j] = math.Max(0, math.Min(1, tempR[i][j]))
			tempG[i][j] = math.Max(0, math.Min(1, tempG[i][j]))
			tempB[i][j] = math.Max(0, math.Min(1, tempB[i][j]))
		}
	}

	// Copy back to atoms
	for i := halfKernel; i < h-halfKernel; i++ {
		for j := halfKernel; j < w-halfKernel; j++ {
			atoms[i][j].R = tempR[i][j]
			atoms[i][j].G = tempG[i][j]
			atoms[i][j].B = tempB[i][j]
			atoms[i][j].Intensity = (tempR[i][j] + tempG[i][j] + tempB[i][j]) / 3.0
		}
	}
}

// ApplyRichardsonLucyDeconvolution implements the Richardson-Lucy algorithm
// More accurate than unsharp mask, converges to true deconvolution
// Formula: I^(n+1) = I^(n) · [(I_blur / (I^(n) * PSF)) * PSF_mirror]
func ApplyRichardsonLucyDeconvolution(atoms [][]PixelAtomV2, sigma float64, iterations int) {
	h := len(atoms)
	w := len(atoms[0])
	if h < 3 || w < 3 {
		return
	}

	// Build PSF (Point Spread Function = Gaussian kernel)
	kernelSize := int(math.Ceil(3 * sigma))
	if kernelSize%2 == 0 {
		kernelSize++
	}
	halfKernel := kernelSize / 2

	kernel := make([][]float64, kernelSize)
	kernelSum := 0.0
	for i := 0; i < kernelSize; i++ {
		kernel[i] = make([]float64, kernelSize)
		for j := 0; j < kernelSize; j++ {
			x := float64(i - halfKernel)
			y := float64(j - halfKernel)
			kernel[i][j] = math.Exp(-(x*x + y*y) / (2 * sigma * sigma))
			kernelSum += kernel[i][j]
		}
	}
	for i := 0; i < kernelSize; i++ {
		for j := 0; j < kernelSize; j++ {
			kernel[i][j] /= kernelSum
		}
	}

	// Store original blurred image
	blurredR := make([][]float64, h)
	blurredG := make([][]float64, h)
	blurredB := make([][]float64, h)
	for i := 0; i < h; i++ {
		blurredR[i] = make([]float64, w)
		blurredG[i] = make([]float64, w)
		blurredB[i] = make([]float64, w)
		for j := 0; j < w; j++ {
			blurredR[i][j] = atoms[i][j].R
			blurredG[i][j] = atoms[i][j].G
			blurredB[i][j] = atoms[i][j].B
		}
	}

	// Richardson-Lucy iterations
	for iter := 0; iter < iterations; iter++ {
		// Step 1: Convolve current estimate with PSF → I^(n) * PSF
		convolvedR := convolve2D(atoms, kernel, 0)
		convolvedG := convolve2D(atoms, kernel, 1)
		convolvedB := convolve2D(atoms, kernel, 2)

		// Step 2: Compute ratio I_blur / (I^(n) * PSF)
		ratioR := make([][]float64, h)
		ratioG := make([][]float64, h)
		ratioB := make([][]float64, h)
		for i := 0; i < h; i++ {
			ratioR[i] = make([]float64, w)
			ratioG[i] = make([]float64, w)
			ratioB[i] = make([]float64, w)
			for j := 0; j < w; j++ {
				// Avoid division by zero
				epsilon := 1e-6
				ratioR[i][j] = blurredR[i][j] / (convolvedR[i][j] + epsilon)
				ratioG[i][j] = blurredG[i][j] / (convolvedG[i][j] + epsilon)
				ratioB[i][j] = blurredB[i][j] / (convolvedB[i][j] + epsilon)
			}
		}

		// Step 3: Convolve ratio with mirrored PSF → ratio * PSF_mirror
		correctionR := convolve2DRatio(ratioR, kernel)
		correctionG := convolve2DRatio(ratioG, kernel)
		correctionB := convolve2DRatio(ratioB, kernel)

		// Step 4: Update estimate I^(n+1) = I^(n) · correction
		for i := 0; i < h; i++ {
			for j := 0; j < w; j++ {
				atoms[i][j].R *= correctionR[i][j]
				atoms[i][j].G *= correctionG[i][j]
				atoms[i][j].B *= correctionB[i][j]

				// Clamp to [0, 1]
				atoms[i][j].R = math.Max(0, math.Min(1, atoms[i][j].R))
				atoms[i][j].G = math.Max(0, math.Min(1, atoms[i][j].G))
				atoms[i][j].B = math.Max(0, math.Min(1, atoms[i][j].B))
			}
		}
	}
}

// convolve2D performs 2D convolution for a specific RGB channel
func convolve2D(atoms [][]PixelAtomV2, kernel [][]float64, channel int) [][]float64 {
	h := len(atoms)
	w := len(atoms[0])
	kernelSize := len(kernel)
	halfKernel := kernelSize / 2

	result := make([][]float64, h)
	for i := 0; i < h; i++ {
		result[i] = make([]float64, w)
		for j := 0; j < w; j++ {
			var sum float64
			for ki := 0; ki < kernelSize; ki++ {
				for kj := 0; kj < kernelSize; kj++ {
					ii := i + ki - halfKernel
					jj := j + kj - halfKernel

					// Handle boundaries (reflect)
					if ii < 0 {
						ii = -ii
					}
					if ii >= h {
						ii = 2*h - ii - 1
					}
					if jj < 0 {
						jj = -jj
					}
					if jj >= w {
						jj = 2*w - jj - 1
					}

					var value float64
					switch channel {
					case 0:
						value = atoms[ii][jj].R
					case 1:
						value = atoms[ii][jj].G
					case 2:
						value = atoms[ii][jj].B
					}

					sum += value * kernel[ki][kj]
				}
			}
			result[i][j] = sum
		}
	}
	return result
}

// convolve2DRatio performs convolution on a ratio matrix
func convolve2DRatio(ratio [][]float64, kernel [][]float64) [][]float64 {
	h := len(ratio)
	w := len(ratio[0])
	kernelSize := len(kernel)
	halfKernel := kernelSize / 2

	result := make([][]float64, h)
	for i := 0; i < h; i++ {
		result[i] = make([]float64, w)
		for j := 0; j < w; j++ {
			var sum float64
			for ki := 0; ki < kernelSize; ki++ {
				for kj := 0; kj < kernelSize; kj++ {
					ii := i + ki - halfKernel
					jj := j + kj - halfKernel

					// Handle boundaries
					if ii < 0 {
						ii = -ii
					}
					if ii >= h {
						ii = 2*h - ii - 1
					}
					if jj < 0 {
						jj = -jj
					}
					if jj >= w {
						jj = 2*w - jj - 1
					}

					sum += ratio[ii][jj] * kernel[ki][kj]
				}
			}
			result[i][j] = sum
		}
	}
	return result
}

// generateStructuredNoise creates coherent noise pattern for texture synthesis
// Uses pseudo-Perlin noise based on position for realistic texture
func generateStructuredNoise(i, j, channel int) float64 {
	// Simple hash-based noise (Perlin-like)
	seed := float64(i*7919 + j*6553 + channel*4999)

	// Multiple octaves for natural-looking texture
	freq1 := math.Sin(seed*0.1) * 0.5
	freq2 := math.Sin(seed*0.3) * 0.3
	freq3 := math.Sin(seed*0.7) * 0.2

	noise := freq1 + freq2 + freq3

	// Normalize to [-1, 1]
	return noise
}

// ============================================================================
// GRADIENT COMPUTATION & STORAGE
// ============================================================================

// GradientField stores gradient information (∇I) for an image patch
type GradientField struct {
	GradX [][]float64 // Horizontal gradients
	GradY [][]float64 // Vertical gradients
	Mag   [][]float64 // Gradient magnitude
}

// ComputeGradientField calculates ∇I for a patch using Sobel operator
func ComputeGradientField(atoms [][]PixelAtomV2) *GradientField {
	h := len(atoms)
	w := len(atoms[0])

	field := &GradientField{
		GradX: make([][]float64, h),
		GradY: make([][]float64, h),
		Mag:   make([][]float64, h),
	}

	for i := 0; i < h; i++ {
		field.GradX[i] = make([]float64, w)
		field.GradY[i] = make([]float64, w)
		field.Mag[i] = make([]float64, w)
	}

	// Sobel kernels
	// Gx = [-1 0 1; -2 0 2; -1 0 1]
	// Gy = [-1 -2 -1; 0 0 0; 1 2 1]
	for i := 1; i < h-1; i++ {
		for j := 1; j < w-1; j++ {
			// Horizontal gradient (Gx)
			gx := -atoms[i-1][j-1].Intensity + atoms[i-1][j+1].Intensity +
				-2*atoms[i][j-1].Intensity + 2*atoms[i][j+1].Intensity +
				-atoms[i+1][j-1].Intensity + atoms[i+1][j+1].Intensity
			gx /= 8.0

			// Vertical gradient (Gy)
			gy := -atoms[i-1][j-1].Intensity - 2*atoms[i-1][j].Intensity - atoms[i-1][j+1].Intensity +
				atoms[i+1][j-1].Intensity + 2*atoms[i+1][j].Intensity + atoms[i+1][j+1].Intensity
			gy /= 8.0

			field.GradX[i][j] = gx
			field.GradY[i][j] = gy
			field.Mag[i][j] = math.Sqrt(gx*gx + gy*gy)
		}
	}

	return field
}

// ComputeGradientFieldRGB calculates gradients for RGB channels separately
func ComputeGradientFieldRGB(atoms [][]PixelAtomV2) [3]*GradientField {
	h := len(atoms)
	w := len(atoms[0])

	fields := [3]*GradientField{
		{GradX: make([][]float64, h), GradY: make([][]float64, h), Mag: make([][]float64, h)},
		{GradX: make([][]float64, h), GradY: make([][]float64, h), Mag: make([][]float64, h)},
		{GradX: make([][]float64, h), GradY: make([][]float64, h), Mag: make([][]float64, h)},
	}

	for c := 0; c < 3; c++ {
		for i := 0; i < h; i++ {
			fields[c].GradX[i] = make([]float64, w)
			fields[c].GradY[i] = make([]float64, w)
			fields[c].Mag[i] = make([]float64, w)
		}
	}

	for i := 1; i < h-1; i++ {
		for j := 1; j < w-1; j++ {
			// R channel
			gxR := -atoms[i-1][j-1].R + atoms[i-1][j+1].R +
				-2*atoms[i][j-1].R + 2*atoms[i][j+1].R +
				-atoms[i+1][j-1].R + atoms[i+1][j+1].R
			gyR := -atoms[i-1][j-1].R - 2*atoms[i-1][j].R - atoms[i-1][j+1].R +
				atoms[i+1][j-1].R + 2*atoms[i+1][j].R + atoms[i+1][j+1].R
			fields[0].GradX[i][j] = gxR / 8.0
			fields[0].GradY[i][j] = gyR / 8.0
			fields[0].Mag[i][j] = math.Sqrt(gxR*gxR+gyR*gyR) / 8.0

			// G channel
			gxG := -atoms[i-1][j-1].G + atoms[i-1][j+1].G +
				-2*atoms[i][j-1].G + 2*atoms[i][j+1].G +
				-atoms[i+1][j-1].G + atoms[i+1][j+1].G
			gyG := -atoms[i-1][j-1].G - 2*atoms[i-1][j].G - atoms[i-1][j+1].G +
				atoms[i+1][j-1].G + 2*atoms[i+1][j].G + atoms[i+1][j+1].G
			fields[1].GradX[i][j] = gxG / 8.0
			fields[1].GradY[i][j] = gyG / 8.0
			fields[1].Mag[i][j] = math.Sqrt(gxG*gxG+gyG*gyG) / 8.0

			// B channel
			gxB := -atoms[i-1][j-1].B + atoms[i-1][j+1].B +
				-2*atoms[i][j-1].B + 2*atoms[i][j+1].B +
				-atoms[i+1][j-1].B + atoms[i+1][j+1].B
			gyB := -atoms[i-1][j-1].B - 2*atoms[i-1][j].B - atoms[i-1][j+1].B +
				atoms[i+1][j-1].B + 2*atoms[i+1][j].B + atoms[i+1][j+1].B
			fields[2].GradX[i][j] = gxB / 8.0
			fields[2].GradY[i][j] = gyB / 8.0
			fields[2].Mag[i][j] = math.Sqrt(gxB*gxB+gyB*gyB) / 8.0
		}
	}

	return fields
}

// ============================================================================
// ENERGY TERM: E_SHARPEN (GRADIENT MATCHING)
// ============================================================================

// CalculateSharpenEnergy computes E_sharpen = Σ ||∇I_recon - ∇I_target||²
// This forces the reconstruction to match the gradient structure of the target
func CalculateSharpenEnergy(recon [][]PixelAtomV2, target *GradientField) float64 {
	reconGrad := ComputeGradientField(recon)
	h := len(recon)
	w := len(recon[0])

	var totalError float64
	count := 0

	for i := 1; i < h-1; i++ {
		for j := 1; j < w-1; j++ {
			// Compare gradient magnitudes
			diffX := reconGrad.GradX[i][j] - target.GradX[i][j]
			diffY := reconGrad.GradY[i][j] - target.GradY[i][j]
			totalError += diffX*diffX + diffY*diffY
			count++
		}
	}

	if count > 0 {
		return totalError / float64(count)
	}
	return 0
}

// CalculateSharpenEnergyRGB computes gradient matching for RGB channels
func CalculateSharpenEnergyRGB(recon [][]PixelAtomV2, targetRGB [3]*GradientField) float64 {
	reconRGB := ComputeGradientFieldRGB(recon)
	h := len(recon)
	w := len(recon[0])

	var totalError float64
	count := 0

	for i := 1; i < h-1; i++ {
		for j := 1; j < w-1; j++ {
			// Sum error across all 3 channels
			for c := 0; c < 3; c++ {
				diffX := reconRGB[c].GradX[i][j] - targetRGB[c].GradX[i][j]
				diffY := reconRGB[c].GradY[i][j] - targetRGB[c].GradY[i][j]
				totalError += diffX*diffX + diffY*diffY
			}
			count++
		}
	}

	if count > 0 {
		return totalError / float64(count*3)
	}
	return 0
}

// ============================================================================
// EDGE ENHANCEMENT ENERGY
// ============================================================================

// CalculateEdgeEnhancementEnergy computes E_sharpness = -λ Σ ||∇I||²
// Negative because we SUBTRACT this (gradient descent MINIMIZES -E, i.e., MAXIMIZES E)
// The system is "attracted" to regions with high gradients (sharp edges)
func CalculateEdgeEnhancementEnergy(atoms [][]PixelAtomV2, lambda float64) float64 {
	h := len(atoms)
	w := len(atoms[0])
	if h < 3 || w < 3 {
		return 0
	}

	var totalEnergy float64
	count := 0

	// Simple Sobel for each pixel
	for i := 1; i < h-1; i++ {
		for j := 1; j < w-1; j++ {
			// Sobel gradient magnitude using intensity
			gx := atoms[i+1][j].Intensity - atoms[i-1][j].Intensity
			gy := atoms[i][j+1].Intensity - atoms[i][j-1].Intensity

			gradMag := gx*gx + gy*gy // ||∇I||²

			// Negative: rewards high gradients (system minimizes -E, maximizes E)
			totalEnergy -= lambda * gradMag
			count++
		}
	}

	if count > 0 {
		return totalEnergy / float64(count)
	}
	return 0
}

// ComputeLocalEdgeStrength returns the gradient magnitude at position (i,j)
// Used for local energy gradient calculations
func ComputeLocalEdgeStrength(atoms [][]PixelAtomV2, i, j int) float64 {
	h := len(atoms)
	w := len(atoms[0])
	if i < 1 || i >= h-1 || j < 1 || j >= w-1 {
		return 0
	}

	// Simplified Sobel
	gx := atoms[i+1][j].Intensity - atoms[i-1][j].Intensity
	gy := atoms[i][j+1].Intensity - atoms[i][j-1].Intensity

	return gx*gx + gy*gy // ||∇I||²
}

// ComputeLocalEdgeGradient returns how much each RGB channel should move
// to increase the local gradient magnitude
// Result: [R, G, B] gradient contributions
func ComputeLocalEdgeGradient(atoms [][]PixelAtomV2, i, j int, lambda float64) [3]float64 {
	h := len(atoms)
	w := len(atoms[0])
	var grad [3]float64

	if i < 1 || i >= h-1 || j < 1 || j >= w-1 {
		return grad
	}

	// Sobel gradients
	gx := atoms[i+1][j].Intensity - atoms[i-1][j].Intensity
	gy := atoms[i][j+1].Intensity - atoms[i][j-1].Intensity

	// How to increase gradient magnitude?
	// d/dI(gx²+gy²) = 2·gx·(∂gx/∂I) + 2·gy·(∂gy/∂I)

	// For pixel (i,j):
	// ∂gx/∂I_{i,j} = 0 (not directly in Sobel formula)
	// ∂gy/∂I_{i,j} = 0
	// But for neighbors:
	// ∂gx/∂I_{i+1,j} = 1,  ∂gx/∂I_{i-1,j} = -1
	// ∂gy/∂I_{i,j+1} = 1,  ∂gy/∂I_{i,j-1} = -1

	// Local contribution: move in direction that increases gradient
	// If gx > 0 and gy > 0, moving in positive direction helps
	// Coefficient for R,G,B channels (all equal for intensity-based)
	coeff := lambda * (gx + gy) / 100.0 // Normalized

	grad[0] = coeff
	grad[1] = coeff
	grad[2] = coeff

	return grad
}

// ============================================================================
// BLUR DETECTION & ENERGY CALCULATION
// ============================================================================

// BlurLevel contains blur metrics for a patch
type BlurLevel struct {
	GradientStdDev    float64 // σ = std_dev(∇I)
	LaplacianVariance float64 // ∇² variance (texture detail)
	EdgeDensity       float64 // fraction of pixels with strong gradients
	IsBlurry          bool    // σ < threshold
	DetailLevel       float64 // [0,1] - 0=blurry, 1=sharp
	RequiredAtomBoost float64 // multiplier for atom allocation
	EstimatedSigma    float64 // Estimated Gaussian blur σ (for deconvolution)
}

// CalculateBlurLevel detects blur using gradient standard deviation
func CalculateBlurLevel(atoms [][]PixelAtomV2) *BlurLevel {
	h := len(atoms)
	w := len(atoms[0])
	if h < 2 || w < 2 {
		return &BlurLevel{GradientStdDev: 0.1, IsBlurry: true, DetailLevel: 0, RequiredAtomBoost: 2.0}
	}

	// Calculate gradient magnitude at each pixel
	gradients := make([]float64, 0)
	var edgeCount int

	for i := 0; i < h; i++ {
		for j := 0; j < w; j++ {
			// Sobel-like gradient (simplified)
			var gx, gy float64

			// Horizontal gradient
			if j > 0 && j < w-1 {
				left := atoms[i][j-1].Intensity
				right := atoms[i][j+1].Intensity
				gx = (right - left) / 2.0
			}

			// Vertical gradient
			if i > 0 && i < h-1 {
				top := atoms[i-1][j].Intensity
				bottom := atoms[i+1][j].Intensity
				gy = (bottom - top) / 2.0
			}

			grad := math.Sqrt(gx*gx + gy*gy)
			gradients = append(gradients, grad)

			// Count edges (strong gradients)
			if grad > 0.05 {
				edgeCount++
			}
		}
	}

	// Calculate std_dev of gradients
	mean := 0.0
	for _, g := range gradients {
		mean += g
	}
	mean /= float64(len(gradients))

	variance := 0.0
	for _, g := range gradients {
		diff := g - mean
		variance += diff * diff
	}
	variance /= float64(len(gradients))
	stdDev := math.Sqrt(variance)

	// Calculate Laplacian variance (texture detail)
	laplacianVar := 0.0
	laplacianCount := 0
	for i := 1; i < h-1; i++ {
		for j := 1; j < w-1; j++ {
			// 2nd derivative (Laplacian approximation)
			center := atoms[i][j].Intensity
			lap := 4*center -
				atoms[i-1][j].Intensity - atoms[i+1][j].Intensity -
				atoms[i][j-1].Intensity - atoms[i][j+1].Intensity
			laplacianVar += lap * lap
			laplacianCount++
		}
	}
	if laplacianCount > 0 {
		laplacianVar /= float64(laplacianCount)
	}

	// Detect blur threshold
	blurThreshold := 0.08
	isBlurry := stdDev < blurThreshold
	detailLevel := math.Min(1.0, stdDev/blurThreshold)

	// Boost atoms for blurry zones
	atomBoost := 1.0
	if isBlurry {
		// Inverse allocation: more atoms where blur is stronger
		atomBoost = 1.0 + (1.0-detailLevel)*2.0 // 1-3x boost
	}

	return &BlurLevel{
		GradientStdDev:    stdDev,
		LaplacianVariance: laplacianVar,
		EdgeDensity:       float64(edgeCount) / float64(len(gradients)),
		IsBlurry:          isBlurry,
		DetailLevel:       detailLevel,
		RequiredAtomBoost: atomBoost,
	}
}

// ============================================================================
// ENERGY TERMS FOR DEBLURRING
// ============================================================================

// CalculateStructureEnergy measures detail presence (variance of gradients)
func CalculateStructureEnergy(atoms [][]PixelAtomV2) float64 {
	h := len(atoms)
	w := len(atoms[0])
	if h < 2 || w < 2 {
		return 0
	}

	var totalEnergy float64

	for i := 0; i < h; i++ {
		for j := 0; j < w; j++ {
			var gx, gy float64

			if j > 0 && j < w-1 {
				gx = (atoms[i][j+1].Intensity - atoms[i][j-1].Intensity) / 2.0
			}
			if i > 0 && i < h-1 {
				gy = (atoms[i+1][j].Intensity - atoms[i-1][j].Intensity) / 2.0
			}

			grad := math.Sqrt(gx*gx + gy*gy)
			totalEnergy += grad
		}
	}

	return totalEnergy / float64(h*w)
}

// CalculateInteractionEnergy measures coherence with neighbors
func CalculateInteractionEnergy(atoms [][]PixelAtomV2, neighbor [][]PixelAtomV2) float64 {
	h := len(atoms)
	w := len(atoms[0])
	nh := len(neighbor)
	nw := len(neighbor[0])

	if h != nh || w != nw {
		return 0
	}

	var totalError float64

	for i := 0; i < h; i++ {
		for j := 0; j < w; j++ {
			colorDiff := math.Sqrt(
				math.Pow(atoms[i][j].R-neighbor[i][j].R, 2) +
					math.Pow(atoms[i][j].G-neighbor[i][j].G, 2) +
					math.Pow(atoms[i][j].B-neighbor[i][j].B, 2),
			)
			totalError += colorDiff
		}
	}

	return totalError / float64(h*w)
}

// CalculateSharpnessEnergy encourages high-frequency content (gradient variance)
func CalculateSharpnessEnergy(atoms [][]PixelAtomV2) float64 {
	blur := CalculateBlurLevel(atoms)
	// Low σ = high energy penalty
	return math.Max(0, 0.1-blur.GradientStdDev)
}

// CalculateTextureEnergy measures Laplacian (2nd derivative details)
func CalculateTextureEnergy(atoms [][]PixelAtomV2) float64 {
	h := len(atoms)
	w := len(atoms[0])
	if h < 3 || w < 3 {
		return 0
	}

	var totalLaplacian float64
	count := 0

	for i := 1; i < h-1; i++ {
		for j := 1; j < w-1; j++ {
			center := atoms[i][j].Intensity
			lap := 4*center -
				atoms[i-1][j].Intensity - atoms[i+1][j].Intensity -
				atoms[i][j-1].Intensity - atoms[i][j+1].Intensity
			totalLaplacian += math.Abs(lap)
			count++
		}
	}

	if count > 0 {
		return totalLaplacian / float64(count)
	}
	return 0
}

// ============================================================================
// MULTI-PHASE DEBLUR RELAXATION
// ============================================================================

// DeblurPhase represents one phase of the multi-phase deblur process
type DeblurPhase struct {
	Name             string  // "Coarse", "Medium", "Fine"
	PatchSize        int     // 16, 8, or 4
	Alpha            float64 // weight for E_structure
	Beta             float64 // weight for E_interaction
	Lambda           float64 // weight for E_sharp
	Mu               float64 // weight for E_texture
	Iterations       int     // iterations for this phase
	MinimumBlurRatio float64 // only process patches with σ < this ratio
}

// NewDeblurPhases creates the three-phase pipeline
func NewDeblurPhases() []*DeblurPhase {
	return []*DeblurPhase{
		{
			Name:             "Coarse",
			PatchSize:        16,
			Alpha:            0.6, // Strong structure focus
			Beta:             0.4, // Some interaction
			Lambda:           0.1, // Minimal sharpness
			Mu:               0.0, // No texture yet
			Iterations:       30,
			MinimumBlurRatio: 1.0, // Process all patches
		},
		{
			Name:             "Medium",
			PatchSize:        8,
			Alpha:            0.4, // Balanced structure
			Beta:             0.3, // Moderate interaction
			Lambda:           0.3, // Increase sharpness focus
			Mu:               0.0, // Still no texture
			Iterations:       40,
			MinimumBlurRatio: 0.9, // Focus on slightly blurry areas
		},
		{
			Name:             "Fine",
			PatchSize:        4,
			Alpha:            0.3, // Light structure
			Beta:             0.2, // Minimal interaction
			Lambda:           0.3, // Maintain sharpness
			Mu:               0.2, // Add texture refinement
			Iterations:       50,
			MinimumBlurRatio: 0.7, // Focus on very blurry areas
		},
	}
}

// ============================================================================
// BLUR MASK FOR SELECTIVE RELAXATION
// ============================================================================

// BlurMask identifies which patches need deblurring
type BlurMask struct {
	Levels       [][]*BlurLevel // Blur metrics per patch
	NeedsRelax   [][]bool       // Which patches should be relaxed
	BlurryCount  int            // Count of blurry patches
	TotalPatches int
}

// NewBlurMask creates and analyzes blur in patch grid
func NewBlurMask(patches [][]OptimizedPatch) *BlurMask {
	h := len(patches)
	w := len(patches[0])

	mask := &BlurMask{
		Levels:       make([][]*BlurLevel, h),
		NeedsRelax:   make([][]bool, h),
		TotalPatches: h * w,
	}

	for i := 0; i < h; i++ {
		mask.Levels[i] = make([]*BlurLevel, w)
		mask.NeedsRelax[i] = make([]bool, w)

		for j := 0; j < w; j++ {
			blur := CalculateBlurLevel(patches[i][j].Atoms)
			mask.Levels[i][j] = blur
			mask.NeedsRelax[i][j] = blur.IsBlurry

			if blur.IsBlurry {
				mask.BlurryCount++
			}
		}
	}

	return mask
}

// AdaptiveAtomAllocationForDeblur allocates atoms based on blur
func (mask *BlurMask) AdaptiveAtomAllocationForDeblur(baseAtoms int, scaleFactor float64) int {
	h := len(mask.Levels)
	w := len(mask.Levels[0])

	totalAtoms := 0
	for i := 0; i < h; i++ {
		for j := 0; j < w; j++ {
			blur := mask.Levels[i][j]
			// n_i,j = ceil(k / (σ + ε))
			atomFactor := scaleFactor / (blur.GradientStdDev + 0.01)
			atomFactor = math.Min(atomFactor, 4.0) // Cap at 4x
			allocatedAtoms := int(math.Ceil(float64(baseAtoms) * atomFactor))
			totalAtoms += allocatedAtoms
		}
	}

	return totalAtoms
}

// PrintBlurAnalysis displays blur detection results
func (mask *BlurMask) PrintBlurAnalysis() {
	fmt.Printf("\n[BLUR ANALYSIS]\n")
	fmt.Printf("  • Blurry patches detected: %d / %d\n", mask.BlurryCount, mask.TotalPatches)
	fmt.Printf("  • Blur ratio: %.1f%%\n", float64(mask.BlurryCount)*100/float64(mask.TotalPatches))

	// Find most/least blurry
	var minBlur, maxBlur float64 = 1e10, 0
	for i := 0; i < len(mask.Levels); i++ {
		for j := 0; j < len(mask.Levels[0]); j++ {
			sigma := mask.Levels[i][j].GradientStdDev
			if sigma < minBlur {
				minBlur = sigma
			}
			if sigma > maxBlur {
				maxBlur = sigma
			}
		}
	}

	fmt.Printf("  • σ range: %.4f - %.4f\n", minBlur, maxBlur)
	fmt.Printf("  • Focus zones: selective relaxation enabled\n")
}

// ============================================================================
// DEBLUR EXECUTOR
// ============================================================================

// ExecuteDeblurPipeline runs the multi-phase deblurring
func (grid *OptimizedPatchGrid) ExecuteDeblurPipeline(totalIterations int) {
	fmt.Printf("\n[DEBLURRING PIPELINE]\n")

	// Phase 1: Blur detection
	blurMask := NewBlurMask(grid.Patches)
	blurMask.PrintBlurAnalysis()

	// Phase 2: Multi-phase relaxation
	phases := NewDeblurPhases()
	iterPerPhase := totalIterations / len(phases)

	for phaseIdx, phase := range phases {
		fmt.Printf("\n[PHASE %d: %s]\n", phaseIdx+1, phase.Name)
		fmt.Printf("  • Patch size: %d×%d\n", phase.PatchSize, phase.PatchSize)
		fmt.Printf("  • Energy weights: α=%.2f, β=%.2f, λ=%.2f, μ=%.2f\n",
			phase.Alpha, phase.Beta, phase.Lambda, phase.Mu)
		fmt.Printf("  • Iterations: %d\n", iterPerPhase)

		// Run relaxation for this phase
		grid.RelaxPhaseAdaptive(phase, blurMask, iterPerPhase)
	}

	fmt.Printf("\n[DEBLURRING COMPLETE]\n")
	fmt.Printf("  ✓ Blur-aware optimization finished\n")
}

// RelaxPhaseAdaptive runs one phase with adaptive weights and deconvolution
func (grid *OptimizedPatchGrid) RelaxPhaseAdaptive(phase *DeblurPhase, mask *BlurMask, iterations int) {
	gridH := len(grid.Patches)
	gridW := len(grid.Patches[0])
	params := NewDefaultDeconvolutionParams()

	for iter := 0; iter < iterations; iter++ {
		for i := 0; i < gridH; i++ {
			for j := 0; j < gridW; j++ {
				// Only relax if patch is blurry
				if !mask.NeedsRelax[i][j] {
					continue
				}

				patch := &grid.Patches[i][j]
				blurLevel := mask.Levels[i][j]

				// Apply phase-specific energy minimization
				patch.RelaxWithEnergyTerms(
					phase.Alpha, phase.Beta, phase.Lambda, phase.Mu,
					grid.LearningRate, blurLevel,
				)

				// Apply Gaussian deconvolution every 10 iterations (moderate sharpening)
				if iter > 0 && iter%10 == 0 {
					ApplyGaussianDeconvolution(patch.Atoms, params.GaussianSigma, params.SharpeningStrength)
				}
			}
		}
	}

	// 🔥 FINAL DECONVOLUTION: Richardson-Lucy for true detail hallucination
	fmt.Printf("  • Applying Richardson-Lucy deconvolution (k=%.2f, %d iterations)...\n", params.GradientAmplification, 5)
	for i := 0; i < gridH; i++ {
		for j := 0; j < gridW; j++ {
			if mask.NeedsRelax[i][j] {
				// Richardson-Lucy: more accurate deconvolution
				ApplyRichardsonLucyDeconvolution(grid.Patches[i][j].Atoms, params.GaussianSigma, 5)

				// Then amplify high-frequencies for extra sharpness
				ApplyGaussianDeconvolution(grid.Patches[i][j].Atoms, params.GaussianSigma, params.SharpeningStrength*1.5)
			}
		}
	}
}

// RelaxWithEnergyTerms applies multi-term energy minimization with gradient matching
// RelaxWithEnergyTerms applies multi-term energy minimization with gradient amplification
func (patch *OptimizedPatch) RelaxWithEnergyTerms(alpha, beta, lambda, mu, learningRate float64, blur *BlurLevel) {
	// 🔥 ADAPTIVE: k increases with blur severity
	params := NewAdaptiveDeconvolutionParams(blur.EstimatedSigma)
	k := params.GradientAmplification          // k ∈ [1.5, 3.0] based on blur level
	lambdaEdge := params.EdgeEnhancementLambda // λ for edge enhancement

	// Compute energy components
	structureEnergy := CalculateStructureEnergy(patch.Atoms)
	sharpnessEnergy := CalculateSharpnessEnergy(patch.Atoms)
	textureEnergy := CalculateTextureEnergy(patch.Atoms)

	// 🌟 NEW: Edge enhancement energy (rewards high gradients)
	// E_edge = -λ Σ ||∇I||² (negative, so system minimizes -E = maximizes E)
	edgeEnergy := CalculateEdgeEnhancementEnergy(patch.Atoms, lambdaEdge)

	// KEY: E_sharpen = Σ||∇I_recon - k·∇I_blur||² where k > 1
	// This AMPLIFIES gradients to hallucinate missing details
	sharpenEnergy := 0.0
	if patch.OriginalGradient[0] != nil {
		sharpenEnergy = CalculateSharpenEnergyAmplified(patch.Atoms, patch.OriginalGradient, k)
	}
	_ = sharpenEnergy // Used for monitoring

	// Total energy = α*E_structure + β*E_constraint + γ*E_interaction + λ*E_edge
	// Weighted: E_total = α*E_struct + λ*E_sharp + μ*E_texture + 0.5*E_edge
	totalWeightedEnergy := alpha*structureEnergy + lambda*sharpnessEnergy + mu*textureEnergy + 0.5*edgeEnergy

	// Gradient: ∇E with amplification boost
	gradientBoost := blur.RequiredAtomBoost * learningRate

	for i := range patch.Atoms {
		for j := range patch.Atoms[i] {
			// Compute local gradient contribution
			baseGradient := totalWeightedEnergy

			// Add sharpen gradient with k-amplification (DETAIL HALLUCINATION)
			localSharpenGrad := computeLocalSharpenGradientAmplified(patch.Atoms, patch.OriginalGradient, i, j, k)

			// 🌟 NEW: Add edge enhancement gradient (pushes towards high-gradient regions)
			edgeGrad := ComputeLocalEdgeGradient(patch.Atoms, i, j, lambdaEdge)

			// Update RGB with gradient descent
			// SUBTRACT gradients to MINIMIZE total energy
			// This forces I_recon to have STRONGER gradients than original AND rewards edges
			patch.Atoms[i][j].R -= gradientBoost * (baseGradient - 0.5*localSharpenGrad[0] - 0.3*edgeGrad[0])
			patch.Atoms[i][j].G -= gradientBoost * (baseGradient - 0.5*localSharpenGrad[1] - 0.3*edgeGrad[1])
			patch.Atoms[i][j].B -= gradientBoost * (baseGradient - 0.5*localSharpenGrad[2] - 0.3*edgeGrad[2])

			// Clamp to [0, 1]
			patch.Atoms[i][j].R = math.Max(0, math.Min(1, patch.Atoms[i][j].R))
			patch.Atoms[i][j].G = math.Max(0, math.Min(1, patch.Atoms[i][j].G))
			patch.Atoms[i][j].B = math.Max(0, math.Min(1, patch.Atoms[i][j].B))

			patch.Atoms[i][j].Intensity = (patch.Atoms[i][j].R + patch.Atoms[i][j].G + patch.Atoms[i][j].B) / 3.0
		}
	}
}

// CalculateSharpenEnergyAmplified computes E_sharpen with k-amplification
// E = Σ ||∇I_recon - k·∇I_blur||² where k > 1 encourages stronger gradients
func CalculateSharpenEnergyAmplified(recon [][]PixelAtomV2, targetRGB [3]*GradientField, k float64) float64 {
	reconRGB := ComputeGradientFieldRGB(recon)
	h := len(recon)
	w := len(recon[0])

	var totalError float64
	count := 0

	for i := 1; i < h-1; i++ {
		for j := 1; j < w-1; j++ {
			// Sum error across all 3 channels
			for c := 0; c < 3; c++ {
				// Compare ∇I_recon with k·∇I_blur (amplified target)
				targetX := k * targetRGB[c].GradX[i][j]
				targetY := k * targetRGB[c].GradY[i][j]

				diffX := reconRGB[c].GradX[i][j] - targetX
				diffY := reconRGB[c].GradY[i][j] - targetY
				totalError += diffX*diffX + diffY*diffY
			}
			count++
		}
	}

	if count > 0 {
		return totalError / float64(count*3)
	}
	return 0
}

// computeLocalSharpenGradientAmplified calculates gradient with k-amplification
func computeLocalSharpenGradientAmplified(atoms [][]PixelAtomV2, targetGrad [3]*GradientField, i, j int, k float64) [3]float64 {
	h := len(atoms)
	w := len(atoms[0])

	// If at boundary or no target gradient, return zero
	if i <= 0 || i >= h-1 || j <= 0 || j >= w-1 || targetGrad[0] == nil {
		return [3]float64{0, 0, 0}
	}

	// Compute current gradients at (i,j)
	currentGrad := ComputeGradientFieldRGB(atoms)

	grad := [3]float64{0, 0, 0}

	for c := 0; c < 3; c++ {
		if currentGrad[c] != nil && targetGrad[c] != nil {
			// Compare with k-amplified target gradients
			targetX := k * targetGrad[c].GradX[i][j]
			targetY := k * targetGrad[c].GradY[i][j]

			diffX := currentGrad[c].GradX[i][j] - targetX
			diffY := currentGrad[c].GradY[i][j] - targetY
			grad[c] = 2.0 * (diffX + diffY) // Gradient of squared error
		}
	}

	return grad
}
