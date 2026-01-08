package database

import (
	"math"
	"sort"
	"strings"
)

// ============================================================================
// COMPRESSION REWARD SYSTEM
// ============================================================================
// Système de récompense pour atteindre exactement le ratio de compression demandé
// Extrait graduellement jusqu'à atteindre le target

// CompressionTarget contient le target et les métriques
type CompressionTarget struct {
	TargetChars  int     // Nombre de chars demandés
	TargetRatio  float64 // Ratio 0.1 = 10% du texte original
	RewardScore  float64 // Score de proximité (0-1)
	ActualRatio  float64 // Ratio réel atteint
	DeltaPercent float64 // Écart en pourcentage
	FinalSummary string  // Résumé final
}

// ExtractWithCompressionReward extrait le texte en visant un ratio exact
// Applique une "récompense" si on atteint le ratio demandé
func ExtractWithCompressionReward(sourceText string, targetRatio float64) CompressionTarget {
	result := CompressionTarget{
		TargetRatio: targetRatio,
		TargetChars: int(float64(len(sourceText)) * targetRatio),
	}

	// Diviser le texte en phrases
	sentences := SplitSentencesForExtraction(sourceText)

	if len(sentences) == 0 {
		result.FinalSummary = sourceText[:intMin(len(sourceText), result.TargetChars)]
		result.ActualRatio = float64(len(result.FinalSummary)) / float64(len(sourceText))
		result.RewardScore = 0.1
		return result
	}

	// Scorer les phrases par importance
	type scoredSentence struct {
		text  string
		score float64
	}

	var scored []scoredSentence
	for _, sent := range sentences {
		score := scoreImportance(sent, sourceText)
		scored = append(scored, scoredSentence{text: sent, score: score})
	}

	// Trier par score décroissant
	sort.Slice(scored, func(i, j int) bool {
		return scored[i].score > scored[j].score
	})

	// Extraire les phrases jusqu'à atteindre le target
	var selected []string
	currentLength := 0

	for _, sent := range scored {
		candidateLength := currentLength + len(sent.text) + 1 // +1 pour l'espace

		// Si on dépasse le target, évaluer si c'est ok
		if candidateLength > result.TargetChars {
			// Vérifier si on est proche du target
			if float64(currentLength)/float64(result.TargetChars) > 0.85 {
				// On est assez proche (85%+), ajouter quand même
				selected = append(selected, sent.text)
				currentLength = candidateLength
			} else if float64(candidateLength)/float64(result.TargetChars) < 1.2 {
				// Seulement 20% au-dessus du target, acceptable
				selected = append(selected, sent.text)
				currentLength = candidateLength
			}
			// Sinon on arrête
			if currentLength > result.TargetChars {
				break
			}
		} else {
			// Bien sous le target, ajouter
			selected = append(selected, sent.text)
			currentLength = candidateLength
		}
	}

	// Construire le résumé final (préserver ordre original)
	finalSentences := preserveOrderInExtraction(sourceText, selected)
	result.FinalSummary = strings.Join(finalSentences, " ")

	// Calculer les métriques
	result.ActualRatio = float64(len(result.FinalSummary)) / float64(len(sourceText))
	result.DeltaPercent = math.Abs((result.ActualRatio - result.TargetRatio) / result.TargetRatio * 100)

	// Calculer le RewardScore (proximité du target)
	// 1.0 = exact, 0.5 = ±50%, 0.0 = très loin
	if result.DeltaPercent < 5 {
		result.RewardScore = 1.0 // Excellent
	} else if result.DeltaPercent < 15 {
		result.RewardScore = 0.9 // Très bon
	} else if result.DeltaPercent < 25 {
		result.RewardScore = 0.8 // Bon
	} else if result.DeltaPercent < 50 {
		result.RewardScore = 0.6 // Acceptable
	} else {
		result.RewardScore = 0.4 // Mauvais
	}

	return result
}

// scoreImportance calcule l'importance d'une phrase
func scoreImportance(sentence string, context string) float64 {
	score := 0.0

	// 1. Longueur (phrases longues = plus d'info)
	wordCount := float64(len(strings.Fields(sentence)))
	score += math.Log(wordCount+1) * 0.1

	// 2. Fréquence des termes clés
	contextWords := strings.Fields(strings.ToLower(context))
	sentenceWords := strings.Fields(strings.ToLower(sentence))

	frequency := make(map[string]int)
	for _, word := range contextWords {
		frequency[word]++
	}

	for _, word := range sentenceWords {
		if freq, exists := frequency[word]; exists && freq > 1 {
			score += 0.05 * float64(freq)
		}
	}

	// 3. Présence de chiffres/dates (important pour factuel)
	if strings.ContainsAny(sentence, "0123456789") {
		score += 0.3
	}

	// 4. Pas trop de termes rares
	rareCount := 0
	for _, word := range sentenceWords {
		if frequency[word] == 0 {
			rareCount++
		}
	}
	if rareCount > 0 {
		score -= float64(rareCount) * 0.02
	}

	return math.Max(score, 0.1) // Score minimum
}

// SplitSentencesForExtraction divise le texte en phrases
func SplitSentencesForExtraction(text string) []string {
	sentences := strings.Split(text, ".")
	var result []string

	for _, sent := range sentences {
		sent = strings.TrimSpace(sent)
		if len(sent) > 20 { // Ignorer les phrases trop courtes
			result = append(result, sent)
		}
	}

	return result
}

// preserveOrderInExtraction réordonne les phrases extraites dans leur ordre original
func preserveOrderInExtraction(original string, extracted []string) []string {
	// Créer une map des phrases extraites pour lookup rapide
	extractedSet := make(map[string]bool)
	for _, sent := range extracted {
		extractedSet[strings.TrimSpace(sent)] = true
	}

	// Parcourir l'original et sélectionner dans l'ordre
	sentences := strings.Split(original, ".")
	var ordered []string

	for _, sent := range sentences {
		sent = strings.TrimSpace(sent)
		if extractedSet[sent] {
			ordered = append(ordered, sent)
		}
	}

	return ordered
}

// intMin retourne le min entre deux ints
func intMin(a, b int) int {
	if a < b {
		return a
	}
	return b
}
