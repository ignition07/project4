// This pulls up a video that should only pick up on the color of the glove

package main

import (
	"gocv.io/x/gocv"
	"image"
)

func Thresholded() {

	size := image.Point{X: 600, Y: 600}
	blur := image.Point{X: 11, Y: 11}
	wt := gocv.NewWindow("Just the Hand")
	img := gocv.NewMat()

	wt.ResizeWindow(1400, 1400)
	wt.MoveWindow(0, 0)
	video, _ := gocv.OpenVideoCapture(0)
	defer video.Close()

	for {
		if !video.Read(&img) {
			break
		}

		gocv.Flip(img, &img, 1)
		gocv.Resize(img, &img, size, 0, 0, gocv.InterpolationLinear)
		gocv.GaussianBlur(img, &frame, blur, 0, 0, gocv.BorderReflect101)
		gocv.CvtColor(frame, &hsv, gocv.ColorBGRToHSV)
		gocv.InRangeWithScalar(hsv, lhsv, hhsv, &mask)
		gocv.Erode(mask, &mask, kernel)
		gocv.Dilate(mask, &mask, kernel)

		if ImShow(mask, wt) {
			break
		}
	}
}
