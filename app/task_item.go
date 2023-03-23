package app

import (
	"fmt"
	"log"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/mitoteam/mttools"
)

type TaskItem struct {
	Name string
	Path string

	Streams []FfStream
}

func (task_item *TaskItem) SelectStreams() {
	fmt.Println()
	fmt.Printf("*** FILE: %s\n", task_item.Name)
	fmt.Println("Running ffprobe...")

	stream_list := FfGetStreamList(task_item.Path)

	//Prepare default selection
	default_selected := make([]int, 0)
	var has_video, has_rus_audio, has_eng_audio, has_rus_st, has_eng_st bool

	for i := 0; i < len(stream_list); i++ {
		stream := stream_list[i]

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
	for i := 0; i < len(stream_list); i++ {
		options_list = append(options_list, stream_list[i].Name)
	}

	var default_choice string
	for i := 0; i < len(default_selected); i++ {
		if default_choice != "" {
			default_choice += " "
		}
		default_choice += strconv.Itoa(default_selected[i] + 1)

		options_list[default_selected[i]] = "* " + options_list[default_selected[i]]
	}

	selected := AskUserChoice(
		"Please select streams to include to output",
		"Your choice (default: "+default_choice+"): ",
		options_list,
	)

	if len(selected) == 0 {
		selected = default_selected
	}

	//clear
	task_item.Streams = make([]FfStream, len(selected))

	for i := 0; i < len(selected); i++ {
		task_item.Streams[i] = stream_list[selected[i]]
	}
}

func (task_item *TaskItem) Convert() {
	new_filename := filepath.Base(task_item.Path)
	new_filename = strings.TrimSuffix(new_filename, filepath.Ext(new_filename))
	new_filename = new_filename + "_" + AppSettings.DefaultAudioCodec + ".mkv"
	new_filename = filepath.Join(filepath.Dir(task_item.Path), new_filename)

	args := make([]string, 0, 10)

	//generic options
	args = append(args, "-y")                 //overwrite DST file silently
	args = append(args, "-hide_banner")       //do not print FFPMEG intro banner
	args = append(args, "-loglevel", "error") //be silent
	args = append(args, "-stats", "-stats_period", "5")

	//input file
	args = append(args, "-i", task_item.Path)

	//selected streams
	for i := 0; i < len(task_item.Streams); i++ {
		stream := task_item.Streams[i]
		args = append(args, "-map", "0:"+strconv.Itoa(stream.Index))

		if stream.Data.CodecType == "video" || stream.Data.CodecType == "subtitle" {
			args = append(args, "-c:"+strconv.Itoa(i), "copy")
		}

		if stream.Data.CodecType == "audio" {
			if stream.Data.CodecName == AppSettings.DefaultAudioCodec || stream.Data.CodecName == "ac3" {
				args = append(args, "-c:"+strconv.Itoa(i), "copy")
			} else {
				args = append(args, "-c:"+strconv.Itoa(i), AppSettings.DefaultAudioCodec)
				args = append(args, "-b:"+strconv.Itoa(i), "640k")
			}
		}
	}

	//output file
	args = append(args, new_filename)

	fmt.Println("\nStarting ffmpeg for", task_item.Name)
	//fmt.Println("ARGS:", args)

	//call ffmpeg
	//fmt.Print(args)
	if err := mttools.ExecCmdWaitPrint(AppSettings.FfmpegPath, args); err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println("Done...")
}
