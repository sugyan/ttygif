package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/sugyan/ttygif"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

var version = "0.0.2"

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	input := flag.String("in", "ttyrecord", "input ttyrec file")
	output := flag.String("out", "tty.gif", "output gif file")
	speed := flag.Float64("s", 1.0, "play speed")
	help := flag.Bool("help", false, "usage")
	v := flag.Bool("v", false, "version")
	flag.Parse()
	if *help {
		flag.Usage()
		os.Exit(0)
	}
	if *v {
		fmt.Println(version)
		os.Exit(0)
	}

	err := validateInputFile(*input)
	if err != nil {
		fmt.Fprintln(os.Stderr, "input error:", err)
		flag.Usage()
		os.Exit(1)
	}

	generator := ttygif.NewGifGenerator()
	generator.Speed(*speed)
	err = generator.Generate(*input, *output)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	absPath, err := filepath.Abs(*output)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Printf("%s created!\n", absPath)
}

func validateInputFile(filename string) (err error) {
	file, err := os.Open(filename)
	if err != nil {
		return
	}
	defer file.Close()

	var timestamp int32
	now := time.Now()
	reader := ttygif.NewTtyReader(file)
	for {
		var data *ttygif.TtyData
		data, err = reader.ReadData()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return
			}
		}
		if data.TimeVal.Sec > int32(now.Unix()) || data.TimeVal.Sec < timestamp {
			return errors.New("invalid file")
		}
		timestamp = data.TimeVal.Sec
	}
	return nil
}
