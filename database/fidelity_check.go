package database

import (
	"regexp"
	"sort"
	"strings"
)

// ============================================================================
// FIDÉLITÉ DU RÉSUMÉ: Vérification que R ⊆ C(T)
// ============================================================================
// Problème: Phase 15 génère du contenu hallucin (invente des concepts)
// Solution: Mesurer fidélité Ff(R,T) et forcer fallback extractif si Ff < τ

// KeyTermsSet extrait les concepts clés du texte source
// C(T) = ensemble des mots/concepts importants dans T
func ExtractKeyTerms(text string) map[string]bool {
	keyTerms := make(map[string]bool)

	// Termes techniques importants à détecter
	technicalPatterns := []string{
		// IA Atomique
		"IA atomique", "atome", "computationnel", "résonance", "T.R.A",
		"asynchrone", "asynchronism", "asynchronisme",
		"moteur d'inférence", "inférence", "moteur",
		"architecture distribuée", "distribuée", "décentralisée",
		"interaction locale", "local", "interactions",
		"comportement émergent", "émergence", "émergent",
		"atomes", "réseau", "nœud",

		// Concepts clés
		"synchrone", "synchrone", "centralisé", "centralisée",
		"micro-interaction", "resonance", "couplage",
		"plasticité", "poids de connexion", "apprentissage",
		"capteur", "robot", "IoT", "urbain",
		"système multi-agent", "SMA", "agent",

		// Concepts NÉGATIFS (ce qu'on ne devrait PAS voir)
		// (Ces mots dans résumé = hallucination potentielle)
	}

	textLower := strings.ToLower(text)

	// Extraire termes techniques
	for _, term := range technicalPatterns {
		if strings.Contains(textLower, strings.ToLower(term)) {
			// Normaliser et ajouter
			normalized := strings.ToLower(term)
			keyTerms[normalized] = true

			// Extraire aussi les variantes
			parts := strings.Fields(normalized)
			for _, part := range parts {
				if len(part) > 4 { // Ignorer articles/prépositions courts
					keyTerms[part] = true
				}
			}
		}
	}

	// Extraire mots fréquents du texte (top 30)
	words := extractFrequentWords(textLower, 30)
	for word := range words {
		keyTerms[word] = true
	}

	return keyTerms
}

// CalculateFidelity mesure Ff(R,T) = |C(R) ∩ C(T)| / |C(R)|
// Retourne score entre 0 et 1
// Score élevé = résumé fidèle au texte original
// Score faible = résumé hallucine
func CalculateFidelity(summary string, sourceText string) float64 {
	if summary == "" || sourceText == "" {
		return 0.0
	}

	// Extraire concepts clés du source
	sourceTerms := ExtractKeyTerms(sourceText)

	// Extraire mots du résumé
	summaryWords := extractWordsFromText(strings.ToLower(summary))

	if len(summaryWords) == 0 {
		return 0.0
	}

	// Calculer intersection
	intersection := 0
	for word := range summaryWords {
		if sourceTerms[word] {
			intersection++
		}
	}

	// Fidélité = mots du résumé présents dans source / total mots résumé
	fidelity := float64(intersection) / float64(len(summaryWords))

	return fidelity
}

// ExtractiveResume génère résumé extractif basé sur phrases clés
// Stratégie A: Sélectionner les n meilleures phrases contenant concepts clés
func ExtractiveResume(text string, compressionRatio float64) string {
	sentences := SplitSentences(text)

	if len(sentences) == 0 {
		return ""
	}

	// Calculer nombre de phrases à garder
	targetSentences := max(1, int(float64(len(sentences))*compressionRatio))

	// Extraire termes clés
	keyTerms := ExtractKeyTerms(text)

	// Scorer chaque phrase
	type SentenceScore struct {
		sentence string
		score    float64
	}

	var scores []SentenceScore

	for _, sentence := range sentences {
		score := scoreSentenceByKeyTerms(sentence, keyTerms)
		if score > 0 {
			scores = append(scores, SentenceScore{sentence, score})
		}
	}

	// Trier par score (décroissant)
	sort.Slice(scores, func(i, j int) bool {
		return scores[i].score > scores[j].score
	})

	// Prendre top N phrases
	if len(scores) > targetSentences {
		scores = scores[:targetSentences]
	}

	// Reconstruire dans ordre original
	// (Trier par position dans texte original)
	sort.Slice(scores, func(i, j int) bool {
		posI := strings.Index(text, scores[i].sentence)
		posJ := strings.Index(text, scores[j].sentence)
		return posI < posJ
	})

	// Joindre phrases
	result := ""
	for _, ss := range scores {
		result += ss.sentence + " "
	}

	return strings.TrimSpace(result)
}

// ============================================================================
// Fonctions helper
// ============================================================================

// SplitSentences divise texte en phrases
func SplitSentences(text string) []string {
	// Régex pour points, points d'interrogation, etc.
	re := regexp.MustCompile(`[.!?]+`)
	sentences := re.Split(text, -1)

	var result []string
	for _, s := range sentences {
		trimmed := strings.TrimSpace(s)
		if len(trimmed) > 10 { // Ignorer fragments trop courts
			result = append(result, trimmed)
		}
	}

	return result
}

// extractWordsFromText extrait mots d'un texte
func extractWordsFromText(text string) map[string]bool {
	words := make(map[string]bool)

	// Supprimer ponctuation
	re := regexp.MustCompile(`[^a-zàâäæçéèêëïîôùûüœ0-9\s]`)
	cleaned := re.ReplaceAllString(text, " ")

	parts := strings.Fields(cleaned)

	// Filtrer stopwords et mots trop courts
	stopwords := getStopwords()

	for _, word := range parts {
		if len(word) > 3 && !stopwords[word] {
			words[word] = true
		}
	}

	return words
}

// extractFrequentWords extrait mots les plus fréquents
func extractFrequentWords(text string, topN int) map[string]bool {
	words := extractWordsFromText(text)

	result := make(map[string]bool)
	count := 0

	for word := range words {
		result[word] = true
		count++
		if count >= topN {
			break
		}
	}

	return result
}

// scoreSentenceByKeyTerms donne score à une phrase basée sur termes clés
func scoreSentenceByKeyTerms(sentence string, keyTerms map[string]bool) float64 {
	sentenceLower := strings.ToLower(sentence)
	score := 0.0

	for term := range keyTerms {
		if strings.Contains(sentenceLower, term) {
			score += 1.0
		}
	}

	return score
}

// getStopwords retourne liste de stopwords français
func getStopwords() map[string]bool {
	return map[string]bool{
		"le": true, "la": true, "les": true, "de": true, "du": true,
		"des": true, "un": true, "une": true, "et": true, "ou": true,
		"en": true, "à": true, "au": true, "que": true, "qui": true,
		"se": true, "pour": true, "dans": true, "par": true, "avec": true,
		"est": true, "sont": true, "avoir": true, "être": true, "ce": true,
		"ça": true, "il": true, "elle": true, "mais": true, "donc": true,
		"cependant": true, "ainsi": true, "soit": true, "tout": true,
	}
}

// max retourne le maximum de deux entiers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// HybridResume implémente stratégie C: Hybridation extractif + génératif
// Si résumé génératif Rg a fidélité < seuil → utiliser extractif Re
func HybridResume(generatedSummary string, sourceText string, fidelityThreshold float64) (string, float64, string) {
	// Calculer fidélité du résumé généré
	fidelity := CalculateFidelity(generatedSummary, sourceText)

	// Compression ratio estimé
	compressionRatio := float64(len(generatedSummary)) / float64(len(sourceText))
	if compressionRatio > 1.0 {
		compressionRatio = 0.1 // Par défaut 10% si résumé plus long
	}

	var finalSummary string
	var mode string

	if fidelity >= fidelityThreshold {
		// Résumé généré est fidèle → le garder
		finalSummary = generatedSummary
		mode = "GÉNÉRATIF (fidèle)"
	} else {
		// Résumé généré hallucine → utiliser extractif
		finalSummary = ExtractiveResume(sourceText, compressionRatio)
		mode = "EXTRACTIF (hallucination détectée)"
	}

	return finalSummary, fidelity, mode
}
