package database

import (
	"fmt"
	"sort"
	"strings"
)

// AnalysePhrases analyse un texte phrase par phrase
type AnalysePhrases struct {
	Phrases []AnalysePhrase
	Resume  map[int]ResumeCategorie // Résumé par catégorie
}

type AnalysePhrase struct {
	Texte          string
	CategorieID    int
	Categorie      string
	Confiance      float64
	MotsClés       []string
	Score          float64
	VerbePrincipal string // Le verbe détecté (séparé)
}

type ResumeCategorie struct {
	Categorie     string
	NbPhrases     int
	Confiance     float64
	MotsCles      []string
	Phrases       []string
	EnergieTotale float64
}

// AnalyserParPhrases traite le texte phrase par phrase
func AnalyserParPhrases(texte string) AnalysePhrases {
	phrases := SplitterParPhrases(texte)

	analyse := AnalysePhrases{
		Phrases: []AnalysePhrase{},
		Resume:  make(map[int]ResumeCategorie),
	}

	// Analyser chaque phrase
	for _, phrase := range phrases {
		if len(strings.TrimSpace(phrase)) < 10 {
			continue
		}

		// Tokeniser
		tokens := TokeniserTexte(phrase)

		// Détecter le verbe principal (à part)
		verbePrincipal := ""
		for _, token := range tokens {
			if word, ok := Words[token]; ok && word.Categorie == 6 {
				verbePrincipal = token
				break
			}
		}

		// Trouver la catégorie (en ignorant les verbes)
		catActivation := make(map[int]int)
		for _, token := range tokens {
			// Chercher le mot exact
			if word, ok := Words[token]; ok && word.Categorie > 0 && word.Categorie != 6 {
				// Ignorer les verbes (cat 6) pour la classification principale
				catActivation[word.Categorie]++
			} else if len(token) > 2 && strings.HasSuffix(token, "s") {
				// Essayer singulier (retirer le 's' final)
				singulier := token[:len(token)-1]
				if word, ok := Words[singulier]; ok && word.Categorie > 0 && word.Categorie != 6 {
					catActivation[word.Categorie]++
				}
			}
		}

		// Catégorie principale
		var catMain int
		var scoreMain int
		for cat, count := range catActivation {
			if count > scoreMain {
				scoreMain = count
				catMain = cat
			}
		}

		// Si aucune catégorie trouvée, chercher un verbe
		if catMain == 0 {
			for _, token := range tokens {
				if word, ok := Words[token]; ok && word.Categorie == 6 {
					// C'est un verbe, on classe selon le contexte
					catMain = 6
					scoreMain = 1
					break
				}
			}
		}

		// Extraire mots clés pour cette phrase
		motsCles := extraireMotsPhrase(tokens, catMain)

		// Calculer confiance
		confiance := 0.0
		if len(tokens) > 0 {
			confiance = float64(scoreMain) / float64(len(tokens))
		}

		// Score d'importance
		score := ScorerPhrase(phrase)

		analysePhrase := AnalysePhrase{
			Texte:          phrase,
			CategorieID:    catMain,
			Categorie:      NumeroVersCategorie(catMain),
			Confiance:      confiance,
			MotsClés:       motsCles,
			Score:          score,
			VerbePrincipal: verbePrincipal, // Ajouter le verbe détecté
		}

		analyse.Phrases = append(analyse.Phrases, analysePhrase)

		// Agréger dans le résumé
		if catMain > 0 {
			resume := analyse.Resume[catMain]
			resume.Categorie = NumeroVersCategorie(catMain)
			resume.NbPhrases++
			resume.Confiance = (resume.Confiance + confiance) / 2
			resume.EnergieTotale += score

			// Ajouter la phrase si importante
			if score > 0.5 {
				resume.Phrases = append(resume.Phrases, phrase)
			}

			// Ajouter mots clés uniques
			for _, mot := range motsCles {
				found := false
				for _, m := range resume.MotsCles {
					if m == mot {
						found = true
						break
					}
				}
				if !found {
					resume.MotsCles = append(resume.MotsCles, mot)
				}
			}

			analyse.Resume[catMain] = resume
		}
	}

	return analyse
}

// extraireMotsPhrase extrait les mots clés d'une phrase
func extraireMotsPhrase(tokens []string, catID int) []string {
	var motsCles []string

	for _, token := range tokens {
		if word, ok := Words[token]; ok && word.Categorie == catID {
			motsCles = append(motsCles, token)
		}
	}

	// Limiter à 5
	if len(motsCles) > 5 {
		motsCles = motsCles[:5]
	}

	return motsCles
}

// AfficherAnalyseDetaillee affiche l'analyse complète
func AfficherAnalyseDetaillee(analyse AnalysePhrases) string {
	output := ""

	// Résumé global
	totalPhrases := len(analyse.Phrases)
	output += fmt.Sprintf("\n[ANALYSE PAR PHRASE]\nTotal phrases: %d\n\n", totalPhrases)

	// Afficher les phrases par catégorie
	for catID := 1; catID <= 6; catID++ {
		if resume, ok := analyse.Resume[catID]; ok && resume.NbPhrases > 0 {
			output += fmt.Sprintf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
			output += fmt.Sprintf("[%s] - %d phrases (énergie: %.2f)\n",
				resume.Categorie, resume.NbPhrases, resume.EnergieTotale)
			output += fmt.Sprintf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")

			// Mots clés
			if len(resume.MotsCles) > 0 {
				motsClesListe := resume.MotsCles
				if len(motsClesListe) > 5 {
					motsClesListe = motsClesListe[:5]
				}
				output += fmt.Sprintf("Mots clés: %s\n", strings.Join(motsClesListe, " | "))
			}

			// Phrases avec verbes détectés
			if len(resume.Phrases) > 0 {
				output += fmt.Sprintf("\nPhrases clés:\n")
				for i, phrase := range resume.Phrases {
					if i >= 3 { // Max 3 phrases par catégorie
						break
					}
					// Chercher le verbe pour cette phrase
					verbe := ""
					for _, p := range analyse.Phrases {
						if p.Texte == phrase && p.VerbePrincipal != "" {
							verbe = " [verbe: " + p.VerbePrincipal + "]"
							break
						}
					}
					output += fmt.Sprintf("  • %s%s\n", phrase, verbe)
				}
			}
			output += "\n"
		}
	}

	return output
}

// StatistiquesAnalyse retourne des statistiques
func StatistiquesAnalyse(analyse AnalysePhrases) map[string]interface{} {
	stats := make(map[string]interface{})

	// Comptage par catégorie
	comptage := make(map[string]int)
	energieTotale := 0.0
	confianceMoyenne := 0.0

	for _, phrase := range analyse.Phrases {
		comptage[phrase.Categorie]++
		energieTotale += phrase.Score
		confianceMoyenne += phrase.Confiance
	}

	if len(analyse.Phrases) > 0 {
		confianceMoyenne /= float64(len(analyse.Phrases))
	}

	stats["comptage"] = comptage
	stats["energie_totale"] = energieTotale
	stats["confiance_moyenne"] = confianceMoyenne
	stats["nb_phrases"] = len(analyse.Phrases)

	return stats
}

// CategoriesDominantes retourne les catégories les plus présentes
func CategoriesDominantes(analyse AnalysePhrases) []map[string]interface{} {
	comptage := make(map[string]int)

	for _, phrase := range analyse.Phrases {
		if phrase.CategorieID > 0 {
			comptage[phrase.Categorie]++
		}
	}

	type catResult struct {
		Categorie   string
		Nombre      int
		Pourcentage float64
	}

	var resultats []catResult
	total := len(analyse.Phrases)

	for cat, nb := range comptage {
		resultats = append(resultats, catResult{
			Categorie:   cat,
			Nombre:      nb,
			Pourcentage: float64(nb) * 100 / float64(total),
		})
	}

	// Trier par nombre décroissant
	sort.Slice(resultats, func(i, j int) bool {
		return resultats[i].Nombre > resultats[j].Nombre
	})

	// Convertir en interface{}
	var output []map[string]interface{}
	for _, r := range resultats {
		output = append(output, map[string]interface{}{
			"categorie":   r.Categorie,
			"nombre":      r.Nombre,
			"pourcentage": r.Pourcentage,
		})
	}

	return output
}
