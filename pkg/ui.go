package pkg

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"
	"github.com/fatih/color"
)

// Color functions for terminal output
var (
	CyanBold   = color.New(color.FgCyan, color.Bold).SprintFunc()
	GreenBold  = color.New(color.FgGreen, color.Bold).SprintFunc()
	YellowBold = color.New(color.FgYellow, color.Bold).SprintFunc()
	BlueBold   = color.New(color.FgBlue, color.Bold).SprintFunc()
	RedBold    = color.New(color.FgRed, color.Bold).SprintFunc()
	White      = color.New(color.FgWhite).SprintFunc()
	Cyan       = color.New(color.FgCyan).SprintFunc()
	Blue       = color.New(color.FgBlue).SprintFunc()
	Green      = color.New(color.FgGreen).SprintFunc()
	Magenta    = color.New(color.FgMagenta).SprintFunc()
)

// UI styling constants
var SpinnerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("69"))

// UIModel represents the UI state for rendering
type UIModel struct {
	Spinner  spinner.Model
	Fallback bool
	Result   FormattedResponse
}

// CreateSpinner creates and configures a new spinner
func CreateSpinner() spinner.Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = SpinnerStyle
	return s
}

// RenderLoadingView renders the loading state with spinner
func RenderLoadingView(uiModel UIModel) string {
	if uiModel.Fallback {
		return RedBold("\nLost in the tunnel of knowledge. Please try again later.\n")
	}
	return fmt.Sprintf("\n %s %s\n\n", uiModel.Spinner.View(), CyanBold("THINKING.."))
}

// RenderResultView renders the formatted research results
func RenderResultView(result FormattedResponse) string {
	b := strings.Builder{}
	b.WriteString("\n" + CyanBold("✨ === PHOTON RESEARCH RESULTS === ✨") + "\n")
	b.WriteString("\n" + YellowBold("✨ SUMMARY:") + "\n")
	b.WriteString(White(result.Summary) + "\n")

	if len(result.KeyPoints) > 0 {
		b.WriteString("\n" + GreenBold("💡 KEY POINTS:") + "\n")
		for i, point := range result.KeyPoints {
			b.WriteString(fmt.Sprintf("%s %d. %s\n", Cyan("➤"), i+1, White(point)))
		}
	}

	b.WriteString("\n" + CyanBold("✨ ========================== ✨") + "\n")
	return b.String()
}

// PrintFormattedResearch prints research results directly to console
func PrintFormattedResearch(research FormattedResponse) {
	fmt.Println(RenderResultView(research))
}
