package ttygif

import (
	"image"
	"image/color/palette"
	"image/gif"
	"io/ioutil"
	"os"
)

// TtyPlayCaptureProcessor type, which is the implementation of TtyPlayProcessor
type TtyPlayCaptureProcessor struct {
	timestamp TimeVal
	images    []*image.Paletted
	delays    []int
	tempDir   string
}

// Process captures and append images
func (t *TtyPlayCaptureProcessor) Process(diff TimeVal) (err error) {
	delay := int((diff.Sec*1000000 + diff.Usec) / 10000)
	if delay == 0 {
		return nil
	}
	t.delays = append(t.delays, delay)
	t.timestamp = t.timestamp.Add(diff)

	// capture and append to images
	img, err := CaptureImage(t.tempDir, t.timestamp)
	if err != nil {
		return
	}
	paletted := image.NewPaletted(img.Bounds(), palette.WebSafe)
	for x := paletted.Rect.Min.X; x < paletted.Rect.Max.X; x++ {
		for y := paletted.Rect.Min.Y; y < paletted.Rect.Max.Y; y++ {
			paletted.Set(x, y, img.At(x, y))
		}
	}
	t.images = append(t.images, paletted)
	return nil
}

// NewTtyPlayCaptureProcessor returns TtyPlayCaptureProcessor instance
func NewTtyPlayCaptureProcessor() (t *TtyPlayCaptureProcessor, err error) {
	tempDir, err := ioutil.TempDir("", "ttygif")
	if err != nil {
		return
	}
	return &TtyPlayCaptureProcessor{
		tempDir: tempDir,
	}, nil
}

// RemoveTempDirectory remove tempDir
func (t *TtyPlayCaptureProcessor) RemoveTempDirectory() (err error) {
	if err = os.RemoveAll(t.tempDir); err != nil {
		return
	}
	return nil
}

// GenerateGif creates gif animation
func GenerateGif(filename string) (err error) {
	capturer, err := NewTtyPlayCaptureProcessor()
	if err != nil {
		return err
	}
	defer capturer.RemoveTempDirectory()

	player := NewTtyPlayer()
	player.Processor(capturer)
	player.Play(filename)

	outFile, err := os.Create("out.gif")
	if err != nil {
		return
	}
	defer outFile.Close()

	err = gif.EncodeAll(outFile, &gif.GIF{
		Image:     capturer.images,
		Delay:     capturer.delays,
		LoopCount: -1,
	})
	if err != nil {
		return
	}
	return nil
}
