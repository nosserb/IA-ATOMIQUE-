package database

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"os"
	"strings"
)

// PatternMathematical représente un pattern par ses composantes mathématiques
// Pattern = f(x,y) = Σ αk · gk(x,y)
// Où gk sont les fonctions de base et αk les coefficients
type PatternMathematical struct {
	PatternID        string             // "sunset_001"
	Width, Height    int                // Dimensions origines
	BasisFunctions   int                // N = nombre de composantes (ex: 20)
	Coefficients     []float64          // αk pour chaque fonction (longueur: BasisFunctions × 3 pour RGB)
	BasisType        string             // "fourier", "gaussian", "polynomial"
	Reconstruction   float64            // Erreur quadratique moyenne (MSE) après reconstruction
	IntensityProfile []float64          // Profil d'intensité 1D pour patterns simples
	Tags             []string           // NEW: ["forest", "green", "nature"]
	Keywords         map[string]float64 // NEW: {"tree": 0.9, "leaf": 0.8} - confidence scores
}

// BasisFunctionEvaluator calcule gk(x,y) pour une position donnée
type BasisFunctionEvaluator struct {
	Type         string // "fourier", "gaussian", "polynomial", "mixed"
	BasisCount   int    // N = nombre de fonctions
	Width        int
	Height       int
	Coefficients []float64 // Stocké pour accès rapide
}

// ExtractPatternFromImage décompose une image en bases de Fourier
// Retourne les coefficients αk qui représentent le pattern
func ExtractPatternFromImage(imagePath string, basisCount int) (*PatternMathematical, error) {
	file, err := os.Open(imagePath)
	if err != nil {
		return nil, fmt.Errorf("cannot open image: %v", err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("cannot decode image: %v", err)
	}

	bounds := img.Bounds()
	width := bounds.Max.X - bounds.Min.X
	height := bounds.Max.Y - bounds.Min.Y

	pattern := &PatternMathematical{
		PatternID:      "extracted_pattern",
		Width:          width,
		Height:         height,
		BasisFunctions: basisCount,
		BasisType:      "fourier",
		Coefficients:   make([]float64, basisCount*3), // 3 channels RGB
	}

	// Étape 1: Normaliser l'image en [0, 1] pour chaque pixel et canal
	imageData := make([][][3]float64, height)
	for y := 0; y < height; y++ {
		imageData[y] = make([][3]float64, width)
		for x := 0; x < width; x++ {
			r, g, b, _ := img.At(bounds.Min.X+x, bounds.Min.Y+y).RGBA()
			// RGBA retourne 16-bit values, normaliser à [0, 1]
			imageData[y][x][0] = float64(r) / 65535.0
			imageData[y][x][1] = float64(g) / 65535.0
			imageData[y][x][2] = float64(b) / 65535.0
		}
	}

	// Étape 2: Appliquer décomposition Fourier pour chaque canal
	for channel := 0; channel < 3; channel++ {
		coeffs := DecomposeFourierBasis(imageData, width, height, basisCount, channel)
		copy(pattern.Coefficients[channel*basisCount:(channel+1)*basisCount], coeffs)
	}

	// Étape 3: Calculer l'erreur de reconstruction (MSE)
	reconstruction := ReconstructImage(pattern, width, height)
	mse := 0.0
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			for c := 0; c < 3; c++ {
				diff := imageData[y][x][c] - reconstruction[y*width+x][c]
				mse += diff * diff
			}
		}
	}
	pattern.Reconstruction = math.Sqrt(mse / float64(width*height*3))

	return pattern, nil
}

// DecomposeFourierBasis extrait les coefficients de Fourier pour un canal
func DecomposeFourierBasis(imageData [][][3]float64, width, height, basisCount, channel int) []float64 {
	coefficients := make([]float64, basisCount)

	// Meilleure décomposition Fourier 2D
	// gk(x,y) = cos(2π·kx·x/W) × cos(2π·ky·y/H)
	// Plus de détail capturé avec normalisation appropriée

	idx := 0
	sqrtBasis := int(math.Sqrt(float64(basisCount)) + 0.5)

	for ky := 0; ky < sqrtBasis; ky++ {
		for kx := 0; kx < sqrtBasis; kx++ {
			if idx >= basisCount {
				break
			}

			// Coefficient pour cette fréquence (intégrale)
			coeff := 0.0
			normalization := 0.0

			for y := 0; y < height; y++ {
				for x := 0; x < width; x++ {
					pixelValue := imageData[y][x][channel]

					// Basis functions: cosinus séparable 2D
					basisX := math.Cos(2 * math.Pi * float64(kx) * float64(x) / float64(width))
					basisY := math.Cos(2 * math.Pi * float64(ky) * float64(y) / float64(height))
					basisValue := basisX * basisY

					coeff += pixelValue * basisValue
					normalization += basisValue * basisValue
				}
			}

			// Normalisation robuste
			if normalization > 1e-10 {
				coefficients[idx] = coeff / normalization
			} else {
				coefficients[idx] = 0
			}

			idx++
		}
	}

	return coefficients
}

// ReconstructImage utilise les coefficients pour reconstruire une approximation de l'image
func ReconstructImage(pattern *PatternMathematical, width, height int) [][]float64 {
	reconstructed := make([][]float64, width*height)
	for i := range reconstructed {
		reconstructed[i] = make([]float64, 3)
	}

	evaluator := NewBasisFunctionEvaluator(pattern.BasisType, pattern.BasisFunctions, width, height)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			idx := y*width + x

			// Pour chaque canal RGB
			for channel := 0; channel < 3; channel++ {
				reconstructed[idx][channel] = 0.0

				// Sommation des composantes
				for k := 0; k < pattern.BasisFunctions; k++ {
					basis := evaluator.Evaluate(k, x, y)
					coeff := pattern.Coefficients[channel*pattern.BasisFunctions+k]
					reconstructed[idx][channel] += coeff * basis
				}

				// Clamp [0, 1]
				reconstructed[idx][channel] = math.Max(0, math.Min(1, reconstructed[idx][channel]))
			}
		}
	}

	return reconstructed
}

// NewBasisFunctionEvaluator crée un évaluateur de fonctions de base
func NewBasisFunctionEvaluator(basisType string, basisCount, width, height int) *BasisFunctionEvaluator {
	return &BasisFunctionEvaluator{
		Type:       basisType,
		BasisCount: basisCount,
		Width:      width,
		Height:     height,
	}
}

// Evaluate retourne gk(x,y) pour la k-ème fonction de base à position (x,y)
func (b *BasisFunctionEvaluator) Evaluate(k, x, y int) float64 {
	switch b.Type {
	case "fourier":
		return b.evaluateFourier(k, x, y)
	case "gaussian":
		return b.evaluateGaussian(k, x, y)
	case "polynomial":
		return b.evaluatePolynomial(k, x, y)
	default:
		return b.evaluateFourier(k, x, y)
	}
}

// evaluateFourier retourne cos(2π·kx·x/W) × cos(2π·ky·y/H)
func (b *BasisFunctionEvaluator) evaluateFourier(k, x, y int) float64 {
	sqrtN := int(math.Sqrt(float64(b.BasisCount)) + 0.5)
	ky := k / sqrtN
	kx := k % sqrtN

	return math.Cos(2*math.Pi*float64(kx)*float64(x)/float64(b.Width)) *
		math.Cos(2*math.Pi*float64(ky)*float64(y)/float64(b.Height))
}

// evaluateGaussian retourne exp(-(((x-cx)²+(y-cy)²)/σ²))
func (b *BasisFunctionEvaluator) evaluateGaussian(k, x, y int) float64 {
	sqrtN := int(math.Sqrt(float64(b.BasisCount)))
	cx := (k / sqrtN) * b.Width / sqrtN
	cy := (k % sqrtN) * b.Height / sqrtN
	sigma := float64(int(math.Max(float64(b.Width), float64(b.Height))) / sqrtN)

	dx := float64(x - cx)
	dy := float64(y - cy)
	distSq := dx*dx + dy*dy
	sigmaSq := sigma * sigma

	if sigmaSq <= 0 {
		return 0
	}
	return math.Exp(-distSq / (2 * sigmaSq))
}

// evaluatePolynomial retourne (x/W)^kx * (y/H)^ky
func (b *BasisFunctionEvaluator) evaluatePolynomial(k, x, y int) float64 {
	sqrtN := int(math.Sqrt(float64(b.BasisCount)))
	kx := k / sqrtN
	ky := k % sqrtN

	normX := float64(x) / float64(b.Width)
	normY := float64(y) / float64(b.Height)

	return math.Pow(normX, float64(kx)) * math.Pow(normY, float64(ky))
}

// ApplyPatternToAtomicNetwork applique le pattern comme contrainte à un réseau atomique
// Chaque atome reçoit: Citarget = f(xi, yi) basé sur le pattern
func ApplyPatternToAtomicNetwork(network *AtomicImageNetwork, pattern *PatternMathematical) {
	if network == nil || pattern == nil || len(network.Atoms) == 0 {
		return
	}

	evaluator := NewBasisFunctionEvaluator(pattern.BasisType, pattern.BasisFunctions, pattern.Width, pattern.Height)

	for i := 0; i < len(network.Atoms); i++ {
		// Accéder au champ Color du réseau s'il existe
		// TODO: Adapter selon la structure réelle d'AtomicImageNetwork
		_ = i
		_ = evaluator
		// Placeholder - à adapter selon la vraie structure
	}
}

// CombinePatterns combine plusieurs patterns avec des poids
// ResultColor[i] = Σ weight[p] * pattern[p].evaluate(i)
func CombinePatterns(patterns []*PatternMathematical, weights []float64, width, height int) [][][3]float64 {
	if len(patterns) == 0 {
		return nil
	}

	result := make([][][3]float64, height)
	for y := range result {
		result[y] = make([][3]float64, width)
	}

	// Normaliser les poids
	weightSum := 0.0
	for _, w := range weights {
		weightSum += w
	}
	if weightSum <= 0 {
		return result
	}

	for p, pattern := range patterns {
		evaluator := NewBasisFunctionEvaluator(pattern.BasisType, pattern.BasisFunctions, pattern.Width, pattern.Height)

		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				weight := weights[p] / weightSum

				for channel := 0; channel < 3; channel++ {
					for k := 0; k < pattern.BasisFunctions; k++ {
						basis := evaluator.Evaluate(k, x, y)
						coeff := pattern.Coefficients[channel*pattern.BasisFunctions+k]
						result[y][x][channel] += weight * coeff * basis
					}
				}
			}
		}
	}

	// Clamp final [0, 1]
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			for c := 0; c < 3; c++ {
				result[y][x][c] = math.Max(0, math.Min(1, result[y][x][c]))
			}
		}
	}

	return result
}

// PrintPatternAnalysis affiche l'analyse d'un pattern
func PrintPatternAnalysis(pattern *PatternMathematical) {
	fmt.Printf("\n" + strings.Repeat("=", 70) + "\n")
	fmt.Printf("📊 PATTERN MATHEMATICAL ANALYSIS\n")
	fmt.Printf(strings.Repeat("=", 70) + "\n")

	fmt.Printf("Pattern ID:          %s\n", pattern.PatternID)
	fmt.Printf("Dimensions:          %d×%d pixels\n", pattern.Width, pattern.Height)
	fmt.Printf("Basis Type:          %s\n", pattern.BasisType)
	fmt.Printf("Basis Functions:     %d\n", pattern.BasisFunctions)
	fmt.Printf("Reconstruction MSE:  %.6f (lower is better)\n", pattern.Reconstruction)

	fmt.Printf("\nCoefficients (αk) - Top 10 by magnitude:\n")
	type coeff struct {
		idx   int
		val   float64
		color string
	}
	coeffs := []coeff{}

	for c := 0; c < 3; c++ {
		channelName := []string{"Red", "Green", "Blue"}[c]
		for k := 0; k < pattern.BasisFunctions; k++ {
			coeffs = append(coeffs, coeff{
				idx:   k,
				val:   math.Abs(pattern.Coefficients[c*pattern.BasisFunctions+k]),
				color: channelName,
			})
		}
	}

	// Bubble sort (simple)
	for i := 0; i < len(coeffs); i++ {
		for j := i + 1; j < len(coeffs); j++ {
			if coeffs[j].val > coeffs[i].val {
				coeffs[i], coeffs[j] = coeffs[j], coeffs[i]
			}
		}
	}

	for i := 0; i < 10 && i < len(coeffs); i++ {
		fmt.Printf("  [%d] %s: α[%d] = %.6f\n", i+1, coeffs[i].color, coeffs[i].idx, coeffs[i].val)
	}

	fmt.Printf(strings.Repeat("=", 70) + "\n\n")
}

// PatternSimilarity calcule la similarité entre deux patterns
// Retourne 1.0 si identiques, 0.0 si totalement différents
func PatternSimilarity(p1, p2 *PatternMathematical) float64 {
	if p1 == nil || p2 == nil {
		return 0.0
	}

	if p1.BasisFunctions != p2.BasisFunctions {
		return 0.0
	}

	// Comparer les coefficients
	dotProduct := 0.0
	norm1 := 0.0
	norm2 := 0.0

	for i := 0; i < len(p1.Coefficients); i++ {
		dotProduct += p1.Coefficients[i] * p2.Coefficients[i]
		norm1 += p1.Coefficients[i] * p1.Coefficients[i]
		norm2 += p2.Coefficients[i] * p2.Coefficients[i]
	}

	if norm1 == 0 || norm2 == 0 {
		return 0.0
	}

	cosine := dotProduct / (math.Sqrt(norm1) * math.Sqrt(norm2))
	return math.Max(0, cosine) // Retourner [0, 1]
}

// SearchPatternsByKeyword cherche tous les patterns avec un mot-clé donné
// Retourne les patterns et leur score de pertinence
func SearchPatternsByKeyword(patterns []*PatternMathematical, keyword string) []*PatternMathematical {
	var results []*PatternMathematical
	keyword = strings.ToLower(keyword)

	for _, p := range patterns {
		if p == nil {
			continue
		}

		// Vérifier dans Tags
		for _, tag := range p.Tags {
			if strings.ToLower(tag) == keyword || strings.Contains(strings.ToLower(tag), keyword) {
				results = append(results, p)
				break
			}
		}

		// Vérifier dans Keywords
		for kw := range p.Keywords {
			if strings.ToLower(kw) == keyword || strings.Contains(strings.ToLower(kw), keyword) {
				found := false
				for _, r := range results {
					if r.PatternID == p.PatternID {
						found = true
						break
					}
				}
				if !found {
					results = append(results, p)
				}
				break
			}
		}
	}

	return results
}

// CombinePatternsByKeyword combine intelligemment tous les patterns matchant un keyword
// Poids = confidence score du keyword dans chaque pattern
func CombinePatternsByKeyword(patterns []*PatternMathematical, keyword string, width, height int) [][][3]float64 {
	// Chercher tous les patterns matchant
	matched := SearchPatternsByKeyword(patterns, keyword)

	if len(matched) == 0 {
		// Aucun match → créer image grise
		result := make([][][3]float64, height)
		for y := 0; y < height; y++ {
			result[y] = make([][3]float64, width)
			for x := 0; x < width; x++ {
				result[y][x] = [3]float64{0.5, 0.5, 0.5}
			}
		}
		return result
	}

	// Normaliser les poids selon la confiance du keyword
	weights := make([]float64, len(matched))
	totalWeight := 0.0

	for i, p := range matched {
		// Score par défaut = 1.0
		score := 1.0

		// Si le keyword est dans Keywords, utiliser sa confidence
		keyword_lower := strings.ToLower(keyword)
		for kw, conf := range p.Keywords {
			if strings.ToLower(kw) == keyword_lower {
				score = conf
				break
			}
		}

		weights[i] = score
		totalWeight += score
	}

	// Normaliser (Σ poids = 1)
	if totalWeight > 0 {
		for i := range weights {
			weights[i] /= totalWeight
		}
	}

	// Créer résultat vierge
	result := make([][][3]float64, height)
	for y := 0; y < height; y++ {
		result[y] = make([][3]float64, width)
	}

	// Combiner tous les patterns
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			color := [3]float64{0, 0, 0}

			for patIdx, p := range matched {
				weight := weights[patIdx]

				// Évaluer f_p(x, y)
				for c := 0; c < 3; c++ {
					for k := 0; k < p.BasisFunctions; k++ {
						evaluator := &BasisFunctionEvaluator{
							Type:       p.BasisType,
							BasisCount: p.BasisFunctions,
							Width:      p.Width,
							Height:     p.Height,
						}
						basis := evaluator.Evaluate(k, x, y)
						coeff := p.Coefficients[c*p.BasisFunctions+k]
						color[c] += weight * coeff * basis
					}
				}
			}

			// Clamp [0, 1]
			for c := 0; c < 3; c++ {
				if color[c] < 0 {
					color[c] = 0
				} else if color[c] > 1 {
					color[c] = 1
				}
			}
			result[y][x] = color
		}
	}

	return result
}
