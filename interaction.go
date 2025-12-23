package main

import (
	"IA-ATOMIQUE/database"
	"bufio"
	"fmt"
	"os"
	"strings"
)

// TraiterFichier traite un fichier texte
func TraiterFichier(cheminFichier string) {
	contenu, err := os.ReadFile(cheminFichier)
	if err != nil {
		fmt.Printf("[ERREUR] Impossible de lire le fichier: %v\n", err)
		return
	}

	texte := string(contenu)
	TraiterTexte(texte, cheminFichier)
}

// TraiterTexte analyse et traite un texte complet - Approche phrase par phrase
func TraiterTexte(texte string, source string) {
	fmt.Printf("\n╔════════════════════════════════════════╗\n")
	fmt.Printf("║  TRAITEMENT - %s\n", source)
	fmt.Printf("╚════════════════════════════════════════╝\n\n")

	// Nouvelle approche : analyser phrase par phrase
	analyse := database.AnalyserParPhrases(texte)

	// Statistiques globales
	stats := database.StatistiquesAnalyse(analyse)
	comptage := stats["comptage"].(map[string]int)
	energieTotale := stats["energie_totale"].(float64)
	confianceMoyenne := stats["confiance_moyenne"].(float64)
	nbPhrases := stats["nb_phrases"].(int)

	fmt.Printf("[STATISTIQUES GLOBALES]\n")
	fmt.Printf("• Phrases analysées: %d\n", nbPhrases)
	fmt.Printf("• Catégories détectées: %d\n", len(comptage))
	fmt.Printf("• Énergie totale: %.2f\n", energieTotale)
	fmt.Printf("• Confiance moyenne: %.1f%%\n\n", confianceMoyenne*100)

	// Distribution par catégorie
	dominantes := database.CategoriesDominantes(analyse)
	if len(dominantes) > 0 {
		fmt.Printf("[DISTRIBUTION]\n")
		for _, catData := range dominantes {
			cat := catData["categorie"].(string)
			nombre := catData["nombre"].(int)
			pourcentage := catData["pourcentage"].(float64)
			barLen := int(pourcentage / 5)
			bar := ""
			for i := 0; i < barLen; i++ {
				bar += "█"
			}
			fmt.Printf("  %-12s %s %.0f%% (%d phrases)\n",
				cat, bar, pourcentage, nombre)
		}
		fmt.Println()
	}

	// Afficher l'analyse détaillée par catégorie
	fmt.Print(database.AfficherAnalyseDetaillee(analyse))

	// RÉSUMÉ SYNTHÉTIQUE
	fmt.Printf("\n╔════════════════════════════════════════╗\n")
	fmt.Printf("║  RÉSUMÉ SYNTHÉTIQUE (ses propres mots)\n")
	fmt.Printf("╚════════════════════════════════════════╝\n\n")
	resumeSynthetique := database.GenererResumeSynthetique(analyse)
	fmt.Printf("%s\n\n", resumeSynthetique)

	// Apprentissage multi-catégories
	tokens := database.TokeniserTexte(texte)

	// Apprendre pour chaque phrase dans sa catégorie
	for _, phraseAnalyse := range analyse.Phrases {
		if phraseAnalyse.CategorieID > 0 {
			phrasTokens := database.TokeniserTexte(phraseAnalyse.Texte)
			for _, token := range phrasTokens {
				database.Apprendre(token, phraseAnalyse.CategorieID)
			}

			// Activer les neurones
			for i := range database.Neurones {
				if database.Neurones[i].CategorieID == phraseAnalyse.CategorieID {
					database.Neurones[i].Valeur += phraseAnalyse.Score
				}
			}
		}
	}

	fmt.Printf("[APPRENTISSAGE]\n✓ %d tokens analysés et mémorisés\n\n", len(tokens))

	VisualiserStats()
	fmt.Println("✓ Dashboard mis à jour")
}

// InteractionInteractive permet à l'utilisateur d'interagir
func InteractionInteractive() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf(`
╔═══════════════════════════════════════╗
║   IA-ATOMIQUE v4.2 - Interface       ║
║   Analyse Phrase par Phrase           ║
╚═══════════════════════════════════════╝

Commandes:
  file <chemin>     - Traiter un fichier
  text              - Entrer du texte libre
  ask               - Poser une question
  stats             - Voir les statistiques
  quit              - Quitter

`)

	for {
		fmt.Print("\n> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "" {
			continue
		}

		parts := strings.Fields(input)
		commande := parts[0]

		switch commande {
		case "quit", "exit":
			fmt.Println("Au revoir!")
			return

		case "file":
			if len(parts) > 1 {
				TraiterFichier(parts[1])
			} else {
				fmt.Println("[ERREUR] Syntaxe: file <chemin>")
			}

		case "text":
			fmt.Print("Entrez votre texte (terminez par une ligne vide):\n> ")
			var texte string
			for {
				ligne, _ := reader.ReadString('\n')
				if strings.TrimSpace(ligne) == "" {
					break
				}
				texte += ligne
			}
			if texte != "" {
				TraiterTexte(texte, "entrée utilisateur")
			}

		case "ask":
			if len(parts) > 1 {
				question := strings.Join(parts[1:], " ")
				fmt.Printf("\n[QUESTION]\n%s\n", question)

				catAct, mots, conf := database.ProcesserTexte(question)
				var cat int
				var score int
				for c, s := range catAct {
					if s > score {
						score = s
						cat = c
					}
				}

				if cat > 0 {
					reponse := database.GenererReponse(cat, nil)
					fmt.Printf("\n[REPONSE]\n%s\n", reponse)
					if len(mots) > 0 {
						maxMots := 3
						if len(mots) < 3 {
							maxMots = len(mots)
						}
						fmt.Printf("\nMots clés: %s\n", strings.Join(mots[:maxMots], ", "))
					}
					fmt.Printf("Confiance: %.0f%%\n", conf*100)
				}
			} else {
				fmt.Println("[ERREUR] Syntaxe: ask <question>")
			}

		case "stats":
			VisualiserStats()
			fmt.Println("\n✓ Dashboard généré")

		default:
			TraiterTexte(input, "entrée directe")
		}
	}
}
