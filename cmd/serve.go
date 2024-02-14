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
		addr := cmd.Flag("address").Value.String()
		dbPath := cmd.Flag("database").Value.String()

		db, err := repository.New(dbPath)
		if err != nil {
			return fmt.Errorf("could not create database: %w", err)
		}
		defer db.Close()

		if err = db.Migrate(); err != nil {
			return fmt.Errorf("could not migrate database: %w", err)
		}

		return server.New(server.WithRepository(db), server.WithAddress(addr)).Start()
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	serverCmd.Flags().StringP("address", "a", ":8080", "address to listen to")
	serverCmd.Flags().StringP("database", "d", ":memory:", "path to the database")
}
