package database

import (
	"encoding/json"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// PatternMetadata stores information about a discovered pattern
type PatternMetadata struct {
	ID              string         `json:"id"`       // Unique identifier
	Filename        string         `json:"filename"` // Original file
	Width           int            `json:"width"`    // Pattern dimensions
	Height          int            `json:"height"`
	Categories      map[string]int `json:"categories"`      // Activated neuron categories
	AverageColor    [3]float64     `json:"average_color"`   // RGB average
	Complexity      float64        `json:"complexity"`      // Pattern complexity (0-1)
	ContentSummary  string         `json:"content_summary"` // AI-generated description
	Keywords        []string       `json:"keywords"`        // Extracted keywords
	Confidence      float64        `json:"confidence"`      // Analysis confidence
	CreatedAt       string         `json:"created_at"`      // When indexed
	PatternDataHash string         `json:"pattern_hash"`    // Hash for pattern data
}

// PatternDatabase stores and manages indexed patterns
type PatternDatabase struct {
	Version  string                     `json:"version"`
	Patterns map[string]PatternMetadata `json:"patterns"`
	Index    []string                   `json:"index"` // Ordered list of pattern IDs
}

// PatternIndexer analyzes images and creates pattern database
type PatternIndexer struct {
	patterns      *PatternDatabase
	inputPath     string
	dbPath        string
	patternEngine *PatternEmergenceEngine
}

// NewPatternIndexer creates a new pattern indexer
func NewPatternIndexer(inputPath, dbPath string) *PatternIndexer {
	return &PatternIndexer{
		patterns:  &PatternDatabase{Version: "1.0", Patterns: make(map[string]PatternMetadata)},
		inputPath: inputPath,
		dbPath:    dbPath,
	}
}

// LoadDatabase loads existing pattern database from file
func (pi *PatternIndexer) LoadDatabase() error {
	data, err := os.ReadFile(pi.dbPath)
	if err != nil {
		if os.IsNotExist(err) {
			// Create new database if doesn't exist
			pi.patterns = &PatternDatabase{
				Version:  "1.0",
				Patterns: make(map[string]PatternMetadata),
				Index:    []string{},
			}
			return nil
		}
		return err
	}

	err = json.Unmarshal(data, pi.patterns)
	return err
}

// SaveDatabase saves pattern database to file
func (pi *PatternIndexer) SaveDatabase() error {
	data, err := json.MarshalIndent(pi.patterns, "", "  ")
	if err != nil {
		return err
	}

	// Create directory if needed
	dir := filepath.Dir(pi.dbPath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, 0755)
	}

	return os.WriteFile(pi.dbPath, data, 0644)
}

// IndexDirectory scans input directory and indexes all images
func (pi *PatternIndexer) IndexDirectory() error {
	fmt.Println("\n📚 PATTERN INDEXING ENGINE")
	fmt.Println("═══════════════════════════════════════════════════")
	fmt.Printf("Scanning input directory: %s\n", pi.inputPath)

	// Get list of image files
	entries, err := os.ReadDir(pi.inputPath)
	if err != nil {
		return fmt.Errorf("cannot read input directory: %w", err)
	}

	imageFiles := []string{}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		filename := entry.Name()
		ext := ""
		if idx := strings.LastIndex(filename, "."); idx >= 0 {
			ext = strings.ToLower(filename[idx:])
		}
		if ext == ".png" || ext == ".jpg" || ext == ".jpeg" {
			imageFiles = append(imageFiles, filename)
		}
	}

	if len(imageFiles) == 0 {
		fmt.Println("⚠ No image files found in input directory")
		return nil
	}

	fmt.Printf("Found %d image(s) to index\n\n", len(imageFiles))

	// Index each image
	indexedCount := 0
	for i, filename := range imageFiles {
		fmt.Printf("[%d/%d] Processing: %s\n", i+1, len(imageFiles), filename)

		filepath := filepath.Join(pi.inputPath, filename)
		metadata, err := pi.indexImage(filepath, filename)
		if err != nil {
			fmt.Printf("       ⚠ Error analyzing: %v\n", err)
			continue
		}

		pi.patterns.Patterns[metadata.ID] = metadata
		pi.patterns.Index = append(pi.patterns.Index, metadata.ID)
		indexedCount++

		// Display pattern info
		fmt.Printf("       ✓ ID: %s\n", metadata.ID)
		fmt.Printf("       Size: %dx%d\n", metadata.Width, metadata.Height)
		if len(metadata.Categories) > 0 {
			fmt.Printf("       Categories: ")
			for cat, count := range metadata.Categories {
				fmt.Printf("%s(%d) ", cat, count)
			}
			fmt.Println()
		}
		fmt.Printf("       Complexity: %.2f | Confidence: %.1f%%\n", metadata.Complexity, metadata.Confidence*100)
		fmt.Printf("       Keywords: %s\n", strings.Join(metadata.Keywords, ", "))
		fmt.Println()
	}

	fmt.Printf("═══════════════════════════════════════════════════\n")
	fmt.Printf("✓ Indexed %d pattern(s)\n", indexedCount)

	// Save database
	err = pi.SaveDatabase()
	if err != nil {
		fmt.Printf("⚠ Could not save database: %v\n", err)
		return nil
	}

	fmt.Printf("✓ Database saved to: %s\n\n", pi.dbPath)
	return nil
}

// indexImage analyzes a single image and creates metadata
func (pi *PatternIndexer) indexImage(filepath, filename string) (PatternMetadata, error) {
	// Load image
	file, err := os.Open(filepath)
	if err != nil {
		return PatternMetadata{}, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return PatternMetadata{}, err
	}

	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// Analyze colors and patterns
	filename_no_ext := filename
	if idx := strings.LastIndex(filename, "."); idx >= 0 {
		filename_no_ext = filename[:idx]
	}

	metadata := PatternMetadata{
		ID:         filename_no_ext,
		Filename:   filename,
		Width:      width,
		Height:     height,
		Categories: make(map[string]int),
		Keywords:   []string{},
		CreatedAt:  fmt.Sprintf("%d", os.Getenv("TIMESTAMP")),
	}

	// Extract colors and analyze
	var totalR, totalG, totalB float64
	colorFreq := make(map[[3]uint32]int)
	pixelCount := 0

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			if a > 0 {
				totalR += float64(r)
				totalG += float64(g)
				totalB += float64(b)
				colorFreq[[3]uint32{r, g, b}]++
				pixelCount++
			}
		}
	}

	// Calculate average color
	if pixelCount > 0 {
		metadata.AverageColor = [3]float64{
			totalR / float64(pixelCount),
			totalG / float64(pixelCount),
			totalB / float64(pixelCount),
		}
	}

	// Calculate complexity (color diversity)
	complexity := float64(len(colorFreq)) / float64(pixelCount)
	if complexity > 1.0 {
		complexity = 1.0
	}
	if complexity > 0.1 {
		complexity = complexity / 10.0
	}
	metadata.Complexity = complexity

	// Analyze content from filename and extract keywords
	metadata.Keywords = extractKeywords(filename)

	// Simulate category activation based on image analysis
	categories := analyzeImageContent(metadata.AverageColor, metadata.Complexity, filename)
	metadata.Categories = categories

	// Generate confidence based on analysis
	metadata.Confidence = 0.75 + (metadata.Complexity * 0.25)

	// Create content summary
	metadata.ContentSummary = generateContentSummary(metadata)

	// Create pattern hash
	metadata.PatternDataHash = generateHash(metadata)

	return metadata, nil
}

// extractKeywords extracts keywords from filename and creates related terms
func extractKeywords(filename string) []string {
	name := ""
	if idx := strings.LastIndex(filename, "."); idx >= 0 {
		name = filename[:idx]
	} else {
		name = filename
	}
	name = strings.ToLower(name)

	// Split by common separators
	separators := []string{"_", "-", " ", "."}
	keywords := []string{}
	current := name

	for _, sep := range separators {
		if strings.Contains(current, sep) {
			parts := strings.Split(current, sep)
			keywords = append(keywords, parts...)
			break
		}
	}

	// Clean up keywords (remove empty, filter short)
	cleaned := []string{}
	for _, kw := range keywords {
		kw = strings.TrimSpace(kw)
		if len(kw) > 2 {
			cleaned = append(cleaned, kw)
		}
	}

	return cleaned
}

// analyzeImageContent determines activated neuron categories from image
func analyzeImageContent(avgColor [3]float64, complexity float64, filename string) map[string]int {
	categories := make(map[string]int)

	// Analyze based on dominant color
	r, g, b := avgColor[0], avgColor[1], avgColor[2]

	// Map colors to categories
	if r > g && r > b {
		// Red/warm tones → HISTOIRE, BUSINESS
		categories["HISTOIRE"] = 3
		categories["BUSINESS"] = 2
	} else if g > r && g > b {
		// Green tones → ALIMENTATION, SANTÉ
		categories["ALIMENTATION"] = 3
		categories["SANTÉ"] = 2
	} else if b > r && b > g {
		// Blue tones → TECH
		categories["TECH"] = 4
	}

	// High complexity → more tech/business
	if complexity > 0.6 {
		oldVal := categories["TECH"]
		if oldVal < 2 {
			categories["TECH"] = 2
		}
	}

	// Analyze filename for keywords
	filename = strings.ToLower(filename)
	if strings.Contains(filename, "tech") || strings.Contains(filename, "robot") {
		if oldVal, ok := categories["TECH"]; !ok || oldVal < 4 {
			categories["TECH"] = 4
		}
	}
	if strings.Contains(filename, "food") || strings.Contains(filename, "meal") {
		if oldVal, ok := categories["ALIMENTATION"]; !ok || oldVal < 4 {
			categories["ALIMENTATION"] = 4
		}
	}
	if strings.Contains(filename, "historic") || strings.Contains(filename, "old") {
		if oldVal, ok := categories["HISTOIRE"]; !ok || oldVal < 4 {
			categories["HISTOIRE"] = 4
		}
	}

	return categories
}

// generateContentSummary creates a description of the pattern
func generateContentSummary(metadata PatternMetadata) string {
	summary := fmt.Sprintf("Pattern '%s' (%dx%d, complexity %.2f%%)",
		metadata.ID, metadata.Width, metadata.Height, metadata.Complexity*100)

	if len(metadata.Categories) > 0 {
		summary += " with categories: "
		cats := []string{}
		for cat, count := range metadata.Categories {
			cats = append(cats, fmt.Sprintf("%s(%d)", cat, count))
		}
		sort.Strings(cats)
		summary += strings.Join(cats, ", ")
	}

	if len(metadata.Keywords) > 0 {
		summary += fmt.Sprintf(" [%s]", strings.Join(metadata.Keywords, ", "))
	}

	return summary
}

// generateHash creates a hash of pattern data for caching
func generateHash(metadata PatternMetadata) string {
	hashStr := fmt.Sprintf("%s_%d_%d_%.2f",
		metadata.ID, metadata.Width, metadata.Height, metadata.Complexity)
	// Simple hash: first 16 chars of concatenation
	if len(hashStr) > 16 {
		hashStr = hashStr[:16]
	}
	return hashStr
}

// FindPatternByID retrieves a pattern by its ID
func (pi *PatternIndexer) FindPatternByID(id string) (PatternMetadata, bool) {
	pattern, exists := pi.patterns.Patterns[id]
	return pattern, exists
}

// FindPatternsByCategory finds all patterns with a specific category
func (pi *PatternIndexer) FindPatternsByCategory(category string) []PatternMetadata {
	results := []PatternMetadata{}
	for _, id := range pi.patterns.Index {
		if pattern, exists := pi.patterns.Patterns[id]; exists {
			if _, hasCategory := pattern.Categories[category]; hasCategory {
				results = append(results, pattern)
			}
		}
	}
	return results
}

// GetPatternDatabase returns the entire pattern database
func (pi *PatternIndexer) GetPatternDatabase() *PatternDatabase {
	return pi.patterns
}

// PrintPatternStats displays database statistics
func (pi *PatternIndexer) PrintPatternStats() {
	fmt.Println("\n📊 PATTERN DATABASE STATISTICS")
	fmt.Println("═══════════════════════════════════════════════════")
	fmt.Printf("Total patterns indexed: %d\n", len(pi.patterns.Patterns))

	if len(pi.patterns.Patterns) == 0 {
		fmt.Println("No patterns in database")
		return
	}

	// Category stats
	catFreq := make(map[string]int)
	totalComplexity := 0.0
	totalConfidence := 0.0

	for _, pattern := range pi.patterns.Patterns {
		totalComplexity += pattern.Complexity
		totalConfidence += pattern.Confidence
		for cat := range pattern.Categories {
			catFreq[cat]++
		}
	}

	fmt.Printf("\nAverage complexity: %.2f\n", totalComplexity/float64(len(pi.patterns.Patterns)))
	fmt.Printf("Average confidence: %.1f%%\n", (totalConfidence/float64(len(pi.patterns.Patterns)))*100)

	fmt.Println("\nCategory distribution:")
	categories := []string{}
	for cat := range catFreq {
		categories = append(categories, cat)
	}
	sort.Strings(categories)
	for _, cat := range categories {
		fmt.Printf("  %s: %d patterns\n", cat, catFreq[cat])
	}

	fmt.Println("═══════════════════════════════════════════════════\n")
}

// PatternMinInt returns minimum of two ints
func PatternMinInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// PatternMaxInt returns maximum of two ints
func PatternMaxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
