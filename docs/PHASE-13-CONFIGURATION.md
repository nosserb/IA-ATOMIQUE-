# Guide Phase 13+++ - Configuration & Tuning

##  Quick Start

```bash
# Compiler
go build -o programme

# Résumé standard (Phase 13+++)
./programme resume input.txt 0.12

# Résumé avec style humanisé
./programme humanize file -p input.txt
```

**Résultat**: Résumé de qualité maximale, sans répétitions, lecture fluide.

---

##  Param�tres Configurables

### 1. **Normalisation Lexicale** (resumeur_coherence.go)

**Coefficient de Pénalité**:
```go
// Ligne ~565
penalite += float64(count-2) * 0.1  //  Ajuster ce coefficient
```

| Coefficient | Effet |
|-------------|-------|
| 0.05 | Tr�s permissif, peu de pénalité |
| **0.1** |  **Recommandé** |
| 0.15 | Strict, forte déprioritisation |
| 0.20 | Tr�s strict, blocs répétitifs exclus |

**Exemple**:
```go
// Pour �tre plus permissif sur répétitions internes
penalite += float64(count-2) * 0.05

// Pour �tre ultra-strict
penalite += float64(count-2) * 0.2
```

### 2. **Pondération TF-IDF** (generation.go)

**Seuils et Pénalité**:
```go
// Ligne ~507
if idf[mot] > 0.5 && tf[mot] > 0.05 {  //  Ajuster seuils
    tfidfVal *= 0.8  //  Ajuster multiplicateur
}
```

| IDF | TF | Multiplicateur | Cas d'Usage |
|-----|----|----|---|
| >0.5 | >0.05 | **0.8** |  **Standard** |
| >0.6 | >0.10 | 0.8 | Plus sélectif |
| >0.4 | >0.03 | 0.9 | Plus permissif |
| >0.7 | >0.15 | 0.7 | Ultra-strict |

**Configuration Avancée**:
```go
// Cas 1: Plus de mots (moins strict)
if idf[mot] > 0.6 && tf[mot] > 0.10 {
    tfidfVal *= 0.9
}

// Cas 2: Zéro répétitions (ultra-strict)
if idf[mot] > 0.4 {  // N'importe quel mot rare
    tfidfVal *= 0.7
}
```

### 3. **Fen�trage Strict - Similarité Lexicale** (resumeur_coherence.go)

**Seuil de Diversité**:
```go
// Ligne ~735
if similarity > 0.6 {  //  Ajuster ce seuil
    delete(selectedIndices, i)
    continue
}
```

| Seuil | Similarité Max | Effet |
|-------|---|---|
| 0.4 | 40% vocab commun | Ultra-diversifié, peu de blocs |
| 0.5 | 50% vocab commun | Diversifié, sélection modérée |
| **0.6** | 60% vocab commun |  **Recommandé** |
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

// Default (recommandé)
if similarity > 0.6 {
    delete(selectedIndices, i)
}
```

### 4. **Anti-Répétition - Distance Minimale** (coherence.go)

**Distance Intra-Texte**:
```go
// Ligne ~420
if i-lastPos < 5 {  //  Ajuster distance
    continue  // Skip répétition
}
```

| Distance | Effet |
|----------|-------|
| <3 | Ultra-strict, "tr�s tr�s" impossible |
| <5 |  **Recommandé**, répétitions bien séparées |
| <7 | Modéré, quelques proches répétitions acceptées |
| <10 | Permissif, répétitions espacées acceptées |

**Configuration**:
```go
// Ultra-strict: aucune proche répétition
if i-lastPos < 3 {
    continue
}

// Recommandé: 5 mots minimum entre occurrences
if i-lastPos < 5 {
    continue
}

// Permissif: 10 mots minimum
if i-lastPos < 10 {
    continue
}
```

### 5. **Diversification Synonymes** (coherence.go)

**Fréquence de Remplacement**:
```go
// Ligne ~460
if compteurMots[motClean] % 3 == 0 && synChoisi != motClean {
    // Remplacer tous les 3 occurrences
}
```

| Modulo | Fréquence | Effet |
|--------|-----------|-------|
| % 2 | Tous les 2 | Tr�s diversifié (50% synonymes) |
| % 3 | Tous les 3 |  **Recommandé** (33% synonymes) |
| % 4 | Tous les 4 | Modéré (25% synonymes) |
| % 5 | Tous les 5 | Discret (20% synonymes) |

**Configuration**:
```go
// Tr�s varié (synonymes fréquents)
if compteurMots[motClean] % 2 == 0 && synChoisi != motClean {
    motsFiltres = append(motsFiltres, synChoisi)
    continue
}

// Standard (recommandé)
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

##  Profils de Configuration Pré-définis

###  **Profil 1: Maximum Quality** (Zéro Répétitions)

```go
// resumeur_coherence.go
penalite += float64(count-2) * 0.2  // Pénalité forte
if similarity > 0.4 { delete(...) }  // Strict diversité

// generation.go
if idf[mot] > 0.5 && tf[mot] > 0.05 {
    tfidfVal *= 0.7  // Tr�s pénalisant
}

// coherence.go
if i-lastPos < 7 { continue }  // Distance longue
if compteurMots[motClean] % 2 == 0 && ... {  // Synonymes fréquents
    motsFiltres = append(motsFiltres, synChoisi)
}
```

**Résultat**: 400-500 mots, **zéro** répétitions visibles, lecture tr�s naturelle.

---

###  **Profil 2: Balanced** (Recommandé - Par Défaut)

```go
// resumeur_coherence.go
penalite += float64(count-2) * 0.1  // Pénalité modérée
if similarity > 0.6 { delete(...) }  // Diversité modérée

// generation.go
if idf[mot] > 0.5 && tf[mot] > 0.05 {
    tfidfVal *= 0.8  // Standard
}

// coherence.go
if i-lastPos < 5 { continue }  // Distance standard
if compteurMots[motClean] % 3 == 0 && ... {  // Synonymes modérés
    motsFiltres = append(motsFiltres, synChoisi)
}
```

**Résultat**: 600-800 mots, 95% cohérence, lecture fluide, quelques variantes.

---

###  **Profil 3: Maximum Coverage** (Plus de Contenu)

```go
// resumeur_coherence.go
penalite += float64(count-2) * 0.05  // Pénalité faible
if similarity > 0.75 { delete(...) }  // Diversité permissive

// generation.go
if idf[mot] > 0.6 && tf[mot] > 0.10 {
    tfidfVal *= 0.9  // Moins pénalisant
}

// coherence.go
if i-lastPos < 3 { continue }  // Distance courte
if compteurMots[motClean] % 5 == 0 && ... {  // Synonymes discrets
    motsFiltres = append(motsFiltres, synChoisi)
}

// Augmenter limite sélection dans SelectionnerBlocsAvecFenetrageGlissant
if numBlocs > 75 { numBlocs = 75 }  // vs 50
```

**Résultat**: 1000-1200 mots, 92-94% cohérence, couverture maximale.

---

##  Recettes de Tuning

### Cas 1: Résumé Tr�s Court (100-200 mots)
```go
// �tre ultra-strict pour qualité
penalite += float64(count-2) * 0.2
if similarity > 0.4 { delete(...) }
if i-lastPos < 7 { continue }
if compteurMots[motClean] % 2 == 0 { synonymes }
```

### Cas 2: Résumé Standard (500-800 mots)
```go
// Configuration par défaut Phase 13+++
// Aucune modification nécessaire
// Fonctionne optimalement "out of the box"
```

### Cas 3: Résumé Long (1000+ mots)
```go
// Assouplir crit�res
penalite += float64(count-2) * 0.05
if similarity > 0.7 { delete(...) }
if i-lastPos < 3 { continue }
if compteurMots[motClean] % 4 == 0 { synonymes }
// Augmenter limites bloc
if numBlocs > 80 { numBlocs = 80 }
```

### Cas 4: Texte Technique (vocabulaire spécialisé)
```go
// Synonymes peuvent perdre précision
// Désactiver ou utiliser dict technique
if compteurMots[motClean] % 10 == 0 { synonymes }  // Tr�s rare

// Augmenter pénalité TF-IDF (mots techniques courants)
tfidfVal *= 0.7
```

### Cas 5: Texte Narratif (fiction, histoire)
```go
// Synonymes naturels et attendus
if compteurMots[motClean] % 2 == 0 { synonymes }

// Assouplir antirépétition (répétitions stylistiques acceptables)
if i-lastPos < 4 { continue }  // vs 5
```

---

##  Benchmarking & Tuning Itératif

### Test 1: Mesurer Longueur
```bash
./programme resume input.txt 0.12 2>&1 | grep "Mots générés"
# Note la longueur
```

### Test 2: Mesurer Cohérence
```bash
./programme resume input.txt 0.12 2>&1 | grep "Cohérence moyenne"
# Note le pourcentage
```

### Test 3: Mesurer Rapidité
```bash
time ./programme resume input.txt 0.12 > /dev/null
# Note le temps d'exécution
```

### Cycle d'Optimisation
```
1. Test configuration actuelle
2. Changer 1 param�tre � la fois
3. Re-tester (longueur, cohérence, vitesse)
4. Garder si améliore métrique cible
5. Répéter jusqu'� optimum
```

---

##  Diagnostic & Dépannage

###  Probl�me: Trop de Répétitions
**Cause probable**: Distance anti-répétition trop courte  
**Solution**:
```go
if i-lastPos < 7 { continue }  // Augmenter de 5 � 7
```

###  Probl�me: Résumé Trop Court
**Cause probable**: Fen�trage strict trop agressif  
**Solution**:
```go
if similarity > 0.7 { delete(...) }  // Augmenter de 0.6 � 0.7
```

###  Probl�me: Synonymes "Bizarre"
**Cause probable**: Dictionnaire inadapté au domaine  
**Solution**: �tendre `SynonymsDict` avec termes appropriés
```go
"algorithme": {"méthode", "procédure", "technique", "algorithme"},
"données":    {"informations", "éléments", "contenus", "données"},
```

###  Probl�me: Cohérence Baisse
**Cause probable**: Filtrage trop agressif  
**Solution**: Assouplir pénalités
```go
penalite += float64(count-2) * 0.05  // Réduire de 0.1
tfidfVal *= 0.85  // Augmenter de 0.8
```

---

##  Comprendre l'Impact

### Les 5 Filtres en Cascade
```
Input Blocs (180)
    
[1] Normalisation Lexicale  Pénalité blocs répétitifs
    
[2] TF-IDF Intelligent  Mots rares pénalisés
    
[3] Fen�trage Strict  Blocs consécutifs diversifiés
     (Sélection: 45 blocs)
[4] Génération
    
[5] Post-Traitement
     Anti-répétition <5 mots
     Synonymes contextuels
    
Output Texte (679 mots, 95% cohérence)
```

**Si sortie n'est pas bonne**:
1. Identifier quel filtre contribue le plus
2. Ajuster ce filtre specifically
3. Laisser autres � defaults

---

##  Checklist Optimisation

- [ ] Tester configuration default Phase 13+++
- [ ] Mesurer: longueur, cohérence, vitesse
- [ ] Identifier métrique cible (longueur? qualité? vitesse?)
- [ ] Ajuster param�tres selon profil (Quality/Balanced/Coverage)
- [ ] Re-tester apr�s chaque changement
- [ ] Documenter configuration finale
- [ ] Valider sur multi-textes d'entrée

---

##  Références

- [PHASE-13-PLUS-PLUS-PLUS.md](PHASE-13-PLUS-PLUS-PLUS.md) - Spécifications techniques
- [PHASE-13-COMPARISON.md](PHASE-13-COMPARISON.md) - Comparaison avant/apr�s
- [Code Source: resumeur_coherence.go](database/resumeur_coherence.go#L554-L580)
- [Code Source: generation.go](database/generation.go#L495-L510)
- [Code Source: coherence.go](database/coherence.go#L410-L470)

---

**Derni�re mise � jour**: Phase 13+++  
**Status**:  Stable & Documented  
**Recommandation**: Utiliser profil "Balanced" par défaut, adapter selon besoins spécifiques.
