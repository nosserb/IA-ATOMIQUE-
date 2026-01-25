#  IA-ATOMIQUE: 5-Phase Atomic Image Generation Pipeline

## Overview

This document describes the complete **5-phase atomic image generation system** based on **Atomic Resonance Technology (T.R.A.)**. This approach generates images through local interactions between autonomous computational atoms, without centralized processing.

**Key Innovation**: Intelligence emerges from bottom-up local resonance rather than top-down hierarchical neural networks.

---

## Phase 1: Multi-Scale Structuration (Structuration Multi-àchelle)

### Objective
Transform isolated pixels into coherent patterns through **local resonance** at multiple scales (micro, meso, macro).

### Mathematical Foundation

For each pixel at position $(i,j)$:

$$c_{i,j}(t+1) = c_{i,j}(t) + \alpha \sum_{k \in N(p_{i,j})} w_{i,j,k} \cdot (c_k(t) - c_{i,j}(t))$$

Where:
- $c_{i,j}(t)$ = RGB color of pixel at time $t$
- $\alpha$ = coupling coefficient (influence strength)  0.7
- $w_{i,j,k}$ = weight of influence from neighbor $k$
- $N(p_{i,j})$ = 8-neighborhood of pixel
- $c_k(t) - c_{i,j}(t)$ = color difference driving alignment

### Key Features

1. **Multi-Scale Layers**
   - **Micro**: Individual pixels (scale = 1)
   - **Meso**: 4à4 or 8à8 blocks (scale = 4, 8)
   - **Macro**: Full image patterns (scale = 16+)

2. **Resonance Calculation**
   $$R(s_i, s_j) = \exp\left(-\frac{\|s_i - s_j\|^2}{2\sigma^2}\right)$$
   - Gaussian similarity function
   - High resonance  colors more similar
   - Drives color alignment

3. **Cross-Scale Propagation**
   - Finer scales (micro) compute first
   - Patterns propagate upward to coarser scales
   - Ensures consistency across resolutions

### Implementation

```go
// Run Phase 1 for 50 iterations with multi-scale resonance
network.PhaseOne_StructurationMultiEchelle(50)
```

### Usage Example

```bash
# Generate with Phase 1 only
./programme image phase1 512 512 100 8
# Output: phase1_output.png
```

### Expected Results

-  Individual pixels cluster into color regions
-  Smooth transitions between regions
-  Elimination of random noise
-  Formation of basic textures

---

## Phase 2: Shape Emergence (àmergence de Forme)

### Objective
Make **primitive shapes** appear (contours, lines, curves) through **capsule-based resonance**.

### Mathematical Foundation

For each block/capsule with state vector $s_m$ (texture, intensity, orientation, etc.):

$$s_m(t+1) = s_m(t) + \gamma \sum_{n \in N(B_m)} R(s_n(t), s_m(t))$$

Where:
- $s_m$ = 6D state vector: [texture, intensity, orientation, coherence, entropy, stability]
- $\gamma$ = reinforcement factor  0.15
- $R(a,b)$ = **motif compatibility function** (Gaussian similarity)
- Capsules are blocks of atoms organized in a 2D grid

### Key Features

1. **Resonance Capsules**
   - Each block/patch has a state vector
   - States encode visual properties
   - Capsules interact only with 8 neighbors

2. **Compatibility Function**
   $$\text{compat}(s_a, s_b) = \exp\left(-\frac{\sum_i (s_a^i - s_b^i)^2}{2\sigma^2}\right)$$
   - Measures similarity in all 6 dimensions
   - Drives convergence toward local consensus

3. **Shape Recognition**
   - Edge detection: high intensity + orientation changes
   - Pattern stabilization through reinforcement
   - Emerges without explicit feature training

### Implementation

```go
// Create shape emergence engine
engine := database.NewShapeEmergenceEngine(width, height, blockSize)

// Run Phase 2 for 50 iterations
network.PhaseTwo_ShapeEmergence(engine, 50)
```

### Usage Example

```bash
./programme image phase2 512 512 100 8
# Output: phase2_output.png
```

### Expected Results

-  Contours and edges clearly defined
-  Connected components form recognizable shapes
-  Orientation fields emerge naturally
-  Texture coherence within regions

---

## Phase 3: Prompt Conditioning (Conditionnement par Prompt)

### Objective
**Steer image generation** toward user's creative intention through natural language prompts.

### Mathematical Foundation

Modify the Phase 1 equation with external guidance:

$$c_{i,j}(t+1) = c_{i,j}(t) + \alpha \sum_{k} w_{i,j,k}(c_k - c_{i,j}) + \beta \cdot G_{i,j}(P)$$

Where:
- $G_{i,j}(P)$ = guidance function based on prompt $P$
- $\beta$ = guidance strength coefficient  0.3
- Guidance includes:
  - Color targets (e.g., "red"  RGB(0.8, -0.4, -0.4))
  - Brightness modifiers ("bright", "dark", "night")
  - Texture properties ("smooth", "rough", "detailed")

### Key Features

1. **Prompt Parsing**
   - Extract color terms: "red", "blue", "golden", etc.
   - Extract mood/lighting: "dark", "bright", "gloomy", "sunny"
   - Extract texture: "smooth", "rough", "sharp", "blurry"

2. **Spatial Guidance**
   - Gaussian fall-off from central focus point
   - Stronger effect near center, weaker at edges
   - Allows object localization

3. **Style Vector**
   - Global RGB modifier $[r, g, b]$ from prompt
   - Blends multiple color terms
   - Clamped to $[-1, 1]$ range

### Implementation

```go
// Parse prompt into guide
guide := network.ParsePromptToGuide("red sunset over mountains")

// Apply conditioning for 50 iterations
network.PhaseThree_PromptConditioning(guide, 50)
```

### Usage Example

```bash
./programme image phase3 512 512 100 8 "dark forest with mysterious atmosphere"
# Output: phase3_output.png
```

### Supported Prompt Keywords

**Colors**: red, blue, green, yellow, purple, orange, pink, cyan, golden, silver, white, black

**Mood/Lighting**: dark, night, bright, light, sunny, dim, glowing, shadowy, ominous, peaceful

**Texture**: rough, smooth, detailed, blurry, sharp, noisy, clean, grainy, fine

### Example Prompts

- "bright sunny beach with palm trees"
- "dark misty forest"
- "golden sunset over calm ocean"
- "sharp detailed abstract geometric patterns"
- "smooth dreamy pastel landscape"

---

## Phase 4: Iterative Refinement (Affinement Itératif)

### Objective
Add **fine details** and **realistic texture** through local smoothing and noise adjustment.

### Mathematical Foundation

For each iteration:

$$c_{i,j}^{(n+1)} = c_{i,j}^{(n)} + \delta \cdot \text{Laplacian}(c_{i,j}^{(n)}) + \epsilon \cdot \text{NoiseAdjust}(c_{i,j}^{(n)})$$

Where:
- **Laplacian operator**: $\nabla^2 c = c_{xx} + c_{yy}$ (second derivative, smooths/sharpens)
- $\delta$ = smoothing strength  0.1
- $\epsilon$ = noise strength  0.08
- Laplacian computed from neighbors: $(c_{up} + c_{down} + c_{left} + c_{right} - 4c_{center})$

### Key Features

1. **Laplacian Smoothing**
   - Reduces blocky artifacts
   - Smooths color transitions
   - Prevents disconnected regions
   - Implementation: 2D discrete Laplacian

2. **Texture Addition**
   - Prevents overly "flat" rendered images
   - Uses procedural noise: $\sin(t \cdot 0.3 + x \cdot 0.7) \cdot \cos(y \cdot 0.5)$
   - Noise varies spatially and temporally
   - Creates natural variation

3. **Multi-Pass Refinement**
   - Laplacian applied every iteration
   - Noise applied with time-varying magnitude
   - Gradual accumulation of fine details

### Implementation

```go
// Run Phase 4 for 50 iterations
// Laplacian strength = 0.1, Noise strength = 0.08
network.PhaseFour_IterativeRefinement(0.1, 0.08, 50)
```

### Usage Example

```bash
./programme image phase4 512 512 100 8
# Output: phase4_output.png
```

### Expected Results

-  Smooth color gradients without visible blocks
-  Natural texture variation
-  Fine details and patterns
-  Reduced aliasing and artifacts
-  Overall visual quality improved

---

## Phase 5: Coherence Verification (Vérification de Cohérence Atomique)

### Objective
Perform **quality control** and **automatic repair** of low-coherence regions.

### Mathematical Foundation

For each atom, compute local coherence:

$$\text{coherence}_{i,j} = 1 - \frac{\sum_{k \in N(p_{i,j})} \|c_{i,j} - c_k\|}{\text{max\_possible\_diff}}$$

Where:
- Measures color consistency with neighbors
- Range: [0, 1] (1 = perfect coherence with neighbors)
- max_possible_diff = $3 \sqrt{3}$ (maximum RGB distance)

### Key Features

1. **Coherence Computation**
   - For each atom, compute average color difference to neighbors
   - Flag atoms with coherence < 0.3 as "faulty"
   - Generate coherence map showing quality distribution

2. **Global Health Score**
   - Percentage of atoms with good coherence
   - Indicates overall image quality
   - Target: > 95%

3. **Automatic Repair**
   - Identify low-coherence atoms
   - Blend with neighborhood average (40% weight)
   - Improves local color consistency
   - Preserves important features

4. **Coherence Map Visualization**
   - Red = low coherence (problematic)
   - Green = high coherence (good)
   - Blue channel for contrast
   - Helps identify problem areas

### Implementation

```go
// Run verification and get detailed report
report := network.PhaseFive_CoherenceVerification()

// Access results
fmt.Printf("Global Coherence: %.3f\n", report.GlobalCoherence)
fmt.Printf("Health Score: %.1f%%\n", report.OverallHealthScore * 100)
fmt.Printf("Repaired atoms: %d\n", report.RepairCount)

// Render coherence visualization
coherenceImg := report.RenderCoherenceMap(width, height)
```

### Usage Example

```bash
./programme image phase5 512 512 100 8
# Output: phase5_output.png
```

### Report Structure

```
CoherenceReport {
  GlobalCoherence: 0.87          // Overall coherence [0, 1]
  CoherenceMap: [][]float64      // Per-atom coherence
  FaultyAtoms: []               // Low-coherence atoms detected
  RepairCount: 23               // Atoms automatically fixed
  OverallHealthScore: 0.96      // Percentage of healthy atoms
}
```

---

## Complete 5-Phase Pipeline

### Sequential Execution

```bash
./programme image pipeline 512 512 200 8 "golden sunset over calm ocean"
```

### Pipeline Order

1. **Phase 1** (iterations / 5): Multi-scale structuration
2. **Phase 2** (iterations / 5): Shape emergence
3. **Phase 3** (iterations / 5): Prompt conditioning
4. **Phase 4** (iterations / 5): Iterative refinement
5. **Phase 5** (iterations / 5): Coherence verification

### Recommended Iteration Counts

| Use Case | Width | Height | Total Iterations | Patch Size |
|----------|-------|--------|------------------|------------|
| Quick preview | 256 | 256 | 50 | 8 |
| Standard quality | 512 | 512 | 200 | 8 |
| High quality | 1024 | 1024 | 400 | 8 |
| Ultra detailed | 1024 | 1024 | 800 | 4 |
| Experimental | 2048 | 2048 | 1000 | 16 |

---

## Usage Guide

### Full Pipeline Example

```bash
# High-quality image generation with prompt
./programme image pipeline 512 512 200 8 "mystical purple forest with glowing mushrooms"

# Output files:
# - generated_image_pipeline.png  (final image)
# - coherence_map.png            (quality visualization)
```

### Phase-by-Phase Testing

```bash
# Test individual phases for debugging/experimentation

./programme image phase1 256 256 50 8
# Test multi-scale structuration

./programme image phase3 256 256 50 8 "blue ocean waves"
# Test prompt-guided generation

./programme image phase4 256 256 50 8
# Test texture refinement
```

### CLI Commands Reference

```bash
# Complete pipeline
./programme image pipeline <w> <h> <iter> <patch> "prompt"

# Individual phases
./programme image phase1 <w> <h> <iter> <patch>
./programme image phase2 <w> <h> <iter> <patch>
./programme image phase3 <w> <h> <iter> <patch> "prompt"
./programme image phase4 <w> <h> <iter> <patch>
./programme image phase5 <w> <h> <iter> <patch>

# Traditional commands
./programme image generate <w> <h> <iter> <patch> "prompt"
./programme image prompt "description"
```

---

## Technical Details

### Computational Complexity

- **Per atom, per iteration**: O(8) neighborhood reads + resonance computation
- **Total per iteration**: O(N) where N = (width/patch) à (height/patch)
- **Full pipeline**: O(5 à N à iterations)
- **Memory**: O(N) atoms + auxiliary structures

### Parallelization

- Phase 1: Parallel across all atoms (8-way neighborhood read-only)
- Phase 2: Parallel capsule updates (independent state updates)
- Phase 3: Parallel atom guidance (independent external field)
- Phase 4: Parallel Laplacian computation (local stencil)
- Phase 5: Parallel coherence checks + sequential repair

### Energy Efficiency

- Frozen atoms consume minimal energy
- Resonance-based updates reduce redundant computation
- Natural convergence detection via coherence metrics

---

## Advanced Topics

### Tuning Parameters

```go
// Phase 1: Coupling coefficient (higher = faster alignment)
layer.Alpha = 0.7  // [0.1 to 0.9]

// Phase 2: Reinforcement factor (higher = stronger stabilization)
engine.Gamma = 0.15  // [0.05 to 0.3]

// Phase 3: Guidance strength (higher = more prompt influence)
guide.Strength = 0.3  // [0.1 to 0.8]

// Phase 4: Smoothing & noise
// Laplacian strength: 0.1
// Noise strength: 0.08
// Both in range [0.0 to 0.3]

// Phase 5: Coherence threshold
// Flag atoms below 0.3 as faulty
// Repair blend factor: 0.4
```

### Custom Prompt Parsing

Extend `ParsePrompt()` to add new keywords:

```go
// Example: Add "metallic" texture
textureMappings["metallic"] = 0.6
```

### Integration with Text System

Image-to-text bridge:
1. Generate image with Phase 1-4
2. Analyze capsule patterns (Phase 2)
3. Map shapes to semantic concepts
4. Generate descriptive text
5. Use as prompt feedback loop

---

## Experimental Results

### Convergence Behavior

- **Micro scale**: Converges in 20-50 iterations
- **Macro scale**: Requires 200+ iterations for stability
- **Shape emergence**: 50-100 iterations optimal
- **Prompt strength**: Visible within 10 iterations

### Quality Metrics

- **Global Coherence**: 0.7-0.9 range (0.85+ excellent)
- **Health Score**: 90%+ indicates minimal artifacts
- **Active Atoms**: 60-90% typical
- **Frozen Atoms**: Increases over time (efficiency)

---

## Troubleshooting

### Problem: Blurry or overly smooth output
- **Solution**: Reduce Phase 4 Laplacian strength or iterations
- **Alternative**: Increase patch size (faster but coarser)

### Problem: Image doesn't match prompt
- **Solution**: Increase Phase 3 guidance strength or iterations
- **Try**: More specific color/texture terms in prompt

### Problem: Artifacts or disconnected regions
- **Solution**: Run Phase 5 coherence verification
- **Increase**: Phase 1 coupling coefficient for better alignment

### Problem: Very slow generation
- **Solution**: Increase patch size (8 or 16 instead of 4)
- **Alternative**: Reduce total iterations or image dimensions

---

## Mathematical Papers & References

This implementation is based on:
- Atomic Resonance Technology (T.R.A.) principles
- Local-only computation without global coordination
- Resonance equations: $R(s_i, s_j) = \exp(-||s_i - s_j||^2 / 2\sigma^2)$
- Asynchronous distributed algorithms

---

## Future Enhancements

- [ ] GPU acceleration for large images
- [ ] Adaptive parameter tuning based on content
- [ ] Multi-prompt blending (multiple guiding fields)
- [ ] Real-time interactive generation
- [ ] Integration with pre-trained style embeddings
- [ ] 3D volumetric generation (Phase 1-5 in 3D)
- [ ] Video generation (temporal coherence)

---

**Generated**: January 9, 2026
**Project**: IA-ATOMIQUE v4.0
**Technology**: Atomic Resonance Technology (T.R.A.)
