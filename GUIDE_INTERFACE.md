# IA-ATOMIQUE - Interface Graphique

Interface graphique moderne pour le système de traitement du langage naturel IA-ATOMIQUE.

## 📋 Caractéristiques

- ✨ Interface moderne et responsive
- 📝 Analyse de texte en temps réel
- 📊 Statistiques du réseau neuronal
- 🧠 Entraînement du modèle
- ⚙️ Paramètres configurables
- 🎨 Design épuré et intuitif

## 🚀 Installation et Démarrage

### Prérequis

- Go 1.16+
- Un navigateur moderne (Chrome, Firefox, Safari, Edge)

### Démarrage Rapide

1. **Compile et lance le serveur:**

```bash
go run server.go main.go interaction.go
```

Le serveur démarre sur `http://localhost:8080`

2. **Accède à l'interface:**

Ouvre ton navigateur et va sur `http://localhost:8080/app.html`

## 📖 Guide d'Utilisation

### 1. Analyse de Texte

- Écris ou colle un texte dans la zone d'entrée
- Clique sur "Analyser"
- L'IA traite le texte et affiche :
  - Catégories détectées
  - Entités trouvées
  - Résumé de l'analyse

### 2. Statistiques

- Vois les stats en temps réel du réseau neuronal:
  - Neurones actifs
  - Énergie totale
  - Nombre de catégories
  - Phrases traitées
- Clique "Actualiser" pour rafraîchir

### 3. Entraînement

- Sélectionne un fichier texte
- Définis le nombre d'itérations
- Lance l'entraînement
- Le système apprend et améliore sa précision

### 4. Paramètres

- **Taux d'apprentissage**: Contrôle la vitesse d'apprentissage (0.01-1.0)
- **Seuil de confiance**: Confiance minimale pour accepter une prédiction (0-1)
- Les paramètres sont sauvegardés localement

## 🏗️ Architecture

```
├── app.html          # Interface graphique
├── app.js            # Logique côté client
├── server.go         # API serveur Go
├── main.go           # Noyau IA
└── database/         # Base de données NLP
    ├── data.go
    ├── language.go
    ├── nlp.go
    └── phrase_analysis.go
```

## 🔌 API REST

### `POST /api/analyze`
Analyse un texte

```bash
curl -X POST http://localhost:8080/api/analyze \
  -H "Content-Type: application/json" \
  -d '{"text":"Votre texte ici"}'
```

### `GET /api/stats`
Récupère les statistiques du réseau

```bash
curl http://localhost:8080/api/stats
```

### `POST /api/train`
Entraîne le modèle

```bash
curl -X POST http://localhost:8080/api/train \
  -F "file=@fichier.txt" \
  -F "epochs=10"
```

### `GET /api/health`
Vérifie l'état du serveur

```bash
curl http://localhost:8080/api/health
```

## 🎨 Interface

L'interface est construite avec:
- **HTML5** pour la structure
- **CSS3** avec gradients et animations
- **JavaScript Vanilla** pour l'interactivité
- Design responsive (mobile, tablette, desktop)

## 📱 Responsive Design

L'interface s'adapte à tous les appareils:
- Desktop: 2 colonnes
- Tablette: 1 colonne
- Mobile: Interface optimisée

## 🔐 Sécurité

- CORS activé pour la communication cross-origin
- Validation côté serveur
- Pas de dépendances externes (vanilla JS)

## 🐛 Troubleshooting

### Erreur de connexion au serveur
- Assure-toi que le serveur Go est lancé
- Vérifie que le port 8080 est disponible
- Regarde la console du navigateur pour les erreurs

### Les stats ne se chargent pas
- Attends quelques secondes
- Clique sur "Actualiser les Stats"
- Vérifie que l'API est accessible

### L'analyse prend du temps
- C'est normal pour les textes longs
- Attends la fin de l'analyse
- Tu peux analyser d'autres textes en parallèle

## 📝 Fichiers de Configuration

Les paramètres sont sauvegardés dans `localStorage`:
```javascript
{
  "learningRate": 0.1,
  "threshold": 0.5
}
```

## 🚀 Améliorations Futures

- [ ] Export des résultats (PDF, JSON)
- [ ] Historique d'analyse
- [ ] Graphiques statistiques avancés
- [ ] Mode sombre
- [ ] Support multilingue
- [ ] Intégration avec bases de données externes

## 📄 Licence

Voir LICENSE

## 👤 Support

Pour des questions ou des problèmes, ouvre une issue sur le dépôt.

---

**IA-ATOMIQUE v4.0** - Système de Traitement du Langage Naturel
