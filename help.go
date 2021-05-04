package main

import (
	"fmt"
	"gocv.io/x/gocv"
	"image"
	"io"
	"log"
)

func Help() {
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
		thumb := GetDefectCount(img, c, 80, 100)
		RobotMovement(fingers, thumb)

		if ImShow(img, wt) {
			break
		}
	}
}
