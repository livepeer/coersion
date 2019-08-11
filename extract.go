package main

import (
	"io"
	"os/exec"
)

// ExtractImageTask uses ffmpeg to perform the extraction of a still image from a video source
type ExtractImageTask struct {
	source string
	cmd    *exec.Cmd
}

// NewExtractImageTask creates a new image extraction task
func NewExtractImageTask(source string) *ExtractImageTask {

	// TO FILE:
	// return &ExtractImageTask{
	// 	source: source,
	// 	dest:   dest,
	// 	cmd: exec.Command(ffmpegCmd, []string{
	// 		"-ss", "0",
	// 		"-i", source,
	// 		"-t", "1",
	// 		"-q:v", "2",
	// 		"-vf", `select="eq(pict_type\,PICT_TYPE_I)"`,
	// 		"-vsync", "0",
	// 		fmt.Sprintf("%s%%03d.jpg", dest),
	// 	}...),
	// }
	return &ExtractImageTask{
		source: source,
		cmd: exec.Command(ffmpegCmd, []string{
			"-ss", "0",
			"-i", source,
			"-vframes", "1",
			"-q:v", "2",
			"-vf", `select="eq(pict_type\,PICT_TYPE_I)"`,
			"-vsync", "0",
			"-f", "image2pipe",
			"-vcodec", "jpg",
			"-",
		}...),
	}
}

// Start begins the task
func (j *ExtractImageTask) Start() error {
	return j.cmd.Start()
}

// StderrPipe returns the task's StderrPipe
func (j *ExtractImageTask) StderrPipe() (io.ReadCloser, error) {
	return j.cmd.StderrPipe()
}

// StdoutPipe returns the task's StdoutPipe
func (j *ExtractImageTask) StdoutPipe() (io.ReadCloser, error) {
	return j.cmd.StdoutPipe()
}
