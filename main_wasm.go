//go:build js && wasm
// +build js,wasm

package main

import (
	"github.com/nosserb/IA-ATOMIQUE-/database"
	"syscall/js"
	"time"
)

func setupWebAssembly() {
	// Expose analyserTexte à JavaScript
	js.Global().Set("analyserTexte", js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) < 1 {
			return map[string]interface{}{
				"erreur": "Pas de texte fourni",
			}
		}

		texte := args[0].String()

		// Appelle la fonction d'analyse existante
		TraiterTexte(texte, "WebAssembly")

		// Retourne les stats
		actifs := 0
		totalEnergie := 0.0
		for _, n := range database.Neurones {
			if n.Valeur > 0 {
				actifs++
				totalEnergie += n.Valeur
			}
		}

		return map[string]interface{}{
			"message":   "Analyse complétée",
			"texte":     texte,
			"neurones":  len(database.Neurones),
			"actifs":    actifs,
			"energie":   totalEnergie,
			"timestamp": time.Now().String(),
		}
	}))

	// Expose getStats pour consulter l'état du réseau neural
	js.Global().Set("getStats", js.FuncOf(func(this js.Value, args []js.Value) any {
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

		return map[string]interface{}{
			"total_neurones":  len(database.Neurones),
			"neurones_actifs": actifs,
			"pourcentage":     float64(actifs) * 100 / float64(len(database.Neurones)),
			"energie_totale":  totalEnergie,
			"energie_max":     maxValeur,
			"timestamp":       time.Now().String(),
		}
	}))
}
