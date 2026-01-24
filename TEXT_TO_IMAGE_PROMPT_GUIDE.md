# 🧠 Text-to-Image Generation via Neural Prompts

## Overview

Generate images directly from **natural language prompts** using the neural network system. The system:

1. **Analyzes** your text prompt using the neuron network
2. **Extracts** keywords and activates relevant neuron categories  
3. **Creates** a pattern from neural activations
4. **Generates** a final image through atomic resonance

This bridges language understanding and visual generation—what you describe in words becomes visual structure through neural activation.

## How It Works

### Step 1: Prompt Analysis 🧠
```
Input: "une forêt mystérieuse avec des arbres luminescents"
         ↓
Tokenization → [forêt, mystérieuse, arbres, luminescents]
         ↓
Neuron Activation → Category 1 (TECH): 2 neurons
                  → Category 2 (HISTOIRE): 0 neurons
                  → etc.
         ↓
Semantic Confidence: 100%
```

The neural system:
- Tokenizes the prompt text
- Removes stopwords (a, de, avec, etc.)
- Extracts meaningful keywords
- Maps words to neuron categories
- Calculates confidence score

### Step 2: Pattern Injection 📌
Neural activations become spatial color patterns:
- Each activated neuron influences a color region
- Activation strength determines saturation
- Category ID determines hue
- Confidence determines lightness
- Distributed across the grid with variation

### Step 3: Atomic Generation ⚛️
Standard atomic resonance with neural-injected pattern:
- Atoms propagate states locally (resonance α)
- Follow neural pattern (pattern β)
- Smooth colors locally (smoothing γ)
- Stabilize to final image

### Step 4: Output 🖼️
Final PNG image reflecting the semantic content of your prompt.

## Command Syntax

```bash
./programme generate from-prompt <width> <height> <iterations> "<prompt_text>"
```

### Parameters

| Parameter | Range | Example |
|-----------|-------|---------|
| **width** | 128-1024 | 512 |
| **height** | 128-1024 | 512 |
| **iterations** | 50-500 | 200 |
| **prompt** | Any French/English text | "forêt mystérieuse" |

## Examples

### Fast Generation (Low Detail)
```bash
./programme generate from-prompt 256 256 100 "une montagne enneigée"
# Time: ~2-3 seconds
# Output: output/atomic_prompt_256x256_100iter.png
```

### Quality Generation (High Detail)
```bash
./programme generate from-prompt 512 512 300 "océan tempétueux avec des vagues"
# Time: ~12-15 seconds
# Output: output/atomic_prompt_512x512_300iter.png
```

### Complex Scene
```bash
./programme generate from-prompt 512 512 200 "château médiéval dans la brume avec des torches"
# Activates: HISTOIRE (castle), mixed semantic content
# Output: Complex spatial patterns reflecting multiple concepts
```

### Tech-Focused
```bash
./programme generate from-prompt 512 512 200 "technologie futuriste avec algorithmes et circuits"
# Activates: TECH category (multiple neurons)
# Output: Patterns reflecting computational concepts
```

## Neuron Categories & Keywords

The system recognizes these semantic categories:

| Category | ID | Keywords |
|----------|----|----|
| **TECH** | 1 | technologie, algorithme, circuit, code, digital, futuriste, robot |
| **HISTOIRE** | 2 | histoire, château, medieval, roi, empire, politique, ancien |
| **BUSINESS** | 3 | commerce, économie, marché, entreprise, vente, affaires |
| **ALIMENTATION** | 4 | nourriture, fruit, cuisine, restaurant, manger, pain, gâteau |
| **SANTÉ** | 5 | médecine, santé, hôpital, docteur, maladie, cure |
| **VERBE** | 6 | actions detected from verb forms |

**Keywords are extracted** from your prompt and mapped to these categories. More keywords in a category = stronger activation.

## Output Format

The generation shows:

```
🧠 Phase 1: Analyze Prompt
   Input: "your prompt text"
   Keywords extracted: [word1, word2, ...]
   Semantic confidence: XX%
   N categories activated
   
   Top activations:
      • CATEGORY_NAME: N neurons

🔧 Phase 2: Initialize Grid (256x256)
   65536 atoms initialized

📌 Phase 3: Inject Pattern
   Neural activations injected as pattern

⚛️ Phase 4: Generation (100 iterations)
   Progress at 25%, 50%, 75%, 100%

✅ Complete
   Output: output/atomic_prompt_256x256_100iter.png
```

## Tips for Better Results

### Keywords Matter
✅ **Good**: "forêt sombre avec des arbres anciens"
   - Multiple descriptive keywords
   - Clear visual concepts
   
❌ **Less Good**: "la forêt"
   - Too generic, few keywords extracted

### Be Descriptive
✅ **Good**: "océan tempétueux avec vagues géantes et tempête"
   - More activation from multiple keywords
   - Stronger pattern injection
   
❌ **Less Good**: "mer"
   - Single word, minimal activation

### Mixing Categories
✅ "technologie futuriste dans un château médiéval"
   - Mixes TECH + HISTOIRE
   - Creates complex pattern interplay
   
✅ "restaurant avec nourriture saine"
   - Mixes ALIMENTATION + SANTÉ
   - Blended semantic influence

### Emotion & Adjectives
✅ "forêt sombre, mystérieuse, ancienne"
   - Adjectives create stronger activation
   - Mood influences pattern density
   
✅ "technologie brillante et futuriste"
   - Color descriptors enhance pattern

## Parameter Tuning for Prompts

### Iterations
- **Fast test** (80-100): Quick feedback, less detail
- **Standard** (150-200): Good balance
- **Quality** (300+): Fine details, slow

### Resolution
- **Small** (256×256): 2-3 seconds, fast iteration
- **Medium** (512×512): 10-15 seconds, good quality
- **Large** (1024×1024): 30s+, very detailed

### Typical Workflow
```bash
# 1. Quick test with descriptive prompt
./programme generate from-prompt 256 256 80 "your description"

# 2. If good, scale up to 512x512 and 150+ iterations
./programme generate from-prompt 512 512 150 "your description"

# 3. For final output, use 300+ iterations
./programme generate from-prompt 512 512 300 "your description"
```

## Understanding Output Quality

The generated image is determined by:

1. **Prompt Keywords** → Neuron activation pattern
2. **Category Mix** → Which neuron categories light up
3. **Iterations** → How much atoms stabilize
4. **Resolution** → Detail level possible
5. **Atomic Parameters** → Local resonance/smoothing

If result seems weak:
- Add more descriptive words to prompt
- Increase iterations (use 200+ instead of 100)
- Use larger resolution (512×512 instead of 256×256)

If result is too uniform/smooth:
- Your prompt may activate only one category
- Try mixing different semantic concepts
- Add more adjectives

## Advanced Usage

### Combining with Other Modes

**Prompt → Pattern → Atomic Generation**
```bash
# 1. Generate from prompt
./programme generate from-prompt 512 512 150 "forêt mystérieuse"

# 2. Use output as pattern for more iterations
./programme generate pattern 512 512 200 output/atomic_prompt_512x512_150iter.png

# 3. Optionally refine with feedback
./programme generate with-feedback 512 512 200 input/image/reference.jpg
```

### Iterative Refinement
```bash
# 1. First generation
./programme generate from-prompt 256 256 100 "concept initial"

# 2. Like it? Scale up
./programme generate from-prompt 512 512 250 "same concept"

# 3. Final polish
./programme generate from-prompt 1024 1024 400 "same concept"
```

### Experimenting with Keywords
```bash
# Test 1: Just location
./programme generate from-prompt 256 256 100 "forêt"

# Test 2: Location + mood
./programme generate from-prompt 256 256 100 "forêt sombre mystérieuse"

# Test 3: Location + mood + elements  
./programme generate from-prompt 256 256 100 "forêt sombre avec arbres anciens"

# Compare outputs to see how keywords influence pattern
```

## Technical Details

### Neural Activation Algorithm
1. **Tokenization**: Split prompt on whitespace, remove stopwords
2. **Keyword Extraction**: Keep meaningful words
3. **Category Mapping**: Each keyword maps to neuron categories via predefined dictionary
4. **Activation Count**: Each word increments activation counter for its category
5. **Confidence**: Calculated from number and quality of activated keywords

### Pattern Generation
- For each neuron category with activation:
  - Hue ← Category ID (maps to 0-1 range)
  - Saturation ← Activation strength (0-1)
  - Lightness ← Overall confidence
  - Distributed spatially with pseudo-random variation (deterministic)

### Atomic Parameters (Fixed)
- α (resonance) = 0.3 (moderate cooperation)
- β (pattern) = 0.5 (balance pattern/emergence)
- γ (smoothing) = 0.2 (moderate smoothing)
- ε (damping) = 0.9 (stable convergence)

These are optimized for prompt-based generation. Adjust in code if needed.

## Common Patterns & Results

### Nature Prompts
```bash
./programme generate from-prompt 512 512 200 "montagne enneigée sous la lune"
→ Likely activates: HISTOIRE (medieval connotation)
   Result: Spatial structures reflecting mountain concept
```

### Tech Prompts
```bash
./programme generate from-prompt 512 512 200 "réseau de neurones et algorithmes"
→ Likely activates: TECH
   Result: Pattern reflecting computational structures
```

### Food Prompts
```bash
./programme generate from-prompt 256 256 100 "festin médiéval avec pain et fruits"
→ Likely activates: ALIMENTATION + HISTOIRE
   Result: Mixed patterns from both categories
```

## Troubleshooting

**Q: Output looks like noise**
- A: Prompt may have no activating keywords
- Try: Add more descriptive words

**Q: Output is uniform gray**
- A: Only one neuron category activated with low intensity
- Try: Mix multiple semantic concepts, add adjectives

**Q: All outputs look similar**
- A: Atomic parameters are fixed for stability
- Try: Vary prompt keywords significantly, increase iterations

**Q: Very slow**
- A: Likely using large resolution
- Try: Test with 256×256 first, scale up later

**Q: Keywords not extracted**
- A: Your words may be stopwords (a, de, le, etc.)
- Try: Use more content words, fewer function words

## File Locations

```
output/
  atomic_prompt_256x256_100iter.png    (from-prompt output)
  atomic_prompt_512x512_200iter.png    (higher quality)
  atomic_generated_*.png               (pattern mode)
  atomic_feedback_*.png                (feedback mode)
```

---

**Feature Status**: ✅ Complete & Tested
**Latest Test**: January 9, 2025
**Language Support**: French & English keywords
