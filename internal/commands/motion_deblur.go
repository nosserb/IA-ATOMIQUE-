package commands

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"math"
	"os"
)

// MotionDeblurKernel représente un kernel de motion blur estim
type MotionDeblurKernel struct {
	Length int     // Longueur du kernel (pixels du motion blur)
	Angle  float64 // Angle de direction (0-360 deg)
	Data   []float64
	Width  int
	Height int
}

// EstimateMotionBlurKernel estime le kernel de motion blur par autocorrélation
// Analyse la texture de l'image pour détecter la direction et longueur du blur
func EstimateMotionBlurKernel(imgPath string) (*MotionDeblurKernel, error) {
	file, err := os.Open(imgPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var img image.Image
	var cfg image.Config

	ext := imgPath[len(imgPath)-4:]
	if ext == ".jpg" || ext == "jpeg" {
		cfg, err = jpeg.DecodeConfig(file)
		if err != nil {
			return nil, err
		}
		file.Seek(0, 0)
		img, err = jpeg.Decode(file)
		if err != nil {
			return nil, err
		}
	} else {
		cfg, err = png.DecodeConfig(file)
		if err != nil {
			return nil, err
		}
		file.Seek(0, 0)
		img, err = png.Decode(file)
		if err != nil {
			return nil, err
		}
	}

	// Convertir en grayscale pour l'analyse
	bounds := img.Bounds()
	gray := image.NewGray(bounds)
	draw.Draw(gray, bounds, img, image.Point{}, draw.Src)

	// Détecter motion blur par analyse de gradient
	// Motion blur = gradients faibles perpendiculaires à la direction
	kernel := &MotionDeblurKernel{
		Width:  cfg.Width,
		Height: cfg.Height,
	}

	// Analyser les bords pour estimer direction du motion blur
	// Utiliser Sobel pour détecter direction
	dx, dy := computeGradients(gray)

	// Estimer angle et longueur
	kernel.Length, kernel.Angle = analyzeBlurDirection(dx, dy, cfg.Width, cfg.Height)

	// Créer kernel linéaire 1D
	kernel.Data = createMotionBlurKernel1D(kernel.Length)

	fmt.Printf("   📊 Motion Blur Analysis:\n")
	fmt.Printf("      • Detected length: %d pixels\n", kernel.Length)
	fmt.Printf("      • Detected angle: %.1f°\n", kernel.Angle)

	return kernel, nil
}

// computeGradients calcule les gradients Sobel
func computeGradients(gray *image.Gray) ([][]float64, [][]float64) {
	bounds := gray.Bounds()
	w, h := bounds.Dx(), bounds.Dy()

	dx := make([][]float64, h)
	dy := make([][]float64, h)
	for i := range dx {
		dx[i] = make([]float64, w)
		dy[i] = make([]float64, w)
	}

	for y := 1; y < h-1; y++ {
		for x := 1; x < w-1; x++ {
			// Sobel X
			gx := float64(gray.GrayAt(x-1, y-1).Y) * -1
			gx += float64(gray.GrayAt(x-1, y).Y) * -2
			gx += float64(gray.GrayAt(x-1, y+1).Y) * -1
			gx += float64(gray.GrayAt(x+1, y-1).Y) * 1
			gx += float64(gray.GrayAt(x+1, y).Y) * 2
			gx += float64(gray.GrayAt(x+1, y+1).Y) * 1
			dx[y][x] = gx / 8.0

			// Sobel Y
			gy := float64(gray.GrayAt(x-1, y-1).Y) * -1
			gy += float64(gray.GrayAt(x, y-1).Y) * -2
			gy += float64(gray.GrayAt(x+1, y-1).Y) * -1
			gy += float64(gray.GrayAt(x-1, y+1).Y) * 1
			gy += float64(gray.GrayAt(x, y+1).Y) * 2
			gy += float64(gray.GrayAt(x+1, y+1).Y) * 1
			dy[y][x] = gy / 8.0
		}
	}

	return dx, dy
}

// analyzeBlurDirection estime la direction et longueur du motion blur
func analyzeBlurDirection(dx, dy [][]float64, width, height int) (int, float64) {
	// Calculer magnitude et angle des gradients
	var angleHistogram [360]int

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if y < len(dy) && x < len(dy[0]) && y < len(dx) && x < len(dx[0]) {
				gx, gy := dx[y][x], dy[y][x]
				mag := math.Sqrt(gx*gx + gy*gy)

				// Seuil: ignorer les gradients faibles
				if mag > 10 {
					// Angle perpendiculaire au gradient = direction du blur
					angle := math.Atan2(-gx, gy) * 180 / math.Pi
					if angle < 0 {
						angle += 360
					}
					angleHistogram[int(angle)]++
				}
			}
		}
	}

	// Trouver angle dominant
	maxCount := 0
	dominantAngle := 0
	for a := 0; a < 360; a++ {
		if angleHistogram[a] > maxCount {
			maxCount = angleHistogram[a]
			dominantAngle = a
		}
	}

	// Estimer longueur du blur (par défaut 10-50 pixels pour motion blur)
	// Ajuster selon la force du flou
	blurLength := 25
	if maxCount < 1000 {
		blurLength = 15
	} else if maxCount > 5000 {
		blurLength = 40
	}

	return blurLength, float64(dominantAngle)
}

// createMotionBlurKernel1D crée un kernel linéaire normalisé
func createMotionBlurKernel1D(length int) []float64 {
	kernel := make([]float64, length)
	for i := 0; i < length; i++ {
		kernel[i] = 1.0 / float64(length)
	}
	return kernel
}

// LucyRichardsonDeconvolve applique unsharp masking (simple et efficace)
func LucyRichardsonDeconvolve(imgPath string, outputPath string, iterations int) error {
	file, err := os.Open(imgPath)
	if err != nil {
		return err
	}
	defer file.Close()

	var img image.Image

	ext := imgPath[len(imgPath)-4:]
	if ext == ".jpg" || ext == "jpeg" {
		_, err = jpeg.DecodeConfig(file)
		file.Seek(0, 0)
		img, err = jpeg.Decode(file)
	} else {
		_, err = png.DecodeConfig(file)
		file.Seek(0, 0)
		img, err = png.Decode(file)
	}
	if err != nil {
		return err
	}

	bounds := img.Bounds()
	w, h := bounds.Dx(), bounds.Dy()

	// Simple unsharp masking (rehausse contours sans artefacts)
	fmt.Printf("   📊 Motion Blur Analysis: Unsharp masking for edge enhancement\n")
	fmt.Printf("   ⚡ Processing at original resolution...\n")

	imgRGB := imageToRGB(img, w, h)

	fmt.Printf("   🔄 Unsharp masking (%d iterations)...\n", iterations)

	// Itérations d'unsharp masking
	for iter := 0; iter < iterations; iter++ {
		imgRGB = unsharpMaskPass(imgRGB, w, h, 0.4) // 40% enhancement per pass (increased)

		if (iter+1)%3 == 0 {
			fmt.Printf("      [%d/%d] Sharpened\n", iter+1, iterations)
		}
	}

	// Post-processing: bilateral denoise
	fmt.Printf("   🎨 Post-processing (bilateral denoise)...\n")
	imgRGB = bilateralDenoiseAndNormalize(imgRGB, w, h)

	// Convertir back to image
	deblurredImg := rgbToImage(imgRGB, w, h)

	outFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	if len(outputPath) > 4 && outputPath[len(outputPath)-3:] == "png" {
		err = png.Encode(outFile, deblurredImg)
	} else {
		err = jpeg.Encode(outFile, deblurredImg, &jpeg.Options{Quality: 95})
	}

	fmt.Printf("   ✅ Sharpening complete\n")
	return err
}

// unsharpMaskPass applique une passe d'unsharp masking
func unsharpMaskPass(img [][][3]float64, w, h int, strength float64) [][][3]float64 {
	result := make([][][3]float64, h)
	for i := range result {
		result[i] = make([][3]float64, w)
	}

	// Gaussian blur 3x3 (rapide)
	blurred := make([][][3]float64, h)
	for i := range blurred {
		blurred[i] = make([][3]float64, w)
	}

	weights := []float64{
		0.0625, 0.1250, 0.0625,
		0.1250, 0.2500, 0.1250,
		0.0625, 0.1250, 0.0625,
	}

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			for c := 0; c < 3; c++ {
				sum := 0.0
				idx := 0
				for dy := -1; dy <= 1; dy++ {
					for dx := -1; dx <= 1; dx++ {
						iy := y + dy
						ix := x + dx
						if iy >= 0 && iy < h && ix >= 0 && ix < w {
							sum += img[iy][ix][c] * weights[idx]
						}
						idx++
					}
				}
				blurred[y][x][c] = sum
			}
		}
	}

	// Unsharp masking: original + (original - blurred) * strength
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			for c := 0; c < 3; c++ {
				delta := (img[y][x][c] - blurred[y][x][c]) * strength
				result[y][x][c] = math.Min(1.0, math.Max(0.0, img[y][x][c]+delta))
			}
		}
	}

	return result
}

// medianFilterQuick applique un median filter 3x3 rapide
func medianFilterQuick(img [][][3]float64, w, h int) [][][3]float64 {
	result := make([][][3]float64, h)
	for i := range result {
		result[i] = make([][3]float64, w)
	}

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			for c := 0; c < 3; c++ {
				if y == 0 || y == h-1 || x == 0 || x == w-1 {
					// Bords: copier
					result[y][x][c] = img[y][x][c]
				} else {
					// Médiane 3x3
					var vals [9]float64
					idx := 0
					for dy := -1; dy <= 1; dy++ {
						for dx := -1; dx <= 1; dx++ {
							vals[idx] = img[y+dy][x+dx][c]
							idx++
						}
					}

					// Trier (simple bubble sort pour 9 éléments)
					for i := 0; i < 9; i++ {
						for j := i + 1; j < 9; j++ {
							if vals[i] > vals[j] {
								vals[i], vals[j] = vals[j], vals[i]
							}
						}
					}
					result[y][x][c] = vals[4] // Médiane
				}
			}
		}
	}

	return result
}

// bilateralDenoiseAndNormalize applique bilateral filter + normalize les couleurs
func bilateralDenoiseAndNormalize(img [][][3]float64, w, h int) [][][3]float64 {
	// Clamp les valeurs
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			for c := 0; c < 3; c++ {
				if img[y][x][c] < 0 {
					img[y][x][c] = 0
				} else if img[y][x][c] > 1.0 {
					img[y][x][c] = 1.0
				}
			}
		}
	}

	result := make([][][3]float64, h)
	for i := range result {
		result[i] = make([][3]float64, w)
	}

	sigmaColor := 0.12 // Réduit pour mieux nettoyer les points
	sigmaSpace := 1.2

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			for c := 0; c < 3; c++ {
				sum := 0.0
				weightSum := 0.0
				centerVal := img[y][x][c]

				for dy := -2; dy <= 2; dy++ {
					for dx := -2; dx <= 2; dx++ {
						iy := y + dy
						ix := x + dx

						if iy >= 0 && iy < h && ix >= 0 && ix < w {
							neighborVal := img[iy][ix][c]

							spatialDist := math.Sqrt(float64(dy*dy+dx*dx)) / sigmaSpace
							spatialWeight := math.Exp(-spatialDist * spatialDist / 2.0)

							colorDist := math.Abs(neighborVal-centerVal) / sigmaColor
							colorWeight := math.Exp(-colorDist * colorDist / 2.0)

							weight := spatialWeight * colorWeight
							sum += neighborVal * weight
							weightSum += weight
						}
					}
				}

				if weightSum > 0 {
					result[y][x][c] = sum / weightSum
				} else {
					result[y][x][c] = img[y][x][c]
				}

				if result[y][x][c] < 0 {
					result[y][x][c] = 0
				} else if result[y][x][c] > 1.0 {
					result[y][x][c] = 1.0
				}
			}
		}
	}

	return result
}

// denoiseLucyRichardsonArtifacts applique un median filter + slight gaussian blur
// pour nettoyer les artefacts RGB du Lucy-Richardson
func denoiseLucyRichardsonArtifacts(img [][][3]float64, w, h int) [][][3]float64 {
	// Median filter (3x3) pour enlever les pixels étranges
	result := make([][][3]float64, h)
	for i := range result {
		result[i] = make([][3]float64, w)
	}

	for y := 1; y < h-1; y++ {
		for x := 1; x < w-1; x++ {
			for c := 0; c < 3; c++ {
				// Récupérer les 9 pixels (3x3)
				var vals [9]float64
				idx := 0
				for dy := -1; dy <= 1; dy++ {
					for dx := -1; dx <= 1; dx++ {
						vals[idx] = img[y+dy][x+dx][c]
						idx++
					}
				}

				// Trier et prendre la médiane
				for i := 0; i < 9; i++ {
					for j := i + 1; j < 9; j++ {
						if vals[i] > vals[j] {
							vals[i], vals[j] = vals[j], vals[i]
						}
					}
				}
				result[y][x][c] = vals[4] // Médiane (element central après tri)
			}
		}
	}

	// Copier les bords
	for x := 0; x < w; x++ {
		result[0][x] = img[0][x]
		result[h-1][x] = img[h-1][x]
	}
	for y := 0; y < h; y++ {
		result[y][0] = img[y][0]
		result[y][w-1] = img[y][w-1]
	}

	// Gaussian blur léger (sigma=0.8) pour lisser
	result2 := make([][][3]float64, h)
	for i := range result2 {
		result2[i] = make([][3]float64, w)
	}

	// Kernel gaussien 5x5 simplifié (approximation)
	for y := 2; y < h-2; y++ {
		for x := 2; x < w-2; x++ {
			for c := 0; c < 3; c++ {
				// Gaussian 5x5 weights (normalized)
				sum := 0.0
				weight := 0.0

				weights := []float64{
					0.0625, 0.1250, 0.1563, 0.1250, 0.0625,
					0.1250, 0.2500, 0.3125, 0.2500, 0.1250,
					0.1563, 0.3125, 0.3906, 0.3125, 0.1563,
					0.1250, 0.2500, 0.3125, 0.2500, 0.1250,
					0.0625, 0.1250, 0.1563, 0.1250, 0.0625,
				}

				idx := 0
				for dy := -2; dy <= 2; dy++ {
					for dx := -2; dx <= 2; dx++ {
						sum += result[y+dy][x+dx][c] * weights[idx]
						weight += weights[idx]
						idx++
					}
				}
				result2[y][x][c] = sum / weight
			}
		}
	}

	// Copier les bords
	for x := 0; x < w; x++ {
		result2[0][x] = result[0][x]
		result2[1][x] = result[1][x]
		result2[h-2][x] = result[h-2][x]
		result2[h-1][x] = result[h-1][x]
	}
	for y := 0; y < h; y++ {
		result2[y][0] = result[y][0]
		result2[y][1] = result[y][1]
		result2[y][w-2] = result[y][w-2]
		result2[y][w-1] = result[y][w-1]
	}

	return result2
}

// imageToRGB convertit image.Image en array RGB
func imageToRGB(img image.Image, w, h int) [][][3]float64 {
	rgb := make([][][3]float64, h)
	for y := 0; y < h; y++ {
		rgb[y] = make([][3]float64, w)
		for x := 0; x < w; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			rgb[y][x][0] = float64(r>>8) / 255.0
			rgb[y][x][1] = float64(g>>8) / 255.0
			rgb[y][x][2] = float64(b>>8) / 255.0
		}
	}
	return rgb
}

// rgbToImage convertit array RGB back en image.Image
func rgbToImage(rgb [][][3]float64, w, h int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			r := uint8(math.Min(1.0, math.Max(0.0, rgb[y][x][0])) * 255)
			g := uint8(math.Min(1.0, math.Max(0.0, rgb[y][x][1])) * 255)
			b := uint8(math.Min(1.0, math.Max(0.0, rgb[y][x][2])) * 255)
			img.SetRGBA(x, y, color.RGBA{r, g, b, 255})
		}
	}
	return img
}

// createRotated2DKernel crée un kernel 2D rotaté selon l'angle
func createRotated2DKernel(kernel1D []float64, angle float64, size int) [][]float64 {
	result := make([][]float64, size)
	for i := range result {
		result[i] = make([]float64, size)
	}

	// Convertir kernel 1D en 2D rotaté
	rad := angle * math.Pi / 180
	cosA := math.Cos(rad)
	sinA := math.Sin(rad)

	center := float64(size) / 2

	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			// Coordonnées relatives
			dy, dx := float64(i)-center, float64(j)-center

			// Rotation inverse
			projX := dx*cosA + dy*sinA

			// Projeter sur l'axe du blur
			idx := int(math.Abs(projX) + 0.5)
			if idx < len(kernel1D) {
				result[i][j] = kernel1D[idx]
			}
		}
	}

	// Normaliser
	sum := 0.0
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			sum += result[i][j]
		}
	}
	if sum > 0 {
		for i := 0; i < size; i++ {
			for j := 0; j < size; j++ {
				result[i][j] /= sum
			}
		}
	}

	return result
}

// lucyRichardsonIteration effectue une itération Lucy-Richardson OPTIMISÉE
func lucyRichardsonIteration(img [][][3]float64, kernel [][]float64, w, h int) [][][3]float64 {
	ksize := len(kernel)
	khalf := ksize / 2

	// Utiliser un kernel simplifié (1D au lieu de 2D) pour la vitesse
	// Appliquer séparément horizontalement et verticalement
	result := make([][][3]float64, h)
	for i := range result {
		result[i] = make([][3]float64, w)
	}

	// Copier l'image originale
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			result[y][x] = img[y][x]
		}
	}

	// Appliquer Lucy-Richardson simplifié: resampler + ratio
	// Pour chaque pixel, ajuster légèrement basé sur le contexte local
	strength := 0.5 // Réduit pour minimiser les artefacts

	for y := khalf; y < h-khalf; y++ {
		for x := khalf; x < w-khalf; x++ {
			for c := 0; c < 3; c++ {
				// Calculer moyenne locale
				localSum := 0.0
				for ky := 0; ky < ksize; ky++ {
					for kx := 0; kx < ksize; kx++ {
						iy := y + ky - khalf
						ix := x + kx - khalf
						localSum += img[iy][ix][c] * kernel[ky][kx]
					}
				}

				// Sharpening ratio
				if localSum > 1e-6 {
					ratio := img[y][x][c] / localSum
					delta := (img[y][x][c] - localSum) * strength
					result[y][x][c] = math.Min(1.0, math.Max(0.0, img[y][x][c]+delta*ratio))
				}
			}
		}
	}

	return result
}
