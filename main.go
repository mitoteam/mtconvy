package main

import (
	"log"
	"os"

	"github.com/mitoteam/mtconvy/app"
)

func main() {
	//Deal with settings
	app.AppSettings.Load()
	app.AppSettings.Print()

	if !app.AppSettings.Check() {
		os.Exit(-1)
		return
	}

	//Create task
	current_path, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Current path: %s", current_path)

	task := app.NewTask(current_path)

	task.SelectFiles()
}
