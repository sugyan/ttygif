package ttygif

import (
	"fmt"
	"github.com/sugyan/ttygif/image/xwd"
	"image"
	"image/png"
	"os"
	"os/exec"
	"path/filepath"
)

// CaptureImage take a screen shot of terminal
// TODO: Terminal/iTerm or X-Window only?
func CaptureImage(dir string, filename string) (image.Image, error) {
	switch os.Getenv("WINDOWID") {
	case "":
		return captureByScreencapture(dir, filename)
	default:
		return captureByXwd(dir, filename)
	}
}

func captureByScreencapture(dir string, filename string) (img image.Image, err error) {
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
	path := filepath.Join(dir, filename)
	if err = exec.Command("screencapture", "-l", string(windowID), "-o", "-m", "-t", "png", path).Run(); err != nil {
		return
	}

	// open saved capture image
	file, err := os.Open(path)
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

func captureByXwd(dir string, filename string) (img image.Image, err error) {
	path := filepath.Join(dir, filename)
	if err = exec.Command("xwd", "-id", os.Getenv("WINDOWID"), "-out", path).Run(); err != nil {
		return
	}
	// open saved capture image
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	img, err = xwd.Decode(file)
	if err != nil {
		return
	}
	return img, nil
}
