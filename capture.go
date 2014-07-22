package ttygif

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
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
	err = exec.Command("screencapture", "-l", string(windowID), "-o", "-m", "-t", "png", path).Run()
	if err != nil {
		return
	}
	// resize image if high resolution (retina display)
	getProperty := func(key string) (result float64, err error) {
		sips := exec.Command("sips", "-g", key, path)
		awk := exec.Command("awk", "/:/ { print $2 }")
		sipsOut, err := sips.StdoutPipe()
		if err != nil {
			return
		}
		awk.Stdin = sipsOut
		sips.Start()
		output, err := awk.Output()
		if err != nil {
			return
		}
		err = sips.Wait()
		if err != nil {
			return
		}
		str := string(output)
		result, err = strconv.ParseFloat(str[:len(str)-1], 32)
		if err != nil {
			return
		}
		return result, nil
	}
	properties, err := func() (results map[string]float64, err error) {
		results = make(map[string]float64)
		for _, key := range []string{"pixelHeight", "pixelWidth", "dpiHeight", "dpiWidth"} {
			var property float64
			property, err = getProperty(key)
			if err != nil {
				return
			}
			results[key] = property
		}
		return results, nil
	}()
	if err != nil {
		return
	}
	if properties["dpiHeight"] > 72.0 && properties["dpiWidth"] > 72.0 {
		pixelHeight := int(properties["pixelHeight"] * 72.0 / properties["dpiHeight"])
		pixelWidth := int(properties["pixelWidth"] * 72.0 / properties["dpiWidth"])
		err = exec.Command("sips",
			"-s", "dpiWidth", "72.0", "-s", "dpiHeight", "72.0",
			"-z", strconv.Itoa(pixelHeight), strconv.Itoa(pixelWidth),
			path,
		).Run()
		if err != nil {
			return
		}
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
