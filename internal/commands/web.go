package commands

import (
	"IA-ATOMIQUE/database"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strings"
)

// SummarizeRequest structure pour la requête de résumé
type SummarizeRequest struct {
	Text string `json:"text"`
}

// SummarizeResponse structure pour la réponse
type SummarizeResponse struct {
	KeyPhrases []string               `json:"keyPhrases"`
	KeyIdeas   []string               `json:"keyIdeas"`
	Summary    string                 `json:"summary"`
	Stats      map[string]interface{} `json:"stats"`
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

	// Extraire les phrases clés et idées
	keyPhrases, keyIdeas, summary := generateDetailedSummary(analyse)

	// Préparer les statistiques
	comptage := stats["comptage"].(map[string]int)
	confianceMoyenne := stats["confiance_moyenne"].(float64)
	nbPhrases := stats["nb_phrases"].(int)

	response := SummarizeResponse{
		KeyPhrases: keyPhrases,
		KeyIdeas:   keyIdeas,
		Summary:    summary,
		Stats: map[string]interface{}{
			"phrases":    nbPhrases,
			"confidence": fmt.Sprintf("%.1f", confianceMoyenne*100),
			"categories": len(comptage),
		},
	}

	json.NewEncoder(w).Encode(response)
}

// generateDetailedSummary génère les phrases clés résumées, idées clés et résumé
func generateDetailedSummary(analyse database.AnalysePhrases) ([]string, []string, string) {
	type phraseImportante struct {
		texte     string
		score     float64
		confiance float64
	}

	var phrasesTriees []phraseImportante

	// Collecter toutes les phrases avec un score décent
	for _, phrase := range analyse.Phrases {
		if phrase.Score > 0 || phrase.Confiance > 0.3 {
			phrasesTriees = append(phrasesTriees, phraseImportante{
				texte:     phrase.Texte,
				score:     phrase.Score,
				confiance: phrase.Confiance,
			})
		}
	}

	// Trier par score décroissant
	sort.Slice(phrasesTriees, func(i, j int) bool {
		if phrasesTriees[i].score == phrasesTriees[j].score {
			return phrasesTriees[i].confiance > phrasesTriees[j].confiance
		}
		return phrasesTriees[i].score > phrasesTriees[j].score
	})

	// 1. PHRASES CLÉS = résumer les meilleures phrases (max 100 chars)
	var keyPhrases []string
	for i, p := range phrasesTriees {
		if i >= 8 { // Max 8 phrases clés
			break
		}
		texte := strings.TrimSpace(p.texte)

		// Résumer si trop long
		if len(texte) > 120 {
			// Garder jusqu'au dernier point/virgule avant 120 chars
			if idx := strings.LastIndex(texte[:120], "."); idx > 0 {
				texte = texte[:idx+1]
			} else if idx := strings.LastIndex(texte[:120], ","); idx > 0 {
				texte = texte[:idx]
			} else {
				texte = texte[:120] + "..."
			}
		}

		keyPhrases = append(keyPhrases, strings.TrimSpace(texte))
	}

	// 2. IDÉES CLÉS = mots-clés extraits des top phrases + catégories
	var keyIdeas []string

	// Extraire les mots-clés des top phrases
	ideasMap := make(map[string]bool)
	for i := 0; i < 3 && i < len(phrasesTriees); i++ {
		phrase := phrasesTriees[i]
		words := strings.Fields(phrase.texte)
		for _, word := range words {
			// Garder les mots > 4 caractères (sauf articles)
			if len(word) > 4 && !isStopWord(word) {
				ideasMap[word] = true
			}
		}
	}

	// Ajouter les catégories dominantes
	dominantes := database.CategoriesDominantes(analyse)
	for i, catData := range dominantes {
		if i >= 3 { // Top 3 catégories
			break
		}
		cat := catData["categorie"].(string)
		pct := catData["pourcentage"].(float64)
		ideasMap[fmt.Sprintf("%s (%.0f%%)", cat, pct)] = true
	}

	// Convertir en slice
	for idea := range ideasMap {
		keyIdeas = append(keyIdeas, idea)
		if len(keyIdeas) >= 6 {
			break
		}
	}

	// 3. RÉSUMÉ = les 2-3 meilleures phrases
	var summaryText string
	if len(phrasesTriees) >= 2 {
		summaryText = strings.TrimSpace(phrasesTriees[0].texte) + " " + strings.TrimSpace(phrasesTriees[1].texte)
	} else if len(phrasesTriees) > 0 {
		summaryText = strings.TrimSpace(phrasesTriees[0].texte)
	} else {
		summaryText = "Impossible de générer un résumé pour ce texte."
	}

	// Ajouter ponctuation si manquante
	if len(summaryText) > 0 && !strings.HasSuffix(summaryText, ".") {
		summaryText += "."
	}

	return keyPhrases, keyIdeas, summaryText
}

// isStopWord vérifie si un mot est un stop word français
func isStopWord(word string) bool {
	stopWords := map[string]bool{
		"le": true, "la": true, "les": true, "de": true, "du": true, "des": true,
		"un": true, "une": true, "et": true, "ou": true, "mais": true, "est": true,
		"ont": true, "être": true, "avoir": true, "par": true, "pour": true, "que": true,
		"qui": true, "ce": true, "cet": true, "dans": true, "sur": true, "avec": true,
		"sans": true, "sous": true, "entre": true, "vers": true, "pendant": true,
	}
	wordLower := strings.ToLower(strings.TrimSpace(word))
	wordLower = strings.Trim(wordLower, ".,;:!?()[]{}\"'")
	return stopWords[wordLower]
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
