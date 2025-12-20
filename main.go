package main

import (
	"IA-ATOMIQUE/database"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

var stopTrain = false

func VisualiserStats() {
	f, _ := os.Create("dashboard")
	defer f.Close()

	f.WriteString("===========================================\n")
	f.WriteString("   IA-ATOMIQUE v4.0 - NEURAL MONITOR\n")
	f.WriteString("===========================================\n\n")

	// STATISTIQUES NEURONES
	actifs := 0
	totalEnergie := 0.0
	maxValeur := 0.0
	for _, n := range database.Neurones {
		if n.Valeur > 0 {
			actifs++
			totalEnergie += n.Valeur
		}
		if n.Valeur > maxValeur {
			maxValeur = n.Valeur
		}
	}

	f.WriteString(fmt.Sprintf("[RESEAU NEURAL]\n"))
	f.WriteString(fmt.Sprintf("Total neurones: 1000\n"))
	f.WriteString(fmt.Sprintf("Neurones actifs: %d (%.1f%%)\n", actifs, float64(actifs)*100/1000))
	f.WriteString(fmt.Sprintf("Energie totale: %.2f\n", totalEnergie))
	f.WriteString(fmt.Sprintf("Energie max: %.2f\n\n", maxValeur))

	// DISTRIBUTION PAR CATEGORIE
	f.WriteString("[DISTRIBUTION CATEGORIES]\n")
	catStats := make(map[int]map[string]interface{})
	for _, n := range database.Neurones {
		if _, ok := catStats[n.CategorieID]; !ok {
			catStats[n.CategorieID] = map[string]interface{}{"count": 0, "energy": 0.0, "active": 0}
		}
		catStats[n.CategorieID]["count"] = catStats[n.CategorieID]["count"].(int) + 1
		catStats[n.CategorieID]["energy"] = catStats[n.CategorieID]["energy"].(float64) + n.Valeur
		if n.Valeur > 0 {
			catStats[n.CategorieID]["active"] = catStats[n.CategorieID]["active"].(int) + 1
		}
	}

	for id := 1; id <= 50; id++ {
		if stats, ok := catStats[id]; ok {
			count := stats["count"].(int)
			energy := stats["energy"].(float64)
			activeCount := stats["active"].(int)
			f.WriteString(fmt.Sprintf("CAT %2d: ", id))
			// Barre d'energie
			barLen := int(energy / maxValeur * 30)
			if barLen > 30 {
				barLen = 30
			}
			for i := 0; i < barLen; i++ {
				f.WriteString("█")
			}
			for i := barLen; i < 30; i++ {
				f.WriteString("░")
			}
			f.WriteString(fmt.Sprintf(" [%3d/%3d] Energy: %.2f\n", activeCount, count, energy))
		}
	}

	f.WriteString("\n[TOP NEURONES ACTIFS]\n")
	type NeuroneInfo struct {
		id  int
		cat int
		val float64
	}
	var topNeurones []NeuroneInfo
	for _, n := range database.Neurones {
		if n.Valeur > 0 {
			topNeurones = append(topNeurones, NeuroneInfo{n.ID, n.CategorieID, n.Valeur})
		}
	}
	// Sort by value
	for i := 0; i < len(topNeurones) && i < 20; i++ {
		maxIdx := i
		for j := i + 1; j < len(topNeurones); j++ {
			if topNeurones[j].val > topNeurones[maxIdx].val {
				maxIdx = j
			}
		}
		topNeurones[i], topNeurones[maxIdx] = topNeurones[maxIdx], topNeurones[i]
	}

	for i := 0; i < len(topNeurones) && i < 20; i++ {
		n := topNeurones[i]
		f.WriteString(fmt.Sprintf("N%04d [CAT %2d] %.2f ", n.id, n.cat, n.val))
		barLen := int(n.val / maxValeur * 25)
		if barLen > 25 {
			barLen = 25
		}
		for j := 0; j < barLen; j++ {
			f.WriteString("▓")
		}
		f.WriteString("\n")
	}

	f.WriteString("\n[ETAT SYSTEM]\n")
	f.WriteString(fmt.Sprintf("Timestamp: %s\n", time.Now().Format("2006-01-02 15:04:05")))
	f.WriteString(fmt.Sprintf("Status: OPERATIONAL\n"))
}

func AutoIntrospection() {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for !stopTrain {
		<-ticker.C

		// Degrader neurones
		for i := range database.Neurones {
			database.Neurones[i].Valeur *= 0.95
		}

		// Aleatoire activation
		randCat := rand.Intn(50) + 1
		for i := range database.Neurones {
			if database.Neurones[i].CategorieID == randCat {
				database.Neurones[i].Valeur += float64(rand.Intn(3))
			}
		}

		// Sauvegarder stats
		VisualiserStats()
	}
}

func main() {
	go AutoIntrospection()

	if len(os.Args) < 2 {
		fmt.Println("[IA-ATOMIQUE] Prêt pour analyse...")
		for {
			time.Sleep(100 * time.Millisecond)
		}
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

	// SEUIL DE MIGRATION - TRÈS RÉDUIT POUR APPRENDRE AGRESSIVEMENT
	if premierCat != 0 && certitude >= 1.0 {
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
		fmt.Printf("\n> %s\n", database.GenererReponse(premierCat, tokens))
	}
	fmt.Println("---------------")

	VisualiserStats()
	fmt.Println("\n[OK] dashboard généré")

	database.RegenererNeurones()
}
