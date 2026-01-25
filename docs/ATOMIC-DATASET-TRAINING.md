#  IA-ATOMIQUE: Atomic Dataset Training Guide

## Overview

This guide explains how to create a dataset, extract target atomic states from reference images, and train an atomic neural network to generate realistic images that match your dataset's characteristics.

**Key Innovation**: The training process is fully **atomistic** - each pixel learns independently through local gradient descent, enabling distributed learning without central coordination.

---

## Concept Overview

### From Images to Atomic States

```
Reference Image (JPEG/PNG)
        
    Preprocessing
        
    Extract Target States:
    - Pixel colors (RGB)
    - Gradients (/x, /y)
    - Edge strength (magnitude)
    - Curvature (Laplacian)
    - Texture energy
    - Orientation
        
    Store as TargetAtomicState
        
    Use for Training
```

### Loss Function

The training minimizes:

$$L = \sum_{i,j} \|c_{i,j}^{\text{gen}} - C_{i,j}^*\|^2 + \lambda \sum_{i,j} \sum_{k \in N(p_{i,j})} \|c_{i,j} - c_k\|^2$$

Where:
- **First term**: Color fidelity (how close generated image is to target)
- **Second term**: Coherence penalty (encourages smooth, locally correlated colors)
- **�**: Balance parameter (typically 0.5)

### Atomic Gradient Descent

Each atom updates independently:

$$c_{i,j} \leftarrow c_{i,j} - \eta \cdot \frac{\partial L}{\partial c_{i,j}}$$

Where:
- **�**: Learning rate (typically 0.01)
- **�L**: Gradient computed locally (no global synchronization needed)

---

## Step 1: Prepare Your Dataset

### Directory Structure

```
training_images/
 image1.jpg
 image2.png
 image3.jpg
 image4.png
```

Requirements:
- **Format**: PNG or JPEG
- **Size**: Any resolution (recommended: 256x256 to 1024x1024)
- **Color**: RGB or grayscale
- **Count**: At least 5-10 images for meaningful training

### Example Dataset Scenarios

**Scenario 1: Natural Photos**
```
landscapes/
 sunset_1.jpg
 sunset_2.jpg
 mountain_1.jpg
 ocean_1.jpg
 forest_1.jpg
```

**Scenario 2: Abstract Art**
```
abstract_art/
 painting_1.jpg
 painting_2.jpg
 modern_1.png
 texture_1.jpg
```

**Scenario 3: Synthetic Patterns**
```
patterns/
 grid_1.png
 circles_1.png
 fractals_1.jpg
 noise_1.png
```

---

## Step 2: Load and Analyze Dataset

### Load Images

```bash
./programme train dataset-load ./training_images my_dataset 8
```

Output:
```
 Loading Dataset

Directory: ./training_images
Dataset: my_dataset
Patch Size: 8

Loading images...
 Loaded 5 images

Computing statistics...

 Dataset Statistics: my_dataset

Number of images: 5
Average size: 512.0 � 512.0
Total pixels: 1,310,720

 Image Statistics:
  Avg gradient magnitude: 0.1234
  Avg edge strength: 0.0856
  Avg texture energy: 0.0923

 Color Distribution: 4521 unique colors

```

### View Dataset Statistics

```bash
./programme train dataset-stats ./training_images my_dataset
```

### Understand the Metrics

| Metric | Meaning | Range |
|--------|---------|-------|
| **Avg gradient magnitude** | How much colors change locally | 0.0-1.0 |
| **Avg edge strength** | Strength of edges/boundaries | 0.0-1.0 |
| **Avg texture energy** | Overall texture complexity | 0.0-1.0 |
| **Color distribution** | Diversity of colors used | # unique colors |

**Interpretation**:
- **High gradient**: Image has sharp transitions, detailed boundaries
- **High edge strength**: Many clear edges/objects
- **High texture energy**: Complex, detailed textures

---

## Step 3: Visualize Target States

### Export Visualizations

```bash
# Export color target
./programme train dataset-export ./training_images my_dataset 0 color

# Export gradient visualization
./programme train dataset-export ./training_images my_dataset 0 gradient

# Export edge detection
./programme train dataset-export ./training_images my_dataset 0 edge

# Export texture energy
./programme train dataset-export ./training_images my_dataset 0 texture

# Export orientation field
./programme train dataset-export ./training_images my_dataset 0 orientation
```

### Understanding Visualizations

**Color Mode**
- Shows original image colors
- Red/Green/Blue intensities correspond to RGB channels

**Gradient Mode**
- Bright = strong color changes (high gradient)
- Dark = smooth areas (low gradient)
- Useful to see where the network should learn transitions

**Edge Mode**
- Bright = strong edges/boundaries
- Dark = smooth regions
- Indicates structural elements the network should recognize

**Texture Mode**
- Bright = high texture energy (complex/detailed)
- Dark = smooth, simple regions
- Shows where fine details exist

**Orientation Mode**
- Color encodes direction of dominant gradient
- Red  horizontal edges
- Blue  vertical edges
- Purple  diagonal edges

---

## Step 4: Train the Model

### Basic Training

```bash
./programme train train ./training_images my_dataset 512 512
```

Defaults:
- Epochs: 20
- Learning rate: 0.01
- Lambda: 0.5

### Advanced Training Parameters

```bash
./programme train train ./training_images my_dataset <width> <height> [epochs] [learning_rate] [lambda]
```

**Example: Standard quality**
```bash
./programme train train ./training_images my_dataset 512 512 20 0.01 0.5
```

**Example: High-quality training**
```bash
./programme train train ./training_images my_dataset 512 512 50 0.005 0.7
```

**Example: Fast, coarse learning**
```bash
./programme train train ./training_images my_dataset 256 256 10 0.05 0.3
```

### Understanding Parameters

#### Epochs (Training Iterations)
```
10-20:  Quick training, basic learning
30-50:  Good convergence for most datasets
100+:   Very detailed learning, risk of overfitting
```

#### Learning Rate (�)
```
0.001:  Very conservative, slow convergence
0.01:   Recommended default, stable learning
0.05:   Faster, but may oscillate
0.1+:   Very fast, risk of divergence
```

**Adaptive adjustment**: If loss plateaus, learning rate reduces automatically (0.95x)

#### Coherence Lambda (�)
```
0.1:   Ignore local coherence, focus on color match
0.3:   Light coherence enforcement
0.5:   Balanced (recommended)
0.7:   Strong local structure preservation
1.0+:  Maximize coherence, may lose color fidelity
```

### Training Output

```
 ATOMIC MODEL TRAINING

Loading dataset: my_dataset
Loaded 5 images

 Initializing Network
Dimensions: 512x512

 Starting Training

[Epoch   1/20] Loss: 0.523461 | Best: 0.523461 (ep 0) | Time: 2.34s
[Epoch   2/20] Loss: 0.412834 | Best: 0.412834 (ep 1) | Time: 2.31s
[Epoch   3/20] Loss: 0.328956 | Best: 0.328956 (ep 2) | Time: 2.29s
...
[Epoch  20/20] Loss: 0.089234 | Best: 0.087123 (ep 18) | Time: 2.25s

 Training Complete

Final Loss: 0.089234
Best Loss: 0.087123 (epoch 18)
Total Time: 47s
Weights Updated: 131072


 Model exported: model_my_dataset_trained.txt
 Generated image: trained_output_my_dataset.png
```

### Interpreting Loss Convergence

**Healthy convergence**:
```
Loss: 0.500  0.400  0.350  0.300  0.280  0.270
```

**Slow convergence** (learning rate too low):
```
Loss: 0.500  0.498  0.496  0.495  0.494
```
 Increase learning rate

**Divergence** (learning rate too high):
```
Loss: 0.500  0.600  1.200  NaN
```
 Decrease learning rate

**Plateauing** (no improvement):
```
Loss: 0.200  0.195  0.195  0.195  0.195
```
 Increase epochs or adjust lambda

---

## Step 5: Validate the Model

### Run Validation

```bash
./programme train validate ./training_images my_dataset 512 512
```

Output:
```
 Validating on Dataset

  [1/5] Loss: 0.1234
  [2/5] Loss: 0.0987
  [3/5] Loss: 0.1102
  [4/5] Loss: 0.0945
  [5/5] Loss: 0.1087

Average Validation Loss: 0.1071
```

### Quality Assessment

| Validation Loss | Quality |
|-----------------|---------|
| < 0.1 | Excellent - Very close to training targets |
| 0.1 - 0.2 | Good - Reasonable match |
| 0.2 - 0.5 | Fair - Significant differences |
| > 0.5 | Poor - Major discrepancies |

---

## Complete Workflow Example

### Create a Natural Photos Dataset

```bash
# 1. Prepare directory and images
mkdir training_images
cp ~/Pictures/sunset1.jpg training_images/
cp ~/Pictures/sunset2.jpg training_images/
cp ~/Pictures/ocean1.jpg training_images/
cp ~/Pictures/forest1.jpg training_images/

# 2. Load and analyze
./programme train dataset-load ./training_images landscape_set 8

# 3. View statistics
./programme train dataset-stats ./training_images landscape_set

# 4. Visualize targets
./programme train dataset-export ./training_images landscape_set 0 gradient
./programme train dataset-export ./training_images landscape_set 0 edge

# 5. Train for good quality
./programme train train ./training_images landscape_set 512 512 50 0.01 0.5

# 6. Validate
./programme train validate ./training_images landscape_set 512 512

# 7. Output files
# - trained_output_landscape_set.png (generated image)
# - model_landscape_set_trained.txt (model info)
```

---

## Advanced Topics

### Multi-Scale Training

For large datasets, train at multiple scales:

```bash
# Coarse structure (large patches)
./programme train train ./images dataset 256 256 20 0.05 0.3

# Fine details (smaller patches)
./programme train train ./images dataset 512 512 50 0.01 0.7

# Ultra-fine (very small patches)
./programme train train ./images dataset 1024 1024 100 0.005 0.8
```

### Fine-Tuning Pre-trained Models

```bash
# Very small learning rate for adjustments
./programme train train ./images dataset 512 512 10 0.001 0.6
```

### Handling Different Image Types

**Photography**
```bash
./programme train train ./images photos 512 512 40 0.01 0.5
```

**Art/Paintings**
```bash
./programme train train ./images art 512 512 60 0.008 0.6
```

**Synthetic/Technical**
```bash
./programme train train ./images technical 512 512 30 0.015 0.4
```

**Abstract**
```bash
./programme train train ./images abstract 512 512 70 0.005 0.7
```

---

## Troubleshooting

### Problem: Training very slow
**Solutions**:
- Reduce image dimensions: `256 256` instead of `512 512`
- Increase patch size: `16` instead of `8`
- Reduce epochs: `10-15` instead of `50+`
- Fewer images in dataset

### Problem: Loss not decreasing
**Causes & Solutions**:
- Learning rate too low  increase to 0.05 or 0.1
- Lambda too high  reduce from 0.5 to 0.3
- Dataset too small  add more reference images
- Images too different  use more consistent dataset

### Problem: Loss diverging (becoming larger)
**Solutions**:
- Learning rate too high  reduce to 0.001 or 0.005
- Check image formats (must be PNG/JPEG)
- Verify image quality and colors

### Problem: Generated image doesn't match dataset
**Solutions**:
- Increase epochs (more training time)
- Increase lambda (emphasize coherence matching)
- Use more/better reference images
- Ensure all images in dataset have similar style

---

## Mathematical Details

### Target State Extraction

For each pixel at position $(x, y)$:

1. **Color**: Directly from RGB values, normalized to [0, 1]
2. **Gradients**: 
   - $G_x = \frac{I(x+1,y) - I(x-1,y)}{2}$
   - $G_y = \frac{I(x,y+1) - I(x,y-1)}{2}$
3. **Edge Strength**: $\|G\| = \sqrt{G_x^2 + G_y^2}$
4. **Curvature**: $\nabla^2 I = I_{xx} + I_{yy}$
5. **Texture Energy**: $\alpha \cdot \|G\| + (1-\alpha) \cdot |\nabla^2 I|$ (�=0.7)
6. **Orientation**: $\theta = \arctan\left(\frac{G_y}{G_x}\right)$

### Gradient Computation

```
L/c = 2(c - C*) + ��2�Σ_neighbors(c - c_neighbor)
```

Clipped to prevent explosion:
```
if |L/c| > clip_value: L/c = clip_value � sign(L/c)
```

### Weight Update

Connection weights also update via:

```
w_ij  w_ij � (1 - ��decay) + gamma�coherence
```

Where coherence increases when neighboring atoms align well.

---

## Performance Metrics

### Training Time (Typical)

| Dimensions | Epochs | Time |
|------------|--------|------|
| 256�256 | 20 | 5-10s |
| 512�512 | 20 | 30-45s |
| 1024�1024 | 20 | 2-3 min |

### Memory Usage

```
Atoms: (width/patch) � (height/patch)
512�512 with patch=8  4,096 atoms
Memory per atom: ~500 bytes
Total: ~2 MB
```

### Convergence Speed

```
Dataset size:    5 images  convergence in 20 epochs
Dataset size:   10 images  convergence in 30 epochs
Dataset size:   50 images  convergence in 50-100 epochs
```

---

## Exporting and Using Trained Models

### Model Export

```bash
./programme train export-model landscape_model model.txt
```

Creates `model.txt` with:
- Training parameters
- Final loss values
- Network statistics
- Training duration

### Using Trained Network

The trained network state is automatically saved in the atomic network structure. To use it:

```bash
# Generate new images with trained model
./programme image generate 512 512 100 8 "your prompt"
```

The network will use learned weights and color patterns.

---

## Best Practices

1. **Dataset Consistency**: Use images with similar style/content for better results
2. **Diverse Angles**: Include images from different perspectives
3. **Start Small**: Test with 256�256 before scaling to 1024�1024
4. **Monitor Loss**: Watch for convergence patterns
5. **Validate Regularly**: Check against validation set
6. **Adjust Parameters**: Use metrics to tune learning rate and lambda
7. **Save Checkpoints**: Export models at different epochs

---

## Advanced Configuration

### Extreme Quality (Slow)
```bash
./programme train train ./images dataset 1024 1024 200 0.003 0.8
```
- Very high resolution
- Many epochs for detailed learning
- Low learning rate for stability
- High coherence emphasis

### Fast Training (Low Quality)
```bash
./programme train train ./images dataset 256 256 5 0.1 0.2
```
- Low resolution
- Few epochs
- High learning rate
- Low coherence emphasis

### Balanced (Recommended)
```bash
./programme train train ./images dataset 512 512 30 0.01 0.5
```
- Good resolution
- Reasonable time
- Stable learning
- Balanced coherence

---

**Generated**: January 9, 2026
**Project**: IA-ATOMIQUE v4.0
**Technology**: Atomic Resonance Technology with Distributed Gradient Descent
