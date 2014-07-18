package main

import (
	"flag"
	"fmt"
	"github.com/sugyan/ttygif"
	"os"
	"path/filepath"
	"runtime"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	input := flag.String("in", "ttyrecord", "input ttyrec file")
	output := flag.String("out", "tty.gif", "output gif file")
	speed := flag.Float64("s", 1.0, "play speed")
	help := flag.Bool("help", false, "usage")
	flag.Parse()
	if *help {
		flag.Usage()
		os.Exit(0)
	}

	generator := ttygif.NewGifGenerator()
	generator.Speed(*speed)
	err := generator.Generate(*input, *output)
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
