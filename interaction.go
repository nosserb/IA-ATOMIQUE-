package main

import (
	"IA-ATOMIQUE/database"
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// DetecterLangue identifie la langue d'un texte de manière robuste
func DetecterLangue(texte string) string {
	lower := strings.ToLower(texte)

	// Si trop court, on ne peut pas décider
	if len(strings.Fields(texte)) < 2 {
		return "unknown"
	}

	// Compter les caractères accentués français
	accentsFR := strings.Count(lower, "é") + strings.Count(lower, "è") + strings.Count(lower, "ê") +
		strings.Count(lower, "à") + strings.Count(lower, "ù") + strings.Count(lower, "û") +
		strings.Count(lower, "ô") + strings.Count(lower, "ç") + strings.Count(lower, "œ") +
		strings.Count(lower, "â") + strings.Count(lower, "ï") + strings.Count(lower, "î")

	// Caractères allemands
	accentDE := strings.Count(lower, "ä") + strings.Count(lower, "ö") + strings.Count(lower, "ü") + strings.Count(lower, "ß")

	// Mots fréquents français (stopwords)
	motsFR := []string{" le ", " la ", " de ", " et ", " que ", " du ", " un ", " une ", " en ", " a ",
		" les ", " des ", " pour ", " par ", " avec ", " sur ", " ce ", " est ", " au ", " pas ",
		" ou ", " qui ", " ne ", " se ", " plus ", " dans ", " lui ", " son ", " elle ", " nous ",
		" vous ", " s ", " cet ", " cette ", " qu ", " l ", " d "}
	comptFR := 0
	for _, mot := range motsFR {
		if strings.Contains(" "+lower+" ", mot) {
			comptFR += 2
		}
	}

	// Mots typiquement allemands
	motsDE := []string{" der ", " die ", " das ", " und ", " zu ", " in ", " von ", " den ", " mit ",
		" für ", " ist ", " ein ", " eine ", " sich ", " nicht ", " worden ", " hat ", " kann "}
	comptDE := 0
	for _, mot := range motsDE {
		if strings.Contains(" "+lower+" ", mot) {
			comptDE += 2
		}
	}

	// Mots typiquement anglais
	motsEN := []string{" the ", " and ", " to ", " of ", " a ", " in ", " is ", " you ", " that ",
		" he ", " was ", " for ", " on ", " are ", " be ", " have ", " this ", " but ", " his "}
	comptEN := 0
	for _, mot := range motsEN {
		if strings.Contains(" "+lower+" ", mot) {
			comptEN += 2
		}
	}

	// Mots typiquement espagnols
	motsES := []string{" el ", " la ", " de ", " que ", " y ", " a ", " en ", " un ", " ser ", " se ",
		" no ", " por ", " con ", " una ", " está ", " han "}
	comptES := 0
	for _, mot := range motsES {
		if strings.Contains(" "+lower+" ", mot) {
			comptES += 2
		}
	}

	// Patterns caractéristiques
	if accentDE > 0 {
		comptDE += accentDE * 5 // Les accents allemands sont très distinctifs
	}

	// Décider basé sur les scores
	scores := map[string]int{
		"fr": comptFR + (accentsFR * 3),
		"de": comptDE,
		"en": comptEN,
		"es": comptES,
	}

	maxScore := 0
	var topLang string
	var secondScore int

	// Trouver les deux meilleurs scores
	for lang, score := range scores {
		if score > maxScore {
			secondScore = maxScore
			maxScore = score
			topLang = lang
		}
	}

	// Si pas assez de confiance, marquer comme "unknown"
	if maxScore < 3 {
		return "unknown"
	}

	// Si le français est très proche de l'autre langue, douter
	frScore := scores["fr"]
	if topLang != "fr" && frScore > 0 && (maxScore-secondScore) < 5 {
		return "mixed"
	}

	return topLang
}

// FiltrerParLangue sépare le texte par langue et garde uniquement le français
func FiltrerParLangue(texte string) (string, map[string][]string) {
	phrases := strings.Split(texte, ".")
	var texteFR strings.Builder
	langues := make(map[string][]string)

	for _, phrase := range phrases {
		phrase = strings.TrimSpace(phrase)
		if len(phrase) < 5 {
			continue
		}

		langue := DetecterLangue(phrase)

		// Enregistrer chaque phrase détectée
		langues[langue] = append(langues[langue], phrase)

		// Ajouter UNIQUEMENT le français au texte filtré
		if langue == "fr" {
			texteFR.WriteString(phrase + ". ")
		}
	}

	return texteFR.String(), langues
}

// ExtraireArticles sépare le HTML en articles distincts
func ExtraireArticles(html string) []string {
	// D'abord nettoyer les éléments inutiles
	scriptRegex := regexp.MustCompile(`(?is)<script[^>]*>.*?</script>`)
	html = scriptRegex.ReplaceAllString(html, " ")

	styleRegex := regexp.MustCompile(`(?is)<style[^>]*>.*?</style>`)
	html = styleRegex.ReplaceAllString(html, " ")

	navRegex := regexp.MustCompile(`(?is)<nav[^>]*>.*?</nav>`)
	html = navRegex.ReplaceAllString(html, " ")

	footerRegex := regexp.MustCompile(`(?is)<footer[^>]*>.*?</footer>`)
	html = footerRegex.ReplaceAllString(html, " ")

	asideRegex := regexp.MustCompile(`(?is)<aside[^>]*>.*?</aside>`)
	html = asideRegex.ReplaceAllString(html, " ")

	headerRegex := regexp.MustCompile(`(?is)<header[^>]*>.*?</header>`)
	html = headerRegex.ReplaceAllString(html, " ")

	// Extraire les articles
	articleRegex := regexp.MustCompile(`(?is)<article[^>]*>(.*?)</article>`)
	matches := articleRegex.FindAllStringSubmatch(html, -1)
	if len(matches) > 0 {
		var articles []string
		for _, match := range matches {
			if len(match) > 1 {
				articles = append(articles, match[1])
			}
		}
		if len(articles) > 0 {
			return articles
		}
	}

	// Fallback: chercher <main>
	mainRegex := regexp.MustCompile(`(?is)<main[^>]*>(.*?)</main>`)
	mainMatches := mainRegex.FindStringSubmatch(html)
	if len(mainMatches) > 1 && strings.TrimSpace(mainMatches[1]) != "" {
		return []string{mainMatches[1]}
	}

	// Fallback: chercher <div class="content">
	contentRegex := regexp.MustCompile(`(?is)<div[^>]*class=["\'].*?content["\'][^>]*>(.*?)</div>`)
	contentMatches := contentRegex.FindStringSubmatch(html)
	if len(contentMatches) > 1 && strings.TrimSpace(contentMatches[1]) != "" {
		return []string{contentMatches[1]}
	}

	// Dernière chance: le HTML entier
	return []string{html}
}

// NettoyerArticle nettoie un article HTML en texte pur (utilise la nouvelle extraction avancée)
func NettoyerArticle(html string) string {
	// ÉTAPES 1-3: Utiliser la nouvelle fonction d'extraction améliorée
	contenuPrincipal := database.ExtraireContenuHTMLPrincipal(html)

	// Si le contenu est vide, fallback à l'ancienne méthode
	if strings.TrimSpace(contenuPrincipal) == "" {
		// Fallback: ancienne méthode simple
		tagRegex := regexp.MustCompile(`<[^>]+>`)
		texte := tagRegex.ReplaceAllString(html, " ")

		// Enlever les entités HTML
		texte = strings.ReplaceAll(texte, "&nbsp;", " ")
		texte = strings.ReplaceAll(texte, "&lt;", "<")
		texte = strings.ReplaceAll(texte, "&gt;", ">")
		texte = strings.ReplaceAll(texte, "&amp;", "&")
		texte = strings.ReplaceAll(texte, "&quot;", "\"")
		texte = strings.ReplaceAll(texte, "&apos;", "'")

		// Espaces multiples
		spaceRegex := regexp.MustCompile(`\s+`)
		texte = spaceRegex.ReplaceAllString(texte, " ")

		return strings.TrimSpace(texte)
	}

	// ÉTAPE 4: Segmenter en phrases (déjà fait par ExtraireContenuHTMLPrincipal)
	// ÉTAPE 5: Filtrer le bruit avec la blacklist existante
	result := FiltrerLignesRepetitives(contenuPrincipal)

	return result
}

// FiltrerLignesRepetitives supprime les séquences répétitives et inutiles
func FiltrerLignesRepetitives(texte string) string {
	lines := strings.Split(texte, ".")
	var result []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) < 3 {
			continue
		}

		// Ignorer si trop de ponctuation
		punctCount := strings.Count(line, ",") + strings.Count(line, ";") + strings.Count(line, "-")
		if punctCount > len(line)/3 {
			continue
		}

		// Ignorer patterns inutiles
		if DansBlocageAvance(line) {
			continue
		}

		// [NOUVEAU] Filtrer les phrases non-conceptuelles (disclaimers, génériques)
		if !EstPhrasConceptuelle(line) {
			continue
		}

		result = append(result, line)
	}

	return strings.Join(result, ". ") + "."
}

// DansBlocageAvance vérifie si une ligne doit être ignorée
func DansBlocageAvance(texte string) bool {
	lower := strings.ToLower(texte)

	// Patterns à ignorer - Disclaimers légaux
	blocages := []string{
		"voir tous les commentaires",
		"journal du jour",
		"lire la suite",
		"afficher plus",
		"charger plus",
		"nouveau commentaire",
		"partager",
		"réagir",
		"publié",
		"mis à jour",
		"tous droits réservés",
		"copyright",
		"mentions légales",
		"politique de confidentialité",
		"conditions d'utilisation",
		"inscription",
		"connexion",
		"mon compte",
	}

	for _, blocage := range blocages {
		if strings.Contains(lower, blocage) {
			return true
		}
	}

	// Ignorer si c'est juste des dates/numéros
	if regexp.MustCompile(`^\d+[/-]\d+`).MatchString(lower) {
		return true
	}

	return false
}

// FiltrerPhrasesNonConceptuelles élimine les phrases sans contenu conceptuel
func FiltrerPhrasesNonConceptuelles(texte string) string {
	phrases := strings.Split(texte, ".")
	var resultat []string

	for _, phrase := range phrases {
		phrase = strings.TrimSpace(phrase)
		if len(phrase) < 5 {
			continue
		}

		// Vérifier si la phrase est conceptuelle
		if !EstPhrasConceptuelle(phrase) {
			continue
		}

		resultat = append(resultat, phrase)
	}

	if len(resultat) > 0 {
		return strings.Join(resultat, ". ") + "."
	}
	return texte
}

// EstPhrasConceptuelle évalue si une phrase a du contenu unique et conceptuel
func EstPhrasConceptuelle(phrase string) bool {
	lower := strings.ToLower(phrase)

	// 1. Vérifier les disclaimers/avertissements légaux
	disclaimers := []string{
		"résultats peuvent varier",
		"ne constitue pas",
		"avis médical",
		"consultation professionnelle",
		"consulter un médecin",
		"consulter un docteur",
		"cliquez ici",
		"en savoir plus",
		"témoignages peuvent être fictionnels",
		"fictionnels",
		"personnes et noms fictifs",
		"toute ressemblance",
		"pur hasard",
		"droits réservés",
		"copyright",
		"mentions légales",
		"politique de confidentialité",
		"tous les droits",
		"reproduction interdite",
		"utilisation personnelle",
		"conditions d'utilisation",
		"prix nobel",
		"lauréat nobel",
		"découverte révolutionnaire",
	}

	for _, disclaimer := range disclaimers {
		if strings.Contains(lower, disclaimer) {
			return false
		}
	}

	// 2. Vérifier les patterns de phrases génériques/vides
	phrasesVides := []string{
		"les études montrent",
		"les recherches montrent",
		"des recherches récentes",
		"selon les études",
		"selon les experts",
		"les experts affirment",
		"il a été découvert",
		"on sait que",
		"il est connu que",
		"il y a",
		"c'est important",
		"c'est essentiel",
		"remarquable",
		"intéressant",
		"fascinant",
	}

	for _, pattern := range phrasesVides {
		if strings.Contains(lower, pattern) {
			// Vérifier que la phrase a au moins 15 mots (contenu après le pattern)
			nbMots := len(strings.Fields(phrase))
			if nbMots < 15 {
				return false
			}
		}
	}

	// 3. Calculer le score conceptuel basé sur les mots-clés uniques
	scoreConceptuel := CalculerScoreConceptuel(phrase)
	if scoreConceptuel < 0.3 { // Minimum 30% de contenu unique
		return false
	}

	// 4. Vérifier si la phrase n'est que des pronoms/articles/connecteurs génériques
	mots := strings.Fields(phrase)
	if len(mots) < 3 {
		return false
	}

	motsGeneriques := map[string]bool{
		"le": true, "la": true, "les": true, "un": true, "une": true, "des": true,
		"et": true, "ou": true, "mais": true, "car": true, "donc": true, "cependant": true,
		"alors": true, "ainsi": true, "d'ailleurs": true, "aussi": true, "peut": true,
		"peut-être": true, "semble": true, "paraît": true, "il": true, "elle": true,
		"ça": true, "ce": true, "cela": true, "celui": true, "quelque": true,
		"quelquefois": true, "souvent": true, "jamais": true, "toujours": true,
	}

	motsConcepts := 0
	for _, mot := range mots {
		motClean := strings.ToLower(strings.TrimSpace(mot))
		// Enlever la ponctuation
		motClean = regexp.MustCompile(`[.,!?;:\-'"]`).ReplaceAllString(motClean, "")
		if len(motClean) > 0 && !motsGeneriques[motClean] {
			motsConcepts++
		}
	}

	// Au moins 40% des mots doivent être des concepts
	percentConcepts := float64(motsConcepts) / float64(len(mots))
	if percentConcepts < 0.4 {
		return false
	}

	// 5. Vérifier si ce n'est pas juste une publicité ou un appel à clic
	publicites := []string{
		"acheter",
		"commander",
		"cliquez",
		"visitez",
		"consultez",
		"inscrivez",
		"téléchargez",
		"rejoignez",
		"abonnez",
		"profitez",
		"offre limitée",
		"promo",
		"réduction",
		"gratuit",
		"sans frais",
		"code promo",
		"lien",
		"url",
	}

	pubCount := 0
	for _, pubPattern := range publicites {
		if strings.Contains(lower, pubPattern) {
			pubCount++
		}
	}

	// Si trop de mots publicitaires, rejeter
	if pubCount > len(mots)/5 {
		return false
	}

	return true
}

// CalculerScoreConceptuel évalue la densité de concepts dans une phrase
func CalculerScoreConceptuel(phrase string) float64 {
	mots := strings.Fields(phrase)

	if len(mots) == 0 {
		return 0
	}

	// Mots très génériques (stopwords)
	stopwords := map[string]bool{
		"le": true, "la": true, "les": true, "un": true, "une": true, "des": true,
		"et": true, "ou": true, "mais": true, "car": true, "donc": true, "cependant": true,
		"alors": true, "ainsi": true, "d'ailleurs": true, "aussi": true, "peut": true,
		"il": true, "elle": true, "ça": true, "ce": true, "cela": true,
		"à": true, "au": true, "du": true, "de": true, "par": true, "pour": true,
		"avec": true, "sans": true, "dans": true, "sur": true, "sous": true,
		"avant": true, "après": true, "pendant": true, "entre": true, "parmi": true,
		"qui": true, "que": true, "quoi": true, "quel": true, "quelle": true,
		"est": true, "sont": true, "avoir": true, "aller": true, "faire": true,
		"venir": true, "pouvoir": true, "vouloir": true, "devoir": true, "sembler": true,
		"paraître": true, "devenir": true, "rester": true, "penser": true, "savoir": true,
		"dire": true, "donner": true, "prendre": true, "laisser": true, "mettre": true,
		"pas": true, "non": true, "plus": true, "moins": true, "très": true, "bien": true,
		"beaucoup": true, "peu": true, "quelque": true, "plusieurs": true, "tous": true,
		"tant": true, "si": true,
		"c": true, "s": true, "d": true, "l": true, "qu": true, "n": true,
	}

	// Compter les mots non-génériques
	conceptCount := 0
	for _, mot := range mots {
		motClean := strings.ToLower(strings.TrimSpace(mot))
		motClean = regexp.MustCompile(`[.,!?;:\-'"]`).ReplaceAllString(motClean, "")
		if len(motClean) > 2 && !stopwords[motClean] {
			conceptCount++
		}
	}

	// Bonus pour les mots longs (généralement plus conceptuels)
	for _, mot := range mots {
		motClean := strings.ToLower(strings.TrimSpace(mot))
		motClean = regexp.MustCompile(`[.,!?;:\-'"]`).ReplaceAllString(motClean, "")
		if len(motClean) > 8 {
			conceptCount += 2 // Mots longs comptent plus
		}
	}

	// Bonus pour les mots avec préfixes/suffixes scientifiques
	sciPattern := regexp.MustCompile(`(?i)(bio|chimio|cardio|pneumo|gastro|neuro|psycho|patho|itis|osis|emia|phagia|sclerosis|ectomy|plasty|graphy|metry|ology|itis|emia)`)
	if sciPattern.MatchString(phrase) {
		conceptCount += 3
	}

	return float64(conceptCount) / float64(len(mots))
}

// NettoyerHTML enlève les balises HTML inutiles et extrait le vrai contenu
func NettoyerHTML(texte string) string {
	// Enlever les éléments de bruit de haut niveau
	// Scripts et styles
	scriptRegex := regexp.MustCompile(`(?is)<script[^>]*>.*?</script>`)
	texte = scriptRegex.ReplaceAllString(texte, " ")

	styleRegex := regexp.MustCompile(`(?is)<style[^>]*>.*?</style>`)
	texte = styleRegex.ReplaceAllString(texte, " ")

	// Navigation et structure de page
	navRegex := regexp.MustCompile(`(?is)<nav[^>]*>.*?</nav>`)
	texte = navRegex.ReplaceAllString(texte, " ")

	footerRegex := regexp.MustCompile(`(?is)<footer[^>]*>.*?</footer>`)
	texte = footerRegex.ReplaceAllString(texte, " ")

	asideRegex := regexp.MustCompile(`(?is)<aside[^>]*>.*?</aside>`)
	texte = asideRegex.ReplaceAllString(texte, " ")

	// Mentions légales et copyright
	disclaimerRegex := regexp.MustCompile(`(?is)<div[^>]*class=["\'].*?(disclaimer|legal|copyright|privacy)["\'][^>]*>.*?</div>`)
	texte = disclaimerRegex.ReplaceAllString(texte, " ")

	// Formulaires
	formRegex := regexp.MustCompile(`(?is)<form[^>]*>.*?</form>`)
	texte = formRegex.ReplaceAllString(texte, " ")

	// Publicités
	adRegex := regexp.MustCompile(`(?is)<div[^>]*class=["\'].*?(ad|advertisement|sponsor)["\'][^>]*>.*?</div>`)
	texte = adRegex.ReplaceAllString(texte, " ")

	// Commentaires HTML
	commentRegex := regexp.MustCompile(`<!--.*?-->`)
	texte = commentRegex.ReplaceAllString(texte, " ")

	// Métadonnées et éléments vides
	metaRegex := regexp.MustCompile(`(?is)<meta[^>]*>`)
	texte = metaRegex.ReplaceAllString(texte, " ")

	linkRegex := regexp.MustCompile(`(?is)<link[^>]*>`)
	texte = linkRegex.ReplaceAllString(texte, " ")

	baseRegex := regexp.MustCompile(`(?is)<base[^>]*>`)
	texte = baseRegex.ReplaceAllString(texte, " ")

	titleRegex := regexp.MustCompile(`(?is)<title[^>]*>.*?</title>`)
	texte = titleRegex.ReplaceAllString(texte, " ")

	headRegex := regexp.MustCompile(`(?is)<head[^>]*>.*?</head>`)
	texte = headRegex.ReplaceAllString(texte, " ")

	// Essayer d'extraire le contenu principal
	// Priorité 1: <main>
	mainRegex := regexp.MustCompile(`(?is)<main[^>]*>(.*?)</main>`)
	mainMatches := mainRegex.FindStringSubmatch(texte)
	if len(mainMatches) > 1 && strings.TrimSpace(mainMatches[1]) != "" {
		texte = mainMatches[1]
	} else {
		// Priorité 2: <article>
		articleRegex := regexp.MustCompile(`(?is)<article[^>]*>(.*?)</article>`)
		articleMatches := articleRegex.FindStringSubmatch(texte)
		if len(articleMatches) > 1 && strings.TrimSpace(articleMatches[1]) != "" {
			texte = articleMatches[1]
		} else {
			// Priorité 3: <div class="content"> ou similaire
			contentRegex := regexp.MustCompile(`(?is)<div[^>]*class=["\'].*?content["\'][^>]*>(.*?)</div>`)
			contentMatches := contentRegex.FindStringSubmatch(texte)
			if len(contentMatches) > 1 && strings.TrimSpace(contentMatches[1]) != "" {
				texte = contentMatches[1]
			}
		}
	}

	// Enlever les balises HTML restantes
	tagRegex := regexp.MustCompile(`<[^>]+>`)
	texte = tagRegex.ReplaceAllString(texte, " ")

	// Enlever les entités HTML
	texte = strings.ReplaceAll(texte, "&nbsp;", " ")
	texte = strings.ReplaceAll(texte, "&lt;", "<")
	texte = strings.ReplaceAll(texte, "&gt;", ">")
	texte = strings.ReplaceAll(texte, "&amp;", "&")
	texte = strings.ReplaceAll(texte, "&quot;", "\"")
	texte = strings.ReplaceAll(texte, "&apos;", "'")

	// Enlever les espaces multiples
	spaceRegex := regexp.MustCompile(`\s+`)
	texte = spaceRegex.ReplaceAllString(texte, " ")

	// Enlever les caractères spéciaux de code
	codeRegex := regexp.MustCompile(`[{}\[\];:|/\\]`)
	texte = codeRegex.ReplaceAllString(texte, " ")

	return strings.TrimSpace(texte)
}

// TraiterFichier traite un fichier texte avec segmentation par article
func TraiterFichier(cheminFichier string) {
	contenu, err := os.ReadFile(cheminFichier)
	if err != nil {
		fmt.Printf("[ERREUR] Impossible de lire le fichier: %v\n", err)
		return
	}

	texte := string(contenu)

	// Vérifier si c'est du HTML
	isHTML := strings.HasSuffix(strings.ToLower(cheminFichier), ".html") || strings.Contains(texte, "<html") || strings.Contains(texte, "<article")

	if isHTML {
		fmt.Printf("🧹 Extraction des articles du HTML en cours...\n")
		// Extraire les articles
		articles := ExtraireArticles(texte)
		fmt.Printf("📰 %d article(s) détecté(s)\n\n", len(articles))

		// Traiter chaque article
		for i, article := range articles {
			articleTexte := NettoyerArticle(article)
			articleTexte = FiltrerLignesRepetitives(articleTexte)

			// Filtrer par langue - garder uniquement le français
			articleTexte, languesDetectees := FiltrerParLangue(articleTexte)

			if strings.TrimSpace(articleTexte) == "" {
				continue
			}

			// Afficher les langues détectées
			fmt.Printf("🌐 Langues détectées: ")
			for lang, count := range languesDetectees {
				if lang != "fr" {
					fmt.Printf("%s(%d) ", strings.ToUpper(lang), len(count))
				}
			}
			fmt.Printf("\n")

			fmt.Printf("═══════════════════════════════════════\n")
			fmt.Printf("Article %d\n", i+1)
			fmt.Printf("═══════════════════════════════════════\n\n")

			TraiterTexte(articleTexte, fmt.Sprintf("Article #%d", i+1))
		}
	} else {
		// Pour les fichiers texte, aussi filtrer par langue
		texteFiltré, languesDetectees := FiltrerParLangue(texte)
		if len(languesDetectees) > 0 {
			fmt.Printf("🌐 Langues détectées dans le document:\n")
			for lang, phrases := range languesDetectees {
				if lang != "fr" {
					fmt.Printf("  • %s: %d phrase(s)\n", strings.ToUpper(lang), len(phrases))
				}
			}
			fmt.Printf("  • FR: %d phrase(s) conservées pour analyse\n\n", len(languesDetectees["fr"]))
		}
		TraiterTexte(texteFiltré, cheminFichier)
	}
}

// TraiterTexte analyse et traite un texte complet - Approche phrase par phrase
func TraiterTexte(texte string, source string) {
	fmt.Printf("\n╔════════════════════════════════════════╗\n")
	fmt.Printf("║  TRAITEMENT - %s\n", source)
	fmt.Printf("╚════════════════════════════════════════╝\n\n")

	// Nouvelle approche : analyser phrase par phrase
	analyse := database.AnalyserParPhrases(texte)

	// Statistiques globales
	stats := database.StatistiquesAnalyse(analyse)
	comptage := stats["comptage"].(map[string]int)
	energieTotale := stats["energie_totale"].(float64)
	confianceMoyenne := stats["confiance_moyenne"].(float64)
	nbPhrases := stats["nb_phrases"].(int)

	fmt.Printf("[STATISTIQUES GLOBALES]\n")
	fmt.Printf("• Phrases analysées: %d\n", nbPhrases)
	fmt.Printf("• Catégories détectées: %d\n", len(comptage))
	fmt.Printf("• Énergie totale: %.2f\n", energieTotale)
	fmt.Printf("• Confiance moyenne: %.1f%%\n\n", confianceMoyenne*100)

	// Distribution par catégorie
	dominantes := database.CategoriesDominantes(analyse)
	if len(dominantes) > 0 {
		fmt.Printf("[DISTRIBUTION]\n")
		for _, catData := range dominantes {
			cat := catData["categorie"].(string)
			nombre := catData["nombre"].(int)
			pourcentage := catData["pourcentage"].(float64)
			barLen := int(pourcentage / 5)
			bar := ""
			for i := 0; i < barLen; i++ {
				bar += "█"
			}
			fmt.Printf("  %-12s %s %.0f%% (%d phrases)\n",
				cat, bar, pourcentage, nombre)
		}
		fmt.Println()
	}

	// Afficher l'analyse détaillée par catégorie
	fmt.Print(database.AfficherAnalyseDetaillee(analyse))
	// Apprendre pour chaque phrase dans sa catégorie (silencieusement, pas de logs)
	for _, phraseAnalyse := range analyse.Phrases {
		if phraseAnalyse.CategorieID > 0 {
			phrasTokens := database.TokeniserTexte(phraseAnalyse.Texte)
			for _, token := range phrasTokens {
				database.Apprendre(token, phraseAnalyse.CategorieID)
			}

			// Activer les neurones
			for i := range database.Neurones {
				if database.Neurones[i].CategorieID == phraseAnalyse.CategorieID {
					database.Neurones[i].Valeur += phraseAnalyse.Score
				}
			}
		}
	}

	// RÉSUMÉ SYNTHÉTIQUE - À LA FIN (sans logs)
	fmt.Printf("\n╔════════════════════════════════════════╗\n")
	fmt.Printf("║  RÉSUMÉ SYNTHÉTIQUE (ses propres mots)\n")
	fmt.Printf("╚════════════════════════════════════════╝\n\n")
	resumeSynthetique := database.GenererResumeSynthetique(analyse)
	fmt.Printf("%s\n\n", resumeSynthetique)
}

// InteractionInteractive permet à l'utilisateur d'interagir
func InteractionInteractive() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf(`
╔═══════════════════════════════════════╗
║   IA-ATOMIQUE v4.2 - Interface       ║
║   Analyse Phrase par Phrase           ║
╚═══════════════════════════════════════╝

Commandes:
  file <chemin>     - Traiter un fichier
  text              - Entrer du texte libre
  ask               - Poser une question
  stats             - Voir les statistiques
  quit              - Quitter

`)

	for {
		fmt.Print("\n> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "" {
			continue
		}

		parts := strings.Fields(input)
		commande := parts[0]

		switch commande {
		case "quit", "exit":
			fmt.Println("Au revoir!")
			return

		case "file":
			if len(parts) > 1 {
				TraiterFichier(parts[1])
			} else {
				fmt.Println("[ERREUR] Syntaxe: file <chemin>")
			}

		case "text":
			fmt.Print("Entrez votre texte (terminez par une ligne vide):\n> ")
			var texte string
			for {
				ligne, _ := reader.ReadString('\n')
				if strings.TrimSpace(ligne) == "" {
					break
				}
				texte += ligne
			}
			if texte != "" {
				TraiterTexte(texte, "entrée utilisateur")
			}

		case "ask":
			if len(parts) > 1 {
				question := strings.Join(parts[1:], " ")
				fmt.Printf("\n[QUESTION]\n%s\n", question)

				catAct, mots, conf := database.ProcesserTexte(question)
				var cat int
				var score int
				for c, s := range catAct {
					if s > score {
						score = s
						cat = c
					}
				}

				if cat > 0 {
					reponse := database.GenererReponse(cat, nil)
					fmt.Printf("\n[REPONSE]\n%s\n", reponse)
					if len(mots) > 0 {
						maxMots := 3
						if len(mots) < 3 {
							maxMots = len(mots)
						}
						fmt.Printf("\nMots clés: %s\n", strings.Join(mots[:maxMots], ", "))
					}
					fmt.Printf("Confiance: %.0f%%\n", conf*100)
				}
			} else {
				fmt.Println("[ERREUR] Syntaxe: ask <question>")
			}

		case "stats":
			VisualiserStats()
			fmt.Println("\n✓ Dashboard généré")

		default:
			TraiterTexte(input, "entrée directe")
		}
	}
}
