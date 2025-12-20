package database

import (
	"bufio"
	"fmt"
	"math/rand"
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

func init() {
	for i := 0; i < 1000; i++ {
		Neurones = append(Neurones, Neurone{ID: i, CategorieID: rand.Intn(60) + 1})
	}
	Words = make(map[string]Word)

	// Catégorie 0 : NEUTRE
	Injecter(0, 0.1, "bonjour", "salut", "merci", "svp", "donc", "alors", "mais", "avec", "que", "qui", "est")

	Injecter(1, 5.0, "ia", "robot", "ordinateur", "technologie", "code")
	Injecter(4, 5.0, "vendre", "entreprise", "projet", "business", "argent")
	Injecter(50, 5.0, "manger", "faim", "nourriture", "pizza", "pates")
	Injecter(6, 5.0, "mal", "santé", "douleur", "hopital")

	ChargerLexique("lexique.txt")
	ChargerProbation("temp.txt")
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
	categories := map[int]string{0: "NEUTRE", 1: "TECH", 4: "BUSINESS", 6: "SANTE", 50: "BESOIN"}
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
