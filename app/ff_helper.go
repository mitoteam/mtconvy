package app

import (
	"encoding/json"
	"log"
	"strconv"
	"strings"

	"github.com/mitoteam/mttools"
)

type FfStream struct {
	Index    int
	Name     string
	Language string
	Data     jsonStream
}

type jsonFull struct {
	Streams []jsonStream `json:"streams"`
}

type jsonStream struct {
	Index       int               `json:"index"`
	CodecName   string            `json:"codec_name"`
	CodecType   string            `json:"codec_type"`
	Width       int               `json:"width"`
	Height      int               `json:"height"`
	Disposition map[string]string `json:"disposition"`
	Tags        map[string]string `json:"tags"`
}

func FfGetStreamList(path string) []FfStream {
	args := make([]string, 0)
	args = append(args, "-v", "quiet", "-print_format", "json", "-show_streams", path)

	json_str, err := mttools.ExecCmd(AppSettings.FfprobePath, args)
	if err != nil {
		log.Fatal(err)
	}

	data := jsonFull{}

	json.Unmarshal([]byte(json_str), &data)
	//fmt.Println(data)

	list := make([]FfStream, 0, len(data.Streams))

	for _, streamData := range data.Streams {
		stream := FfStream{
			Index: streamData.Index,
			Data:  streamData,
		}

		stream.Name = streamData.CodecType + "/" + streamData.CodecName

		// resolution
		if streamData.CodecType == "video" {
			stream.Name += " " + strconv.Itoa(streamData.Width) + "x" + strconv.Itoa(streamData.Height)
		}

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

		//title
		if streamData.CodecType == "audio" || streamData.CodecType == "subtitle" {
			if v, exists := streamData.Tags["title"]; exists {
				stream.Name += " " + v
			}
		}

		list = append(list, stream)
	}

	return list
}
