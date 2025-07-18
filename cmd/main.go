package main

import (
	"strings"
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

// model stores TUI state, including the effective model ID
type model struct {
	spinner      spinner.Model
	loadingState state
	question     string
	modelID      string
	fallback     bool
	result       pkg.FormattedResponse
}

// initialModel creates a new TUI model with the question and chosen modelID
func initialModel(question, modelID string) model {
	return model{
		spinner:      pkg.CreateSpinner(),
		loadingState: stateLoading,
		question:     question,
		modelID:      modelID,
		fallback:     false,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		timeoutCmd(),
		getLLMResearchCmd(m.question, m.modelID),
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
		// Derive display name from m.modelID
		var displayName string
		if mdl, err := pkg.GetModel(m.modelID); err == nil {
			displayName = mdl.Name
		} else if strings.HasPrefix(m.modelID, "__online__") {
			displayName = m.modelID[len("__online__"):]
		} else {
			displayName = m.modelID
		}
		return pkg.RenderResultView(m.result, displayName)
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

// getLLMResearchCmd returns a command that fetches research using the specified modelID
func getLLMResearchCmd(question, modelID string) tea.Cmd {
	return func() tea.Msg {
		research := pkg.FormatWithModel(question, modelID)
		return llmResultMsg{Research: research}
	}
}

func main() {
	Execute()
}
