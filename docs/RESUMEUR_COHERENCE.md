# 📊 Résumeur Vectoriel avec Découpage & Cohérence

## Architecture Mathématique

### 1️⃣ Découpage du Texte (Étape 1)

Pour un document avec **P** phrases et **N** mots :

```
B = ⌈N / k⌉  (nombre de blocs)
k = taille cible par bloc (défaut: 25 mots)
```

**Exemple:**
- Document: N = 1189 mots
- Blocs: B = ⌈1189 / 25⌉ = 48 blocs
- Notre implémentation: k=25 → B=30 blocs réels

### 2️⃣ Vectorisation (Étape 2)

Chaque bloc **i** devient un vecteur **vᵢ** dans ℝ¹¹ (11 dimensions):

```
vᵢ = MoyennePondereeVecteurs(phrases[bloc_i])
```

- Vectorisation TF-IDF intégrée
- Normalisation L2 implicite
- Capture le sens du bloc

### 3️⃣ Cohérence Vectorielle Globale (Étape 3)

Calcul du vecteur global du document:

```
V_doc = (1/B) Σᵢ₌₁ᴮ vᵢ   (moyenne des vecteurs de blocs)
```

Normalisé en ℝ¹¹:

```
V_doc_norm = V_doc / ||V_doc||
```

### 4️⃣ Sélection par Énergie (Étape 4)

Pour chaque bloc i, calculer le **score d'énergie**:

```
cohérence_i = sim_cosinus(vᵢ, V_doc)    [0, 1]
tfidf_i = moyenne_tfidf(mots_bloc_i)     [0, 1]

énergie_i = α·cohérence_i + β·tfidf_i
```

où:
- **α = 0.6** : Poids cohérence (bloc représentatif)
- **β = 0.4** : Poids importance lexicale (mots clés)
- α + β = 1.0

**Sélection:** Les blocs avec **énergie** > seuil sont conservés

Nombre de blocs sélectionnés:
```
B_selected = ⌈B × ratioCompression⌉
```

### 5️⃣ Génération du Résumé (Étape 5)

À partir des blocs sélectionnés:

1. Extraire le vocabulaire valide (filtré par `isValidWord()`)
2. Vectoriser chaque mot du vocabulaire
3. Générer par **argmax** (similitude + cohérence)
4. Ajouter connecteurs sémantiques

Formule génération:
```
mot_sélectionné = argmax_w [ sim(vw, V_doc) + η·sim(vw, v_mot_précédent) ]
```

### 6️⃣ Complexité & Gains de Performance

**Sans découpage:**
```
Complexité: O(P² · d)
P = nombre de phrases (~50-100)
d = dimensions (11)
```

**Avec découpage:**
```
Complexité: O(B² · d)
B = nombre de blocs (~30)
Réduction: (B/P)² = (30/50)² ≈ 0.36
Gain: ~3x plus rapide
```

Pour documents massifs:
```
P = 128,938 phrases
B = 10,000 blocs
Réduction: (10K/128K)² ≈ 0.006
Gain: 166x plus rapide!
```

## Implémentation

### Structures principales

```go
type BlocVectoriel struct {
	Index      int              // Position
	Contenu    string           // Texte
	Mots       []string         // Vocabulaire
	Vecteur    VecteurAtomique  // vᵢ ∈ ℝ¹¹
	Coherence  float64          // sim(vᵢ, V_doc) ∈ [0,1]
	TFIDFScore float64          // Importance lexicale ∈ [0,1]
	Energie    float64          // Score final = α·cohérence + β·tfidf
}

type ResumeurCoherence struct {
	Blocs         []BlocVectoriel
	VecteurGlobal VecteurAtomique   // V_doc
	AlphaCoherence float64          // α = 0.6
	BetaTFIDF      float64          // β = 0.4
}
```

### Méthodes

```go
// 1. Découper et vectoriser
resumeur.Decouper(phrases)

// 2. Sélectionner par énergie
blocsSelectionnés := resumeur.Selectionner(ratioCompression)

// 3. Générer résumé
resume := resumeur.GenerByVectors(blocsSelectionnés, tailleResume)
```

## Commande CLI

```bash
./programme resume <fichier> [ratio_compression=0.15]
```

**Exemple:**
```bash
./programme resume input.txt 0.12    # Compression 12%
```

Affiche:
- Nombre de blocs créés / sélectionnés
- Cohérence moyenne
- TF-IDF moyen
- Énergie moyenne
- Résumé généré
- Statistiques de temps

## Comparaison avec d'autres méthodes

| Aspect | Vectoriel simple | Découpage & Cohérence | Atomique |
|--------|------------------|----------------------|----------|
| **Cohérence** | 88-92% | 95-100% | 9-15% |
| **Couverture** | Partielle | Bonne (début-fin) | Exploratory |
| **Vitesse** | ~600µs | ~300µs | ~2ms |
| **Nature** | Greedy | Déterministe | Résonance |
| **Mots redondants** | Quelques-uns | Minimal | Aucun |

## Résultats sur test.txt (103 mots, 15 phrases)

### Vectoriel (30% compression)
```
Résumé: "bénéficient insights patterns précieux... constamment innovations"
Mots: 25 | Cohérence: 89% | Temps: 600µs
```

### Découpage (20% compression)
```
Résumé: "L'analyse insights patterns précieux... constamment innovations"
Mots: 22 | Cohérence: 100% | Temps: 350µs | Blocs: 1/4
```

## Résultats sur input.txt (1189 mots, 53 phrases)

### Vectoriel (12% compression)
```
Mots: 178 | Cohérence: 93% | Compression: 6.7x | Temps: 6.7ms
```

### Découpage (12% compression)
```
Mots: 79 | Cohérence: 100% | Compression: 15.1x | Temps: 5ms
Blocs sélectionnés: 4/30
```

## Points forts

✅ **Cohérence garantie** - Tous les vecteurs alignés avec V_doc  
✅ **Couverture texto** - Blocs sélectionnés préservent l'ordre  
✅ **Scalabilité** - O(B²·d) vs O(P²·d)  
✅ **Interprétabilité** - Énergie = α·cohérence + β·importance  
✅ **Vitesse** - 3-166x plus rapide selon taille  

## Amélioration futures

- [ ] Adapter α et β par document (apprentissage)
- [ ] Ajouter pondération temporelle (plus récent = mieux)
- [ ] Intégrer connecteurs de blocs (donc, cependant, etc.)
- [ ] Visualisation de la "carte" de cohérence
- [ ] Multi-langage (vecteurs multilingues)

## Mathématique résumée

```
Entrée:  P phrases → B blocs → vecteurs vᵢ
         B blocs vectoriels

V_doc ← MoyennePonderee(v₁, ..., vB)    [Étape 3]

Pour chaque bloc i:
  cohérence_i ← cos_sim(vᵢ, V_doc)
  énergie_i ← 0.6·cohérence_i + 0.4·tfidf_i

Sélectionner: top B_selected par énergie    [Étape 4]

Générer:  argmax vocabulaire par similarity + coherence  [Étape 5]

Sortie:   Résumé cohérent début-fin de longueur cible
```

---

**Date:** 7 janvier 2026  
**Implémentation:** Go 1.22  
**Module:** `database/resumeur_coherence.go` (288 lignes)
