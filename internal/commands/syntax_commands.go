package commands

import (
	"flag"
	"fmt"
	"strings"

	"github.com/nosserb/IA-ATOMIQUE-/database"
)

// syntaxCommand gère les commandes de syntaxe avancée
func SyntaxCommand(args []string) {
	if len(args) < 2 {
		fmt.Println("Usage: ./programme syntax <subcommand> <text>")
		fmt.Println("\nSubcommands:")
		fmt.Println("  analyze <text>           - Analyze syntax of a sentence")
		fmt.Println("  enhance <text>           - Generate better version of sentence")
		fmt.Println("  paragraph <text>         - Analyze paragraph structure")
		fmt.Println("  batch <file>             - Process entire file with enhancement")
		fmt.Println("  pos <text>               - Show POS tagging")
		return
	}

	subcommand := args[0]
	text := strings.Join(args[1:], " ")

	switch subcommand {
	case "analyze":
		cmdAnalyzeSyntax(text)
	case "enhance":
		cmdEnhanceSentence(text)
	case "paragraph":
		cmdAnalyzeParagraph(text)
	case "batch":
		cmdBatchProcess(text)
	case "pos":
		cmdShowPOS(text)
	default:
		fmt.Printf("Unknown subcommand: %s\n", subcommand)
	}
}

// cmdAnalyzeSyntax affiche l'analyse syntaxique
func cmdAnalyzeSyntax(sentence string) {
	fmt.Println("╔════════════════════════════════════════════════╗")
	fmt.Println("║  SYNTACTIC ANALYSIS - PHRASE                  ║")
	fmt.Println("╚════════════════════════════════════════════════╝")
	fmt.Printf("\nOriginal: %s\n\n", sentence)

	parsed := database.AnalyzeSyntax(sentence)

	fmt.Println("COMPONENTS:")
	fmt.Println("──────────────────────────────────────────────────")
	if parsed.MainVerb != nil {
		fmt.Printf("Main Verb:  %s (POS: %s)\n", parsed.MainVerb.Word, parsed.MainVerb.POS)
	} else {
		fmt.Println("Main Verb:  NONE")
	}

	if parsed.Subject != nil {
		fmt.Printf("Subject:    %s (POS: %s)\n", parsed.Subject.Word, parsed.Subject.POS)
	} else {
		fmt.Println("Subject:    NONE")
	}

	if len(parsed.Objects) > 0 {
		fmt.Print("Objects:    ")
		for i, obj := range parsed.Objects {
			if i > 0 {
				fmt.Print(", ")
			}
			fmt.Printf("%s", obj.Word)
		}
		fmt.Println()
	}

	if len(parsed.Modifiers) > 0 {
		fmt.Print("Modifiers:  ")
		for i, mod := range parsed.Modifiers {
			if i > 0 {
				fmt.Print(", ")
			}
			fmt.Printf("%s", mod.Word)
		}
		fmt.Println()
	}

	fmt.Println("\nSCORES:")
	fmt.Println("──────────────────────────────────────────────────")
	fmt.Printf("Grammar Score:  %.2f%% - ", parsed.GrammarScore*100)
	if parsed.GrammarScore > 0.8 {
		fmt.Println("✅ EXCELLENT")
	} else if parsed.GrammarScore > 0.6 {
		fmt.Println("✅ GOOD")
	} else {
		fmt.Println("⚠️  WEAK")
	}

	fmt.Printf("Style Score:    %.2f%% - ", parsed.StyleScore*100)
	if parsed.StyleScore > 0.7 {
		fmt.Println("✨ VIVID & VARIED")
	} else if parsed.StyleScore > 0.4 {
		fmt.Println("Moderate variation")
	} else {
		fmt.Println("⚠️  Repetitive")
	}

	fmt.Printf("FINAL Score:    %.2f%%\n", parsed.FinalScore*100)

	fmt.Println("\nDEPENDENCIES:")
	fmt.Println("──────────────────────────────────────────────────")
	for _, dep := range parsed.Dependencies {
		gov := parsed.Words[dep.Governor].Word
		dep_word := parsed.Words[dep.Dependent].Word
		fmt.Printf("  %s --[%s]--> %s (confidence: %.1f%%)\n",
			gov, dep.RelType, dep_word, dep.Confidence*100)
	}

	fmt.Println()
}

// cmdEnhanceSentence améliore une phrase
func cmdEnhanceSentence(sentence string) {
	fmt.Println("╔════════════════════════════════════════════════╗")
	fmt.Println("║  SENTENCE ENHANCEMENT - OPTIMIZATION           ║")
	fmt.Println("╚════════════════════════════════════════════════╝")
	fmt.Printf("\n📝 Original:  %s\n", sentence)

	enhanced, score := database.GenerateBetterSentence(sentence)

	fmt.Printf("✨ Enhanced:  %s\n", enhanced)
	fmt.Printf("📊 Score:     %.2f%%\n\n", score*100)

	// Analyser les deux versions
	origParsed := database.AnalyzeSyntax(sentence)
	enhParsed := database.AnalyzeSyntax(enhanced)

	fmt.Println("COMPARISON:")
	fmt.Println("──────────────────────────────────────────────────")
	fmt.Printf("Original  - Grammar: %.2f%%, Style: %.2f%%\n",
		origParsed.GrammarScore*100, origParsed.StyleScore*100)
	fmt.Printf("Enhanced  - Grammar: %.2f%%, Style: %.2f%%\n",
		enhParsed.GrammarScore*100, enhParsed.StyleScore*100)

	improvements := ""
	if enhParsed.GrammarScore > origParsed.GrammarScore {
		improvements += fmt.Sprintf("  ✅ Grammar improved by %.1f%%\n",
			(enhParsed.GrammarScore-origParsed.GrammarScore)*100)
	}
	if enhParsed.StyleScore > origParsed.StyleScore {
		improvements += fmt.Sprintf("  ✅ Style improved by %.1f%%\n",
			(enhParsed.StyleScore-origParsed.StyleScore)*100)
	}
	if len(enhParsed.Objects) > len(origParsed.Objects) {
		improvements += fmt.Sprintf("  ✅ Added %d more objects/complements\n",
			len(enhParsed.Objects)-len(origParsed.Objects))
	}

	if improvements != "" {
		fmt.Println("\nIMPROVEMENTS:")
		fmt.Println(improvements)
	}

	fmt.Println()
}

// cmdAnalyzeParagraph analyse un paragraphe
func cmdAnalyzeParagraph(paragraph string) {
	fmt.Println("╔════════════════════════════════════════════════╗")
	fmt.Println("║  PARAGRAPH STRUCTURE ANALYSIS                  ║")
	fmt.Println("╚════════════════════════════════════════════════╝")

	result := database.AnalyzeParagraphStructure(paragraph)

	fmt.Printf("\nSentences analyzed: %v\n", result["sentence_count"])
	fmt.Printf("Average Grammar Score: %.2f%%\n",
		result["avg_grammar"].(float64)*100)
	fmt.Printf("Average Style Score: %.2f%%\n",
		result["avg_style"].(float64)*100)
	fmt.Printf("Overall Coherence: %.2f%%\n\n",
		result["overall_score"].(float64)*100)

	analyses := result["sentences"].([]database.ParsedSentence)
	fmt.Println("INDIVIDUAL SENTENCES:")
	fmt.Println("──────────────────────────────────────────────────")

	for i, parsed := range analyses {
		status := "✅"
		if parsed.GrammarScore < 0.6 {
			status = "⚠️"
		}

		fmt.Printf("%d. %s Grammar:%.1f%%, Style:%.1f%%\n",
			i+1, status,
			parsed.GrammarScore*100,
			parsed.StyleScore*100)
		fmt.Printf("   %s\n", parsed.Original)
		fmt.Println()
	}
}

// cmdBatchProcess traite un fichier entier
func cmdBatchProcess(filepath string) {
	// Implémenter après
	fmt.Printf("Processing: %s\n", filepath)
	fmt.Println("(Implementation pending - requires file I/O integration)")
}

// cmdShowPOS affiche le POS tagging
func cmdShowPOS(sentence string) {
	fmt.Println("╔════════════════════════════════════════════════╗")
	fmt.Println("║  POS TAGGING - PART-OF-SPEECH                  ║")
	fmt.Println("╚════════════════════════════════════════════════╝")
	fmt.Printf("\nText: %s\n\n", sentence)

	tagged := database.TagMots(sentence)

	fmt.Println("TOKENS:")
	fmt.Println("──────────────────────────────────────────────────")
	fmt.Printf("%-15s %-12s %-15s %-8s\n", "Word", "POS", "Lemma", "Score")
	fmt.Println("──────────────────────────────────────────────────")

	for _, w := range tagged {
		fmt.Printf("%-15s %-12s %-15s %.1f%%\n",
			w.Word, w.POS, w.Lemma, w.Score*100)
	}

	fmt.Println()

	// Statistiques
	posCount := make(map[database.POSTag]int)
	for _, w := range tagged {
		posCount[w.POS]++
	}

	fmt.Println("POS DISTRIBUTION:")
	fmt.Println("──────────────────────────────────────────────────")
	for pos, count := range posCount {
		fmt.Printf("%-12s: %d\n", pos, count)
	}
	fmt.Println()
}

// Flag helper
func syntaxFlags(flagSet *flag.FlagSet) {
	// Placeholder pour future options
}
