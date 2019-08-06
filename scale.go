package main

import (
	"fmt"
	"os/exec"
)

// ScaleImageTask uses ffmpeg to perform the scaling of a still image from a video source
type ScaleImageTask struct {
	source string
	cmd    *exec.Cmd
}

// NewScaleImageTask creates a new image Scaleion task
func NewScaleImageTask(source string, w, h uint) *ScaleImageTask {
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
	return &ScaleImageTask{
		source: source,
		cmd: exec.Command(ffmpegCmd, []string{
			"-i", source,
			"-vf", fmt.Sprintf("scale=%d:%d", w, h),
			"-f", "image2pipe",
			"-vcodec", "png",
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
