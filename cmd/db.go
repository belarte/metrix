package cmd

import (
	"fmt"
	"log"

	"github.com/belarte/metrix/repository"
	"github.com/spf13/cobra"
)

var dbCmd = &cobra.Command{
	Use:   "db",
	Short: "Database operations",
}

var dbMigrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate the database",
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Println("Migrating database", database)

		db, err := repository.New(database)
		if err != nil {
			return fmt.Errorf("error opening database: %w", err)
		}

		err = db.Migrate()
		if err != nil {
			return fmt.Errorf("error migrating database: %w", err)
		}

		return err
	},
}

func init() {
	dbCmd.AddCommand(dbMigrateCmd)
	rootCmd.AddCommand(dbCmd)
}
