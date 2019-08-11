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

func (j *BitmapTask) getCmd() *exec.Cmd {
	return j.cmd
}
