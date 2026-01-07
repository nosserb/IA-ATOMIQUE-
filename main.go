package main

import (
	"IA-ATOMIQUE/database"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

var stopTrain = false

func VisualiserStats() {
	f, _ := os.Create("dashboard")
	defer f.Close()

	f.WriteString("===========================================\n")
	f.WriteString("   IA-ATOMIQUE v4.0 - NEURAL MONITOR\n")
	f.WriteString("===========================================\n\n")

	// STATISTIQUES NEURONES
	actifs := 0
	totalEnergie := 0.0
	maxValeur := 0.0
	for _, n := range database.Neurones {
		if n.Valeur > 0 {
			actifs++
			totalEnergie += n.Valeur
		}
		if n.Valeur > maxValeur {
			maxValeur = n.Valeur
		}
	}

	f.WriteString(fmt.Sprintf("[RESEAU NEURAL]\n"))
	f.WriteString(fmt.Sprintf("Total neurones: 1000\n"))
	f.WriteString(fmt.Sprintf("Neurones actifs: %d (%.1f%%)\n", actifs, float64(actifs)*100/1000))
	f.WriteString(fmt.Sprintf("Energie totale: %.2f\n", totalEnergie))
	f.WriteString(fmt.Sprintf("Energie max: %.2f\n\n", maxValeur))

	// DISTRIBUTION PAR CATEGORIE
	f.WriteString("[DISTRIBUTION CATEGORIES]\n")
	catStats := make(map[int]map[string]interface{})
	for _, n := range database.Neurones {
		if _, ok := catStats[n.CategorieID]; !ok {
			catStats[n.CategorieID] = map[string]interface{}{"count": 0, "energy": 0.0, "active": 0}
		}
		catStats[n.CategorieID]["count"] = catStats[n.CategorieID]["count"].(int) + 1
		catStats[n.CategorieID]["energy"] = catStats[n.CategorieID]["energy"].(float64) + n.Valeur
		if n.Valeur > 0 {
			catStats[n.CategorieID]["active"] = catStats[n.CategorieID]["active"].(int) + 1
		}
	}

	for id := 1; id <= 50; id++ {
		if stats, ok := catStats[id]; ok {
			count := stats["count"].(int)
			energy := stats["energy"].(float64)
			activeCount := stats["active"].(int)
			f.WriteString(fmt.Sprintf("CAT %2d: ", id))
			// Barre d'energie
			barLen := int(energy / maxValeur * 30)
			if barLen > 30 {
				barLen = 30
			}
			for i := 0; i < barLen; i++ {
				f.WriteString("█")
			}
			for i := barLen; i < 30; i++ {
				f.WriteString("░")
			}
			f.WriteString(fmt.Sprintf(" [%3d/%3d] Energy: %.2f\n", activeCount, count, energy))
		}
	}

	f.WriteString("\n[TOP NEURONES ACTIFS]\n")
	type NeuroneInfo struct {
		id  int
		cat int
		val float64
	}
	var topNeurones []NeuroneInfo
	for _, n := range database.Neurones {
		if n.Valeur > 0 {
			topNeurones = append(topNeurones, NeuroneInfo{n.ID, n.CategorieID, n.Valeur})
		}
	}
	// Sort by value
	for i := 0; i < len(topNeurones) && i < 20; i++ {
		maxIdx := i
		for j := i + 1; j < len(topNeurones); j++ {
			if topNeurones[j].val > topNeurones[maxIdx].val {
				maxIdx = j
			}
		}
		topNeurones[i], topNeurones[maxIdx] = topNeurones[maxIdx], topNeurones[i]
	}

	for i := 0; i < len(topNeurones) && i < 20; i++ {
		n := topNeurones[i]
		f.WriteString(fmt.Sprintf("N%04d [CAT %2d] %.2f ", n.id, n.cat, n.val))
		barLen := int(n.val / maxValeur * 25)
		if barLen > 25 {
			barLen = 25
		}
		for j := 0; j < barLen; j++ {
			f.WriteString("▓")
		}
		f.WriteString("\n")
	}

	f.WriteString("\n[ETAT SYSTEM]\n")
	f.WriteString(fmt.Sprintf("Timestamp: %s\n", time.Now().Format("2006-01-02 15:04:05")))
	f.WriteString(fmt.Sprintf("Status: OPERATIONAL\n"))
}

func AutoIntrospection() {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for !stopTrain {
		<-ticker.C

		// Degrader neurones
		for i := range database.Neurones {
			database.Neurones[i].Valeur *= 0.95
		}

		// Aleatoire activation
		randCat := rand.Intn(50) + 1
		for i := range database.Neurones {
			if database.Neurones[i].CategorieID == randCat {
				database.Neurones[i].Valeur += float64(rand.Intn(3))
			}
		}

		// Sauvegarder stats
		VisualiserStats()
	}
}

func afficherAide() {
	fmt.Println(`
╔═══════════════════════════════════════════════════════════════╗
║          IA-ATOMIQUE v4.1 - Neural Text Analysis             ║
║    Analyse de texte, Apprentissage & Génération de résumés   ║
╚═══════════════════════════════════════════════════════════════╝

📖 UTILISATION:

  1. MODE INTERACTIF (par défaut)
     $ ./programme
     ou
     $ ./programme interactive

  2. ANALYSER UN FICHIER
     $ ./programme file <chemin_fichier>
     
     Exemples:
       $ ./programme file wikipedia.txt
       $ ./programme file document.md
       $ ./programme file /home/user/texte.txt

  3. HUMANISER UN TEXTE
     $ ./programme humanize [file|-s|-p|-a] <chemin_fichier>
     
     Styles disponibles:
       -s : Standard (naturel et fluide) - défaut
       -p : Professionnel (formel et technique)
       -a : Avancé (analyse de style + paraphrase intelligente)
     
     Exemples (tous les formats sont supportés):
       $ ./programme humanize file document.txt
       $ ./programme humanize file -s document.txt
       $ ./programme humanize file -p document.txt
       $ ./programme humanize file -a document.txt
       $ ./programme humanize -s file document.txt
       $ ./programme humanize -p file document.txt
       $ ./programme humanize -a file document.txt
     
     Résultat: Crée un fichier "_humanized.txt", "_humanized_prof.txt" ou "_humanized_avance.txt"

  4. ANALYSER UN TEXTE DIRECT
     $ ./programme text <votre texte>
     
     Exemples:
       $ ./programme text "L'IA est la technologie du futur"
       $ ./programme text "Einstein a découvert la relativité en 1905"

  5. MODE CLASSIQUE (hérité)
     $ ./programme <texte quelconque>
     
     Exemples:
       $ ./programme bonjour
       $ ./programme l'informatique progresse rapidement

═══════════════════════════════════════════════════════════════

⚙️  FONCTIONNALITÉS:

  ✓ Analyse phrase par phrase
  ✓ Classification en 6 catégories:
    - TECH: Technologie, informatique, numérique
    - HISTOIRE: Politique, histoire, événements
    - BUSINESS: Commerce, économie, affaires
    - ALIMENTATION: Nutrition, gastronomie
    - SANTE: Santé, médecine, bien-être
    - VERBE: Détection d'actions (verbes principaux)

  ✓ Détection de la structure grammaticale
  ✓ Génération de résumés synthétiques
  ✓ Système de modération (blacklist chiffrée AES-256)
  ✓ Apprentissage dynamique des mots
  ✓ Humanisation de texte
    - Améliore la fluidité des phrases
    - Remplace les formulations maladroites
    - Ajoute des connecteurs naturels
    - Optimise la structure et la longueur

═══════════════════════════════════════════════════════════════

📊 SORTIE STANDARD:

  [STATISTIQUES GLOBALES]
    • Phrases analysées
    • Catégories détectées
    • Énergie totale
    • Confiance moyenne

  [DISTRIBUTION]
    • Répartition en pourcentages
    • Graphiques visuels

  [ANALYSE PAR PHRASE]
    • Catégorie détectée
    • Mots clés identifiés
    • Verbes principaux
    • Résumé synthétique

═══════════════════════════════════════════════════════════════

💡 CONSEILS D'UTILISATION:

  • Pour les meilleurs résultats, utilisez des textes français
  • Les fichiers volumineux seront traités progressivement
  • La modération empêche l'apprentissage de contenu offensant
  • Les neurones se régénèrent automatiquement après chaque exécution

═══════════════════════════════════════════════════════════════

❓ AIDE:
  $ ./programme help
  ou
  $ ./programme -h
  ou
  $ ./programme --help

═══════════════════════════════════════════════════════════════
`)
}

func main() {
	// Mode CLI (terminal)
	go AutoIntrospection()
	go AutoIntrospection()

	if len(os.Args) < 2 {
		fmt.Println("[IA-ATOMIQUE v1.0] - Technologie de Résonance Atomique")
		fmt.Println("Interface Interactive - Réseau Atomique Distribué")
		RunAtomicDemo()
		return
	}

	// Vérifier les commandes atomiques en priorité
	commande := os.Args[1]

	// Commandes du réseau atomique
	if commande == "simulate" || commande == "network-stats" || commande == "benchmark" ||
		(commande == "help" && len(os.Args) == 2) {
		ParseSimulationArgs(os.Args)
		return
	}

	// Vérifier les autres commandes
	switch commande {
	case "help", "-h", "--help":
		afficherAide()

	case "file":
		if len(os.Args) > 2 {
			TraiterFichier(os.Args[2])
		} else {
			fmt.Println("[ERREUR] Utilisation: ./programme file <chemin>")
		}

	case "humanize":
		style := "standard" // Style par défaut
		filepath := ""

		// Vérifier les différentes syntaxes possibles
		if len(os.Args) > 3 {
			// Cas 1: humanize file <chemin>
			// Cas 2: humanize file -s <chemin>
			// Cas 3: humanize file -p <chemin>
			// Cas 4: humanize file -a <chemin>
			// Cas 5: humanize file <chemin> -s
			// Cas 6: humanize file <chemin> -p
			// Cas 7: humanize file <chemin> -a
			// Cas 8: humanize -s file <chemin>
			// Cas 9: humanize -p file <chemin>
			// Cas 10: humanize -a file <chemin>

			if os.Args[2] == "file" {
				// Cas 1-7: "file" est le second argument
				if len(os.Args) > 3 {
					if os.Args[3] == "-s" || os.Args[3] == "-p" || os.Args[3] == "-a" {
						// Cas 2-4: flag avant le chemin
						if len(os.Args) > 4 {
							if os.Args[3] == "-p" {
								style = "professionnel"
							} else if os.Args[3] == "-a" {
								style = "avance"
							}
							filepath = os.Args[4]
						}
					} else {
						// Cas 1 ou 5-7: chemin avant le flag
						filepath = os.Args[3]
						if len(os.Args) > 4 && (os.Args[4] == "-s" || os.Args[4] == "-p" || os.Args[4] == "-a") {
							if os.Args[4] == "-p" {
								style = "professionnel"
							} else if os.Args[4] == "-a" {
								style = "avance"
							}
						}
					}
				}
			} else if (os.Args[2] == "-s" || os.Args[2] == "-p" || os.Args[2] == "-a") && os.Args[3] == "file" {
				// Cas 8-10: flag avant "file"
				if os.Args[2] == "-p" {
					style = "professionnel"
				} else if os.Args[2] == "-a" {
					style = "avance"
				}
				if len(os.Args) > 4 {
					filepath = os.Args[4]
				}
			}
		}

		if filepath != "" {
			TraiterFichierHumanize(filepath, style)
		} else {
			fmt.Println("[ERREUR] Utilisation: ./programme humanize [file|-s|-p] <chemin>")
			fmt.Println("  Formats supportés:")
			fmt.Println("    ./programme humanize file <chemin>")
			fmt.Println("    ./programme humanize file -s <chemin>")
			fmt.Println("    ./programme humanize file -p <chemin>")
			fmt.Println("    ./programme humanize -s file <chemin>")
			fmt.Println("    ./programme humanize -p file <chemin>")
		}

	case "text":
		if len(os.Args) > 2 {
			texte := strings.Join(os.Args[2:], " ")
			TraiterTexte(texte, "entrée CLI")
		} else {
			fmt.Println("[ERREUR] Utilisation: ./programme text <votre texte>")
		}

	case "interactive", "inter":
		InteractionInteractive()

	default:
		// Mode hérité - traiter comme phrase simple
		phrase := strings.Join(os.Args[1:], " ")
		TraiterTexte(phrase, "mode classique")
	}

	database.RegenererNeurones()
}
