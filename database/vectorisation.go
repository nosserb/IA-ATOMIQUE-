package database

import (
	"math"
	"strings"
)

// VecteurAtomique représente un vecteur dans l'espace sémantique
// Dimensions = catégories + fréquence + position
type VecteurAtomique struct {
	Dimensions []float64 // d ≈ 20 (catégories + meta)
	Energie    float64   // Norme L2 du vecteur
	Contenu    string    // Texte original (phrase ou mot)
}

// DIM_VECTORIEL = nombre de dimensions (6 catégories + 5 meta)
const DIM_VECTORIEL = 11

// VectoriserPhrase transforme une phrase en vecteur atomique
// Basé sur: fréquence mots-clés, catégories, longueur, position
func VectoriserPhrase(phrase *Phrase) VecteurAtomique {
	vec := VecteurAtomique{
		Dimensions: make([]float64, DIM_VECTORIEL),
		Contenu:    phrase.Contenu,
	}

	// Dimension 0-5: Activation des catégories
	// Category IDs: 1=TECH, 2=HISTOIRE, 3=BUSINESS, 4=ALIMENTATION, 5=SANTÉ, 6=VERBE
	categorieMap := map[int]int{
		1: 0, // TECH -> dim 0
		2: 1, // HISTOIRE -> dim 1
		3: 2, // BUSINESS -> dim 2
		4: 3, // ALIMENTATION -> dim 3
		5: 4, // SANTÉ -> dim 4
		6: 5, // VERBE -> dim 5
	}

	// Compter les mots par catégorie
	for _, mot := range phrase.MotsClés {
		if word, ok := Words[strings.ToLower(mot)]; ok {
			if idx, found := categorieMap[word.Categorie]; found {
				vec.Dimensions[idx] += word.Poids
			}
		}
	}

	// Dimension 6: Longueur normalisée (plus de mots = plus d'info)
	longueur := float64(len(phrase.Mots))
	vec.Dimensions[6] = math.Min(longueur/50.0, 1.0) // Normaliser à [0,1]

	// Dimension 7: Énergie intrinsèque de la phrase
	vec.Dimensions[7] = phrase.Energie

	// Dimension 8: Densité de mots-clés (mots_clés / total mots)
	if len(phrase.Mots) > 0 {
		densiteMotsClés := float64(len(phrase.MotsClés)) / float64(len(phrase.Mots))
		vec.Dimensions[8] = densiteMotsClés
	}

	// Dimension 9: Facteur de traduction (confiance)
	vec.Dimensions[9] = phrase.FacteurConfiance
	if vec.Dimensions[9] == 0 {
		vec.Dimensions[9] = 1.0 // Default: français original
	}

	// Dimension 10: Position normalisée (début/fin du doc)
	vec.Dimensions[10] = 0.5 // Défaut: neutre (à adapter par contexte)

	// Calculer l'énergie (norme L2)
	for _, d := range vec.Dimensions {
		vec.Energie += d * d
	}
	vec.Energie = math.Sqrt(vec.Energie)

	// Normaliser le vecteur
	if vec.Energie > 0 {
		for i := range vec.Dimensions {
			vec.Dimensions[i] /= vec.Energie
		}
	}

	return vec
}

// VectoriserMot crée un vecteur pour un mot individuel
func VectoriserMot(mot string) VecteurAtomique {
	vec := VecteurAtomique{
		Dimensions: make([]float64, DIM_VECTORIEL),
		Contenu:    mot,
	}

	categorieMap := map[int]int{
		1: 0, // TECH -> dim 0
		2: 1, // HISTOIRE -> dim 1
		3: 2, // BUSINESS -> dim 2
		4: 3, // ALIMENTATION -> dim 3
		5: 4, // SANTÉ -> dim 4
		6: 5, // VERBE -> dim 5
	}

	// Chercher le mot dans le lexique
	if word, ok := Words[strings.ToLower(mot)]; ok {
		if idx, found := categorieMap[word.Categorie]; found {
			vec.Dimensions[idx] = word.Poids
		}
	} else {
		// Mot non trouvé: activation neutre
		vec.Dimensions[0] = 0.1 // Petit poids par défaut
	}

	// Longueur du mot normalisée
	vec.Dimensions[6] = float64(len(mot)) / 20.0
	if vec.Dimensions[6] > 1 {
		vec.Dimensions[6] = 1.0
	}

	// Densité (pour un mot, c'est 1.0)
	vec.Dimensions[8] = 1.0

	// Calculer l'énergie
	for _, d := range vec.Dimensions {
		vec.Energie += d * d
	}
	vec.Energie = math.Sqrt(vec.Energie)

	// Normaliser
	if vec.Energie > 0 {
		for i := range vec.Dimensions {
			vec.Dimensions[i] /= vec.Energie
		}
	}

	return vec
}

// SimilariteCosinus calcule la similarité cosinus entre deux vecteurs
func SimilariteCosinus(v1, v2 VecteurAtomique) float64 {
	if v1.Energie == 0 || v2.Energie == 0 {
		return 0
	}

	produitScalaire := 0.0
	for i := range v1.Dimensions {
		if i < len(v2.Dimensions) {
			produitScalaire += v1.Dimensions[i] * v2.Dimensions[i]
		}
	}

	return produitScalaire / (v1.Energie * v2.Energie)
}

// CombinerVecteurs fusionne plusieurs vecteurs avec pondération
func CombinerVecteurs(vecteurs []VecteurAtomique, poids []float64) VecteurAtomique {
	resultat := VecteurAtomique{
		Dimensions: make([]float64, DIM_VECTORIEL),
	}

	totalPoids := 0.0
	for i, v := range vecteurs {
		w := 1.0
		if i < len(poids) {
			w = poids[i]
		}
		totalPoids += w

		for j := range resultat.Dimensions {
			if j < len(v.Dimensions) {
				resultat.Dimensions[j] += w * v.Dimensions[j]
			}
		}
	}

	// Normaliser par poids total
	if totalPoids > 0 {
		for i := range resultat.Dimensions {
			resultat.Dimensions[i] /= totalPoids
		}
	}

	// Calculer l'énergie
	for _, d := range resultat.Dimensions {
		resultat.Energie += d * d
	}
	resultat.Energie = math.Sqrt(resultat.Energie)

	// Normaliser le vecteur final
	if resultat.Energie > 0 {
		for i := range resultat.Dimensions {
			resultat.Dimensions[i] /= resultat.Energie
		}
	}

	return resultat
}

// MoyennePondereeVecteurs crée un vecteur cible VR
// VR = Σ E(Pi)·vi / Σ E(Pi)
func MoyennePondereeVecteurs(phrases []Phrase) VecteurAtomique {
	vecteurs := make([]VecteurAtomique, len(phrases))
	poids := make([]float64, len(phrases))

	totalEnergie := 0.0
	for i, p := range phrases {
		vecteurs[i] = VectoriserPhrase(&p)
		poids[i] = p.Energie
		if poids[i] == 0 {
			poids[i] = 1.0 // Éviter division par zéro
		}
		totalEnergie += poids[i]
	}

	// Normaliser les poids
	for i := range poids {
		poids[i] /= totalEnergie
	}

	return CombinerVecteurs(vecteurs, poids)
}

// MettreAJourVecteur adapte VR après génération d'un mot
// VR(k+1) = VR(k) + γ·(vrk - VR(k))
func MettreAJourVecteur(vCourant VecteurAtomique, vMotGeneré VecteurAtomique, gamma float64) VecteurAtomique {
	resultat := VecteurAtomique{
		Dimensions: make([]float64, DIM_VECTORIEL),
	}

	for i := range resultat.Dimensions {
		// Appliquer l'équation: VR(k+1) = VR(k) + γ·(vrk - VR(k))
		resultat.Dimensions[i] = vCourant.Dimensions[i] + gamma*(vMotGeneré.Dimensions[i]-vCourant.Dimensions[i])
	}

	// Recalculer l'énergie
	for _, d := range resultat.Dimensions {
		resultat.Energie += d * d
	}
	resultat.Energie = math.Sqrt(resultat.Energie)

	// Normaliser
	if resultat.Energie > 0 {
		for i := range resultat.Dimensions {
			resultat.Dimensions[i] /= resultat.Energie
		}
	}

	return resultat
}

// ExtraireVocabulaire retourne tous les mots uniques des phrases
func ExtraireVocabulaire(phrases []Phrase) []string {
	vocabMap := make(map[string]bool)
	for _, p := range phrases {
		for _, mot := range p.Mots {
			motLower := strings.ToLower(mot)
			vocabMap[motLower] = true
		}
	}

	vocab := make([]string, 0, len(vocabMap))
	for mot := range vocabMap {
		vocab = append(vocab, mot)
	}
	return vocab
}
