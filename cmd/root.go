/*
Copyright © 2025 Renzi Febriandika <Renzifebriandika923@gmail.com>
*/
package cmd

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)


var version = "1.0.0"

var rootCmd = &cobra.Command{
    Use:   "panteks",
    Short: "Panteks is a CLI tool for web scraping",
    Long: `Panteks is a command-line interface (CLI) tool designed to scrape data from websites efficiently.
It extracts the important information from web pages and converts it into plain text, making it easy
to process, analyze, or store. Panteks supports various scraping options and output formats,
enabling developers and data analysts to automate the collection of textual content from websites.

Features include:
  - Extract plain text content from HTML pages
  - Save results in text or HTML formats
  - Support for target URLs`,
    Version: version,
}



func Execute() {
 	 if err := godotenv.Load(); err != nil {
        log.Println("Warning: .env file not found, using system env")
    }

    if err := rootCmd.Execute(); err != nil {
        log.Fatal(err)
    }
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}


