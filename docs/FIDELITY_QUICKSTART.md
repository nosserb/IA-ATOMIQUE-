# GUIDE UTILISATEUR - MODULE ANTI-HALLUCINATION PHASE 15

##  Demarrage rapide

### 1. Compiler le projet
```bash
cd "/home/student/autre projets/IA-ATOMIQUE-"
go build -o programme
```

### 2. Tester les fonctionnalites

#### Test automatique complet
```bash
./programme fidelity test
```

**Resultat attendu** :
```
[TEST: Technique Simple]
  Source length: 28 words
  Generated length: 10 words
  Fidelity: 100.00%  

[TEST: Texte Encyclopedique]
  Source length: 29 words
  Generated length: 10 words
  Fidelity: 100.00%  
```

#### Analyser un fichier specifique
```bash
./programme fidelity file input.txt
```

**Sortie** :
1.  Termes cles extraits
2. Resume Phase 13+++
3. Score fidelite (Ff) calcule
4. Decision : G�N�RATIF ou EXTRACTIF
5. Resume final
6. Rapport sauvegarde

#### Comparer les strategies
```bash
./programme fidelity compare test_atomique_technique.txt
```

**Affiche** : Tableau comparatif des 3 strategies principales

#### Test hybridation avec seuil custom
```bash
./programme fidelity hybrid input.txt 0.85
```

---

##  Resultat attendus

### Cas 1 : Texte technique ( BONNE FID�LIT�)

```
Texte original (167 mots):
"Un atome computationnel est une unite autonome du reseau atomique T.R.A..."

[ANALYSE FID�LIT�]
Coverage (Ff): 38-42%
Mode selectionne: EXTRACTIF

Raison: Score < 80%  Basculer automatiquement sur extraction TF-IDF
```

### Cas 2 : Texte general ( HALLUCINATION D�TECT�E)

```
Texte original (8080 mots):
"La Mesange huppee (Lophophanes cristatus)..."

[ANALYSE FID�LIT�]
Coverage (Ff): 1.42%
Mode selectionne: EXTRACTIF

Raison: Termes techniques IA-ATOMIQUE ne correspondent pas
```

### Cas 3 : Texte bien aligne ( TEXTE G�N�R� ACCEPT�)

```
[ANALYSE FID�LIT�]
Coverage (Ff): 85%
Mode selectionne: G�N�RATIF (fid�le)

Raison: Score >= 80%  Texte genere est assez fid�le
```

---

##  Interpretation des scores

| Score Ff | Interpretation | Action |
|---|---|---|
| ** 90%** | Excellent |  Garder resume genere |
| **80-90%** | Bon |  Garder resume genere |
| **70-80%** | Acceptable |  Garder avec vigilance |
| **60-70%** | Faible |  Utiliser extractif |
| **< 60%** | Critique |  FORCER extractif |

---

##  Strategies utilisables

### Strategie A : Extraction pure (TF-IDF)
**Approche** : Selectionner les meilleures phrases du texte original

```go
extractedSummary := database.ExtractiveResume(sourceText, compressionRatio)
```

**Fidelite** : 100% (aucune hallucination possible)
**Avantage** : Garantie de coherence
**Inconvenient** : Moins naturel que la generation

---

### Strategie B : Filtrage post-generation
**Approche** : Supprimer les mots inventes du resume genere

```go
filteredSummary := database.FilterForFidelity(generated, sourceVocab)
```

**Fidelite** : Depend du filtrage
**Avantage** : Garde la generation Phase 15 quand possible
**Inconvenient** : Peut creer des textes fragmentes

---

### Strategie C : Hybridation (RECOMMAND�E)
**Approche** : Utiliser generation si fid�le, sinon extractif

```go
finalSummary, fidelity, mode := database.HybridResume(
    generatedSummary,
    sourceText,
    0.80, // seuil fidelite
)

// mode = "G�N�RATIF (fid�le)" ou "EXTRACTIF (hallucination detectee)"
```

**Fidelite** : Garantie  80%
**Avantage** : Meilleur des deux mondes
**Inconvenient** : Aucun !

---

##  Mesures mathematiques

### Formule fidelite

$$F_f(R,T) = \frac{|\text{mots du resume en commun avec source}|}{|\text{total mots du resume}|}$$

**Exemple** :
- Resume genere: "Le reseau converge via resonance"
- Source contient: { reseau, converge, resonance, ... }
- Fidelite = 3/5 = 60%

### Seuil hybride

- **Par defaut** : � = 0.80 (80%)
- **Ajustable** : Augmentez � 0.85-0.90 pour domaines critiques

---

##  Depannage

### "Coverage tr�s bas (< 10%)"

**Cause probable** : Vocabulaire source incompatible avec termes techniques du projet

**Solution** : 
1. Verifier que texte source utilise terminologie IA-ATOMIQUE
2. Enrichir `database/fidelity_check.go` avec nouveaux termes techniques

### "Resume fragmente apr�s filtrage"

**Cause** : Trop de mots rejetes par strategie B

**Solution** : Utiliser strategie C (hybridation) � la place

### "Mode EXTRACTIF quand je veux generation"

**Cause** : Fidelite < seuil

**Solution** : Soit
1. Rel�cher le seuil : `database.HybridResume(..., 0.70)`
2. Soit enrichir le vocabulaire source pour meilleure couverture

---

##  Checklist de validation

- [ ] Compiler le projet : `go build -o programme`
- [ ] Tester basique : `./programme fidelity test`
- [ ] Tester sur fichier : `./programme fidelity file test_atomique_technique.txt`
- [ ] Verifier decision hybride : Mode = EXTRACTIF ou G�N�RATIF ?
- [ ] Consulter rapport : `test_atomique_technique_fidelity_report_*.txt`

---

##  Fichiers cles

| Fichier | Description |
|---|---|
| `database/fidelity_check.go` | Implementation du scoring Ff et strategies |
| `fidelity_commands.go` | CLI pour tester les strategies |
| `PHASE-15-ANTI-HALLUCINATION.md` | Documentation mathematique compl�te |

---

**Statut** :  Operationnel  
**Derni�re mise � jour** : 8 janvier 2026  
**Contact** : IA-ATOMIQUE Project Team
