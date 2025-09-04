/*
Copyright © 2025 output HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"panteks/internal/validator"
	"panteks/internal/web"
	"time"

	"github.com/spf13/cobra"
)

type saveFunc func([]byte, ...string) (*string, int, error)

var outputs = map[string]saveFunc{
    "html": web.SaveToHTML,
    // "json": web.SaveToJSON,
    // "txt":  web.SaveToTXT,
}


var scrapeCmd = &cobra.Command{
	Use:   "scrape",
	Short: "Scrape data from a specified target",
	Long: `Scrape data from a specified target. Use the --target or -t flag to specify the target.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		target, _ := cmd.Flags().GetString("target")
		selector, _ := cmd.Flags().GetString("selector")
		html, _ := cmd.Flags().GetBool("html")
		output,_ := cmd.Flags().GetString("output")
		jsonFlag, _ := cmd.Flags().GetBool("json")
		txt, _ := cmd.Flags().GetBool("txt")

		if jsonFlag && selector == "" {
            return fmt.Errorf("--selector (-s) is required when using --json")
        }
	
	if output == "" {
   	  switch {
   	 		case html:
    	    	output = "output.html"
    		case jsonFlag:
      			  output = "output.json"
    		case txt:
      			  output = "output.txt"
    		default:
       			 output = "output.txt" 
  	  	}
	}

		if err := validator.ValidateTarget(target); err != nil {
			return err
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		data, err := web.Scrape(ctx, target)
		if err != nil {
			return fmt.Errorf("error occurred: %w", err)
		}

		flags := map[string]bool{
			"html": html,
			"json": jsonFlag,
			"txt":  txt,
		}
		anyFlag := false

		for k, enabled := range flags {
		if enabled {
			path, n, err := outputs[k](data, output)
			if err != nil {
				return fmt.Errorf("failed to save %s: %w", k, err)
			}
			fmt.Printf("Data successfully saved in %s (%d bytes)\n", *path, n)
			anyFlag = true
		}
}

		if !anyFlag {
			path, _, err := web.SaveToHTML(data, output)
			if err != nil {
				return fmt.Errorf("error occurred: %w", err)
			}
			fmt.Printf("Data successfully saved in %s\n", *path)
		}

				return nil
			},
		}

func init() {
	rootCmd.AddCommand(scrapeCmd)

	scrapeCmd.Flags().StringP("target", "t", "", "Target URL or data source to scrape")
	scrapeCmd.Flags().StringP("selector", "s", "", "CSS selector for the data to extract")
	scrapeCmd.Flags().StringP("output", "o", "", "Output location")
	scrapeCmd.Flags().BoolP("html", "H", false, "Generate HTML output")
	scrapeCmd.Flags().BoolP("json", "J", false, "Generate JSON output")
	scrapeCmd.Flags().BoolP("txt", "T", false, "Generate TXT output")

	scrapeCmd.MarkFlagRequired("target")

	// Here you can define additional flags and configuration settings.

	// Cobra supports Persistent Flags, which will work for this command
	// and all subcommands, e.g.:
	// scrapeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// scrapeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
