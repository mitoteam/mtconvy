# mtconvy

mtconvy - ffmpeg command-line helper utility to convert DTS audio tracks in video files to AC3 or AAC ones using ffmpeg utility.

LG TVs does no support DTS codecs in videos since 2020. So conversion to AC3 is required.
It is well done with ffmpeg. But there two problems: 1) ffmpeg has very complicated command-line syntax hard to keep in memory
2) you should manually explore available tracks with `ffprobe` or `mediainfo` to know whicj tracks to convert.

This utility makes selection of tracks and conversion very simple and easy.
