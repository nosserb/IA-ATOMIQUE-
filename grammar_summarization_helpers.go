package main

import (
	"strings"

	"IA-ATOMIQUE/database"
)

// validateSummaryAgainstDomain filtre les phrases qui ne sont pas dans le domaine
func validateSummaryAgainstDomain(summary string, domainSpace *database.DomainSpace) string {
	sentences := strings.Split(summary, ".")
	var validSentences []string

	for _, sent := range sentences {
		sent = strings.TrimSpace(sent)
		if sent == "" {
			continue
		}

		// Vérifier si la phrase est valide dans ce domaine
		valid, _ := domainSpace.ValidateSentenceInDomain(sent, 0.75)
		if valid {
			validSentences = append(validSentences, sent)
		}
	}

	return strings.Join(validSentences, ". ") + "."
}
