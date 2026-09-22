package app

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/mitoteam/goapp"
	"github.com/mitoteam/mttools"
)

type FfStream struct {
	Index    int
	Name     string
	Language string
	Size     int64
	Data     jsonStream
}

type jsonFull struct {
	Streams []jsonStream `json:"streams"`
}

type jsonStream struct {
	Index         int               `json:"index"`
	CodecName     string            `json:"codec_name"`
	CodecType     string            `json:"codec_type"`
	ChannelLayout string            `json:"channel_layout"`
	Channels      int               `json:"channels"`
	Width         int               `json:"width"`
	Height        int               `json:"height"`
	Disposition   map[string]string `json:"disposition"`
	Tags          map[string]string `json:"tags"`
}

func FfGetStreamList(path string) ([]FfStream, error) {
	args := make([]string, 0)
	args = append(args, "-hide_banner", "-v", "quiet", "-print_format", "json", "-show_streams", path)

	goapp.PrintDev("ffprobe cmd: " + AppSettings.FfprobePath + " " + strings.Join(args, " "))

	json_str, err := mttools.ExecCmd(AppSettings.FfprobePath, args)
	if err != nil {
		return make([]FfStream, 0), fmt.Errorf("Error running ffprobe for %s: %s", filepath.Base(path), err.Error())
	}

	data := jsonFull{}

	json.Unmarshal([]byte(json_str), &data)
	//fmt.Println(data)

	list := make([]FfStream, 0, len(data.Streams))

	for i := 0; i < len(data.Streams); i++ {
		streamData := data.Streams[i]

		stream := FfStream{
			Index: streamData.Index,
			Data:  streamData,
		}

		stream.Name = streamData.CodecType + "/" + streamData.CodecName

		//language
		if streamData.CodecType == "audio" || streamData.CodecType == "subtitle" {
			if v, exists := streamData.Tags["language"]; exists {
				stream.Language = strings.ToUpper(v)

				if stream.Language == "RU" {
					stream.Language = "RUS"
				} else if stream.Language == "EN" {
					stream.Language = "ENG"
				}

				stream.Name += " " + stream.Language
			}
		}

		// resolution
		if streamData.CodecType == "video" {
			stream.Name += " " + strconv.Itoa(streamData.Width) + "x" + strconv.Itoa(streamData.Height)
		}

		// channels
		if streamData.CodecType == "audio" {
			if streamData.ChannelLayout != "" {
				stream.Name += " " + streamData.ChannelLayout
			}
		}

		//title or name
		if streamData.CodecType == "audio" || streamData.CodecType == "subtitle" {
			if v, exists := streamData.Tags["title"]; exists {
				stream.Name += " \"" + v + "\""
			} else if v, exists := streamData.Tags["name"]; exists {
				stream.Name += " \"" + v + "\""
			}
		}

		//stream size
		// checking tags for name with NUMBER_OF_BYTES prefix
		for key, val := range streamData.Tags {
			if strings.HasPrefix(key, "NUMBER_OF_BYTES") {
				// convert to int64
				if parsedSize, err := strconv.ParseInt(val, 10, 64); err == nil {
					stream.Size = parsedSize
					break
				}
			}
		}

		if stream.Size > 0 {
			stream.Name += ", " + mttools.FormatFileSize(int64(stream.Size))
		}

		// append to list
		list = append(list, stream)
	}

	return list, nil
}
