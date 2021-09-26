package main

import (
	"fmt"
	"os"

	"github.com/masibw/violint/usecase"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"gocv.io/x/gocv"
	"log"
)

func main() {

	app := cli.NewApp()
	app.Name = "violint"
	app.Usage = "violint [filename]"

	app.Action = func(c *cli.Context) error {
		if c.NArg() < 1 {
			fmt.Println("Please specify a video file: violint [filename]")
			os.Exit(1)
		}
		file := c.Args().Get(0)

		video, err := gocv.VideoCaptureFile(file)
		if err != nil {
			return errors.Wrapf(err, "Error opening video capture file: %s", file)
		}
		defer video.Close()

		err = usecase.CheckContrast(video)
		if err != nil {
			return errors.Wrapf(err, "file : %s", file)
		}
		return nil
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatalf("failed to run app: %v", err)
	}
}
