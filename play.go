package main

import (
	"github.com/sugyan/ttyread"
	"io"
	"os"
	"os/exec"
)

// Play reads ttyrecord file and play (write to STDOUT)
func Play(filename string, capture func(ttyread.TimeVal) error) (err error) {
	file, err := os.Open(filename)
	if err != nil {
		return
	}
	defer file.Close()
	defer clearScreen()

	var prevTv *ttyread.TimeVal
	reader := ttyread.NewTtyReader(file)
	for {
		// read
		var data *ttyread.TtyData
		data, err = reader.ReadData()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return
			}
		}
		// calc delay
		var diff ttyread.TimeVal
		if prevTv == nil {
			err = clearScreen()
			if err != nil {
				return
			}
		} else {
			diff = data.TimeVal.Subtract(*prevTv)
		}
		prevTv = &data.TimeVal
		// capture
		err = capture(diff)
		if err != nil {
			return
		}
		// play
		_, err = os.Stdout.Write(*data.Buffer)
		if err != nil {
			return
		}
	}
	return nil
}

func clearScreen() (err error) {
	bytes, err := exec.Command("clear").Output()
	if err != nil {
		return
	}
	os.Stdout.Write(bytes)
	return nil
}
