// Package database - Hellaswag Benchmark
// Test de raisonnement de bon sens et continuation logique
//
// PRINCIPE : Compléter des actions quotidiennes
// - Contexte : Début d'une action (ex: "Il prend une casserole...")
// - 4 suites possibles, 1 seule logique
// - Test la compréhension des séquences causales
//
// AVANTAGE ATOMIQUE :
// - Perplexité 1.05 = détection immédiate de la suite stable
// - Structure atomique cohérente pour actions logiques

package database

import (
	"fmt"
	"math"
	"strings"
)

// HellaswagQuestion représente une question Hellaswag
type HellaswagQuestion struct {
	ID       int
	Context  string   // Contexte de l'action
	Endings  []string // 4 fins possibles
	Answer   int      // Index de la fin correcte
	Activity string   // Type d'activité (ex: "cuisine", "sport")
}

// HellaswagResult résultat pour une question
type HellaswagResult struct {
	Question       *HellaswagQuestion
	SelectedIndex  int
	IsCorrect      bool
	Confidence     float64
	EndingScores   []float64
	PerplexityDiff float64 // Différence de perplexité entre meilleur et pire choix
}

// HellaswagBenchmark résultats globaux
type HellaswagBenchmark struct {
	TotalQuestions    int
	CorrectAnswers    int
	Score             float64
	AverageConfidence float64
	AveragePerplexity float64
	Results           []*HellaswagResult
}

// HellaswagEngine moteur d'évaluation optimisé
type HellaswagEngine struct {
	perplexityCalc       *PerplexityCalculator
	ContextMemory        map[string][]string      // Mémoire des continuations logiques
	TemporalPatterns     map[string]float64       // Patterns temporels
	CausalPatterns       map[string]float64       // Patterns causals
	AdaptiveLearning     bool                     // Apprentissage adaptatif
	ContextEngine        *ContextEngine           // Moteur de contexte enrichi
	DeepSemanticAnalyzer *DeepSemanticAnalyzer    // Analyseur sémantique profond
	HybridEngine         *HybridAtomicProbability // Moteur hybride probabilité+stabilité
	Knowledge            *LightweightKB           // Base de connaissances factuelles
}

// NewHellaswagEngine crée un nouveau moteur optimisé
func NewHellaswagEngine() *HellaswagEngine {
	return &HellaswagEngine{
		perplexityCalc:       NewPerplexityCalculator(),
		ContextMemory:        make(map[string][]string),
		TemporalPatterns:     initializeTemporalPatterns(),
		CausalPatterns:       initializeCausalPatterns(),
		AdaptiveLearning:     true,
		ContextEngine:        NewContextEngine(),
		DeepSemanticAnalyzer: NewDeepSemanticAnalyzer(),
		HybridEngine:         NewHybridAtomicProbability(),
		Knowledge:            LoadDefaultKB(),
	}
}

// initializeTemporalPatterns initialise patterns temporels
func initializeTemporalPatterns() map[string]float64 {
	return map[string]float64{
		"puis":            0.8,
		"ensuite":         0.8,
		"après":           0.7,
		"avant":           0.7,
		"pendant":         0.6,
		"enfin":           0.9,
		"maintenant":      0.7,
		"alors":           0.8,
		"finalement":      0.85,
		"progressivement": 0.75,
	}
}

// initializeCausalPatterns initialise patterns causals
func initializeCausalPatterns() map[string]float64 {
	return map[string]float64{
		"donc":        0.9,
		"car":         0.8,
		"parce que":   0.85,
		"à cause de":  0.8,
		"grâce à":     0.75,
		"pour":        0.7,
		"afin de":     0.75,
		"résultat":    0.85,
		"conséquence": 0.9,
	}
}

// EvaluateQuestion évalue avec analyse multi-dimensionnelle
func (engine *HellaswagEngine) EvaluateQuestion(question *HellaswagQuestion) *HellaswagResult {
	result := &HellaswagResult{
		Question:     question,
		EndingScores: make([]float64, len(question.Endings)),
	}

	// Analyse profonde du contexte avec enrichissement
	contextTokens := TokeniserTexte(question.Context)
	contextCategories := ActiverCategoriesParTexte(contextTokens)
	contextKeywords := ExtraireMotsClés(contextTokens)

	// NOUVEAU: Contexte enrichi multi-niveau
	contextEnriched := engine.ContextEngine.EnrichContext(question.Context, contextCategories)
	expandedContext := engine.ContextEngine.ExpandContextWithHistory(contextCategories, 0.4)

	bestScore := -999.0
	bestIndex := 0
	minPerplexity := 999.0
	maxPerplexity := 0.0

	// Évaluer chaque fin avec 8 critères (avec contexte enrichi)
	for i, ending := range question.Endings {
		fullText := question.Context + " " + ending

		// 1. Perplexité du texte complet
		perplexityResult := engine.perplexityCalc.CalculatePerplexity(fullText)
		perplexity := perplexityResult.GlobalPerplexity

		endingTokens := TokeniserTexte(ending)
		endingCategories := ActiverCategoriesParTexte(endingTokens)
		endingKeywords := ExtraireMotsClés(endingTokens)

		// NOUVEAU: Contexte enrichi pour chaque fin
		endingEnriched := engine.ContextEngine.EnrichContext(ending, endingCategories)

		// 2. Cohérence sémantique avancée (avec contexte étendu)
		coherence := engine.calculateAdvancedCoherence(expandedContext, endingCategories)

		// 3. NOUVEAU: Contexte enrichi (n-grammes, concepts, profondeur narrative)
		contextScore := engine.ContextEngine.CompareEnrichedContexts(contextEnriched, endingEnriched)

		// 3.5 NOUVEAU: Analyse sémantique profonde (cause-effet, séquences logiques)
		deepSemanticScore := engine.DeepSemanticAnalyzer.AnalyzeSemanticCoherence(question.Context, ending)

		// 4. Continuité lexicale (mots en commun)
		lexicalContinuity := engine.analyzeLexicalContinuity(contextKeywords, endingKeywords)

		// 5. Patterns temporels (puis, ensuite, après...)
		temporalScore := engine.detectTemporalPatterns(ending)

		// 6. Patterns causaux (donc, parce que...)
		causalScore := engine.detectCausalPatterns(ending)

		// 7. Cohérence d'action (verbes compatibles)
		actionCoherence := engine.analyzeActionCoherence(question.Context, ending)

		// 7bis. Continuité des verbes d'action (réutilisation du même geste)
		verbContinuity := engine.actionVerbContinuity(question.Context, ending)

		// 7ter. Alignement avec l'activité attendue
		activityScore := engine.activityAlignment(question.Activity, ending)

		// 8. Fluidité narrative (transitions naturelles)
		narrativeFlow := engine.assessNarrativeFlow(fullText)

		// 8bis. Pénalité de rupture de sujet (fin qui introduit un thème étranger)
		topicShift := engine.topicShiftPenalty(question.Context, ending)

		// 8ter. Pénalité d'absurdité (actions manifestement hors-sujet)
		absurdPenalty := engine.absurdityPenalty(question.Activity, ending)

		// Score de perplexité inversé (plus bas = meilleur)
		perplexityScore := 1.0 / (perplexity + 0.05)

		// 9. NOUVEAU: Score HYBRIDE (Probabilité + Stabilité Atomique)
		hybridScore := engine.HybridEngine.HybridScore(question.Context, ending)

		// 9bis. BOOST CONNAISSANCES FACTUELLES (si KB disponible)
		kbBoost := 0.0
		if engine.Knowledge != nil {
			kbBoost = engine.Knowledge.ConfidenceBoost(question.Context + " " + ending)
			if kbBoost > 0.35 {
				kbBoost = 0.35
			}
		}

		// Score BASE (sans contexte enrichi) - poids originaux optimaux
		baseScore := (perplexityScore*0.18 +
			coherence*0.16 +
			deepSemanticScore*0.18 +
			hybridScore*0.20 + // NOUVEAU: 20% pour score hybride!
			lexicalContinuity*0.10 +
			temporalScore*0.10 +
			causalScore*0.06 +
			actionCoherence*0.04 +
			verbContinuity*0.08 + // Favorise la continuité du geste amorcé
			activityScore*0.14 + // Aligne la fin avec l'activité déclarée
			narrativeFlow*0.01)

		// Pénalité douce si la fin part sur un autre sujet (limite la divagation)
		if topicShift > 0 {
			baseScore -= topicShift * 0.10
		}

		// Pénalité si action absurde / hors-activité
		if absurdPenalty > 0 {
			baseScore -= absurdPenalty * 0.12
		}
		if engine.Knowledge != nil {
			factBonus := engine.Knowledge.FactBonus(question.Context + " " + ending)
			baseScore += factBonus
		}
		// Appliquer boost KB (mix additif + multiplicatif)
		if kbBoost > 0 {
			baseScore += 0.04
			baseScore *= (1 + kbBoost)
		}
		// BOOST du contexte enrichi (multiplicatif, jusqu'à +12%)
		contextBoost := 1.0 + (contextScore * 0.12)

		// Score final avec contexte en boost multiplicatif
		finalScore := baseScore * contextBoost

		// Réinitialiser réseau atomique pour prochain test
		engine.HybridEngine.ResetNetwork()

		result.EndingScores[i] = finalScore

		if perplexity < minPerplexity {
			minPerplexity = perplexity
		}
		if perplexity > maxPerplexity {
			maxPerplexity = perplexity
		}

		if finalScore > bestScore {
			bestScore = finalScore
			bestIndex = i
		}
	}

	// Apprentissage adaptatif
	if engine.AdaptiveLearning {
		engine.updateContextMemory(question.Context, question.Endings[question.Answer]) // NOUVEAU: Mise à jour historique contextuel
		engine.ContextEngine.UpdateHistory(contextCategories)
	}

	result.SelectedIndex = bestIndex
	result.IsCorrect = (bestIndex == question.Answer)
	result.Confidence = bestScore
	result.PerplexityDiff = maxPerplexity - minPerplexity

	return result
}

// calculateAdvancedCoherence cohérence avancée avec pondération
func (engine *HellaswagEngine) calculateAdvancedCoherence(cat1, cat2 map[int]int) float64 {
	if len(cat1) == 0 || len(cat2) == 0 {
		return 0.0
	}

	dotProduct := 0.0
	for cat, val1 := range cat1 {
		if val2, exists := cat2[cat]; exists {
			categoryWeight := 1.0
			if cat >= 1 && cat <= 6 {
				categoryWeight = 1.4 // Catégories d'action plus importantes
			}
			dotProduct += float64(val1*val2) * categoryWeight
		}
	}

	norm1 := 0.0
	for cat, val := range cat1 {
		categoryWeight := 1.0
		if cat >= 1 && cat <= 6 {
			categoryWeight = 1.4
		}
		norm1 += float64(val*val) * categoryWeight
	}
	norm1 = math.Sqrt(norm1)

	norm2 := 0.0
	for cat, val := range cat2 {
		categoryWeight := 1.0
		if cat >= 1 && cat <= 6 {
			categoryWeight = 1.4
		}
		norm2 += float64(val*val) * categoryWeight
	}
	norm2 = math.Sqrt(norm2)

	if norm1 == 0 || norm2 == 0 {
		return 0.0
	}

	cosineSim := dotProduct / (norm1 * norm2)

	// Bonus pour continuité de domaine
	overlapCount := 0
	for cat := range cat1 {
		if _, exists := cat2[cat]; exists {
			overlapCount++
		}
	}
	overlapBonus := math.Min(float64(overlapCount)*0.08, 0.4)

	return math.Min(cosineSim+overlapBonus, 1.0)
}

// analyzeLexicalContinuity analyse la continuité lexicale
func (engine *HellaswagEngine) analyzeLexicalContinuity(contextKW, endingKW []string) float64 {
	if len(contextKW) == 0 || len(endingKW) == 0 {
		return 0.0
	}

	kwMap := make(map[string]bool)
	for _, kw := range contextKW {
		kwMap[strings.ToLower(kw)] = true
	}

	matchCount := 0
	for _, kw := range endingKW {
		if kwMap[strings.ToLower(kw)] {
			matchCount++
		}
	}

	// Pénaliser trop de répétition (redondance)
	overlapRatio := float64(matchCount) / float64(len(endingKW))
	if overlapRatio > 0.7 {
		return 0.5 // Redondant
	} else if overlapRatio > 0.3 && overlapRatio <= 0.7 {
		return 0.8 // Équilibré
	} else if overlapRatio > 0.1 {
		return 0.6 // Peu de lien
	}
	return 0.3 // Très peu de lien
}

// detectTemporalPatterns détecte les marqueurs temporels
func (engine *HellaswagEngine) detectTemporalPatterns(text string) float64 {
	lowerText := strings.ToLower(text)
	score := 0.0
	count := 0

	for pattern, weight := range engine.TemporalPatterns {
		if strings.Contains(lowerText, pattern) {
			score += weight
			count++
		}
	}

	if count > 0 {
		return math.Min(score/float64(count), 1.0)
	}
	return 0.3 // Pas de marqueur temporel
}

// detectCausalPatterns détecte les relations causales
func (engine *HellaswagEngine) detectCausalPatterns(text string) float64 {
	lowerText := strings.ToLower(text)
	score := 0.0
	count := 0

	for pattern, weight := range engine.CausalPatterns {
		if strings.Contains(lowerText, pattern) {
			score += weight
			count++
		}
	}

	if count > 0 {
		return math.Min(score/float64(count), 1.0)
	}
	return 0.4 // Pas de causalité explicite
}

// analyzeActionCoherence analyse la cohérence des actions (verbes)
func (engine *HellaswagEngine) analyzeActionCoherence(context, ending string) float64 {
	contextVerbs := extractVerbs(context)
	endingVerbs := extractVerbs(ending)

	if len(contextVerbs) == 0 || len(endingVerbs) == 0 {
		return 0.5
	}

	// Vérifier compatibilité des actions
	score := 0.0

	// Actions progressives (commencer → continuer → finir)
	if containsAction(contextVerbs, []string{"commence", "démarre", "entame"}) &&
		containsAction(endingVerbs, []string{"continue", "poursuit", "finit"}) {
		score += 0.8
	}

	// Actions contradictoires (monter vs descendre)
	if detectContradiction(contextVerbs, endingVerbs) {
		score = 0.1
	} else {
		score += 0.5
	}

	return math.Min(score, 1.0)
}

// assessNarrativeFlow évalue la fluidité narrative globale
func (engine *HellaswagEngine) assessNarrativeFlow(fullText string) float64 {
	tokens := TokeniserTexte(fullText)
	if len(tokens) == 0 {
		return 0.0
	}

	// 1. Longueur appropriée
	lengthScore := 0.0
	if len(tokens) >= 15 && len(tokens) <= 50 {
		lengthScore = 0.8
	} else if len(tokens) > 50 {
		lengthScore = 0.5
	} else {
		lengthScore = 0.4
	}

	// 2. Diversité lexicale
	uniqueWords := make(map[string]bool)
	for _, token := range tokens {
		uniqueWords[strings.ToLower(token)] = true
	}
	diversity := float64(len(uniqueWords)) / float64(len(tokens))

	// 3. Présence de connecteurs
	hasConnectors := containsConnectors(fullText)
	connectorScore := 0.0
	if hasConnectors {
		connectorScore = 0.7
	} else {
		connectorScore = 0.4
	}

	// Score final
	return (lengthScore*0.4 + diversity*0.3 + connectorScore*0.3)
}

// updateContextMemory met à jour la mémoire contextuelle
func (engine *HellaswagEngine) updateContextMemory(context, correctEnding string) {
	maxLen := 50
	if len(context) < maxLen {
		maxLen = len(context)
	}
	contextKey := strings.ToLower(context[:maxLen])

	if _, exists := engine.ContextMemory[contextKey]; !exists {
		engine.ContextMemory[contextKey] = []string{}
	}

	maxEndLen := 50
	if len(correctEnding) < maxEndLen {
		maxEndLen = len(correctEnding)
	}

	engine.ContextMemory[contextKey] = append(
		engine.ContextMemory[contextKey],
		strings.ToLower(correctEnding[:maxEndLen]),
	)

	// Limiter la mémoire à 100 entrées
	if len(engine.ContextMemory) > 100 {
		// Supprimer entrée la plus ancienne
		for k := range engine.ContextMemory {
			delete(engine.ContextMemory, k)
			break
		}
	}
}

// Helper functions
func extractVerbs(text string) []string {
	commonVerbs := []string{
		"prend", "met", "fait", "va", "vient", "sort", "entre",
		"commence", "finit", "continue", "regarde", "écoute",
		"parle", "mange", "boit", "court", "marche", "monte",
		"descend", "ouvre", "ferme", "attend", "reste",
	}

	lowerText := strings.ToLower(text)
	var foundVerbs []string

	for _, verb := range commonVerbs {
		if strings.Contains(lowerText, verb) {
			foundVerbs = append(foundVerbs, verb)
		}
	}

	return foundVerbs
}

func containsAction(verbs []string, actions []string) bool {
	for _, verb := range verbs {
		for _, action := range actions {
			if strings.Contains(verb, action) || strings.Contains(action, verb) {
				return true
			}
		}
	}
	return false
}

func detectContradiction(verbs1, verbs2 []string) bool {
	contradictions := map[string][]string{
		"monte":    {"descend"},
		"entre":    {"sort"},
		"ouvre":    {"ferme"},
		"commence": {"finit", "termine"},
		"allume":   {"éteint"},
	}

	for _, v1 := range verbs1 {
		if opposites, exists := contradictions[v1]; exists {
			for _, v2 := range verbs2 {
				for _, opposite := range opposites {
					if strings.Contains(v2, opposite) {
						return true
					}
				}
			}
		}
	}

	return false
}

func containsConnectors(text string) bool {
	connectors := []string{
		"puis", "ensuite", "alors", "donc", "car",
		"parce", "ainsi", "enfin", "finalement",
	}

	lowerText := strings.ToLower(text)
	for _, conn := range connectors {
		if strings.Contains(lowerText, conn) {
			return true
		}
	}

	return false
}

// RunBenchmark exécute le benchmark complet
func (engine *HellaswagEngine) RunBenchmark(questions []*HellaswagQuestion) *HellaswagBenchmark {
	benchmark := &HellaswagBenchmark{
		TotalQuestions: len(questions),
		Results:        make([]*HellaswagResult, 0),
	}

	totalConfidence := 0.0
	totalPerplexity := 0.0

	for _, question := range questions {
		result := engine.EvaluateQuestion(question)
		benchmark.Results = append(benchmark.Results, result)

		if result.IsCorrect {
			benchmark.CorrectAnswers++
		}
		totalConfidence += result.Confidence
		totalPerplexity += result.PerplexityDiff
	}

	if benchmark.TotalQuestions > 0 {
		benchmark.Score = float64(benchmark.CorrectAnswers) / float64(benchmark.TotalQuestions) * 100
		benchmark.AverageConfidence = totalConfidence / float64(benchmark.TotalQuestions)
		benchmark.AveragePerplexity = totalPerplexity / float64(benchmark.TotalQuestions)
	}

	return benchmark
}

// CreateSampleHellaswag crée un échantillon de questions
func CreateSampleHellaswag() []*HellaswagQuestion {
	questions := []*HellaswagQuestion{
		{
			ID:       1,
			Context:  "Une femme entre dans une cuisine. Elle prend une casserole et la remplit d'eau.",
			Activity: "cuisine",
			Endings: []string{
				"Elle commence à jouer au football dans le salon.",
				"Elle met la casserole sur le feu et attend que l'eau bouille.",
				"Elle va au cinéma avec des amis.",
				"Elle ouvre une fenêtre et regarde le ciel.",
			},
			Answer: 1, // Suite logique de cuisine
		},
		{
			ID:       2,
			Context:  "Un homme met ses chaussures de sport et sort de chez lui. Il commence à courir dans la rue.",
			Activity: "sport",
			Endings: []string{
				"Il s'arrête devant un magasin et achète des légumes.",
				"Il rentre chez lui et prend un bain chaud.",
				"Il continue de courir pendant 30 minutes puis ralentit progressivement.",
				"Il téléphone à son avocat pour discuter d'un contrat.",
			},
			Answer: 2, // Suite logique de jogging
		},
		{
			ID:       3,
			Context:  "Un étudiant ouvre son ordinateur et lance un traitement de texte. Il commence à taper.",
			Activity: "travail",
			Endings: []string{
				"Il saute par la fenêtre avec un parachute.",
				"Il rédige son rapport pendant deux heures puis le sauvegarde.",
				"Il va nager dans la piscine municipale.",
				"Il mange une pizza surgelée.",
			},
			Answer: 1, // Suite logique d'écriture
		},
		{
			ID:       4,
			Context:  "Une fille prend son vélo dans le garage. Elle met son casque.",
			Activity: "sport",
			Endings: []string{
				"Elle part faire une balade à vélo dans le parc.",
				"Elle cuisine un gâteau au chocolat.",
				"Elle répare la télévision du salon.",
				"Elle lit un livre de philosophie.",
			},
			Answer: 0, // Suite logique de vélo
		},
		{
			ID:       5,
			Context:  "Un médecin entre dans la salle d'examen. Il salue le patient et sort son stéthoscope.",
			Activity: "médecine",
			Endings: []string{
				"Il commence à danser le tango.",
				"Il écoute les battements du cœur et les poumons du patient.",
				"Il va au supermarché acheter du pain.",
				"Il repeint les murs de son bureau.",
			},
			Answer: 1, // Suite logique d'examen médical
		},
		{
			ID:       6,
			Context:  "Un chef cuisinier prépare une omelette. Il casse trois œufs dans un bol.",
			Activity: "cuisine",
			Endings: []string{
				"Il bat les œufs avec une fourchette puis verse le mélange dans une poêle.",
				"Il va réparer sa voiture dans le garage.",
				"Il écrit un article de journal.",
				"Il regarde un match de football à la télévision.",
			},
			Answer: 0, // Suite logique de cuisine
		},
		{
			ID:       7,
			Context:  "Une personne entre dans une bibliothèque. Elle cherche un livre dans les rayons.",
			Activity: "lecture",
			Endings: []string{
				"Elle trouve le livre, s'assoit à une table et commence à lire.",
				"Elle fait du skateboard dans les allées.",
				"Elle plante des arbres dans le jardin.",
				"Elle nage dans la piscine intérieure.",
			},
			Answer: 0, // Suite logique de bibliothèque
		},
		{
			ID:       8,
			Context:  "Un mécanicien ouvre le capot d'une voiture. Il regarde le moteur avec attention.",
			Activity: "mécanique",
			Endings: []string{
				"Il vérifie le niveau d'huile et inspecte les courroies.",
				"Il joue au tennis sur un court extérieur.",
				"Il fait cuire un rôti dans le four.",
				"Il étudie les mathématiques avancées.",
			},
			Answer: 0, // Suite logique de mécanique
		},
		{
			ID:       9,
			Context:  "Une personne allume son téléphone et ouvre l'application météo.",
			Activity: "quotidien",
			Endings: []string{
				"Elle consulte les prévisions pour la semaine et décide de ses vêtements.",
				"Elle va escalader une montagne.",
				"Elle répare un ordinateur.",
				"Elle écrit un poème romantique.",
			},
			Answer: 0, // Suite logique de consultation météo
		},
		{
			ID:       10,
			Context:  "Un jardinier prend ses outils et sort dans le jardin. Il se met à genoux près des fleurs.",
			Activity: "jardinage",
			Endings: []string{
				"Il arrache les mauvaises herbes et arrose les plantes.",
				"Il regarde un film d'action.",
				"Il fait ses devoirs de mathématiques.",
				"Il prépare un exposé sur l'histoire.",
			},
			Answer: 0, // Suite logique de jardinage
		},
	}

	return questions
}

// actionVerbContinuity mesure si la fin réemploie les mêmes gestes que le contexte
func (engine *HellaswagEngine) actionVerbContinuity(context, ending string) float64 {
	verbs := []string{
		"prendre", "mettre", "sortir", "remplir", "verser", "chauffer", "porter", "soulever",
		"ouvrir", "fermer", "cuisiner", "mixer", "couper", "marcher", "courir", "sauter",
		"attraper", "lancer", "respirer", "travailler", "écrire", "taper", "manger", "boire",
	}
	verbSet := make(map[string]struct{}, len(verbs))
	for _, v := range verbs {
		verbSet[v] = struct{}{}
	}

	contextTokens := TokeniserTexte(context)
	endingTokens := TokeniserTexte(ending)
	ctxVerbCount := 0
	for _, tok := range contextTokens {
		tok = strings.ToLower(tok)
		if _, ok := verbSet[tok]; ok {
			ctxVerbCount++
		}
	}

	if ctxVerbCount == 0 {
		return 0
	}

	continuity := 0
	for _, tok := range endingTokens {
		tok = strings.ToLower(tok)
		if _, ok := verbSet[tok]; ok {
			continuity++
		}
	}

	score := float64(continuity) / float64(ctxVerbCount)
	if score > 1 {
		score = 1
	}
	return score
}

// topicShiftPenalty pénalise les fins qui introduisent trop de nouveau vocabulaire
func (engine *HellaswagEngine) topicShiftPenalty(context, ending string) float64 {
	contextTokens := TokeniserTexte(context)
	endingTokens := TokeniserTexte(ending)
	contextSet := make(map[string]struct{}, len(contextTokens))
	for _, tok := range contextTokens {
		tok = strings.ToLower(tok)
		if len(tok) <= 3 {
			continue
		}
		contextSet[tok] = struct{}{}
	}

	novel := 0
	valid := 0
	for _, tok := range endingTokens {
		tok = strings.ToLower(tok)
		if len(tok) <= 3 {
			continue
		}
		valid++
		if _, ok := contextSet[tok]; !ok {
			novel++
		}
	}

	if valid == 0 {
		return 0
	}
	ratio := float64(novel) / float64(valid)
	if ratio < 0.3 {
		return 0
	}
	if ratio > 0.6 {
		ratio = 0.6
	}
	return ratio
}

// activityAlignment favorise les fins qui contiennent des actions liées à l'activité déclarée
func (engine *HellaswagEngine) activityAlignment(activity, ending string) float64 {
	keywords := map[string][]string{
		"cuisine":      {"cuisine", "poêle", "casserole", "cuire", "verser", "chauffer", "omelette", "œufs", "couteau", "fourchette"},
		"sport":        {"courir", "vélo", "balade", "casque", "jogging", "chaussures", "parc"},
		"travail":      {"écrire", "rapport", "document", "ordinateur", "sauvegarde", "taper"},
		"médecine":     {"stéthoscope", "cœur", "poumons", "patient", "examiner", "écouter"},
		"bibliothèque": {"livre", "lire", "rayons", "table", "silence"},
		"mécanique":    {"moteur", "huile", "courroies", "capot", "vérifier"},
		"quotidien":    {"prévisions", "météo", "vêtements", "consulter"},
		"jardinage":    {"arroser", "fleurs", "mauvaises", "herbes", "plantes", "outil"},
		"lecture":      {"livre", "lire", "page", "rayons"},
	}

	activity = strings.ToLower(activity)
	list, ok := keywords[activity]
	if !ok {
		return 0
	}
	endingTokens := TokeniserTexte(ending)
	set := make(map[string]struct{}, len(list))
	for _, w := range list {
		set[w] = struct{}{}
	}
	matches := 0
	for _, tok := range endingTokens {
		tok = strings.ToLower(tok)
		if _, ok := set[tok]; ok {
			matches++
		}
	}
	if len(list) == 0 {
		return 0
	}
	score := float64(matches) / float64(len(list))
	if score > 1 {
		score = 1
	}
	return score
}

// absurdityPenalty pénalise les fins avec des actions hors-sujet pour l'activité
func (engine *HellaswagEngine) absurdityPenalty(activity, ending string) float64 {
	absurdWords := []string{
		"parachute", "tango", "danser", "escalader", "skateboard", "piscine", "tennis", "peindre", "poème", "film", "supermarché", "voiture",
	}
	allowed := map[string]struct{}{}
	for _, w := range absurdWords {
		allowed[w] = struct{}{}
	}

	// Autorisations spécifiques par activité (évite faux positifs)
	if strings.ToLower(activity) == "sport" {
		delete(allowed, "tennis")
		delete(allowed, "piscine")
		delete(allowed, "skateboard")
	}

	endingTokens := TokeniserTexte(ending)
	penalty := 0.0
	for _, tok := range endingTokens {
		tok = strings.ToLower(tok)
		if _, ok := allowed[tok]; ok {
			penalty += 0.15
		}
	}
	if penalty > 0.6 {
		penalty = 0.6
	}
	return penalty
}

// PrintResults affiche les résultats
func (benchmark *HellaswagBenchmark) PrintResults() {
	sep := strings.Repeat("=", 70)
	fmt.Println("\n" + sep)
	fmt.Println("  RÉSULTATS BENCHMARK HELLASWAG")
	fmt.Println(sep)

	fmt.Printf("\n[SCORE GLOBAL]\n")
	fmt.Printf("  • Questions: %d\n", benchmark.TotalQuestions)
	fmt.Printf("  • Réponses correctes: %d\n", benchmark.CorrectAnswers)
	fmt.Printf("  • Score: %.2f%%\n", benchmark.Score)
	fmt.Printf("  • Confiance moyenne: %.3f\n", benchmark.AverageConfidence)
	fmt.Printf("  • Écart perplexité moyen: %.3f\n", benchmark.AveragePerplexity)

	// Quelques exemples
	fmt.Printf("\n[EXEMPLES DE RAISONNEMENT]\n")
	maxExamples := 3
	if len(benchmark.Results) < maxExamples {
		maxExamples = len(benchmark.Results)
	}

	for i := 0; i < maxExamples; i++ {
		result := benchmark.Results[i]
		q := result.Question

		status := "✗"
		if result.IsCorrect {
			status = "✓"
		}

		fmt.Printf("\n[Q%d] %s\n", q.ID, status)
		fmt.Printf("  Contexte: %s\n", q.Context)
		fmt.Printf("  Suite choisie: %s\n", q.Endings[result.SelectedIndex])
		fmt.Printf("  Confiance: %.3f | Écart perplexité: %.3f\n",
			result.Confidence, result.PerplexityDiff)
		if !result.IsCorrect {
			fmt.Printf("  Suite correcte: %s\n", q.Endings[q.Answer])
		}
	}

	fmt.Println("\n" + sep + "\n")
}
