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

// Run begins the task (via Start()) and waits for completion
func (j *ExtractImageTask) Run() error {
	if err := j.Start(); err != nil {
		return err
	}
	return j.cmd.Wait()
}

// Output returns the task's output
func (j *ExtractImageTask) Output() ([]byte, error) {
	return j.cmd.Output()
}

// StderrPipe returns the task's StderrPipe
func (j *ExtractImageTask) StderrPipe() (io.ReadCloser, error) {
	return j.cmd.StderrPipe()
}

// Promise runs the task asynchronously and returns a channel
// that will emit the image extraction task's status when it completes
func (j *ExtractImageTask) Promise() chan error {
	ch := make(chan error)
	go func() {
		if err := j.Start(); err != nil {
			ch <- err
		}
		if err := j.cmd.Wait(); err != nil {
			ch <- err
		}
		ch <- nil
	}()
	return ch
}
