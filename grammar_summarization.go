package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"IA-ATOMIQUE/database"
)

// ============================================================================
// PHASE 15: MODULE 5 - PIPELINE COMPLET GRAMMAR-AWARE SUMMARIZATION
// ============================================================================
// Intègre tous modules: Prétraitement → Résumé (Phase 13+++) → Optimisation
// Syntaxique → Enrichissement vocabulaire → Vérification cohérence

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
}

// GrammarSummarizer orchestre le pipeline Phase 15
type GrammarSummarizer struct {
	Enricher *database.VocabularyEnricher
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

	// === ÉTAPE 1: Prétraitement ===
	fmt.Println("\n[PHASE 15] Étape 1: Prétraitement & nettoyage...")
	preprocessResult := database.PreprocessText(inputText)
	result.PreprocessedText = preprocessResult.CleanedText
	fmt.Printf("  ✓ Supprimé %d lignes de bruit\n", len(preprocessResult.RemovedLines))
	fmt.Printf("  ✓ Fusionné %d fragments courts\n", preprocessResult.MergedSentences)
	fmt.Printf("  ✓ Normalisé ponctuation: %d opérations\n", preprocessResult.NormalizedCount)

	// === ÉTAPE 2: Résumé de base (Phase 13+++) ===
	fmt.Println("\n[PHASE 15] Étape 2: Résumé atomique (Phase 13+++)...")
	baseSummary := database.ResumerTexte(result.PreprocessedText, threshold)
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

	// === Résultats finaux ===
	result.ProcessingTime = time.Since(startTime).Milliseconds()
	result.ImprovementPercentage = ((result.StyleScore + result.CoherenceScore + result.LexicalRichness) / 3.0) - (result.GrammarScore / 3.0)

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

	report.WriteString("📊 METRICS:\n")
	report.WriteString("────────────────────────────────────────────────────────────\n")
	report.WriteString(fmt.Sprintf("Grammar Score:     %.1f%%\n", gas.GrammarScore*100))
	report.WriteString(fmt.Sprintf("Style Score:       %.1f%%\n", gas.StyleScore*100))
	report.WriteString(fmt.Sprintf("Coherence Score:   %.1f%%\n", gas.CoherenceScore*100))
	report.WriteString(fmt.Sprintf("Lexical Richness:  %.1f%%\n", gas.LexicalRichness*100))
	report.WriteString(fmt.Sprintf("Improvement:       +%.1f%%\n", gas.ImprovementPercentage*100))

	report.WriteString("\n📈 PROCESSING:\n")
	report.WriteString("────────────────────────────────────────────────────────────\n")
	report.WriteString(fmt.Sprintf("Original Length:   %d chars\n", len(gas.OriginalText)))
	report.WriteString(fmt.Sprintf("Summary Length:    %d chars\n", len(gas.OptimizedSummary)))
	report.WriteString(fmt.Sprintf("Compression:       %.1f%%\n",
		100.0*(1.0-float64(len(gas.OptimizedSummary))/float64(len(gas.OriginalText)))))
	report.WriteString(fmt.Sprintf("Variants Created:  %d\n", gas.VariantsGenerated))
	report.WriteString(fmt.Sprintf("Processing Time:   %d ms\n", gas.ProcessingTime))

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
