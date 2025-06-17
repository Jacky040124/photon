package pkg

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// keyMap defines the keybindings for the model selector
type keyMap struct {
	Up       key.Binding
	Down     key.Binding
	Select   key.Binding
	Quit     key.Binding
	Help     key.Binding
	Details  key.Binding
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	Select: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc"),
		key.WithHelp("q/esc", "quit"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "help"),
	),
	Details: key.NewBinding(
		key.WithKeys(" "),
		key.WithHelp("space", "details"),
	),
}

// ModelSelectorModel represents the state of the model selector TUI
type ModelSelectorModel struct {
	models       map[string]Model
	modelOrder   []string
	currentModel string
	cursor       int
	selected     string
	showDetails  bool
	showHelp     bool
	width        int
	height       int
}

// NewModelSelector creates a new model selector
func NewModelSelector(currentModel string) ModelSelectorModel {
	models := GetAvailableModels()
	modelOrder := []string{"deepseek-r1", "deepseek-v3", "llama-4", "gemini", "mistral"}
	
	// Find current model index for initial cursor position
	cursor := 0
	for i, id := range modelOrder {
		if id == currentModel {
			cursor = i
			break
		}
	}

	return ModelSelectorModel{
		models:       models,
		modelOrder:   modelOrder,
		currentModel: currentModel,
		cursor:       cursor,
		width:        80,
		height:       20,
	}
}

// Init initializes the model selector
func (m ModelSelectorModel) Init() tea.Cmd {
	return nil
}

// Update handles keyboard input and updates the model state
func (m ModelSelectorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Quit):
			return m, tea.Quit

		case key.Matches(msg, keys.Up):
			if m.cursor > 0 {
				m.cursor--
			} else {
				m.cursor = len(m.modelOrder) - 1
			}

		case key.Matches(msg, keys.Down):
			if m.cursor < len(m.modelOrder)-1 {
				m.cursor++
			} else {
				m.cursor = 0
			}

		case key.Matches(msg, keys.Select):
			m.selected = m.modelOrder[m.cursor]
			return m, tea.Quit

		case key.Matches(msg, keys.Details):
			m.showDetails = !m.showDetails

		case key.Matches(msg, keys.Help):
			m.showHelp = !m.showHelp
		}
	}

	return m, nil
}

// View renders the model selector interface
func (m ModelSelectorModel) View() string {
	var b strings.Builder

	// Header
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("86")).
		MarginBottom(1)
	
	b.WriteString(headerStyle.Render("🤖 Select AI Model:"))
	b.WriteString("\n\n")

	// Model list
	for i, modelID := range m.modelOrder {
		model := m.models[modelID]
		
		// Style based on selection and current model
		var style lipgloss.Style
		var prefix string
		var suffix string
		
		if i == m.cursor {
			// Highlighted/selected item
			style = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("15")).
				Background(lipgloss.Color("63")).
				Padding(0, 1)
			prefix = ">"
		} else {
			// Regular item
			style = lipgloss.NewStyle().
				Foreground(lipgloss.Color("15")).
				Padding(0, 1)
			prefix = " "
		}
		
		// Mark current model
		if modelID == m.currentModel {
			suffix = GreenBold(" (current)")
		}
		
		// Model name line
		modelLine := fmt.Sprintf("%s %s%s", prefix, model.Name, suffix)
		b.WriteString(style.Render(modelLine))
		b.WriteString("\n")
		
		// Description line (always shown, but styled differently for selected)
		var descStyle lipgloss.Style
		if i == m.cursor {
			descStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("245")).
				Italic(true).
				MarginLeft(2)
		} else {
			descStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("240")).
				MarginLeft(2)
		}
		
		b.WriteString(descStyle.Render(model.Description))
		b.WriteString("\n")
		
		// Show additional details for selected model if requested
		if i == m.cursor && m.showDetails {
			detailStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color("117")).
				MarginLeft(2).
				Italic(true)
			
			details := fmt.Sprintf("Provider: %s | Context: %d tokens | Features: %s",
				model.Provider,
				model.ContextLen,
				strings.Join(model.Features, ", "))
			
			b.WriteString(detailStyle.Render(details))
			b.WriteString("\n")
		}
		
		b.WriteString("\n")
	}

	// Help section
	if m.showHelp {
		helpStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("244")).
			Border(lipgloss.RoundedBorder()).
			Padding(1).
			MarginTop(1)
		
		helpText := "Controls:\n" +
			"↑/k: Move up    ↓/j: Move down\n" +
			"Enter: Select   Space: Toggle details\n" +
			"q/Esc: Quit     ?: Toggle help"
		
		b.WriteString(helpStyle.Render(helpText))
		b.WriteString("\n")
	}

	// Footer
	footerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("244")).
		MarginTop(1)
	
	footerText := "Press Enter to select, q to quit"
	if !m.showHelp {
		footerText += ", ? for help"
	}
	
	b.WriteString(footerStyle.Render(footerText))

	return b.String()
}

// GetSelectedModel returns the selected model ID, or empty string if cancelled
func (m ModelSelectorModel) GetSelectedModel() string {
	return m.selected
}

// RunModelSelector runs the interactive model selector and returns the selected model ID
func RunModelSelector(currentModel string) (string, error) {
	m := NewModelSelector(currentModel)
	program := tea.NewProgram(m)
	
	finalModel, err := program.Run()
	if err != nil {
		return "", err
	}
	
	if selectorModel, ok := finalModel.(ModelSelectorModel); ok {
		return selectorModel.GetSelectedModel(), nil
	}
	
	return "", nil
}