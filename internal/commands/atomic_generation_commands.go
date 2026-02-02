package commands

import (
	"github.com/nosserb/IA-ATOMIQUE-/database"
	"fmt"
	"os"
	"strings"
	"time"
)

// GenererAvecAtomesCommand génère un résumé via résonance atomique autonome
func GenererAvecAtomesCommand(fichier string, compression float64) {
	contenu, err := os.ReadFile(fichier)
	if err != nil {
		fmt.Printf("[ERREUR] Impossible de lire le fichier: %v\n", err)
		return
	}

	texte := string(contenu)
	fmt.Printf("\n╔════════════════════════════════════════╗\n")
	fmt.Printf("║  GÉNÉRATION ATOMIQUE AUTONOME        ║\n")
	fmt.Printf("╚════════════════════════════════════════╝\n\n")
	fmt.Printf("📄 Fichier: %s\n", fichier)
	fmt.Printf("📊 Compression cible: %.0f%%\n", compression*100)
	fmt.Printf("📝 Taille originale: %d mots\n\n", len(strings.Fields(texte)))

	debut := time.Now()

	// Parser phrases
	phrasesRaw := strings.Split(texte, ".")
	var phrases []database.Phrase

	for i, p := range phrasesRaw {
		p = strings.TrimSpace(p)
		if len(p) < 5 {
			continue
		}

		phrase := database.Phrase{
			Contenu: p,
			Mots:    strings.Fields(p),
			Index:   i,
		}
		phrases = append(phrases, phrase)
	}

	// Traduire si nécessaire
	phrases = database.DetecterEtTraduirePhrases(phrases)

	fmt.Printf("[1/3] ✓ %d phrases parsées et traduites\n", len(phrases))

	// Générer avec résonance atomique

	// Stratégie:
	// - Texte petit: une seule passe
	// - Texte moyen: plusieurs sous-réseaux
	// - Texte gros: sous-réseaux avec résonance inter-réseaux

	var resultat string
	if len(phrases) < 10 {
		// Simple: un seul réseau
		resultat = database.GenererAvecReseauAtomique(phrases, 200, compression)
	} else if len(phrases) < 50 {
		// Moyen: plusieurs chunks
		resultat = database.GenererAvecSousReseaux(phrases, 10, compression)
	} else {
		// Gros: chunk plus large avec résonance globale
		resultat = database.GenererAvecSousReseaux(phrases, 20, compression)
	}

	temps := time.Since(debut)

	// Afficher résultat
	fmt.Printf("╔════════════════════════════════════════╗\n")
	fmt.Printf("║  RÉSUMÉ GÉNÉRÉ PAR RÉSONANCE         ║\n")
	fmt.Printf("╚════════════════════════════════════════╝\n\n")

	// Wrapper le texte pour meilleure lisibilité
	words := strings.Fields(resultat)
	ligne := ""
	for _, word := range words {
		if len(ligne)+len(word)+1 > 80 {
			fmt.Println(ligne)
			ligne = word
		} else {
			if ligne == "" {
				ligne = word
			} else {
				ligne += " " + word
			}
		}
	}
	if ligne != "" {
		fmt.Println(ligne)
	}

	fmt.Printf("\n╔════════════════════════════════════════╗\n")
	fmt.Printf("║  ANALYSE                              ║\n")
	fmt.Printf("╚════════════════════════════════════════╝\n\n")

	tailleResume := len(strings.Fields(resultat))
	compression_reelle := float64(len(strings.Fields(texte))) / float64(tailleResume)

	fmt.Printf("• Mots générés:      %d\n", tailleResume)
	fmt.Printf("• Compression:       %.1fx\n", compression_reelle)
	fmt.Printf("• Phrases source:    %d\n", len(phrases))
	fmt.Printf("• Temps total:       %v\n", temps)
	fmt.Printf("• Architecture:      Résonance atomique distribuée\n")
	fmt.Printf("• Type de chunks:    ")

	if len(phrases) < 10 {
		fmt.Printf("Mono-réseau\n")
	} else if len(phrases) < 50 {
		fmt.Printf("Multi-réseaux (chunks=10)\n")
	} else {
		fmt.Printf("Muli-réseaux + résonance globale (chunks=20)\n")
	}

	fmt.Printf("\n✓ Génération réussie - Réseau autonome\n\n")
}

// GenererComparatifCommand compare vectoriel vs atomique
func GenererComparatifCommand(fichier string, compression float64) {
	contenu, err := os.ReadFile(fichier)
	if err != nil {
		fmt.Printf("[ERREUR] Impossible de lire: %v\n", err)
		return
	}

	texte := string(contenu)

	// Parser phrases
	phrasesRaw := strings.Split(texte, ".")
	var phrases []database.Phrase

	for i, p := range phrasesRaw {
		p = strings.TrimSpace(p)
		if len(p) < 5 {
			continue
		}
		phrase := database.Phrase{
			Contenu: p,
			Mots:    strings.Fields(p),
			Index:   i,
		}
		phrases = append(phrases, phrase)
	}

	phrases = database.DetecterEtTraduirePhrases(phrases)

	fmt.Printf("\n╔════════════════════════════════════════╗\n")
	fmt.Printf("║  COMPARAISON: VECTORIEL vs ATOMIQUE   ║\n")
	fmt.Printf("╚════════════════════════════════════════╝\n\n")

	// Méthode vectorielle
	fmt.Printf("1️⃣  MÉTHODE VECTORIELLE (argmax)\n")
	fmt.Printf("─────────────────────────────────\n")
	tDebut := time.Now()
	params := database.ParamsParDefaut()
	params.TailleCible = int(float64(len(strings.Fields(texte))) * compression)
	state := database.GenerationAdaptive(phrases, params)
	textVectoriel := strings.Join(state.MotsGenerés, " ")
	tVectoriel := time.Since(tDebut)

	coherenceVect := database.FiltrerCoherence(state)
	scoringVect := database.EvaluerQualite(state)

	fmt.Printf("Texte: %s\n\n", wrapText(textVectoriel, 80))
	fmt.Printf("⏱️  Temps: %v\n", tVectoriel)
	fmt.Printf("🎯 Cohérence: %.2f%%\n", coherenceVect*100)
	fmt.Printf("🔀 Diversité: %.2f%%\n", scoringVect.DiversiteVocab*100)
	fmt.Printf("✨ Qualité: %.2f%%\n\n", scoringVect.QualiteGlobale*100)

	// Méthode atomique
	fmt.Printf("2️⃣  MÉTHODE ATOMIQUE (résonance)\n")
	fmt.Printf("───────────────────────────────\n")
	tDebut = time.Now()
	textAtomique := database.GenererAvecReseauAtomique(phrases, 200, compression)
	tAtomique := time.Since(tDebut)

	// Calculer cohérence atomique
	words := strings.Fields(textAtomique)
	totalCoh := 0.0
	for _, w := range words {
		v := database.VectoriserMot(w)
		vTarget := database.MoyennePondereeVecteurs(phrases)
		totalCoh += database.SimilariteCosinus(v, vTarget)
	}
	coherenceAtom := totalCoh / float64(len(words))

	// Diversité
	uniqueWords := make(map[string]bool)
	for _, w := range words {
		uniqueWords[strings.ToLower(w)] = true
	}
	diversiteAtom := float64(len(uniqueWords)) / float64(len(words))

	fmt.Printf("Texte: %s\n\n", wrapText(textAtomique, 80))
	fmt.Printf("⏱️  Temps: %v\n", tAtomique)
	fmt.Printf("🎯 Cohérence: %.2f%%\n", coherenceAtom*100)
	fmt.Printf("🔀 Diversité: %.2f%%\n", diversiteAtom*100)
	fmt.Printf("✨ Qualité: %.2f%%\n\n", (coherenceAtom+diversiteAtom)/2*100)

	// Comparaison
	fmt.Printf("📊 COMPARAISON\n")
	fmt.Printf("──────────────\n")
	speedup := float64(tVectoriel) / float64(tAtomique)
	fmt.Printf("Speedup: %.1fx (Atomique %s Vectoriel)\n",
		speedup,
		map[bool]string{true: "PLUS RAPIDE", false: "PLUS LENT"}[speedup > 1])

	fmt.Printf("Vectoriel: %d mots | Atomique: %d mots\n",
		len(strings.Fields(textVectoriel)),
		len(strings.Fields(textAtomique)))

	fmt.Printf("\n✓ Comparaison complète\n\n")
}

// wrapText enveloppe le texte à une certaine largeur
func wrapText(texte string, width int) string {
	words := strings.Fields(texte)
	lines := []string{}
	currentLine := ""

	for _, word := range words {
		if len(currentLine)+len(word)+1 > width {
			if currentLine != "" {
				lines = append(lines, currentLine)
			}
			currentLine = word
		} else {
			if currentLine == "" {
				currentLine = word
			} else {
				currentLine += " " + word
			}
		}
	}

	if currentLine != "" {
		lines = append(lines, currentLine)
	}

	return strings.Join(lines, "\n")
}
