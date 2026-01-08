package database

import (
	"math"
	"sort"
	"strings"
)

// ============================================================================
// DOMAIN CONSTRAINT SYSTEM
// ============================================================================
// Empêche hallucinations EN AMONT par contraintes sémantiques strictes
// Formellement: ∀s∈R, ∃c∈VT tel que Sim(s,c)≥α

// DomainSpace représente l'espace conceptuel du texte source
type DomainSpace struct {
	CoreConcepts      []string           // VT = concepts dominants du texte
	TechTerms         map[string]float64 // Termes techniques avec scores
	ForbiddenConcepts []string           // BT = concepts interdits (hors domaine)
	TechDensity       float64            // Densité technique du texte source
	DomainMode        string             // "technical", "narrative", "social" etc.
}

// ExtractDomainConcepts construit VT (espace conceptuel source)
// C(T) = ensemble des concepts importants dans T
func ExtractDomainConcepts(sourceText string) DomainSpace {
	result := DomainSpace{
		TechTerms: make(map[string]float64),
	}

	// 1. Extraction TF-IDF technique
	result.TechTerms = ExtractTechnicalTermsWithScores(sourceText)

	// 2. Concepts principaux (mots-clés avec haute fréquence)
	concepts := ExtractPrimaryConceptsFromText(sourceText)
	result.CoreConcepts = concepts

	// 3. Densité technique
	result.TechDensity = CalculateTechDensity(sourceText, result.TechTerms)

	// 4. Détecter mode de domaine
	result.DomainMode = DetectDomainMode(sourceText, result.TechDensity)

	// 5. Construire BT (concepts interdits = top concepts d'autres domaines)
	result.ForbiddenConcepts = BuildForbiddenConceptList(sourceText, result.DomainMode)

	return result
}

// ExtractTechnicalTermsWithScores extrait termes techniques avec scores TF-IDF
func ExtractTechnicalTermsWithScores(text string) map[string]float64 {
	terms := make(map[string]float64)

	// Patterns techniques reconnus
	technicalPatterns := map[string]float64{
		// Ornithologie (pour input.txt)
		"mésange": 1.0, "oiseau": 0.8, "plumage": 0.9, "huppe": 1.0,
		"nid": 0.7, "reproduction": 0.8, "nidification": 0.9,
		"alimentation": 0.7, "insectes": 0.8, "graines": 0.7,
		"habitat": 0.8, "forêt": 0.6, "conifères": 0.8,
		"territoire": 0.7, "ponte": 0.9, "fledgling": 0.8,
		"prédateur": 0.8, "migration": 0.7, "hibernation": 0.7,
		"espèce": 0.8, "taxinomie": 0.9, "sous-espèce": 0.9,
		"paridé": 1.0, "lophophanes": 1.0, "cristatus": 1.0,

		// IA-ATOMIQUE patterns
		"atome": 1.0, "réseau": 0.9, "résonance": 0.95,
		"asynchrone": 1.0, "moteur": 0.7, "inférence": 0.9,
		"distribuée": 0.85, "interaction": 0.8, "émergence": 0.9,

		// Patterns génériques
		"étude": 0.5, "données": 0.6, "analyse": 0.6,
	}

	textLower := strings.ToLower(text)

	for term, baseScore := range technicalPatterns {
		if strings.Contains(textLower, term) {
			// Calculer TF-IDF
			freq := float64(strings.Count(textLower, term))
			// IDF simple (terme technique = haute valeur)
			tfidf := baseScore * math.Log(freq+1)
			terms[term] = tfidf
		}
	}

	return terms
}

// ExtractPrimaryConceptsFromText extrait concepts principaux
func ExtractPrimaryConceptsFromText(text string) []string {
	words := strings.Fields(text)
	freq := make(map[string]int)

	for _, word := range words {
		// Normaliser
		w := strings.ToLower(word)
		w = strings.TrimRight(w, ".,;:!?()[]{}")

		// Skip stopwords
		if IsStopwordFR(w) || len(w) < 4 {
			continue
		}

		freq[w]++
	}

	// Trier par fréquence
	type kv struct {
		word string
		freq int
	}

	var sorted []kv
	for w, f := range freq {
		if f >= 2 { // Au moins 2 occurrences
			sorted = append(sorted, kv{w, f})
		}
	}

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].freq > sorted[j].freq
	})

	// Top 30 concepts
	var concepts []string
	for i := 0; i < 30 && i < len(sorted); i++ {
		concepts = append(concepts, sorted[i].word)
	}

	return concepts
}

// CalculateTechDensity = termes techniques / mots totaux
func CalculateTechDensity(text string, techTerms map[string]float64) float64 {
	if len(techTerms) == 0 {
		return 0.0
	}

	words := strings.Fields(text)
	techCount := 0
	textLower := strings.ToLower(text)

	for term := range techTerms {
		techCount += strings.Count(textLower, term)
	}

	return float64(techCount) / float64(len(words))
}

// DetectDomainMode détecte le mode de domaine (avec sous-catégories)
func DetectDomainMode(text string, techDensity float64) string {
	textLower := strings.ToLower(text)

	// Scorer chaque domaine
	ornithologyScore := strings.Count(textLower, "oiseau") +
		strings.Count(textLower, "mésange")*2 +
		strings.Count(textLower, "plumage") +
		strings.Count(textLower, "huppe")

	aiScore := strings.Count(textLower, "atome")*2 +
		strings.Count(textLower, "réseau")*2 +
		strings.Count(textLower, "résonnance") +
		strings.Count(textLower, "asynchrone")*2 +
		strings.Count(textLower, "distributed")*2 +
		strings.Count(textLower, "concept")*1 +
		strings.Count(textLower, "architecture")*1

	socialScore := strings.Count(textLower, "violence") +
		strings.Count(textLower, "économique") +
		strings.Count(textLower, "société") +
		strings.Count(textLower, "oppression")

	// Domaine vainqueur
	if ornithologyScore > aiScore && ornithologyScore > socialScore {
		return "ornithology"
	} else if aiScore > socialScore {
		// IA détectée - distinguer hard vs conceptual
		// Si beaucoup de code/équations -> hard, sinon conceptual
		codePatterns := strings.Count(textLower, "code") +
			strings.Count(textLower, "function") +
			strings.Count(textLower, "return") +
			strings.Count(textLower, "if ") +
			strings.Count(textLower, "loop")

		equationPatterns := strings.Count(textLower, "∀") +
			strings.Count(textLower, "∃") +
			strings.Count(textLower, "⊆") +
			strings.Count(textLower, "arg max") +
			strings.Count(textLower, "λ")*2

		if codePatterns > 5 || equationPatterns > 3 || techDensity > 0.25 {
			return "technical-ai-hard"
		} else {
			return "technical-ai-conceptual" // ← NOUVEAU
		}
	} else if techDensity > 0.15 {
		return "technical"
	} else if socialScore > 0 {
		return "social"
	}
	return "narrative"
}

// BuildForbiddenConceptList construit BT (concepts interdits)
func BuildForbiddenConceptList(sourceText string, domainMode string) []string {
	forbidden := []string{}

	// Concepts à TOUJOURS interdire
	alwaysForbidden := []string{
		"hallucination", "invention", "fabulation",
	}
	forbidden = append(forbidden, alwaysForbidden...)

	// Concepts interdits par domaine
	switch domainMode {
	case "ornithology":
		// Si domaine = ornithologie, interdire concepts sociaux
		forbidden = append(forbidden,
			"violence", "économique", "sociologie", "oppression",
			"classe", "inégalité", "politique", "pauvreté")
	case "technical-ai":
		// Si domaine = IA, interdire concepts ornithologiques
		forbidden = append(forbidden,
			"oiseau", "plumage", "huppe", "nid", "vol",
			"mésange", "reproduction", "nidification")
	case "social":
		// Si domaine = social, interdire concepts ornithologiques
		forbidden = append(forbidden,
			"oiseau", "plumage", "huppe", "nid", "vol",
			"mésange", "reproduction", "nidification")
	}

	return forbidden
}

// ValidateSentenceInDomain vérifie si une phrase est dans VT et pas dans BT
// ∀s∈R, ∃c∈VT tel que Sim(s,c)≥α ET s ∉ BT
func (ds *DomainSpace) ValidateSentenceInDomain(sentence string, alphaThreshold float64) (valid bool, reason string) {
	// 1. Vérifier pas dans BT
	sentenceLower := strings.ToLower(sentence)

	for _, forbidden := range ds.ForbiddenConcepts {
		if strings.Contains(sentenceLower, forbidden) {
			return false, "concept_forbidden"
		}
	}

	// 2. Pour textes conceptuels-AI: utiliser similarité au prototype
	if ds.DomainMode == "technical-ai-conceptual" {
		// Calcul: cos(phrase→, prototype_domaine→) ≥ 0.65
		prototypeScore := CalculateSentenceDomainSimilarity(sentence, ds)
		if prototypeScore >= 0.65 {
			return true, "valid_conceptual"
		} else if prototypeScore >= 0.45 {
			// Borderline - accepter si au moins un concept principal
			for _, concept := range ds.CoreConcepts {
				if strings.Contains(sentenceLower, concept) {
					return true, "valid_conceptual_borderline"
				}
			}
		}
		return false, "low_domain_similarity"
	}

	// 3. Pour textes techniques: vérifier au moins un terme technique
	techTermCount := 0
	for term := range ds.TechTerms {
		if strings.Contains(sentenceLower, term) {
			techTermCount++
		}
	}

	if techTermCount == 0 && ds.TechDensity > 0.2 {
		return false, "no_tech_terms"
	}

	// 4. Vérifier alignment avec concepts principaux
	matchedConcepts := 0
	for _, concept := range ds.CoreConcepts {
		if strings.Contains(sentenceLower, concept) {
			matchedConcepts++
		}
	}

	if matchedConcepts > 0 {
		return true, "valid"
	}

	// Si tech dense et la phrase parle de tech, c'est ok
	if ds.TechDensity > 0.2 && techTermCount > 0 {
		return true, "valid_tech"
	}

	return false, "low_alignment"
}

// CalculateMultiObjectiveScore = λ1*Coverage + λ2*TechDensity - λ3*Drift
func CalculateMultiObjectiveScore(
	summary string,
	sourceText string,
	lambda1 float64, // weight Coverage
	lambda2 float64, // weight TechDensity
	lambda3 float64, // weight Drift penalty
) float64 {

	// 1. Coverage(R,T) = |Units(R) ∩ Units(T)| / |Units(R)|
	coverage := CalculateFidelity(summary, sourceText)

	// 2. TechDensity(R)
	techTerms := ExtractTechnicalTermsWithScores(sourceText)
	techDensity := CalculateTechDensity(summary, techTerms)

	// 3. Drift(R,T) = 1 - cos(R→, T→)
	drift := CalculateDrift(summary, sourceText)

	// Score multi-objectif
	score := lambda1*coverage + lambda2*techDensity - lambda3*drift

	return score
}

// CalculateDrift = 1 - cosine_similarity(summary, source)
func CalculateDrift(summary string, sourceText string) float64 {
	// Vecteurs simples: occurrence des termes clés
	terms := ExtractTechnicalTermsWithScores(sourceText)

	summaryVec := make(map[string]float64)
	sourceVec := make(map[string]float64)

	summaryLower := strings.ToLower(summary)
	sourceLower := strings.ToLower(sourceText)

	// Compter occurrences
	for term := range terms {
		summaryVec[term] = float64(strings.Count(summaryLower, term))
		sourceVec[term] = float64(strings.Count(sourceLower, term))
	}

	// Cosine similarity
	dotProduct := 0.0
	normSum := 0.0
	normSource := 0.0

	for term := range terms {
		dotProduct += summaryVec[term] * sourceVec[term]
		normSum += summaryVec[term] * summaryVec[term]
		normSource += sourceVec[term] * sourceVec[term]
	}

	if normSum == 0 || normSource == 0 {
		return 1.0 // Drift maximal
	}

	cosineSim := dotProduct / (math.Sqrt(normSum) * math.Sqrt(normSource))
	drift := 1.0 - cosineSim

	return math.Max(0, math.Min(1, drift))
}

// IsStopwordFR vérifie si un mot est un stopword français
func IsStopwordFR(word string) bool {
	stopwords := map[string]bool{
		"le": true, "la": true, "les": true, "de": true, "des": true,
		"un": true, "une": true, "et": true, "ou": true, "dans": true,
		"pour": true, "par": true, "est": true, "à": true, "son": true,
		"sa": true, "ses": true, "ce": true, "que": true, "qui": true,
		"avec": true, "du": true, "en": true, "on": true, "il": true,
		"elle": true, "nous": true, "vous": true, "ils": true, "elles": true,
		"au": true, "aux": true, "mais": true, "peut": true, "se": true,
		"tu": true, "te": true, "moi": true, "toi": true, "lui": true,
		"leur": true, "même": true, "bien": true, "très": true, "plus": true,
	}
	return stopwords[strings.ToLower(word)]
}
