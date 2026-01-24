package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"IA-ATOMIQUE/database"
)

// ============================================================================
// ANTI-HALLUCINATION COMMAND MODULE
// Commandes CLI pour tester et valider les stratégies anti-hallucination
// ============================================================================

// ProcessAntiHallucination orchestre le test complet des stratégies
func ProcessAntiHallucination(args []string) {
	if len(args) < 2 {
		fmt.Println("Usage: programme fidelity <file|test> [compression_ratio]")
		fmt.Println("  fidelity file <input.txt>     - Analyser un fichier avec tous les critères")
		fmt.Println("  fidelity test                 - Lancer suite de tests complets")
		fmt.Println("  fidelity compare <file>       - Comparer stratégies A-E sur un fichier")
		return
	}

	strategy := args[1]

	switch strategy {
	case "file":
		if len(args) < 3 {
			fmt.Println("Usage: programme fidelity file <input.txt>")
			return
		}
		ProcessFidelityAnalysis(args[2])

	case "test":
		RunCompleteFidelityTests()

	case "compare":
		if len(args) < 3 {
			fmt.Println("Usage: programme fidelity compare <input.txt>")
			return
		}
		CompareAllStrategies(args[2])

	case "hybrid":
		if len(args) < 3 {
			fmt.Println("Usage: programme fidelity hybrid <input.txt> [tau=0.8]")
			return
		}
		tau := 0.8
		if len(args) > 3 {
			fmt.Sscanf(args[3], "%f", &tau)
		}
		TestHybridApproach(args[2], tau)

	default:
		fmt.Printf("Stratégie inconnue: %s\n", strategy)
	}
}

// ProcessFidelityAnalysis analyse un texte avec tous les critères de fidélité
func ProcessFidelityAnalysis(filename string) {
	// Lire le texte
	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("❌ Erreur lecture: %v\n", err)
		return
	}

	sourceText := string(content)

	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("PHASE 15 + FIDELITY CHECK")
	fmt.Println(strings.Repeat("=", 80))

	// Étape 1: Extraire vocabulaire source
	fmt.Println("\n[1] EXTRACTION VOCABULAIRE SOURCE")
	fmt.Println(strings.Repeat("-", 80))

	keyTerms := database.ExtractKeyTerms(sourceText)
	fmt.Printf("  ✓ Termes clés identifiés: %d\n", len(keyTerms))

	// Étape 2: Générer résumé Phase 13+++
	fmt.Println("\n[2] RÉSUMÉ PHASE 13+++")
	fmt.Println(strings.Repeat("-", 80))

	compressionRatio := 0.3 // Compression de 30%
	baseSummary := database.ResumerTexte(sourceText, compressionRatio)
	fmt.Printf("  Texte original: %d mots\n", len(strings.Fields(sourceText)))
	fmt.Printf("  Résumé généré:  %d mots\n", len(strings.Fields(baseSummary)))
	fmt.Printf("  Compression:    %.1f%%\n\n", compressionRatio*100)
	if len(baseSummary) > 200 {
		fmt.Printf("  %s...\n\n", baseSummary[:200])
	} else {
		fmt.Printf("  %s\n\n", baseSummary)
	}

	// Étape 3: Calculer fidelité
	fmt.Println("[3] ANALYSE FIDÉLITÉ")
	fmt.Println(strings.Repeat("-", 80))

	fidelity := database.CalculateFidelity(baseSummary, sourceText)
	fmt.Printf("  Coverage (Ff): %.2f%%\n", fidelity*100)
	fmt.Printf("  Threshold acceptable: 0.80 (80%%)\n")

	// Étape 4: Décision hybride
	fmt.Println("\n[4] DÉCISION HYBRIDE")
	fmt.Println(strings.Repeat("-", 80))

	tau := 0.80
	finalSummary, finalFidelity, mode := database.HybridResume(baseSummary, sourceText, tau)

	fmt.Printf("  Mode sélectionné: %s\n", mode)
	fmt.Printf("  Fidelité: %.2f%%\n", finalFidelity*100)

	fmt.Println("\n[5] RÉSUMÉ FINAL")
	fmt.Println(strings.Repeat("-", 80))
	if len(finalSummary) > 200 {
		fmt.Printf("%s...\n\n", finalSummary[:200])
	} else {
		fmt.Printf("%s\n\n", finalSummary)
	}

	// Sauvegarde
	timestamp := "20260108"
	baseName := strings.TrimSuffix(filepath.Base(filename), filepath.Ext(filename))

	outputFile := fmt.Sprintf("%s_fidelity_report_%s.txt", baseName, timestamp)
	reportContent := fmt.Sprintf("PHASE 15 FIDELITY ANALYSIS REPORT\n"+
		"================================\n\n"+
		"Original text length: %d words\n"+
		"Generated summary: %d words\n"+
		"Compression ratio: %.1f%%\n\n"+
		"FIDELITY SCORE: %.2f%%\n"+
		"Mode: %s\n\n"+
		"FINAL SUMMARY:\n%s\n",
		len(strings.Fields(sourceText)),
		len(strings.Fields(baseSummary)),
		compressionRatio*100,
		finalFidelity*100,
		mode,
		finalSummary)

	err = os.WriteFile(outputFile, []byte(reportContent), 0644)
	if err != nil {
		fmt.Printf("❌ Erreur sauvegarde: %v\n", err)
	} else {
		fmt.Printf("\n✓ Rapport sauvegardé: %s\n", outputFile)
	}
}

// CompareAllStrategies compare les stratégies sur un même texte
func CompareAllStrategies(filename string) {
	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("❌ Erreur lecture: %v\n", err)
		return
	}

	sourceText := string(content)

	fmt.Println("\n" + strings.Repeat("=", 100))
	fmt.Println("COMPARAISON DES STRATÉGIES ANTI-HALLUCINATION")
	fmt.Println(strings.Repeat("=", 100))

	compressionRatio := 0.3
	baseSummary := database.ResumerTexte(sourceText, compressionRatio)

	// STRATÉGIE A: Extraction par phrases clés
	fmt.Println("\n[STRATÉGIE A] Extraction par phrases clés")
	fmt.Println(strings.Repeat("-", 100))

	strategyAPhrases := database.ExtractiveResume(sourceText, compressionRatio)
	strategyAFidelity := database.CalculateFidelity(strategyAPhrases, sourceText)

	fmt.Printf("  Fidelity: %.2f%%\n", strategyAFidelity*100)

	// STRATÉGIE B: Filtrage post-génération
	fmt.Println("\n[STRATÉGIE B] Filtrage post-génération")
	fmt.Println(strings.Repeat("-", 100))

	// Simuler un résumé avec du bruit
	noisySummary := baseSummary + " avec concepts inventés"
	strategyBFidelity := database.CalculateFidelity(noisySummary, sourceText)

	fmt.Printf("  Fidelity avant filtrage: %.2f%%\n", strategyBFidelity*100)

	// STRATÉGIE C: Hybridation
	fmt.Println("\n[STRATÉGIE C] Hybridation extractif + génératif")
	fmt.Println(strings.Repeat("-", 100))

	finalSummary, strategyCFidelity, mode := database.HybridResume(baseSummary, sourceText, 0.80)

	fmt.Printf("  Mode: %s\n", mode)
	fmt.Printf("  Fidelity: %.2f%%\n", strategyCFidelity*100)
	fmt.Printf("  Final summary: %s\n", finalSummary[:min(len(finalSummary), 100)])

	// Tableau comparatif
	fmt.Println("\n" + strings.Repeat("=", 100))
	fmt.Println("TABLEAU COMPARATIF")
	fmt.Println(strings.Repeat("=", 100))
	fmt.Printf("%-20s | %-15s | %-40s\n", "STRATÉGIE", "FIDÉLITÉ", "MODE")
	fmt.Println(strings.Repeat("-", 100))
	fmt.Printf("%-20s | %14.2f%% | %s\n", "A (Extractive)", strategyAFidelity*100, "Extraction TF-IDF")
	fmt.Printf("%-20s | %14.2f%% | %s\n", "Base (Phase 13+++)", database.CalculateFidelity(baseSummary, sourceText)*100, "Génératif")
	fmt.Printf("%-20s | %14.2f%% | %s\n", "C (Hybrid)", strategyCFidelity*100, mode)
}

// TestHybridApproach teste l'approche hybridation
func TestHybridApproach(filename string, tau float64) {
	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("❌ Erreur: %v\n", err)
		return
	}

	sourceText := string(content)

	fmt.Printf("\n[HYBRID APPROACH TEST] (tau=%.2f)\n", tau)
	fmt.Println(strings.Repeat("=", 80))

	generatedSummary := database.ResumerTexte(sourceText, 0.3)

	fmt.Println("Résumé généré (Phase 13+++):")
	fmt.Printf("Fidelity: %.2f%%\n\n", database.CalculateFidelity(generatedSummary, sourceText)*100)

	finalSummary, finalFidelity, mode := database.HybridResume(generatedSummary, sourceText, tau)

	fmt.Println("\nDécision hybride:")
	fmt.Printf("Mode: %s\n", mode)
	fmt.Printf("Final fidelity: %.2f%%\n", finalFidelity*100)
	fmt.Printf("Final summary preview: %s\n", finalSummary[:min(len(finalSummary), 100)])
}

// RunCompleteFidelityTests exécute une suite complète de tests
func RunCompleteFidelityTests() {
	fmt.Println("\n" + strings.Repeat("=", 100))
	fmt.Println("SUITE COMPLÈTE: TESTS ANTI-HALLUCINATION")
	fmt.Println(strings.Repeat("=", 100))

	// Créer des textes de test simples
	testCases := []struct {
		name   string
		source string
	}{
		{
			"Technique Simple",
			"Un atome computationnel est une unité autonome. Il perçoit ses voisins. Il possède un état interne. Tous les atomes forment un réseau asynchrone. Le réseau converge via résonance.",
		},
		{
			"Texte Encyclopédique",
			"L'intelligence artificielle atomique utilise des structures décentralisées. Chaque atome calcule localement. La communication est asynchrone. Le réseau émerge sans coordinateur central. Les applications incluent la NLP et l'analyse sémantique.",
		},
	}

	for _, tc := range testCases {
		fmt.Printf("\n[TEST: %s]\n", tc.name)
		fmt.Println(strings.Repeat("-", 100))

		generated := database.ResumerTexte(tc.source, 0.3)
		fidelity := database.CalculateFidelity(generated, tc.source)

		fmt.Printf("  Source length: %d words\n", len(strings.Fields(tc.source)))
		fmt.Printf("  Generated length: %d words\n", len(strings.Fields(generated)))
		fmt.Printf("  Fidelity: %.2f%%\n", fidelity*100)

		if fidelity < 0.80 {
			fmt.Printf("  ⚠️  Coverage < 80%% → Would use extractive\n")
		}
	}
}
