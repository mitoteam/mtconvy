package app

import (
	"fmt"
	"strconv"
)

type TaskItem struct {
	Name string
	Path string

	Streams []FfStream
}

func (task_item TaskItem) SelectStreams() {
	fmt.Println()
	fmt.Printf("*** FILE: %s\n", task_item.Name)
	fmt.Println("Running ffprobe...")

	stream_list := FfGetStreamList(task_item.Path)

	//Prepare default selection
	default_selected := make([]int, 0)
	var has_video, has_rus_audio, has_eng_audio, has_rus_st, has_eng_st bool

	for i, stream := range stream_list {
		if !has_video && stream.Data.CodecType == "video" {
			default_selected = append(default_selected, i)
			has_video = true
		}

		if stream.Data.CodecType == "audio" {
			if stream.Language == "RUS" && !has_rus_audio {
				default_selected = append(default_selected, i)
				has_rus_audio = true
			}

			if stream.Language == "ENG" && !has_eng_audio {
				default_selected = append(default_selected, i)
				has_eng_audio = true
			}
		}

		if stream.Data.CodecType == "subtitle" {
			if stream.Language == "RUS" && !has_rus_st {
				default_selected = append(default_selected, i)
				has_rus_st = true
			}

			if stream.Language == "ENG" && !has_eng_st {
				default_selected = append(default_selected, i)
				has_eng_st = true
			}
		}
	}

	options_list := make([]string, 0, len(stream_list))
	for _, stream := range stream_list {
		options_list = append(options_list, stream.Name)
	}

	var default_choice string
	for _, i := range default_selected {
		if default_choice != "" {
			default_choice += " "
		}
		default_choice += strconv.Itoa(i + 1)

		options_list[i] = "* " + options_list[i]
	}

	selected := AskUserChoice(
		"Select streams",
		"Your choice (default: "+default_choice+"): ",
		options_list,
	)

	if len(selected) == 0 {
		selected = default_selected
	}

	task_item.Streams = make([]FfStream, len(selected))

	for i, v := range selected {
		task_item.Streams[i] = stream_list[v]
	}
}
