package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	if len(os.Args) < 2 {
		fmt.Printf("usage: %s <filename>\n", os.Args[0])
		os.Exit(1)
	}
	filename := os.Args[1]
	log.Println(filename)

	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	play(file)
}

func play(file *os.File) {
	var (
		prevTv TimeVal
		first  = true
	)
	ch := TtyRead(file)
	for data := range ch {
		if first {
			first = false
		} else {
			diff := data.Header.tv.Subtract(&prevTv)
			<-time.After(time.Microsecond * time.Duration(diff.sec*1000000+diff.usec))
		}
		prevTv = data.Header.tv

		print(string(*data.Buffer))
	}
}
