package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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

var spinnerLine = spinner.Line
var spinnerDot = spinner.Dot
var spinnerMiniDot = spinner.MiniDot
var spinnerJump = spinner.Jump
var spinnerPulse = spinner.Pulse
var spinnerPoints = spinner.Points
var spinnerGlobe = spinner.Globe
var spinnerMoon = spinner.Moon
var spinnerMonkey = spinner.Monkey

var spinners = []spinner.Spinner{
	spinnerLine,
	spinnerDot,
	spinnerMiniDot,
	spinnerJump,
	spinnerPulse,
	spinnerPoints,
	spinnerGlobe,
	spinnerMoon,
	spinnerMonkey,
}

var textStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("252")).Render
var spinnerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("69"))

// Loading state
const (
	stateLoading = iota
	stateResult
)

type state int

type model struct {
	spinner      spinner.Model
	loadingState state
	question     string
	spinnerIndex int
	fallback     bool
	result       FormattedResearch
}

// LLM result message
type llmResultMsg struct {
	Research FormattedResearch
}

func initialModel(question string) model {
	m := model{
		spinner:      spinner.New(),
		loadingState: stateLoading,
		question:     question,
		spinnerIndex: 1,
		fallback:     false,
	}
	m.spinner.Spinner = spinners[m.spinnerIndex]
	m.spinner.Style = spinnerStyle
	return m
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		timeoutCmd(),
		getLLMResearchCmd(m.question),
	)
}

func timeoutCmd() tea.Cmd {
	return func() tea.Msg {
		time.Sleep(15 * time.Second)
		return fallbackMsg{}
	}
}

type fallbackMsg struct{}

func getLLMResearchCmd(question string) tea.Cmd {
	return func() tea.Msg {
		research := FormatResearch(question)
		return llmResultMsg{Research: research}
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case fallbackMsg:
		if m.loadingState == stateLoading {
			m.fallback = true
			return m, tea.Quit
		}
		return m, nil
	case llmResultMsg:
		if m.loadingState == stateLoading && !m.fallback {
			m.result = msg.Research
			m.loadingState = stateResult
			return m, tea.Quit
		}
		return m, nil
	}
	return m, nil
}

func loadingView(m model) string {
	gap := " "
	if m.spinnerIndex == 1 {
		gap = ""
	}
	if m.fallback {
		return textStyle("\nLost in the tunnel of knowledge. Please try again later.\n")
	}
	return fmt.Sprintf("\n %s%s%s\n\n", m.spinner.View(), gap, textStyle("Digging with LLM..."))
}

func resultView(m model) string {
	b := strings.Builder{}
	b.WriteString("\n✨ === DIG RESEARCH RESULTS === ✨\n")
	b.WriteString("\n✨ SUMMARY:\n")
	b.WriteString(textStyle(m.result.Summary) + "\n")
	if len(m.result.KeyPoints) > 0 {
		b.WriteString("\n💡 KEY POINTS:\n")
		for i, point := range m.result.KeyPoints {
			b.WriteString(fmt.Sprintf("➤ %d. %s\n", i+1, point))
		}
	}
	if len(m.result.SourceLinks) > 0 {
		b.WriteString("\n🔗 SOURCES:\n")
		for _, source := range m.result.SourceLinks {
			b.WriteString("➤ " + source + "\n")
		}
	}
	b.WriteString("\n✨ ========================== ✨\n")
	return b.String()
}

func (m model) View() string {
	switch m.loadingState {
	case stateResult:
		return resultView(m)
	default:
		return loadingView(m)
	}
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

func OpenAI(question string) (string, error) {
	godotenv.Load()
	var openaiPayload map[string]interface{}

	openaiKey := os.Getenv("OPENAI_API_KEY")
	openaiPayload = map[string]interface{}{
		"model": "gpt-3.5-turbo",
		"messages": []map[string]string{
			{"role": "system", "content": "You are a research assistant that provides concise, factual information with source links. Use subtle and aesthetically pleasing emojis where appropriate to enhance readability and engagement. Format the output in terminal friendly format, avoiding markdown or excessive indentation."},
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
	if len(os.Args) < 2 {
		fmt.Println("Please provide a querry")
		return
	}
	godotenv.Load()
	question := os.Args[1]
	m := initialModel(question)

	_, err := tea.NewProgram(m).Run()
	if err != nil {
		fmt.Println("could not run program:", err) 
		os.Exit(1)
	}
}
