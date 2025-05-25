package main

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/joho/godotenv"

	"github.com/jacky/dig/pkg"
)

// Loading state constants
const (
	stateLoading = iota
	stateResult
)

type state int

type model struct {
	spinner      spinner.Model
	loadingState state
	question     string
	fallback     bool
	result       pkg.FormattedResearch
}

type llmResultMsg struct {
	Research pkg.FormattedResearch
}

func initialModel(question string) model {
	return model{
		spinner:      pkg.CreateSpinner(),
		loadingState: stateLoading,
		question:     question,
		fallback:     false,
	}
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
		research := pkg.FormatResearch(question)
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

func (m model) View() string {
	switch m.loadingState {
	case stateResult:
		return pkg.RenderResultView(m.result)
	default:
		uiModel := pkg.UIModel{
			Spinner:  m.spinner,
			Fallback: m.fallback,
			Result:   m.result,
		}
		return pkg.RenderLoadingView(uiModel)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println(pkg.RedBold("Please provide a query"))
		return
	}

	// Load config from configs directory
	godotenv.Load("../configs/.env")

	question := os.Args[1]
	m := initialModel(question)

	_, err := tea.NewProgram(m).Run()
	if err != nil {
		fmt.Println(pkg.RedBold("could not run program: ") + err.Error())
		os.Exit(1)
	}
}
