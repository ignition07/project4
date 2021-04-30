package main

import (
	"gocv.io/x/gocv"
	"image/color"
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
	*/
	lhsv = gocv.Scalar{Val1: 101, Val2: 133, Val3: 71}
	hhsv = gocv.Scalar{Val1: 123, Val2: 255, Val3: 255}

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

	// Colors
	blue  = color.RGBA{B: 255}
	red   = color.RGBA{R: 255}
	green = color.RGBA{G: 255}
)
