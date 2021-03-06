// Shared functions go here:

// ConnectDrone
// GetBiggestContour
// GetPosFloat
// ImShow
// GetInstruction
// DoInstruction
// Trim
// GetMode

package main

import (
	"fmt"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/dji/tello"
	"gocv.io/x/gocv"
	"image/color"
	"io"
	"math"
	"os"
	"os/exec"
	"time"
)

func ConnectDrone() (io.ReadCloser, *tello.Driver) {
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

	return ffmpegOut, drone
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
	//status := fmt.Sprintf("defectCount: %d", defectCount+1)

	rect := gocv.BoundingRect(c)
	gocv.Rectangle(&img, rect, color.RGBA{R: 255, G: 255, B: 255}, 2)
	//gocv.PutText(&img, status, image.Pt(10, 20), gocv.FontHersheyPlain, 1.2, red, 2)

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

// GetInstruction checks hand position and returns an int
func GetInstruction(fingers int, thumb int) int {

	if !flipping {
		if thumb == 1 { // flip front
			if !flipFront {
				status2 = fmt.Sprintf("I see 1 thumb - flip forward")
				flipLeft = false
				flipRight = false
				flipBack = false
				flipFront = true
				//return 0
			}
		} else {
			if fingers == 1 { // hover
				return 5
			} else if fingers == 2 { // flip left
				if !flipLeft {
					status2 = fmt.Sprintf("I see 2 fingers - flip left")
					flipLeft = true
					flipRight = false
					flipBack = false
					flipFront = false
					return 1
				}
			} else if fingers == 3 { // flip back
				if !flipBack {
					status2 = fmt.Sprintf("I see 3 fingers - flip back")
					flipLeft = false
					flipRight = false
					flipBack = true
					flipFront = false
					return 2
				}
			} else if fingers == 4 { // flip right
				if !flipRight {
					status2 = fmt.Sprintf("I see 4 fingers - flip right")
					flipLeft = false
					flipRight = true
					flipBack = false
					flipFront = false
					return 3
				}
			} else if fingers == 5 { // land
				status2 = fmt.Sprintf("I see 5 fingers - landing drone")
				return 4
			} else if fingers > 5 || thumb > 1 { // ABORT
				return 4
			} else {
				status2 = fmt.Sprintf("Waiting for instruction")
				return 5 // hover
			}
		}
	}
	return 5
}

//////////////////////////////////////////////////////////////////////

// DoInstruction takes an int and executes movement
func DoInstruction(instruction int, drone *tello.Driver) {

	if instruction < 4 {
		if !flipping {
			flipping = true
			drone.Flip(tello.FlipType(instruction))
		}

	} else if instruction == 4 {
		status1 = fmt.Sprintf("Landing. Goodbye!")
		fmt.Println("LANDING")
		drone.Land()
		fmt.Println("GOOD BYE")
		os.Exit(0)

	} else {
		flipping = false
		drone.Hover()
	}
}

//////////////////////////////////////////////////////////////////////

// Trim takes a slice and returns a new slice with the last 2 values
func Trim(s []int) []int {
	s[0] = s[len(s)-2]
	s[0] = s[len(s)-2]
	return s[:2]
}

//////////////////////////////////////////////////////////////////////

// GetMode from github user Napear: https://gist.github.com/Napear/df41f13bfb5c10566102
func GetMode(s []int) (mode int) {
	//	Create a map and populated it with each value in the slice
	//	mapped to the number of times it occurs
	countMap := make(map[int]int)
	for _, value := range s {
		countMap[value] += 1
	}

	//	Find the smallest item with greatest number of occurrence in
	//	the input slice
	max := 0
	for _, key := range s {
		freq := countMap[key]
		if freq > max {
			mode = key
			max = freq
		}
	}
	return
}
