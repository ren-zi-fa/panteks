/*
Copyright © 2025 <renzifebriandika923@gmail.com>
*/
package cmd

import (
	"context"
	"fmt"
	"os"
	"panteks/internal/utils"
	"panteks/internal/validator"

	"github.com/spf13/cobra"
)

var (
    ApiKey = os.Getenv("API_KEY")
    ApiURL = os.Getenv("API_URL")
)


var scrapeCmd = &cobra.Command{
	Use:   "scrape",
	Short: "Scrape data from a specified target",
	Long:  `Scrape data from a specified target. Use the --target or -t flag to specify the target.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		target, _ := cmd.Flags().GetString("target")
		output, _ := cmd.Flags().GetString("output")
		html, _ := cmd.Flags().GetBool("html")

		if err := validator.ValidateTarget(target); err != nil {
			return err
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		switch {
		case html:
			resultScrape, err := utils.Scrape(ctx, target)
			if err != nil {
				return err
			}
			resultHTML, _, err := utils.SaveToHTML(resultScrape, output)
			if err != nil {
				return err
			}
			fmt.Printf("The result is available at %s\n", *resultHTML)

		default:
			resultScrape, err := utils.Scrape(ctx, target)
			if err != nil {
				return err
			}

			chunks := utils.SplitContent(string(resultScrape), 3000)

			var combinedResult string
			for i, chunk := range chunks {
				fmt.Printf("➡️ Processing chunk %d/%d...\n", i+1, len(chunks))
				content := "Extract the important information from the following HTML and return ONLY plain text:\n\n" + chunk

				result, err := utils.CallAPIWithRetry(ApiKey, ApiURL, content)
				if err != nil {
					panic(err)
				}
				combinedResult += result + "\n"
			}

			resultTXT, _, err := utils.SaveToTXT([]byte(combinedResult), output)
			if err != nil {
				return err
			}
			fmt.Printf("The result is available at %s\n", *resultTXT)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(scrapeCmd)

	scrapeCmd.Flags().StringP("target", "t", "", "Target URL or data source to scrape")
	scrapeCmd.Flags().StringP("output", "o", "", "Output location")
	scrapeCmd.Flags().BoolP("html", "H", false, "Generate HTML output")

	scrapeCmd.MarkFlagRequired("target")

	// Here you can define additional flags and configuration settings.

	// Cobra supports Persistent Flags, which will work for this command
	// and all subcommands, e.g.:
	// scrapeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// scrapeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
