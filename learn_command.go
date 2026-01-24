package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ============================================================================
// COMMANDES D'APPRENTISSAGE AUTOMATIQUE
// ============================================================================

// GlobalKnowledgeBase - Instance globale de la base de connaissances
var GlobalKnowledgeBase *KnowledgeBase

// initKnowledgeBase initialise la base de connaissances globale
func initKnowledgeBase() {
	if GlobalKnowledgeBase == nil {
		GlobalKnowledgeBase = NewKnowledgeBase()

		// Charger la base existante si elle existe
		if _, err := os.Stat("knowledge_base.json"); err == nil {
			err := GlobalKnowledgeBase.LoadFromFile("knowledge_base.json")
			if err != nil {
				fmt.Printf("⚠️ Impossible de charger knowledge_base.json: %v\n", err)
			} else {
				fmt.Println("📚 Base de connaissances chargée (knowledge_base.json)")
			}
		}
	}
}

// LearnCommand - Apprend à partir d'un fichier ou dossier
func LearnCommand(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: ./programme learn <fichier.txt|dossier>")
		fmt.Println("\nExemples:")
		fmt.Println("  ./programme learn histoire.txt")
		fmt.Println("  ./programme learn corpus/histoire/")
		fmt.Println("  ./programme learn medecine.txt")
		return
	}

	initKnowledgeBase()

	path := args[0]
	info, err := os.Stat(path)
	if err != nil {
		fmt.Printf("❌ Erreur: %v\n", err)
		return
	}

	if info.IsDir() {
		// Traiter tous les fichiers .txt du dossier
		learnFromDirectory(path)
	} else {
		// Traiter un seul fichier
		learnFromFile(path)
	}

	// Afficher statistiques
	GlobalKnowledgeBase.PrintStats()

	// Sauvegarder automatiquement
	err = GlobalKnowledgeBase.SaveToFile("knowledge_base.json")
	if err != nil {
		fmt.Printf("⚠️ Erreur sauvegarde: %v\n", err)
	} else {
		fmt.Println("💾 Base de connaissances sauvegardée dans knowledge_base.json")
	}
}

// learnFromFile apprend à partir d'un seul fichier
func learnFromFile(filepath string) {
	fmt.Printf("\n📖 Lecture de %s...\n", filepath)

	err := GlobalKnowledgeBase.LearnFromTextFile(filepath)
	if err != nil {
		fmt.Printf("❌ Erreur: %v\n", err)
		return
	}

	fmt.Println("✅ Apprentissage terminé!")
}

// learnFromDirectory apprend à partir de tous les fichiers .txt d'un dossier
func learnFromDirectory(dirpath string) {
	fmt.Printf("\n📚 Lecture du dossier %s...\n", dirpath)

	var count int
	err := filepath.Walk(dirpath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Ne traiter que les fichiers .txt
		if !info.IsDir() && strings.HasSuffix(strings.ToLower(path), ".txt") {
			fmt.Printf("  • Traitement: %s\n", filepath.Base(path))
			err := GlobalKnowledgeBase.LearnFromTextFile(path)
			if err != nil {
				fmt.Printf("    ⚠️ Erreur: %v\n", err)
			} else {
				count++
			}
		}

		return nil
	})

	if err != nil {
		fmt.Printf("❌ Erreur parcours dossier: %v\n", err)
		return
	}

	fmt.Printf("\n✅ %d fichiers traités!\n", count)
}

// KnowledgeCommand - Affiche les connaissances sur un terme
func KnowledgeCommand(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: ./programme knowledge <terme>")
		fmt.Println("\nExemple:")
		fmt.Println("  ./programme knowledge Napoléon")
		fmt.Println("  ./programme knowledge 1815")
		fmt.Println("  ./programme knowledge coeur")
		return
	}

	initKnowledgeBase()

	term := args[0]
	fmt.Println("============================================================")

	info := GlobalKnowledgeBase.GetFactualInfo(term)

	if len(info) == 0 {
		fmt.Println("❌ Aucune connaissance trouvée.")
		fmt.Println("\n💡 Utilisez 'learn' pour apprendre à partir de textes:")
		fmt.Println("   ./programme learn histoire.txt")
		return
	}

	// Définition
	if def, exists := info["definition"]; exists {
		fmt.Printf("\n📖 DÉFINITION:\n   %s\n", def)
	}

	// Dates
	if dates, exists := info["dates"]; exists {
		fmt.Printf("\n📅 ÉVÉNEMENTS:\n")
		for _, event := range dates.([]string) {
			fmt.Printf("   • %s\n", event)
		}
	}

	// Causes/Effets
	if causes, exists := info["causes"]; exists {
		fmt.Printf("\n🔗 RELATIONS CAUSALES:\n")
		for _, effect := range causes.([]string) {
			fmt.Printf("   → %s\n", effect)
		}
	}

	// Lieux
	if locs, exists := info["locations"]; exists {
		fmt.Printf("\n🌍 LOCALISATION:\n")
		for _, loc := range locs.([]string) {
			fmt.Printf("   • %s\n", loc)
		}
	}

	// Concepts liés
	if related, exists := info["related"]; exists {
		fmt.Printf("\n🔗 CONCEPTS LIÉS:\n")
		relatedList := related.([]string)
		// Limiter à 10
		if len(relatedList) > 10 {
			relatedList = relatedList[:10]
		}
		for _, concept := range relatedList {
			fmt.Printf("   • %s\n", concept)
		}
	}

	fmt.Println("============================================================\n")
}

// StatsCommand - Affiche les statistiques de la base de connaissances
func StatsCommand(args []string) {
	initKnowledgeBase()

	if GlobalKnowledgeBase.TotalTextsProcessed == 0 {
		fmt.Println("\n❌ Base de connaissances vide.")
		fmt.Println("\n💡 Utilisez 'learn' pour apprendre à partir de textes:")
		fmt.Println("   ./programme learn histoire.txt")
		return
	}

	GlobalKnowledgeBase.PrintStats()

	// Exemples de connaissances
	fmt.Println("EXEMPLES DE CONNAISSANCES ACQUISES:")
	fmt.Println("------------------------------------------------------------")

	// 3 définitions
	count := 0
	for term, def := range GlobalKnowledgeBase.DefinitionFacts {
		if count >= 3 {
			break
		}
		fmt.Printf("• %s: %s\n", term, def)
		count++
	}

	// 3 faits datés
	count = 0
	for year, events := range GlobalKnowledgeBase.DateFacts {
		if count >= 3 {
			break
		}
		if len(events) > 0 {
			fmt.Printf("• %s: %s\n", year, events[0])
		}
		count++
	}

	fmt.Println("============================================================\n")
}

// TestKnowledgeCommand - Teste les connaissances avec une question
func TestKnowledgeCommand(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: ./programme test-knowledge <question>")
		fmt.Println("\nExemple:")
		fmt.Println("  ./programme test-knowledge \"Quand a eu lieu la bataille de Waterloo?\"")
		return
	}

	initKnowledgeBase()

	question := strings.Join(args[0:], " ")
	fmt.Println("============================================================")

	// Extraire mots-clés de la question
	words := strings.Fields(strings.ToLower(question))

	var foundKnowledge bool
	for _, word := range words {
		word = strings.Trim(word, ".,;:!?()[]{}\"'")
		if len(word) >= 4 {
			info := GlobalKnowledgeBase.GetFactualInfo(word)
			if len(info) > 0 {
				foundKnowledge = true
				fmt.Printf("\n🔍 Connaissances sur '%s':\n", word)

				// Afficher les infos pertinentes
				if def, exists := info["definition"]; exists {
					fmt.Printf("   • %s\n", def)
				}
				if dates, exists := info["dates"]; exists {
					for _, event := range dates.([]string) {
						fmt.Printf("   • %s\n", event)
					}
				}
				if causes, exists := info["causes"]; exists {
					for _, effect := range causes.([]string) {
						fmt.Printf("   • cause: %s\n", effect)
					}
				}
			}
		}
	}

	if !foundKnowledge {
		fmt.Println("\n❌ Aucune connaissance trouvée pour répondre à cette question.")
		fmt.Println("\n💡 La base de connaissances doit être enrichie avec:")
		fmt.Println("   ./programme learn <fichier_pertinent.txt>")
	}

	// Calculer le boost de confiance
	boost := GlobalKnowledgeBase.GetConfidenceBoost(question)
	fmt.Printf("\n📊 Boost de confiance: %.1f%%\n", boost*100)

	fmt.Println("============================================================\n")
}
