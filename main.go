package main

import (
	"fmt"
	"os"

	"github.com/masibw/violint/usecase"
	"gocv.io/x/gocv"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("How to run:\n\tviolint [filename]")
		return
	}
	// parse args
	file := os.Args[1]

	video, err := gocv.VideoCaptureFile(file)
	if err != nil {
		fmt.Printf("Error opening video capture file: %s\n", file)
		return
	}
	defer video.Close()

	err = usecase.CheckContrast(video)
	if err != nil {
		fmt.Printf("file : %s, Error: %s\n", file, err.Error())
		return
	}
}
