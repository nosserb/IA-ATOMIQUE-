package commands

import (
	"bufio"
	"fmt"
	"strings"

	"IA-ATOMIQUE/database"
	"os"
)

// StartLLMMode launches a chat-like loop that can run core commands or free-text analysis.
func StartLLMMode() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("\n==============================")
	fmt.Println(" IA-ATOMIQUE LLM chat mode")
	fmt.Println("==============================")
	fmt.Println("Type /help for commands. Type text without a slash for direct analysis. Type /quit to exit.")

	for {
		fmt.Print("llm> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("[INFO] End of input, leaving chat mode.")
			return
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		lower := strings.ToLower(line)
		if lower == "/quit" || lower == "quit" || lower == "exit" {
			fmt.Println("[INFO] Chat mode closed.")
			return
		}

		switch {
		case lower == "/help" || lower == "help":
			fmt.Print(`
Slash commands:
	/file <path>                analyze a file
	/text <content>             analyze inline text
	/resume <file> [ratio]      hybrid summary atomique+proba (default 0.55)
	/compare <file> [ratio]     compare summarizers
	/generate <file> [ratio]    generate summary (vector/atomic)
	/extract <file> [ratio]     extract key sentences
	/humanize <file> [-s|-p|-a] rewrite file (standard|pro|avance)
	/atomic <cmd> [...]         atomic engine (simulate, network-stats, benchmark)
	/ask <question>             quick Q&A using current pipeline
	/help                       show this help
	/quit                       leave chat mode
Free text without a slash is analyzed directly.
`)
			continue

		case strings.HasPrefix(lower, "/file "):
			path := strings.TrimSpace(line[len("/file "):])
			if path == "" {
				fmt.Println("Usage: /file <path>")
				continue
			}
			TraiterFichier(path)
			continue

		case strings.HasPrefix(lower, "/text "):
			text := strings.TrimSpace(line[len("/text "):])
			if text == "" {
				fmt.Println("Usage: /text <content>")
				continue
			}
			TraiterTexte(text, "llm-text")
			continue

		case strings.HasPrefix(lower, "/resume "):
			parts := strings.Fields(line)
			if len(parts) < 2 {
				fmt.Println("Usage: /resume <file> [ratio]")
				continue
			}
			ratio := "0.55"
			if len(parts) > 2 {
				ratio = parts[2]
			}
			ResumeHybridCommand([]string{parts[1], ratio})
			continue

		case strings.HasPrefix(lower, "/compare "):
			parts := strings.Fields(line)
			if len(parts) < 2 {
				fmt.Println("Usage: /compare <file> [compression]")
				continue
			}
			compression := 0.3
			if len(parts) > 2 {
				fmt.Sscanf(parts[2], "%f", &compression)
			}
			GenererComparatifCommand(parts[1], compression)
			continue

		case strings.HasPrefix(lower, "/generate "):
			parts := strings.Fields(line)
			if len(parts) < 2 {
				fmt.Println("Usage: /generate <file> [ratio]")
				continue
			}
			GenerateCommand(parts[1:])
			continue

		case strings.HasPrefix(lower, "/extract "):
			parts := strings.Fields(line)
			if len(parts) < 2 {
				fmt.Println("Usage: /extract <file> [ratio]")
				continue
			}
			ExtrairePhrasesClésCommand(parts[1:])
			continue

		case strings.HasPrefix(lower, "/humanize "):
			parts := strings.Fields(line)
			if len(parts) < 2 {
				fmt.Println("Usage: /humanize <file> [-s|-p|-a]")
				continue
			}
			style := "standard"
			path := ""
			for _, p := range parts[1:] {
				switch p {
				case "-s":
					style = "standard"
				case "-p":
					style = "professionnel"
				case "-a":
					style = "avance"
				default:
					path = p
				}
			}
			if path == "" {
				fmt.Println("Usage: /humanize <file> [-s|-p|-a]")
				continue
			}
			TraiterFichierHumanize(path, style)
			continue

		case strings.HasPrefix(lower, "/atomic"):
			parts := strings.Fields(line)
			if len(parts) < 2 {
				fmt.Println("Usage: /atomic simulate <iter> [atoms] | /atomic network-stats | /atomic benchmark")
				continue
			}
			args := append([]string{"programme"}, parts[1:]...)
			ParseSimulationArgs(args)
			continue

		case strings.HasPrefix(lower, "/ask "):
			question := strings.TrimSpace(line[len("/ask "):])
			if question == "" {
				fmt.Println("Usage: /ask <question>")
				continue
			}
			catAct, mots, conf := database.ProcesserTexte(question)
			var cat int
			var score int
			for c, s := range catAct {
				if s > score {
					score = s
					cat = c
				}
			}
			if cat > 0 {
				reponse := database.GenererReponse(cat, nil)
				fmt.Printf("\n[REPONSE]\n%s\n", reponse)
				if len(mots) > 0 {
					maxMots := 3
					if len(mots) < 3 {
						maxMots = len(mots)
					}
					fmt.Printf("Mots cles: %s\n", strings.Join(mots[:maxMots], ", "))
				}
				fmt.Printf("Confiance: %.0f%%\n", conf*100)
			} else {
				fmt.Println("Aucune categorie dominante detectee.")
			}
			continue
		}

		// Default: treat as free text for analysis
		TraiterTexte(line, "llm-free-text")
	}
}
