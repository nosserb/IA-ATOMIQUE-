package database

import (
	"math"
	"strings"
)

// CalculateSentenceDomainSimilarity = cos(phrase→, prototype_domaine→)
// Mesure la similarité sémantique d'une phrase au domaine
func CalculateSentenceDomainSimilarity(sentence string, ds *DomainSpace) float64 {
	if len(ds.CoreConcepts) == 0 {
		return 0.5 // Défaut neutre
	}

	sentenceLower := strings.ToLower(sentence)

	// Compter les matches avec concepts principaux
	matches := 0.0
	for _, concept := range ds.CoreConcepts {
		if strings.Contains(sentenceLower, concept) {
			matches += 1.0
		}
	}

	// Similarité = proportion de concepts trouvés
	similarity := matches / float64(len(ds.CoreConcepts))

	// Bonus pour présence de termes techniques (même peu denses)
	for term := range ds.TechTerms {
		if strings.Contains(sentenceLower, term) {
			similarity = math.Min(1.0, similarity+0.15)
		}
	}

	return similarity
}

// CalculateAllowedAbstraction = α*(1-ρtech) + β*C
// Détermine le niveau d'abstraction autorisé
// Plus c'est conceptuel, plus on autorise l'abstraction
func CalculateAllowedAbstraction(techDensity float64, conceptualScore float64) float64 {
	alpha := 0.6 // Poids de la faible densité technique
	beta := 0.4  // Poids du score conceptuel

	// Si très technique -> peu d'abstraction
	// Si peu technique et très conceptuel -> beaucoup d'abstraction
	allowedAbstraction := alpha*(1.0-techDensity) + beta*conceptualScore

	return math.Max(0.3, math.Min(0.9, allowedAbstraction))
}

// CalculateConceptualScore mesure le caractère conceptuel/théorique du texte
func CalculateConceptualScore(text string) float64 {
	textLower := strings.ToLower(text)

	// Patterns de textes conceptuels/théoriques
	conceptualPatterns := map[string]int{
		"concept": 1, "théorie": 2, "abstrait": 1, "principe": 1,
		"architecture": 1, "design": 1, "pattern": 1, "cadre": 1,
		"approche": 1, "méthode": 1, "framework": 1, "modèle": 1,
		"abstraction": 2, "niveau": 1, "couche": 1, "interface": 1,
		"propriété": 1, "caractéristique": 1, "définition": 1,
		"suppose": 1, "considère": 1, "envisage": 1, "propose": 1,
	}

	totalScore := 0.0
	for pattern, weight := range conceptualPatterns {
		count := float64(strings.Count(textLower, pattern))
		totalScore += count * float64(weight)
	}

	// Normaliser par longueur du texte
	words := float64(len(strings.Fields(text)))
	if words == 0 {
		return 0.0
	}

	conceptualScore := math.Min(1.0, totalScore/words*10)
	return conceptualScore
}

// ============================================================================
// ACTION 3: Adapter compression minimale selon type
// ============================================================================

// GetMinCompressionForDomain retourne la compression minimale par domaine
func GetMinCompressionForDomain(domainMode string) float64 {
	minCompressions := map[string]float64{
		"ornithology":             0.1,  // Peut compresser agressivement (encyclopédique)
		"technical-ai-hard":       0.3,  // Code/équations: moins de compression
		"technical-ai-conceptual": 0.4,  // Architecture/théorie: modéré
		"technical":               0.25, // Technique standard
		"narrative":               0.2,  // Narratif
		"social":                  0.3,  // Social
	}

	if compression, exists := minCompressions[domainMode]; exists {
		return compression
	}
	return 0.5 // Défaut
}
