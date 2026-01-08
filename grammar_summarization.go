package main

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"IA-ATOMIQUE/database"
)

// ============================================================================
// PHASE 15: MODULE 5 - PIPELINE COMPLET GRAMMAR-AWARE SUMMARIZATION
// ============================================================================
// Intègre tous modules: Prétraitement → Résumé (Phase 13+++) → Optimisation
// Syntaxique → Enrichissement vocabulaire → Vérification cohérence

// TextType représente le type de texte détecté
type TextType int

const (
	ENCYCLOPEDIC TextType = iota
	NARRATIVE
	CONCEPTUAL
)

func (t TextType) String() string {
	switch t {
	case ENCYCLOPEDIC:
		return "Encyclopédique"
	case NARRATIVE:
		return "Narratif"
	case CONCEPTUAL:
		return "Conceptuel"
	default:
		return "Inconnu"
	}
}

// SystemStats contient les metrics système
type SystemStats struct {
	RAMUsedMB       uint64
	RAMTotalMB      uint64
	RAMPercent      float64
	GoroutinesCount int
	AllocationsMB   float64
	CPUUsagePercent float64
}

// GrammarAwareSummary contient le résumé optimisé
type GrammarAwareSummary struct {
	OriginalText          string
	PreprocessedText      string
	BaseSummary           string // Phase 13+++
	OptimizedSummary      string // Phase 15
	GrammarScore          float64
	StyleScore            float64
	CoherenceScore        float64
	LexicalRichness       float64
	ProcessingTime        int64 // en ms
	VariantsGenerated     int
	ImprovementPercentage float64
	SystemStats           SystemStats
	TextType              TextType // Type de texte détecté
	SkipAbstraction       bool     // Si true, saute Phase X+1
}

// GrammarSummarizer orchestre le pipeline Phase 15
type GrammarSummarizer struct {
	Enricher *database.VocabularyEnricher
}

// CaptureSystemStats capture les metrics système
func CaptureSystemStats() SystemStats {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	stats := SystemStats{
		RAMUsedMB:       m.Alloc / 1024 / 1024,
		RAMTotalMB:      m.TotalAlloc / 1024 / 1024,
		RAMPercent:      float64(m.Alloc) / float64(m.TotalAlloc),
		GoroutinesCount: runtime.NumGoroutine(),
		AllocationsMB:   float64(m.Alloc) / 1024 / 1024,
	}

	// Essayer de récupérer le CPU usage via /proc/self/stat (Linux)
	if data, err := os.ReadFile("/proc/self/stat"); err == nil {
		fields := strings.Fields(string(data))
		if len(fields) > 14 {
			if utime, err := strconv.ParseFloat(fields[13], 64); err == nil {
				if stime, err := strconv.ParseFloat(fields[14], 64); err == nil {
					stats.CPUUsagePercent = (utime + stime) / 100.0 // Approximation
				}
			}
		}
	}

	return stats
}

// ============================================================================
// DÉTECTION DE TYPE DE TEXTE
// ============================================================================

// DetectTextType détecte le type de texte (encyclopédique, narratif, conceptuel)
func DetectTextType(text string) TextType {
	lower := strings.ToLower(text)
	wordCount := len(strings.Fields(text))

	// Calculer scores pour chaque type
	var encyclopedicScore float64
	var narrativeScore float64
	var conceptualScore float64

	// === KEYWORDS ENCYCLOPÉDIQUES ===
	encyclopedicKeywords := []string{
		"définition", "est un", "est une", "aussi appelé", "également connu",
		"mesure", "pèse", "diamètre", "longueur", "hauteur", "largeur",
		"date", "siècle", "année", "mois", "jour", "né", "mort", "décédé",
		"localisation", "situé", "région", "pays", "continent", "ville",
		"habitat", "habitat naturel", "environnement", "écosystème",
		"espèce", "genre", "famille", "classification", "taxonomie",
		"propriété", "caractéristique", "trait", "feature",
		"composé", "élément", "substance", "matière",
		"fonction", "rôle", "utilité", "servir",
		"découverte", "inventé", "créé", "développé",
		"structure", "organisation", "constitution",
		"nombres", "pourcentage", "statistique", "données",
		"formule", "équation", "théorie", "loi scientifique",
	}

	// === KEYWORDS NARRATIFS ===
	narrativeKeywords := []string{
		"personnage", "héros", "protagoniste", "antagoniste",
		"raconte", "histoire", "conte", "légende", "saga", "épopée",
		"dialogue", "dit", "demanda", "répondit", "cria", "murmura",
		"personnage", "action", "action dramatique", "coup de théâtre",
		"scène", "acte", "tableau", "chapitre",
		"voyage", "aventure", "quête", "mission", "expédition",
		"rencontre", "amitié", "conflit", "confrontation",
		"émotion", "sentiment", "amour", "haine", "peur", "courage",
		"destin", "fatalité", "révélation", "secret",
		"description", "décor", "paysage", "atmosphère",
		"il y avait", "autrefois", "il était une fois", "jadis",
		"soudain", "brusquement", "tout à coup",
	}

	// === KEYWORDS CONCEPTUELS ===
	conceptualKeywords := []string{
		"concept", "idée", "principe", "théorie", "approche",
		"système", "structure", "logique", "mécanisme",
		"cause", "conséquence", "relation", "lien", "connexion",
		"analyse", "perspective", "vision", "interprétation",
		"abstraction", "généralisation", "universalité",
		"social", "politique", "économique", "culturel", "existentiel",
		"valeur", "sens", "signification", "implication",
		"justice", "injustice", "égalité", "inégalité",
		"liberté", "pouvoir", "autorité", "domination",
		"réflexion", "philosophie", "métaphysique",
		"paradoxe", "dialectique", "contradiction",
		"contexte", "cadre", "perspective", "angle",
		"implicite", "sous-entendu", "allusion",
	}

	// Compter occurrences de chaque type
	for _, keyword := range encyclopedicKeywords {
		encyclopedicScore += float64(strings.Count(lower, keyword))
	}
	for _, keyword := range narrativeKeywords {
		narrativeScore += float64(strings.Count(lower, keyword))
	}
	for _, keyword := range conceptualKeywords {
		conceptualScore += float64(strings.Count(lower, keyword))
	}

	// Analyser structure de phrases
	sentences := strings.Split(text, ".")
	avgSentenceLength := float64(len(text)) / float64(len(sentences))

	// Phrases longues (>60 chars) = narratif/conceptuel
	if avgSentenceLength > 60 {
		narrativeScore += 10
		conceptualScore += 10
	}

	// Phrases courtes (<40 chars) = encyclopédique
	if avgSentenceLength < 40 {
		encyclopedicScore += 10
	}

	// Compter les connecteurs logiques
	logicalConnectors := []string{"car", "donc", "ainsi", "par conséquent", "en résumé", "conclusion"}
	logicalCount := 0
	for _, connector := range logicalConnectors {
		logicalCount += strings.Count(lower, connector)
	}
	if logicalCount > 5 {
		conceptualScore += float64(logicalCount * 2)
	}

	// Normaliser par nombre de mots
	if wordCount > 0 {
		encyclopedicScore /= float64(wordCount) / 100.0
		narrativeScore /= float64(wordCount) / 100.0
		conceptualScore /= float64(wordCount) / 100.0
	}

	// Déterminer le type dominant
	if encyclopedicScore > narrativeScore && encyclopedicScore > conceptualScore {
		return ENCYCLOPEDIC
	}
	if narrativeScore > encyclopedicScore && narrativeScore > conceptualScore {
		return NARRATIVE
	}
	return CONCEPTUAL
}

// GetCompressionForType retourne la compression recommandée selon le type
func GetCompressionForType(textType TextType, userCompression float64) float64 {
	switch textType {
	case ENCYCLOPEDIC:
		// Pour encyclopédique, limiter la compression (perte d'infos trop élevée)
		if userCompression > 0.6 {
			return 0.6
		}
		return userCompression
	case NARRATIVE:
		// Pour narratif, compression modérée
		return userCompression
	case CONCEPTUAL:
		// Pour conceptuel, compression plus agressive possible
		if userCompression < 0.8 {
			return 0.8
		}
		return userCompression
	default:
		return userCompression
	}
}

// ShouldSkipAbstractionForType détermine si on saute Phase X+1 selon le type
func ShouldSkipAbstractionForType(textType TextType) bool {
	// Skip Phase X+1 pour textes encyclopédiques (besoin de faits concrets)
	return textType == ENCYCLOPEDIC
}

// NewGrammarSummarizer crée un nouveau pipeline
func NewGrammarSummarizer() *GrammarSummarizer {
	return &GrammarSummarizer{
		Enricher: database.NewVocabularyEnricher(),
	}
}

// ProcessWithPhase15 exécute le pipeline complet Phase 15
func (gs *GrammarSummarizer) ProcessWithPhase15(inputText string, threshold float64) *GrammarAwareSummary {
	startTime := time.Now()
	result := &GrammarAwareSummary{
		OriginalText: inputText,
	}

	// === ÉTAPE 0: Détection du type de texte ===
	fmt.Println("\n[PHASE 15] Étape 0: Détection du type de texte...")
	textType := DetectTextType(inputText)
	result.TextType = textType
	result.SkipAbstraction = ShouldSkipAbstractionForType(textType)

	// Adapter la compression selon le type
	adjustedThreshold := GetCompressionForType(textType, threshold)
	if adjustedThreshold != threshold {
		fmt.Printf("  ℹ️  Type: %s (compression ajustée %.2f → %.2f)\n", textType.String(), threshold, adjustedThreshold)
	} else {
		fmt.Printf("  ℹ️  Type: %s\n", textType.String())
	}
	if result.SkipAbstraction {
		fmt.Printf("  ⚠️  Phase X+1 désactivée pour ce type\n")
	}

	// === ÉTAPE 1: Prétraitement ===
	fmt.Println("\n[PHASE 15] Étape 1: Prétraitement & nettoyage...")
	preprocessResult := database.PreprocessText(inputText)
	result.PreprocessedText = preprocessResult.CleanedText
	fmt.Printf("  ✓ Supprimé %d lignes de bruit\n", len(preprocessResult.RemovedLines))
	fmt.Printf("  ✓ Fusionné %d fragments courts\n", preprocessResult.MergedSentences)
	fmt.Printf("  ✓ Normalisé ponctuation: %d opérations\n", preprocessResult.NormalizedCount)

	// === ÉTAPE 2: Résumé de base (Phase 13+++) ===
	fmt.Println("\n[PHASE 15] Étape 2: Résumé atomique (Phase 13+++)...")
	baseSummary := database.ResumerTexte(result.PreprocessedText, adjustedThreshold)
	result.BaseSummary = baseSummary
	fmt.Printf("  ✓ Résumé généré: %d caractères\n", len(baseSummary))

	// === ÉTAPE 3: Analyse syntaxique ===
	fmt.Println("\n[PHASE 15] Étape 3: Analyse syntaxique...")
	sentences := strings.Split(baseSummary, ".")
	var syntaxScores []float64

	for _, sent := range sentences {
		sent = strings.TrimSpace(sent)
		if sent == "" {
			continue
		}
		// Analyse syntaxique de chaque phrase - score basé sur richesse lexicale
		syntaxScores = append(syntaxScores, gs.Enricher.CalculateLexicalScore(sent))
	}

	avgGrammarScore := 0.7 // Score par défaut
	if len(syntaxScores) > 0 {
		for _, score := range syntaxScores {
			avgGrammarScore += score
		}
		avgGrammarScore /= float64(len(syntaxScores) + 1)
	}
	result.GrammarScore = avgGrammarScore
	fmt.Printf("  ✓ Score syntaxe moyen: %.2f%%\n", avgGrammarScore*100)

	// === ÉTAPE 4: Enrichissement vocabulaire ===
	fmt.Println("\n[PHASE 15] Étape 4: Enrichissement vocabulaire...")
	enrichedSummary := gs.EnrichSummary(baseSummary)
	fmt.Printf("  ✓ Vocabulaire enrichi (style naturel)\n")

	// === ÉTAPE 5: Génération de variantes optimisées ===
	fmt.Println("\n[PHASE 15] Étape 5: Génération de variantes...")
	variants := gs.GenerateOptimizedVariants(enrichedSummary, 5)
	result.VariantsGenerated = len(variants)
	fmt.Printf("  ✓ Généré %d variantes\n", len(variants))

	// === ÉTAPE 6: Sélection de la meilleure ===
	fmt.Println("\n[PHASE 15] Étape 6: Sélection optimale...")
	bestVariant, bestScore := gs.SelectBestVariant(variants)
	result.OptimizedSummary = bestVariant
	result.StyleScore = bestScore.StyleScore
	result.LexicalRichness = bestScore.LexicalScore
	fmt.Printf("  ✓ Meilleure variante: style %.2f%%, richesse %.2f%%\n",
		bestScore.StyleScore*100, bestScore.LexicalScore*100)

	// === ÉTAPE 7: Vérification cohérence ===
	fmt.Println("\n[PHASE 15] Étape 7: Vérification cohérence...")
	coherenceScore := gs.AnalyzeFinalCoherence(result.OptimizedSummary)
	result.CoherenceScore = coherenceScore
	fmt.Printf("  ✓ Score cohérence: %.2f%%\n", coherenceScore*100)

	// === ÉTAPE 8: PHASE X+1 - SEMANTIC ABSTRACTION LAYER ===
	if result.SkipAbstraction {
		fmt.Println("\n[PHASE X+1] Étape 8: Abstraction sémantique (SKIPPÉE pour type encyclopédique)...")
		fmt.Printf("  ℹ️  Texte encyclopédique: conservation des faits concrets\n")
	} else {
		fmt.Println("\n[PHASE X+1] Étape 8: Abstraction sémantique forcée...")

		// Analyser le texte original pour extraire les concepts
		phrasesOriginales := strings.Split(inputText, ".")
		var phrasesListe []string
		for _, p := range phrasesOriginales {
			p = strings.TrimSpace(p)
			if len(p) > 10 {
				phrasesListe = append(phrasesListe, p)
			}
		}

		analyseSemantique := database.AnalyserSemantiquement(inputText, phrasesListe)
		scoreAbstraction := database.EvaluerAbstraction(result.OptimizedSummary, analyseSemantique)

		fmt.Printf("  → Score d'abstraction: %.1f%%\n", scoreAbstraction.ScoreGlobal)

		// Appliquer abstraction forcée si score < 60%
		if scoreAbstraction.ScoreGlobal < 60.0 {
			fmt.Printf("  ⚠️  Score < 60%% → Réécrit en phrases conceptuelles\n")
			phrasesAbstraites := database.GenererPhrasesConceptuelles(analyseSemantique)

			// === ÉTAPE 9: PHASE X+3 - NATURAL SYNTAX LAYER ===
			fmt.Println("\n[PHASE X+3] Étape 9: Humanisation syntaxique (pas de connecteurs explicites)...")
			resumeHumain := database.HumanizeStructure(phrasesAbstraites)
			result.OptimizedSummary = resumeHumain
			fmt.Printf("  ✓ Syntaxe naturelle appliquée (subordination, ponctuation, rythme)\n")
		} else {
			fmt.Printf("  ✓ Score acceptable: pas de réécriture\n")
		}
	}

	// === Résultats finaux ===
	result.ProcessingTime = time.Since(startTime).Milliseconds()
	result.ImprovementPercentage = ((result.StyleScore + result.CoherenceScore + result.LexicalRichness) / 3.0) - (result.GrammarScore / 3.0)

	// Capturer les stats système finales
	result.SystemStats = CaptureSystemStats()

	return result
}

// EnrichSummary enrichit un résumé avec vocabulaire varié
func (gs *GrammarSummarizer) EnrichSummary(summary string) string {
	sentences := strings.Split(summary, ".")
	var enriched []string

	for i, sent := range sentences {
		sent = strings.TrimSpace(sent)
		if sent == "" {
			continue
		}

		// Enrichir avec vocabulaire naturel
		enrichedSent := gs.Enricher.EnrichSentence(sent, "naturel")

		// Ajouter connecteur transitoire si ce n'est pas la première phrase
		if i > 0 && i%2 == 0 {
			enrichedSent = gs.Enricher.AddTransitionConnector(enrichedSent)
		}

		enriched = append(enriched, enrichedSent)
	}

	return strings.Join(enriched, ". ") + "."
}

// VariantScore contient le score d'une variante
type VariantScore struct {
	Variant        string
	GrammarScore   float64
	StyleScore     float64
	LexicalScore   float64
	CoherenceScore float64
	FinalScore     float64
}

// GenerateOptimizedVariants génère N variantes optimisées
func (gs *GrammarSummarizer) GenerateOptimizedVariants(summary string, n int) []VariantScore {
	styles := []string{"naturel", "formel", "littéraire"}
	var variants []VariantScore

	// Variante 1: Originale
	original := VariantScore{
		Variant: summary,
	}
	original.CalculateScores(gs.Enricher)
	variants = append(variants, original)

	// Générer N-1 variantes
	for i := 1; i < n; i++ {
		style := styles[i%len(styles)]
		variant := gs.Enricher.EnrichParagraph(summary, style)

		vs := VariantScore{
			Variant: variant,
		}
		vs.CalculateScores(gs.Enricher)
		variants = append(variants, vs)
	}

	return variants
}

// CalculateScores calcule tous les scores pour une variante
func (vs *VariantScore) CalculateScores(enricher *database.VocabularyEnricher) {
	// Score lexical
	vs.LexicalScore = enricher.CalculateLexicalScore(vs.Variant)

	// Score style (richesse + fluidité)
	vs.StyleScore = vs.LexicalScore

	// Score grammaire (basé sur structure)
	vs.GrammarScore = 0.8 // Score par défaut pour structure correcte

	// Score cohérence
	vs.CoherenceScore = 0.7 // Placeholder - intégrer analyse cohérence réelle

	// Score final (pondéré)
	vs.FinalScore = (vs.GrammarScore * 0.25) + (vs.StyleScore * 0.35) +
		(vs.LexicalScore * 0.25) + (vs.CoherenceScore * 0.15)
}

// SelectBestVariant sélectionne la meilleure variante
func (gs *GrammarSummarizer) SelectBestVariant(variants []VariantScore) (string, VariantScore) {
	if len(variants) == 0 {
		return "", VariantScore{}
	}

	best := variants[0]
	for _, v := range variants[1:] {
		if v.FinalScore > best.FinalScore {
			best = v
		}
	}

	return best.Variant, best
}

// AnalyzeFinalCoherence analyse la cohérence du résumé final
func (gs *GrammarSummarizer) AnalyzeFinalCoherence(summary string) float64 {
	sentences := strings.Split(summary, ".")
	if len(sentences) < 2 {
		return 1.0
	}

	// Score cohérence simple: connecteurs + transitions
	connectorCount := 0
	totalSentences := 0

	connectors := []string{
		"de plus", "en outre", "cependant", "néanmoins", "toutefois",
		"par conséquent", "dès lors", "ainsi", "donc", "puis",
		"ensuite", "alors", "finalement", "en conclusion",
	}

	for _, sent := range sentences {
		sent = strings.ToLower(sent)
		if sent == "" {
			continue
		}
		totalSentences++

		for _, conn := range connectors {
			if strings.Contains(sent, conn) {
				connectorCount++
				break
			}
		}
	}

	if totalSentences == 0 {
		return 0.5
	}

	return float64(connectorCount) / float64(totalSentences)
}

// GetSummaryReport génère un rapport complet
func (gas *GrammarAwareSummary) GetSummaryReport() string {
	var report strings.Builder

	report.WriteString("\n╔════════════════════════════════════════════════════════════╗\n")
	report.WriteString("║         PHASE 15: GRAMMAR-AWARE SUMMARIZATION REPORT        ║\n")
	report.WriteString("╚════════════════════════════════════════════════════════════╝\n\n")

	report.WriteString("\n📊 METRICS:\n")
	report.WriteString("────────────────────────────────────────────────────────────\n")
	report.WriteString(fmt.Sprintf("Grammar Score:     %.1f%%\n", gas.GrammarScore*100))
	report.WriteString(fmt.Sprintf("Style Score:       %.1f%%\n", gas.StyleScore*100))
	report.WriteString(fmt.Sprintf("Coherence Score:   %.1f%%\n", gas.CoherenceScore*100))
	report.WriteString(fmt.Sprintf("Lexical Richness:  %.1f%%\n", gas.LexicalRichness*100))
	report.WriteString(fmt.Sprintf("Improvement:       +%.1f%%\n", gas.ImprovementPercentage*100))

	report.WriteString("\n📝 TEXT ANALYSIS:\n")
	report.WriteString("────────────────────────────────────────────────────────────\n")
	report.WriteString(fmt.Sprintf("Text Type:         %s\n", gas.TextType.String()))
	if gas.SkipAbstraction {
		report.WriteString(fmt.Sprintf("Abstraction:       SKIPPED (facts preservation)\n"))
	} else {
		report.WriteString(fmt.Sprintf("Abstraction:       APPLIED\n"))
	}

	report.WriteString("\n📈 PROCESSING:\n")
	report.WriteString("────────────────────────────────────────────────────────────\n")
	report.WriteString(fmt.Sprintf("Original Length:   %d chars\n", len(gas.OriginalText)))
	report.WriteString(fmt.Sprintf("Summary Length:    %d chars\n", len(gas.OptimizedSummary)))
	report.WriteString(fmt.Sprintf("Compression:       %.1f%%\n",
		100.0*(1.0-float64(len(gas.OptimizedSummary))/float64(len(gas.OriginalText)))))
	report.WriteString(fmt.Sprintf("Variants Created:  %d\n", gas.VariantsGenerated))
	report.WriteString(fmt.Sprintf("Processing Time:   %d ms\n", gas.ProcessingTime))

	report.WriteString("\n💾 SYSTEM RESOURCES:\n")
	report.WriteString("────────────────────────────────────────────────────────────\n")
	report.WriteString(fmt.Sprintf("RAM Used:          %.1f MB / %.1f MB (%.1f%%)\n",
		float64(gas.SystemStats.RAMUsedMB), float64(gas.SystemStats.RAMTotalMB), gas.SystemStats.RAMPercent*100))
	report.WriteString(fmt.Sprintf("Go Routines:       %d\n", gas.SystemStats.GoroutinesCount))
	report.WriteString(fmt.Sprintf("Memory Alloc:      %.1f MB\n", gas.SystemStats.AllocationsMB))
	report.WriteString(fmt.Sprintf("CPU Usage:         %.2f%%\n", gas.SystemStats.CPUUsagePercent))

	report.WriteString("\n💬 SUMMARY:\n")
	report.WriteString("────────────────────────────────────────────────────────────\n")
	report.WriteString(gas.OptimizedSummary)
	report.WriteString("\n\n")

	return report.String()
}

// SaveOptimizedSummary sauvegarde le résumé optimisé
func (gas *GrammarAwareSummary) SaveOptimizedSummary(baseFilename string) error {
	optimizedFile := strings.TrimSuffix(baseFilename, ".txt") + "_optimized_phase15.txt"
	return os.WriteFile(optimizedFile, []byte(gas.GetSummaryReport()), 0644)
}
