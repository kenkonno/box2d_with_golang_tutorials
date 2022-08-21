package main

import (
	"bytes"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	_ "image/png"
	"log"

	b2d "github.com/E4/box2d"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/images"
)

/*
参考リンク。
この辺で hello worldして物理エンジン動かしてみる。
https://qiita.com/zenwerk/items/d15ee04335e1d1b8217b
https://github.com/E4/box2d
https://box2d.org/documentation/md__d_1__git_hub_box2d_docs_hello.html
https://github.com/zenwerk/ebiten-example/commit/141fd83be850e7a89b0be12204edc3c13480d31f
*/
const (
	screenWidth  = 640
	screenHeight = 480
)

var (
	ebitenImage *ebiten.Image
)

type Game struct {
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "Nearest Filter (default) VS Linear Filter")

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(4, 4)
	op.GeoM.Translate(64, 64)
	// By default, nearest filter is used.
	screen.DrawImage(ebitenImage, op)

	op = &ebiten.DrawImageOptions{}
	op.GeoM.Scale(4, 4)
	op.GeoM.Translate(64, 64+240)
	// Specify linear filter.
	op.Filter = ebiten.FilterLinear
	screen.DrawImage(ebitenImage, op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

type System struct {
	World b2d.B2World
}

var s System

func init() {
	
	// https://box2d.org/documentation/md__d_1__git_hub_box2d_docs_hello.html

	//////////////////////////////////////////////
	// Creating a World
	//////////////////////////////////////////////
	// Step.1
	gravity := b2d.NewB2Vec2(0.0, -10.0)
	// step.2
	s.World = b2d.MakeB2World(*gravity)

	//////////////////////////////////////////////
	// Creating a Ground Box
	//////////////////////////////////////////////
	// Step.1
	groundBodyDef := b2d.NewB2BodyDef()
	groundBodyDef.Position.Set(0.0, -10.0)

	// Step.2
	groundBody := s.World.CreateBody(groundBodyDef)

	// Step.3
	groundBox := b2d.NewB2PolygonShape()
	groundBox.SetAsBox(50.0, 10.0)

	// Step.4
	groundBody.CreateFixture(groundBox, 0.0)

	//////////////////////////////////////////////
	// Creating a Ground Box
	//////////////////////////////////////////////
	// Step.1
	bodyDef := b2d.NewB2BodyDef()
	bodyDef.Type = b2d.B2BodyType.B2_dynamicBody
	bodyDef.Position.Set(0.0, 4.0)
	body := s.World.CreateBody(bodyDef)

	// Step.2
	dynamicBox := b2d.NewB2PolygonShape()
	dynamicBox.SetAsBox(1.0, 1.0)

	// Step.3
	fixtureDef := b2d.B2FixtureDef{
		Shape:    dynamicBox,
		Friction: 0.3,
		Density:  1.0,
	}

	// Step.4
	body.CreateFixtureFromDef(&fixtureDef)

	//////////////////////////////////////////////
	// Simulating the World
	//////////////////////////////////////////////
	timeStep := 1.0 / 60.0

	velocityIterations := 6
	positionIterations := 2

	for i := 0; i < 60; i++ {
		s.World.Step(timeStep, velocityIterations, positionIterations)
		position := body.GetPosition()
		angle := body.GetAngle()
		fmt.Printf("%4.2f %4.2f %4.2f\n", position.X, position.Y, angle)
	}

}

func main() {

	// Decode an image from the image file's byte slice.
	// Now the byte slice is generated with //go:generate for Go 1.15 or older.
	// If you use Go 1.16 or newer, it is strongly recommended to use //go:embed to embed the image file.
	// See https://pkg.go.dev/embed for more details.
	img, _, err := image.Decode(bytes.NewReader(images.Ebiten_png))
	if err != nil {
		log.Fatal(err)
	}

	ebitenImage = ebiten.NewImageFromImage(img)

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Filter (Ebiten Demo)")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
