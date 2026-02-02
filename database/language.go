package database

import (
	"fmt"
	"regexp"
	"strings"
)

// Types pour tri optimisé
type phraseScore struct {
	texte string
	score float64
}

type motFreq struct {
	mot  string
	freq int
}

// Types pour analyse avancée (3 couches)
type TypePhrase int

const (
	FAIT TypePhrase = iota
	INTERPRETATION
	CONSEQUENCE
	NEUTRE
)

type AxeSemantique string

// ============= NOUVEAU: Système des 5 piliers scientifiques =============

// FonctionScientifique catégorise chaque phrase selon sa fonction logique
type FonctionScientifique int

const (
	CONTEXTE FonctionScientifique = iota
	PROBLEME
	OBJECTIF
	APPROCHE
	APPORT
	NON_CLASSE
)

func (f FonctionScientifique) String() string {
	switch f {
	case CONTEXTE:
		return "CONTEXTE"
	case PROBLEME:
		return "PROBLÈME"
	case OBJECTIF:
		return "OBJECTIF"
	case APPROCHE:
		return "APPROCHE"
	case APPORT:
		return "APPORT"
	default:
		return "NON_CLASSÉ"
	}
}

// PhraseFonctionnelle représente une phrase avec sa fonction scientifique
type PhraseFonctionnelle struct {
	Texte    string
	Fonction FonctionScientifique
	Score    float64
	Source   string // Phrase originale
}

// ResumePiliers structure le résumé selon les 5 piliers
type ResumePiliers struct {
	Contexte []PhraseFonctionnelle
	Probleme []PhraseFonctionnelle
	Objectif []PhraseFonctionnelle
	Approche []PhraseFonctionnelle
	Apport   []PhraseFonctionnelle
	ScoreSc  float64 // Score de résumabilité scientifique
}

const (
	EMOTIONNEL    AxeSemantique = "Émotionnel"
	EDUCATIF      AxeSemantique = "Éducatif"
	PSYCHOLOGIQUE AxeSemantique = "Psychologique"
	RELATIONNEL   AxeSemantique = "Relationnel"
	MORAL         AxeSemantique = "Moral"
	PRAGMATIQUE   AxeSemantique = "Pragmatique"
	COGNITIF      AxeSemantique = "Cognitif"
)

type PhraseAnalysee struct {
	Texte          string
	Type           TypePhrase
	AxesPrincipaux []AxeSemantique
	Score          float64
}

// ExtraireContenuHTMLPrincipal isole le contenu principal du HTML (étapes 1-3)
// Élimine menus, pubs, scripts et balises non-pertinentes
func ExtraireContenuHTMLPrincipal(htmlBrut string) string {
	contenu := htmlBrut

	// ÉTAPE 1: Supprimer les balises de bruit
	tagsASupprimer := []string{
		`(?i)<script[^>]*>.*?</script>`,                           // Scripts
		`(?i)<style[^>]*>.*?</style>`,                             // Styles
		`(?i)<nav[^>]*>.*?</nav>`,                                 // Navigation
		`(?i)<aside[^>]*>.*?</aside>`,                             // Barre latérale
		`(?i)<footer[^>]*>.*?</footer>`,                           // Pied de page
		`(?i)<header[^>]*>.*?</header>`,                           // En-tête
		`(?i)<sup[^>]*>.*?</sup>`,                                 // Superscript (références)
		`(?i)<sub[^>]*>.*?</sub>`,                                 // Subscript
		`(?i)<!--.*?-->`,                                          // Commentaires HTML
		`(?i)<noscript[^>]*>.*?</noscript>`,                       // Noscript
		`(?i)<meta[^>]*>`,                                         // Meta tags
		`(?i)<link[^>]*>`,                                         // Links
		`(?i)<img[^>]*>`,                                          // Images
		`(?i)class="(mw-ref|reference|reflist)"[^>]*>.*?</[^>]*>`, // Références wiki
	}

	for _, pattern := range tagsASupprimer {
		re := regexp.MustCompile(pattern)
		contenu = re.ReplaceAllString(contenu, "")
	}

	// ÉTAPE 2: Extraire le contenu principal (priorité: <main>, <article>, <div class="content">)
	reMain := regexp.MustCompile(`(?i)<main[^>]*>(.*?)</main>`)
	matchesMain := reMain.FindStringSubmatch(contenu)
	if len(matchesMain) > 1 {
		contenu = matchesMain[1]
	} else {
		reArticle := regexp.MustCompile(`(?i)<article[^>]*>(.*?)</article>`)
		matchesArticle := reArticle.FindStringSubmatch(contenu)
		if len(matchesArticle) > 1 {
			contenu = matchesArticle[1]
		}
	}

	// ÉTAPE 3: Nettoyer les balises HTML restantes
	contenu = nettoyerBalises(contenu)

	// ÉTAPE 4: Normaliser les espaces et retours à la ligne
	contenu = regexp.MustCompile(`\s+`).ReplaceAllString(contenu, " ")
	contenu = strings.TrimSpace(contenu)

	return contenu
}

// nettoyerBalises supprime les balises HTML tout en préservant le texte
func nettoyerBalises(texte string) string {
	// Remplacer les balises de paragraphes/sections par des séparateurs
	texte = regexp.MustCompile(`(?i)</p>`).ReplaceAllString(texte, ". ")
	texte = regexp.MustCompile(`(?i)</div>`).ReplaceAllString(texte, " ")
	texte = regexp.MustCompile(`(?i)<br[^>]*>`).ReplaceAllString(texte, ". ")
	texte = regexp.MustCompile(`(?i)</li>`).ReplaceAllString(texte, ". ")
	texte = regexp.MustCompile(`(?i)<li[^>]*>`).ReplaceAllString(texte, "• ")

	// Supprimer toutes les autres balises
	texte = regexp.MustCompile(`<[^>]+>`).ReplaceAllString(texte, "")

	// Supprimer les entités HTML
	texte = strings.ReplaceAll(texte, "&nbsp;", " ")
	texte = strings.ReplaceAll(texte, "&amp;", "&")
	texte = strings.ReplaceAll(texte, "&lt;", "<")
	texte = strings.ReplaceAll(texte, "&gt;", ">")
	texte = strings.ReplaceAll(texte, "&quot;", "\"")

	return texte
}

// IdentifierConceptsClésHTML analyse le contenu HTML pour extraire les concepts importants
// Utilise les titres (<h2>, <h3>) et la fréquence des termes
func IdentifierConceptsClésHTML(contenuHTML string) []string {
	conceptes := make(map[string]int)

	// Extraire les titres comme concepts clés (haute priorité)
	reTitres := regexp.MustCompile(`(?i)<h[2-4][^>]*>([^<]+)</h[2-4]>`)
	matches := reTitres.FindAllStringSubmatch(contenuHTML, -1)
	for _, match := range matches {
		if len(match) > 1 {
			titre := strings.TrimSpace(match[1])
			mots := strings.Fields(titre)
			for _, mot := range mots {
				if len(mot) > 4 { // Ignorer les mots courts
					conceptes[strings.ToLower(mot)] += 5 // Poids élevé pour les titres
				}
			}
		}
	}

	// Extraire du contenu principal via les balises <strong>, <b>, <em>
	reImportant := regexp.MustCompile(`(?i)<(strong|b|em)[^>]*>([^<]+)</\1>`)
	matchesImp := reImportant.FindAllStringSubmatch(contenuHTML, -1)
	for _, match := range matchesImp {
		if len(match) > 2 {
			texte := strings.TrimSpace(match[2])
			mots := strings.Fields(texte)
			for _, mot := range mots {
				if len(mot) > 4 {
					conceptes[strings.ToLower(mot)] += 3
				}
			}
		}
	}

	// Trier par fréquence et retourner top 10
	type conceptScore struct {
		concept string
		score   int
	}
	var sorted []conceptScore
	for concept, score := range conceptes {
		sorted = append(sorted, conceptScore{concept, score})
	}

	// Bubble sort
	for i := 0; i < len(sorted)-1; i++ {
		for j := 0; j < len(sorted)-1-i; j++ {
			if sorted[j].score < sorted[j+1].score {
				sorted[j], sorted[j+1] = sorted[j+1], sorted[j]
			}
		}
	}

	var result []string
	limit := 10
	if len(sorted) < 10 {
		limit = len(sorted)
	}
	for i := 0; i < limit; i++ {
		result = append(result, sorted[i].concept)
	}
	return result
}

// ClassifierPhraseScientiifque assigne une fonction scientifique à une phrase
func ClassifierPhraseScientiifque(texte string, motsCles []string) FonctionScientifique {
	texte = strings.ToLower(texte)

	// Mots-clés pour CONTEXTE
	contexteKeywords := []string{
		"domaine", "champ", "sujet", "topic", "étude", "research",
		"actuellement", "currently", "traditionnellement", "typically",
		"existant", "existing", "connu", "known", "établi", "established",
		"dans", "in the context", "au fil des années", "depuis",
	}

	// Mots-clés pour PROBLÈME
	problemeKeywords := []string{
		"cependant", "however", "mais", "but", "limite", "limitation",
		"problème", "problem", "défi", "challenge", "insuffisant", "insufficient",
		"nécessaire", "required", "manque", "lacking", "difficile", "difficult",
		"obstacle", "bottleneck", "critique", "critical", "reste", "remains",
		"n'", "ne ", "sans", "without", "absence", "lack",
	}

	// Mots-clés pour OBJECTIF
	objectifKeywords := []string{
		"proposer", "propose", "objectif", "objective", "but", "goal",
		"visée", "aim", "intention", "intent", "cherche", "seek", "cherche",
		"visant", "aiming", "pour", "to", "dans le but", "in order to",
		"développer", "develop", "créer", "create", "concevoir", "design",
		"améliorer", "improve", "résoudre", "solve", "aborder", "address",
		"permettre", "enable", "faciliter", "facilitate", "établir", "establish",
	}

	// Mots-clés pour APPROCHE
	approacheKeywords := []string{
		"méthode", "method", "approche", "approach", "technique", "technique",
		"utiliser", "use", "basé", "based", "fondé", "founded",
		"architecture", "architecture", "algorithme", "algorithm", "processus", "process",
		"étapes", "steps", "procédure", "procedure", "implémentation", "implementation",
		"framework", "framework", "système", "system", "stratégie", "strategy",
		"analyse", "analyze", "examiner", "examine", "évaluer", "evaluate",
		"tester", "test", "valider", "validate", "mesurer", "measure",
	}

	// Mots-clés pour APPORT
	apportKeywords := []string{
		"résultat", "result", "contribution", "contribution", "apport", "contribution",
		"nouveau", "novel", "original", "original", "unique", "unique",
		"amélioration", "improvement", "avantage", "advantage", "bénéfice", "benefit",
		"efficacité", "efficiency", "rapidité", "speed", "performance", "performance",
		"permet", "enables", "facilite", "facilitates", "montre", "shows",
		"démontre", "demonstrates", "prouves", "proves", "révèle", "reveals",
		"supérieur", "superior", "meilleur", "better", "remarquable", "remarkable",
		"significatif", "significant", "important", "important", "majeur", "major",
	}

	// Compter les correspondances
	scoreContexte := float64(0)
	scoreProbleme := float64(0)
	scoreObjectif := float64(0)
	scoreApproche := float64(0)
	scoreApport := float64(0)

	for _, keyword := range contexteKeywords {
		if strings.Contains(texte, keyword) {
			scoreContexte += 1.0
		}
	}

	for _, keyword := range problemeKeywords {
		if strings.Contains(texte, keyword) {
			scoreProbleme += 1.5 // Poids plus élevé
		}
	}

	for _, keyword := range objectifKeywords {
		if strings.Contains(texte, keyword) {
			scoreObjectif += 1.3
		}
	}

	for _, keyword := range approacheKeywords {
		if strings.Contains(texte, keyword) {
			scoreApproche += 1.2
		}
	}

	for _, keyword := range apportKeywords {
		if strings.Contains(texte, keyword) {
			scoreApport += 1.3
		}
	}

	// Trouver le score maximal
	maxScore := scoreContexte
	fonction := CONTEXTE

	if scoreProbleme > maxScore {
		maxScore = scoreProbleme
		fonction = PROBLEME
	}
	if scoreObjectif > maxScore {
		maxScore = scoreObjectif
		fonction = OBJECTIF
	}
	if scoreApproche > maxScore {
		maxScore = scoreApproche
		fonction = APPROCHE
	}
	if scoreApport > maxScore {
		maxScore = scoreApport
		fonction = APPORT
	}

	// Si aucun score significatif, retourner NON_CLASSE
	if maxScore < 0.5 {
		return NON_CLASSE
	}

	return fonction
}

// ClassifierPhrasesDuResume assigne une fonction à chaque phrase importante
func ClassifierPhrasesDuResume(phrases []string, analyse AnalysePhrases) []PhraseFonctionnelle {
	resultat := make([]PhraseFonctionnelle, 0, len(phrases))

	for _, phrase := range phrases {
		// Trouver la phrase source pour obtenir ses mots-clés
		var motsCles []string
		for _, p := range analyse.Phrases {
			if strings.Contains(p.Texte, phrase) || strings.Contains(phrase, p.Texte) {
				motsCles = p.MotsClés
				break
			}
		}

		fonction := ClassifierPhraseScientiifque(phrase, motsCles)
		if fonction != NON_CLASSE {
			resultat = append(resultat, PhraseFonctionnelle{
				Texte:    phrase,
				Fonction: fonction,
				Score:    1.0,
				Source:   phrase,
			})
		}
	}

	return resultat
}

// StructurerParPiliers organise les phrases selon les 5 piliers
func StructurerParPiliers(phrasesClassifiees []PhraseFonctionnelle) ResumePiliers {
	resume := ResumePiliers{
		Contexte: make([]PhraseFonctionnelle, 0),
		Probleme: make([]PhraseFonctionnelle, 0),
		Objectif: make([]PhraseFonctionnelle, 0),
		Approche: make([]PhraseFonctionnelle, 0),
		Apport:   make([]PhraseFonctionnelle, 0),
		ScoreSc:  0.0,
	}

	for _, pf := range phrasesClassifiees {
		switch pf.Fonction {
		case CONTEXTE:
			resume.Contexte = append(resume.Contexte, pf)
		case PROBLEME:
			resume.Probleme = append(resume.Probleme, pf)
		case OBJECTIF:
			resume.Objectif = append(resume.Objectif, pf)
		case APPROCHE:
			resume.Approche = append(resume.Approche, pf)
		case APPORT:
			resume.Apport = append(resume.Apport, pf)
		}
	}

	// Calculer le score de résumabilité scientifique
	// Bonus si tous les 5 piliers sont présents
	piliersPrésents := 0
	if len(resume.Contexte) > 0 {
		piliersPrésents++
	}
	if len(resume.Probleme) > 0 {
		piliersPrésents++
	}
	if len(resume.Objectif) > 0 {
		piliersPrésents++
	}
	if len(resume.Approche) > 0 {
		piliersPrésents++
	}
	if len(resume.Apport) > 0 {
		piliersPrésents++
	}

	resume.ScoreSc = float64(piliersPrésents) / 5.0 * 100.0

	return resume
}

// FormaterResumePiliers génère le texte final structuré
func FormaterResumePiliers(piliers ResumePiliers) string {
	output := ""

	output += fmt.Sprintf("╔════════════════════════════════════════════════════════════╗\n")
	output += fmt.Sprintf("║  RÉSUMÉ SCIENTIFIQUE (%.1f%% complet)                       ║\n", piliers.ScoreSc)
	output += fmt.Sprintf("╚════════════════════════════════════════════════════════════╝\n\n")

	// CONTEXTE
	if len(piliers.Contexte) > 0 {
		output += fmt.Sprintf("[CONTEXTE]\n")
		for i, pf := range piliers.Contexte {
			if i >= 2 {
				break
			}
			output += fmt.Sprintf("  • %s\n", pf.Texte)
		}
		output += "\n"
	}

	// PROBLÈME
	if len(piliers.Probleme) > 0 {
		output += fmt.Sprintf("[PROBLÈME IDENTIFIÉ]\n")
		for i, pf := range piliers.Probleme {
			if i >= 2 {
				break
			}
			output += fmt.Sprintf("  • %s\n", pf.Texte)
		}
		output += "\n"
	}

	// OBJECTIF
	if len(piliers.Objectif) > 0 {
		output += fmt.Sprintf("[OBJECTIF / BUT]\n")
		for i, pf := range piliers.Objectif {
			if i >= 2 {
				break
			}
			output += fmt.Sprintf("  • %s\n", pf.Texte)
		}
		output += "\n"
	}

	// APPROCHE
	if len(piliers.Approche) > 0 {
		output += fmt.Sprintf("[APPROCHE MÉTHODOLOGIQUE]\n")
		for i, pf := range piliers.Approche {
			if i >= 2 {
				break
			}
			output += fmt.Sprintf("  • %s\n", pf.Texte)
		}
		output += "\n"
	}

	// APPORT
	if len(piliers.Apport) > 0 {
		output += fmt.Sprintf("[APPORT / RÉSULTATS]\n")
		for i, pf := range piliers.Apport {
			if i >= 2 {
				break
			}
			output += fmt.Sprintf("  • %s\n", pf.Texte)
		}
		output += "\n"
	}

	return output
}
func GenererResumeSynthetique(analyse AnalysePhrases) string {
	if len(analyse.Phrases) == 0 {
		return "Aucun contenu à résumer."
	}

	// ÉTAPE 1: Extraire les phrases principales
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

	// Récupérer les phrases les plus importantes
	phrasesPrincipales := extrairePhrasesImportantes(analyse, catMain)

	// ÉTAPE 2: Classifier chaque phrase selon sa fonction scientifique
	phrasesClassifiees := ClassifierPhrasesDuResume(phrasesPrincipales, analyse)

	// ÉTAPE 3: Structurer selon les 5 piliers
	piliers := StructurerParPiliers(phrasesClassifiees)

	// ÉTAPE 4: Générer le résumé formaté
	resumeStructure := FormaterResumePiliers(piliers)

	// ÉTAPE 5: Ajouter un résumé narratif enrichi
	motsCles := extraireMotsClesUniques(analyse)
	resumeNarratif := genererNarratifPiliers(piliers, motsCles)

	return resumeStructure + "\n" + resumeNarratif
}

// genererNarratifPiliers crée un texte narratif cohérent à partir des piliers
func genererNarratifPiliers(piliers ResumePiliers, motsCles []string) string {
	output := ""
	output += fmt.Sprintf("╔════════════════════════════════════════════════════════════╗\n")
	output += fmt.Sprintf("║  RÉCIT SCIENTIFIQUE CONTINU                               ║\n")
	output += fmt.Sprintf("╚════════════════════════════════════════════════════════════╝\n\n")

	// Construire un narratif fluide qui relie les piliers
	phrases := make([]string, 0)

	// 1. CONTEXTE seul
	if len(piliers.Contexte) > 0 {
		contexte := piliers.Contexte[0].Texte
		contexte = strings.TrimSuffix(contexte, ".")
		phrases = append(phrases, fmt.Sprintf("Le domaine considéré : %s.", strings.ToLower(contexte)))
	}

	// 2. PROBLÈME (indépendant)
	if len(piliers.Probleme) > 0 {
		probleme := piliers.Probleme[0].Texte
		probleme = strings.TrimSuffix(probleme, ".")
		phrases = append(phrases, fmt.Sprintf("Cependant, il existe une limitation fondamentale : %s.", strings.ToLower(probleme)))
	}

	// 3. OBJECTIF (ce qu'on veut faire)
	if len(piliers.Objectif) > 0 {
		objectif := piliers.Objectif[0].Texte
		objectif = strings.TrimSuffix(objectif, ".")
		phrases = append(phrases, fmt.Sprintf("Pour remédier à cette situation, ce travail vise à %s.", strings.ToLower(objectif)))
	}

	// 4. APPROCHE (comment on le fait)
	if len(piliers.Approche) > 0 {
		approche := piliers.Approche[0].Texte
		approche = strings.TrimSuffix(approche, ".")
		phrases = append(phrases, fmt.Sprintf("La stratégie employée consiste à : %s.", strings.ToLower(approche)))
	}

	// 5. APPORT (ce que ça change)
	if len(piliers.Apport) > 0 {
		apport := piliers.Apport[0].Texte
		apport = strings.TrimSuffix(apport, ".")
		phrases = append(phrases, fmt.Sprintf("Il en résultat que : %s.", strings.ToLower(apport)))
	}

	for _, p := range phrases {
		output += p + "\n"
	}

	return output
}

// contexts est une fonction utilitaire
func contexts(s string) string {
	return strings.TrimPrefix(strings.TrimPrefix(s, "Le "), "La ")
}

// extrairePhrasesImportantes récupère les phrases les plus pertinentes (OPTIMISÉ: QuickSort + limite adaptative)
func extrairePhrasesImportantes(analyse AnalysePhrases, catID int) []string {
	phrases := make([]phraseScore, 0, len(analyse.Phrases)/4) // Pré-allocate

	for _, p := range analyse.Phrases {
		if p.CategorieID == catID && len(p.MotsClés) > 0 {
			textNettoye := nettoyerPhraseWiki(p.Texte)
			if strings.TrimSpace(textNettoye) == "" {
				continue
			}
			phrases = append(phrases, phraseScore{
				texte: textNettoye,
				score: p.Score * p.Confiance,
			})
		}
	}

	// QuickSort pour O(n log n) au lieu de O(n²)
	quickSortPhrasesImpl(phrases, 0, len(phrases)-1)

	// OPTIMISÉ: Adapter limite selon nombre total de phrases
	limit := 10 // Réduit de 15 à 10 pour performance texte
	if len(analyse.Phrases) > 200 {
		limit = 15 // Mais garder 15 pour HTML complet
	}
	if len(phrases) < limit {
		limit = len(phrases)
	}

	top := make([]string, 0, limit)
	for i := 0; i < limit; i++ {
		top = append(top, phrases[i].texte)
	}
	return top
}

// extraireMotsClesUniques récupère les mots clés les plus importants (OPTIMISÉ: QuickSort)
func extraireMotsClesUniques(analyse AnalysePhrases) []string {
	freqMots := make(map[string]int)

	for _, phrase := range analyse.Phrases {
		for _, mot := range phrase.MotsClés {
			freqMots[mot]++
		}
	}

	mots := make([]motFreq, 0, len(freqMots)) // Pré-allocate
	for mot, freq := range freqMots {
		mots = append(mots, motFreq{mot, freq})
	}

	// QuickSort pour O(n log n) au lieu de O(n²)
	quickSortMotsImpl(mots, 0, len(mots)-1)

	limit := 15
	if len(mots) < 15 {
		limit = len(mots)
	}
	top15 := make([]string, 0, limit)
	for i := 0; i < limit; i++ {
		top15 = append(top15, mots[i].mot)
	}
	return top15
}

// formulerResume construit un résumé narratif complet basé sur les données réelles
func formulerResume(phrases []string, motsCles []string, catID int, totalPhrases int) string {
	// ⚠️ Filtrer les phrases pour garder seulement les vraies phrases (>=5 mots maintenant)
	var phrasesFiltrees []string
	for _, p := range phrases {
		phraseTrim := strings.TrimSpace(p)
		motsList := strings.Fields(phraseTrim)
		// Garder seulement si la phrase a au moins 5 mots (pas une simple énumération)
		if len(motsList) >= 5 {
			phrasesFiltrees = append(phrasesFiltrees, phraseTrim)
		}
	}

	// Si aucune vraie phrase, essayer avec au moins 3 mots
	if len(phrasesFiltrees) == 0 {
		for _, p := range phrases {
			phraseTrim := strings.TrimSpace(p)
			mots := strings.Fields(phraseTrim)
			if len(mots) >= 3 {
				phrasesFiltrees = append(phrasesFiltrees, phraseTrim)
			}
		}
	}

	// Si toujours rien, utiliser les phrases brutes
	if len(phrasesFiltrees) == 0 {
		phrasesFiltrees = phrases
	}

	if len(phrasesFiltrees) == 0 {
		return ""
	}

	// OPTIMISÉ: Utiliser strings.Builder pour meilleure performance
	var builder strings.Builder
	builder.Grow(len(phrasesFiltrees) * 100) // Pré-allocate mémoire

	// Commencer directement par la première phrase importante
	phraseClé := strings.ToLower(phrasesFiltrees[0])
	builder.WriteString(phraseClé)
	builder.WriteString(". ")

	// Ajouter BEAUCOUP de phrases principales pour un résumé très détaillé
	if len(phrasesFiltrees) > 1 {
		connecteurs := []string{
			"De plus, ",
			"En outre, ",
			"Par ailleurs, ",
			"D'autre part, ",
			"Il est également important de noter que ",
			"Également, ",
			"Ainsi, ",
			"À cet égard, ",
			"Concernant cet aspect, ",
			"Notablement, ",
			"Dans ce contexte, ",
			"Il convient de relever que ",
			"Sur ce point, ",
			"Relativement à cela, ",
			"Sous un autre angle, ",
			"Il faut souligner que ",
			"À ce titre, ",
		}
		// OPTIMISÉ: Réduire phrases intégrées pour texte long
		maxPhrases := 12
		if len(phrasesFiltrees) < 30 {
			maxPhrases = 6 // Texte court = moins de phrases
		}

		for i := 1; i < maxPhrases && i < len(phrasesFiltrees); i++ {
			phraseAdditionnelle := strings.ToLower(phrasesFiltrees[i])
			builder.WriteString(connecteurs[i%len(connecteurs)])
			builder.WriteString(phraseAdditionnelle)
			builder.WriteString(". ")
		}
	}

	// Observation finale basée sur le contenu réel
	if len(phrasesFiltrees) > 1 {
		builder.WriteString(genererObservationDynamique(len(phrasesFiltrees), catID))
	}

	return builder.String()
}

// enrichirResume améliore le résumé initial avec plus de détails clairs et cohérents (PASSE 2)

// enrichirResumeAvance intègre les 3 couches d'analyse avancée
func enrichirResumeAvance(resumePasse1 string, phrasesPrincipales []string, motsCles []string,
	analyse AnalysePhrases, ideeDominante string, phrasesAnalysees []PhraseAnalysee) string {

	if strings.TrimSpace(resumePasse1) == "" {
		return resumePasse1
	}

	var result strings.Builder
	result.Grow(len(resumePasse1) + 2000)

	// ============================================================
	// COUCHE 1: IDÉE DOMINANTE UNIQUE
	// ============================================================
	result.WriteString("\n╔════════════════════════════════════════╗\n")
	result.WriteString("║  IDÉE DOMINANTE\n")
	result.WriteString("╚════════════════════════════════════════╝\n\n")
	result.WriteString(ideeDominante)
	result.WriteString("\n")

	// ============================================================
	// COUCHE 2: FAIT / INTERPRÉTATION / CONSÉQUENCE
	// ============================================================
	result.WriteString("\n╔════════════════════════════════════════╗\n")
	result.WriteString("║  STRUCTURE ANALYTIQUE\n")
	result.WriteString("╚════════════════════════════════════════╝\n\n")

	var faits []string
	var interpretations []string
	var consequences []string

	for _, pa := range phrasesAnalysees {
		switch pa.Type {
		case FAIT:
			if len(faits) < 3 {
				faits = append(faits, pa.Texte)
			}
		case INTERPRETATION:
			if len(interpretations) < 3 {
				interpretations = append(interpretations, pa.Texte)
			}
		case CONSEQUENCE:
			if len(consequences) < 3 {
				consequences = append(consequences, pa.Texte)
			}
		}
	}

	// Afficher les faits
	if len(faits) > 0 {
		result.WriteString("FAITS OBSERVÉS:\n")
		for i, fait := range faits {
			result.WriteString(fmt.Sprintf("  %d. %s\n", i+1, fait))
		}
		result.WriteString("\n")
	}

	// Afficher les interprétations
	if len(interpretations) > 0 {
		result.WriteString("INTERPRÉTATIONS:\n")
		for i, interp := range interpretations {
			result.WriteString(fmt.Sprintf("  %d. %s\n", i+1, interp))
		}
		result.WriteString("\n")
	}

	// Afficher les conséquences
	if len(consequences) > 0 {
		result.WriteString("CONSÉQUENCES:\n")
		for i, cons := range consequences {
			result.WriteString(fmt.Sprintf("  %d. %s\n", i+1, cons))
		}
		result.WriteString("\n")
	}

	// ============================================================
	// COUCHE 3: AXES SÉMANTIQUES
	// ============================================================
	result.WriteString("\n╔════════════════════════════════════════╗\n")
	result.WriteString("║  AXES SÉMANTIQUES\n")
	result.WriteString("╚════════════════════════════════════════╝\n\n")

	// Compter les axes sémantiques
	axesCounts := make(map[AxeSemantique]int)
	for _, pa := range phrasesAnalysees {
		for _, axe := range pa.AxesPrincipaux {
			axesCounts[axe]++
		}
	}

	// Afficher les axes dominants
	type axeCount struct {
		axe   AxeSemantique
		count int
	}
	var sorted []axeCount
	for axe, count := range axesCounts {
		if count > 0 {
			sorted = append(sorted, axeCount{axe, count})
		}
	}

	// Trier par fréquence
	for i := 0; i < len(sorted)-1; i++ {
		for j := 0; j < len(sorted)-1-i; j++ {
			if sorted[j].count < sorted[j+1].count {
				sorted[j], sorted[j+1] = sorted[j+1], sorted[j]
			}
		}
	}

	for i, ac := range sorted {
		if i < 4 { // Top 4 axes
			result.WriteString(fmt.Sprintf("• %s (%d phrases)\n", ac.axe, ac.count))
		}
	}

	// ============================================================
	// SECTIONS CLASSIQUES (simplifiées)
	// ============================================================
	result.WriteString("\n╔════════════════════════════════════════╗\n")
	result.WriteString("║  RÉSUMÉ NARRATIF\n")
	result.WriteString("╚════════════════════════════════════════╝\n\n")
	result.WriteString(resumePasse1)

	return result.String()
}

// enrichirResume intègre les 3 couches d'analyse avancée
func enrichirResume(resumePasse1 string, phrasesPrincipales []string, motsCles []string, analyse AnalysePhrases) string {
	if strings.TrimSpace(resumePasse1) == "" {
		return resumePasse1
	}

	// OPTIMISÉ: Utiliser strings.Builder
	var result strings.Builder
	result.Grow(len(resumePasse1) + 1000)
	result.WriteString(resumePasse1)

	// SECTION 1: Ajouter les détails supplémentaires restants (LARGEMENT ÉTENDU)
	if len(phrasesPrincipales) > 3 {
		phrasesRestantes := phrasesPrincipales[3:]
		if len(phrasesRestantes) > 0 {
			result.WriteString("\n\nDimensions complémentaires détaillées: ")
			transitionsEnrichissement := []string{
				"D'une part, ",
				"D'autre part, ",
				"Aussi, ",
				"Enfin, ",
				"Par ailleurs, ",
				"En complément, ",
				"Sous cet angle, ",
				"À cet égard, ",
				"Il convient d'ajouter que ",
				"De manière significative, ",
			}

			// Ajouter les phrases restantes de manière très complète (7 au lieu de 3)
			for i, phrase := range phrasesRestantes {
				if i < 7 { // Augmenté de 3 à 7 phrases supplémentaires
					mots := strings.Fields(phrase)
					if len(mots) >= 4 { // Réduit de 5 à 4 pour plus de complétude
						phraseEnrichie := strings.ToLower(phrase)
						result.WriteString(transitionsEnrichissement[i%len(transitionsEnrichissement)])
						result.WriteString(phraseEnrichie)
						result.WriteString(". ")
					}
				}
			}
		}
	}

	// SECTION 2: Mots-clés et concepts principaux (ÉTENDU)
	if len(motsCles) > 0 {
		result.WriteString("\n\nConcepts clés et thèmes majeurs: ")
		maxMots := 10 // Augmenté de 5 à 10
		if len(motsCles) < 10 {
			maxMots = len(motsCles)
		}
		for i := 0; i < maxMots; i++ {
			if i > 0 && i < maxMots-1 {
				result.WriteString(", ")
			} else if i == maxMots-1 && i > 0 {
				result.WriteString(" et ")
			}
			result.WriteString(motsCles[i])
		}
		result.WriteString(".")
	}

	// SECTION 3: Analyse détaillée des relations (NOUVELLE)
	result.WriteString("\n\nRelations et interconnexions: Cette analyse révèle comment les concepts ")
	result.WriteString("s'articulent entre eux, formant un réseau complexe de dépendances et d'influences mutuelles. ")
	result.WriteString("Les éléments identifiés ne sont pas isolés mais constituent plutôt ")
	result.WriteString("les éléments d'une structure cohésive plus large.")

	// SECTION 4: Synthèse conclusive (ÉTENDUE)
	result.WriteString("\n\nSynthèse conclusive: Cet analyse synthétique consolide une compréhension holistique du sujet, ")
	result.WriteString("révélant les interdépendances entre les différentes dimensions du contenu. ")
	result.WriteString("Les concepts majeurs s'articulent pour former un ensemble complet et nuancé, ")
	result.WriteString("permettant une appréhension profonde des enjeux et des dynamiques en cause.")

	return result.String()
}

// ContientPhrasesReelles vérifie si le texte contient des phrases vraies et cohérentes
// (pas juste une liste de mots isolés)
func ContientPhrasesReelles(texte string) bool {
	// Rejeter si le texte est vide
	if strings.TrimSpace(texte) == "" {
		return false
	}

	// Rejeter si des mots de corruption typiques présents
	motsCorruption := []string{
		"modifier",         // Signe de wiki/article editing
		"article détaillé", // Signe de mauvaise extraction
		"catégorie",        // Signe de structure wiki
		",   ",             // Espaces anormaux
	}

	texteClean := strings.ToLower(texte)
	for _, motCorrup := range motsCorruption {
		if strings.Count(texteClean, strings.ToLower(motCorrup)) > 2 {
			return false // Trop de corruption détectée
		}
	}

	phrases := strings.Split(texte, ".")
	phrasesReelles := 0

	verbsCommuns := []string{
		"être", "avoir", "faire", "aller", "pouvoir", "devoir", "vouloir",
		"dire", "donner", "prendre", "venir", "voir", "savoir", "croire",
		"trouver", "montrer", "demander", "obtenir", "produire", "mettre",
		"penser", "utiliser", "laisser", "jouer", "joindre", "changer",
		"est", "sont", "a", "ont", "fait", "vont", "peut", "doit", "veut",
		"dit", "donne", "prend", "vient", "voit", "sait", "croit", "trouve",
		"montre", "demande", "obtient", "produit", "met", "pense", "utilise",
		"laisse", "joue", "joint", "change", "contient", "représente", "indique",
		"montre", "révèle", "propose", "analyse", "explore", "traite", "aborde",
		"examine", "considère", "affirme", "soutient", "appuie", "confirme",
		"démontre", "établit", "explicite", "implique", "suppose", "constate",
	}

	for _, phrase := range phrases {
		phraseTrim := strings.TrimSpace(phrase)
		if phraseTrim == "" {
			continue
		}

		mots := strings.Fields(phraseTrim)
		// Une vraie phrase doit avoir au minimum 8 mots
		if len(mots) < 8 {
			continue
		}

		// Rejeter si trop de singletons (liste de mots isolés)
		// Exemple: "les premiers sont ibm , microsoft , computer" = trop de virgules
		virgulesCount := strings.Count(phraseTrim, ",")
		if virgulesCount > len(mots)/3 {
			// Trop de listes = pas une phrase réelle
			continue
		}

		// Chercher si la phrase contient un verbe
		motLower := strings.ToLower(phraseTrim)
		contientVerbe := false
		for _, verbe := range verbsCommuns {
			if strings.Contains(motLower, verbe) {
				contientVerbe = true
				break
			}
		}

		if contientVerbe {
			phrasesReelles++
		}
	}

	// Besoin d'au minimum 2 phrases réelles avec verbe (pas juste 1)
	return phrasesReelles >= 2
}

// genererObservationDynamique crée une observation basée sur le nombre de phrases
func genererObservationDynamique(nbPhrases int, catID int) string {
	// Adaptative: plus de phrases = plus détaillé
	if nbPhrases < 2 {
		return "Ce contenu est bref mais significatif."
	}
	if nbPhrases < 5 {
		return "L'analyse révèle plusieurs dimensions du sujet traité."
	}
	if nbPhrases < 10 {
		return "La profondeur du contenu montre une exploration complète du thème."
	}

	// Pour 10+ phrases, générer une observation dynamique
	observations := []string{
		"L'ampleur du contenu démontre une couverture exhaustive du sujet.",
		"La richesse des éléments présentés révèle une perspective multidimensionnelle.",
		"L'ensemble des informations s'articule autour d'une vision cohérente et structurée.",
		"La multiplicité des angles explorés indique une compréhension profonde et nuancée.",
		"La densité des contenus traités suggère une expertise ou une connaissance approfondie.",
	}
	return observations[nbPhrases%len(observations)]
}

// nettoyerPhraseWiki supprime les artefacts wiki/non-pertinents d'une phrase (OPTIMISÉ)
func nettoyerPhraseWiki(phrase string) string {
	resultat := strings.ToLower(phrase)

	// OPTIMISÉ: Faire un seul pass avec strings.Builder au lieu de 17+ ReplaceAll
	var builder strings.Builder
	builder.Grow(len(resultat))

	skipPatterns := map[string]bool{
		"↑": true, "https": true, "http": true, "www.": true,
		"le code": true, "modifier": true, "article détaillé": true,
		"article connexe": true, "voir aussi": true, "catégorie": true,
		"wikip": true, "lien externe": true, "liens externes": true,
		"portail:": true, "références": true, "source:": true,
		"lire en ligne": true,
	}

	words := strings.Fields(resultat)
	for i, word := range words {
		// Sauter si le mot est dans la liste de suppression
		if skipPatterns[word] {
			continue
		}

		// Ajouter mot avec espaces
		if i > 0 {
			builder.WriteString(" ")
		}
		builder.WriteString(word)
	}

	resultat = builder.String()

	// Nettoyer les ponctuations mal placées (minimal)
	resultat = strings.ReplaceAll(resultat, "  ,", ",")
	resultat = strings.ReplaceAll(resultat, " .", ".")
	resultat = strings.ReplaceAll(resultat, "  ", " ")

	return strings.TrimSpace(resultat)
	resultat = strings.ReplaceAll(resultat, " ,", ",")

	// Nettoyer les caractères restants
	resultat = strings.TrimSpace(resultat)

	// Rejeter si trop court ou que du bruit
	if len(resultat) < 15 || strings.Contains(resultat, "↑") || strings.Contains(resultat, "http") {
		return ""
	}

	return resultat
}

// GenererReponseAvancee crée une réponse avec context
func GenererReponseAvancee(catID int, motsClés []string, texteOriginal string) string {
	base := GenererReponse(catID, motsClés)

	// Adapter le ton selon la forme de la question
	intro := ""
	question := strings.ToLower(strings.TrimSpace(texteOriginal))
	switch {
	case strings.Contains(question, "pourquoi"):
		intro = "Explication rapide : "
	case strings.Contains(question, "comment"):
		intro = "Voici une manière simple d'aborder cela : "
	case strings.Contains(question, "quand"):
		intro = "Repère chronologique : "
	case strings.Contains(question, "où"):
		intro = "Contexte / localisation : "
	case strings.Contains(question, "qui"):
		intro = "Concernant l'acteur principal : "
	case strings.Contains(question, "quoi") || strings.Contains(question, "qu'est-ce"):
		intro = "Définition concise : "
	case strings.Contains(question, "combien"):
		intro = "Ordre de grandeur : "
	}

	response := base
	if intro != "" {
		response = intro + base
	}

	// Ajouter contexte basé sur les mots clés
	if len(motsClés) > 0 {
		limit := 3
		if len(motsClés) < 3 {
			limit = len(motsClés)
		}
		response = fmt.Sprintf("%s\nPoints clés: %s", response, strings.Join(motsClés[:limit], ", "))
	}

	return response
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

// VerifierQualiteResume évalue la qualité et pertinence du résumé généré
func VerifierQualiteResume(resume string, analyse AnalysePhrases, phrasesCandidates []string) float64 {
	score := 0.0
	maxScore := 100.0

	// 🔴 REJET IMMÉDIAT 1: Si le résumé ne contient que des mots sans phrases réelles
	if !ContientPhrasesReelles(resume) {
		return 0.0 // Score 0 = Rejeter immédiatement
	}

	// 🔴 REJET IMMÉDIAT 2: Vérifier que la première phrase est réelle
	// (pas une liste de mots comme "les premiers sont ibm , microsoft...")
	premierePhrase := strings.Split(resume, ".")[0]
	premierePhraseMots := strings.Fields(strings.TrimSpace(premierePhrase))
	virgulesDansPremiere := strings.Count(premierePhrase, ",")

	// Une phrase avec plus de 3 virgules ET moins de 20 mots = probablement une liste
	if virgulesDansPremiere > 2 && len(premierePhraseMots) < 20 {
		return 0.0
	}

	// Ou si elle commence par "les premiers", "sont", "IBM", etc (énumération)
	if len(premierePhraseMots) > 0 {
		premierMot := strings.ToLower(premierePhraseMots[0])
		patternsList := []string{"les premiers", "les", "sont", "ibm", "microsoft", "computer", "digital"}
		deuxiemeMot := ""
		if len(premierePhraseMots) > 1 {
			deuxiemeMot = strings.ToLower(premierePhraseMots[1])
		}

		for _, pattern := range patternsList {
			if pattern == premierMot || pattern == deuxiemeMot {
				// Si start with énumération + virgules = liste = bad
				if virgulesDansPremiere > 1 {
					return 0.0
				}
			}
		}
	}

	// ✅ CRITÈRE 1: Longueur minimum (20 mots minimum)
	mots := strings.Fields(resume)
	if len(mots) >= 20 {
		score += 20.0
	} else if len(mots) >= 10 {
		score += 10.0
	}

	// Préparer resumeLower pour les critères suivants
	resumeLower := strings.ToLower(resume)

	// ✅ CRITÈRE 2: Nombre de phrases minimum (au moins 3 phrases réelles)
	phrases := strings.Split(resume, ".")
	nbPhrasesResume := 0
	phrasesMoyLongueur := 0
	for _, p := range phrases {
		if strings.TrimSpace(p) != "" {
			nbPhrasesResume++
			phrasesMoyLongueur += len(strings.Fields(p))
		}
	}
	if nbPhrasesResume >= 3 {
		score += 15.0
	} else if nbPhrasesResume >= 2 {
		score += 7.5
	}

	// Vérifier que les phrases ont une longueur minimale (5+ mots par phrase)
	if nbPhrasesResume > 0 {
		longueurMoyenne := phrasesMoyLongueur / nbPhrasesResume
		if longueurMoyenne < 5 {
			// Les phrases sont trop courtes = résumé fragile
			return 0.0
		}
	}

	// ✅ CRITÈRE 3bis: Vérifier absence de mots de corruption
	// (Cela rejette les résumés contaminés par des données wiki/malformées)
	motsCorruptionGraves := []string{
		"modifier",
		"article détaillé",
		"catégorie",
		"wikip",
	}
	nbMotsCorruption := 0
	for _, motCorrup := range motsCorruptionGraves {
		nbMotsCorruption += strings.Count(resumeLower, motCorrup)
	}
	if nbMotsCorruption > 1 {
		// Trop de corruption = résumé invalide
		return 0.0
	}

	// ✅ CRITÈRE 3: Densité conceptuelle (pas de phrases génériques)

	// Phrases génériques à éviter
	phrasesVides := []string{
		"c'est important",
		"il est important",
		"selon les experts",
		"les études montrent",
		"il apparaît que",
		"il convient de noter",
		"à cet égard",
	}

	contientGenerique := false
	for _, generic := range phrasesVides {
		if strings.Contains(resumeLower, generic) {
			contientGenerique = true
			break
		}
	}

	// Rejeter si trop d'espaces (signe de mauvaise structure)
	espacesMultiples := strings.Count(resume, "   ")
	if espacesMultiples > 3 {
		contientGenerique = true
	}

	if !contientGenerique {
		score += 20.0
	} else {
		score += 5.0
	}

	// ✅ CRITÈRE 4: Présence de mots clés du contenu original
	motsClesAnalyse := extraireMotsClesUniques(analyse)
	motsClesMatches := 0
	for _, mot := range motsClesAnalyse {
		if strings.Contains(resumeLower, strings.ToLower(mot)) {
			motsClesMatches++
		}
	}

	if len(motsClesAnalyse) > 0 {
		pourcentageMatch := float64(motsClesMatches) / float64(len(motsClesAnalyse))
		score += pourcentageMatch * 20.0
	}

	// ✅ CRITÈRE 5: Diversité (variation des connecteurs)
	connecteurs := []string{
		"en outre,",
		"de plus,",
		"par ailleurs,",
		"d'autre part,",
		"il est également",
		"la profondeur",
		"l'ampleur",
	}

	connecteursPresents := 0
	for _, conn := range connecteurs {
		if strings.Contains(resumeLower, conn) {
			connecteursPresents++
		}
	}

	if connecteursPresents >= 2 {
		score += 15.0
	} else if connecteursPresents >= 1 {
		score += 7.5
	}

	// ✅ CRITÈRE 6: Pas de répétitions excessives
	mots = strings.Fields(strings.ToLower(resume))
	freqMots := make(map[string]int)
	maxFreq := 0
	for _, mot := range mots {
		freqMots[mot]++
		if freqMots[mot] > maxFreq {
			maxFreq = freqMots[mot]
		}
	}

	// Si un mot se répète trop (max 30% du résumé = mauvais)
	if maxFreq < len(mots)/3 {
		score += 10.0
	}

	// ✅ CRITÈRE 7: Couvre au moins 2 phrases du source
	if len(phrasesCandidates) >= 2 {
		score += 5.0
	}

	// Normaliser le score en pourcentage (0.0 - 1.0)
	finalScore := score / maxScore
	if finalScore > 1.0 {
		finalScore = 1.0
	}

	return finalScore
}

// quickSortPhrasesImpl effectue un QuickSort sur les phrases (O(n log n) vs O(n²) pour bubble sort)
func quickSortPhrasesImpl(phrases []phraseScore, low, high int) {
	if low < high {
		pi := partitionPhrasesImpl(phrases, low, high)
		quickSortPhrasesImpl(phrases, low, pi-1)
		quickSortPhrasesImpl(phrases, pi+1, high)
	}
}

func partitionPhrasesImpl(phrases []phraseScore, low, high int) int {
	pivot := phrases[high].score
	i := low - 1

	for j := low; j < high; j++ {
		if phrases[j].score > pivot { // Descending order
			i++
			phrases[i], phrases[j] = phrases[j], phrases[i]
		}
	}
	phrases[i+1], phrases[high] = phrases[high], phrases[i+1]
	return i + 1
}

// quickSortMotsImpl effectue un QuickSort sur les mots-clés (O(n log n) vs O(n²) pour bubble sort)
func quickSortMotsImpl(mots []motFreq, low, high int) {
	if low < high {
		pi := partitionMotsImpl(mots, low, high)
		quickSortMotsImpl(mots, low, pi-1)
		quickSortMotsImpl(mots, pi+1, high)
	}
}

func partitionMotsImpl(mots []motFreq, low, high int) int {
	pivot := mots[high].freq
	i := low - 1

	for j := low; j < high; j++ {
		if mots[j].freq > pivot { // Descending order
			i++
			mots[i], mots[j] = mots[j], mots[i]
		}
	}
	mots[i+1], mots[high] = mots[high], mots[i+1]
	return i + 1
}

// ============================================================
// COUCHE 1: IDÉE DOMINANTE UNIQUE
// ============================================================

// IdentifierIdeeDominante extrait la phrase-clé synthétisant le cœur de l'article
func IdentifierIdeeDominante(phrases []string, motsCles []string) string {
	if len(phrases) == 0 {
		return ""
	}

	// Chercher la phrase contenant le plus de mots-clés majeurs
	maxMotsCles := 0
	var phraseDominante string

	for _, phrase := range phrases {
		phraseLower := strings.ToLower(phrase)
		compteur := 0

		// Compter les mots-clés dans cette phrase
		for i := 0; i < len(motsCles) && i < 5; i++ {
			if strings.Contains(phraseLower, strings.ToLower(motsCles[i])) {
				compteur++
			}
		}

		if compteur > maxMotsCles {
			maxMotsCles = compteur
			phraseDominante = phrase
		}
	}

	if phraseDominante == "" && len(phrases) > 0 {
		phraseDominante = phrases[0] // Fallback: première phrase
	}

	return fmt.Sprintf("Le cœur de cet article est que : %s", strings.ToLower(phraseDominante))
}

// ============================================================
// COUCHE 2: CLASSIFICATION FAIT / INTERPRÉTATION / CONSÉQUENCE
// ============================================================

// ClassifierPhrase détermine si une phrase est un Fait, Interprétation ou Conséquence
func ClassifierPhrase(phrase string) TypePhrase {
	lower := strings.ToLower(phrase)

	// Patterns FAIT: énoncés objectifs, observations
	patternsHorizontalFait := []string{
		"est ", "sont ", "a ", "ont ", "existe", "se trouve",
		"mesure ", "contient", "comprend", "représente",
		"selon ", "d'après ", "indique", "montre", "révèle",
		"en ", "pendant ", "lors de ", "au cours de",
		"constitue ", "forme ", "compose ", "structure",
	}

	// Patterns INTERPRÉTATION: jugements, sens, interprétation
	patternsInterpretation := []string{
		"signifie ", "implique ", "suggère", "révèle",
		"reflète ", "traduit ", "exprime", "démontre",
		"semble ", "paraît ", "apparaît", "considéré",
		"constitue ", "représente une", "est une forme",
		"bien que", "tandis que", "alors que",
		"en réalité", "en vérité", "vraiment", "effectivement",
	}

	// Patterns CONSÉQUENCE: résultats, impacts, suites
	patternsConsequence := []string{
		"conduit à", "mène à", "engendre", "provoque", "cause",
		"résulte en", "aboutit", "génère", "crée", "entraîne",
		"pour cette raison", "c'est pourquoi", "donc", "ainsi",
		"en conséquence", "par conséquent", "du coup", "finalement",
		"cela explique", "cela démontre", "ce qui", "ce qui en résulte",
	}

	for _, pattern := range patternsConsequence {
		if strings.Contains(lower, pattern) {
			return CONSEQUENCE
		}
	}

	for _, pattern := range patternsInterpretation {
		if strings.Contains(lower, pattern) {
			return INTERPRETATION
		}
	}

	for _, pattern := range patternsHorizontalFait {
		if strings.Contains(lower, pattern) {
			return FAIT
		}
	}

	return NEUTRE
}

// ============================================================
// COUCHE 3: AXES SÉMANTIQUES
// ============================================================

// IdentifierAxesSemantiques détecte les axes sémantiques d'une phrase
func IdentifierAxesSemantiques(phrase string) []AxeSemantique {
	lower := strings.ToLower(phrase)
	var axes []AxeSemantique

	// Mots-clés par axe sémantique
	motsCles := map[AxeSemantique][]string{
		EMOTIONNEL: {"culpabilité", "déception", "trahison", "honte", "peur", "colère", "amour", "haine", "joie", "tristesse", "angoisse", "confiance", "doute", "regret", "jalousie", "fierté", "humiliation", "soulagement", "frustration"},

		EDUCATIF: {"mensonge", "vérité", "croyance", "apprentissage", "enseignement", "connaissance", "ignorance", "savoir", "comprendre", "expliquer", "éducation", "formation", "école", "maître", "discipline", "correction", "instruction", "leçon"},

		PSYCHOLOGIQUE: {"construction du réel", "perception", "représentation", "inconscient", "conscient", "mémoire", "trauma", "défense", "mécanisme", "psyché", "esprit", "mental", "cognitif", "pensée", "emotion", "sentiment", "affect", "fusion", "identification"},

		RELATIONNEL: {"parent", "enfant", "relation", "lien", "rupture", "séparation", "union", "distance", "proximité", "intimité", "communication", "dialogue", "silence", "secret", "confiance", "betrayal", "dépendance", "autonomie"},

		MORAL: {"bien", "mal", "juste", "injuste", "droit", "devoir", "responsabilité", "culpabilité", "innocence", "pardon", "vengeance", "justice", "éthique", "moralité", "conscience", "valeur", "principe"},

		PRAGMATIQUE: {"solution", "problème", "résoudre", "gérer", "action", "faire", "réaliser", "conséquence", "impact", "effet", "résultat", "outcome", "pratique", "efficacité", "utilité", "efficace"},

		COGNITIF: {"comprendre", "savoir", "apprendre", "logique", "raison", "analyse", "synthèse", "abstraction", "concept", "idée", "théorie", "système", "structure", "pattern", "règle", "loi"},
	}

	// Parcourir chaque axe et compter les occurrences
	scores := make(map[AxeSemantique]int)
	for axe, mots := range motsCles {
		for _, mot := range mots {
			if strings.Contains(lower, mot) {
				scores[axe]++
			}
		}
	}

	// Trier les axes par score (keep only top 3)
	type axeScore struct {
		axe   AxeSemantique
		score int
	}
	var sorted []axeScore
	for axe, score := range scores {
		if score > 0 {
			sorted = append(sorted, axeScore{axe, score})
		}
	}

	// Bubble sort simple par score (descending)
	for i := 0; i < len(sorted)-1; i++ {
		for j := 0; j < len(sorted)-1-i; j++ {
			if sorted[j].score < sorted[j+1].score {
				sorted[j], sorted[j+1] = sorted[j+1], sorted[j]
			}
		}
	}

	// Ajouter les top 3 axes
	limit := 3
	if len(sorted) < 3 {
		limit = len(sorted)
	}
	for i := 0; i < limit; i++ {
		axes = append(axes, sorted[i].axe)
	}

	return axes
}

// AnalyserPhrasesAvancee retourne les phrases analysées avec types et axes
func AnalyserPhrasesAvancee(phrases []string) []PhraseAnalysee {
	var analysees []PhraseAnalysee

	for _, phrase := range phrases {
		if strings.TrimSpace(phrase) == "" {
			continue
		}

		analysees = append(analysees, PhraseAnalysee{
			Texte:          phrase,
			Type:           ClassifierPhrase(phrase),
			AxesPrincipaux: IdentifierAxesSemantiques(phrase),
		})
	}

	return analysees
}

// FormatTypePhrase convertit le type en string lisible
func FormatTypePhrase(t TypePhrase) string {
	switch t {
	case FAIT:
		return "[FAIT]"
	case INTERPRETATION:
		return "[INTERPRÉTATION]"
	case CONSEQUENCE:
		return "[CONSÉQUENCE]"
	default:
		return "[NEUTRE]"
	}
}
