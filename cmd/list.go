package cmd

import (
	"kurocfer/lil-url/shortner"

	"github.com/spf13/cobra"
)

var numLines int

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Show a list of shortened links with their original links",
	Long: `list, Shows a list of the links that have been shortened,
showing the shortened links and their original links. By default list a maximum
of 10 links this can be modified with the -n flag.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		storageFile, err := cmd.Flags().GetString("file")
		if err != nil {
			return err
		}

		return listCmdRun(storageFile, numLines)
	},
}

func listCmdRun(storageFile string, numLines int) error {
	shortner := shortner.NewShortner(storageFile)
	return shortner.List(numLines)
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().IntVarP(&numLines, "num", "n", 10, "Number of links listed")
}
