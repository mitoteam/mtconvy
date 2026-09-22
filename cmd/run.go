package cmd

import (
	"fmt"
	"os"

	"github.com/mitoteam/goapp"
	"github.com/mitoteam/mtconvy/app"
	"github.com/mitoteam/mttools"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run [/path/to/directory]",
	Short: "Runs conversion in given directory",
	Long:  "Runs conversion in given directory. If no path is given current directory is used.",

	RunE: func(cmd *cobra.Command, args []string) error {
		var path string

		// read path from first argument if given or use current directory
		if len(args) > 0 {
			path = args[0]
		} else {
			var err error
			path, err = os.Getwd()

			if err != nil {
				return err
			}
		}

		return doConversion(path)
	},
}

func init() {
	// cmd.Flags().BoolVar(
	// 	&app.JobRuntimeOptions.Solid, "solid", false,
	// 	"Create solid archives.",
	// )

	rootCmd.AddCommand(runCmd)
}

func doConversion(path string) error {
	fmt.Println(app.AppName + " v" + goapp.BuildVersion + " - " + app.AppDescription)
	fmt.Println(goapp.MOTTO)
	fmt.Println()

	if !mttools.IsDirExists(path) {
		return fmt.Errorf("Directory %s does not exist", path)
	}

	//Deal with settings
	app.AppSettings.Load(path)
	app.AppSettings.Print()

	if err := app.AppSettings.Check(); err != nil {
		return err
	}

	fmt.Printf("Current directory: %s\n", path)

	//Create task
	task := app.NewTask(path)

	if err := task.SelectFiles(); err != nil {
		return err
	}

	if err := task.SelectStreams(); err != nil {
		return err
	}

	if err := task.Convert(); err != nil {
		return err
	}

	return nil //no errors
}
