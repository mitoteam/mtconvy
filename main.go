package main

import (
	"log"
	"os"

	"github.com/mitoteam/mtconvy/app"
	"github.com/mitoteam/mttools"
)

func main() {
	var current_path string

	//Read path from args if available
	if len(os.Args) > 1 {
		current_path = os.Args[1]
	} else {
		var err error
		current_path, err = os.Getwd()

		if err != nil {
			log.Fatal(err)
		}
	}

	if !mttools.IsDirExists(current_path) {
		log.Fatalf("Directory %s does not exist", current_path)
		os.Exit(-1)
	}

	//Deal with settings
	app.AppSettings.Load(current_path)
	app.AppSettings.Print()

	if !app.AppSettings.Check() {
		os.Exit(-1)
		return
	}

	//Create task
	log.Printf("Current path: %s", current_path)

	task := app.NewTask(current_path)

	task.SelectFiles()
	task.SelectStreams()
}
