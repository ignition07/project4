// Shared functions go here:

// ConnectDrone()
// GetBiggestContour()
// GetPosFloat()
// ImShow()

package main

import (
	"fmt"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/dji/tello"
	"gocv.io/x/gocv"
	"image"
	"image/color"
	"io"
	"math"
	"os/exec"
	"time"
)

func ConnectDrone() io.ReadCloser {
	ffmpeg := exec.Command("ffmpeg", "-i", "pipe:0", "-pix_fmt", "bgr24", "-vcodec", "rawvideo",
		"-an", "-sn", "-s", "960x720", "-f", "rawvideo", "pipe:1")
	ffmpegIn, _ := ffmpeg.StdinPipe()
	ffmpegOut, _ := ffmpeg.StdoutPipe()
	if err := ffmpeg.Start(); err != nil {
		fmt.Println(err)
	}

	drone := tello.NewDriver("8890")
	drone.On(tello.ConnectedEvent, func(data interface{}) {
		fmt.Println("Connected")
		drone.StartVideo()
		drone.SetExposure(1)
		drone.SetVideoEncoderRate(4)

		gobot.Every(100*time.Millisecond, func() {
			drone.StartVideo()
		})
	})
	drone.On(tello.VideoFrameEvent, func(data interface{}) {
		pkt := data.([]byte)
		if _, err := ffmpegIn.Write(pkt); err != nil {
			fmt.Println(err)
		}
	})
	robot := gobot.NewRobot("tello",
		[]gobot.Connection{},
		[]gobot.Device{drone},
	)
	robot.Start(false)
	return ffmpegOut
}

//////////////////////////////////////////////////////////////////////

// GetBiggestContour is part of gocv sample programs
func GetBiggestContour(contours gocv.PointsVector) gocv.PointVector {
	var area float64
	index := 0
	for i := 0; i < contours.Size(); i++ {
		newArea := gocv.ContourArea(contours.At(i))
		if newArea > area {
			area = newArea
			index = i
		}
	}
	return contours.At(index)
}

//////////////////////////////////////////////////////////////////////

// GetDefectCount checks for angles between min and max
func GetDefectCount(img gocv.Mat, c gocv.PointVector, min float64, max float64) int {

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
		if angle >= min && angle <= max {
			defectCount++
			gocv.Circle(&img, far, 1, red, 2)
		}

	}
	status := fmt.Sprintf("defectCount: %d", defectCount+1)

	rect := gocv.BoundingRect(c)
	gocv.Rectangle(&img, rect, color.RGBA{R: 255, G: 255, B: 255}, 2)
	gocv.PutText(&img, status, image.Pt(10, 20), gocv.FontHersheyPlain, 1.2, red, 2)

	return defectCount
}

// GetPosFloat is part of gocv sample programs
func GetPosFloat(t *gocv.Trackbar) float64 {
	return float64(t.GetPos())
}

//////////////////////////////////////////////////////////////////////

// ImShow opens image window for specified gocv.Mat
func ImShow(img gocv.Mat, window *gocv.Window) bool {
	window.IMShow(img)
	return window.WaitKey(1) == 27
}

//////////////////////////////////////////////////////////////////////

func RobotMovement(fingers int, thumb int) {
	/*
		Here are the following scenarios you should adhere to:

		Hover -- no defects detected, closed fist - palm facing you
		Move Left -- create an "L" shape with your hand, should detect a large angle as a defect and no normal defects
		Move Right -- Peace symbol with two fingers - detects 1 red defect with normal angle
		Move Backward --
		Move Forward --


	*/
	if fingers == 0 && thumb == 0 { // don't show hand or make a closed back fist to camera (former is better though)
		fmt.Println("HOVER")
	} else if fingers == 0 && thumb == 1 { // "L" shape hand
		moveLeft = true
		fmt.Println("MOVE LEFT!")
	} else if fingers == 1 && thumb == 1 { // "Peace" symbol - has one big angle and one small
		moveRight = true
		fmt.Println("MOVE RIGHT!")
	} else if fingers == 4 && thumb == 0 {
		moveBackward = true
		fmt.Println("MOVE BACKWARD")
	}
}
