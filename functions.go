package main

import "gocv.io/x/gocv"

// ImShow opens image window for specified gocv.Mat
func ImShow(img gocv.Mat, window *gocv.Window) bool {
	window.IMShow(img)
	return window.WaitKey(1) == 27
}

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

func GetPosFloat(t *gocv.Trackbar) float64 {
	return float64(t.GetPos())
}
