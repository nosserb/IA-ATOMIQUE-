# ⚛️ Atomic Generation - 5-Minute Quick Start

## What It Does

Transforms abstract **wave patterns** → **complete images** using autonomous atomic resonance.

Each pixel is an "atom" that:
- Listens to its neighbors
- Follows an injected pattern
- Converts internal state to color
- Iteratively stabilizes

Result: Complex images emerge from simple local interactions.

## Installation

Already included in the build. Just use:
```bash
./programme generate <command>
```

## Quick Start: 3 Steps

### 1️⃣ Create a Pattern (Optional)

If you have a reference image, extract its pattern:
```bash
./programme pattern emerge 256 256 100 input/image/reference.jpg 0.15
# Output: output/pattern_final_emerged.png
```

Or use an existing pattern/test image.

### 2️⃣ Generate Image from Pattern

**Basic (pattern-only):**
```bash
./programme generate pattern 256 256 150 output/pattern_final_emerged.png
```

**Guided (with target feedback):**
```bash
./programme generate with-feedback 256 256 200 input/image/reference.jpg
```

### 3️⃣ View Results

```bash
# Check output
ls -lh output/atomic_*.png

# View the image
# Linux: eog output/atomic_generated_256x256_150iter.png
# macOS: open output/atomic_generated_256x256_150iter.png
```

## Common Commands

```bash
# Show help
./programme generate

# Show all parameters explained
./programme generate parameters

# Run speed benchmark
./programme generate benchmark

# 256×256 quick test
./programme generate pattern 256 256 50 output/test.png

# 512×512 quality output
./programme generate pattern 512 512 300 output/pattern.png

# With target guidance
./programme generate with-feedback 512 512 400 input/image/target.jpg
```

## Key Parameters (Tweak These)

| Parameter | Range | Default | Effect |
|-----------|-------|---------|--------|
| **α** Resonance | 0.1-0.5 | 0.3 | How much atoms cooperate (higher = smoother) |
| **β** Pattern | 0.2-0.8 | 0.5 | How strictly atoms follow pattern (higher = more faithful) |
| **γ** Smoothing | 0.05-0.3 | 0.2 | Local color blending (higher = less blocky) |
| **δ** Feedback | 0.0-1.0 | 0.3 | Target guidance (higher = closer to target) |
| **ε** Damping | 0.8-1.0 | 0.9 | Momentum smoothing (higher = more stable) |

**Quick tuning**:
- Too blocky? → Increase γ to 0.25
- Doesn't follow pattern? → Increase β to 0.7
- Too uniform? → Increase pattern better or increase iterations
- Too chaotic? → Increase α to 0.4

## Complete Workflow Example

```bash
# 1. Create pattern from reference
./programme pattern emerge 512 512 150 input/image/my_photo.jpg 0.15

# 2. Generate image from pattern
./programme generate pattern 512 512 200 output/pattern_final_emerged.png

# 3. Refine with target guidance
./programme generate with-feedback 512 512 300 input/image/my_photo.jpg

# 4. Compare results
ls -lh output/pattern_final_emerged.png \
       output/atomic_generated_512x512_200iter.png \
       output/atomic_feedback_512x512_300iter.png
```

## Understanding the Output

**Loss value** (~0.08): How different atoms' states are from the injected pattern. Lower is better but even 0.07+ is fine.

**Iterations**: More = better convergence, but with diminishing returns. 200-300 is usually optimal.

**Size**: 512×512 is sweet spot. Larger is slower but higher quality.

## Typical Use Cases

### Quick Artistic Generation
```bash
./programme generate pattern 256 256 100 pattern.png
# ~2-3 seconds, decent results
```

### Quality Output
```bash
./programme generate pattern 512 512 300 pattern.png
# ~10-15 seconds, high quality
```

### Faithful Reproduction
```bash
./programme generate with-feedback 512 512 400 target.jpg
# ~15-20 seconds, closely matches target
```

### Exploration/Experimentation
```bash
# Vary resonance (α)
./programme generate pattern 256 256 100 pattern.png  # α=0.3 default
# Increase cooperation
./programme generate pattern 256 256 100 pattern.png  # α=0.5

# Vary pattern adherence (β)
./programme generate pattern 256 256 100 pattern.png  # β=0.5 default
# Strict pattern
./programme generate pattern 256 256 100 pattern.png  # β=0.8
```

## Troubleshooting

**Q: Image looks like noise**
- A: Pattern input might be weak. Try stronger pattern or increase β to 0.7

**Q: Image is too uniform/gray**
- A: Pattern might be too faint. Increase iterations to 300+

**Q: Very slow**
- A: Use smaller size (256×256) or fewer iterations (100)

**Q: Blocky appearance**
- A: Increase γ (smoothing) to 0.25-0.3

**Q: Different every run**
- A: This is normal for complex systems. Results are deterministic within a run.

## Next Steps

- 📖 Read [ATOMIC_GENERATION_GUIDE.md](ATOMIC_GENERATION_GUIDE.md) for complete documentation
- 🔬 Experiment with parameters: `./programme generate parameters`
- 📊 Check performance: `./programme generate benchmark`
- 🎨 Create patterns: `./programme pattern --help`

## File Locations

| Type | Location |
|------|----------|
| Input patterns | `output/pattern_*.png` |
| Generated images | `output/atomic_generated_*.png` |
| With-feedback outputs | `output/atomic_feedback_*.png` |
| Help & docs | `*.md` files in root |

---

**Status**: ✅ Ready to use
**Last tested**: January 2025
