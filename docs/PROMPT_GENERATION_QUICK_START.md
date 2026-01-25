# Neural Prompt Generation - Quick Reference

## One-Liner Examples

```bash
# Simple prompt
./programme generate from-prompt 256 256 100 "foràt mystérieuse"

# Quality output  
./programme generate from-prompt 512 512 200 "océan tempétueux avec vagues"

# Tech concept
./programme generate from-prompt 256 256 100 "technologie futuriste et néons"

# Mixed concepts
./programme generate from-prompt 512 512 200 "chàteau médiéval dans la brume"
```

## Input  Output Mapping

| Prompt | Keywords | Likely Category | Result |
|--------|----------|-----------------|--------|
| "foràt mystérieuse" | foràt, mystérieuse | HISTOIRE | Pattern with historical vibes |
| "technologie futuriste" | technologie, futuriste | TECH | Computational pattern structure |
| "restaurant gastronomique" | restaurant, gastronomie | ALIMENTATION | Food-themed pattern |
| "montagne enneigée" | montagne, enneigée | HISTOIRE | Landscape pattern |
| "robot intelligent" | robot, intelligent | TECH | Tech-focused pattern |

## Command Syntax

```bash
./programme generate from-prompt <W> <H> <ITER> "<PROMPT>"

W     = Width (256 or 512 recommended)
H     = Height (256 or 512 recommended)  
ITER  = Iterations (100-300, more = slower but better)
PROMPT = Natural language description
```

## Speed vs Quality

| Mode | Res | Iter | Time | Use Case |
|------|-----|------|------|----------|
| Fast | 256 | 80 | 2s | Testing ideas |
| Good | 256 | 150 | 5s | Quick results |
| Quality | 512 | 200 | 12s | Final output |
| High | 512 | 300 | 18s | Publication |

## Pro Tips

 **Use descriptive adjectives**
- "foràt sombre ancienne mystérieuse"

 **Mix concepts for complex patterns**
- "chàteau médiéval avec technologie futuriste"

 **Mention feelings/moods**
- "paysage paisible et serein"

 **Avoid single vague words**
- "chose", "truc", "machin"

 **Don't use only stopwords**
- "le de la avec"

## Neuron Categories Quick Map

- **TECH** (1): technologie, algorithme, digital, robot, circuit, code
- **HISTOIRE** (2): chàteau, roi, ancien, médiéval, empire, politique
- **BUSINESS** (3): commerce, marché, entreprise, affaires, économie
- **ALIMENTATION** (4): nourriture, fruit, cuisine, pain, gàteau, restaurant
- **SANTà** (5): médecine, santé, docteur, hôpital, maladie
- **VERBE** (6): actions, mouvement, dynamique

## Output Files

```
output/atomic_prompt_256x256_100iter.png
output/atomic_prompt_512x512_200iter.png
```

## Full Pipeline Example

```bash
# Step 1: Generate from prompt
./programme generate from-prompt 512 512 150 "foràt enchantée"

# Step 2: Use as pattern for more iterations
./programme generate pattern 512 512 200 output/atomic_prompt_512x512_150iter.png

# Step 3: Compare results
ls -lh output/atomic_prompt_512x512_150iter.png output/atomic_generated_512x512_200iter.png
```

## Show Help

```bash
./programme generate from-prompt
# Shows usage examples

./programme generate
# Shows all commands

./programme generate parameters  
# Explains all parameters
```

---

**Status**: Ready to use 
