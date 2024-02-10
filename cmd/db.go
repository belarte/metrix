package cmd

import (
	"fmt"

	"github.com/belarte/metrix/repository"
	"github.com/spf13/cobra"
)

const file = "dev.sqlite"

var dbCmd = &cobra.Command{
	Use:   "db",
	Short: "Database operations",
	RunE: func(cmd *cobra.Command, args []string) error {
		db, err := repository.New(file)
		if err != nil {
			return fmt.Errorf("error opening database: %w", err)
		}

		metrics, err := db.GetMetrics()
		if err != nil {
			return fmt.Errorf("error getting metrics: %w", err)
		}

		for _, m := range metrics {
			fmt.Printf("%s (%s): %s\n", m.Title, m.Unit, m.Description)
		}

		return err
	},
}

func init() {
	rootCmd.AddCommand(dbCmd)
}
