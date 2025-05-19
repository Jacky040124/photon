package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type OpenAIResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

type FormattedResearch struct {
	Summary     string
	KeyPoints   []string
	SourceLinks []string
}

func OpenAI(question string) (string, error) {
	godotenv.Load()
	var openaiPayload map[string]interface{}
	var openaiKey string

	openaiKey = os.Getenv("OPENAI_API_KEY")
	openaiPayload = map[string]interface{}{
		"model": "gpt-3.5-turbo",
		"messages": []map[string]string{
			{"role": "system", "content": "You are a research assistant that provides concise, factual information with source links. Use subtle and aesthetically pleasing emojis where appropriate to enhance readability and engagement. Format the output as plain text suitable for a terminal, avoiding markdown or excessive indentation."},
			{"role": "user", "content": question + "\n\nProvide a concise summary using relevant emojis, 3-5 key points, and credible source links. Ensure the output is plain text and terminal-friendly."},
		},
	}

	openaiBody, _ := json.Marshal(openaiPayload)
	openaiReq, _ := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(openaiBody))
	openaiReq.Header.Set("Authorization", "Bearer "+openaiKey)
	openaiReq.Header.Set("Content-Type", "application/json")
	openaiResp, err := http.DefaultClient.Do(openaiReq)

	if err != nil {
		return "", err
	}

	defer openaiResp.Body.Close()
	body, _ := io.ReadAll(openaiResp.Body)

	var response OpenAIResponse
	err = json.Unmarshal(body, &response)
	if err != nil || len(response.Choices) == 0 {
		return string(body), nil // Return raw response if parsing fails
	}

	return response.Choices[0].Message.Content, nil
}

func FormatResearch(query string) FormattedResearch {
	content, err := OpenAI(query)
	if err != nil {
		return FormattedResearch{
			Summary: "Error fetching research: " + err.Error(),
		}
	}

	var result FormattedResearch

	// Simple parsing of the response
	lines := strings.Split(content, "\n")
	inSummary := false
	inKeyPoints := false
	inSources := false

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if strings.Contains(strings.ToLower(line), "summary") {
			inSummary = true
			inKeyPoints = false
			inSources = false
			continue
		} else if strings.Contains(strings.ToLower(line), "key point") {
			inSummary = false
			inKeyPoints = true
			inSources = false
			continue
		} else if strings.Contains(strings.ToLower(line), "source") {
			inSummary = false
			inKeyPoints = false
			inSources = true
			continue
		}

		if line == "" {
			continue
		}

		if inSummary && result.Summary == "" {
			result.Summary = line
		} else if inKeyPoints {
			// Remove numbers and bullets at start
			cleanLine := strings.TrimLeft(line, "0123456789-.*• ")
			if cleanLine != "" {
				result.KeyPoints = append(result.KeyPoints, cleanLine)
			}
		} else if inSources {
			if strings.Contains(line, "http") {
				result.SourceLinks = append(result.SourceLinks, line)
			}
		}
	}

	// If parsing failed, use the content as summary
	if result.Summary == "" {
		result.Summary = content
	}

	return result
}

func PrintFormattedResearch(research FormattedResearch) {
	fmt.Println("\n✨ === DIG RESEARCH RESULTS === ✨")
	fmt.Println("\n✨ SUMMARY:")
	fmt.Println(research.Summary)

	if len(research.KeyPoints) > 0 {
		fmt.Println("\n💡 KEY POINTS:")
		for i, point := range research.KeyPoints {
			fmt.Printf("➤ %d. %s\n", i+1, point)
		}
	}

	if len(research.SourceLinks) > 0 {
		fmt.Println("\n🔗 SOURCES:")
		for _, source := range research.SourceLinks {
			fmt.Println("➤ " + source)
		}
	}

	fmt.Println("\n✨ ========================== ✨")
}

func Perplexity(question string) {
	godotenv.Load()
	perplexityKey := os.Getenv("PERPLEXITY_API_KEY")

	perplexityPayload := map[string]interface{}{
		"model": "sonar",
		"messages": []map[string]string{
			{"role": "user", "content": question},
		},
	}

	perplexityBody, _ := json.Marshal(perplexityPayload)
	perplexityReq, _ := http.NewRequest("POST", "https://api.perplexity.ai/chat/completions", bytes.NewBuffer(perplexityBody))
	perplexityReq.Header.Set("Authorization", "Bearer "+perplexityKey)
	perplexityReq.Header.Set("Content-Type", "application/json")
	perplexityResp, err := http.DefaultClient.Do(perplexityReq)
	if err != nil {
		fmt.Println("Perplexity error:", err)
	} else {
		defer perplexityResp.Body.Close()
		body, _ := io.ReadAll(perplexityResp.Body)
		fmt.Println("Perplexity response:")
		fmt.Println(string(body))
	}
}

func main() {
	var question string

	if len(os.Args) < 2 {
		fmt.Println("Please provide a querry")
		return
	}

	godotenv.Load()
	question = os.Args[1]

	research := FormatResearch(question)
	PrintFormattedResearch(research)
}
