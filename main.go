package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("usage: %s <ttyrec file>\n", os.Args[0])
		os.Exit(1)
	}
	GenerateGif(os.Args[1])
}

func play(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var (
		first  = true
		prevTv TimeVal
	)
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
		if first {
			first = false
		} else {
			diff := data.Header.tv.Subtract(prevTv)
			time.Sleep(time.Microsecond * time.Duration(diff.sec*1000000+diff.usec))
		}
		prevTv = data.Header.tv
		print(string(*data.Buffer))
	}
}
