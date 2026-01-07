package database

import (
	"fmt"
	"math/rand"
	"strings"
)

// ═══════════════════════════════════════════════════════════════════════════
// PHASE 14: Système de Génération Grammaticale Avancée avec Renforcement
// ═══════════════════════════════════════════════════════════════════════════

// ─────────────────────────────────────────────────────────────────────────
// PARTIE 1: POS Tagging (Part-of-Speech) pour Français
// ─────────────────────────────────────────────────────────────────────────

type POSTag string

const (
	POS_NOUN    POSTag = "NOUN"  // Nom
	POS_VERB    POSTag = "VERB"  // Verbe
	POS_ADJ     POSTag = "ADJ"   // Adjectif
	POS_ADV     POSTag = "ADV"   // Adverbe
	POS_DET     POSTag = "DET"   // Déterminant (le, la, un, une)
	POS_PREP    POSTag = "PREP"  // Préposition (de, à, en, pour)
	POS_CCONJ   POSTag = "CCONJ" // Conjonction de coordination (et, ou, mais)
	POS_SCONJ   POSTag = "SCONJ" // Conjonction de subordination (que, si, parce que)
	POS_PRON    POSTag = "PRON"  // Pronom (je, tu, il, elle, on)
	POS_PUNCT   POSTag = "PUNCT" // Ponctuation
	POS_UNKNOWN POSTag = "X"     // Inconnu
)

// TaggedWord représente un mot avec son étiquette POS
type TaggedWord struct {
	Word  string
	POS   POSTag
	Lemma string // Forme canonique du mot
	Index int
	Score float64
}

// DependencyRelation représente une relation de dépendance syntaxique
type DependencyRelation struct {
	Governor   int    // Index du mot gouvernant
	Dependent  int    // Index du mot dépendant
	RelType    string // Type de relation (sujet, objet, attribut, etc.)
	Confidence float64
}

// ParsedSentence représente une phrase analysée
type ParsedSentence struct {
	Original     string
	Words        []TaggedWord
	Dependencies []DependencyRelation
	MainVerb     *TaggedWord
	Subject      *TaggedWord
	Objects      []*TaggedWord
	Modifiers    []*TaggedWord
	GrammarScore float64 // 0.0 à 1.0
	StyleScore   float64 // Vivacité/variation
	FinalScore   float64 // Grammar × (1 + 0.5×Style)
}

// ─────────────────────────────────────────────────────────────────────────
// PARTIE 2: Dictionnaires POS pour Français
// ─────────────────────────────────────────────────────────────────────────

var FrenchPOSDictionary = map[string]POSTag{
	// Articles/Déterminants
	"le": POS_DET, "la": POS_DET, "les": POS_DET,
	"un": POS_DET, "une": POS_DET, "des": POS_DET,
	"ce": POS_DET, "cette": POS_DET, "cet": POS_DET, "ces": POS_DET,
	"mon": POS_DET, "ma": POS_DET, "mes": POS_DET,
	"ton": POS_DET, "ta": POS_DET, "tes": POS_DET,
	"son": POS_DET, "sa": POS_DET, "ses": POS_DET,
	"notre": POS_DET, "nos": POS_DET,
	"votre": POS_DET, "vos": POS_DET,
	"leur": POS_DET, "leurs": POS_DET,

	// Pronoms
	"je": POS_PRON, "tu": POS_PRON, "il": POS_PRON, "elle": POS_PRON,
	"on": POS_PRON, "nous": POS_PRON, "vous": POS_PRON, "ils": POS_PRON, "elles": POS_PRON,
	"me": POS_PRON, "te": POS_PRON, "se": POS_PRON,
	"moi": POS_PRON, "toi": POS_PRON, "lui": POS_PRON, "eux": POS_PRON,

	// Prépositions communes
	"de": POS_PREP, "à": POS_PREP, "en": POS_PREP, "pour": POS_PREP,
	"par": POS_PREP, "avec": POS_PREP, "sans": POS_PREP, "sous": POS_PREP,
	"sur": POS_PREP, "entre": POS_PREP, "chez": POS_PREP, "vers": POS_PREP,
	"dans": POS_PREP, "pendant": POS_PREP, "avant": POS_PREP, "après": POS_PREP,

	// Conjonctions de coordination
	"et": POS_CCONJ, "ou": POS_CCONJ, "mais": POS_CCONJ,
	"donc": POS_CCONJ, "car": POS_CCONJ, "ni": POS_CCONJ,

	// Conjonctions de subordination
	"que": POS_SCONJ, "si": POS_SCONJ, "parce": POS_SCONJ,
	"quand": POS_SCONJ, "où": POS_SCONJ, "comment": POS_SCONJ,
	"pourquoi": POS_SCONJ, "bien": POS_SCONJ,

	// Verbes courants (forme infinitive)
	"être": POS_VERB, "avoir": POS_VERB, "aller": POS_VERB,
	"faire": POS_VERB, "pouvoir": POS_VERB, "vouloir": POS_VERB,
	"devoir": POS_VERB, "falloir": POS_VERB, "savoir": POS_VERB,
}

var FrenchLemmas = map[string]string{
	// Conjugaisons vers infinitif
	"suis": "être", "es": "être", "est": "être",
	"sommes": "être", "êtes": "être", "sont": "être",
	"ai": "avoir", "as": "avoir", "avons": "avoir", "avez": "avoir", "ont": "avoir",
	"vais": "aller", "vas": "aller", "allons": "aller", "allez": "aller", "vont": "aller",
	"fais": "faire", "fait": "faire", "faisons": "faire", "faites": "faire", "font": "faire",
}

// ─────────────────────────────────────────────────────────────────────────
// PARTIE 3: Tagging POS avec Heuristiques
// ─────────────────────────────────────────────────────────────────────────

// TagMots effectue le tagging POS d'une phrase
func TagMots(phrase string) []TaggedWord {
	words := strings.Fields(phrase)
	tagged := []TaggedWord{}

	for i, word := range words {
		wordClean := strings.ToLower(strings.TrimSuffix(strings.TrimSuffix(word, "."), ","))

		// Chercher dans dictionnaire
		if pos, exists := FrenchPOSDictionary[wordClean]; exists {
			lemma := wordClean
			if l, hasLemma := FrenchLemmas[wordClean]; hasLemma {
				lemma = l
			}
			tagged = append(tagged, TaggedWord{
				Word:  word,
				POS:   pos,
				Lemma: lemma,
				Index: i,
				Score: 0.95,
			})
			continue
		}

		// Heuristiques pour mots inconnus
		pos := ClassifyUnknownWord(wordClean)
		tagged = append(tagged, TaggedWord{
			Word:  word,
			POS:   pos,
			Lemma: wordClean,
			Index: i,
			Score: 0.6, // Confiance plus basse pour mots inconnus
		})
	}

	return tagged
}

// ClassifyUnknownWord utilise l'heuristique pour classer un mot inconnu
func ClassifyUnknownWord(word string) POSTag {
	// Suffixes typiques
	if strings.HasSuffix(word, "tion") || strings.HasSuffix(word, "sion") {
		return POS_NOUN
	}
	if strings.HasSuffix(word, "ment") {
		return POS_ADV
	}
	if strings.HasSuffix(word, "eux") || strings.HasSuffix(word, "eux") {
		return POS_ADJ
	}
	if strings.HasSuffix(word, "er") || strings.HasSuffix(word, "ir") {
		return POS_VERB
	}
	if strings.HasSuffix(word, ".") || strings.HasSuffix(word, ",") {
		return POS_PUNCT
	}

	// Par défaut: nom (stratégie conservative)
	return POS_NOUN
}

// ─────────────────────────────────────────────────────────────────────────
// PARTIE 4: Analyse de Dépendances Syntaxiques
// ─────────────────────────────────────────────────────────────────────────

// AnalyzeDependencies identifie les relations syntaxiques
func AnalyzeDependencies(tagged []TaggedWord) ParsedSentence {
	parsed := ParsedSentence{
		Words:        tagged,
		Dependencies: []DependencyRelation{},
		Objects:      []*TaggedWord{},
		Modifiers:    []*TaggedWord{},
	}

	// Trouver le verbe principal
	mainVerbIdx := -1
	for i, w := range tagged {
		if w.POS == POS_VERB {
			mainVerbIdx = i
			parsed.MainVerb = &tagged[i]
			break
		}
	}

	if mainVerbIdx == -1 {
		parsed.GrammarScore = 0.3 // Phrase sans verbe = grammaticalement faible
		return parsed
	}

	// Identifier le sujet (pronom/nom avant verbe)
	for i := mainVerbIdx - 1; i >= 0; i-- {
		w := tagged[i]
		if w.POS == POS_NOUN || w.POS == POS_PRON {
			parsed.Subject = &tagged[i]
			parsed.Dependencies = append(parsed.Dependencies, DependencyRelation{
				Governor:   mainVerbIdx,
				Dependent:  i,
				RelType:    "nsubj", // sujet nominal
				Confidence: 0.9,
			})
			break
		}
	}

	// Identifier les objets (noms après verbe)
	for i := mainVerbIdx + 1; i < len(tagged); i++ {
		w := tagged[i]
		if w.POS == POS_NOUN {
			parsed.Objects = append(parsed.Objects, &tagged[i])
			parsed.Dependencies = append(parsed.Dependencies, DependencyRelation{
				Governor:   mainVerbIdx,
				Dependent:  i,
				RelType:    "obj", // objet
				Confidence: 0.85,
			})
		} else if w.POS == POS_ADJ || w.POS == POS_ADV {
			// Adjectifs/adverbes modifient le verbe ou les noms
			parsed.Modifiers = append(parsed.Modifiers, &tagged[i])
		}
	}

	// Calculer score grammatical
	parsed.GrammarScore = CalculateGrammarScore(parsed)

	return parsed
}

// ─────────────────────────────────────────────────────────────────────────
// PARTIE 5: Scoring Grammatical & Stylistique
// ─────────────────────────────────────────────────────────────────────────

// CalculateGrammarScore évalue la validité grammaticale
func CalculateGrammarScore(parsed ParsedSentence) float64 {
	score := 1.0

	// Critère 1: Présence de verbe principal (-0.3 si absent)
	if parsed.MainVerb == nil {
		score -= 0.3
	}

	// Critère 2: Cohérence sujet-verbe (-0.2 si sujet absent)
	if parsed.Subject == nil && parsed.MainVerb != nil {
		score -= 0.2
	}

	// Critère 3: Présence d'au moins un complément (-0.15 si absent)
	if len(parsed.Objects) == 0 && parsed.MainVerb != nil {
		score -= 0.15
	}

	// Critère 4: Ordre des mots (SVO = Subject-Verb-Object)
	if parsed.Subject != nil && parsed.MainVerb != nil && len(parsed.Objects) > 0 {
		subjIdx := parsed.Subject.Index
		verbIdx := parsed.MainVerb.Index
		objIdx := parsed.Objects[0].Index

		// Ordre correct: sujet < verbe < objet
		if subjIdx < verbIdx && verbIdx < objIdx {
			score += 0.2 // Bonus pour ordre SVO naturel
		}
	}

	// Clamp entre 0 et 1
	if score < 0 {
		score = 0
	}
	if score > 1 {
		score = 1
	}

	return score
}

// CalculateStyleScore évalue la vivacité et variation
func CalculateStyleScore(parsed ParsedSentence) float64 {
	score := 0.5 // Base neutre

	// +0.1 par modificateur (adjectifs/adverbes)
	score += float64(len(parsed.Modifiers)) * 0.1

	// +0.15 si plusieurs objets (richesse)
	if len(parsed.Objects) > 1 {
		score += 0.15
	}

	// +0.1 si présence d'adverbe (nuance)
	for _, w := range parsed.Words {
		if w.POS == POS_ADV {
			score += 0.1
			break
		}
	}

	// Clamp
	if score > 1 {
		score = 1
	}

	return score
}

// ─────────────────────────────────────────────────────────────────────────
// PARTIE 6: Variation Lexicale Intelligente
// ─────────────────────────────────────────────────────────────────────────

// LexicalVariant représente une variante d'une phrase
type LexicalVariant struct {
	Text          string
	Substitutions int // Nombre de synonymes utilisés
	StyleChange   float64
	Grammar       float64
	FinalScore    float64
}

// GenerateLexicalVariants crée des variantes d'une phrase par synonymes
func GenerateLexicalVariants(original string, maxVariants int) []LexicalVariant {
	variants := []LexicalVariant{}
	tagged := TagMots(original)

	// Identifier les mots significatifs pouvant être variés
	significantWords := []int{}
	for i, w := range tagged {
		if (w.POS == POS_NOUN || w.POS == POS_VERB || w.POS == POS_ADJ) &&
			len(w.Word) > 3 && w.Score > 0.8 {
			significantWords = append(significantWords, i)
		}
	}

	if len(significantWords) == 0 {
		variants = append(variants, LexicalVariant{
			Text:          original,
			Substitutions: 0,
			StyleChange:   0,
			Grammar:       1.0,
			FinalScore:    1.0,
		})
		return variants
	}

	// Générer variantes par substitution de synonymes
	for v := 0; v < maxVariants && v < 3; v++ {
		words := strings.Fields(original)
		numSubs := rand.Intn(min(len(significantWords), 3)) + 1 // 1-3 substitutions

		// Sélectionner mots à remplacer
		toReplace := make(map[int]bool)
		for i := 0; i < numSubs; i++ {
			idx := significantWords[rand.Intn(len(significantWords))]
			toReplace[idx] = true
		}

		// Effectuer substitutions
		substitutions := 0
		for idx := range toReplace {
			word := strings.ToLower(strings.TrimSuffix(strings.TrimSuffix(words[idx], "."), ","))
			if syn, exists := SynonymsDict[word]; exists && len(syn) > 0 {
				synonym := syn[rand.Intn(len(syn))]
				if synonym != word {
					words[idx] = synonym
					substitutions++
				}
			}
		}

		variantText := strings.Join(words, " ")

		// Analyser variante
		variantParsed := AnalyzeDependencies(TagMots(variantText))
		variantParsed.StyleScore = CalculateStyleScore(variantParsed)
		variantParsed.FinalScore = variantParsed.GrammarScore * (1.0 + 0.5*variantParsed.StyleScore)

		variants = append(variants, LexicalVariant{
			Text:          variantText,
			Substitutions: substitutions,
			StyleChange:   variantParsed.StyleScore,
			Grammar:       variantParsed.GrammarScore,
			FinalScore:    variantParsed.FinalScore,
		})
	}

	return variants
}

// ─────────────────────────────────────────────────────────────────────────
// PARTIE 7: Système de Renforcement Progressif
// ─────────────────────────────────────────────────────────────────────────

// ReinforcementAgent apprend par génération et évaluation
type ReinforcementAgent struct {
	SuccessfulPatterns map[string]float64 // Patterns → score récompense
	FailurePatterns    map[string]float64 // Patterns → pénalité
	LearningRate       float64
	EpisodeCount       int
}

// NewReinforcementAgent crée un nouvel agent
func NewReinforcementAgent() *ReinforcementAgent {
	return &ReinforcementAgent{
		SuccessfulPatterns: make(map[string]float64),
		FailurePatterns:    make(map[string]float64),
		LearningRate:       0.1,
		EpisodeCount:       0,
	}
}

// EvaluateSentence évalue une phrase et fournit un score de récompense
func (ra *ReinforcementAgent) EvaluateSentence(sentence string) (float64, string) {
	parsed := AnalyzeDependencies(TagMots(sentence))

	grammarReward := 0.0
	styleReward := 0.0
	naturalReward := 0.0

	// Récompense: Grammaire correcte (+1.0)
	if parsed.MainVerb != nil && parsed.Subject != nil {
		grammarReward = parsed.GrammarScore
	}

	// Récompense: Style vivant (+0.5 max)
	parsed.StyleScore = CalculateStyleScore(parsed)
	styleReward = parsed.StyleScore * 0.5

	// Récompense: Naturel (structure SVO, pas de répétitions)
	if parsed.Subject != nil && parsed.MainVerb != nil {
		if parsed.Subject.Index < parsed.MainVerb.Index {
			naturalReward = 0.5
		}
	}

	totalReward := grammarReward + styleReward + naturalReward

	// Extraction du pattern de la phrase
	pattern := ExtractPattern(parsed)

	// Mise à jour des patterns
	if totalReward > 0.7 {
		ra.SuccessfulPatterns[pattern] += totalReward * ra.LearningRate
	} else {
		ra.FailurePatterns[pattern] += (1.0 - totalReward) * ra.LearningRate
	}

	ra.EpisodeCount++

	feedback := fmt.Sprintf("Grammar:%.2f, Style:%.2f, Natural:%.2f, Total:%.2f",
		grammarReward, styleReward, naturalReward, totalReward)

	return totalReward, feedback
}

// ExtractPattern extrait le pattern structurel d'une phrase
func ExtractPattern(parsed ParsedSentence) string {
	pattern := ""

	if parsed.Subject != nil {
		pattern += "SUBJ:"
	}
	if parsed.MainVerb != nil {
		pattern += "VERB:"
	}
	if len(parsed.Objects) > 0 {
		pattern += "OBJ:"
	}
	if len(parsed.Modifiers) > 0 {
		pattern += "MOD:"
	}

	return pattern
}

// SelectBestVariant choisit la meilleure variante
func (ra *ReinforcementAgent) SelectBestVariant(variants []LexicalVariant) LexicalVariant {
	best := variants[0]
	bestScore := variants[0].FinalScore

	for _, v := range variants[1:] {
		score := v.FinalScore

		// Bonus si pattern appris avec succès
		pattern := "VARIANT"
		if successScore, exists := ra.SuccessfulPatterns[pattern]; exists {
			score += successScore * 0.1
		}

		if score > bestScore {
			bestScore = score
			best = v
		}
	}

	return best
}

// ─────────────────────────────────────────────────────────────────────────
// PARTIE 8: Pipeline Complète de Génération Améliorée
// ─────────────────────────────────────────────────────────────────────────

// EnhancedSentenceGeneration pipeline complète
type EnhancedSentenceGeneration struct {
	Agent    *ReinforcementAgent
	History  []ParsedSentence
	TopScore float64
}

// NewEnhancedGenerator crée un nouveau générateur
func NewEnhancedGenerator() *EnhancedSentenceGeneration {
	return &EnhancedSentenceGeneration{
		Agent:    NewReinforcementAgent(),
		History:  []ParsedSentence{},
		TopScore: 0.0,
	}
}

// GenerateAndOptimize génère et optimise une phrase
func (eg *EnhancedSentenceGeneration) GenerateAndOptimize(baseSentence string, iterations int) string {
	best := baseSentence
	bestScore := 0.0

	// Parse la phrase de base
	baseParsed := AnalyzeDependencies(TagMots(baseSentence))
	baseParsed.StyleScore = CalculateStyleScore(baseParsed)
	baseParsed.FinalScore = baseParsed.GrammarScore * (1.0 + 0.5*baseParsed.StyleScore)

	if baseParsed.FinalScore > bestScore {
		bestScore = baseParsed.FinalScore
		best = baseSentence
	}

	// Générer et tester variantes
	for i := 0; i < iterations; i++ {
		variants := GenerateLexicalVariants(best, 3)
		selectedVariant := eg.Agent.SelectBestVariant(variants)

		// Évaluer
		reward, _ := eg.Agent.EvaluateSentence(selectedVariant.Text)

		if reward > bestScore {
			bestScore = reward
			best = selectedVariant.Text
		}
	}

	// Enregistrer dans historique
	finalParsed := AnalyzeDependencies(TagMots(best))
	eg.History = append(eg.History, finalParsed)
	eg.TopScore = bestScore

	return best
}

// GetReport retourne un rapport sur la phrase générée
func (eg *EnhancedSentenceGeneration) GetReport() string {
	if len(eg.History) == 0 {
		return "No sentences generated yet"
	}

	last := eg.History[len(eg.History)-1]
	return fmt.Sprintf(
		"Grammar:%.2f, Style:%.2f, Final:%.2f, Episodes:%d",
		last.GrammarScore, last.StyleScore, eg.TopScore, eg.Agent.EpisodeCount,
	)
}

// ─────────────────────────────────────────────────────────────────────────
// Utilitaires
// ─────────────────────────────────────────────────────────────────────────

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// BatchProcessParagraph traite un paragraphe entier
func BatchProcessParagraph(paragraph string) []string {
	sentences := strings.Split(paragraph, ". ")
	enhanced := []string{}

	generator := NewEnhancedGenerator()

	for _, sent := range sentences {
		if len(strings.TrimSpace(sent)) < 5 {
			continue
		}

		optimized := generator.GenerateAndOptimize(sent, 5)
		enhanced = append(enhanced, optimized)
	}

	return enhanced
}

// ─────────────────────────────────────────────────────────────────────────
// Export Public Functions
// ─────────────────────────────────────────────────────────────────────────

// AnalyzeSyntax effectue analyse syntaxique complète
func AnalyzeSyntax(sentence string) ParsedSentence {
	tagged := TagMots(sentence)
	parsed := AnalyzeDependencies(tagged)
	parsed.Original = sentence
	parsed.StyleScore = CalculateStyleScore(parsed)
	parsed.FinalScore = parsed.GrammarScore * (1.0 + 0.5*parsed.StyleScore)
	return parsed
}

// GenerateBetterSentence génère une meilleure version d'une phrase
func GenerateBetterSentence(sentence string) (string, float64) {
	generator := NewEnhancedGenerator()
	optimized := generator.GenerateAndOptimize(sentence, 10)
	return optimized, generator.TopScore
}

// AnalyzeParagraphStructure analyse un paragraphe complet
func AnalyzeParagraphStructure(paragraph string) map[string]interface{} {
	sentences := strings.Split(paragraph, ". ")
	result := make(map[string]interface{})

	totalGrammar := 0.0
	totalStyle := 0.0
	sentenceCount := 0

	analyses := []ParsedSentence{}
	for _, sent := range sentences {
		if len(strings.TrimSpace(sent)) < 5 {
			continue
		}

		parsed := AnalyzeSyntax(sent)
		analyses = append(analyses, parsed)
		totalGrammar += parsed.GrammarScore
		totalStyle += parsed.StyleScore
		sentenceCount++
	}

	if sentenceCount > 0 {
		result["avg_grammar"] = totalGrammar / float64(sentenceCount)
		result["avg_style"] = totalStyle / float64(sentenceCount)
		result["sentence_count"] = sentenceCount
		result["overall_score"] = result["avg_grammar"].(float64) *
			(1.0 + 0.5*result["avg_style"].(float64))
	}

	result["sentences"] = analyses

	return result
}
