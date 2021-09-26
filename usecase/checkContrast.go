package usecase

import (
	"fmt"
	"os"
	"time"

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
	fmt.Println("start")
	for {
		if ok := video.Read(&frame); !ok {
			break
		}
		if frame.Empty() {
			break
		}

		wd, err := os.Getwd()
		if err != nil {
			break
		}

		// frame.ToBytes()だとgosseractで読み込めないので一旦書き出す
		//TODO: どう考えてもIO処理の無駄なので書き出さずに読み出せないか
		filePath := wd + "/images/" + time.Now().String() + ".png"
		fmt.Println(filePath)
		ok := gocv.IMWrite(filePath, frame)
		if !ok {
			break
		}

		err = client.SetImage(filePath)
		if err != nil {
			return fmt.Errorf("image set error: %w", err)
		}
		client.SetLanguage([]string{"eng"}...)

		outs, err := client.GetBoundingBoxes(gosseract.RIL_WORD)
		if err != nil {
			break
		}
		for _, out := range outs {
			fmt.Println("out: ", out)
			//TODO: Contrastチェックをする
			imgImage, err := frame.ToImage()
			if err != nil {
				return fmt.Errorf("image to image error: %w", err)
			}
			frame.Channels()
			// 文字の中心の色を取る
			//TODO: ○ みたいな場合(文字じゃないが) 色を正確に取れないので別の方法も組み合わせる
			textColor := imgImage.At((out.Box.Min.X+out.Box.Max.X)/2, (out.Box.Min.Y+out.Box.Max.Y)/2)
			tr, tg, tb, ta := textColor.RGBA()

			// 右から8bit取られるのでおそらく大丈夫なはず
			trgba, err := colors.RGBA(uint8(tr), uint8(tg), uint8(tb), float64(uint8(ta)/255))
			if err != nil {
				fmt.Println(err)
				continue
			}
			trgb := trgba.ToRGB()
			// とりあえず文字の外側-10くらいを取る
			//TODO: 画像外参照する可能性もあるので改善する
			if out.Box.Min.X < 10 || out.Box.Min.Y < 10 {
				fmt.Println("out of the range: ", out.Box.Min.X, out.Box.Min.Y)
				continue
			}
			backgroundColor := imgImage.At(out.Box.Min.X-10, out.Box.Min.Y-10)
			br, bg, bb, ba := backgroundColor.RGBA()
			//fmt.Println(br, bg, bb, ba)
			brgba, err := colors.RGBA(uint8(br), uint8(bg), uint8(bb), float64(uint8(ba)/255))
			brgb := brgba.ToRGB()
			ratio := contrast.GetRatio(contrast.NewRGBColor(trgb.R, trgb.G, trgb.B), contrast.NewRGBColor(brgb.R, brgb.G, brgb.B))
			fmt.Println("ratio: ", ratio)
			fmt.Println("level: ", contrast.GetLevel(ratio))

		}

		//TODO: 何かしらで結果を表示する
		window.IMShow(frame)
		if window.WaitKey(1) >= 0 {
			break
		}
		os.Remove(filePath)
	}

	fmt.Println("finished")
	return nil
}
