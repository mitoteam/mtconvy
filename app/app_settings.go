package app

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/mitoteam/mttools"
)

type appSettingsType struct {
	FfmpegPath  string `yaml:"ffmpeg_path"`
	FfprobePath string `yaml:"ffprobe_path"`

	DefaultAudioCodec string `yaml:"default_audio_codec"`
}

const DefaultSettingsFilename = ".mtconvy.yml"

var AppSettings *appSettingsType

func init() {
	AppSettings = getDefaultAppSettings()
}

func getDefaultAppSettings() *appSettingsType {
	return &appSettingsType{
		FfmpegPath:        "ffmpeg",
		FfprobePath:       "ffprobe",
		DefaultAudioCodec: "eac3",
	}
}

func (s *appSettingsType) Load() {
	var filename string

	//1) look in current directory
	settingspath, err := os.Getwd()
	//log.Println(settingspath)

	if err == nil {
		filename = filepath.Join(settingspath, DefaultSettingsFilename)
	}

	//2) look near executable
	if !mttools.IsFileExists(filename) {
		settingspath, err = os.Executable()

		if err == nil {
			settingspath = filepath.Dir(settingspath)
			//log.Println(settingspath)
			filename = filepath.Join(settingspath, DefaultSettingsFilename)
		}
	}

	//3) look in homedir
	if !mttools.IsFileExists(filename) {
		settingspath, err = os.UserHomeDir()
		//log.Println(settingspath)

		if err == nil {
			filename = filepath.Join(settingspath, DefaultSettingsFilename)
		}
	}

	// Load settings
	if mttools.IsFileExists(filename) {
		log.Println("Settings file loaded: " + filename)

		mttools.LoadYamlSettingFromFile(filename, s)
	} else {
		log.Println("No " + DefaultSettingsFilename + " file found. Using default settings.")
	}
}

func (s *appSettingsType) Print() {
	mttools.PrintYamlSettings(s)
}

func (s *appSettingsType) Check() bool {
	//Check FFMPEG
	out, err := mttools.ExecCmd(s.FfmpegPath, []string{"-version"})

	if err != nil {
		log.Fatalln(err)
	}

	//read first line
	scanner := bufio.NewReader(strings.NewReader(out))
	out, _ = scanner.ReadString('\n')

	log.Print("FFmpeg found: " + out)

	//Check FFPROBE
	out, err = mttools.ExecCmd(s.FfprobePath, []string{"-version"})

	if err != nil {
		log.Fatalln(err)
	}

	//read first line
	scanner = bufio.NewReader(strings.NewReader(out))
	out, _ = scanner.ReadString('\n')

	log.Print("FFprobe found: " + out)

	return true
}
