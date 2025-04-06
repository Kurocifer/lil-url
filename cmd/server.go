/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"kurocfer/lil-url/server"
	"strings"

	"github.com/spf13/cobra"
)

var port string

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start a server that listens to reqeust from shortened urls and redirects to the orginal url",
	RunE: func(cmd *cobra.Command, args []string) error {
		storageFile, err := cmd.Flags().GetString("file")
		if err != nil {
			return err
		}

		return serverCmdRun(storageFile, port)
	},
}

func serverCmdRun(storageFile, port string) error {
	if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}

	return server.Start(storageFile, port)
}

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.Flags().StringVarP(&port, "port", "p", ":8080", "server port number")
}
