// Package database - Hybrid Atomic-Probability Engine
// Combine la stabilité atomique (T.R.A.) avec la probabilité linguistique
//
// INNOVATION : Au lieu de choisir entre probabilité OU stabilité,
// on les MIXE pour obtenir le meilleur des deux mondes:
// - Probabilité → capture les patterns statistiques du langage
// - Stabilité atomique → capture la cohérence sémantique profonde
//
// FORMULE HYBRIDE:
// Score = α·P(probabilité) + β·S(stabilité) + γ·P(prob)·S(stab)
// Terme multiplicatif γ·P·S capture les synergies!

package database

import (
	"fmt"
	"math"
)

// HybridAtomicProbability moteur hybride
type HybridAtomicProbability struct {
	AlphaProb     float64 // Poids probabilité pure (0.3-0.5)
	BetaStability float64 // Poids stabilité pure (0.3-0.5)
	GammaSynergy  float64 // Poids synergie multiplicative (0.1-0.3)
	TauThreshold  float64 // Seuil métacognitif pour décider Atomique vs Hybride

	// Réseau atomique pour stabilité
	AtomicNetwork *AtomicNetwork

	// Calculateur de perplexité pour probabilité
	PerplexityCalc *PerplexityCalculator
}

// NewHybridAtomicProbability crée un moteur hybride
func NewHybridAtomicProbability() *HybridAtomicProbability {
	return &HybridAtomicProbability{
		AlphaProb:     0.38, // Probabilité un peu renforcée pour factuels
		BetaStability: 0.44, // Stabilité toujours majoritaire mais réduite
		GammaSynergy:  0.18, // Synergie maintenue
		TauThreshold:  0.80, // Seuil métacognitif plus bas → hybrider plus tôt

		AtomicNetwork:  NewAtomicNetwork(300), // Réseau de 300 atomes (3× plus grand)
		PerplexityCalc: NewPerplexityCalculator(),
	}
}

// HybridScore calcule le score hybride combinant probabilité ET stabilité
func (hap *HybridAtomicProbability) HybridScore(context string, continuation string) float64 {
	decision := hap.MetacognitiveHybridScore(context, continuation)
	return decision.FinalScore
}

// MetacognitiveDecision retourne la décision du système de gating
type MetacognitiveDecision struct {
	Mode             string  // "atomic" ou "hybrid"
	StabilitySa      float64 // Mesure Sa (stabilité atomique)
	Entropy          float64 // Entropie des catégories
	ProbabilityScore float64 // Score probabiliste
	StabilityScore   float64 // Score stabilité réseau
	FinalScore       float64 // Score retenu après gating
	Threshold        float64 // Seuil τ courant
	FusionAlpha      float64 // Coefficient de fusion atomique/proba
}

// MetacognitiveHybridScore applique le gating métacognitif avant de fusionner
func (hap *HybridAtomicProbability) MetacognitiveHybridScore(context, continuation string) *MetacognitiveDecision {
	// 1) Préparer les tokens et catégories
	combinedText := context + " " + continuation
	tokens := TokeniserTexte(combinedText)
	categories := ActiverCategoriesParTexte(tokens)

	// 2) Calcul des composantes
	probScore := hap.calculateProbabilityScore(context, continuation)
	stabilityScore := hap.calculateAtomicStability(context, continuation)
	resonanceScore := hap.computeAtomicResonanceScore(tokens, categories)
	entropy := hap.computeDecisionEntropy(categories)

	// 3) Gating selon τ (stabilité) et entropie (conflit)
	mode := "hybrid"
	finalScore := 0.0
	fusionAlpha := 0.0

	if resonanceScore >= hap.TauThreshold && entropy < 1.0 {
		// Structure claire → Atomique pur (gain de temps)
		mode = "atomic"
		finalScore = stabilityScore
		fusionAlpha = 1.0
	} else {
		// Conflit ou incertitude → activer hybridation
		fusionAlpha = hap.BetaStability / (hap.AlphaProb + hap.BetaStability)
		baseHybrid := (hap.AlphaProb*probScore +
			hap.BetaStability*stabilityScore +
			hap.GammaSynergy*probScore*stabilityScore)
		finalScore = fusionAlpha*stabilityScore + (1-fusionAlpha)*probScore
		// Injecter la synergie pour garder l'effet multiplicatif
		finalScore = (finalScore*0.7 + baseHybrid*0.3)
		mode = "hybrid"
	}

	// 4) Auto-ajustement du seuil τ
	hap.adaptThreshold(resonanceScore, entropy, probScore, stabilityScore, mode)

	return &MetacognitiveDecision{
		Mode:             mode,
		StabilitySa:      resonanceScore,
		Entropy:          entropy,
		ProbabilityScore: probScore,
		StabilityScore:   stabilityScore,
		FinalScore:       clamp01(finalScore),
		Threshold:        hap.TauThreshold,
		FusionAlpha:      fusionAlpha,
	}
}

// calculateProbabilityScore calcule le score de probabilité (0-1)
func (hap *HybridAtomicProbability) calculateProbabilityScore(context, continuation string) float64 {
	// Utiliser perplexité inversée comme score de probabilité
	fullText := context + " " + continuation
	perplexityResult := hap.PerplexityCalc.CalculatePerplexity(fullText)

	// Perplexité basse = probabilité haute
	// Normaliser: perplexité [1, 20] → score [1.0, 0.0]
	perplexity := perplexityResult.GlobalPerplexity

	// Score = 1 / (1 + perplexité/10)
	probScore := 1.0 / (1.0 + perplexity/10.0)

	return probScore
}

// calculateAtomicStability calcule la stabilité du réseau atomique (0-1)
func (hap *HybridAtomicProbability) calculateAtomicStability(context, continuation string) float64 {
	// 1. Activer les atomes avec le contexte
	contextTokens := TokeniserTexte(context)
	contextCategories := ActiverCategoriesParTexte(contextTokens)

	// Injecter le contexte dans le réseau
	hap.activateNetworkWithCategories(contextCategories)

	// 2. Mesurer cohérence AVANT ajout continuation
	initialCoherence := hap.AtomicNetwork.GetNetworkCoherence()

	// 3. Ajouter la continuation
	contTokens := TokeniserTexte(continuation)
	contCategories := ActiverCategoriesParTexte(contTokens)

	hap.activateNetworkWithCategories(contCategories)

	// 4. Faire converger le réseau (plus d'itérations pour meilleure stabilité)
	for i := 0; i < 10; i++ {
		hap.AtomicNetwork.IterateNetwork()
	}

	// 5. Mesurer cohérence APRÈS
	finalCoherence := hap.AtomicNetwork.GetNetworkCoherence()

	// 6. Stabilité = cohérence finale avec bonus/pénalité selon évolution
	stabilityScore := finalCoherence

	// BONUS si cohérence a augmenté (continuation renforce le contexte)
	if finalCoherence > initialCoherence {
		improvementBonus := (finalCoherence - initialCoherence) * 0.8 // Plus de bonus
		stabilityScore += improvementBonus
	}

	// PÉNALITÉ si cohérence a chuté (continuation perturbe)
	if finalCoherence < initialCoherence {
		degradationPenalty := (initialCoherence - finalCoherence) * 1.2 // Plus de pénalité
		stabilityScore -= degradationPenalty
	}

	// Clamp [0, 1]
	return math.Max(0.0, math.Min(1.0, stabilityScore))
}

// activateNetworkWithCategories active le réseau avec des catégories
func (hap *HybridAtomicProbability) activateNetworkWithCategories(categories map[int]int) {
	// Distribuer l'activation aux atomes selon les catégories
	numAtoms := len(hap.AtomicNetwork.Atoms)

	for category, count := range categories {
		// Chaque catégorie active un sous-ensemble d'atomes
		startIdx := (category * numAtoms / 50) % numAtoms
		endIdx := ((category + 1) * numAtoms / 50) % numAtoms

		if endIdx < startIdx {
			endIdx = numAtoms
		}

		activation := float64(count) / 10.0
		if activation > 1.0 {
			activation = 1.0
		}

		for i := startIdx; i < endIdx && i < numAtoms; i++ {
			// Ajouter perception directement
			hap.AtomicNetwork.Atoms[i].Perceptions[category] = activation
		}
	}
}

// CalculateDetailedHybridScore calcule un score détaillé avec décomposition
func (hap *HybridAtomicProbability) CalculateDetailedHybridScore(context, continuation string) *HybridScoreDetail {
	probScore := hap.calculateProbabilityScore(context, continuation)
	stabilityScore := hap.calculateAtomicStability(context, continuation)
	synergyScore := probScore * stabilityScore

	hybridScore := (hap.AlphaProb*probScore +
		hap.BetaStability*stabilityScore +
		hap.GammaSynergy*synergyScore)

	return &HybridScoreDetail{
		ProbabilityScore: probScore,
		StabilityScore:   stabilityScore,
		SynergyScore:     synergyScore,
		HybridScore:      hybridScore,
		AlphaWeight:      hap.AlphaProb,
		BetaWeight:       hap.BetaStability,
		GammaWeight:      hap.GammaSynergy,
	}
}

// HybridScoreDetail détails du score hybride
type HybridScoreDetail struct {
	ProbabilityScore float64 // Score probabilité pure [0-1]
	StabilityScore   float64 // Score stabilité atomique [0-1]
	SynergyScore     float64 // Score synergie P×S [0-1]
	HybridScore      float64 // Score hybride final [0-1]

	// Poids utilisés
	AlphaWeight float64
	BetaWeight  float64
	GammaWeight float64
}

// String représentation textuelle du score détaillé
func (hsd *HybridScoreDetail) String() string {
	return fmt.Sprintf("Hybrid[P:%.3f S:%.3f Syn:%.3f → %.3f]",
		hsd.ProbabilityScore,
		hsd.StabilityScore,
		hsd.SynergyScore,
		hsd.HybridScore)
}

// AdaptWeights ajuste les poids selon les performances
func (hap *HybridAtomicProbability) AdaptWeights(probCorrect, stabilityCorrect, hybridCorrect int, total int) {
	if total == 0 {
		return
	}

	// Calculer taux de succès de chaque approche
	probRate := float64(probCorrect) / float64(total)
	stabRate := float64(stabilityCorrect) / float64(total)
	_ = hybridCorrect // Utilisé pour monitoring futur

	// Réajuster les poids en faveur de l'approche la plus performante
	totalRate := probRate + stabRate
	if totalRate > 0 {
		// Normaliser les poids Alpha et Beta
		newAlpha := probRate / totalRate * 0.7 // 70% divisé selon performance
		newBeta := stabRate / totalRate * 0.7
		newGamma := 0.3 // 30% reste pour synergie

		// Appliquer graduellement (learning rate 0.3)
		learningRate := 0.3
		hap.AlphaProb = hap.AlphaProb*(1-learningRate) + newAlpha*learningRate
		hap.BetaStability = hap.BetaStability*(1-learningRate) + newBeta*learningRate
		hap.GammaSynergy = hap.GammaSynergy*(1-learningRate) + newGamma*learningRate

		// Normaliser pour que α+β+γ = 1
		total := hap.AlphaProb + hap.BetaStability + hap.GammaSynergy
		hap.AlphaProb /= total
		hap.BetaStability /= total
		hap.GammaSynergy /= total
	}
}

// ResetNetwork réinitialise le réseau atomique
func (hap *HybridAtomicProbability) ResetNetwork() {
	// Réinitialiser tous les états atomiques à 0.5
	for i := range hap.AtomicNetwork.Atoms {
		hap.AtomicNetwork.Atoms[i].InternalState = 0.5
		hap.AtomicNetwork.Atoms[i].Perceptions = make(map[int]float64)
	}
}

// GetWeights retourne les poids actuels
func (hap *HybridAtomicProbability) GetWeights() (alpha, beta, gamma float64) {
	return hap.AlphaProb, hap.BetaStability, hap.GammaSynergy
}

// SetWeights définit manuellement les poids
func (hap *HybridAtomicProbability) SetWeights(alpha, beta, gamma float64) {
	// Normaliser
	total := alpha + beta + gamma
	hap.AlphaProb = alpha / total
	hap.BetaStability = beta / total
	hap.GammaSynergy = gamma / total
}

// computeAtomicResonanceScore calcule Sa (stabilité atomique) via résonance
// Approche pragmatique: plus une catégorie est présente, plus la résonance est forte
func (hap *HybridAtomicProbability) computeAtomicResonanceScore(tokens []string, categories map[int]int) float64 {
	if len(tokens) == 0 {
		return 0
	}

	totalTokens := float64(len(tokens))
	score := 0.0

	for _, count := range categories {
		weight := float64(count) / totalTokens // ω_i
		// ρ approximée par une courbe saturante (1 - e^{-x})
		resonance := 1.0 - math.Exp(-(float64(count)/totalTokens)*1.5)
		score += weight * resonance
	}

	return clamp01(score)
}

// computeDecisionEntropy calcule l'entropie H sur la distribution des catégories
func (hap *HybridAtomicProbability) computeDecisionEntropy(categories map[int]int) float64 {
	if len(categories) == 0 {
		return 0
	}

	total := 0.0
	for _, c := range categories {
		total += float64(c)
	}

	entropy := 0.0
	for _, c := range categories {
		p := float64(c) / total
		if p > 0 {
			entropy += -p * math.Log2(p)
		}
	}

	return entropy
}

// adaptThreshold ajuste dynamiquement τ selon la clarté ou l'incertitude
func (hap *HybridAtomicProbability) adaptThreshold(sa, entropy, probScore, stabilityScore float64, mode string) {
	current := hap.TauThreshold
	target := current

	if mode == "atomic" && entropy < 0.8 && math.Abs(probScore-stabilityScore) < 0.12 {
		target = current + 0.02 // Confiance confirmée → seuil plus ambitieux
	} else if mode == "hybrid" && math.Abs(probScore-stabilityScore) > 0.20 {
		target = current - 0.02 // Désaccord fort → seuil plus bas pour hybrider plus tôt
	}

	// Appliquer lissage
	learningRate := 0.3
	newTau := current + (target-current)*learningRate
	hap.TauThreshold = clamp(newTau, 0.60, 0.95)
}

// Helpers
func clamp01(v float64) float64 {
	return clamp(v, 0.0, 1.0)
}

func clamp(v, min, max float64) float64 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}
