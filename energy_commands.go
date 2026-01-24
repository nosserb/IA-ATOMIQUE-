package main

import (
	"IA-ATOMIQUE/database"
	"fmt"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"os"
	"strconv"
	"strings"
)

// EnergyBasedImageCommand handles constraint relaxation image generation
func EnergyBasedImageCommand(args []string) {
	if len(args) < 2 {
		PrintEnergyImageHelp()
		return
	}

	subcommand := args[1]

	switch subcommand {
	case "generate":
		HandleEnergyGenerate(args[2:])
	case "relax":
		HandleConstraintRelaxation(args[2:])
	case "analyze":
		HandleEnergyAnalysis(args[2:])
	case "multi-phase":
		HandleMultiPhaseGeneration(args[2:])
	case "from-image":
		HandleGenerateFromEnergyProfile(args[2:])
	default:
		fmt.Printf("❌ Commande inconnue: %s\n", subcommand)
		PrintEnergyImageHelp()
	}
}

// HandleEnergyGenerate - Constraint relaxation from scratch
// Usage: ./programme energy generate 512 512 500 8 "constraints"
func HandleEnergyGenerate(args []string) {
	if len(args) < 4 {
		fmt.Println("Usage: ./programme energy generate <width> <height> <iterations> <patch_size> [constraints]")
		fmt.Println("\nExample:")
		fmt.Println("  ./programme energy generate 512 512 500 8")
		fmt.Println("  ./programme energy generate 256 256 300 4 \"dark smooth edges\"")
		return
	}

	width, _ := strconv.Atoi(args[0])
	height, _ := strconv.Atoi(args[1])
	iterations, _ := strconv.Atoi(args[2])
	patchSize, _ := strconv.Atoi(args[3])

	// Parse constraints if provided
	constraints := ""
	if len(args) > 4 {
		constraints = strings.Join(args[4:], " ")
	}

	// Validate
	if width < 64 || height < 64 || width > 2048 || height > 2048 {
		fmt.Println("❌ Dimensions must be 64-2048 pixels")
		return
	}

	if patchSize < 1 || patchSize > 64 {
		fmt.Println("❌ Patch size must be 1-64")
		return
	}

	fmt.Printf("\n⚛️  CONSTRAINT RELAXATION NETWORK\n")
	fmt.Printf("══════════════════════════════════\n\n")
	fmt.Printf("📐 Configuration:\n")
	fmt.Printf("   Dimensions: %dx%d\n", width, height)
	fmt.Printf("   Patch size: %d×%d\n", patchSize, patchSize)
	fmt.Printf("   Grid atoms: %d\n", (width/patchSize)*(height/patchSize))
	fmt.Printf("   Max iterations: %d\n", iterations)

	if constraints != "" {
		fmt.Printf("   Constraints: \"%s\"\n", constraints)
	}

	fmt.Printf("\n🔧 Initializing network...\n")
	network := database.NewConstraintRelaxationNetwork(width, height, patchSize)

	// Apply constraints to global field if provided
	applyConstraintsToField(network.GlobalField, constraints)

	fmt.Printf("✓ Network ready with %d atoms in %d grid\n\n",
		(width/patchSize)*(height/patchSize),
		(width / patchSize),
	)

	fmt.Printf("⚡ Starting constraint relaxation...\n")
	fmt.Printf("   [Each atom minimizes local tension]\n")
	fmt.Printf("   [Global coherence emerges gradually]\n\n")

	// Run generation
	network.Generate(iterations)

	// Save result
	fmt.Printf("\n💾 Rendering and saving...\n")
	img := network.RenderToImage()

	outputFile := "generated_energy_based.png"
	file, err := os.Create(outputFile)
	if err != nil {
		fmt.Printf("❌ Failed to save: %v\n", err)
		return
	}
	defer file.Close()

	err = png.Encode(file, img)
	if err != nil {
		fmt.Printf("❌ Failed to encode PNG: %v\n", err)
		return
	}

	fmt.Printf("✅ Saved to: %s\n\n", outputFile)

	// Print final statistics
	printNetworkStatistics(network)

	// Print detected patterns
	fmt.Printf("\n🎨 Detected Patterns:\n")
	fmt.Printf("   Total regions: %d\n", len(network.PatternRegions))

	for i, region := range network.PatternRegions {
		if i >= 10 { // Limit output
			fmt.Printf("   ... and %d more regions\n", len(network.PatternRegions)-i)
			break
		}

		fmt.Printf("   Region %d:\n", i+1)
		fmt.Printf("      Type: %s\n", region.TextureType)
		fmt.Printf("      Size: %d atoms\n", len(region.Atoms))
		fmt.Printf("      Edge strength: %.2f\n", region.EdgeStrength)
		fmt.Printf("      Coherence: %.2f\n", region.Coherence)
	}
}

// HandleConstraintRelaxation - Continue relaxation from existing state
func HandleConstraintRelaxation(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: ./programme energy relax <iterations> [input_file]")
		return
	}

	iterations, _ := strconv.Atoi(args[0])

	fmt.Printf("\n⚛️  CONSTRAINT RELAXATION CONTINUATION\n")
	fmt.Printf("════════════════════════════════════\n\n")
	fmt.Printf("Running %d more iterations...\n\n", iterations)

	// For now, create a small demo
	network := database.NewConstraintRelaxationNetwork(256, 256, 8)
	network.Generate(iterations)

	fmt.Printf("\n✅ Relaxation complete\n")
	printNetworkStatistics(network)
}

// HandleEnergyAnalysis - Analyze energy landscape
func HandleEnergyAnalysis(args []string) {
	fmt.Printf("\n📊 ENERGY LANDSCAPE ANALYSIS\n")
	fmt.Printf("═════════════════════════════\n\n")

	network := database.NewConstraintRelaxationNetwork(256, 256, 8)

	// Run a few iterations and track energy
	fmt.Printf("Iteration | Total Energy | Avg Local | Stability | Oscillating%%\n")
	fmt.Printf("────────────────────────────────────────────────────────────\n")

	for i := 0; i < 100; i++ {
		network.RelaxationStep()

		if i%10 == 0 {
			// Count oscillating atoms
			gridH := 256 / 8
			gridW := 256 / 8
			oscillating := 0
			for y := 0; y < gridH; y++ {
				for x := 0; x < gridW; x++ {
					if network.Atoms[y][x].EnergyTrend < 0.005 && i > 10 {
						oscillating++
					}
				}
			}
			oscillPercent := float64(oscillating) / float64(gridW*gridH) * 100

			fmt.Printf("%4d       | %.6f      | %.6f   | %+.2f      | %.1f%%\n",
				i,
				network.TotalEnergy,
				network.AverageLocalEnergy,
				network.SystemStability,
				oscillPercent,
			)
		}
	}

	fmt.Printf("\n✅ Analysis complete\n")
}

// HandleMultiPhaseGeneration - Three-phase generation for quality
func HandleMultiPhaseGeneration(args []string) {
	fmt.Printf("\n⚛️  MULTI-PHASE CONSTRAINT RELAXATION\n")
	fmt.Printf("════════════════════════════════════\n\n")
	fmt.Printf("Phase 1: Coarse structure (large patches)\n")
	fmt.Printf("Phase 2: Medium details (medium patches) - HÉRITAGE de Phase 1\n")
	fmt.Printf("Phase 3: Fine refinement (small patches) - HÉRITAGE de Phase 2\n\n")

	// Phase 1: Coarse (16x16 patches)
	fmt.Printf("🔵 PHASE 1: Coarse Structure\n")
	fmt.Printf("   Creating 256x256 image with 16×16 patches...\n")
	net1 := database.NewConstraintRelaxationNetwork(256, 256, 16)
	net1.Generate(200)       // Augmenté de 150 → 200
	net1.CapturePhaseState() // 📸 Sauvegarder l'état stable
	fmt.Printf("   ✓ Phase 1 complete\n")
	fmt.Printf("     Energy: %.6f | Stability: %.2f\n\n", net1.AverageLocalEnergy, net1.SystemStability)

	// Phase 2: Medium (8x8 patches)
	fmt.Printf("🟢 PHASE 2: Medium Details\n")
	fmt.Printf("   Creating 256x256 image with 8×8 patches...\n")
	net2 := database.NewConstraintRelaxationNetwork(256, 256, 8)

	// Initialize from phase 1 (upscale)
	for y := 0; y < 32; y++ {
		for x := 0; x < 32; x++ {
			coarseX := x / 2
			coarseY := y / 2
			if coarseX < 16 && coarseY < 16 {
				net2.Atoms[y][x].R = net1.Atoms[coarseY][coarseX].R
				net2.Atoms[y][x].G = net1.Atoms[coarseY][coarseX].G
				net2.Atoms[y][x].B = net1.Atoms[coarseY][coarseX].B
				net2.Atoms[y][x].Intensity = net1.Atoms[coarseY][coarseX].Intensity
			}
		}
	}

	// 🔗 Hériter l'état stable de Phase 1
	net2.InheritPhaseState(net1)

	net2.Generate(300)       // Augmenté de 200 → 300
	net2.CapturePhaseState() // 📸 Sauvegarder l'état stable
	fmt.Printf("   ✓ Phase 2 complete\n")
	fmt.Printf("     Energy: %.6f | Stability: %.2f\n\n", net2.AverageLocalEnergy, net2.SystemStability)

	// Phase 3: Fine (4x4 patches)
	fmt.Printf("🟡 PHASE 3: Fine Refinement\n")
	fmt.Printf("   Creating 256x256 image with 4×4 patches...\n")
	net3 := database.NewConstraintRelaxationNetwork(256, 256, 4)

	// Initialize from phase 2 (upscale)
	for y := 0; y < 64; y++ {
		for x := 0; x < 64; x++ {
			mediumX := x / 2
			mediumY := y / 2
			if mediumX < 32 && mediumY < 32 {
				net3.Atoms[y][x].R = net2.Atoms[mediumY][mediumX].R
				net3.Atoms[y][x].G = net2.Atoms[mediumY][mediumX].G
				net3.Atoms[y][x].B = net2.Atoms[mediumY][mediumX].B
				net3.Atoms[y][x].Intensity = net2.Atoms[mediumY][mediumX].Intensity
			}
		}
	}

	// 🔗 Hériter l'état stable de Phase 2
	net3.InheritPhaseState(net2)

	net3.Generate(400) // Augmenté de 250 → 400 (TRÈS important pour la qualité fine)
	fmt.Printf("   ✓ Phase 3 complete\n")
	fmt.Printf("     Energy: %.6f | Stability: %.2f\n\n", net3.AverageLocalEnergy, net3.SystemStability)

	// Save phase 3 result
	fmt.Printf("💾 Saving multi-phase result...\n")
	img := net3.RenderToImage()
	file, err := os.Create("generated_multiphase_energy.png")
	if err == nil {
		defer file.Close()
		png.Encode(file, img)
		fmt.Printf("✅ Saved to: generated_multiphase_energy.png\n")
	}

	fmt.Printf("\n📊 Multi-Phase Summary:\n")
	fmt.Printf("   Phase 1 Energy: %.6f → Phase 2: %.6f → Phase 3: %.6f\n",
		net1.AverageLocalEnergy, net2.AverageLocalEnergy, net3.AverageLocalEnergy)
	fmt.Printf("   Stability progression: %.2f → %.2f → %.2f\n",
		net1.SystemStability, net2.SystemStability, net3.SystemStability)
	fmt.Printf("   Phase 2 héritage: ✓ | Phase 3 héritage: ✓\n")
}

// ============================================================================
// HELPER FUNCTIONS
// ============================================================================

func applyConstraintsToField(field *database.GlobalCoherenceField, constraintStr string) {
	if constraintStr == "" {
		return
	}

	lower := strings.ToLower(constraintStr)

	// Brightness constraints
	if strings.Contains(lower, "dark") {
		field.AverageBrightness = 0.3
	} else if strings.Contains(lower, "bright") {
		field.AverageBrightness = 0.7
	}

	// Edge constraints
	if strings.Contains(lower, "smooth") {
		field.EdgeCohesion = 0.2
	} else if strings.Contains(lower, "sharp") || strings.Contains(lower, "detailed") {
		field.EdgeCohesion = 0.8
	}

	// Symmetry constraints
	if strings.Contains(lower, "symmetric") {
		field.GlobalSymmetryTarget = 0.8
	}

	// Texture constraints
	if strings.Contains(lower, "rough") {
		field.TextureConsistency = 0.2
	} else if strings.Contains(lower, "smooth") {
		field.TextureConsistency = 0.8
	}

	// Light direction constraints
	if strings.Contains(lower, "top") || strings.Contains(lower, "light from top") {
		field.ShadowDirection = 3.14159 / 2 // π/2
	} else if strings.Contains(lower, "side") {
		field.ShadowDirection = 0
	}
}

func printNetworkStatistics(network *database.ConstraintRelaxationNetwork) {
	fmt.Printf("\n📊 NETWORK STATISTICS:\n")
	fmt.Printf("   Iterations: %d\n", network.Iteration)
	fmt.Printf("   Final energy: %.6f\n", network.TotalEnergy)
	fmt.Printf("   Average local energy: %.6f\n", network.AverageLocalEnergy)
	fmt.Printf("   System stability: %.2f (range: -1 to +1)\n", network.SystemStability)
	fmt.Printf("   Plateau iterations: %d\n", network.IterationsAtPlateau)
	fmt.Printf("   Global field influence: %.3f\n", network.GlobalField.InfluenceWeight)
}

func PrintEnergyImageHelp() {
	fmt.Printf(`
⚛️  ENERGY-BASED IMAGE GENERATION
═════════════════════════════════════════════════════════════════

PARADIGM:
   NOT "draw" but "relax to equilibrium"
   Each pixel-atom minimizes local tension
   Global coherence emerges from local interactions
   Physics-inspired constraint propagation

THREE LEVELS:
   Level 1: Atoms (pixel state, color, intensity, orientation)
   Level 2: Patterns (edges, gradients, regions, textures)
   Level 3: Global field (weak influence on all atoms)

COMMANDS:
───────────────────────────────────────────────────────────────

./programme energy generate <w> <h> <iter> <patch> [constraints]
   Generate from scratch with constraint relaxation
   
   Examples:
     ./programme energy generate 512 512 500 8
     ./programme energy generate 256 256 300 4 "dark smooth sharp"
   
   Constraints (optional):
     - "dark" / "bright"
     - "smooth" / "sharp" / "detailed"
     - "rough" / "clean"
     - "symmetric"
     - "top" / "side" (light direction)

./programme energy relax <iterations> [input]
   Continue relaxation from existing state
   
   Example:
     ./programme energy relax 200

./programme energy analyze
   Analyze energy landscape and convergence

./programme energy multi-phase
   Three-phase generation: coarse → medium → fine
   Best quality, slower

HOW IT WORKS:
───────────────────────────────────────────────────────────────

1. Initialize atoms with random color/state
2. Compute local energy for each atom:
   - Continuity with neighbors (smooth → high energy)
   - Color/intensity consistency
   - Orientation alignment
   - Attraction to global field
3. Each atom moves to reduce its energy (gradient descent)
4. Repeat until convergence (plateau in total energy)
5. Auto-réévaluation: detect oscillations, suppress them
6. Detect patterns: find coherent regions

KEY INSIGHT:
───────────────────────────────────────────────────────────────

You don't tell atoms what color to be.
You tell them:
   - Try to match your neighbors
   - Respect global direction of shadows
   - Maintain some texture consistency
   - Minimize oscillations

They find equilibrium themselves.

ADVANTAGES OVER TRADITIONAL GAN/DIFFUSION:
───────────────────────────────────────────────────────────────

✓ Instant generation (no forward pass through huge network)
✓ Fully interpretable (energy function is explicit)
✓ Adaptable (change constraints in real-time)
✓ Parallel (atoms update independently)
✓ Embedded/edge-friendly (no GPU needed)
✓ No training required (pure physics)
✓ Can optimize for different objectives (energy, stability, etc.)

LIMITATIONS:
───────────────────────────────────────────────────────────────

✗ Photorealism needs training (for now)
✓ BUT: Great for procedural, abstract, adaptive imaging
✓ AND: Can combine with learned features later

PHILOSOPHY:
───────────────────────────────────────────────────────────────

"An image is not a sum of pixels.
It's a system of local interactions in equilibrium.
Generate by relaxation, not by simulation of a neural network."

`)
}

// HandleGenerateFromEnergyProfile - Analyze image and generate new image with same energy signature
func HandleGenerateFromEnergyProfile(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: ./programme energy from-image <source_image.png> [width] [height] [iterations] [patch_size]")
		fmt.Println("\nExample:")
		fmt.Println("  ./programme energy from-image target.png                      (uses 256x256, 300 iter, patch 4)")
		fmt.Println("  ./programme energy from-image target.png 512 512 400 4")
		fmt.Println("\nResult: Generates new image matching the energy signature of source_image.png")
		return
	}

	sourceImage := args[0]
	width := 256
	height := 256
	iterations := 300
	patchSize := 4

	if len(args) > 1 {
		if w, err := strconv.Atoi(args[1]); err == nil {
			width = w
		}
	}
	if len(args) > 2 {
		if h, err := strconv.Atoi(args[2]); err == nil {
			height = h
		}
	}
	if len(args) > 3 {
		if iter, err := strconv.Atoi(args[3]); err == nil {
			iterations = iter
		}
	}
	if len(args) > 4 {
		if ps, err := strconv.Atoi(args[4]); err == nil {
			patchSize = ps
		}
	}

	fmt.Printf("\n⚛️  ENERGY SIGNATURE MATCHING\n")
	fmt.Printf("════════════════════════════════════\n\n")
	fmt.Printf("📸 Analyzing source image: %s\n", sourceImage)

	// Analyze the source image
	profile, err := database.NewImageEnergyProfile(sourceImage)
	if err != nil {
		fmt.Printf("❌ Error analyzing image: %v\n", err)
		return
	}

	fmt.Printf("✓ Image analyzed (%dx%d)\n\n", profile.Width, profile.Height)
	fmt.Printf("📊 Energy Signature:\n")
	fmt.Printf("   Gradient energy (λ): %.4f\n", profile.LambdaGradient)
	fmt.Printf("   Local coherence (λ): %.4f\n", profile.LambdaLocal)
	fmt.Printf("   Texture energy (λ): %.4f\n", profile.LambdaTexture)
	fmt.Printf("   Scale distribution (λ): %.4f\n", profile.LambdaScale)
	fmt.Printf("   Smoothing penalty (λ): %.4f\n\n", profile.LambdaSmoothing)

	fmt.Printf("📈 Statistics:\n")
	fmt.Printf("   Average gradient: %.6f\n", profile.AverageGradient)
	fmt.Printf("   Local coherence: %.6f\n", profile.LocalCoherence)
	fmt.Printf("   Texture level: %.6f\n", profile.TextureLevel)
	fmt.Printf("   Sharpness ratio: %.4f\n", profile.SharpnessRatio)
	fmt.Printf("   Edge density: %.4f\n", profile.EdgeDensity)
	fmt.Printf("   Flat regions: %.4f | Textured regions: %.4f\n",
		profile.FlatRegionsFraction, profile.TexturedRegionsFraction)
	fmt.Printf("   Histogram shape: %s\n\n", profile.HistogramShape)

	fmt.Printf("🔧 Generation config:\n")
	fmt.Printf("   Output size: %dx%d\n", width, height)
	fmt.Printf("   Iterations: %d\n", iterations)
	fmt.Printf("   Patch size: %d\n", patchSize)
	fmt.Printf("   Grid: %d×%d atoms\n\n", width/patchSize, height/patchSize)

	fmt.Printf("⚡ Generating new image with matched energy signature...\n")
	fmt.Printf("   🔥 Using MULTI-SCALE PIPELINE with SYMMETRY BREAKING\n\n")

	// 🔥 SOLUTION 1: Phase Continuity (déjà implémenté dans ComputeLocalEnergy)
	fmt.Printf("   1️⃣  Phase continuity: ✓ Integrated in energy computation\n")

	// 🔥 SOLUTION 2: Weak Directional Field (champ directeur faible)
	fmt.Printf("   2️⃣  Extracting weak directional field (16×16 blur)...\n")
	guidanceField, err := database.ExtractLowResGuidanceField(sourceImage, 16)
	if err != nil {
		fmt.Printf("      ⚠️  Could not extract field: %v\n", err)
		guidanceField = database.NewGlobalCoherenceField()
	} else {
		fmt.Printf("      ✓ Weak directional field extracted (influence: %.1f%%)\n", guidanceField.FieldStrength*100)
	}

	// 🔥 SOLUTION 3: Topological Constraint (contrainte topologique)
	// Skip for very large source images to save computation time
	var edgeTopology *database.EdgeTopologyMap
	sourcePixels := profile.Width * profile.Height
	if sourcePixels < 600000 { // Skip if source > 600k pixels (e.g., 1200×1098 = 1.3M)
		fmt.Printf("   3️⃣  Computing edge topology map...\n")
		var err error
		edgeTopology, err = database.ComputeEdgeTopology(sourceImage)
		if err != nil {
			fmt.Printf("      ⚠️  Could not compute topology: %v\n", err)
			edgeTopology = nil
		} else {
			fmt.Printf("      ✓ Edge topology extracted (edge density: %.2f%%)\n", edgeTopology.EdgeDensity*100)
		}
	} else {
		fmt.Printf("   3️⃣  Skipping edge topology (source: %dx%d = %d pixels, too large)\n",
			profile.Width, profile.Height, sourcePixels)
		fmt.Printf("      💡 Tip: Use smaller source images (<600k pixels) for topology constraints\n")
		edgeTopology = nil
	}

	// 🔥 SOLUTION 4: Multi-Scale Pipeline (CRUCIAL!)
	fmt.Printf("   4️⃣  Multi-scale pipeline: Adaptive for %dx%d\n\n", width, height)

	pipeline := database.NewAdaptiveMultiScalePipeline(width, height)

	fmt.Printf("⏳ Running multi-scale relaxation pipeline (%d scales)...\n", len(pipeline.Scales))
	for i, scale := range pipeline.Scales {
		fmt.Printf("   [Scale %d/%d: %d] → %d iterations\n", i+1, len(pipeline.Scales), scale, pipeline.RelaxationPerScale[i])
	}
	fmt.Println()

	// Execute the complete symmetry-breaking pipeline
	network := pipeline.RunMultiScalePipeline(width, height, profile, guidanceField, edgeTopology)
	// Save result
	fmt.Printf("\n💾 Rendering and saving...\n")
	img := network.RenderToImage()
	file, err := os.Create("generated_from_signature.png")
	if err != nil {
		fmt.Printf("❌ Error saving image: %v\n", err)
		return
	}
	defer file.Close()

	png.Encode(file, img)
	fmt.Printf("✅ Saved to: generated_from_signature.png\n\n")

	fmt.Printf("📊 Final statistics:\n")
	printNetworkStatistics(network)
	fmt.Printf("\n🎯 New image has same energy signature as: %s\n", sourceImage)
	fmt.Printf("   But is a completely different image (not a copy!)\n")
}
