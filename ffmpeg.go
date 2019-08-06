package main

const ffmpegCmd = "ffmpeg"

type ffmpegTask interface {
	// Start begins the task. Start() may be called internally by Run().
	// If called from the outside, then wait for status using Promise().
	Start() error

	// Run begins the task (via Start()) and waits for completion.
	Run() error

	// Promise runs the task asynchronously and returns a channel
	// that will emit thetask's status upon completion.
	Promise() chan error
}
