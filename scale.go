package main

import (
	"fmt"
	"io"
	"os/exec"
	"time"
)

// ScaleImageTask uses ffmpeg to perform the scaling of a still image from a video source
type ScaleImageTask struct {
	source string
	cmd    *exec.Cmd
}

func fmtDuration(d time.Duration) string {
	d = d.Round(time.Second)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second
	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}

// NewScaleImageTask creates a new image Scaleion task
func NewScaleImageTask(source string, w, h uint, o time.Duration) *ScaleImageTask {
	// TO FILE:
	// return &ScaleImageTask{
	// 	source: source,
	// 	dest:   dest,
	// 	cmd: exec.Command(ffmpegCmd, []string{
	// 		"-i", source,
	// 		"-vf", fmt.Sprintf("scale=%d:%d", w, h),
	// 		fmt.Sprintf("%s.jpg", dest),
	// 	}...),
	// }
	offset := fmtDuration(o)
	return &ScaleImageTask{
		source: source,
		cmd: exec.Command(ffmpegCmd, []string{
			"-ss", offset,
			"-i", source,
			"-vframes", "1",
			"-vf", fmt.Sprintf("scale=%d:%d", w, h),
			"-f", "image2pipe",
			"-vcodec", "png",
			"-q:v", "2",
			"-",
		}...),
	}
}

// Start begins the task
func (j *ScaleImageTask) Start() error {
	return j.cmd.Start()
}

// Run begins the task (via Start()) and waits for completion
func (j *ScaleImageTask) Run() error {
	if err := j.Start(); err != nil {
		return err
	}
	return j.cmd.Wait()
}

// Output returns the task's output
func (j *ScaleImageTask) Output() ([]byte, error) {
	return j.cmd.Output()
}

// StdoutPipe returns the task's StdoutPipe
func (j *ScaleImageTask) StdoutPipe() (io.ReadCloser, error) {
	return j.cmd.StdoutPipe()
}

// Promise runs the task asynchronously and returns a channel
// that will emit the image scaling task's status when it completes
func (j *ScaleImageTask) Promise() chan error {
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
