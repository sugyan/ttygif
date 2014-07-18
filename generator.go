package ttygif

import (
	"fmt"
	"image/gif"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

// TtyPlayCaptureProcessor type, which is the implementation of TtyPlayProcessor
type TtyPlayCaptureProcessor struct {
	timestamp TimeVal
	delays    []int
	tempDir   string
	worker    *Worker
	speed     float64
}

// Process captures and append images
func (t *TtyPlayCaptureProcessor) Process(diff TimeVal) (err error) {
	delay := int(float64(diff.Sec*1000000+diff.Usec)/t.speed) / 10000
	if delay == 0 {
		return nil
	}
	t.delays = append(t.delays, delay)
	t.timestamp = t.timestamp.Add(diff)

	// capture and append to images
	imgPath := filepath.Join(t.tempDir, fmt.Sprintf("%03d_%06d", t.timestamp.Sec, t.timestamp.Usec))
	fileType, err := CaptureImage(imgPath)
	if err != nil {
		return
	}
	t.worker.AddTargetFile(imgPath, fileType)
	return nil
}

// NewTtyPlayCaptureProcessor returns TtyPlayCaptureProcessor instance
func NewTtyPlayCaptureProcessor(speed float64) (t *TtyPlayCaptureProcessor, err error) {
	tempDir, err := ioutil.TempDir("", "ttygif")
	if err != nil {
		return
	}
	return &TtyPlayCaptureProcessor{
		tempDir: tempDir,
		speed:   speed,
		worker:  NewWorker(),
	}, nil
}

// RemoveTempDirectory remove tempDir
func (t *TtyPlayCaptureProcessor) RemoveTempDirectory() (err error) {
	if err = os.RemoveAll(t.tempDir); err != nil {
		return
	}
	return nil
}

// GifGenerator type
type GifGenerator struct {
	speed float64
}

// NewGifGenerator returns GifGenerator instance
func NewGifGenerator() *GifGenerator {
	return &GifGenerator{speed: 1.0}
}

// Speed sets the speed
func (g *GifGenerator) Speed(speed float64) {
	g.speed = speed
}

// Generate writes to outFile an animated GIF
func (g *GifGenerator) Generate(inFile string, outFile string) (err error) {
	capturer, err := NewTtyPlayCaptureProcessor(g.speed)
	if err != nil {
		return err
	}
	defer capturer.RemoveTempDirectory()
	// play and capture
	player := NewTtyPlayer()
	player.Processor(capturer)
	err = player.Play(inFile)
	if err != nil {
		return
	}
	// get paletted images from capture fiels
	progress := make(chan struct{})
	go func() {
	Loop:
		for {
			select {
			case _, ok := <-progress:
				if !ok {
					break Loop
				}
				print(".")
			case <-time.After(time.Second):
				print(".")
			}
		}
		print("\r")
	}()
	images, err := capturer.worker.GetAllImages(progress)
	if err != nil {
		return
	}
	// generate GIF file
	file, err := os.Create(outFile)
	if err != nil {
		return
	}
	defer file.Close()
	err = gif.EncodeAll(file, &gif.GIF{
		Image:     images,
		Delay:     capturer.delays,
		LoopCount: -1,
	})
	if err != nil {
		return
	}
	return nil
}
