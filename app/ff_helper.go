package app

import (
	"log"

	"github.com/mitoteam/mttools"
)

type FfStream struct {
	Name string
}

func FfGetStreamList(path string) []FfStream {
	args := make([]string, 0)
	args = append(args, "-v", "quiet", "-print_format", "json", "-show_streams", path)

	_, err := mttools.ExecCmd(AppSettings.FfprobePath, args)

	if err != nil {
		log.Fatal(err)
	}

	list := make([]FfStream, 0)

	return list
}
