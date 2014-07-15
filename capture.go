package ttygif

import (
	"fmt"
	"github.com/sugyan/ttygif/image/xwd"
	"image"
	"image/png"
	"io"
	"os"
	"os/exec"
)

// CapturedImage type
type CapturedImage struct {
	path    string
	decoder func(r io.Reader) (image.Image, error)
}

// CaptureImage take a screen shot of terminal
// TODO: Terminal/iTerm or X-Window only?
func CaptureImage(path string) (result *CapturedImage, err error) {
	switch os.Getenv("WINDOWID") {
	case "":
		return captureByScreencapture(path)
	default:
		return captureByXwd(path)
	}
}

// func captureByScreencapture(dir string, filename string) (img image.Image, err error) {
func captureByScreencapture(path string) (result *CapturedImage, err error) {
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
	if err = exec.Command("screencapture", "-l", string(windowID), "-o", "-m", "-t", "png", path).Run(); err != nil {
		return
	}
	return &CapturedImage{
		path:    path,
		decoder: png.Decode,
	}, nil
}

func captureByXwd(path string) (result *CapturedImage, err error) {
	if err = exec.Command("xwd", "-id", os.Getenv("WINDOWID"), "-out", path).Run(); err != nil {
		return
	}
	return &CapturedImage{
		path:    path,
		decoder: xwd.Decode,
	}, nil
}
