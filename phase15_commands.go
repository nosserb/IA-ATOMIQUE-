package main

import (
	"IA-ATOMIQUE/database"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// ============================================================================
// PHASE 15: COMMANDES CLI - GRAMMAR-AWARE SUMMARIZATION
// ============================================================================

func resumeOptimizedCommand(args []string) {
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

func compareSummariesCommand(args []string) {
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

func analyzePreprocessingCommand(args []string) {
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

func analyzeVocabularyCommand(args []string) {
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
	resumeOptimizedCommand(args)
}
