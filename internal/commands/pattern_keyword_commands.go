package commands

import (
	"IA-ATOMIQUE/database"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"strconv"
)

// HandleGenerateWithKeyword génère une image en cherchant tous les patterns avec un keyword
// Usage: ./programme generate with-keyword 512 512 100 "forest"
func HandleGenerateWithKeyword(args []string) {
	if len(args) < 4 {
		fmt.Println("Usage: ./programme generate with-keyword <WIDTH> <HEIGHT> <ITERATIONS> \"<KEYWORD>\"")
		return
	}

	width, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Println("Invalid width:", err)
		return
	}

	height, err := strconv.Atoi(args[1])
	if err != nil {
		fmt.Println("Invalid height:", err)
		return
	}

	iterations, err := strconv.Atoi(args[2])
	if err != nil {
		fmt.Println("Invalid iterations:", err)
		return
	}

	keyword := args[3]

	fmt.Printf("\n🔍 Génération avec keyword '%s'...\n", keyword)
	fmt.Printf("   Résolution: %d×%d\n", width, height)
	fmt.Printf("   Itérations: %d\n", iterations)

	// Charger les patterns
	patterns, err := LoadPatternDatabase()
	if err != nil {
		fmt.Println("❌ Cannot load patterns:", err)
		return
	}

	if len(patterns) == 0 {
		fmt.Println("❌ No patterns found in database")
		return
	}

	fmt.Printf("   Patterns disponibles: %d\n", len(patterns))

	// Chercher les patterns matchant
	matched := database.SearchPatternsByKeyword(patterns, keyword)

	if len(matched) == 0 {
		fmt.Printf("❌ No patterns found for keyword '%s'\n", keyword)
		fmt.Println("Available keywords:")
		for _, p := range patterns {
			fmt.Printf("   Pattern '%s':\n", p.PatternID)
			fmt.Printf("     Tags: %v\n", p.Tags)
			fmt.Printf("     Keywords: %v\n", p.Keywords)
		}
		return
	}

	fmt.Printf("✅ Found %d matching patterns:\n", len(matched))
	for _, p := range matched {
		fmt.Printf("   - %s\n", p.PatternID)
	}

	// Combiner les patterns
	fmt.Println("\n📐 Combining patterns...")
	combinedImage := database.CombinePatternsByKeyword(patterns, keyword, width, height)

	// Convertir en image PNG
	fmt.Println("🖼️  Converting to PNG...")
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			c := combinedImage[y][x]
			r := uint8(c[0] * 255)
			g := uint8(c[1] * 255)
			b := uint8(c[2] * 255)
			img.Set(x, y, color.RGBA{r, g, b, 255})
		}
	}

	// Sauvegarder PNG
	fmt.Println("💾 Exporting to PNG...")
	file, err := os.Create("generated_image_keyword.png")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	err = png.Encode(file, img)
	if err != nil {
		fmt.Println("Error encoding PNG:", err)
		return
	}

	fmt.Println("✅ Done!")
	fmt.Println("   Output: generated_image_keyword.png")
}

// LoadPatternDatabase charge tous les patterns depuis le répertoire input/ ou utilise des défauts
func LoadPatternDatabase() ([]*database.PatternMathematical, error) {
	// Chercher les images dans le répertoire input/
	imageDir := "input"

	// Si le répertoire n'existe pas, retourner patterns par défaut
	if _, err := os.Stat(imageDir); os.IsNotExist(err) {
		fmt.Printf("⚠️  Répertoire '%s' non trouvé\n", imageDir)
		return CreateDefaultPatterns(), nil
	}

	// Charger depuis les fichiers images
	patterns, err := database.LoadPatternDatabaseFromImages(imageDir)
	if err != nil {
		fmt.Printf("⚠️  Erreur lors du chargement des images: %v\n", err)
		fmt.Println("   Utilisation des patterns par défaut...")
		return CreateDefaultPatterns(), nil
	}

	if len(patterns) == 0 {
		fmt.Println("⚠️  Aucun pattern trouvé dans 'input/'")
		fmt.Println("   Utilisation des patterns par défaut...")
		return CreateDefaultPatterns(), nil
	}

	fmt.Printf("✅ %d patterns chargés depuis les fichiers images\n", len(patterns))
	database.PrintPatternDatabaseInfo(patterns)
	return patterns, nil
}

// CreateDefaultPatterns crée les patterns par défaut si aucune image n'existe
func CreateDefaultPatterns() []*database.PatternMathematical {
	patterns := make([]*database.PatternMathematical, 0)

	// Pattern 1: Sunset (orange/red)
	p1 := &database.PatternMathematical{
		PatternID:      "sunset",
		Width:          512,
		Height:         512,
		BasisFunctions: 20,
		BasisType:      "fourier",
		Tags:           []string{"sunset", "orange", "warm", "landscape"},
		Keywords: map[string]float64{
			"sunset":  1.0,
			"sun":     0.9,
			"orange":  0.8,
			"warm":    0.7,
			"evening": 0.8,
		},
		Coefficients: make([]float64, 60),
	}
	// Remplir avec coefficients d'exemple (orangé)
	for i := 0; i < 20; i++ {
		p1.Coefficients[i] = 0.8 + float64(i%5)*0.02    // R: haut
		p1.Coefficients[20+i] = 0.5 + float64(i%5)*0.02 // G: moyen
		p1.Coefficients[40+i] = 0.2 + float64(i%5)*0.01 // B: bas
	}
	patterns = append(patterns, p1)

	// Pattern 2: Forest (green)
	p2 := &database.PatternMathematical{
		PatternID:      "forest",
		Width:          512,
		Height:         512,
		BasisFunctions: 20,
		BasisType:      "fourier",
		Tags:           []string{"forest", "green", "nature", "trees"},
		Keywords: map[string]float64{
			"forest": 1.0,
			"green":  0.9,
			"trees":  0.8,
			"nature": 0.8,
			"leaf":   0.7,
		},
		Coefficients: make([]float64, 60),
	}
	// Remplir avec coefficients d'exemple (vert)
	for i := 0; i < 20; i++ {
		p2.Coefficients[i] = 0.2 + float64(i%5)*0.02    // R: bas
		p2.Coefficients[20+i] = 0.7 + float64(i%5)*0.03 // G: haut
		p2.Coefficients[40+i] = 0.2 + float64(i%5)*0.02 // B: bas
	}
	patterns = append(patterns, p2)

	// Pattern 3: Ocean (blue)
	p3 := &database.PatternMathematical{
		PatternID:      "ocean",
		Width:          512,
		Height:         512,
		BasisFunctions: 20,
		BasisType:      "fourier",
		Tags:           []string{"ocean", "blue", "water", "sea"},
		Keywords: map[string]float64{
			"ocean": 1.0,
			"blue":  0.95,
			"water": 0.9,
			"sea":   0.85,
			"wave":  0.7,
		},
		Coefficients: make([]float64, 60),
	}
	// Remplir avec coefficients d'exemple (bleu)
	for i := 0; i < 20; i++ {
		p3.Coefficients[i] = 0.1 + float64(i%5)*0.02    // R: très bas
		p3.Coefficients[20+i] = 0.3 + float64(i%5)*0.02 // G: bas-moyen
		p3.Coefficients[40+i] = 0.8 + float64(i%5)*0.03 // B: haut
	}
	patterns = append(patterns, p3)

	// Pattern 4: Sky (light blue)
	p4 := &database.PatternMathematical{
		PatternID:      "sky",
		Width:          512,
		Height:         512,
		BasisFunctions: 20,
		BasisType:      "fourier",
		Tags:           []string{"sky", "cloud", "blue", "air"},
		Keywords: map[string]float64{
			"sky":   1.0,
			"cloud": 0.8,
			"blue":  0.7,
			"light": 0.8,
			"air":   0.6,
		},
		Coefficients: make([]float64, 60),
	}
	// Remplir avec coefficients d'exemple (bleu clair)
	for i := 0; i < 20; i++ {
		p4.Coefficients[i] = 0.6 + float64(i%5)*0.02     // R: moyen-haut
		p4.Coefficients[20+i] = 0.7 + float64(i%5)*0.02  // G: moyen-haut
		p4.Coefficients[40+i] = 0.95 + float64(i%5)*0.01 // B: très haut
	}
	patterns = append(patterns, p4)

	fmt.Println("✅ Patterns par défaut chargés (4 patterns)")
	return patterns
}
