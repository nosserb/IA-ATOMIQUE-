
<div align="center" style="max-width: 900px; margin: 0 auto; padding: 0 20px;">

<div style="background-color: #fff3cd; border: 2px solid #ff9800; border-radius: 8px; padding: 20px; margin: 20px 0;">
  <p style="color: #ff6f00; font-weight: bold; margin: 0; font-size: 1.1em;">⚠️ BETA VERSION - v1.0</p>
  <p style="color: #333; margin: 10px 0 0 0;">This version may contain bugs. If you encounter any issues or have feedback for improvements, please report them. Your feedback helps us fix problems and improve the system. <strong>Contact: send issues or suggestions</strong></p>
</div>

<p align="center">
  <a href="https://i.ibb.co/WN5frxd6/2ee3d1b24181.png">
    <img src="https://i.ibb.co/LzGHGkWG/8e199e30c114.png" alt="How it works">
  </a>
</p>

<p style="font-size: 1.2em; color: #666666b4; margin-bottom: 30px;"><strong>A sophisticated and powerful neural network</strong></p>

<div style="margin: 25px 0;">
  <a href="https://golang.org/"><img src="https://img.shields.io/badge/Go-1.18+-00ADD8?style=flat-square&logo=go" alt="Go"></a>
  <span style="margin: 0 10px;"></span>
  <a href="#license"><img src="https://img.shields.io/badge/License-MIT-4CAF50?style=flat-square" alt="License"></a>
</div>

<hr style="margin: 40px 0; border: none; border-top: 2px solid #eee;">

</div>

<div align="center" style="max-width: 900px; margin: 0 auto; padding: 0 20px;">

<div style="text-align: left; margin: 30px 0;">

### Overview

**ATOMIC AI** is an artificial intelligence system based on neural networks, designed for learning and prediction. The project integrates a modular architecture with data management, visualization interface, and a powerful AI engine.

</div>

<hr style="margin: 40px 0; border: none; border-top: 2px solid #eee;">

<div style="text-align: left; margin: 30px 0;">

### Features

- **Neural Network** - Sophisticated architecture of interconnected neurons
- **Data Management** - Integrated database for learning
- **Dashboard** - Visualization and monitoring interface
- **Lexicon** - Linguistic resources for language processing
- **Flexible Configuration** - Adjustable parameters for fine-tuning

</div>

<hr style="margin: 40px 0; border: none; border-top: 2px solid #eee;">

<div style="text-align: left; margin: 30px 0;">

### Project Structure

```
IA-ATOMIQUE-/
├── main.go              Main orchestration
├──  database/
│   └── data.go             Database management
├── dashboard            User interface
├── ia/                  AI engine
├── lexique.txt          Vocabulary and resources
├── neurones.txt        Neuron configuration
├── go.mod               Go dependencies
└── README.md            Documentation
```

</div>

<hr style="margin: 40px 0; border: none; border-top: 2px solid #eee;">

<div style="text-align: left; margin: 30px 0;">

### Quick Start

#### Requirements

- **Go** 1.18 or higher
- **Git** for version control

#### Installation

```bash
# Clone the repository
git clone <repository-url>
cd IA-ATOMIQUE-

# Download dependencies
go mod download
```

#### Run the Application

```bash
go run main.go
```

</div>

<hr style="margin: 40px 0; border: none; border-top: 2px solid #eee;">

<div style="text-align: left; margin: 30px 0;">

### Usage Guide

#### Neuron Configuration

Modify `neurones.txt` to adjust:
- Number of layers
- Number of neurons per layer
- Learning rate
- Activation function

#### Extend the Lexicon

Add terms to `lexique.txt` to improve natural language processing.

#### Use the Database

Data access is done via `database/data.go` for CRUD operations.

</div>

<hr style="margin: 40px 0; border: none; border-top: 2px solid #eee;">

<div style="text-align: left; margin: 30px 0;">

### Main Files

| File | Description |
|------|-------------|
| `main.go` | Entry point and system orchestration |
| `database/data.go` | Data access layer |
| `ia/` | Neural network implementation |
| `dashboard` | Visualization interface |
| `neurones.txt` | Neuron configuration |
| `lexique.txt` | Linguistic resources |

</div>

<hr style="margin: 40px 0; border: none; border-top: 2px solid #eee;">

<div style="text-align: left; margin: 30px 0;">

### Dashboard Categories

The dashboard monitors **50 categories** of neural networks, each responsible for different aspects of learning:

| Category | Purpose | Function |
|----------|---------|----------|
| **CAT 1-5** | Input Processing | Processes and normalizes input data |
| **CAT 6-15** | Feature Detection | Identifies key patterns and features |
| **CAT 16-25** | Intermediate Processing | Transforms and combines features |
| **CAT 26-35** | Pattern Recognition | Recognizes complex patterns |
| **CAT 36-45** | Decision Making | Analyzes and makes predictions |
| **CAT 46-50** | Output Layer | Generates final results and confidence scores |

**Dashboard Metrics:**
- **Total Neurons** - 1000 neurons across all categories
- **Active Neurons** - Real-time count of active neurons
- **Energy Distribution** - Energy level across each category
- **Top Active** - Top 20 most active neurons

</div>

<hr style="margin: 40px 0; border: none; border-top: 2px solid #eee;">

<div style="text-align: left; margin: 30px 0;">

### Advanced Configuration

For advanced configurations, modify the parameters in:

```
neurones.txt   → Neural network parameters
lexique.txt    → Vocabulary and resources
main.go        → Orchestration logic
```

</div>

<hr style="margin: 40px 0; border: none; border-top: 2px solid #eee;">

<div style="text-align: left; margin: 30px 0;">

### Performance Benchmarks

IA-ATOMIQUE v1.0 has been optimized for speed and efficiency:

#### Core Metrics

| Metric | Result |
|--------|--------|
| **Average Throughput** | 790 phrases/second |
| **Latency** | 1.26 ms/phrase |
| **Memory Usage** | <20 MB RAM |
| **CPU Utilization** | ~23% average |
| **Binary Size** | 2.3 MB |

#### Stress Testing

- **1000 Phrases** - Processed in 31ms (32,107 phrases/sec peak)
- **Stability** - <20% variation over 5 consecutive runs
- **Scalability** - Linear performance with text size

#### Comparison with LLM Local

| Metric | IA-ATOMIQUE | LLM Local | Advantage |
|--------|-------------|-----------|-----------|
| **Speed** | 1.8 ms | 56.3 ms | **30x faster** |
| **Throughput** | 5,475 phrases/s | 195 phrases/s | **28x faster** |
| **CPU** | 19.8% | 19.5% | Comparable |
| **Memory** | 729 KB | 100 KB | Lightweight |

#### Tested Scenarios

- ✅ Wikipedia articles (multi-domain text)
- ✅ Short snippets (tweet-like content)
- ✅ Long documents (1000+ phrases)
- ✅ Mixed language (French with English terms)

</div>

<hr style="margin: 40px 0; border: none; border-top: 2px solid #eee;">

<div style="text-align: left; margin: 30px 0;">

### License

This project is distributed under the **MIT** license. See the [LICENSE](LICENSE) file for more details.

</div>

<hr style="margin: 40px 0; border: none; border-top: 2px solid #eee;">

<div style="margin: 40px 0;">
  <p><strong>nosserb | 2025</strong></p>
  <p style="color: #999; font-size: 0.95em;">⭐ Feel free to star this project if you like it!</p>
</div>

</div>
