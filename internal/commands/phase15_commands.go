package commands

import (
	"IA-ATOMIQUE/database"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// ============================================================================
// PHASE 15: COMMANDES CLI - GRAMMAR-AWARE SUMMARIZATION
// ============================================================================

func ResumeOptimizedCommand(args []string) {
	if len(args) == 0 {
		fmt.Println("[INFO] Utilisation: ./programme resume-optimized <file.txt> <threshold> [output_file]")
		fmt.Println("       Paramètres:")
		fmt.Println("       - file.txt: Fichier à résumer")
		fmt.Println("       - threshold: Seuil de résumé (0.05-0.20)")
		fmt.Println("       - output_file: (optionnel) Fichier de sortie")
		return
	}

	filename := args[0]
	threshold := 0.10 // Défaut

	if len(args) > 1 {
		if t, err := strconv.ParseFloat(args[1], 64); err == nil {
			threshold = t
		}
	}

	// Lire le fichier
	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("[ERREUR] Impossible de lire %s: %v\n", filename, err)
		return
	}

	text := string(content)

	// === Créer et exécuter le pipeline ===
	summarizer := NewGrammarSummarizer()
	fmt.Println("╔════════════════════════════════════════════════════════════╗")
	fmt.Println("║     PHASE 15: GRAMMAR-AWARE SUMMARIZATION PIPELINE         ║")
	fmt.Println("╚════════════════════════════════════════════════════════════╝")

	result := summarizer.ProcessWithPhase15(text, threshold)

	// Afficher le rapport
	fmt.Print(result.GetSummaryReport())

	// Sauvegarder le résumé
	outputFile := filename
	if len(args) > 2 {
		outputFile = args[2]
	}

	if err := result.SaveOptimizedSummary(outputFile); err != nil {
		fmt.Printf("[ERREUR] Impossible de sauvegarder: %v\n", err)
	} else {
		savedFile := strings.TrimSuffix(outputFile, ".txt") + "_optimized_phase15.txt"
		fmt.Printf("\n✓ Résumé optimisé sauvegardé: %s\n", savedFile)
	}
}

// ResumeHybridCommand produit un résumé long, fluide et hybride (atomique + proba)
func ResumeHybridCommand(args []string) {
	// Timing et mémoire
	startTime := time.Now()
	var startMem runtime.MemStats
	runtime.ReadMemStats(&startMem)

	if len(args) == 0 {
		fmt.Println("[INFO] Utilisation: ./programme resume <file.txt> [ratio=0.55]")
		fmt.Println("       Mode hybride (couverture élevée) atomique + pondération proba")
		return
	}

	filename := args[0]
	targetRatio := 0.55

	if len(args) > 1 {
		if t, err := strconv.ParseFloat(args[1], 64); err == nil {
			targetRatio = t
		}
	}

	targetRatio = clampFloat(targetRatio, 0.35, 0.90)

	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("[ERREUR] Impossible de lire %s: %v\n", filename, err)
		return
	}

	text := string(content)

	fmt.Println("╔════════════════════════════════════════════════════════════╗")
	fmt.Println("║     HYBRID FULL SUMMARY - ATOMIQUE + PROBABILISTE         ║")
	fmt.Println("╚════════════════════════════════════════════════════════════╝")
	fmt.Printf("🎯 Mode: Résumé 10 lignes max\n")
	fmt.Println("⚙️  Pipeline: extraction atomique + scoring proba + reformulation")

	llmRatio := clampFloat(targetRatio*0.65, 0.30, 0.60)

	summarizer := NewGrammarSummarizer()
	fmt.Println("\n[MODULE PROBA] Phase 15 pour la réécriture fluide...")
	phaseResult := summarizer.ProcessWithPhase15(text, llmRatio)

	// Limiter à 10 lignes maximum
	summary := limitTo10Lines(phaseResult.OptimizedSummary)

	outputFile := strings.TrimSuffix(filename, ".txt") + "_hybrid_full.txt"
	report := buildHybridReport(summary, targetRatio, targetRatio, 0.85, "hybrid", phaseResult)

	fmt.Println("\n📝 Résumé (4 paragraphes de 10 lignes max):")
	fmt.Println("────────────────────────────────────────────────────────────")
	fmt.Println(summary)

	if err := os.WriteFile(outputFile, []byte(report), 0644); err != nil {
		fmt.Printf("[ERREUR] Impossible de sauvegarder: %v\n", err)
		return
	}

	// Calculer les stats finales
	endTime := time.Now()
	var endMem runtime.MemStats
	runtime.ReadMemStats(&endMem)

	duration := endTime.Sub(startTime)
	ramUsed := endMem.Alloc - startMem.Alloc
	inputSize := float64(len(text)) / 1024.0     // KB
	outputSize := float64(len(summary)) / 1024.0 // KB
	compressionRatio := (float64(len(text)) / float64(len(summary)))

	// Afficher les logs détaillés
	fmt.Println("\n╔════════════════════════════════════════════════════════════╗")
	fmt.Println("║              PERFORMANCE & RESOURCE LOGS                  ║")
	fmt.Println("╚════════════════════════════════════════════════════════════╝")
	fmt.Printf("\n⏱️  TIMING:\n")
	fmt.Printf("   Durée totale:    %v\n", duration)
	fmt.Printf("   Durée Phase 15:  %d ms\n", phaseResult.ProcessingTime)

	fmt.Printf("\n💾 MÉMOIRE:\n")
	fmt.Printf("   RAM utilisée:    %.2f MB\n", float64(ramUsed)/(1024*1024))
	fmt.Printf("   RAM totale:      %.2f MB\n", float64(endMem.Alloc)/(1024*1024))
	fmt.Printf("   Goroutines:      %d\n", runtime.NumGoroutine())

	fmt.Printf("\n📊 TAILLE:\n")
	fmt.Printf("   Input:           %.2f KB\n", inputSize)
	fmt.Printf("   Output:          %.2f KB\n", outputSize)
	fmt.Printf("   Ratio:           %.2f:1\n", compressionRatio)

	fmt.Printf("\n✅ Résumé sauvegardé: %s\n", outputFile)
}

func CompareSummariesCommand(args []string) {
	if len(args) == 0 {
		fmt.Println("[INFO] Utilisation: ./programme compare-summaries <file.txt> <threshold>")
		fmt.Println("       Compare Phase 13+++ vs Phase 15")
		return
	}

	filename := args[0]
	threshold := 0.10

	if len(args) > 1 {
		if t, err := strconv.ParseFloat(args[1], 64); err == nil {
			threshold = t
		}
	}

	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("[ERREUR] Impossible de lire %s: %v\n", filename, err)
		return
	}

	text := string(content)

	// === Phase 13++ (basique) ===
	fmt.Println("\n╔════════════════════════════════════════════════════════════╗")
	fmt.Println("║           COMPARISON: Phase 13+++ vs Phase 15              ║")
	fmt.Println("╚════════════════════════════════════════════════════════════╝\n")

	fmt.Println("📊 Phase 13+++ (Baseline):")
	fmt.Println("────────────────────────────────────────────────────────────")
	basicSummary := database.ResumerTexte(text, threshold)
	fmt.Printf("Length: %d chars\n", len(basicSummary))
	fmt.Printf("Summary: %s...\n\n", basicSummary[:minInt(100, len(basicSummary))])

	// === Phase 15 (optimisé) ===
	fmt.Println("✨ Phase 15 (Grammar-Aware):")
	fmt.Println("────────────────────────────────────────────────────────────")
	summarizer := NewGrammarSummarizer()
	result := summarizer.ProcessWithPhase15(text, threshold)

	fmt.Printf("Grammar Score:   %.1f%%\n", result.GrammarScore*100)
	fmt.Printf("Style Score:     %.1f%%\n", result.StyleScore*100)
	fmt.Printf("Coherence:       %.1f%%\n", result.CoherenceScore*100)
	fmt.Printf("Improvement:     +%.1f%%\n", result.ImprovementPercentage*100)
	fmt.Printf("Variants Tested: %d\n", result.VariantsGenerated)
	fmt.Printf("Processing:      %d ms\n\n", result.ProcessingTime)

	fmt.Println("📝 Optimized Summary:")
	fmt.Println("────────────────────────────────────────────────────────────")
	fmt.Println(result.OptimizedSummary)
}

func AnalyzePreprocessingCommand(args []string) {
	if len(args) == 0 {
		fmt.Println("[INFO] Utilisation: ./programme analyze-preprocessing <file.txt>")
		fmt.Println("       Analyse l'impact du prétraitement")
		return
	}

	filename := args[0]
	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("[ERREUR] Impossible de lire %s: %v\n", filename, err)
		return
	}

	text := string(content)

	fmt.Println("\n╔════════════════════════════════════════════════════════════╗")
	fmt.Println("║        MODULE 1: PREPROCESSING ANALYSIS                   ║")
	fmt.Println("╚════════════════════════════════════════════════════════════╝\n")

	// Appliquer le prétraitement
	result := database.PreprocessText(text)

	fmt.Println("📋 ORIGINAL TEXT:")
	fmt.Printf("Length: %d chars\n", len(text))
	fmt.Printf("Sentences: %d\n", result.OriginalLength)
	fmt.Printf("Lines: %d\n\n", len(strings.Split(text, "\n")))

	fmt.Println("🧹 AFTER CLEANING:")
	fmt.Printf("Length: %d chars\n", len(result.CleanedText))
	fmt.Printf("Sentences: %d\n", result.CleanedLength)
	fmt.Printf("Reduction: %.1f%%\n\n", 100.0*(1.0-float64(len(result.CleanedText))/float64(len(text))))

	fmt.Println("⚙️ OPERATIONS:")
	fmt.Printf("Noise lines removed: %d\n", len(result.RemovedLines))
	fmt.Printf("Sentences merged: %d\n", result.MergedSentences)
	fmt.Printf("Normalizations: %d\n\n", result.NormalizedCount)

	if len(result.RemovedLines) > 0 {
		fmt.Println("🗑️ NOISE EXAMPLES (first 3):")
		for i, line := range result.RemovedLines {
			if i >= 3 {
				break
			}
			if len(line) > 80 {
				fmt.Printf("  %s...\n", line[:80])
			} else {
				fmt.Printf("  %s\n", line)
			}
		}
	}
}

func AnalyzeVocabularyCommand(args []string) {
	if len(args) == 0 {
		fmt.Println("[INFO] Utilisation: ./programme analyze-vocabulary <sentence>")
		fmt.Println("       Analyse l'enrichissement lexical")
		return
	}

	sentence := strings.Join(args, " ")

	fmt.Println("\n╔════════════════════════════════════════════════════════════╗")
	fmt.Println("║      MODULE 3: VOCABULARY ENRICHMENT ANALYSIS             ║")
	fmt.Println("╚════════════════════════════════════════════════════════════╝\n")

	enricher := database.NewVocabularyEnricher()

	fmt.Println("📝 Original Sentence:")
	fmt.Printf("  %s\n\n", sentence)

	// Générer variantes
	variants := enricher.GenerateVariants(sentence, 4, []string{"naturel", "formel", "littéraire"})

	fmt.Println("✨ Generated Variants:")
	fmt.Println("────────────────────────────────────────────────────────────")
	for i, variant := range variants {
		lexScore := enricher.CalculateLexicalScore(variant)
		richness := enricher.RichnessBonus

		if i == 0 {
			fmt.Printf("%d. [ORIGINAL] %s\n", i+1, variant)
			fmt.Printf("   Score: %.1f%% | Richness: %.1f%%\n\n", lexScore*100, richness*100)
		} else {
			fmt.Printf("%d. [VARIANT %d] %s\n", i+1, i, variant)
			fmt.Printf("   Score: %.1f%% | Richness: %.1f%%\n\n", lexScore*100, richness*100)
		}
	}

	// Montrer synonymes disponibles
	fmt.Println("📚 Available Synonyms:")
	fmt.Println("────────────────────────────────────────────────────────────")
	words := strings.Fields(sentence)
	foundSynonyms := false

	for _, word := range words {
		synonyms := enricher.GetSynonymsForWord(word)
		if len(synonyms) > 0 {
			foundSynonyms = true
			fmt.Printf("  %s: %s\n", word, strings.Join(synonyms[:minInt(3, len(synonyms))], ", "))
		}
	}

	if !foundSynonyms {
		fmt.Println("  (No enrichment dictionary entries for this sentence)")
	}
}

// Helper function
func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// clampFloat bounds a value between min and max
func clampFloat(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// buildHybridReport generates the final hybrid summary report
func buildHybridReport(finalSummary string, targetRatio, actualRatio, fidelityScore float64, fidelityMode string, phaseResult *GrammarAwareSummary) string {
	var b strings.Builder

	b.WriteString("╔════════════════════════════════════════════════════════════╗\n")
	b.WriteString("║        HYBRID FULL SUMMARY - ATOMIQUE + PROBA REPORT       ║\n")
	b.WriteString("╚════════════════════════════════════════════════════════════╝\n\n")

	b.WriteString(fmt.Sprintf("Cible: %.1f%% | Réel: %.1f%%\n", targetRatio*100, actualRatio*100))
	b.WriteString(fmt.Sprintf("Fidélité hybride: %.1f%% (%s)\n", fidelityScore*100, fidelityMode))
	b.WriteString(fmt.Sprintf("Phase 15 - Style: %.1f%% | Cohérence: %.1f%% | Grammaire: %.1f%%\n",
		phaseResult.StyleScore*100,
		phaseResult.CoherenceScore*100,
		phaseResult.GrammarScore*100,
	))

	b.WriteString("\n[Résumé complet]\n")
	b.WriteString("────────────────────────────────────────────────────────────\n")
	b.WriteString(finalSummary)
	b.WriteString("\n")

	return b.String()
}

// cleanNoiseFromText supprime le bruit et les caractères parasites du texte
func cleanNoiseFromText(text string) string {
	// Remplacer les caractères spéciaux bizarres
	replacements := map[string]string{
		"º": "",
		"ñ": "",
		"ü": "",
		"ö": "",
		"ø": "",
		"§": "",
		"¶": "",
		"†": "",
		"‡": "",
		"|": " ",
		"}": " ",
		"{": " ",
		"[": " ",
		"]": " ",
		"№": "",
		"™": "",
		"®": "",
		"©": "",
	}

	result := text
	for old, new := range replacements {
		result = strings.ReplaceAll(result, old, new)
	}

	// Supprimer les lignes qui ne contiennent que des caractères spéciaux ou des numéros
	lines := strings.Split(result, "\n")
	var cleanLines []string

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Ignorer les lignes vides
		if line == "" {
			continue
		}

		// Ignorer les lignes qui sont principalement des numéros/caractères spéciaux
		alphaCount := 0
		for _, r := range line {
			if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
				alphaCount++
			}
		}

		// Garder seulement si au moins 30% de caractères alphabétiques
		if len(line) > 0 && float64(alphaCount)/float64(len(line)) > 0.3 {
			cleanLines = append(cleanLines, line)
		}
	}

	// Rejoindre et nettoyer les espacements multiples
	result = strings.Join(cleanLines, " ")

	// Supprimer les espacements multiples
	for strings.Contains(result, "  ") {
		result = strings.ReplaceAll(result, "  ", " ")
	}

	return strings.TrimSpace(result)
}

// limitTo10Lines limite un texte à 4 paragraphes de 10 lignes max (40 lignes total)
func limitTo10Lines(text string) string {
	// Nettoyer le bruit d'abord
	text = cleanNoiseFromText(text)

	// Séparer par phrases (à la fin d'une phrase)
	sentences := strings.FieldsFunc(text, func(r rune) bool {
		return r == '.' || r == '!' || r == '?'
	})

	// 4 paragraphes x 10 phrases = 40 phrases max
	maxParagraphes := 4
	phrasesParParagraphe := 10
	maxSentences := maxParagraphes * phrasesParParagraphe

	// Limiter aux N premières phrases
	if len(sentences) > maxSentences {
		sentences = sentences[:maxSentences]
	}

	// Organiser en paragraphes
	var paragraphes []string
	for i := 0; i < len(sentences); i += phrasesParParagraphe {
		end := i + phrasesParParagraphe
		if end > len(sentences) {
			end = len(sentences)
		}

		para := strings.Join(sentences[i:end], ". ")
		if para != "" {
			para = para + "."
		}
		if para != "" {
			paragraphes = append(paragraphes, para)
		}
	}

	// Joindre les paragraphes avec double saut de ligne
	result := strings.Join(paragraphes, "\n\n")
	return result
}

// phaseCommand dispatcher
func phaseCommand(args []string) {
	if len(args) == 0 {
		fmt.Println("\n[PHASE 15] Grammar-Aware Summarization Commands:")
		fmt.Println("─────────────────────────────────────────────────")
		fmt.Println("  ./programme resume-optimized <file> <threshold>")
		fmt.Println("    → Full pipeline: preprocess → summarize → optimize")
		fmt.Println("")
		fmt.Println("  ./programme compare-summaries <file> <threshold>")
		fmt.Println("    → Compare Phase 13+++ vs Phase 15")
		fmt.Println("")
		fmt.Println("  ./programme analyze-preprocessing <file>")
		fmt.Println("    → Show preprocessing impact (Module 1)")
		fmt.Println("")
		fmt.Println("  ./programme analyze-vocabulary <sentence>")
		fmt.Println("    → Show vocabulary enrichment (Module 3)")
		return
	}

	// main.go appelle directement les commandes spécifiques
	ResumeOptimizedCommand(args)
}
