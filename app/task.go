package app

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"

	"github.com/mitoteam/mttools"
)

type Task struct {
	items []TaskItem

	path string
}

func NewTask(path string) *Task {
	if !mttools.IsDirExists(path) {
		log.Fatalf("Path %s does not exists", path)
	}

	t := &Task{}

	t.path = path

	return t
}

func (t *Task) SelectFiles() {
	files_list, err := os.ReadDir(t.path)
	if err != nil {
		log.Fatalln(err)
	}

	options_list := make([]string, 0, len(files_list))

	re := regexp.MustCompile(`^(.+)\.(mkv)$`)

	for _, value := range files_list {
		//skip directories
		if value.IsDir() {
			continue
		}

		//check by regex
		if !re.MatchString(value.Name()) {
			continue
		}

		options_list = append(options_list, value.Name())
	}

	//sort by name
	sort.Strings(options_list)

	//TODO: no files found!

	numbers_list := AskUserChoice(
		"Please select files to process",
		"Enter file numbers separated by space or comma and press Enter. Empty input means \"All Files\".\nYour choice: ",
		options_list,
	)

	//all files
	if len(numbers_list) == 0 {
		for i := 0; i < len(options_list); i++ {
			numbers_list = append(numbers_list, i)
		}
	}

	fmt.Println("Selected files:")
	for _, v := range numbers_list {
		fmt.Println(options_list[v])
	}
}
