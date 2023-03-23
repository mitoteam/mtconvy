package app

import "fmt"

type TaskItem struct {
	Name string
	Path string
}

func (i TaskItem) SelectStreams() {
	fmt.Printf("*** FILE: %s\n", i.Name)
	fmt.Println("Running ffprobe...")

	stream_list := FfGetStreamList(i.Path)

	options_list := make([]string, 0, len(stream_list))
	for _, stream := range stream_list {
		options_list = append(options_list, stream.Name)
	}

	AskUserChoice(
		"Select streams",
		"Your choice: ",
		options_list,
	)
}
