package database

// Word représente un mot avec sa catégorie et son numéro
type Word struct {
	Mot       string
	Categorie string
	Numero    int
}

// Liste de mots et leur classification
var Words = []Word{
	{"voiture", "Ingénierie", 1},
	{"moteur", "Ingénierie", 1},
	{"ordinateur", "Programmation", 5},
	{"code", "Programmation", 5},
	{"physique", "Science", 2},
	{"chimie", "Science", 2},
	{"forêt", "Nature", 3},
	{"arbre", "Nature", 3},
	{"peinture", "Art", 4},
	{"sculpture", "Art", 4},
	{"robot", "Ingénierie", 1},
	{"algorithme", "Programmation", 5},
	{"biologie", "Science", 2},
	{"montagne", "Nature", 3},
	{"dessin", "Art", 4},
	{"voilier", "Ingénierie", 1},
	{"réseau", "Programmation", 5},
	{"physicien", "Science", 2},
	{"plante", "Nature", 3},
	{"musique", "Art", 4},
	{"camion", "Ingénierie", 1},
	{"serveur", "Programmation", 5},
	{"astronomie", "Science", 2},
	{"rivière", "Nature", 3},
	{"cinéma", "Art", 4},
	{"train", "Ingénierie", 1},
	{"microprocesseur", "Programmation", 5},
	{"mathématique", "Science", 2},
	{"lac", "Nature", 3},
	{"théâtre", "Art", 4},
	{"drone", "Ingénierie", 1},
	{"langage", "Programmation", 5},
	{"chimiste", "Science", 2},
	{"plage", "Nature", 3},
	{"danse", "Art", 4},
	{"pont", "Ingénierie", 1},
	{"script", "Programmation", 5},
	{"géologie", "Science", 2},
	{"jardin", "Nature", 3},
	{"peintre", "Art", 4},
	{"fusée", "Ingénierie", 1},
	{"application", "Programmation", 5},
	{"biologiste", "Science", 2},
	{"forêt tropicale", "Nature", 3},
	{"théâtre classique", "Art", 4},
	{"voilier de course", "Ingénierie", 1},
	{"intelligence artificielle", "Programmation", 5},
	{"astronaute", "Science", 2},
	{"plante carnivore", "Nature", 3},
	{"orchestre", "Art", 4},
}
