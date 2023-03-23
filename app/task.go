package app

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"

	"github.com/mitoteam/mttools"
)

type Task struct {
	items []*TaskItem

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

	for i := 0; i < len(files_list); i++ {
		file_entry := files_list[i]

		//skip directories
		if file_entry.IsDir() {
			continue
		}

		//check by regex
		if !re.MatchString(file_entry.Name()) {
			continue
		}

		options_list = append(options_list, file_entry.Name())
	}

	//sort by name
	sort.Strings(options_list)

	//TODO: no files found!

	fmt.Println()
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

	//Create task items
	for i := 0; i < len(numbers_list); i++ {
		task_item := TaskItem{
			Name: options_list[numbers_list[i]],
			Path: filepath.Join(t.path, options_list[numbers_list[i]]),
		}

		t.items = append(t.items, &task_item)
	}
}

func (t *Task) SelectStreams() {
	for i := 0; i < len(t.items); i++ {
		t.items[i].SelectStreams()
	}
}

func (t *Task) Convert() {
	for i := 0; i < len(t.items); i++ {
		t.items[i].Convert()
	}
}
