package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	screenWidth         = 640 // 画面の幅
	screenHeight        = 480 // 画面の高さ
	charSize            = 20  // キャラクターのサイズ
	logicalScreenWidth  = 320 // 論理的な画面の幅
	logicalScreenHeight = 240 // 論理的な画面の高さ
)

type Game struct {
	x, y                 float64 // キャラクターの位置
	obstacleX, obstacleY float64 // 障害物の位置
	obstacleSpeed        float64 // 障害物の速度
}

func (g *Game) Update() error {
	// キーボード入力に応じて位置を更新
	newX, newY := g.x, g.y
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		newY--
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		newY++
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		newX--
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		newX++
	}

	// 画面の範囲内にキャラクターが収まるように調整
	g.x = max(0, min(newX, logicalScreenWidth-charSize))
	g.y = max(0, min(newY, logicalScreenHeight-charSize))

	// 障害物を左に移動
	g.obstacleX -= g.obstacleSpeed

	// 画面の左端に到達したら、右端から再スタート
	if g.obstacleX < -20 { // 20は障害物の大きさを想定した値
		g.obstacleX = logicalScreenWidth
	}

	return nil
}

// min は2つの値のうち小さい方を返します。
func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

// max は2つの値のうち大きい方を返します。
func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DrawRect(screen, g.x, g.y, 20, 20, color.RGBA{255, 0, 0, 255})

	radius := 10.0
	obstacleColor := color.RGBA{0, 128, 0, 255} // 緑色

	// 円を描画するための画像を作成
	obstacleImage := ebiten.NewImage(int(radius*2), int(radius*2))
	for y := -radius; y < radius; y++ {
		for x := -radius; x < radius; x++ {
			if x*x+y*y <= radius*radius {
				obstacleImage.Set(int(x+radius), int(y+radius), obstacleColor)
			}
		}
	}

	// 画像を画面に描画
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(g.obstacleX, g.obstacleY)
	screen.DrawImage(obstacleImage, opts)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return logicalScreenWidth, logicalScreenHeight
}

func main() {
	game := &Game{
		x:             50,                      // キャラクターの初期位置X
		y:             logicalScreenHeight / 2, // キャラクターの初期位置Y
		obstacleX:     logicalScreenWidth,      // 画面の右端から開始
		obstacleY:     logicalScreenHeight / 2, // 画面の中央の高さ
		obstacleSpeed: 2,                       // 移動速度を適宜設定
	}
	// ...（Ebitenの設定）
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Avoid Game")

	// ゲームを実行し、gameインスタンスを渡します。
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}