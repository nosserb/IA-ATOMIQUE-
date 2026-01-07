package database

import (
	"math"
	"sort"
	"strings"
)

// AtomeGeneration représente un atome autonome pour la génération de texte
// Chaque atome = une idée/concept avec son propre état d'activation
type AtomeGeneration struct {
	ID            int
	Concept       string          // Le mot/concept qu'il représente
	Activation    float64         // aj(t) - state [0,1]
	Voisins       map[int]float64 // Couplages αjk vers autres atomes
	Categorie     int             // Catégorie lexicale
	Style         string          // Style sémantique (tech, narrative, formal, etc)
	Intensite     float64         // Intensité sémantique [0,1]
	Connecteur    string          // Connecteur optionnel ("En outre", "Cependant", etc)
	Energie       float64         // Énergie dépensée
	EtatFreeze    bool            // Frozen state
	IterationsGel int             // Itérations depuis freeze
}

// ReseauAtomiquePourGeneration réseau distribué pour génération
type ReseauAtomiquePourGeneration struct {
	Atomes           map[int]*AtomeGeneration
	VecteurCible     VecteurAtomique // VR = idée centrale
	Couplage         float64         // α par défaut
	InfluenceCible   float64         // β - influence de VR
	SeuilResonance   float64         // σ_wake - seuil d'émission
	TauxDecroissance float64         // Décroissance naturelle de l'activation
	FreezeThreshold  float64         // Activation threshold pour freeze
	IterationID      int             // Compteur d'itérations
	MotsEmis         []string        // Mots générés
	AtoméActivesEmis []int           // IDs des atomes qui ont émis
}

// NouveauReseauGeneration crée un réseau atomique pour génération
func NouveauReseauGeneration(vecteurCible VecteurAtomique, phrases []Phrase) *ReseauAtomiquePourGeneration {
	reseau := &ReseauAtomiquePourGeneration{
		Atomes:           make(map[int]*AtomeGeneration),
		VecteurCible:     vecteurCible,
		Couplage:         0.25, // α (réduit pour moins de bruit)
		InfluenceCible:   0.75, // β (augmenté pour plus de pertinence)
		SeuilResonance:   0.70, // σ_wake (augmenté pour moins d'émissions)
		TauxDecroissance: 0.15, // Décroissance plus rapide
		FreezeThreshold:  0.20,
		IterationID:      0,
		MotsEmis:         []string{},
	}

	// Créer les atomes à partir du vocabulaire
	vocab := ExtraireVocabulaire(phrases)
	for i, mot := range vocab {
		atome := &AtomeGeneration{
			ID:            i,
			Concept:       mot,
			Activation:    0.3 + (0.2 * math.Sin(float64(i))), // Variation initiale
			Voisins:       make(map[int]float64),
			Style:         DeterminStyles(mot),
			Intensite:     CalculerIntensiteSemantique(mot),
			Connecteur:    "",
			Energie:       0.0,
			EtatFreeze:    false,
			IterationsGel: 0,
		}

		// Assigner catégorie
		if word, ok := Words[strings.ToLower(mot)]; ok {
			atome.Categorie = word.Categorie
		}

		reseau.Atomes[i] = atome
	}

	// Initialiser les couplages (voisinage basé sur similarité sémantique)
	reseau.InitialiserCouplages()

	return reseau
}

// InitialiserCouplages établit les connexions entre atomes voisins
func (r *ReseauAtomiquePourGeneration) InitialiserCouplages() {
	atomesList := make([]*AtomeGeneration, 0, len(r.Atomes))
	for _, atome := range r.Atomes {
		atomesList = append(atomesList, atome)
	}

	// Pour chaque atome, connecter aux k voisins sémantiquement proches
	const kVoisins = 5
	for i, atome1 := range atomesList {
		// Calculer similarités
		type SimilariteAtome struct {
			id  int
			sim float64
		}
		similarities := make([]SimilariteAtome, 0)

		v1 := VectoriserMot(atome1.Concept)
		for j, atome2 := range atomesList {
			if i == j {
				continue
			}
			v2 := VectoriserMot(atome2.Concept)
			sim := SimilariteCosinus(v1, v2)
			if sim > 0.1 { // Seuil minimum
				similarities = append(similarities, SimilariteAtome{j, sim})
			}
		}

		// Trier et prendre les k meilleurs
		sort.Slice(similarities, func(a, b int) bool {
			return similarities[a].sim > similarities[b].sim
		})

		for j := 0; j < len(similarities) && j < kVoisins; j++ {
			sim := similarities[j].sim
			atome1.Voisins[similarities[j].id] = 0.2 + (0.8 * sim) // α_jk ∈ [0.2, 1.0]
		}
	}
}

// IterationAtomique exécute une itération de résonance
// Mise à jour: aj(t+1) = f(aj(t), Σ αjk*ak(t), βj*VR)
func (r *ReseauAtomiquePourGeneration) IterationAtomique() []string {
	r.IterationID++
	nouvellsActivations := make(map[int]float64)
	motsCettIteration := []string{}

	// Copier état courant
	etatCourant := make(map[int]float64)
	for id, atome := range r.Atomes {
		etatCourant[id] = atome.Activation
	}

	// Mettre à jour chaque atome
	for id, atome := range r.Atomes {
		if atome.EtatFreeze {
			// Atome gelé: activation très lente
			nouvellsActivations[id] = atome.Activation * 0.95
			atome.IterationsGel++

			// Dégel si isolation terminée
			if atome.IterationsGel > 10 {
				atome.EtatFreeze = false
				atome.IterationsGel = 0
			}
			continue
		}

		// Influence des voisins: Σ αjk*ak(t)
		influenceVoisins := 0.0
		for voisinID, couplage := range atome.Voisins {
			influenceVoisins += couplage * etatCourant[voisinID]
		}

		// Influence du vecteur cible: βj*VR
		vAtome := VectoriserMot(atome.Concept)
		influenceCible := r.InfluenceCible * SimilariteCosinus(vAtome, r.VecteurCible)

		// Fonction de mise à jour sigmoïde
		aj := etatCourant[id]
		input := aj + influenceVoisins + influenceCible
		nouvelleActivation := Sigmoide(input)

		// Appliquer décroissance naturelle
		nouvelleActivation *= (1.0 - r.TauxDecroissance)

		nouvellsActivations[id] = nouvelleActivation

		// Vérifier si atome dépasse le seuil de résonance σ_wake
		if nouvelleActivation >= r.SeuilResonance {
			motsCettIteration = append(motsCettIteration, atome.Concept)
			r.AtoméActivesEmis = append(r.AtoméActivesEmis, id)

			// Geler temporairement pour éviter répétition immédiate
			atome.EtatFreeze = true
			atome.IterationsGel = 0
			nouvelleActivation = 0.2 // Reset après émission

			// Consommer de l'énergie
			atome.Energie += nouvelleActivation
		}

		// Gel si activation trop basse
		if nouvelleActivation < r.FreezeThreshold {
			atome.EtatFreeze = true
		}

		atome.Activation = nouvelleActivation
	}

	r.MotsEmis = append(r.MotsEmis, motsCettIteration...)
	return motsCettIteration
}

// Sigmoide fonction d'activation
func Sigmoide(x float64) float64 {
	return 1.0 / (1.0 + math.Exp(-x))
}

// DeterminStyles détermine le style sémantique d'un mot
func DeterminStyles(mot string) string {
	lower := strings.ToLower(mot)

	// Mots techniques
	if strings.Contains(lower, "algo") || strings.Contains(lower, "code") ||
		strings.Contains(lower, "systeme") || strings.Contains(lower, "data") {
		return "technical"
	}

	// Mots narratifs/story
	if strings.Contains(lower, "histoire") || strings.Contains(lower, "jour") ||
		strings.Contains(lower, "temps") || strings.Contains(lower, "monde") {
		return "narrative"
	}

	// Mots formels
	if strings.Contains(lower, "considera") || strings.Contains(lower, "signifi") ||
		strings.Contains(lower, "analyse") {
		return "formal"
	}

	return "neutral"
}

// CalculerIntensiteSemantique évalue l'importance sémantique d'un mot
func CalculerIntensiteSemantique(mot string) float64 {
	if word, ok := Words[strings.ToLower(mot)]; ok {
		return math.Min(word.Poids/10.0, 1.0) // Normaliser poids en [0,1]
	}
	return 0.3 // Poids par défaut
}

// GenererAvecReseauAtomique génère du texte via résonance distribuée
func GenererAvecReseauAtomique(phrases []Phrase, iterationsMax int, cibleCompression float64) string {
	// Créer le réseau
	vecteurCible := MoyennePondereeVecteurs(phrases)
	reseau := NouveauReseauGeneration(vecteurCible, phrases)

	// Itérer jusqu'à target ou max iterations
	tailleTarget := int(float64(len(ExtraireVocabulaire(phrases))) * cibleCompression)
	if tailleTarget < 20 {
		tailleTarget = 20
	}

	for iter := 0; iter < iterationsMax && len(reseau.MotsEmis) < tailleTarget; iter++ {
		_ = reseau.IterationAtomique()
	}

	// Assembler le texte final
	texte := strings.Join(reseau.MotsEmis, " ")
	texte = strings.TrimSpace(texte)

	// Ajouter ponctuation
	if len(texte) > 0 {
		texte = strings.ToUpper(texte[:1]) + texte[1:]
		if !strings.HasSuffix(texte, ".") && !strings.HasSuffix(texte, "!") && !strings.HasSuffix(texte, "?") {
			texte += "."
		}
	}

	return texte
}

// DecouperEnSousReseaux divise un ensemble de phrases en chunks gérés indépendamment
type SousReseau struct {
	Phrases []Phrase
	Reseau  *ReseauAtomiquePourGeneration
	Texture string // Texte généré par ce sous-réseau
}

// GenererAvecSousReseaux génère sur longs textes avec résonance inter-réseaux
func GenererAvecSousReseaux(phrases []Phrase, tailleChunk int, cibleCompression float64) string {
	// Diviser en chunks
	sousReseaux := make([]*SousReseau, 0)

	for i := 0; i < len(phrases); i += tailleChunk {
		fin := i + tailleChunk
		if fin > len(phrases) {
			fin = len(phrases)
		}

		chunk := phrases[i:fin]
		sr := &SousReseau{
			Phrases: chunk,
			Texture: "",
		}
		sousReseaux = append(sousReseaux, sr)
	}

	// Générer chaque sous-réseau
	for _, sr := range sousReseaux {
		vTarget := MoyennePondereeVecteurs(sr.Phrases)
		sr.Reseau = NouveauReseauGeneration(vTarget, sr.Phrases)

		// Itérer
		for iter := 0; iter < 50; iter++ {
			_ = sr.Reseau.IterationAtomique()
		}

		sr.Texture = strings.Join(sr.Reseau.MotsEmis, " ")
	}

	// Fusionner avec résonance: vecteur cible global influence chaque sous-réseau
	vGlobal := MoyennePondereeVecteurs(phrases)
	texturesFinal := make([]string, 0)

	for _, sr := range sousReseaux {
		// Augmenter activation de atomes proches du vecteur global
		for _, atome := range sr.Reseau.Atomes {
			vAtome := VectoriserMot(atome.Concept)
			simGlobal := SimilariteCosinus(vAtome, vGlobal)
			atome.Activation = math.Min(atome.Activation+0.2*simGlobal, 1.0)
		}

		// Réitérer avec influence globale
		for iter := 0; iter < 30; iter++ {
			mots := sr.Reseau.IterationAtomique()
			if len(mots) > 0 {
				texturesFinal = append(texturesFinal, mots...)
			}
		}
	}

	// Assembler final
	texte := strings.Join(texturesFinal, " ")
	texte = strings.TrimSpace(texte)

	if len(texte) > 0 {
		texte = strings.ToUpper(texte[:1]) + texte[1:]
		if !strings.HasSuffix(texte, ".") {
			texte += "."
		}
	}

	return texte
}

// CalculerCoherenceReseau évalue la cohérence globale du réseau
func (r *ReseauAtomiquePourGeneration) CalculerCoherenceReseau() float64 {
	if len(r.MotsEmis) == 0 {
		return 0.0
	}

	totalCoherence := 0.0
	for _, mot := range r.MotsEmis {
		vMot := VectoriserMot(mot)
		coherence := SimilariteCosinus(vMot, r.VecteurCible)
		totalCoherence += coherence
	}

	return totalCoherence / float64(len(r.MotsEmis))
}
