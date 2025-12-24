package database

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	cryptorand "crypto/rand"
	"crypto/sha256"
	"encoding/base64"
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

	// Catégorie 1 : TECH - Très spécifique + technologie générale
	Injecter(1, 6.0, "ia", "robot", "ordinateur", "code", "logiciel", "programme", "serveur", "python", "javascript", "golang",
		"api", "base de données", "algorithme", "machine learning", "deep learning", "neural network", "cpu", "gpu", "cloud",
		"database", "application", "software", "hardware", "processor", "memory", "cache", "encryption", "cybersecurity",
		"technologie", "technique", "outils", "digital", "numérique", "informatique", "électronique", "automatisation",
		"innovation", "système", "méthode", "processus", "développement", "engineering", "computeur", "données", "internet")

	// Catégorie 2 : HISTOIRE - Châteaux, monuments, politique
	Injecter(2, 6.0, "château", "donjon", "fortification", "tour", "enceinte", "monument", "médiéval", "siècle", "construction",
		"historique", "ancien", "ruine", "rempart", "muraille", "forteresse", "citadelle", "pièce", "seigneur", "roi", "reine",
		"cour", "courtyard", "france", "période", "époque", "gouvernement", "politique", "parlement", "parlementaire", "dissolution",
		"législature", "éléctions", "anticipées", "constitution", "loi", "légal", "légalité", "régime", "démocratique", "mandat",
		"députés", "assemblée", "sénat", "vote", "suffrage", "maire", "commune", "département", "région",
		"histoire", "culture", "établissement", "géographie", "sculpture", "art", "ancienne", "conservation", "longévité",
		"antique", "calligraphie", "accessoire", "locomotive", "écriture", "observation", "animal", "zoologie",
		"époque", "construction", "fruit", "boisson", "pain", "cuisine", "époque")

	// Catégorie 3 : BUSINESS - Commerce, affaires
	Injecter(3, 6.0, "vendre", "entreprise", "business", "argent", "profit", "commerce", "client", "marché", "stratégie",
		"vente", "achat", "prix", "revenue", "startup", "compagnie", "négociant", "transaction", "contrat", "accord",
		"affaires", "transport", "industrie", "mode", "cosmétique", "voyage", "manufacturier")

	// Catégorie 4 : ALIMENTATION - Nourriture
	Injecter(4, 6.0, "manger", "nourriture", "pizza", "pates", "pâtes", "aliment", "cuisine", "restaurant", "recette", "faim",
		"cuire", "sauce", "fromage", "pain", "viande", "légume", "fruit", "boisson", "café", "thé", "vin", "plat", "assiette",
		"récolte", "saveur", "farine", "blé", "alimentation", "herbivore", "écologie",
		"cuisine", "boulanger", "aliments", "repas", "épices", "goût", "recettes")

	// Catégorie 5 : SANTÉ - Médecine
	Injecter(5, 6.0, "santé", "maladie", "médecin", "hôpital", "patient", "traitement", "douleur", "mal", "symptôme",
		"cure", "remède", "médecine", "pharmacie", "allergie", "virus", "infection", "diagnostic", "test", "vaccin")

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
	ChargerBlacklistChiffrée("blacklist.enc") // Charger la blacklist chiffrée
}

func Injecter(cat int, poids float64, mots ...string) {
	for _, m := range mots {
		Words[strings.ToLower(m)] = Word{Mot: m, Categorie: cat, Poids: poids}
	}
}

func Apprendre(mot string, catID int) {
	mot = strings.ToLower(mot)
	if StopWords[mot] || len(mot) <= 2 {
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
			fmt.Printf("[MIGRATION] '%s' quitte %s pour %s\n", mot, NumeroVersCategorie(w.Categorie), NumeroVersCategorie(catID))
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
				fmt.Printf("[ADOPTION] '%s' gravé dans %s !\n", mot, NumeroVersCategorie(catID))
			} else {
				fmt.Printf("[APPRENTISSAGE] '%s' (%d/2) dans %s\n", mot, val.Compteur, NumeroVersCategorie(catID))
			}
		} else {
			val.Compteur--
			if val.Compteur <= 0 {
				delete(LexiqueTemp, mot)
			}
		}
	} else {
		LexiqueTemp[mot] = &MotEnAttente{Categorie: catID, Compteur: 1}
		fmt.Printf("[NOUVEAU] '%s' (1/2) dans %s\n", mot, NumeroVersCategorie(catID))
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

	fmt.Printf("[BLACKLIST] %d mots interdits chargés\n", compteur)
}

// ═══════════════════════════════════════════════════════════════════════════
// FONCTIONS DE CHIFFREMENT - Blacklist sécurisée AES-256
// ═══════════════════════════════════════════════════════════════════════════

// Clé de chiffrement (32 bytes pour AES-256) - À garder secrète
var encryptionKey = []byte{
	0x4e, 0x4f, 0x53, 0x53, 0x45, 0x52, 0x42, 0x2d, // NOSSERB-
	0x49, 0x41, 0x2d, 0x41, 0x54, 0x4f, 0x4d, 0x49, // IA-ATOMI
	0x51, 0x55, 0x45, 0x2d, 0x42, 0x4c, 0x41, 0x43, // QUE-BLAC
	0x4b, 0x4c, 0x49, 0x53, 0x54, 0x2d, 0x53, 0x31, // KLIS-S1
}

// ChiffrerBlacklist chiffre le contenu de la blacklist avec AES-256
func ChiffrerBlacklist(texteOriginal string) (string, error) {
	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(cryptorand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(texteOriginal), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DéchiffrerBlacklist déchiffre le contenu de la blacklist
func DéchiffrerBlacklist(texteChiffré string) (string, error) {
	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	ciphertext, err := base64.StdEncoding.DecodeString(texteChiffré)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", fmt.Errorf("ciphertext trop court")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// VérifierIntéritéBlacklist() - Vérifie que le fichier n'a pas été déplacé, renommé ou modifié
func VérifierIntéritéBlacklist(nomFichier string) bool {
	file, err := os.Open(nomFichier)
	if err != nil {
		fmt.Printf("❌ [SÉCURITÉ] Fichier blacklist.enc introuvable ou inaccessible\n")
		fmt.Printf("   Le fichier doit rester dans le répertoire racine du projet\n")
		fmt.Printf("   ⚠️  NE DOIT PAS être déplacé, renommé ou supprimé\n")
		return false
	}
	defer file.Close()

	// Lire le fichier et calculer le hash SHA256
	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		fmt.Printf("❌ [SÉCURITÉ] Impossible de vérifier l'intégrité du fichier\n")
		return false
	}

	fichierHash := hex.EncodeToString(hash.Sum(nil))

	// Vérifier que le hash correspond
	if fichierHash != expectedBlacklistHash {
		fmt.Printf("❌ [SÉCURITÉ] ALERTE - Fichier blacklist.enc modifié ou corrompu!\n")
		fmt.Printf("   Hash attendu:  %s\n", expectedBlacklistHash)
		fmt.Printf("   Hash détecté:  %s\n", fichierHash)
		fmt.Printf("   ⚠️  Le fichier a peut-être été déplacé, renommé ou modifié\n")
		fmt.Printf("   ✓ Restaurez le fichier original depuis git\n")
		return false
	}

	return true
}

// ChargerBlacklistChiffrée charge et déchiffre la blacklist depuis un fichier .enc
func ChargerBlacklistChiffrée(nomFichier string) {
	// Vérifier d'abord l'intégrité du fichier
	if !VérifierIntéritéBlacklist(nomFichier) {
		fmt.Printf("❌ Impossible de charger la blacklist - vérification de sécurité échouée\n")
		fmt.Printf("❌ LE PROGRAMME S'ARRÊTE - Restore le fichier blacklist.enc\n")
		os.Exit(1)
	}

	// Essayer d'abord le fichier chiffré
	file, err := os.Open(nomFichier)
	if err == nil {
		defer file.Close()

		scanner := bufio.NewScanner(file)
		if scanner.Scan() {
			texteChiffré := scanner.Text()
			texteOriginal, err := DéchiffrerBlacklist(texteChiffré)
			if err != nil {
				fmt.Printf("[ERREUR] Impossible de déchiffrer la blacklist: %v\n", err)
				return
			}

			// Charger depuis le texte déchiffré
			compteur := 0
			for _, ligne := range strings.Split(texteOriginal, "\n") {
				ligne = strings.TrimSpace(ligne)
				if ligne == "" || strings.HasPrefix(ligne, "#") {
					continue
				}

				mot := strings.ToLower(ligne)
				Blacklist[mot] = true
				compteur++
			}

			fmt.Printf("[BLACKLIST] ✓ %d mots interdits chargés (chiffrement AES-256)\n", compteur)
			return
		}
	}

	// Sinon charger depuis le fichier texte
	fmt.Printf("[BLACKLIST] Fichier chiffré introuvable, chargement depuis texte clair...\n")
	ChargerBlacklist("blacklist.txt")
}
