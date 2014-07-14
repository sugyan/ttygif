package main

import (
	"fmt"
	"github.com/sugyan/ttygif"
	"log"
	"os"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("usage: %s <ttyrec file>\n", os.Args[0])
		os.Exit(1)
	}
	ttygif.GenerateGif(os.Args[1])
}
