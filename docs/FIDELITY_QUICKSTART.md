# GUIDE UTILISATEUR - MODULE ANTI-HALLUCINATION PHASE 15

## Démarrage rapide

### 1. Compiler le projet
```bash
cd "/home/student/autre projets/IA-ATOMIQUE-"
go build -o programme
```

### 2. Tester les fonctionnalités

#### Test automatique complet
```bash
./programme fidelity test
```

**Résultat attendu** :
```
[TEST: Technique Simple]
  Source length: 28 words
  Generated length: 10 words
  Fidelity: 100.00%  

[TEST: Texte Encyclopédique]
  Source length: 29 words
  Generated length: 10 words
  Fidelity: 100.00%  
```

#### Analyser un fichier spécifique
```bash
./programme fidelity file input.txt
```

**Sortie** :
1.  Termes clés extraits
2. Résumé Phase 13+++
3. Score fidélité (Ff) calculé
4. Décision : GπNπRATIF ou EXTRACTIF
5. Résumé final
6. Rapport sauvegardé

#### Comparer les stratégies
```bash
./programme fidelity compare test_atomique_technique.txt
```

**Affiche** : Tableau comparatif des 3 stratégies principales

#### Test hybridation avec seuil custom
```bash
./programme fidelity hybrid input.txt 0.85
```

---

## Résultat attendus

### Cas 1 : Texte technique ( BONNE FIDπLITπ)

```
Texte original (167 mots):
"Un atome computationnel est une unité autonome du réseau atomique T.R.A..."

[ANALYSE FIDπLITπ]
Coverage (Ff): 38-42%
Mode sélectionné: EXTRACTIF

Raison: Score < 80%  Basculer automatiquement sur extraction TF-IDF
```

### Cas 2 : Texte général ( HALLUCINATION DπTECTπE)

```
Texte original (8080 mots):
"La Mésange huppée (Lophophanes cristatus)..."

[ANALYSE FIDπLITπ]
Coverage (Ff): 1.42%
Mode sélectionné: EXTRACTIF

Raison: Termes techniques IA-ATOMIQUE ne correspondent pas
```

### Cas 3 : Texte bien aligné ( TEXTE GπNπRπ ACCEPTπ)

```
[ANALYSE FIDπLITπ]
Coverage (Ff): 85%
Mode sélectionné: GπNπRATIF (fidπle)

Raison: Score >= 80%  Texte généré est assez fidπle
```

---

## Interprétation des scores

| Score Ff | Interprétation | Action |
|---|---|---|
| ** 90%** | Excellent |  Garder résumé généré |
| **80-90%** | Bon |  Garder résumé généré |
| **70-80%** | Acceptable |  Garder avec vigilance |
| **60-70%** | Faible |  Utiliser extractif |
| **< 60%** | Critique |  FORCER extractif |

---

## Stratégies utilisables

### Stratégie A : Extraction pure (TF-IDF)
**Approche** : Sélectionner les meilleures phrases du texte original

```go
extractedSummary := database.ExtractiveResume(sourceText, compressionRatio)
```

**Fidélité** : 100% (aucune hallucination possible)
**Avantage** : Garantie de cohérence
**Inconvénient** : Moins naturel que la génération

---

### Stratégie B : Filtrage post-génération
**Approche** : Supprimer les mots inventés du résumé généré

```go
filteredSummary := database.FilterForFidelity(generated, sourceVocab)
```

**Fidélité** : Dépend du filtrage
**Avantage** : Garde la génération Phase 15 quand possible
**Inconvénient** : Peut créer des textes fragmentés

---

### Stratégie C : Hybridation (RECOMMANDπE)
**Approche** : Utiliser génération si fidπle, sinon extractif

```go
finalSummary, fidelity, mode := database.HybridResume(
    generatedSummary,
    sourceText,
    0.80, // seuil fidélité
)

// mode = "GπNπRATIF (fidπle)" ou "EXTRACTIF (hallucination détectée)"
```

**Fidélité** : Garantie  80%
**Avantage** : Meilleur des deux mondes
**Inconvénient** : Aucun !

---

## Mesures mathématiques

### Formule fidélité

$$F_f(R,T) = \frac{|\text{mots du résumé en commun avec source}|}{|\text{total mots du résumé}|}$$

**Exemple** :
- Résumé généré: "Le réseau converge via résonance"
- Source contient: { réseau, converge, résonance, ... }
- Fidélité = 3/5 = 60%

### Seuil hybride

- **Par défaut** : π = 0.80 (80%)
- **Ajustable** : Augmentez π 0.85-0.90 pour domaines critiques

---

## Dépannage

### "Coverage trπs bas (< 10%)"

**Cause probable** : Vocabulaire source incompatible avec termes techniques du projet

**Solution** : 
1. Vérifier que texte source utilise terminologie IA-ATOMIQUE
2. Enrichir `database/fidelity_check.go` avec nouveaux termes techniques

### "Résumé fragmenté aprπs filtrage"

**Cause** : Trop de mots rejetés par stratégie B

**Solution** : Utiliser stratégie C (hybridation) π la place

### "Mode EXTRACTIF quand je veux génération"

**Cause** : Fidélité < seuil

**Solution** : Soit
1. Relπcher le seuil : `database.HybridResume(..., 0.70)`
2. Soit enrichir le vocabulaire source pour meilleure couverture

---

## Checklist de validation

- [ ] Compiler le projet : `go build -o programme`
- [ ] Tester basique : `./programme fidelity test`
- [ ] Tester sur fichier : `./programme fidelity file test_atomique_technique.txt`
- [ ] Vérifier décision hybride : Mode = EXTRACTIF ou GπNπRATIF ?
- [ ] Consulter rapport : `test_atomique_technique_fidelity_report_*.txt`

---

## Fichiers clés

| Fichier | Description |
|---|---|
| `database/fidelity_check.go` | Implémentation du scoring Ff et stratégies |
| `fidelity_commands.go` | CLI pour tester les stratégies |
| `PHASE-15-ANTI-HALLUCINATION.md` | Documentation mathématique complπte |

---

**Statut** :  Opérationnel  
**Derniπre mise π jour** : 8 janvier 2026  
**Contact** : IA-ATOMIQUE Project Team
