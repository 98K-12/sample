package main

import (
	"fmt"
	"image/color"
	_ "image/png"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var flg_click, flg_move int8
var x float64 = 200.0               //初期X座標
var y float64 = 500.0               //初期Y座標
var x_1st, y_1st float64 = 0.0, 0.0 //マウスが初めにクリックされたX,Y座標
var Rad float64 = -1.4              //ラジアン
var Rad_pre float64 = 0.0           //ラジアン前回値
var dx, dy float64 = 0.0, 0.0       //sin, cos値
var Speed float64 = 5.0

var flg_touch int8 = 0

type game struct {
	img     *ebiten.Image
	lineimg *ebiten.Image
}

func newGame() (*game, error) {
	g := &game{}

	// 画像を読み込む
	img, _, err := ebitenutil.NewImageFromFile("koronesuki.png")
	if err != nil {
		return nil, err
	}
	g.img = img

	lineimg, _, err := ebitenutil.NewImageFromFile("Line.png")
	if err != nil {
		return nil, err
	}
	g.lineimg = lineimg

	return g, nil
}

func (g *game) Update() error {
	return nil
}

func (g *game) Draw(screen *ebiten.Image) {

	// 画像のOption定義
	op := &ebiten.DrawImageOptions{}

	//タッチ判定用
	touchIDs := ebiten.AppendTouchIDs(nil)

	// クリック位置保存
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		Cur1stx, Cur1sty := ebiten.CursorPosition()
		x_1st = float64(Cur1stx)
		y_1st = float64(Cur1sty)
	}

	if flg_touch == 0 {
		if len(touchIDs) != 0 {
			Cur1stx, Cur1sty := ebiten.TouchPosition(touchIDs[0])
			x_1st = float64(Cur1stx)
			y_1st = float64(Cur1sty)
		}
	}

	//ボタンを押している間
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		flg_click = 1
		Speed = 5.0
		// 画像の回転
		// クリックしている間画像の中心点とクリック位置の2点を結ぶ直線方向に回転
		bounds := g.img.Bounds()
		op.GeoM.Translate(-float64(bounds.Dx())/2, -float64(bounds.Dy())/2)
		Curx, Cury := ebiten.CursorPosition()
		Rad = math.Atan2((y_1st - float64(Cury)), (x_1st - float64(Curx)))

		// クリック時の回転方向前回値更新
		if Rad == 0 {
			Rad = Rad_pre
		}

		op.GeoM.Rotate(Rad)

		// 画像とクリック位置の2点間の距離
		//d = math.Sqrt((x-float64(dx))*(x-float64(dx)) + (y-float64(dy))*(y-float64(dy)))

		// 画像描画
		op.GeoM.Translate(x, y)
		screen.DrawImage(g.img, op)
		vector.StrokeLine(screen, float32(x_1st), float32(y_1st), float32(Curx), float32(Cury), 1, color.RGBA{0xff, 0xff, 0xff, 0xff}, true)

	} else if len(touchIDs) != 0 {
		flg_touch = 1
		flg_click = 1
		Speed = 5.0
		// 画像の回転
		// タッチしている間画像の中心点とタッチ位置の2点を結ぶ直線方向に回転
		bounds := g.img.Bounds()
		op.GeoM.Translate(-float64(bounds.Dx())/2, -float64(bounds.Dy())/2)
		Curx, Cury := ebiten.TouchPosition(touchIDs[0])
		Rad = math.Atan2((y_1st - float64(Cury)), (x_1st - float64(Curx)))

		// タッチ時の回転方向前回値更新
		if Rad == 0 {
			Rad = Rad_pre
		}

		op.GeoM.Rotate(Rad)

		// 画像とタッチ位置の2点間の距離
		//d = math.Sqrt((x-float64(dx))*(x-float64(dx)) + (y-float64(dy))*(y-float64(dy)))

		// 画像描画
		op.GeoM.Translate(x, y)
		screen.DrawImage(g.img, op)
	} else { //ボタンを離しているとき
		if flg_click == 1 {
			flg_move = 1
			flg_click = 0
		}
		flg_touch = 0

		bounds := g.img.Bounds()
		op.GeoM.Translate(-float64(bounds.Dx())/2, -float64(bounds.Dy())/2)

		if flg_move == 1 {
			dx = math.Cos(Rad)
			dy = math.Sin(Rad)

			Speed += 0.3

			x = x + dx*Speed
			y = y + dy*Speed
		}
		Rad_pre = Rad
		op.GeoM.Rotate(Rad)
		op.GeoM.Translate(x, y)
		screen.DrawImage(g.img, op)
	}

	if x < 0 {
		flg_move = 0
		x = 0
	} else if 460 < x {
		flg_move = 0
		x = 460
	}

	if y < 0 {
		flg_move = 0
		y = 0
	} else if 720 < y {
		flg_move = 0
		y = 720
	}

	var deg = Rad * 180 / math.Pi
	s := fmt.Sprintln(Rad, x, y, dx, dy, deg, touchIDs)
	ebitenutil.DebugPrint(screen, s)
}

func (g *game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func main() {
	ebiten.SetWindowTitle("sample")
	ebiten.SetWindowSize(460, 720)
	g, err := newGame()
	if err != nil {
		panic(err)
	}
	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
}
