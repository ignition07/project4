// This file includes the drone functions:

// GetHand()
// DetectBlueHand()

package main

import (
	"fmt"
	"gocv.io/x/gocv"
	"image"
	"io"
	"log"
)

// GetHand is used to get HSV values for glove with drone camera
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

	ffmpegOut := ConnectDrone()

	for {
		buf := make([]byte, frameSize)
		if _, err := io.ReadFull(ffmpegOut, buf); err != nil {
			fmt.Println(err)
			continue
		}

		img, err := gocv.NewMatFromBytes(720, 960, gocv.MatTypeCV8UC3, buf)
		if err != nil {
			log.Print(err)
			continue
		}
		if img.Empty() {
			continue
		}

		thresholded := gocv.NewMat()
		gocv.CvtColor(img, &img, gocv.ColorBGRToHSV)
		gocv.InRangeWithScalar(img,
			gocv.Scalar{Val1: GetPosFloat(lh), Val2: GetPosFloat(ls), Val3: GetPosFloat(lv)},
			gocv.Scalar{Val1: GetPosFloat(hh), Val2: GetPosFloat(hs), Val3: GetPosFloat(hv)},
			&thresholded)

		if ImShow(img, wi) || ImShow(thresholded, wt) {
			break
		}

	}
}

///////////////////////////////////////////////////////////////////////////////////////////

// DetectBlueHand detects glove using drone camera and ffmpeg
func DetectBlueHand() {
	size := image.Point{X: 600, Y: 600}
	blur := image.Point{X: 11, Y: 11}
	wt := gocv.NewWindow("Just the Hand")
	defer wt.Close()
	img := gocv.NewMat()
	defer img.Close()

	wt.ResizeWindow(1400, 1400)
	wt.MoveWindow(0, 0)

	ffmpegOut := ConnectDrone()

	for {
		buf := make([]byte, frameSize)
		if _, err := io.ReadFull(ffmpegOut, buf); err != nil {
			fmt.Println(err)
			continue
		}

		img, err := gocv.NewMatFromBytes(720, 960, gocv.MatTypeCV8UC3, buf)
		if err != nil {
			log.Print(err)
			continue
		}
		if img.Empty() {
			continue
		}

		// cleaning up the image
		gocv.Flip(img, &img, 1)
		gocv.Resize(img, &img, size, 0, 0, gocv.InterpolationLinear)
		gocv.GaussianBlur(img, &frame, blur, 0, 0, gocv.BorderReflect101)
		gocv.CvtColor(frame, &hsv, gocv.ColorBGRToHSV)
		gocv.InRangeWithScalar(hsv, lhsv, hhsv, &mask)
		gocv.Erode(mask, &mask, kernel)
		gocv.Dilate(mask, &mask, kernel)

		// hand detection stuff from gocv sample
		contours := gocv.FindContours(mask, gocv.RetrievalExternal, gocv.ChainApproxSimple)
		if contours.Size() <= 0 {
			ImShow(img, wt)
			continue
		}

		c := GetBiggestContour(contours)
		gocv.ConvexHull(c, &hull, true, false)
		gocv.ConvexityDefects(c, hull, &defects)
		fingers := GetDefectCount(img, c, 0, 70)
		fmt.Println(fingers)

		if ImShow(img, wt) {
			break
		}
	}
}
