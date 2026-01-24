#!/bin/bash

# Script pour faire un commit et push facile

echo "╔════════════════════════════════════════╗"
echo "║  IA-ATOMIQUE - Git Helper              ║"
echo "╚════════════════════════════════════════╝"
echo ""

# Vérifier la branche actuelle
CURRENT_BRANCH=$(git rev-parse --abbrev-ref HEAD)
echo "📍 Branche actuelle: $CURRENT_BRANCH"
echo ""

# Vérifier le statut
CHANGES=$(git status --porcelain | wc -l)
if [ $CHANGES -eq 0 ]; then
    echo "✅ Aucun changement à committer"
    exit 0
fi

echo "📝 Changements détectés: $CHANGES fichier(s)"
echo ""

# Afficher les changements
git status --short | head -10
echo ""

# Demander un message de commit
read -p "📝 Message de commit: " -r COMMIT_MSG

if [ -z "$COMMIT_MSG" ]; then
    echo "❌ Message vide, abandon"
    exit 1
fi

# Stage, commit et push
echo ""
echo "🔄 Traitement..."
git add .
git commit -m "$COMMIT_MSG"

# Demander si on veut pusher
echo ""
read -p "🚀 Pusher sur $CURRENT_BRANCH? (y/n) " -n 1 -r
echo ""

if [[ $REPLY =~ ^[Yy]$ ]]; then
    git push origin $CURRENT_BRANCH
    echo "✅ Push réussi!"
else
    echo "⏸️  Push annulé"
fi

echo ""
echo "📊 Statut final:"
git status --short
