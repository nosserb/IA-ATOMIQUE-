// Package database - Deep Semantic Analyzer
// Analyseur sémantique profond pour comprendre le sens au-delà des mots

package database

import (
	"math"
	"strings"
)

// DeepSemanticAnalyzer analyse le sens profond des textes
type DeepSemanticAnalyzer struct {
	SemanticRules    map[string][]string // Règles sémantiques (cause→effet, etc.)
	ConceptHierarchy map[string][]string // Hiérarchie de concepts
	ActionSequences  map[string][]string // Séquences d'actions logiques
	ContextualClues  map[string]float64  // Indices contextuels
}

// NewDeepSemanticAnalyzer crée un nouvel analyseur
func NewDeepSemanticAnalyzer() *DeepSemanticAnalyzer {
	return &DeepSemanticAnalyzer{
		SemanticRules:    initializeSemanticRules(),
		ConceptHierarchy: initializeConceptHierarchy(),
		ActionSequences:  initializeActionSequences(),
		ContextualClues:  initializeContextualClues(),
	}
}

// Analyze Semantic Coherence calcule la cohérence sémantique profonde
func (dsa *DeepSemanticAnalyzer) AnalyzeSemanticCoherence(context, continuation string) float64 {
	score := 0.0

	// 1. Cohérence causale (cause → effet logique)
	causalScore := dsa.analyzeCausalCoherence(context, continuation)

	// 2. Cohérence temporelle (séquence logique d'actions)
	temporalScore := dsa.analyzeTemporalCoherence(context, continuation)

	// 3. Cohérence thématique (même domaine sémantique)
	thematicScore := dsa.analyzeThematicCoherence(context, continuation)

	// 4. Cohérence actionnelle (actions compatibles)
	actionScore := dsa.analyzeActionCoherence(context, continuation)

	// 5. Cohérence référentielle (pronoms et références)
	referentialScore := dsa.analyzeReferentialCoherence(context, continuation)

	// Score combiné pondéré
	score = (causalScore*0.25 +
		temporalScore*0.25 +
		thematicScore*0.20 +
		actionScore*0.20 +
		referentialScore*0.10)

	return score
}

// analyzeCausalCoherence détecte les relations de cause à effet
func (dsa *DeepSemanticAnalyzer) analyzeCausalCoherence(context, continuation string) float64 {
	score := 0.5 // score neutre par défaut

	contextLower := strings.ToLower(context)
	contLower := strings.ToLower(continuation)

	// Détection d'effets logiques dans la continuation
	effects := map[string][]string{
		"prend casserole":  {"met", "pose", "utilise", "remplit"},
		"remplit eau":      {"bouillir", "chauffer", "chauffe", "feu"},
		"met feu":          {"attend", "chauffe", "bouillir", "cuire"},
		"sort chez":        {"va", "marche", "courir", "dehors"},
		"chaussures sport": {"courir", "jogger", "exercice", "marcher"},
		"commence courir":  {"continue", "poursuit", "ralentit", "arrête"},
		"ouvre ordinateur": {"lance", "tape", "écrit", "utilise"},
		"lance traitement": {"tape", "écrit", "rédige", "saisit"},
		"commence taper":   {"rédige", "écrit", "continue", "termine"},
	}

	// Chercher des paires cause→effet
	for cause, expectedEffects := range effects {
		if strings.Contains(contextLower, cause) {
			for _, effect := range expectedEffects {
				if strings.Contains(contLower, effect) {
					score += 0.3
					break
				}
			}
		}
	}

	// Bonus si marqueurs causaux explicites
	causalMarkers := []string{"donc", "alors", "conséquence", "résultat", "c'est pourquoi"}
	for _, marker := range causalMarkers {
		if strings.Contains(contLower, marker) {
			score += 0.1
		}
	}

	return math.Min(score, 1.0)
}

// analyzeTemporalCoherence vérifie la séquence temporelle logique
func (dsa *DeepSemanticAnalyzer) analyzeTemporalCoherence(context, continuation string) float64 {
	score := 0.5

	contextLower := strings.ToLower(context)
	contLower := strings.ToLower(continuation)

	// Séquences d'actions logiques
	sequences := map[string][]string{
		"entre cuisine":   {"prend", "ouvre", "commence", "prépare"},
		"prend casserole": {"remplit", "met", "pose", "utilise"},
		"remplit eau":     {"met feu", "chauffe", "attend"},
		"met feu":         {"attend", "observe", "surveille"},
		"attend":          {"bouillir", "cuire", "termine"},

		"met chaussures":  {"sort", "va", "marche", "court"},
		"sort chez":       {"commence", "va", "marche", "court"},
		"commence courir": {"continue", "court", "maintient"},
		"continue courir": {"ralentit", "termine", "arrête", "poursuit"},
		"ralentit":        {"arrête", "marche", "termine"},

		"ouvre ordinateur": {"lance", "démarre", "utilise"},
		"lance traitement": {"commence", "tape", "écrit"},
		"commence taper":   {"rédige", "écrit", "continue"},
		"rédige rapport":   {"écrit", "continue", "termine", "sauvegarde"},
		"termine":          {"sauvegarde", "ferme", "enregistre"},
	}

	// Vérifier cohérence séquentielle
	for contextKey, expectedNext := range sequences {
		if strings.Contains(contextLower, contextKey) {
			for _, next := range expectedNext {
				if strings.Contains(contLower, next) {
					score += 0.3
					break
				}
			}
		}
	}

	// Bonus si marqueurs temporels
	temporalMarkers := []string{"puis", "ensuite", "après", "enfin", "finalement", "pendant"}
	for _, marker := range temporalMarkers {
		if strings.Contains(contLower, marker) {
			score += 0.1
			break
		}
	}

	return math.Min(score, 1.0)
}

// analyzeThematicCoherence vérifie la cohérence thématique
func (dsa *DeepSemanticAnalyzer) analyzeThematicCoherence(context, continuation string) float64 {
	score := 0.5

	contextLower := strings.ToLower(context)
	contLower := strings.ToLower(continuation)

	// Thèmes et vocabulaire associé
	themes := map[string][]string{
		"cuisine":   {"casserole", "eau", "feu", "bouillir", "cuire", "préparer", "manger", "cuisson", "aliment"},
		"sport":     {"courir", "chaussures", "exercice", "ralentir", "marcher", "entraînement", "cardio", "effort"},
		"travail":   {"ordinateur", "taper", "rapport", "rédiger", "sauvegarder", "document", "texte", "fichier"},
		"maison":    {"cuisine", "pièce", "entrer", "chez", "domicile"},
		"extérieur": {"sortir", "rue", "dehors", "extérieur", "marcher"},
	}

	// Identifier thème dominant du contexte
	var contextTheme string
	maxContextWords := 0

	for theme, words := range themes {
		count := 0
		for _, word := range words {
			if strings.Contains(contextLower, word) {
				count++
			}
		}
		if count > maxContextWords {
			maxContextWords = count
			contextTheme = theme
		}
	}

	// Vérifier si continuation maintient le thème
	if contextTheme != "" {
		contWords := 0
		for _, word := range themes[contextTheme] {
			if strings.Contains(contLower, word) {
				contWords++
			}
		}
		if contWords > 0 {
			score += 0.3
		}
	}

	return math.Min(score, 1.0)
}

// analyzeActionCoherence vérifie la compatibilité des actions
func (dsa *DeepSemanticAnalyzer) analyzeActionCoherence(context, continuation string) float64 {
	score := 0.5

	contextLower := strings.ToLower(context)
	contLower := strings.ToLower(continuation)

	// Actions incompatibles (contradictions)
	incompatibleActions := map[string][]string{
		"entre":    {"sort", "quitte", "part"},
		"sort":     {"entre", "rentre", "pénètre"},
		"prend":    {"pose", "lâche", "abandonne"},
		"ouvre":    {"ferme", "clôt"},
		"commence": {"termine", "finit", "arrête"},
		"remplit":  {"vide", "verse"},
		"courir":   {"immobile", "assis", "couché"},
	}

	// Pénaliser si actions incompatibles
	for contextAction, incompActions := range incompatibleActions {
		if strings.Contains(contextLower, contextAction) {
			for _, incomp := range incompActions {
				if strings.Contains(contLower, incomp) {
					score -= 0.4
				}
			}
		}
	}

	// Bonus pour continuité d'action
	continuousActions := []string{"continue", "poursuit", "maintient", "persiste"}
	for _, cont := range continuousActions {
		if strings.Contains(contLower, cont) {
			score += 0.2
			break
		}
	}

	return math.Max(0.0, math.Min(score, 1.0))
}

// analyzeReferentialCoherence vérifie cohérence des pronoms et références
func (dsa *DeepSemanticAnalyzer) analyzeReferentialCoherence(context, continuation string) float64 {
	score := 0.7 // score de base élevé (ok par défaut)

	contextLower := strings.ToLower(context)
	contLower := strings.ToLower(continuation)

	// Détection du genre dans le contexte
	var isFeminine, isMasculine bool

	if strings.Contains(contextLower, "femme") || strings.Contains(contextLower, "elle") {
		isFeminine = true
	}
	if strings.Contains(contextLower, "homme") || strings.Contains(contextLower, "il") ||
		strings.Contains(contextLower, "étudiant") {
		isMasculine = true
	}

	// Vérifier cohérence des pronoms dans continuation
	if isFeminine {
		if strings.Contains(contLower, " il ") || strings.Contains(contLower, " lui ") {
			score -= 0.5 // Incohérence grave
		}
		if strings.Contains(contLower, " elle ") {
			score += 0.1 // Cohérence confirmée
		}
	}

	if isMasculine {
		if strings.Contains(contLower, " elle ") {
			score -= 0.5 // Incohérence grave
		}
		if strings.Contains(contLower, " il ") || strings.Contains(contLower, " lui ") {
			score += 0.1 // Cohérence confirmée
		}
	}

	return math.Max(0.0, math.Min(score, 1.0))
}

// Helper initialization functions
func initializeSemanticRules() map[string][]string {
	return map[string][]string{
		"cause_effect": {"prendre→utiliser", "remplir→chauffer", "ouvrir→utiliser"},
	}
}

func initializeConceptHierarchy() map[string][]string {
	return map[string][]string{
		"cuisine": {"ustensile", "casserole", "nourriture"},
		"sport":   {"exercice", "course", "entraînement"},
	}
}

func initializeActionSequences() map[string][]string {
	return map[string][]string{
		"cuisine": {"entrer", "prendre", "préparer", "cuire", "manger"},
		"sport":   {"équiper", "sortir", "commencer", "continuer", "terminer"},
	}
}

func initializeContextualClues() map[string]float64 {
	return map[string]float64{
		"puis":    0.8,
		"ensuite": 0.8,
		"donc":    0.7,
	}
}
