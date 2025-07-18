package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"github.com/Jacky040124/photon/pkg"
)

var onlineModel string

var rootCmd = &cobra.Command{
	Use:   "ptn [query]",
	Short: "Packets of pure knowledge at light speed",
	Long:  "Photon is a lightning-fast terminal research tool that delivers packets of pure knowledge at light speed.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Load and validate config
		config, err := LoadConfig()
		if err != nil {
			fmt.Println(pkg.RedBold("Error loading config: ") + err.Error())
			os.Exit(1)
		}

		err = config.Validate()
		if err != nil {
			fmt.Println(pkg.RedBold("Configuration error: ") + err.Error())
			fmt.Println("Please set PHOTON_OPEN_ROUTER_KEY environment variable")
			fmt.Println("Example: export PHOTON_OPEN_ROUTER_KEY=\"your-api-key\"")
			os.Exit(1)
		}

		question := args[0]
		// If --online is set, override the model for this run
		if onlineModel != "" {
			config.CurrentModel = "__online__" + onlineModel
		}
		// Initialize the TUI with the chosen model ID
		m := initialModel(question, config.GetCurrentModel())

		_, err = tea.NewProgram(m).Run()
		if err != nil {
			fmt.Println(pkg.RedBold("could not run program: ") + err.Error())
			os.Exit(1)
		}
	},
}

func Execute() {
	// Add model subcommand
	rootCmd.AddCommand(modelCmd)
	// Add --online flag
	rootCmd.PersistentFlags().StringVar(&onlineModel, "online", "", "Use any OpenRouter model by API name (bypasses local list)")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(pkg.RedBold("Error: ") + err.Error())
		os.Exit(1)
	}
}
