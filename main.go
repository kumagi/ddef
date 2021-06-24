package main

import (
	"image/color"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	screenWidth  = 640
	screenHeight = 480
	scale        = 64
	starsNum     = 1024
)

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

type Star struct {
	fromx, fromy, tox, toy, brightness float64
}

func (s *Star) Init() {
	s.tox = rand.Float64() * screenWidth * scale
	s.fromx = s.tox
	s.toy = rand.Float64() * screenHeight * scale
	s.fromy = s.toy
	s.brightness = rand.Float64() * 0xff
}

func (s *Star) Update(x, y float64) {
	s.fromx = s.tox
	s.fromy = s.toy
	s.tox += (s.tox - x) / 32
	s.toy += (s.toy - y) / 32
	s.brightness += 1
	if 0xff < s.brightness {
		s.brightness = 0xff
	}
	if s.fromx < 0 || screenWidth*scale < s.fromx || s.fromy < 0 || screenHeight*scale < s.fromy {
		s.Init()
	}
}

func (s *Star) Draw(screen *ebiten.Image) {
	color := color.RGBA{uint8(0xbb * s.brightness / 0xff),
		uint8(0xdd * s.brightness / 0xff),
		uint8(0xff * s.brightness / 0xff),
		0xff}
	ebitenutil.DrawLine(screen, s.fromx/scale, s.fromy/scale, s.tox/scale, s.toy/scale, color)
}

type Game struct {
	stars [starsNum]Star
}

func NewGame() *Game {
	g := &Game{}
	for i := 0; i < starsNum; i++ {
		g.stars[i].Init()
	}
	return g
}

func (g *Game) Update() error {
	x, y := ebiten.CursorPosition()
	for i := 0; i < starsNum; i++ {
		g.stars[i].Update(float64(x*scale), float64(y*scale))
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for i := 0; i < starsNum; i++ {
		g.stars[i].Draw(screen)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	rand.Seed(time.Now().UnixNano())
	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Stars (Ebiten Demo)")
	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
