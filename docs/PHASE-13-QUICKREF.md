# 🚀 Phase 13+++ Quick Reference

## One-Liner
**5 stratégies imbriquées** (normalisation lexicale, TF-IDF intelligent, fenêtrage strict, anti-répétition, synonymes) → **95% cohérence** sans répétitions, **86% plus rapide**.

---

## The 5 Strategies

| # | Nom | Fichier | Effet | Code |
|---|-----|---------|-------|------|
| 1️⃣ | Normalisation Lexicale | resumeur_coherence.go | Pénalité blocs répétitifs | `penalite = (count-2) × 0.1` |
| 2️⃣ | TF-IDF Intelligent | generation.go | Mots rares moins influents | `tfidf *= 0.8` si IDF>0.5 |
| 3️⃣ | Fenêtrage Strict | resumeur_coherence.go | Blocs consécutifs diversifiés | `skip si similarity > 0.6` |
| 4️⃣ | Anti-Répétition | coherence.go | Zéro répétition <5 mots | `skip si position_gap < 5` |
| 5️⃣ | Synonymes | coherence.go | Vocabulaire varié | `remplace 1/3 occurrences` |

---

## Key Metrics

```
Phase 13++   →  Phase 13+++
──────────────────────────
1297 mots   →  679-847 mots
50 blocs    →  45 blocs
94.83%      →  95.00% cohérence
1384ms      →  219ms (6.3x faster) ⚡
```

---

## Configuration Profiles

### 🎯 Maximum Quality (Zéro répétitions)
```go
penalite += count-2 * 0.2           // Fort
similarity > 0.4                     // Strict
tfidf *= 0.7                         // Très pénalisant
distance < 7                         // Long
% 2 synonymes                        // Fréquent
→ 400-500 mots, ultra-qualité
```

### ⚡ Balanced (Recommandé)
```go
penalite += count-2 * 0.1           // Standard
similarity > 0.6                     // Modéré
tfidf *= 0.8                         // Standard
distance < 5                         // Standard
% 3 synonymes                        // Modéré
→ 600-800 mots, qualité optimale
```

### 📖 Maximum Coverage (Plus contenu)
```go
penalite += count-2 * 0.05          // Faible
similarity > 0.75                    // Permissif
tfidf *= 0.9                         // Moins pénalisant
distance < 3                         // Court
% 5 synonymes                        // Discret
→ 1000-1200 mots, couverture max
```

---

## Tuning Checklist

- [ ] Test config default
- [ ] Mesurer: longueur, cohérence, vitesse
- [ ] Choisir métrique cible
- [ ] Sélectionner profil (Quality/Balanced/Coverage)
- [ ] Ajuster 1 param à la fois
- [ ] Re-tester après chaque change
- [ ] Valider sur multi-textes
- [ ] Documenter config finale

---

## Files Modified

```
database/resumeur_coherence.go  (+150 lines)
├─ RepetitionsBloc field
├─ NormaliserRepetitionsBlocs()
├─ Scoring formula update
├─ Fenêtrage strict
└─ CalculerSimilarityVocabLexical()

database/generation.go           (+15 lines)
└─ TF-IDF penalty 0.8x

database/coherence.go            (+80 lines)
├─ SynonymsDict
├─ Anti-répétition <5 mots
└─ Synonym replacement
```

---

## Test Results

```
✅ input.txt ratio 12%: 679 mots, 95.00%, 187.9ms
✅ input.txt ratio 15%: 847 mots, 95.00%, 219.1ms
✅ test.txt: 12 mots, 95.00%, 977µs
✅ Build: SUCCESS
✅ All filters: ACTIVE
```

---

## Common Issues & Fixes

| Issue | Cause | Fix |
|-------|-------|-----|
| Too many repetitions | Distance too short | Increase distance from 5 to 7 |
| Summary too short | Fenêtrage too strict | Increase similarity from 0.6 to 0.7 |
| Weird synonyms | Wrong dictionary | Extend SynonymsDict for domain |
| Coherence drops | Filtering too aggressive | Reduce penalty 0.1 → 0.05 |
| Slow execution | Too many blocks | Reduce max blocks from 50 to 30 |

---

## Performance Gains

```
Speed:  6.3x faster (219ms vs 1384ms)
Quality: 95% coherence maintained
Reps:   ~0 detections (was many)
Vocab:  Natural variation (was repetitive)
```

---

## Recommended Usage

```bash
# Default (Balanced, recommended)
./programme resume input.txt 0.12

# If need more quality
./programme resume input.txt 0.10  # Compress more → better curation

# If need more content
./programme resume input.txt 0.20  # Compress less → wider selection
```

---

## Documentation Files

- **[PHASE-13-PLUS-PLUS-PLUS.md](PHASE-13-PLUS-PLUS-PLUS.md)** - Technical specs
- **[PHASE-13-COMPARISON.md](PHASE-13-COMPARISON.md)** - Before/After
- **[PHASE-13-CONFIGURATION.md](PHASE-13-CONFIGURATION.md)** - Tuning guide
- **[PHASE-13-VALIDATION.md](PHASE-13-VALIDATION.md)** - Full validation

---

## Next Steps (Optional)

1. **Extend SynonymsDict**: 20 → 50+ entries
2. **Context-aware synonyms**: Vary by category
3. **Bigram checking**: Check word pairs too
4. **Lemmatization**: donné/donnent/donnée = same
5. **Incremental caching**: Pre-compute for recurrent texts

---

**Status**: ✅ PRODUCTION-READY  
**Recommendation**: Use default "Balanced" profile  
**Deployment**: Ready for production use
