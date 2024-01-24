package cmd

import (
	"github.com/belarte/metrix/database"
	"github.com/belarte/metrix/server"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the server",
	RunE: func(cmd *cobra.Command, args []string) error {
		db := database.NewInMemory()
		return server.New(db).Start(":8080")
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
