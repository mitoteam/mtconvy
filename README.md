# mtconvy

![GitHub](https://img.shields.io/github/license/mitoteam/mtconvy)
[![GitHub Version](https://img.shields.io/github/v/release/mitoteam/mtconvy?logo=github)](https://github.com/mitoteam/mtconvy)
[![GitHub Release Date](https://img.shields.io/github/release-date/mitoteam/mtconvy)](https://github.com/mitoteam/mtconvy/releases)
[![GitHub contributors](https://img.shields.io/github/contributors-anon/mitoteam/mtconvy)](https://github.com/mitoteam/mtconvy/graphs/contributors)
[![GitHub commit activity](https://img.shields.io/github/commit-activity/y/mitoteam/mtconvy)](https://github.com/mitoteam/mtconvy/commits)
[![GitHub downloads](https://img.shields.io/github/downloads/mitoteam/mtconvy/total)](https://github.com/mitoteam/mtconvy/releases)
[![Build&Tests](https://github.com/mitoteam/mtconvy/actions/workflows/tests.yml/badge.svg)](https://github.com/mitoteam/mtconvy/actions/workflows/tests.yml)

mtconvy - ffmpeg command-line helper utility to convert DTS audio tracks in video files to AC3 or AAC ones using ffmpeg utility.

LG dropped DTS support for the 2020–2022 models of their TVs and removed it once again starting with the 2025 models.
So if you downloaded movie wuth high-quailty DTS sound you have to coinvert it to AC3 or AAC codec.

This is done well with `ffmpeg`. But there two problems:

1) ffmpeg has very complicated command-line syntax hard to keep in memory
2) you should manually explore available tracks with `ffprobe` or `mediainfo` to know which tracks to convert

`mtconvy` utility makes selection of tracks and conversion very simple and easy. It also support file renaming.

Works well on PCs under Linux/Windows and also on Linux-based NASes.

## Installation

### Using Scoop (Windows)

Scoop is useful command-line installer and updater for Windows.

* Install `scoop` (_if you have not already_): https://scoop.sh

* Add bucket (_if you have not already_):

```sh
scoop bucket add mitoteam https://github.com/mitoteam/scoop-bucket
```

* Install:

```sh
scoop install mitoteam/mtconvy
```

* Update:

```sh
scoop update mtconvy
```

### Manual installation (Windows or Linux)

* Download latest release from [Releases](https://github.com/mitoteam/mtconvy/releases) page.
* Unpack with 7-zip.
* Add path to mtconvy to system's or user's PATH variable.

## Usage

Just run `mtconvy` in directory you want to convert some file(s). It will ask you to choose files to convert first. Then it will ask you what streams to keep for each file.

You can adjust some options in config file `.mtconvy.yml`.

## Feedback

Please feel free to create issues to discuss new fetures or improve utilty. Bugreports are also always welcomed.
