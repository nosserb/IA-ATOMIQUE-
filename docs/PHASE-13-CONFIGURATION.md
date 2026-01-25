# Guide Phase 13+++ - Configuration & Tuning

##  Quick Start

```bash
# Compiler
go build -o programme

# Resume standard (Phase 13+++)
./programme resume input.txt 0.12

# Resume avec style humanise
./programme humanize file -p input.txt
```

**Resultat**: Resume de qualite maximale, sans repetitions, lecture fluide.

---

##  Param�tres Configurables

### 1. **Normalisation Lexicale** (resumeur_coherence.go)

**Coefficient de Penalite**:
```go
// Ligne ~565
penalite += float64(count-2) * 0.1  //  Ajuster ce coefficient
```

| Coefficient | Effet |
|-------------|-------|
| 0.05 | Tr�s permissif, peu de penalite |
| **0.1** |  **Recommande** |
| 0.15 | Strict, forte deprioritisation |
| 0.20 | Tr�s strict, blocs repetitifs exclus |

**Exemple**:
```go
// Pour �tre plus permissif sur repetitions internes
penalite += float64(count-2) * 0.05

// Pour �tre ultra-strict
penalite += float64(count-2) * 0.2
```

### 2. **Ponderation TF-IDF** (generation.go)

**Seuils et Penalite**:
```go
// Ligne ~507
if idf[mot] > 0.5 && tf[mot] > 0.05 {  //  Ajuster seuils
    tfidfVal *= 0.8  //  Ajuster multiplicateur
}
```

| IDF | TF | Multiplicateur | Cas d'Usage |
|-----|----|----|---|
| >0.5 | >0.05 | **0.8** |  **Standard** |
| >0.6 | >0.10 | 0.8 | Plus selectif |
| >0.4 | >0.03 | 0.9 | Plus permissif |
| >0.7 | >0.15 | 0.7 | Ultra-strict |

**Configuration Avancee**:
```go
// Cas 1: Plus de mots (moins strict)
if idf[mot] > 0.6 && tf[mot] > 0.10 {
    tfidfVal *= 0.9
}

// Cas 2: Zero repetitions (ultra-strict)
if idf[mot] > 0.4 {  // N'importe quel mot rare
    tfidfVal *= 0.7
}
```

### 3. **Fen�trage Strict - Similarite Lexicale** (resumeur_coherence.go)

**Seuil de Diversite**:
```go
// Ligne ~735
if similarity > 0.6 {  //  Ajuster ce seuil
    delete(selectedIndices, i)
    continue
}
```

| Seuil | Similarite Max | Effet |
|-------|---|---|
| 0.4 | 40% vocab commun | Ultra-diversifie, peu de blocs |
| 0.5 | 50% vocab commun | Diversifie, selection moderee |
| **0.6** | 60% vocab commun |  **Recommande** |
| 0.7 | 70% vocab commun | Assez permissif, plus de blocs |
| 0.8 | 80% vocab commun | Tr�s permissif, couverture max |

**Exemple Pratique**:
```go
// Pour maximal diversity (moins de blocs)
if similarity > 0.4 {
    delete(selectedIndices, i)
}

// Pour couverture maximale
if similarity > 0.75 {
    delete(selectedIndices, i)
}

// Default (recommande)
if similarity > 0.6 {
    delete(selectedIndices, i)
}
```

### 4. **Anti-Repetition - Distance Minimale** (coherence.go)

**Distance Intra-Texte**:
```go
// Ligne ~420
if i-lastPos < 5 {  //  Ajuster distance
    continue  // Skip repetition
}
```

| Distance | Effet |
|----------|-------|
| <3 | Ultra-strict, "tr�s tr�s" impossible |
| <5 |  **Recommande**, repetitions bien separees |
| <7 | Modere, quelques proches repetitions acceptees |
| <10 | Permissif, repetitions espacees acceptees |

**Configuration**:
```go
// Ultra-strict: aucune proche repetition
if i-lastPos < 3 {
    continue
}

// Recommande: 5 mots minimum entre occurrences
if i-lastPos < 5 {
    continue
}

// Permissif: 10 mots minimum
if i-lastPos < 10 {
    continue
}
```

### 5. **Diversification Synonymes** (coherence.go)

**Frequence de Remplacement**:
```go
// Ligne ~460
if compteurMots[motClean] % 3 == 0 && synChoisi != motClean {
    // Remplacer tous les 3 occurrences
}
```

| Modulo | Frequence | Effet |
|--------|-----------|-------|
| % 2 | Tous les 2 | Tr�s diversifie (50% synonymes) |
| % 3 | Tous les 3 |  **Recommande** (33% synonymes) |
| % 4 | Tous les 4 | Modere (25% synonymes) |
| % 5 | Tous les 5 | Discret (20% synonymes) |

**Configuration**:
```go
// Tr�s varie (synonymes frequents)
if compteurMots[motClean] % 2 == 0 && synChoisi != motClean {
    motsFiltres = append(motsFiltres, synChoisi)
    continue
}

// Standard (recommande)
if compteurMots[motClean] % 3 == 0 && synChoisi != motClean {
    motsFiltres = append(motsFiltres, synChoisi)
    continue
}

// Discret (synonymes rares)
if compteurMots[motClean] % 5 == 0 && synChoisi != motClean {
    motsFiltres = append(motsFiltres, synChoisi)
    continue
}
```

---

##  Profils de Configuration Pre-definis

###  **Profil 1: Maximum Quality** (Zero Repetitions)

```go
// resumeur_coherence.go
penalite += float64(count-2) * 0.2  // Penalite forte
if similarity > 0.4 { delete(...) }  // Strict diversite

// generation.go
if idf[mot] > 0.5 && tf[mot] > 0.05 {
    tfidfVal *= 0.7  // Tr�s penalisant
}

// coherence.go
if i-lastPos < 7 { continue }  // Distance longue
if compteurMots[motClean] % 2 == 0 && ... {  // Synonymes frequents
    motsFiltres = append(motsFiltres, synChoisi)
}
```

**Resultat**: 400-500 mots, **zero** repetitions visibles, lecture tr�s naturelle.

---

###  **Profil 2: Balanced** (Recommande - Par Defaut)

```go
// resumeur_coherence.go
penalite += float64(count-2) * 0.1  // Penalite moderee
if similarity > 0.6 { delete(...) }  // Diversite moderee

// generation.go
if idf[mot] > 0.5 && tf[mot] > 0.05 {
    tfidfVal *= 0.8  // Standard
}

// coherence.go
if i-lastPos < 5 { continue }  // Distance standard
if compteurMots[motClean] % 3 == 0 && ... {  // Synonymes moderes
    motsFiltres = append(motsFiltres, synChoisi)
}
```

**Resultat**: 600-800 mots, 95% coherence, lecture fluide, quelques variantes.

---

###  **Profil 3: Maximum Coverage** (Plus de Contenu)

```go
// resumeur_coherence.go
penalite += float64(count-2) * 0.05  // Penalite faible
if similarity > 0.75 { delete(...) }  // Diversite permissive

// generation.go
if idf[mot] > 0.6 && tf[mot] > 0.10 {
    tfidfVal *= 0.9  // Moins penalisant
}

// coherence.go
if i-lastPos < 3 { continue }  // Distance courte
if compteurMots[motClean] % 5 == 0 && ... {  // Synonymes discrets
    motsFiltres = append(motsFiltres, synChoisi)
}

// Augmenter limite selection dans SelectionnerBlocsAvecFenetrageGlissant
if numBlocs > 75 { numBlocs = 75 }  // vs 50
```

**Resultat**: 1000-1200 mots, 92-94% coherence, couverture maximale.

---

##  Recettes de Tuning

### Cas 1: Resume Tr�s Court (100-200 mots)
```go
// �tre ultra-strict pour qualite
penalite += float64(count-2) * 0.2
if similarity > 0.4 { delete(...) }
if i-lastPos < 7 { continue }
if compteurMots[motClean] % 2 == 0 { synonymes }
```

### Cas 2: Resume Standard (500-800 mots)
```go
// Configuration par defaut Phase 13+++
// Aucune modification necessaire
// Fonctionne optimalement "out of the box"
```

### Cas 3: Resume Long (1000+ mots)
```go
// Assouplir crit�res
penalite += float64(count-2) * 0.05
if similarity > 0.7 { delete(...) }
if i-lastPos < 3 { continue }
if compteurMots[motClean] % 4 == 0 { synonymes }
// Augmenter limites bloc
if numBlocs > 80 { numBlocs = 80 }
```

### Cas 4: Texte Technique (vocabulaire specialise)
```go
// Synonymes peuvent perdre precision
// Desactiver ou utiliser dict technique
if compteurMots[motClean] % 10 == 0 { synonymes }  // Tr�s rare

// Augmenter penalite TF-IDF (mots techniques courants)
tfidfVal *= 0.7
```

### Cas 5: Texte Narratif (fiction, histoire)
```go
// Synonymes naturels et attendus
if compteurMots[motClean] % 2 == 0 { synonymes }

// Assouplir antirepetition (repetitions stylistiques acceptables)
if i-lastPos < 4 { continue }  // vs 5
```

---

##  Benchmarking & Tuning Iteratif

### Test 1: Mesurer Longueur
```bash
./programme resume input.txt 0.12 2>&1 | grep "Mots generes"
# Note la longueur
```

### Test 2: Mesurer Coherence
```bash
./programme resume input.txt 0.12 2>&1 | grep "Coherence moyenne"
# Note le pourcentage
```

### Test 3: Mesurer Rapidite
```bash
time ./programme resume input.txt 0.12 > /dev/null
# Note le temps d'execution
```

### Cycle d'Optimisation
```
1. Test configuration actuelle
2. Changer 1 param�tre � la fois
3. Re-tester (longueur, coherence, vitesse)
4. Garder si ameliore metrique cible
5. Repeter jusqu'� optimum
```

---

##  Diagnostic & Depannage

###  Probl�me: Trop de Repetitions
**Cause probable**: Distance anti-repetition trop courte  
**Solution**:
```go
if i-lastPos < 7 { continue }  // Augmenter de 5 � 7
```

###  Probl�me: Resume Trop Court
**Cause probable**: Fen�trage strict trop agressif  
**Solution**:
```go
if similarity > 0.7 { delete(...) }  // Augmenter de 0.6 � 0.7
```

###  Probl�me: Synonymes "Bizarre"
**Cause probable**: Dictionnaire inadapte au domaine  
**Solution**: �tendre `SynonymsDict` avec termes appropries
```go
"algorithme": {"methode", "procedure", "technique", "algorithme"},
"donnees":    {"informations", "elements", "contenus", "donnees"},
```

###  Probl�me: Coherence Baisse
**Cause probable**: Filtrage trop agressif  
**Solution**: Assouplir penalites
```go
penalite += float64(count-2) * 0.05  // Reduire de 0.1
tfidfVal *= 0.85  // Augmenter de 0.8
```

---

##  Comprendre l'Impact

### Les 5 Filtres en Cascade
```
Input Blocs (180)
    
[1] Normalisation Lexicale  Penalite blocs repetitifs
    
[2] TF-IDF Intelligent  Mots rares penalises
    
[3] Fen�trage Strict  Blocs consecutifs diversifies
     (Selection: 45 blocs)
[4] Generation
    
[5] Post-Traitement
     Anti-repetition <5 mots
     Synonymes contextuels
    
Output Texte (679 mots, 95% coherence)
```

**Si sortie n'est pas bonne**:
1. Identifier quel filtre contribue le plus
2. Ajuster ce filtre specifically
3. Laisser autres � defaults

---

##  Checklist Optimisation

- [ ] Tester configuration default Phase 13+++
- [ ] Mesurer: longueur, coherence, vitesse
- [ ] Identifier metrique cible (longueur? qualite? vitesse?)
- [ ] Ajuster param�tres selon profil (Quality/Balanced/Coverage)
- [ ] Re-tester apr�s chaque changement
- [ ] Documenter configuration finale
- [ ] Valider sur multi-textes d'entree

---

##  References

- [PHASE-13-PLUS-PLUS-PLUS.md](PHASE-13-PLUS-PLUS-PLUS.md) - Specifications techniques
- [PHASE-13-COMPARISON.md](PHASE-13-COMPARISON.md) - Comparaison avant/apr�s
- [Code Source: resumeur_coherence.go](database/resumeur_coherence.go#L554-L580)
- [Code Source: generation.go](database/generation.go#L495-L510)
- [Code Source: coherence.go](database/coherence.go#L410-L470)

---

**Derni�re mise � jour**: Phase 13+++  
**Status**:  Stable & Documented  
**Recommandation**: Utiliser profil "Balanced" par defaut, adapter selon besoins specifiques.
