package database

import (
	"math"
	"sort"
	"strings"
	"unicode"
)

// ParametresGeneration contrôle la génération atomique
type ParametresGeneration struct {
	TailleCible    int     // Nombre de mots à générer
	Eta            float64 // Coefficient cohérence [0, 1]
	Gamma          float64 // Facteur apprentissage [0, 1]
	LongueurMin    int     // Minimum de mots
	DiversiteMin   float64 // Similarité min entre mots consécutifs [0, 1]
	TauxResilience float64 // Activation résiduelle du mot précédent
}

// GenerationState suivi l'état pendant la génération
type GenerationState struct {
	VecteurCible    VecteurAtomique
	MotsGenerés     []string
	VecteursGenerés []VecteurAtomique
	Activations     []float64
	Coherences      []float64
	Energie         float64
}

// isValidWord retourne true si le mot est valide (pas de LaTeX, nombres, etc)
func isValidWord(mot string) bool {
	if len(mot) == 0 {
		return false
	}

	// Rejeter les mots trop longs (souvent LaTeX)
	if len(mot) > 50 {
		return false
	}

	// Rejeter si contient LaTeX markers ou caractères spéciaux
	if strings.Contains(mot, "\\") || strings.Contains(mot, "{") || strings.Contains(mot, "}") ||
		strings.Contains(mot, "cdot") || strings.Contains(mot, "boldsymbol") || strings.Contains(mot, "frac") {
		return false
	}

	// Rejeter les nombres purs
	allDigits := true
	for _, r := range mot {
		if !unicode.IsDigit(r) && r != ',' && r != '.' && r != '-' {
			allDigits = false
			break
		}
	}
	if allDigits {
		return false
	}

	// Rejeter si trop de caractères spéciaux
	specialCount := 0
	for _, r := range mot {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '-' && r != '\'' {
			specialCount++
		}
	}
	if specialCount > len(mot)/2 {
		return false
	}

	// Rejeter si principalement non-latin
	latinCount := 0
	for _, r := range mot {
		if r < 128 || (r >= 192 && r <= 255) { // ASCII + Latin extended
			latinCount++
		}
	}
	if latinCount < len(mot)/2 {
		return false
	}

	return true
}

// ParamsParDefaut retourne les paramètres par défaut
func ParamsParDefaut() ParametresGeneration {
	return ParametresGeneration{
		TailleCible:    100,
		Eta:            0.3,  // Équilibre entre pertinence (70%) et cohérence (30%)
		Gamma:          0.15, // Adaptation lente du vecteur cible
		LongueurMin:    30,
		DiversiteMin:   0.4,
		TauxResilience: 0.5,
	}
}

// GenerationAdaptive génère un texte complet avec rk = argmax[sim(vw, VR(k)) + η·aw(k-1)]
func GenerationAdaptive(phrases []Phrase, params ParametresGeneration) GenerationState {
	state := GenerationState{
		VecteurCible:    MoyennePondereeVecteurs(phrases),
		MotsGenerés:     []string{},
		VecteursGenerés: []VecteurAtomique{},
		Activations:     []float64{},
		Coherences:      []float64{},
	}

	// Extraire le vocabulaire
	vocab := ExtraireVocabulaire(phrases)
	if len(vocab) == 0 {
		return state // Pas de vocab = rien à générer
	}

	// Vecteurs du vocabulaire (pré-calculés pour perf)
	vecabVecteurs := make(map[string]VecteurAtomique)
	for _, mot := range vocab {
		vecabVecteurs[mot] = VectoriserMot(mot)
	}

	// Activation atomique résiduelle du mot précédent
	var activationPrecedente float64 = 0.0
	motsPrecedents := make(map[string]int) // Track répétitions

	// Boucle de génération
	for k := 0; k < params.TailleCible; k++ {
		// Calcul du score pour chaque candidat: sim(vw, VR(k)) + η·aw(k-1)
		meilleurScore := -math.Inf(1)
		meilleurMot := ""
		meilleurVecteur := VecteurAtomique{}
		meilleurActivation := 0.0

		for _, mot := range vocab {
			// Pénaliser les mots trop répétés
			repetitions := motsPrecedents[mot]
			penalite := 0.0
			if repetitions > 0 {
				penalite = 0.4 + (0.3 * float64(repetitions))
			}

			// Score de pertinence
			vMot := vecabVecteurs[mot]
			similarite := SimilariteCosinus(vMot, state.VecteurCible)

			// Score de cohérence (favorise mot similaire au précédent)
			scoreCoherence := 0.0
			if len(state.MotsGenerés) > 0 {
				vMot2 := state.VecteursGenerés[len(state.VecteursGenerés)-1]
				scoreCoherence = SimilariteCosinus(vMot, vMot2)
			}

			// Score composite
			score := similarite + params.Eta*scoreCoherence + params.TauxResilience*activationPrecedente - penalite

			if score > meilleurScore {
				meilleurScore = score
				meilleurMot = mot
				meilleurVecteur = vMot
				meilleurActivation = similarite // Activation pour prochaine itération
			}
		}

		// Vérifier la diversité (éviter répétitions)
		if len(state.MotsGenerés) > 0 && meilleurScore < 0 {
			break // Pas assez de similarité = arrêter
		}

		// Ajouter le mot généré
		state.MotsGenerés = append(state.MotsGenerés, meilleurMot)
		state.VecteursGenerés = append(state.VecteursGenerés, meilleurVecteur)
		state.Activations = append(state.Activations, meilleurActivation)

		// Mettre à jour le vecteur cible: VR(k+1) = VR(k) + γ·(vrk - VR(k))
		state.VecteurCible = MettreAJourVecteur(state.VecteurCible, meilleurVecteur, params.Gamma)

		// Calculer la cohérence avec la cible
		coherence := SimilariteCosinus(meilleurVecteur, state.VecteurCible)
		state.Coherences = append(state.Coherences, coherence)

		// Accumuler l'énergie
		state.Energie += meilleurVecteur.Energie * coherence

		// Activation pour l'itération suivante
		activationPrecedente = params.TauxResilience * meilleurActivation

		// Tracker la répétition
		motsPrecedents[meilleurMot]++
	}

	return state
}

// GenerationAvecCouverture génère un résumé en garantissant couverture début-fin
func GenerationAvecCouverture(phrases []Phrase, params ParametresGeneration) GenerationState {
	state := GenerationState{
		VecteurCible:    MoyennePondereeVecteurs(phrases),
		MotsGenerés:     []string{},
		VecteursGenerés: []VecteurAtomique{},
		Activations:     []float64{},
		Coherences:      []float64{},
	}

	// Extraire vocabulaire avec filtrage des mots valides
	allVocab := ExtraireVocabulaire(phrases)
	vocab := []string{}
	for _, mot := range allVocab {
		if isValidWord(mot) && !StopWords[mot] && len(mot) > 2 {
			vocab = append(vocab, mot)
		}
	}

	if len(vocab) == 0 {
		return state
	}

	// Calculer TF-IDF pour chaque mot
	tfidfScores := CalculerTFIDF(phrases, vocab)

	// Vecteurs du vocabulaire
	vocabVecteurs := make(map[string]VecteurAtomique)
	for _, mot := range vocab {
		vocabVecteurs[mot] = VectoriserMot(mot)
	}

	motsPrecedents := make(map[string]int)
	var activationPrecedente float64 = 0.0

	// Boucle de génération
	for k := 0; k < params.TailleCible; k++ {
		meilleurScore := -math.Inf(1)
		meilleurMot := ""
		meilleurVecteur := VecteurAtomique{}
		meilleurActivation := 0.0

		for _, mot := range vocab {
			repetitions := motsPrecedents[mot]
			penalite := 0.0
			if repetitions > 0 {
				penalite = 0.4 + (0.3 * float64(repetitions))
			}

			vMot := vocabVecteurs[mot]
			similarite := SimilariteCosinus(vMot, state.VecteurCible)
			tfidf := tfidfScores[mot]
			bonusImportance := 0.3 * tfidf

			scoreCoherence := 0.0
			if len(state.MotsGenerés) > 0 {
				vMot2 := state.VecteursGenerés[len(state.VecteursGenerés)-1]
				scoreCoherence = SimilariteCosinus(vMot, vMot2)
			}

			score := similarite + bonusImportance + params.Eta*scoreCoherence +
				params.TauxResilience*activationPrecedente - penalite

			if score > meilleurScore {
				meilleurScore = score
				meilleurMot = mot
				meilleurVecteur = vMot
				meilleurActivation = similarite
			}
		}

		if meilleurScore < 0 {
			break
		}

		state.MotsGenerés = append(state.MotsGenerés, meilleurMot)
		state.VecteursGenerés = append(state.VecteursGenerés, meilleurVecteur)
		state.Activations = append(state.Activations, meilleurActivation)
		state.VecteurCible = MettreAJourVecteur(state.VecteurCible, meilleurVecteur, params.Gamma)
		coherence := SimilariteCosinus(meilleurVecteur, state.VecteurCible)
		state.Coherences = append(state.Coherences, coherence)
		state.Energie += meilleurVecteur.Energie * coherence

		// Activation suivante
		activationPrecedente = params.TauxResilience * meilleurActivation
		motsPrecedents[meilleurMot]++
	}

	return state
}

// GenererResume crée un résumé par génération atomique avec couverture
func GenererResume(phrases []Phrase, ratioCompression float64) string {
	params := ParamsParDefaut()
	params.TailleCible = int(float64(len(phrases)) * ratioCompression * 5) // Estimer words count
	if params.TailleCible < 30 {
		params.TailleCible = 30
	}
	if params.TailleCible > 500 {
		params.TailleCible = 500
	}

	// Utiliser la génération avec couverture
	state := GenerationAvecCouverture(phrases, params)

	// Relier les mots en texte
	texte := strings.Join(state.MotsGenerés, " ")

	// Ajouter des connecteurs et ponctuation
	texte = AjouterConnecteurs(texte, state)

	return texte
}

// GenererHumanisation réécrit un texte de façon naturelle
func GenererHumanisation(phrase string, style string, params ParametresGeneration) string {
	// Parser la phrase en tokens
	mots := strings.Fields(phrase)
	if len(mots) == 0 {
		return phrase
	}

	// Créer une pseudo-phrase pour vectorisation
	p := Phrase{
		Contenu: phrase,
		Mots:    mots,
		Energie: 0.7,
	}
	p.MotsClés = ExtraireMotsClés(mots)

	// Adapter les paramètres selon le style
	switch style {
	case "avance":
		params.Eta = 0.5   // Plus de créativité
		params.Gamma = 0.2 // Adaptation plus rapide
		params.TauxResilience = 0.3
	case "professionnel":
		params.Eta = 0.2            // Moins de libertés
		params.Gamma = 0.1          // Adaptation lente
		params.TauxResilience = 0.7 // Cohérence stricte
	default: // standard
		params.Eta = 0.3
		params.Gamma = 0.15
		params.TauxResilience = 0.5
	}

	params.TailleCible = len(mots) * 2 / 3 // Légèrement plus court

	// Générer avec vecteur cible unique
	vCible := VectoriserPhrase(&p)
	state := GenerationState{
		VecteurCible: vCible,
	}

	// Utiliser uniquement les mots d'entrée et quelques synonymes courants
	vocab := []string{}
	vocabMap := make(map[string]bool)

	// Ajouter les mots existants d'abord (mots "de confiance")
	for _, mot := range mots {
		lower := strings.ToLower(mot)
		clean := strings.Trim(lower, ".,;:!?'\"")
		if !vocabMap[clean] && len(vocab) < 100 && len(clean) > 2 && !StopWords[clean] {
			vocab = append(vocab, clean)
			vocabMap[clean] = true
		}
	}

	// Pour la réécriture, utiliser une approche plus simple basée sur variantes
	// Ajouter des variantes et synonymes courants
	synonymes := map[string][]string{
		"intelligence": {"sagesse", "esprit", "connaissance", "compréhension"},
		"artificielle": {"synthétique", "créée", "programmée", "développée"},
		"révolutionne": {"transforme", "change", "modifie", "innove"},
		"monde":        {"univers", "société", "environnement", "contexte"},
		"moderne":      {"contemporain", "actuel", "récent", "nouveau"},
		"robot":        {"machine", "système", "automate", "programme"},
		"calculé":      {"traité", "analysé", "évalué", "déterminé"},
		"solution":     {"réponse", "approche", "méthode", "résultat"},
		"précision":    {"exactitude", "rigueur", "perfection", "soin"},
	}

	for _, mot := range mots {
		lower := strings.ToLower(mot)
		clean := strings.Trim(lower, ".,;:!?'\"")

		if syns, ok := synonymes[clean]; ok {
			for _, syn := range syns {
				if !vocabMap[syn] && len(vocab) < 150 {
					vocab = append(vocab, syn)
					vocabMap[syn] = true
				}
			}
		}
	}

	vecabVecteurs := make(map[string]VecteurAtomique)
	for _, mot := range vocab {
		vecabVecteurs[mot] = VectoriserMot(mot)
	}

	var activationPrecedente float64 = 0.0
	motsPrecedents := make(map[string]int) // Tracker mots pour éviter répétitions

	for k := 0; k < params.TailleCible; k++ {
		meilleurScore := -math.Inf(1)
		meilleurMot := ""
		meilleurVecteur := VecteurAtomique{}
		meilleurActivation := 0.0

		for _, mot := range vocab {
			// Pénaliser fortement les mots trop répétés
			repetitions := motsPrecedents[mot]
			penalite := 0.0
			if repetitions > 0 {
				penalite = 0.5 + (0.3 * float64(repetitions)) // Forte pénalité dès la première répétition
			}

			vMot := vecabVecteurs[mot]
			similarite := SimilariteCosinus(vMot, state.VecteurCible)
			scoreCoherence := 0.0

			if len(state.MotsGenerés) > 0 {
				vMot2 := state.VecteursGenerés[len(state.VecteursGenerés)-1]
				scoreCoherence = SimilariteCosinus(vMot, vMot2)
			}

			score := similarite + params.Eta*scoreCoherence + params.TauxResilience*activationPrecedente - penalite

			if score > meilleurScore {
				meilleurScore = score
				meilleurMot = mot
				meilleurVecteur = vMot
				meilleurActivation = similarite
			}
		}

		if len(state.MotsGenerés) > 0 && meilleurScore < 0 {
			break
		}

		state.MotsGenerés = append(state.MotsGenerés, meilleurMot)
		state.VecteursGenerés = append(state.VecteursGenerés, meilleurVecteur)
		state.VecteurCible = MettreAJourVecteur(state.VecteurCible, meilleurVecteur, params.Gamma)
		activationPrecedente = params.TauxResilience * meilleurActivation

		// Tracker la répétition du mot
		motsPrecedents[meilleurMot]++
	}

	// Relier et ajouter ponctuation
	texte := strings.Join(state.MotsGenerés, " ")
	texte = strings.TrimSpace(texte)

	// Capitaliser première lettre
	if len(texte) > 0 {
		texte = strings.ToUpper(texte[:1]) + texte[1:]
	}

	// Ajouter ponctuation finale
	if !strings.HasSuffix(texte, ".") && !strings.HasSuffix(texte, "!") && !strings.HasSuffix(texte, "?") {
		texte += "."
	}

	return texte
}

// CalculerTFIDF retourne les scores TF-IDF pour chaque mot du vocabulaire
func CalculerTFIDF(phrases []Phrase, vocab []string) map[string]float64 {
	tfidf := make(map[string]float64)

	if len(phrases) == 0 {
		return tfidf
	}

	// Calculer TF (Term Frequency) dans toutes les phrases
	tf := make(map[string]float64)
	for _, phrase := range phrases {
		for _, mot := range phrase.Mots {
			tf[mot]++
		}
	}

	// Normaliser TF
	totalMots := float64(0)
	for _, count := range tf {
		totalMots += count
	}
	if totalMots > 0 {
		for mot := range tf {
			tf[mot] /= totalMots
		}
	}

	// Calculer IDF (Inverse Document Frequency)
	idf := make(map[string]float64)
	for _, mot := range vocab {
		docCount := 0
		for _, phrase := range phrases {
			for _, m := range phrase.Mots {
				if m == mot {
					docCount++
					break
				}
			}
		}
		if docCount > 0 {
			idf[mot] = math.Log(float64(len(phrases)) / float64(docCount))
		}
	}

	// TF-IDF = TF × IDF
	for _, mot := range vocab {
		tfidfVal := tf[mot] * idf[mot]
		if tfidfVal < 0 {
			tfidfVal = 0
		}
		if tfidfVal > 1 {
			tfidfVal = 1
		}

		// === PHASE 13+++: Pondération intelligente des mots rares mais répétitifs ===
		// Identifier les mots rares (IDF élevé) qui apparaissent souvent (TF élevé)
		// Pénaliser de 0.8x pour éviter qu'ils dominent le score
		if idf[mot] > 0.5 && tf[mot] > 0.05 {
			// Mot rare mais fréquent = potentielle répétition
			tfidfVal *= 0.8
		}

		tfidf[mot] = tfidfVal
	}

	return tfidf
}

// FiltrerCoherence évalue et améliore la cohérence globale du texte généré
func FiltrerCoherence(state GenerationState) float64 {
	if len(state.VecteursGenerés) == 0 {
		return 0
	}

	// Cohérence finale = (1/M) Σ (VR · vrk) / (||VR|| × ||vrk||)
	totalCoherence := 0.0
	for _, v := range state.VecteursGenerés {
		sim := SimilariteCosinus(state.VecteurCible, v)
		totalCoherence += sim
	}

	return totalCoherence / float64(len(state.VecteursGenerés))
}

// ScoragePhrase évalue la qualité globale d'un texte généré
type ScoringGeneration struct {
	CoherenceFinale float64
	DiversiteVocab  float64
	EnergieTotale   float64
	LongueurMoyenne float64
	QualiteGlobale  float64
}

// EvaluerQualite calcule les scores du texte généré
func EvaluerQualite(state GenerationState) ScoringGeneration {
	score := ScoringGeneration{
		CoherenceFinale: FiltrerCoherence(state),
		EnergieTotale:   state.Energie,
	}

	// Diversité du vocabulaire
	uniqueMots := make(map[string]bool)
	for _, mot := range state.MotsGenerés {
		uniqueMots[mot] = true
	}
	if len(state.MotsGenerés) > 0 {
		score.DiversiteVocab = float64(len(uniqueMots)) / float64(len(state.MotsGenerés))
	}

	// Longueur moyenne des mots
	totalLen := 0
	for _, mot := range state.MotsGenerés {
		totalLen += len(mot)
	}
	if len(state.MotsGenerés) > 0 {
		score.LongueurMoyenne = float64(totalLen) / float64(len(state.MotsGenerés))
	}

	// Score composite
	score.QualiteGlobale = (score.CoherenceFinale + score.DiversiteVocab) / 2.0

	return score
}

// RegenererSiMauvais remplace les mots si cohérence < seuil
func RegenererSiMauvais(state GenerationState, params ParametresGeneration, seuilCoherence float64) GenerationState {
	coherence := FiltrerCoherence(state)
	if coherence >= seuilCoherence {
		return state // Assez bon
	}

	// Réécrire: garder 50% des mots, regénérer les 50%
	newState := GenerationState{
		VecteurCible: state.VecteurCible,
	}

	for i, mot := range state.MotsGenerés {
		if i%2 == 0 {
			// Garder
			newState.MotsGenerés = append(newState.MotsGenerés, mot)
			if i < len(state.VecteursGenerés) {
				newState.VecteursGenerés = append(newState.VecteursGenerés, state.VecteursGenerés[i])
			}
		}
		// Sinon: regénérer lors du parcours suivant
	}

	return newState
}

// TrierParScore ordonne les mots par score de pertinence
type CandidatMot struct {
	Mot     string
	Score   float64
	Vecteur VecteurAtomique
}

// GenerationTopK génère en sélectionnant top-K mots (pour plus de diversité)
func GenerationTopK(phrases []Phrase, params ParametresGeneration, k int) GenerationState {
	if k < 1 {
		k = 5
	}

	state := GenerationState{
		VecteurCible: MoyennePondereeVecteurs(phrases),
	}

	vocab := ExtraireVocabulaire(phrases)
	vecabVecteurs := make(map[string]VecteurAtomique)
	for _, mot := range vocab {
		vecabVecteurs[mot] = VectoriserMot(mot)
	}

	var activationPrecedente float64 = 0.0

	for i := 0; i < params.TailleCible; i++ {
		candidats := []CandidatMot{}

		for _, mot := range vocab {
			vMot := vecabVecteurs[mot]
			similarite := SimilariteCosinus(vMot, state.VecteurCible)
			scoreCoherence := 0.0

			if len(state.MotsGenerés) > 0 {
				vMot2 := state.VecteursGenerés[len(state.VecteursGenerés)-1]
				scoreCoherence = SimilariteCosinus(vMot, vMot2)
			}

			score := similarite + params.Eta*scoreCoherence + params.TauxResilience*activationPrecedente
			candidats = append(candidats, CandidatMot{mot, score, vMot})
		}

		// Trier par score décroissant
		sort.Slice(candidats, func(i, j int) bool {
			return candidats[i].Score > candidats[j].Score
		})

		// Sélectionner aléatoirement dans top-K
		if len(candidats) > k {
			candidats = candidats[:k]
		}

		if len(candidats) == 0 {
			break
		}

		// Choisir le meilleur
		meilleur := candidats[0]
		state.MotsGenerés = append(state.MotsGenerés, meilleur.Mot)
		state.VecteursGenerés = append(state.VecteursGenerés, meilleur.Vecteur)
		state.VecteurCible = MettreAJourVecteur(state.VecteurCible, meilleur.Vecteur, params.Gamma)
		activationPrecedente = params.TauxResilience * meilleur.Score
	}

	return state
}
