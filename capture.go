package ttygif

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// CaptureImage take a screen shot of terminal
// TODO: Linux
func CaptureImage(dir string, data *TtyData) (filename string, err error) {
	// TODO: Terminal/iTerm only?
	var program string
	switch os.Getenv("TERM_PROGRAM") {
	case "iTerm.app":
		program = "iTerm"
	case "Apple_Terminal":
		program = "Terminal"
	default:
		return "", fmt.Errorf("Can't get TERM_PROGRAM")
	}
	// get window id
	windowID, err := exec.Command("osascript", "-e",
		fmt.Sprintf("tell app \"%s\" to id of window 1", program),
	).Output()
	if err != nil {
		return
	}
	// get screen capture
	// TODO: resize image if high resolution (retina display)
	filename = filepath.Join(dir, fmt.Sprintf("%d_%d.png", data.TimeVal.Sec, data.TimeVal.Usec))
	if err = exec.Command("screencapture", "-l", string(windowID), "-o", "-m", "-t", "png", filename).Run(); err != nil {
		return
	}

	return filename, nil
}
