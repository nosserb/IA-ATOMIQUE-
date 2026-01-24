# IA-ATOMIQUE v4.1

> Advanced Neural Network System with Atomic Resonance Technology

[![Go 1.22+](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat-square&logo=go)](https://golang.org/)
[![License MIT](https://img.shields.io/badge/License-MIT-4CAF50?style=flat-square)](LICENSE)
[![Security AES-256](https://img.shields.io/badge/Security-AES--256-blue?style=flat-square)](docs/SECURITY.md)

## Overview

**IA-ATOMIQUE** is a sophisticated Go-based neural network implementing **Atomic Resonance Technology (T.R.A.)** for advanced text analysis, language learning, and intelligent summarization. It's ~30x faster than equivalent local LLMs while consuming minimal resources.

### Core Concept

Intelligence emerges from **local interactions between autonomous computational atoms** rather than centralized computation. Each atom updates independently using only neighbor data—no global synchronization needed.

## ⚡ Quick Start

```bash
# Build
go build -o programme

# Interactive LLM chat mode
./programme start

# Analyze text
./programme text "Your text here"

# Analyze file
./programme file document.txt

# Generate summary
./programme resume document.txt 0.1

# Atomic network simulation
./programme atomic simulate 1000 500
```

## Key Features

### Neural Processing
- **1000-neuron network** with 50+ categories
- **Atomic Resonance Technology** - decentralized, asynchronous computation
- **6-category classification**: TECH, HISTOIRE, BUSINESS, ALIMENTATION, SANTÉ, VERBE
- **Subject-Verb-Complement (SVC)** structure detection
- **Multi-language support**: FR/EN/DE/ES

### Text Analysis & Summarization
- **Grammar-aware summarization** with 5 scientific pillars
- **Three humanization modes**: standard, professional, advanced
- **Keyword extraction** with category weighting
- **Automatic learning** from input text
- **Content moderation** with AES-256 encrypted blacklist

### Image Generation & Processing
- **Energy-based image generation** from text prompts
- **Multi-scale pipeline** (32×32 → 256×256)
- **Symmetry breaking** for coherent image synthesis
- **Ultra deblur** with 327,000+ atoms
- **Motion blur removal** (Lucy-Richardson deconvolution)

### Advanced Capabilities
- **Pattern emergence** detection
- **Cellular relaxation** networks
- **Academic benchmarks** (MMLU, HellaSwag samples)
- **Knowledge base** extraction and query
- **Needle-in-haystack** semantic search

## Interactive LLM Mode

```bash
./programme start
```

Chat-like interface with slash commands:
- `/file <path>` - Analyze a file
- `/text <content>` - Analyze text
- `/resume <file> [threshold]` - Generate summary
- `/humanize <file> [-s|-p|-a]` - Rewrite file
- `/atomic simulate <iter> [atoms]` - Run simulation
- `/ask <question>` - Quick Q&A
- `/help` - Show all commands

Type text without slash for direct analysis.

## Architecture

### Atomic Network Structure
```
ComputationalAtom (autonomous unit)
  ├─ InternalState [0, 1]
  ├─ LocalRules (β parameter)
  ├─ Perceptions
  └─ Neighbors (local connections)

AtomicNetwork (distributed)
  ├─ Asynchronous iteration
  ├─ Resonance equation: R(si, sj) = exp(-||si - sj||²/2σ²)
  ├─ Weight update: dwij/dt = γ*coherence - δ*wij
  └─ Freeze system (energy efficiency)
```

### Key Parameters
- **α (coupling)**: Neighbor influence strength
- **β (local rules)**: Internal rule weight
- **γ (reinforcement)**: Weight update strength
- **δ (decay)**: Connection decay rate
- **σ (resonance sensitivity)**: Resonance sharpness

## Commands Reference

### Text Processing
```bash
# Analyze file
./programme file document.txt

# Humanize (3 styles)
./programme humanize file document.txt -s    # Standard
./programme humanize file document.txt -p    # Professional
./programme humanize file document.txt -a    # Advanced

# Generate summary
./programme resume document.txt 0.10         # 10% compression
./programme resume document.txt 0.30         # 30% compression

# Compare summarizers
./programme compare document.txt
```

### Atomic Network
```bash
# Simulation (N iterations, M atoms)
./programme atomic simulate 1000 500

# Network statistics
./programme atomic network-stats

# Benchmark performance
./programme atomic benchmark
```

### Image Generation
```bash
# Generate from prompt
./programme image generate "sunset over ocean"

# Ultra deblur (4K)
./programme deblur ultra blurry.jpg output.png

# Combo (motion + atomic)
./programme combo blurry.jpg clean.png
```

### Learning & Knowledge
```bash
# Learn from text
./programme learn document.txt

# Query knowledge
./programme knowledge "Napoleon"

# Knowledge base stats
./programme stats-kb
```

### Academic Tests
```bash
# Run all benchmarks
./programme academic all

# Specific tests
./programme academic mmlu
./programme academic hellaswag
```

## Use Cases

1. **Text Analysis** - Analyze documents, extract keywords, classify content
2. **Summarization** - Generate intelligent summaries with grammar awareness
3. **Learning** - Extract knowledge from corpus, build knowledge graphs
4. **Humanization** - Rewrite text in natural, professional, or advanced style
5. **Image Generation** - Create images from text prompts using energy paradigm
6. **Deblurring** - Remove blur from images with atomic relaxation
7. **Research** - Test atomic resonance principles, benchmark performance

## Documentation

- [ Complete Documentation Index](docs/INDEX.md)
- [ Atomic Implementation Guide](docs/ATOMIC-IMPLEMENTATION.md)
- [ Image Generation Guide](docs/IMAGE-GENERATION-GUIDE.md)
- [ Cellular Emergence](docs/CELLULAR_EMERGENCE_README.md)
- [ Benchmarks & Tests](docs/TESTS_GUIDE.md)
- [ Installation Guide](docs/INSTALL.md)
- [ Web Interface Setup](docs/SETUP_WEB.md)

## Installation

### Prerequisites
- Go 1.22 or higher
- Git

### Linux/macOS
```bash
git clone https://github.com/yourusername/IA-ATOMIQUE.git
cd IA-ATOMIQUE
go build -o programme
./programme start
```

### Windows
```bash
git clone https://github.com/yourusername/IA-ATOMIQUE.git
cd IA-ATOMIQUE
go build -o programme.exe
programme.exe start
```

See [INSTALL.md](docs/INSTALL.md) for detailed instructions.

## Technical Highlights

### Atomic Resonance Technology
- **Decentralized**: Each atom operates independently
- **Asynchronous**: No global clock or synchronization
- **Resonance-based**: Atoms align via state compatibility
- **Adaptive**: Connections strengthen/weaken dynamically
- **Energy-efficient**: Freeze system for inactive atoms
- **Emergent**: Global intelligence from local interactions

### Performance
- **3.96M words/sec** text processing
- **30x faster** than equivalent local LLMs
- **O(n) complexity** for text chunking
- **327,000+ atoms** for ultra deblur
- **Multi-scale pipeline** for image generation

## Contributing

Contributions are welcome! Please read [CONTRIBUTING.md](CONTRIBUTING.md) before submitting PRs.

## License

This project is licensed under the MIT License - see [LICENSE](LICENSE) for details.

## Known Issues

This is a **beta version (v5.2)**. Known issues:
- Some edge cases in multi-language detection
- Image generation can be slow on low-end hardware
- Web interface requires manual setup

Please report bugs via GitHub Issues.

## Acknowledgments

- Built with Go 1.22+
- Inspired by cellular automata and emergence theory
- Uses AES-256 encryption for content moderation
- Atomic resonance concept adapted from physics principles

## Contact

For questions, suggestions, or bug reports:
- GitHub Issues: [Report here](https://github.com/nosserb/IA-ATOMIQUE/issues)
- Email: guylann.bresson.gb@gmail.com

---

**Made with by the IA-ATOMIQUE team**
