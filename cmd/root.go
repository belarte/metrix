package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var database string

const inMemoryDatabase = ":memory:"

var rootCmd = &cobra.Command{
	Use: "metrix",
}

func init() {
	viper.BindEnv("METRIX_DB")
	viper.SetDefault("METRIX_DB", inMemoryDatabase)
	rootCmd.PersistentFlags().StringVarP(&database, "database", "d", viper.GetString("METRIX_DB"), "path to the database")
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
