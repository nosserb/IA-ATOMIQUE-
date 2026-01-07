package database

import (
	"regexp"
	"strings"
	"unicode"
)

// ============================================================================
// PHASE 15: MODULE 1 - PRÉTRAITEMENT & NETTOYAGE INTELLIGENT
// ============================================================================
// Supprime bruit (Gutenberg, URLs, métadonnées), normalise phrases,
// regroupe fragments courts, nettoie ponctuation bizarre

// PreprocessingResult contient le texte nettoyé et métadonnées
type PreprocessingResult struct {
	CleanedText     string   // Texte nettoyé
	OriginalLength  int      // Nombre de phrases originales
	CleanedLength   int      // Nombre de phrases après regroupement
	RemovedLines    []string // Lignes supprimées (bruit)
	NormalizedCount int      // Nombre de normalisations effectuées
	MergedSentences int      // Nombre de phrases regroupées
}

// NoisePatterns contient les regex pour détecter le bruit
var NoisePatterns = struct {
	ProjectGutenberg *regexp.Regexp
	URLs             *regexp.Regexp
	Metadata         *regexp.Regexp
	EmptyLines       *regexp.Regexp
	PageNumbers      *regexp.Regexp
	Credits          *regexp.Regexp
	Citations        *regexp.Regexp
}{
	ProjectGutenberg: regexp.MustCompile(`(?i)(Project Gutenberg|gutenberg\s+(etext|tm)|produced by|etext\s+\d+)`),
	URLs:             regexp.MustCompile(`https?://[^\s]+|www\.[^\s]+|\[http[^\]]*\]`),
	Metadata:         regexp.MustCompile(`(?i)(ISBN|ISSN|\*\*\*|===|---)`),
	EmptyLines:       regexp.MustCompile(`\n\s*\n+`),
	PageNumbers:      regexp.MustCompile(`^\s*\d+\s*$|^Page\s+\d+|^\[\d+\]`),
	Credits:          regexp.MustCompile(`(?i)(Copyright|©|All rights reserved|License|©\s*\d{4})`),
	Citations:        regexp.MustCompile(`^\s*[\[\(]\w+[\]\)]\s*`),
}

// PunctuationNormalization contient les replacements de ponctuation
var PunctuationNormalization = map[string]string{
	"\u201C": "\"",
	"\u201D": "\"",
	"\u00AB": "\"",
	"\u00BB": "\"",
	"\u201F": "\"",
	"\u201E": "\"",
	"\u2013": "-",
	"\u2014": "-",
	"\u2026": "...",
	"\u00B4": "'",
	"\u0060": "'",
	"\u0301": "'",
	"\u201A": "'",
	"\u203A": "'",
	"\u2039": "'",
}

// CommonStopwords pour identifier fragments courts triviaux
var FragmentStopwords = map[string]bool{
	"et":     true,
	"ou":     true,
	"mais":   true,
	"donc":   true,
	"cela":   true,
	"celui":  true,
	"celle":  true,
	"qui":    true,
	"que":    true,
	"dont":   true,
	"duquel": true,
	"auquel": true,
}

// PreprocessText effectue le nettoyage complet
func PreprocessText(text string) *PreprocessingResult {
	result := &PreprocessingResult{
		OriginalLength:  0,
		CleanedLength:   0,
		RemovedLines:    []string{},
		NormalizedCount: 0,
		MergedSentences: 0,
	}

	// Étape 1: Supprimer le bruit
	cleaned := RemoveNoise(text, result)

	// Étape 2: Normaliser la ponctuation
	cleaned = NormalizePunctuation(cleaned, result)

	// Étape 3: Nettoyer espaces et sauts de ligne
	cleaned = CleanWhitespace(cleaned)

	// Étape 4: Segmenter en phrases
	sentences := SegmentSentences(cleaned)
	result.OriginalLength = len(sentences)

	// Étape 5: Regrouper fragments courts
	mergedSentences := MergeShortFragments(sentences, result)
	result.CleanedLength = len(mergedSentences)

	// Étape 6: Reconstruire le texte
	result.CleanedText = strings.Join(mergedSentences, " ")

	return result
}

// RemoveNoise supprime tous les patterns de bruit détectés
func RemoveNoise(text string, result *PreprocessingResult) string {
	lines := strings.Split(text, "\n")
	var cleaned []string

	for _, line := range lines {
		// Vérifier tous les patterns de bruit
		if NoisePatterns.ProjectGutenberg.MatchString(line) {
			result.RemovedLines = append(result.RemovedLines, line)
			continue
		}
		if NoisePatterns.URLs.MatchString(line) {
			line = NoisePatterns.URLs.ReplaceAllString(line, "")
		}
		if NoisePatterns.Metadata.MatchString(line) {
			result.RemovedLines = append(result.RemovedLines, line)
			continue
		}
		if NoisePatterns.PageNumbers.MatchString(strings.TrimSpace(line)) {
			result.RemovedLines = append(result.RemovedLines, line)
			continue
		}
		if NoisePatterns.Credits.MatchString(line) {
			result.RemovedLines = append(result.RemovedLines, line)
			continue
		}

		// Garder la ligne si elle a du contenu
		if strings.TrimSpace(line) != "" {
			cleaned = append(cleaned, line)
		}
	}

	return strings.Join(cleaned, "\n")
}

// NormalizePunctuation normalise les variantes de ponctuation
func NormalizePunctuation(text string, result *PreprocessingResult) string {
	original := text

	// Appliquer les replacements
	for oldChar, newChar := range PunctuationNormalization {
		if strings.Contains(text, oldChar) {
			text = strings.ReplaceAll(text, oldChar, newChar)
			result.NormalizedCount++
		}
	}

	// Nettoyer ponctuation bizarre en fin de ligne
	text = regexp.MustCompile(`([.!?]){2,}`).ReplaceAllString(text, "$1")

	// Ajouter espace après ponctuation si manquant
	text = regexp.MustCompile(`([.!?])([A-Za-z])`).ReplaceAllString(text, "$1 $2")

	if text != original {
		result.NormalizedCount++
	}

	return text
}

// CleanWhitespace nettoie les espaces et sauts de ligne
func CleanWhitespace(text string) string {
	// Normaliser sauts de ligne multiples
	text = strings.ReplaceAll(text, "\r\n", "\n")
	text = strings.ReplaceAll(text, "\r", "\n")

	// Supprimer espaces multiples
	text = regexp.MustCompile(`  +`).ReplaceAllString(text, " ")

	// Supprimer sauts de ligne multiples
	text = regexp.MustCompile(`\n\n+`).ReplaceAllString(text, "\n")

	// Trim le début et fin
	text = strings.TrimSpace(text)

	return text
}

// SegmentSentences découpe le texte en phrases
func SegmentSentences(text string) []string {
	// Pattern pour détecter fin de phrase
	endPattern := regexp.MustCompile(`([.!?])\s+`)
	sentences := endPattern.Split(text, -1)

	var result []string
	for _, s := range sentences {
		s = strings.TrimSpace(s)
		if len(s) > 0 {
			result = append(result, s)
		}
	}

	return result
}

// MergeShortFragments regroupe les fragments courts avec la phrase suivante
func MergeShortFragments(sentences []string, result *PreprocessingResult) []string {
	if len(sentences) == 0 {
		return sentences
	}

	var merged []string
	var buffer string

	for _, sentence := range sentences {
		words := strings.Fields(sentence)

		// Si phrase très courte ET commence par stopword
		if len(words) <= 3 && len(words) > 0 {
			firstWord := strings.ToLower(words[0])
			if _, isStopword := FragmentStopwords[firstWord]; isStopword {
				// Fusionner avec la phrase précédente
				if len(merged) > 0 {
					merged[len(merged)-1] = merged[len(merged)-1] + " " + sentence
					result.MergedSentences++
					continue
				}
			}
		}

		// Sinon, ajouter comme phrase normale
		if buffer != "" {
			merged = append(merged, buffer)
			buffer = ""
		}
		merged = append(merged, sentence)
	}

	if buffer != "" {
		merged = append(merged, buffer)
	}

	return merged
}

// RemoveDuplicateWords supprime les mots répétés dans une phrase
func RemoveDuplicateWords(sentence string) string {
	words := strings.Fields(sentence)
	if len(words) <= 1 {
		return sentence
	}

	seen := make(map[string]int)
	var result []string
	maxRepeat := 2 // Permettre max 2 occurrences

	for _, word := range words {
		// Normaliser pour comparaison (lowercased, sans ponctuation)
		normalized := strings.ToLower(strings.TrimFunc(word, func(r rune) bool {
			return !unicode.IsLetter(r) && !unicode.IsNumber(r)
		}))

		seen[normalized]++

		// Garder le mot si < maxRepeat occurrences
		if seen[normalized] <= maxRepeat {
			result = append(result, word)
		}
	}

	return strings.Join(result, " ")
}

// CleanMultiplePunctuation nettoie la ponctuation excessive
func CleanMultiplePunctuation(text string) string {
	// Remplacer ... par ...
	text = regexp.MustCompile(`\.{4,}`).ReplaceAllString(text, "...")

	// Remplacer ??? par ?
	text = regexp.MustCompile(`\?{2,}`).ReplaceAllString(text, "?")
	text = regexp.MustCompile(`!{2,}`).ReplaceAllString(text, "!")

	// Remplacer ?! ou !? par ?
	text = regexp.MustCompile(`[?!]{2,}`).ReplaceAllString(text, "?")

	return text
}

// NormalizeContractions normalise les contractions françaises
func NormalizeContractions(text string) string {
	contractions := map[string]string{
		"d'un":     "de un",
		"d'une":    "de une",
		"d'le":     "de le",
		"d'la":     "de la",
		"d'les":    "de les",
		"l'un":     "le un",
		"l'une":    "le une",
		"l'autre":  "le autre",
		"l'auteur": "le auteur",
		"c'est":    "cela est",
		"s'il":     "si il",
		"s'ils":    "si ils",
		"m'a":      "me a",
		"t'a":      "te a",
		"j'ai":     "je ai",
		"n'a":      "ne a",
		"qu'un":    "que un",
		"qu'une":   "que une",
	}

	for old, new := range contractions {
		pattern := regexp.MustCompile(`(?i)\b` + regexp.QuoteMeta(old) + `\b`)
		text = pattern.ReplaceAllString(text, new)
	}

	return text
}

// IsValidSentence vérifie si une phrase a une structure valide
func IsValidSentence(sentence string) bool {
	words := strings.Fields(sentence)

	// Au minimum 2 mots
	if len(words) < 2 {
		return false
	}

	// Ne commence pas par ponctuation
	if len(words[0]) > 0 && !unicode.IsLetter(rune(words[0][0])) {
		return false
	}

	// Contient au moins une lettre
	hasLetter := false
	for _, word := range words {
		for _, char := range word {
			if unicode.IsLetter(char) {
				hasLetter = true
				break
			}
		}
		if hasLetter {
			break
		}
	}

	return hasLetter
}

// GetPreprocessingSummary retourne un résumé des opérations
func (pr *PreprocessingResult) GetPreprocessingSummary() string {
	var summary strings.Builder
	summary.WriteString("=== PREPROCESSING SUMMARY ===\n")
	summary.WriteString("Original sentences: " + string(rune(pr.OriginalLength)) + "\n")
	summary.WriteString("After cleanup: " + string(rune(pr.CleanedLength)) + "\n")
	summary.WriteString("Merged sentences: " + string(rune(pr.MergedSentences)) + "\n")
	summary.WriteString("Normalizations: " + string(rune(pr.NormalizedCount)) + "\n")
	summary.WriteString("Noise lines removed: " + string(rune(len(pr.RemovedLines))) + "\n")
	return summary.String()
}
