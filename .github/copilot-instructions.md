# Copilot Instructions for IA-ATOMIQUE

## Project Overview
**IA-ATOMIQUE** is a sophisticated Go-based neural network system implementing **Atomic Resonance Technology (T.R.A.)** for advanced text analysis, language learning, and intelligent summarization. It's ~30x faster than equivalent local LLMs while consuming minimal resources.

Core concept: Intelligence emerges from **local interactions between autonomous computational atoms** rather than centralized computation.

## Architecture Overview

### Core Components

1. **`database/atomic.go`** - The Atomic Resonance Engine
   - `ComputationalAtom`: Elementary autonomous unit with internal state, local rules, perceptions, and neighbor connections
   - `AtomicNetwork`: Distributed network implementing asynchronous iteration
   - Resonance equation: `R(si, sj) = exp(-||si - sj||²/2σ²)`
   - Key parameters: α (coupling), β (local rules), γ (reinforcement), δ (decay), σ (resonance sensitivity)
   - **Critical pattern**: Each atom updates independently using only neighbor data + local rules; no central coordinator

2. **`database/` module** - Supporting NLP & Data Systems
   - `data.go`: 1000-neuron network initialization, word categorization, lexique management
   - `nlp.go`: Text tokenization, keyword extraction, category activation
   - `language.go`: Multi-language detection (FR/EN/DE/ES), filtering for French content
   - `phrase_analysis.go`: Subject-Verb-Complement (SVC) structure detection

3. **Top-level CLI Components**
   - `main.go`: REPL, dashboard generation (neural monitor output)
   - `atomic_cli.go`: Atomic simulation commands (`simulate`, `network-stats`, `benchmark`)
   - `interaction.go`: Text analysis, humanization, language-specific processing

### 6-Category Classification System
Words are categorized at initialization (see `database/data.go:Injecter()`):
- **Cat 1: TECH** (weight 7.0) - Programming, algorithms, hardware, digital systems
- **Cat 2: HISTOIRE** (weight 6.5) - History, politics, governmental structures
- **Cat 3: BUSINESS** (weight 6.0) - Commerce, finance, employment
- **Cat 4: ALIMENTATION** (weight 6.0) - Food, nutrition, recipes
- **Cat 5: SANTÉ** - Health, medicine, wellness
- **Cat 6: VERBE** - Action/verb detection

## Critical Development Patterns

### Asynchronous Atom Iteration
```go
// In atomic.go: IterateNetwork()
// Each atom updates independently:
// 1. Compute resonance with neighbors
// 2. Calculate neighbor influence (α parameter)
// 3. Apply local rules/perceptions (β parameter)
// 4. Clamp state to [0, 1] range
// 5. Update connection weights via: dwij/dt = γ*coherence - δ*wij
```
**Key insight**: No global synchronization required; atoms naturally align via resonance.

### Energy Efficiency & Freeze System
Atoms enter "freeze" state when inactive (<ϵ threshold for T iterations) to reduce energy consumption. Woken by resonance threshold σ_wake. See `atomic.go` for freeze state management.

### Text Processing Pipeline
1. **Language Detection** (`interaction.go:DetecterLangue`) - Counts accents and stopwords for FR/EN/DE/ES
2. **Tokenization** (`nlp.go:TokeniserTexte`) - Removes stopwords (predefined in `data.go`)
3. **Keyword Extraction** (`nlp.go:ExtraireMotsClés`) - Scores by frequency × category weight
4. **Category Activation** (`nlp.go:ActiverCategoriesParTexte`) - Maps words to neurons by category
5. **Confidence Calculation** - Based on activation distribution

### Humanization (Multi-mode)
Three style modes in `interaction.go`:
- `-s` (Standard): Natural, fluent French
- `-p` (Professional): Formal, technical language
- `-a` (Advanced): Style analysis + intelligent paraphrasing

Each mode applies language-specific transformations (conjugations, register adjustments).

## Build & Execution

```bash
# Build
go build -o programme

# Run interactive mode (default)
./programme

# Simulate atomic network (N iterations, M atoms)
./programme atomic simulate 1000 500

# Analyze file
./programme file document.txt

# Humanize with style
./programme humanize file -p document.txt

# Analyze text directly
./programme text "Your text here"
```

**Dashboard output** (`dashboard` file): Displays neuron statistics, energy distribution, active neurones, and emergent behaviors.

## Integration Points & Data Flow

1. **Blacklist Security** (`blacklist.enc`)
   - AES-256 encrypted content moderation list
   - SHA-256 integrity check (hash in `data.go`)

2. **Dynamic Lexique** (`LexiqueTemp`)
   - Two-stage probation for unknown words
   - Transition from probation to confirmed vocabulary

3. **Neuron-Atom Bridge**
   - Legacy neuron network (1000 neurons × 50+ categories) coexists with new AtomicNetwork
   - Both used for analysis; atomic system used for simulations

4. **File Input/Output**
   - Input: Plain text, markdown, formatted documents
   - Output: `_humanized.txt`, `_humanized_prof.txt`, `_humanized_avance.txt`

## Key Files for Reference

- **Architecture decisions**: [README-ARTICLE.md](README-ARTICLE.md), [ATOMIC-IMPLEMENTATION.md](ATOMIC-IMPLEMENTATION.md)
- **Implementation details**: [IMPLEMENTATION-SUMMARY.md](IMPLEMENTATION-SUMMARY.md)
- **Project status**: [IMPLEMENTATION-STATUS.txt](IMPLEMENTATION-STATUS.txt)

## Common Modifications

- **Add category weights**: Modify `Injecter()` calls in `database/data.go`
- **Tune atomic parameters**: Edit `α, β, γ, δ, σ` in `database.NewAtomicNetwork()`
- **Change freeze behavior**: Adjust `FreezeThreshold`, `FreezeIterations`, `WakeThreshold` in `AtomicNetwork`
- **Extend language support**: Add language detection patterns in `interaction.go:DetecterLangue()`

## Important Conventions

- **Thread safety**: Use `mutex.Lock/Unlock` for concurrent atom updates
- **State ranges**: Atom states and resonance values are clamped to [0, 1]
- **Energy tracking**: Each atom tracks `EnergyConsumption`; network accumulates `TotalEnergy`
- **Asynchronous safety**: Atoms iterate independently; only neighbor states are read (no locking needed for perception)
