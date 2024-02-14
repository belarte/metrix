package cmd

import (
	"fmt"

	"github.com/belarte/metrix/repository"
	"github.com/belarte/metrix/server"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the server",
	RunE: func(cmd *cobra.Command, args []string) error {
		db, err := repository.New(":memory:")
		if err != nil {
			return fmt.Errorf("could not create database: %w", err)
		}
		defer db.Close()

		if err = db.Migrate(); err != nil {
			return fmt.Errorf("could not migrate database: %w", err)
		}

		return server.New(db).Start(":8080")
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
