// This program will help you get the hsv values for the glove

package main

import (
	"gocv.io/x/gocv"
)

func GetHand() {
	wi := gocv.NewWindow("normal")
	wt := gocv.NewWindow("threshold")
	wt.ResizeWindow(1400, 1400)
	wt.MoveWindow(0, 0)
	wi.MoveWindow(1400, 0)
	wi.ResizeWindow(1400, 1400)

	lh := wi.CreateTrackbar("Low H", 360/2)
	hh := wi.CreateTrackbar("High H", 255)
	ls := wi.CreateTrackbar("Low S", 255)
	hs := wi.CreateTrackbar("High S", 255)
	lv := wi.CreateTrackbar("Low V", 255)
	hv := wi.CreateTrackbar("High V", 255)

	video, _ := gocv.OpenVideoCapture(0)
	img := gocv.NewMat()

	for {
		video.Read(&img)
		gocv.CvtColor(img, &img, gocv.ColorBGRToHSV)

		thresholded := gocv.NewMat()
		gocv.InRangeWithScalar(img,
			gocv.Scalar{Val1: GetPosFloat(lh), Val2: GetPosFloat(ls), Val3: GetPosFloat(lv)},
			gocv.Scalar{Val1: GetPosFloat(hh), Val2: GetPosFloat(hs), Val3: GetPosFloat(hv)},
			&thresholded)

		if ImShow(img, wi) || ImShow(thresholded, wt) {
			break
		}

	}
}
