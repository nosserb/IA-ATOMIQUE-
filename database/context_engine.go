// Package database - Context Engine
// Moteur de contexte enrichi pour analyse multi-niveau
// Résout le problème du manque de contexte dans les benchmarks

package database

import (
	"math"
	"strings"
)

// ContextEngine moteur de contexte enrichi
type ContextEngine struct {
	NGrams           map[string]float64  // N-grammes contextuels
	ConceptGraph     map[string][]string // Graphe de concepts
	SemanticHistory  []map[int]int       // Historique sémantique
	ContextWindow    int                 // Taille de la fenêtre
	DeepAnalysisMode bool                // Mode analyse profonde
}

// NewContextEngine crée un nouveau moteur de contexte
func NewContextEngine() *ContextEngine {
	return &ContextEngine{
		NGrams:           make(map[string]float64),
		ConceptGraph:     initializeConceptGraph(),
		SemanticHistory:  make([]map[int]int, 0),
		ContextWindow:    5,
		DeepAnalysisMode: true,
	}
}

// initializeConceptGraph initialise le graphe de concepts
func initializeConceptGraph() map[string][]string {
	// Utiliser le graphe étendu
	return initializeExtendedConceptGraph()
}

// EnrichContext enrichit le contexte par analyse multi-niveau
func (ce *ContextEngine) EnrichContext(text string, categories map[int]int) *EnrichedContext {
	tokens := TokeniserTexte(text)

	enriched := &EnrichedContext{
		OriginalText:    text,
		Tokens:          tokens,
		Categories:      categories,
		BiGrams:         ce.extractBiGrams(tokens),
		TriGrams:        ce.extractTriGrams(tokens),
		ConceptClusters: ce.buildConceptClusters(tokens),
		SemanticDepth:   ce.calculateSemanticDepth(tokens),
		ContextualScore: 0.0,
	}

	// Calculer score contextuel global
	enriched.ContextualScore = ce.calculateContextScore(enriched)

	return enriched
}

// EnrichedContext contexte enrichi
type EnrichedContext struct {
	OriginalText    string
	Tokens          []string
	Categories      map[int]int
	BiGrams         []string
	TriGrams        []string
	ConceptClusters [][]string
	SemanticDepth   float64
	ContextualScore float64
}

// extractBiGrams extrait les bi-grammes
func (ce *ContextEngine) extractBiGrams(tokens []string) []string {
	bigrams := make([]string, 0)
	for i := 0; i < len(tokens)-1; i++ {
		bigram := strings.ToLower(tokens[i] + " " + tokens[i+1])
		bigrams = append(bigrams, bigram)

		// Stocker dans NGrams avec score
		if _, exists := ce.NGrams[bigram]; !exists {
			ce.NGrams[bigram] = 1.0
		} else {
			ce.NGrams[bigram] += 0.1
		}
	}
	return bigrams
}

// extractTriGrams extrait les tri-grammes
func (ce *ContextEngine) extractTriGrams(tokens []string) []string {
	trigrams := make([]string, 0)
	for i := 0; i < len(tokens)-2; i++ {
		trigram := strings.ToLower(tokens[i] + " " + tokens[i+1] + " " + tokens[i+2])
		trigrams = append(trigrams, trigram)

		// Stocker avec score plus élevé (plus de contexte)
		if _, exists := ce.NGrams[trigram]; !exists {
			ce.NGrams[trigram] = 1.5
		} else {
			ce.NGrams[trigram] += 0.15
		}
	}
	return trigrams
}

// buildConceptClusters construit des clusters de concepts liés
func (ce *ContextEngine) buildConceptClusters(tokens []string) [][]string {
	clusters := make([][]string, 0)

	for _, token := range tokens {
		tokenLower := strings.ToLower(token)
		if relatedConcepts, exists := ce.ConceptGraph[tokenLower]; exists {
			cluster := []string{tokenLower}
			cluster = append(cluster, relatedConcepts...)
			clusters = append(clusters, cluster)
		}
	}

	return clusters
}

// calculateSemanticDepth calcule la profondeur sémantique
func (ce *ContextEngine) calculateSemanticDepth(tokens []string) float64 {
	if len(tokens) == 0 {
		return 0.0
	}

	// Plus de tokens uniques = plus de profondeur
	uniqueTokens := make(map[string]bool)
	for _, token := range tokens {
		uniqueTokens[strings.ToLower(token)] = true
	}

	diversity := float64(len(uniqueTokens)) / float64(len(tokens))

	// Présence dans le graphe de concepts = plus de profondeur
	conceptCount := 0
	for token := range uniqueTokens {
		if _, exists := ce.ConceptGraph[token]; exists {
			conceptCount++
		}
	}

	conceptDensity := float64(conceptCount) / float64(len(uniqueTokens))

	return (diversity*0.4 + conceptDensity*0.6)
}

// calculateContextScore score contextuel global
func (ce *ContextEngine) calculateContextScore(enriched *EnrichedContext) float64 {
	score := 0.0

	// 1. Score basé sur n-grammes connus
	ngramScore := 0.0
	for _, bigram := range enriched.BiGrams {
		if weight, exists := ce.NGrams[bigram]; exists {
			ngramScore += weight
		}
	}
	for _, trigram := range enriched.TriGrams {
		if weight, exists := ce.NGrams[trigram]; exists {
			ngramScore += weight * 1.5
		}
	}
	ngramScore = math.Min(ngramScore/10.0, 1.0)

	// 2. Score basé sur clusters de concepts
	clusterScore := float64(len(enriched.ConceptClusters)) / 10.0
	clusterScore = math.Min(clusterScore, 1.0)

	// 3. Profondeur sémantique
	depthScore := enriched.SemanticDepth

	// Score final pondéré
	score = (ngramScore*0.3 + clusterScore*0.3 + depthScore*0.4)

	return score
}

// CompareEnrichedContexts compare deux contextes enrichis avec pondération améliorée
func (ce *ContextEngine) CompareEnrichedContexts(ctx1, ctx2 *EnrichedContext) float64 {
	if ctx1 == nil || ctx2 == nil {
		return 0.0
	}

	score := 0.0

	// 1. Overlap de bi-grammes
	bigramOverlap := ce.calculateNGramOverlap(ctx1.BiGrams, ctx2.BiGrams)

	// 2. Overlap de tri-grammes (plus important)
	trigramOverlap := ce.calculateNGramOverlap(ctx1.TriGrams, ctx2.TriGrams)

	// 3. Overlap de concepts (critique pour la compréhension)
	conceptOverlap := ce.calculateConceptOverlap(ctx1.ConceptClusters, ctx2.ConceptClusters)

	// 4. Similarité de profondeur sémantique
	depthSimilarity := 1.0 - math.Abs(ctx1.SemanticDepth-ctx2.SemanticDepth)

	// 5. NOUVEAU: Overlap de tokens simples (pour contextes courts)
	tokenOverlap := ce.calculateTokenOverlap(ctx1.Tokens, ctx2.Tokens)

	// 6. NOUVEAU: Similarité de catégories
	categoryOverlap := ce.calculateCategoryOverlap(ctx1.Categories, ctx2.Categories)

	// Score combiné avec bonus pour contextes courts
	if len(ctx1.Tokens) < 5 || len(ctx2.Tokens) < 5 {
		// Pour contextes courts: plus de poids aux tokens et catégories
		score = (bigramOverlap*0.15 +
			trigramOverlap*0.20 +
			conceptOverlap*0.25 +
			depthSimilarity*0.10 +
			tokenOverlap*0.20 +
			categoryOverlap*0.10)
	} else {
		// Pour contextes longs: plus de poids aux n-grammes
		score = (bigramOverlap*0.20 +
			trigramOverlap*0.35 +
			conceptOverlap*0.30 +
			depthSimilarity*0.15)
	}

	return score
}

// calculateNGramOverlap calcule l'overlap de n-grammes
func (ce *ContextEngine) calculateNGramOverlap(ngrams1, ngrams2 []string) float64 {
	if len(ngrams1) == 0 || len(ngrams2) == 0 {
		return 0.0
	}

	ngramSet := make(map[string]bool)
	for _, ng := range ngrams1 {
		ngramSet[ng] = true
	}

	overlapCount := 0
	for _, ng := range ngrams2 {
		if ngramSet[ng] {
			overlapCount++
		}
	}

	return float64(overlapCount) / float64(len(ngrams2))
}

// calculateConceptOverlap calcule l'overlap de concepts
func (ce *ContextEngine) calculateConceptOverlap(clusters1, clusters2 [][]string) float64 {
	if len(clusters1) == 0 || len(clusters2) == 0 {
		return 0.0
	}

	// Aplatir les clusters en set
	conceptSet1 := make(map[string]bool)
	for _, cluster := range clusters1 {
		for _, concept := range cluster {
			conceptSet1[concept] = true
		}
	}

	conceptSet2 := make(map[string]bool)
	for _, cluster := range clusters2 {
		for _, concept := range cluster {
			conceptSet2[concept] = true
		}
	}

	// Calculer intersection
	overlapCount := 0
	for concept := range conceptSet1 {
		if conceptSet2[concept] {
			overlapCount++
		}
	}

	// Jaccard similarity
	unionSize := len(conceptSet1) + len(conceptSet2) - overlapCount
	if unionSize == 0 {
		return 0.0
	}

	return float64(overlapCount) / float64(unionSize)
}

// calculateTokenOverlap calcule l'overlap de tokens simples
func (ce *ContextEngine) calculateTokenOverlap(tokens1, tokens2 []string) float64 {
	if len(tokens1) == 0 || len(tokens2) == 0 {
		return 0.0
	}

	tokenSet := make(map[string]bool)
	for _, token := range tokens1 {
		tokenSet[token] = true
	}

	overlapCount := 0
	for _, token := range tokens2 {
		if tokenSet[token] {
			overlapCount++
		}
	}

	return float64(overlapCount) / float64(len(tokens2))
}

// calculateCategoryOverlap calcule l'overlap de catégories
func (ce *ContextEngine) calculateCategoryOverlap(cat1, cat2 map[int]int) float64 {
	if len(cat1) == 0 || len(cat2) == 0 {
		return 0.0
	}

	// Calcul de similarité cosinus pour catégories
	var dotProduct, norm1, norm2 float64

	for cat, count1 := range cat1 {
		if count2, exists := cat2[cat]; exists {
			dotProduct += float64(count1 * count2)
		}
		norm1 += float64(count1 * count1)
	}

	for _, count2 := range cat2 {
		norm2 += float64(count2 * count2)
	}

	if norm1 == 0 || norm2 == 0 {
		return 0.0
	}

	return dotProduct / (math.Sqrt(norm1) * math.Sqrt(norm2))
}

// UpdateHistory met à jour l'historique sémantique
func (ce *ContextEngine) UpdateHistory(categories map[int]int) {
	ce.SemanticHistory = append(ce.SemanticHistory, categories)

	// Garder seulement les N dernières entrées
	maxHistory := 20
	if len(ce.SemanticHistory) > maxHistory {
		ce.SemanticHistory = ce.SemanticHistory[len(ce.SemanticHistory)-maxHistory:]
	}
}

// GetHistoricalContext obtient le contexte historique moyen
func (ce *ContextEngine) GetHistoricalContext() map[int]int {
	if len(ce.SemanticHistory) == 0 {
		return make(map[int]int)
	}

	// Fusionner les catégories historiques avec décroissance temporelle
	historicalContext := make(map[int]float64)

	for i, categories := range ce.SemanticHistory {
		// Plus récent = plus de poids
		weight := float64(i+1) / float64(len(ce.SemanticHistory))

		for cat, count := range categories {
			historicalContext[cat] += float64(count) * weight
		}
	}

	// Convertir en int
	result := make(map[int]int)
	for cat, val := range historicalContext {
		result[cat] = int(math.Round(val))
	}

	return result
}

// ExpandContextWithHistory enrichit avec l'historique
func (ce *ContextEngine) ExpandContextWithHistory(categories map[int]int, weight float64) map[int]int {
	historical := ce.GetHistoricalContext()
	expanded := make(map[int]int)

	// Copier categories actuelles
	for cat, count := range categories {
		expanded[cat] = count
	}

	// Ajouter contexte historique avec pondération
	for cat, count := range historical {
		historicalContribution := int(float64(count) * weight)
		if existing, exists := expanded[cat]; exists {
			expanded[cat] = existing + historicalContribution
		} else {
			expanded[cat] = historicalContribution
		}
	}

	return expanded
}
