package database

import (
	"math"
	"strings"
)

// ============================================================================
// WEIGHTED FIDELITY METRIC
// ============================================================================
// Ff_weighted = Σ w_i * overlap_i
// où w: concepts/équations=5, texte narratif=1
// Évite de pénaliser les reformulations intelligentes

// WeightedFidelityReport contient les détails du calcul de fidélité pondérée
type WeightedFidelityReport struct {
	TotalScore       float64 // Score pondéré final
	ConceptScore     float64 // w=3 pour concepts clés
	EquationScore    float64 // w=5 pour équations/symboles
	NarrativeScore   float64 // w=1 pour texte normal
	WeightedScore    float64 // Score final [0,1]
	ConceptsMatched  int
	ConceptsTotal    int
	EquationsMatched int
	EquationsTotal   int
}

// CalculateWeightedFidelity implémente Ff_weighted
func CalculateWeightedFidelity(summary string, sourceText string) WeightedFidelityReport {
	report := WeightedFidelityReport{}

	// 1. Extraire les concepts clés du source
	sourceConcepts := ExtractPrimaryConceptsFromText(sourceText)
	report.ConceptsTotal = len(sourceConcepts)

	// Chercher lesquels sont dans le résumé
	summaryLower := strings.ToLower(summary)
	conceptMatches := 0
	for _, concept := range sourceConcepts {
		if strings.Contains(summaryLower, concept) {
			conceptMatches++
		}
	}
	report.ConceptsMatched = conceptMatches
	report.ConceptScore = float64(conceptMatches) / float64(math.Max(1, float64(report.ConceptsTotal)))

	// 2. Extraire les équations/symboles mathématiques
	sourceEquations := ExtractEquationsAndSymbols(sourceText)
	report.EquationsTotal = len(sourceEquations)

	equationMatches := 0
	for _, eq := range sourceEquations {
		if strings.Contains(summaryLower, eq) {
			equationMatches++
		}
	}
	report.EquationsMatched = equationMatches
	report.EquationScore = float64(equationMatches) / float64(math.Max(1, float64(report.EquationsTotal)))

	// 3. Score narratif (mots communs)
	report.NarrativeScore = CalculateFidelity(summary, sourceText) // Simple overlap

	// 4. Score pondéré
	// w_concepts = 3, w_equations = 5, w_narrative = 1
	weights := map[string]float64{
		"concept":   3.0,
		"equation":  5.0,
		"narrative": 1.0,
	}

	weightedSum := (report.ConceptScore * weights["concept"]) +
		(report.EquationScore * weights["equation"]) +
		(report.NarrativeScore * weights["narrative"])

	totalWeight := weights["concept"] + weights["equation"] + weights["narrative"]

	report.WeightedScore = weightedSum / totalWeight

	return report
}

// ExtractEquationsAndSymbols extrait les équations et symboles mathématiques
func ExtractEquationsAndSymbols(text string) []string {
	var equations []string

	// Patterns courants
	mathSymbols := []string{
		"∀", "∃", "∈", "∉", "⊆", "⊂", "∪", "∩",
		"λ", "μ", "σ", "ρ", "τ", "ω", "δ", "ε",
		"→", "↦", "⟹", "⟺", "≤", "≥", "≠",
		"arg max", "arg min", "Σ", "∫", "∂",
		"cos(", "sin(", "exp(", "log(",
	}

	for _, symbol := range mathSymbols {
		if strings.Contains(text, symbol) {
			equations = append(equations, symbol)
		}
	}

	// Extraire les "définitions" (lignes avec "=" ou ":=")
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		if strings.Contains(line, "=") || strings.Contains(line, ":=") {
			// Simplifier et ajouter
			trimmed := strings.TrimSpace(line)
			if len(trimmed) > 10 && len(trimmed) < 200 {
				equations = append(equations, trimmed)
			}
		}
	}

	return equations
}

// ============================================================================
// ABSTRACTION BLENDING: R = λ*R_extractive + (1-λ)*R_abstractive
// ============================================================================

// CalculateBlendingLambda = max(0, 1 - Ff)
// Si Ff élevée (bon alignement) → λ petit (plus d'abstraction)
// Si Ff basse (divergence) → λ grand (plus d'extraction)
func CalculateBlendingLambda(fidelityScore float64) float64 {
	lambda := math.Max(0, 1.0-fidelityScore)
	return math.Min(1.0, lambda)
}

// BlendSummaries fusionne extractif et abstractif avec poids λ
func BlendSummaries(extractiveSummary string, abstractiveSummary string, lambda float64) string {
	if lambda > 0.5 {
		// Extractif dominant
		return extractiveSummary
	} else if lambda < 0.3 {
		// Abstractif dominant
		return abstractiveSummary
	}

	// Blend: mélanger les deux
	// Stratégie: prendre les phrases du extractif pour les concepts clés
	// et les phrases du abstractif pour les transitions

	extractiveLines := strings.Split(extractiveSummary, ".")
	abstractiveLines := strings.Split(abstractiveSummary, ".")

	var blended []string

	// Alterner: 1 extractif, 1 abstractif
	maxLines := len(extractiveLines)
	if len(abstractiveLines) > maxLines {
		maxLines = len(abstractiveLines)
	}

	for i := 0; i < maxLines; i++ {
		if i < len(extractiveLines) && i%2 == 0 {
			blended = append(blended, extractiveLines[i])
		} else if i < len(abstractiveLines) {
			blended = append(blended, abstractiveLines[i])
		}
	}

	return strings.Join(blended, ". ")
}

// ============================================================================
// EQUATION PROTECTION
// ============================================================================

// TagEquationsForProtection entoure les équations de tags spéciaux
func TagEquationsForProtection(text string) string {
	lines := strings.Split(text, "\n")
	var protected []string

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Détecteur d'équation simple: contient symboles math + "="
		isEquation := false
		mathSymbols := []string{"∀", "∃", "arg max", "∈", "⊆", "λ", "→", ":="}
		for _, sym := range mathSymbols {
			if strings.Contains(line, sym) && strings.Contains(line, "=") {
				isEquation = true
				break
			}
		}

		if isEquation {
			protected = append(protected, "[[MATH_PROTECTED]]"+trimmed+"[[/MATH_PROTECTED]]")
		} else {
			protected = append(protected, line)
		}
	}

	return strings.Join(protected, "\n")
}

// UntagEquations enlève les tags de protection
func UntagEquations(text string) string {
	text = strings.ReplaceAll(text, "[[MATH_PROTECTED]]", "")
	text = strings.ReplaceAll(text, "[[/MATH_PROTECTED]]", "")
	return text
}

// ExtractProtectedEquations extrait les équations protégées
func ExtractProtectedEquations(text string) []string {
	var equations []string

	start := 0
	for {
		idx := strings.Index(text[start:], "[[MATH_PROTECTED]]")
		if idx == -1 {
			break
		}

		startIdx := start + idx + len("[[MATH_PROTECTED]]")
		endIdx := strings.Index(text[startIdx:], "[[/MATH_PROTECTED]]")
		if endIdx == -1 {
			break
		}

		equation := text[startIdx : startIdx+endIdx]
		equations = append(equations, equation)

		start = startIdx + endIdx + len("[[/MATH_PROTECTED]]")
	}

	return equations
}
