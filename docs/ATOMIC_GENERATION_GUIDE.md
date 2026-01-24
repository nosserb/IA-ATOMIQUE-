# ⚛️ Atomic Generation System - Complete Guide

## Overview

The **Atomic Generation System** transforms abstract wave patterns into complete, recognizable images through **local resonance interactions**. Each pixel is modeled as an autonomous "atom" with:

- **Internal state vector** (5 dimensions): intensity, orientation, phase, frequency, coherence
- **Local connections** to neighboring atoms (8-connected neighborhood)
- **Pattern influence** from injected wave patterns
- **Optional feedback** toward a target image

The system emerges complex global structures from simple local interactions, creating images without centralized computation.

## Mathematical Foundation

### Core Equations

#### 1. Local State Propagation
$$s_{ij}^{(t+1)} = \alpha \cdot \frac{1}{|N|}\sum_{(k,l) \in N} s_{kl}^{(t)} + \beta \cdot P_{ij} + \epsilon \cdot V_{ij}^{(t)}$$

Where:
- $\alpha$ (resonance): How much atoms respond to neighbors [0.1-0.5]
- $\beta$ (pattern weight): Strength of injected pattern [0.2-0.8]
- $\epsilon$ (velocity damping): Momentum smoothing [0.8-1.0]
- $P_{ij}$: Pattern value at position (i,j)
- $V_{ij}^{(t)}$: Velocity/momentum term

**Physical interpretation**: Each atom listens to its neighbors (resonance), follows a guiding pattern, and maintains momentum for smooth transitions.

#### 2. State-to-Color Conversion
$$C_{ij} = f(s_{ij}) = \text{HSL-to-RGB}(h, s, l)$$

Where state components map to HSL color space:
- $h$ (hue) ← state[1] (orientation)
- $s$ (saturation) ← state[4] (coherence)
- $l$ (lightness) ← state[0] (intensity)

**Result**: Internal state automatically converts to vibrant RGB colors.

#### 3. Local Color Smoothing
$$C'_{ij} = (1 - \gamma) \cdot C_{ij} + \gamma \cdot \frac{1}{|N|}\sum_{(k,l) \in N} C_{kl}$$

Where:
- $\gamma$ (smoothing strength): Local blending [0.05-0.3]

**Purpose**: Prevents blocky artifacts while preserving local patterns.

#### 4. Optional Target Feedback
$$s_{ij}^{(t)} \leftarrow s_{ij}^{(t)} + \delta \cdot (C_{\text{target}} - C_{\text{current}})$$

Where:
- $\delta$ (feedback weight): Guidance strength [0.0-1.0]

**Effect**: Gradually guides generation toward target image appearance.

## Command Reference

### Generate from Pattern Only
```bash
./programme generate pattern <width> <height> <iterations> [pattern.png]
```

**Example**: Fast emergence guided by wave patterns
```bash
./programme generate pattern 512 512 200 output/pattern_final_emerged.png
```

**Output**: `output/atomic_generated_512x512_200iter.png`

**Best for**:
- Quick generation without target constraints
- Exploring emergent patterns
- Low iteration count (100-200)

### Generate with Target Feedback
```bash
./programme generate with-feedback <width> <height> <iterations> <target.png>
```

**Example**: Guided generation toward specific image
```bash
./programme generate with-feedback 512 512 300 input/image/face.png
```

**Output**: `output/atomic_feedback_512x512_300iter.png`

**Best for**:
- Constrained generation (recreate specific images)
- Higher quality output
- More iterations (300-500)
- Starting from reference image

### Explain Parameters
```bash
./programme generate parameters
```

Shows detailed explanation of all parameters and their ranges.

### Benchmark
```bash
./programme generate benchmark
```

Tests generation speed at different resolutions for performance tuning.

## Parameter Tuning Guide

### Resonance Alpha (α): 0.1 - 0.5

**Controls how much atoms listen to neighbors**

| Value | Effect | Use Case |
|-------|--------|----------|
| 0.1 | Atoms mostly independent, strong pattern adherence | Sharp details, local patterns |
| 0.2 | Balanced independence & cooperation | Fine detail with cohesion |
| 0.3 | **Default** - Good balance | Most applications |
| 0.4 | Strong cooperation, smooth transitions | Smooth blends |
| 0.5 | Maximum cooperation, highly collective | Emergent global structures |

**Recommendation**: Start with 0.3, increase for smoother results, decrease for more varied details.

### Pattern Beta (β): 0.2 - 0.8

**Controls how strictly atoms follow the injected pattern**

| Value | Effect | Use Case |
|-------|--------|----------|
| 0.2 | Pattern is loose guideline, atoms improvise | Creative variation |
| 0.3 | Balanced creativity vs. pattern adherence | Artistic generation |
| 0.5 | **Default** - Good balance | Most applications |
| 0.6 | Pattern is strong constraint | Preserve pattern structure |
| 0.8 | Atoms closely follow pattern | Strict reproduction |

**Recommendation**: Use 0.5-0.6 to preserve pattern structure while allowing atom interactions.

### Smoothing Gamma (γ): 0.05 - 0.3

**Controls local color smoothing to prevent block artifacts**

| Value | Effect | Use Case |
|-------|--------|----------|
| 0.05 | Minimal smoothing, sharp edges | High-frequency detail |
| 0.10 | Light smoothing, details preserved | Fine grain |
| 0.15 | **Default** - Light smoothing | Most applications |
| 0.20 | Moderate smoothing | Balanced appearance |
| 0.30 | Heavy smoothing, smooth gradients | Soft, blended look |

**Recommendation**: 0.10-0.20 for most cases. Increase if seeing blocky artifacts.

### Feedback Weight (δ): 0.0 - 1.0

**Controls strength of target image guidance (with-feedback mode only)**

| Value | Effect |
|-------|--------|
| 0.0 | No feedback (pattern-only generation) |
| 0.2 | Light guidance, preserves emergence |
| 0.3 | **Default** - Gentle guidance |
| 0.5 | Balanced guidance |
| 1.0 | Maximum guidance, tight target adherence |

**Recommendation**: 0.2-0.4 for artistic generation, 0.5+ for faithful reproduction.

### Velocity Damping (ε): 0.8 - 1.0

**Controls momentum smoothing in state updates**

| Value | Effect |
|-------|--------|
| 0.80 | Responsive, may oscillate |
| 0.90 | **Default** - Good balance |
| 0.95 | Smooth, stable convergence |
| 1.00 | No damping, can become chaotic |

**Recommendation**: Use 0.90-0.95. Adjust only if seeing instability.

## Complete Workflow

### Step 1: Create Pattern from Reference
```bash
# Extract wave pattern from reference image
./programme pattern emerge 512 512 200 input/image/reference.jpg 0.15
# Output: output/pattern_final_emerged.png
```

### Step 2: Generate from Pattern
```bash
# Basic generation guided by pattern
./programme generate pattern 512 512 300 output/pattern_final_emerged.png
# Output: output/atomic_generated_512x512_300iter.png
```

### Step 3: Refine with Target Feedback (Optional)
```bash
# Generate with guidance toward original
./programme generate with-feedback 512 512 400 input/image/reference.jpg
# Output: output/atomic_feedback_512x512_400iter.png
```

### Step 4: Compare Results
```bash
# Compare original → pattern → generated
ls -lh output/pattern_final_emerged.png output/atomic_generated_512x512_300iter.png output/atomic_feedback_512x512_400iter.png
```

## Configuration Examples

### Fine Detail Configuration
```
α = 0.2   (limited cooperation)
β = 0.6   (strong pattern adherence)
γ = 0.10  (minimal smoothing)
```
**Result**: Sharp edges, detailed local patterns, preserves texture

### Smooth Blend Configuration
```
α = 0.4   (strong cooperation)
β = 0.4   (balanced pattern)
γ = 0.25  (heavy smoothing)
```
**Result**: Smooth gradients, natural-looking, cohesive

### Wave-Centric Configuration
```
α = 0.1   (minimal cooperation)
β = 0.8   (strict pattern)
γ = 0.05  (sharp details)
```
**Result**: Pattern dominates, atomic cooperation minimal, faithful reproduction

### Emergent Configuration
```
α = 0.5   (maximum cooperation)
β = 0.3   (loose pattern)
γ = 0.20  (moderate smoothing)
```
**Result**: Complex global structures, creative emergence, less pattern-dependent

## Performance Guidelines

### Iteration Count

| Iterations | Time | Quality | Best For |
|-----------|------|---------|----------|
| 50 | ~2-5s | Low convergence | Quick tests |
| 100 | ~4-10s | Medium | Pattern-only mode |
| 200 | ~8-20s | Good | Standard use |
| 300 | ~12-30s | High | Feedback mode |
| 500+ | ~20s+ | Very high | Publication-quality |

**Recommendation**: Start with 100, increase to 300 for feedback mode.

### Resolution vs. Time

| Size | Atoms | Approx. Time (100 iter) |
|------|-------|----------------------|
| 256² | 65K | ~1-2s |
| 512² | 262K | ~4-8s |
| 1024² | 1M | ~15-30s |

**Note**: Time scales roughly with atom count. Use 512×512 for standard work.

## Troubleshooting

### Problem: Blocky appearance in output
**Solution**: Increase smoothing gamma (γ) to 0.15-0.25

### Problem: Generation looks nothing like pattern
**Solution**: Increase pattern beta (β) to 0.6-0.8, decrease resonance alpha (α) to 0.2

### Problem: Loss not decreasing
**Solution**: This is normal for this system. Loss reflects state-pattern divergence, not image quality.

### Problem: Too uniform/gray
**Solution**: Increase pattern beta (β), inject better-defined pattern

### Problem: Too much noise/chaos
**Solution**: Increase resonance alpha (α) for more cooperation, increase smoothing gamma (γ)

## Advanced Usage

### Iterative Refinement
```bash
# Generate once
./programme generate pattern 512 512 100 pattern.png

# Use output as new pattern input
./programme generate pattern 512 512 100 output/atomic_generated_512x512_100iter.png
```

### Multi-Scale Generation
```bash
# Generate small version first
./programme generate pattern 256 256 200 pattern.png

# Upscale and refine
# (Create 512x512 from 256x256 result)
./programme generate pattern 512 512 200 output/atomic_generated_256x256_200iter.png
```

### Ensemble Averaging
Generate multiple times with different parameters, then average the results for smoother output.

## Understanding the Physics

The atomic generation system models **coupled oscillators with pattern guidance**:

1. **Resonance** (α term): Atoms oscillate in sync with neighbors
2. **Pattern Injection** (β term): External pattern drives atoms
3. **Smoothing**: Local color averaging prevents phase separation
4. **Feedback**: Optional attractor pulling toward target

This creates a **self-organizing system** where:
- Global patterns emerge from local interactions
- No centralized controller needed
- Convergence is natural (atoms settle to low-energy states)
- Results are deterministic (same seed gives same output)

## Output Interpretation

Each generated image shows the **stable state** of the atomic system after:
- All atoms have updated their local state using resonance equation
- State vectors converted to RGB colors
- Local smoothing applied to remove artifacts
- Optional target feedback applied

The **loss value** reported represents average divergence between current state and injected pattern - not a quality metric.

## Further Reading

- [ATOMIC-IMPLEMENTATION.md](ATOMIC-IMPLEMENTATION.md) - Low-level implementation details
- [PATTERN_EMERGENCE_GUIDE.md](PATTERN_EMERGENCE_GUIDE.md) - How to create input patterns
- [README-ATOMIQUE.md](README-ATOMIQUE.md) - Full system overview

## Quick Reference Card

```
PATTERN-ONLY GENERATION:
  ./programme generate pattern 512 512 200 pattern.png

FEEDBACK-GUIDED GENERATION:
  ./programme generate with-feedback 512 512 300 target.png

PARAMETER RANGES:
  α (resonance):  0.1-0.5    (default: 0.3)
  β (pattern):    0.2-0.8    (default: 0.5)
  γ (smoothing):  0.05-0.3   (default: 0.2)
  δ (feedback):   0.0-1.0    (default: 0.3)
  ε (damping):    0.8-1.0    (default: 0.9)

TYPICAL CONFIGS:
  Fine Detail:    α=0.2, β=0.6, γ=0.10
  Smooth Blend:   α=0.4, β=0.4, γ=0.25
  Wave-Centric:   α=0.1, β=0.8, γ=0.05
  Emergent:       α=0.5, β=0.3, γ=0.20

ITERATION GUIDE:
  Pattern-only:   100-200 iterations
  With-feedback:  300-500 iterations
  Resolution:     512×512 recommended
```

---

**Last Updated**: January 2025
**Version**: 1.0
**Status**: Production Ready ✅
