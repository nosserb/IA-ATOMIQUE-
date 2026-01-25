#  Pattern Database - Quick Start (5 Minutes)

## What You Can Do Now

- **Index images**  Get metadata
- **Search patterns**  Find by category or keyword
- **Use patterns**  Influence image generation

---

## In 30 Seconds

```bash
# 1. Index your images
./programme pattern index input/image

# 2. View results
./programme pattern list

# 3. Done! Use with generation
./programme generate from-prompt 512 512 200 "forest"
```

---

## One-Liner Examples

### Index

```bash
# Index images in input/image folder
./programme pattern index input/image
```

### View

```bash
# List all patterns
./programme pattern list

# Get stats
./programme pattern stats

# Details about one
./programme pattern info input
```

### Search

```bash
# Find all TECH patterns
./programme pattern search category TECH

# Find patterns with keywords
./programme pattern search keyword forest
```

### Generate (with patterns)

```bash
# Generate 512x512 image in 200 iterations
./programme generate from-prompt 512 512 200 "dark forest"

# The system automatically:
# - Loads patterns.db
# - Finds forest/dark related patterns
# - Uses them to guide generation
```

---

## Workflow

### Step 1: Prepare Images

```
Put images in: input/image/
Examples:
  - forest.png
  - technology.jpg
  - ocean.jpg
  - food.png
```

### Step 2: Index

```bash
./programme pattern index input/image
```

**Output:**
```
 PATTERN INDEXING ENGINE
Scanning input directory: input/image
Found 7 image(s) to index

[1/7] Processing: forest.png
   ID: forest
  Size: 512x512
  Categories: HISTOIRE(4) ALIMENTATION(1)
  Complexity: 0.65 | Confidence: 91.2%

[Results shown for all images...]

 Indexed 7 pattern(s)
 Database saved to: patterns.db
```

### Step 3: Review

```bash
./programme pattern list
```

Shows all indexed patterns with summaries.

### Step 4: Generate

```bash
./programme generate from-prompt 512 512 200 "dark forest landscape"
```

System uses patterns automatically!

---

## Common Tasks

### Find TECH patterns

```bash
./programme pattern search category TECH
```

### Find food-related patterns

```bash
./programme pattern search keyword food
```

### See pattern details

```bash
./programme pattern info forest
```

Shows:
- Dimensions
- Colors
- Complexity
- Keywords
- Categories
- Confidence

### Rebuild database

```bash
rm patterns.db
./programme pattern index input/image
```

---

## Speed Reference

| Operation | Time |
|-----------|------|
| Index 7 images | ~10 seconds |
| List patterns | < 100ms |
| Get pattern info | < 50ms |
| Search category | < 100ms |
| Generate with patterns | Same as normal |

---

## Category Quick Map

| Category | Typical Content | Color |
|----------|-----------------|-------|
| TECH | Technology, circuits, digital | Blue |
| HISTOIRE | Historical, medieval, ancient | Warm |
| BUSINESS | Markets, commerce, enterprise | Mixed |
| ALIMENTATION | Food, fruits, meals | Green |
| SANTà | Health, medicine, hospitals | Green/Blue |

---

## Keyboard Commands Cheatsheet

```
# Index all in input/image
./programme pattern index input/image

# List all
./programme pattern list

# Stats
./programme pattern stats

# Search TECH
./programme pattern search category TECH

# Search keyword
./programme pattern search keyword forest

# Details
./programme pattern info <id>

# Generate (automatic pattern use)
./programme generate from-prompt 512 512 200 "description"
```

---

## Pro Tips

1. **Name images clearly**
   ```
   Good: dark_forest.png, tech_circuit.jpg
   Bad: img1.png, photo.jpg
   ```

2. **Keep images organized**
   ```
   input/image/
    forest.png
    ocean.jpg
    technology.png
   ```

3. **Index after adding images**
   ```bash
   cp new_image.png input/image/
   ./programme pattern index input/image  # Refresh database
   ```

4. **Check results**
   ```bash
   ./programme pattern stats  # Quick overview
   ./programme pattern list   # Full listing
   ```

5. **Use in generation**
   ```bash
   ./programme generate from-prompt 512 512 200 "your description"
   # System finds matching patterns automatically
   ```

---

## FAQ

**Q: How do I add new images?**
```
A: Copy to input/image/, then:
   ./programme pattern index input/image
```

**Q: Can I use different folder?**
```
A: Yes!
   ./programme pattern index custom_folder patterns.db
```

**Q: What if database gets corrupted?**
```
A: Rebuild it:
   rm patterns.db
   ./programme pattern index input/image
```

**Q: Does indexing slow down generation?**
```
A: No! Pattern lookup is instant.
```

**Q: Can I search by color?**
```
A: Not directly. Use categories or keywords instead.
```

---

## What Next?

-  Read **[PATTERN_DATABASE_GUIDE.md](PATTERN_DATABASE_GUIDE.md)** for complete reference
-  Try generating images with prompts (uses patterns!)
-  Explore your patterns with search commands
-  Check statistics and metadata

---

**5-minute tutorial complete!** 

Start indexing:
```bash
./programme pattern index input/image
```

List patterns:
```bash
./programme pattern list
```

Generate with patterns:
```bash
./programme generate from-prompt 512 512 200 "forest"
```
