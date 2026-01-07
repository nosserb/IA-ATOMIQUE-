package main

import (
	"IA-ATOMIQUE/database"
	"fmt"
	"os"
	"strings"
	"time"
)

// GenererResumeCommand génère un résumé par génération vectorielle
func GenererResumeCommand(fichier string, ratio float64) {
	// Lire le fichier
	contenu, err := os.ReadFile(fichier)
	if err != nil {
		fmt.Printf("[ERREUR] Impossible de lire le fichier: %v\n", err)
		return
	}

	texte := string(contenu)
	fmt.Printf("[GÉNÉRATION DE RÉSUMÉ]\n")
	fmt.Printf("Fichier: %s\n", fichier)
	fmt.Printf("Ratio de compression: %.1f%%\n", ratio*100)
	fmt.Printf("Taille originale: %d mots\n\n", len(strings.Fields(texte)))

	debut := time.Now()

	// Découper en phrases (même logique que TraiterFichier)
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

	fmt.Printf("[1/3] Découpage: %d phrases identifiées\n", len(phrases))

	// Traduction si nécessaire
	phrases = database.DetecterEtTraduirePhrases(phrases)
	fmt.Printf("[2/3] Traduction: langues diversifiées converties en FR\n")

	// Vectorisation
	fmt.Printf("[3/3] Vectorisation et génération en cours...\n")

	// Paramètres de génération
	params := database.ParamsParDefaut()
	params.TailleCible = int(float64(len(strings.Fields(texte))) * ratio)

	// Générer résumé
	state := database.GenerationAdaptive(phrases, params)

	// Évaluer qualité
	scoring := database.EvaluerQualite(state)

	// Temps total
	temps := time.Since(debut)

	// Afficher résultats
	resume := strings.Join(state.MotsGenerés, " ")
	if !strings.HasSuffix(resume, ".") {
		resume += "."
	}

	// Connecteurs sémantiques
	resume = database.AjouterConnecteurs(resume, state)

	fmt.Printf("\n╔════════════════════════════════════════╗\n")
	fmt.Printf("║  RÉSUMÉ GÉNÉRÉ (VECTORIEL)              ║\n")
	fmt.Printf("╚════════════════════════════════════════╝\n\n")
	fmt.Printf("%s\n\n", resume)

	fmt.Printf("╔════════════════════════════════════════╗\n")
	fmt.Printf("║  MÉTRIQUES DE QUALITÉ                  ║\n")
	fmt.Printf("╚════════════════════════════════════════╝\n\n")
	fmt.Printf("• Cohérence finale:    %.2f%%\n", scoring.CoherenceFinale*100)
	fmt.Printf("• Diversité lexicale:  %.2f%%\n", scoring.DiversiteVocab*100)
	fmt.Printf("• Qualité globale:     %.2f%%\n", scoring.QualiteGlobale*100)
	fmt.Printf("• Énergie totale:      %.2f\n", scoring.EnergieTotale)

	fmt.Printf("\n╔════════════════════════════════════════╗\n")
	fmt.Printf("║  STATISTIQUES                         ║\n")
	fmt.Printf("╚════════════════════════════════════════╝\n\n")
	fmt.Printf("• Mots générés:        %d (cible: %d)\n", len(state.MotsGenerés), params.TailleCible)
	fmt.Printf("• Compression:         %.1fx\n", float64(len(strings.Fields(texte)))/float64(len(state.MotsGenerés)))
	fmt.Printf("• Phrases d'entrée:    %d\n", len(phrases))
	fmt.Printf("• Temps total:         %v\n", temps)

	fmt.Printf("\n✓ Génération réussie - Résumé créé par résonance vectorielle\n\n")
}

// GenererResumeDecoupageCommand génère un résumé avec découpage et cohérence vectorielle
func GenererResumeDecoupageCommand(fichier string, ratio float64) {
	// Lire le fichier
	contenu, err := os.ReadFile(fichier)
	if err != nil {
		fmt.Printf("[ERREUR] Impossible de lire le fichier: %v\n", err)
		return
	}

	texte := string(contenu)
	fmt.Printf("[RÉSUMÉ AVEC DÉCOUPAGE & COHÉRENCE]\n")
	fmt.Printf("Fichier: %s\n", fichier)
	fmt.Printf("Ratio de compression: %.1f%%\n", ratio*100)
	motsTotaux := len(strings.Fields(texte))
	fmt.Printf("Taille originale: %d mots\n\n", motsTotaux)

	debut := time.Now()

	// Découper en phrases
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

	fmt.Printf("[1/5] Découpage: %d phrases identifiées\n", len(phrases))

	// Traduction
	phrases = database.DetecterEtTraduirePhrases(phrases)
	fmt.Printf("[2/5] Traduction: langues diversifiées converties en FR\n")

	// Découper en blocs et vectoriser
	fmt.Printf("[3/5] Découpage en blocs vectoriels...\n")
	resumeur := database.ResumerAvecDecoupage(phrases, ratio)
	fmt.Printf("      %d blocs créés, vecteur global calculé\n", len(resumeur.Blocs))

	// Sélectionner les meilleurs blocs
	// === PHASE 13: Utiliser la sélection avec fenêtrage glissant ===
	ratioSelection := ratio
	if ratioSelection < 0.25 {
		ratioSelection = 0.25 // Minimum 25% de blocs pour diversité
	}
	fmt.Printf("[4/5] Sélection par cohérence, fenêtrage glissant & normalisation...\n")

	// Utiliser le fenêtrage glissant (window size 15 par défaut)
	blocsSelectionnes := resumeur.SelectionnerBlocsAvecFenetrageGlissant(ratioSelection, 15)

	fmt.Printf("      %d blocs sélectionnés (fenêtrage glissant, cohérence: %.2f%%)\n",
		len(blocsSelectionnes),
		calculerCoherenceMoyenne(blocsSelectionnes)*100)

	// Générer résumé à partir des blocs sélectionnés
	fmt.Printf("[5/5] Génération du résumé...\n")
	// === PHASE 13++: Augmenter cible pour ~1000+ mots ===
	// 50 blocs × 25 mots = 1250 mots cible
	motsParBloc := 25
	tailleResume := len(blocsSelectionnes) * motsParBloc
	if tailleResume < 300 {
		tailleResume = 300
	}
	if tailleResume > int(float64(motsTotaux)*ratio) { // Ne pas dépasser ratio théorique
		tailleResume = int(float64(motsTotaux) * ratio)
	}
	resume := resumeur.GenerByVectorsWithGrammar(blocsSelectionnes, tailleResume)

	// === PHASE 13: Post-traitement pour fluidité ===
	resume = database.PostTraiterResume(resume)

	temps := time.Since(debut)

	// Afficher résultats
	if !strings.HasSuffix(resume, ".") {
		resume += "."
	}

	fmt.Printf("\n╔════════════════════════════════════════╗\n")
	fmt.Printf("║  RÉSUMÉ (DÉCOUPAGE & COHÉRENCE)        ║\n")
	fmt.Printf("╚════════════════════════════════════════╝\n\n")

	// Wrapper le texte pour affichage
	lignes := strings.Split(wrapText(resume, 80), "\n")
	for _, ligne := range lignes {
		fmt.Printf("%s\n", ligne)
	}

	fmt.Printf("\n╔════════════════════════════════════════╗\n")
	fmt.Printf("║  ANALYSE VECTORIELLE                  ║\n")
	fmt.Printf("╚════════════════════════════════════════╝\n\n")

	cohMoyenne := calculerCoherenceMoyenne(blocsSelectionnes)
	tfidfMoyenne := calculerTFIDFMoyenne(blocsSelectionnes)
	energieMoyenne := calculerEnergieMoyenne(blocsSelectionnes)

	fmt.Printf("• Cohérence moyenne:   %.2f%%\n", cohMoyenne*100)
	fmt.Printf("• TF-IDF moyen:        %.2f%%\n", tfidfMoyenne*100)
	fmt.Printf("• Énergie moyenne:     %.4f\n", energieMoyenne)
	fmt.Printf("• Blocs sélectionnés:  %d / %d\n", len(blocsSelectionnes), len(resumeur.Blocs))

	fmt.Printf("\n╔════════════════════════════════════════╗\n")
	fmt.Printf("║  STATISTIQUES                         ║\n")
	fmt.Printf("╚════════════════════════════════════════╝\n\n")
	fmt.Printf("• Mots générés:        %d (cible: %d)\n", len(strings.Fields(resume)), tailleResume)
	fmt.Printf("• Compression réelle:  %.1fx\n", float64(motsTotaux)/float64(len(strings.Fields(resume))))
	fmt.Printf("• Phrases d'entrée:    %d\n", len(phrases))
	fmt.Printf("• Temps total:         %v\n", temps)

	fmt.Printf("\n✓ Résumé généré avec découpage vectoriel et cohérence globale\n\n")
}

// Fonctions helper
func calculerCoherenceMoyenne(blocs []database.BlocVectoriel) float64 {
	if len(blocs) == 0 {
		return 0
	}
	somme := 0.0
	for _, b := range blocs {
		somme += b.Coherence
	}
	return somme / float64(len(blocs))
}

func calculerTFIDFMoyenne(blocs []database.BlocVectoriel) float64 {
	if len(blocs) == 0 {
		return 0
	}
	somme := 0.0
	for _, b := range blocs {
		somme += b.TFIDFScore
	}
	return somme / float64(len(blocs))
}

func calculerEnergieMoyenne(blocs []database.BlocVectoriel) float64 {
	if len(blocs) == 0 {
		return 0
	}
	somme := 0.0
	for _, b := range blocs {
		somme += b.Energie
	}
	return somme / float64(len(blocs))
}

// ReecrireCommand réécrit un texte dans différents styles
func ReecrireCommand(texte string, style string) {
	fmt.Printf("[RÉÉCRITURE DE TEXTE]\n")
	fmt.Printf("Style: %s\n", style)
	fmt.Printf("Texte original: %s\n\n", texte)

	debut := time.Now()

	// Paramétriser selon le style
	params := database.ParamsParDefaut()

	switch style {
	case "professionnel":
		params.Eta = 0.2 // Moins créatif
		params.Gamma = 0.1
		fmt.Printf("Mode professionnel: Priorité à la cohérence et la précision\n\n")
	case "avance":
		params.Eta = 0.5 // Plus créatif
		params.Gamma = 0.2
		fmt.Printf("Mode avancé: Priorité à la créativité et la variété\n\n")
	default: // standard
		params.Eta = 0.3
		params.Gamma = 0.15
		fmt.Printf("Mode standard: Équilibre cohérence/créativité\n\n")
	}

	// Créer phrase virtuelle
	mots := strings.Fields(texte)
	phrase := database.Phrase{
		Contenu: texte,
		Mots:    mots,
		Energie: 0.8,
	}
	phrase.MotsClés = database.ExtraireMotsClés(mots)

	// Vectorisation
	vec := database.VectoriserPhrase(&phrase)
	fmt.Printf("Vectorisation: %d dimensions sémantiques\n", len(vec.Dimensions))
	fmt.Printf("Énergie du vecteur: %.3f\n\n", vec.Energie)

	// Générer version réécrite
	texteReecrit := database.ReecrirePhrase(texte, style)

	// Évaluer
	temps := time.Since(debut)

	fmt.Printf("╔════════════════════════════════════════╗\n")
	fmt.Printf("║  TEXTE RÉÉCRIT                         ║\n")
	fmt.Printf("╚════════════════════════════════════════╝\n\n")
	fmt.Printf("%s\n\n", texteReecrit)

	fmt.Printf("╔════════════════════════════════════════╗\n")
	fmt.Printf("║  ANALYSE                              ║\n")
	fmt.Printf("╚════════════════════════════════════════╝\n\n")

	originalWords := len(mots)
	rewriteWords := len(strings.Fields(texteReecrit))
	uniqueOriginal := make(map[string]bool)
	for _, w := range mots {
		uniqueOriginal[strings.ToLower(w)] = true
	}
	uniqueRewrite := make(map[string]bool)
	for _, w := range strings.Fields(texteReecrit) {
		uniqueRewrite[strings.ToLower(w)] = true
	}

	fmt.Printf("• Longueur originale:  %d mots\n", originalWords)
	fmt.Printf("• Longueur réécrite:   %d mots\n", rewriteWords)
	fmt.Printf("• Variation lexicale:  %.1f%% unique\n", float64(len(uniqueRewrite))/float64(rewriteWords)*100)
	fmt.Printf("• Temps de réécriture: %v\n", temps)
	fmt.Printf("• Paramètres:          Eta=%.2f, Gamma=%.2f\n", params.Eta, params.Gamma)

	fmt.Printf("\n✓ Réécriture réussie - Texte transformé en mode '%s'\n\n", style)
}

// TestGenerationCombinee teste la génération sur plusieurs styles
func TestGenerationCombinee() {
	texteTest := "L'intelligence artificielle révolutionne le monde moderne en fournissant des solutions innovantes pour les entreprises."

	fmt.Printf("\n╔════════════════════════════════════════╗\n")
	fmt.Printf("║  TEST GÉNÉRATION MULTI-STYLE          ║\n")
	fmt.Printf("╚════════════════════════════════════════╝\n\n")

	styles := []string{"standard", "professionnel", "avance"}

	for _, style := range styles {
		fmt.Printf("\n--- Style: %s ---\n", style)
		texteReecrit := database.ReecrirePhrase(texteTest, style)
		fmt.Printf("%s\n", texteReecrit)
	}

	fmt.Printf("\n✓ Test multi-style complété\n\n")
}
