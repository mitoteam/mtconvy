package cmd

import (
	"fmt"

	"github.com/mitoteam/goapp"
	"github.com/mitoteam/mtconvy/app"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Prints the raw version number of " + app.AppName,

		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(goapp.BuildVersion)
		},
	})
}
