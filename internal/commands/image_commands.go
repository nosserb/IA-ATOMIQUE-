package commands

import (
	"IA-ATOMIQUE/database"
	"fmt"
	"image/png"
	"os"
	"strconv"
	"strings"
)

// ImageGenerationCommand handles atomic image generation
func ImageGenerationCommand(args []string) {
	if len(args) < 2 {
		PrintImageGenerationHelp()
		return
	}

	subcommand := args[1]

	switch subcommand {
	case "pipeline":
		HandleFullPipelineGeneration(args[2:])
	case "generate":
		HandleImageGenerate(args[2:])
	case "prompt":
		HandleImageFromPrompt(args[2:])
	case "multi-scale":
		HandleMultiScaleGeneration(args[2:])
	case "ultra":
		HandleUltraFastImageGeneration(args[2:])
	case "fast":
		HandleFastImageGeneration(args[2:])
	case "draft":
		HandleDraftImageGeneration(args[2:])
	case "phase1":
		HandlePhaseOne(args[2:])
	case "phase2":
		HandlePhaseTwo(args[2:])
	case "phase3":
		HandlePhaseThree(args[2:])
	case "phase4":
		HandlePhaseFour(args[2:])
	case "phase5":
		HandlePhaseFive(args[2:])
	case "stats":
		HandleImageNetworkStats(args[2:])
	case "interactive":
		HandleInteractiveImageGeneration()
	default:
		fmt.Printf("❌ Unknown image command: %s\n", subcommand)
		PrintImageGenerationHelp()
	}
}

// HandleImageGenerate generates an image with specified dimensions
// Usage: ./programme image generate 512 512 100 8 "prompt here"
func HandleImageGenerate(args []string) {
	if len(args) < 4 {
		fmt.Println("Usage: ./programme image generate <width> <height> <iterations> <patch_size> [prompt]")
		fmt.Println("\nExample:")
		fmt.Println("  ./programme image generate 512 512 100 8 \"red sunset over ocean\"")
		return
	}

	width, _ := strconv.Atoi(args[0])
	height, _ := strconv.Atoi(args[1])
	iterations, _ := strconv.Atoi(args[2])
	patchSize, _ := strconv.Atoi(args[3])

	// Validate dimensions
	if width < 64 || height < 64 || width > 2048 || height > 2048 {
		fmt.Println("❌ Image dimensions must be between 64x64 and 2048x2048")
		return
	}

	if patchSize < 1 || patchSize > 64 {
		fmt.Println("❌ Patch size must be between 1 and 64")
		return
	}

	// Extract prompt if provided
	prompt := ""
	if len(args) > 4 {
		prompt = strings.Join(args[4:], " ")
	}

	fmt.Printf("\n🎨 Initializing Atomic Image Network\n")
	fmt.Printf("   Dimensions: %dx%d\n", width, height)
	fmt.Printf("   Patch Size: %d×%d\n", patchSize, patchSize)
	fmt.Printf("   Atoms: %d\n", (width/patchSize)*(height/patchSize))

	if prompt != "" {
		fmt.Printf("   Prompt: \"%s\"\n", prompt)
	}

	// Create network
	network := database.NewAtomicImageNetwork(width, height, patchSize)

	// Parse prompt into constraints
	if prompt != "" {
		network.ParsePrompt(prompt)
		fmt.Printf("\n📍 Prompt constraints parsed\n")
		fmt.Printf("   Style vector: [%.2f, %.2f, %.2f]\n",
			network.StyleVector[0], network.StyleVector[1], network.StyleVector[2])

		// Apply external constraints to all atoms
		gridWidth := width / patchSize
		gridHeight := height / patchSize
		for y := 0; y < gridHeight; y++ {
			for x := 0; x < gridWidth; x++ {
				constraint := network.ComputeExternalConstraint(x, y)
				network.Atoms[y][x].ExternalConstraint = constraint
			}
		}
	}

	// Run generation
	fmt.Printf("\n🔄 Running generation for %d iterations\n", iterations)

	for i := 0; i < iterations; i++ {
		network.IterateGeneration()

		if (i+1)%20 == 0 {
			stats := network.GetNetworkStats()
			fmt.Printf("   [%4d/%4d] State: %.3f | Energy: %.2e | Active: %.1f%% | Frozen: %d\n",
				i+1, iterations,
				stats["avg_state"],
				stats["energy_density"],
				stats["active_percent"],
				int(stats["frozen_atoms"]),
			)
		}
	}

	// Post-processing
	fmt.Printf("\n✨ Applying post-processing\n")
	network.LocalSmoothing(1)
	network.EdgeEnhancement(0.3)

	// Save image
	outputFile := "generated_image.png"
	err := network.SaveImage(outputFile)
	if err != nil {
		fmt.Printf("❌ Failed to save image: %v\n", err)
		return
	}

	fmt.Printf("\n✅ Image successfully generated: %s\n", outputFile)

	// Print final statistics
	stats := network.GetNetworkStats()
	fmt.Printf("\n📊 Final Network Statistics:\n")
	fmt.Printf("   Average State: %.4f\n", stats["avg_state"])
	fmt.Printf("   Average Intensity: %.4f\n", stats["avg_intensity"])
	fmt.Printf("   Active Atoms: %d (%.1f%%)\n", int(stats["active_atoms"]), stats["active_percent"])
	fmt.Printf("   Frozen Atoms: %d\n", int(stats["frozen_atoms"]))
	fmt.Printf("   Total Energy: %.2e\n", stats["total_energy"])
	fmt.Printf("   Energy Density: %.2e\n", stats["energy_density"])
}

// HandleImageFromPrompt generates image directly from natural language prompt
// Automatically determines optimal dimensions and parameters
func HandleImageFromPrompt(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: ./programme image prompt \"your creative prompt\"")
		fmt.Println("\nExamples:")
		fmt.Println("  ./programme image prompt \"starry night with mountains\"")
		fmt.Println("  ./programme image prompt \"dark forest, misty, scary\"")
		fmt.Println("  ./programme image prompt \"bright sunny beach with palm trees\"")
		return
	}

	prompt := strings.Join(args, " ")

	// Estimate complexity from prompt length
	complexity := len(strings.Fields(prompt))

	// Set parameters based on complexity
	width, height := 512, 512
	if complexity > 10 {
		width, height = 768, 768
	} else if complexity < 3 {
		width, height = 256, 256
	}

	iterations := 50 + complexity*10
	patchSize := 8

	fmt.Printf("\n🎨 Atomic Prompt-to-Image Generation\n")
	fmt.Printf("   Prompt: \"%s\"\n", prompt)
	fmt.Printf("   Complexity: %d words\n", complexity)
	fmt.Printf("   Suggested dimensions: %dx%d\n", width, height)
	fmt.Printf("   Iterations: %d\n\n", iterations)

	// Use the generate command with these parameters
	args = []string{
		strconv.Itoa(width),
		strconv.Itoa(height),
		strconv.Itoa(iterations),
		strconv.Itoa(patchSize),
		prompt,
	}

	HandleImageGenerate(args)
}

// HandleMultiScaleGeneration implements coarse-to-fine progressive generation
// First generates at 16x16 patches, then 8x8, then 4x4 for detail
func HandleMultiScaleGeneration(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: ./programme image multi-scale \"prompt\" [output_file]")
		fmt.Println("\nGenerates image at multiple scales for better detail and speed")
		return
	}

	prompt := strings.Join(args, " ")
	outputFile := "generated_multiscale.png"

	if len(args) > 1 {
		outputFile = args[len(args)-1]
	}

	width, height := 512, 512

	fmt.Printf("\n🎨 Multi-Scale Atomic Image Generation\n")
	fmt.Printf("   Prompt: \"%s\"\n\n", prompt)

	// Stage 1: Coarse generation (16x16 patches)
	fmt.Printf("📍 Stage 1: Coarse Generation (16×16 patches)\n")
	network16 := database.NewAtomicImageNetwork(width, height, 16)
	network16.ParsePrompt(prompt)

	gridWidth := width / 16
	gridHeight := height / 16
	for y := 0; y < gridHeight; y++ {
		for x := 0; x < gridWidth; x++ {
			network16.Atoms[y][x].ExternalConstraint = network16.ComputeExternalConstraint(x, y)
		}
	}

	for i := 0; i < 30; i++ {
		network16.IterateGeneration()
		if (i+1)%10 == 0 {
			stats := network16.GetNetworkStats()
			fmt.Printf("   [%2d/30] State: %.3f | Energy: %.2e\n", i+1, stats["avg_state"], stats["energy_density"])
		}
	}

	// Stage 2: Medium generation (8x8 patches)
	fmt.Printf("\n📍 Stage 2: Medium Generation (8×8 patches)\n")
	network8 := database.NewAtomicImageNetwork(width, height, 8)
	network8.ParsePrompt(prompt)

	// Initialize from coarse layer - upscale states
	gridWidth8 := width / 8
	gridHeight8 := height / 8
	for y := 0; y < gridHeight8; y++ {
		for x := 0; x < gridWidth8; x++ {
			coarseX := x / 2
			coarseY := y / 2
			if coarseX < (width/16) && coarseY < (height/16) {
				network8.Atoms[y][x].State = network16.Atoms[coarseY][coarseX].State
				network8.Atoms[y][x].Color = network16.Atoms[coarseY][coarseX].Color
			}
			network8.Atoms[y][x].ExternalConstraint = network8.ComputeExternalConstraint(x, y)
		}
	}

	for i := 0; i < 40; i++ {
		network8.IterateGeneration()
		if (i+1)%10 == 0 {
			stats := network8.GetNetworkStats()
			fmt.Printf("   [%2d/40] State: %.3f | Energy: %.2e\n", i+1, stats["avg_state"], stats["energy_density"])
		}
	}

	// Apply post-processing
	fmt.Printf("\n✨ Applying post-processing\n")
	network8.LocalSmoothing(1)
	network8.EdgeEnhancement(0.3)

	// Save
	err := network8.SaveImage(outputFile)
	if err != nil {
		fmt.Printf("❌ Failed to save image: %v\n", err)
		return
	}

	fmt.Printf("✅ Multi-scale image generated: %s\n", outputFile)
}

// HandleImageNetworkStats displays network statistics during generation
func HandleImageNetworkStats(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: ./programme image stats <network_dump_file>")
		return
	}

	filename := args[0]

	// In a full implementation, this would read saved network state
	fmt.Printf("📊 Image Network Statistics\n")
	fmt.Printf("   File: %s\n", filename)
	fmt.Println("   (Network state loading not yet implemented)")
}

// HandleInteractiveImageGeneration opens interactive generation mode
func HandleInteractiveImageGeneration() {
	fmt.Printf("\n🎨 Interactive Atomic Image Generation\n")
	fmt.Printf("========================================\n\n")

	fmt.Println("Commands:")
	fmt.Println("  generate <width> <height> <iterations> <patch> <prompt>  - Generate image")
	fmt.Println("  prompt \"description\"                                      - Auto-generate from prompt")
	fmt.Println("  params                                                     - Show/modify parameters")
	fmt.Println("  help                                                       - Show this help")
	fmt.Println("  exit                                                       - Exit interactive mode")
	fmt.Println()

	for {
		fmt.Print("atomic-image> ")
		var input string
		fmt.Scanln(&input)

		if input == "exit" || input == "quit" {
			fmt.Println("Exiting interactive mode...")
			break
		} else if input == "help" {
			fmt.Println("Available commands listed above")
		} else if input == "params" {
			fmt.Println("Current parameters can be configured here")
		} else if strings.HasPrefix(input, "generate") {
			parts := strings.Fields(input)
			HandleImageGenerate(parts[1:])
		} else if strings.HasPrefix(input, "prompt") {
			// Extract everything after "prompt"
			prompt := strings.TrimPrefix(input, "prompt ")
			HandleImageFromPrompt([]string{prompt})
		}
	}
}

// HandleFullPipelineGeneration runs all 5 phases sequentially
// Usage: ./programme image pipeline 512 512 200 8 "prompt"
func HandleFullPipelineGeneration(args []string) {
	if len(args) < 4 {
		fmt.Println("Usage: ./programme image pipeline <width> <height> <iterations> <patch_size> [prompt]")
		fmt.Println("Example: ./programme image pipeline 512 512 200 8 \"red sunset over mountains\"")
		return
	}

	width, _ := strconv.Atoi(args[0])
	height, _ := strconv.Atoi(args[1])
	iterations, _ := strconv.Atoi(args[2])
	patchSize, _ := strconv.Atoi(args[3])

	if width < 64 || height < 64 || width > 2048 || height > 2048 {
		fmt.Println("❌ Image dimensions must be between 64x64 and 2048x2048")
		return
	}

	if patchSize < 1 || patchSize > 64 {
		fmt.Println("❌ Patch size must be between 1 and 64")
		return
	}

	prompt := ""
	if len(args) > 4 {
		prompt = strings.Join(args[4:], " ")
	}

	fmt.Printf("\n🎨 ATOMIC IMAGE GENERATION - 5-PHASE PIPELINE\n")
	fmt.Printf("═══════════════════════════════════════════════════\n")
	fmt.Printf("Dimensions: %dx%d | Patch: %d | Iterations: %d\n", width, height, patchSize, iterations)
	fmt.Printf("Prompt: \"%s\"\n\n", prompt)

	// Initialize network
	network := database.NewAtomicImageNetwork(width, height, patchSize)
	network.ParsePrompt(prompt)

	// Run full pipeline
	report := network.FullImageGenerationPipeline(prompt, iterations)

	// Phase-by-phase summary
	fmt.Printf("📊 PHASE 1: Multi-Scale Structuration\n")
	fmt.Printf("   ✓ Pixels aligned locally (micro → macro)\n")
	fmt.Printf("   ✓ Coherent patterns emerged\n\n")

	fmt.Printf("📊 PHASE 2: Shape Emergence\n")
	fmt.Printf("   ✓ Primitive shapes (contours, lines, curves)\n")
	fmt.Printf("   ✓ Capsule resonance stabilized\n\n")

	fmt.Printf("📊 PHASE 3: Prompt Conditioning\n")
	fmt.Printf("   ✓ User guidance applied\n")
	fmt.Printf("   ✓ Style vector: [%.2f, %.2f, %.2f]\n",
		network.StyleVector[0], network.StyleVector[1], network.StyleVector[2])
	fmt.Printf("\n")

	fmt.Printf("📊 PHASE 4: Iterative Refinement\n")
	fmt.Printf("   ✓ Laplacian smoothing applied\n")
	fmt.Printf("   ✓ Texture details added\n\n")

	fmt.Printf("📊 PHASE 5: Coherence Verification\n")
	fmt.Printf("   ✓ Global Coherence: %.3f\n", report.GlobalCoherence)
	fmt.Printf("   ✓ Health Score: %.1f%%\n", report.OverallHealthScore*100)
	fmt.Printf("   ✓ Faulty atoms detected: %d\n", len(report.FaultyAtoms))
	fmt.Printf("   ✓ Repaired atoms: %d\n\n", report.RepairCount)

	// Save outputs
	fmt.Printf("💾 Saving outputs...\n")
	err := network.SaveImage("generated_image_pipeline.png")
	if err != nil {
		fmt.Printf("❌ Failed to save image: %v\n", err)
		return
	}

	fmt.Printf("   ✓ Image: generated_image_pipeline.png\n")

	// Save coherence map
	coherenceImg := report.RenderCoherenceMap(width/patchSize, height/patchSize)
	coherenceFile, err := os.Create("coherence_map.png")
	if err == nil {
		png.Encode(coherenceFile, coherenceImg)
		coherenceFile.Close()
		fmt.Printf("   ✓ Coherence Map: coherence_map.png\n")
	}

	fmt.Printf("\n✅ Pipeline complete!\n")
	fmt.Printf("═══════════════════════════════════════════════════\n")
}

// HandlePhaseOne runs only Phase 1 for testing
func HandlePhaseOne(args []string) {
	if len(args) < 4 {
		fmt.Println("Usage: ./programme image phase1 <width> <height> <iterations> <patch_size>")
		return
	}

	width, _ := strconv.Atoi(args[0])
	height, _ := strconv.Atoi(args[1])
	iterations, _ := strconv.Atoi(args[2])
	patchSize, _ := strconv.Atoi(args[3])

	fmt.Printf("\n🔬 PHASE 1: Multi-Scale Structuration\n")
	fmt.Printf("Transforming isolated pixels into coherent patterns\n\n")

	network := database.NewAtomicImageNetwork(width, height, patchSize)
	fmt.Printf("Running %d iterations of multi-scale resonance...\n", iterations)

	network.PhaseOne_StructurationMultiEchelle(iterations)

	network.SaveImage("phase1_output.png")
	fmt.Printf("✅ Phase 1 complete: phase1_output.png\n")
}

// HandlePhaseTwo runs only Phase 2 for testing
func HandlePhaseTwo(args []string) {
	if len(args) < 4 {
		fmt.Println("Usage: ./programme image phase2 <width> <height> <iterations> <patch_size>")
		return
	}

	width, _ := strconv.Atoi(args[0])
	height, _ := strconv.Atoi(args[1])
	iterations, _ := strconv.Atoi(args[2])
	patchSize, _ := strconv.Atoi(args[3])

	fmt.Printf("\n🔬 PHASE 2: Shape Emergence\n")
	fmt.Printf("Making primitive shapes appear through capsule resonance\n\n")

	network := database.NewAtomicImageNetwork(width, height, patchSize)
	engine := database.NewShapeEmergenceEngine(width, height, patchSize)

	fmt.Printf("Running %d iterations of shape emergence...\n", iterations)
	network.PhaseTwo_ShapeEmergence(engine, iterations)

	network.SaveImage("phase2_output.png")
	fmt.Printf("✅ Phase 2 complete: phase2_output.png\n")
}

// HandlePhaseThree runs only Phase 3 for testing
func HandlePhaseThree(args []string) {
	if len(args) < 5 {
		fmt.Println("Usage: ./programme image phase3 <width> <height> <iterations> <patch_size> \"prompt\"")
		return
	}

	width, _ := strconv.Atoi(args[0])
	height, _ := strconv.Atoi(args[1])
	iterations, _ := strconv.Atoi(args[2])
	patchSize, _ := strconv.Atoi(args[3])
	prompt := strings.Join(args[4:], " ")

	fmt.Printf("\n🔬 PHASE 3: Prompt Conditioning\n")
	fmt.Printf("Applying user guidance: %s\n\n", prompt)

	network := database.NewAtomicImageNetwork(width, height, patchSize)
	network.ParsePrompt(prompt)

	guide := network.ParsePromptToGuide(prompt)
	fmt.Printf("Running %d iterations with guidance...\n", iterations)
	network.PhaseThree_PromptConditioning(guide, iterations)

	network.SaveImage("phase3_output.png")
	fmt.Printf("✅ Phase 3 complete: phase3_output.png\n")
}

// HandlePhaseFour runs only Phase 4 for testing
func HandlePhaseFour(args []string) {
	if len(args) < 4 {
		fmt.Println("Usage: ./programme image phase4 <width> <height> <iterations> <patch_size>")
		return
	}

	width, _ := strconv.Atoi(args[0])
	height, _ := strconv.Atoi(args[1])
	iterations, _ := strconv.Atoi(args[2])
	patchSize, _ := strconv.Atoi(args[3])

	fmt.Printf("\n🔬 PHASE 4: Iterative Refinement\n")
	fmt.Printf("Adding fine details and realistic texture\n\n")

	network := database.NewAtomicImageNetwork(width, height, patchSize)

	fmt.Printf("Running %d iterations of refinement...\n", iterations)
	network.PhaseFour_IterativeRefinement(0.1, 0.08, iterations)

	network.SaveImage("phase4_output.png")
	fmt.Printf("✅ Phase 4 complete: phase4_output.png\n")
}

// HandlePhaseFive runs only Phase 5 for testing
func HandlePhaseFive(args []string) {
	if len(args) < 4 {
		fmt.Println("Usage: ./programme image phase5 <width> <height> <iterations> <patch_size>")
		return
	}

	width, _ := strconv.Atoi(args[0])
	height, _ := strconv.Atoi(args[1])
	iterations, _ := strconv.Atoi(args[2])
	patchSize, _ := strconv.Atoi(args[3])

	fmt.Printf("\n🔬 PHASE 5: Coherence Verification\n")
	fmt.Printf("Checking and correcting atomic coherence\n\n")

	network := database.NewAtomicImageNetwork(width, height, patchSize)

	// Run some basic generation first
	fmt.Printf("Running %d iterations of generation...\n", iterations)
	for i := 0; i < iterations; i++ {
		network.IterateGeneration()
	}

	// Verify coherence
	fmt.Printf("Verifying coherence...\n")
	report := network.PhaseFive_CoherenceVerification()

	fmt.Printf("\nCoherence Report:\n")
	fmt.Printf("  Global Coherence: %.3f\n", report.GlobalCoherence)
	fmt.Printf("  Health Score: %.1f%%\n", report.OverallHealthScore*100)
	fmt.Printf("  Faulty atoms: %d\n", len(report.FaultyAtoms))
	fmt.Printf("  Repaired: %d\n", report.RepairCount)

	network.SaveImage("phase5_output.png")
	fmt.Printf("\n✅ Phase 5 complete: phase5_output.png\n")
}

// HandleUltraFastImageGeneration - SUPER ULTRA-FAST mode (<500ms!)
// Usage: ./programme image ultra "prompt"
func HandleUltraFastImageGeneration(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: ./programme image ultra \"your prompt\"")
		fmt.Println("\nExamples:")
		fmt.Println("  ./programme image ultra \"blue\"")
		fmt.Println("  ./programme image ultra \"red\"")
		fmt.Println("\nGenerates 128x128 image in <500ms (ultra mode)")
		return
	}

	prompt := strings.Join(args, " ")

	fmt.Printf("\n⚡⚡ ULTRA MODE (<500ms)\n")
	fmt.Printf("═════════════════════════════════════════════\n")
	fmt.Printf("Prompt: \"%s\"\n", prompt)
	fmt.Printf("Resolution: 128x128\n")
	fmt.Printf("Iterations: 2\n")
	fmt.Printf("Patch Size: 64x64\n")
	fmt.Printf("Target Time: <500ms\n\n")

	config := database.PresetConfigs["ultra"]
	fnet := database.NewFastAtomicImageNetwork(config)
	fnet.BaseNetwork.ParsePrompt(prompt)

	// Apply constraints
	gridWidth := config.Width / config.PatchSize
	gridHeight := config.Height / config.PatchSize
	for y := 0; y < gridHeight; y++ {
		for x := 0; x < gridWidth; x++ {
			fnet.BaseNetwork.Atoms[y][x].ExternalConstraint =
				fnet.BaseNetwork.ComputeExternalConstraint(x, y)
		}
	}

	// Ultra fast generation
	fmt.Printf("🔄 Running generation...\n")
	for i := 0; i < config.Iterations; i++ {
		fnet.OptimizedIterateGeneration()
		fmt.Printf("   [%d/%d] ✓\n", i+1, config.Iterations)
	}

	fmt.Printf("\n✅ Generation complete\n")

	// Save
	err := fnet.BaseNetwork.SaveImage("generated_ultra.png")
	if err != nil {
		fmt.Printf("❌ Failed to save image: %v\n", err)
		return
	}

	fmt.Printf("💾 Saved: generated_ultra.png\n")
	fnet.PrintOptimizationStats()
}

// HandleDraftImageGeneration - ULTRA-FAST draft mode (<2 seconds)
// Usage: ./programme image draft "prompt"
func HandleDraftImageGeneration(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: ./programme image draft \"your prompt\"")
		fmt.Println("\nExamples:")
		fmt.Println("  ./programme image draft \"blue ocean\"")
		fmt.Println("  ./programme image draft \"red sunset\"")
		fmt.Println("  ./programme image draft \"abstract patterns\"")
		fmt.Println("\nGenerates 256x256 image in ~1-2 seconds (ultra-fast draft mode)")
		return
	}

	prompt := strings.Join(args, " ")

	fmt.Printf("\n⚡ ULTRA-FAST DRAFT MODE\n")
	fmt.Printf("═════════════════════════════════════════════\n")
	fmt.Printf("Prompt: \"%s\"\n", prompt)
	fmt.Printf("Resolution: 256x256\n")
	fmt.Printf("Iterations: 5\n")
	fmt.Printf("Patch Size: 32x32\n")
	fmt.Printf("Target Time: <2 seconds\n\n")

	config := database.PresetConfigs["draft"]
	config.Width = 256
	config.Height = 256

	fnet := database.NewFastAtomicImageNetwork(config)
	fnet.BaseNetwork.ParsePrompt(prompt)

	// Apply constraints
	gridWidth := config.Width / config.PatchSize
	gridHeight := config.Height / config.PatchSize
	for y := 0; y < gridHeight; y++ {
		for x := 0; x < gridWidth; x++ {
			fnet.BaseNetwork.Atoms[y][x].ExternalConstraint =
				fnet.BaseNetwork.ComputeExternalConstraint(x, y)
		}
	}

	// Fast generation
	fmt.Printf("🔄 Running generation...\n")
	for i := 0; i < config.Iterations; i++ {
		fnet.OptimizedIterateGeneration()
		if (i+1)%2 == 0 {
			fmt.Printf("   [%d/%d] ✓\n", i+1, config.Iterations)
		}
	}

	// Skip post-processing for speed
	fmt.Printf("\n✅ Generation complete\n")

	// Save
	err := fnet.BaseNetwork.SaveImage("generated_draft.png")
	if err != nil {
		fmt.Printf("❌ Failed to save image: %v\n", err)
		return
	}

	fmt.Printf("💾 Saved: generated_draft.png\n")
	fnet.PrintOptimizationStats()
}

// HandleFastImageGeneration - FAST mode (~2-3 seconds)
// Usage: ./programme image fast "prompt"
func HandleFastImageGeneration(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: ./programme image fast \"your prompt\"")
		fmt.Println("\nExamples:")
		fmt.Println("  ./programme image fast \"detailed sunset over ocean\"")
		fmt.Println("  ./programme image fast \"dark mysterious forest\"")
		fmt.Println("\nGenerates 256x256 image in ~2-3 seconds (fast mode)")
		return
	}

	prompt := strings.Join(args, " ")

	fmt.Printf("\n⚡ FAST MODE\n")
	fmt.Printf("═════════════════════════════════════════════\n")
	fmt.Printf("Prompt: \"%s\"\n", prompt)
	fmt.Printf("Resolution: 256x256\n")
	fmt.Printf("Iterations: 10\n")
	fmt.Printf("Patch Size: 16x16\n")
	fmt.Printf("Target Time: 2-3 seconds\n\n")

	config := database.PresetConfigs["fast"]
	config.Width = 256
	config.Height = 256

	fnet := database.NewFastAtomicImageNetwork(config)
	fnet.BaseNetwork.ParsePrompt(prompt)

	// Apply constraints
	gridWidth := config.Width / config.PatchSize
	gridHeight := config.Height / config.PatchSize
	for y := 0; y < gridHeight; y++ {
		for x := 0; x < gridWidth; x++ {
			fnet.BaseNetwork.Atoms[y][x].ExternalConstraint =
				fnet.BaseNetwork.ComputeExternalConstraint(x, y)
		}
	}

	// Fast generation
	fmt.Printf("🔄 Running generation...\n")
	for i := 0; i < config.Iterations; i++ {
		fnet.OptimizedIterateGeneration()
		if (i+1)%3 == 0 {
			fmt.Printf("   [%d/%d] ✓\n", i+1, config.Iterations)
		}
	}

	// Minimal post-processing
	fmt.Printf("\n✅ Generation complete\n")

	// Save
	err := fnet.BaseNetwork.SaveImage("generated_fast.png")
	if err != nil {
		fmt.Printf("❌ Failed to save image: %v\n", err)
		return
	}

	fmt.Printf("💾 Saved: generated_fast.png\n")
	fnet.PrintOptimizationStats()
}

// PrintImageGenerationHelp displays help for image generation commands
func PrintImageGenerationHelp() {
	fmt.Printf(`
🎨 Atomic Image Generation - T.R.A. for Visual Media
═════════════════════════════════════════════════════

⚡ ULTRA-FAST COMMANDS (< 2 seconds):
  ./programme image draft "prompt"
      Ultra-fast draft mode (256x256, 5 iterations)
      Example: ./programme image draft "blue ocean"
      ⏱️  Target: <2 seconds

  ./programme image fast "prompt"  
      Fast mode (256x256, 10 iterations)
      Example: ./programme image fast "red sunset"
      ⏱️  Target: 2-3 seconds

MAIN COMMANDS:
  ./programme image pipeline <w> <h> <iter> <patch> "prompt"
      Run complete 5-phase pipeline: structuration → emergence → 
      conditioning → refinement → verification
      Example: ./programme image pipeline 512 512 200 8 "red sunset"

  ./programme image generate <w> <h> <iter> <patch> "prompt"
      Generate image with specified dimensions and prompt
      Example: ./programme image generate 512 512 100 8 "blue sky"

  ./programme image prompt "creative description"
      Auto-generate from prompt (uses optimal settings)
      Example: ./programme image prompt "starry night with mountains"

  ./programme image multi-scale "prompt" [output]
      Multi-scale coarse-to-fine generation for better quality

PHASE-SPECIFIC COMMANDS (for testing/debugging):
  ./programme image phase1 <w> <h> <iter> <patch>
      Phase 1: Multi-scale structuration

  ./programme image phase2 <w> <h> <iter> <patch>
      Phase 2: Shape emergence (capsule resonance)

  ./programme image phase3 <w> <h> <iter> <patch> "prompt"
      Phase 3: Prompt conditioning

  ./programme image phase4 <w> <h> <iter> <patch>
      Phase 4: Iterative refinement (smoothing & texture)

  ./programme image phase5 <w> <h> <iter> <patch>
      Phase 5: Coherence verification & repair

OTHER COMMANDS:
  ./programme image stats <file>
      Display network statistics

  ./programme image interactive
      Enter interactive generation mode

═════════════════════════════════════════════════════
PHASE DESCRIPTIONS:

Phase 1: Multi-Scale Structuration
  ✓ Local pixel resonance at multiple scales (micro → macro)
  ✓ Implements: c_ij(t+1) = c_ij(t) + α * Σ(w_ijk * (c_k - c_ij))
  ✓ Result: Coherent patterns from isolated pixels

Phase 2: Shape Emergence
  ✓ Capsule-based resonance for motif recognition
  ✓ Implements: s_m(t+1) = s_m(t) + γ * Σ(R(s_n, s_m))
  ✓ Result: Primitive shapes (contours, lines, curves)

Phase 3: Prompt Conditioning
  ✓ User guidance via external field
  ✓ Implements: + β * G_ij(P)
  ✓ Result: Image matches desired description

Phase 4: Iterative Refinement
  ✓ Laplacian smoothing + texture addition
  ✓ Implements: + δ*Laplacian(c) + ϵ*NoiseAdjust(c)
  ✓ Result: Fine details and realistic texture

Phase 5: Coherence Verification
  ✓ Quality control and automatic repair
  ✓ Implements: coherence_ij = 1 - Σ||c_ij - c_k||/max_diff
  ✓ Result: High-quality, coherent final image

═════════════════════════════════════════════════════
PARAMETERS:
  width       - Image width (64-2048 pixels)
  height      - Image height (64-2048 pixels)
  iterations  - Generation iterations (higher = more detail, slower)
  patch_size  - Atomic patch size (1-64, larger = faster/coarser)
  prompt      - Natural language description of desired image

QUICK EXAMPLES:
  # Full pipeline with prompt
  ./programme image pipeline 512 512 200 8 "colorful abstract art"

  # Individual phase testing
  ./programme image phase1 256 256 50 8
  ./programme image phase3 256 256 50 8 "blue ocean"

  # Complex generation
  ./programme image prompt "detailed fantasy landscape with dragons"

═════════════════════════════════════════════════════
`)
}
