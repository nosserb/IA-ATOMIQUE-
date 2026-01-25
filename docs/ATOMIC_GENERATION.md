# IA-ATOMIQUE v5.0 - Generation par Resonance Atomique Autonome

##  Architecture

### Trois methodes de generation coexistent:

#### 1. **Generation Vectorielle (Originale)**
- Selection greedy: `rk = argmax[sim(vw, VR(k)) + ��aw(k-1)]`
- Vecteur cible VR = Σ E(Pi)�vi / Σ E(Pi)
- Avantages: Tr�s coherent (91%), rapide (237µs), deterministe
- Parfait pour: Resumes serres, textes formels

#### 2. **Generation Atomique Autonome (NEW!)**
- Resonance distribuee: `aj(t+1) = f(aj(t), Σ �jk�ak(t), betaj�VR)`
- Chaque atome = idee autonome avec activation propre
- �mission quand activation > �_wake (seuil de resonance)
- Avantages: Exploratif, moins repetitif, naturel
- Parfait pour: Textes creatifs, exploration semantique

#### 3. **Generation Multi-reseaux avec Resonance Globale**
- Divise le texte en sous-reseaux (~10-20 phrases chacun)
- Chaque sous-reseau = resonance locale autonome
- Resonance inter-reseaux = influence vectorielle globale
- Avantages: Scalable pour gros textes, econome en memoire

---

##  Formules Mathematiques

### Mise � jour atomique
```
aj(t+1) = �(aj(t) + Σ_{kvoisins} �jk�ak(t) + betaj�VR)
```
- � = sigmo�de (activation)
- �jk  [0.2, 1.0] = couplage semblant aux voisins
- betaj = influence du vecteur cible (beta  0.75)

### �mission d'un mot
```
Si aj(t+1)  �_wake  alors  rj = concept(j)
```
- �_wake  0.70 = seuil de resonance
- Apr�s emission: aj  0.2 (reset) et EtatFreeze=true

### �nergie et freeze
```
Energie(j) += Activation(j)
Si Activation(j) < FreezeThreshold  alors  EtatFreeze(j)=true
```
- Reduit la consommation d'energie pour atomes inactifs
- Reveille progressivement via resonance des voisins

### Vecteur cible d'un sous-reseau
```
VR_chunk = Σ_{i=1}^N E(Pi)�Ai / Σ_{i=1}^N E(Pi)
```

---

##  Commandes

### Generer avec methode vectorielle (existant)
```bash
./programme generate <fichier> [compression=0.3]
```

### Generer avec resonance atomique autonome (NEW!)
```bash
./programme atomic <fichier> [compression=0.3]
```

### Comparer les deux methodes
```bash
./programme compare <fichier> [compression=0.3]
```

---

##  Resultats Comparatifs (test.txt)

| Metrique | Vectoriel | Atomique |
|----------|-----------|----------|
| **Coherence** | 91.47% | 9.19% |
| **Diversite** | 100% | 100% |
| **Qualite Globale** | 95.73% | 54.60% |
| **Temps** | 237µs | 2.0ms |
| **Mots generes** | 20 | 25 |
| **Nature** | Deterministe | Exploratoire |

---

##  Quand utiliser quelle methode?

### Utiliser **VECTORIEL** (generate):
-  Resumes de documents techniques
-  Textes � forte coherence requise
-  Performance critique (temps reel)
-  Contexte formel/professionnel

### Utiliser **ATOMIQUE** (atomic):
-  Generation creative, exploration semantique
-  Textes avec variete thematique
-  Apprentissage des patterns distribues
-  Textes longs (multi-reseaux + resonance globale)

### Utiliser **COMPARAISON** (compare):
-  Tester differentes strategies
-  Analyser le comportement du syst�me
-  Benchmark performance

---

##  Insight Theorique

La generation atomique implemente le concept fondamental de T.R.A. (Technologie de Resonance Atomique):

> **L'intelligence emerge de l'interaction locale entre entites autonomes, pas d'une computation centrale.**

Chaque atome:
- Maintient son propre etat d'activation
- N'interagit qu'avec ses voisins (couplage �jk)
- Est influence par l'idee centrale (betaj�VR)
- �met quand la resonance atteint �_wake

Le texte emerge comme un motif stable dans le reseau, plutot que d'�tre construit word-by-word de mani�re greedy.

---

##  Param�tres Configurables

```go
type ReseauAtomiquePourGeneration struct {
    Couplage         float64 // � = 0.25 (interaction locale)
    InfluenceCible   float64 // beta = 0.75 (pertinence globale)
    SeuilResonance   float64 // �_wake = 0.70 (seuil d'emission)
    TauxDecroissance float64 // = 0.15 (oubli)
    FreezeThreshold  float64 // = 0.20 (gel automatique)
}
```

---

##  Integration avec le syst�me existant

-  Compatible avec vectorisation 11D existante
-  Utilise le lexique et categories du syst�me
-  Respecte la hierarchie des phrases (energie E(Pi))
-  Peut combiner avec translation/humanization

---

##  Exemple d'utilisation

```bash
# Generer resume 30% avec resonance atomique
./programme atomic document.txt 0.3

# Comparer vectoriel vs atomique sur m�me fichier
./programme compare document.txt 0.2

# Puis humaniser le resume atomique
./programme rewrite "resume genere" professionnel
```

---

**Status**:  Implemente et fonctionnel
**Performance**: Optimise (2ms pour 15 phrases)
**Scalabilite**: Multi-reseaux pour 1000+ phrases
