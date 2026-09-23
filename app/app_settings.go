package app

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mitoteam/mttools"
)

type appSettingsType struct {
	FfmpegPath  string `yaml:"ffmpeg_path"`
	FfprobePath string `yaml:"ffprobe_path"`

	Conversions  map[string]string `yaml:"conversions"`
	AudioBitrate string            `yaml:"audio_bitrate"`

	Suffix string `yaml:"suffix"`

	ReplaceOriginal bool              `yaml:"replace_original"`
	OriginalSuffix  string            `yaml:"original_suffix"`
	StrReplace      map[string]string `yaml:"str_replace"`

	Languages []string `yaml:"languages"`
}

const DefaultSettingsFilename = ".mtconvy.yml"

var AppSettings *appSettingsType

func init() {
	AppSettings = getDefaultAppSettings()
}

func getDefaultAppSettings() *appSettingsType {
	settings := appSettingsType{
		FfmpegPath:  "ffmpeg",
		FfprobePath: "ffprobe",

		Conversions:  map[string]string{"dts": "eac3"},
		AudioBitrate: "640k",

		Suffix:         "CONVERTED",
		OriginalSuffix: "ORIGINAL",
		StrReplace:     map[string]string{},

		Languages: []string{"ENG", "RUS"},
	}

	return &settings
}

func (s *appSettingsType) Load(path string) {
	//1) look in current directory
	settings_file_path := filepath.Join(path, DefaultSettingsFilename)

	//2) look near executable
	if !mttools.IsFileExists(settings_file_path) {
		if dir, err := os.Executable(); err == nil {
			dir = filepath.Dir(dir)
			//log.Println(dir)
			settings_file_path = filepath.Join(dir, DefaultSettingsFilename)
		}
	}

	//3) look in homedir
	if !mttools.IsFileExists(settings_file_path) {
		if dir, err := os.UserHomeDir(); err == nil {
			//log.Println(dir)
			settings_file_path = filepath.Join(dir, DefaultSettingsFilename)
		}
	}

	// Load settings
	if mttools.IsFileExists(settings_file_path) {
		fmt.Println("Settings file loaded: " + settings_file_path)

		mttools.LoadYamlSettingFromFile(settings_file_path, s)
	} else {
		fmt.Println("No " + DefaultSettingsFilename + " file found. Using default settings.")
	}
}

func (s *appSettingsType) Print() {
	mttools.PrintYamlSettings(s)
}

func (s *appSettingsType) Check() error {
	//Check FFMPEG
	out, err := mttools.ExecCmd(s.FfmpegPath, []string{"-version"})

	if err != nil {
		return err
	}

	//read first line
	scanner := bufio.NewReader(strings.NewReader(out))
	out, _ = scanner.ReadString('\n')

	fmt.Print("FFmpeg found: " + out)

	//Check FFPROBE
	out, err = mttools.ExecCmd(s.FfprobePath, []string{"-version"})

	if err != nil {
		return err
	}

	//read first line
	scanner = bufio.NewReader(strings.NewReader(out))
	out, _ = scanner.ReadString('\n')

	fmt.Print("FFprobe found: " + out)

	return nil
}
