package main

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/Jacky040124/photon/pkg"
)

const (
	stateLoading = iota
	stateResult
)

type state int

type fallbackMsg struct{}

type llmResultMsg struct {
	Research pkg.FormattedResponse
}

type model struct {
	spinner      spinner.Model
	loadingState state
	question     string
	fallback     bool
	result       pkg.FormattedResponse
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

func timeoutCmd() tea.Cmd {
	return func() tea.Msg {
		time.Sleep(15 * time.Second)
		return fallbackMsg{}
	}
}

func getLLMResearchCmd(question string) tea.Cmd {
	return func() tea.Msg {
		// Load config to get current model
		config, err := LoadConfig()
		if err != nil {
			return llmResultMsg{Research: pkg.FormattedResponse{
				Summary: fmt.Sprintf("Error loading config: %s", err.Error()),
			}}
		}
		
		// Use the configured model
		research := pkg.FormatWithModel(question, config.GetCurrentModel())
		return llmResultMsg{Research: research}
	}
}

func main() {
	Execute()
}
