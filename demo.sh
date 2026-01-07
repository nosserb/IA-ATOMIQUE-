#!/bin/bash
echo "╔═══════════════════════════════════════════════════════════╗"
echo "║     IA-ATOMIQUE v5.0 - Résumé Dual-Mode Demo             ║"
echo "╚═══════════════════════════════════════════════════════════╝"
echo ""

echo "📄 Fichier source: test.txt (103 mots)"
echo ""

echo "─────────────────────────────────────────────────────────────"
echo "1️⃣  GÉNÉRATION VECTORIELLE (argmax + VR)"
echo "─────────────────────────────────────────────────────────────"
./programme generate test.txt 0.25 2>&1 | grep -A 2 "║  RÉSUMÉ"

echo ""
echo "─────────────────────────────────────────────────────────────"
echo "2️⃣  GÉNÉRATION ATOMIQUE (résonance distribuée)"
echo "─────────────────────────────────────────────────────────────"
./programme atomic test.txt 0.25 2>&1 | grep -A 2 "║  RÉSUMÉ"

echo ""
echo "─────────────────────────────────────────────────────────────"
echo "✅ Demo complète - Deux paradigmes, une seule plateforme!"
echo "─────────────────────────────────────────────────────────────"
