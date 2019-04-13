package main

import (
	"log"
	"math"
	"unicode/utf8"

	"github.com/fogleman/gg"
)

const PageSize = 1000
const Radius = 480

type Point struct {
	X, Y float64
}

func Polygon(n int, x, y, r float64) []Point {
	result := make([]Point, n)
	for i := 0; i < n; i++ {
		a := float64(i)*2*math.Pi/float64(n) - math.Pi/2
		result[i] = Point{x + r*math.Cos(a), y + r*math.Sin(a)}
	}
	return result
}

func position(radius, angle float64) Point {
	p := Point{}
	p.X = radius * math.Cos(angle*math.Pi/180.0)
	p.Y = radius * math.Sin(angle*math.Pi/180.0)
	return p
}

func drawText(dc *gg.Context, text string, xpos, ypos, rotation float64) {
	if rotation != 0 {
		dc.RotateAbout(rotation, xpos, ypos)
	}
	dc.DrawStringAnchored(text, xpos, ypos, 0.5, 0.5)
	if rotation != 0 {
		dc.RotateAbout(-rotation, xpos, ypos)
	}
}

func Draw(message string) {

	dc := gg.NewContext(PageSize, PageSize)

	// background
	dc.DrawRectangle(0, 0, PageSize, PageSize)
	dc.SetRGB(1, 1, 1)
	dc.Fill()

	// cover
	dc.DrawCircle(PageSize/2, PageSize/2, Radius+16)
	dc.SetRGB(1, 0, 0)
	dc.Fill()

	// main
	dc.DrawCircle(PageSize/2, PageSize/2, Radius)
	dc.SetRGB(1, 1, 1)
	dc.Fill()

	// draw string in circle
	err := dc.LoadFontFace("wqy-microhei.ttc", 64)
	if err != nil {
		log.Fatalf(err.Error())
	}
	dc.SetRGB(1, 0, 0)
	var startAngle, stepAngel float64
	strLen := utf8.RuneCountInString(message)
	// 0 90
	// 90 180
	// 180 270
	// 120 210
	// 270 0
	startAngle = 180 + (-45)
	stepAngel = 320.0 / float64(strLen)
	var testAngle float64 = startAngle + 90
	for _, c := range message {
		p := position(Radius-32, startAngle)
		drawText(dc, string(c), p.X+PageSize/2, p.Y+PageSize/2, gg.Radians(testAngle))
		startAngle += stepAngel
		testAngle += stepAngel
	}

	// star
	n := 5
	points := Polygon(n, PageSize/2, PageSize/2, Radius*0.6)
	for i := 0; i < n+1; i++ {
		index := (i * 2) % n
		p := points[index]
		dc.LineTo(p.X, p.Y)
	}
	dc.Fill()

	dc.SavePNG("out.png")
}

func main() {
	Draw("深圳市零成本科技股份有限公司")
}
