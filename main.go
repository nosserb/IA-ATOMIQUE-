package main

import (
	"IA-ATOMIQUE/database"
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		return
	}
	phrase := strings.Join(os.Args[1:], " ")
	tokens := strings.Fields(phrase)

	stats := make(map[int]int)

	fmt.Print("[SCAN] ")
	for _, t := range tokens {
		word, clean := database.MotProche(t)

		if word.Mot != "" {
			if word.Categorie == 0 {
				fmt.Print("/")
			} else {
				fmt.Print("|")
				stats[word.Categorie]++
				for i := range database.Neurones {
					if database.Neurones[i].CategorieID == word.Categorie {
						database.Neurones[i].Valeur += word.Poids
					}
				}
			}
		} else if temp, ok := database.LexiqueTemp[clean]; ok {
			fmt.Print("?")
			stats[temp.Categorie]++
			for i := range database.Neurones {
				if database.Neurones[i].CategorieID == temp.Categorie {
					database.Neurones[i].Valeur += 1.5
				}
			}
		} else if len(clean) > 2 && !database.StopWords[clean] {
			fmt.Print(".")
		}
	}
	fmt.Println(" [OK]")

	pop := make(map[int]int)
	ener := make(map[int]float64)
	for _, n := range database.Neurones {
		pop[n.CategorieID]++
		ener[n.CategorieID] += n.Valeur
	}

	var premierCat int
	var premierScore, deuxiemeScore float64

	for id, e := range ener {
		if id == 0 || pop[id] <= 0 {
			continue
		}
		score := (e / float64(pop[id])) * float64(stats[id]+1)

		if score > premierScore {
			deuxiemeScore = premierScore
			premierScore = score
			premierCat = id
		} else if score > deuxiemeScore {
			deuxiemeScore = score
		}
	}

	certitude := 10.0
	if deuxiemeScore > 0 {
		certitude = premierScore / deuxiemeScore
	}

	// SEUIL DE MIGRATION : On l'autorise dès que la certitude est claire (>1.5)
	if premierCat != 0 && certitude >= 1.5 {
		for _, t := range tokens {
			_, clean := database.MotProche(t)
			database.Apprendre(clean, premierCat)
		}
	}

	fmt.Printf("\n--- SPECTRE ---\n")
	if premierCat == 0 {
		fmt.Println("ÉTAT : INERTE")
	} else {
		fmt.Printf("1. %-10s : %.2f\n", database.NumeroVersCategorie(premierCat), premierScore)
		fmt.Printf("Certitude : %.2f\n", certitude)
	}
	fmt.Println("---------------")

	database.RegenererNeurones()
}
