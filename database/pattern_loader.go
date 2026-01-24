package database

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// LoadPatternDatabaseFromImages charge les patterns depuis des fichiers images
// Fichier "rose.png" → Pattern avec tags ["rose", "flower", "red"]
// Fichier "coeur.png" → Pattern avec tags ["coeur", "heart", "love"]
// Etc.
func LoadPatternDatabaseFromImages(imageDir string) ([]*PatternMathematical, error) {
	var patterns []*PatternMathematical

	files, err := os.ReadDir(imageDir)
	if err != nil {
		return nil, fmt.Errorf("cannot read directory: %v", err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		filename := file.Name()
		ext := strings.ToLower(filepath.Ext(filename))

		// Accepter PNG et JPG
		if ext != ".png" && ext != ".jpg" && ext != ".jpeg" {
			continue
		}

		// Extraire le nom sans extension
		baseName := strings.TrimSuffix(filename, ext)

		imagePath := filepath.Join(imageDir, filename)

		fmt.Printf("📷 Loading pattern from: %s\n", filename)

		// Extraire coefficients Fourier depuis l'image
		pattern, err := ExtractPatternFromImage(imagePath, 20)
		if err != nil {
			fmt.Printf("   ⚠️  Error extracting pattern: %v (using defaults)\n", err)
			// Créer pattern par défaut avec coefficients placeholder
			pattern = &PatternMathematical{
				PatternID:      baseName,
				Width:          512,
				Height:         512,
				BasisFunctions: 20,
				BasisType:      "fourier",
				Tags:           GenerateTagsFromFilename(baseName),
				Keywords:       GenerateKeywordsFromFilename(baseName),
				Coefficients:   make([]float64, 60),
			}
			// Remplir avec coefficients d'exemple
			for i := 0; i < len(pattern.Coefficients); i++ {
				pattern.Coefficients[i] = 0.5
			}
		} else {
			// Pattern extrait avec succès - mettre à jour l'ID et ajouter les tags sémantiques
			pattern.PatternID = baseName
			pattern.Tags = GenerateTagsFromFilename(baseName)
			pattern.Keywords = GenerateKeywordsFromFilename(baseName)
		}

		patterns = append(patterns, pattern)
		fmt.Printf("   ✅ Added pattern '%s' with tags: %v\n", baseName, pattern.Tags)
	}

	return patterns, nil
}

// GenerateTagsFromFilename génère les tags à partir du nom du fichier
// "rose" → ["rose", "flower", "red"]
// "coeur" → ["coeur", "heart", "love", "red"]
// "soleil" → ["soleil", "sun", "yellow", "warm"]
func GenerateTagsFromFilename(filename string) []string {
	filename = strings.ToLower(filename)

	tagMap := map[string][]string{
		// Fleurs
		"rose":      {"rose", "flower", "red", "nature"},
		"tulip":     {"tulip", "flower", "colorful"},
		"daisy":     {"daisy", "flower", "white"},
		"sunflower": {"sunflower", "flower", "yellow"},

		// Formes / Objets
		"coeur":  {"coeur", "heart", "love", "red"},
		"heart":  {"heart", "love", "red"},
		"star":   {"star", "bright", "yellow"},
		"moon":   {"moon", "night", "white"},
		"sun":    {"sun", "warm", "yellow", "bright"},
		"soleil": {"soleil", "sun", "warm", "yellow"},

		// Paysages
		"ocean":    {"ocean", "water", "blue", "wave"},
		"ciel":     {"ciel", "sky", "blue", "air"},
		"sky":      {"sky", "blue", "air", "cloud"},
		"forest":   {"forest", "green", "trees", "nature"},
		"foret":    {"foret", "forest", "green", "trees"},
		"mountain": {"mountain", "nature", "gray"},
		"desert":   {"desert", "sand", "warm", "yellow"},

		// Ambiances
		"sunset":  {"sunset", "orange", "warm", "evening"},
		"coucher": {"coucher", "sunset", "orange", "warm"},
		"night":   {"night", "dark", "blue", "stars"},
		"nuit":    {"nuit", "night", "dark", "blue"},
	}

	if tags, exists := tagMap[filename]; exists {
		return tags
	}

	// Par défaut: utiliser le nom du fichier comme tag
	return []string{filename}
}

// GenerateKeywordsFromFilename génère les keywords avec confidence scores
// "rose" → {"rose": 1.0, "flower": 0.9, "red": 0.8}
func GenerateKeywordsFromFilename(filename string) map[string]float64 {
	filename = strings.ToLower(filename)

	keywordMap := map[string]map[string]float64{
		// Fleurs
		"rose": {
			"rose":   1.0,
			"flower": 0.9,
			"red":    0.8,
			"petal":  0.7,
			"nature": 0.6,
		},
		"sunflower": {
			"sunflower": 1.0,
			"flower":    0.9,
			"yellow":    0.85,
			"bright":    0.7,
			"nature":    0.6,
		},

		// Formes / Objets
		"coeur": {
			"coeur": 1.0,
			"heart": 0.95,
			"love":  0.8,
			"red":   0.7,
			"shape": 0.6,
		},
		"heart": {
			"heart": 1.0,
			"love":  0.8,
			"red":   0.7,
			"shape": 0.6,
		},
		"star": {
			"star":   1.0,
			"bright": 0.8,
			"yellow": 0.7,
			"night":  0.5,
		},
		"sun": {
			"sun":    1.0,
			"warm":   0.9,
			"yellow": 0.85,
			"bright": 0.8,
		},
		"soleil": {
			"soleil": 1.0,
			"sun":    0.95,
			"warm":   0.9,
			"yellow": 0.85,
			"bright": 0.8,
		},
		"moon": {
			"moon":  1.0,
			"night": 0.8,
			"white": 0.7,
			"dark":  0.6,
		},

		// Paysages
		"ocean": {
			"ocean":  1.0,
			"water":  0.95,
			"blue":   0.9,
			"wave":   0.8,
			"nature": 0.7,
		},
		"ciel": {
			"ciel":  1.0,
			"sky":   0.95,
			"blue":  0.85,
			"air":   0.7,
			"light": 0.8,
		},
		"sky": {
			"sky":   1.0,
			"blue":  0.9,
			"cloud": 0.7,
			"air":   0.6,
			"light": 0.7,
		},
		"forest": {
			"forest": 1.0,
			"green":  0.95,
			"tree":   0.85,
			"nature": 0.8,
			"leaf":   0.7,
		},
		"foret": {
			"foret":  1.0,
			"forest": 0.95,
			"green":  0.95,
			"tree":   0.85,
			"nature": 0.8,
		},

		// Ambiances
		"sunset": {
			"sunset":    1.0,
			"orange":    0.9,
			"warm":      0.85,
			"evening":   0.7,
			"landscape": 0.6,
		},
		"coucher": {
			"coucher": 1.0,
			"sunset":  0.95,
			"orange":  0.9,
			"warm":    0.85,
			"evening": 0.7,
		},
		"night": {
			"night": 1.0,
			"dark":  0.9,
			"blue":  0.7,
			"stars": 0.8,
		},
		"nuit": {
			"nuit":  1.0,
			"night": 0.95,
			"dark":  0.9,
			"blue":  0.7,
		},
	}

	if keywords, exists := keywordMap[filename]; exists {
		return keywords
	}

	// Par défaut: utiliser le nom du fichier avec confiance 1.0
	return map[string]float64{
		filename: 1.0,
	}
}

// PrintPatternDatabaseInfo affiche les patterns chargés
func PrintPatternDatabaseInfo(patterns []*PatternMathematical) {
	fmt.Println("\n📚 PATTERN DATABASE INFO")
	fmt.Println("═══════════════════════════════════════════════")
	fmt.Printf("Total patterns: %d\n\n", len(patterns))

	for _, p := range patterns {
		fmt.Printf("📷 Pattern: %s\n", p.PatternID)
		fmt.Printf("   Tags:     %v\n", p.Tags)
		fmt.Printf("   Keywords: ")
		for kw, conf := range p.Keywords {
			fmt.Printf("%s(%.1f) ", kw, conf)
		}
		fmt.Printf("\n\n")
	}
}
