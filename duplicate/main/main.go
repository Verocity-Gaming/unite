package main

import (
	"image/png"
	"os"
	"strconv"
	"time"

	"github.com/verocity-gaming/unitehud/duplicate"
	"github.com/verocity-gaming/unitehud/team"
	"gocv.io/x/gocv"
)

func matrix(file string) gocv.Mat {
	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	img, err := png.Decode(f)
	if err != nil {
		panic(err)
	}

	matrix, err := gocv.ImageToMatRGB(img)
	if err != nil {
		panic(err)
	}

	return matrix
}

func save(m gocv.Mat) {
	img, err := m.ToImage()
	if err != nil {
		panic(err)
	}

	f, err := os.Create(time.Now().Format(strconv.Itoa(int(time.Now().Unix())) + ".png"))
	if err != nil {
		panic(err)
	}

	err = png.Encode(f, img)
	if err != nil {
		panic(err)
	}
}

func main() {
	d1 := duplicate.New(0, team.Self, matrix("50_2.png"))
	d2 := duplicate.New(0, team.Self, matrix("50_1.png"))
	save(d1.Region())
	println(d1.Is(d2))
}
