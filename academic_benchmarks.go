package main

import (
	"IA-ATOMIQUE/database"
	"fmt"
	"strings"
	"time"
)

// RunMMLUBenchmark exécute le benchmark MMLU
func RunMMLUBenchmark() {
	sep := strings.Repeat("=", 70)
	fmt.Println("\n" + sep)
	fmt.Println("  BENCHMARK MMLU - Massive Multitask Language Understanding")
	fmt.Println(sep)

	fmt.Println("\n[INITIALISATION]")
	fmt.Println("  • Chargement des questions test...")

	// Créer les questions de test
	questions := database.CreateSampleMMLU()
	fmt.Printf("  • Questions chargées: %d\n", len(questions))
	fmt.Printf("  • Sujets: Histoire, Médecine, Droit, Mathématiques, Sciences, Littérature\n")

	// Créer le moteur
	engine := database.NewMMLUEngine()

	// Lancer le benchmark
	fmt.Println("\n[EXÉCUTION]")
	startTime := time.Now()

	benchmark := engine.RunBenchmark(questions)

	elapsed := time.Since(startTime)

	// Afficher les résultats
	benchmark.PrintResults()

	// Performances
	fmt.Printf("[PERFORMANCES]\n")
	fmt.Printf("  • Temps total: %v\n", elapsed)
	fmt.Printf("  • Temps par question: %v\n", elapsed/time.Duration(len(questions)))

	if len(questions) > 0 {
		questionsPerSec := float64(len(questions)) / elapsed.Seconds()
		fmt.Printf("  • Questions/seconde: %.0f\n", questionsPerSec)

		// Extrapolation pour 16 000 questions
		totalTime := time.Duration(float64(elapsed) * (16000.0 / float64(len(questions))))
		fmt.Printf("\n  • Extrapolation pour 16 000 questions: %v\n", totalTime)
	}

	// Comparaison
	fmt.Printf("\n[COMPARAISON STANDARDS]\n")
	fmt.Println("  Scores de référence:")
	fmt.Println("    • Humain expert:     ~90%")
	fmt.Println("    • GPT-4:             ~86%")
	fmt.Println("    • GPT-3.5:           ~70%")
	fmt.Printf("    • IA-ATOMIQUE:       %.2f%% ", benchmark.Score)

	if benchmark.Score >= 85 {
		fmt.Println("✓ (Niveau GPT-4+)")
	} else if benchmark.Score >= 70 {
		fmt.Println("✓ (Niveau GPT-3.5+)")
	} else if benchmark.Score >= 60 {
		fmt.Println("(Niveau correct)")
	} else {
		fmt.Println("(À améliorer)")
	}

	fmt.Println("\n" + sep + "\n")
}

// RunHellaswagBenchmark exécute le benchmark Hellaswag
func RunHellaswagBenchmark() {
	sep := strings.Repeat("=", 70)
	fmt.Println("\n" + sep)
	fmt.Println("  BENCHMARK HELLASWAG - Raisonnement de Bon Sens")
	fmt.Println(sep)

	fmt.Println("\n[INITIALISATION]")
	fmt.Println("  • Chargement des questions test...")

	// Créer les questions
	questions := database.CreateSampleHellaswag()
	fmt.Printf("  • Questions chargées: %d\n", len(questions))
	fmt.Printf("  • Activités: Cuisine, Sport, Travail, Médecine, etc.\n")

	// Créer le moteur
	engine := database.NewHellaswagEngine()

	// Lancer le benchmark
	fmt.Println("\n[EXÉCUTION]")
	startTime := time.Now()

	benchmark := engine.RunBenchmark(questions)

	elapsed := time.Since(startTime)

	// Afficher les résultats
	benchmark.PrintResults()

	// Performances
	fmt.Printf("[PERFORMANCES]\n")
	fmt.Printf("  • Temps total: %v\n", elapsed)
	fmt.Printf("  • Temps par question: %v\n", elapsed/time.Duration(len(questions)))

	if len(questions) > 0 {
		questionsPerSec := float64(len(questions)) / elapsed.Seconds()
		fmt.Printf("  • Questions/seconde: %.0f\n", questionsPerSec)
	}

	// Comparaison
	fmt.Printf("\n[COMPARAISON STANDARDS]\n")
	fmt.Println("  Scores de référence:")
	fmt.Println("    • Humain:            ~95%")
	fmt.Println("    • GPT-4:             ~95%")
	fmt.Println("    • GPT-3:             ~85%")
	fmt.Println("    • BERT:              ~75%")
	fmt.Printf("    • IA-ATOMIQUE:       %.2f%% ", benchmark.Score)

	if benchmark.Score >= 90 {
		fmt.Println("✓ (Niveau Humain/GPT-4)")
	} else if benchmark.Score >= 80 {
		fmt.Println("✓ (Niveau GPT-3+)")
	} else if benchmark.Score >= 70 {
		fmt.Println("✓ (Niveau BERT+)")
	} else {
		fmt.Println("(À améliorer)")
	}

	fmt.Println("\n" + sep + "\n")
}

// RunAllAcademicBenchmarks exécute tous les benchmarks académiques
func RunAllAcademicBenchmarks() {
	sep := strings.Repeat("=", 70)
	fmt.Println("\n" + sep)
	fmt.Println("  SUITE COMPLÈTE DE BENCHMARKS ACADÉMIQUES")
	fmt.Println(sep)

	startTotal := time.Now()

	// MMLU
	fmt.Println("\n[1/2] MMLU - Culture Générale Multi-Domaines")
	RunMMLUBenchmark()

	// Hellaswag
	fmt.Println("\n[2/2] Hellaswag - Raisonnement de Bon Sens")
	RunHellaswagBenchmark()

	totalTime := time.Since(startTotal)

	// Résumé final
	fmt.Println("\n" + sep)
	fmt.Println("  RÉSUMÉ FINAL")
	fmt.Println(sep)
	fmt.Printf("\n  • Temps total: %v\n", totalTime)
	fmt.Println("\n  ✓ Tous les benchmarks académiques terminés!")
	fmt.Println("\n  Voir BENCHMARK_RESULTS.md pour le tableau comparatif complet.")
	fmt.Println("\n" + sep + "\n")
}

// HandleAcademicBenchmarks gère les commandes de benchmarks académiques
func HandleAcademicBenchmarks(args []string) {
	if len(args) < 1 {
		fmt.Println("Benchmarks académiques disponibles:")
		fmt.Println("  ./programme academic mmlu      - Test MMLU (culture générale)")
		fmt.Println("  ./programme academic hellaswag - Test Hellaswag (bon sens)")
		fmt.Println("  ./programme academic all       - Tous les benchmarks")
		fmt.Println("\nExemples:")
		fmt.Println("  ./programme academic mmlu")
		fmt.Println("  ./programme academic all")
		return
	}

	benchType := args[0]

	switch benchType {
	case "mmlu":
		RunMMLUBenchmark()

	case "hellaswag", "hella":
		RunHellaswagBenchmark()

	case "all", "complete", "full":
		RunAllAcademicBenchmarks()

	default:
		fmt.Printf("Benchmark inconnu: %s\n", benchType)
		fmt.Println("Benchmarks disponibles: mmlu, hellaswag, all")
	}
}
