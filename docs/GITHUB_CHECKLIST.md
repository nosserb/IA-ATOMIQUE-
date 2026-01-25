#  Checklist Deploiement GitHub

## � Avant de pusher

### Code & Compilation
- [x] `go build` compile sans erreur
- [x] `go mod tidy` execute
- [x] Pas de fichiers `.DS_Store` ou temporaires
- [x] `.gitignore` � jour

### Fichiers Web
- [x] `web/index.html` existe
- [x] `web/style.css` existe  
- [x] `web/script.js` existe
- [x] Port 8080 utilise par defaut

### Configuration
- [x] `go.mod` configure
- [x] `Dockerfile` inclus
- [x] `docker-compose.yml` inclus
- [x] `.github/workflows/build.yml` configure

### Documentation
- [x] `README.md` mis � jour avec section web
- [x] `INSTALL.md` cree avec instructions
- [x] `SETUP_WEB.md` cree
- [x] `WEB_README.md` cree
- [x] `Makefile` cree

### Scripts
- [x] `start-web.sh` executable

---

##  Commands pour verifier avant push

```bash
# Verifier la compilation
go build -o programme && echo " Build OK"

# Verifier les fichiers web
test -f web/index.html && test -f web/style.css && test -f web/script.js && echo " Web files OK"

# Verifier Git
git status
git add .
git commit -m "feat: Add web interface for text summarization"

# Verifier les fichiers qui vont �tre pushes
git diff --cached --name-status

# Push
git push origin chore/siteweb
```

---

## � Test avec Docker (optionnel)

```bash
# Build l'image
docker build -t ia-atomique .

# Lance le conteneur
docker run -p 8080:8080 ia-atomique

# Teste l'interface
# Ouvre http://localhost:8080
```

---

##  Fichiers importants � verifier

```
 README.md          - Section web ajoutee
 INSTALL.md         - Instructions de setup
 WEB_README.md      - Docs interface web
 SETUP_WEB.md       - Configuration avancee
 Dockerfile         - Deploiement Docker
 docker-compose.yml - Compose config
 .github/workflows/ - CI/CD
 Makefile           - Build automation
 .gitignore         - Fichiers ignores
 go.mod            - Dependances Go
 web/               - Interface web compl�te
```

---

##  Apr�s le push

1. GitHub Actions va tester la compilation automatiquement
2. L'interface web sera accessible une fois clonee
3. Les instructions INSTALL.md guideront les utilisateurs

---

##  Pour les utilisateurs qui clonent le repo

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

 **Pr�t pour GitHub!**
