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
	directory_files_list, err := os.ReadDir(t.path)
	if err != nil {
		log.Fatalln(err)
	}

	//options to display to user
	options_list := make([]string, 0, len(directory_files_list))
	//raw file names
	files_list := make([]string, 0, len(directory_files_list))

	extensions := "mkv|mp4|avi|m4v"

	re := regexp.MustCompile(`^(.+)\.(` + extensions + `)$`)

	for i := 0; i < len(directory_files_list); i++ {
		file_entry := directory_files_list[i]

		//skip directories
		if file_entry.IsDir() {
			continue
		}

		//check by regex
		if !re.MatchString(file_entry.Name()) {
			continue
		}

		option := file_entry.Name()

		info, err := file_entry.Info()
		if err == nil {
			option += ", " + mttools.FormatFileSize(info.Size())
		}

		options_list = append(options_list, option)
		files_list = append(files_list, file_entry.Name())
	}

	//sort by name
	sort.Strings(options_list)

	if len(options_list) > 0 {
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
				Path: filepath.Join(t.path, files_list[numbers_list[i]]),
			}

			t.items = append(t.items, &task_item)
		}
	} else {
		log.Printf("No %s files found in current directory.", extensions)
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
