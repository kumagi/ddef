
package main

import (
	"log"
	"image/color"
	"time"
	"strconv"
	"math/rand"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  = 640
	screenHeight = 480

	frameOX     = 0
	frameOY     = 32
	frameWidth  = 32
	frameHeight = 32
	frameNum    = 8
	scale = 64
	points = 1024
)

type Star struct {
	fromx, fromy, tox, toy int
}

func (s *Star) Init() {
	s.tox = rand.Intn(screenWidth * scale)
	s.fromx = s.tox
	s.toy = rand.Intn(screenHeight * scale)
	s.fromy = s.toy
}

func (s *Star) Out() bool {
	return s.tox < 0 || screenWidth * scale < s.tox || s.toy < 0 || screenHeight * scale < s.toy
}

func (s *Star) Update(x, y float64) {
	s.fromx = s.tox
	s.fromy = s.toy
	s.tox += int((float64(s.tox) - x * scale) / 32)
	s.toy += int((float64(s.toy) - y * scale) / 32)
	if s.Out() {
		s.Init()
	}
}

func (s *Star) Pos() (int, int, int, int) {
	return s.fromx / scale, s.fromy / scale, s.tox / scale, s.toy / scale
}

type Game struct {
	stars [points]Star
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func DrawLine(img *ebiten.Image, fromx, fromy, tox, toy int, color color.RGBA) {
	steep := abs(toy - fromy) > abs(tox - fromx)
	if steep {
		fromy, fromx = fromx, fromy
		toy, tox = tox, toy
	}
	if fromx > tox {
		fromx, tox = tox, fromx
		fromy, toy = toy, fromy
	}
	if tox == toy {
		if toy < fromy {
			toy, fromy = fromy, toy
		}
		for y := fromy; y <= toy; y++ {
			img.Set(fromx, y, color)
		}
		return
	}
	dx := tox - fromx
	dy := abs(toy - fromy)
	error := dx / 2
	var step int
	if fromy < toy {
		step = 1
	} else {
		step = -1
	}
	for x, y := fromx, fromy; x <= tox; x++ {
		if steep {
			img.Set(y, x, color)	
		} else {
			img.Set(x, y, color)
		}
		error -= dy
		if error < 0 {
			y += step
			error += dx
		}
	}
}

func NewGame() *Game {
	g := new(Game)
	for i := 0; i < points; i++ {
		g.stars[i].Init()
	}
	return g
}

func (g *Game) Update() error {
	x, y := ebiten.CursorPosition()
	ebiten.SetWindowTitle(strconv.Itoa(x) + ":" + strconv.Itoa(y))
	for i := 0; i < points; i++ {
		g.stars[i].Update(float64(x), float64(y))
	}
	return nil
}

func (g *Game) Draw(img *ebiten.Image) {
	for i := 0; i < points; i++ {
		s := &g.stars[i]
		fx, fy, tx, ty := s.Pos()
		DrawLine(img, fx, fy, tx, ty, color.RGBA{0xbb, 0xdd, 0xff, 0xfe})
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	rand.Seed(time.Now().UnixNano())
	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Animation (Ebiten Demo)")
	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
