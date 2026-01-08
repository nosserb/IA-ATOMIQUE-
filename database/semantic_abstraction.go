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

	// Connecteurs explicites à réduire (Phase X+3: Natural Syntax Layer)
	ConnecteursExplicites = []string{
		// Connecteurs temporels
		"alors", "ensuite", "puis", "maintenant", "avant", "après",
		// Connecteurs d'addition
		"de plus", "en outre", "d'ailleurs", "également", "aussi",
		// Connecteurs de conséquence
		"donc", "en conséquence", "par suite", "dès lors", "c'est pourquoi", "ainsi",
		// Connecteurs de contraste
		"cependant", "néanmoins", "pourtant", "toutefois", "or",
		// Connecteurs de révélation
		"révélant ainsi", "révélant", "montrant ainsi",
	}
)

// ============= PHASE X+3: NATURAL SYNTAX LAYER =============
// Humanise la syntaxe en réduisant connecteurs explicites et variant les structures

// RephraseSyntax réécrit le résumé avec syntaxe naturelle
func RephraseSyntax(texte string) string {
	// 1. Détecter les connecteurs explicites
	phrases := strings.Split(texte, ".")
	if len(phrases) <= 1 {
		return texte
	}

	// 2. Nettoyer et restructurer chaque phrase
	var phrasesRestructurees []string
	for i, phrase := range phrases {
		phrase = strings.TrimSpace(phrase)
		if len(phrase) < 5 {
			continue
		}

		// Restructurer selon la position et la longueur
		restructuree := RestructurerPhrase(phrase, i, len(phrases))
		if restructuree != "" {
			phrasesRestructurees = append(phrasesRestructurees, restructuree)
		}
	}

	// 3. Fusionner phrases courtes avec participes
	result := FusionnerPhrasesAvecParticipes(phrasesRestructurees)

	return result
}

// RestructurerPhrase reformule une phrase pour réduire connecteurs explicites
func RestructurerPhrase(phrase string, index int, total int) string {
	phraseLower := strings.ToLower(phrase)

	// Supprimer les connecteurs explicites au début (case-insensitive)
	for _, connecteur := range ConnecteursExplicites {
		// "Alors le profit" ou "alors le profit" → "le profit"
		if strings.HasPrefix(phraseLower, connecteur+" ") || strings.HasPrefix(phraseLower, strings.ToLower(connecteur)+" ") {
			// Trouver la vraie position du connecteur dans la phrase d'origine
			idx := strings.Index(strings.ToLower(phrase), strings.ToLower(connecteur)+" ")
			if idx == 0 {
				phrase = strings.TrimSpace(phrase[len(connecteur):])
				break
			}
		}
	}

	return strings.TrimSpace(phrase)
}

// TransformerEnRelative convertit "L'inégalité est programmée" en relative
func TransformerEnRelative(phrase string) string {
	// "L'inégalité n'est pas accidentelle mais programmée"
	// → "qui n'est pas accidentelle mais programmée" (pour fusionner après)

	// Si commence par "la" ou "le", le transformer en relative
	if strings.HasPrefix(strings.ToLower(phrase), "la ") {
		// Garder le "qui" pour fusion later
		return "qui " + strings.TrimPrefix(strings.ToLower(phrase), "la ")
	}
	if strings.HasPrefix(strings.ToLower(phrase), "le ") {
		return "qui " + strings.TrimPrefix(strings.ToLower(phrase), "le ")
	}
	if strings.HasPrefix(strings.ToLower(phrase), "l'") {
		return "qui " + strings.TrimPrefix(strings.ToLower(phrase), "l'")
	}

	return phrase
}

// FusionnerPhrasesAvecParticipes fusionne phrases courtes avec participes
func FusionnerPhrasesAvecParticipes(phrases []string) string {
	if len(phrases) == 0 {
		return ""
	}

	if len(phrases) == 1 {
		return phrases[0] + "."
	}

	var result []string
	i := 0

	for i < len(phrases) {
		current := phrases[i]

		// Chercher si la phrase suivante est une relative (commence par "qui")
		if i+1 < len(phrases) && strings.HasPrefix(strings.ToLower(phrases[i+1]), "qui ") {
			// Fusionner: "L'inégalité, qui n'est pas accidentelle..."
			next := phrases[i+1]
			fusion := current + ", " + next
			result = append(result, fusion)
			i += 2
			continue
		}

		// Chercher si la phrase suivante peut être convertie en participe
		if i+1 < len(phrases) && CanConvertToParticiple(phrases[i+1]) {
			next := phrases[i+1]
			participle := ConvertToParticiple(next)
			fusion := current + ", " + participle
			result = append(result, fusion)
			i += 2
			continue
		}

		// Sinon, phrase normale
		result = append(result, current)
		i++
	}

	// Joindre avec variations de ponctuation
	return JoinWithNaturalRhythm(result)
}

// CanConvertToParticiple vérifie si une phrase peut être convertie en participe
func CanConvertToParticiple(phrase string) bool {
	phraseLower := strings.ToLower(phrase)
	// Peut être convertie si elle parle de conséquence ou transformation
	return strings.Contains(phraseLower, "se construit") ||
		strings.Contains(phraseLower, "transformant") ||
		strings.Contains(phraseLower, "révélant") ||
		strings.Contains(phraseLower, "forçant") ||
		strings.Contains(phraseLower, "impose")
}

// ConvertToParticiple convertit une phrase en participe présent
func ConvertToParticiple(phrase string) string {
	phraseLower := strings.ToLower(phrase)

	// "Le profit se construit sur la précarité"
	// → "transformant le profit en construction sur la précarité"

	// Chercher le verbe principal
	if strings.Contains(phraseLower, "se construit") {
		return "transformant " + strings.ToLower(strings.TrimPrefix(phrase, "le ")) +
			", " + strings.ToLower(strings.TrimPrefix(phrase, "le "))
	}

	// "Le renoncement forcé est une violence"
	// → "révélant le renoncement forcé comme une forme de violence"
	if strings.Contains(phraseLower, " est ") {
		return "révélant " + strings.ToLower(phrase)
	}

	return phrase
}

// JoinWithNaturalRhythm joint les phrases avec variation de ponctuation
func JoinWithNaturalRhythm(phrases []string) string {
	if len(phrases) == 0 {
		return ""
	}

	var output string

	for i, phrase := range phrases {
		phrase = strings.TrimSpace(phrase)

		if i == 0 {
			// Première phrase: capitaliser
			output = CapitalizeFirst(phrase)
		} else if i%3 == 1 {
			// Variation 1: point-virgule (liaison logique)
			output += "; " + strings.ToLower(phrase)
		} else if i%3 == 2 {
			// Variation 2: nouvelle phrase (respiration)
			output += ". " + CapitalizeFirst(phrase)
		} else {
			// Variation 3: virgule (continuation)
			output += ", " + strings.ToLower(phrase)
		}
	}

	if !strings.HasSuffix(output, ".") {
		output += "."
	}

	return output
}

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

// ContientCitation vérifie si un texte contient une citation (explicite ou implicite)
func ContientCitation(texte string) bool {
	// Citations explicites (guillemets, tirets)
	for _, pattern := range QuotationPatterns {
		if pattern.MatchString(texte) {
			return true
		}
	}

	// Citations implicites: si une phrase contient trop de mots du texte original
	// C'est une heuristique mais efficace pour détecter les résumés trop littéraux
	// Par défaut on retourne false pour les patterns de guillemets
	// mais cette fonction sera améliorée en Phase X+2
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

// ============= PHASE X+1: CONCEPT → PHRASE ABSTRAITE =============

// PhrasesConceptuellesPar concept → phrases abstraites génériques mais puissantes
var PhrasesConceptuelles = map[string][]string{
	"misère": {
		"La misère économique contraint l'individu à des choix destructeurs.",
		"La pauvreté structure les comportements de survie.",
		"L'indigence impose des sacrifices inévitables.",
		"La dénuement engendre des stratégies de désespoir.",
	},
	"exploitation": {
		"Le système social exploite la vulnérabilité des plus faibles.",
		"L'exploitation économique est un mécanisme de domination systémique.",
		"Le profit se construit sur la précarité d'autrui.",
		"L'extraction de valeur repose sur l'asymétrie des pouvoirs.",
	},
	"injustice": {
		"L'injustice structurelle maintient les inégalités générationnelles.",
		"Les systèmes institutionnels reproduisent les discriminations.",
		"L'inégalité n'est pas accidentelle mais programmée.",
		"La discrimination est inscrite dans les mécanismes de pouvoir.",
	},
	"oppression": {
		"L'oppression repose sur la normalisation de la domination.",
		"Le système oppressif rend invisible sa propre violence.",
		"L'asservissement social se perpétue par des mécanismes implicites.",
		"La domination s'exerce par l'internationalisation des normes.",
	},
	"sacrifice": {
		"Le sacrifice devient un mode de survie imposé par les circonstances.",
		"L'abnégation est exigée de ceux qui n'ont rien à donner.",
		"Le renoncement forcé est une forme de violence silencieuse.",
		"Le dévouement demandé aux exploités justifie leur exploitation.",
	},
	"violence": {
		"La violence structurelle s'exerce sans visage et sans trace.",
		"La brutalité systémique se cache sous l'apparence de normalité.",
		"La violence économique détruit aussi sûrement que la violence physique.",
		"La maltraitance institutionnelle est la forme moderne de la servitude.",
	},
	"rôle": {
		"Les rôles assignés figent les trajectoires sociales.",
		"L'identité sociale détermine les libertés et les chaînes.",
		"Les rôles imposés légitiment les hiérarchies.",
		"La condition sociale réduit l'individu à une fonction.",
	},
}

// GenererPhrasesConceptuelles crée un résumé basé UNIQUEMENT sur concepts abstraits
func GenererPhrasesConceptuelles(analyse AnalyseSemantique) []string {
	phrasesAbstraites := make([]string, 0)

	// Pour chaque concept détecté, générer une phrase abstraite
	for _, concept := range analyse.Concepts {
		if templates, ok := PhrasesConceptuelles[concept.Mot]; ok {
			// Choisir une phrase template basée sur score du concept
			idx := int(concept.Score*10) % len(templates)
			phrasesAbstraites = append(phrasesAbstraites, templates[idx])
		}
	}

	// Si pas assez de phrases, ajouter des phrases génériques basées sur les thèmes
	if len(phrasesAbstraites) < 2 && len(analyse.Themes) > 0 {
		for _, theme := range analyse.Themes {
			phrasesAbstraites = append(phrasesAbstraites,
				fmt.Sprintf("Le thème central de %s révèle les mécanismes sous-jacents de la société.", theme))
		}
	}

	return phrasesAbstraites
}

// AppliquerAbstractionForced remplace les phrases concrètes par des abstraites si score < 60%
func AppliquerAbstractionForcee(resumeOriginal string, scoreAbstraction float64, analyse AnalyseSemantique) string {
	// Si score est bon, garder le résumé original
	if scoreAbstraction >= 60.0 {
		return resumeOriginal
	}

	// Sinon, remplacer par phrases conceptuelles
	phrasesAbstraites := GenererPhrasesConceptuelles(analyse)

	if len(phrasesAbstraites) == 0 {
		return resumeOriginal // Fallback si pas de concepts
	}

	output := ""
	output += fmt.Sprintf("╔════════════════════════════════════════════════════════════╗\n")
	output += fmt.Sprintf("║  ⚠️  RÉSUMÉ RÉÉCRIT EN ABSTRACTION FORCÉE                  ║\n")
	output += fmt.Sprintf("║  (Score %.1f%% < 60%% | Phrases réécrites conceptuellement)  ║\n", scoreAbstraction)
	output += fmt.Sprintf("╚════════════════════════════════════════════════════════════╝\n\n")

	for i, phrase := range phrasesAbstraites {
		output += fmt.Sprintf("%d. %s\n", i+1, phrase)
	}

	return output
}

// ============= PHASE X+2: CONCEPTUAL LINKING LAYER =============
// Relie les phrases conceptuelles avec une micro-structure discursive

// ConnecteursLiaison = connecteurs faibles pour lier les concepts
var ConnecteursLiaison = map[string][]string{
	"contrast": {
		"non comme", "mais comme", "contrairement à", "cependant",
		"or", "pourtant", "néanmoins", "toutefois",
	},
	"consequence": {
		"en conséquence", "dès lors", "partant", "donc",
		"par suite", "c'est pourquoi", "ainsi", "il en résulte que",
	},
	"cadre": {
		"dans ce cadre", "dans ce contexte", "de ce fait", "en ce sens",
		"sous cet angle", "de cette manière", "de ce point de vue",
	},
	"addition": {
		"de plus", "par ailleurs", "en outre", "d'ailleurs",
		"également", "aussi", "de surcroît", "ainsi que",
	},
	"revelation": {
		"révélant ainsi", "ce qui révèle", "mettant en lumière", "exposant",
		"dévoilant", "mettant en évidence", "illustrant", "témoignant de",
	},
}

// LierPhrasesConceptuelles regroupe et lie les phrases avec structure discursive
func LierPhrasesConceptuelles(phrases []string) string {
	if len(phrases) == 0 {
		return ""
	}

	if len(phrases) == 1 {
		p := strings.TrimSuffix(phrases[0], ".")
		return p + "."
	}

	// Nettoyer les phrases (enlever ponctuation inutile)
	var phrasesPropres []string
	for _, p := range phrases {
		p = strings.TrimSuffix(p, ".")
		p = strings.TrimSpace(p)
		phrasesPropres = append(phrasesPropres, p)
	}

	// Structure: Thèse + Mécanismes + Conclusion

	// === THÈSE (première phrase) ===
	these := phrasesPropres[0] + ","

	// === MÉCANISMES (grouper par 2-3) ===
	var mecanismes []string
	if len(phrasesPropres) > 2 {
		phrasesMiddle := phrasesPropres[1 : len(phrasesPropres)-1]

		for i := 0; i < len(phrasesMiddle); i += 2 {
			// Première phrase du groupe
			groupe := phrasesMiddle[i]

			// Ajouter la suivante avec connecteur si elle existe
			if i+1 < len(phrasesMiddle) {
				connecteurs := ConnecteursLiaison["addition"]
				connecteur := connecteurs[i%len(connecteurs)]
				groupe += " " + connecteur + " " + strings.ToLower(phrasesMiddle[i+1])
			}

			mecanismes = append(mecanismes, groupe)
		}
	}

	// === CONCLUSION (dernière phrase) ===
	conclusion := phrasesPropres[len(phrasesPropres)-1]

	// === ASSEMBLER LE DISCOURS AVEC CONNECTEURS ===
	output := ""

	// Ajouter thèse
	output += these

	// Ajouter mécanismes avec connecteurs appropriés
	if len(mecanismes) > 0 {
		for i, meca := range mecanismes {
			output += " "
			if i == 0 {
				// Premier mécanisme: cadre de contexte
				output += "dans ce cadre, " + strings.ToLower(meca) + "."
			} else {
				// Autres: conséquences
				connecteurs := ConnecteursLiaison["consequence"]
				connecteur := connecteurs[(i-1)%len(connecteurs)]
				output += connecteur + " " + strings.ToLower(meca) + "."
			}
		}
	}

	// Ajouter conclusion avec connecteur de révélation
	output += " "
	connecteurs := ConnecteursLiaison["revelation"]
	connecteur := connecteurs[0]
	output += connecteur + ", " + strings.ToLower(conclusion) + "."

	return strings.TrimSpace(output)
}

// ============= PHASE X+3: NATURAL SYNTAX LAYER =============
// Reformule les phrases avec subordination/ponctuation, pas de connecteurs explicites

// HumanizeStructure convertit une structure avec connecteurs explicites en syntaxe naturelle
func HumanizeStructure(phrasesAbstraites []string) string {
	if len(phrasesAbstraites) == 0 {
		return ""
	}

	if len(phrasesAbstraites) == 1 {
		p := strings.TrimSuffix(phrasesAbstraites[0], ".")
		return p + "."
	}

	// Nettoyer
	var phrases []string
	for _, p := range phrasesAbstraites {
		p = strings.TrimSuffix(p, ".")
		p = strings.TrimSpace(p)
		phrases = append(phrases, p)
	}

	// === STRATÉGIE: Fusionner avec subordination plutôt que connecteurs ===
	// Objectif: réduire de 40-60% l'usage de connecteurs explicites
	// Les remplacer par: subordination ("car"), relatives ("qui"), ponctuation, reformulation

	// Première phrase (thèse)
	these := strings.ToLower(phrases[0])
	these = CapitalizeFirst(these)

	// Cas 2 phrases: simple subordination
	if len(phrases) == 2 {
		second := strings.ToLower(phrases[1])
		output := these + ", car " + second + "."
		return output
	}

	// Cas 3+ phrases: fusionner progressivement avec variation de syntaxe
	var segments []string

	// Segment 1: Thèse + 1ère mécanique avec subordination
	if len(phrases) > 1 {
		second := strings.ToLower(phrases[1])
		// Remplacer "est" par "qui rend/qui crée/qui impose"
		second = AddRelativeClause(second)
		seg1 := these + ", car " + second + "."
		segments = append(segments, seg1)
	}

	// Segments intermédiaires: transformer en participes ou nouvelles phrases
	for i := 2; i < len(phrases); i++ {
		phrase := strings.ToLower(phrases[i])

		if i == len(phrases)-1 {
			// Dernière phrase: participiale explicite
			phrase = "cette logique révèle " + strings.TrimPrefix(phrase, "la ")
			segments = append(segments, CapitalizeFirst(phrase)+".")
		} else {
			// Phrases du milieu: transformer en participes
			// "Le profit se construit" → "Le profit se construit alors sur..."
			phrase = AddTransformParticiple(phrase)
			segments = append(segments, CapitalizeFirst(phrase)+".")
		}
	}

	// Joindre puis appliquer Phase X+3: humanisation syntaxique
	resultat := strings.Join(segments, "\n")

	// Appliquer l'humanisation syntaxique (réduire connecteurs, varier rythme)
	resultat = RephraseSyntax(resultat)

	return resultat
}

// AddRelativeClause ajoute une relative ("qui", "que", "dont") à une phrase
func AddRelativeClause(phrase string) string {
	// "l'inégalité n'est pas accidentelle mais programmée"
	// → "l'inégalité qui n'est pas accidentelle mais programmée"
	// ou en gardant la structure: "qui rend l'inégalité structurelle"

	// Si la phrase contient "est" ou "n'est"
	if strings.Contains(phrase, " est ") || strings.Contains(phrase, " n'est") {
		// Chercher le sujet (avant "est")
		parts := strings.Split(phrase, " est ")
		if len(parts) == 2 {
			subject := strings.TrimPrefix(parts[0], "la ")

			// "l'inégalité n'est pas accidentelle" → "qui rend l'inégalité structurelle plutôt qu'accidentelle"
			return "qui rend " + subject + " structurelle plutôt qu'accidentelle"
		}
	}

	// Si la phrase contient "s'exerce"
	if strings.Contains(phrase, "s'exerce") {
		// "la domination s'exerce par..." → "qui s'exerce par..."
		return "qui " + strings.TrimPrefix(phrase, "la ")
	}

	return phrase
}

// AddTransformParticiple ajoute "alors" ou reformule en participe
func AddTransformParticiple(phrase string) string {
	// Simplement ajouter "alors" pour marquer une nouvelle assertion sans connecteur explicite
	// Pas de transformation complexe ici
	phrase = strings.TrimSpace(phrase)
	if strings.HasPrefix(phrase, "le ") || strings.HasPrefix(phrase, "la ") {
		return "alors " + phrase
	}
	return phrase
}

// CapitalizeFirst met en majuscule le premier caractère
func CapitalizeFirst(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
