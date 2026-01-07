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
	if ratio <= 0 || ratio >= 1 {
		ratio = 0.3 // 30% par défaut
	}

	// Splitter par phrases
	phrases := SplitterParPhrases(texte)

	if len(phrases) == 0 {
		return ""
	}

	// Calculer le nombre de phrases à garder
	nbPhrasesGardees := int(float64(len(phrases)) * ratio)
	if nbPhrasesGardees < 1 {
		nbPhrasesGardees = 1
	}

	// Scorer chaque phrase
	type phraseScore struct {
		phrase string
		index  int
		score  float64
	}

	var scores []phraseScore

	for i, phrase := range phrases {
		score := ScorerPhrase(phrase)
		scores = append(scores, phraseScore{phrase, i, score})
	}

	// Trier par score
	sort.Slice(scores, func(i, j int) bool {
		return scores[i].score > scores[j].score
	})

	// Garder les meilleures phrases et les réordonner
	var phrasesGardees []phraseScore
	for i := 0; i < nbPhrasesGardees && i < len(scores); i++ {
		phrasesGardees = append(phrasesGardees, scores[i])
	}

	// Trier par index original
	sort.Slice(phrasesGardees, func(i, j int) bool {
		return phrasesGardees[i].index < phrasesGardees[j].index
	})

	// Construire le résumé
	var resume []string
	for _, ps := range phrasesGardees {
		resume = append(resume, strings.TrimSpace(ps.phrase))
	}

	return strings.Join(resume, " ")
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
