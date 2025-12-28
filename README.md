
<div align="center" style="max-width: 900px; margin: 0 auto; padding: 0 20px;">

<div style="background-color: #fff3cd; border: 2px solid #ff9800; border-radius: 8px; padding: 20px; margin: 20px 0;">
  <p style="color: #ff6f00; font-weight: bold; margin: 0; font-size: 1.1em;">⚠️ BETA VERSION - v4.1</p>
  <p style="color: #333; margin: 10px 0 0 0;">This version may contain bugs. If you encounter any issues or have feedback for improvements, please report them. Your feedback helps us fix problems and improve the system. <strong>Contact: send issues or suggestions</strong></p>
</div>

<p align="center">
  <a href="https://i.ibb.co/WN5frxd6/2ee3d1b24181.png">
    <img src="https://i.ibb.co/LzGHGkWG/8e199e30c114.png" alt="How it works">
  </a>
</p>

<p style="font-size: 1.2em; color: #666666b4; margin-bottom: 30px;"><strong>A sophisticated and powerful neural network for text analysis, learning and summary generation</strong></p>

<div style="margin: 25px 0;">
  <a href="https://golang.org/"><img src="https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat-square&logo=go" alt="Go"></a>
  <span style="margin: 0 10px;"></span>
  <a href="#license"><img src="https://img.shields.io/badge/License-MIT-4CAF50?style=flat-square" alt="License"></a>
  <span style="margin: 0 10px;"></span>
  <img src="https://img.shields.io/badge/Security-AES--256--Encryption-blue?style=flat-square" alt="Security">
</div>

<hr style="margin: 40px 0; border: none; border-top: 2px solid #eee;">

</div>

<div align="center" style="max-width: 900px; margin: 0 auto; padding: 0 20px;">

<div style="text-align: left; margin: 30px 0;">

### Overview

**IA-ATOMIQUE** is a sophisticated neural network system designed for advanced text analysis, language learning, and intelligent summary generation. It analyzes text at the phrase level, classifies content across 6 specialized categories, detects grammatical structures, and learns from input while maintaining ethical content moderation.

**Key Advantage:** 30x faster than equivalent local LLMs while consuming minimal resources.

</div>

<hr style="margin: 40px 0; border: none; border-top: 2px solid #eee;">

<div style="text-align: left; margin: 30px 0;">

### Features

- **Neural Network Analysis** - 1000 sophisticated neurons with phrase-level analysis
- **6-Category Classification**:
  - **TECH** - Technology, programming, digital systems
  - **HISTOIRE** - History, politics, current events
  - **BUSINESS** - Commerce, economy, trade
  - **ALIMENTATION** - Food, nutrition, gastronomy
  - **SANTE** - Health, medicine, wellness
  - **VERBE** - Action detection (principal verbs)

- **Grammatical Structure Detection** - Subject-Verb-Complement analysis
- **Dynamic Learning** - Two-stage probation system for new words
- **Intelligent Summaries** - Context-aware synthesis generation
- **Content Moderation** - Blacklist with AES-256 encryption
- **Real-time Dashboard** - Neural network visualization and monitoring
- **30x Performance** - Massive speed advantage over local LLMs

</div>

<hr style="margin: 40px 0; border: none; border-top: 2px solid #eee;">

<div style="text-align: left; margin: 30px 0;">

### Quick Start

#### Requirements

- **Go** 1.22 or higher
- **Git** for version control

#### Installation

```bash
# Clone the repository
git clone <repository-url>
cd "IA-ATOMIQUE-"

# Download dependencies (already configured)
go mod download

# Build the executable
go build -o programme
```

#### First Run

```bash
# Display help
./programme help

# Run interactive mode
./programme

# Analyze a text file
./programme file wikipedia.txt

# Analyze inline text
./programme text "L'intelligence artificielle progresse rapidement"
```

</div>

<hr style="margin: 40px 0; border: none; border-top: 2px solid #eee;">

<div style="text-align: left; margin: 30px 0;">

### Command Reference

#### Get Help

Display complete usage guide with examples:

```bash
./programme help
# or
./programme -h
./programme --help
```

---

#### Mode 1: Interactive Mode (Default)

Perfect for exploring the system interactively with step-by-step learning:

```bash
# Default interactive
./programme

# Explicit interactive
./programme interactive
./programme inter
```

**What happens:**
- Prompts for input
- Analyzes each phrase
- Shows classification and confidence
- Learns from your input
- Generates summaries

**Example session:**
```
[IA-ATOMIQUE v4.1] Mode Interactif
Entrez du texte (ou 'quit' pour arrêter):
> Einstein a découvert la théorie de la relativité
[HISTOIRE] - Confiance: 87% - Énergie: 5.23
Résumé: Einstein découverte théorie relativité
```

---

#### Mode 2: File Analysis

Analyze entire files (documents, articles, reports):

```bash
./programme file <path_to_file>
```

**Examples:**

```bash
# Analyze a Wikipedia article
./programme file wikipedia.txt

# Analyze a Markdown document
./programme file document.md

# Absolute path
./programme file /home/user/documents/article.txt
```

**Output includes:**
- Global statistics (phrases analyzed, total energy, confidence)
- Category distribution with percentages and bar charts
- Per-phrase analysis with key words
- Detected verbs in `[verbe: ...]` format
- Synthetic summaries for each category

**Supported formats:**
- Plain text (.txt)
- Markdown (.md)
- HTML (.html)
- Any text-based format

---

#### Mode 3: Inline Text Analysis

Quick analysis of text from command line:

```bash
./programme text "<your text here>"
```

**Examples:**

```bash
# Single sentence
./programme text "Les neurones artificiels imitent le cerveau humain"

# Multiple sentences (quoted as single argument)
./programme text "Python est un langage populaire. L'IA utilise le machine learning."

# Current events
./programme text "Le changement climatique affecte la biodiversité"
```

**Output:**
- Immediate classification
- Energy score
- Key terms identified
- Synthetic summary

---

#### Mode 4: Classic Mode (Legacy)

Simple text processing without arguments structure:

```bash
./programme <text>
./programme bonjour
./programme "l'histoire de l'informatique"
./programme python javascript golang
```

This mode treats all arguments as a single phrase for analysis.

</div>

<hr style="margin: 40px 0; border: none; border-top: 2px solid #eee;">

<div style="text-align: left; margin: 30px 0;">

### Output Explanation

#### Global Statistics

```
[STATISTIQUES GLOBALES]
• Phrases analysées: 7
• Catégories détectées: 2
• Énergie totale: 15.72
• Confiance moyenne: 62.7%
```

**Metrics:**
- **Phrases analyzed** - Total number of sentences processed
- **Categories detected** - How many categories are represented
- **Total energy** - Cumulative neuron activation energy
- **Average confidence** - Mean confidence across all classifications

#### Category Distribution

```
[DISTRIBUTION]
  HISTOIRE     █████████████████ 86% (6 phrases)
  TECH         ██ 14% (1 phrases)
```

Visual bar chart showing:
- Category name
- Percentage of phrases in that category
- Number of phrases
- Visual proportion

#### Per-Phrase Analysis

```
[TECH] - 1 phrases (énergie: 2.50)
Mots clés: son | formalisme | mathématique | efficacité

Phrases clés:
  • Si son formalisme mathématique est d'une efficacité inégalée...

Résumé: Formalisme mathématique efficacité
```

Shows:
- Category and phrase count
- Total energy for that category
- Key terms (mots clés)
- Original phrases
- Synthetic summary

#### Verb Detection

Verbs are detected separately and displayed:

```
[verbe: découvert]
[verbe: analyse]
```

This prevents action words from polluting the content classification.

</div>

<hr style="margin: 40px 0; border: none; border-top: 2px solid #eee;">

<div style="text-align: left; margin: 30px 0;">

### How It Works

#### 1. Phrase-Level Analysis

Text is broken down into individual phrases (not just words), allowing independent categorization without one category dominating others.

#### 2. Neural Classification

Each phrase activates relevant neurons across the 6 categories based on keyword matches and learned patterns. Energy scores indicate confidence.

#### 3. Learning System

**Two-stage probation:**
- Stage 1: New word observed, added to temporary learning list
- Stage 2: Word seen 2+ times, confirmed in lexicon with permanent category
- Prevents misclassification from single occurrences

#### 4. Content Moderation

The system maintains an encrypted blacklist of 130+ offensive words:
- Prevents learning inappropriate content
- Blocks these words during analysis
- Logs `[BLOQUÉ]` when encounter restricted words
- Encryption ensures system integrity and ethical alignment

#### 5. Summary Generation

Context-aware synthesis extracts key terms while preserving meaning:
- Removes stop words (le, la, de, etc.)
- Maintains important nouns and verbs
- Orders by significance
- Creates readable, compact summaries

</div>

<hr style="margin: 40px 0; border: none; border-top: 2px solid #eee;">

<div style="text-align: left; margin: 30px 0;">

### Performance Benchmarks

#### Core Metrics

| Metric | Result |
|--------|--------|
| **Throughput** | 790 phrases/second |
| **Latency** | 1.26 ms/phrase |
| **Memory** | <20 MB RAM |
| **CPU Avg** | ~23% utilization |
| **Binary Size** | 2.3 MB |

#### Stress Testing

| Test | Result |
|------|--------|
| **1000 Phrases** | 31ms (32,107 phrases/sec peak) |
| **Stability** | <20% variation over 5 runs |
| **Scalability** | Linear with text size |

#### Comparison vs Local LLM

| Metric | IA-ATOMIQUE | LLM Local | Advantage |
|--------|-------------|-----------|-----------|
#### Real-World Performance (Measured)

| Case | Time | Remark |
|------|------|--------|
| **Short text** (1 phrase) | **6ms** | Extremely fast |
| **Long text** (~40 phrases) | **16ms** | ~0.4ms per phrase |
| **HTML file** | ~150ms | HTML parsing included |

**Measured results** :
```bash
# Short text
$ time ./programme text "This is a simple test..."
real    0m0.006s

# Full TXT file
$ time ./programme file input.txt
real    0m0.016s
```

**Observed performance** : **16x faster** than conservative estimates (30ms)! 

</div>

<hr style="margin: 40px 0; border: none; border-top: 2px solid #eee;">

<div style="text-align: left; margin: 30px 0;">

### Security & Ethics

#### Encrypted Blacklist

The content moderation blacklist is secured with:
- **File**: `blacklist.enc` (encrypted, unreadable)
- **Protection**: Prevents unauthorized modification
- **Runtime**: Transparently decrypted during execution
- **Integrity**: Ethical values enforcement at system level

#### Why Encryption?

> "If someone deletes words from here, it no longer aligns with the values at the time of publication on GitHub by nosserb."

The encryption ensures the content moderation system remains integral and trustworthy.

#### What's Blocked

The blacklist prevents learning of:
- Obscene language
- Hate speech
- Offensive slurs
- Discriminatory terms
- Explicit content

This ensures IA-ATOMIQUE maintains ethical standards aligned with GitHub's community policies.

</div>

<hr style="margin: 40px 0; border: none; border-top: 2px solid #eee;">

<div style="text-align: left; margin: 30px 0;">

### Project Structure

```
IA-ATOMIQUE-/
├── main.go                    Entry point & command routing
├── interaction.go             Interactive mode & file processing
├── database/
│   ├── data.go               Neurons, lexicon, encryption
│   └── phrase_analysis.go    Phrase-level analysis engine
├── go.mod                     Go dependencies
├── go.sum                     Dependency checksums
├── blacklist.enc              Encrypted content moderation (AES-256)
├── dashboard                  Neural network visualization
├── programme                  Compiled executable
└── README.md                  This documentation
```

</div>

<hr style="margin: 40px 0; border: none; border-top: 2px solid #eee;">

<div style="text-align: left; margin: 30px 0;">

### Usage Tips

#### For Best Results

1. **Use French text** - Optimized for French language patterns
2. **Natural sentences** - Phrase-based analysis works better with grammatically correct text
3. **Diverse topics** - Mix content to see category distribution
4. **Iterative learning** - Run multiple times with similar texts to strengthen patterns
5. **Large documents** - Files with 100+ phrases show better statistical significance

#### Performance Tips

```bash
# Fast analysis of large file
time ./programme file large_document.txt

# Profile with multiple categories
./programme file wikipedia.txt

# Build release binary for speed
go build -ldflags="-s -w" -o programme
```

#### Debugging

The system logs:
- `[BLACKLIST]` messages for content moderation
- `[BLOQUÉ]` when blocked words are encountered
- Energy levels and confidence scores
- Category distributions

</div>

<hr style="margin: 40px 0; border: none; border-top: 2px solid #eee;">

<div style="text-align: left; margin: 30px 0;">

### Examples

#### Example 1: Analyze Technology Article

```bash
./programme text "Python est un langage de programmation puissant. L'IA utilise Python pour le machine learning. TensorFlow et PyTorch sont des frameworks populaires."
```

**Expected Output:**
```
[STATISTIQUES GLOBALES]
• Phrases analysées: 3
• Catégories détectées: 1
• Énergie totale: 18.5
• Confiance moyenne: 94.2%

[DISTRIBUTION]
  TECH  ██████████████████ 100% (3 phrases)

[TECH] - 3 phrases (énergie: 18.5)
Mots clés: Python | programmation | IA | machine learning | TensorFlow

Résumé: Python langage programmation IA machine learning
```

---

#### Example 2: Mixed Historical & Business Content

```bash
./programme file wikipedia.txt
```

**Expected Output:**
```
[DISTRIBUTION]
  HISTOIRE  ████████████████ 60% (9 phrases)
  BUSINESS  ██████ 40% (6 phrases)

[ANALYSE PAR PHRASE]
[HISTOIRE] Napoléon a transformé l'Europe...
[BUSINESS] Le commerce s'est développé...
```

---

#### Example 3: Interactive Learning Session

```bash
./programme interactive
```

**Session:**
```
Entrez du texte:
> Einstein a développé la théorie de la relativité en 1905
[HISTOIRE] - Confiance: 89%
Verbe détecté: [verbe: développé]
Résumé: Einstein théorie relativité

Entrez du texte:
> quit
```

</div>

<hr style="margin: 40px 0; border: none; border-top: 2px solid #eee;">

<div style="text-align: left; margin: 30px 0;">

### Development

#### Build from Source

```bash
go build -o programme
```

#### Run Tests (if applicable)

```bash
go test ./...
```

#### Clean Build

```bash
go clean
go build -o programme
```

</div>

<hr style="margin: 40px 0; border: none; border-top: 2px solid #eee;">

<div style="text-align: left; margin: 30px 0;">

### Development Phases

| Phase | Feature | Impact |
|-------|---------|--------|
| **1-2** | 7-criteria validation + auto-regeneration (3 attempts, 80% threshold) | Guaranteed summary quality |
| **3-4** | Two-pass architecture (extract → enrich) + log masking | Contextual summaries, clean output |
| **5** | 8-step HTML pipeline (main content extraction) | Reliable web page processing |
| **6** | Phrase extraction expansion (7→15) and integration (6→12) | More complete summaries |
| **7** | **Performance optimization** : QuickSort + Builder (84x faster!) | TXT: 2.5s → 30ms |
| **8** | Optimized cleaning (17 ReplaceAll → map lookup) | Stable performance |
| **9** | SHA256 integrity check (instead of AES-256) | Transparent security |
| **10** | **3-layer semantic analysis** | Dominant idea + Classification + Semantic axes |
| **11** | Cleanup (21 unnecessary files removed) | Lean codebase |

</div>

<hr style="margin: 40px 0; border: none; border-top: 2px solid #eee;">

<div style="text-align: left; margin: 30px 0;">

### License

This project is distributed under the **MIT** license. See the [LICENSE](LICENSE) file for more details.

</div>

<hr style="margin: 40px 0; border: none; border-top: 2px solid #eee;">

<div style="margin: 40px 0;">
  <p><strong>nosserb | 2025</strong></p>
  <p style="color: #999; font-size: 0.95em;">⭐ Star this project if you find it useful!</p>
</div>

</div>

