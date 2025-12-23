package database

import (
	"fmt"
	"math/rand"
	"strings"
)

// GenererResumeSynthetique crée un résumé unique en propres mots
func GenererResumeSynthetique(analyse AnalysePhrases) string {
	if len(analyse.Phrases) == 0 {
		return "Aucun contenu à résumer."
	}

	// Trouver la catégorie dominante
	var catMain int
	var maxPhrases int
	for catID, resume := range analyse.Resume {
		if resume.NbPhrases > maxPhrases {
			maxPhrases = resume.NbPhrases
			catMain = catID
		}
	}

	if catMain == 0 {
		return "Impossible de générer un résumé."
	}

	// Extraire les meilleures phrases du texte
	phrasesPrincipales := extrairePhrasesImportantes(analyse, catMain)
	motsCles := extraireMotsClesUniques(analyse)

	// Générer le résumé dynamique
	resume := formulerResume(phrasesPrincipales, motsCles, catMain, len(analyse.Phrases))
	return resume
}

// extrairePhrasesImportantes récupère les phrases les plus pertinentes
func extrairePhrasesImportantes(analyse AnalysePhrases, catID int) []string {
	type phraseScore struct {
		texte string
		score float64
	}
	var phrases []phraseScore

	for _, p := range analyse.Phrases {
		if p.CategorieID == catID && len(p.MotsClés) > 0 {
			phrases = append(phrases, phraseScore{
				texte: p.Texte,
				score: p.Score * p.Confiance,
			})
		}
	}

	// Trier par score
	for i := 0; i < len(phrases)-1; i++ {
		for j := 0; j < len(phrases)-1-i; j++ {
			if phrases[j].score < phrases[j+1].score {
				phrases[j], phrases[j+1] = phrases[j+1], phrases[j]
			}
		}
	}

	// Retourner top 3
	var top []string
	for i := 0; i < min(3, len(phrases)); i++ {
		top = append(top, phrases[i].texte)
	}
	return top
}

// extraireMotsClesUniques récupère les mots clés les plus importants
func extraireMotsClesUniques(analyse AnalysePhrases) []string {
	freqMots := make(map[string]int)

	for _, phrase := range analyse.Phrases {
		for _, mot := range phrase.MotsClés {
			freqMots[mot]++
		}
	}

	type motFreq struct {
		mot  string
		freq int
	}
	var mots []motFreq
	for mot, freq := range freqMots {
		mots = append(mots, motFreq{mot, freq})
	}

	// Trier par fréquence
	for i := 0; i < len(mots)-1; i++ {
		for j := 0; j < len(mots)-1-i; j++ {
			if mots[j].freq < mots[j+1].freq {
				mots[j], mots[j+1] = mots[j+1], mots[j]
			}
		}
	}

	// Retourner top 5
	var top5 []string
	for i := 0; i < min(5, len(mots)); i++ {
		top5 = append(top5, mots[i].mot)
	}
	return top5
}

// formulerResume construit un résumé narratif unique basé sur les données réelles
func formulerResume(phrases []string, motsCles []string, catID int, totalPhrases int) string {
	var texte string

	// Varier l'introduction selon la catégorie et le contenu
	introductions := map[int][]string{
		1: {
			"En analysant ce contenu, je détecte une forte présence d'éléments technologiques.",
			"Le texte révèle une préoccupation marquée pour les questions informatiques.",
			"L'analyse montre une orientation clairement technologique.",
		},
		2: {
			"Le contenu traite principalement d'enjeux historiques et patrimoniaux.",
			"Cette analyse révèle une dimension historique significative.",
			"Le texte se concentre sur des aspects historiques et culturels.",
		},
		3: {
			"Ce texte aborde des questions commerciales et entrepreneuriales.",
			"L'analyse révèle une préoccupation pour les enjeux d'affaires.",
			"Le contenu se focalise sur des aspects économiques et commerciaux.",
		},
		4: {
			"Le texte porte principalement sur l'alimentation et la gastronomie.",
			"L'analyse détecte une forte présence de sujets culinaires.",
			"Ce contenu traite de questions liées à la nourriture et à la cuisine.",
		},
		5: {
			"Le contenu aborde principalement des questions de santé et de bien-être.",
			"L'analyse révèle une préoccupation centrale pour la santé.",
			"Ce texte se concentre sur des enjeux médicaux et sanitaires.",
		},
	}

	// Sélectionner une introduction aléatoire
	intros := introductions[catID]
	intro := intros[rand.Intn(len(intros))]
	texte += intro + " "

	// Ajouter les mots clés détectés
	if len(motsCles) > 0 {
		texte += fmt.Sprintf("Les éléments marquants sont: %s. ",
			strings.Join(motsCles[:min(3, len(motsCles))], ", "))
	}

	// Extraire et reformuler une phrase clé
	if len(phrases) > 0 {
		phraseClé := phrases[0]
		// Nettoyer la phrase
		phraseClé = strings.TrimSpace(phraseClé)
		phraseClé = strings.ToLower(phraseClé)

		// Varier la façon de présenter
		variations := []string{
			fmt.Sprintf("En particulier, %s", phraseClé),
			fmt.Sprintf("Notamment, %s", phraseClé),
			fmt.Sprintf("Plus concrètement, %s", phraseClé),
			fmt.Sprintf("À titre d'illustration, %s", phraseClé),
		}
		texte += variations[rand.Intn(len(variations))] + ". "
	}

	// Ajouter des observations basées sur les patterns
	if len(phrases) > 1 {
		observations := map[int][]string{
			1: {
				"Cela démontre l'importance croissante de la technologie dans ce domaine.",
				"Ces éléments reflètent l'impact de l'innovation numérique.",
				"Cela illustre comment la technologie transforme cet espace.",
			},
			2: {
				"Cela illustre l'importance du patrimoine historique.",
				"Ces éléments montrent la richesse de notre héritage culturel.",
				"Cela démontre la continuité historique de ces phénomènes.",
			},
			3: {
				"Cela souligne l'importance des stratégies commerciales.",
				"Ces éléments reflètent la dynamique des marchés.",
				"Cela illustre les enjeux entrepreneuriaux contemporains.",
			},
			4: {
				"Cela montre l'importance de la cuisine et de l'alimentation.",
				"Ces éléments révèlent la richesse gastronomique.",
				"Cela illustre la diversité des pratiques culinaires.",
			},
			5: {
				"Cela souligne l'importance du bien-être personnel.",
				"Ces éléments reflètent la priorité accordée à la santé.",
				"Cela démontre l'impact de la santé sur nos vies.",
			},
		}
		obs := observations[catID]
		texte += obs[rand.Intn(len(obs))] + " "
	}

	// Conclusion statistique
	texte += fmt.Sprintf("Au total, %d phrases ont été analysées, avec une concentration importante sur ces thématiques.", totalPhrases)

	return texte
}

// GenererReponseAvancee crée une réponse avec context
func GenererReponseAvancee(catID int, motsClés []string, texteOriginal string) string {
	base := GenererReponse(catID, nil)

	// Ajouter contexte basé sur les mots clés
	if len(motsClés) > 0 {
		return fmt.Sprintf("%s\nConcepts détectés: %s", base, strings.Join(motsClés[:min(3, len(motsClés))], ", "))
	}

	return base
}

// ParlerTexte génère du texte parlé basé sur le contenu
func ParlerTexte(resume string, categorie string) string {
	// Enlever les sauts de lignes
	texte := strings.ReplaceAll(resume, "\n", " ")
	texte = strings.ReplaceAll(texte, "  ", " ")

	// Ajouter la catégorie
	output := fmt.Sprintf("Dans la catégorie %s: %s", categorie, texte)

	return output
}

// helper
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
