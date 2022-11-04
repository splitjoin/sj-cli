/*
Copyright © 2022 SplitJoin <support@splitjoin.com>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "sj",
	Short: "A CLI tool for writing faster commit messages",
	Long: `SplitJoin is a productivity focused service for helping
developers communicate more easily and effectively.

With SplitJoin, you can write your commit messages in
an instant and then edit them as you wish.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	viper.AutomaticEnv()
	rootCmd.CompletionOptions.HiddenDefaultCmd = true
	commitCmd.Flags().StringP("endpoint", "e", viper.GetString("SJ_ENDPOINT"), "API endpoint")
}
