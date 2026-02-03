package commands

import (
	"github.com/nosserb/IA-ATOMIQUE-/database"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"strconv"
	"strings"
)

// PatternCommand handles pattern emergence operations
func PatternCommand(args []string) {
	if len(args) < 1 {
		PrintPatternHelp()
		return
	}

	subcommand := args[0]

	switch subcommand {
	case "create":
		HandlePatternCreate(args[1:])
	case "diffuse":
		HandlePatternDiffuse(args[1:])
	case "reinforce":
		HandlePatternReinforce(args[1:])
	case "seed":
		HandlePatternSeed(args[1:])
	case "emerge":
		HandlePatternEmerge(args[1:])
	case "extract":
		HandlePatternExtract(args[1:])
	case "visualize":
		HandlePatternVisualize(args[1:])
	case "index":
		HandlePatternIndex(args[1:])
	case "list":
		HandlePatternList(args[1:])
	case "info":
		HandlePatternInfo(args[1:])
	case "stats":
		HandlePatternStats(args[1:])
	case "search":
		HandlePatternSearch(args[1:])
	default:
		fmt.Printf("❌ Unknown pattern command: %s\n", subcommand)
		PrintPatternHelp()
	}
}

// HandlePatternCreate initializes a new pattern emergence engine
// Usage: ./programme pattern create <width> <height> [seed_image]
func HandlePatternCreate(args []string) {
	if len(args) < 2 {
		fmt.Println("Usage: ./programme pattern create <width> <height> [seed_image.png]")
		fmt.Println("Example: ./programme pattern create 512 512")
		return
	}

	width, _ := strconv.Atoi(args[0])
	height, _ := strconv.Atoi(args[1])

	if width < 64 || height < 64 || width > 2048 || height > 2048 {
		fmt.Println("❌ Dimensions must be 64-2048")
		return
	}

	fmt.Printf("\n  Creating Pattern Emergence Engine\n")
	fmt.Printf("═══════════════════════════════════════\n")
	fmt.Printf("Dimensions: %dx%d\n", width, height)

	// Create engine
	engine := database.NewPatternEmergenceEngine(width, height)

	fmt.Printf("✓ Uniform pixel connections initialized\n")
	fmt.Printf("✓ Diffusion parameters: α=%.3f, β=%.3f\n",
		engine.DiffusionAlpha, engine.VelocityDamping)
	fmt.Printf("✓ Reinforcement γ=%.4f\n", engine.ReinforcementGamma)

	// Load seed image if provided
	if len(args) > 2 {
		seedPath := args[2]
		sampleDensity := 0.1 // 10% of pixels as seeds
		if len(args) > 3 {
			sampleDensity, _ = strconv.ParseFloat(args[3], 64)
		}

		fmt.Printf("\n Loading seed image: %s\n", seedPath)
		file, err := os.Open(seedPath)
		if err != nil {
			fmt.Printf(" Failed to load image: %v\n", err)
			return
		}
		defer file.Close()

		img, _, err := image.Decode(file)
		if err != nil {
			fmt.Printf("  Failed to decode image: %v\n", err)
			return
		}

		engine.AddSeedsFromImage(img, sampleDensity)
	}

	// Save initial state
	outFile := "output/pattern_initial.png"
	engine.SaveImage(outFile)
	fmt.Printf(" Initial state saved: %s\n", outFile)

	fmt.Printf("═══════════════════════════════════════\n")
	fmt.Printf(" Engine ready for pattern emergence\n\n")

	// Store engine globally for next commands
	// (In real app, use session management)
	_ = engine
}

// HandlePatternDiffuse runs diffusion iterations
// Usage: ./programme pattern diffuse <iterations>
func HandlePatternDiffuse(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: ./programme pattern diffuse <iterations> [visualize_interval]")
		fmt.Println("Example: ./programme pattern diffuse 100 10")
		return
	}

	iterations, _ := strconv.Atoi(args[0])
	vizInterval := 50
	if len(args) > 1 {
		vizInterval, _ = strconv.Atoi(args[1])
	}

	fmt.Printf("\n Pixel Diffusion - P_ij(t+1) = P_ij(t) + α·Σ W·(P_kl - P_ij)\n")
	fmt.Printf("═══════════════════════════════════════\n")
	fmt.Printf("Running %d diffusion steps\n", iterations)
	fmt.Printf("Visualization interval: %d steps\n\n", vizInterval)

	// Create a demo engine for illustration
	engine := database.NewPatternEmergenceEngine(256, 256)

	// Add some seed points for demonstration
	engine.AddSeedPoint(50, 50, 255, 0, 0)   // Red seed
	engine.AddSeedPoint(200, 200, 0, 255, 0) // Green seed
	engine.AddSeedPoint(100, 150, 0, 0, 255) // Blue seed

	fmt.Printf(" Added 3 seed points to anchor waves\n\n")

	for iter := 0; iter < iterations; iter++ {
		engine.DiffuseStep()

		if (iter+1)%vizInterval == 0 {
			outFile := fmt.Sprintf("output/pattern_diffuse_%03d.png", iter+1)
			engine.SaveImage(outFile)
			fmt.Printf("✓ Iteration %d: Loss %.8f | Saved: %s\n",
				iter+1, engine.AverageLoss, outFile)
		}
	}

	fmt.Printf("\n Diffusion complete\n")
	engine.PrintStatistics()
}

// HandlePatternReinforce runs connection reinforcement
// Usage: ./programme pattern reinforce <iterations>
func HandlePatternReinforce(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: ./programme pattern reinforce <iterations>")
		fmt.Println("Example: ./programme pattern reinforce 20")
		return
	}

	iterations, _ := strconv.Atoi(args[0])

	fmt.Printf("\n Connection Reinforcement - W_ij;kl(t+1) = W_ij;kl + γ·exp(-||ΔP||²)\n")
	fmt.Printf("═══════════════════════════════════════\n")
	fmt.Printf("Strengthening connections between similar pixels\n")
	fmt.Printf("Iterations: %d\n\n", iterations)

	engine := database.NewPatternEmergenceEngine(256, 256)

	// Add seed points
	engine.AddSeedPoint(50, 50, 255, 0, 0)
	engine.AddSeedPoint(200, 200, 0, 255, 0)

	// Diffuse first to create patterns
	fmt.Printf("Phase 1: Initial diffusion...\n")
	for i := 0; i < 30; i++ {
		engine.DiffuseStep()
	}

	fmt.Printf("Phase 2: Reinforcing connections...\n")
	for i := 0; i < iterations; i++ {
		engine.ReinforceConnections()
		if (i+1)%5 == 0 {
			fmt.Printf("  ✓ Reinforcement iteration %d/%d\n", i+1, iterations)
		}
	}

	outFile := "output/pattern_reinforced.png"
	engine.SaveImage(outFile)
	fmt.Printf("\n✅ Reinforcement complete\n")
	fmt.Printf("   Saved: %s\n\n", outFile)
}

// HandlePatternSeed manages seed points
// Usage: ./programme pattern seed add <x> <y> <r> <g> <b>
//
//	./programme pattern seed load <image.png> <density>
func HandlePatternSeed(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: ./programme pattern seed <add|load> ...")
		fmt.Println("  add <x> <y> <r> <g> <b>    - Add single seed point")
		fmt.Println("  load <image> <density>      - Load seeds from image (0.0-1.0)")
		return
	}

	action := args[0]
	engine := database.NewPatternEmergenceEngine(512, 512)

	if action == "add" && len(args) >= 6 {
		x, _ := strconv.Atoi(args[1])
		y, _ := strconv.Atoi(args[2])
		r, _ := strconv.Atoi(args[3])
		g, _ := strconv.Atoi(args[4])
		b, _ := strconv.Atoi(args[5])

		fmt.Printf("\n📌 Adding Seed Point\n")
		fmt.Printf("═══════════════════════════════════════\n")
		fmt.Printf("Position: (%d, %d)\n", x, y)
		fmt.Printf("Color: RGB(%d, %d, %d)\n", r, g, b)

		engine.AddSeedPoint(x, y, float64(r), float64(g), float64(b))
		fmt.Printf("✅ Seed point added\n\n")

	} else if action == "load" && len(args) >= 3 {
		imagePath := args[1]
		density, _ := strconv.ParseFloat(args[2], 64)

		fmt.Printf("\n📌 Loading Seeds from Image\n")
		fmt.Printf("═══════════════════════════════════════\n")
		fmt.Printf("Image: %s\n", imagePath)
		fmt.Printf("Density: %.1f%%\n\n", density*100)

		file, err := os.Open(imagePath)
		if err != nil {
			fmt.Printf("❌ Failed to load image: %v\n", err)
			return
		}
		defer file.Close()

		img, _, err := image.Decode(file)
		if err != nil {
			fmt.Printf("❌ Failed to decode image: %v\n", err)
			return
		}

		engine.AddSeedsFromImage(img, density)
		outFile := "output/pattern_seeds_loaded.png"
		engine.SaveImage(outFile)
		fmt.Printf("✅ Seeds loaded and saved: %s\n\n", outFile)
	}
}

// HandlePatternEmerge runs full emergence cycle
// Usage: ./programme pattern emerge <width> <height> <iterations> [seed_image] [density]
func HandlePatternEmerge(args []string) {
	if len(args) < 3 {
		fmt.Println("Usage: ./programme pattern emerge <w> <h> <iterations> [seed_image] [density]")
		fmt.Println("Example: ./programme pattern emerge 512 512 200 input/image/test.png 0.15")
		return
	}

	width, _ := strconv.Atoi(args[0])
	height, _ := strconv.Atoi(args[1])
	iterations, _ := strconv.Atoi(args[2])
	seedImage := ""
	seedDensity := 0.1

	if len(args) > 3 {
		seedImage = args[3]
	}
	if len(args) > 4 {
		seedDensity, _ = strconv.ParseFloat(args[4], 64)
	}

	fmt.Printf("\n✨ PATTERN EMERGENCE CYCLE\n")
	fmt.Printf("═══════════════════════════════════════\n")
	fmt.Printf("Transforming abstract waves → recognizable structures\n\n")

	fmt.Printf("🎨 Phase 1: Initialize Engine (%dx%d)\n", width, height)
	engine := database.NewPatternEmergenceEngine(width, height)
	fmt.Printf("   ✓ Created pixel grid with uniform connections\n")

	// Load seed image if provided
	if seedImage != "" && seedImage != "-" {
		fmt.Printf("\n📌 Phase 2: Load Reference Seeds (density: %.1f%%)\n", seedDensity*100)
		file, err := os.Open(seedImage)
		if err != nil {
			fmt.Printf("   ❌ Could not load seed image\n")
		} else {
			defer file.Close()
			img, _, err := image.Decode(file)
			if err != nil {
				fmt.Printf("   ❌ Could not decode image\n")
			} else {
				engine.AddSeedsFromImage(img, seedDensity)
				outFile := "output/pattern_seeds_visual.png"
				engine.SaveImage(outFile)
				fmt.Printf("   ✓ Seeds saved to visualization\n")
			}
		}
	}

	// Run emergence
	fmt.Printf("\n🌊 Phase 3: Pixel Diffusion (%d iterations)\n", iterations)
	fmt.Printf("   P_ij(t+1) = P_ij(t) + α·Σ W·(neighbor_colors - P_ij)\n")

	// Save samples at intervals
	sampleInterval := iterations / 4
	if sampleInterval < 1 {
		sampleInterval = 1
	}

	for step := 0; step < iterations; step++ {
		engine.DiffuseStep()

		// Reinforce every 10 steps
		if step%10 == 0 && step > 0 {
			engine.ReinforceConnections()
		}

		// Visualize progress
		if (step+1)%sampleInterval == 0 || step == 0 {
			outFile := fmt.Sprintf("output/pattern_emerge_%04d.png", step+1)
			engine.SaveImage(outFile)
			fmt.Printf("   ✓ Iter %d: Loss %.8f\n", step+1, engine.AverageLoss)
		}
	}

	fmt.Printf("\n💪 Phase 4: Connection Reinforcement (10 cycles)\n")
	for i := 0; i < 10; i++ {
		engine.ReinforceConnections()
	}
	fmt.Printf("   ✓ Weights strengthened for stable patterns\n")

	// Final save
	finalFile := "output/pattern_final_emerged.png"
	engine.SaveImage(finalFile)
	fmt.Printf("\n✅ EMERGENCE COMPLETE\n")
	fmt.Printf("═══════════════════════════════════════\n")
	engine.PrintStatistics()
	fmt.Printf("Final image: %s\n\n", finalFile)
}

// HandlePatternExtract runs full emergence cycle but keeps ONLY final image
// Cleans up intermediate files automatically
// Usage: ./programme pattern extract <w> <h> <iterations> [seed_image] [density]
func HandlePatternExtract(args []string) {
	if len(args) < 3 {
		fmt.Println("Usage: ./programme pattern extract <w> <h> <iterations> [seed_image] [density]")
		fmt.Println("Example: ./programme pattern extract 512 512 200 input/image/face.png 0.15")
		fmt.Println("\nThis command runs the FULL emergence pipeline but keeps only the final image.")
		return
	}

	width, _ := strconv.Atoi(args[0])
	height, _ := strconv.Atoi(args[1])
	iterations, _ := strconv.Atoi(args[2])
	seedImage := ""
	seedDensity := 0.1

	if len(args) > 3 {
		seedImage = args[3]
	}
	if len(args) > 4 {
		seedDensity, _ = strconv.ParseFloat(args[4], 64)
	}

	fmt.Printf("\n⚡ FAST PATTERN EXTRACTION (Final Image Only)\n")
	fmt.Printf("═══════════════════════════════════════════════════\n")
	fmt.Printf("Running complete emergence pipeline...\n")
	fmt.Printf("Intermediate files: WILL BE DELETED\n")
	fmt.Printf("Final result: KEPT in output/\n\n")

	// Create engine
	fmt.Printf("🎨 Phase 1: Initialize Engine (%dx%d)\n", width, height)
	engine := database.NewPatternEmergenceEngine(width, height)
	fmt.Printf("   ✓ Created pixel grid with uniform connections\n")

	// Load seed image if provided
	if seedImage != "" && seedImage != "-" {
		fmt.Printf("\n📌 Phase 2: Load Reference Seeds (density: %.1f%%)\n", seedDensity*100)
		file, err := os.Open(seedImage)
		if err != nil {
			fmt.Printf("   ⚠ Could not load seed image (continuing without seeds)\n")
		} else {
			defer file.Close()
			img, _, err := image.Decode(file)
			if err != nil {
				fmt.Printf("   ⚠ Could not decode image (continuing without seeds)\n")
			} else {
				engine.AddSeedsFromImage(img, seedDensity)
				fmt.Printf("   ✓ Seeds loaded\n")
			}
		}
	}

	// Run emergence WITHOUT saving intermediate files
	fmt.Printf("\n🌊 Phase 3: Pixel Diffusion (%d iterations - Silent Mode)\n", iterations)
	fmt.Printf("   Processing without intermediate saves...\n")

	lastProgressTime := 0
	for step := 0; step < iterations; step++ {
		engine.DiffuseStep()

		// Reinforce every 10 steps
		if step%10 == 0 && step > 0 {
			engine.ReinforceConnections()
		}

		// Show progress every 25% of total
		progressInterval := iterations / 4
		if progressInterval < 1 {
			progressInterval = 1
		}

		if (step+1)%progressInterval == 0 && step+1 != lastProgressTime {
			percent := (step + 1) * 100 / iterations
			fmt.Printf("   ✓ %3d%% complete (Iter %d, Loss: %.8f)\n",
				percent, step+1, engine.AverageLoss)
			lastProgressTime = step + 1
		}
	}

	// Final reinforcement
	fmt.Printf("\n💪 Phase 4: Connection Reinforcement (10 cycles)\n")
	for i := 0; i < 10; i++ {
		engine.ReinforceConnections()
	}
	fmt.Printf("   ✓ Weights strengthened\n")

	// Generate unique filename
	outputFile := fmt.Sprintf("output/pattern_extracted_%dx%d_%diter.png", width, height, iterations)
	if seedImage != "" && seedImage != "-" {
		outputFile = fmt.Sprintf("output/pattern_%s_%diter.png",
			getFilenameWithoutExt(seedImage), iterations)
	}

	// Save ONLY final image
	fmt.Printf("\n💾 Saving final image...\n")
	engine.SaveImage(outputFile)
	fmt.Printf("   ✓ Saved: %s\n", outputFile)

	// Clean up intermediate files if they exist
	fmt.Printf("\n🧹 Cleaning up intermediate files...\n")
	cleanupCount := 0
	for i := 1; i <= iterations; i++ {
		tempFile := fmt.Sprintf("output/pattern_emerge_%04d.png", i)
		if _, err := os.Stat(tempFile); err == nil {
			os.Remove(tempFile)
			cleanupCount++
		}
	}

	// Also clean up seed visualization
	seedVizFile := "output/pattern_seeds_visual.png"
	if _, err := os.Stat(seedVizFile); err == nil {
		os.Remove(seedVizFile)
		cleanupCount++
	}

	fmt.Printf("   ✓ Deleted %d intermediate file(s)\n", cleanupCount)

	// Final stats
	fmt.Printf("\n✅ EXTRACTION COMPLETE\n")
	fmt.Printf("═══════════════════════════════════════════════════\n")
	fmt.Printf("📊 Statistics:\n")
	fmt.Printf("   Iterations:     %d\n", iterations)
	fmt.Printf("   Seed Points:    %d\n", len(engine.Seeds))
	fmt.Printf("   Final Loss:     %.8f\n", engine.AverageLoss)
	fmt.Printf("   Output File:    %s\n", outputFile)
	fmt.Printf("═══════════════════════════════════════════════════\n\n")
}

// Helper function to get filename without extension
func getFilenameWithoutExt(filepath string) string {
	// Get just the filename from path
	name := filepath
	if idx := findLastIndex(name, '/'); idx >= 0 {
		name = name[idx+1:]
	}
	// Remove extension
	if idx := findLastIndex(name, '.'); idx >= 0 {
		name = name[:idx]
	}
	return name
}

// Helper function to find last index of a string
func findLastIndex(s string, c byte) int {
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == c {
			return i
		}
	}
	return -1
}

// HandlePatternVisualize shows current pattern state
func HandlePatternVisualize(args []string) {
	fmt.Printf("\n🎨 Pattern Visualization Guide\n")
	fmt.Printf("═══════════════════════════════════════\n")
	fmt.Printf("Check output/ directory for:\n\n")
	fmt.Printf("  pattern_initial.png         - Starting state (uniform gray)\n")
	fmt.Printf("  pattern_diffuse_*.png       - Diffusion progression\n")
	fmt.Printf("  pattern_reinforced.png      - After weight strengthening\n")
	fmt.Printf("  pattern_seeds_*.png         - Seed points visualization\n")
	fmt.Printf("  pattern_emerge_*.png        - Full emergence sequence\n")
	fmt.Printf("  pattern_final_emerged.png   - Final recognizable pattern\n\n")
}

// PrintPatternHelp displays pattern command help
func PrintPatternHelp() {
	fmt.Printf(`
✨ Pattern Emergence - From Abstract Waves to Recognizable Images
═════════════════════════════════════════════════════════════════

Mathematical Framework:

1️⃣ LOCAL DIFFUSION (P_ij models pixels):
   P_ij(t+1) = P_ij(t) + α·Σ W_ij;kl·(P_kl(t) - P_ij(t))
   
   Each pixel influences neighbors through weighted connections
   α: diffusion coefficient (0.0-1.0)
   W: connection weights between adjacent pixels

2️⃣ SEED POINTS (anchor waves to reality):
   P_ij = P_real  if (i,j) ∈ {seed_points}
   
   Known pixels guide pattern emergence and prevent drift
   Seeds "lock in" expected values, forcing structure

3️⃣ CONNECTION REINFORCEMENT (strengthen correct patterns):
   W_ij;kl(t+1) = W_ij;kl(t) + γ·exp(-||P_ij - P_kl||²)
   
   γ: reinforcement rate
   Similar pixels strengthen their mutual influence
   Incorrect connections weaken naturally

COMMANDS:

  ./programme pattern create <width> <height> [seed.png] [density]
      Initialize pattern engine
      Default: 512x512, no seeds
      With seeds: anchors waves to reference image

  ./programme pattern diffuse <iterations> [viz_interval]
      Run pure diffusion (waves spread)
      Watch abstract waves form patterns
      Example: ./programme pattern diffuse 100 10

  ./programme pattern reinforce <iterations>
      Strengthen connections between similar pixels
      Creates more stable, recognizable patterns
      Example: ./programme pattern reinforce 20

  ./programme pattern seed add <x> <y> <r> <g> <b>
      Add single anchor point
      Example: ./programme pattern seed add 50 50 255 0 0

  ./programme pattern seed load <image.png> <density>
      Load seeds from reference image
      Density: 0.05 (5%) to 0.5 (50%)
      Example: ./programme pattern seed load ref.png 0.15

  ./programme pattern emerge <w> <h> <iterations> [seed.png] [density]
      FULL CYCLE: init → diffuse → reinforce
      Transforms waves into recognizable structures
      Saves ALL intermediate images for visualization
      Example: ./programme pattern emerge 512 512 200 input/image/face.png 0.2

  ./programme pattern extract <w> <h> <iterations> [seed.png] [density]
      FAST EXTRACTION: Full pipeline → Final image ONLY
      Runs complete emergence but keeps ONLY final result
      Automatically deletes all intermediate files
      Perfect for quick results without clutter
      Example: ./programme pattern extract 512 512 200 input/image/face.png 0.2

═════════════════════════════════════════════════════════════════
WORKFLOW EXAMPLE:

1. Create engine with seed reference:
   ./programme pattern emerge 512 512 300 input/image/test.png 0.15

2. Check output/ directory:
   pattern_emerge_0001.png  → Starting state (pure waves)
   pattern_emerge_0100.png  → Diffusion in progress
   pattern_emerge_0200.png  → Pattern structure emerges
   pattern_final_emerged.png → Final recognizable image

═════════════════════════════════════════════════════════════════
PARAMETERS & TUNING:

Width/Height:
  64-256:     Quick experiments (few seconds)
  512-1024:   Production quality (minutes)
  2048+:      Ultra high-res (hours)

Iterations:
  50-100:     Basic pattern formation
  200-500:    Good detail emergence
  1000+:      Very refined patterns

Seed Density:
  0.05:       Sparse anchoring (fewer constraints)
  0.15:       Moderate guidance (balanced)
  0.5:        Dense seeds (strong structure)

Alpha (α, diffusion):
  Low (0.05):  Slow, gradual waves
  Medium (0.15): Balanced emergence
  High (0.3):  Fast, may destabilize

Gamma (γ, reinforcement):
  Low (0.01):  Weak pattern memory
  Medium (0.05): Normal learning
  High (0.1):  Strong pattern fixation

═════════════════════════════════════════════════════════════════
EXPECTED PROGRESSION:

Iteration 0:    Uniform gray (no structure)
Iteration 50:   Vague color waves spreading
Iteration 100:  Regions forming, seeds anchoring
Iteration 200:  Clear patterns emerging
Iteration 300:  Recognizable structures appear
Iteration 500+: High-detail, stable images

═════════════════════════════════════════════════════════════════
PATTERN INDEXING & DATABASE:

NEW COMMANDS (for pattern discovery & reuse):

  ./programme pattern index [input_dir] [db_path]
      Scan /input for images and create pattern database
      Analyzes colors, complexity, keywords
      Creates patterns.db with all metadata
      
  ./programme pattern list [db_path]
      Display all indexed patterns with summaries
      
  ./programme pattern info <pattern_id> [db_path]
      Show detailed information about one pattern
      
  ./programme pattern stats [db_path]
      Display database statistics
      
  ./programme pattern search <type> <query> [db_path]
      Search patterns by category or keyword

═════════════════════════════════════════════════════════════════
`)
}

// HandlePatternIndex scans input directory and indexes all patterns
func HandlePatternIndex(args []string) {
	inputPath := "input"
	dbPath := "patterns.db"

	if len(args) > 0 {
		inputPath = args[0]
	}
	if len(args) > 1 {
		dbPath = args[1]
	}

	indexer := database.NewPatternIndexer(inputPath, dbPath)

	err := indexer.LoadDatabase()
	if err != nil && !os.IsNotExist(err) {
		fmt.Printf("⚠ Warning: could not load existing database: %v\n", err)
	}

	err = indexer.IndexDirectory()
	if err != nil {
		fmt.Printf("❌ Indexing failed: %v\n", err)
		return
	}

	indexer.PrintPatternStats()
}

// HandlePatternList displays all indexed patterns
func HandlePatternList(args []string) {
	dbPath := "patterns.db"
	if len(args) > 0 {
		dbPath = args[0]
	}

	indexer := database.NewPatternIndexer("input", dbPath)
	err := indexer.LoadDatabase()
	if err != nil {
		fmt.Printf("❌ Cannot load pattern database: %v\n", err)
		return
	}

	db := indexer.GetPatternDatabase()
	if len(db.Patterns) == 0 {
		fmt.Println("No patterns in database. Run: ./programme pattern index")
		return
	}

	fmt.Println("\n📋 INDEXED PATTERNS")
	fmt.Println("═══════════════════════════════════════════════════")
	fmt.Printf("Total: %d patterns\n\n", len(db.Patterns))

	for i, id := range db.Index {
		if pattern, exists := db.Patterns[id]; exists {
			fmt.Printf("[%d] %s\n", i+1, pattern.ID)
			fmt.Printf("    File: %s\n", pattern.Filename)
			fmt.Printf("    Size: %dx%d\n", pattern.Width, pattern.Height)
			fmt.Printf("    Complexity: %.2f | Confidence: %.1f%%\n",
				pattern.Complexity, pattern.Confidence*100)

			if len(pattern.Categories) > 0 {
				fmt.Printf("    Categories: ")
				for cat, count := range pattern.Categories {
					fmt.Printf("%s(%d) ", cat, count)
				}
				fmt.Println()
			}

			if len(pattern.Keywords) > 0 {
				fmt.Printf("    Keywords: %s\n", strings.Join(pattern.Keywords, ", "))
			}
			fmt.Println()
		}
	}
	fmt.Println("═══════════════════════════════════════════════════\n")
}

// HandlePatternInfo displays detailed information about a pattern
func HandlePatternInfo(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: ./programme pattern info <pattern_id> [patterns.db]")
		return
	}

	patternID := args[0]
	dbPath := "patterns.db"
	if len(args) > 1 {
		dbPath = args[1]
	}

	indexer := database.NewPatternIndexer("input", dbPath)
	err := indexer.LoadDatabase()
	if err != nil {
		fmt.Printf("Cannot load pattern database: %v\n", err)
		return
	}

	pattern, exists := indexer.FindPatternByID(patternID)
	if !exists {
		fmt.Printf(" Pattern not found: %s\n", patternID)
		return
	}

	fmt.Println("\n PATTERN DETAILS")
	fmt.Println("═══════════════════════════════════════════════════")
	fmt.Printf("ID: %s\n", pattern.ID)
	fmt.Printf("File: %s\n", pattern.Filename)
	fmt.Printf("Dimensions: %d × %d pixels\n", pattern.Width, pattern.Height)
	fmt.Printf("Area: %d pixels\n", pattern.Width*pattern.Height)
	fmt.Printf("\nComplexity: %.2f (0=simple, 1=complex)\n", pattern.Complexity)
	fmt.Printf("Confidence: %.1f%%\n", pattern.Confidence*100)
	fmt.Printf("\nAverage Color (RGB): %.0f, %.0f, %.0f\n",
		pattern.AverageColor[0]/256, pattern.AverageColor[1]/256, pattern.AverageColor[2]/256)

	fmt.Println("\nActivated Categories:")
	if len(pattern.Categories) == 0 {
		fmt.Println("  (none)")
	} else {
		for cat, count := range pattern.Categories {
			fmt.Printf("  • %s: %d neurons activated\n", cat, count)
		}
	}

	fmt.Println("\nKeywords:")
	if len(pattern.Keywords) == 0 {
		fmt.Println("  (none)")
	} else {
		for _, kw := range pattern.Keywords {
			fmt.Printf("  • %s\n", kw)
		}
	}

	fmt.Printf("\nContent: %s\n", pattern.ContentSummary)
	fmt.Printf("Hash: %s\n", pattern.PatternDataHash)
	fmt.Println("═══════════════════════════════════════════════════\n")
}

// HandlePatternStats displays database statistics
func HandlePatternStats(args []string) {
	dbPath := "patterns.db"
	if len(args) > 0 {
		dbPath = args[0]
	}

	indexer := database.NewPatternIndexer("input", dbPath)
	err := indexer.LoadDatabase()
	if err != nil {
		fmt.Printf("❌ Cannot load pattern database: %v\n", err)
		return
	}

	indexer.PrintPatternStats()
}

// HandlePatternSearch searches patterns by keyword or category
func HandlePatternSearch(args []string) {
	if len(args) < 2 {
		fmt.Println("Usage: ./programme pattern search <type> <query> [patterns.db]")
		fmt.Println("  type: category | keyword")
		fmt.Println("  Example: ./programme pattern search category TECH")
		fmt.Println("  Example: ./programme pattern search keyword forest")
		return
	}

	searchType := strings.ToLower(args[0])
	query := strings.ToLower(args[1])
	dbPath := "patterns.db"
	if len(args) > 2 {
		dbPath = args[2]
	}

	indexer := database.NewPatternIndexer("input", dbPath)
	err := indexer.LoadDatabase()
	if err != nil {
		fmt.Printf("❌ Cannot load pattern database: %v\n", err)
		return
	}

	fmt.Println("\n🔍 PATTERN SEARCH")
	fmt.Println("═══════════════════════════════════════════════════")
	fmt.Printf("Type: %s | Query: %s\n\n", searchType, query)

	var results []database.PatternMetadata

	if searchType == "category" {
		results = indexer.FindPatternsByCategory(strings.ToUpper(query))
		fmt.Printf("Found %d patterns with category %s\n\n", len(results), strings.ToUpper(query))
	} else if searchType == "keyword" {
		db := indexer.GetPatternDatabase()
		for _, id := range db.Index {
			if pattern, exists := db.Patterns[id]; exists {
				for _, kw := range pattern.Keywords {
					if strings.Contains(strings.ToLower(kw), query) {
						results = append(results, pattern)
						break
					}
				}
			}
		}
		fmt.Printf("Found %d patterns with keyword '%s'\n\n", len(results), query)
	} else {
		fmt.Printf("❌ Unknown search type: %s\n", searchType)
		return
	}

	if len(results) == 0 {
		fmt.Println("No patterns found matching your query.")
		fmt.Println("═══════════════════════════════════════════════════\n")
		return
	}

	for i, pattern := range results {
		fmt.Printf("[%d] %s (%dx%d)\n", i+1, pattern.ID, pattern.Width, pattern.Height)
		fmt.Printf("    Complexity: %.2f | Confidence: %.1f%%\n",
			pattern.Complexity, pattern.Confidence*100)
		if len(pattern.Keywords) > 0 {
			fmt.Printf("    Keywords: %s\n", strings.Join(pattern.Keywords, ", "))
		}
		fmt.Println()
	}
	fmt.Println("═══════════════════════════════════════════════════\n")
}
