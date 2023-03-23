# mtconvy

[![Go Report Card](https://goreportcard.com/badge/github.com/mitoteam/mtconvy)](https://goreportcard.com/report/github.com/mitoteam/mtconvy)
![GitHub](https://img.shields.io/github/license/mitoteam/mtconvy)

[![GitHub Version](https://img.shields.io/github/v/release/mitoteam/mtconvy?logo=github)](https://github.com/mitoteam/mtconvy)
[![GitHub Release Date](https://img.shields.io/github/release-date/mitoteam/mtconvy)](https://github.com/mitoteam/mtconvy/releases)
![GitHub code size in bytes](https://img.shields.io/github/languages/code-size/mitoteam/mtconvy)
[![GitHub contributors](https://img.shields.io/github/contributors-anon/mitoteam/mtconvy)](https://github.com/mitoteam/mtconvy/graphs/contributors)
[![GitHub commit activity](https://img.shields.io/github/commit-activity/y/mitoteam/mtconvy)](https://github.com/mitoteam/mtconvy/commits)
[![GitHub downloads](https://img.shields.io/github/downloads/mitoteam/mtconvy/total)](https://github.com/mitoteam/mtconvy/releases)

mtconvy - ffmpeg command-line helper utility to convert DTS audio tracks in video files to AC3 or AAC ones using ffmpeg utility.

LG TVs does no support DTS codecs in videos since 2020. So conversion to AC3 is required.
It is well done with ffmpeg. But there two problems: 1) ffmpeg has very complicated command-line syntax hard to keep in memory
2) you should manually explore available tracks with `ffprobe` or `mediainfo` to know which tracks to convert.

This utility makes selection of tracks and conversion very simple and easy.
