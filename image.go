package main

import (
	"image"
	"image/color/palette"
	"image/gif"
	"image/png"
	"log"
	"os"
	"os/exec"
)

// GenerateGif creates gif animation
func GenerateGif(file *os.File) {

	var (
		prevTv  TimeVal
		first   = true
		tmpfile = "out.png"
		delays  []int
		images  []*image.Paletted
	)
	ch := TtyRead(file)
	for data := range ch {
		// play
		if first {
			first = false
		} else {
			diff := data.Header.tv.Subtract(&prevTv)
			delays = append(delays, int((diff.sec*1000000+diff.usec)/10000))
		}
		prevTv = data.Header.tv
		print(string(*data.Buffer))

		// screen capture (for Mac)
		windowID, err := exec.Command("osascript", "-e", "tell app \"iTerm\" to id of window 1").Output()
		if err != nil {
			log.Fatal(err)
		}
		if err := exec.Command("screencapture", "-l", string(windowID), "-o", "-m", tmpfile).Run(); err != nil {
			log.Fatal(err)
		}

		// read file
		file, err := os.Open(tmpfile)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		img, err := png.Decode(file)
		// add image
		p := image.NewPaletted(img.Bounds(), palette.WebSafe)
		for x := p.Rect.Min.X; x < p.Rect.Max.X; x++ {
			for y := p.Rect.Min.Y; y < p.Rect.Max.Y; y++ {
				p.Set(x, y, img.At(x, y))
			}
		}
		images = append(images, p)
	}
	if len(images) > len(delays) {
		delays = append(delays, 200)
	}

	outFile, err := os.Create("out.gif")
	if err != nil {
		log.Fatal(err)
	}
	defer outFile.Close()
	if err = gif.EncodeAll(outFile, &gif.GIF{
		Image:     images,
		Delay:     delays,
		LoopCount: -1,
	}); err != nil {
		log.Fatal(err)
	}
}
