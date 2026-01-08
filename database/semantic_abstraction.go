package database

import (
	"fmt"
	"regexp"
	"strings"
)

// ============= PHASE X: SEMANTIC ABSTRACTION LAYER =============

// ConceptAbstrait représente un concept extrait du texte
type ConceptAbstrait struct {
	Mot       string
	Type      string // "action", "cause", "theme", "critique", "objectif"
	Score     float64
	Exemples  []string // phrases qui illustrent ce concept
	Frequency int
}

// AnalyseSemantique analyse en profondeur pour extraire les concepts
type AnalyseSemantique struct {
	Concepts      []ConceptAbstrait
	Themes        []string       // thèmes principaux
	Causes        []string       // causes identifiées
	Objectifs     []string       // objectifs/buts
	Critiques     []string       // critiques sociales
	Abstractions  map[string]int // mots abstraits trouvés
	ScoreAbstrait float64        // 0-100% d'abstraction
}

// DictionnairesAbstraction contient les mappings pour l'abstraction
var (
	// Concepts sociaux
	ConceptsSociaux = map[string][]string{
		"misère":       {"pauvreté", "pauvre", "indigence", "dénuement", "besoin", "manque", "faim", "froid"},
		"oppression":   {"esclavage", "domination", "tyrannie", "servitude", "assujettissement", "subjugation"},
		"injustice":    {"inégalité", "discrimination", "iniquité", "préjugé", "abus", "mauvais traitement"},
		"exploitation": {"trafiquer", "exploiter", "profiter", "abuser", "extorquer", "rançonner"},
		"sacrifice":    {"renoncer", "abandon", "abnégation", "dévouement", "offrande", "renoncement"},
		"violence":     {"maltraitance", "cruauté", "brutalité", "sévérité", "agression", "coup"},
		"rôle":         {"mère", "enfant", "père", "femme", "homme", "riche", "pauvre", "maître", "esclave"},
	}

	// Verbes d'abstraction (indiquent qu'on monte en abstraction)
	VerbsAbstraction = []string{
		"illustrer", "incarner", "représenter", "symboliser", "dénoncer",
		"révéler", "montrer", "exposer", "critiquer", "condemner",
		"manifester", "témoigner", "exprimer", "signifier", "refléter",
		"exemplifier", "démontrer", "prouver", "etablir",
	}

	// Patterns de citations à éviter
	QuotationPatterns = []*regexp.Regexp{
		regexp.MustCompile(`«[^»]+»`),
		regexp.MustCompile(`"[^"]+"`),
		regexp.MustCompile(`'[^']+'`),
		regexp.MustCompile(`—\s+[^.!?]+[.!?]`), // tiret + citation
	}
)

// AnalyserSemantiquement extrait concepts et abstraction d'un texte
func AnalyserSemantiquement(texte string, phrases []string) AnalyseSemantique {
	analyse := AnalyseSemantique{
		Concepts:      make([]ConceptAbstrait, 0),
		Themes:        make([]string, 0),
		Causes:        make([]string, 0),
		Objectifs:     make([]string, 0),
		Critiques:     make([]string, 0),
		Abstractions:  make(map[string]int),
		ScoreAbstrait: 0.0,
	}

	texteLower := strings.ToLower(texte)

	// 1. Extraire les concepts sociaux
	for concept, synonymes := range ConceptsSociaux {
		score := 0.0
		exemples := make([]string, 0)
		frequency := 0

		// Chercher le concept principal
		if strings.Contains(texteLower, concept) {
			score += 2.0
			frequency++
		}

		// Chercher les synonymes
		for _, syn := range synonymes {
			count := strings.Count(texteLower, syn)
			if count > 0 {
				score += float64(count) * 1.5
				frequency += count
				// Trouver la phrase qui contient le synonyme
				for _, phrase := range phrases {
					if strings.Contains(strings.ToLower(phrase), syn) && len(exemples) < 2 {
						exemples = append(exemples, phrase)
					}
				}
			}
		}

		if score > 0 {
			analyse.Concepts = append(analyse.Concepts, ConceptAbstrait{
				Mot:       concept,
				Type:      "theme",
				Score:     score,
				Exemples:  exemples,
				Frequency: frequency,
			})
			analyse.Abstractions[concept] = frequency
		}
	}

	// 2. Détecter les verbes d'abstraction
	verbsCount := 0
	for _, verb := range VerbsAbstraction {
		if strings.Contains(texteLower, verb) {
			count := strings.Count(texteLower, verb)
			analyse.Abstractions[verb] = count
			verbsCount += count
		}
	}

	// 3. Détecter les causes (phrases avec "car", "parce que", "à cause de")
	causeCues := []string{"car", "parce que", "à cause de", "du fait que", "en raison de", "causé par"}
	for _, phrase := range phrases {
		phraseLower := strings.ToLower(phrase)
		for _, cue := range causeCues {
			if strings.Contains(phraseLower, cue) {
				analyse.Causes = append(analyse.Causes, phrase)
				break
			}
		}
	}

	// 4. Détecter les objectifs (phrases avec "pour", "afin de", "visant à")
	objectiveCues := []string{"pour", "afin de", "visant à", "en vue de", "dans le but de", "chercher"}
	for _, phrase := range phrases {
		phraseLower := strings.ToLower(phrase)
		for _, cue := range objectiveCues {
			if strings.Contains(phraseLower, cue) {
				analyse.Objectifs = append(analyse.Objectifs, phrase)
				break
			}
		}
	}

	// 5. Détecter les critiques (mots péjoratifs, négatifs)
	critiqueWords := []string{"injuste", "immoral", "cruel", "inhumain", "scandaleux", "honteux", "révoltant", "abominable", "dégradant", "avilissant"}
	for _, phrase := range phrases {
		phraseLower := strings.ToLower(phrase)
		for _, word := range critiqueWords {
			if strings.Contains(phraseLower, word) {
				analyse.Critiques = append(analyse.Critiques, phrase)
				break
			}
		}
	}

	// 6. Calculer le score d'abstraction (pourcentage de phrases abstraites)
	phrasesAbstraites := 0
	for _, phrase := range phrases {
		if EstAbstraite(phrase) {
			phrasesAbstraites++
		}
	}

	if len(phrases) > 0 {
		analyse.ScoreAbstrait = float64(phrasesAbstraites) / float64(len(phrases)) * 100.0
	}

	return analyse
}

// EstAbstraite vérifie si une phrase est abstraite (contient des concepts généraux)
func EstAbstraite(phrase string) bool {
	phraseLower := strings.ToLower(phrase)

	// Vérifier la présence de verbes d'abstraction
	for _, verb := range VerbsAbstraction {
		if strings.Contains(phraseLower, verb) {
			return true
		}
	}

	// Vérifier la présence de concepts abstraits
	for concept := range ConceptsSociaux {
		if strings.Contains(phraseLower, concept) {
			return true
		}
	}

	// Vérifier la présence de mots généralisants
	abstractWords := []string{
		"représente", "symbolise", "incarne", "manifeste",
		"révèle", "montre", "illustre", "démontre",
		"révèlement", "témoignage", "expression",
		"universel", "général", "fondamental", "essentiel",
	}

	for _, word := range abstractWords {
		if strings.Contains(phraseLower, word) {
			return true
		}
	}

	return false
}

// FiltrerCitations élimine les citations brutes du texte
func FiltrerCitations(texte string) string {
	resultat := texte

	// Supprimer chaque pattern de citation
	for _, pattern := range QuotationPatterns {
		resultat = pattern.ReplaceAllString(resultat, "")
	}

	// Nettoyer les espaces en double
	resultat = regexp.MustCompile(`\s+`).ReplaceAllString(resultat, " ")
	resultat = strings.TrimSpace(resultat)

	return resultat
}

// ContientCitation vérifie si un texte contient une citation
func ContientCitation(texte string) bool {
	for _, pattern := range QuotationPatterns {
		if pattern.MatchString(texte) {
			return true
		}
	}
	return false
}

// TransformerEnAbstrait monte une phrase en abstraction
func TransformerEnAbstrait(phrase string, analyse AnalyseSemantique) string {
	// Si déjà abstraite, ne pas modifier
	if EstAbstraite(phrase) {
		return phrase
	}

	// Essayer de mapper à un concept trouvé
	_ = strings.ToLower(phrase) // phraseLower non utilisé actuellement

	for _, concept := range analyse.Concepts {
		if concept.Frequency > 0 {
			// Remplacer des éléments concrets par des concepts abstraits
			switch concept.Type {
			case "theme":
				return fmt.Sprintf("Cet extrait illustre le thème de la %s : %s",
					concept.Mot, strings.TrimSuffix(phrase, "."))
			}
		}
	}

	return phrase
}

// ScoreAbstraction évalue le niveau d'abstraction d'un résumé
type ScoreAbstraction struct {
	PresenceConceptsAbstraits float64 // 0-100%
	AbsenceCitations          float64 // 0-100%
	VerbsAbstraction          float64 // 0-100%
	PresenceThemes            float64 // 0-100%
	ScoreGlobal               float64 // 0-100%
}

// EvaluerAbstraction calcule le score global d'abstraction
func EvaluerAbstraction(resume string, analyse AnalyseSemantique) ScoreAbstraction {
	score := ScoreAbstraction{}

	resumeLower := strings.ToLower(resume)

	// 1. Présence de concepts abstraits
	conceptsPresents := 0
	for concept := range analyse.Abstractions {
		if strings.Contains(resumeLower, concept) {
			conceptsPresents++
		}
	}
	if len(analyse.Abstractions) > 0 {
		score.PresenceConceptsAbstraits = float64(conceptsPresents) / float64(len(analyse.Abstractions)) * 100.0
	}

	// 2. Absence de citations (bonus)
	if !ContientCitation(resume) {
		score.AbsenceCitations = 100.0
	} else {
		score.AbsenceCitations = 0.0
	}

	// 3. Présence de verbes d'abstraction
	verbsFound := 0
	for _, verb := range VerbsAbstraction {
		if strings.Contains(resumeLower, verb) {
			verbsFound++
		}
	}
	score.VerbsAbstraction = float64(verbsFound) / float64(len(VerbsAbstraction)) * 100.0

	// 4. Présence de thèmes
	themesFound := 0
	for _, theme := range analyse.Themes {
		if strings.Contains(resumeLower, strings.ToLower(theme)) {
			themesFound++
		}
	}
	if len(analyse.Themes) > 0 {
		score.PresenceThemes = float64(themesFound) / float64(len(analyse.Themes)) * 100.0
	}

	// Score global = moyenne pondérée
	score.ScoreGlobal = (score.PresenceConceptsAbstraits*0.35 +
		score.AbsenceCitations*0.25 +
		score.VerbsAbstraction*0.25 +
		score.PresenceThemes*0.15)

	return score
}

// AfficherAnalyseSemantique affiche les résultats de l'analyse
func AfficherAnalyseSemantique(analyse AnalyseSemantique) string {
	output := ""

	output += fmt.Sprintf("\n╔════════════════════════════════════════════════════════════╗\n")
	output += fmt.Sprintf("║  ANALYSE SÉMANTIQUE - ABSTRACTION (%.1f%% abstrait)        ║\n", analyse.ScoreAbstrait)
	output += fmt.Sprintf("╚════════════════════════════════════════════════════════════╝\n\n")

	// Concepts trouvés
	if len(analyse.Concepts) > 0 {
		output += fmt.Sprintf("[CONCEPTS ABSTRAITS DÉTECTÉS]\n")
		for _, concept := range analyse.Concepts {
			output += fmt.Sprintf("  • %s (score: %.1f, fréquence: %d)\n",
				concept.Mot, concept.Score, concept.Frequency)
			if len(concept.Exemples) > 0 {
				for i, ex := range concept.Exemples {
					if i >= 1 {
						break
					}
					output += fmt.Sprintf("    → \"%s\"\n", ex)
				}
			}
		}
		output += "\n"
	}

	// Causes
	if len(analyse.Causes) > 0 {
		output += fmt.Sprintf("[CAUSES IDENTIFIÉES]\n")
		for i, cause := range analyse.Causes {
			if i >= 2 {
				break
			}
			output += fmt.Sprintf("  • %s\n", cause)
		}
		output += "\n"
	}

	// Objectifs
	if len(analyse.Objectifs) > 0 {
		output += fmt.Sprintf("[OBJECTIFS DÉTECTÉS]\n")
		for i, obj := range analyse.Objectifs {
			if i >= 2 {
				break
			}
			output += fmt.Sprintf("  • %s\n", obj)
		}
		output += "\n"
	}

	// Critiques
	if len(analyse.Critiques) > 0 {
		output += fmt.Sprintf("[CRITIQUES / DÉNONCIATIONS]\n")
		for i, crit := range analyse.Critiques {
			if i >= 2 {
				break
			}
			output += fmt.Sprintf("  • %s\n", crit)
		}
		output += "\n"
	}

	return output
}
