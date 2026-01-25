#  Neural Prompt  Image Generation - Quick Card

## Command
```bash
./programme generate from-prompt <W> <H> <ITER> "<PROMPT>"
```

## Quick Examples
```bash
# Fast (2-3s)
./programme generate from-prompt 256 256 100 "for�t mysterieuse"

# Quality (10-15s)
./programme generate from-prompt 512 512 200 "ocean tempetueux"

# High quality (15-20s)
./programme generate from-prompt 512 512 300 "ch�teau medieval avec technologie"
```

## How It Works
1. **Analyze** prompt with neuron network
2. **Extract** keywords (for�t, arbre, etc.)
3. **Activate** neuron categories (HISTOIRE, TECH, etc.)
4. **Create** pattern from neural activation
5. **Generate** image through atomic resonance
6. **Save** to PNG file

## Output File
- Format: `output/atomic_prompt_WxH_ITERiter.png`
- Example: `output/atomic_prompt_512x512_200iter.png`

## Best Prompts
 Use adjectives: "for�t **sombre** **ancienne** **mysterieuse**"
 Mix concepts: "ch�teau **et** technologie"
 Be specific: "ocean **tempetueux** avec **vagues geantes**"

## Categories
- **TECH** (1): technologie, algorithme, robot, circuit, code
- **HISTOIRE** (2): ch�teau, roi, medieval, ancien, empire
- **BUSINESS** (3): commerce, marche, entreprise
- **ALIMENTATION** (4): nourriture, fruit, cuisine, pain
- **SANT�** (5): medecine, docteur, hopital
- **VERBE** (6): action, mouvement

## Speed Reference
| Size | Iter | Time | Use |
|------|------|------|-----|
| 256^2 | 80 | 2s | Test |
| 256^2 | 150 | 5s | Good |
| 512^2 | 200 | 12s | Quality |
| 512^2 | 300 | 18s | Premium |

## Show Help
```bash
./programme generate              # All commands
./programme generate from-prompt  # This command only
./programme generate parameters   # Tuning guide
```

## Tips
- More keywords = stronger pattern
- Longer descriptions = better results
- Mix categories for complexity
- Increase iterations for detail
- Use 512�512 for standard work

## Full Documentation
 **TEXT_TO_IMAGE_PROMPT_GUIDE.md** - Complete reference
 **PROMPT_GENERATION_QUICK_START.md** - 5-min guide
