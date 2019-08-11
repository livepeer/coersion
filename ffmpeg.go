package main

import (
	"io"
	"os/exec"
)

const ffmpegCmd = "ffmpeg"

type ffmpegTask interface {
	// getCmd gets the internal exec.cmd
	getCmd() *exec.Cmd
}

// Start begins the task.
func Start(f ffmpegTask) error {
	return f.getCmd().Start()
}

// StdoutPipe returns the task's StdoutPipe
func StdoutPipe(f ffmpegTask) (io.ReadCloser, error) {
	return f.getCmd().StdoutPipe()
}
