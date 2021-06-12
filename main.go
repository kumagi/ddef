
package main

import (
	"log"
	"image/color"
	"time"
	"strconv"
	"math"
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
	points = 1024
)

type Game struct {
	count int
	fromx, fromy, tox, toy [points]int
}

func abs(a int) int {
    return int(math.Abs(float64(a)))
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
		g.tox[i] = rand.Intn(screenWidth)
		g.fromx[i] = g.tox[i]
		g.toy[i] = rand.Intn(screenHeight)
		g.fromy[i] = g.toy[i]
	}
	return g
}

func (g *Game) Update() error {
	g.count++
	x, y := ebiten.CursorPosition()
	ebiten.SetWindowTitle(strconv.Itoa(x) + ":" + strconv.Itoa(y))
	for i := 0; i < points; i++ {
		g.fromx[i] = g.tox[i]
		g.fromy[i] = g.toy[i]
		g.tox[i] += (g.tox[i] - x) / 10
		g.toy[i] += (g.toy[i] - y) / 10
		if g.tox[i] < 0 || screenWidth < g.tox[i] || g.toy[i] < 0 || screenHeight < g.toy[i] {
			g.tox[i] = rand.Intn(screenWidth)
			g.toy[i] = rand.Intn(screenHeight)
			g.fromx[i] = g.tox[i]
			g.fromy[i] = g.toy[i]
		}
		
	}
	return nil
}

func (g *Game) Draw(img *ebiten.Image) {
	for i := 0; i < points; i++ {
		DrawLine(img, g.fromx[i], g.fromy[i], g.tox[i], g.toy[i], color.RGBA{0xbb, 0xdd, 0xff, 0xfe})
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
