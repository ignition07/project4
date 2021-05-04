//// this file includes the functions:
//// GetHandLaptop()
//// DetectBlueHandLaptop()
//// DetectHand()
//
package main

//
//import (
//	"fmt"
//	"gocv.io/x/gocv"
//	"image"
//	"image/color"
//	"math"
//)
//
//// GetHandLaptop is used to get HSV values for glove using computer's camera
//func GetHandLaptop() {
//	wi := gocv.NewWindow("normal")
//	wt := gocv.NewWindow("threshold")
//	wt.ResizeWindow(1400, 1400)
//	wt.MoveWindow(0, 0)
//	wi.MoveWindow(1400, 0)
//	wi.ResizeWindow(1400, 1400)
//
//	lh := wi.CreateTrackbar("Low H", 360/2)
//	hh := wi.CreateTrackbar("High H", 255)
//	ls := wi.CreateTrackbar("Low S", 255)
//	hs := wi.CreateTrackbar("High S", 255)
//	lv := wi.CreateTrackbar("Low V", 255)
//	hv := wi.CreateTrackbar("High V", 255)
//
//	video, _ := gocv.OpenVideoCapture(0)
//	img := gocv.NewMat()
//
//	for {
//		thresholded := gocv.NewMat()
//		video.Read(&img)
//		gocv.CvtColor(img, &img, gocv.ColorBGRToHSV)
//
//		gocv.InRangeWithScalar(img,
//			gocv.Scalar{Val1: GetPosFloat(lh), Val2: GetPosFloat(ls), Val3: GetPosFloat(lv)},
//			gocv.Scalar{Val1: GetPosFloat(hh), Val2: GetPosFloat(hs), Val3: GetPosFloat(hv)},
//			&thresholded)
//
//		if ImShow(img, wi) || ImShow(thresholded, wt) {
//			break
//		}
//
//	}
//}
//
////////////////////////////////////////////////////////////////
//
//// DetectBlueHandLaptop detects glove using computer's camera
//func DetectBlueHandLaptop() {
//
//	deviceID := 0
//	size := image.Point{X: 600, Y: 600}
//	blur := image.Point{X: 11, Y: 11}
//	wt := gocv.NewWindow("Just the Hand")
//	img := gocv.NewMat()
//
//	wt.ResizeWindow(1400, 1400)
//	wt.MoveWindow(0, 0)
//
//	video, _ := gocv.OpenVideoCapture(deviceID)
//
//	for {
//		if !video.Read(&img) {
//			break
//		}
//
//		// cleaning up the image
//		gocv.Flip(img, &img, 1)
//		gocv.Resize(img, &img, size, 0, 0, gocv.InterpolationLinear)
//		gocv.GaussianBlur(img, &frame, blur, 0, 0, gocv.BorderReflect101)
//		gocv.CvtColor(frame, &hsv, gocv.ColorBGRToHSV)
//		gocv.InRangeWithScalar(hsv, lhsv, hhsv, &mask)
//		gocv.Erode(mask, &mask, kernel)
//		gocv.Dilate(mask, &mask, kernel)
//
//		/////////////////////////////
//		// hand detection stuff
//		contours := gocv.FindContours(mask, gocv.RetrievalExternal, gocv.ChainApproxSimple)
//		if contours.Size() <= 0 {
//			ImShow(img, wt)
//			continue
//		}
//
//		c := GetBiggestContour(contours)
//		gocv.ConvexHull(c, &hull, true, false)
//		gocv.ConvexityDefects(c, hull, &defects)
//
//		fingers := GetDefectCount(img, c, 0, 40) + 1
//		thumb := GetDefectCount(img, c, 41, 100)
//		fingerCount = append(fingerCount, fingers)
//		fingersMode := GetMode(fingerCount)
//		thumbCount = append(thumbCount, thumb)
//		thumbMode := GetMode(thumbCount)
//		fmt.Println("fingers:", fingers, "   thumbs:", thumb)
//		fmt.Println("fingersMode:", fingersMode, "   thumbsMode:", thumbMode)
//		fmt.Println("fingerCount size:", len(fingerCount), "   thumbCount size:", len(thumbCount))
//
//
//		if len(fingerCount) > 50 {
//			fingerCount = Trim(fingerCount)
//		}
//		if len(thumbCount) > 50 {
//			thumbCount = Trim(thumbCount)
//		}
//
//		instruction := GetInstruction(fingersMode, thumbMode)
//		DoInstruction(instruction)
//		if ImShow(img, wt) {
//			break
//		}
//	}
//}
//
////////////////////////////////////////////////////////////////
//
//// DetectHand sample from gocv
//func DetectHand() {
//
//	deviceID := 0
//	webcam, err := gocv.OpenVideoCapture(deviceID)
//	if err != nil {
//		fmt.Printf("Error opening video capture device: %v\n", deviceID)
//		return
//	}
//	defer webcam.Close()
//
//	window := gocv.NewWindow("Hand Gestures")
//	defer window.Close()
//
//	window.ResizeWindow(1400, 1400)
//
//	img := gocv.NewMat()
//	defer img.Close()
//
//	fmt.Printf("Start reading device: %v\n", deviceID)
//	for {
//		if ok := webcam.Read(&img); !ok {
//			fmt.Printf("Device closed: %v\n", deviceID)
//			return
//		}
//		if img.Empty() {
//			continue
//		}
//
//		gocv.CvtColor(img, &hsv, gocv.ColorBGRToHSV)
//		gocv.InRangeWithScalar(hsv, lhsv, hhsv, &mask)
//
//		// cleaning up image
//		gocv.Flip(img, &img, 1)
//		gocv.CvtColor(img, &imgGrey, gocv.ColorBGRToGray)
//		gocv.GaussianBlur(imgGrey, &imgBlur, image.Pt(35, 35), 0, 0, gocv.BorderDefault)
//		gocv.Threshold(imgBlur, &imgThresh, 0, 255, gocv.ThresholdBinaryInv+gocv.ThresholdOtsu)
//
//		// now find biggest contour
//		contours := gocv.FindContours(imgThresh, gocv.RetrievalExternal, gocv.ChainApproxSimple)
//		c := GetBiggestContour(contours)
//
//		gocv.ConvexHull(c, &hull, true, false)
//		gocv.ConvexityDefects(c, hull, &defects)
//
//		var angle float64
//		defectCount := 0
//		for i := 0; i < defects.Rows(); i++ {
//			start := c.At(int(defects.GetIntAt(i, 0)))
//			end := c.At(int(defects.GetIntAt(i, 1)))
//			far := c.At(int(defects.GetIntAt(i, 2)))
//
//			a := math.Sqrt(math.Pow(float64(end.X-start.X), 2) + math.Pow(float64(end.Y-start.Y), 2))
//			b := math.Sqrt(math.Pow(float64(far.X-start.X), 2) + math.Pow(float64(far.Y-start.Y), 2))
//			c := math.Sqrt(math.Pow(float64(end.X-far.X), 2) + math.Pow(float64(end.Y-far.Y), 2))
//
//			// apply cosine rule here
//			angle = math.Acos((math.Pow(b, 2)+math.Pow(c, 2)-math.Pow(a, 2))/(2*b*c)) * 57
//
//			// ignore angles > 90 and highlight rest with dots
//			if angle <= 50 {
//				defectCount++
//				gocv.Circle(&img, far, 1, green, 2)
//			}
//		}
//
//		status := fmt.Sprintf("defectCount: %d", defectCount+1)
//
//		rect := gocv.BoundingRect(c)
//		gocv.Rectangle(&img, rect, color.RGBA{255, 255, 255, 0}, 2)
//
//		gocv.PutText(&img, status, image.Pt(10, 20), gocv.FontHersheyPlain, 1.2, red, 2)
//
//		if ImShow(img, window) {
//			break
//		}
//	}
//}
