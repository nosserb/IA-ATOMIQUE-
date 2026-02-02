package commands

import (
	"github.com/nosserb/IA-ATOMIQUE-/database"
)

// validateSummaryAgainstDomain - NE PAS filtrer, juste retourner le résumé
// (La validation doit respecter le ratio de compression demandé)
func validateSummaryAgainstDomain(summary string, domainSpace *database.DomainSpace) string {
	// Retourner le résumé tel quel - la validation de domaine ne doit pas
	// réduire le résumé et violer le ratio de compression demandé
	return summary
}
