package usecase

import (
	"fmt"
	"image"
	"image/color"

	"github.com/pkg/errors"

	"github.com/masibw/violint/pkg/contrast"
	"github.com/otiai10/gosseract/v2"
	"gopkg.in/go-playground/colors.v1"

	"gocv.io/x/gocv"
)

func CheckContrast(video *gocv.VideoCapture) error {
	window := gocv.NewWindow("Contrast")
	defer window.Close()

	frame := gocv.NewMat()
	defer frame.Close()

	client := gosseract.NewClient()
	defer client.Close()
	client.SetLanguage([]string{"eng"}...)

	fmt.Println("start")
	for {
		if ok := video.Read(&frame); !ok || frame.Empty() {
			break
		}

		buf, err := gocv.IMEncode(gocv.PNGFileExt, frame)
		if err != nil {
			return errors.Wrapf(err, "failed to Encode to png from %v", frame)
		}

		err = client.SetImageFromBytes(buf.GetBytes())
		if err != nil {
			return errors.Wrapf(err, "image set error")
		}

		outs, err := client.GetBoundingBoxes(gosseract.RIL_SYMBOL)
		if err != nil {
			return errors.Wrapf(err, "failed to get bounding boxes")
		}

		imgImage, err := frame.ToImage()
		if err != nil {
			return errors.Wrapf(err, "failed to convert Mat to image")
		}
		for _, out := range outs {
			fmt.Println("out: ", out)

			// Take the color of the center of the text
			//TODO: In case of something like n, you can't get the exact color, so you need to combine different methods.
			textColor := imgImage.At((out.Box.Min.X+out.Box.Max.X)/2, (out.Box.Min.Y+out.Box.Max.Y)/2)
			tr, tg, tb, ta := textColor.RGBA()

			// It takes 8 bits from the right, so it should probably be okay.
			trgba, err := colors.RGBA(uint8(tr), uint8(tg), uint8(tb), float64(uint8(ta)/255))
			if err != nil {
				return errors.Wrapf(err, "failed to convert colors.RGBAColor from %s", textColor)
			}
			textRGB := trgba.ToRGB()

			// For now, take about -10 outside the character.
			//TODO: There is a possibility of out-of-image referencing, which need to be improved.
			if out.Box.Min.X < 10 || out.Box.Min.Y < 10 {
				return errors.Wrapf(err, "the coordinate is out of the image")
			}
			backgroundColor := imgImage.At(out.Box.Min.X-10, out.Box.Min.Y-10)
			br, bg, bb, ba := backgroundColor.RGBA()

			brgba, err := colors.RGBA(uint8(br), uint8(bg), uint8(bb), float64(uint8(ba)/255))
			if err != nil {
				return errors.Wrapf(err, "failed to convert colors.RGBAColor from %s", backgroundColor)
			}
			backgroundRGB := brgba.ToRGB()
			ratio := contrast.GetRatio(textRGB, backgroundRGB)
			fmt.Println("ratio: ", ratio)
			level := contrast.GetLevel(ratio)
			fmt.Println("level: ", level)
			r := out.Box
			blue := color.RGBA{B: 255}
			gocv.Rectangle(&frame, r, blue, 3)
			if level == "-" || level == "A" {
				r := out.Box
				red := color.RGBA{R: 255}
				gocv.Rectangle(&frame, r, red, 3)
				size := gocv.GetTextSize("Low Contrast", gocv.FontHersheyPlain, 1.2, 2)
				pt := image.Pt((r.Min.X+r.Max.X)/2-(size.X/2), r.Min.Y-2)
				gocv.PutText(&frame, "Low Contrast", pt, gocv.FontHersheyPlain, 1.2, red, 2)
			}

		}

		//TODO: Display the results in some way
		window.IMShow(frame)
		if window.WaitKey(1) >= 0 {
			break
		}
	}

	fmt.Println("finished")
	return nil
}
