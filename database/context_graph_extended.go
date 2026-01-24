// Package database - Extended Concept Graph
// Graphe de concepts étendu pour améliorer le contexte enrichi

package database

// initializeExtendedConceptGraph crée un graphe de concepts très étendu
func initializeExtendedConceptGraph() map[string][]string {
	return map[string][]string{
		// Histoire - très élargi
		"napoléon":   {"waterloo", "empire", "france", "empereur", "corse", "austerlitz", "joséphine", "1804", "1815", "bonaparte"},
		"révolution": {"1789", "bastille", "robespierre", "terreur", "girondins", "jacobins", "louis", "monarchie"},
		"guerre":     {"bataille", "conflit", "armée", "soldat", "victoire", "défaite", "paix", "combat"},
		"empire":     {"empereur", "colonie", "territoire", "puissance", "royaume", "napoléon"},
		"roi":        {"monarchie", "couronne", "château", "noblesse", "louis", "versailles", "règne"},
		"1789":       {"révolution", "bastille", "monarchie", "république", "louis"},
		"waterloo":   {"napoléon", "belgique", "défaite", "wellington", "1815", "bataille"},
		"1804":       {"napoléon", "empereur", "couronnement", "empire", "sacre"},
		"austerlitz": {"napoléon", "victoire", "1805", "bataille", "trois empereurs"},
		"bastille":   {"révolution", "1789", "prison", "prise", "juillet"},
		"louis":      {"roi", "versailles", "monarchie", "france", "règne"},
		"versailles": {"château", "roi", "louis", "france", "palais"},

		// Médecine - très élargi
		"hépatite":   {"foie", "virus", "maladie", "hépatique", "infection", "jaunisse", "inflammat"},
		"foie":       {"hépatite", "organe", "digestion", "bile", "cirrhose", "détoxification"},
		"cœur":       {"cardiaque", "artère", "circulation", "pompe", "ventricule", "battement", "sang"},
		"poumon":     {"respiration", "air", "oxygène", "alvéole", "bronche", "thorax"},
		"cerveau":    {"neurologie", "neurone", "pensée", "intelligence", "mémoire", "cortex"},
		"maladie":    {"symptôme", "diagnostic", "traitement", "médecine", "santé", "pathologie"},
		"virus":      {"infection", "contagion", "vaccin", "immunité", "microbe"},
		"traitement": {"médicament", "thérapie", "soin", "guérison", "dose", "prescription"},
		"symptôme":   {"signe", "manifestation", "maladie", "diagnostic", "clinique"},
		"médecin":    {"docteur", "soigner", "patient", "diagnostic", "prescription"},
		"patient":    {"malade", "médecin", "soigner", "traitement", "consultation"},
		"organe":     {"foie", "cœur", "poumon", "corps", "biologie", "anatomie"},

		// Mathématiques - très élargi
		"équation":  {"algèbre", "variable", "solution", "résoudre", "inconnue", "x", "y"},
		"racine":    {"carré", "mathématique", "radical", "puissance", "√", "nombre"},
		"théorème":  {"pythagore", "démonstration", "preuve", "mathématique", "géométrie", "postulat"},
		"carré":     {"racine", "puissance", "deux", "multiplication", "quadrat", "²"},
		"144":       {"racine", "12", "carré", "multiplication", "nombre"},
		"12":        {"racine", "144", "carré", "douze", "multiple"},
		"pythagore": {"théorème", "triangle", "hypoténuse", "géométrie", "a²+b²=c²"},
		"triangle":  {"géométrie", "pythagore", "côté", "angle", "polygone"},
		"algèbre":   {"équation", "variable", "mathématique", "résoudre", "calcul"},
		"géométrie": {"triangle", "cercle", "angle", "figure", "espace"},
		"nombre":    {"mathématique", "chiffre", "quantité", "calcul", "entier"},

		// Sciences - très élargi
		"atome":    {"électron", "proton", "neutron", "noyau", "particule", "chimie", "molécule"},
		"cellule":  {"biologie", "organisme", "vivant", "division", "membrane", "adn", "noyau"},
		"énergie":  {"physique", "travail", "force", "puissance", "joule", "watt", "conservation"},
		"lumière":  {"photon", "onde", "optique", "couleur", "vitesse", "visible"},
		"électron": {"atome", "charge", "orbital", "particule", "électrique", "négatif"},
		"proton":   {"atome", "charge", "positif", "noyau", "particule"},
		"neutron":  {"atome", "noyau", "masse", "particule", "neutre"},
		"molécule": {"atome", "chimie", "liaison", "composé", "substance"},
		"adn":      {"génétique", "cellule", "hérédité", "chromosome", "gène"},
		"physique": {"science", "énergie", "force", "matière", "loi"},
		"chimie":   {"atome", "molécule", "réaction", "élément", "substance"},

		// Actions de cuisine - très élargi
		"cuisine":   {"casserole", "cuire", "manger", "préparer", "recette", "cuisson", "ingrédient", "plat"},
		"casserole": {"cuisine", "cuire", "eau", "feu", "bouillir", "pot", "ustensile"},
		"eau":       {"liquide", "bouillir", "boire", "h2o", "chauffer", "vapeur"},
		"feu":       {"chaleur", "cuisson", "cuisiner", "brûler", "flamme", "chauffer"},
		"bouillir":  {"eau", "température", "vapeur", "ébullition", "100°", "chauffer"},
		"cuire":     {"cuisine", "cuisson", "chaleur", "préparer", "feu", "four"},
		"manger":    {"nourriture", "repas", "aliment", "cuisine", "ingérer"},
		"préparer":  {"cuisine", "recette", "ingrédient", "faire", "confectionner"},
		"recette":   {"cuisine", "ingrédient", "préparer", "étapes", "plat"},
		"attendre":  {"patienter", "temps", "durée", "rester"},

		// Actions de sport - très élargi
		"sport":           {"courir", "entraînement", "exercice", "fitness", "musculaire", "cardio", "athlète"},
		"courir":          {"sport", "course", "jogger", "marcher", "cardio", "pieds", "vitesse"},
		"chaussures":      {"sport", "pieds", "marcher", "courir", "running", "paire"},
		"exercice":        {"sport", "musculation", "cardio", "entraînement", "mouvement"},
		"entraînement":    {"sport", "exercice", "pratique", "préparation", "session"},
		"marcher":         {"courir", "pieds", "déplacer", "pas", "mouvement"},
		"ralentir":        {"vitesse", "diminuer", "freiner", "décélérer", "progressif"},
		"progressivement": {"graduellement", "peu à peu", "lentement", "étape"},

		// Actions de travail - très élargi
		"travail":     {"bureau", "ordinateur", "projet", "tâche", "collègue", "réunion", "emploi"},
		"ordinateur":  {"informatique", "écran", "clavier", "travail", "pc", "souris"},
		"taper":       {"clavier", "écrire", "texte", "ordinateur", "saisir"},
		"rapport":     {"document", "travail", "rédiger", "texte", "compte-rendu"},
		"sauvegarder": {"fichier", "ordinateur", "enregistrer", "données", "backup"},
		"étudiant":    {"université", "cours", "examen", "diplôme", "apprendre", "école"},
		"rédiger":     {"écrire", "texte", "document", "composer", "rapport"},
		"heures":      {"temps", "durée", "minutes", "période", "horaire"},

		// Personnes et pronoms
		"femme":    {"personne", "elle", "dame", "féminin", "madame"},
		"homme":    {"personne", "il", "monsieur", "masculin", "homme"},
		"personne": {"individu", "gens", "quelqu'un", "être humain"},
		"il":       {"homme", "masculin", "personne", "pronom"},
		"elle":     {"femme", "féminin", "personne", "pronom"},

		// Actions et mouvements
		"entrer":    {"pénétrer", "aller dans", "accéder", "franchir"},
		"sortir":    {"quitter", "partir", "extraire", "dehors"},
		"prendre":   {"saisir", "attraper", "récupérer", "tenir"},
		"mettre":    {"placer", "poser", "installer", "déposer"},
		"ouvrir":    {"déverrouiller", "accéder", "démarrer", "débuter"},
		"commencer": {"débuter", "entamer", "initier", "démarrer"},
		"continuer": {"poursuivre", "maintenir", "persévérer", "persister"},
		"lancer":    {"démarrer", "initier", "ouvrir", "commencer"},

		// Temps et durée
		"pendant": {"durée", "au cours de", "durant", "temps"},
		"minutes": {"temps", "durée", "heures", "secondes"},
		"puis":    {"ensuite", "après", "alors", "séquence"},
		"ensuite": {"puis", "après", "alors", "suivant"},
		"après":   {"ensuite", "puis", "ultérieurement", "suivant"},
		"alors":   {"ensuite", "puis", "donc", "conséquence"},

		// Lieux
		"pièce":  {"maison", "chambre", "espace", "intérieur"},
		"rue":    {"route", "chemin", "voie", "extérieur"},
		"chez":   {"maison", "domicile", "habitation", "résidence"},
		"maison": {"domicile", "habitation", "logement", "résidence"},
	}
}
