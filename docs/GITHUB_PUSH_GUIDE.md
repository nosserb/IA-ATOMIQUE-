# π Guide: Créer et Pusher sur GitHub

## πtape 1: Créer le repo sur GitHub

1. Va sur [github.com](https://github.com)
2. Clique sur **"New"** (ou **"+"** en haut π droite)
3. Crée un nouveau repository:
   - **Repository name**: `IA-ATOMIQUE`
   - **Description**: `Neural network for text analysis and summarization with web interface`
   - **Public** ou **Private** (π toi de choisir)
   - NE coche PAS "Initialize with README" (tu en as déjπ un)
   - Clique **"Create repository"**

---

## πtape 2: Configuration Git Locale

### Si tu n'as pas encore configuré Git globalement:

```bash
git config --global user.name "Ton Nom"
git config --global user.email "ton.email@example.com"
```

### Dans le répertoire du projet:

```bash
cd "/home/student/autre projets/IA-ATOMIQUE-"

# Initialiser le repo si ce n'est pas déjπ fait
git init

# Ajouter l'URL du repo GitHub
# (remplace USERNAME par ton user GitHub et ia-atomique par le nom du repo)
git remote add origin https://github.com/USERNAME/ia-atomique.git

# Ou avec SSH (si tu as configuré les clés SSH):
# git remote add origin git@github.com:USERNAME/ia-atomique.git
```

---

## πtape 3: Ajouter les Fichiers

```bash
# Vérifier le statut
git status

# Ajouter tous les fichiers
git add .

# Vérifier ce qui sera commit
git status

# Commit initial
git commit -m "feat: Add web interface for text summarization

- Add modern web interface with text summarization
- Add HTTP server with /api/summarize endpoint
- Include Docker and docker-compose configuration
- Add comprehensive documentation and installation guide
- Setup GitHub Actions CI/CD
- Add Makefile for easy build and run commands"
```

---

## πtape 4: Configurer la Branche

### Option A: Push vers main directement

```bash
# Renommer la branche locale (si elle s'appelle master)
git branch -M main

# Push vers GitHub
git push -u origin main
```

### Option B: Créer une PR (recommandé)

Ton code est déjπ sur la branche `chore/siteweb`, donc:

```bash
# Vérifier la branche actuelle
git branch

# Si tu es pas sur chore/siteweb
git switch chore/siteweb

# Push la branche
git push -u origin chore/siteweb

# Puis va sur GitHub et crée une Pull Request
# de chore/siteweb vers main
```

---

## πtape 5: Vérification

Aprπs le push, va sur GitHub et vérifie:

 Tous les fichiers sont présents  
 README.md s'affiche bien  
 Les workflows GitHub Actions tournent  
 Les tests passent  

---

## Aprπs le Push

### Partager le Repo

```bash
# Afficher l'URL du repo
git remote -v

# Copier l'URL pour partager
# https://github.com/USERNAME/ia-atomique
```

### Pour les Utilisateurs qui Clonent

Ils vont faire:

```bash
git clone https://github.com/USERNAME/ia-atomique.git
cd ia-atomique
make web
# Ouvrir http://localhost:8080
```

---

## Configuration SSH (optionnel mais recommandé)

Pour éviter de taper le mot de passe π chaque fois:

```bash
# Générer une clé SSH
ssh-keygen -t ed25519 -C "ton.email@example.com"

# Copier la clé publique
cat ~/.ssh/id_ed25519.pub

# Va sur GitHub Settings > SSH and GPG keys
# Clique "New SSH key"
# Colle la clé et valide
```

Puis utilise l'URL SSH au lieu de HTTPS:

```bash
git remote set-url origin git@github.com:USERNAME/ia-atomique.git
```

---

## π Dépannage

### Erreur: "fatal: not a git repository"

```bash
# Tu dois πtre dans le répertoire du projet
cd "/home/student/autre projets/IA-ATOMIQUE-"

# Puis initialiser git
git init
```

### Erreur: "Permission denied (publickey)"

C'est un problπme SSH. Soit:
1. Utilise HTTPS π la place
2. Configure les clés SSH correctement

### Erreur: "failed to authenticate"

```bash
# Utilise HTTPS avec un Personal Access Token
git remote set-url origin https://YOUR_USERNAME:YOUR_TOKEN@github.com/USERNAME/ia-atomique.git
```

Génπre un token sur GitHub: Settings > Developer settings > Personal access tokens

---

## Fichiers Importants

```
 web/                    - Interface web
 Dockerfile              - Docker
 docker-compose.yml      - Compose
 README.md              - Documentation principale
 INSTALL.md             - Installation
 Makefile               - Build helpers
 .github/workflows/     - CI/CD automatisé
 go.mod                 - Dépendances
 *.go                   - Code source
```

Tout est prπt! 

---

## Message de Commit Détaillé

Si tu veux un message plus détaillé:

```bash
git commit -m "feat: Add web interface for text summarization

## Changes
- Add modern responsive web interface (HTML/CSS/JavaScript)
- Implement HTTP server with /api/summarize endpoint
- Include Docker and docker-compose configuration
- Setup GitHub Actions CI/CD pipeline
- Add comprehensive documentation:
  * INSTALL.md - Installation guide
  * WEB_README.md - Web interface documentation
  * SETUP_WEB.md - Advanced configuration
  * GITHUB_CHECKLIST.md - Pre-deployment checklist

## New Files
- web.go - Web server implementation
- web/ - Frontend files (HTML, CSS, JS)
- Dockerfile - Container configuration
- docker-compose.yml - Service orchestration
- .github/workflows/build.yml - CI/CD workflow
- Makefile - Build automation
- Various documentation files

## Features
- One-click deployment
- Docker support
- Modern UI with real-time analysis
- Fully documented and tested

Closes #XX (if you have an issue)"
```

---

**Tu es prπt! Bon push! **
