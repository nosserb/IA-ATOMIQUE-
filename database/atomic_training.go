// Package database - Atomic Training System
// Implements loss function, gradient descent optimization, and model training
// Tunes atomic weights and states to match target images from dataset

package database

import (
	"fmt"
	"math"
	"os"
	"sync"
	"time"
)

// TrainingConfig contains hyperparameters for training
type TrainingConfig struct {
	LearningRate          float64 // η (eta) - step size for gradient descent
	CoherenceLambda       float64 // λ (lambda) - coherence weight
	Epochs                int     // Number of training passes through dataset
	BatchSize             int     // Images per batch
	GradientClipValue     float64 // Prevent gradient explosion
	WeightDecay           float64 // L2 regularization
	EarlyStoppingPatience int     // Epochs without improvement before stopping
	LogInterval           int     // Print loss every N epochs
}

// TrainingState tracks training progress
type TrainingState struct {
	Epoch                    int
	BatchIdx                 int
	TotalBatches             int
	CurrentLoss              float64
	AvgLoss                  float64
	BestLoss                 float64
	BestEpoch                int
	EpochsWithoutImprovement int
	StartTime                time.Time
	ElapsedTime              time.Duration
	GradientNorm             float64
}

// TrainingMetrics collects training statistics
type TrainingMetrics struct {
	EpochLosses        []float64
	ValidationLosses   []float64
	GradientNorms      []float64
	LearningRates      []float64
	BestValidationLoss float64
	TrainingTime       time.Duration
	AtomsAffected      int
	WeightsUpdated     int
}

// AtomicTrainer manages the training process
type AtomicTrainer struct {
	Network        *AtomicImageNetwork
	Dataset        *AtomicDataset
	Config         TrainingConfig
	State          TrainingState
	Metrics        TrainingMetrics
	GradientBuffer map[int][3]float64 // Gradients per atom (ID → gradient)
	mutex          sync.RWMutex
}

// NewAtomicTrainer creates a new trainer
func NewAtomicTrainer(
	network *AtomicImageNetwork,
	dataset *AtomicDataset,
	config TrainingConfig,
) *AtomicTrainer {
	return &AtomicTrainer{
		Network:        network,
		Dataset:        dataset,
		Config:         config,
		State:          TrainingState{},
		Metrics:        TrainingMetrics{},
		GradientBuffer: make(map[int][3]float64),
	}
}

// ComputeLoss calculates the loss between generated and target states
// L = Σ_ij ||c_ij_generated - C_ij*||^2 + λ Σ_ij Σ_k |c_ij - c_k|^2
func (trainer *AtomicTrainer) ComputeLoss(imageIdx int) float64 {
	targets, err := trainer.Dataset.GetTargetStates(imageIdx)
	if err != nil {
		return math.Inf(1)
	}

	gridHeight := len(trainer.Network.Atoms)
	gridWidth := 0
	if gridHeight > 0 {
		gridWidth = len(trainer.Network.Atoms[0])
	}

	loss := 0.0

	// First term: color fidelity loss
	for y := 0; y < gridHeight && y < len(targets)/gridWidth; y++ {
		for x := 0; x < gridWidth && x*gridWidth+y < len(targets); x++ {
			atom := &trainer.Network.Atoms[y][x]
			targetIdx := y*gridWidth + x
			if targetIdx >= len(targets) {
				continue
			}

			target := targets[targetIdx]

			// Color distance
			for i := 0; i < 3; i++ {
				diff := atom.Color[i] - target.ColorTarget[i]
				loss += diff * diff
			}

			// Gradient distance (micro-pattern matching)
			loss += math.Pow(atom.Features["orientation"]-target.Orientation, 2) * 0.1
		}
	}

	// Second term: coherence loss (local resonance penalty)
	coherenceLoss := 0.0
	for y := 0; y < gridHeight; y++ {
		for x := 0; x < gridWidth; x++ {
			atom := &trainer.Network.Atoms[y][x]

			// Penalize color difference with neighbors
			directions := [][2]int{
				{-1, -1}, {0, -1}, {1, -1},
				{-1, 0}, {1, 0},
				{-1, 1}, {0, 1}, {1, 1},
			}

			for _, dir := range directions {
				ny, nx := y+dir[0], x+dir[1]
				if nx >= 0 && nx < gridWidth && ny >= 0 && ny < gridHeight {
					neighbor := &trainer.Network.Atoms[ny][nx]
					for i := 0; i < 3; i++ {
						diff := atom.Color[i] - neighbor.Color[i]
						coherenceLoss += diff * diff
					}
				}
			}
		}
	}

	// Combine losses
	totalLoss := loss + trainer.Config.CoherenceLambda*coherenceLoss

	return totalLoss / float64(gridHeight*gridWidth)
}

// ComputeGradients calculates gradients for each atom
// ∂L/∂c_ij = 2(c_ij_generated - C_ij*) + λ * (coherence terms)
func (trainer *AtomicTrainer) ComputeGradients(imageIdx int) {
	trainer.GradientBuffer = make(map[int][3]float64)

	targets, err := trainer.Dataset.GetTargetStates(imageIdx)
	if err != nil {
		return
	}

	gridHeight := len(trainer.Network.Atoms)
	gridWidth := 0
	if gridHeight > 0 {
		gridWidth = len(trainer.Network.Atoms[0])
	}

	// Color fidelity gradients
	for y := 0; y < gridHeight; y++ {
		for x := 0; x < gridWidth; x++ {
			atom := &trainer.Network.Atoms[y][x]
			atomID := y*gridWidth + x

			targetIdx := y*gridWidth + x
			if targetIdx >= len(targets) {
				continue
			}

			target := targets[targetIdx]
			grad := [3]float64{}

			// Color gradient: ∂L/∂c = 2(c - target)
			for i := 0; i < 3; i++ {
				grad[i] = 2.0 * (atom.Color[i] - target.ColorTarget[i])
			}

			// Coherence gradient from neighbors
			directions := [][2]int{
				{-1, -1}, {0, -1}, {1, -1},
				{-1, 0}, {1, 0},
				{-1, 1}, {0, 1}, {1, 1},
			}

			for _, dir := range directions {
				ny, nx := y+dir[0], x+dir[1]
				if nx >= 0 && nx < gridWidth && ny >= 0 && ny < gridHeight {
					neighbor := &trainer.Network.Atoms[ny][nx]
					for i := 0; i < 3; i++ {
						colorDiff := atom.Color[i] - neighbor.Color[i]
						grad[i] += 2.0 * trainer.Config.CoherenceLambda * colorDiff
					}
				}
			}

			// Clip gradient
			for i := 0; i < 3; i++ {
				if grad[i] > trainer.Config.GradientClipValue {
					grad[i] = trainer.Config.GradientClipValue
				} else if grad[i] < -trainer.Config.GradientClipValue {
					grad[i] = -trainer.Config.GradientClipValue
				}
			}

			trainer.GradientBuffer[atomID] = grad
		}
	}
}

// ApplyGradients updates atomic states using computed gradients
// c_ij ← c_ij - η * ∂L/∂c_ij
func (trainer *AtomicTrainer) ApplyGradients() {
	gridHeight := len(trainer.Network.Atoms)
	gridWidth := 0
	if gridHeight > 0 {
		gridWidth = len(trainer.Network.Atoms[0])
	}

	updatedCount := 0

	for y := 0; y < gridHeight; y++ {
		for x := 0; x < gridWidth; x++ {
			atom := &trainer.Network.Atoms[y][x]
			atomID := y*gridWidth + x

			if grad, ok := trainer.GradientBuffer[atomID]; ok {
				// Gradient descent update
				for i := 0; i < 3; i++ {
					atom.Color[i] -= trainer.Config.LearningRate * grad[i]

					// Clamp to [0, 1]
					if atom.Color[i] > 1.0 {
						atom.Color[i] = 1.0
					} else if atom.Color[i] < 0.0 {
						atom.Color[i] = 0.0
					}
				}

				// Weight decay (L2 regularization)
				for neighborID := range atom.ConnectionWeights {
					atom.ConnectionWeights[neighborID] *=
						(1.0 - trainer.Config.WeightDecay*trainer.Config.LearningRate)

					if atom.ConnectionWeights[neighborID] > 1.0 {
						atom.ConnectionWeights[neighborID] = 1.0
					} else if atom.ConnectionWeights[neighborID] < 0.0 {
						atom.ConnectionWeights[neighborID] = 0.0
					}
				}

				updatedCount++
			}
		}
	}

	trainer.Metrics.WeightsUpdated += updatedCount
}

// TrainEpoch runs one epoch through the entire dataset
func (trainer *AtomicTrainer) TrainEpoch(epochNum int) float64 {
	trainer.State.Epoch = epochNum
	totalLoss := 0.0
	datasetSize := trainer.Dataset.GetDatasetSize()

	if datasetSize == 0 {
		return 0.0
	}

	for imgIdx := 0; imgIdx < datasetSize; imgIdx++ {
		// Compute loss
		loss := trainer.ComputeLoss(imgIdx)
		totalLoss += loss

		// Compute gradients
		trainer.ComputeGradients(imgIdx)

		// Apply gradients
		trainer.ApplyGradients()

		// Resonance iteration to propagate changes
		trainer.Network.IterateGeneration()

		// Progress
		if (imgIdx+1)%10 == 0 && trainer.Config.LogInterval > 0 {
			fmt.Printf("  [%d/%d] Loss: %.6f\n", imgIdx+1, datasetSize, loss)
		}
	}

	avgLoss := totalLoss / float64(datasetSize)
	return avgLoss
}

// Train runs the complete training loop
func (trainer *AtomicTrainer) Train() error {
	trainer.State.StartTime = time.Now()
	trainer.Metrics.EpochLosses = make([]float64, 0)
	trainer.Metrics.GradientNorms = make([]float64, 0)

	if trainer.Dataset.GetDatasetSize() == 0 {
		return fmt.Errorf("dataset is empty")
	}

	fmt.Printf("\n🎓 Starting Atomic Training\n")
	fmt.Printf("═══════════════════════════════════════\n")
	fmt.Printf("Epochs: %d\n", trainer.Config.Epochs)
	fmt.Printf("Learning Rate: %.6f\n", trainer.Config.LearningRate)
	fmt.Printf("Coherence Lambda: %.4f\n", trainer.Config.CoherenceLambda)
	fmt.Printf("Dataset Size: %d images\n", trainer.Dataset.GetDatasetSize())
	fmt.Printf("═══════════════════════════════════════\n\n")

	bestLoss := math.Inf(1)
	trainer.State.EpochsWithoutImprovement = 0

	for epoch := 0; epoch < trainer.Config.Epochs; epoch++ {
		epochStartTime := time.Now()

		// Train one epoch
		avgLoss := trainer.TrainEpoch(epoch)
		trainer.Metrics.EpochLosses = append(trainer.Metrics.EpochLosses, avgLoss)
		trainer.State.AvgLoss = avgLoss

		// Check for improvement
		if avgLoss < bestLoss {
			bestLoss = avgLoss
			trainer.State.BestLoss = bestLoss
			trainer.State.BestEpoch = epoch
			trainer.State.EpochsWithoutImprovement = 0
		} else {
			trainer.State.EpochsWithoutImprovement++
		}

		// Log progress
		if (epoch+1)%trainer.Config.LogInterval == 0 || epoch == 0 {
			elapsed := time.Since(epochStartTime)
			fmt.Printf("[Epoch %3d/%d] Loss: %.6f | Best: %.6f (ep %d) | Time: %v\n",
				epoch+1,
				trainer.Config.Epochs,
				avgLoss,
				bestLoss,
				trainer.State.BestEpoch,
				elapsed,
			)
		}

		// Early stopping
		if trainer.Config.EarlyStoppingPatience > 0 &&
			trainer.State.EpochsWithoutImprovement >= trainer.Config.EarlyStoppingPatience {
			fmt.Printf("\n⚠️  Early stopping at epoch %d (no improvement for %d epochs)\n",
				epoch+1, trainer.Config.EarlyStoppingPatience)
			break
		}
	}

	trainer.State.ElapsedTime = time.Since(trainer.State.StartTime)

	fmt.Printf("\n✅ Training Complete\n")
	fmt.Printf("═══════════════════════════════════════\n")
	fmt.Printf("Final Loss: %.6f\n", trainer.State.AvgLoss)
	fmt.Printf("Best Loss: %.6f (epoch %d)\n", trainer.State.BestLoss, trainer.State.BestEpoch)
	fmt.Printf("Total Time: %v\n", trainer.State.ElapsedTime)
	fmt.Printf("Weights Updated: %d\n", trainer.Metrics.WeightsUpdated)
	fmt.Printf("═══════════════════════════════════════\n")

	return nil
}

// ValidateOnDataset evaluates model on entire dataset
func (trainer *AtomicTrainer) ValidateOnDataset() float64 {
	totalLoss := 0.0
	datasetSize := trainer.Dataset.GetDatasetSize()

	if datasetSize == 0 {
		return 0.0
	}

	fmt.Printf("\n📊 Validating on Dataset\n")

	for imgIdx := 0; imgIdx < datasetSize; imgIdx++ {
		loss := trainer.ComputeLoss(imgIdx)
		totalLoss += loss

		if (imgIdx+1)%5 == 0 {
			fmt.Printf("  [%d/%d] Loss: %.6f\n", imgIdx+1, datasetSize, loss)
		}
	}

	avgLoss := totalLoss / float64(datasetSize)
	fmt.Printf("\nAverage Validation Loss: %.6f\n", avgLoss)

	return avgLoss
}

// PrintTrainingStatus displays current training state
func (trainer *AtomicTrainer) PrintTrainingStatus() {
	fmt.Printf("\n📈 Training Status\n")
	fmt.Printf("═══════════════════════════════════════\n")
	fmt.Printf("Epoch: %d\n", trainer.State.Epoch)
	fmt.Printf("Current Loss: %.6f\n", trainer.State.CurrentLoss)
	fmt.Printf("Average Loss: %.6f\n", trainer.State.AvgLoss)
	fmt.Printf("Best Loss: %.6f (epoch %d)\n", trainer.State.BestLoss, trainer.State.BestEpoch)
	fmt.Printf("Epochs Without Improvement: %d\n", trainer.State.EpochsWithoutImprovement)
	fmt.Printf("Elapsed Time: %v\n", trainer.State.ElapsedTime)
	fmt.Printf("═══════════════════════════════════════\n")
}

// ExportModel saves the trained network state
func (trainer *AtomicTrainer) ExportModel(filePath string) error {
	// For now, save network statistics
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	file.WriteString("IA-ATOMIQUE Trained Model\n")
	file.WriteString("═══════════════════════════════════════\n")
	file.WriteString(fmt.Sprintf("Dataset: %s\n", trainer.Dataset.DatasetName))
	file.WriteString(fmt.Sprintf("Epochs Trained: %d\n", trainer.State.Epoch))
	file.WriteString(fmt.Sprintf("Final Loss: %.6f\n", trainer.State.AvgLoss))
	file.WriteString(fmt.Sprintf("Best Loss: %.6f\n", trainer.State.BestLoss))
	file.WriteString(fmt.Sprintf("Training Time: %v\n", trainer.State.ElapsedTime))

	// Add network statistics
	stats := trainer.Network.GetNetworkStats()
	file.WriteString("\nNetwork Statistics:\n")
	for key, val := range stats {
		file.WriteString(fmt.Sprintf("  %s: %.6f\n", key, val))
	}

	return nil
}

// DefaultTrainingConfig returns sensible defaults
func DefaultTrainingConfig() TrainingConfig {
	return TrainingConfig{
		LearningRate:          0.01,  // η - conservative learning rate
		CoherenceLambda:       0.5,   // λ - balance between fidelity and coherence
		Epochs:                20,    // Reasonable default
		BatchSize:             4,     // Small batches
		GradientClipValue:     1.0,   // Prevent explosion
		WeightDecay:           0.001, // Light regularization
		EarlyStoppingPatience: 5,     // Stop if no improvement for 5 epochs
		LogInterval:           1,     // Log every epoch
	}
}

// AdaptiveLearningRate adjusts learning rate based on progress
func (trainer *AtomicTrainer) AdaptiveLearningRate() {
	if len(trainer.Metrics.EpochLosses) < 5 {
		return
	}

	recent := trainer.Metrics.EpochLosses[len(trainer.Metrics.EpochLosses)-5:]
	avgRecent := 0.0
	for _, loss := range recent {
		avgRecent += loss
	}
	avgRecent /= float64(len(recent))

	// If loss plateaus, reduce learning rate
	if trainer.State.EpochsWithoutImprovement > 2 {
		trainer.Config.LearningRate *= 0.95
		fmt.Printf("📉 Adjusting learning rate to %.6f\n", trainer.Config.LearningRate)
	}
}
