// Package database - Perplexity Calculation
// Calcul de la perplexité basée sur la résonance atomique
//
// PRINCIPE : Perplexité = mesure de "surprise" de l'IA face à un texte
// - Perplexité BASSE = texte prévisible, structure cohérente
// - Perplexité HAUTE = texte surprenant, structure chaotique
//
// IMPLÉMENTATION ATOMIQUE :
// - Perplexité = Instabilité énergétique du réseau atomique
// - Réseau stable (haute cohérence) = faible perplexité
// - Réseau instable (basse cohérence) = haute perplexité
//
// AVANTAGE : Mesure physique directe, pas de probabilités artificielles

package database

import (
	"math"
)

// PerplexityCalculator calculateur de perplexité atomique
type PerplexityCalculator struct {
	SampleSize    int     // Taille des échantillons pour calcul
	MinPerplexity float64 // Perplexité minimale (texte parfait)
	MaxPerplexity float64 // Perplexité maximale (chaos total)
}

// PerplexityResult résultat du calcul de perplexité
type PerplexityResult struct {
	GlobalPerplexity   float64     // Perplexité globale
	LocalPerplexities  []float64   // Perplexités locales par segment
	AverageCoherence   float64     // Cohérence moyenne du réseau
	EnergyVariance     float64     // Variance énergétique
	StabilityScore     float64     // Score de stabilité [0-1]
	CategoryDistrib    map[int]int // Distribution des catégories
	SurpriseMoments    []int       // Positions des moments de surprise
	InterpretedQuality string      // Interprétation qualitative
}

// NewPerplexityCalculator crée un nouveau calculateur
func NewPerplexityCalculator() *PerplexityCalculator {
	return &PerplexityCalculator{
		SampleSize:    100,  // Échantillons de 100 mots
		MinPerplexity: 1.0,  // Log2 des prédictions parfaites
		MaxPerplexity: 15.0, // Log2 du chaos (2^15 = 32768 possibilités)
	}
}

// CalculatePerplexity calcule la perplexité atomique d'un texte
func (calc *PerplexityCalculator) CalculatePerplexity(texte string) *PerplexityResult {
	// Phase 1 : Tokenisation
	tokens := TokeniserTexte(texte)
	if len(tokens) == 0 {
		return &PerplexityResult{
			GlobalPerplexity:   calc.MaxPerplexity,
			InterpretedQuality: "Texte vide ou invalide",
		}
	}

	result := &PerplexityResult{
		LocalPerplexities: make([]float64, 0),
		CategoryDistrib:   make(map[int]int),
		SurpriseMoments:   make([]int, 0),
	}

	// Phase 2 : Découpage en segments
	segments := calc.createSegments(tokens)

	// Phase 3 : Calcul de perplexité locale pour chaque segment
	coherenceSum := 0.0
	energySum := 0.0
	energies := make([]float64, 0)

	for i, segment := range segments {
		// Activer le réseau pour ce segment
		categories := ActiverCategoriesParTexte(segment)

		// Accumuler distribution globale
		for catID, count := range categories {
			result.CategoryDistrib[catID] += count
		}

		// Calculer cohérence locale (entropie inversée)
		coherence := calc.calculateSegmentCoherence(categories)
		coherenceSum += coherence

		// Calculer énergie locale
		energy := 0.0
		for _, count := range categories {
			energy += float64(count)
		}
		energySum += energy
		energies = append(energies, energy)

		// Convertir cohérence en perplexité
		// Perplexité = 2^(-coherence * facteur)
		// Cohérence haute → perplexité basse
		localPerplexity := calc.coherenceToPerplexity(coherence)
		result.LocalPerplexities = append(result.LocalPerplexities, localPerplexity)

		// Détecter moments de surprise (perplexité anormalement haute)
		if localPerplexity > calc.MaxPerplexity*0.7 {
			result.SurpriseMoments = append(result.SurpriseMoments, i*calc.SampleSize)
		}
	}

	// Phase 4 : Agrégation statistiques globales
	if len(segments) > 0 {
		result.AverageCoherence = coherenceSum / float64(len(segments))

		// Variance énergétique (mesure d'instabilité)
		meanEnergy := energySum / float64(len(energies))
		varianceSum := 0.0
		for _, e := range energies {
			diff := e - meanEnergy
			varianceSum += diff * diff
		}
		result.EnergyVariance = varianceSum / float64(len(energies))

		// Score de stabilité = cohérence normalisée
		result.StabilityScore = result.AverageCoherence

		// Perplexité globale = moyenne géométrique des perplexités locales
		result.GlobalPerplexity = calc.geometricMean(result.LocalPerplexities)
	}

	// Phase 5 : Interprétation qualitative
	result.InterpretedQuality = calc.interpretPerplexity(result.GlobalPerplexity)

	return result
}

// createSegments découpe les tokens en segments de taille fixe
func (calc *PerplexityCalculator) createSegments(tokens []string) [][]string {
	segments := make([][]string, 0)

	for i := 0; i < len(tokens); i += calc.SampleSize {
		end := i + calc.SampleSize
		if end > len(tokens) {
			end = len(tokens)
		}
		segment := tokens[i:end]
		segments = append(segments, segment)
	}

	return segments
}

// calculateSegmentCoherence calcule la cohérence d'un segment
// Basé sur l'entropie de Shannon de la distribution des catégories
func (calc *PerplexityCalculator) calculateSegmentCoherence(categories map[int]int) float64 {
	if len(categories) == 0 {
		return 0.0
	}

	// Calculer total
	total := 0.0
	for _, count := range categories {
		total += float64(count)
	}

	if total < 0.001 {
		return 0.0
	}

	// Entropie de Shannon
	entropy := 0.0
	for _, count := range categories {
		if count > 0 {
			p := float64(count) / total
			entropy -= p * math.Log2(p)
		}
	}

	// Normaliser par entropie maximale
	maxEntropy := math.Log2(float64(len(categories)))
	if maxEntropy < 0.001 {
		return 1.0
	}

	normalizedEntropy := entropy / maxEntropy

	// Cohérence = 1 - entropie normalisée
	coherence := 1.0 - normalizedEntropy

	return coherence
}

// coherenceToPerplexity convertit la cohérence en perplexité
func (calc *PerplexityCalculator) coherenceToPerplexity(coherence float64) float64 {
	// Formule : Perplexity = MinPerplexity * 2^((1-coherence) * facteur)
	// Coherence = 1.0 → Perplexity = MinPerplexity
	// Coherence = 0.0 → Perplexity = MaxPerplexity

	factor := math.Log2(calc.MaxPerplexity / calc.MinPerplexity)
	exponent := (1.0 - coherence) * factor

	perplexity := calc.MinPerplexity * math.Pow(2, exponent)

	// Clamper dans les bornes
	if perplexity < calc.MinPerplexity {
		perplexity = calc.MinPerplexity
	}
	if perplexity > calc.MaxPerplexity {
		perplexity = calc.MaxPerplexity
	}

	return perplexity
}

// geometricMean calcule la moyenne géométrique
func (calc *PerplexityCalculator) geometricMean(values []float64) float64 {
	if len(values) == 0 {
		return calc.MaxPerplexity
	}

	product := 1.0
	for _, v := range values {
		product *= v
	}

	return math.Pow(product, 1.0/float64(len(values)))
}

// interpretPerplexity interprète la perplexité en texte
func (calc *PerplexityCalculator) interpretPerplexity(perplexity float64) string {
	if perplexity < 2.0 {
		return "EXCELLENT - Texte très cohérent et prévisible"
	} else if perplexity < 4.0 {
		return "BON - Texte bien structuré"
	} else if perplexity < 6.0 {
		return "MOYEN - Texte correct mais avec variations"
	} else if perplexity < 9.0 {
		return "FAIBLE - Texte peu cohérent"
	} else {
		return "MAUVAIS - Texte chaotique ou incohérent"
	}
}

// CompareTexts compare la perplexité de deux textes
func (calc *PerplexityCalculator) CompareTexts(texte1, texte2 string) (result1, result2 *PerplexityResult, comparison string) {
	result1 = calc.CalculatePerplexity(texte1)
	result2 = calc.CalculatePerplexity(texte2)

	diff := result1.GlobalPerplexity - result2.GlobalPerplexity

	if math.Abs(diff) < 0.5 {
		comparison = "Les deux textes ont une qualité similaire"
	} else if diff < 0 {
		comparison = "Le premier texte est plus cohérent"
	} else {
		comparison = "Le second texte est plus cohérent"
	}

	return result1, result2, comparison
}
