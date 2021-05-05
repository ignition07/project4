// This file includes the drone functions:

// GetHand()
// DetectBlueHand()

package main

import (
	"fmt"
	"gobot.io/x/gobot"
	"gocv.io/x/gocv"
	"image"
	"io"
	"log"
	"time"
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

	ffmpegOut, drone := ConnectDrone()

	fmt.Println(drone)

	//gobot.After(3*time.Second, func() {
	//	drone.TakeOff()
	//	fmt.Println("Tello Taking Off...")
	//	time.Sleep(time.Second * 3)
	//})

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
		gocv.Flip(img, &img, 1)
		gocv.Flip(thresholded, &thresholded, 1)
		if ImShow(img, wi) || ImShow(thresholded, wt) {
			break
		}

	}
}

///////////////////////////////////////////////////////////////////////////////////////////

// DetectBlueHand detects glove using drone camera and ffmpeg
func DetectBlueHand() {

	xmlFile := "haarcascade_frontalface_default.xml"
	// load classifier to recognize faces
	classifier := gocv.NewCascadeClassifier()
	defer classifier.Close()
	if !classifier.Load(xmlFile) {
		fmt.Printf("Error reading cascade file: %v\n", xmlFile)
		return
	}

	size := image.Point{X: 600, Y: 600}
	blur := image.Point{X: 11, Y: 11}

	wt := gocv.NewWindow("Tello Drone")

	wt.ResizeWindow(2000, 2000)
	wt.MoveWindow(0, 0)

	ffmpegOut, drone := ConnectDrone()

	//TakeOff the Drone
	gobot.After(3*time.Second, func() {
		drone.TakeOff()
		fmt.Println("Tello Taking Off...")
		time.Sleep(time.Second * 3)
	})

	gobot.After(60*time.Second, func() {
		drone.Land()
	})

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

		rects := classifier.DetectMultiScaleWithParams(img, 1.1, 1, 0, image.Pt(50, 50), image.Pt(800, 800))
		//rects := classifier.DetectMultiScale(img)

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

		fingers := GetDefectCount(img, c, 0, 80) + 1
		thumb := GetDefectCount(img, c, 90, 110)
		fingerCount = append(fingerCount, fingers)
		fingersMode := GetMode(fingerCount)
		thumbCount = append(thumbCount, thumb)
		thumbMode := GetMode(thumbCount)

		if !heightReached {
			status1 = fmt.Sprintf("No Face Detected")
			drone.Up(15)
		}

		if len(fingerCount) > 50 {

			if len(rects) > 0 {
				status1 = fmt.Sprintf("I see you!")
				heightReached = true
				instruction := GetInstruction(fingersMode, thumbMode)
				DoInstruction(instruction, drone)
			} else {
				status1 = fmt.Sprintf("No Face Detected")
			}

			fingerCount = Trim(fingerCount)
			thumbCount = Trim(thumbCount)
		}

		gocv.PutText(&img, status1, image.Pt(10, 20), gocv.FontHersheyPlain, 2, red, 2)
		gocv.PutText(&img, status2, image.Pt(10, 50), gocv.FontHersheyPlain, 2, blue, 2)

		if ImShow(img, wt) {
			break
		}
	}
}
