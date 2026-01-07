package database

import (
	"math"
	"sort"
	"strings"
)

// BlocVectoriel représente un segment du texte avec son vecteur et énergie atomique
type BlocVectoriel struct {
	Index              int             // Position dans le texte
	Contenu            string          // Texte du bloc
	Mots               []string        // Mots du bloc
	Vecteur            VecteurAtomique // Vecteur du bloc
	Coherence          float64         // sim(vi, vj) avec autres blocs
	TFIDFScore         float64         // Importance lexicale
	Energie            float64         // Score final = α·coherence + β·tfidf
	EnergyAtomic       float64         // Ei = α·Σcos(vi,vj) + β·importance
	EtatFreeze         bool            // Vrai si EnergyAtomic < FreezeThreshold
	CoherenceSum       float64         // Σj≠i cos(vi,vj) pour calcul énergie
	CoherenceContexte  float64         // Cohérence avec voisins (avant/après)
	InformationDensite float64         // Densité d'information (mots-clés uniques)
	PenaliteRepetition float64         // Pénalité pour mots répétés
	// === PHASE 13+++: Normalisation lexicale ===
	RepetitionsBloc map[string]int // Compte des répétitions au sein du bloc
}

// ResumeurCoherence découpe le texte et sélectionne les blocs les plus pertinents
type ResumeurCoherence struct {
	Blocs             []BlocVectoriel    // Tous les blocs
	VecteurGlobal     VecteurAtomique    // Vecteur moyen du document
	TailleBlocsMin    int                // Taille min par bloc (10)
	TailleBlocsMax    int                // Taille max par bloc (20)
	AlphaCoherence    float64            // Poids α pour cohérence locale
	BetaTFIDF         float64            // Poids β pour TF-IDF
	BetaContexte      float64            // Poids pour cohérence contextuelle
	GammaInformation  float64            // Poids γ pour densité d'information
	FreezeThreshold   float64            // Seuil ε pour freeze (0.15)
	BigrammeProbs     map[string]float64 // P(wt|wt-1) pour cohérence grammaticale
	RepetitionsGlobal map[string]int     // Compte global des répétitions
}

// NouveauResumeurCoherence crée un résumeur avec découpage amélioré (10-20 mots)
func NouveauResumeurCoherence(tailleMin, tailleMax int) *ResumeurCoherence {
	if tailleMin < 5 {
		tailleMin = 5
	}
	if tailleMax < tailleMin+5 {
		tailleMax = tailleMin + 5
	}
	return &ResumeurCoherence{
		Blocs:             []BlocVectoriel{},
		TailleBlocsMin:    tailleMin, // 10 mots min
		TailleBlocsMax:    tailleMax, // 20 mots max
		AlphaCoherence:    0.4,       // Poids cohérence locale
		BetaTFIDF:         0.3,       // Poids importance lexicale
		BetaContexte:      0.2,       // Poids cohérence contextuelle
		GammaInformation:  0.1,       // Poids densité d'information
		FreezeThreshold:   0.15,      // Seuil freeze (Ei < 0.15)
		BigrammeProbs:     make(map[string]float64),
		RepetitionsGlobal: make(map[string]int),
	}
}

// countUniqueWords compte le nombre de mots uniques dans une liste
func countUniqueWords(mots []string) int {
	unique := make(map[string]bool)
	for _, mot := range mots {
		unique[strings.ToLower(mot)] = true
	}
	return len(unique)
}

// Decouper divise le texte en blocs optimisés (10-20 mots) et calcule l'énergie atomique
func (r *ResumeurCoherence) Decouper(phrases []Phrase) {
	r.Blocs = []BlocVectoriel{}

	// Fusionner les phrases en blocs de 10-20 mots
	blocsRaw := [][]Phrase{}
	blocActuel := []Phrase{}
	motsDansBloc := 0

	for _, p := range phrases {
		motsDansBloc += len(p.Mots)
		blocActuel = append(blocActuel, p)

		// Si on dépasse TailleBlocsMax, créer un nouveau bloc
		if motsDansBloc >= r.TailleBlocsMax {
			blocsRaw = append(blocsRaw, blocActuel)
			blocActuel = []Phrase{}
			motsDansBloc = 0
		}
	}
	if len(blocActuel) > 0 {
		blocsRaw = append(blocsRaw, blocActuel)
	}

	// Vectoriser chaque bloc
	allMots := []string{}
	for _, bloc := range blocsRaw {
		mots := []string{}
		for _, p := range bloc {
			mots = append(mots, p.Mots...)
		}
		allMots = append(allMots, mots...)
	}
	tfidf := CalculerTFIDF(phrases, allMots)

	for idx, bloc := range blocsRaw {
		mots := []string{}
		contenu := ""
		for _, p := range bloc {
			mots = append(mots, p.Mots...)
			contenu += p.Contenu + " "
		}
		contenu = strings.TrimSpace(contenu)

		// Vectoriser le bloc
		vecteur := MoyennePondereeVecteurs(bloc)

		// TF-IDF moyen du bloc
		tfidfMoyen := 0.0
		for _, mot := range mots {
			tfidfMoyen += tfidf[mot]
		}
		if len(mots) > 0 {
			tfidfMoyen /= float64(len(mots))
		}

		b := BlocVectoriel{
			Index:      idx,
			Contenu:    contenu,
			Mots:       mots,
			Vecteur:    vecteur,
			TFIDFScore: tfidfMoyen,
		}
		r.Blocs = append(r.Blocs, b)
	}

	// === PHASE 13+++: Normalisation lexicale - compter répétitions par bloc ===
	r.NormaliserRepetitionsBlocs()

	// Calculer la cohérence entre blocs
	r.CalculerCoherenceBlocs()

	// Calculer l'énergie atomique et le freeze
	r.CalculerEnergieAtomique()

	// Calculer le vecteur global à partir des blocs non-gelés
	r.CalculerVecteurGlobal()

	// Calculer les probabilités de bigrammes
	r.CalculerBigrammeProbs(phrases)

	// Calculer la cohérence contextuelle (avec voisins)
	r.CalculerCoherenceContexte()

	// Calculer la densité d'information
	r.CalculerInformationDensite()
}

// CalculerCoherenceContexte calcule la cohérence avec les voisins immédiats
func (r *ResumeurCoherence) CalculerCoherenceContexte() {
	for i := range r.Blocs {
		sommeCoherence := 0.0
		count := 0

		// Voisin avant (i-1)
		if i > 0 {
			sim := SimilariteCosinus(r.Blocs[i].Vecteur, r.Blocs[i-1].Vecteur)
			sommeCoherence += sim
			count++
		}

		// Voisin après (i+1)
		if i < len(r.Blocs)-1 {
			sim := SimilariteCosinus(r.Blocs[i].Vecteur, r.Blocs[i+1].Vecteur)
			sommeCoherence += sim
			count++
		}

		if count > 0 {
			r.Blocs[i].CoherenceContexte = sommeCoherence / float64(count)
		}
	}
}

// CalculerInformationDensite calcule la densité d'information (mots-clés uniques)
func (r *ResumeurCoherence) CalculerInformationDensite() {
	for i := range r.Blocs {
		// Compter les mots uniques non-stopwords
		unique := make(map[string]bool)
		for _, mot := range r.Blocs[i].Mots {
			motNorm := strings.ToLower(mot)
			if !StopWords[motNorm] && len(motNorm) > 2 && isValidWord(motNorm) {
				unique[motNorm] = true
			}
		}
		// Densité = mots uniques / total
		if len(r.Blocs[i].Mots) > 0 {
			r.Blocs[i].InformationDensite = float64(len(unique)) / float64(len(r.Blocs[i].Mots))
		}
	}
}

// CalculerRepetitionPenalty calcule la pénalité basée sur les mots déjà utilisés
func (r *ResumeurCoherence) CalculerRepetitionPenalty(resumeEnCours []string) {
	// Compter la fréquence des mots dans le résumé en cours
	motCount := make(map[string]int)
	for _, mot := range resumeEnCours {
		motNorm := strings.ToLower(mot)
		motCount[motNorm]++
	}

	// Pour chaque bloc, calculer la pénalité
	for i := range r.Blocs {
		penalty := 0.0
		for _, mot := range r.Blocs[i].Mots {
			motNorm := strings.ToLower(mot)
			if motCount[motNorm] > 0 {
				// Pénalité = λ · f(wi) où f(wi) = fréquence du mot dans le résumé
				penalty += float64(motCount[motNorm])
			}
		}
		r.Blocs[i].PenaliteRepetition = penalty / float64(len(r.Blocs[i].Mots)+1)
	}
}

// CalculerCoherenceBlocs calcule cos(vi, vj) pour chaque paire
func (r *ResumeurCoherence) CalculerCoherenceBlocs() {
	for i := range r.Blocs {
		sumCoherence := 0.0
		for j := range r.Blocs {
			if i != j {
				sim := SimilariteCosinus(r.Blocs[i].Vecteur, r.Blocs[j].Vecteur)
				if sim < 0 {
					sim = 0
				}
				if sim > 1 {
					sim = 1
				}
				sumCoherence += sim
			}
		}
		r.Blocs[i].CoherenceSum = sumCoherence
		r.Blocs[i].Coherence = sumCoherence / float64(len(r.Blocs)-1)
	}
}

// CalculerEnergieAtomique calcule Ei = α·Σcos + β·tfidf et freeze
func (r *ResumeurCoherence) CalculerEnergieAtomique() {
	for i := range r.Blocs {
		// Ei = α·Σj≠i cos(vi,vj) + β·importance_lexicale
		r.Blocs[i].EnergyAtomic = r.AlphaCoherence*r.Blocs[i].CoherenceSum +
			r.BetaTFIDF*r.Blocs[i].TFIDFScore

		// Pénaliser les répétitions au sein du bloc: Ei_corrigé = Ei / (1 + λ·répétitions)
		repetitions := float64(len(r.Blocs[i].Mots) - countUniqueWords(r.Blocs[i].Mots))
		lambda := 0.3 // Facteur de pénalisation
		r.Blocs[i].EnergyAtomic = r.Blocs[i].EnergyAtomic / (1.0 + lambda*repetitions)

		// Caper la cohérence à 95% max pour éviter que un bloc domine complètement
		if r.Blocs[i].Coherence > 0.95 {
			r.Blocs[i].Coherence = 0.95
		}

		// Freeze si Ei < ε (seuil 0.15)
		r.Blocs[i].EtatFreeze = r.Blocs[i].EnergyAtomic < r.FreezeThreshold

		// Aussi utiliser l'ancienne formule pour sélection
		r.Blocs[i].Energie = r.Blocs[i].EnergyAtomic
	}
}

// CalculerVecteurGlobal calcule le vecteur moyen du document (uniquement blocs non-gelés)
// Vdoc = Σ(Ei·vi) / Σ(Ei) pour blocs où Ei >= FreezeThreshold
func (r *ResumeurCoherence) CalculerVecteurGlobal() {
	if len(r.Blocs) == 0 {
		r.VecteurGlobal = VecteurAtomique{}
		return
	}

	// Filtrer les blocs non-gelés
	blocsActifs := []BlocVectoriel{}
	for _, b := range r.Blocs {
		if !b.EtatFreeze {
			blocsActifs = append(blocsActifs, b)
		}
	}

	// Si tous gelés, utiliser quand même les moins gelés
	if len(blocsActifs) == 0 {
		blocsActifs = r.Blocs
	}

	// Moyenne pondérée par l'énergie atomique
	somme := make([]float64, 11)
	sumEnergie := 0.0

	for _, b := range blocsActifs {
		weight := 1.0
		if b.EnergyAtomic > 0 {
			weight = b.EnergyAtomic
		}
		sumEnergie += weight

		for i := 0; i < 11; i++ {
			somme[i] += weight * b.Vecteur.Dimensions[i]
		}
	}

	// Normaliser par la somme des poids
	if sumEnergie > 0 {
		for i := 0; i < 11; i++ {
			somme[i] /= sumEnergie
		}
	}

	// Normaliser le vecteur global
	norm := 0.0
	for i := 0; i < 11; i++ {
		norm += somme[i] * somme[i]
	}
	norm = math.Sqrt(norm)
	if norm > 0 {
		for i := 0; i < 11; i++ {
			somme[i] /= norm
		}
	}

	r.VecteurGlobal.Dimensions = somme
	r.VecteurGlobal.Energie = 0.8
}

// CalculerBigrammeProbs calcule P(wt|wt-1) pour cohérence grammaticale
func (r *ResumeurCoherence) CalculerBigrammeProbs(phrases []Phrase) {
	r.BigrammeProbs = make(map[string]float64)

	// Compter les bigrammes
	bigrammeCounts := make(map[string]int)
	prevWordCounts := make(map[string]int)

	for _, p := range phrases {
		for i := 0; i < len(p.Mots)-1; i++ {
			w1 := strings.ToLower(p.Mots[i])
			w2 := strings.ToLower(p.Mots[i+1])
			bigram := w1 + "|" + w2

			// Filtrer les mots valides
			if !isValidWord(w1) || !isValidWord(w2) || StopWords[w1] || StopWords[w2] {
				continue
			}

			bigrammeCounts[bigram]++
			prevWordCounts[w1]++
		}
	}

	// Calculer les probabilités P(wt|wt-1)
	for bigram, count := range bigrammeCounts {
		parts := strings.Split(bigram, "|")
		if len(parts) == 2 {
			w1 := parts[0]
			if prevWordCounts[w1] > 0 {
				prob := float64(count) / float64(prevWordCounts[w1])
				r.BigrammeProbs[bigram] = prob
			}
		}
	}
}

// Selectionner retourne les meilleurs blocs selon l'énergie
func (r *ResumeurCoherence) Selectionner(ratioCompression float64) []BlocVectoriel {
	// Nombre de blocs à garder: sélectionner plus de blocs pour éviter la répétition
	// Au lieu de ceil, utiliser un minimum de 3-4 blocs pour diversité suffisante
	numBlocs := int(math.Ceil(float64(len(r.Blocs)) * ratioCompression))

	// Forcer minimum 3 blocs sauf si moins de 3 blocs disponibles
	if numBlocs < 3 && len(r.Blocs) >= 3 {
		numBlocs = 3 // Forcer au minimum 3 blocs pour diversité
	} else if numBlocs < 2 && len(r.Blocs) >= 2 {
		numBlocs = 2
	} else if numBlocs < 1 {
		numBlocs = 1
	}
	if numBlocs > len(r.Blocs) {
		numBlocs = len(r.Blocs)
	}

	// Trier par énergie
	selected := make([]BlocVectoriel, len(r.Blocs))
	copy(selected, r.Blocs)
	sort.Slice(selected, func(i, j int) bool {
		return selected[i].Energie > selected[j].Energie
	})

	// Prendre les top N
	selected = selected[:numBlocs]

	// Réordonner selon l'index pour préserver l'ordre du texte original
	sort.Slice(selected, func(i, j int) bool {
		return selected[i].Index < selected[j].Index
	})

	return selected
}

// GenerByVectors génère un résumé à partir des blocs sélectionnés
func (r *ResumeurCoherence) GenerByVectors(blocsSelectionnes []BlocVectoriel, tailleResume int) string {
	// Créer un super-vecteur cible à partir des blocs sélectionnés
	vetorCible := MoyennePondereeVecteurs([]Phrase{}) // Vide au départ
	for _, b := range blocsSelectionnes {
		vetorCible.Dimensions = b.Vecteur.Dimensions
	}

	// Normaliser
	norm := 0.0
	for i := 0; i < 11; i++ {
		norm += vetorCible.Dimensions[i] * vetorCible.Dimensions[i]
	}
	norm = math.Sqrt(norm)
	if norm > 0 {
		for i := 0; i < 11; i++ {
			vetorCible.Dimensions[i] /= norm
		}
	}

	// Extraire tous les mots des blocs sélectionnés
	vocab := []string{}
	vocabSet := make(map[string]bool)
	for _, b := range blocsSelectionnes {
		for _, mot := range b.Mots {
			if !vocabSet[mot] && isValidWord(mot) && !StopWords[mot] && len(mot) > 2 {
				vocab = append(vocab, mot)
				vocabSet[mot] = true
			}
		}
	}

	// Générer le résumé
	params := ParametresGeneration{
		TailleCible:    tailleResume,
		Eta:            0.4,
		Gamma:          0.15,
		LongueurMin:    20,
		DiversiteMin:   0.3,
		TauxResilience: 0.5,
	}

	state := GenerationState{
		VecteurCible:    vetorCible,
		MotsGenerés:     []string{},
		VecteursGenerés: []VecteurAtomique{},
		Activations:     []float64{},
		Coherences:      []float64{},
	}

	// Vectoriser le vocabulaire
	vocabVecteurs := make(map[string]VecteurAtomique)
	for _, mot := range vocab {
		vocabVecteurs[mot] = VectoriserMot(mot)
	}

	motsPrecedents := make(map[string]int)
	var activationPrecedente float64 = 0.0

	// Boucle de génération
	for k := 0; k < params.TailleCible && len(vocab) > 0; k++ {
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

			scoreCoherence := 0.0
			if len(state.MotsGenerés) > 0 {
				vMot2 := state.VecteursGenerés[len(state.VecteursGenerés)-1]
				scoreCoherence = SimilariteCosinus(vMot, vMot2)
			}

			score := similarite + params.Eta*scoreCoherence +
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

		activationPrecedente = params.TauxResilience * meilleurActivation
		motsPrecedents[meilleurMot]++
	}

	texte := strings.Join(state.MotsGenerés, " ")
	return AjouterConnecteurs(texte, state)
}

// scoreWord pour la génération avec bigrammes
type scoreWord struct {
	mot   string
	score float64
}

// === PHASE 13: Advanced Normalization & Sliding Window Vectorization ===

// CalculerRepetitionsBloc compte le nombre de mots répétés dans un bloc
func CalculerRepetitionsBloc(bloc BlocVectoriel) int {
	wordCount := make(map[string]int)
	for _, mot := range bloc.Mots {
		motNorm := strings.ToLower(mot)
		if isValidWord(motNorm) && !StopWords[motNorm] && len(motNorm) > 2 {
			wordCount[motNorm]++
		}
	}

	repetitions := 0
	for _, count := range wordCount {
		if count > 1 {
			repetitions += count - 1
		}
	}
	return repetitions
}

// === PHASE 13+++: Normalisation lexicale dans les blocs ===

// NormaliserRepetitionsBlocs compte les répétitions au sein de chaque bloc
func (r *ResumeurCoherence) NormaliserRepetitionsBlocs() {
	for i := range r.Blocs {
		r.Blocs[i].RepetitionsBloc = make(map[string]int)

		for _, mot := range r.Blocs[i].Mots {
			motNorm := strings.ToLower(mot)
			if isValidWord(motNorm) && !StopWords[motNorm] && len(motNorm) > 2 {
				r.Blocs[i].RepetitionsBloc[motNorm]++
			}
		}

		// Calculer la pénalité de répétition = Σ(repetitions > 2)
		penalite := 0.0
		for mot, count := range r.Blocs[i].RepetitionsBloc {
			if count > 2 {
				// Pénalité progressive: 0.1 par répétition excédentaire
				penalite += float64(count-2) * 0.1
			}
			// Log pour debug: mots très répétés
			if count > 3 && len(mot) > 4 {
				// Mot répété >3x dans le bloc (potentiellement artifact)
			}
		}
		r.Blocs[i].PenaliteRepetition = penalite
	}
}

// NormaliserScore applique la normalisation de récompense: Score / (1 + λ·Répétitions_ratio)
func NormaliserScore(score float64, repetitionsRatio float64, lambda float64) float64 {
	if lambda == 0 {
		lambda = 0.5
	}
	return score / (1.0 + lambda*repetitionsRatio)
}

// CalculerFenetreGlissante calcule les vecteurs pour fenêtres glissantes
// de k blocs successifs et retourne une map des vecteurs par position
func (r *ResumeurCoherence) CalculerFenetreGlissante(tailleWindow int) map[int]VecteurAtomique {
	vecteursWindow := make(map[int]VecteurAtomique)

	if len(r.Blocs) == 0 || tailleWindow <= 0 {
		return vecteursWindow
	}

	// Pour chaque position de début de fenêtre
	for start := 0; start < len(r.Blocs); start++ {
		end := start + tailleWindow
		if end > len(r.Blocs) {
			end = len(r.Blocs)
		}

		// Moyenne pondérée des vecteurs dans la fenêtre
		windowVect := VecteurAtomique{
			Dimensions: make([]float64, 11),
		}

		sumWeight := 0.0
		for i := start; i < end; i++ {
			weight := r.Blocs[i].EnergyAtomic
			if weight <= 0 {
				weight = 1.0
			}
			for j := 0; j < 11; j++ {
				windowVect.Dimensions[j] += weight * r.Blocs[i].Vecteur.Dimensions[j]
			}
			sumWeight += weight
		}

		// Normaliser
		if sumWeight > 0 {
			for i := 0; i < 11; i++ {
				windowVect.Dimensions[i] /= sumWeight
			}
		}

		// Normaliser le vecteur
		norm := 0.0
		for i := 0; i < 11; i++ {
			norm += windowVect.Dimensions[i] * windowVect.Dimensions[i]
		}
		norm = math.Sqrt(norm)
		if norm > 0 {
			for i := 0; i < 11; i++ {
				windowVect.Dimensions[i] /= norm
			}
		}

		vecteursWindow[start] = windowVect
	}

	return vecteursWindow
}

// SelectionnerBlocsAvecFenetrageGlissant sélectionne les blocs en utilisant
// la cohérence locale par rapport aux fenêtres glissantes
func (r *ResumeurCoherence) SelectionnerBlocsAvecFenetrageGlissant(ratio float64, tailleWindow int) []BlocVectoriel {
	if len(r.Blocs) == 0 || ratio <= 0 {
		return []BlocVectoriel{}
	}

	if tailleWindow <= 0 {
		tailleWindow = 15 // Par défaut
	}

	// Calculer les fenêtres glissantes
	vecteursWindow := r.CalculerFenetreGlissante(tailleWindow)

	// Calculer le score de chaque bloc par rapport à sa fenêtre locale
	blocsAvecScore := make([]struct {
		bloc  BlocVectoriel
		score float64
	}, 0)

	for i, bloc := range r.Blocs {
		// Trouver la fenêtre contenant ce bloc
		windowStart := i
		if windowStart > len(r.Blocs)-tailleWindow {
			windowStart = len(r.Blocs) - tailleWindow
		}
		if windowStart < 0 {
			windowStart = 0
		}

		windowVect, exists := vecteursWindow[windowStart]
		if !exists {
			continue
		}

		// Cohérence locale: similarité avec la fenêtre
		localCoherence := SimilariteCosinus(bloc.Vecteur, windowVect)

		// Cohérence globale (avec le vecteur document)
		globalCoherence := SimilariteCosinus(bloc.Vecteur, r.VecteurGlobal)

		// Multi-level: α·locale + (1-α)·globale (α≈0.6 privileges local)
		alpha := 0.6
		multiLevelCoherence := alpha*localCoherence + (1-alpha)*globalCoherence

		// === PHASE 13+++: Pénalité de répétition du bloc ===
		// Pénalité de répétition: score / (1 + λ·repetitions_ratio)
		repetitions := CalculerRepetitionsBloc(bloc)
		repetitionsRatio := 0.0
		if len(bloc.Mots) > 0 {
			repetitionsRatio = float64(repetitions) / float64(len(bloc.Mots))
		}

		normalizedScore := NormaliserScore(multiLevelCoherence, repetitionsRatio, 0.5)

		// Ajouter l'énergie atomique comme multiplicateur
		// Aussi appliquer la pénalité de répétition du bloc
		finalScore := normalizedScore * bloc.EnergyAtomic * (1.0 - bloc.PenaliteRepetition)

		blocsAvecScore = append(blocsAvecScore, struct {
			bloc  BlocVectoriel
			score float64
		}{bloc, finalScore})
	}

	// Trier par score décroissant
	sort.Slice(blocsAvecScore, func(i, j int) bool {
		return blocsAvecScore[i].score > blocsAvecScore[j].score
	})

	// Sélectionner minimum 3-5 blocs ou ratio*len(blocs)
	// === PHASE 13+: Augmenter limite à 40-50 pour plus de contenu ===
	numBlocs := int(math.Ceil(ratio * float64(len(r.Blocs))))
	if numBlocs < 5 {
		numBlocs = 5 // Min 5 blocs
	}
	// Limiter à max 40-50 blocs pour documents longs (vs 20 avant)
	if numBlocs > 50 {
		numBlocs = 50
	}
	if numBlocs > len(blocsAvecScore) {
		numBlocs = len(blocsAvecScore)
	}

	// Conserver dans l'ordre original
	selectedIndices := make(map[int]bool)
	for i := 0; i < numBlocs && i < len(blocsAvecScore); i++ {
		for j, b := range r.Blocs {
			if b.Index == blocsAvecScore[i].bloc.Index {
				selectedIndices[j] = true
				break
			}
		}
	}

	// === PHASE 13+++: Fenêtrage strict - Forcer diversité lexicale entre blocs consécutifs ===
	// Vérifier que blocs consécutifs n'ont pas >60% de vocabulaire identique
	result := []BlocVectoriel{}
	for i, bloc := range r.Blocs {
		if selectedIndices[i] {
			// Vérifier similarity avec dernier bloc sélectionné
			if len(result) > 0 {
				lastBloc := result[len(result)-1]
				similarity := CalculerSimilarityVocabLexical(lastBloc.Mots, bloc.Mots)

				// Si >60% similarité, skip ce bloc
				if similarity > 0.6 {
					delete(selectedIndices, i)
					continue
				}
			}

			result = append(result, bloc)
		}
	}

	return result
}

// GenerByVectorsWithGrammar génère un résumé avec cohérence grammaticale (bigrammes)
func (r *ResumeurCoherence) GenerByVectorsWithGrammar(blocsSelectionnes []BlocVectoriel, tailleResume int) string {
	// Vecteur cible pondéré
	vetorCible := MoyennePondereeVecteurs([]Phrase{})
	sumEnergie := 0.0
	for i := 0; i < 11; i++ {
		vetorCible.Dimensions[i] = 0
	}

	for _, b := range blocsSelectionnes {
		weight := 1.0
		if b.EnergyAtomic > 0 {
			weight = b.EnergyAtomic
		}
		sumEnergie += weight
		for i := 0; i < 11; i++ {
			vetorCible.Dimensions[i] += weight * b.Vecteur.Dimensions[i]
		}
	}

	if sumEnergie > 0 {
		for i := 0; i < 11; i++ {
			vetorCible.Dimensions[i] /= sumEnergie
		}
	}

	// Normaliser
	norm := 0.0
	for i := 0; i < 11; i++ {
		norm += vetorCible.Dimensions[i] * vetorCible.Dimensions[i]
	}
	norm = math.Sqrt(norm)
	if norm > 0 {
		for i := 0; i < 11; i++ {
			vetorCible.Dimensions[i] /= norm
		}
	}

	// === PHASE 13++: Extraire vocabulaire de TOUS les blocs, pas juste sélectionnés ===
	// Cela donne 10x plus de vocabulaire disponible
	vocab := []string{}
	vocabSet := make(map[string]bool)
	for _, b := range r.Blocs { // Utiliser tous les blocs du document
		for _, mot := range b.Mots {
			motNorm := strings.ToLower(mot)
			// Inclure les mots d'au moins 2 caractères
			if !vocabSet[motNorm] && isValidWord(motNorm) && !StopWords[motNorm] && len(motNorm) > 1 {
				vocab = append(vocab, motNorm)
				vocabSet[motNorm] = true
			}
		}
	}

	// Si vocab vide, retourner un texte simple
	if len(vocab) == 0 {
		return ""
	}

	// Vectoriser vocabulaire
	vocabVecteurs := make(map[string]VecteurAtomique)
	for _, mot := range vocab {
		vocabVecteurs[mot] = VectoriserMot(mot)
	}

	// Génération mot par mot avec TF-IDF dynamique et cohérence multi-niveau
	mots := []string{}
	prevWord := ""
	motsDéjàUtilisés := make(map[string]int) // Compte les répétitions dans ce résumé

	// === PHASE 13: Comptage dynamique pour TF-IDF ===
	comptageMotsResume := make(map[string]int) // Compte les mots dans le résumé en cours

	for k := 0; k < tailleResume && len(vocab) > 0; k++ {
		meilleurScore := -math.Inf(1)
		meilleurMot := ""

		for _, mot := range vocab {
			// === PHASE 13+: Ban strict (1 seule occurrence) pour éviter doublons ===
			if motsDéjàUtilisés[mot] >= 1 {
				continue // Ignorer ce mot après 1ère occurrence
			}

			// === PHASE 13: TF-IDF dynamique ===
			// Score_word = Score_initial × 1/(1 + 0.3·f(w)) [réduit de 1.0 à 0.3 pour moins d'agressivité]
			// où f(w) = compte actuel du mot w dans le résumé en cours
			tfIDFDynamic := 1.0
			if comptageMotsResume[mot] > 0 {
				// Coefficient réduit de 1.0 à 0.3 pour permettre plus de mots
				tfIDFDynamic = 1.0 / (1.0 + 0.3*float64(comptageMotsResume[mot]))
			}

			// Similarité avec le vecteur cible
			sim := SimilariteCosinus(vocabVecteurs[mot], vetorCible)

			// Appliquer TF-IDF dynamique: réduire progressivement le score
			simAjustee := sim * tfIDFDynamic

			// Pénalité si mot déjà utilisé - N/A car ban strict avant ce code
			repetitionPenalty := 1.0

			simAjustee *= repetitionPenalty

			// Probabilité grammaticale P(wt|wt-1)
			gramProb := 0.0
			if prevWord != "" {
				bigram := strings.ToLower(prevWord) + "|" + mot
				if prob, exists := r.BigrammeProbs[bigram]; exists {
					gramProb = prob
				} else {
					gramProb = 0.1 // Pénalité si bigramme rare
				}
			}

			// Bonus de diversité: si mot pas encore utilisé, +40%
			diversityBonus := 1.0
			if motsDéjàUtilisés[mot] == 0 {
				diversityBonus = 1.4 // +40% si nouveau
			}

			// Score composite avec multi-level coherence
			// Alpha=0.6 privilégie la similarité locale (sim), 1-alpha=0.4 pour bigramme
			alpha := 0.6
			score := (alpha*simAjustee + (1-alpha)*(1.0+0.3*gramProb)) * diversityBonus

			// Si pas assez de similarité, utiliser gramProb comme score
			if simAjustee < 0.001 {
				score = gramProb * 10 // Boost pour que ce soit compétitif
			}

			if score > meilleurScore {
				meilleurScore = score
				meilleurMot = mot
			}
		}

		if meilleurScore < 0 || meilleurMot == "" {
			break
		}

		mots = append(mots, meilleurMot)
		motsDéjàUtilisés[meilleurMot]++
		comptageMotsResume[meilleurMot]++ // === PHASE 13: Incrémenter le comptage dynamique ===
		prevWord = meilleurMot
		r.RepetitionsGlobal[meilleurMot]++
	}

	texte := strings.Join(mots, " ")

	// Créer un GenerationState complet pour AjouterConnecteurs
	state := GenerationState{
		MotsGenerés: mots,
		Coherences:  make([]float64, len(mots)),
	}

	// Remplir les cohérences pour chaque mot généré
	for i, mot := range mots {
		vword := VectoriserMot(mot)
		coh := SimilariteCosinus(vword, r.VecteurGlobal)
		if coh > 0.01 {
			state.Coherences[i] = coh
		}
	}

	return AjouterConnecteurs(texte, state)
}

// === PHASE 13+++: CalculerSimilarityVocabLexical - Mesurer similarité lexicale entre deux listes de mots ===
// Retourne le ratio de mots communs (intersection / union)
func CalculerSimilarityVocabLexical(mots1, mots2 []string) float64 {
	if len(mots1) == 0 || len(mots2) == 0 {
		return 0.0
	}

	// Créer des maps pour normalisation
	vocab1 := make(map[string]bool)
	vocab2 := make(map[string]bool)

	for _, mot := range mots1 {
		motClean := strings.ToLower(strings.TrimSuffix(strings.TrimSuffix(mot, "."), ","))
		if len(motClean) > 2 && !StopWords[motClean] {
			vocab1[motClean] = true
		}
	}

	for _, mot := range mots2 {
		motClean := strings.ToLower(strings.TrimSuffix(strings.TrimSuffix(mot, "."), ","))
		if len(motClean) > 2 && !StopWords[motClean] {
			vocab2[motClean] = true
		}
	}

	// Calculer intersection et union
	intersection := 0
	for mot := range vocab1 {
		if vocab2[mot] {
			intersection++
		}
	}

	union := len(vocab1) + len(vocab2) - intersection
	if union == 0 {
		return 0.0
	}

	similarity := float64(intersection) / float64(union)
	return similarity
}

// ResumerAvecDecoupage est la fonction principale - version améliorée (10-20 mots par bloc)
func ResumerAvecDecoupage(phrases []Phrase, ratioCompression float64) ResumeurCoherence {
	resumeur := NouveauResumeurCoherence(10, 20) // Blocs de 10-20 mots
	resumeur.Decouper(phrases)
	return *resumeur
}
