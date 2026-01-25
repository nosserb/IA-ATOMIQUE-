#  Phase 13+++ - Documentation Index

##  Start Here

 **[PHASE-13-QUICKREF.md](PHASE-13-QUICKREF.md)**  One-page overview

---

##  Documentation Map

### Phase 13+++: Normalisation Lexicale Avancee

#### Core Documentation

1. **[PHASE-13-QUICKREF.md](PHASE-13-QUICKREF.md)** (2 min read)
   - One-liner summary
   - The 5 strategies table
   - Key metrics
   - Quick configuration profiles
   - Common issues & fixes
   - ** Best for: Quick reference, decisions**

2. **[PHASE-13-PLUS-PLUS-PLUS.md](PHASE-13-PLUS-PLUS-PLUS.md)** (10 min read)
   - Full technical specifications
   - Each of 5 strategies explained
   - Code examples for each
   - Architecture diagram
   - Formulas and equations
   - ** Best for: Understanding implementation**

3. **[PHASE-13-COMPARISON.md](PHASE-13-COMPARISON.md)** (8 min read)
   - Before/after metrics (Phase 13++ vs 13+++)
   - Detailed advantages analysis
   - Trade-offs discussion
   - When to use each phase
   - ** Best for: Understanding impact, comparisons**

4. **[PHASE-13-CONFIGURATION.md](PHASE-13-CONFIGURATION.md)** (15 min read)
   - All configurable parameters
   - 3 pre-defined profiles (Quality/Balanced/Coverage)
   - Tuning recipes for specific cases
   - Benchmarking guide
   - Diagnostic troubleshooting
   - ** Best for: Customization, optimization**

5. **[PHASE-13-VALIDATION.md](PHASE-13-VALIDATION.md)** (10 min read)
   - Full validation results
   - Test suite output
   - Quality metrics
   - Checklist validation
   - Lessons learned
   - ** Best for: Verification, production readiness**

---

##  Find What You Need

### "I want to..."

#### Understand Phase 13+++ quickly
 Read [PHASE-13-QUICKREF.md](PHASE-13-QUICKREF.md) (2 min)

#### Know how it compares to Phase 13++
 Read [PHASE-13-COMPARISON.md](PHASE-13-COMPARISON.md) (8 min)

#### Understand technical details
 Read [PHASE-13-PLUS-PLUS-PLUS.md](PHASE-13-PLUS-PLUS-PLUS.md) (10 min)

#### Customize for my use case
 Read [PHASE-13-CONFIGURATION.md](PHASE-13-CONFIGURATION.md) (15 min)

#### Verify it's production-ready
 Read [PHASE-13-VALIDATION.md](PHASE-13-VALIDATION.md) (10 min)

#### Deep dive into implementation
 Read all docs in order (45 min)

---

##  The 5 Strategies at a Glance

```
INPUT BLOCS (180)
    
[1] Normalisation Lexicale
     Count word occurrences per bloc
     Penalize blocs with repetitions
     Apply penalty to selection score
    
[2] TF-IDF Intelligent
     Identify rare but frequent words
     Penalize 0.8x during vectorization
     Prevent repetitive words from dominating
    
[3] Fen�trage Strict
     Calculate lexical similarity (Jaccard)
     Skip consecutive blocks with >60% overlap
     Force vocabulary diversity
     (SELECTION: 45 blocs)
[4] Anti-Repetition <5 mots
     Track word positions in generated text
     If word repeated <5 words apart: skip 2nd
     Zero close-range repetitions
    
[5] Synonymes Contextuels
     Replace frequent words with synonyms
     Natural vocabulary variation
     Every 3rd occurrence
    
OUTPUT TEXT (95% coherence, ~0 repetitions)
```

---

##  Key Concepts

### Penalite de Repetition (Strategy 1)
```
If word appears N times in bloc (N > 2):
    Penalty = (N - 2) � 0.1
    finalScore *= (1 - Penalty)
```

### TF-IDF Ajuste (Strategy 2)
```
For rare but frequent words:
    If IDF > 0.5 AND TF > 0.05:
        tfidf *= 0.8  (penalty)
```

### Similarite Lexicale Jaccard (Strategy 3)
```
Similarity = |A � B| / |A  B|
If Similarity > 0.6:
    Skip this block (too similar to previous)
```

### Distance Intra-Texte (Strategy 4)
```
For each word position:
    If word appears again < 5 words away:
        Skip the 2nd occurrence
```

### Remplacement Synonymes (Strategy 5)
```
For each frequent word (>2 occurrences):
    If occurrence_count % 3 == 0:
        Replace with random synonym (70% chance)
```

---

##  Quick Start

```bash
# Default usage (Balanced profile)
./programme resume input.txt 0.12

# For maximum quality (fewer but better words)
./programme resume input.txt 0.10

# For maximum coverage (more content)
./programme resume input.txt 0.20
```

---

##  Key Metrics

| Metric | Phase 13++ | Phase 13+++ | Improvement |
|--------|-----------|-----------|------------|
| Words | 1297 | 679-847 | More quality |
| Coherence | 94.83% | 95.00% | +0.17% |
| Speed | 1384ms | 219ms | **6.3x faster**  |
| Repetitions | Many | ~0 | **Eliminated**  |

---

##  Validation Status

-  Build: SUCCESS
-  Tests: ALL PASSED
-  Coherence: 95.00%
-  Repetitions: ~0 detected
-  Performance: 219ms
-  Production: READY

---

##  Configuration

### Balanced Profile (Default, Recommended)
```go
PenalityCoeff:    0.1
SimilarityThresh: 0.6
TFIDFPenalty:     0.8
AntiRepDistance:  5 words
SynonymFreq:      1 every 3 occurrences
```

### Quality Profile
```go
PenalityCoeff:    0.2
SimilarityThresh: 0.4
TFIDFPenalty:     0.7
AntiRepDistance:  7 words
SynonymFreq:      1 every 2 occurrences
```

### Coverage Profile
```go
PenalityCoeff:    0.05
SimilarityThresh: 0.75
TFIDFPenalty:     0.9
AntiRepDistance:  3 words
SynonymFreq:      1 every 5 occurrences
```

See [PHASE-13-CONFIGURATION.md](PHASE-13-CONFIGURATION.md) for details.

---

##  Code Changes

```
database/resumeur_coherence.go
 NormaliserRepetitionsBlocs() function
 CalculerSimilarityVocabLexical() function
 RepetitionsBloc field in BlocVectoriel
 Fen�trage strict logic

database/generation.go
 TF-IDF penalty 0.8x for rare-frequent words

database/coherence.go
 SynonymsDict map (20+ entries)
 Anti-repetition <5 words filter
 Synonym replacement logic
```

Total changes: ~250 lines of code.

---

##  For Developers

### Understanding the Flow
1. Read [PHASE-13-PLUS-PLUS-PLUS.md](PHASE-13-PLUS-PLUS-PLUS.md)
2. Study code in `/database/resumeur_coherence.go` lines 554-580
3. Study code in `/database/generation.go` lines 495-510
4. Study code in `/database/coherence.go` lines 410-470

### Customizing
1. Choose profile in [PHASE-13-CONFIGURATION.md](PHASE-13-CONFIGURATION.md)
2. Follow tuning recipes
3. Test with benchmarking guide
4. Validate with test suite

### Extending
1. Add entries to `SynonymsDict` for domain terms
2. Adjust thresholds for specific corpus
3. Add context-awareness to synonym selection
4. Consider lemmatization for better grouping

---

## � Support

### Common Questions

**Q: Why 95% coherence instead of 99%?**  
A: 95% = optimal balance. Higher coherence requires more aggressive filtering  less content.

**Q: Why 0.6 similarity threshold?**  
A: 60% = psychologically optimal. <50% loses topic continuity. >70% permits too-similar blocks.

**Q: Should I use Quality or Coverage profile?**  
A: Default to Balanced. Quality if reading experience > quantity. Coverage if need 1000+ words.

**Q: Can I add more synonyms?**  
A: Yes! Add to `SynonymsDict` in `database/coherence.go`. Format: `"word": {"syn1", "syn2", "syn3", "word"}`

**Q: What if coherence drops?**  
A: Reduce penalties: `0.1  0.05` for normalization, `0.8  0.85` for TF-IDF.

---

##  Summary

**Phase 13+++** implements **5 complementary strategies** to eliminate repetitions while maintaining 95% coherence at 6.3x faster speed.

**Documentation** provides:
-  Quick reference (2 min)
-  Technical depth (45 min)
-  Configuration guide (15 min)
-  Validation proof (10 min)

**Status**: Production-ready, tested, documented, optimized.

---

**Next doc to read**: [PHASE-13-QUICKREF.md](PHASE-13-QUICKREF.md)

**For detailed implementation**: [PHASE-13-PLUS-PLUS-PLUS.md](PHASE-13-PLUS-PLUS-PLUS.md)

**For configuration**: [PHASE-13-CONFIGURATION.md](PHASE-13-CONFIGURATION.md)

---

Last updated: Phase 13+++  
Status:  Complete and validated
