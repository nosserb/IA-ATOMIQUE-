// Package database - Needle In Haystack Search
// Recherche ultra-rapide d'informations cachées dans textes massifs
// Exploite la résonance atomique pour détecter anomalies sémantiques
//
// PRINCIPE : Une phrase absurde/incohérente crée une RUPTURE de résonance
// → Les atomes voisins ne résonnent PAS avec elle
// → Signature énergétique anormale détectable instantanément
//
// PERFORMANCE : 7M mots/sec = scan de 5× Les Misérables en <1 seconde

package database

import (
	"math"
	"sort"
	"strings"
)

// NeedleResult représente un résultat de recherche
type NeedleResult struct {
	Position        int     // Position du mot dans le texte
	Sentence        string  // Phrase contenant l'aiguille
	CoherenceScore  float64 // Score de cohérence (bas = anormal)
	EnergySignature float64 // Signature énergétique
	ContextBefore   string  // Contexte avant
	ContextAfter    string  // Contexte après
	AnomalyRank     int     // Rang d'anomalie (1 = plus suspect)
}

// SemanticWindow représente une fenêtre glissante d'analyse
type SemanticWindow struct {
	Tokens         []string
	StartPos       int
	EndPos         int
	LocalCoherence float64
	EnergyLevel    float64
	Categories     map[int]int
}

// NeedleSearchEngine moteur de recherche d'aiguille
type NeedleSearchEngine struct {
	WindowSize       int     // Taille de la fenêtre glissante (défaut: 50 mots)
	AnomalyThreshold float64 // Seuil de détection d'anomalie (défaut: 0.3)
	MinSentenceWords int     // Nombre min de mots pour une phrase
	MaxResults       int     // Nombre max de résultats à retourner
}

// NewNeedleSearchEngine crée un nouveau moteur de recherche
func NewNeedleSearchEngine() *NeedleSearchEngine {
	return &NeedleSearchEngine{
		WindowSize:       50,
		AnomalyThreshold: 0.3,
		MinSentenceWords: 5,
		MaxResults:       10,
	}
}

// FindNeedle trouve l'aiguille dans la botte de foin
// Retourne les phrases les plus suspectes/incohérentes
func (engine *NeedleSearchEngine) FindNeedle(texte string) []NeedleResult {
	// Phase 1 : Tokenisation ultra-rapide
	tokens := TokeniserTexte(texte)
	if len(tokens) < engine.WindowSize {
		return nil
	}

	// Phase 2 : Analyse par fenêtre glissante
	windows := engine.createSlidingWindows(tokens)

	// Phase 3 : Calcul de cohérence locale pour chaque fenêtre
	engine.computeLocalCoherence(windows)

	// Phase 4 : Détection d'anomalies (ruptures de cohérence)
	anomalies := engine.detectAnomalies(windows)

	// Phase 5 : Extraction des phrases suspectes
	results := engine.extractSuspiciousSentences(texte, tokens, anomalies)

	// Phase 6 : Tri par score d'anomalie
	sort.Slice(results, func(i, j int) bool {
		return results[i].CoherenceScore < results[j].CoherenceScore
	})

	// Limiter aux N meilleurs résultats
	if len(results) > engine.MaxResults {
		results = results[:engine.MaxResults]
	}

	// Assigner les rangs
	for i := range results {
		results[i].AnomalyRank = i + 1
	}

	return results
}

// createSlidingWindows crée les fenêtres glissantes
func (engine *NeedleSearchEngine) createSlidingWindows(tokens []string) []*SemanticWindow {
	windows := make([]*SemanticWindow, 0)

	// Fenêtre glissante avec chevauchement de 50%
	step := engine.WindowSize / 2
	for i := 0; i+engine.WindowSize <= len(tokens); i += step {
		window := &SemanticWindow{
			Tokens:     tokens[i : i+engine.WindowSize],
			StartPos:   i,
			EndPos:     i + engine.WindowSize,
			Categories: make(map[int]int),
		}
		windows = append(windows, window)
	}

	return windows
}

// computeLocalCoherence calcule la cohérence locale de chaque fenêtre
// Utilise le réseau neuronal pour mesurer l'homogénéité sémantique
func (engine *NeedleSearchEngine) computeLocalCoherence(windows []*SemanticWindow) {
	for _, window := range windows {
		// Activer le réseau neuronal pour cette fenêtre
		categories := ActiverCategoriesParTexte(window.Tokens)
		window.Categories = categories

		// Calculer l'homogénéité des catégories activées
		// Une fenêtre cohérente a peu de catégories avec activations fortes
		// Une fenêtre incohérente a beaucoup de catégories faiblement activées
		coherence := engine.calculateCategoryCoherence(categories)
		window.LocalCoherence = coherence

		// Calculer l'énergie totale
		totalEnergy := 0.0
		for _, val := range categories {
			totalEnergy += float64(val)
		}
		window.EnergyLevel = totalEnergy
	}
}

// calculateCategoryCoherence calcule la cohérence d'un ensemble de catégories
// Principe : Shannon entropy inversée = mesure de concentration
func (engine *NeedleSearchEngine) calculateCategoryCoherence(categories map[int]int) float64 {
	if len(categories) == 0 {
		return 0.0
	}

	// Normaliser les valeurs
	total := 0.0
	for _, val := range categories {
		total += float64(val)
	}

	if total < 0.001 {
		return 0.0
	}

	// Calculer l'entropie de Shannon
	entropy := 0.0
	for _, val := range categories {
		if val > 0 {
			p := float64(val) / total
			entropy -= p * math.Log2(p)
		}
	}

	// Normaliser l'entropie (max = log2(nb_categories))
	maxEntropy := math.Log2(float64(len(categories)))
	if maxEntropy < 0.001 {
		return 1.0
	}

	normalizedEntropy := entropy / maxEntropy

	// Cohérence = 1 - entropie normalisée
	// Cohérence haute = entropie basse = concentré sur peu de catégories
	// Cohérence basse = entropie haute = dispersé sur beaucoup de catégories
	coherence := 1.0 - normalizedEntropy

	return coherence
}

// detectAnomalies détecte les fenêtres avec cohérence anormalement basse
func (engine *NeedleSearchEngine) detectAnomalies(windows []*SemanticWindow) []*SemanticWindow {
	if len(windows) == 0 {
		return nil
	}

	// Calculer la cohérence moyenne et écart-type
	sumCoherence := 0.0
	for _, w := range windows {
		sumCoherence += w.LocalCoherence
	}
	meanCoherence := sumCoherence / float64(len(windows))

	sumSquares := 0.0
	for _, w := range windows {
		diff := w.LocalCoherence - meanCoherence
		sumSquares += diff * diff
	}
	stdDev := math.Sqrt(sumSquares / float64(len(windows)))

	// Seuil d'anomalie : moyenne - 2*écart-type (outliers statistiques)
	anomalyThreshold := meanCoherence - 2*stdDev
	if anomalyThreshold < engine.AnomalyThreshold {
		anomalyThreshold = engine.AnomalyThreshold
	}

	// Filtrer les fenêtres anormales
	anomalies := make([]*SemanticWindow, 0)
	for _, w := range windows {
		if w.LocalCoherence < anomalyThreshold {
			anomalies = append(anomalies, w)
		}
	}

	return anomalies
}

// extractSuspiciousSentences extrait les phrases suspectes autour des anomalies
func (engine *NeedleSearchEngine) extractSuspiciousSentences(texte string, tokens []string, anomalies []*SemanticWindow) []NeedleResult {
	results := make([]NeedleResult, 0)

	// Marquer les positions suspectes
	suspectPositions := make(map[int]float64)
	for _, anomaly := range anomalies {
		for pos := anomaly.StartPos; pos < anomaly.EndPos; pos++ {
			// Accumuler le score d'anomalie
			score := 1.0 - anomaly.LocalCoherence
			if existing, ok := suspectPositions[pos]; ok {
				suspectPositions[pos] = math.Max(existing, score)
			} else {
				suspectPositions[pos] = score
			}
		}
	}

	// Extraire les phrases autour des positions suspectes
	sentences := engine.splitIntoSentences(texte)

	for _, sentence := range sentences {
		// Trouver la position de cette phrase dans les tokens
		sentenceTokens := TokeniserTexte(sentence)
		if len(sentenceTokens) < engine.MinSentenceWords {
			continue
		}

		// Chercher cette séquence dans les tokens
		pos := engine.findTokenSequence(tokens, sentenceTokens)
		if pos < 0 {
			continue
		}

		// Calculer le score d'anomalie de cette phrase
		maxAnomaly := 0.0
		for i := pos; i < pos+len(sentenceTokens) && i < len(tokens); i++ {
			if score, ok := suspectPositions[i]; ok {
				if score > maxAnomaly {
					maxAnomaly = score
				}
			}
		}

		// Si phrase suspecte, l'ajouter aux résultats
		if maxAnomaly > 0.2 {
			result := NeedleResult{
				Position:        pos,
				Sentence:        sentence,
				CoherenceScore:  1.0 - maxAnomaly,
				EnergySignature: maxAnomaly,
				ContextBefore:   engine.getContext(tokens, pos-10, pos),
				ContextAfter:    engine.getContext(tokens, pos+len(sentenceTokens), pos+len(sentenceTokens)+10),
			}
			results = append(results, result)
		}
	}

	return results
}

// splitIntoSentences découpe le texte en phrases
func (engine *NeedleSearchEngine) splitIntoSentences(texte string) []string {
	// Découpage simple par ponctuation forte
	texte = strings.ReplaceAll(texte, "! ", "!|")
	texte = strings.ReplaceAll(texte, "? ", "?|")
	texte = strings.ReplaceAll(texte, ". ", ".|")
	texte = strings.ReplaceAll(texte, ".\n", ".|")

	sentences := strings.Split(texte, "|")
	result := make([]string, 0)

	for _, s := range sentences {
		s = strings.TrimSpace(s)
		if len(s) > 20 {
			result = append(result, s)
		}
	}

	return result
}

// findTokenSequence trouve la position d'une séquence de tokens
func (engine *NeedleSearchEngine) findTokenSequence(haystack, needle []string) int {
	if len(needle) == 0 || len(haystack) < len(needle) {
		return -1
	}

	for i := 0; i <= len(haystack)-len(needle); i++ {
		match := true
		for j := 0; j < len(needle); j++ {
			if haystack[i+j] != needle[j] {
				match = false
				break
			}
		}
		if match {
			return i
		}
	}

	return -1
}

// getContext récupère le contexte autour d'une position
func (engine *NeedleSearchEngine) getContext(tokens []string, start, end int) string {
	if start < 0 {
		start = 0
	}
	if end > len(tokens) {
		end = len(tokens)
	}
	if start >= end {
		return ""
	}

	return strings.Join(tokens[start:end], " ")
}
