package cmd

import (
	"fmt"
	"log"

	"github.com/belarte/metrix/repository"
	"github.com/belarte/metrix/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var address string

const defaultAddress = ":8080"

var serverCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the server",
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Println("Using database", database)

		db, err := repository.New(database)
		if err != nil {
			return fmt.Errorf("could not create database: %w", err)
		}
		defer db.Close()

		if database == inMemoryDatabase {
			if err = db.Migrate(); err != nil {
				return fmt.Errorf("could not migrate database: %w", err)
			}
		}

		return server.New(server.WithRepository(db), server.WithAddress(address)).Start()
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	viper.BindEnv("METRIX_ADDRESS")
	viper.SetDefault("METRIX_ADDRESS", defaultAddress)
	serverCmd.Flags().StringVarP(&address, "address", "a", viper.GetString("METRIX_ADDRESS"), "address to listen to")
}
