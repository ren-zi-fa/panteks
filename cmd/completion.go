/*
Copyright © 2025 <renzifebriandika923@gmail.com>
*/

package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "Generate shell completion script for your CLI",
	Long: `Generate a shell completion script for Panteks CLI.
Supports bash, zsh, fish, and PowerShell.

Examples:

  # Bash
  panteks completion bash > ~/.panteks_completion.sh
  source ~/.panteks_completion.sh

  # Zsh
  panteks completion zsh > ~/.panteks_completion.zsh
  source ~/.panteks_completion.zsh

  # Fish
  panteks completion fish > ~/.config/fish/completions/panteks.fish

  # PowerShell
  panteks completion powershell > panteks_completion.ps1
  . ./panteks_completion.ps1`,
	Args: cobra.ExactArgs(1),
	ValidArgs: []string{"bash", "zsh", "fish", "powershell"},
	Run: func(cmd *cobra.Command, args []string) {
		shell := args[0]

		var err error
		switch shell {
		case "bash":
			err = rootCmd.GenBashCompletionFile("panteks_completion.sh")
			fmt.Println("Bash completion script generated: panteks_completion.sh")
		case "zsh":
			err = rootCmd.GenZshCompletionFile("panteks_completion.zsh")
			fmt.Println("Zsh completion script generated: panteks_completion.zsh")
		case "fish":
			err = rootCmd.GenFishCompletionFile("panteks_completion.fish", true)
			fmt.Println("Fish completion script generated: panteks_completion.fish")
		case "powershell":
			err = rootCmd.GenPowerShellCompletionFile("panteks_completion.ps1")
			fmt.Println("PowerShell completion script generated: panteks_completion.ps1")
		default:
			log.Fatalf("Unsupported shell: %s", shell)
		}

		if err != nil {
			log.Fatalf("Failed to generate completion script: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)
}
