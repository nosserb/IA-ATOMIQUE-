package main

import (
	"fmt"
	"math"
	"os"
	"strings"

	"IA-ATOMIQUE/database"
)

// MotVersNumero retourne le numéro/catégorie d'un mot
func MotVersNumero(mot string) int {
	mot = strings.ToLower(mot)
	for _, w := range database.Words {
		if w.Mot == mot {
			return w.Numero
		}
	}
	return 0
}

// InfluenceNeurones met à jour les neurones selon les mots de la phrase
func InfluenceNeurones(tokens []string) {
	for _, t := range tokens {
		num := MotVersNumero(t)
		if num == 0 {
			continue
		}

		// Choisir un neurone non figé
		for i := range database.Neurones {
			if !database.Neurones[i].Fige {
				database.Neurones[i].Valeur += float64(num)
				break
			}
		}
	}
}

// MoyenneNeurones calcule la moyenne des neurones actifs
func MoyenneNeurones() float64 {
	total := 0.0
	count := 0
	for _, n := range database.Neurones {
		if n.Valeur > 0 {
			total += n.Valeur
			count++
		}
	}
	if count == 0 {
		return 0
	}
	return total / float64(count)
}

// Map du numéro de catégorie vers le nom
func NumeroVersCategorie(num int) string {
	switch num {
	case 1:
		return "Ingénierie"
	case 2:
		return "Science"
	case 3:
		return "Nature"
	case 4:
		return "Art"
	case 5:
		return "Programmation"
	default:
		return "Inconnu"
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go \"votre phrase ici\"")
		return
	}

	phrase := strings.Join(os.Args[1:], " ")
	tokens := strings.Fields(phrase)

	InfluenceNeurones(tokens)

	moy := MoyenneNeurones()
	catNum := int(math.Round(moy))
	categorie := NumeroVersCategorie(catNum)

	fmt.Printf("Thème principal détecté : %s\n", categorie)

	// Figer neurones performants
	for _, n := range database.Neurones {
		if n.Valeur >= 0.7 && !n.Fige {
			database.FigerNeurone(n.ID)
		}
	}

	database.RegenererNeurones()
}
