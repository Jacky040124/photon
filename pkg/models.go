package pkg

import (
	"fmt"
	"strings"
)

// Model represents an AI model with its configuration and metadata
type Model struct {
	ID           string
	Name         string
	APIName      string
	Description  string
	Provider     string
	Features     []string
	ContextLen   int
	BestFor      string
	IsThinking   bool
	IsMultimodal bool
}

// GetAvailableModels returns the list of all available free models
func GetAvailableModels() map[string]Model {
	return map[string]Model{
		"deepseek-r1": {
			ID:           "deepseek-r1",
			Name:         "DeepSeek R1",
			APIName:      "deepseek/deepseek-r1:free",
			Description:  "Advanced reasoning model with step-by-step thinking capabilities",
			Provider:     "DeepSeek",
			Features:     []string{"Reasoning", "Problem Solving", "Analysis"},
			ContextLen:   163840,
			BestFor:      "Complex analysis, problem-solving, research tasks",
			IsThinking:   true,
			IsMultimodal: false,
		},
		"deepseek-v3": {
			ID:           "deepseek-v3",
			Name:         "DeepSeek V3 Chat",
			APIName:      "deepseek/deepseek-chat:free",
			Description:  "General purpose model with excellent coding and instruction following",
			Provider:     "DeepSeek",
			Features:     []string{"General Purpose", "Coding", "Instruction Following"},
			ContextLen:   163840,
			BestFor:      "General queries, coding help, conversational tasks",
			IsThinking:   false,
			IsMultimodal: false,
		},
		"llama-4": {
			ID:           "llama-4",
			Name:         "Meta Llama 4 Maverick",
			APIName:      "meta-llama/llama-4-maverick:free",
			Description:  "Multimodal model supporting text and image analysis",
			Provider:     "Meta",
			Features:     []string{"Multimodal", "Text", "Image Analysis"},
			ContextLen:   128000,
			BestFor:      "Image analysis, visual content research",
			IsThinking:   false,
			IsMultimodal: true,
		},
		"kimi": {
			ID:           "kimi",
			Name:         "MoonshotAI Kimi K2",
			APIName:      "moonshotai/kimi-k2:free",
			Description:  "Advanced Chinese AI model with strong reasoning capabilities",
			Provider:     "MoonshotAI",
			Features:     []string{"Strong Reasoning", "Chinese & English", "Code Generation"},
			ContextLen:   200000,
			BestFor:      "Bilingual research, code analysis, logical reasoning",
			IsThinking:   false,
			IsMultimodal: false,
		},
		"mistral": {
			ID:           "mistral",
			Name:         "Mistral Small 3.1",
			APIName:      "mistralai/mistral-small-3.1-24b-instruct:free",
			Description:  "Efficient and fast model with good balance of speed and capability",
			Provider:     "Mistral AI",
			Features:     []string{"Fast", "Efficient", "Balanced"},
			ContextLen:   128000,
			BestFor:      "Quick responses, general research",
			IsThinking:   false,
			IsMultimodal: true,
		},
	}
}

// GetDefaultModel returns the default model ID
func GetDefaultModel() string {
	return "deepseek-v3"
}

// GetModel returns a model by ID, or nil if not found
func GetModel(id string) (*Model, error) {
	models := GetAvailableModels()
	if model, exists := models[id]; exists {
		return &model, nil
	}
	return nil, fmt.Errorf("model '%s' not found", id)
}

// GetModelByAPIName returns a model by its API name
func GetModelByAPIName(apiName string) (*Model, error) {
	models := GetAvailableModels()
	for _, model := range models {
		if model.APIName == apiName {
			return &model, nil
		}
	}
	return nil, fmt.Errorf("model with API name '%s' not found", apiName)
}

// ValidateModel checks if a model ID is valid
func ValidateModel(id string) bool {
	models := GetAvailableModels()
	_, exists := models[id]
	return exists
}

// FormatModelInfo returns a formatted string with model information
func FormatModelInfo(model Model) string {
	var b strings.Builder
	
	b.WriteString(fmt.Sprintf("%s %s\n", CyanBold("📋 Model:"), YellowBold(model.Name)))
	b.WriteString(fmt.Sprintf("%s %s\n", BlueBold("🏢 Provider:"), White(model.Provider)))
	b.WriteString(fmt.Sprintf("%s %s\n", GreenBold("📝 Description:"), White(model.Description)))
	b.WriteString(fmt.Sprintf("%s %s\n", Magenta("🎯 Best For:"), White(model.BestFor)))
	b.WriteString(fmt.Sprintf("%s %s\n", Cyan("🔧 Features:"), White(strings.Join(model.Features, ", "))))
	b.WriteString(fmt.Sprintf("%s %d tokens\n", Blue("📏 Context:"), model.ContextLen))
	
	if model.IsThinking {
		b.WriteString(fmt.Sprintf("%s %s\n", GreenBold("🧠 Special:"), White("Supports reasoning with <think> tokens")))
	}
	if model.IsMultimodal {
		b.WriteString(fmt.Sprintf("%s %s\n", YellowBold("🖼️  Multimodal:"), White("Supports text and images")))
	}
	
	return b.String()
}

// FormatModelList returns a formatted list of all available models
func FormatModelList(currentModel string) string {
	var b strings.Builder
	models := GetAvailableModels()
	
	b.WriteString(CyanBold("✨ Available Models:\n\n"))
	
	for _, id := range []string{"kimi", "deepseek-r1", "deepseek-v3", "llama-4", "mistral"} {
		model := models[id]
		current := ""
		if id == currentModel {
			current = GreenBold(" (current)")
		}
		
		b.WriteString(fmt.Sprintf("%s %s%s\n", 
			YellowBold(fmt.Sprintf("%-12s", id)), 
			White(model.Name), 
			current))
		b.WriteString(fmt.Sprintf("             %s\n", Cyan(model.Description)))
		b.WriteString("\n")
	}
	
	return b.String()
}