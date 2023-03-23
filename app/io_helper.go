package app

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"golang.org/x/exp/maps"
)

func AskUserChoice(
	message, prompt string,
	options_list []string,
) []int {
	fmt.Println("*** " + message)

	for k, v := range options_list {
		fmt.Printf("%2d: %s\n", k+1, v)
	}

	fmt.Print("*** " + prompt)

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	err := scanner.Err()
	if err != nil {
		log.Fatal(err)
	}

	user_input := strings.TrimSpace(scanner.Text())
	user_input = strings.ReplaceAll(user_input, ",", " ")

	number_string_list := strings.Split(user_input, " ")
	numbers_list := make(map[int]int, len(number_string_list))

	if len(user_input) == 0 {
		return []int{} //empty array
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

		if n < 1 || n > len(options_list) {
			continue
		}

		numbers_list[n] = n - 1
	}

	return maps.Values(numbers_list)
}
