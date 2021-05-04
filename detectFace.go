package main

//func DetectFace() {
//
//	xmlFile := "haarcascade_frontalface_default.xml"
//	window := gocv.NewWindow("Face Detect")
//	defer window.Close()
//
//	window.ResizeWindow(1400, 1400)
//
//	// load classifier to recognize faces
//	classifier := gocv.NewCascadeClassifier()
//	defer classifier.Close()
//
//	if !classifier.Load(xmlFile) {
//		fmt.Printf("Error reading cascade file: %v\n", xmlFile)
//		return
//	}
//
//	ffmpegOut, drone := ConnectDrone()
//	fmt.Println(drone)
//
//	for {
//		buf := make([]byte, frameSize)
//		if _, err := io.ReadFull(ffmpegOut, buf); err != nil {
//			fmt.Println(err)
//			continue
//		}
//
//		img, err := gocv.NewMatFromBytes(720, 960, gocv.MatTypeCV8UC3, buf)
//		if err != nil {
//			log.Print(err)
//			continue
//		}
//		if img.Empty() {
//			continue
//		}
//
//		gocv.Flip(img, &img, 1)
//
//		// detect faces
//		//rects := classifier.DetectMultiScale(img)
//		rects := classifier.DetectMultiScaleWithParams(img, 1.1, 3,0,image.Pt(90, 90),image.Pt(160, 160))
//
//		// draw a rectangle around each face on the original image,
//		// along with text identifying as "Human"
//		for _, r := range rects {
//			gocv.Rectangle(&img, r, blue, 3)
//
//			size := gocv.GetTextSize("Human", gocv.FontHersheyPlain, 1.2, 2)
//			pt := image.Pt(r.Min.X+(r.Min.X/2)-(size.X/2), r.Min.Y-2)
//			gocv.PutText(&img, "Human", pt, gocv.FontHersheyPlain, 1.2, blue, 2)
//		}
//
//		if ImShow(img, window) {
//			break
//		}
//	}
//}

//func DetectFaceLaptop() {
//	deviceID, _ := 0, 0
//	xmlFile := "haarcascade_frontalface_default.xml"
//
//	// open webcam
//	webcam, err := gocv.VideoCaptureDevice(int(deviceID))
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	defer webcam.Close()
//
//	// open display window
//	window := gocv.NewWindow("Face Detect")
//	defer window.Close()
//
//	window.ResizeWindow(1400, 1400)
//
//	// prepare image matrix
//	img := gocv.NewMat()
//	defer img.Close()
//
//	// load classifier to recognize faces
//	classifier := gocv.NewCascadeClassifier()
//	defer classifier.Close()
//
//	if !classifier.Load(xmlFile) {
//		fmt.Printf("Error reading cascade file: %v\n", xmlFile)
//		return
//	}
//
//	fmt.Printf("start reading camera device: %v\n", deviceID)
//	for {
//		if ok := webcam.Read(&img); !ok {
//			fmt.Printf("cannot read device %d\n", deviceID)
//			return
//		}
//		if img.Empty() {
//			continue
//		}
//
//		gocv.Flip(img, &img, 1)
//
//		// detect faces
//		//rects := classifier.DetectMultiScale(img)
//
//		rects := classifier.DetectMultiScaleWithParams(img, 1.1, 1,0,image.Pt(200, 200),image.Pt(1000, 1000))
//
//		// draw a rectangle around each face on the original image,
//		// along with text identifying as "Human"
//		for _, r := range rects {
//			gocv.Rectangle(&img, r, blue, 3)
//
//			if ImShow(img, window) {
//				break
//			}
//		}
//
//	}
//}
