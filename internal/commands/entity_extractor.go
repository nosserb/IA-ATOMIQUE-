package commands

import (
	"regexp"
	"strings"
)

// ============================================================================
// EXTRACTEUR INTELLIGENT D'ENTITÉS ET RELATIONS
// ============================================================================

// PatternExtractor détecte et extrait entités avec leurs propriétés
type PatternExtractor struct {
	// Patterns pour détecter types d'entités
	PersonPattern   *regexp.Regexp
	WorkPattern     *regexp.Regexp
	DatePattern     *regexp.Regexp
	RolePattern     *regexp.Regexp
	LocationPattern *regexp.Regexp
}

// NewPatternExtractor crée un nouveau extracteur
func NewPatternExtractor() *PatternExtractor {
	return &PatternExtractor{
		// Patterns pour détecter des personnes (Nom Prénom)
		PersonPattern: regexp.MustCompile(`([A-Z][a-zéèêàùçœ]+(?:\s+[A-Z][a-zéèêàùçœ]+)*)`),

		// Patterns pour détecter des œuvres (guillemets ou ALL CAPS ou "de X")
		WorkPattern: regexp.MustCompile(`«([^»]+)»|"([^"]+)"|([A-Z][A-Z\s]+)|\b(de|du|des|l'|le|la)(\s+[A-Z][a-zéèêàùçœ\s]+)`),

		// Pattern pour dates
		DatePattern: regexp.MustCompile(`(\d{1,2})\s+(janvier|février|mars|avril|mai|juin|juillet|août|septembre|octobre|novembre|décembre)\s+(\d{4})|(\d{4})`),

		// Pattern pour rôles/fonctions
		RolePattern: regexp.MustCompile(`(?:est|sont)\s+(?:un|une|des)?\s+([a-zéèêàùçœ\s]+?)(?:\s+et\s+|,|\.|\s+de|\s+qui)`),

		// Pattern pour lieux
		LocationPattern: regexp.MustCompile(`\b(France|Royaume-Uni|Angleterre|États-Unis|Paris|Londres|Oxford|Cambridge|Porto|Édimbourg|Bristol|Yate|Tutshill|Chepstow)\b`),
	}
}

// ExtractEntitiesFromText analyse un texte et extrait entités et propriétés
func (pe *PatternExtractor) ExtractEntitiesFromText(text string) map[string]map[string]interface{} {
	entities := make(map[string]map[string]interface{})

	// Extraire personnes
	pe.extractPeople(text, entities)

	// Extraire dates et les lier aux personnes
	pe.extractDates(text, entities)

	// Extraire rôles/fonctions
	pe.extractRoles(text, entities)

	// Extraire lieux
	pe.extractLocations(text, entities)

	// Extraire œuvres
	pe.extractWorks(text, entities)

	return entities
}

// extractPeople détecte les personnes et leurs variantes
func (pe *PatternExtractor) extractPeople(text string, entities map[string]map[string]interface{}) {
	// Pattern pour noms propres avec variantes (J.K. Rowling, Joanne Rowling, etc.)
	nameVariantPattern := regexp.MustCompile(
		`([A-Z][a-zéèêàùçœ]+)(?:\s+[A-Z]\.)?(?:\s+[A-Z][a-zéèêàùçœ]+)*(?:\s*,\s*plus connue sous[^,]*([^,]+))?`,
	)

	matches := nameVariantPattern.FindAllStringSubmatch(text, -1)
	for _, match := range matches {
		if len(match) >= 2 {
			mainName := NormalizeID(match[1])
			if mainName != "" {
				if _, exists := entities[mainName]; !exists {
					entities[mainName] = map[string]interface{}{
						"type":    "personne",
						"aliases": []string{},
					}
				}

				// Ajouter alias si présent
				if len(match) > 2 && match[2] != "" {
					aliases := entities[mainName]["aliases"].([]string)
					aliasText := strings.TrimSpace(match[2])
					if aliasText != "" && !containsString(aliases, aliasText) {
						entities[mainName]["aliases"] = append(aliases, aliasText)
					}
				}
			}
		}
	}
}

// extractDates détecte les dates et les associe INTELLIGEMMENT aux personnes
func (pe *PatternExtractor) extractDates(text string, entities map[string]map[string]interface{}) {
	// Pattern 1: "Personne est née/mort le DATE" - LE PLUS FIABLE
	naissancePattern := regexp.MustCompile(
		`([A-Z][a-zéèêàùçœ]+(?:\s+[A-Z][a-zéèêàùçœ]+)*)\s+(?:est\s+)?(?:né|née|naît|naquit)\s+(?:le\s+)?(\d{1,2})\s+(janvier|février|mars|avril|mai|juin|juillet|août|septembre|octobre|novembre|décembre)\s+(\d{4})`,
	)

	// Pattern 1b: "Personne est née en ANNÉE" - SIMPLER VERSION
	naissanceAnneePattern := regexp.MustCompile(
		`([A-Z][a-zéèêàùçœ]+(?:\s+[A-Z][a-zéèêàùçœ]+)*)\s+(?:est\s+)?(?:né|née|naît|naquit)\s+en\s+(\d{4})`,
	)

	// Pattern 2: "Personne (DATE - DATE)" ou "Personne (DATE)" - FORMAT COMPACT
	parenthesesPattern := regexp.MustCompile(
		`([A-Z][a-zéèêàùçœ]+(?:\s+[A-Z][a-zéèêàùçœ]+)*)\s+\((\d{1,2})\s+(janvier|février|mars|avril|mai|juin|juillet|août|septembre|octobre|novembre|décembre)?\s*(\d{4})(?:\s*-\s*[^)]+)?\)`,
	)

	// Pattern 3: "Personne était né le DATE" - VARIANTES PASSÉ
	variationPattern := regexp.MustCompile(
		`([A-Z][a-zéèêàùçœ]+(?:\s+[A-Z][a-zéèêàùçœ]+)*)\s+(?:était|a été)\s+(?:né|née)\s+(?:le\s+)?(\d{1,2})\s+(janvier|février|mars|avril|mai|juin|juillet|août|septembre|octobre|novembre|décembre)\s+(\d{4})`,
	)

	monthMap := map[string]string{
		"janvier": "01", "février": "02", "mars": "03", "avril": "04",
		"mai": "05", "juin": "06", "juillet": "07", "août": "08",
		"septembre": "09", "octobre": "10", "novembre": "11", "décembre": "12",
	}

	// Chercher pattern 1: "X est né le DATE" - TRÈS FIABLE
	matches := naissancePattern.FindAllStringSubmatch(text, -1)
	for _, match := range matches {
		if len(match) >= 5 {
			personID := NormalizeID(match[1])
			if personID != "" {
				if _, exists := entities[personID]; !exists {
					entities[personID] = map[string]interface{}{
						"type": "personne",
					}
				}
				// Format: jour-mois-année ou juste année
				if monthIdx, hasMonth := monthMap[strings.ToLower(match[3])]; hasMonth {
					dateStr := match[4] + "-" + monthIdx + "-" + match[2]
					entities[personID]["naissance"] = dateStr
				}
			}
		}
	}

	// Chercher pattern 1b: "X est née en 1967" - SIMPLE ANNÉE
	matches = naissanceAnneePattern.FindAllStringSubmatch(text, -1)
	for _, match := range matches {
		if len(match) >= 3 {
			personID := NormalizeID(match[1])
			if personID != "" {
				if _, exists := entities[personID]; !exists {
					entities[personID] = map[string]interface{}{
						"type": "personne",
					}
				}
				// Ajouter uniquement si pas déjà défini
				if _, hasNaissance := entities[personID]["naissance"]; !hasNaissance {
					entities[personID]["naissance"] = match[2]
				}
			}
		}
	}

	// Chercher pattern 2: "X (31 juillet 1965)" - COMPACT
	matches = parenthesesPattern.FindAllStringSubmatch(text, -1)
	for _, match := range matches {
		if len(match) >= 5 {
			personID := NormalizeID(match[1])
			if personID != "" {
				if _, exists := entities[personID]; !exists {
					entities[personID] = map[string]interface{}{
						"type": "personne",
					}
				}
				// Format: jour mois année
				year := match[4]
				if month, hasMonth := monthMap[strings.ToLower(match[3])]; match[3] != "" && hasMonth {
					dateStr := year + "-" + month + "-" + match[2]
					entities[personID]["naissance"] = dateStr
				} else if match[3] == "" {
					// Juste l'année: "X (1965)"
					entities[personID]["naissance"] = year
				}
			}
		}
	}

	// Chercher pattern 3: "X était née le DATE"
	matches = variationPattern.FindAllStringSubmatch(text, -1)
	for _, match := range matches {
		if len(match) >= 5 {
			personID := NormalizeID(match[1])
			if personID != "" {
				if _, exists := entities[personID]; !exists {
					entities[personID] = map[string]interface{}{
						"type": "personne",
					}
				}
				if monthIdx, hasMonth := monthMap[strings.ToLower(match[3])]; hasMonth {
					dateStr := match[4] + "-" + monthIdx + "-" + match[2]
					entities[personID]["naissance"] = dateStr
				}
			}
		}
	}
}

// extractRoles détecte les fonctions/rôles avec plusieurs patterns
func (pe *PatternExtractor) extractRoles(text string, entities map[string]map[string]interface{}) {
	// Pattern 1: "[Personne] est [un/une] [role]"
	roleDetailedPattern := regexp.MustCompile(
		`([A-Z][a-zéèêàùçœ]+(?:\s+[A-Z][a-zéèêàùçœ]+)*)\s+est\s+(?:un|une|l')\s+([a-zéèêàùçœ]+(?:\s+et\s+[a-zéèêàùçœ]+)*)\b`,
	)

	// Pattern 2: "[Personne] (commat) [role] de"
	roleSimplePattern := regexp.MustCompile(
		`([A-Z][a-zéèêàùçœ]+(?:\s+[A-Z][a-zéèêàùçœ]+)*),\s+([a-zéèêàùçœ]+)\s+(?:de|et)`,
	)

	// Pattern 3: "auteur|écrivain|romancier|artiste" near person name
	professionKeywords := []string{"auteur", "autrice", "écrivain", "écrivaine", "romancier", "romancière",
		"poète", "poétesse", "compositeur", "compositrice", "artiste", "musicien", "musicienne",
		"peintre", "sculpteur", "sculptrice", "réalisateur", "réalisatrice", "acteur", "actrice",
		"chercheur", "chercheuse", "professeur", "professeure", "médecin", "scientifique", "biologiste"}

	// Chercher pattern 1: "X est un/une [role]"
	matches := roleDetailedPattern.FindAllStringSubmatch(text, -1)
	for _, match := range matches {
		if len(match) >= 3 {
			personID := NormalizeID(match[1])
			role := strings.TrimSpace(match[2])

			if personID != "" && role != "" && len(role) > 2 {
				if _, exists := entities[personID]; !exists {
					entities[personID] = map[string]interface{}{
						"type": "personne",
					}
				}

				// Ajouter rôle(s)
				if existing, hasRole := entities[personID]["fonction"]; hasRole {
					roles := existing.([]string)
					if !containsString(roles, role) {
						entities[personID]["fonction"] = append(roles, role)
					}
				} else {
					entities[personID]["fonction"] = []string{role}
				}
			}
		}
	}

	// Chercher pattern 2: "X, [role] de..."
	matches = roleSimplePattern.FindAllStringSubmatch(text, -1)
	for _, match := range matches {
		if len(match) >= 3 {
			personID := NormalizeID(match[1])
			role := strings.TrimSpace(match[2])

			if personID != "" && role != "" && len(role) > 2 {
				if _, exists := entities[personID]; !exists {
					entities[personID] = map[string]interface{}{
						"type": "personne",
					}
				}

				if existing, hasRole := entities[personID]["fonction"]; hasRole {
					roles := existing.([]string)
					if !containsString(roles, role) {
						entities[personID]["fonction"] = append(roles, role)
					}
				} else {
					entities[personID]["fonction"] = []string{role}
				}
			}
		}
	}

	// Chercher les keywords de profession dans le texte près des noms
	for personID := range entities {
		if entity, exists := entities[personID]; exists {
			// Chercher si cette personne a déjà une fonction
			if _, hasFunc := entity["fonction"]; !hasFunc {
				// Chercher les keywords de profession dans le texte (case-insensitive)
				textLower := strings.ToLower(text)
				personNameLower := strings.ToLower(personID)

				for _, keyword := range professionKeywords {
					// Chercher le keyword dans une fenêtre proche du nom
					personIdx := strings.Index(textLower, personNameLower)
					if personIdx >= 0 {
						// Regarder dans une fenêtre de 100 caractères autour du nom
						start := personIdx - 100
						if start < 0 {
							start = 0
						}
						end := personIdx + len(personNameLower) + 100
						if end > len(textLower) {
							end = len(textLower)
						}

						windowText := textLower[start:end]
						if strings.Contains(windowText, keyword) {
							entities[personID]["fonction"] = []string{keyword}
							break
						}
					}
				}
			}
		}
	}
}

// extractLocations détecte les lieux
func (pe *PatternExtractor) extractLocations(text string, entities map[string]map[string]interface{}) {
	locMatches := pe.LocationPattern.FindAllString(text, -1)
	for _, loc := range locMatches {
		locID := NormalizeID(loc)
		if locID != "" {
			if _, exists := entities[locID]; !exists {
				entities[locID] = map[string]interface{}{
					"type": "lieu",
				}
			}
		}
	}
}

// extractWorks détecte les œuvres (livres, séries, etc.)
func (pe *PatternExtractor) extractWorks(text string, entities map[string]map[string]interface{}) {
	// Pattern: "le/la/les [titre]" suivi de "est/sont"
	workPattern := regexp.MustCompile(
		`(?:le|la|les|l')\s+([A-Z][a-zéèêàùçœ\s]+?)(?:\s+est\s+|\s+sont\s+|,|\.)|«([^»]+)»|"([^"]+)"`,
	)

	matches := workPattern.FindAllStringSubmatch(text, -1)
	for _, match := range matches {
		var workName string
		if match[1] != "" {
			workName = strings.TrimSpace(match[1])
		} else if match[2] != "" {
			workName = strings.TrimSpace(match[2])
		} else if match[3] != "" {
			workName = strings.TrimSpace(match[3])
		}

		if workName != "" {
			workID := NormalizeID(workName)
			if workID != "" && len(workID) >= 4 { // Ignorer trop court
				if _, exists := entities[workID]; !exists {
					entities[workID] = map[string]interface{}{
						"type": "œuvre",
					}
				}
			}
		}
	}
}

// BuildEntityAndAddToKB construit une Entity et l'ajoute à la KB hiérarchisée
func BuildEntityAndAddToKB(entityID string, data map[string]interface{}) {
	if GlobalHierarchicalKB == nil {
		return
	}

	// Convertir map[string]interface{} en Properties valides
	properties := make(map[string]interface{})

	// Copier les propriétés
	for key, value := range data {
		if key != "type" && key != "aliases" {
			properties[key] = value
		}
	}

	// Déterminer le type
	entityType := "concept"
	if typeVal, hasType := data["type"]; hasType {
		entityType = typeVal.(string)
	}

	GlobalHierarchicalKB.AddOrMergeEntity(entityID, entityType, properties)
}
