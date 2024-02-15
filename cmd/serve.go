package cmd

import (
	"fmt"

	"github.com/belarte/metrix/repository"
	"github.com/belarte/metrix/server"
	"github.com/spf13/cobra"
)

const (
	defaultAddress   = ":8080"
	inMemoryDatabase = ":memory:"
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

		if dbPath == inMemoryDatabase {
			fmt.Println("Using in-memory database")
			if err = db.Migrate(); err != nil {
				return fmt.Errorf("could not migrate database: %w", err)
			}
		}

		return server.New(server.WithRepository(db), server.WithAddress(addr)).Start()
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	serverCmd.Flags().StringP("address", "a", defaultAddress, "address to listen to")
	serverCmd.Flags().StringP("database", "d", inMemoryDatabase, "path to the database")
}
