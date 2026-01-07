#!/bin/bash

# Script pour démarrer le serveur web IA-ATOMIQUE

echo "╔════════════════════════════════════════╗"
echo "║  IA-ATOMIQUE - Serveur Web             ║"
echo "║  Démarrage du serveur...               ║"
echo "╚════════════════════════════════════════╝"
echo ""

# Vérifier que le binaire existe
if [ ! -f "./programme" ]; then
    echo "[ERREUR] Le binaire 'programme' n'existe pas."
    echo "Compilation en cours..."
    go build -o programme
    if [ $? -ne 0 ]; then
        echo "[ERREUR] Échec de la compilation"
        exit 1
    fi
fi

# Lancer le serveur
echo "🚀 Démarrage du serveur sur http://localhost:8080"
echo "Appuyez sur Ctrl+C pour arrêter le serveur"
echo ""

./programme web
