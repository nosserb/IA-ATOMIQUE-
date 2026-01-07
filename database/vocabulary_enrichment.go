package database

import (
	"math"
	"math/rand"
	"strings"
)

// ============================================================================
// PHASE 15: MODULE 3 - ENRICHISSEMENT VOCABULAIRE & SYNONYMES CONTEXTUELS
// ============================================================================
// Dictionnaire synonymes contextuels, paraphrases, variantes stylistiques,
// récompense pour richesse lexicale

// SynonymSet contient un mot et ses synonymes contextualisés
type SynonymSet struct {
	Word      string   // Mot original
	Synonyms  []string // Synonymes possibles
	Contexts  []string // Contextes où utiliser (ex: "formel", "littéraire")
	Frequency float64  // Fréquence du mot (0-1)
	Score     float64  // Score d'enrichissement
}

// VocabularyEnricher gère l'enrichissement lexical
type VocabularyEnricher struct {
	SynonymDictionary map[string]*SynonymSet
	Paraphrases       map[string][]string
	StyleVariants     map[string]map[string][]string // mot -> {style -> variantes}
	LexicalScore      float64
	RichnessBonus     float64
}

// VerbConjugations: variantes de verbes courants
var VerbConjugations = map[string][]string{
	"dire":    {"affirmer", "déclarer", "annoncer", "crier", "murmurer", "chuchoter", "hurler", "clamer", "proclamer", "énnoncer"},
	"faire":   {"réaliser", "accomplir", "exécuter", "effectuer", "perpétrer", "commencer", "entreprendre", "produire", "fabriquer", "construire"},
	"aller":   {"se diriger", "se rendre", "se transporter", "marcher", "progresser", "cheminer", "partir", "se déplacer", "voyager", "circuler"},
	"venir":   {"arriver", "se présenter", "accourir", "apparaître", "survenir", "advenir", "se produire", "procéder", "émaner", "découler"},
	"prendre": {"saisir", "attraper", "capturer", "empoigner", "emparer", "semparer", "s'approprier", "confisquer", "dérober", "enlever"},
	"donner":  {"offrir", "accorder", "conférer", "attribuer", "remettre", "transmettre", "léguer", "céder", "octrroyer", "gratifier"},
	"voir":    {"apercevoir", "distinguer", "remarquer", "observer", "contempler", "constater", "découvrir", "percevoir", "visualiser", "entrevoir"},
	"savoir":  {"connaître", "comprendre", "posséder", "maîtriser", "saisir", "ignorer", "apprendre", "découvrir", "discerner", "éprouver"},
	"pouvoir": {"être capable", "avoir la possibilité", "être à même", "pouvoir", "être en état", "disposer de", "avoir le pouvoir", "avoir l'aptitude", "être autorisé", "pouvoir faire"},
	"vouloir": {"désirer", "souhaiter", "aspirer", "convoiter", "ambitionner", "envier", "rêver", "préférer", "réclamer", "exiger"},
	"devoir":  {"être obligé", "être tenu", "falloir", "être censé", "être contraint", "être forcé", "être destiné", "être redevable", "être supposé", "commettre"},
}

// AdjectiveIntensifiers: variantes d'adjectifs avec intensité
var AdjectiveIntensifiers = map[string][]string{
	"grand":   {"immense", "colossal", "énorme", "gigantesque", "prodigieux", "monumental", "titanesque", "brobdingnagien", "cyclopéen", "gargantuesque"},
	"petit":   {"minuscule", "infime", "microscopique", "diminutif", "lilliputien", "imperceptible", "imperceptible", "étriqué", "rachitique", "chétif"},
	"beau":    {"magnifique", "splendide", "superbe", "ravissant", "merveilleux", "exquis", "sublime", "riche", "opulent", "flamboyant"},
	"laid":    {"hideux", "repoussant", "répugnant", "abominable", "affreux", "monstrueux", "nauséabond", "répulsif", "dégoûtant", "ignoble"},
	"bon":     {"excellent", "formidable", "remarquable", "exceptionnel", "magistral", "splendide", "savoureux", "délectable", "succulent", "divin"},
	"mauvais": {"abominable", "affreux", "atroce", "détestable", "exécrable", "horrible", "ignoble", "infâme", "infect", "maudit"},
	"heureux": {"joyeux", "radieux", "enthousiaste", "exultant", "béat", "euphorique", "ravi", "épanoui", "triomphant", "bienveillant"},
	"triste":  {"mélancolique", "morose", "chagriné", "affligé", "déprimé", "désolé", "abattu", "maussade", "sinistre", "lugubre"},
	"rapide":  {"véloce", "précipité", "fulgurant", "éclaireur", "pressé", "hâtif", "brusque", "prestement", "fébrilement", "vivement"},
	"lent":    {"traînard", "lourd", "pesant", "languissant", "flegmatique", "indolent", "morne", "monotone", "terne", "morose"},
}

// NoveltyAdverbs: adverbes pour enrichir style
var NoveltyAdverbs = []string{
	"magnifiquement", "splendidement", "merveilleusement", "fabuleusement",
	"extraordinairement", "incroyablement", "suprêmement", "parfaitement",
	"remarquablement", "exceptionnellement", "singulièrement", "étrangement",
	"curieusement", "bizarrement", "étonnamment", "dès lors", "néanmoins",
	"cependant", "toutefois", "pourtant", "malgré tout", "somme toute",
	"finalement", "ultimement", "en fin de compte", "par conséquent", "ainsi",
}

// TransitionConnectors: connecteurs pour fluidité
var TransitionConnectors = []string{
	"De plus", "En outre", "De surcroît", "Par ailleurs", "Cependant",
	"Néanmoins", "Toutefois", "Pourtant", "Or", "Ainsi", "Donc",
	"Par conséquent", "De ce fait", "Dès lors", "En conséquence",
	"À cet égard", "Dans le même temps", "Simultanément", "Désormais",
	"Ultérieurement", "Antérieurement", "Entre-temps", "Pendant ce temps",
}

// NewVocabularyEnricher crée un nouvel enrichisseur
func NewVocabularyEnricher() *VocabularyEnricher {
	return &VocabularyEnricher{
		SynonymDictionary: InitializeSynonymDictionary(),
		Paraphrases:       InitializeParaphrases(),
		StyleVariants:     InitializeStyleVariants(),
		LexicalScore:      0.0,
		RichnessBonus:     0.0,
	}
}

// InitializeSynonymDictionary crée le dictionnaire de synonymes
func InitializeSynonymDictionary() map[string]*SynonymSet {
	dict := make(map[string]*SynonymSet)

	// Verbes enrichis
	for verb, synonyms := range VerbConjugations {
		dict[verb] = &SynonymSet{
			Word:      verb,
			Synonyms:  synonyms,
			Contexts:  []string{"littéraire", "formel", "naturel"},
			Frequency: 0.5,
			Score:     0.8,
		}
	}

	// Adjectifs enrichis
	for adj, variants := range AdjectiveIntensifiers {
		dict[adj] = &SynonymSet{
			Word:      adj,
			Synonyms:  variants,
			Contexts:  []string{"intensifié", "nuancé", "classique"},
			Frequency: 0.5,
			Score:     0.75,
		}
	}

	return dict
}

// InitializeParaphrases crée les paraphrases contextualisées
func InitializeParaphrases() map[string][]string {
	return map[string][]string{
		"il y a":    {"existe", "se trouve", "subsiste", "persiste", "demeure"},
		"très":      {"fort", "grandement", "bien", "singulièrement", "remarquablement"},
		"donc":      {"ainsi", "par conséquent", "dès lors", "en conséquence", "pour cette raison"},
		"mais":      {"cependant", "néanmoins", "toutefois", "pourtant", "or"},
		"parce que": {"car", "puisque", "vu que", "du fait que", "étant donné que"},
		"même":      {"également", "aussi", "pareillement", "similairement", "de même"},
		"jamais":    {"nunca", "aucunement", "point", "nullement", "absolument pas"},
		"souvent":   {"fréquemment", "régulièrement", "couramment", "généralement", "maintes fois"},
		"toujours":  {"constamment", "perpétuellement", "incessamment", "sans cesse", "continuellement"},
		"peut-être": {"probablement", "vraisemblablement", "possiblement", "hypothétiquement", "éventuellement"},
	}
}

// InitializeStyleVariants crée les variantes stylistiques par contexte
func InitializeStyleVariants() map[string]map[string][]string {
	return map[string]map[string][]string{
		"aller": {
			"formel":     {"se diriger", "se transporter", "cheminer", "se rendre"},
			"littéraire": {"s'en aller", "prendre son chemin", "voyager", "errer"},
			"naturel":    {"partir", "y aller", "avancer", "se déplacer"},
		},
		"beau": {
			"formel":     {"magnifique", "splendide", "remarquable", "notable"},
			"littéraire": {"ravissant", "merveilleux", "sublime", "exquis"},
			"naturel":    {"joli", "sympa", "mignon", "agréable"},
		},
		"aimer": {
			"formel":     {"apprécier", "préférer", "favoriser", "goûter"},
			"littéraire": {"chérir", "révérer", "idolâtrer", "adorer"},
			"naturel":    {"kimer", "bien aimer", "adorer", "aimer bien"},
		},
		"triste": {
			"formel":     {"mélancolique", "morose", "affligé", "peiné"},
			"littéraire": {"désabusé", "lugubre", "sinistre", "désolé"},
			"naturel":    {"déprimé", "pas bien", "pas gai", "malheureux"},
		},
		"chose": {
			"formel":     {"matière", "objet", "élément", "composant"},
			"littéraire": {"denrée", "créature", "vestige", "essence"},
			"naturel":    {"truc", "machin", "bidule", "chose"},
		},
	}
}

// EnrichSentence enrichit une phrase avec des synonymes contextuels
func (ve *VocabularyEnricher) EnrichSentence(sentence string, style string) string {
	words := strings.Fields(sentence)
	if len(words) == 0 {
		return sentence
	}

	var enriched []string
	for _, word := range words {
		normalized := strings.ToLower(word)

		// Chercher une variante si le mot est dans le dictionnaire
		if synonymSet, exists := ve.SynonymDictionary[normalized]; exists {
			// 30% chance de substitution
			if rand.Float64() < 0.3 && len(synonymSet.Synonyms) > 0 {
				idx := rand.Intn(len(synonymSet.Synonyms))
				enriched = append(enriched, synonymSet.Synonyms[idx])
				continue
			}
		}

		// Chercher variante stylistique
		if styleMap, exists := ve.StyleVariants[normalized]; exists {
			if variants, hasStyle := styleMap[style]; hasStyle && len(variants) > 0 {
				if rand.Float64() < 0.2 {
					idx := rand.Intn(len(variants))
					enriched = append(enriched, variants[idx])
					continue
				}
			}
		}

		enriched = append(enriched, word)
	}

	return strings.Join(enriched, " ")
}

// GenerateVariants génère N variantes d'une phrase
func (ve *VocabularyEnricher) GenerateVariants(sentence string, numVariants int, styles []string) []string {
	if len(styles) == 0 {
		styles = []string{"naturel", "formel", "littéraire"}
	}

	var variants []string
	variants = append(variants, sentence) // Version originale

	for i := 1; i < numVariants; i++ {
		style := styles[i%len(styles)]
		variant := ve.EnrichSentence(sentence, style)
		variants = append(variants, variant)
	}

	return variants
}

// CalculateLexicalRichness calcule la richesse lexicale d'une phrase
func (ve *VocabularyEnricher) CalculateLexicalRichness(sentence string) float64 {
	words := strings.Fields(sentence)
	if len(words) == 0 {
		return 0.0
	}

	// Compte unique words
	uniqueWords := make(map[string]bool)
	complexWords := 0
	modifiers := 0

	for _, word := range words {
		normalized := strings.ToLower(word)
		uniqueWords[normalized] = true

		// Chercher synonymes (indique mot riche)
		if _, exists := ve.SynonymDictionary[normalized]; exists {
			complexWords++
		}

		// Chercher adverbes/adjectifs
		for _, adv := range NoveltyAdverbs {
			if strings.Contains(strings.ToLower(word), adv) {
				modifiers++
				break
			}
		}
	}

	// Ratio = mots uniques + bonus pour mots complexes
	typeTokenRatio := float64(len(uniqueWords)) / float64(len(words))
	complexityBonus := float64(complexWords) * 0.1 / float64(len(words))
	modifierBonus := float64(modifiers) * 0.05 / float64(len(words))

	richness := typeTokenRatio + complexityBonus + modifierBonus
	return math.Min(richness, 1.0)
}

// AddVariationAdverb ajoute un adverbe de variation à une phrase
func (ve *VocabularyEnricher) AddVariationAdverb(sentence string) string {
	words := strings.Fields(sentence)
	if len(words) < 2 {
		return sentence
	}

	// Insérer adverbe aléatoire après le verbe principal (position 1-2 généralement)
	insertPos := 1
	if len(words) > 2 {
		insertPos = 1 + rand.Intn(2)
	}

	adverb := NoveltyAdverbs[rand.Intn(len(NoveltyAdverbs))]
	words = append(words[:insertPos], append([]string{adverb}, words[insertPos:]...)...)

	return strings.Join(words, " ")
}

// AddTransitionConnector ajoute un connecteur en début de phrase
func (ve *VocabularyEnricher) AddTransitionConnector(sentence string) string {
	connector := TransitionConnectors[rand.Intn(len(TransitionConnectors))]
	return connector + " " + sentence
}

// CalculateLexicalScore calcule le score lexical enrichi
func (ve *VocabularyEnricher) CalculateLexicalScore(sentence string) float64 {
	richness := ve.CalculateLexicalRichness(sentence)
	wordCount := float64(len(strings.Fields(sentence)))
	lengthBonus := math.Min(wordCount/20.0, 1.0) // Bonus pour phrases plus longues

	ve.LexicalScore = (richness * 0.7) + (lengthBonus * 0.3)
	ve.RichnessBonus = richness

	return ve.LexicalScore
}

// GetSynonymsForWord retourne tous les synonymes d'un mot
func (ve *VocabularyEnricher) GetSynonymsForWord(word string) []string {
	normalized := strings.ToLower(word)
	if synSet, exists := ve.SynonymDictionary[normalized]; exists {
		return synSet.Synonyms
	}
	return []string{}
}

// GetParaphraseOptions retourne les options de paraphrase
func (ve *VocabularyEnricher) GetParaphraseOptions(phrase string) []string {
	normalized := strings.ToLower(phrase)
	if options, exists := ve.Paraphrases[normalized]; exists {
		return options
	}
	return []string{}
}

// IsRichVocabulary vérifie si le vocabulaire est riche
func (ve *VocabularyEnricher) IsRichVocabulary(sentence string) bool {
	return ve.CalculateLexicalScore(sentence) > 0.6
}

// EnrichParagraph enrichit un paragraphe entier
func (ve *VocabularyEnricher) EnrichParagraph(paragraph string, style string) string {
	sentences := strings.Split(paragraph, ".")
	var enriched []string

	for i, sent := range sentences {
		sent = strings.TrimSpace(sent)
		if sent == "" {
			continue
		}

		// Ajouter connecteur de transition (sauf première phrase)
		if i > 0 && rand.Float64() < 0.4 {
			sent = ve.AddTransitionConnector(sent)
		}

		// Enrichir la phrase
		sent = ve.EnrichSentence(sent, style)

		enriched = append(enriched, sent)
	}

	return strings.Join(enriched, ". ")
}
