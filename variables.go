package main

import (
	"fmt"
	"gocv.io/x/gocv"
	"image/color"
)

const (
	frameSize = 960 * 720 * 3
)

var (
	/*  HSV values
	--------------
	These values seem good:
		lhsv   = gocv.Scalar{Val1: 106, Val2: 92, Val3: 191}
		hhsv   = gocv.Scalar{Val1: 109, Val2: 255, Val3: 255}
	These too:
		lhsv   = gocv.Scalar{Val1: 101, Val2: 133, Val3: 71}
		hhsv   = gocv.Scalar{Val1: 123, Val2: 255, Val3: 255}
	Good for drone:
		lhsv   = gocv.Scalar{Val1: 109, Val2: 128, Val3: 31}
		hhsv   = gocv.Scalar{Val1: 128, Val2: 255, Val3: 255}
	Also good for drone:
		lhsv = gocv.Scalar{Val1: 108, Val2: 130, Val3: 0}
		hhsv = gocv.Scalar{Val1: 130, Val2: 193, Val3: 255}
	*/

	lhsv = gocv.Scalar{Val1: 99, Val2: 95, Val3: 0}
	hhsv = gocv.Scalar{Val1: 122, Val2: 255, Val3: 210}

	// Mats
	mask      = gocv.NewMat()
	hsv       = gocv.NewMat()
	frame     = gocv.NewMat()
	kernel    = gocv.NewMat()
	imgGrey   = gocv.NewMat()
	imgBlur   = gocv.NewMat()
	imgThresh = gocv.NewMat()
	hull      = gocv.NewMat()
	defects   = gocv.NewMat()

	// Movement
	fingerCount   []int
	thumbCount    []int
	flipLeft      = false
	flipRight     = false
	flipBack      = false
	flipFront     = false
	flipping      = false
	heightReached = false

	status1 = fmt.Sprintf("")
	status2 = fmt.Sprintf("")

	// Colors
	blue  = color.RGBA{B: 255}
	red   = color.RGBA{R: 255}
	green = color.RGBA{G: 255}
	black = color.RGBA{}
)
