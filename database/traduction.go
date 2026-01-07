package database

import (
	"strings"
)

// TraductionMap contient les traductions simples EN → FR
// Pour une production, utiliser LibreTranslate, DeepL ou Google Translate
var TraductionMap = map[string]map[string]string{
	"en": {
		// Fichiers & structure
		"file":     "fichier",
		"document": "document",
		"page":     "page",
		"ebook":    "livre électronique",
		"book":     "livre",
		"chapter":  "chapitre",

		// Production & data
		"produced":    "produit",
		"production":  "production",
		"from":        "de",
		"images":      "images",
		"data":        "données",
		"information": "information",
		"analysis":    "analyse",
		"content":     "contenu",

		// Metadata
		"title":       "titre",
		"author":      "auteur",
		"version":     "version",
		"date":        "date",
		"time":        "heure",
		"release":     "sortie",
		"published":   "publié",
		"publication": "publication",
		"original":    "original",

		// Scientific terms
		"essay":          "essai",
		"view":           "point de vue",
		"art":            "art",
		"dying":          "teinture",
		"painting":       "peinture",
		"combustion":     "combustion",
		"hypothesis":     "hypothèse",
		"hypotheses":     "hypothèses",
		"erroneous":      "erroné",
		"proven":         "prouvé",
		"phlogistic":     "phlogistique",
		"antiphlogistic": "antiphlogistique",
		"theory":         "théorie",
		"demonstration":  "démonstration",
		"experiment":     "expérience",
		"experiments":    "expériences",
		"result":         "résultat",
		"results":        "résultats",
		"effect":         "effet",
		"effects":        "effets",
		"cause":          "cause",
		"causes":         "causes",
		"principle":      "principe",
		"principles":     "principes",
		"substance":      "substance",
		"substancees":    "substances",
		"matter":         "matière",
		"motion":         "mouvement",
		"force":          "force",
		"forces":         "forces",
		"reaction":       "réaction",
		"reactions":      "réactions",
		"process":        "processus",
		"processes":      "processus",
		"method":         "méthode",
		"methods":        "méthodes",
		"observation":    "observation",
		"observations":   "observations",
		"conclusion":     "conclusion",
		"conclusions":    "conclusions",

		// Descriptive
		"new":     "nouveau",
		"first":   "premier",
		"second":  "deuxième",
		"third":   "troisième",
		"upon":    "sur",
		"between": "entre",
		"towards": "vers",
		"without": "sans",
		"within":  "dans",
		"wherein": "où",

		// Actions & metadata
		"warning":   "avertissement",
		"error":     "erreur",
		"example":   "exemple",
		"note":      "note",
		"important": "important",
		"please":    "s'il vous plaît",
		"thank":     "merci",
		"help":      "aide",
		"support":   "support",
		"contact":   "contact",
		"follow":    "suivre",
		"next":      "suivant",
		"previous":  "précédent",
		"home":      "accueil",
		"about":     "à propos",
		"search":    "recherche",
		"filter":    "filtre",
		"sort":      "trier",
		"download":  "télécharger",
		"upload":    "téléverser",
		"delete":    "supprimer",
		"create":    "créer",
		"edit":      "éditer",
		"save":      "enregistrer",
		"cancel":    "annuler",
		"ok":        "ok",
		"yes":       "oui",
		"no":        "non",
		"back":      "retour",
		"forward":   "avant",
		"close":     "fermer",
		"open":      "ouvrir",
		"success":   "succès",
		"failed":    "échoué",

		// Common words
		"the":   "le",
		"a":     "un",
		"of":    "de",
		"and":   "et",
		"to":    "à",
		"in":    "en",
		"is":    "est",
		"that":  "que",
		"it":    "c'est",
		"for":   "pour",
		"on":    "sur",
		"with":  "avec",
		"by":    "par",
		"be":    "être",
		"this":  "ce",
		"all":   "tous",
		"are":   "sont",
		"was":   "était",
		"were":  "étaient",
		"been":  "été",
		"their": "leur",
	},
	"de": {
		"datei":    "fichier",
		"daten":    "données",
		"analyse":  "analyse",
		"dokument": "document",
		"seite":    "page",
		"inhalt":   "contenu",
		"titel":    "titre",
		"autor":    "auteur",
		"version":  "version",
		"datum":    "date",
		"zeit":     "heure",
		"warnung":  "avertissement",
		"fehler":   "erreur",
		"beispiel": "exemple",
		"notiz":    "note",
		"wichtig":  "important",
		"bitte":    "s'il vous plaît",
		"danke":    "merci",
		"hilfe":    "aide",
		"kontakt":  "contact",
	},
	"es": {
		"archivo":     "fichier",
		"datos":       "données",
		"análisis":    "analyse",
		"documento":   "document",
		"página":      "page",
		"contenido":   "contenu",
		"título":      "titre",
		"autor":       "auteur",
		"versión":     "version",
		"fecha":       "date",
		"hora":        "heure",
		"advertencia": "avertissement",
		"error":       "erreur",
		"ejemplo":     "exemple",
		"nota":        "note",
		"importante":  "important",
		"por favor":   "s'il vous plaît",
		"gracias":     "merci",
		"ayuda":       "aide",
		"contacto":    "contact",
	},
}

// TraduireSiNecessaire traduit une phrase vers le FR si elle n'y est pas
// Retourne la phrase traduite et le facteur de confiance γi
func TraduireSiNecessaire(phrase *Phrase, langue string) (string, bool, float64) {
	// Si déjà en FR, pas de traduction
	if langue == "fr" {
		return phrase.Contenu, false, 1.0 // γi = 1.0 (confiance totale)
	}

	// Traduction simple: remplacer mots individuels
	contenuTraduit := TraduireMotsPar(phrase.Contenu, langue)

	// Facteur de confiance pour traductions auto
	// γi = 0.8 pour traductions (moins confiantes que texte original)
	// γi = 0.7 pour traductions longues (plus d'erreurs potentielles)
	confidenceFacteur := 0.8
	if len(strings.Fields(phrase.Contenu)) > 10 {
		confidenceFacteur = 0.7
	}

	return contenuTraduit, true, confidenceFacteur
}

// TraduireMotsPar traduit les mots d'une phrase selon la langue source
func TraduireMotsPar(texte string, languageSource string) string {
	traductions := TraductionMap[languageSource]
	if traductions == nil {
		return texte // Pas de traduction disponible
	}

	mots := strings.Fields(texte)
	for i, mot := range mots {
		motLower := strings.ToLower(mot)
		// Enlever la ponctuation pour chercher
		motClean := strings.TrimFunc(motLower, func(r rune) bool {
			return !isLetterOrNumber(r)
		})

		if traduit, ok := traductions[motClean]; ok {
			// Replacer avec la première lettre en majuscule si original l'était
			if len(mot) > 0 && isUpperCase(mot[0]) {
				mots[i] = strings.ToUpper(traduit[:1]) + traduit[1:]
			} else {
				mots[i] = traduit
			}
		}
	}

	return strings.Join(mots, " ")
}

// DetecterEtTraduirePhrases détecte la langue de chaque phrase et traduit en FR
// Rfinal = Fusion({Pi traduites et filtrées})
func DetecterEtTraduirePhrases(phrases []Phrase) []Phrase {
	for i := range phrases {
		// Déterminer la langue de chaque phrase
		// Note: Pour une meilleure détection, importer DetecterLangue de interaction.go
		// Pour maintenant, utiliser une détection simple par mots-clés
		langue := DetecterLanguePhrase(phrases[i].Contenu)
		phrases[i].Langue = langue

		// Traduire si nécessaire
		if langue != "fr" {
			contenuTraduit, estTraduit, gamma := TraduireSiNecessaire(&phrases[i], langue)
			phrases[i].Contenu = contenuTraduit
			phrases[i].Mots = strings.Fields(contenuTraduit)
			phrases[i].MotsClés = ExtraireMotsClés(phrases[i].Mots)
			phrases[i].EstTraduire = estTraduit
			phrases[i].FacteurConfiance = gamma
		} else {
			phrases[i].EstTraduire = false
			phrases[i].FacteurConfiance = 1.0 // Pas de traduction = confiance totale
		}
	}

	return phrases
}

// DetecterLanguePhrase détecte la langue d'une phrase simple
// Cherche des mots-clés typiques de chaque langue
func DetecterLanguePhrase(phrase string) string {
	phraseLower := strings.ToLower(phrase)

	// Mots-clés typiques FR
	motsFR := []string{"le", "la", "les", "un", "une", "des", "de", "et", "est", "que", "qui", "où", "à"}
	// Mots-clés typiques EN
	motsEN := []string{"the", "is", "and", "or", "be", "of", "in", "to", "a"}
	// Mots-clés typiques DE
	motsDE := []string{"der", "die", "das", "den", "dem", "des", "ein", "eine", "und", "ist", "zu"}
	// Mots-clés typiques ES
	motsES := []string{"el", "la", "los", "las", "un", "una", "unos", "unas", "de", "y", "es", "que"}

	countFR := countKeywords(phraseLower, motsFR)
	countEN := countKeywords(phraseLower, motsEN)
	countDE := countKeywords(phraseLower, motsDE)
	countES := countKeywords(phraseLower, motsES)

	// Retourner la langue avec le score le plus élevé
	if countFR > countEN && countFR > countDE && countFR > countES {
		return "fr"
	} else if countEN > countDE && countEN > countES {
		return "en"
	} else if countDE > countES {
		return "de"
	} else if countES > 0 {
		return "es"
	}

	return "fr" // Défaut: supposer FR
}

// countKeywords compte le nombre de mots-clés présents dans le texte
// Utilise word boundaries pour éviter les faux positifs
func countKeywords(texte string, keywords []string) int {
	count := 0
	words := strings.FieldsFunc(texte, func(r rune) bool {
		return !isLetterOrNumber(r)
	})

	for _, kw := range keywords {
		for _, word := range words {
			if strings.ToLower(word) == kw {
				count++
			}
		}
	}
	return count
}

// isLetterOrNumber vérifie si un rune est une lettre ou un chiffre
func isLetterOrNumber(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == 'é' || r == 'è' || r == 'ê' || r == 'ô' || r == 'ù' || r == 'ç'
}

// isUpperCase vérifie si un byte est en majuscule
func isUpperCase(b byte) bool {
	return b >= 'A' && b <= 'Z'
}
