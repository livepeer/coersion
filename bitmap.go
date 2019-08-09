package main

import (
	"fmt"
	"os/exec"
	"time"
)

// BitmapTask uses ffmpeg to perform the scaling of a still image from a video source
type BitmapTask struct {
	source string
	cmd    *exec.Cmd
}

// NewBitmapTask creates a new image Scaleion task
func NewBitmapTask(source string, w, h uint, o time.Duration) *BitmapTask {
	// TO FILE:
	// return &BitmapTask{
	// 	source: source,
	// 	dest:   dest,
	// 	cmd: exec.Command(ffmpegCmd, []string{
	// 		"-i", source,
	// 		"-vf", fmt.Sprintf("scale=%d:%d", w, h),
	// 		fmt.Sprintf("%s.jpg", dest),
	// 	}...),
	// }
	offset := fmtDuration(o)
	return &BitmapTask{
		source: source,
		cmd: exec.Command(ffmpegCmd, []string{
			"-ss", offset,
			"-i", source,
			"-vframes", "1",
			"-vf", fmt.Sprintf("scale=%d:%d", w, h),
			"-pix_fmt", "bgr8",
			"-f", "image2pipe",
			"-vcodec", "bmp",
			"-q:v", "2",
			"-",
		}...),
	}
}

// Start begins the task
func (j *BitmapTask) Start() error {
	return j.cmd.Start()
}

// Run begins the task (via Start()) and waits for completion
func (j *BitmapTask) Run() error {
	if err := j.Start(); err != nil {
		return err
	}
	return j.cmd.Wait()
}

// Output returns the task's output
func (j *BitmapTask) Output() ([]byte, error) {
	return j.cmd.Output()
}

// Promise runs the task asynchronously and returns a channel
// that will emit the image scaling task's status when it completes
func (j *BitmapTask) Promise() chan error {
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
