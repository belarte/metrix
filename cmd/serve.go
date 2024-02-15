package cmd

import (
	"fmt"

	"github.com/belarte/metrix/repository"
	"github.com/belarte/metrix/server"
	"github.com/spf13/cobra"
)

var address string

const defaultAddress = ":8080"

var serverCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the server",
	RunE: func(cmd *cobra.Command, args []string) error {
		db, err := repository.New(database)
		if err != nil {
			return fmt.Errorf("could not create database: %w", err)
		}
		defer db.Close()

		if database == inMemoryDatabase {
			fmt.Println("Using in-memory database")
			if err = db.Migrate(); err != nil {
				return fmt.Errorf("could not migrate database: %w", err)
			}
		}

		return server.New(server.WithRepository(db), server.WithAddress(address)).Start()
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	serverCmd.Flags().StringVarP(&address, "address", "a", defaultAddress, "address to listen to")
}
