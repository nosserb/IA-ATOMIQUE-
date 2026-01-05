package main

import (
	"IA-ATOMIQUE/database"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

// SummarizeRequest structure pour la requête de résumé
type SummarizeRequest struct {
	Text string `json:"text"`
}

// SummarizeResponse structure pour la réponse
type SummarizeResponse struct {
	Summary string                 `json:"summary"`
	Stats   map[string]interface{} `json:"stats"`
}

// StartWebServer démarre le serveur web local
func StartWebServer(port string) {
	// Servir les fichiers statiques
	fs := http.FileServer(http.Dir("./web"))
	http.Handle("/", fs)

	// API endpoint pour le résumé
	http.HandleFunc("/api/summarize", handleSummarize)

	fmt.Printf("\n╔════════════════════════════════════════╗\n")
	fmt.Printf("║  Serveur Web Démarré\n")
	fmt.Printf("║  URL: http://localhost:%s\n", port)
	fmt.Printf("╚════════════════════════════════════════╝\n\n")

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Printf("[ERREUR] Impossible de démarrer le serveur: %v\n", err)
	}
}

// handleSummarize traite les requêtes de résumé
func handleSummarize(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	var req SummarizeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Erreur de décodage JSON", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(req.Text) == "" {
		http.Error(w, "Texte vide", http.StatusBadRequest)
		return
	}

	// Analyser le texte avec votre IA
	analyse := database.AnalyserParPhrases(req.Text)
	stats := database.StatistiquesAnalyse(analyse)

	// Générer le résumé
	summary := generateSummary(analyse, req.Text)

	// Préparer les statistiques
	comptage := stats["comptage"].(map[string]int)
	confianceMoyenne := stats["confiance_moyenne"].(float64)
	nbPhrases := stats["nb_phrases"].(int)

	response := SummarizeResponse{
		Summary: summary,
		Stats: map[string]interface{}{
			"phrases":    nbPhrases,
			"confidence": fmt.Sprintf("%.1f", confianceMoyenne*100),
			"categories": len(comptage),
		},
	}

	json.NewEncoder(w).Encode(response)
}

// generateSummary génère un résumé basé sur l'analyse
func generateSummary(analyse database.AnalysePhrases, texteOriginal string) string {
	// Extraire les catégories dominantes
	dominantes := database.CategoriesDominantes(analyse)

	if len(dominantes) == 0 {
		return "Impossible de générer un résumé pour ce texte."
	}

	// Construire un résumé basé sur les catégories principales
	var resumeParts []string

	for i, catData := range dominantes {
		if i >= 3 { // Limiter à 3 catégories principales
			break
		}
		cat := catData["categorie"].(string)
		pct := catData["pourcentage"].(float64)
		resumeParts = append(resumeParts, fmt.Sprintf("• %s (%.1f%%)", cat, pct))
	}

	summary := fmt.Sprintf(
		"📊 Analyse du texte:\n\n%s\n\nLe texte se concentre principalement sur ces catégories.",
		strings.Join(resumeParts, "\n"),
	)

	// Ajouter des phrases clés du texte original
	phrases := strings.Split(texteOriginal, ".")
	if len(phrases) > 3 {
		summary += "\n\n🔑 Phrases clés:\n"
		for i := 0; i < 3 && i < len(phrases); i++ {
			phrase := strings.TrimSpace(phrases[i])
			if len(phrase) > 20 {
				if len(phrase) > 150 {
					phrase = phrase[:150] + "..."
				}
				summary += fmt.Sprintf("\n• %s.", phrase)
			}
		}
	}

	return summary
}

// InitWebInterface initialise et démarre l'interface web
func InitWebInterface() {
	// Créer le répertoire web s'il n'existe pas
	webDir := "./web"
	if _, err := os.Stat(webDir); os.IsNotExist(err) {
		os.MkdirAll(webDir, 0755)
	}

	// Récupérer le port depuis les variables d'environnement (pour Railway, Heroku, etc.)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Port par défaut pour développement local
	}

	// Démarrer le serveur
	go StartWebServer(port)

	// Garder le programme actif
	select {}
}
