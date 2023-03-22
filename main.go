package main

import (
	"os"

	"github.com/mitoteam/mtconvy/app"
)

func main() {
	app.AppSettings.Load()
	app.AppSettings.Print()

	if !app.AppSettings.Check() {
		os.Exit(-1)
		return
	}
}
