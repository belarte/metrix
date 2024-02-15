package cmd

import (
	"fmt"

	"github.com/belarte/metrix/repository"
	"github.com/spf13/cobra"
)

const (
	file = "bin/dev.sqlite"
)

var dbCmd = &cobra.Command{
	Use:   "db",
	Short: "Database operations",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("no command specified")
		}

		switch args[0] {
		case "migrate":
			db, err := repository.New(file)
			if err != nil {
				return fmt.Errorf("error opening database: %w", err)
			}

			err = db.Migrate()
			if err != nil {
				return fmt.Errorf("error migrating database: %w", err)
			}

			return err
		default:
			return fmt.Errorf("unknown command: %s", args[0])
		}
	},
}

func init() {
	rootCmd.AddCommand(dbCmd)
}
