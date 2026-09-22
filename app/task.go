package app

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	huh "charm.land/huh/v2"
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

func (t *Task) SelectFiles() error {
	directory_files_list, err := os.ReadDir(t.path)
	if err != nil {
		return err
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
		var numbers_list []int //result
		var huh_options []huh.Option[int]

		for index, fileName := range options_list {
			huh_options = append(huh_options, huh.NewOption(fileName, index))
		}

		// Build the multi-select form
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewMultiSelect[int]().
					Title("Please select files to process:").
					Description("Just press ENTER to select all files").
					Options(huh_options...).
					Height(GetHuhMultiselectHeight(len(huh_options))).
					Value(&numbers_list), // Binds selected values here
			),
		)

		if err := form.Run(); err != nil {
			return err
		}

		//nothing selected = all files
		if len(numbers_list) == 0 {
			for i := 0; i < len(options_list); i++ {
				numbers_list = append(numbers_list, i)
			}
		}

		//Create task items
		for i := 0; i < len(numbers_list); i++ {
			file_name := files_list[numbers_list[i]]

			task_item := TaskItem{
				Name:     options_list[numbers_list[i]],
				Path:     filepath.Join(t.path, file_name),
				BaseName: strings.TrimSuffix(file_name, filepath.Ext(file_name)),
				Ext:      filepath.Ext(file_name),
			}

			t.items = append(t.items, &task_item)
		}
	} else {
		fmt.Printf("No %s files found in current directory.", extensions)
	}

	return nil
}

func (t *Task) SelectStreams() error {
	for i := 0; i < len(t.items); i++ {
		if err := t.items[i].SelectStreams(); err != nil {
			return err
		}
	}

	return nil
}

func (t *Task) Convert() error {
	for i := 0; i < len(t.items); i++ {
		if err := t.items[i].Convert(); err != nil {
			return err
		}
	}

	return nil
}
