package database

import (
	"fmt"
	"math/rand"
	"time"
)

// Neurone représente un neurone avec un ID, une valeur et des connexions
type Neurone struct {
	ID         string
	Valeur     float64
	Connexions map[string]float64
	Fige       bool // true si le neurone est figé
}

// Liste globale des neurones
var Neurones []Neurone

func init() {
	rand.Seed(time.Now().UnixNano())

	// Crée 40 neurones avec valeurs et connexions aléatoires
	for i := 1; i <= 40; i++ {
		n := Neurone{
			ID:         fmt.Sprintf("N%d", i),
			Valeur:     rand.Float64(), // valeur entre 0 et 1
			Connexions: make(map[string]float64),
			Fige:       false,
		}

		// Ajouter 2 à 5 connexions aléatoires vers d'autres neurones
		nbConnexions := rand.Intn(4) + 2
		for j := 0; j < nbConnexions; j++ {
			target := fmt.Sprintf("N%d", rand.Intn(40)+1)
			if target != n.ID {
				n.Connexions[target] = rand.Float64()
			}
		}

		Neurones = append(Neurones, n)
	}
}

// Fonction pour figer un neurone qui a bien marché
func FigerNeurone(id string) {
	for i := range Neurones {
		if Neurones[i].ID == id {
			Neurones[i].Fige = true
			return
		}
	}
}

// Fonction pour regénérer les neurones non figés
func RegenererNeurones() {
	for i := range Neurones {
		if !Neurones[i].Fige {
			Neurones[i].Valeur = rand.Float64()
			Neurones[i].Connexions = make(map[string]float64)
			nbConnexions := rand.Intn(4) + 2
			for j := 0; j < nbConnexions; j++ {
				target := fmt.Sprintf("N%d", rand.Intn(40)+1)
				if target != Neurones[i].ID {
					Neurones[i].Connexions[target] = rand.Float64()
				}
			}
		}
	}
}
