package database

import (
	"fmt"
	"strings"
)

// ============================================================================
// MATHEMATICAL INTEGRITY: Équations comme entités atomiques non compressibles
// ============================================================================
// Règle: ∀e∈Equations, e⊆R (toute équation doit être présente intégralement)

// MathBlock représente un bloc mathématique protégé
type MathBlock struct {
	ID      int
	RawText string
	LaTeX   string
}

// ProtectedContent contient texte avec équations marquées
type ProtectedContent struct {
	MarkedText string      // Texte avec [[MATH:id]] placeholders
	MathBlocks []MathBlock // Équations sauvegardées
	BlockCount int
}

// ============================================================================
// EXTRACTION AMÉLIORÉE: Identifier PHRASES contenant équations (pas juste lignes)
// ============================================================================

// ExtractAndProtectEquations V2: Extrait phrases + équations (stratégie robuste)
func ExtractAndProtectEquations(text string) ProtectedContent {
	result := ProtectedContent{
		MarkedText: text,
		MathBlocks: []MathBlock{},
		BlockCount: 0,
	}

	// Découper en phrases (délimitées par . ! ?)
	sentences := strings.FieldsFunc(text, func(r rune) bool {
		return r == '.' || r == '!' || r == '?'
	})

	var processedSentences []string
	blockID := 0

	for _, sentence := range sentences {
		trimmed := strings.TrimSpace(sentence)

		// Vérifier si la phrase contient une NOTATION MATHÉMATIQUE
		if ContainsMathNotation(trimmed) {
			// C'est une phrase avec équation ou notation formelle
			block := MathBlock{
				ID:      blockID,
				RawText: trimmed,
				LaTeX:   ExtractLaTeX(trimmed),
			}
			result.MathBlocks = append(result.MathBlocks, block)

			// Remplacer par placeholder
			placeholder := fmt.Sprintf("[[MATH:%d]]", blockID)
			processedSentences = append(processedSentences, placeholder)

			blockID++
		} else {
			processedSentences = append(processedSentences, trimmed)
		}
	}

	// Reconstruire le texte avec points
	var sb strings.Builder
	for i, sent := range processedSentences {
		sb.WriteString(sent)
		if i < len(processedSentences)-1 {
			sb.WriteString(". ")
		} else {
			sb.WriteString(".")
		}
	}

	result.MarkedText = sb.String()
	result.BlockCount = blockID

	return result
}

// ContainsMathNotation détecte si une phrase contient une notation mathématique
func ContainsMathNotation(text string) bool {
	// Symboles mathématiques "forts"
	strongPatterns := []string{
		"∀", "∃", "∈", "⊆", "⊂", "⊃", "∪", "∩", "∖",
		"λ", "μ", "σ", "ρ", "τ", "ω", "δ", "ε", "θ",
		"→", "↦", "⟹", "⟺", "≤", "≥", "≠", "≈", "~",
		"arg max", "arg min", "||", "||²",
		"cos(", "sin(", "exp(", "log(", "sqrt(", "tan(", "ln(",
	}

	// Vérifier patterns forts
	for _, pattern := range strongPatterns {
		if strings.Contains(text, pattern) {
			return true
		}
	}

	// Patterns faibles mais combinés
	hasOperator := strings.Contains(text, "=") ||
		strings.Contains(text, ":=") ||
		strings.Contains(text, "←")

	// Compter les indices/exposants (caractère Unicode de contrôle spécial)
	hasIndices := strings.Count(text, "​") >= 2 || // Zero-width space (indice caché)
		(strings.Count(text, "(") >= 2 && strings.Contains(text, "i")) ||
		(strings.Count(text, "(") >= 2 && strings.Contains(text, "j"))

	// Heuristique: si opérateur + au moins 2 variables indexées
	if hasOperator && hasIndices {
		return true
	}

	// Chercher patterns d'équation textuelles (formules verbalisées)
	eqTexts := []string{
		"peut être formalisée",
		"formule suivante",
		"équation",
		"fonction",
		"coefficient de", // Variables nommées
		"paramètre",      // Dans contexte d'équation
	}

	for _, eqText := range eqTexts {
		if strings.Contains(strings.ToLower(text), eqText) {
			// Double-check: doit avoir au moins un opérateur ou symbole
			if hasOperator || strings.Count(text, "​") > 0 {
				return true
			}
		}
	}

	return false
}

// IsEquationLine est une fonction auxiliaire (utilise plutôt ContainsMathNotation)
func IsEquationLine(line string) bool {
	return ContainsMathNotation(line)
}

// ExtractLaTeX récupère la forme LaTeX (simplifié)
func ExtractLaTeX(line string) string {
	// Pour l'instant, c'est la ligne elle-même
	// Dans une vraie implémentation, convertir vers LaTeX
	return line
}

// ============================================================================
// ÉTAPE B: Préservation des équations dans le résumé
// ============================================================================

// PreserveEquationsInSummary garantit que toutes les équations sont présentes
func PreserveEquationsInSummary(summary string, originalProtected ProtectedContent) string {
	result := summary

	// Pour chaque équation, vérifier sa présence
	missingEquations := []MathBlock{}

	for _, block := range originalProtected.MathBlocks {
		placeholder := fmt.Sprintf("[[MATH:%d]]", block.ID)

		if !strings.Contains(result, placeholder) &&
			!strings.Contains(result, block.RawText) {
			missingEquations = append(missingEquations, block)
		}
	}

	// Si équations manquantes, les ajouter à la fin
	if len(missingEquations) > 0 {
		result += "\n\n[Équations critiques]\n"
		for _, block := range missingEquations {
			result += fmt.Sprintf("[[MATH:%d]] %s\n", block.ID, block.RawText)
		}
	}

	return result
}

// RestoreEquationsFromPlaceholders remplace [[MATH:id]] par équations réelles
func RestoreEquationsFromPlaceholders(markedText string, mathBlocks []MathBlock) string {
	result := markedText

	for _, block := range mathBlocks {
		placeholder := fmt.Sprintf("[[MATH:%d]]", block.ID)
		result = strings.ReplaceAll(result, placeholder, block.RawText)
	}

	return result
}

// ============================================================================
// ÉTAPE C: Métrique binaire pour équations (1 si toutes présentes, 0 sinon)
// ============================================================================

// CalculateEquationIntegrityScore retourne 1.0 si toutes les équations sont présentes
func CalculateEquationIntegrityScore(summary string, mathBlocks []MathBlock) float64 {
	if len(mathBlocks) == 0 {
		return 1.0 // Pas d'équations = pas de problème
	}

	summaryLower := strings.ToLower(summary)

	foundCount := 0
	for _, block := range mathBlocks {
		// Chercher soit le placeholder, soit le contenu brut
		placeholder := fmt.Sprintf("[[MATH:%d]]", block.ID)
		blockLower := strings.ToLower(block.RawText)

		if strings.Contains(summary, placeholder) ||
			strings.Contains(summaryLower, blockLower) {
			foundCount++
		}
	}

	// Métrique binaire: 1 si 100%, 0 sinon
	if foundCount == len(mathBlocks) {
		return 1.0
	}
	return 0.0

	// Alternative: score progressif (si on veut être moins strict)
	// return float64(foundCount) / float64(len(mathBlocks))
}

// ============================================================================
// Nouvelle fidélité pondérée avec contrainte mathématique
// ============================================================================

// CalculateWeightedFidelityWithMathConstraint = α*Concept + β*Equation(binaire) + γ*Text
func CalculateWeightedFidelityWithMathConstraint(
	summary string,
	sourceText string,
	mathBlocks []MathBlock,
) float64 {
	// Scores individuels
	conceptScore := CalculateConceptualFidelity(summary, sourceText)
	equationScore := CalculateEquationIntegrityScore(summary, mathBlocks)
	textScore := CalculateFidelity(summary, sourceText)

	// Poids: équations très importants (β ≥ 0.4)
	alpha := 0.3 // Concepts
	beta := 0.5  // Équations (CRITIQUE)
	gamma := 0.2 // Texte narratif

	weightedScore := (alpha * conceptScore) +
		(beta * equationScore) +
		(gamma * textScore)

	return weightedScore
}

// CalculateConceptualFidelity mesure la préservation des concepts
func CalculateConceptualFidelity(summary string, sourceText string) float64 {
	sourceConcepts := ExtractPrimaryConceptsFromText(sourceText)

	if len(sourceConcepts) == 0 {
		return 1.0
	}

	summaryLower := strings.ToLower(summary)
	foundCount := 0

	for _, concept := range sourceConcepts {
		if strings.Contains(summaryLower, concept) {
			foundCount++
		}
	}

	return float64(foundCount) / float64(len(sourceConcepts))
}
