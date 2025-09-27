package main

import (
	"fmt"
	_ "image/png"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var flg_x, flg_y, flg_j int8
var vy float64 = 0.0
var y float64 = 500.0
var x float64 = 540.0
var t float64 = 0.0
var d float64 = 0.0

type game struct {
	img *ebiten.Image
}

func newGame() (*game, error) {
	g := &game{}

	// 画像を読み込む
	img, _, err := ebitenutil.NewImageFromFile("koronesuki.png")
	if err != nil {
		return nil, err
	}
	g.img = img

	return g, nil
}

func (g *game) Update() error {
	return nil
}

func (g *game) Draw(screen *ebiten.Image) {

	// 画像のOption定義
	op := &ebiten.DrawImageOptions{}

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		// ジャンプ用フラグ管理
		flg_y = 1

		// 画像の回転
		bounds := g.img.Bounds()
		op.GeoM.Translate(-float64(bounds.Dx())/2, -float64(bounds.Dy())/2)
		dx, dy := ebiten.CursorPosition()
		t = math.Atan2((y-float64(dy)), (x-float64(dx))) + math.Pi*2/4
		op.GeoM.Rotate(t)

		// 画像との距離
		d = math.Sqrt((x-float64(dx))*(x-float64(dx)) + (y-float64(dy))*(y-float64(dy)))

		// 画像描画
		op.GeoM.Translate(x, y)
		screen.DrawImage(g.img, op)
	} else {
		if flg_y == 1 {
			flg_j = 1
		}
		flg_y = 0
		bounds := g.img.Bounds()
		op.GeoM.Translate(-float64(bounds.Dx())/2, -float64(bounds.Dy())/2)
		op.GeoM.Rotate(t)
		op.GeoM.Translate(x, y)
		screen.DrawImage(g.img, op)
	}

	if flg_j == 1 {
		if d < 100 {
			vy = -5 // 小ジャンプ
		} else if 100 < d && d < 200 {
			vy = -15 // 中ジャンプ
		} else {
			vy = -30 // 大ジャンプ
		}
	}
	vy += 0.5 // 速度に加速度を足す（重力）
	y += vy   // 位置に速度を足す

	if 500 < y {
		y = 500
		flg_j = 0
	}
	if y < 0 {
		y = 0
		flg_j = 0
	}
	/*
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			flg_y = 1
			flg_x = 1
		}
		if flg_y == 1 {
			vy = -10 // ジャンプ
		}
		if flg_x == 1 {
			x += 2
		}
		vy += 0.5 // 速度に加速度を足す（重力）
		y += vy   // 位置に速度を足す

		if y < 180 {
			flg_y = 0
		}

		if 300 < y {
			y = 300
			flg_x = 0
		}
	*/
	s := fmt.Sprintln(t+math.Pi*2/4, x, y, d)
	ebitenutil.DebugPrint(screen, s)
}

func (g *game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func main() {
	ebiten.SetWindowTitle("sample")
	ebiten.SetWindowSize(1280, 720)
	g, err := newGame()
	if err != nil {
		panic(err)
	}
	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
}
