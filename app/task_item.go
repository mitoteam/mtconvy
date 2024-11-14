package app

import (
	"fmt"
	"log"
	"path/filepath"
	"strconv"
	"strings"
	"time"

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

	if len(stream_list) > 0 {
		//Prepare default selection
		default_selected := make([]int, 0)
		default_stream_map := make(map[string]bool)

		for i := 0; i < len(stream_list); i++ {
			stream := stream_list[i]

			// at least one video stream
			if stream.Data.CodecType == "video" {
				_, has_stream := default_stream_map["video"]

				if !has_stream {
					default_selected = append(default_selected, i)
					default_stream_map["video"] = true
				}
			}

			if stream.Data.CodecType == "audio" || stream.Data.CodecType == "subtitle" {
				for j := 0; j < len(AppSettings.Languages); j++ {
					language := AppSettings.Languages[j]
					key := stream.Data.CodecType + "_" + language

					_, has_stream := default_stream_map[key]

					if !has_stream && stream.Language == strings.ToUpper(language) {
						default_selected = append(default_selected, i)
						default_stream_map[key] = true
					}
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

		//fill with selected streams
		task_item.Streams = make([]FfStream, len(selected))

		for i := 0; i < len(selected); i++ {
			task_item.Streams[i] = stream_list[selected[i]]
		}
	} else {
		//clear
		task_item.Streams = make([]FfStream, 0)
	}
}

func (task_item *TaskItem) Convert() {
	//only if something was selected
	if len(task_item.Streams) > 0 {
		new_filename := filepath.Base(task_item.Path)
		new_filename = strings.TrimSuffix(new_filename, filepath.Ext(new_filename))
		new_filename = new_filename + "_" + AppSettings.Suffix + ".mkv"
		new_filename = filepath.Join(filepath.Dir(task_item.Path), new_filename)

		args := make([]string, 0, 10)

		//generic options
		args = append(args, "-y")                 //overwrite DST file silently
		args = append(args, "-hide_banner")       //do not print ffmpeg intro banner
		args = append(args, "-loglevel", "error") //be silent
		args = append(args, "-stats", "-stats_period", "5")

		//input file
		args = append(args, "-i", task_item.Path)

		//selected streams
		for i := 0; i < len(task_item.Streams); i++ {
			stream := task_item.Streams[i]

			//add stream from source
			args = append(args, "-map", "0:"+strconv.Itoa(stream.Index))

			selector := ":" + strconv.Itoa(i)

			if new_codec, exists := AppSettings.Conversions[stream.Data.CodecName]; exists {
				args = append(args, "-c"+selector, new_codec)

				if stream.Data.CodecType == "audio" {
					args = append(args, "-b"+selector, AppSettings.AudioBitrate)
				}
			} else {
				//copy stream without re-encoding
				args = append(args, "-c"+selector, "copy")
			}
		}

		//output file
		args = append(args, new_filename)

		fmt.Println("\nStarting ffmpeg for", task_item.Name)
		//fmt.Println("ARGS:", args)

		//call ffmpeg
		//fmt.Print(args)
		start := time.Now()
		if err := mttools.ExecCmdWaitPrint(AppSettings.FfmpegPath, args); err != nil {
			log.Fatal(err.Error())
		}

		elapsed := time.Since(start).Round(time.Second)

		fmt.Printf("Done. Took %s.\n", elapsed)
	}
}
