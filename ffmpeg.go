package main

const ffmpegCmd = "ffmpeg"

type ffmpegTask interface {
	// Start begins the task.
	Start() error
}
