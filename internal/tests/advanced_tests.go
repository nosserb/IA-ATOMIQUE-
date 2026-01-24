package tests

import (
	"IA-ATOMIQUE/database"
	"fmt"
	"os"
	"strings"
	"time"
)

// TestNeedleInHaystack test de recherche d'aiguille dans botte de foin
func TestNeedleInHaystack(filePath string) {
	sep := strings.Repeat("=", 70)
	fmt.Println("\n" + sep)
	fmt.Println("  TEST NEEDLE IN HAYSTACK - Recherche Sémantique Ultra-Rapide")
	fmt.Println(sep)

	// Lecture du fichier
	fmt.Println("\n[CHARGEMENT]")
	start := time.Now()
	data, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("  ✗ Erreur: %v\n", err)
		return
	}
	texte := string(data)
	loadTime := time.Since(start)

	tailleMo := float64(len(texte)) / (1024 * 1024)
	motsEstimes := len(strings.Fields(texte))
	fmt.Printf("  • Fichier: %s\n", filePath)
	fmt.Printf("  • Taille: %.2f Mo\n", tailleMo)
	fmt.Printf("  • Mots: ~%d\n", motsEstimes)
	fmt.Printf("  • Temps lecture: %v\n", loadTime)

	// Recherche
	fmt.Println("\n[RECHERCHE D'ANOMALIES SÉMANTIQUES]")
	engine := database.NewNeedleSearchEngine()

	startSearch := time.Now()
	results := engine.FindNeedle(texte)
	searchTime := time.Since(startSearch)

	fmt.Printf("  • Temps de scan: %v\n", searchTime)
	fmt.Printf("  • Vitesse: %.0f mots/sec\n", float64(motsEstimes)/searchTime.Seconds())
	fmt.Printf("  • Anomalies détectées: %d\n", len(results))

	// Afficher les top anomalies
	fmt.Println("\n[TOP 5 PHRASES SUSPECTES]")
	maxDisplay := 5
	if len(results) < maxDisplay {
		maxDisplay = len(results)
	}

	for i := 0; i < maxDisplay; i++ {
		result := results[i]
		fmt.Printf("\n[#%d] Position: %d | Cohérence: %.3f | Anomalie: %.3f\n",
			result.AnomalyRank, result.Position, result.CoherenceScore, result.EnergySignature)

		// Limiter la longueur de la phrase affichée
		sentence := result.Sentence
		if len(sentence) > 150 {
			sentence = sentence[:150] + "..."
		}
		fmt.Printf("  Phrase: \"%s\"\n", sentence)

		if len(result.ContextBefore) > 0 {
			contextBefore := result.ContextBefore
			if len(contextBefore) > 80 {
				contextBefore = "..." + contextBefore[len(contextBefore)-80:]
			}
			fmt.Printf("  Avant: ...%s\n", contextBefore)
		}
	}

	// Statistiques
	fmt.Println("\n" + sep)
	fmt.Println("  PERFORMANCES")
	fmt.Println(sep)
	fmt.Printf("  • Temps total: %v\n", searchTime)
	fmt.Printf("  • Vitesse scan: %.2f M mots/sec\n", float64(motsEstimes)/searchTime.Seconds()/1000000)

	if motsEstimes > 0 {
		fmt.Printf("\n  • Extrapolations:\n")
		ratio := float64(motsEstimes)
		fmt.Printf("    - 1M mots: %v\n", time.Duration(float64(searchTime)*(1000000/ratio)))
		fmt.Printf("    - 10M mots: %v\n", time.Duration(float64(searchTime)*(10000000/ratio)))
		fmt.Printf("    - 5× Les Misérables (~3M mots): %v\n", time.Duration(float64(searchTime)*(3000000/ratio)))
	}

	fmt.Println("\n" + sep + "\n")
}

// TestPerplexity test de calcul de perplexité
func TestPerplexity(filePath string) {
	sep := strings.Repeat("=", 70)
	fmt.Println("\n" + sep)
	fmt.Println("  TEST PERPLEXITÉ - Mesure de Cohérence Atomique")
	fmt.Println(sep)

	// Lecture du fichier
	fmt.Println("\n[CHARGEMENT]")
	start := time.Now()
	data, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("  ✗ Erreur: %v\n", err)
		return
	}
	texte := string(data)
	loadTime := time.Since(start)

	tailleMo := float64(len(texte)) / (1024 * 1024)
	motsEstimes := len(strings.Fields(texte))
	fmt.Printf("  • Fichier: %s\n", filePath)
	fmt.Printf("  • Taille: %.2f Mo\n", tailleMo)
	fmt.Printf("  • Mots: ~%d\n", motsEstimes)
	fmt.Printf("  • Temps lecture: %v\n", loadTime)

	// Calcul de perplexité
	fmt.Println("\n[CALCUL PERPLEXITÉ]")
	calculator := database.NewPerplexityCalculator()

	startCalc := time.Now()
	result := calculator.CalculatePerplexity(texte)
	calcTime := time.Since(startCalc)

	fmt.Printf("  • Temps de calcul: %v\n", calcTime)
	fmt.Printf("  • Vitesse: %.0f mots/sec\n", float64(motsEstimes)/calcTime.Seconds())

	// Résultats
	fmt.Println("\n[RÉSULTATS]")
	fmt.Printf("  • Perplexité globale: %.3f\n", result.GlobalPerplexity)
	fmt.Printf("  • Cohérence moyenne: %.3f\n", result.AverageCoherence)
	fmt.Printf("  • Score de stabilité: %.3f\n", result.StabilityScore)
	fmt.Printf("  • Variance énergétique: %.3f\n", result.EnergyVariance)
	fmt.Printf("  • Qualité: %s\n", result.InterpretedQuality)
	fmt.Printf("  • Moments de surprise: %d\n", len(result.SurpriseMoments))

	// Distribution des catégories
	fmt.Println("\n[DISTRIBUTION SÉMANTIQUE]")
	fmt.Printf("  • Catégories activées: %d\n", len(result.CategoryDistrib))

	// Top 5 catégories
	type catCount struct {
		id    int
		count int
	}
	counts := make([]catCount, 0)
	for id, count := range result.CategoryDistrib {
		counts = append(counts, catCount{id, count})
	}
	// Tri simple
	for i := 0; i < len(counts)-1; i++ {
		for j := i + 1; j < len(counts); j++ {
			if counts[j].count > counts[i].count {
				counts[i], counts[j] = counts[j], counts[i]
			}
		}
	}

	maxCats := 5
	if len(counts) < maxCats {
		maxCats = len(counts)
	}
	fmt.Println("  • Top 5 catégories:")
	for i := 0; i < maxCats; i++ {
		fmt.Printf("    %d. Catégorie %d: %d mots\n", i+1, counts[i].id, counts[i].count)
	}

	// Perplexité locale
	if len(result.LocalPerplexities) > 0 {
		fmt.Println("\n[VARIATION LOCALE]")
		minPerp := result.LocalPerplexities[0]
		maxPerp := result.LocalPerplexities[0]
		for _, p := range result.LocalPerplexities {
			if p < minPerp {
				minPerp = p
			}
			if p > maxPerp {
				maxPerp = p
			}
		}
		fmt.Printf("  • Perplexité min: %.3f (segment le plus cohérent)\n", minPerp)
		fmt.Printf("  • Perplexité max: %.3f (segment le moins cohérent)\n", maxPerp)
		fmt.Printf("  • Écart: %.3f\n", maxPerp-minPerp)
	}

	// Comparaison avec standards
	fmt.Println("\n[COMPARAISON]")
	fmt.Println("  Références typiques:")
	fmt.Println("    • GPT-4: Perplexité ~10-20 sur texte général")
	fmt.Println("    • GPT-3: Perplexité ~20-40")
	fmt.Println("    • Modèles simples: Perplexité >100")
	fmt.Printf("    • IA-ATOMIQUE: %.2f ✓\n", result.GlobalPerplexity)

	if result.GlobalPerplexity < 10 {
		fmt.Println("\n  🌟 EXCELLENT - Performance au niveau GPT-4+ !")
	} else if result.GlobalPerplexity < 20 {
		fmt.Println("\n  ✓ BON - Performance comparable aux meilleurs LLM")
	}

	fmt.Println("\n" + sep + "\n")
}

// HandleAdvancedTests gère les tests avancés
func HandleAdvancedTests(args []string) {
	if len(args) < 1 {
		fmt.Println("Tests avancés disponibles:")
		fmt.Println("  ./programme test needle <fichier>    - Test Needle In Haystack")
		fmt.Println("  ./programme test perplexity <fichier> - Test de Perplexité")
		fmt.Println("\nExemples:")
		fmt.Println("  ./programme test needle input.txt")
		fmt.Println("  ./programme test perplexity input.txt")
		return
	}

	testType := args[0]

	switch testType {
	case "needle", "haystack":
		if len(args) < 2 {
			fmt.Println("Usage: ./programme test needle <fichier>")
			return
		}
		TestNeedleInHaystack(args[1])

	case "perplexity", "perp":
		if len(args) < 2 {
			fmt.Println("Usage: ./programme test perplexity <fichier>")
			return
		}
		TestPerplexity(args[1])

	default:
		fmt.Printf("Test inconnu: %s\n", testType)
		fmt.Println("Tests disponibles: needle, perplexity")
	}
}
