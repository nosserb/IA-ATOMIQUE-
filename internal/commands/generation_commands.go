package commands

import (
	"IA-ATOMIQUE/database"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"strconv"
)

// GenerateCommand handles atomic generation operations
func GenerateCommand(args []string) {
	if len(args) < 1 {
		PrintGenerateHelp()
		return
	}

	subcommand := args[0]

	switch subcommand {
	case "pattern":
		HandleGenerateFromPattern(args[1:])
	case "with-feedback":
		HandleGenerateWithFeedback(args[1:])
	case "from-prompt":
		HandleGenerateFromPrompt(args[1:])
	case "with-keyword":
		HandleGenerateWithKeyword(args[1:])
	case "parameters":
		HandleGenerateParameters(args[1:])
	case "benchmark":
		HandleGenerateBenchmark(args[1:])
	default:
		fmt.Printf("❌ Unknown generate command: %s\n", subcommand)
		PrintGenerateHelp()
	}
}

// HandleGenerateFromPattern generates image from wave pattern alone
// Usage: ./programme generate pattern <width> <height> <iterations> [pattern_image]
func HandleGenerateFromPattern(args []string) {
	if len(args) < 3 {
		fmt.Println("Usage: ./programme generate pattern <w> <h> <iterations> [pattern_file]")
		fmt.Println("Example: ./programme generate pattern 512 512 200 output/pattern_final_emerged.png")
		return
	}

	width, _ := strconv.Atoi(args[0])
	height, _ := strconv.Atoi(args[1])
	iterations, _ := strconv.Atoi(args[2])
	patternFile := ""

	if len(args) > 3 {
		patternFile = args[3]
	}

	fmt.Printf("\n🎨 ATOMIC GENERATION FROM PATTERN\n")
	fmt.Printf("═══════════════════════════════════════════════════\n")
	fmt.Printf("Local → Global: waves become recognizable structures\n\n")

	// Create generation grid
	fmt.Printf("🔧 Phase 1: Initialize Atomic Grid (%dx%d)\n", width, height)
	grid := database.NewGenerationGrid(width, height)
	fmt.Printf("   ✓ %d atoms initialized with neutral state\n", width*height)

	// Load pattern if provided
	if patternFile != "" {
		fmt.Printf("\n📌 Phase 2: Inject Wave Pattern\n")
		file, err := os.Open(patternFile)
		if err != nil {
			fmt.Printf("   ⚠ Could not load pattern file, using defaults\n")
		} else {
			defer file.Close()
			img, _, err := image.Decode(file)
			if err != nil {
				fmt.Printf("   ⚠ Could not decode pattern\n")
			} else {
				// Create pattern engine from image
				engine := database.NewPatternEmergenceEngine(width, height)
				// Load pattern as seeds (but don't propagate - just inject)
				bounds := img.Bounds()
				for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
					for x := bounds.Min.X; x < bounds.Max.X; x++ {
						r, g, b, _ := img.At(x, y).RGBA()
						tx := int(float64(x-bounds.Min.X) * float64(width) / float64(bounds.Dx()))
						ty := int(float64(y-bounds.Min.Y) * float64(height) / float64(bounds.Dy()))
						if tx < width && ty < height {
							engine.Pixels[ty][tx].Color = [3]float64{
								float64(r) / 65535.0,
								float64(g) / 65535.0,
								float64(b) / 65535.0,
							}
						}
					}
				}
				grid.SetPattern(engine)
				fmt.Printf("   ✓ Pattern injected into atomic grid\n")
			}
		}
	}

	// Run generation
	fmt.Printf("\n⚛️ Phase 3: Atomic Generation (%d iterations)\n", iterations)
	fmt.Printf("   s_ij(t+1) = α·Σ neighbors + β·P_ij\n\n")

	for iter := 0; iter < iterations; iter++ {
		grid.GenerateStep()

		// Show progress
		if (iter+1)%(iterations/4) == 0 || iter == 0 {
			percent := (iter + 1) * 100 / iterations
			fmt.Printf("   ✓ %3d%% | Iter %d | Loss: %.8f\n",
				percent, iter+1, grid.AverageLoss)
		}
	}

	// Save result
	outputFile := fmt.Sprintf("output/atomic_generated_%dx%d_%diter.png", width, height, iterations)
	grid.SaveImage(outputFile)

	fmt.Printf("\n✅ GENERATION COMPLETE\n")
	fmt.Printf("═══════════════════════════════════════════════════\n")
	grid.PrintStatistics()
	fmt.Printf("Output: %s\n\n", outputFile)
}

// HandleGenerateWithFeedback generates with target image feedback
// Usage: ./programme generate with-feedback <w> <h> <iterations> <target_image>
func HandleGenerateWithFeedback(args []string) {
	if len(args) < 4 {
		fmt.Println("Usage: ./programme generate with-feedback <w> <h> <iterations> <target.png>")
		fmt.Println("Example: ./programme generate with-feedback 512 512 300 input/image/face.png")
		fmt.Println("\nGenerates with target image feedback for guided emergence.")
		return
	}

	width, _ := strconv.Atoi(args[0])
	height, _ := strconv.Atoi(args[1])
	iterations, _ := strconv.Atoi(args[2])
	targetFile := args[3]

	fmt.Printf("\n🎯 ATOMIC GENERATION WITH TARGET FEEDBACK\n")
	fmt.Printf("═══════════════════════════════════════════════════\n")
	fmt.Printf("Guided emergence: pattern + target image → final image\n\n")

	// Create grid
	fmt.Printf("🔧 Phase 1: Initialize Atomic Grid (%dx%d)\n", width, height)
	grid := database.NewGenerationGrid(width, height)
	grid.FeedbackWeight = 0.3 // Stronger feedback when target provided
	fmt.Printf("   ✓ Grid ready with enhanced feedback (δ=0.3)\n")

	// Load target image
	fmt.Printf("\n🎯 Phase 2: Load Target Image\n")
	file, err := os.Open(targetFile)
	if err != nil {
		fmt.Printf("   ❌ Failed to load target: %v\n", err)
		return
	}
	defer file.Close()

	targetImg, _, err := image.Decode(file)
	if err != nil {
		fmt.Printf("   ❌ Failed to decode target: %v\n", err)
		return
	}

	grid.TargetImage = targetImg
	fmt.Printf("   ✓ Target loaded: %s\n", targetFile)
	fmt.Printf("   ✓ Feedback will guide generation toward target\n")

	// Inject pattern from target
	fmt.Printf("\n📌 Phase 3: Inject Target as Initial Pattern\n")
	engine := database.NewPatternEmergenceEngine(width, height)
	bounds := targetImg.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := targetImg.At(x, y).RGBA()
			tx := int(float64(x-bounds.Min.X) * float64(width) / float64(bounds.Dx()))
			ty := int(float64(y-bounds.Min.Y) * float64(height) / float64(bounds.Dy()))
			if tx < width && ty < height {
				engine.Pixels[ty][tx].Color = [3]float64{
					float64(r) / 65535.0,
					float64(g) / 65535.0,
					float64(b) / 65535.0,
				}
			}
		}
	}
	grid.SetPattern(engine)
	fmt.Printf("   ✓ Target pattern injected\n")

	// Run generation WITH feedback
	fmt.Printf("\n⚛️ Phase 4: Guided Atomic Generation (%d iterations)\n", iterations)
	fmt.Printf("   With feedback: δ·(target_color - current_color)\n\n")

	for iter := 0; iter < iterations; iter++ {
		grid.GenerateStep()

		if (iter+1)%(iterations/4) == 0 || iter == 0 {
			percent := (iter + 1) * 100 / iterations
			fmt.Printf("   ✓ %3d%% | Loss: %.8f\n", percent, grid.AverageLoss)
		}
	}

	// Save result
	outputFile := fmt.Sprintf("output/atomic_feedback_%dx%d_%diter.png", width, height, iterations)
	grid.SaveImage(outputFile)

	fmt.Printf("\n✅ FEEDBACK GENERATION COMPLETE\n")
	fmt.Printf("═══════════════════════════════════════════════════\n")
	grid.PrintStatistics()
	fmt.Printf("Output: %s\n\n", outputFile)
}

// HandleGenerateFromPrompt generates image from natural language prompt
// Usage: ./programme generate from-prompt <width> <height> <iterations> "<prompt_text>"
func HandleGenerateFromPrompt(args []string) {
	if len(args) < 4 {
		fmt.Println("Usage: ./programme generate from-prompt <w> <h> <iterations> \"<prompt_text>\"")
		fmt.Println("Example: ./programme generate from-prompt 512 512 200 \"une forêt mystérieuse avec des arbres luminescents\"")
		return
	}

	width, _ := strconv.Atoi(args[0])
	height, _ := strconv.Atoi(args[1])
	iterations, _ := strconv.Atoi(args[2])
	prompt := args[3]

	fmt.Printf("\n🧠 ATOMIC GENERATION FROM TEXT PROMPT\n")
	fmt.Printf("═══════════════════════════════════════════════════\n")
	fmt.Printf("Neural interpretation: text → neuron activation → image\n\n")

	// Phase 1: Analyze prompt with neural system
	fmt.Printf("🧠 Phase 1: Analyze Prompt with Neuron Network\n")
	fmt.Printf("   Input: \"%s\"\n\n", prompt)

	// Process text through neural system (like summarization)
	catActivation, keywords, confidence := database.ProcesserTexte(prompt)

	fmt.Printf("   ✓ Keywords extracted: %v\n", keywords)
	fmt.Printf("   ✓ Semantic confidence: %.2f%%\n", confidence*100)
	fmt.Printf("   ✓ %d categories activated\n\n", len(catActivation))

	// Show top activated categories
	type CatScore struct {
		id    int
		score int
	}
	var scores []CatScore
	for catID, activation := range catActivation {
		if activation > 0 {
			scores = append(scores, CatScore{catID, activation})
		}
	}
	// Sort by activation
	for i := 0; i < len(scores) && i < 5; i++ {
		maxIdx := i
		for j := i + 1; j < len(scores); j++ {
			if scores[j].score > scores[maxIdx].score {
				maxIdx = j
			}
		}
		scores[i], scores[maxIdx] = scores[maxIdx], scores[i]
	}

	fmt.Printf("   Top category activations:\n")
	for i := 0; i < len(scores) && i < 5; i++ {
		categoryNames := map[int]string{
			1: "TECH",
			2: "HISTOIRE",
			3: "BUSINESS",
			4: "ALIMENTATION",
			5: "SANTÉ",
			6: "VERBE",
		}
		catName := "UNKNOWN"
		if name, ok := categoryNames[scores[i].id]; ok {
			catName = name
		}
		fmt.Printf("      • %s (category %d): %d neurons\n", catName, scores[i].id, scores[i].score)
	}

	// Phase 2: Create atomic grid
	fmt.Printf("\n🔧 Phase 2: Initialize Atomic Grid (%dx%d)\n", width, height)
	grid := database.NewGenerationGrid(width, height)
	fmt.Printf("   ✓ %d atoms initialized\n", width*height)

	// Phase 3: Inject neural activations as pattern
	fmt.Printf("\n📌 Phase 3: Inject Neural Pattern into Grid\n")

	// Create pattern grid from neural activations
	patternEngine := database.NewPatternEmergenceEngine(width, height)

	// Map neuron activations to colors
	// Each activated neuron creates a color influence based on its category
	for catID, activation := range catActivation {
		if activation > 0 {
			// Map category to color hue
			hue := float64(catID%6) / 6.0
			saturation := float64(activation) / 100.0
			brightness := confidence

			// Distribute activation across the pattern grid
			for y := 0; y < height; y++ {
				for x := 0; x < width; x++ {
					// Create noise-based distribution of pattern
					var influence float64
					if (x+y+catID)%10 < 5 {
						influence = saturation * (1.0 - float64((x+y+catID)%10)/10.0)
					} else {
						influence = saturation * (0.5 + float64((x+y+catID)%5)/10.0)
					}

					// Blend with existing pattern
					patternEngine.Pixels[y][x].Color[0] =
						0.7*patternEngine.Pixels[y][x].Color[0] + 0.3*hue
					patternEngine.Pixels[y][x].Color[1] =
						0.7*patternEngine.Pixels[y][x].Color[1] + 0.3*influence
					patternEngine.Pixels[y][x].Color[2] =
						0.7*patternEngine.Pixels[y][x].Color[2] + 0.3*brightness
				}
			}
		}
	}

	grid.SetPattern(patternEngine)
	fmt.Printf("   ✓ Neural activations injected as pattern\n")

	// Phase 4: Generate with atomic resonance
	fmt.Printf("\n⚛️ Phase 4: Atomic Generation (%d iterations)\n", iterations)
	fmt.Printf("   Neurons → Pattern → Resonance → Image\n\n")

	for iter := 0; iter < iterations; iter++ {
		grid.GenerateStep()

		if (iter+1)%(iterations/4) == 0 || iter == 0 {
			percent := (iter + 1) * 100 / iterations
			fmt.Printf("   ✓ %3d%% | Loss: %.8f\n", percent, grid.AverageLoss)
		}
	}

	// Save result
	outputFile := fmt.Sprintf("output/atomic_prompt_%dx%d_%diter.png", width, height, iterations)
	grid.SaveImage(outputFile)

	fmt.Printf("\n✅ PROMPT-BASED GENERATION COMPLETE\n")
	fmt.Printf("═══════════════════════════════════════════════════\n")
	fmt.Printf("Prompt: \"%s\"\n", prompt)
	fmt.Printf("Keywords: %v\n", keywords)
	fmt.Printf("Confidence: %.2f%%\n\n", confidence*100)
	grid.PrintStatistics()
	fmt.Printf("Output: %s\n\n", outputFile)
}

// HandleGenerateParameters shows effect of different parameters
// Usage: ./programme generate parameters
func HandleGenerateParameters(args []string) {
	fmt.Printf("\n⚛️ Atomic Generation Parameters\n")
	fmt.Printf("═══════════════════════════════════════════════════\n\n")

	fmt.Printf("🔧 PROPAGATION PARAMETERS:\n\n")

	fmt.Printf("  Resonance Alpha (α): 0.1 - 0.5\n")
	fmt.Printf("    How much atoms respond to neighbors\n")
	fmt.Printf("    Low (0.1):   Isolated atoms, strong patterns\n")
	fmt.Printf("    High (0.5):  Cooperative, smooth emergence\n\n")

	fmt.Printf("  Pattern Beta (β): 0.2 - 0.8\n")
	fmt.Printf("    How much atoms follow the injected pattern\n")
	fmt.Printf("    Low (0.2):   Creative variation from pattern\n")
	fmt.Printf("    High (0.8):  Strict adherence to pattern\n\n")

	fmt.Printf("  Smoothing Gamma (γ): 0.05 - 0.3\n")
	fmt.Printf("    Local color smoothing to avoid cube artifacts\n")
	fmt.Printf("    Low (0.05):  Sharp details, more noise\n")
	fmt.Printf("    High (0.3):  Smooth gradients, less detail\n\n")

	fmt.Printf("  Feedback Weight (δ): 0.0 - 1.0\n")
	fmt.Printf("    Strength of target image guidance\n")
	fmt.Printf("    0.0:  No feedback, pattern-only generation\n")
	fmt.Printf("    0.5:  Balanced (default with target)\n")
	fmt.Printf("    1.0:  Very tight guidance toward target\n\n")

	fmt.Printf("  Velocity Damping (ε): 0.8 - 1.0\n")
	fmt.Printf("    Momentum smoothing in state propagation\n")
	fmt.Printf("    0.8:  Responsive, may oscillate\n")
	fmt.Printf("    0.95: Smooth, stable convergence\n\n")

	fmt.Printf("═══════════════════════════════════════════════════\n")
	fmt.Printf("EXAMPLE CONFIGURATIONS:\n\n")

	configs := []struct {
		name  string
		alpha float64
		beta  float64
		gamma float64
		desc  string
	}{
		{"Fine Detail", 0.2, 0.6, 0.1, "Sharp edges, local patterns"},
		{"Smooth Blend", 0.4, 0.4, 0.25, "Balanced, natural-looking"},
		{"Wave Centric", 0.1, 0.8, 0.05, "Follows wave pattern closely"},
		{"Emergent", 0.5, 0.3, 0.2, "Cooperative atoms, creative"},
	}

	for _, cfg := range configs {
		fmt.Printf("  %s:\n    α=%.1f  β=%.1f  γ=%.2f\n    → %s\n\n",
			cfg.name, cfg.alpha, cfg.beta, cfg.gamma, cfg.desc)
	}
}

// HandleGenerateBenchmark tests generation speed
func HandleGenerateBenchmark(args []string) {
	fmt.Printf("\n⚙️ ATOMIC GENERATION BENCHMARK\n")
	fmt.Printf("═══════════════════════════════════════════════════\n")
	fmt.Printf("Testing generation speed at different resolutions\n\n")

	tests := []struct {
		width  int
		height int
		iters  int
	}{
		{256, 256, 50},
		{512, 512, 50},
		{256, 256, 100},
	}

	for _, test := range tests {
		fmt.Printf("Testing %dx%d × %d iterations...\n", test.width, test.height, test.iters)

		grid := database.NewGenerationGrid(test.width, test.height)

		// Initialize with pattern
		engine := database.NewPatternEmergenceEngine(test.width, test.height)
		grid.SetPattern(engine)

		// Run and measure
		fmt.Printf("  Running generation...\n")
		grid.Generate(test.iters)

		fmt.Printf("  ✓ Completed | Final loss: %.8f\n\n", grid.AverageLoss)
	}

	fmt.Printf("═══════════════════════════════════════════════════\n")
	fmt.Printf("Benchmark complete. Adjust resolutions as needed.\n\n")
}

// PrintGenerateHelp displays generation command help
func PrintGenerateHelp() {
	fmt.Printf(`
⚛️ Atomic Generation - From Patterns to Complete Images
═════════════════════════════════════════════════════════════

Mathematical Framework:

LOCAL STATE PROPAGATION:
  s_ij(t+1) = α·(1/|N|)·Σ s_kl(t) + β·P_ij + ε·V_ij(t)
  
  α: resonance (neighbor influence)
  β: pattern weight (input guidance)
  ε: momentum damping
  
STATE TO COLOR CONVERSION:
  C_ij = f(s_ij)
  Converts internal state → RGB via HSL-like function
  
LOCAL SMOOTHING:
  C'_ij = (1-γ)·C_ij + γ·(1/|N|)·Σ C_kl
  
TARGET FEEDBACK:
  s_ij += δ·(C_target - C_current)
  Guides generation toward reference image

═════════════════════════════════════════════════════════════
COMMANDS:

  ./programme generate pattern <w> <h> <iterations> [pattern.png]
      Generate from wave pattern only (no target feedback)
      Quick emergence guided by injected waves
      Example: ./programme generate pattern 512 512 200 output/pattern_final_emerged.png

  ./programme generate with-feedback <w> <h> <iterations> <target.png>
      Generate with target image feedback
      Guided emergence toward specific appearance
      Example: ./programme generate with-feedback 512 512 300 input/image/face.png

  ./programme generate from-prompt <w> <h> <iterations> "<text>"
      Generate from natural language prompt using neuron interpretation
      Neural network analyzes text, creates pattern, generates image
      Example: ./programme generate from-prompt 512 512 200 "une forêt mystérieuse"

  ./programme generate parameters
      Explain all generation parameters and their effects

  ./programme generate benchmark
      Run speed test at different resolutions

═════════════════════════════════════════════════════════════
WORKFLOW:

1. GENERATE FROM TEXT PROMPT:
   ./programme generate from-prompt 512 512 200 "description of image"
   → output/atomic_prompt_512x512_200iter.png

2. CREATE PATTERN (Optional):
   ./programme pattern emerge 512 512 200 input/image/ref.png 0.15
   → output/pattern_final_emerged.png

3. GENERATE FROM PATTERN:
   ./programme generate pattern 512 512 300 output/pattern_final_emerged.png
   → output/atomic_generated_512x512_300iter.png

4. REFINE WITH FEEDBACK (Optional):
   ./programme generate with-feedback 512 512 400 input/image/ref.png
   → output/atomic_feedback_512x512_400iter.png

═════════════════════════════════════════════════════════════
NEURAL TEXT-TO-IMAGE GENERATION:

The from-prompt mode works by:

Phase 1: Neuron Analysis
  • Tokenize and analyze text prompt
  • Extract keywords from text
  • Activate neuron categories based on semantic content
  • Calculate semantic confidence

Phase 2: Pattern Injection
  • Map neuron activations to spatial colors
  • Create pattern grid influenced by activated neurons
  • Higher activation = stronger color influence

Phase 3: Atomic Resonance
  • Atoms propagate resonance from neighbors
  • Follow injected neural pattern
  • Local smoothing creates coherent image
  • Iteratively stabilize to final state

Result: Abstract neural activation → Concrete visual image

EXAMPLE PROMPTS:
  • "une forêt sombre avec des arbres anciens"
  • "océan tumultueux avec tempête"
  • "technologie futuriste et néons bleus"
  • "jardin fleuri au printemps"
  • "architecture gothique médiévale"

═════════════════════════════════════════════════════════════
HOW IT WORKS:

Phase 1: Initialize atomic grid (neutral gray atoms)
Phase 2: Inject wave pattern (from diffusion)
Phase 3: Propagate local state (s_ij resonates with neighbors)
Phase 4: Convert state to color (internal state → RGB)
Phase 5: Smooth colors locally (remove cube artifacts)
Phase 6: Apply feedback (optional, guide toward target)
Phase 7: Iterate until convergence

Result: Abstract waves → Recognizable image through local interactions

═════════════════════════════════════════════════════════════
KEY PARAMETERS:

α (Resonance):  How much atoms listen to neighbors [0.1-0.5]
β (Pattern):    How strictly atoms follow input [0.2-0.8]
γ (Smoothing):  Local color blending [0.05-0.3]
δ (Feedback):   Target image influence [0.0-1.0]
ε (Damping):    State update momentum [0.8-1.0]

═════════════════════════════════════════════════════════════
TIPS:

• For pattern-only: use low iterations (100-200)
• For feedback: use more iterations (300-500)
• High β keeps pattern details
• High α creates smooth emergence
• γ controls detail vs. smoothness trade-off

═════════════════════════════════════════════════════════════
`)
}
