package ttygif

import (
	"io"
	"os"
	"time"
)

// TtyPlayProcessor interface
type TtyPlayProcessor interface {
	Process(TimeVal) error
}

// TtyPlayWaitProcessor type
type TtyPlayWaitProcessor struct {
	Speed float32
}

// Process waits diff interval
func (t TtyPlayWaitProcessor) Process(diff TimeVal) error {
	time.Sleep(time.Microsecond * time.Duration(float32(diff.Sec*1000000+diff.Usec)/t.Speed))
	return nil
}

// TtyPlayer type
type TtyPlayer struct {
	processor TtyPlayProcessor
}

// NewTtyPlayer returns TtyPlayer instance
// Default TtyPlayProcessor is TtyPlayWaitProcessor.
func NewTtyPlayer() *TtyPlayer {
	return &TtyPlayer{
		processor: &TtyPlayWaitProcessor{
			Speed: 1.0,
		},
	}
}

// Processor sets the processor
func (player *TtyPlayer) Processor(processor TtyPlayProcessor) {
	player.processor = processor
}

// Play read ttyrec file and play
func (player *TtyPlayer) Play(filename string) (err error) {
	file, err := os.Open(filename)
	if err != nil {
		return
	}
	defer file.Close()
	defer clearScreen()

	var (
		first  = true
		prevTv TimeVal
	)
	reader := NewTtyReader(file)
	for {
		var data *TtyData
		data, err = reader.ReadData()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return
			}
		}

		var diff TimeVal
		if first {
			clearScreen()
			first = false
		} else {
			diff = data.TimeVal.Subtract(prevTv)
		}
		prevTv = data.TimeVal

		err = player.processor.Process(diff)
		if err != nil {
			return
		}
		print(string(*data.Buffer))
	}
	return nil
}

func clearScreen() {
	print("\x1b[1;1H\x1b[2J")
}
