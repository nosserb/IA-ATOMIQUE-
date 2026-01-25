# Checklist Déploiement GitHub

## à Avant de pusher

### Code & Compilation
- [x] `go build` compile sans erreur
- [x] `go mod tidy` exécuté
- [x] Pas de fichiers `.DS_Store` ou temporaires
- [x] `.gitignore` à jour

### Fichiers Web
- [x] `web/index.html` existe
- [x] `web/style.css` existe  
- [x] `web/script.js` existe
- [x] Port 8080 utilisé par défaut

### Configuration
- [x] `go.mod` configuré
- [x] `Dockerfile` inclus
- [x] `docker-compose.yml` inclus
- [x] `.github/workflows/build.yml` configuré

### Documentation
- [x] `README.md` mis à jour avec section web
- [x] `INSTALL.md` créé avec instructions
- [x] `SETUP_WEB.md` créé
- [x] `WEB_README.md` créé
- [x] `Makefile` créé

### Scripts
- [x] `start-web.sh` exécutable

---

## Commands pour vérifier avant push

```bash
# Vérifier la compilation
go build -o programme && echo " Build OK"

# Vérifier les fichiers web
test -f web/index.html && test -f web/style.css && test -f web/script.js && echo " Web files OK"

# Vérifier Git
git status
git add .
git commit -m "feat: Add web interface for text summarization"

# Vérifier les fichiers qui vont àtre pushés
git diff --cached --name-status

# Push
git push origin chore/siteweb
```

---

## à Test avec Docker (optionnel)

```bash
# Build l'image
docker build -t ia-atomique .

# Lance le conteneur
docker run -p 8080:8080 ia-atomique

# Teste l'interface
# Ouvre http://localhost:8080
```

---

## Fichiers importants à vérifier

```
 README.md          - Section web ajoutée
 INSTALL.md         - Instructions de setup
 WEB_README.md      - Docs interface web
 SETUP_WEB.md       - Configuration avancée
 Dockerfile         - Déploiement Docker
 docker-compose.yml - Compose config
 .github/workflows/ - CI/CD
 Makefile           - Build automation
 .gitignore         - Fichiers ignorés
 go.mod            - Dépendances Go
 web/               - Interface web complàte
```

---

## Apràs le push

1. GitHub Actions va tester la compilation automatiquement
2. L'interface web sera accessible une fois clonée
3. Les instructions INSTALL.md guideront les utilisateurs

---

## Pour les utilisateurs qui clonent le repo

Ils pourront faire:

```bash
git clone <ton-repo>
cd IA-ATOMIQUE
make web
# Ouvrir http://localhost:8080
```

Ou:

```bash
docker-compose up
# Ouvrir http://localhost:8080
```

Ou:

```bash
go build -o programme && ./programme web
# Ouvrir http://localhost:8080
```

---

 **Pràt pour GitHub!**
