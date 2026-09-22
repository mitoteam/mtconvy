// Package cmd provides CLI commands, flags and arguments handling.
// spf13/cobra based.
package cmd

import (
	"github.com/mitoteam/goapp"
	"github.com/mitoteam/mtconvy/app"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     app.AppName,
	Version: goapp.BuildVersion,
	Long:    app.AppName + " v" + goapp.BuildVersion + " - " + app.AppDescription + "\n",

	CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true}, //disable 'completion' subcommand
	SilenceUsage:      true,                                             // Suppresses the help text output when an error occurs
	SilenceErrors:     false,                                            // Keep this false so Cobra still prints the error message itself
	Args:              cobra.ArbitraryArgs,                              // do not return error for unknown subcommand and just pass it to args

	RunE: func(cmd *cobra.Command, args []string) error {
		//use `run` command if no command given
		return runCmd.RunE(cmd, args)
	},
}

func Root() *cobra.Command {
	return rootCmd
}
