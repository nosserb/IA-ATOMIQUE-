package commands

import (
	"IA-ATOMIQUE/database"
	"fmt"
	"os"
	"time"
)

// ExtrairePhrasesClésCommand extrait et affiche les phrases clés d'un texte.
func ExtrairePhrasesClésCommand(args []string) {
	if len(args) == 0 {
		fmt.Println("[ERREUR] Spécifier un fichier")
		return
	}

	contenu, err := os.ReadFile(args[0])
	if err != nil {
		fmt.Printf("[ERREUR] Impossible de lire %s: %v\n", args[0], err)
		return
	}

	texte := string(contenu)

	ratio := 0.3
	if len(args) > 1 {
		fmt.Sscanf(args[1], "%f", &ratio)
	}

	fmt.Printf("\n╔════════════════════════════════════════════════════════╗\n")
	fmt.Printf("║   EXTRACTION DE PHRASES CLÉS - ÉNERGIE ATOMIQUE       ║\n")
	fmt.Printf("╚════════════════════════════════════════════════════════╝\n\n")

	debut := time.Now()

	phrasesClés := database.ExtrairePhrasesClés(texte, ratio)

	fmt.Printf("[PHRASES CLÉS] (ratio %.0f%%)\n", ratio*100)
	for i, phrase := range phrasesClés {
		fmt.Printf("%2d. %s\n", i+1, phrase)
	}

	fmt.Printf("\n[PERFORMANCE]\n")
	fmt.Printf("  • Temps: %v\n", time.Since(debut))
	fmt.Printf("  • Phrases extraites: %d\n\n", len(phrasesClés))
	fmt.Println("✓ Extraction terminée")
}
