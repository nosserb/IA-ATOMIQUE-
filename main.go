package main

import (
	"IA-ATOMIQUE/database"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
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

func afficherAide() {
	fmt.Println(`
╔═══════════════════════════════════════════════════════════════╗
║          IA-ATOMIQUE v4.1 - Neural Text Analysis             ║
║    Analyse de texte, Apprentissage & Génération de résumés   ║
╚═══════════════════════════════════════════════════════════════╝

📖 UTILISATION:

  1. MODE INTERACTIF (par défaut)
     $ ./programme
     ou
     $ ./programme interactive

  2. ANALYSER UN FICHIER
     $ ./programme file <chemin_fichier>
     
     Exemples:
       $ ./programme file wikipedia.txt
       $ ./programme file document.md
       $ ./programme file /home/user/texte.txt

  3. ANALYSER UN TEXTE DIRECT
     $ ./programme text <votre texte>
     
     Exemples:
       $ ./programme text "L'IA est la technologie du futur"
       $ ./programme text "Einstein a découvert la relativité en 1905"

  4. MODE CLASSIQUE (hérité)
     $ ./programme <texte quelconque>
     
     Exemples:
       $ ./programme bonjour
       $ ./programme l'informatique progresse rapidement

═══════════════════════════════════════════════════════════════

⚙️  FONCTIONNALITÉS:

  ✓ Analyse phrase par phrase
  ✓ Classification en 6 catégories:
    - TECH: Technologie, informatique, numérique
    - HISTOIRE: Politique, histoire, événements
    - BUSINESS: Commerce, économie, affaires
    - ALIMENTATION: Nutrition, gastronomie
    - SANTE: Santé, médecine, bien-être
    - VERBE: Détection d'actions (verbes principaux)

  ✓ Détection de la structure grammaticale
  ✓ Génération de résumés synthétiques
  ✓ Système de modération (blacklist chiffrée AES-256)
  ✓ Apprentissage dynamique des mots

═══════════════════════════════════════════════════════════════

📊 SORTIE STANDARD:

  [STATISTIQUES GLOBALES]
    • Phrases analysées
    • Catégories détectées
    • Énergie totale
    • Confiance moyenne

  [DISTRIBUTION]
    • Répartition en pourcentages
    • Graphiques visuels

  [ANALYSE PAR PHRASE]
    • Catégorie détectée
    • Mots clés identifiés
    • Verbes principaux
    • Résumé synthétique

═══════════════════════════════════════════════════════════════

💡 CONSEILS D'UTILISATION:

  • Pour les meilleurs résultats, utilisez des textes français
  • Les fichiers volumineux seront traités progressivement
  • La modération empêche l'apprentissage de contenu offensant
  • Les neurones se régénèrent automatiquement après chaque exécution

═══════════════════════════════════════════════════════════════

❓ AIDE:
  $ ./programme help
  ou
  $ ./programme -h
  ou
  $ ./programme --help

═══════════════════════════════════════════════════════════════
`)
}

// API Response structure
type APIResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Result  string `json:"result"`
}

// ProcessRequest handles API requests for text processing
func ProcessRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "POST" {
		json.NewEncoder(w).Encode(APIResponse{
			Status:  "error",
			Message: "Method not allowed",
			Result:  "",
		})
		return
	}

	var req struct {
		Text string `json:"text"`
		Type string `json:"type"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(APIResponse{
			Status:  "error",
			Message: "Invalid request",
			Result:  "",
		})
		return
	}

	if req.Text == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(APIResponse{
			Status:  "error",
			Message: "Empty text",
			Result:  "",
		})
		return
	}

	// Capturer la sortie stdout
	oldStdout := os.Stdout
	readPipe, writePipe, _ := os.Pipe()
	os.Stdout = writePipe

	// Process the text
	TraiterTexte(req.Text, "API request")

	// Restaurer stdout et lire la sortie
	writePipe.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	io.Copy(&buf, readPipe)
	output := buf.String()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(APIResponse{
		Status:  "success",
		Message: "Text processed successfully",
		Result:  output,
	})
}

func main() {
	go AutoIntrospection()

	if len(os.Args) < 2 {
		fmt.Println("[IA-ATOMIQUE v4.1] Interface Interactive + Web Server")
		fmt.Println("Compréhension + Apprentissage + Résumé")

		// Lancer le serveur web
		StartWebServer()
		return
	}

	// Vérifier les commandes
	commande := os.Args[1]

	switch commande {
	case "help", "-h", "--help":
		afficherAide()

	case "file":
		if len(os.Args) > 2 {
			TraiterFichier(os.Args[2])
		} else {
			fmt.Println("[ERREUR] Utilisation: ./programme file <chemin>")
		}

	case "text":
		if len(os.Args) > 2 {
			texte := strings.Join(os.Args[2:], " ")
			TraiterTexte(texte, "entrée CLI")
		} else {
			fmt.Println("[ERREUR] Utilisation: ./programme text <votre texte>")
		}

	case "interactive", "inter":
		InteractionInteractive()

	case "web":
		fmt.Println("[IA-ATOMIQUE v4.1] Web Server Mode")
		fmt.Println("API disponible sur http://localhost:8080")
		StartWebServer()
		select {} // Keep the server running

	default:
		// Mode hérité - traiter comme phrase simple
		phrase := strings.Join(os.Args[1:], " ")
		TraiterTexte(phrase, "mode classique")
	}

	database.RegenererNeurones()
}

// StartWebServer starts the HTTP server
func StartWebServer() {
	// Get port from environment variable or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/api/process", ProcessRequest)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("."))))

	addr := ":" + port
	fmt.Println("Serveur lancé sur http://localhost:" + port)

	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal("Erreur serveur:", err)
	}
}
