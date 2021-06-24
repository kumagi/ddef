
package main

import (
	"log"
	"image/color"
	"time"
	"strconv"
	"math/rand"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
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
	fromx, fromy, tox, toy, blightness float64
}

func (s *Star) Init() {
	s.tox = rand.Float64() * screenWidth * scale
	s.fromx = s.tox
	s.toy = rand.Float64() * screenHeight * scale
	s.fromy = s.toy
	s.blightness = rand.Float64() * 0xff
}

func (s *Star) Out() bool {
	return s.fromx < 0 || screenWidth * scale < s.fromx || s.fromy < 0 || screenHeight * scale < s.fromy
}

func (s *Star) Update(x, y float64) {
	s.fromx = s.tox
	s.fromy = s.toy
	s.tox += (s.tox - x) / 32
	s.toy += (s.toy - y) / 32
	s.blightness += 1
	if 0xff < s.blightness {
		s.blightness = 0xff
	}
	if s.Out() {
		s.Init()
	}
}

func (s *Star) Pos() (float64, float64, float64, float64) {
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

func NewGame() *Game {
	g := new(Game)
	for i := 0; i < points; i++ {
		g.stars[i].Init()
	}
	return g
}

func (g *Game) Update() error {
	x, y := ebiten.CursorPosition()
	x *= scale
	y *= scale
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
		ebitenutil.DrawLine(img, fx, fy, tx, ty, color.RGBA{uint8(0xbb * s.blightness / 0xff),
			uint8(0xdd * s.blightness / 0xff), uint8(0xff * s.blightness / 0xff), 0xff})
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
