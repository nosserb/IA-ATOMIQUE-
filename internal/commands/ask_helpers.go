package commands

import (
	"fmt"
	"sort"
	"strings"
	"unicode"
)

// AskCommand traite les questions posées à la base de connaissances
func AskCommand(args []string) {
	if len(args) == 0 {
		fmt.Println("Utilisation: ask <question>")
		fmt.Println("\nExemples:")
		fmt.Println("  ask qui est Joanne Rowling")
		fmt.Println("  ask quand est né Shakespeare")
		fmt.Println("  ask où se trouve Paris")
		return
	}

	question := args[0]
	fmt.Printf("\n[QUESTION]\n%s\n", question)

	// Initialiser les bases de connaissances
	initKnowledgeBase()

	// Extraire les mots-clés de la question
	keywords := extractQuestionKeywords(question)

	// Essayer de répondre avec la base de connaissances
	if answer, ok := buildKnowledgeAnswer(question, keywords); ok {
		fmt.Printf("\n[REPONSE]\n%s\n", answer)
	} else if isFactualQuestion(question) {
		fmt.Println("\n[REPONSE]\nJe n'ai pas de fait fiable pour cette question.")
		fmt.Println("Astuce: utilisez /learn avec un texte pertinent pour enrichir la base.")
	} else {
		fmt.Println("\n[REPONSE]\nCette question ne semble pas être une question factuelle.")
		fmt.Println("Essayez: qui est..., quand..., où..., quoi..., qu'est-ce que...")
	}

	if len(keywords) > 0 {
		maxMots := 3
		if len(keywords) < 3 {
			maxMots = len(keywords)
		}
		fmt.Printf("\nMots clés détectés: %s\n", strings.Join(keywords[:maxMots], ", "))
	}
}

// extractQuestionKeywords extrait les mots-clés importants d'une question
func extractQuestionKeywords(question string) []string {
	words := strings.Fields(strings.ToLower(question))

	// Mots à ignorer
	stopwords := map[string]bool{
		"est": true, "est-ce": true, "qui": true, "quand": true, "où": true,
		"quoi": true, "quels": true, "quelles": true, "quel": true, "qu'est-ce": true,
		"de": true, "du": true, "des": true, "un": true, "une": true,
		"le": true, "la": true, "les": true, "l'": true, "et": true, "ou": true,
		"dans": true, "sur": true, "à": true, "par": true, "pour": true, "avec": true,
		"sans": true, "en": true, "au": true, "aux": true, "donc": true, "car": true,
		"mais": true, "cependant": true, "ainsi": true, "c'est": true, "il": true,
		"elle": true, "on": true, "nous": true, "vous": true, "elles": true, "ils": true,
	}

	var keywords []string
	for _, word := range words {
		word = strings.Trim(word, "?!.,;:")
		if word != "" && !stopwords[word] && len(word) > 2 {
			keywords = append(keywords, word)
		}
	}

	return keywords
}

func buildKnowledgeAnswer(question string, keywords []string) (string, bool) {
	initKnowledgeBase()

	// Essayer d'abord la KB hiérarchisée (nouvelle approche)
	if answer, ok := buildHierarchicalAnswer(question, keywords); ok {
		return answer, true
	}

	// Sinon, utiliser l'ancienne KB
	candidates := extractKnowledgeTerms(question, keywords)
	if len(candidates) == 0 {
		return "", false
	}

	qLower := strings.ToLower(question)
	isWhen := strings.Contains(qLower, "quand") || strings.Contains(qLower, "date") || strings.Contains(qLower, "né") || strings.Contains(qLower, "née") || strings.Contains(qLower, "naissance")
	isWhere := strings.Contains(qLower, "où") || strings.Contains(qLower, "ou ") || strings.Contains(qLower, "lieu")
	isWho := strings.Contains(qLower, "qui")
	isWhat := strings.Contains(qLower, "quoi") || strings.Contains(qLower, "qu'est-ce") || strings.Contains(qLower, "definition")

	for _, term := range candidates {
		info := GlobalKnowledgeBase.GetFactualInfo(term)
		if len(info) == 0 {
			if isWhen {
				if fallback := findDateEventByTerm(term); fallback != "" {
					return fmt.Sprintf("Connaissances sur '%s':\n%s", term, fallback), true
				}
			}
			continue
		}

		lines := []string{}
		header := fmt.Sprintf("Connaissances sur '%s':", term)

		if def, exists := info["definition"]; exists {
			if isWho || isWhat {
				lines = append(lines, fmt.Sprintf("• %s", def))
			} else {
				lines = append(lines, fmt.Sprintf("• Définition: %s", def))
			}
		}

		if dates, exists := info["dates"]; exists {
			dateItems := dates.([]string)
			if len(dateItems) > 0 {
				if isWhen {
					lines = append(lines, fmt.Sprintf("• %s", dateItems[0]))
				} else {
					max := 2
					if len(dateItems) < 2 {
						max = len(dateItems)
					}
					for i := 0; i < max; i++ {
						lines = append(lines, fmt.Sprintf("• %s", dateItems[i]))
					}
				}
			}
		}

		if locs, exists := info["locations"]; exists {
			locItems := locs.([]string)
			if len(locItems) > 0 {
				if isWhere {
					lines = append(lines, fmt.Sprintf("• %s", locItems[0]))
				} else {
					lines = append(lines, fmt.Sprintf("• Lieux: %s", strings.Join(locItems[:minIntAsk(2, len(locItems))], ", ")))
				}
			}
		}

		if causes, exists := info["causes"]; exists {
			causeItems := causes.([]string)
			if len(causeItems) > 0 {
				lines = append(lines, fmt.Sprintf("• Cause: %s", causeItems[0]))
			}
		}

		if related, exists := info["related"]; exists {
			relItems := related.([]string)
			if len(relItems) > 0 {
				lines = append(lines, fmt.Sprintf("• Lié à: %s", strings.Join(relItems[:minIntAsk(3, len(relItems))], ", ")))
			}
		}

		if len(lines) == 0 {
			continue
		}

		return header + "\n" + strings.Join(lines, "\n"), true
	}

	return "", false
}

func findDateEventByTerm(term string) string {
	termLower := normalizeTerm(term)
	if termLower == "" {
		return ""
	}
	GlobalKnowledgeBase.mutex.RLock()
	defer GlobalKnowledgeBase.mutex.RUnlock()

	for year, events := range GlobalKnowledgeBase.DateFacts {
		for _, event := range events {
			if strings.Contains(normalizeTerm(event), termLower) {
				return fmt.Sprintf("• %s: %s", year, event)
			}
		}
	}

	return ""
}

func normalizeTerm(s string) string {
	var b strings.Builder
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(unicode.ToLower(r))
		}
	}
	return b.String()
}

func isFactualQuestion(question string) bool {
	qLower := strings.ToLower(question)
	return strings.Contains(qLower, "quand") ||
		strings.Contains(qLower, "qui") ||
		strings.Contains(qLower, "où") ||
		strings.Contains(qLower, "ou ") ||
		strings.Contains(qLower, "quoi") ||
		strings.Contains(qLower, "qu'est-ce") ||
		strings.Contains(qLower, "combien") ||
		strings.Contains(qLower, "date") ||
		strings.Contains(qLower, "né") ||
		strings.Contains(qLower, "née") ||
		strings.Contains(qLower, "naissance") ||
		strings.Contains(qLower, "definition")
}

func extractKnowledgeTerms(question string, keywords []string) []string {
	seen := make(map[string]bool)
	add := func(term string) {
		term = strings.TrimSpace(strings.ToLower(term))
		if term == "" {
			return
		}
		if !seen[term] {
			seen[term] = true
		}
	}

	// Keywords from analysis
	for _, k := range keywords {
		k = strings.Trim(k, "\"'.,;:!?()[]{}")
		if len(k) >= 2 {
			add(k)
		}
	}

	// Tokens and uppercase sequences
	tokens := strings.Fields(question)
	var current []string
	flush := func() {
		if len(current) > 0 {
			add(strings.Join(current, " "))
			current = nil
		}
	}

	for _, raw := range tokens {
		clean := strings.Trim(raw, "\"'.,;:!?()[]{}")
		if clean == "" {
			flush()
			continue
		}
		isUpperToken := hasUpper(clean) || strings.Contains(clean, ".")
		if isUpperToken {
			current = append(current, clean)
		} else {
			flush()
			if len(clean) >= 4 {
				add(clean)
			}
		}
	}
	flush()

	terms := make([]string, 0, len(seen))
	for t := range seen {
		terms = append(terms, t)
	}

	// Sort by length (longest first) to prioritize multi-word terms like "joanne rowling"
	sort.Slice(terms, func(i, j int) bool {
		return len(terms[i]) > len(terms[j])
	})

	return terms
}

func hasUpper(token string) bool {
	for _, r := range token {
		if unicode.IsUpper(r) {
			return true
		}
	}
	return false
}

func minIntAsk(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// buildHierarchicalAnswer construit une réponse NATURELLE à partir de la KB hiérarchisée
func buildHierarchicalAnswer(question string, keywords []string) (string, bool) {
	if GlobalHierarchicalKB == nil {
		return "", false
	}

	// Extraire les candidats
	candidates := extractKnowledgeTerms(question, keywords)
	if len(candidates) == 0 {
		return "", false
	}

	qLower := strings.ToLower(question)
	isWhen := strings.Contains(qLower, "quand") || strings.Contains(qLower, "né") || strings.Contains(qLower, "née")
	isWhere := strings.Contains(qLower, "où") || strings.Contains(qLower, "lieu")
	isWho := strings.Contains(qLower, "qui")

	// Chercher l'entité dans la KB hiérarchisée
	for _, candidate := range candidates {
		entity := GlobalHierarchicalKB.GetEntity(candidate)
		if entity == nil {
			continue
		}

		// Générer une réponse avec des phrases naturelles
		answer := formatEntityAsNaturalText(entity, isWho, isWhen, isWhere)
		if answer != "" {
			return answer, true
		}
	}

	return "", false
}

// formatEntityAsNaturalText génère des phrases naturelles à partir d'une entité
func formatEntityAsNaturalText(entity *Entity, isWho, isWhen, isWhere bool) string {
	var sentences []string

	// Capitaliser le nom pour la première phrase
	entityName := strings.ToUpper(entity.ID[:1]) + entity.ID[1:]

	// QUAND - priorité à la date de naissance
	if isWhen {
		if naissance, ok := entity.Properties["naissance"]; ok {
			sentences = append(sentences, fmt.Sprintf("%s est né(e) en %v.", entityName, naissance))
		}
		if mort, ok := entity.Properties["mort"]; ok {
			sentences = append(sentences, fmt.Sprintf("Il/Elle est décédé(e) en %v.", mort))
		}
	}

	// QUI - priorité à la fonction/rôle
	if isWho {
		if fonction, ok := entity.Properties["fonction"]; ok {
			if fonctions, isList := fonction.([]string); isList {
				for i, f := range fonctions {
					if i == 0 {
						sentences = append(sentences, fmt.Sprintf("%s est %s.", entityName, f))
					} else {
						sentences = append(sentences, fmt.Sprintf("Il/Elle est aussi %s.", f))
					}
				}
			}
		}
		if naissance, ok := entity.Properties["naissance"]; ok && len(sentences) == 0 {
			sentences = append(sentences, fmt.Sprintf("%s est né(e) en %v.", entityName, naissance))
		}
	}

	// OÙ - priorité à la localisation
	if isWhere {
		if lieu, ok := entity.Properties["lieu"]; ok {
			sentences = append(sentences, fmt.Sprintf("%s est associé(e) à %v.", entityName, lieu))
		}
	}

	// Si pas de question spécifique, générer un résumé général
	if len(sentences) == 0 {
		// Ajouter la fonction si disponible
		if fonction, ok := entity.Properties["fonction"]; ok {
			if fonctions, isList := fonction.([]string); isList && len(fonctions) > 0 {
				sentences = append(sentences, fmt.Sprintf("%s est %s.", entityName, fonctions[0]))
			}
		}

		// Ajouter la naissance si disponible
		if naissance, ok := entity.Properties["naissance"]; ok {
			sentences = append(sentences, fmt.Sprintf("Il/Elle est né(e) en %v.", naissance))
		}

		// Ajouter autres propriétés intéressantes
		if localisation, ok := entity.Properties["lieu"]; ok {
			sentences = append(sentences, fmt.Sprintf("On le/la retrouve à %v.", localisation))
		}
	}

	// Ajouter infos complémentaires si la réponse est courte
	if len(sentences) <= 1 {
		// Ajouter les œuvres si c'est pertinent
		if len(entity.SubEntities) > 0 {
			count := len(entity.SubEntities)
			sentences = append(sentences, fmt.Sprintf("Il/Elle est associé(e) à %d œuvre(s).", count))
		}

		// Ajouter les relations
		for relType, relatedIDs := range entity.Relations {
			if len(relatedIDs) > 0 {
				count := len(relatedIDs)
				sentences = append(sentences, fmt.Sprintf("Il/Elle a %d lien(s) de type '%s'.", count, relType))
			}
		}
	}

	// Joindre les phrases avec un espace
	result := strings.Join(sentences, " ")

	// Si rien n'a été généré, retourner une info basique
	if result == "" {
		return fmt.Sprintf("%s existe dans la base de connaissances.", entityName)
	}

	return result
}
