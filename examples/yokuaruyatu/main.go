package main

import (
	"box2d/examples/yokuaruyatu/objects"
	"fmt"
	b2d "github.com/E4/box2d"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"image"
	"image/color"
	_ "image/png"
	"log"
)

/**
なんかあんまり知識入ってこなかった気がするので、
よくあるクリックすると箱を出すやつを作ってみることにする。
最終的には、ボールをゴールに運んだらクリアっていうよくあるやつを作りたい。
*/

var World *b2d.B2World

type CustomJoint interface {
	GetAnchorA() b2d.B2Vec2
	GetAnchorB() b2d.B2Vec2
}

var (
	joints          []CustomJoint
	boxes           []objects.PolygonObject
	circle          objects.CircleObject
	emptyImage      = ebiten.NewImage(3, 3)
	backgroundImage *ebiten.Image
)

func addBox(box objects.PolygonObject) {
	boxes = append(boxes, box)
}

/*
参考リンク。
この辺で hello worldして物理エンジン動かしてみる。
https://qiita.com/zenwerk/items/d15ee04335e1d1b8217b
https://github.com/E4/box2d
https://box2d.org/documentation/md__d_1__git_hub_box2d_docs_hello.html
https://github.com/zenwerk/ebiten-example/commit/141fd83be850e7a89b0be12204edc3c13480d31f

これもよさそう
http://vivi.dyndns.org/tech/Qt/Box2D_HelloWorld.html

// testbed sources
https://github.com/erincatto/box2d/tree/main/testbed/tests
// testbed
https://flyover.github.io/box2d.ts/testbed/
*/
func init() {
	//////////////////////////////////////////////
	// ebiten
	//////////////////////////////////////////////
	var err error
	backgroundImage, _, err = ebitenutil.NewImageFromFile("examples/yokuaruyatu/assets/background.png")
	if err != nil {
		panic(err)
	}
	emptyImage.Fill(color.White)

	//////////////////////////////////////////////
	// Creating a World
	//////////////////////////////////////////////
	gravity := b2d.NewB2Vec2(0.0, 2)
	w := b2d.MakeB2World(*gravity)
	World = &w
	groundBodyProxy = World.CreateBody(b2d.NewB2BodyDef())

	// 静的な四角をマウスクリックで移動できるようにする
	kinematic := objects.NewDynamicBox(4, 4, 1, 0.2, 1, World)
	kinematic.Body.SetType(b2d.B2BodyType.B2_staticBody)
	kinematic.Body.SetFixedRotation(true)
	// 画面上に設置しているオブジェクトはすべてstaticにするため衝突しない
	addBox(kinematic)

}

const (
	screenWidth  = 640 * 2
	screenHeight = 480 * 2
)

type Game struct {
	count int
}

const (
	MouseNone = iota
	MouseDrag
)

var MouseStatus = MouseNone
var selectedBlock *objects.PolygonObject
var selectedMouseJoint *b2d.B2MouseJointDef
var mouseB2Vec *b2d.B2Vec2
var groundBodyProxy *b2d.B2Body
var mouseJointDef *b2d.B2MouseJointDef
var mouseJoint *b2d.B2MouseJoint

func (g *Game) Update() error {
	//////////////////////////////////////////////
	// Simulating the World
	//////////////////////////////////////////////
	timeStep := 1.0 / 60.0

	velocityIterations := 8 * 2
	positionIterations := 3 * 2
	World.Step(timeStep, velocityIterations, positionIterations)

	x, y := ebiten.CursorPosition()
	// メートル座標に戻す
	mx := float64(x) / objects.SCALE
	my := float64(y) / objects.SCALE

	// ドラッグで移動は意外と面倒だったので、クリック＆クリック方式に変更
	// マウス関連はマジで苦戦した。正解はWorldに登録したMouseJointを保持して Targetを更新していくでした～
	// https://github.com/erincatto/box2d/blob/c6cc3646d1701ab3c0750ef397d2d68fc6dbcff2/testbed/test.cpp
	switch MouseStatus {
	case MouseNone:
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			// ブロックを検索し、初めて衝突が検出されたものを取得、面倒なので後勝ちにする
			if MouseStatus == MouseNone {

				for _, v := range boxes {
					fmt.Printf("point=%f,%f :::: Position=%f,%f \r\n", mx, my, v.Body.GetPosition().X, v.Body.GetPosition().Y)
					if v.Fixture.TestPoint(b2d.MakeB2Vec2(mx, my)) { // ちょっとよくわからないが Shaeじゃなくて Fixtureの TestPointを使うみたい。 やはり TransFormの意味が分かっていない
						selectedBlock = &v
					}
				}
				if selectedBlock != nil {
					selectedBlock.Body.SetType(b2d.B2BodyType.B2_dynamicBody)
					MouseStatus = MouseDrag
					md := b2d.MakeB2MouseJointDef()
					mouseJointDef = &md
					mouseJointDef.SetBodyA(groundBodyProxy)
					mouseJointDef.SetBodyB(selectedBlock.Body)
					mouseJointDef.Target = b2d.MakeB2Vec2(mx, my)
					mouseJointDef.MaxForce = 1000 * selectedBlock.Body.GetMass()
					mouseJoint = World.CreateJoint(mouseJointDef).(*b2d.B2MouseJoint)
				}
			}
		}
	case MouseDrag:
		if mouseJointDef != nil {
			mouseJoint.SetTarget(b2d.MakeB2Vec2(mx, my))
		}
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) && mouseJoint != nil && selectedBlock != nil {
			selectedBlock.Body.SetType(b2d.B2BodyType.B2_staticBody)
			World.DestroyJoint(mouseJoint)
			mouseJoint = nil
			mouseJointDef = nil
			selectedBlock = nil
			MouseStatus = MouseNone
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	{
		op := &ebiten.DrawImageOptions{}
		screen.DrawImage(backgroundImage.SubImage(image.Rect(0, 0, 979, screenHeight)).(*ebiten.Image), op)
	}

	for _, v := range boxes {
		v.Draw(screen)
	}
	if circle.Body != nil {
		circle.Draw(screen)
	}

	for _, v := range joints {
		scale := objects.SCALE
		ebitenutil.DrawLine(screen, v.GetAnchorA().X*scale, v.GetAnchorA().Y*scale, v.GetAnchorB().X*scale, v.GetAnchorB().Y*scale, color.Black)
	}

	{
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(979, 0)
		screen.DrawImage(backgroundImage.SubImage(image.Rect(979, 0, screenWidth, screenHeight)).(*ebiten.Image), op)
	}

	ebitenutil.DebugPrint(screen, fmt.Sprintf("mouseStatus=%d", MouseStatus))

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Animation (Ebiten Demo)")
	ebiten.SetFullscreen(false)
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
