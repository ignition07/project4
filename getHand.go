// This program will help you get the hsv values for the glove

package main

import (
	"fmt"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/dji/tello"
	"gocv.io/x/gocv"
	"io"
	"log"
	"os/exec"
	"time"
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

	// ffmpeg stuff
	ffmpeg := exec.Command("ffmpeg", "-i", "pipe:0", "-pix_fmt", "bgr24", "-vcodec", "rawvideo",
		"-an", "-sn", "-s", "960x720", "-f", "rawvideo", "pipe:1")
	ffmpegIn, _ := ffmpeg.StdinPipe()
	ffmpegOut, _ := ffmpeg.StdoutPipe()
	if err := ffmpeg.Start(); err != nil {
		fmt.Println(err)
		return
	}

	// drone stuff
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

	//video, _ := gocv.OpenVideoCapture(0)
	//img := gocv.NewMat()

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

		//video.Read(&img)
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
