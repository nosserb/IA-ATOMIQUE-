package database

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"math/rand"
	mathrand "math/rand"
	"os"
	"strconv"
	"strings"
	"unicode"
)

type Word struct {
	Mot       string
	Categorie int
	Poids     float64
}

type Neurone struct {
	ID, CategorieID int
	Valeur          float64
}

type MotEnAttente struct {
	Categorie int
	Compteur  int
}

var Neurones []Neurone
var Words map[string]Word
var Categories map[int]string
var Phrases []string
var LexiqueTemp = make(map[string]*MotEnAttente)
var StopWords = map[string]bool{"le": true, "la": true, "un": true, "une": true, "de": true, "je": true, "tu": true, "est": true, "et": true, "du": true, "des": true, "au": true, "les": true, "pour": true, "dans": true}
var Blacklist = make(map[string]bool) // Mots interdits

// Sécurité: Hash SHA256 du fichier blacklist.enc pour détection de déplacement/modification
const expectedBlacklistHash = "f2a21cd9b30c05caf4c5dce9e087dbeee72eb482b38bf0711cbc6c174bf7c421"

func init() {
	for i := 0; i < 1000; i++ {
		Neurones = append(Neurones, Neurone{ID: i, CategorieID: mathrand.Intn(60) + 1})
	}
	Words = make(map[string]Word)

	// Catégorie 0 : NEUTRE
	Injecter(0, 0.1, "bonjour", "salut", "merci", "svp", "donc", "alors", "mais", "avec", "que", "qui")

	// Catégorie 1 : TECH - Très spécifique INFORMATIQUE/DIGITAL UNIQUEMENT (Poids: 7.0)
	Injecter(1, 7.0, "ia", "robot", "ordinateur", "code", "logiciel", "programme", "serveur", "python", "javascript", "golang",
		"api", "base de données", "algorithme", "machine learning", "deep learning", "neural network", "cpu", "gpu", "cloud",
		"database", "application", "software", "hardware", "processor", "memory", "cache", "encryption", "cybersecurity",
		"technologie", "informatique", "électronique", "numérique", "digital", "internet", "développement", "données",
		"computing", "système informatique", "réseau", "streaming", "web", "backend", "frontend", "framework",
		"compilateur", "debug", "debugging", "variable", "fonction", "classe", "objet", "instance", "polymorphe",
		"versionning", "git", "github", "repository", "commit", "branching", "deploy", "serveur web", "requête http",
		"json", "xml", "html", "css", "typescript", "kotlin", "rust", "swift", "objective-c",
		"blockchain", "crypto", "bitcoin", "ethereum", "smart contract", "token", "nft", "web3",
		"sensor", "iot", "arduino", "raspberry", "microcontroller", "circuit", "électronique", "transistor",
		"quantum", "computing", "neural", "network", "tensor", "gpu", "cuda", "opengl", "api graphics",
		"database sql", "mongodb", "postgresql", "mysql", "redis", "cassandra", "elasticsearch",
		"docker", "kubernetes", "container", "virtualization", "cloud computing", "aws", "azure", "gcp",
		"devops", "ci/cd", "pipeline", "teste", "test unitaire", "integration test", "e2e", "selenium")

	// Catégorie 2 : HISTOIRE - STRICTEMENT HISTORIQUE (Poids: 6.5)
	Injecter(2, 6.5, "château", "donjon", "fortification", "tour", "monument", "médiéval", "siècle",
		"historique", "ancien", "ruine", "rempart", "muraille", "forteresse", "citadelle", "seigneur", "roi", "reine",
		"cour royale", "france", "période", "époque", "gouvernement", "politique", "parlement", "parlementaire", "dissolution",
		"législature", "élection", "constitution", "loi", "légal", "légalité", "régime", "démocratique", "mandat",
		"députés", "assemblée", "sénat", "vote", "suffrage", "maire", "commune", "département", "région",
		"histoire", "sculpture", "art", "antique", "calligraphie", "locomotive", "écriture",
		"civilisation", "conquête", "bataille", "général", "empire", "dynastie", "événement", "révolution", "monarque", "féodal",
		"héritage", "chronologie", "archéologie", "musée", "patrimoine", "vestige", "propagande",
		"historien", "chronique", "récit", "testament", "mémoire", "centenaire", "bicentenaire",
		"antiquité", "renaissance", "baroque", "classique", "romantique", "victorien",
		"pharaon", "egypte", "rome", "grèce", "persa", "viking", "normand", "saxon",
		"croisade", "inquisition", "enlightenment", "revolution française", "napoléon", "bonaparte")

	// Catégorie 3 : BUSINESS - Commerce, affaires
	Injecter(3, 6.0, "vendre", "entreprise", "business", "argent", "profit", "commerce", "client", "marché", "stratégie",
		"vente", "achat", "prix", "revenue", "startup", "compagnie", "négociant", "transaction", "contrat", "accord",
		"affaires", "transport", "industrie", "mode", "cosmétique", "voyage", "manufacturier", "commercial", "grossiste",
		"détaillant", "franchise", "partenariat", "fusion", "acquisition", "investissement", "capital", "action", "dividende",
		"bourse", "portefeuille", "rendement", "dépense", "budget", "comptabilité", "facture", "devis", "salaire", "paie",
		"emploi", "métier", "profession", "carrière", "promotion", "augmentation", "congé", "retraite", "syndicat", "productivité")

	// Catégorie 4 : ALIMENTATION - Nourriture (SPÉCIFIQUE!)
	Injecter(4, 6.0, "manger", "nourriture", "pizza", "pates", "pâtes", "aliment", "cuisine", "restaurant", "recette", "faim",
		"cuire", "sauce", "fromage", "pain", "viande", "légume", "fruit", "boisson", "café", "thé", "vin", "plat", "assiette",
		"saveur", "farine", "blé", "alimentation", "herbivore", "boulanger", "aliments", "repas", "épices", "goût", "recettes",
		"bouche", "appétit", "sucre", "sel", "épice", "herbe", "poisson", "poulet", "boeuf", "porc", "charcuterie",
		"dessert", "gâteau", "biscuit", "chocolat", "bonbon", "sucrerie", "entrée", "plat principal", "accompagnement",
		"sauce tomate", "huile", "beurre", "crème", "lait", "fromage blanc", "yaourt", "oeufs", "miel",
		"nutrition", "calorie", "protéine", "glucide", "lipide", "vitamine", "minéral", "régime", "kcal",
		"gastronomie", "chef", "cuisinier", "resto", "barbecue", "pique-nique", "festin", "banquet", "lunch")

	// Catégorie 5 : SANTE - Médecine STRICTEMENT (Poids: 5.5 - réduit)
	Injecter(5, 5.5, "santé", "médecin", "hôpital", "patient", "traitement", "médecine", "pharmacie", "allergie",
		"virus", "infection", "diagnostic", "vaccin", "cure", "remède", "symptôme",
		"hypertension", "arythmie", "angine", "asthme", "cancer", "carcinome", "tumeur", "leucémie",
		"grippe", "pneumonie", "tuberculose", "bronchite", "gastrite", "eczéma", "psoriasis",
		"chirurgie", "opération", "scalpel", "anesthésie", "transfusion", "greffe", "implant",
		"infirmière", "dentiste", "cardiologue", "neurologique", "dermatologue", "urologue", "psychologue",
		"anticorps", "immunité", "antibiotique", "antiviral", "anti-inflammatoire", "analgésique", "sédatif",
		"épidémie", "pandémie", "contagion", "quarantaine", "isolement", "dépistage", "test", "pcr",
		"prévention", "dépistage", "hygiène", "stérilisation", "asepsie", "désinfection",
		"urgence", "triage", "ambulance", "urgentiste", "traumatologie", "orthopédie", "rhumatologie",
		"prothèse", "fauteuil roulant", "pacemaker", "défibrillateur", "endoprothèse",
		"échographie", "radiographie", "irm", "scanner", "tomographie", "résonnance magnétique",
		"pathologie", "histologie", "génétique", "biochimie", "hématologie", "microbiologie")

	// Catégorie 6 : VERBE - Actions et verbes purs (infinitif + conjugaisons)
	Injecter(6, 5.0, "faire", "aller", "venir", "courir", "sauter", "parler", "écouter", "regarder", "voir", "dire",
		"penser", "croire", "savoir", "pouvoir", "vouloir", "devoir", "aimer", "détester", "rire", "pleurer",
		"boire", "dormir", "travailler", "jouer", "danser", "chanter", "construire", "détruire", "créer",
		"ouvrir", "fermer", "entrer", "sortir", "monter", "descendre", "tomber", "lever", "baisser", "tourner",
		"respirer", "marcher", "coucher", "asseoir", "tenir", "prendre", "donner", "recevoir", "envoyer", "chercher",
		// Formes conjuguées courantes du texte
		"aboie", "fleurit", "traite", "traverse", "enregistre", "affiche", "joue", "raconte", "survole", "brille",
		"vole", "ouvre", "livre", "dessine", "chauffe", "descend", "saute", "réchauffe", "flotte", "attire",
		"nage", "rebondit", "prépare", "organise", "balaie", "protègent", "indique", "glisse", "relie",
		"se", "pose", "contient", "part", "sent", "corrige", "sonne", "restaurée", "roule", "garde", "dort",
		"allume", "aide", "pédaler", "dresse", "décolle", "calcule", "mange", "fume", "recharge")

	ChargerLexique("lexique.txt")
	ChargerProbation("temp.txt")
	ChargerBlacklist("blacklist.enc") // Charger la blacklist avec vérification d'intégrité SHA256
}

func Injecter(cat int, poids float64, mots ...string) {
	for _, m := range mots {
		Words[strings.ToLower(m)] = Word{Mot: m, Categorie: cat, Poids: poids}
	}
}

// DeterLangueMot détecte la langue d'un mot isolé
func DeterLangueMot(mot string) string {
	lower := strings.ToLower(mot)

	// Caractères typiques de chaque langue
	if strings.ContainsAny(lower, "àâäæéèêëïîôùûüœç") {
		return "fr" // Caractères français/néerlandais
	}

	// Patterns allemands
	if strings.ContainsAny(lower, "äöüß") ||
		strings.Contains(lower, "sch") ||
		strings.Contains(lower, "tsch") ||
		strings.HasSuffix(lower, "heit") ||
		strings.HasSuffix(lower, "keit") {
		return "de"
	}

	// Patterns anglais
	if strings.HasSuffix(lower, "ing") ||
		strings.HasSuffix(lower, "tion") ||
		strings.HasSuffix(lower, "ness") ||
		strings.HasSuffix(lower, "ment") && !strings.Contains(lower, "é") {
		return "en"
	}

	// Patterns espagnols
	if strings.HasSuffix(lower, "ción") ||
		strings.HasSuffix(lower, "ado") ||
		strings.HasSuffix(lower, "ada") ||
		strings.Contains(lower, "ñ") {
		return "es"
	}

	// Default: français
	return "fr"
}

func Apprendre(mot string, catID int) {
	mot = strings.ToLower(mot)
	if StopWords[mot] || len(mot) <= 2 {
		return
	}

	// Vérifier que ce n'est pas un mot étranger
	if DeterLangueMot(mot) != "fr" {
		return
	}

	// Vérifier la blacklist
	if Blacklist[mot] {
		fmt.Printf("[BLOQUÉ] '%s' est dans la blacklist et ne sera pas appris\n", mot)
		return
	}

	// --- LOGIQUE DE MIGRATION ---
	if w, ok := Words[mot]; ok {
		if w.Categorie != catID && w.Categorie != 0 {
			// fmt.Printf("[MIGRATION] '%s' quitte %s pour %s\n", mot, NumeroVersCategorie(w.Categorie), NumeroVersCategorie(catID))
			Words[mot] = Word{Mot: mot, Categorie: catID, Poids: 3.0}
			MajLexiqueFichier(mot, catID)
			return
		}
		return
	}

	// --- APPRENTISSAGE RAPIDE (2 fois au lieu de 3) ---
	if val, ok := LexiqueTemp[mot]; ok {
		if val.Categorie == catID {
			val.Compteur++
			if val.Compteur >= 2 {
				Words[mot] = Word{Mot: mot, Categorie: catID, Poids: 3.0}
				delete(LexiqueTemp, mot)
				SauvegarderDefinitif(mot, catID)
				// fmt.Printf("[ADOPTION] '%s' gravé dans %s !\n", mot, NumeroVersCategorie(catID))
			} else {
				// fmt.Printf("[APPRENTISSAGE] '%s' (%d/2) dans %s\n", mot, val.Compteur, NumeroVersCategorie(catID))
			}
		} else {
			val.Compteur--
			if val.Compteur <= 0 {
				delete(LexiqueTemp, mot)
			}
		}
	} else {
		LexiqueTemp[mot] = &MotEnAttente{Categorie: catID, Compteur: 1}
		// fmt.Printf("[NOUVEAU] '%s' (1/2) dans %s\n", mot, NumeroVersCategorie(catID))
	}
	SauvegarderProbation()
}

func MajLexiqueFichier(mot string, nouvelleCat int) {
	input, _ := os.ReadFile("lexique.txt")
	lignes := strings.Split(string(input), "\n")
	trouve := false
	for i, ligne := range lignes {
		if strings.Contains(ligne, ":"+mot) {
			lignes[i] = fmt.Sprintf("%d:%s", nouvelleCat, mot)
			trouve = true
			break
		}
	}
	if trouve {
		output := strings.Join(lignes, "\n")
		os.WriteFile("lexique.txt", []byte(output), 0644)
	}
}

func SauvegarderDefinitif(mot string, catID int) {
	f, _ := os.OpenFile("lexique.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()
	f.WriteString(fmt.Sprintf("%d:%s\n", catID, mot))
}

func ChargerLexique(nomFichier string) {
	file, err := os.Open(nomFichier)
	if err != nil {
		return
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), ":")
		if len(parts) == 2 {
			cat, _ := strconv.Atoi(parts[0])
			Injecter(cat, 3.0, strings.TrimSpace(parts[1]))
		}
	}
}

func ChargerProbation(nomFichier string) {
	file, err := os.Open(nomFichier)
	if err != nil {
		return
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), ":")
		if len(parts) == 3 {
			cat, _ := strconv.Atoi(parts[0])
			count, _ := strconv.Atoi(parts[1])
			LexiqueTemp[parts[2]] = &MotEnAttente{Categorie: cat, Compteur: count}
		}
	}
}

func SauvegarderProbation() {
	f, _ := os.Create("temp.txt")
	defer f.Close()
	for mot, data := range LexiqueTemp {
		f.WriteString(fmt.Sprintf("%d:%d:%s\n", data.Categorie, data.Compteur, mot))
	}
}

func MotProche(token string) (Word, string) {
	clean := strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			return r
		}
		return -1
	}, strings.ToLower(token))
	if w, ok := Words[clean]; ok {
		return w, clean
	}
	return Word{}, clean
}

func NumeroVersCategorie(num int) string {
	categories := map[int]string{
		0: "NEUTRE",
		1: "TECH",
		2: "HISTOIRE",
		3: "BUSINESS",
		4: "ALIMENTATION",
		5: "SANTE",
		6: "VERBE",
	}
	if v, ok := categories[num]; ok {
		return v
	}
	return "AUTRE"
}

func RegenererNeurones() {
	for i := range Neurones {
		Neurones[i].Valeur *= 0.01
	}
}

func GenererReponse(catID int, motsPoses []string) string {
	// Récupère les mots connus de cette catégorie
	motsConnus := []string{}
	for mot, word := range Words {
		if word.Categorie == catID && mot != "" && len(mot) > 2 {
			motsConnus = append(motsConnus, word.Mot)
		}
	}

	sujets := map[int][]string{
		1:  {"Je", "L'IA", "Le robot", "L'ordinateur", "Cela", "Je peux"},
		4:  {"L'entreprise", "Le business", "Ce projet", "On", "C'est", "Nous"},
		6:  {"Le corps", "La santé", "On", "C'est", "La vie", "Cela"},
		50: {"Je", "On", "Tu", "Je voudrais", "Personnellement", "C'est bon"},
	}

	verbes := map[int][]string{
		1:  {"traite", "analyse", "utilise", "comprend", "exécute", "calcule", "gère"},
		4:  {"génère", "crée", "produit", "réalise", "livre", "développe", "fournit"},
		6:  {"améliore", "renforce", "maintient", "aide", "protège", "préserve", "favorise"},
		50: {"aime", "préfère", "mange", "consomme", "apprécie", "recommande", "adore"},
	}

	sujetsListe := sujets[catID]
	verbesListe := verbes[catID]

	if len(sujetsListe) == 0 || len(verbesListe) == 0 {
		return "Je comprends cette catégorie."
	}

	sujet := sujetsListe[rand.Intn(len(sujetsListe))]
	verbe := verbesListe[rand.Intn(len(verbesListe))]

	// Utilise les mots qu'elle connaît
	var complements []string
	if len(motsConnus) > 0 {
		// Cherche d'abord un mot qui matche la question
		for _, mp := range motsPoses {
			for _, mk := range motsConnus {
				if len(mp) > 2 && (strings.Contains(strings.ToLower(mk), mp) || strings.Contains(mp, strings.ToLower(mk))) {
					complements = append(complements, mk)
				}
			}
		}

		// Complète avec d'autres mots aléatoires
		for len(complements) < 2 && len(complements) < len(motsConnus) {
			idx := rand.Intn(len(motsConnus))
			m := motsConnus[idx]
			found := false
			for _, c := range complements {
				if c == m {
					found = true
					break
				}
			}
			if !found {
				complements = append(complements, m)
			}
		}
	}

	// Fallback si pas assez de mots
	if len(complements) == 0 {
		defaultComplements := map[int][]string{
			1:  {"l'information", "les données", "les processus"},
			4:  {"de la valeur", "des solutions", "des stratégies"},
			6:  {"le bien-être", "la vitalité", "l'équilibre"},
			50: {"sainement", "qualitativement", "régulièrement"},
		}
		complements = defaultComplements[catID]
	}

	// Construire la phrase
	var reponse string
	if len(complements) >= 2 {
		reponse = fmt.Sprintf("%s %s %s et %s.", sujet, verbe, complements[0], complements[1])
	} else if len(complements) == 1 {
		reponse = fmt.Sprintf("%s %s %s.", sujet, verbe, complements[0])
	} else {
		reponse = fmt.Sprintf("%s %s bien les concepts de cette catégorie.", sujet, verbe)
	}

	return reponse
}

// ChargerBlacklist charge les mots interdits depuis un fichier
func ChargerBlacklist(nomFichier string) {
	file, err := os.Open(nomFichier)
	if err != nil {
		fmt.Printf("[AVERTISSEMENT] Impossible de charger %s : %v\n", nomFichier, err)
		return
	}
	defer file.Close()

	// Vérifier l'intégrité du fichier (SHA256)
	if !VérifierIntéritéBlacklist(nomFichier) {
		fmt.Printf("[AVERTISSEMENT] Blacklist modifiée - utilisation quand même mais vérifiez git\n")
	}

	scanner := bufio.NewScanner(file)
	compteur := 0

	for scanner.Scan() {
		ligne := strings.TrimSpace(scanner.Text())
		// Ignorer les lignes vides et les commentaires
		if ligne == "" || strings.HasPrefix(ligne, "#") {
			continue
		}

		mot := strings.ToLower(ligne)
		Blacklist[mot] = true
		compteur++
	}

	fmt.Printf("[BLACKLIST] ✓ %d mots interdits chargés (intégrité vérifiée)\n", compteur)
}

// VérifierIntéritéBlacklist() - Vérifie que le fichier n'a pas été modifié (SHA256)
func VérifierIntéritéBlacklist(nomFichier string) bool {
	file, err := os.Open(nomFichier)
	if err != nil {
		return false // Fichier absent
	}
	defer file.Close()

	// Calculer le hash SHA256 du fichier
	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return false
	}

	fichierHash := hex.EncodeToString(hash.Sum(nil))

	// Vérifier que le hash correspond (sinon fichier modifié)
	if fichierHash != expectedBlacklistHash {
		fmt.Printf("⚠️  [ATTENTION] Blacklist modifiée! Hash ne correspond pas.\n")
		return false
	}

	return true
}

// InitDatabase initialise la base de données
func InitDatabase() {
	if len(Categories) == 0 {
		Categories = make(map[int]string)
		Categories[0] = "Neutre"
		Categories[1] = "Positif"
		Categories[2] = "Négatif"
		Categories[3] = "Question"
		Categories[4] = "Commande"
		Categories[5] = "Information"
		Categories[6] = "Feedback"
		Categories[7] = "Autre"
		Categories[8] = "Général"
		Categories[9] = "Spécifique"
	}

	if len(Phrases) == 0 {
		Phrases = make([]string, 0)
	}
}
