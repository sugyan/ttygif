package ttygif

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"os/exec"
	"path/filepath"
)

// CaptureImage take a screen shot of terminal
// TODO: Linux
func CaptureImage(dir string, data *TtyData) (img image.Image, err error) {
	// TODO: Terminal/iTerm only?
	var program string
	switch os.Getenv("TERM_PROGRAM") {
	case "iTerm.app":
		program = "iTerm"
	case "Apple_Terminal":
		program = "Terminal"
	default:
		return nil, fmt.Errorf("Can't get TERM_PROGRAM")
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
	filename := filepath.Join(dir, fmt.Sprintf("%d_%d.png", data.TimeVal.Sec, data.TimeVal.Usec))
	if err = exec.Command("screencapture", "-l", string(windowID), "-o", "-m", "-t", "png", filename).Run(); err != nil {
		return
	}

	// open saved capture image
	file, err := os.Open(filename)
	if err != nil {
		return
	}
	defer file.Close()

	img, err = png.Decode(file)
	if err != nil {
		return
	}
	return img, nil
}
