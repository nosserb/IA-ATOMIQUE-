# Atomic Image Generation - Quick Start Guide

**IA-ATOMIQUE v4.0** now includes a complete **5-phase atomic image generation system** using local resonance principles.

## Quick Examples

### Full Pipeline (All 5 Phases)
```bash
# High-quality image from prompt
./programme image pipeline 512 512 200 8 "mystical purple forest with glowing mushrooms"

# Fast preview
./programme image pipeline 256 256 50 8 "blue ocean with sunset"

# Ultra high-quality (takes longer)
./programme image pipeline 1024 1024 400 8 "detailed fantasy landscape with dragons"
```

### Individual Phase Testing

```bash
# Phase 1: Multi-scale structuration (pixels  patterns)
./programme image phase1 256 256 50 8

# Phase 2: Shape emergence (contours & shapes)
./programme image phase2 256 256 50 8

# Phase 3: Prompt conditioning (apply guidance)
./programme image phase3 256 256 50 8 "dark misty forest"

# Phase 4: Refinement (smooth & texture)
./programme image phase4 256 256 50 8

# Phase 5: Coherence verification (quality control)
./programme image phase5 256 256 50 8
```

## Available Prompt Keywords

### Colors
`red`, `blue`, `green`, `yellow`, `purple`, `orange`, `pink`, `cyan`, `golden`, `silver`, `white`, `black`

### Mood/Lighting
`dark`, `night`, `bright`, `light`, `sunny`, `dim`, `glowing`, `shadowy`, `ominous`, `peaceful`

### Texture
`rough`, `smooth`, `detailed`, `blurry`, `sharp`, `noisy`, `clean`, `grainy`, `fine`

## Parameter Guide

| Parameter | Range | Default | Notes |
|-----------|-------|---------|-------|
| Width | 64-2048 | 512 | Image width |
| Height | 64-2048 | 512 | Image height |
| Iterations | 10-1000 | 200 | Per-phase iterations |
| Patch Size | 1-64 | 8 | Atomic block size (larger = faster) |

## Performance Tips

**Fast (30 seconds)**: 256π256, 50 iterations, patch 8
```bash
./programme image pipeline 256 256 50 8 "prompt"
```

**Balanced (2-3 min)**: 512π512, 200 iterations, patch 8
```bash
./programme image pipeline 512 512 200 8 "prompt"
```

**High Quality (5-10 min)**: 512π512, 200 iterations, patch 4
```bash
./programme image pipeline 512 512 200 4 "prompt"
```

**Ultra Quality (15+ min)**: 1024π1024, 400 iterations, patch 8
```bash
./programme image pipeline 1024 1024 400 8 "prompt"
```

## Output Files

After running pipeline, you get:
- **generated_image_pipeline.png** - Final image
- **coherence_map.png** - Quality visualization (red=low, green=high)

## What Are the 5 Phases?

1. **Structuration**: Local pixel resonance creates coherent patterns
2. **Emergence**: Capsules recognize primitive shapes
3. **Conditioning**: Prompt guides color/mood/texture
4. **Refinement**: Smoothing and texture details
5. **Verification**: Quality check and automatic repair

## Try These Prompts

```bash
# Landscapes
./programme image pipeline 512 512 200 8 "golden sunset over calm ocean"
./programme image pipeline 512 512 200 8 "dark mysterious forest with fog"
./programme image pipeline 512 512 200 8 "bright snowy mountain peak"

# Abstract
./programme image pipeline 512 512 200 8 "colorful abstract geometric patterns"
./programme image pipeline 512 512 200 8 "smooth flowing liquid textures"

# Detailed
./programme image pipeline 512 512 200 8 "detailed fantasy landscape with castle"
./programme image pipeline 512 512 200 8 "intricate sharp detailed fractals"

# Mood-based
./programme image pipeline 512 512 200 8 "peaceful serene meditative atmosphere"
./programme image pipeline 512 512 200 8 "chaotic noisy rough turbulent"
./programme image pipeline 512 512 200 8 "glowing ethereal magical dreamlike"
```

## Troubleshooting

**Blurry output?**
- Reduce Laplacian iterations (Phase 4) or increase patch size

**Image doesn't match prompt?**
- Add more specific color/texture terms
- Increase total iterations
- Try Phase 3 independently to debug

**Artifacts or disconnected regions?**
- Run Phase 5 (coherence verification)
- Increase Phase 1 coupling for better alignment

**Slow generation?**
- Increase patch size: 8  16
- Reduce dimensions: 512  256
- Reduce iterations per phase

## Advanced: Help Text

```bash
./programme image
# Shows full help with all commands and descriptions
```

---

**Created**: January 9, 2026  
**Technology**: Atomic Resonance Technology (T.R.A.)  
**Project**: IA-ATOMIQUE v4.0
