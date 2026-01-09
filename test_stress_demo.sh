#!/bin/bash

echo "╔════════════════════════════════════════════════════════════════╗"
echo "║  STRESS TEST DEMONSTRATION - Amdahl & Batching Optimization   ║"
echo "╚════════════════════════════════════════════════════════════════╝"
echo ""

# Test avec différentes tailles
for ops in 100000 500000 1000000; do
    echo ""
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo "🔬 TEST: $ops opérations arithmétiques (big.Int 18-20 chiffres)"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    
    ./programme stest $ops 2>&1 | grep -A 50 "PHASE 1"
done

echo ""
echo "✅ Démonstration complète terminée"
