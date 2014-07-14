package ttygif

import (
	"image"
	"image/color/palette"
	"image/gif"
	"io/ioutil"
	"log"
	"os"
)

// TtyPlayWithCapture type, which is the implementation of TtyPlayer
type TtyPlayWithCapture struct {
	images  []*image.Paletted
	delays  []int
	tempDir string
}

// GetPlayFunc returns Player func
func (t *TtyPlayWithCapture) GetPlayFunc() func(*TtyData) {
	var (
		first  = true
		prevTv TimeVal
	)
	return func(data *TtyData) {
		if first {
			print("\x1b[1;1H\x1b[2J")
			first = false
		} else {
			diff := data.TimeVal.Subtract(prevTv)
			t.delays = append(t.delays, int((diff.Sec*1000000+diff.Usec)/10000))
		}
		print(string(*data.Buffer))
		prevTv = data.TimeVal

		// capture and append to images
		img, err := CaptureImage(t.tempDir, data)
		if err != nil {
			log.Fatal(err)
		}
		paletted := image.NewPaletted(img.Bounds(), palette.WebSafe)
		for x := paletted.Rect.Min.X; x < paletted.Rect.Max.X; x++ {
			for y := paletted.Rect.Min.Y; y < paletted.Rect.Max.Y; y++ {
				paletted.Set(x, y, img.At(x, y))
			}
		}
		t.images = append(t.images, paletted)
	}
}

// GenerateGif creates gif animation
func GenerateGif(filename string) {
	tempDir, err := ioutil.TempDir("", "ttygif")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			log.Fatal(err)
		}
	}()

	capture := &TtyPlayWithCapture{
		tempDir: tempDir,
	}
	TtyPlay(filename, capture)

	// last delay
	if len(capture.images) > len(capture.delays) {
		capture.delays = append(capture.delays, 200)
	}

	outFile, err := os.Create("out.gif")
	if err != nil {
		log.Fatal(err)
	}
	defer outFile.Close()

	if err = gif.EncodeAll(outFile, &gif.GIF{
		Image:     capture.images,
		Delay:     capture.delays,
		LoopCount: -1,
	}); err != nil {
		log.Fatal(err)
	}
}
