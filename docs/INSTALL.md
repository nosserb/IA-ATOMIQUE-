# 🚀 Installation & Utilisation Rapide

## Installation depuis GitHub

### Prérequis
- **Go 1.22+** ou **Docker**
- **Git**

### 1️⃣ Cloner le repo
```bash
git clone https://github.com/ton-user/IA-ATOMIQUE.git
cd IA-ATOMIQUE
```

### 2️⃣ Compiler et Lancer

#### Option A: Avec Go (recommandé)
```bash
# Compiler
go build -o programme

# Lancer l'interface web
./programme web

# Ou lancer le script
./start-web.sh
```

#### Option B: Avec Docker
```bash
# Construire l'image
docker build -t ia-atomique .

# Lancer le conteneur
docker run -p 8080:8080 ia-atomique
```

#### Option C: Avec Docker Compose
```bash
docker-compose up
```

### 3️⃣ Accéder à l'interface
```
http://localhost:8080
```

## 📖 Modes d'utilisation

### Mode Web (Interface)
```bash
./programme web
# Accède à http://localhost:8080
```

### Mode Fichier
```bash
./programme file mon_texte.txt
```

### Mode Texte
```bash
./programme text "Votre texte à analyser"
```

### Mode Interactif
```bash
./programme interactive
```

## 🐛 Troubleshooting

### Port 8080 déjà utilisé
```bash
# Modifier le port dans web.go ligne 29:
# go StartWebServer("9000")
# Puis recompiler
go build -o programme
```

### Erreur de compilation
```bash
# Mettre à jour les dépendances
go mod tidy
go mod download
go build -o programme
```

### Fichiers web manquants
Vérifier que le répertoire `web/` existe avec:
- `index.html`
- `style.css`
- `script.js`

## 📊 Fichiers importants

```
├── web.go              # Serveur HTTP
├── main.go             # Entrée principale
├── interaction.go      # CLI
├── database/
│   ├── phrase_analysis.go
│   ├── nlp.go
│   ├── language.go
│   └── data.go
├── web/                # Interface web
│   ├── index.html
│   ├── style.css
│   └── script.js
├── Dockerfile          # Pour Docker
├── docker-compose.yml  # Compose config
└── go.mod             # Dépendances Go
```

## 🎯 Exemple de texte à tester

```
L'intelligence artificielle transforme le monde. 
Les réseaux neuronaux permettent aux machines d'apprendre.
La technologie progresse rapidement et change notre société.
```

Colle ce texte dans l'interface et clique "Résumer"!

---

📖 Pour plus d'infos: voir [WEB_README.md](WEB_README.md) et [SETUP_WEB.md](SETUP_WEB.md)
