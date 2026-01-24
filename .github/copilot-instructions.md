# Copilot Instructions for IA-ATOMIQUE

## Project Overview
**IA-ATOMIQUE** is a sophisticated Go-based neural network system implementing **Atomic Resonance Technology (T.R.A.)** for advanced text analysis, language learning, and intelligent summarization. It's ~30x faster than equivalent local LLMs while consuming minimal resources.

**Core Concept**: Intelligence emerges from **local interactions between autonomous computational atoms** rather than centralized computation. Each atom updates independently using only neighbor data—no global synchronization needed.

**Version**: v4.1 (Beta) | **Language**: Go 1.22+ | **Security**: AES-256 encryption

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

### Quick Start
```bash
# Build
go build -o programme

# Run interactive mode (default REPL)
./programme

# Atomic network simulation (N iterations, M atoms)
./programme atomic simulate 1000 500

# Text analysis
./programme text "Your text here"
./programme file document.txt

# Humanization (3 styles: standard, professional, advanced)
./programme humanize file -s document.txt    # Standard
./programme humanize file -p document.txt    # Professional  
./programme humanize file -a document.txt    # Advanced
```

### Output Files
Text processing generates suffixed files in current directory:
- `_humanized.txt` - Standard style
- `_humanized_prof.txt` - Professional style
- `_humanized_avance.txt` - Advanced style

**Dashboard**: `./programme` generates `dashboard` file showing neuron statistics, energy distribution, active neurons.

### Testing & Benchmarking
```bash
# Benchmark atomic network performance
./programme atomic benchmark

# Run stress tests (in project root)
go run stress_test_levier3_fusion.go
go run stress_test_ultimate.go
```

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

- **Architecture decisions**: [../README-ARTICLE.md](../README-ARTICLE.md), [../ATOMIC-IMPLEMENTATION.md](../ATOMIC-IMPLEMENTATION.md)
- **Implementation details**: [../IMPLEMENTATION-SUMMARY.md](../IMPLEMENTATION-SUMMARY.md)
- **Project status**: [../IMPLEMENTATION-STATUS.txt](../IMPLEMENTATION-STATUS.txt)
- **Symmetry breaking**: [../database/symmetry_breaking.go](../database/symmetry_breaking.go)

## Common Modifications

- **Add category weights**: Modify `Injecter()` calls in `database/data.go`
- **Tune atomic parameters**: Edit `α, β, γ, δ, σ` in `database.NewAtomicNetwork()`
- **Change freeze behavior**: Adjust `FreezeThreshold`, `FreezeIterations`, `WakeThreshold` in `AtomicNetwork`
- **Extend language support**: Add language detection patterns in `interaction.go:DetecterLangue()`

## Image Generation: Symmetry Breaking

**Critical Insight**: Same energy ≠ Same image (degeneracy problem)

The system can reach correct energy equilibrium but produce structured noise instead of recognizable images. This happens because multiple micro-states share the same energy (like different crystal structures at same temperature/pressure).

### The 4 Solutions (`database/symmetry_breaking.go`)

1. **Phase Continuity Term** - Forces coherent gradient propagation
   - `E_phase = Σ ||∇I(x) - ∇I(x+1)||²`
   - Implemented in `computePhaseContinuity()`, `computeLocalGradient()`
   - Edges propagate continuously instead of appearing randomly

2. **Weak Directional Field** - Low-res orientation guide (16×16 or 32×32)
   - NOT a mask/model - just weak structural hints (5-15% influence)
   - Extracted via `ExtractLowResGuidanceField()`
   - Physically: external weak field that orients but doesn't dictate

3. **Topological Constraint** - Match edge topology, not intensity
   - `rank(|∇I(x)|) ≈ rank(|∇I_ref(x)|)` instead of `I(x) ≈ I_ref(x)`
   - Implemented in `EdgeTopologyMap`, `ComputeEdgeTopology()`
   - Edges appear in correct PLACES with flexible VALUES

4. **Multi-Scale Pipeline** - CRUCIAL coarse-to-fine approach
   - 32×32 → 64×64 → 128×128 → 256×256
   - Without this, energy distributes randomly → chaos
   - Each scale establishes structure for next level
   - Use `MultiScalePipeline.RunMultiScalePipeline()`

### Key Pattern
```go
// Multi-scale with symmetry breaking
pipeline := NewMultiScalePipeline()
energyProfile := NewImageEnergyProfile("source.png")
guidanceField := ExtractLowResGuidanceField("source.png", 16)
edgeTopology := ComputeEdgeTopology("source.png")

network := pipeline.RunMultiScalePipeline(energyProfile, guidanceField, edgeTopology)
```

## Important Conventions

### Thread Safety & Concurrency
- Use `sync.Mutex` for concurrent atom updates (see `atomic.go`)
- Atoms iterate independently; neighbor state reads require NO locking (perception-only)
- Go channels preferred over mutex for high-throughput batch operations (see stress tests)
- Atomic package (`sync/atomic.Int64`, etc.) for counters without contention

### State Ranges & Energy
- Atom states and resonance values MUST be clamped to [0, 1]
- Each atom tracks `EnergyConsumption`; network accumulates `TotalEnergy`
- Freeze state automatically activates when atom inactive <ϵ for T iterations (reduces energy)
- Woken by resonance threshold σ_wake

### Code Organization Pattern
- **Main CLI dispatch** in `main.go:main()` — routes to command handlers
- **Text pipeline** in `interaction.go` — language detection → tokenization → keyword extraction → classification
- **Neural operations** in `database/` — neuron/atom mechanics, energy calculations
- **Specialized commands** split into dedicated files: `atomic_cli.go`, `generation_commands.go`, `image_commands.go`, `pattern_commands.go`
- **Testing/Stress** in `stress_test_*.go` files (separate compile units, not in main binary)

### Symmetry Breaking (Image Generation)
CRITICAL: Always use **multi-scale pipeline** for image generation; single-scale produces structured noise.
- Extract guidance field (16×16 or 32×32): `ExtractLowResGuidanceField()`
- Compute edge topology: `ComputeEdgeTopology()`
- Run coarse-to-fine: 32×32 → 64×64 → 128×128 → 256×256
- Each scale establishes structure for next level

### Parameter Tuning Reference
- **α (coupling)**: Neighbor influence strength; higher = faster consensus, lower = local exploration
- **β (local rules)**: Internal rule weight; higher = rigidity, lower = adaptability  
- **γ (reinforcement)**: Weight update strength for coherent patterns; tuned per use case
- **δ (decay)**: Connection decay rate; prevents weight explosion
- **σ (resonance sensitivity)**: Resonance sharpness; tunes network responsiveness

## Common Modifications

- **Add categories**: Modify `Injecter()` in `../database/data.go` — weights control impact
- **Tune atomic parameters**: Edit α, β, γ, δ, σ in `database.NewAtomicNetwork()` calls
- **Extend language support**: Add detection patterns in `../interaction.go` `DetecterLangue()`
- **Freeze thresholds**: Adjust `FreezeThreshold`, `FreezeIterations`, `WakeThreshold` in `AtomicNetwork`
- **Add humanization styles**: Extend style cases in `../interaction.go` `HumanizeText()` function

## File Structure Reference

- `../database/atomic.go` — `ComputationalAtom`, `AtomicNetwork`, resonance engine
- `../database/nlp.go` — Tokenization, keyword extraction, category activation
- `../database/data.go` — Neuron initialization, category definitions, stopwords
- `../database/language.go` — Multi-language detection (FR/EN/DE/ES)
- `../database/phrase_analysis.go` — SVC structure detection
- `../interaction.go` — CLI text processing, humanization modes
- `../atomic_cli.go` — Atomic network simulation, benchmarking
- `../generation_commands.go` — Image generation pipeline
- `../Makefile` — Build targets and deployment scripts
