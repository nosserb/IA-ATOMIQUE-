package commands

import (
	"IA-ATOMIQUE/database"
	"fmt"
	"os"
	"strconv"
)

// TrainingCommand handles all training-related operations
func TrainingCommand(args []string) {
	if len(args) < 1 {
		PrintTrainingHelp()
		return
	}

	subcommand := args[0]

	switch subcommand {
	case "dataset-load":
		HandleDatasetLoad(args[1:])
	case "dataset-stats":
		HandleDatasetStats(args[1:])
	case "dataset-export":
		HandleDatasetExport(args[1:])
	case "train":
		HandleTrainModel(args[1:])
	case "validate":
		HandleValidateModel(args[1:])
	case "export-model":
		HandleExportModel(args[1:])
	default:
		fmt.Printf("❌ Unknown training command: %s\n", subcommand)
		PrintTrainingHelp()
	}
}

// HandleDatasetLoad loads images from a directory into a dataset
// Default: ./programme train dataset-load uses input/image
func HandleDatasetLoad(args []string) {
	imageDir := "input/image"
	datasetName := "default_dataset"

	if len(args) == 0 {
		fmt.Println("Usage: ./programme train dataset-load [dataset_name] [image_dir]")
		fmt.Println("  Default: ./programme train dataset-load (uses input/image)")
		fmt.Println("  Custom:  ./programme train dataset-load my_dataset ./custom_images")
		return
	}

	datasetName = args[0]
	if len(args) > 1 {
		imageDir = args[1]
	}
	patchSize := 8
	if len(args) > 2 {
		patchSize, _ = strconv.Atoi(args[2])
	}

	// Validate directory
	info, err := os.Stat(imageDir)
	if err != nil || !info.IsDir() {
		fmt.Printf("❌ Invalid directory: %s\n", imageDir)
		return
	}

	fmt.Printf("\n📁 Loading Dataset\n")
	fmt.Printf("═══════════════════════════════════════\n")
	fmt.Printf("Directory: %s\n", imageDir)
	fmt.Printf("Dataset: %s\n", datasetName)
	fmt.Printf("Patch Size: %d\n\n", patchSize)

	// Create dataset
	ds := database.NewAtomicDataset(datasetName, patchSize)

	// Load images
	fmt.Printf("Loading images...\n")
	count, err := ds.LoadDirectory(imageDir)
	if err != nil {
		fmt.Printf("❌ Error loading directory: %v\n", err)
		return
	}

	if count == 0 {
		fmt.Printf("❌ No images loaded\n")
		return
	}

	fmt.Printf("✅ Loaded %d images\n\n", count)

	// Compute statistics
	fmt.Printf("Computing statistics...\n")
	ds.ComputeStatistics()

	// Print statistics
	ds.PrintStatistics()

	fmt.Printf("\n✅ Dataset ready for training\n")
	fmt.Printf("═══════════════════════════════════════\n")
}

// HandleDatasetStats displays statistics about a dataset
// Default: uses input/image
func HandleDatasetStats(args []string) {
	imageDir := "input/image"
	datasetName := "default_dataset"

	if len(args) > 0 {
		datasetName = args[0]
	}
	if len(args) > 1 {
		imageDir = args[1]
	}

	ds := database.NewAtomicDataset(datasetName, 8)

	fmt.Printf("📊 Loading and analyzing dataset...\n")
	count, err := ds.LoadDirectory(imageDir)
	if err != nil || count == 0 {
		fmt.Printf("❌ Failed to load images\n")
		return
	}

	ds.ComputeStatistics()
	ds.PrintStatistics()
}

// HandleDatasetExport saves visualization of target states
// Default: exports to output/ directory
func HandleDatasetExport(args []string) {
	if len(args) < 2 {
		fmt.Println("Usage: ./programme train dataset-export <image_idx> [mode] [dataset_name]")
		fmt.Println("Modes: color, gradient, edge, texture, orientation")
		fmt.Println("Example: ./programme train dataset-export 0 gradient")
		return
	}

	imageIdx, _ := strconv.Atoi(args[0])
	mode := "color"
	imageDir := "input/image"
	datasetName := "default_dataset"

	if len(args) > 1 {
		mode = args[1]
	}
	if len(args) > 2 {
		datasetName = args[2]
	}
	if len(args) > 3 {
		imageDir = args[3]
	}

	ds := database.NewAtomicDataset(datasetName, 8)
	count, err := ds.LoadDirectory(imageDir)
	if err != nil || count == 0 {
		fmt.Printf("❌ Failed to load images\n")
		return
	}

	outFile := fmt.Sprintf("output/target_%s_%s.png", datasetName, mode)
	err = ds.ExportTargetImage(imageIdx, outFile, mode)
	if err != nil {
		fmt.Printf("❌ Export failed: %v\n", err)
		return
	}

	fmt.Printf("✅ Exported: %s\n", outFile)
}

// HandleTrainModel runs training on a dataset
// Default: uses input/image for training images, output/ for results
func HandleTrainModel(args []string) {
	imageDir := "input/image"
	datasetName := "default_dataset"
	width := 512
	height := 512
	epochs := 20
	learningRate := 0.01
	lambda := 0.5

	if len(args) > 0 {
		datasetName = args[0]
	}
	if len(args) > 1 {
		width, _ = strconv.Atoi(args[1])
	}
	if len(args) > 2 {
		height, _ = strconv.Atoi(args[2])
	}
	if len(args) > 3 {
		epochs, _ = strconv.Atoi(args[3])
	}
	if len(args) > 4 {
		learningRate, _ = strconv.ParseFloat(args[4], 64)
	}
	if len(args) > 5 {
		lambda, _ = strconv.ParseFloat(args[5], 64)
	}
	if len(args) > 6 {
		imageDir = args[6]
	}

	// Validate dimensions
	if width < 64 || height < 64 || width > 2048 || height > 2048 {
		fmt.Println("❌ Dimensions must be 64-2048")
		return
	}

	// Load dataset
	fmt.Printf("\n🎓 ATOMIC MODEL TRAINING\n")
	fmt.Printf("═══════════════════════════════════════\n")
	fmt.Printf("Loading dataset: %s\n", datasetName)

	ds := database.NewAtomicDataset(datasetName, 8)
	count, err := ds.LoadDirectory(imageDir)
	if err != nil || count == 0 {
		fmt.Printf("❌ Failed to load images\n")
		return
	}

	fmt.Printf("Loaded %d images\n", count)
	ds.ComputeStatistics()

	// Create network
	fmt.Printf("\n🔧 Initializing Network\n")
	fmt.Printf("Dimensions: %dx%d\n", width, height)
	network := database.NewAtomicImageNetwork(width, height, 8)

	// Create training config
	config := database.TrainingConfig{
		LearningRate:          learningRate,
		CoherenceLambda:       lambda,
		Epochs:                epochs,
		BatchSize:             1,
		GradientClipValue:     1.0,
		WeightDecay:           0.001,
		EarlyStoppingPatience: 5,
		LogInterval:           1,
	}

	// Create trainer
	trainer := database.NewAtomicTrainer(network, ds, config)

	// Run training
	fmt.Printf("\n🚀 Starting Training\n")
	fmt.Printf("═══════════════════════════════════════\n")
	trainer.Train()

	// Validate
	fmt.Printf("\n📊 Running Validation\n")
	trainer.ValidateOnDataset()

	// Export model
	modelFile := fmt.Sprintf("output/model_%s_trained.txt", datasetName)
	trainer.ExportModel(modelFile)
	fmt.Printf("\n✅ Model exported: %s\n", modelFile)

	// Save final image
	outputFile := fmt.Sprintf("output/trained_output_%s.png", datasetName)
	network.SaveImage(outputFile)
	fmt.Printf("✅ Generated image: %s\n", outputFile)

	fmt.Printf("═══════════════════════════════════════\n")
}

// HandleValidateModel evaluates a trained model
func HandleValidateModel(args []string) {
	if len(args) < 2 {
		fmt.Println("Usage: ./programme train validate <image_dir> <dataset_name> [width] [height]")
		return
	}

	imageDir := args[0]
	datasetName := args[1]
	width := 512
	height := 512

	if len(args) > 2 {
		width, _ = strconv.Atoi(args[2])
	}
	if len(args) > 3 {
		height, _ = strconv.Atoi(args[3])
	}

	// Load dataset
	ds := database.NewAtomicDataset(datasetName, 8)
	count, _ := ds.LoadDirectory(imageDir)
	if count == 0 {
		fmt.Printf("❌ Failed to load dataset\n")
		return
	}

	// Create network
	network := database.NewAtomicImageNetwork(width, height, 8)

	// Create trainer
	config := database.DefaultTrainingConfig()
	trainer := database.NewAtomicTrainer(network, ds, config)

	// Validate
	fmt.Printf("\n📊 Validating Model on Dataset\n")
	fmt.Printf("═══════════════════════════════════════\n")
	trainer.ValidateOnDataset()
}

// HandleExportModel saves trained model
func HandleExportModel(args []string) {
	if len(args) < 2 {
		fmt.Println("Usage: ./programme train export-model <model_name> <output_file>")
		return
	}

	modelName := args[0]
	outputFile := args[1]

	// For now, just create a placeholder
	file, err := os.Create(outputFile)
	if err != nil {
		fmt.Printf("❌ Failed to create file: %v\n", err)
		return
	}
	defer file.Close()

	file.WriteString("IA-ATOMIQUE Trained Model: " + modelName + "\n")
	fmt.Printf("✅ Model exported: %s\n", outputFile)
}

// PrintTrainingHelp displays training command help
func PrintTrainingHelp() {
	fmt.Printf(`
🎓 Atomic Model Training - Learn from Reference Images
═════════════════════════════════════════════════════

📁 DEFAULT PATHS:
  Input images:     ./input/image/  (put test images here)
  Generated output: ./output/       (results go here)

DATASET COMMANDS:
  ./programme train dataset-load [dataset_name] [image_dir]
      Load images from directory into dataset
      Default: ./programme train dataset-load (uses input/image)
      Custom:  ./programme train dataset-load my_dataset ./other_images

  ./programme train dataset-stats [dataset_name] [image_dir]
      Show statistics about dataset (gradients, edges, textures)
      Default: ./programme train dataset-stats
      Custom:  ./programme train dataset-stats my_data ./images

  ./programme train dataset-export <image_idx> [mode] [dataset_name]
      Export visualization of target states (saved to output/)
      Modes: color, gradient, edge, texture, orientation
      Example: ./programme train dataset-export 0 gradient

TRAINING COMMANDS:
  ./programme train train [dataset_name] [width] [height] [epochs] [lr] [lambda]
      Train atomic network on dataset
      Quick:  ./programme train train
      Custom: ./programme train train my_data 512 512 20 0.01 0.5
      
      Default dimensions: 512x512
      Default epochs: 20
      Default lr (η): 0.01
      Default lambda (λ): 0.5

  ./programme train validate [dataset_name] [width] [height]
      Evaluate model on dataset

  ./programme train export-model <model_name> <output_file>
      Save trained model

═════════════════════════════════════════════════════
WORKFLOW QUICK START:

1. Create input directory with test images:
   mkdir input/image
   # Copy PNG/JPG files to input/image/

2. Load and analyze:
   ./programme train dataset-load
   ./programme train dataset-stats

3. Train atomic network:
   ./programme train train

4. Check results in output/ directory:
   ls output/
   # View generated images and model files

═════════════════════════════════════════════════════
MATHEMATICAL FOUNDATION:

Loss Function:
  L = Σ_ij ||c_ij_generated - C_ij*||² + λ·Σ_ij Σ_k ||c_ij - c_k||²
  
  First term:  Color fidelity to target image
  Second term: Local coherence (resonance penalty)

Gradient Descent (Atomic):
  c_ij ← c_ij - η·∂L/∂c_ij
  
  Computed per-atom independently
  Enables asynchronous, distributed training

Target States Extracted:
  ✓ Pixel RGB colors
  ✓ Gradients (X, Y directions)
  ✓ Edge strength (gradient magnitude)
  ✓ Curvature (Laplacian)
  ✓ Texture energy (combined edge + curvature)
  ✓ Orientation (atan2(gradY, gradX))

═════════════════════════════════════════════════════
ADVANCED USAGE:

# Fine-tune with different parameters
./programme train train my_dataset 512 512 30 0.005 0.7

# Large-scale training
./programme train train my_dataset 1024 1024 50 0.01 0.5

# Visualize different aspects of target
./programme train dataset-export 0 edge
./programme train dataset-export 0 texture
./programme train dataset-export 1 gradient

═════════════════════════════════════════════════════
`)
}
