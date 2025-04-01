/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"kurocfer/lil-url/shortner"

	"github.com/spf13/cobra"
)

// shortenCmd represents the shorten command
var shortenCmd = &cobra.Command{
	Use:   "shorten <url>",
	Short: "shorten url",
	RunE: func(cmd *cobra.Command, args []string) error {
		file, err := cmd.Flags().GetString("file")
		if err != nil {
			return err
		}
		return shortenRun(file, args[0])
	},
}

func shortenRun(storageFile, longURL string) error {
	shortner := shortner.NewShortner(storageFile)

	shortURL, err := shortner.Shorten(longURL)
	if err != nil {
		return err
	}

	fmt.Printf("Shortened URL: %s\n", shortURL)
	return nil
}

func init() {
	rootCmd.AddCommand(shortenCmd)
}
