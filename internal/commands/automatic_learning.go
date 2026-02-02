package commands

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/nosserb/IA-ATOMIQUE-/database"
)

// ============================================================================
// SYSTÈME D'APPRENTISSAGE AUTOMATIQUE À PARTIR DE TEXTES
// ============================================================================
// Apprend automatiquement des connaissances factuelles en lisant des corpus
// Combine 3 approches:
// 1. Extraction de patterns (dates, relations causales, lieux)
// 2. Co-occurrences (associations statistiques)
// 3. Entraînement atomique (clusters émergents)

// KnowledgeBase - Base de connaissances apprises automatiquement
type KnowledgeBase struct {
	// Faits extraits par patterns
	DateFacts       map[string][]string // "1815" -> ["Waterloo", "défaite Napoléon"]
	CausalFacts     map[string][]string // "Waterloo" -> ["fin empire", "exil"]
	LocationFacts   map[string][]string // "Paris" -> ["capitale", "France"]
	DefinitionFacts map[string]string   // "coeur" -> "organe muscle pompe sang"

	// Co-occurrences (distance < 5 mots)
	CoOccurrences map[string]map[string]int // "Napoléon" -> {"Waterloo": 15, "1815": 20}

	// Réseau atomique entraîné
	ConceptNetwork *database.AtomicNetwork
	ConceptMapping map[string]int // "Napoléon" -> atom_id

	// Statistiques
	TotalTextsProcessed int
	TotalWordsProcessed int
	TotalFactsExtracted int

	mutex sync.RWMutex
}

// NewKnowledgeBase crée une nouvelle base de connaissances
func NewKnowledgeBase() *KnowledgeBase {
	kb := &KnowledgeBase{
		DateFacts:       make(map[string][]string),
		CausalFacts:     make(map[string][]string),
		LocationFacts:   make(map[string][]string),
		DefinitionFacts: make(map[string]string),
		CoOccurrences:   make(map[string]map[string]int),
		ConceptMapping:  make(map[string]int),
	}

	// Réseau atomique pour apprentissage émergent (500 atomes)
	kb.ConceptNetwork = database.NewAtomicNetwork(500)

	return kb
}

// LearnFromTextFile apprend à partir d'un fichier texte
func (kb *KnowledgeBase) LearnFromTextFile(filepath string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("erreur ouverture fichier: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var fullText strings.Builder

	// Lire tout le texte
	for scanner.Scan() {
		line := scanner.Text()
		fullText.WriteString(line)
		fullText.WriteString(" ")
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("erreur lecture fichier: %v", err)
	}

	text := fullText.String()
	kb.LearnFromText(text)
	kb.TotalTextsProcessed++

	return nil
}

// LearnFromText apprend à partir d'un texte brut
func (kb *KnowledgeBase) LearnFromText(text string) {
	kb.mutex.Lock()
	defer kb.mutex.Unlock()

	// 1. Extraction de patterns
	kb.extractDatePatterns(text)
	kb.extractCausalPatterns(text)
	kb.extractLocationPatterns(text)
	kb.extractDefinitionPatterns(text)

	// 2. Co-occurrences
	kb.buildCoOccurrences(text)

	// 3. Entraînement atomique
	kb.trainAtomicNetwork(text)

	// Compter les mots
	words := strings.Fields(text)
	kb.TotalWordsProcessed += len(words)
}

// extractDatePatterns extrait les faits avec dates
// Patterns: "X en YYYY", "en YYYY X", "X (YYYY)"
func (kb *KnowledgeBase) extractDatePatterns(text string) {
	// Pattern: "bataille de Waterloo en 1815"
	pattern1 := regexp.MustCompile(`([A-ZÉÈÊÀ][a-zéèêàùç\s]+(?:de|d')[A-ZÉÈÊÀ][a-zéèêàùç\s]+)\s+en\s+(\d{4})`)
	matches := pattern1.FindAllStringSubmatch(text, -1)
	for _, match := range matches {
		if len(match) >= 3 {
			event := strings.TrimSpace(match[1])
			year := match[2]
			kb.DateFacts[year] = append(kb.DateFacts[year], event)
			kb.TotalFactsExtracted++
		}
	}

	// Pattern: "en 1815, défaite de Napoléon"
	pattern2 := regexp.MustCompile(`en\s+(\d{4})[,\s]+([a-zéèêàùç\s]+(?:de|d')[A-ZÉÈÊÀ][a-zéèêàùç\s]+)`)
	matches = pattern2.FindAllStringSubmatch(text, -1)
	for _, match := range matches {
		if len(match) >= 3 {
			year := match[1]
			event := strings.TrimSpace(match[2])
			kb.DateFacts[year] = append(kb.DateFacts[year], event)
			kb.TotalFactsExtracted++
		}
	}

	// Pattern: "Napoléon (1769-1821)"
	pattern3 := regexp.MustCompile(`([A-ZÉÈÊÀ][a-zéèêàùç\s]+)\s*\((\d{4})-(\d{4})\)`)
	matches = pattern3.FindAllStringSubmatch(text, -1)
	for _, match := range matches {
		if len(match) >= 4 {
			person := strings.TrimSpace(match[1])
			birthYear := match[2]
			deathYear := match[3]
			kb.DateFacts[birthYear] = append(kb.DateFacts[birthYear], "naissance de "+person)
			kb.DateFacts[deathYear] = append(kb.DateFacts[deathYear], "mort de "+person)
			kb.TotalFactsExtracted += 2
		}
	}
}

// extractCausalPatterns extrait les relations causales
// Patterns: "X a causé Y", "X provoque Y", "à cause de X, Y"
func (kb *KnowledgeBase) extractCausalPatterns(text string) {
	// Pattern: "X a causé Y"
	pattern1 := regexp.MustCompile(`([A-ZÉÈÊÀ][a-zéèêàùç\s]+)\s+a\s+causé\s+([a-zéèêàùç\s]+)`)
	matches := pattern1.FindAllStringSubmatch(text, -1)
	for _, match := range matches {
		if len(match) >= 3 {
			cause := strings.TrimSpace(match[1])
			effect := strings.TrimSpace(match[2])
			kb.CausalFacts[cause] = append(kb.CausalFacts[cause], effect)
			kb.TotalFactsExtracted++
		}
	}

	// Pattern: "X provoque Y"
	pattern2 := regexp.MustCompile(`([A-ZÉÈÊÀ][a-zéèêàùç\s]+)\s+provoque\s+([a-zéèêàùç\s]+)`)
	matches = pattern2.FindAllStringSubmatch(text, -1)
	for _, match := range matches {
		if len(match) >= 3 {
			cause := strings.TrimSpace(match[1])
			effect := strings.TrimSpace(match[2])
			kb.CausalFacts[cause] = append(kb.CausalFacts[cause], effect)
			kb.TotalFactsExtracted++
		}
	}

	// Pattern: "à cause de X, Y"
	pattern3 := regexp.MustCompile(`à cause de\s+([a-zéèêàùç\s]+),\s+([a-zéèêàùç\s]+)`)
	matches = pattern3.FindAllStringSubmatch(text, -1)
	for _, match := range matches {
		if len(match) >= 3 {
			cause := strings.TrimSpace(match[1])
			effect := strings.TrimSpace(match[2])
			kb.CausalFacts[cause] = append(kb.CausalFacts[cause], effect)
			kb.TotalFactsExtracted++
		}
	}
}

// extractLocationPatterns extrait les faits géographiques
// Patterns: "X est situé en Y", "X capitale de Y", "X se trouve en Y"
func (kb *KnowledgeBase) extractLocationPatterns(text string) {
	// Pattern: "Paris capitale de la France"
	pattern1 := regexp.MustCompile(`([A-ZÉÈÊÀ][a-zéèêàùç]+)\s+capitale\s+(?:de|du|d')\s+(?:la|le|l')\s*([A-ZÉÈÊÀ][a-zéèêàùç]+)`)
	matches := pattern1.FindAllStringSubmatch(text, -1)
	for _, match := range matches {
		if len(match) >= 3 {
			city := strings.TrimSpace(match[1])
			country := strings.TrimSpace(match[2])
			kb.LocationFacts[city] = append(kb.LocationFacts[city], "capitale de "+country)
			kb.TotalFactsExtracted++
		}
	}

	// Pattern: "X est situé en Y"
	pattern2 := regexp.MustCompile(`([A-ZÉÈÊÀ][a-zéèêàùç\s]+)\s+(?:est situé|se trouve)\s+en\s+([A-ZÉÈÊÀ][a-zéèêàùç\s]+)`)
	matches = pattern2.FindAllStringSubmatch(text, -1)
	for _, match := range matches {
		if len(match) >= 3 {
			place := strings.TrimSpace(match[1])
			location := strings.TrimSpace(match[2])
			kb.LocationFacts[place] = append(kb.LocationFacts[place], "situé en "+location)
			kb.TotalFactsExtracted++
		}
	}
}

// extractDefinitionPatterns extrait les définitions
// Patterns: "X est Y", "X: Y", "X, Y,"
func (kb *KnowledgeBase) extractDefinitionPatterns(text string) {
	// Pattern: "Le cœur est un organe musculaire"
	pattern1 := regexp.MustCompile(`(?:Le|La|L')\s+([a-zéèêàùçœ]+)\s+est\s+un(?:e)?\s+([a-zéèêàùçœ\s]+(?:qui|dont|de)[a-zéèêàùçœ\s]+)`)
	matches := pattern1.FindAllStringSubmatch(text, -1)
	for _, match := range matches {
		if len(match) >= 3 {
			term := strings.TrimSpace(match[1])
			definition := strings.TrimSpace(match[2])
			// Limiter la définition à 100 caractères
			if len(definition) > 100 {
				definition = definition[:100]
			}
			kb.DefinitionFacts[term] = definition
			kb.TotalFactsExtracted++
		}
	}

	// Pattern: "cœur: organe qui pompe le sang"
	pattern2 := regexp.MustCompile(`([a-zéèêàùçœ]+):\s+([a-zéèêàùçœ\s]+)`)
	matches = pattern2.FindAllStringSubmatch(text, -1)
	for _, match := range matches {
		if len(match) >= 3 {
			term := strings.TrimSpace(match[1])
			definition := strings.TrimSpace(match[2])
			if len(definition) > 100 {
				definition = definition[:100]
			}
			kb.DefinitionFacts[term] = definition
			kb.TotalFactsExtracted++
		}
	}
}

// buildCoOccurrences construit les co-occurrences (mots à distance < 5)
func (kb *KnowledgeBase) buildCoOccurrences(text string) {
	// Tokeniser
	words := strings.Fields(strings.ToLower(text))

	// Filtrer mots courts et stopwords
	var filteredWords []string
	for _, word := range words {
		// Nettoyer ponctuation
		word = strings.Trim(word, ".,;:!?()[]{}\"'")
		if len(word) >= 4 && !isStopword(word) {
			filteredWords = append(filteredWords, word)
		}
	}

	// Fenêtre glissante de 5 mots
	windowSize := 5
	for i := 0; i < len(filteredWords); i++ {
		word1 := filteredWords[i]

		if kb.CoOccurrences[word1] == nil {
			kb.CoOccurrences[word1] = make(map[string]int)
		}

		// Mots dans la fenêtre
		for j := i + 1; j < i+windowSize && j < len(filteredWords); j++ {
			word2 := filteredWords[j]
			kb.CoOccurrences[word1][word2]++

			// Symétrique
			if kb.CoOccurrences[word2] == nil {
				kb.CoOccurrences[word2] = make(map[string]int)
			}
			kb.CoOccurrences[word2][word1]++
		}
	}
}

// trainAtomicNetwork entraîne le réseau atomique sur le texte
func (kb *KnowledgeBase) trainAtomicNetwork(text string) {
	// Extraire concepts clés (mots capitalisés + mots fréquents)
	words := strings.Fields(text)
	wordFreq := make(map[string]int)

	for _, word := range words {
		word = strings.Trim(word, ".,;:!?()[]{}\"'")
		if len(word) >= 4 {
			wordFreq[word]++
		}
	}

	// Sélectionner top concepts (au moins 3 occurrences)
	var concepts []string
	for word, freq := range wordFreq {
		if freq >= 3 {
			concepts = append(concepts, word)
		}
	}

	// Limiter à 100 concepts max par texte
	if len(concepts) > 100 {
		concepts = concepts[:100]
	}

	// Assigner chaque concept à un atome (ou réutiliser)
	for _, concept := range concepts {
		if _, exists := kb.ConceptMapping[concept]; !exists {
			// Trouver un atome libre (< 500)
			atomID := len(kb.ConceptMapping)
			if atomID < 500 {
				kb.ConceptMapping[concept] = atomID
			}
		}
	}

	// Activer les atomes correspondants
	for _, concept := range concepts {
		if atomID, exists := kb.ConceptMapping[concept]; exists {
			if atomID < len(kb.ConceptNetwork.Atoms) {
				kb.ConceptNetwork.Atoms[atomID].Perceptions[1] = 1.0 // Activation externe
			}
		}
	}

	// 50 itérations de convergence
	for iter := 0; iter < 50; iter++ {
		kb.ConceptNetwork.IterateNetwork()
	}
}

// GetRelatedConcepts retourne les concepts liés à un mot clé
func (kb *KnowledgeBase) GetRelatedConcepts(keyword string) []string {
	kb.mutex.RLock()
	defer kb.mutex.RUnlock()

	keyword = strings.ToLower(keyword)
	var related []string

	// Co-occurrences fortes (> 5 occurrences)
	if coOccurs, exists := kb.CoOccurrences[keyword]; exists {
		for word, count := range coOccurs {
			if count >= 5 {
				related = append(related, word)
			}
		}
	}

	return related
}

// findClosestMatch cherche le mot le plus proche dans la KB (avec gestion des synonymes)
func (kb *KnowledgeBase) findClosestMatch(term string) string {
	termLower := strings.ToLower(term)

	// D'abord: chercher correspondance exacte
	if _, exists := kb.DefinitionFacts[termLower]; exists {
		return termLower
	}

	// Deuxième: chercher dans les co-occurrences (mots souvent liés)
	if coOccur, exists := kb.CoOccurrences[termLower]; exists {
		// Retourner le mot le plus fréquemment associé
		var bestMatch string
		var bestCount int
		for word, count := range coOccur {
			if count > bestCount {
				bestCount = count
				bestMatch = word
			}
		}
		if bestMatch != "" {
			return strings.ToLower(bestMatch)
		}
	}

	// Troisième: chercher par similarité Levenshtein (distance d'édition)
	var closestMatch string
	minDistance := 5 // Maximum 5 différences de caractères

	for knownTerm := range kb.DefinitionFacts {
		dist := levenshteinDistance(termLower, knownTerm)
		if dist < minDistance {
			minDistance = dist
			closestMatch = knownTerm
		}
	}

	return closestMatch
}

// levenshteinDistance calcule la distance d'édition entre deux mots
func levenshteinDistance(a, b string) int {
	if len(a) == 0 {
		return len(b)
	}
	if len(b) == 0 {
		return len(a)
	}

	// Matrice de programmation dynamique
	d := make([][]int, len(a)+1)
	for i := range d {
		d[i] = make([]int, len(b)+1)
		d[i][0] = i
	}
	for j := range d[0] {
		d[0][j] = j
	}

	for i := 1; i <= len(a); i++ {
		for j := 1; j <= len(b); j++ {
			cost := 0
			if a[i-1] != b[j-1] {
				cost = 1
			}
			d[i][j] = min(d[i-1][j]+1, min(d[i][j-1]+1, d[i-1][j-1]+cost))
		}
	}
	return d[len(a)][len(b)]
}

// min retourne le minimum de deux entiers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// GetFactualInfo retourne les informations factuelles sur un terme
func (kb *KnowledgeBase) GetFactualInfo(term string) map[string]interface{} {
	kb.mutex.RLock()
	defer kb.mutex.RUnlock()

	info := make(map[string]interface{})

	// Utiliser le matching amélioré pour trouver le meilleur terme
	matchedTerm := kb.findClosestMatch(term)
	termLower := strings.ToLower(term)

	// Définition - utiliser le terme matché si trouvé
	var defTerm string
	if _, exists := kb.DefinitionFacts[termLower]; exists {
		defTerm = termLower
	} else if matchedTerm != "" {
		defTerm = matchedTerm
	}

	if defTerm != "" {
		if def, exists := kb.DefinitionFacts[defTerm]; exists {
			info["definition"] = def
		}
	}

	// Dates associées
	if dates, exists := kb.DateFacts[term]; exists {
		info["dates"] = dates
	}

	// Relations causales
	if causes, exists := kb.CausalFacts[term]; exists {
		info["causes"] = causes
	}

	// Lieux
	if locs, exists := kb.LocationFacts[term]; exists {
		info["locations"] = locs
	}

	// Concepts liés
	related := kb.GetRelatedConcepts(strings.ToLower(term))
	if len(related) > 0 {
		info["related"] = related
	}

	return info
}

// HasFactualKnowledge vérifie si la KB a des infos sur un terme
func (kb *KnowledgeBase) HasFactualKnowledge(term string) bool {
	info := kb.GetFactualInfo(term)
	return len(info) > 0
}

// GetConfidenceBoost retourne un boost de confiance basé sur les connaissances
func (kb *KnowledgeBase) GetConfidenceBoost(text string) float64 {
	kb.mutex.RLock()
	defer kb.mutex.RUnlock()

	// Tokeniser
	words := strings.Fields(strings.ToLower(text))

	var matchCount int
	var totalWords int

	for _, word := range words {
		word = strings.Trim(word, ".,;:!?()[]{}\"'")
		if len(word) >= 4 {
			totalWords++
			if kb.HasFactualKnowledge(word) {
				matchCount++
			}
		}
	}

	if totalWords == 0 {
		return 0
	}

	// Boost proportionnel au nombre de mots connus
	ratio := float64(matchCount) / float64(totalWords)
	return ratio * 0.3 // Max 30% de boost
}

// PrintStats affiche les statistiques d'apprentissage
func (kb *KnowledgeBase) PrintStats() {
	kb.mutex.RLock()
	defer kb.mutex.RUnlock()

	fmt.Println("\n============================================================")
	fmt.Println("  STATISTIQUES D'APPRENTISSAGE AUTOMATIQUE")
	fmt.Println("============================================================")
	fmt.Printf("Textes traités:     %d\n", kb.TotalTextsProcessed)
	fmt.Printf("Mots analysés:      %d\n", kb.TotalWordsProcessed)
	fmt.Printf("Faits extraits:     %d\n", kb.TotalFactsExtracted)
	fmt.Printf("\nBASE DE CONNAISSANCES:\n")
	fmt.Printf("  • Faits datés:     %d\n", len(kb.DateFacts))
	fmt.Printf("  • Relations causales: %d\n", len(kb.CausalFacts))
	fmt.Printf("  • Lieux:           %d\n", len(kb.LocationFacts))
	fmt.Printf("  • Définitions:     %d\n", len(kb.DefinitionFacts))
	fmt.Printf("  • Co-occurrences:  %d mots\n", len(kb.CoOccurrences))
	fmt.Printf("  • Concepts atomiques: %d\n", len(kb.ConceptMapping))
	fmt.Println("============================================================\n")
}

// SaveToFile sauvegarde la base de connaissances dans un fichier JSON
func (kb *KnowledgeBase) SaveToFile(filepath string) error {
	kb.mutex.RLock()
	defer kb.mutex.RUnlock()

	// Créer structure sérialisable (sans mutex et réseau atomique)
	type SerializableKB struct {
		DateFacts           map[string][]string
		CausalFacts         map[string][]string
		LocationFacts       map[string][]string
		DefinitionFacts     map[string]string
		CoOccurrences       map[string]map[string]int
		ConceptMapping      map[string]int
		TotalTextsProcessed int
		TotalWordsProcessed int
		TotalFactsExtracted int
	}

	skb := SerializableKB{
		DateFacts:           kb.DateFacts,
		CausalFacts:         kb.CausalFacts,
		LocationFacts:       kb.LocationFacts,
		DefinitionFacts:     kb.DefinitionFacts,
		CoOccurrences:       kb.CoOccurrences,
		ConceptMapping:      kb.ConceptMapping,
		TotalTextsProcessed: kb.TotalTextsProcessed,
		TotalWordsProcessed: kb.TotalWordsProcessed,
		TotalFactsExtracted: kb.TotalFactsExtracted,
	}

	data, err := json.MarshalIndent(skb, "", "  ")
	if err != nil {
		return fmt.Errorf("erreur sérialisation: %v", err)
	}

	err = os.WriteFile(filepath, data, 0644)
	if err != nil {
		return fmt.Errorf("erreur écriture fichier: %v", err)
	}

	return nil
}

// LoadFromFile charge la base de connaissances depuis un fichier JSON
func (kb *KnowledgeBase) LoadFromFile(filepath string) error {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return fmt.Errorf("erreur lecture fichier: %v", err)
	}

	type SerializableKB struct {
		DateFacts           map[string][]string
		CausalFacts         map[string][]string
		LocationFacts       map[string][]string
		DefinitionFacts     map[string]string
		CoOccurrences       map[string]map[string]int
		ConceptMapping      map[string]int
		TotalTextsProcessed int
		TotalWordsProcessed int
		TotalFactsExtracted int
	}

	var skb SerializableKB
	err = json.Unmarshal(data, &skb)
	if err != nil {
		return fmt.Errorf("erreur désérialisation: %v", err)
	}

	kb.mutex.Lock()
	defer kb.mutex.Unlock()

	kb.DateFacts = skb.DateFacts
	kb.CausalFacts = skb.CausalFacts
	kb.LocationFacts = skb.LocationFacts
	kb.DefinitionFacts = skb.DefinitionFacts
	kb.CoOccurrences = skb.CoOccurrences
	kb.ConceptMapping = skb.ConceptMapping
	kb.TotalTextsProcessed = skb.TotalTextsProcessed
	kb.TotalWordsProcessed = skb.TotalWordsProcessed
	kb.TotalFactsExtracted = skb.TotalFactsExtracted

	// Réinitialiser les maps nil pour éviter les panics
	if kb.CoOccurrences == nil {
		kb.CoOccurrences = make(map[string]map[string]int)
	}
	if kb.DateFacts == nil {
		kb.DateFacts = make(map[string][]string)
	}
	if kb.CausalFacts == nil {
		kb.CausalFacts = make(map[string][]string)
	}
	if kb.LocationFacts == nil {
		kb.LocationFacts = make(map[string][]string)
	}
	if kb.DefinitionFacts == nil {
		kb.DefinitionFacts = make(map[string]string)
	}
	if kb.ConceptMapping == nil {
		kb.ConceptMapping = make(map[string]int)
	}

	return nil
}

// isStopword vérifie si un mot est un stopword
func isStopword(word string) bool {
	stopwords := map[string]bool{
		"le": true, "la": true, "les": true, "un": true, "une": true, "des": true,
		"et": true, "ou": true, "mais": true, "donc": true, "car": true,
		"dans": true, "sur": true, "sous": true, "avec": true, "sans": true,
		"pour": true, "par": true, "vers": true, "chez": true,
		"être": true, "avoir": true, "faire": true, "dire": true, "aller": true,
		"ce": true, "cet": true, "cette": true, "ces": true,
		"qui": true, "que": true, "quoi": true, "dont": true, "où": true,
		"son": true, "sa": true, "ses": true, "leur": true, "leurs": true,
		"mon": true, "ma": true, "mes": true, "ton": true, "ta": true, "tes": true,
	}
	return stopwords[word]
}
