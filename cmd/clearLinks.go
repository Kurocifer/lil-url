/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"kurocfer/lil-url/shortner"

	"github.com/spf13/cobra"
)

// clearLinksCmd represents the clearLinks command
var clearLinksCmd = &cobra.Command{
	Use:   "clearLinks",
	Short: "Clear all shortned links",
	RunE: func(cmd *cobra.Command, args []string) error {
		storageFile, err := cmd.Flags().GetString("file")
		if err != nil {
			return err
		}

		return RunClearLinks(storageFile)
	},
}

func RunClearLinks(storageFile string) error {
	shortner := shortner.NewShortner(storageFile)
	return shortner.Clear()
}

func init() {
	rootCmd.AddCommand(clearLinksCmd)
}
