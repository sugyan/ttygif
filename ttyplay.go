package ttygif

import (
	"io"
	"log"
	"os"
	"time"
)

// TtyPlayer interface
type TtyPlayer interface {
	GetPlayFunc() func(*TtyData)
}

// TtyPlayWithWait is the impementation of TtyPlayer
var TtyPlayWithWait ttyPlayWithWait

type ttyPlayWithWait struct{}

func (ttyPlayWithWait) GetPlayFunc() func(*TtyData) {
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
			time.Sleep(time.Microsecond * time.Duration(diff.Sec*1000000+diff.Usec))
		}
		prevTv = data.TimeVal
		print(string(*data.Buffer))
	}
}

// TtyPlay plays
func TtyPlay(filename string, player TtyPlayer) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	process := player.GetPlayFunc()
	reader := NewTtyReader(file)
	for {
		data, err := reader.ReadData()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				log.Fatal(err)
			}
		}
		process(data)
	}
}
