// this is a "scratch" file. I'm using this space to test and edit
// functions without messing up what's already working

package main

import (
	"fmt"
	"gocv.io/x/gocv"
	"image"
	"image/color"
	"math"
)

func Help() {

	deviceID := 0
	size := image.Point{X: 600, Y: 600}
	blur := image.Point{X: 11, Y: 11}
	wt := gocv.NewWindow("Just the Hand")
	defer wt.Close()
	img := gocv.NewMat()
	defer img.Close()

	wt.ResizeWindow(1400, 1400)
	wt.MoveWindow(0, 0)

	video, _ := gocv.OpenVideoCapture(deviceID)
	defer video.Close()

	for {
		if !video.Read(&img) {
			break
		}

		// cleaning up the image
		gocv.Flip(img, &img, 1)
		gocv.Resize(img, &img, size, 0, 0, gocv.InterpolationLinear)
		gocv.GaussianBlur(img, &frame, blur, 0, 0, gocv.BorderReflect101)
		gocv.CvtColor(frame, &hsv, gocv.ColorBGRToHSV)
		gocv.InRangeWithScalar(hsv, lhsv, hhsv, &mask)
		gocv.Erode(mask, &mask, kernel)
		gocv.Dilate(mask, &mask, kernel)

		/////////////////////////////
		// hand detection stuff
		contours := gocv.FindContours(mask, gocv.RetrievalExternal, gocv.ChainApproxSimple)
		c := GetBiggestContour(contours)
		gocv.ConvexHull(c, &hull, true, false)
		gocv.ConvexityDefects(c, hull, &defects)

		var angle float64
		defectCount := 0
		for i := 0; i < defects.Rows(); i++ {
			start := c.At(int(defects.GetIntAt(i, 0)))
			end := c.At(int(defects.GetIntAt(i, 1)))
			far := c.At(int(defects.GetIntAt(i, 2)))

			a := math.Sqrt(math.Pow(float64(end.X-start.X), 2) + math.Pow(float64(end.Y-start.Y), 2))
			b := math.Sqrt(math.Pow(float64(far.X-start.X), 2) + math.Pow(float64(far.Y-start.Y), 2))
			c := math.Sqrt(math.Pow(float64(end.X-far.X), 2) + math.Pow(float64(end.Y-far.Y), 2))

			// apply cosine rule here
			angle = math.Acos((math.Pow(b, 2)+math.Pow(c, 2)-math.Pow(a, 2))/(2*b*c)) * 57

			// ignore angles > 90 and highlight rest with dots
			if angle <= 65 {
				defectCount++
				gocv.Circle(&img, far, 1, red, 2)
			}
		}
		status := fmt.Sprintf("defectCount: %d", defectCount+1)

		rect := gocv.BoundingRect(c)
		gocv.Rectangle(&img, rect, color.RGBA{R: 255, G: 255, B: 255}, 2)

		gocv.PutText(&img, status, image.Pt(10, 20), gocv.FontHersheyPlain, 1.2, red, 2)
		/////////////////////////////////

		if ImShow(img, wt) {
			break
		}
	}
}
