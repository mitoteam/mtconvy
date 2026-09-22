package app

import (
	"os"

	"github.com/charmbracelet/x/term"
)

func GetHuhMultiselectHeight(options_count int) int {
	_, height, err := term.GetSize(os.Stdout.Fd())
	if err != nil {
		height = 20 // default height
	}

	if height > options_count+2 {
		height = options_count + 2
	}

	return height
}
