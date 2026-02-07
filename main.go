package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/nosserb/IA-ATOMIQUE-/database"
	"github.com/nosserb/IA-ATOMIQUE-/internal/commands"
	"github.com/nosserb/IA-ATOMIQUE-/internal/tests"
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

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🎯 COMMANDES PRINCIPALES
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

1. MODE INTERACTIF (par défaut)
   $ ./programme
   $ ./programme interactive
   → Lance l'interface interactive avec menu

2. ANALYSER UN FICHIER
   $ ./programme file <chemin_fichier>
   
   Exemples:
     $ ./programme file document.txt
     $ ./programme file /path/to/file.md

3. ANALYSER UN TEXTE DIRECT
   $ ./programme text <votre texte>
   
   Exemples:
     $ ./programme text "L'IA est l'avenir"
     $ ./programme text "Bonjour le monde"

4. MODE HÉRITÉ (texte simple)
   $ ./programme <votre texte>
   
   Exemples:
     $ ./programme bonjour
     $ ./programme l'informatique progresse

5. INTERFACE WEB
	$ ./programme web
	→ Lance l'interface web locale

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🎓 APPRENTISSAGE AUTOMATIQUE
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

1. APPRENDRE À PARTIR DE TEXTES
   $ ./programme learn <fichier.txt|dossier>
   
   Exemples:
     $ ./programme learn histoire.txt
     $ ./programme learn corpus/sciences/
     $ ./programme learn medecine.txt
   
   → Extrait automatiquement des connaissances factuelles
   → Patterns: dates, relations causales, définitions
   → Co-occurrences: associations statistiques

2. CONSULTER LES CONNAISSANCES
   $ ./programme knowledge <terme>
   
   Exemples:
     $ ./programme knowledge Napoléon
     $ ./programme knowledge 1815
     $ ./programme knowledge coeur
   
   → Affiche toutes les connaissances sur un terme

3. STATISTIQUES D'APPRENTISSAGE
   $ ./programme stats-kb
   $ ./programme kb-stats
   
   → Affiche les statistiques de la base de connaissances
   → Nombre de faits extraits, définitions, dates, etc.

4. TESTER LES CONNAISSANCES
   $ ./programme test-knowledge <question>
   
   Exemple:
     $ ./programme test-knowledge "Quand a eu lieu Waterloo?"
   
   → Cherche les connaissances pertinentes pour répondre

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
📝 RÉSUMÉS & GÉNÉRATION
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

1. RÉSUMÉ HYBRIDE (Atomique + Proba, Phase 15)
	 $ ./programme resume <fichier> [ratio=0.55]
   
	 Exemples:
		 $ ./programme resume document.txt         (complet, 55%)
		 $ ./programme resume document.txt 0.70    (très complet, 70%)
		 $ ./programme resume document.txt 0.40    (plus condensé, 40%)
   
	 → Fusion atomique + scoring probabiliste + réécriture fluide

2. GÉNÉRATION DE RÉSUMÉ (simple)
   $ ./programme generate <fichier> [ratio=0.3]
   
   Exemples:
     $ ./programme generate document.txt
     $ ./programme generate document.txt 0.5
   
   → Génère un résumé par extraction vectorielle

3. EXTRACTION DES PHRASES CLÉS
   $ ./programme extract <fichier> [ratio=0.3]
   
   Exemples:
     $ ./programme extract document.txt
     $ ./programme extract document.txt 0.5
   
   → Identifie et extrait les phrases les plus importantes

4. GÉNÉRATION AVEC ATOMES
   $ ./programme atomic <fichier> [compression=0.3]
   
   Exemples:
     $ ./programme atomic document.txt
     $ ./programme atomic document.txt 0.5
   
   → Utilise le réseau atomique autonome pour générer

5. COMPARATIF DES MÉTHODES
	$ ./programme compare <fichier> [compression=0.3]
	$ ./programme compare-summaries <fichier>
	→ Compare vectoriel vs atomique vs grammatical

6. ANALYSES INTERNES
	$ ./programme analyze-preprocessing <fichier>
	$ ./programme analyze-vocabulary <fichier>
	→ Affiche les détails du pipeline et du vocabulaire

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
✍️  HUMANISATION & RÉÉCRITURE
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

RÉÉCRIRE UN TEXTE
$ ./programme rewrite "<texte>" [style]

Styles disponibles:
  - standard    : Naturel et fluide (par défaut)
  - professionnel : Formel et technique
  - avance      : Analyse de style + paraphrase intelligente

Exemples:
  $ ./programme rewrite "C'est pas bon"
  $ ./programme rewrite "C'est pas bon" standard
  $ ./programme rewrite "C'est pas bon" professionnel
  $ ./programme rewrite "C'est pas bon" avance

HUMANISER UN FICHIER
$ ./programme humanize file <chemin> [-s|-p|-a]

Exemples:
  $ ./programme humanize file document.txt
  $ ./programme humanize file document.txt -s
  $ ./programme humanize file document.txt -p
  $ ./programme humanize file document.txt -a

Résultats: _humanized.txt, _humanized_prof.txt, _humanized_avance.txt

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🔬 ANALYSE AVANCÉE
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

RÉSUMÉ OPTIMISÉ AVANCÉ
$ ./programme resume-optimized <fichier> [threshold=0.10]

Exemples:
  $ ./programme resume-optimized document.txt
  $ ./programme resume-optimized document.txt 0.3

COMPARER LES APPROCHES DE RÉSUMÉ
$ ./programme compare-summaries <fichier>

Exemples:
  $ ./programme compare-summaries document.txt

→ Compare vectoriel vs atomique vs grammatical

ANALYSER LE PRÉTRAITEMENT
$ ./programme analyze-preprocessing <fichier>

Exemples:
  $ ./programme analyze-preprocessing document.txt

→ Détail du nettoyage, tokenization, normalisation

ANALYSER LE VOCABULAIRE
$ ./programme analyze-vocabulary <fichier>

Exemples:
  $ ./programme analyze-vocabulary document.txt

→ Statistiques sur les mots, fréquences, catégories

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🧪 RÉSEAU ATOMIQUE (T.R.A. - Atomic Resonance Technology)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

LANCER UNE SIMULATION
$ ./programme atomic simulate <itérations> [atomes=500]

Exemples:
  $ ./programme atomic simulate 100
  $ ./programme atomic simulate 500 1000
  $ ./programme atomic simulate 1000 2000

Métriques affichées:
  ✓ Cohérence initiale/finale/moyenne
  ✓ Activation moyenne des atomes
  ✓ Consommation énergétique
  ✓ Freeze system (hibernation)
  ✓ Comportements émergents
  ✓ Performance (iter/sec)

STATISTIQUES RÉSEAU
$ ./programme atomic network-stats

→ Distribution d'énergie, atomes actifs, métriques globales

BENCHMARK DES ATOMES
$ ./programme atomic benchmark <itérations> [atomes]

→ Performance, latence, throughput

AUTRES SOUS-COMMANDES ATOMIQUES
$ ./programme cellular <options>
$ ./programme relaxation <options>
$ ./programme relax-opt <options>
→ Commandes avancées (voir documentation)

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
📊 FONCTIONNALITÉS D'ANALYSE
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

✓ Analyse phrase par phrase
✓ Classification en 6 catégories:
  - TECH       : Technologie, informatique, numérique
  - HISTOIRE   : Politique, histoire, événements
  - BUSINESS   : Commerce, économie, affaires
  - ALIMENTATION: Nutrition, gastronomie
  - SANTE      : Santé, médecine, bien-être
  - VERBE      : Détection d'actions (verbes)

✓ Structure grammaticale (SVC - Subject-Verb-Complement)
✓ Résumés avec 5 piliers scientifiques:
  - CONTEXTE : De quoi parle-t-on?
  - PROBLÈME : Pourquoi c'est insuffisant?
  - OBJECTIF : Qu'essaie-t-on de faire?
  - APPROCHE : Comment c'est fait?
  - APPORT   : Pourquoi c'est nouveau?

✓ Modération (blacklist AES-256 chiffrée)
✓ Apprentissage dynamique
✓ Humanisation intelligente
✓ Réseau de neurones (1000 neurones × 50+ catégories)
✓ Réseau atomique autonome (T.R.A.)

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚙️  DÉCOUPAGE ULTRA-RAPIDE O(n)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

TESTER LE DÉCOUPAGE RAPIDE
$ ./programme split <fichier> [mots_par_bloc=100]

Exemples:
  $ ./programme split document.txt
  $ ./programme split document.txt 50
  $ ./programme split document.txt 200

Affiche:
  ✓ Nombre de blocs générés
  ✓ Statistiques (mots, taille min/max)
  ✓ Confiance moyenne
  ✓ Vitesse (mots/seconde)

Performance: O(n) - Instantané même pour 50+ MB de texte

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
💡 CONSEILS D'UTILISATION
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

• Meilleurs résultats avec du texte français
• Les fichiers volumineux sont traités progressivement
• La modération empêche l'apprentissage de contenu offensant
• Les neurones se régénèrent automatiquement
• Les résumés scientifiques recherchent 100% de complétude

EXEMPLES RAPIDES:
  $ ./programme file document.txt
  $ ./programme resume document.txt 0.1
  $ ./programme text "Votre texte ici"
  $ ./programme rewrite "C'est pas bon" professionnel
  $ ./programme atomic simulate 100
  $ ./programme compare-summaries document.txt

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🏁 BENCHMARKS
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

BENCHMARKS ACADÉMIQUES
$ ./programme academic all
$ ./programme academic mmlu
$ ./programme academic hellaswag
$ ./programme bench-academic
→ Lance les tests MMLU/Hellaswag (échantillons intégrés)

BENCHMARK TEXTE MASSIF
$ ./programme benchmark-text
$ ./programme bench-1m
→ Traite 1M de mots pour mesurer la vitesse

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚛️  TECHNOLOGIE DE RÉSONANCE ATOMIQUE (T.R.A.)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

• Réseau décentralisé asynchrone
• Chaque atome opère indépendamment
• Résonance entre atomes voisins
• Dynamique adaptative des poids (apprentissage)
• Système de freeze pour sobriété énergétique
• Émergence de comportements globaux

Paramètres principaux:
  α (coupling)      : Force d'interaction entre atomes
  β (local rules)   : Poids des règles locales
  γ (reinforcement) : Renforcement des poids
  δ (decay)         : Décroissance des poids
  σ (resonance)     : Sensibilité de résonance

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
❓ HELP
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

$ ./programme help
$ ./programme -h
$ ./programme --help

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🖼️  IMAGE & DÉFLOUTAGE
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

GÉNÉRATION D'IMAGES
$ ./programme image <subcommand> [...]
→ Voir "image help" pour les sous-commandes disponibles

DÉFL_OUTAGE (ATOMIC)
$ ./programme deblur ultra <image> [options]
$ ./programme deblur fast <image> [options]
$ ./programme deblur draft <image> [options]
$ ./programme deblur help

MOTION BLUR (Lucy-Richardson)
$ ./programme motion <imagePath> [output] [iterations]

COMBO (Motion + Ultra Deblur)
$ ./programme combo <imagePath> [output]

═══════════════════════════════════════════════════════════════
`)
}

func main() {
	// Mode CLI (terminal)
	// Désactiver les goroutines pour commands deblur
	// go AutoIntrospection()
	// go AutoIntrospection()

	if len(os.Args) < 2 {
		fmt.Println("[IA-ATOMIQUE v1.0] - Technologie de Résonance Atomique")
		fmt.Println("Interface Interactive - Réseau Atomique Distribué")
		commands.RunAtomicDemo()
		return
	}

	// Vérifier les commandes atomiques en priorité
	commande := os.Args[1]

	// DEBUG
	if commande == "ultra_pure" {
		fmt.Printf("DEBUG: ULTRA_PURE detected! os.Args = %v\n", os.Args)
	}

	// Commandes de génération d'images
	if commande == "image" {
		commands.ImageGenerationCommand(os.Args[1:])
		return
	}

	// Commandes d'entraînement du modèle
	if commande == "train" {
		commands.TrainingCommand(os.Args[2:])
		return
	}

	// Commandes d'émergence de patterns
	if commande == "pattern" {
		commands.PatternCommand(os.Args[2:])
		return
	}

	// Commandes de génération atomique
	if commande == "generate" {
		commands.GenerateCommand(os.Args[2:])
		return
	}

	// Commandes du réseau atomique
	if commande == "simulate" || commande == "network-stats" || commande == "benchmark" ||
		commande == "cellular" || commande == "relaxation" || commande == "relax-opt" ||
		(commande == "deblur" && (len(os.Args) < 3 || (os.Args[2] != "ultra" && os.Args[2] != "draft" && os.Args[2] != "fast"))) ||
		(commande == "help" && len(os.Args) == 2) {
		commands.ParseSimulationArgs(os.Args)
		return
	}

	// Benchmark traitement de texte 1M de mots
	if commande == "benchmark-text" || commande == "bench-1m" {
		fmt.Println("[INFO] Cette fonction a été déplacée ou supprimée")
		return
	}

	// Tests avancés (Needle In Haystack, Perplexity)
	if commande == "test" {
		tests.HandleAdvancedTests(os.Args[2:])
		return
	}

	// Benchmarks académiques (MMLU, Hellaswag)
	if commande == "academic" || commande == "bench-academic" {
		if len(os.Args) > 2 {
			tests.HandleAcademicBenchmarks(os.Args[2:])
		} else {
			tests.HandleAcademicBenchmarks([]string{})
		}
		return
	}

	// Commandes d'apprentissage automatique
	if commande == "learn" {
		commands.LearnCommand(os.Args[2:])
		return
	}
	if commande == "knowledge" {
		commands.KnowledgeCommand(os.Args[2:])
		return
	}
	if commande == "stats-kb" || commande == "kb-stats" {
		commands.StatsCommand(os.Args[2:])
		return
	}
	if commande == "test-knowledge" {
		commands.TestKnowledgeCommand(os.Args[2:])
		return
	}

	// Commande ask - poser une question à la base de connaissances
	if commande == "ask" {
		if len(os.Args) > 2 {
			question := strings.Join(os.Args[2:], " ")
			commands.AskCommand([]string{question})
		} else {
			fmt.Println("[ERREUR] Syntaxe: ask <question>")
		}
		return
	}
	if commande == "deblur" && len(os.Args) > 2 {
		mode := os.Args[2]
		switch mode {
		case "ultra":
			commands.HandleUltraFastDeblur(os.Args[3:])
			return
		case "draft":
			commands.HandleDraftFastDeblur(os.Args[3:])
			return
		case "fast":
			commands.HandleFastDeblur(os.Args[3:])
			return
		case "help", "-h":
			commands.PrintDeblurModesHelp()
			return
		}
	}

	// Commande de motion blur removal (déconvolution Lucy-Richardson)
	if commande == "motion" && len(os.Args) > 2 {
		if len(os.Args) < 3 {
			fmt.Println("Usage: ./programme motion <imagePath> [output] [iterations]")
			fmt.Println("Example: ./programme motion blurry.jpg deblurred.png 10")
			return
		}

		imagePath := os.Args[2]
		outputPath := "deblurred_motion.png"
		iterations := 15

		if len(os.Args) > 3 {
			outputPath = os.Args[3]
		}
		if len(os.Args) > 4 {
			fmt.Sscanf(os.Args[4], "%d", &iterations)
		}

		fmt.Println("\n🎬 MOTION BLUR REMOVAL (Lucy-Richardson Deconvolution)")
		fmt.Println("═════════════════════════════════════════════════════════")
		err := commands.LucyRichardsonDeconvolve(imagePath, outputPath, iterations)
		if err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			return
		}
		fmt.Printf("✅ Motion blur removed: %s\n", outputPath)
		return
	}

	// Commande combo: motion blur removal + ultra deblur (best quality)
	if commande == "combo" && len(os.Args) > 2 {
		if len(os.Args) < 3 {
			fmt.Println("Usage: ./programme combo <imagePath> [output]")
			fmt.Println("Example: ./programme combo blurry.jpg deblurred_final.png")
			fmt.Println("\nCombines motion deconvolution + ultra atomic deblurring")
			fmt.Println("Total time: ~30-40 seconds for best quality")
			return
		}

		imagePath := os.Args[2]
		finalOutput := "deblurred_combo.png"
		if len(os.Args) > 3 {
			finalOutput = os.Args[3]
		}

		// Verify image exists
		if _, err := os.Stat(imagePath); err != nil {
			fmt.Printf("❌ Image not found: %s\n", imagePath)
			return
		}

		tempFile := ".combo_motion_temp.png"

		fmt.Println("\n⚡💎 COMBO: MOTION DECONVOLUTION + ATOMIC ULTRA DEBLUR")
		fmt.Println("═══════════════════════════════════════════════════════")
		fmt.Println("Step 1/2: Motion Blur Removal (Lucy-Richardson)")
		fmt.Println("Step 2/2: Ultra Atomic Deblurring (327,000+ atoms)")
		fmt.Printf("\nProcessing: %s\n", imagePath)

		// Step 1: Motion deconvolution
		fmt.Println("\n🎬 [STEP 1] Motion deconvolution...")
		err := commands.LucyRichardsonDeconvolve(imagePath, tempFile, 18)
		if err != nil {
			fmt.Printf("❌ Motion deconv failed: %v\n", err)
			return
		}

		// Step 2: Ultra deblur on the result
		fmt.Println("\n⚛️  [STEP 2] Ultra atomic deblurring...")
		fmt.Printf("═══════════════════════════════════════════════════════\n")

		// Use the existing ultra deblur pipeline on the motion-corrected image
		commands.RunUltraDeblurPipeline4K(tempFile, 4, 4, 1000, finalOutput)

		os.Remove(tempFile)

		fmt.Printf("\n✅ COMBO COMPLETE!\n")
		fmt.Printf("💾 Saved: %s\n", finalOutput)
		fmt.Printf("✨ Motion blur removed + Atomic deblurred + Enhanced quality!\n")
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
			commands.ExtrairePhrasesClésCommand(os.Args[2:])
		} else {
			fmt.Println("[ERREUR] Utilisation: ./programme extract <fichier> [ratio_conservation]")
		}
		return
	}
	// Commande de génération d'images atomiques (remplace l'ancienne "generate")
	if commande == "generate" {
		commands.GenerateCommand(os.Args[2:])
		return
	}

	// Commande de résumé - NOW USING PHASE 15 (Grammar-Aware Summarization)
	if commande == "resume" {
		if len(os.Args) > 2 {
			filepath := os.Args[2]
			ratio := 0.55 // résumé complet par défaut
			if len(os.Args) > 3 {
				fmt.Sscanf(os.Args[3], "%f", &ratio)
			}
			// Mode hybride: atomique + proba pour un résumé fluide et complet
			commands.ResumeHybridCommand([]string{filepath, fmt.Sprintf("%f", ratio)})
		} else {
			fmt.Println("[ERREUR] Utilisation: ./programme resume <fichier> [ratio=0.55]")
		}
		return
	}

	// Commande de réécriture (humanisation)
	if commande == "rewrite" {
		if len(os.Args) > 2 {
			// This command is deprecated - use 'humanize' instead
			fmt.Println("[ATTENTION] La commande 'rewrite' est dépréciée. Utilisez 'humanize' à la place.")
		} else {
			fmt.Println("[ERREUR] Utilisation: ./programme humanize <fichier> [standard|professionnel|avance]")
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
			commands.GenererAvecAtomesCommand(filepath, compression)
		} else {
			fmt.Println("[ERREUR] Utilisation: ./programme atomic <fichier> [compression=0.3]")
		}
		return
	}
	// === PHASE 15: Grammar-Aware Summarization ===
	if commande == "resume-optimized" {
		commands.ResumeOptimizedCommand(os.Args[2:])
		return
	}
	if commande == "compare-summaries" {
		commands.CompareSummariesCommand(os.Args[2:])
		return
	}
	if commande == "analyze-preprocessing" {
		commands.AnalyzePreprocessingCommand(os.Args[2:])
		return
	}
	if commande == "analyze-vocabulary" {
		commands.AnalyzeVocabularyCommand(os.Args[2:])
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
			commands.GenererComparatifCommand(filepath, compression)
		} else {
			fmt.Println("[ERREUR] Utilisation: ./programme compare <fichier> [compression=0.3]")
		}
		return
	}

	// === PHASE 14: Commande syntaxe avancée avec renforcement ===
	if commande == "syntax" {
		if len(os.Args) > 2 {
			commands.SyntaxCommand(os.Args[2:])
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

	case "web":
		commands.InitWebInterface()

	case "file":
		if len(os.Args) > 2 {
			commands.TraiterFichier(os.Args[2])
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
			commands.TraiterFichierHumanize(filepath, style)
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
			commands.TraiterTexte(texte, "entrée CLI")
		} else {
			fmt.Println("[ERREUR] Utilisation: ./programme text <votre texte>")
		}

	case "interactive", "inter":
		commands.InteractionInteractive()

	case "start":
		commands.StartLLMMode()
		return

	case "fidelity":
		commands.ProcessAntiHallucination(os.Args[1:])

	case "atomic-optimized":
		commands.ProcessOptimizedAtomicCommand(os.Args[1:])

	case "stest", "stress-test":
		fmt.Println("[INFO] Stress tests retirés du binaire principal")

	case "stest-batch":
		fmt.Println("[INFO] Stress tests retirés du binaire principal")

	case "stest-l3":
		fmt.Println("[INFO] Stress tests retirés du binaire principal")

	case "stest-l3-demo":
		fmt.Println("[INFO] Stress tests retirés du binaire principal")

	case "stest-ultra":
		fmt.Println("[INFO] Stress tests retirés du binaire principal")

	case "ultra":
		ProcessUltraLiteMode(os.Args[2:])

	case "energy":
		commands.EnergyBasedImageCommand(os.Args[1:])

	case "ultra_pure":
		if len(os.Args) < 3 {
			fmt.Println("Usage: ./programme ultra_pure <imagePath> [output]")
			fmt.Println("Example: ./programme ultra_pure blurry.jpg deblurred_pure.png")
			fmt.Println("\nPURE atomic deblurring with 5000 iterations (no motion preprocessing)")
			fmt.Println("Total time: ~120-150 seconds for MAXIMUM quality")
			return
		}

		imagePath := os.Args[2]
		finalOutput := "deblurred_ultra_pure.png"
		if len(os.Args) > 3 {
			finalOutput = os.Args[3]
		}

		// Verify image exists
		if _, err := os.Stat(imagePath); err != nil {
			fmt.Printf(" Image not found: %s\n", imagePath)
			return
		}

		fmt.Printf("\n ULTRA PURE ATOMIC DEBLURRING (5000 ITERATIONS, 10px patches)\n")
		fmt.Printf("═══════════════════════════════════════════════════════\n")
		fmt.Printf("Processing: %s\n", imagePath)
		fmt.Printf("Expected time: ~120-150 seconds\n\n")

		// Direct ultra deblurring with 5000 iterations - NO motion preprocessing
		// Use 10px patch size (2x coarser) for more aggressive deblurring
		commands.RunUltraDeblurPipeline4K(imagePath, 4, 4, 5000, finalOutput)

		fmt.Printf("\n ULTRA PURE COMPLETE!\n")
		fmt.Printf(" Saved: %s\n", finalOutput)
		fmt.Printf(" Maximum quality atomic deblurring applied!\n")

	default:
		// Mode hérité - traiter comme phrase simple
		phrase := strings.Join(os.Args[1:], " ")
		commands.TraiterTexte(phrase, "mode classique")
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
	confTotal := 0.0
	for _, c := range allConfiances {
		confTotal += c
	}
	if len(allConfiances) > 0 {
		fmt.Printf("  • Confiance moyenne: %.2f%%\n", (confTotal/float64(len(allConfiances)))*100)
	}

	fmt.Printf("\n[PERFORMANCE]\n")
	fmt.Printf("  • Complexité: O(n)\n")
	fmt.Printf("  • Temps de découpage + analyse: %v\n", temps)
	fmt.Printf("  • Vitesse: %.0f mots/sec\n", float64(stats["total_mots"].(int))/temps.Seconds())

	fmt.Printf("\n✓ Découpage instantané réussi - Compatible avec réseau atomique distribué\n\n")
}
