package app

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"time"

	huh "charm.land/huh/v2"
	"github.com/mitoteam/goapp"
	"github.com/mitoteam/mttools"
)

type TaskItem struct {
	Name     string //task name (filename + size)
	BaseName string //filename without extension
	Ext      string //file extension
	Path     string //full path with filename

	skipTask bool // do not convert this file at all

	Streams []FfStream
}

func (task_item *TaskItem) SelectStreams() error {
	//clear list
	task_item.Streams = make([]FfStream, 0)

	// show available streams and ask user
	fmt.Println()
	fmt.Println("Running ffprobe for " + task_item.Name + "...")

	stream_list, err := FfGetStreamList(task_item.Path)

	if err != nil {
		return err
	}

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

		var selected []int
		var huh_options []huh.Option[int]

		for index, stream := range stream_list {
			var selected = slices.Contains(default_selected, index)
			var o = huh.NewOption(stream.Name, index).Selected(selected)

			huh_options = append(huh_options, o)
		}

		var skip_option = huh.NewOption("* Skip this file (and ignore streams selection)", -1)
		huh_options = append(huh_options, skip_option)

		// Build the multi-select form
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewMultiSelect[int]().
					Title(task_item.Name).
					Description("Please select streams to include to output:").
					Options(huh_options...).
					Height(GetHuhMultiselectHeight(len(huh_options))).
					Value(&selected), // Binds selected values here
			),
		)

		if err := form.Run(); err != nil {
			return err
		}

		//fmt.Printf("DBG selection: %v", selected)

		if slices.Contains(selected, -1) {
			task_item.skipTask = true
			return nil
		}

		//fill with selected streams
		task_item.Streams = make([]FfStream, len(selected))

		for i := 0; i < len(selected); i++ {
			task_item.Streams[i] = stream_list[selected[i]]
		}
	}

	return nil
}

func (task_item *TaskItem) Convert() error {
	if task_item.skipTask {
		fmt.Println("\nSkipping conversion for " + task_item.Name)
		return nil
	}

	//only if something was selected
	if len(task_item.Streams) > 0 {
		new_filename := filepath.Join(filepath.Dir(task_item.Path), task_item.BaseName+"_"+AppSettings.Suffix+task_item.Ext)

		args := make([]string, 0, 10)

		//generic options
		args = append(args, "-y")                           //overwrite DST file silently
		args = append(args, "-hide_banner")                 //do not print ffmpeg intro banner
		args = append(args, "-loglevel", "error")           //be silent
		args = append(args, "-stats", "-stats_period", "5") //update progress every 5 seconds

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
		goapp.PrintDev("ffmpeg command: " + AppSettings.FfmpegPath + " " + strings.Join(args, " "))

		//call ffmpeg
		//fmt.Print(args)
		start := time.Now()
		if _, err := mttools.ExecCmdWaitAndPrint(AppSettings.FfmpegPath, args); err != nil {
			return err
		}

		elapsed := time.Since(start).Round(time.Second)

		fmt.Printf("Done. Took %s.\n", elapsed)

		if AppSettings.ReplaceOriginal {
			//rename original file
			original_filename := filepath.Join(filepath.Dir(task_item.Path), task_item.BaseName+"_"+AppSettings.OriginalSuffix+task_item.Ext)

			if err := os.Rename(task_item.Path, original_filename); err != nil {
				return fmt.Errorf("Failed to rename original file: %v", err)
			}

			//rename converted file to original name
			if err := os.Rename(new_filename, task_item.Path); err != nil {
				return fmt.Errorf("Failed to rename converted file: %v", err)
			}
		}
	}

	return nil
}
