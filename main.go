package main

import (
	"os"

	"github.com/mitoteam/mtconvy/cmd"
)

func main() {
	//cli application - we just let cobra to do it job
	if err := cmd.Root().Execute(); err != nil {
		//log.Fatalln(err)
		os.Exit(1)
	}
}
