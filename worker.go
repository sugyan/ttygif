package ttygif

import (
	"errors"
	"github.com/sugyan/ttygif/image/xwd"
	"image"
	"image/color/palette"
	"image/png"
	"io"
	"os"
	"sync"
)

// WorkerInput type
type WorkerInput struct {
	index    int
	filePath string
	fileType string
}

// WorkerOutput type
type WorkerOutput struct {
	index    int
	paletted *image.Paletted
	err      error
}

// Worker type
type Worker struct {
	inputs []WorkerInput
}

// NewWorker returns Worker instance
func NewWorker() *Worker {
	return &Worker{}
}

// AddTargetFile adds input
func (w *Worker) AddTargetFile(filePath string, fileType string) {
	index := len(w.inputs)
	w.inputs = append(w.inputs, WorkerInput{
		index:    index,
		filePath: filePath,
		fileType: fileType,
	})
}

// GetAllImages waits and returns all images
func (w *Worker) GetAllImages(progress chan<- struct{}) ([]*image.Paletted, error) {
	done := make(chan struct{})
	defer func() {
		close(done)
		close(progress)
	}()
	inputs, errc := w.getInputChannel(done)
	output := make(chan *WorkerOutput)

	var (
		wg         sync.WaitGroup
		numWorkers = 10
	)
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			worker(inputs, output, done)
			wg.Done()
		}()
	}
	go func() {
		wg.Wait()
		close(output)
	}()

	results := make([]*image.Paletted, len(w.inputs))
Loop:
	for {
		select {
		case output, ok := <-output:
			if !ok {
				break Loop
			}
			if output.err != nil {
				return nil, output.err
			}
			results[output.index] = output.paletted
			progress <- struct{}{}
		case err := <-errc:
			if err != nil {
				return nil, err
			}
		}
	}
	return results, nil
}

func (w *Worker) getInputChannel(done <-chan struct{}) (<-chan WorkerInput, <-chan error) {
	inputs := make(chan WorkerInput)
	errc := make(chan error, 1)
	go func() {
		defer close(inputs)
		errc <- func(walkFunc func(WorkerInput) error) error {
			for _, input := range w.inputs {
				err := walkFunc(input)
				if err != nil {
					return err
				}
			}
			return nil
		}(func(input WorkerInput) error {
			select {
			case inputs <- input:
			case <-done:
				return errors.New("Canceled")
			}
			return nil
		})
	}()
	return inputs, errc
}

func worker(inputs <-chan WorkerInput, output chan<- *WorkerOutput, done <-chan struct{}) {
	for input := range inputs {
		paletted, err := decode(input.filePath, input.fileType)
		select {
		case output <- &WorkerOutput{index: input.index, paletted: paletted, err: err}:
		case <-done:
			return
		}
	}
}

func decode(filePath string, fileType string) (paletted *image.Paletted, err error) {
	var decoder func(io.Reader) (image.Image, error)
	switch fileType {
	case "png":
		decoder = png.Decode
	case "xwd":
		decoder = xwd.Decode
	default:
		return nil, errors.New("Unsupported file type")
	}
	// open file
	var file *os.File
	file, err = os.Open(filePath)
	if err != nil {
		return
	}
	defer file.Close()
	// decode
	img, err := decoder(file)
	if err != nil {
		return
	}
	paletted = image.NewPaletted(img.Bounds(), palette.WebSafe)
	for x := paletted.Rect.Min.X; x < paletted.Rect.Max.X; x++ {
		for y := paletted.Rect.Min.Y; y < paletted.Rect.Max.Y; y++ {
			paletted.Set(x, y, img.At(x, y))
		}
	}
	return paletted, nil
}
