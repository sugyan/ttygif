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

	generator := ttygif.NewGifGenerator()
	generator.Speed(1.0)
	err := generator.Generate(os.Args[1], "out.gif")
	if err != nil {
		log.Fatal(err)
	}
	/* play only */
	// player := ttygif.NewTtyPlayer()
	// player.Play(os.Args[1])
}
