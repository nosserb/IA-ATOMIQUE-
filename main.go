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

🧪 SIMULATION RÉSEAU ATOMIQUE (T.R.A.):

  1. LANCER UNE SIMULATION ATOMIQUE
     $ ./programme simulate <itérations> [nombre_atomes]
     
     Exemples:
       $ ./programme simulate 100           (100 itérations, 500 atomes)
       $ ./programme simulate 500 1000      (500 itérations, 1000 atomes)
       $ ./programme simulate 1000 2000     (1000 itérations, 2000 atomes)
     
     Affiche:
       ✓ Cohérence du réseau (initiale/finale/moyenne)
       ✓ Activation moyenne des atomes
       ✓ Consommation énergétique
       ✓ Système de freeze (atomes en hibernation)
       ✓ Comportements émergents
       ✓ Performance (itérations/seconde)

  2. STATISTIQUES RÉSEAU
     $ ./programme network-stats
     
     Affiche les métriques détaillées du dernier réseau simulé

  3. BENCHMARK DE PERFORMANCE
     $ ./programme benchmark
     
     Teste les performances avec différentes tailles de réseau

═══════════════════════════════════════════════════════════════

⚡ DÉCOUPAGE ULTRA-RAPIDE O(n):

  1. TESTER LE DÉCOUPAGE RAPIDE
     $ ./programme split <fichier> [mots_par_bloc]
     
     Exemples:
       $ ./programme split document.txt           (blocs de 100 mots)
       $ ./programme split document.txt 50        (blocs de 50 mots)
       $ ./programme split document.txt 200       (blocs de 200 mots)
     
     Affiche:
       ✓ Nombre de blocs générés
       ✓ Statistiques (mots, taille min/max)
       ✓ Confiance moyenne
       ✓ Vitesse (mots/seconde)
     
     Performance: O(n) - Instantané même pour 50+ MB de texte
     Exemple: 57 000 mots en 27ms (2,1M mots/sec)

═══════════════════════════════════════════════════════════════

⚛️  TECHNOLOGIE DE RÉSONANCE ATOMIQUE (T.R.A.):

  • Réseau décentralisé asynchrone
  • Chaque atome opère indépendamment
  • Résonance entre atomes voisins
  • Dynamique adaptative des poids (apprentissage)
  • Système de freeze pour sobriété énergétique
  • Émergence de comportements globaux

  Équations principales:
    - Résonance: R(si,sj) = exp(-||si-sj||²/2σ²)
    - État: si(t+1) = si(t) + α·Σ(wij·Rij) + β·(Ri+pi)
    - Poids: dwij/dt = γ·cohérence(si,sj) - δ·wij
    - Freeze: Alocal = wa|Δa| + wc|Δc| + wR·R + wE|ΔE|

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

	// Commande de découpage rapide
	if commande == "split" || commande == "decoupe" {
		if len(os.Args) > 2 {
			TestDecoupageRapide(os.Args[2:])
		} else {
			fmt.Println("[ERREUR] Utilisation: ./programme split <fichier> [mots_par_bloc]")
		}
		return
	}

	// Commande d'extraction de phrases clés
	if commande == "extract" || commande == "phrases" {
		if len(os.Args) > 2 {
			ExtrairePhrasesClésCommand(os.Args[2:])
		} else {
			fmt.Println("[ERREUR] Utilisation: ./programme extract <fichier> [ratio_conservation]")
		}
		return
	}

	// Commande de génération de résumé
	if commande == "generate" {
		if len(os.Args) > 2 {
			filepath := os.Args[2]
			ratio := 0.3 // 30% par défaut
			if len(os.Args) > 3 {
				fmt.Sscanf(os.Args[3], "%f", &ratio)
			}
			GenererResumeCommand(filepath, ratio)
		} else {
			fmt.Println("[ERREUR] Utilisation: ./programme generate <fichier> [ratio_compression=0.3]")
		}
		return
	}

	// Commande de résumé - NOW USING PHASE 15 (Grammar-Aware Summarization)
	if commande == "resume" {
		if len(os.Args) > 2 {
			filepath := os.Args[2]
			threshold := 0.10 // 10% default
			if len(os.Args) > 3 {
				fmt.Sscanf(os.Args[3], "%f", &threshold)
			}
			// Use Phase 15 optimized pipeline
			resumeOptimizedCommand([]string{filepath, fmt.Sprintf("%f", threshold)})
		} else {
			fmt.Println("[ERREUR] Utilisation: ./programme resume <fichier> [threshold=0.10]")
		}
		return
	}

	// Commande de réécriture (humanisation)
	if commande == "rewrite" {
		style := "standard"
		if len(os.Args) > 2 {
			texte := os.Args[2]
			if len(os.Args) > 3 && (os.Args[3] == "standard" || os.Args[3] == "professionnel" || os.Args[3] == "avance") {
				style = os.Args[3]
			}
			ReecrireCommand(texte, style)
		} else {
			fmt.Println("[ERREUR] Utilisation: ./programme rewrite \"<texte>\" [standard|professionnel|avance]")
		}
		return
	}

	// Commande de génération atomique autonome (NEW!)
	if commande == "atomic" {
		if len(os.Args) > 2 {
			filepath := os.Args[2]
			compression := 0.3
			if len(os.Args) > 3 {
				fmt.Sscanf(os.Args[3], "%f", &compression)
			}
			GenererAvecAtomesCommand(filepath, compression)
		} else {
			fmt.Println("[ERREUR] Utilisation: ./programme atomic <fichier> [compression=0.3]")
		}
		return
	}
	// === PHASE 15: Grammar-Aware Summarization ===
	if commande == "resume-optimized" {
		resumeOptimizedCommand(os.Args[2:])
		return
	}
	if commande == "compare-summaries" {
		compareSummariesCommand(os.Args[2:])
		return
	}
	if commande == "analyze-preprocessing" {
		analyzePreprocessingCommand(os.Args[2:])
		return
	}
	if commande == "analyze-vocabulary" {
		analyzeVocabularyCommand(os.Args[2:])
		return
	}
	// Commande de comparaison vectoriel vs atomique (NEW!)
	if commande == "compare" {
		if len(os.Args) > 2 {
			filepath := os.Args[2]
			compression := 0.3
			if len(os.Args) > 3 {
				fmt.Sscanf(os.Args[3], "%f", &compression)
			}
			GenererComparatifCommand(filepath, compression)
		} else {
			fmt.Println("[ERREUR] Utilisation: ./programme compare <fichier> [compression=0.3]")
		}
		return
	}

	// === PHASE 14: Commande syntaxe avancée avec renforcement ===
	if commande == "syntax" {
		if len(os.Args) > 2 {
			syntaxCommand(os.Args[2:])
		} else {
			fmt.Println("[INFO] Utilisation: ./programme syntax <subcommand> <text>")
			fmt.Println("Subcommands: analyze, enhance, paragraph, pos")
		}
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

// TestDecoupageRapide teste et affiche les résultats du découpage ultra-rapide
func TestDecoupageRapide(args []string) {
	if len(args) == 0 {
		fmt.Println("[ERREUR] Spécifier un fichier")
		return
	}

	// Lire le fichier
	contenu, err := os.ReadFile(args[0])
	if err != nil {
		fmt.Printf("[ERREUR] Impossible de lire %s: %v\n", args[0], err)
		return
	}

	texte := string(contenu)

	// Déterminer la taille des blocs
	motsParBloc := 100
	if len(args) > 1 {
		fmt.Sscanf(args[1], "%d", &motsParBloc)
	}

	fmt.Printf("\n╔════════════════════════════════════════════════════════╗\n")
	fmt.Printf("║     DÉCOUPAGE ULTRA-RAPIDE O(n) - ANALYSE DE BLOCS    ║\n")
	fmt.Printf("╚════════════════════════════════════════════════════════╝\n\n")

	// Lancer le chrono
	debut := time.Now()

	// Découper le texte
	blocs := database.DecouperTexteRapide(texte, motsParBloc)

	// Analyser chaque bloc
	var allResultats []map[int]int
	var allConfiances []float64

	for i, bloc := range blocs {
		resultat, motsCles, confiance := database.AnalyserBloc(bloc)
		allResultats = append(allResultats, resultat)
		allConfiances = append(allConfiances, confiance)

		// Limiter l'affichage des mots-clés
		affiChage := strings.Join(motsCles, ", ")
		if len(affiChage) > 50 {
			affiChage = affiChage[:50] + "..."
		}

		if i < 5 || i >= len(blocs)-2 { // Afficher les 5 premiers et 2 derniers
			fmt.Printf("[Bloc %d] %d mots | Confiance: %.2f%% | Mots-clés: %s\n",
				bloc.NumeroBloc,
				bloc.Taille,
				confiance*100,
				affiChage,
			)
		} else if i == 5 {
			fmt.Printf("...\n")
		}
	}

	// Fusionner les résultats
	catGlobale := database.FusionnerResultatsBlocs(blocs, allResultats, allConfiances)

	// Afficher les statistiques
	stats := database.StatistiquesBlocs(blocs)
	temps := time.Since(debut)

	fmt.Printf("\n[STATISTIQUES DE DÉCOUPAGE]\n")
	fmt.Printf("  • Nombre de blocs: %d\n", stats["nombre_blocs"])
	fmt.Printf("  • Total de mots: %d\n", stats["total_mots"])
	fmt.Printf("  • Taille moyenne: %d mots\n", stats["taille_moyenne"])
	fmt.Printf("  • Taille min/max: %d/%d mots\n", stats["taille_min"], stats["taille_max"])

	fmt.Printf("\n[RÉSULTATS FUSIONNÉS]\n")
	fmt.Printf("  • Catégories détectées: %d\n", len(catGlobale))
	fmt.Printf("  • Confiance moyenne: %.2f%%\n", (sum(allConfiances)/float64(len(allConfiances)))*100)

	fmt.Printf("\n[PERFORMANCE]\n")
	fmt.Printf("  • Complexité: O(n)\n")
	fmt.Printf("  • Temps de découpage + analyse: %v\n", temps)
	fmt.Printf("  • Vitesse: %.0f mots/sec\n", float64(stats["total_mots"].(int))/temps.Seconds())

	fmt.Printf("\n✓ Découpage instantané réussi - Compatible avec réseau atomique distribué\n\n")
}

func sum(vals []float64) float64 {
	total := 0.0
	for _, v := range vals {
		total += v
	}
	return total
}

// ExtrairePhrasesClésCommand extrait et affiche les phrases clés d'un texte
func ExtrairePhrasesClésCommand(args []string) {
	if len(args) == 0 {
		fmt.Println("[ERREUR] Spécifier un fichier")
		return
	}

	// Lire le fichier
	contenu, err := os.ReadFile(args[0])
	if err != nil {
		fmt.Printf("[ERREUR] Impossible de lire %s: %v\n", args[0], err)
		return
	}

	texte := string(contenu)

	// Déterminer le ratio de conservation
	ratio := 0.3
	if len(args) > 1 {
		fmt.Sscanf(args[1], "%f", &ratio)
	}

	fmt.Printf("\n╔════════════════════════════════════════════════════════╗\n")
	fmt.Printf("║   EXTRACTION DE PHRASES CLÉS - ÉNERGIE ATOMIQUE       ║\n")
	fmt.Printf("╚════════════════════════════════════════════════════════╝\n\n")

	// Lancer le chrono
	debut := time.Now()

	// Extraire les phrases clés
	phrasesClés := database.ExtrairePhrasesClés(texte, ratio)

	// Analyser les résultats
	temps := time.Since(debut)

	// Statistiques
	totalPhrasesOriginal := strings.Count(texte, ".") + strings.Count(texte, "!") + strings.Count(texte, "?")
	ratioRéel := float64(len(phrasesClés)) / float64(totalPhrasesOriginal) * 100

	fmt.Printf("[RÉSUMÉ INTELLIGENT]\n\n")

	for i, phrase := range phrasesClés {
		fmt.Printf("%d. %s\n", i+1, phrase.Contenu)
		fmt.Printf("   Énergie: %.2f | Cohérence: %.2f | Importance: %.1f%%\n\n",
			phrase.Energie,
			phrase.EnergieTotal,
			(phrase.EnergieTotal/2.0)*100,
		)
	}

	fmt.Printf("\n[STATISTIQUES]\n")
	fmt.Printf("  • Phrases originales: ~%d\n", totalPhrasesOriginal)
	fmt.Printf("  • Phrases conservées: %d (%.1f%%)\n", len(phrasesClés), ratioRéel)
	fmt.Printf("  • Ratio demandé: %.1f%%\n", ratio*100)
	fmt.Printf("  • Densité d'information: HIGH\n")

	fmt.Printf("\n[PERFORMANCE]\n")
	fmt.Printf("  • Temps d'extraction: %v\n", temps)
	fmt.Printf("  • Pipeline: Découpage → Énergie → Cohérence → Filtrage → Fusion\n")

	// Calculer le taux de compression
	tailleOriginale := len(strings.Fields(texte))
	tailleRésumé := 0
	for _, p := range phrasesClés {
		tailleRésumé += len(p.Mots)
	}
	compression := float64(tailleOriginale) / float64(tailleRésumé)

	fmt.Printf("  • Taux de compression: %.1fx (de %d à %d mots)\n", compression, tailleOriginale, tailleRésumé)

	fmt.Printf("\n✓ Extraction réussie - Phrases clés sélectionnées par énergie atomique\n\n")
}
