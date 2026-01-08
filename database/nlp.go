package database

import (
	"fmt"
	"sort"
	"strings"
	"unicode"
)

// ProcesserTexte traite un texte d'entrée et l'analyse
func ProcesserTexte(texte string) (map[int]int, []string, float64) {
	// Tokenisation et nettoyage
	tokens := TokeniserTexte(texte)

	// Extraction des mots clés
	motsClés := ExtraireMotsClés(tokens)

	// Activation des catégories
	catActivation := ActiverCategoriesParTexte(tokens)

	// Calcul de la confiance générale
	confiance := CalculerConfiance(catActivation)

	return catActivation, motsClés, confiance
}

// TokeniserTexte divise un texte en tokens
func TokeniserTexte(texte string) []string {
	var tokens []string

	// Convertir en minuscules et nettoyer
	texte = strings.ToLower(texte)

	// Splitter par espaces et ponctuations
	fields := strings.Fields(texte)

	for _, field := range fields {
		// Enlever la ponctuationdu début et fin
		mot := strings.TrimFunc(field, func(r rune) bool {
			return !unicode.IsLetter(r) && !unicode.IsNumber(r)
		})

		if len(mot) > 0 && !StopWords[mot] {
			tokens = append(tokens, mot)
		}
	}

	return tokens
}

// ExtraireMotsClés extrait les mots les plus significatifs
func ExtraireMotsClés(tokens []string) []string {
	// Compter les occurrences
	freq := make(map[string]int)
	poids := make(map[string]float64)

	for _, token := range tokens {
		freq[token]++

		// Chercher dans le lexique
		if word, ok := Words[token]; ok {
			poids[token] += word.Poids
		} else {
			poids[token] += 1.0 // Poids par défaut
		}
	}

	// Trier par fréquence * poids
	type motScore struct {
		mot   string
		score float64
	}

	var scores []motScore
	for mot, p := range poids {
		scores = append(scores, motScore{mot, float64(freq[mot]) * p})
	}

	// Sort descendant
	sort.Slice(scores, func(i, j int) bool {
		return scores[i].score > scores[j].score
	})

	// Retourner top 10
	var resultat []string
	for i := 0; i < len(scores) && i < 10; i++ {
		resultat = append(resultat, scores[i].mot)
	}

	return resultat
}

// ActiverCategoriesParTexte active les catégories basées sur le texte
func ActiverCategoriesParTexte(tokens []string) map[int]int {
	catActivation := make(map[int]int)

	for _, token := range tokens {
		if word, ok := Words[token]; ok && word.Categorie > 0 {
			catActivation[word.Categorie]++
		}
	}

	return catActivation
}

// CalculerConfiance calcule le score de confiance général
func CalculerConfiance(catActivation map[int]int) float64 {
	if len(catActivation) == 0 {
		return 0.0
	}

	total := 0
	for _, count := range catActivation {
		total += count
	}

	if total == 0 {
		return 0.0
	}

	// Confiance basée sur la distribution
	var maxCount int
	for _, count := range catActivation {
		if count > maxCount {
			maxCount = count
		}
	}

	return float64(maxCount) / float64(total)
}

// ResumerTexte crée un résumé du texte
func ResumerTexte(texte string, ratio float64) string {
	// Accepter ratio entre 0 et 1 (0.1 = 10%, 0.5 = 50%, etc)
	if ratio <= 0 {
		ratio = 0.3 // 30% par défaut
	}
	if ratio > 1 {
		ratio = 1 // Max 100%
	}

	// Splitter par phrases (le texte prétraité a été splitté par points, donc on re-split par espaces)
	phrases := strings.Split(texte, " ")

	// Filtrer les phrases vides
	var phrasesNonVides []string
	for _, p := range phrases {
		if strings.TrimSpace(p) != "" {
			phrasesNonVides = append(phrasesNonVides, strings.TrimSpace(p))
		}
	}
	phrases = phrasesNonVides

	if len(phrases) == 0 {
		return texte
	}

	// Calculer le nombre de mots à garder (basé sur ratio)
	nbMotsGardes := int(float64(len(phrases)) * ratio)

	// Minimum 20% du texte
	minMots := len(phrases) / 5
	if minMots < 10 {
		minMots = 10
	}
	if nbMotsGardes < minMots {
		nbMotsGardes = minMots
	}
	if nbMotsGardes > len(phrases) {
		nbMotsGardes = len(phrases)
	}

	// Scorer chaque mot
	type wordScore struct {
		word  string
		index int
		score float64
	}

	var scores []wordScore
	for i, word := range phrases {
		score := 1.0
		if w, ok := Words[word]; ok {
			score = w.Poids
		}
		scores = append(scores, wordScore{word, i, score})
	}

	// Trier par score décroissant
	sort.Slice(scores, func(i, j int) bool {
		return scores[i].score > scores[j].score
	})

	// Garder les meilleurs mots
	var motsGardes []wordScore
	for i := 0; i < nbMotsGardes && i < len(scores); i++ {
		motsGardes = append(motsGardes, scores[i])
	}

	// Réordonner par index original
	sort.Slice(motsGardes, func(i, j int) bool {
		return motsGardes[i].index < motsGardes[j].index
	})

	// Reconstruire le résumé
	var resume []string
	for _, ms := range motsGardes {
		resume = append(resume, ms.word)
	}

	return strings.Join(resume, " ")
}

// ============================================================================
// PHASE X+4: PARAPHRASE ENCYCLOPÉDIQUE
// ============================================================================
// Reformule un résumé atomisé en phrases encyclopédiques cohérentes

// GenerateEncyclopedicSummary reformule des mots clés en phrases lisibles
func GenerateEncyclopedicSummary(mots []string) string {
	if len(mots) == 0 {
		return ""
	}

	// Grouper les mots par segments de 5-7 mots pour former des phrases
	const segmentSize = 6
	var phrases []string

	for i := 0; i < len(mots); i += segmentSize {
		end := i + segmentSize
		if end > len(mots) {
			end = len(mots)
		}

		segment := mots[i:end]
		phrase := reformulerSegment(segment)
		if phrase != "" {
			phrases = append(phrases, phrase)
		}
	}

	// Joindre les phrases avec des points
	return strings.Join(phrases, ". ") + "."
}

// reformulerSegment transforme un segment de mots en phrase grammaticale
func reformulerSegment(mots []string) string {
	if len(mots) == 0 {
		return ""
	}

	// Structures de phrase encyclopédique:
	// "Le sujet est...", "Il comprend...", "Ses éléments..."

	texte := strings.Join(mots, " ")

	// Heuristiques pour créer des phrases grammaticales
	switch {
	// Pattern: determinant + nom + verbe
	case len(mots) >= 3 && (mots[0] == "la" || mots[0] == "le" || mots[0] == "les"):
		// Garder comme est: "La Mésange huppée est..."
		return capitalizeFirst(texte)

	// Pattern: nombre + substantif
	case isNumber(mots[0]):
		return capitalizeFirst(texte)

	// Pattern: adjectif/description
	default:
		// Ajouter un déterminant simple
		return capitalizeFirst("C'est " + texte)
	}
}

// capitalizeFirst met en majuscule le premier caractère
func capitalizeFirst(s string) string {
	if len(s) == 0 {
		return s
	}
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

// isNumber vérifie si une chaîne est un nombre
func isNumber(s string) bool {
	for _, c := range s {
		if !unicode.IsDigit(c) && c != ',' {
			return false
		}
	}
	return len(s) > 0
}

// SplitterParPhrases divise un texte en phrases
func SplitterParPhrases(texte string) []string {
	// Splitter par . ! ?
	var phrases []string
	var phrase string

	for _, char := range texte {
		phrase += string(char)

		if char == '.' || char == '!' || char == '?' {
			p := strings.TrimSpace(phrase)
			if len(p) > 10 {
				phrases = append(phrases, p)
			}
			phrase = ""
		}
	}

	if len(phrase) > 10 {
		phrases = append(phrases, strings.TrimSpace(phrase))
	}

	return phrases
}

// ScorerPhrase calcule le score d'une phrase
func ScorerPhrase(phrase string) float64 {
	tokens := TokeniserTexte(phrase)

	score := 0.0

	for _, token := range tokens {
		if word, ok := Words[token]; ok {
			score += word.Poids
		} else {
			score += 1.0 // Poids par défaut
		}
	}

	// Normaliser par la longueur
	if len(tokens) > 0 {
		score = score / float64(len(tokens))
	}

	return score
}

// GenererResumeDetaille crée un résumé détaillé avec l'analyse IA
func GenererResumeDetaille(texte string, catMain int, motsClés []string) string {
	resume := ResumerTexte(texte, 0.3)

	categorie := NumeroVersCategorie(catMain)

	resumeDetaille := fmt.Sprintf(`
=== ANALYSE IA-ATOMIQUE ===

[CATEGORIE PRINCIPALE]
%s

[MOTS CLÉS DETECTS]
%s

[RESUME]
%s

================
`, categorie, strings.Join(motsClés, " | "), resume)

	return resumeDetaille
}

// ============================================================================
// DÉCOUPAGE ULTRA-RAPIDE O(n) - Pour traitement instantané de grands textes
// ============================================================================

// BlocTexte représente un bloc de texte découpé
type BlocTexte struct {
	Contenu    string
	Mots       []string
	NumeroBloc int
	Taille     int // Nombre de mots
}

// DecouperTexteRapide divise un texte en blocs optimisés (O(n))
// Stratégie:
// 1. Découpe par paragraphes (sauts de ligne doubles)
// 2. Puis chaque paragraphe en blocs de ~100 mots
// 3. Complexité: O(n) où n = nombre de mots
func DecouperTexteRapide(texte string, motsParBloc int) []BlocTexte {
	if motsParBloc <= 0 {
		motsParBloc = 100 // Défaut: 100 mots par bloc
	}

	var blocs []BlocTexte
	numeroBloc := 0

	// Étape 1: Découper par paragraphes (sauts de ligne doubles)
	paragraphes := strings.Split(texte, "\n\n")

	for _, paragraphe := range paragraphes {
		// Ignorer les paragraphes vides
		paragraphe = strings.TrimSpace(paragraphe)
		if paragraphe == "" {
			continue
		}

		// Étape 2: Découper le paragraphe en mots
		mots := strings.Fields(paragraphe) // O(n) - split ultra rapide

		// Étape 3: Créer des blocs de taille fixe
		for i := 0; i < len(mots); i += motsParBloc {
			fin := i + motsParBloc
			if fin > len(mots) {
				fin = len(mots)
			}

			blocMots := mots[i:fin]
			contenu := strings.Join(blocMots, " ")

			blocs = append(blocs, BlocTexte{
				Contenu:    contenu,
				Mots:       blocMots,
				NumeroBloc: numeroBloc,
				Taille:     len(blocMots),
			})
			numeroBloc++
		}
	}

	return blocs
}

// DecouperAdaptatif ajuste la taille des blocs selon la longueur totale
// Pour très grands textes: augmenter motsParBloc pour limiter le nombre de blocs
func DecouperAdaptatif(texte string, maxBlocs int) []BlocTexte {
	// Compter le nombre de mots total
	mots := strings.Fields(texte)
	totalMots := len(mots)

	if totalMots == 0 {
		return []BlocTexte{}
	}

	// Calculer la taille adaptée: assurer au max maxBlocs blocs
	motsParBloc := totalMots / maxBlocs
	if motsParBloc < 50 {
		motsParBloc = 50 // Minimum 50 mots par bloc
	}
	if motsParBloc > 200 {
		motsParBloc = 200 // Maximum 200 mots par bloc
	}

	return DecouperTexteRapide(texte, motsParBloc)
}

// AnalyserBloc traite un bloc texte indépendamment
// Retourne: catégories, mots-clés, confiance
func AnalyserBloc(bloc BlocTexte) (map[int]int, []string, float64) {
	// Réutiliser la pipeline existante
	return ProcesserTexte(bloc.Contenu)
}

// FusionnerResultatsBlocs combine les résultats de tous les blocs
// Stratégie: moyenne pondérée par taille de bloc
func FusionnerResultatsBlocs(blocs []BlocTexte, resultats []map[int]int, confiances []float64) map[int]int {
	catGlobale := make(map[int]int)

	totalMots := 0
	for _, bloc := range blocs {
		totalMots += bloc.Taille
	}

	if totalMots == 0 {
		return catGlobale
	}

	// Fusionner avec poids proportionnel à la taille
	for i, resultat := range resultats {
		poids := float64(blocs[i].Taille) / float64(totalMots)

		for cat, val := range resultat {
			// Moyenne pondérée
			catGlobale[cat] += int(float64(val) * poids)
		}
	}

	return catGlobale
}

// StatistiquesBlocs retourne des stats rapides sur les blocs
func StatistiquesBlocs(blocs []BlocTexte) map[string]interface{} {
	totalMots := 0
	minTaille := 0
	maxTaille := 0

	if len(blocs) == 0 {
		return map[string]interface{}{
			"nombre_blocs":   0,
			"total_mots":     0,
			"taille_moyenne": 0,
		}
	}

	minTaille = blocs[0].Taille
	maxTaille = blocs[0].Taille

	for _, bloc := range blocs {
		totalMots += bloc.Taille
		if bloc.Taille < minTaille {
			minTaille = bloc.Taille
		}
		if bloc.Taille > maxTaille {
			maxTaille = bloc.Taille
		}
	}

	moyenneSize := totalMots / len(blocs)

	return map[string]interface{}{
		"nombre_blocs":   len(blocs),
		"total_mots":     totalMots,
		"taille_min":     minTaille,
		"taille_max":     maxTaille,
		"taille_moyenne": moyenneSize,
	}
}

// ============================================================================
// EXTRACTION DE PHRASES CLÉS - Pipeline d'énergie atomique
// ============================================================================

// Phrase représente une phrase avec son analyse énergétique et linguistique
type Phrase struct {
	Contenu          string
	Mots             []string
	Index            int
	Energie          float64 // E(Pi) - énergie intrinsèque
	EnergieTotal     float64 // Etotal(Pi) - énergie + cohérence + traduction
	Score            float64 // Score de sélection final
	MotsClés         []string
	EstFiltrée       bool
	Langue           string  // Langue détectée (FR, EN, DE, ES)
	EstTraduire      bool    // True si la phrase a été traduite
	FacteurConfiance float64 // γi ∈ [0.7, 1.0] - confiance traduction
}

// ExtrairePhrasesClés extrait les phrases les plus importantes d'un texte
// Pipeline: Découpage → Traduction → Énergie → Cohérence → Filtrage → Fusion
func ExtrairePhrasesClés(texte string, ratioConservation float64) []Phrase {
	// Étape 1: Découper en phrases
	phrases := DécouperEnPhrases(texte)
	if len(phrases) == 0 {
		return []Phrase{}
	}

	// Étape 1.5: Détecter langue et traduire en FR si nécessaire
	phrases = DetecterEtTraduirePhrases(phrases)

	// Étape 2: Calculer l'énergie intrinsèque de chaque phrase
	for i := range phrases {
		phrases[i].Energie = CalculerEnergiePrhase(&phrases[i])
	}

	// Étape 3: Ajouter la cohérence inter-phrases
	for i := range phrases {
		phrases[i].EnergieTotal = AjouterCoherence(&phrases[i], phrases)
	}

	// Étape 4: Filtrer les phrases de faible énergie (seuil adaptatif)
	// Accepter toutes les phrases avec au moins 2 mots pour test
	phrasesFiltrées := []Phrase{}
	for i := range phrases {
		if len(phrases[i].Mots) >= 2 {
			phrasesFiltrées = append(phrasesFiltrées, phrases[i])
		}
	}

	// Si vide après filtrage, retourner au moins quelque chose
	if len(phrasesFiltrées) == 0 {
		return []Phrase{}
	}

	// Étape 5: Déterminer le nombre de phrases à conserver
	// Toujours utiliser au minimum le ratio demandé
	ratio := ratioConservation
	if len(phrasesFiltrées) > 0 {
		// Adapter le ratio si possible
		ratioAdaptatif := CalculerRatioAdaptatif(phrasesFiltrées)
		if ratioAdaptatif > ratioConservation {
			ratio = ratioAdaptatif
		}
	}
	nombreAConserver := int(float64(len(phrasesFiltrées)) * ratio)
	if nombreAConserver < 1 && len(phrasesFiltrées) > 0 {
		nombreAConserver = 1 // Conserver au minimum 1 phrase
	}
	if nombreAConserver > len(phrasesFiltrées) {
		nombreAConserver = len(phrasesFiltrées)
	}

	// Étape 6: Trier par énergie et sélectionner les top phrases
	sort.Slice(phrasesFiltrées, func(i, j int) bool {
		return phrasesFiltrées[i].EnergieTotal > phrasesFiltrées[j].EnergieTotal
	})

	// Retourner dans l'ordre original pour préserver la cohérence textuelle
	resultat := phrasesFiltrées[:nombreAConserver]
	sort.Slice(resultat, func(i, j int) bool {
		return resultat[i].Index < resultat[j].Index
	})

	return resultat
}

// DécouperEnPhrases divise un texte en phrases
func DécouperEnPhrases(texte string) []Phrase {
	var phrases []Phrase

	// Découper par points, points-virgules, points d'exclamation
	delimiteurs := []string{".", "!", "?"}
	texte = strings.TrimSpace(texte)

	// Découpage simple mais efficace
	for _, delim := range delimiteurs {
		parts := strings.Split(texte, delim)
		texte = strings.Join(parts[:len(parts)-1], delim+"\n") // Garder le délimiteur
		if len(parts) > 1 && strings.TrimSpace(parts[len(parts)-1]) != "" {
			texte += delim + "\n" + strings.TrimSpace(parts[len(parts)-1])
		}
	}

	// Découper par sauts de ligne
	lignes := strings.Split(texte, "\n")
	index := 0

	for _, ligne := range lignes {
		ligne = strings.TrimSpace(ligne)
		if ligne == "" {
			continue
		}

		mots := strings.Fields(ligne)
		if len(mots) > 0 {
			phrases = append(phrases, Phrase{
				Contenu:  ligne,
				Mots:     mots,
				Index:    index,
				MotsClés: ExtraireMotsClés(mots),
			})
			index++
		}
	}

	return phrases
}

// CalculerEnergiePrhase calcule l'énergie intrinsèque E(Pi) d'une phrase
// E(Pi) = Σ αk * f(wk) où αk dépend du rôle syntaxique
func CalculerEnergiePrhase(phrase *Phrase) float64 {
	if len(phrase.Mots) == 0 {
		return 0.0
	}

	// Pénalité sévère pour les très courtes phrases (bruits de séparation)
	if len(phrase.Mots) == 1 && len(phrase.Contenu) <= 3 {
		return 0.0 // "." "/" "-" etc sont ignorées
	}

	if len(phrase.Mots) <= 2 && len(strings.FieldsFunc(phrase.Contenu, func(r rune) bool { return r == '.' || r == ',' })) == 0 {
		return 0.1 // Très courtes phrases: pénalité sérieuse
	}

	energie := 0.0
	motUtilsCount := 0

	// Mots avec poids syntaxique
	rolesLourds := map[int]float64{
		0: 1.5, // Premier mot (sujet probable)
		1: 1.3, // Deuxième mot (verbe probable)
	}

	for i, mot := range phrase.Mots {
		// Ignorer les mots entièrement vides
		if estMotVide(mot) {
			continue
		}

		motUtilsCount++

		// Score du mot: présent dans les mots-clés?
		motScore := 0.6 // Score réduit (pas 0.5)
		for _, motClé := range phrase.MotsClés {
			if strings.EqualFold(mot, motClé) {
				motScore = 1.0 // Boost si mot-clé
				break
			}
		}

		// Pondération selon position (rôle syntaxique approx)
		alpha := 1.0
		if poids, ok := rolesLourds[i]; ok {
			alpha = poids
		}

		energie += alpha * motScore
	}

	// Si phrase entièrement composée de mots vides
	if motUtilsCount == 0 {
		return 0.0
	}

	// Normaliser par mots utiles seulement
	return energie / float64(motUtilsCount)
}

// AjouterCoherence ajoute un terme d'influence inter-phrases avec facteur de confiance
// Etotal(Pi) = E(Pi) * γi + β * Σ sim(Pi, Pj)
// où γi ∈ [0.7, 1.0] est le facteur de confiance de traduction
func AjouterCoherence(phrase *Phrase, toutesLesPhotrases []Phrase) float64 {
	const beta = 0.2 // Coefficient de propagation
	coherence := 0.0

	// Calculer la similarité avec les autres phrases
	for j, autre := range toutesLesPhotrases {
		if j == phrase.Index {
			continue
		}

		// Similarité = nombre de mots-clés partagés
		sim := CalculerSimilarite(phrase.MotsClés, autre.MotsClés)
		coherence += sim
	}

	// Appliquer le facteur de confiance: γi (défaut 1.0 si pas traduite)
	facteur := phrase.FacteurConfiance
	if facteur == 0 {
		facteur = 1.0 // Défaut: confiance totale si pas traduite
	}

	// Etotal(Pi) = E(Pi) * γi + β * Σ sim(Pi, Pj)
	return phrase.Energie*facteur + beta*coherence
}

// CalculerSimilarite mesure la similarité entre deux ensembles de mots-clés
func CalculerSimilarite(cles1, cles2 []string) float64 {
	if len(cles1) == 0 || len(cles2) == 0 {
		return 0.0
	}

	intersection := 0
	for _, c1 := range cles1 {
		for _, c2 := range cles2 {
			if strings.EqualFold(c1, c2) {
				intersection++
				break
			}
		}
	}

	// Indice de Jaccard
	union := len(cles1) + len(cles2) - intersection
	if union == 0 {
		return 0.0
	}

	return float64(intersection) / float64(union)
}

// CalculerSeuilFiltering détermine le seuil énergétique pour filtrer le bruit
func CalculerSeuilFiltering(phrases []Phrase) float64 {
	if len(phrases) == 0 {
		return 0.0
	}

	// Seuil = moyenne - 0.5*écart-type (plus strict)
	var sum, sumSq float64
	for _, p := range phrases {
		sum += p.Energie
		sumSq += p.Energie * p.Energie
	}

	moyenne := sum / float64(len(phrases))
	variance := (sumSq / float64(len(phrases))) - (moyenne * moyenne)
	if variance < 0 {
		variance = 0
	}

	ecartType := 0.0
	if variance > 0 {
		ecartType = 0.5 * (moyenne * 0.3) // Écart-type approximé
	}

	// Seuil = moyenne - 0.5*écart-type (élimine bruits)
	seuil := moyenne - ecartType
	if seuil < 0.15 {
		seuil = 0.15 // Seuil minimum plus élevé
	}

	return seuil
}

// CalculerRatioAdaptatif ajuste le ratio de conservation selon la densité d'information
func CalculerRatioAdaptatif(phrases []Phrase) float64 {
	if len(phrases) == 0 {
		return 0.3
	}

	// Densité d'énergie moyenne
	var totalEnergie, countNonZero float64
	for _, p := range phrases {
		if p.EnergieTotal > 0.1 {
			totalEnergie += p.EnergieTotal
			countNonZero++
		}
	}

	if countNonZero == 0 {
		return 0.25
	}

	densité := totalEnergie / countNonZero

	// Ratio adaptatif: plus l'énergie est élevée, plus on conserve
	ratio := 0.2 + (densité * 0.15) // Entre 0.2 et 0.35
	if ratio > 0.35 {
		ratio = 0.35 // Maximum 35%
	}

	return ratio
}

// estMotVide vérifie si un mot est un mot vide (peu informatif)
func estMotVide(mot string) bool {
	motsVides := []string{
		"le", "la", "les", "un", "une", "des",
		"et", "ou", "mais", "donc", "car",
		"de", "à", "en", "par", "pour", "avec",
		"être", "avoir", "aller", "faire", "dire",
		"que", "qui", "où", "quand", "comment",
		"très", "plus", "moins", "bien", "mal",
		"ce", "celui", "cela", "ça", "il", "elle", "on",
	}

	motLower := strings.ToLower(mot)
	for _, vide := range motsVides {
		if motLower == vide {
			return true
		}
	}
	return false
}

// RéécrireSimplifiée simplifie une phrase tout en conservant l'énergie
// Retourne la phrase avec un score de lisibilité
func RéécrireSimplifiée(phrase *Phrase) (string, float64) {
	// Stratégie simple: garder les mots-clés et réorganiser
	motsClesTrouves := []string{}
	for _, mot := range phrase.Mots {
		for _, cleé := range phrase.MotsClés {
			if strings.EqualFold(mot, cleé) {
				motsClesTrouves = append(motsClesTrouves, mot)
				break
			}
		}
	}

	// Score de lisibilité: ratio mots-clés / longueur
	lisibilité := float64(len(motsClesTrouves)) / float64(len(phrase.Mots))
	if lisibilité > 1.0 {
		lisibilité = 1.0
	}

	// Phrase simplifiée = mots-clés + connecteurs
	phraseSimplifée := strings.Join(motsClesTrouves, " ")
	if phraseSimplifée == "" {
		phraseSimplifée = phrase.Contenu
		lisibilité = 0.5
	}

	return phraseSimplifée, lisibilité
}
