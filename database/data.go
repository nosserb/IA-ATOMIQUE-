package database

// Word représente un mot avec sa catégorie et son numéro
type Word struct {
	Mot       string
	Categorie string
	Numero    int
}

// Liste enrichie de mots
var Words = []Word{
	// Ingénierie
	{"voiture", "Ingénierie", 1},
	{"moteur", "Ingénierie", 1},
	{"robot", "Ingénierie", 1},
	{"fusée", "Ingénierie", 1},
	{"pont", "Ingénierie", 1},
	{"camion", "Ingénierie", 1},
	{"train", "Ingénierie", 1},
	{"vélo", "Ingénierie", 1},
	{"voilier", "Ingénierie", 1},
	{"drone", "Ingénierie", 1},
	{"réacteur", "Ingénierie", 1},
	{"turbine", "Ingénierie", 1},
	{"machine", "Ingénierie", 1},
	{"engin", "Ingénierie", 1},
	{"grue", "Ingénierie", 1},
	{"ascenseur", "Ingénierie", 1},

	// Science
	{"physique", "Science", 2},
	{"chimie", "Science", 2},
	{"biologie", "Science", 2},
	{"astronomie", "Science", 2},
	{"géologie", "Science", 2},
	{"mathématique", "Science", 2},
	{"écologie", "Science", 2},
	{"astrophysique", "Science", 2},
	{"biologiste", "Science", 2},
	{"chimiste", "Science", 2},
	{"physicien", "Science", 2},
	{"microbiologie", "Science", 2},
	{"neurosciences", "Science", 2},
	{"médicament", "Science", 2},
	{"médicaments", "Science", 2},
	{"pharmacie", "Science", 2},
	{"médecin", "Science", 2},
	{"hôpital", "Science", 2},
	{"infirmier", "Science", 2},
	{"soin", "Science", 2},

	// Nature
	{"arbre", "Nature", 3},
	{"arbres", "Nature", 3},
	{"forêt", "Nature", 3},
	{"plante", "Nature", 3},
	{"plantes", "Nature", 3},
	{"montagne", "Nature", 3},
	{"rivière", "Nature", 3},
	{"lac", "Nature", 3},
	{"plage", "Nature", 3},
	{"jardin", "Nature", 3},
	{"forêt tropicale", "Nature", 3},
	{"plante carnivore", "Nature", 3},
	{"fleur", "Nature", 3},
	{"herbe", "Nature", 3},

	// Art
	{"peinture", "Art", 4},
	{"peintre", "Art", 4},
	{"dessin", "Art", 4},
	{"sculpture", "Art", 4},
	{"musique", "Art", 4},
	{"orchestre", "Art", 4},
	{"cinéma", "Art", 4},
	{"danse", "Art", 4},
	{"théâtre", "Art", 4},
	{"théâtre classique", "Art", 4},
	{"art numérique", "Art", 4},
	{"chorégraphie", "Art", 4},
	{"orchestration", "Art", 4},

	// Programmation
	{"ordinateur", "Programmation", 5},
	{"code", "Programmation", 5},
	{"algorithme", "Programmation", 5},
	{"langage", "Programmation", 5},
	{"script", "Programmation", 5},
	{"serveur", "Programmation", 5},
	{"microprocesseur", "Programmation", 5},
	{"application", "Programmation", 5},
	{"intelligence artificielle", "Programmation", 5},
	{"réseau", "Programmation", 5},
	{"programmation", "Programmation", 5},
	{"API", "Programmation", 5},
	{"base de données", "Programmation", 5},
	{"framework", "Programmation", 5},
	{"compilateur", "Programmation", 5},
	{"interface", "Programmation", 5},
	{"debugger", "Programmation", 5},
	{"fonction", "Programmation", 5},
	{"variable", "Programmation", 5},
	{"boucle", "Programmation", 5},
}
