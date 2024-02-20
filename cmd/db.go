package cmd

import (
	"fmt"
	"log"
	"strconv"

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

var dbListMetricsCmd = &cobra.Command{
	Use:   "list-metrics",
	Short: "List metrics",
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Println("Listing metrics", database)

		db, err := repository.New(database)
		if err != nil {
			return fmt.Errorf("error opening database: %w", err)
		}

		metrics, err := db.GetMetrics()
		if err != nil {
			return fmt.Errorf("error listing metrics: %w", err)
		}

		fmt.Println(" # Title            Unit     Description")
		for _, metric := range metrics {
			fmt.Printf("%2d %-16s %-8s %s\n", metric.ID, metric.Title, metric.Unit, metric.Description)
		}
		return err
	},
}

var dbListEntriesCmd = &cobra.Command{
	Use:   "list-entries",
	Short: "List entries",
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Println("Listing entries", database)

		if len(args) != 1 {
			return fmt.Errorf("missing metric ID")
		}

		metricID, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid metric ID: %w", err)
		}

		db, err := repository.New(database)
		if err != nil {
			return fmt.Errorf("error opening database: %w", err)
		}

		entries, err := db.GetSortedEntriesForMetric(metricID)
		if err != nil {
			return fmt.Errorf("error listing entries: %w", err)
		}

		fmt.Println("Date       Value")
		for _, entry := range entries {
			fmt.Printf("%s %.2f\n", entry.Date, entry.Value)
		}
		return err
	},
}

func init() {
	dbCmd.AddCommand(dbMigrateCmd)
	dbCmd.AddCommand(dbListMetricsCmd)
	dbCmd.AddCommand(dbListEntriesCmd)
	rootCmd.AddCommand(dbCmd)
}
