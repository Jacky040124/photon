package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/Jacky040124/photon/pkg"
)

var modelCmd = &cobra.Command{
	Use:   "model",
	Short: "Manage AI models",
	Long:  "Manage and configure AI models for research queries",
}

var modelListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available models",
	Long:  "Display all available AI models with their descriptions",
	Run: func(cmd *cobra.Command, args []string) {
		config, err := LoadConfig()
		if err != nil {
			fmt.Println(pkg.RedBold("Error loading config: ") + err.Error())
			os.Exit(1)
		}

		fmt.Print(pkg.FormatModelList(config.GetCurrentModel()))
	},
}

var modelCurrentCmd = &cobra.Command{
	Use:   "current",
	Short: "Show current model",
	Long:  "Display the currently selected AI model",
	Run: func(cmd *cobra.Command, args []string) {
		config, err := LoadConfig()
		if err != nil {
			fmt.Println(pkg.RedBold("Error loading config: ") + err.Error())
			os.Exit(1)
		}

		currentModelID := config.GetCurrentModel()
		model, err := pkg.GetModel(currentModelID)
		if err != nil {
			fmt.Println(pkg.RedBold("Error: ") + err.Error())
			os.Exit(1)
		}

		fmt.Println(pkg.CyanBold("🤖 Current Model:"))
		fmt.Println()
		fmt.Print(pkg.FormatModelInfo(*model))
	},
}

var modelSetCmd = &cobra.Command{
	Use:   "set [model-id]",
	Short: "Set the current model",
	Long:  "Set the AI model to use for research queries. If no model-id is provided, an interactive selection menu will be shown.",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		config, err := LoadConfig()
		if err != nil {
			fmt.Println(pkg.RedBold("Error loading config: ") + err.Error())
			os.Exit(1)
		}

		var modelID string

		if len(args) == 0 {
			// Interactive mode
			modelID = selectModelInteractively(config.GetCurrentModel())
		} else {
			// Direct mode
			modelID = args[0]
		}

		if modelID == "" {
			fmt.Println(pkg.YellowBold("No model selected"))
			return
		}

		// Validate model
		model, err := pkg.GetModel(modelID)
		if err != nil {
			fmt.Println(pkg.RedBold("Error: ") + err.Error())
			os.Exit(1)
		}

		// Set the model
		err = config.SetCurrentModel(modelID)
		if err != nil {
			fmt.Println(pkg.RedBold("Error setting model: ") + err.Error())
			os.Exit(1)
		}

		fmt.Println(pkg.GreenBold("✅ Model set to: ") + pkg.YellowBold(model.Name))
		fmt.Println(pkg.Cyan("Next queries will use this model"))
	},
}

var modelInfoCmd = &cobra.Command{
	Use:   "info <model-id>",
	Short: "Show model details",
	Long:  "Display detailed information about a specific AI model",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		modelID := args[0]
		model, err := pkg.GetModel(modelID)
		if err != nil {
			fmt.Println(pkg.RedBold("Error: ") + err.Error())
			os.Exit(1)
		}

		fmt.Print(pkg.FormatModelInfo(*model))
	},
}

var modelResetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset to default model",
	Long:  "Reset the current model selection to the default model",
	Run: func(cmd *cobra.Command, args []string) {
		config, err := LoadConfig()
		if err != nil {
			fmt.Println(pkg.RedBold("Error loading config: ") + err.Error())
			os.Exit(1)
		}

		defaultModel := pkg.GetDefaultModel()
		err = config.SetCurrentModel(defaultModel)
		if err != nil {
			fmt.Println(pkg.RedBold("Error resetting model: ") + err.Error())
			os.Exit(1)
		}

		model, _ := pkg.GetModel(defaultModel)
		fmt.Println(pkg.GreenBold("✅ Model reset to default: ") + pkg.YellowBold(model.Name))
	},
}

// selectModelInteractively shows an interactive toggle-based model selection
func selectModelInteractively(currentModel string) string {
	selectedModel, err := pkg.RunModelSelector(currentModel)
	if err != nil {
		fmt.Println(pkg.RedBold("Error running model selector: ") + err.Error())
		return ""
	}
	return selectedModel
}

func init() {
	// Add subcommands
	modelCmd.AddCommand(modelListCmd)
	modelCmd.AddCommand(modelCurrentCmd)
	modelCmd.AddCommand(modelSetCmd)
	modelCmd.AddCommand(modelInfoCmd)
	modelCmd.AddCommand(modelResetCmd)
}
