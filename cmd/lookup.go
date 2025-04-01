package cmd

import (
	"kurocfer/lil-url/shortner"

	"github.com/spf13/cobra"
)

// lookupCmd represents the lookup command
var lookupCmd = &cobra.Command{
	Use:   "lookup <short>",
	Short: "lookup and open a shortened url",
	RunE: func(cmd *cobra.Command, args []string) error {
		storageFile, err := cmd.Flags().GetString("file")
		if err != nil {
			return err
		}

		return lookupRun(storageFile, args[0])
	},
}

func lookupRun(storeageFile, shortURL string) error {
	shortner := shortner.NewShortner(storeageFile)

	return shortner.LookupURL(shortURL)
}

func init() {
	rootCmd.AddCommand(lookupCmd)
}
