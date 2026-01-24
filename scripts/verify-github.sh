#!/bin/bash

# Script de vérification avant push GitHub

echo "╔════════════════════════════════════════╗"
echo "║  IA-ATOMIQUE - Pre-push Verification  ║"
echo "╚════════════════════════════════════════╝"
echo ""

ERRORS=0

# 1. Vérifier la compilation
echo "📦 Vérification de la compilation..."
if go build -o programme 2>&1 | grep -q "error"; then
    echo "  ❌ ERREUR: La compilation a échoué"
    ERRORS=$((ERRORS + 1))
else
    echo "  ✓ Compilation OK"
fi

# 2. Vérifier les fichiers web
echo ""
echo "🌐 Vérification des fichiers web..."
for file in web/index.html web/style.css web/script.js; do
    if [ -f "$file" ]; then
        echo "  ✓ $file existe"
    else
        echo "  ❌ ERREUR: $file manquant"
        ERRORS=$((ERRORS + 1))
    fi
done

# 3. Vérifier les fichiers de configuration
echo ""
echo "⚙️  Vérification des fichiers de configuration..."
for file in Dockerfile docker-compose.yml Makefile .gitignore go.mod; do
    if [ -f "$file" ]; then
        echo "  ✓ $file existe"
    else
        echo "  ❌ ERREUR: $file manquant"
        ERRORS=$((ERRORS + 1))
    fi
done

# 4. Vérifier les fichiers README
echo ""
echo "📖 Vérification de la documentation..."
for file in README.md INSTALL.md WEB_README.md SETUP_WEB.md; do
    if [ -f "$file" ]; then
        size=$(wc -c < "$file")
        if [ "$size" -gt 500 ]; then
            echo "  ✓ $file existe ($(($size / 1024))KB)"
        else
            echo "  ⚠️  AVERTISSEMENT: $file est très petit"
        fi
    else
        echo "  ❌ ERREUR: $file manquant"
        ERRORS=$((ERRORS + 1))
    fi
done

# 5. Vérifier le .gitignore
echo ""
echo "🚫 Vérification du .gitignore..."
if grep -q "^programme$" .gitignore; then
    echo "  ✓ Binaire 'programme' dans .gitignore"
else
    echo "  ❌ ERREUR: Binaire 'programme' pas dans .gitignore"
    ERRORS=$((ERRORS + 1))
fi

# 6. Vérifier les fichiers indésirables
echo ""
echo "🧹 Vérification des fichiers indésirables..."
UNWANTED=$(find . -name "programme" -o -name "*.log" -o -name "dashboard" -o -name "temp.txt" | grep -v ".git" | wc -l)
if [ "$UNWANTED" -gt 0 ]; then
    echo "  ⚠️  AVERTISSEMENT: $UNWANTED fichier(s) indésirable(s) détecté(s)"
    find . -name "programme" -o -name "*.log" -o -name "dashboard" -o -name "temp.txt" | grep -v ".git" | sed 's/^/    /'
else
    echo "  ✓ Pas de fichiers indésirables"
fi

# 7. Vérifier GitHub Actions
echo ""
echo "🤖 Vérification de GitHub Actions..."
if [ -f ".github/workflows/build.yml" ]; then
    echo "  ✓ Workflow CI/CD configuré"
else
    echo "  ⚠️  AVERTISSEMENT: Pas de workflow GitHub Actions"
fi

# Résumé
echo ""
echo "════════════════════════════════════════"
if [ $ERRORS -eq 0 ]; then
    echo "✅ Tout est OK! Prêt pour GitHub"
    exit 0
else
    echo "❌ $ERRORS erreur(s) détectée(s)"
    exit 1
fi
