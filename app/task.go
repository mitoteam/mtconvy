package app

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"

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

	//prepend "All" option
	options_list = append([]string{"All files"}, options_list...)

	fmt.Println("\n*** Please select files to process")

	for v, k := range options_list {
		fmt.Printf("%2d: %s\n", v, k)
	}

	fmt.Println("*** Enter file numbers separated by space or comma and press Enter. Default is \"0\".")
	fmt.Print("Files: ")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	err = scanner.Err()
	if err != nil {
		log.Fatal(err)
	}

	user_input := strings.TrimSpace(scanner.Text())
	user_input = strings.ReplaceAll(user_input, ",", " ")

	number_string_list := strings.Split(user_input, " ")
	numbers_list := make([]int, 0, len(number_string_list))

	if len(user_input) == 0 {
		numbers_list = append(numbers_list, 0)
	}

	for _, v := range number_string_list {
		if len(v) < 1 {
			continue
		}

		n, err := strconv.Atoi(v)
		if err != nil {
			log.Printf("Skipping wrong input: %s (%s)", v, err.Error())
			continue
		}

		if n < 0 || n >= len(options_list) {
			continue
		}

		numbers_list = append(numbers_list, n)
	}

	//TODO: remove duplicates

	//all files
	if len(numbers_list) == 1 && numbers_list[0] == 0 {
		numbers_list = make([]int, len(options_list)-1)
		for i := range numbers_list {
			numbers_list[i] = i + 1
		}
	}

	fmt.Printf("Selected files: %v", numbers_list)
}
