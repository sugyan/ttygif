package ttygif

import (
	"fmt"
	"os"
	"os/exec"
)

// CaptureImage take a screen shot of terminal
// TODO: Terminal/iTerm or X-Window only?
func CaptureImage(path string) (string, error) {
	switch os.Getenv("WINDOWID") {
	case "":
		return captureByScreencapture(path)
	default:
		return captureByXwd(path)
	}
}

// func captureByScreencapture(dir string, filename string) (img image.Image, err error) {
func captureByScreencapture(path string) (fileType string, err error) {
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
	err = exec.Command("screencapture", "-l", string(windowID), "-o", "-m", "-t", "png", path).Run()
	if err != nil {
		return
	}
	return "png", nil
}

func captureByXwd(path string) (fileType string, err error) {
	err = exec.Command("xwd", "-id", os.Getenv("WINDOWID"), "-out", path).Run()
	if err != nil {
		return
	}
	return "xwd", nil
}
