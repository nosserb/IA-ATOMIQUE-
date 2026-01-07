package database

import (
	"math"
	"math/rand"
	"strings"
)

// === PHASE 13+++: Dictionnaire de synonymes contextuels ===
var SynonymsDict = map[string][]string{
	"malheureux": {"regrettable", "désolé", "fâcheux", "malheureux"},
	"changer":    {"modifier", "transformer", "altérer", "changer"},
	"important":  {"crucial", "essentiel", "vital", "important"},
	"différent":  {"distinct", "varié", "divers", "différent"},
	"donne":      {"génère", "fournit", "conduit", "donne"},
	"avoir":      {"posséder", "détenir", "disposer de", "avoir"},
	"être":       {"constituer", "représenter", "figurer", "être"},
	"faire":      {"effectuer", "réaliser", "accomplir", "faire"},
	"donner":     {"attribuer", "procurer", "conférer", "donner"},
	"montrer":    {"démontrer", "révéler", "illustrer", "montrer"},
	"monde":      {"univers", "domaine", "sphère", "monde"},
	"jour":       {"époque", "période", "moment", "jour"},
	"suivant":    {"ultérieur", "postérieur", "subséquent", "suivant"},
	"simple":     {"élémentaire", "basique", "rudimentaire", "simple"},
	"nouveau":    {"inédit", "récent", "moderne", "nouveau"},
	"certain":    {"quelque", "divers", "maints", "certain"},
	"cas":        {"situation", "contexte", "occurrence", "cas"},
	"nombre":     {"quantité", "multitude", "plusieurs", "nombre"},
	"fois":       {"occasion", "moment", "instant", "fois"},
	"réalité":    {"vérité", "fait", "actualité", "réalité"},
}

// ConnecteurSemantique relie deux phrases selon axes sémantiques
type ConnecteurSemantique struct {
	Connecteur string
	Categories []string // Categories liées
	Prob       float64  // Probabilité d'utilisation
}

// ConnecteursLogiques par contexte sémantique
var ConnecteursParContexte = map[string][]ConnecteurSemantique{
	"addition": {
		{Connecteur: "En outre", Categories: []string{"TECH", "BUSINESS"}, Prob: 0.8},
		{Connecteur: "De plus", Categories: []string{"HISTOIRE", "SANTÉ"}, Prob: 0.7},
		{Connecteur: "Également", Categories: []string{"ALIMENTATION", "VERBE"}, Prob: 0.6},
		{Connecteur: "Par ailleurs", Categories: []string{"TECH"}, Prob: 0.5},
	},
	"contraste": {
		{Connecteur: "Cependant", Categories: []string{"HISTOIRE", "BUSINESS"}, Prob: 0.85},
		{Connecteur: "Or", Categories: []string{"TECH", "SANTÉ"}, Prob: 0.7},
		{Connecteur: "Néanmoins", Categories: []string{"HISTOIRE"}, Prob: 0.65},
		{Connecteur: "Toutefois", Categories: []string{"BUSINESS", "ALIMENTATION"}, Prob: 0.6},
	},
	"consequence": {
		{Connecteur: "Par conséquent", Categories: []string{"TECH", "BUSINESS"}, Prob: 0.8},
		{Connecteur: "Donc", Categories: []string{"VERBE", "SANTÉ"}, Prob: 0.75},
		{Connecteur: "Ainsi", Categories: []string{"HISTOIRE", "ALIMENTATION"}, Prob: 0.7},
		{Connecteur: "C'est pourquoi", Categories: []string{"HISTOIRE"}, Prob: 0.65},
	},
	"explication": {
		{Connecteur: "En effet", Categories: []string{"TECH", "SANTÉ"}, Prob: 0.8},
		{Connecteur: "C'est-à-dire", Categories: []string{"BUSINESS", "ALIMENTATION"}, Prob: 0.7},
		{Connecteur: "À savoir", Categories: []string{"HISTOIRE", "TECH"}, Prob: 0.65},
		{Connecteur: "En réalité", Categories: []string{"SANTÉ"}, Prob: 0.6},
	},
	"exemple": {
		{Connecteur: "Par exemple", Categories: []string{"TECH", "ALIMENTATION"}, Prob: 0.85},
		{Connecteur: "Notamment", Categories: []string{"BUSINESS", "HISTOIRE"}, Prob: 0.75},
		{Connecteur: "Comme", Categories: []string{"SANTÉ"}, Prob: 0.7},
		{Connecteur: "Prenons le cas de", Categories: []string{"HISTOIRE"}, Prob: 0.65},
	},
}

// MotsClésParAxe détermine l'axe sémantique d'un texte
func DetecterAxeSemantique(mots []string) string {
	// Compter les catégories dominantes
	compteur := make(map[int]int)

	for _, mot := range mots {
		if word, ok := Words[strings.ToLower(mot)]; ok {
			compteur[word.Categorie]++
		}
	}

	// Retrouver catégorie max
	maxCat := 0
	maxCount := 0
	for cat, count := range compteur {
		if count > maxCount {
			maxCount = count
			maxCat = cat
		}
	}

	// Mapper catégorie int → axe sémantique
	// 0=Neutre, 1=TECH, 2=HISTOIRE, 3=BUSINESS, 4=ALIMENTATION, 5=SANTÉ, 6=VERBE
	switch maxCat {
	case 1, 6: // TECH ou VERBE
		return "technique"
	case 2: // HISTOIRE
		return "historical"
	case 3: // BUSINESS
		return "commercial"
	case 5: // SANTÉ
		return "health"
	case 4: // ALIMENTATION
		return "food"
	default:
		return "neutral"
	}
}

// === PHASE 13: Tracking des connecteurs utilisés pour éviter répétitions ===
var connecteursUtilises = make(map[string]int)

// ReinitialiserConnecteurs réinitialise le tracking pour chaque résumé
func ReinitialiserConnecteurs() {
	connecteursUtilises = make(map[string]int)
}

// ChoisirConnecteurAvecDiversite sélectionne un connecteur en évitant répétitions
func ChoisirConnecteurAvecDiversite(axeSemantique string, typeConnexion string, connecteurPrecedent string) string {
	connecteurs, ok := ConnecteursParContexte[typeConnexion]
	if !ok {
		connecteurs = ConnecteursParContexte["addition"] // Default
	}

	// Trouver le meilleur connecteur pour cet axe
	var catTarget int
	switch axeSemantique {
	case "technique":
		catTarget = 1 // TECH
	case "historical":
		catTarget = 2 // HISTOIRE
	case "commercial":
		catTarget = 3 // BUSINESS
	case "health":
		catTarget = 5 // SANTÉ
	case "food":
		catTarget = 4 // ALIMENTATION
	default:
		catTarget = 0
	}

	// Chercher un connecteur adapté et non encore trop utilisé
	meilleursConnecteurs := []string{}
	for _, c := range connecteurs {
		for _, catStr := range c.Categories {
			var catInt int
			switch catStr {
			case "TECH":
				catInt = 1
			case "HISTOIRE":
				catInt = 2
			case "BUSINESS":
				catInt = 3
			case "ALIMENTATION":
				catInt = 4
			case "SANTÉ":
				catInt = 5
			case "VERBE":
				catInt = 6
			}

			if catInt == catTarget {
				// Éviter le connecteur précédent
				if c.Connecteur != connecteurPrecedent {
					meilleursConnecteurs = append(meilleursConnecteurs, c.Connecteur)
				}
			}
		}
	}

	// Si trouvé, sélectionner le moins utilisé
	if len(meilleursConnecteurs) > 0 {
		minUsed := math.MaxInt32
		selectionne := meilleursConnecteurs[0]
		for _, conn := range meilleursConnecteurs {
			if connecteursUtilises[conn] < minUsed {
				minUsed = connecteursUtilises[conn]
				selectionne = conn
			}
		}
		connecteursUtilises[selectionne]++
		return selectionne
	}

	// Fallback: variété de connecteurs
	fallbacks := []string{"Ainsi", "Par ailleurs", "Cependant", "De plus", "Ensuite"}
	minUsed := math.MaxInt32
	selectionne := fallbacks[0]
	for _, conn := range fallbacks {
		if conn != connecteurPrecedent && connecteursUtilises[conn] < minUsed {
			minUsed = connecteursUtilises[conn]
			selectionne = conn
		}
	}
	connecteursUtilises[selectionne]++
	return selectionne
}

// ChoisirConnecteur sélectionne un connecteur adapté au contexte (compatibilité)
func ChoisirConnecteur(axeSemantique string, typeConnexion string) string {
	return ChoisirConnecteurAvecDiversite(axeSemantique, typeConnexion, "")
}

// AjouterConnecteurs insère connecteurs sémantiques dans le texte généré
func AjouterConnecteurs(texte string, state GenerationState) string {
	mots := strings.Fields(texte)
	if len(mots) < 10 {
		return texte // Trop court
	}

	// Réinitialiser le tracking des connecteurs pour ce résumé
	ReinitialiserConnecteurs()

	// === PHASE 13: Diviser en chunks PLUS GRANDS pour moins de connecteurs ===
	const tailleChunk = 40 // Augmenté à 40 pour très peu de connecteurs
	chunks := []string{}

	for i := 0; i < len(mots); i += tailleChunk {
		fin := i + tailleChunk
		if fin > len(mots) {
			fin = len(mots)
		}
		chunks = append(chunks, strings.Join(mots[i:fin], " "))
	}

	// Insérer des connecteurs entre chunks (moins souvent)
	var resultat strings.Builder
	connecteurPrecedent := ""
	for i, chunk := range chunks {
		if i > 0 {
			// === PHASE 13: Utiliser ChoisirConnecteurAvecDiversite pour éviter répétitions ===
			// Déterminer type de connexion selon la cohérence
			var connecteur string
			if i*tailleChunk < len(state.Coherences) && state.Coherences[i*tailleChunk] > 0.7 {
				connecteur = ChoisirConnecteurAvecDiversite(DetecterAxeSemantique(mots), "explication", connecteurPrecedent)
			} else {
				connecteur = ChoisirConnecteurAvecDiversite(DetecterAxeSemantique(mots), "addition", connecteurPrecedent)
			}
			connecteurPrecedent = connecteur

			resultat.WriteString(" ")
			resultat.WriteString(connecteur)
			resultat.WriteString(" ")
		}
		resultat.WriteString(chunk)
	}

	return resultat.String()
}

// ReecrirePhrase améliore une phrase unique via génération
func ReecrirePhrase(phrase string, style string) string {
	// Tokenizer simple
	mots := strings.Fields(strings.TrimSpace(phrase))
	if len(mots) < 3 {
		return phrase // Trop courte
	}

	// Créer pseudo-phrase pour vectorisation
	p := Phrase{
		Contenu: phrase,
		Mots:    mots,
		Energie: 0.8,
	}
	p.MotsClés = ExtraireMotsClés(mots)

	params := ParamsParDefaut()

	// Adapter selon style
	switch style {
	case "avance":
		params.Eta = 0.4
		params.Gamma = 0.2
		params.TailleCible = len(mots) / 2
	case "professionnel":
		params.Eta = 0.15
		params.Gamma = 0.1
		params.TailleCible = len(mots) * 4 / 5
	default:
		params.Eta = 0.3
		params.Gamma = 0.15
		params.TailleCible = len(mots) * 3 / 4
	}

	// Générer version réécrite
	result := GenererHumanisation(phrase, style, params)

	// Ajouter ponctuation si manquante
	if !strings.HasSuffix(result, ".") && !strings.HasSuffix(result, "!") &&
		!strings.HasSuffix(result, "?") && !strings.HasSuffix(result, ",") {
		result += "."
	}

	return result
}

// CombinerPhrasesGenerees fusionne plusieurs phrases générées
func CombinerPhrasesGenerees(states []GenerationState) string {
	allMotsGenerés := []string{}

	for _, state := range states {
		allMotsGenerés = append(allMotsGenerés, state.MotsGenerés...)
	}

	texte := strings.Join(allMotsGenerés, " ")

	// Nettoyer
	texte = strings.TrimSpace(texte)

	// Ajouter ponctuation finale
	if len(texte) > 0 && !strings.HasSuffix(texte, ".") &&
		!strings.HasSuffix(texte, "!") && !strings.HasSuffix(texte, "?") {
		texte += "."
	}

	return texte
}

// FiltrerRepetitions élimine les mots répétés consécutifs
func FiltrerRepetitions(mots []string) []string {
	if len(mots) == 0 {
		return mots
	}

	resultat := []string{mots[0]}
	for i := 1; i < len(mots); i++ {
		if mots[i] != mots[i-1] {
			resultat = append(resultat, mots[i])
		}
	}
	return resultat
}

// AméliorerCoherence améliore la cohérence en réarrangeant mots
func AméliorerCoherence(state *GenerationState, seuilMin float64) float64 {
	coherence := FiltrerCoherence(*state)

	// Si mauvais, essayer d'améliorer en réarrangeant
	if coherence < seuilMin && len(state.MotsGenerés) > 3 {
		// Stratégie simple: garder les 70% meilleurs, regénérer 30%
		nbAGarder := len(state.MotsGenerés) * 70 / 100

		stateAmelioree := GenerationState{
			VecteurCible: state.VecteurCible,
		}

		for i := 0; i < nbAGarder; i++ {
			if i < len(state.MotsGenerés) {
				stateAmelioree.MotsGenerés = append(stateAmelioree.MotsGenerés, state.MotsGenerés[i])
			}
		}

		*state = stateAmelioree
		return FiltrerCoherence(stateAmelioree)
	}

	return coherence
}

// === PHASE 13: Post-traitement pour remplacer connecteurs répétitifs ===

// DiversifierConnecteurs remplace "En outre" et autres connecteurs par variantes naturelles
func DiversifierConnecteurs(texte string) string {
	// Liste de variantes pour chaque connecteur
	variantes := map[string][]string{
		"En outre": {"De plus", "Par ailleurs", "Ainsi", "Ensuite", "Cependant"},
		"outre":    {"plus", "ailleurs", "ainsi", "ensuite", "cependant"},
	}

	mots := strings.Fields(texte)
	if len(mots) < 2 {
		return texte
	}

	// Parcourir et remplacer "En outre" par des variantes
	connecteurIndex := make(map[string]int)
	for i := 0; i < len(mots)-1; i++ {
		// Chercher "En outre" (deux mots)
		if strings.EqualFold(mots[i], "En") && strings.EqualFold(mots[i+1], "outre") {
			// Obtenir les variantes de "En outre"
			if vars, ok := variantes["En outre"]; ok {
				// Sélectionner une variante différente selon occurrence
				idx := connecteurIndex["En outre"] % len(vars)
				connecteurIndex["En outre"]++

				// Remplacer "En outre" par variante
				mots[i] = vars[idx]
				// Supprimer le deuxième mot "outre"
				mots = append(mots[:i+1], mots[i+2:]...)
				i-- // Ajuster l'index après suppression
			}
		}
	}

	return strings.Join(mots, " ")
}

// PostTraiterResume applique les améliorations finales au résumé
func PostTraiterResume(texte string) string {
	// Diversifier les connecteurs
	texte = DiversifierConnecteurs(texte)

	// === PHASE 13+: Améliorations de fluidité ===

	// 1. Nettoyer les espaces multiples
	for strings.Contains(texte, "  ") {
		texte = strings.ReplaceAll(texte, "  ", " ")
	}

	mots := strings.Fields(texte)

	// 2. Supprimer les doublons immédiats (mot mot → mot)
	motsFiltres := []string{}
	prevMot := ""
	for _, mot := range mots {
		motClean := strings.ToLower(strings.TrimSuffix(strings.TrimSuffix(mot, "."), ","))
		prevMotClean := strings.ToLower(strings.TrimSuffix(strings.TrimSuffix(prevMot, "."), ","))

		// Ajouter le mot seulement s'il est différent du précédent
		if motClean != prevMotClean || prevMot == "" {
			motsFiltres = append(motsFiltres, mot)
			prevMot = mot
		}
	}
	mots = motsFiltres

	// === PHASE 13+++: Filtre anti-répétition (mot répété < 5 mots d'écart) ===
	motsFiltres = []string{}
	derniereMention := make(map[string]int) // position de dernière mention du mot
	for i, mot := range mots {
		motClean := strings.ToLower(strings.TrimSuffix(strings.TrimSuffix(mot, "."), ","))

		// Si mot répété à moins de 5 mots d'écart, skip
		if lastPos, exists := derniereMention[motClean]; exists && i-lastPos < 5 {
			continue // Ignorer cette répétition
		}

		derniereMention[motClean] = i
		motsFiltres = append(motsFiltres, mot)
	}
	mots = motsFiltres

	// 3. Ajouter majuscule après point si manquante
	for i := 0; i < len(mots)-1; i++ {
		if strings.HasSuffix(mots[i], ".") && len(mots[i+1]) > 0 {
			mots[i+1] = strings.ToUpper(mots[i+1][:1]) + mots[i+1][1:]
		}
	}

	// 4. Supprimer les répétitions de noms propres (heuristique simple)
	motsCounts := make(map[string]int)
	motsRares := make(map[string]bool)

	for _, mot := range mots {
		motClean := strings.ToLower(mot)
		// Identifier les noms propres (capitalisés et longs)
		if len(mot) > 3 && mot[0] >= 'A' && mot[0] <= 'Z' {
			motsCounts[motClean]++
			if motsCounts[motClean] > 2 {
				motsRares[motClean] = true
			}
		}
	}

	// Filtrer les occurrences > 2 des noms propres rares
	motOccurrences := make(map[string]int)
	motsFiltres = []string{}
	for _, mot := range mots {
		motClean := strings.ToLower(mot)
		if motsRares[motClean] {
			motOccurrences[motClean]++
			if motOccurrences[motClean] <= 2 {
				motsFiltres = append(motsFiltres, mot)
			}
		} else {
			motsFiltres = append(motsFiltres, mot)
		}
	}
	mots = motsFiltres

	// === PHASE 13+++: Diversification par synonymes contextuels ===
	motsFiltres = []string{}
	compteurMots := make(map[string]int) // compte les occurrences de chaque mot

	for _, mot := range mots {
		motClean := strings.ToLower(strings.TrimSuffix(strings.TrimSuffix(mot, "."), ","))
		compteurMots[motClean]++

		// Si mot fréquent (>2 occurrences) et dans le dictionnaire, remplacer par synonyme
		if compteurMots[motClean] > 2 {
			if synonymes, exists := SynonymsDict[motClean]; exists && len(synonymes) > 1 {
				// Choisir un synonyme aléatoire (en excluant le mot original 70% du temps)
				synChoisi := synonymes[rand.Intn(len(synonymes))]
				if compteurMots[motClean]%3 == 0 && synChoisi != motClean {
					// Remplacer le mot avec un synonyme tous les 3 occurrences
					if mot[len(mot)-1] == '.' || mot[len(mot)-1] == ',' {
						// Préserver la ponctuation
						motsFiltres = append(motsFiltres, synChoisi+string(mot[len(mot)-1]))
					} else {
						motsFiltres = append(motsFiltres, synChoisi)
					}
					continue
				}
			}
		}

		motsFiltres = append(motsFiltres, mot)
	}
	mots = motsFiltres

	// 5. Lisser les transitions: remplacer connecteurs mécaniques
	texteSmoothed := strings.Join(mots, " ")
	texteSmoothed = strings.ReplaceAll(texteSmoothed, "; Ainsi ", "; ainsi ")
	texteSmoothed = strings.ReplaceAll(texteSmoothed, ", Ainsi ", ", ainsi ")
	texteSmoothed = strings.ReplaceAll(texteSmoothed, "; De plus ", "; de plus ")
	texteSmoothed = strings.ReplaceAll(texteSmoothed, ", De plus ", ", de plus ")

	return texteSmoothed
}
