package pkg

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

type FormattedResponse struct {
	Summary     string
	KeyPoints   []string
	SourceLinks []string
}

type APIResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func FormatResearch(query string) FormattedResponse {
	content, err := CallLLMAPI(query)
	if err != nil {
		return FormattedResponse{
			Summary: fmt.Sprintf("Error fetching research: %s", err.Error()),
		}
	}

	var result FormattedResponse
	lines := strings.Split(content, "\n")

	// More robust parsing
	currentSection := ""
	summaryLines := []string{}

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if line == "" {
			continue
		}

		// Detect sections
		lowerLine := strings.ToLower(line)
		if strings.Contains(lowerLine, "summary") && strings.Contains(lowerLine, ":") {
			currentSection = "summary"
			continue
		} else if strings.Contains(lowerLine, "key point") && strings.Contains(lowerLine, ":") {
			currentSection = "keypoints"
			continue
		}

		// Process content based on current section
		switch currentSection {
		case "summary":
			// Only collect clean summary lines (no numbered points)
			if !strings.HasPrefix(line, "1.") && !strings.HasPrefix(line, "2.") &&
				!strings.HasPrefix(line, "3.") && !strings.HasPrefix(line, "4.") &&
				!strings.HasPrefix(line, "5.") && !strings.Contains(line, "➤") {
				summaryLines = append(summaryLines, line)
			}
		case "keypoints":
			cleanLine := strings.TrimLeft(line, "0123456789-.*• ➤")
			cleanLine = strings.TrimSpace(cleanLine)
			if cleanLine != "" && !strings.HasPrefix(strings.ToLower(cleanLine), "key point") {
				result.KeyPoints = append(result.KeyPoints, cleanLine)
			}
		default:
			// If no section detected yet, treat as summary
			if result.Summary == "" && currentSection == "" {
				summaryLines = append(summaryLines, line)
			}
		}
	}

	// Join summary lines
	if len(summaryLines) > 0 {
		result.Summary = strings.Join(summaryLines, " ")
	}

	// If parsing failed, use the content as summary but try to clean it
	if result.Summary == "" {
		// Remove obvious key points from summary
		cleanedContent := strings.ReplaceAll(content, "\n\n", " ")
		cleanedContent = strings.ReplaceAll(cleanedContent, "\n", " ")
		result.Summary = cleanedContent
	}

	return result
}

// CallLLMAPI makes a request to OpenRouter API using Mistral 3.1 small 24b and returns the response content
func CallLLMAPI(question string) (string, error) {
	if err := godotenv.Load("configs/.env"); err != nil {
		return "", fmt.Errorf("error loading .env file: %w", err)
	}

	openRouterKey := os.Getenv("OPEN_ROUTER_KEY")
	if openRouterKey == "" {
		return "", fmt.Errorf("OPEN_ROUTER_KEY environment variable is not set")
	}

	payload := map[string]interface{}{
		"model": "mistralai/mistral-small",
		"messages": []map[string]string{
			{"role": "system", "content": "You are a research assistant that provides structured, factual information. Format your response with clear sections using exactly these headers: 'Summary:' and 'Key Points:'. Use emojis sparingly and only where they enhance understanding."},
			{"role": "user", "content": question + "\n\nPlease structure your response as follows:\n\nSummary:\n[Provide a concise 2-3 sentence summary without numbered points]\n\nKey Points:\n1. [First key point]\n2. [Second key point]\n3. [Third key point]"},
		},
	}

	jsonBody, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "https://openrouter.ai/api/v1/chat/completions", bytes.NewBuffer(jsonBody))
	req.Header.Set("Authorization", "Bearer "+openRouterKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("HTTP-Referer", "https://github.com/dig-research-tool")
	req.Header.Set("X-Title", "Dig Research Tool")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	var response APIResponse
	err = json.Unmarshal(body, &response)
	if err != nil || len(response.Choices) == 0 {
		return string(body), nil // Return raw response if parsing fails
	}

	return response.Choices[0].Message.Content, nil
}
