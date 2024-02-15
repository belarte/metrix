package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var database string

const inMemoryDatabase = ":memory:"

var rootCmd = &cobra.Command{
	Use: "metrix",
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&database, "database", "d", inMemoryDatabase, "path to the database")
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
