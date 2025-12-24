
<div align="center" style="max-width: 900px; margin: 0 auto; padding: 0 20px;">

<div style="background-color: #fff3cd; border: 2px solid #ff9800; border-radius: 8px; padding: 20px; margin: 20px 0;">
  <p style="color: #ff6f00; font-weight: bold; margin: 0; font-size: 1.1em;">⚠️ VERSION BETA - v4.1</p>
  <p style="color: #333; margin: 10px 0 0 0;">Cette version peut contenir des bugs. Si vous rencontrez des problèmes ou avez des suggestions d'améliorations, merci de les signaler. Vos retours nous aident à corriger les problèmes et améliorer le système. <strong>Contact: signalez les problèmes ou suggestions</strong></p>
</div>

<p align="center">
  <a href="https://i.ibb.co/WN5frxd6/2ee3d1b24181.png">
    <img src="https://i.ibb.co/LzGHGkWG/8e199e30c114.png" alt="Comment ça marche">
  </a>
</p>

<p style="font-size: 1.2em; color: #666666b4; margin-bottom: 30px;"><strong>Un réseau neuronal sophistiqué et puissant pour l'analyse de texte, l'apprentissage et la génération de résumés</strong></p>

<div style="margin: 25px 0;">
  <a href="https://golang.org/"><img src="https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat-square&logo=go" alt="Go"></a>
  <span style="margin: 0 10px;"></span>
  <a href="#licence"><img src="https://img.shields.io/badge/Licence-MIT-4CAF50?style=flat-square" alt="Licence"></a>
  <span style="margin: 0 10px;"></span>
  <img src="https://img.shields.io/badge/Sécurité-Chiffrement--AES--256-blue?style=flat-square" alt="Sécurité">
</div>

<hr style="margin: 40px 0; border: none; border-top: 2px solid #eee;">

</div>

<div align="center" style="max-width: 900px; margin: 0 auto; padding: 0 20px;">

<div style="text-align: left; margin: 30px 0;">

### Vue d'ensemble

**IA-ATOMIQUE** est un système de réseau neuronal sophistiqué conçu pour l'analyse avancée de texte, l'apprentissage du langage et la génération intelligente de résumés. Il analyse le texte au niveau des phrases, classe le contenu en 6 catégories spécialisées, détecte les structures grammaticales et apprend des entrées tout en maintenant une modération éthique du contenu.

**Avantage clé :** 30 fois plus rapide que les LLM locaux équivalents tout en consommant des ressources minimales.

</div>

<hr style="margin: 40px 0; border: none; border-top: 2px solid #eee;">

<div style="text-align: left; margin: 30px 0;">

### ✨ Fonctionnalités

- **Analyse par Réseau Neuronal** - 1000 neurones sophistiqués avec analyse au niveau des phrases
- **Classification en 6 Catégories**:
  - 🔧 **TECH** - Technologie, programmation, systèmes numériques
  - 📚 **HISTOIRE** - Histoire, politique, événements actuels
  - 💼 **BUSINESS** - Commerce, économie, affaires
  - 🍽️ **ALIMENTATION** - Nourriture, nutrition, gastronomie
  - 🏥 **SANTE** - Santé, médecine, bien-être
  - 🔤 **VERBE** - Détection d'actions (verbes principaux)

- **Détection de la Structure Grammaticale** - Analyse sujet-verbe-complément
- **Apprentissage Dynamique** - Système de probation à deux étapes
- **Résumés Intelligents** - Génération de synthèse contextuelle
- **Modération de Contenu** - Liste noire chiffrée en AES-256
- **Tableau de Bord Temps Réel** - Visualisation et monitoring du réseau neuronal
- **Performance 30x** - Avantage massif de vitesse par rapport aux LLM locaux

</div>

<hr style="margin: 40px 0; border: none; border-top: 2px solid #eee;">

<div style="text-align: left; margin: 30px 0;">

### 🚀 Démarrage Rapide

#### Prérequis

- **Go** 1.22 ou supérieur
- **Git** pour le contrôle de version

#### Installation

```bash
# Cloner le dépôt
git clone <repository-url>
cd "IA-ATOMIQUE-"

# Télécharger les dépendances (déjà configurées)
go mod download

# Compiler l'exécutable
go build -o programme
```

#### Premier Lancement

```bash
# Afficher l'aide
./programme help

# Mode interactif
./programme

# Analyser un fichier
./programme file wikipedia.txt

# Analyser du texte en ligne
./programme text "L'intelligence artificielle progresse rapidement"
```

</div>

<hr style="margin: 40px 0; border: none; border-top: 2px solid #eee;">

<div style="text-align: left; margin: 30px 0;">

### 📖 Guide des Commandes

#### Obtenir de l'aide

Affiche le guide complet avec des exemples:

```bash
./programme help
# ou
./programme -h
./programme --help
```

---

#### Mode 1: Mode Interactif (Par défaut)

Parfait pour explorer le système de façon interactive:

```bash
# Interactif par défaut
./programme

# Interactif explicite
./programme interactive
./programme inter
```

**Ce qui se passe:**
- Demande une entrée
- Analyse chaque phrase
- Affiche la classification et la confiance
- Apprend de votre entrée
- Génère des résumés

**Exemple de session:**
```
[IA-ATOMIQUE v4.1] Mode Interactif
Entrez du texte (ou 'quit' pour arrêter):
> Einstein a découvert la théorie de la relativité
[HISTOIRE] - Confiance: 87% - Énergie: 5.23
Résumé: Einstein découverte théorie relativité
```

---

#### Mode 2: Analyse de Fichier

Analyser des fichiers entiers (documents, articles, rapports):

```bash
./programme file <chemin_fichier>
```

**Exemples:**

```bash
# Analyser un article Wikipedia
./programme file wikipedia.txt

# Analyser un document Markdown
./programme file document.md

# Chemin absolu
./programme file /home/utilisateur/documents/article.txt
```

**La sortie comprend:**
- Statistiques globales (phrases analysées, énergie totale, confiance)
- Distribution des catégories avec pourcentages et graphiques
- Analyse par phrase avec mots clés
- Verbes détectés au format `[verbe: ...]`
- Résumés synthétiques par catégorie

**Formats supportés:**
- Texte brut (.txt)
- Markdown (.md)
- HTML (.html)
- Tout format textuel

---

#### Mode 3: Analyse de Texte en Ligne

Analyse rapide de texte depuis la ligne de commande:

```bash
./programme text "<votre texte ici>"
```

**Exemples:**

```bash
# Une seule phrase
./programme text "Les neurones artificiels imitent le cerveau humain"

# Plusieurs phrases (entre guillemets)
./programme text "Python est un langage populaire. L'IA utilise le machine learning."

# Actualités
./programme text "Le changement climatique affecte la biodiversité"
```

**Résultat:**
- Classification immédiate
- Score d'énergie
- Termes clés identifiés
- Résumé synthétique

---

#### Mode 4: Mode Classique (Hérité)

Traitement simple de texte sans structure d'arguments:

```bash
./programme <texte>
./programme bonjour
./programme "l'histoire de l'informatique"
./programme python javascript golang
```

Ce mode traite tous les arguments comme une seule phrase à analyser.

</div>

<hr style="margin: 40px 0; border: none; border-top: 2px solid #eee;">

<div style="text-align: left; margin: 30px 0;">

### 📊 Explication de la Sortie

#### Statistiques Globales

```
[STATISTIQUES GLOBALES]
• Phrases analysées: 7
• Catégories détectées: 2
• Énergie totale: 15.72
• Confiance moyenne: 62.7%
```

**Métriques:**
- **Phrases analysées** - Nombre total de phrases traitées
- **Catégories détectées** - Nombre de catégories représentées
- **Énergie totale** - Énergie cumulée d'activation des neurones
- **Confiance moyenne** - Moyenne des confiances sur toutes les classifications

#### Distribution des Catégories

```
[DISTRIBUTION]
  HISTOIRE     █████████████████ 86% (6 phrases)
  TECH         ██ 14% (1 phrases)
```

Graphique visuel montrant:
- Nom de la catégorie
- Pourcentage de phrases dans cette catégorie
- Nombre de phrases
- Proportion visuelle

#### Analyse par Phrase

```
[TECH] - 1 phrases (énergie: 2.50)
Mots clés: son | formalisme | mathématique | efficacité

Phrases clés:
  • Si son formalisme mathématique est d'une efficacité inégalée...

Résumé: Formalisme mathématique efficacité
```

Affiche:
- Catégorie et nombre de phrases
- Énergie totale pour cette catégorie
- Termes clés (mots clés)
- Phrases originales
- Résumé synthétique

#### Détection des Verbes

Les verbes sont détectés séparément:

```
[verbe: découvert]
[verbe: analyse]
```

Cela empêche les verbes de polluer la classification du contenu.

</div>

<hr style="margin: 40px 0; border: none; border-top: 2px solid #eee;">

<div style="text-align: left; margin: 30px 0;">

### 🧠 Comment ça Marche

#### 1. Analyse au Niveau des Phrases

Le texte est divisé en phrases individuelles (pas juste des mots), permettant une catégorisation indépendante sans qu'une catégorie domine les autres.

#### 2. Classification Neuronale

Chaque phrase active les neurones pertinents dans les 6 catégories basé sur les correspondances de mots clés et les motifs appris. Les scores d'énergie indiquent la confiance.

#### 3. Système d'Apprentissage

**Probation à deux étapes:**
- Étape 1: Nouveau mot observé, ajouté à la liste d'apprentissage temporaire
- Étape 2: Mot vu 2+ fois, confirmé dans le lexique avec catégorie permanente
- Prévient les mauvaises classifications d'occurrences uniques

#### 4. Modération de Contenu

Le système maintient une liste noire chiffrée (AES-256-GCM) de 130+ mots offensants:
- Prévient l'apprentissage de contenu inapproprié
- Bloque ces mots lors de l'analyse
- Affiche `[BLOQUÉ]` lors de la rencontre de mots restreints
- Le chiffrement assure l'intégrité et l'alignement éthique du système

#### 5. Génération de Résumés

La synthèse contextuelle extrait les termes clés tout en préservant le sens:
- Supprime les mots vides (le, la, de, etc.)
- Maintient les noms et verbes importants
- Ordonne par importance
- Crée des résumés lisibles et compacts

</div>

<hr style="margin: 40px 0; border: none; border-top: 2px solid #eee;">

<div style="text-align: left; margin: 30px 0;">

### ⚡ Benchmarks de Performance

#### Métriques Clés

| Métrique | Résultat |
|----------|----------|
| **Débit** | 790 phrases/seconde |
| **Latence** | 1.26 ms/phrase |
| **Mémoire** | <20 MB RAM |
| **CPU Moy** | ~23% d'utilisation |
| **Taille Binaire** | 2.3 MB |

#### Tests de Stress

| Test | Résultat |
|------|----------|
| **1000 Phrases** | 31ms (32 107 phrases/sec pic) |
| **Stabilité** | <20% variation sur 5 exécutions |
| **Scalabilité** | Linéaire avec la taille du texte |

#### Comparaison vs LLM Local

| Métrique | IA-ATOMIQUE | LLM Local | Avantage |
|----------|-------------|-----------|----------|
| **Vitesse** | 1.8 ms | 56.3 ms | **30x plus rapide** |
| **Débit** | 5 475 phrases/s | 195 phrases/s | **28x plus rapide** |
| **Mémoire** | 12.5 MB | 4000 MB | **320x moins** |
| **CPU** | 19.8% | 94.2% | **4.7x plus efficace** |

</div>

<hr style="margin: 40px 0; border: none; border-top: 2px solid #eee;">

<div style="text-align: left; margin: 30px 0;">

### 🛡️ Sécurité & Éthique

#### Liste Noire Chiffrée

La liste de modération de contenu est sécurisée avec:
- **Algorithme**: AES-256-GCM (Advanced Encryption Standard)
- **Fichier**: `blacklist.enc` (chiffré, illisible)
- **Protection**: Prévient toute modification non autorisée
- **Runtime**: Déchiffré de façon transparente lors de l'exécution
- **Intégrité**: Application des valeurs éthiques au niveau du système

#### Pourquoi le Chiffrement?

> "Si quelqu'un supprime des mots ici, cela ne sera plus en accord avec les valeurs lors de la publication sur GitHub par nosserb."

Le chiffrement assure que le système de modération de contenu reste intègre et fiable.

#### Ce qui est Bloqué

La liste noire prévient l'apprentissage de:
- Langage obscène
- Discours haineux
- Insultes offensantes
- Termes discriminatoires
- Contenu explicite

Cela assure qu'IA-ATOMIQUE maintient les standards éthiques alignés avec les politiques communautaires de GitHub.

</div>

<hr style="margin: 40px 0; border: none; border-top: 2px solid #eee;">

<div style="text-align: left; margin: 30px 0;">

### 📁 Structure du Projet

```
IA-ATOMIQUE-/
├── main.go                    Point d'entrée et routage des commandes
├── interaction.go             Mode interactif et traitement des fichiers
├── database/
│   ├── data.go               Gestion des neurones, lexique, chiffrement
│   └── phrase_analysis.go    Moteur d'analyse au niveau des phrases
├── go.mod                     Dépendances Go
├── go.sum                     Sommes de contrôle des dépendances
├── blacklist.enc              Modération chiffrée (AES-256)
├── dashboard                  Visualisation du réseau neuronal
├── programme                  Exécutable compilé
└── README.md                  Cette documentation
```

</div>

<hr style="margin: 40px 0; border: none; border-top: 2px solid #eee;">

<div style="text-align: left; margin: 30px 0;">

### 💡 Conseils d'Utilisation

#### Pour les Meilleurs Résultats

1. **Utilisez le français** - Optimisé pour les motifs du langage français
2. **Phrases naturelles** - L'analyse par phrase fonctionne mieux avec du texte grammaticalement correct
3. **Contenu varié** - Mélangez les sujets pour voir la distribution des catégories
4. **Apprentissage itératif** - Exécutez plusieurs fois avec des textes similaires pour renforcer les motifs
5. **Documents volumineux** - Les fichiers avec 100+ phrases montrent une meilleure significativité statistique

#### Conseils de Performance

```bash
# Analyse rapide d'un gros fichier
time ./programme file large_document.txt

# Profil avec plusieurs catégories
./programme file wikipedia.txt

# Compiler un binaire de version pour la vitesse
go build -ldflags="-s -w" -o programme
```

#### Débogage

Le système affiche:
- Messages `[BLACKLIST]` pour la modération
- `[BLOQUÉ]` quand des mots bloqués sont rencontrés
- Niveaux d'énergie et scores de confiance
- Distributions des catégories

</div>

<hr style="margin: 40px 0; border: none; border-top: 2px solid #eee;">

<div style="text-align: left; margin: 30px 0;">

### 📝 Exemples

#### Exemple 1: Analyser un Article Technologique

```bash
./programme text "Python est un langage de programmation puissant. L'IA utilise Python pour le machine learning. TensorFlow et PyTorch sont des frameworks populaires."
```

**Résultat Attendu:**
```
[STATISTIQUES GLOBALES]
• Phrases analysées: 3
• Catégories détectées: 1
• Énergie totale: 18.5
• Confiance moyenne: 94.2%

[DISTRIBUTION]
  TECH  ██████████████████ 100% (3 phrases)

[TECH] - 3 phrases (énergie: 18.5)
Mots clés: Python | programmation | IA | machine learning | TensorFlow

Résumé: Python langage programmation IA machine learning
```

---

#### Exemple 2: Contenu Mixte Historique & Affaires

```bash
./programme file wikipedia.txt
```

**Résultat Attendu:**
```
[DISTRIBUTION]
  HISTOIRE  ████████████████ 60% (9 phrases)
  BUSINESS  ██████ 40% (6 phrases)

[ANALYSE PAR PHRASE]
[HISTOIRE] Napoléon a transformé l'Europe...
[BUSINESS] Le commerce s'est développé...
```

---

#### Exemple 3: Session d'Apprentissage Interactif

```bash
./programme interactive
```

**Session:**
```
Entrez du texte:
> Einstein a développé la théorie de la relativité en 1905
[HISTOIRE] - Confiance: 89%
Verbe détecté: [verbe: développé]
Résumé: Einstein théorie relativité

Entrez du texte:
> quit
```

</div>

<hr style="margin: 40px 0; border: none; border-top: 2px solid #eee;">

<div style="text-align: left; margin: 30px 0;">

### 🔧 Développement

#### Compiler depuis la Source

```bash
go build -o programme
```

#### Exécuter les Tests (si applicable)

```bash
go test ./...
```

#### Nettoyage et Recompilation

```bash
go clean
go build -o programme
```

</div>

<hr style="margin: 40px 0; border: none; border-top: 2px solid #eee;">

<div style="text-align: left; margin: 30px 0;">

### 📄 Licence

Ce projet est distribué sous la licence **MIT**. Voir le fichier [LICENSE](LICENSE) pour plus de détails.

</div>

<hr style="margin: 40px 0; border: none; border-top: 2px solid #eee;">

<div style="margin: 40px 0;">
  <p><strong>nosserb | 2025</strong></p>
  <p style="color: #999; font-size: 0.95em;">⭐ Mettez une étoile à ce projet si vous le trouvez utile!</p>
</div>

</div>
