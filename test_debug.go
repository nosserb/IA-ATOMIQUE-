package main

import (
	"fmt"
	"strings"
)

func main() {
	texte := `La photosynthèse est le processus biologique fondamental.
Par lequel les plantes vertes et certains microorganismes convertissent l'énergie lumineuse.
En énergie chimique Ce processus implique la fixation du dioxyde de carbone.
Et la libération d'oxyde comme sous-produit.`

	phrases := strings.Split(texte, ".")
	fmt.Println("Phrases:", len(phrases))
	for i, p := range phrases {
		fmt.Printf("%d: [%s]\n", i, strings.TrimSpace(p))
	}
}
