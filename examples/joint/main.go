package main

import (
	"box2d/examples/yokuaruyatu/objects"
	b2d "github.com/E4/box2d"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"image/color"
	"log"
)

/**
なんかあんまり知識入ってこなかった気がするので、
よくあるクリックすると箱を出すやつを作ってみることにする。
最終的には、ボールをゴールに運んだらクリアっていうよくあるやつを作りたい。
*/

var World b2d.B2World

type CustomJoint interface {
	GetAnchorA() b2d.B2Vec2
	GetAnchorB() b2d.B2Vec2
}

var joints []CustomJoint
var boxes []objects.PolygonObject

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
	// Creating a World
	//////////////////////////////////////////////
	// Step.1
	gravity := b2d.NewB2Vec2(0.0, 2)
	// step.2
	World = b2d.MakeB2World(*gravity)

	{
		// Jointの勉強をする
		// DynamicBodyのほうは普通に作ればいいっぽい
		bodyA := objects.NewDynamicBox(2, 1, 0.4, 0.1, 1.0, &World)
		//bodyA.Body.SetFixedRotation(false)
		bodyA.Body.SetFixedRotation(true) // これで物体の回転自体を制御できる

		bodyB := objects.NewPolygonBox(3, 3, 0.4, 0.1, &World, 0)
		addBox(bodyA)
		addBox(bodyB)

		// Distance Joint
		jointDef := b2d.MakeB2DistanceJointDef()
		// アンカーポイントは世界の座標を指定することに注意
		jointDef.Initialize(bodyA.Body, bodyB.Body, bodyA.Body.GetPosition(), bodyB.Body.GetPosition())
		jointDef.CollideConnected = true
		World.CreateJoint(&jointDef)
		// ジョイントの管理のためにグローバル変数に入れるけど、なんか気持ち悪いなー、種類よって端っこの取り方が違うのかな
		joints = append(joints, b2d.MakeB2DistanceJoint(&jointDef))

		// いろんなJointをしてみる
		// Prismatic Joint
		// よくわからないけどトルク？のジョイントらしい。今回のは無関係かも
		movingBox := objects.NewDynamicBox(2, 1-0.4, 0.2, 0.2, 1.0, &World)
		addBox(movingBox)
		prismaticJointDef := b2d.MakeB2PrismaticJointDef()
		prismaticJointDef.Initialize(bodyA.Body, movingBox.Body, bodyA.Body.GetPosition(), movingBox.Body.GetPosition())
		joints = append(joints, b2d.MakeB2PrismaticJoint(&prismaticJointDef))
	}

	{
		// Jointの勉強をする
		// DynamicBodyのほうは普通に作ればいいっぽい
		bodyA := objects.NewDynamicBox(6, 1, 0.4, 0.1, 1.0, &World)
		//bodyA.Body.SetFixedRotation(false)
		bodyA.Body.SetFixedRotation(true) // これで物体の回転自体を制御できる

		bodyB := objects.NewPolygonBox(7, 3, 0.4, 0.1, &World, 0)
		addBox(bodyA)
		addBox(bodyB)

		// Distance Joint
		jointDef := b2d.MakeB2DistanceJointDef()
		// アンカーポイントは世界の座標を指定することに注意
		jointDef.Initialize(bodyA.Body, bodyB.Body, bodyA.Body.GetPosition(), bodyB.Body.GetPosition())
		jointDef.CollideConnected = true
		World.CreateJoint(&jointDef)
		// ジョイントの管理のためにグローバル変数に入れるけど、なんか気持ち悪いなー、種類よって端っこの取り方が違うのかな
		joints = append(joints, b2d.MakeB2DistanceJoint(&jointDef))

		// いろんなJointをしてみる
		// Distance Joint よくわからないけど固まるw
		movingBox := objects.NewDynamicBox(6.1, 1-0.4, 0.2, 0.2, 1.0, &World)
		addBox(movingBox)
		// Distance Joint
		distanceJointDef := b2d.MakeB2DistanceJointDef()
		// アンカーポイントは世界の座標を指定することに注意
		distanceJointDef.Initialize(bodyA.Body, bodyB.Body, b2d.B2Vec2MulScalar(0.99, bodyA.Body.GetPosition()), movingBox.Body.GetPosition())
		distanceJointDef.CollideConnected = true
		World.CreateJoint(&distanceJointDef)
		joints = append(joints, b2d.MakeB2DistanceJoint(&distanceJointDef))
	}

	{
		// DynamicBodyのほうは普通に作ればいいっぽい
		bodyA := objects.NewDynamicBox(10, 1, 0.4, 0.1, 1.0, &World)
		//bodyA.Body.SetFixedRotation(false)
		bodyA.Body.SetFixedRotation(true) // これで物体の回転自体を制御できる

		bodyB := objects.NewPolygonBox(11, 3, 0.4, 0.1, &World, 0)
		addBox(bodyA)
		addBox(bodyB)

		// WheelJoint  Joint なんかシランがぶっ壊れる
		jointDef := b2d.MakeB2DistanceJointDef()
		// アンカーポイントは世界の座標を指定することに注意
		jointDef.Initialize(bodyA.Body, bodyB.Body, bodyA.Body.GetPosition(), bodyB.Body.GetPosition())
		jointDef.CollideConnected = true
		World.CreateJoint(&jointDef)
		// ジョイントの管理のためにグローバル変数に入れるけど、なんか気持ち悪いなー、種類よって端っこの取り方が違うのかな
		joints = append(joints, b2d.MakeB2DistanceJoint(&jointDef))

		// いろんなJointをしてみる
		// Distance Joint よくわからないけど固まるw
		movingBox := objects.NewDynamicBox(10.1, 1-0.4, 0.2, 0.2, 1.0, &World)
		addBox(movingBox)
		// Distance Joint
		wheelJointDef := b2d.MakeB2WheelJointDef()
		// アンカーポイントは世界の座標を指定することに注意
		wheelJointDef.Initialize(bodyA.Body, bodyB.Body, bodyA.Body.GetPosition(), movingBox.Body.GetPosition())
		World.CreateJoint(&wheelJointDef)
		joints = append(joints, b2d.MakeB2WheelJoint(&wheelJointDef))
	}

	//TODO: なんか上に乗っ型物と同時に動くようにしたい。摩擦を上げればいい？ どうやら物理エンジンの話じゃないっぽい。ほんと～？

	emptyImage.Fill(color.White)

}

const (
	screenWidth  = 640 * 2
	screenHeight = 480 * 2
)

type Game struct {
	count int
}

var (
	emptyImage = ebiten.NewImage(3, 3)
)

func (g *Game) Update() error {
	// 試しにたぶんこうだろうなーっていう衝突判定をする

	//////////////////////////////////////////////
	// Simulating the World
	//////////////////////////////////////////////
	timeStep := 1.0 / 60.0

	velocityIterations := 8 * 2
	positionIterations := 3 * 2
	World.Step(timeStep, velocityIterations, positionIterations)

	// mouse lick
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		// 新しい箱を作成する
		box := objects.NewDynamicBox(float64(float32(x)/objects.SCALE), float64(float32(y)/objects.SCALE), 0.1, 0.1, 0.5, &World)
		box.Fixture.SetRestitution(0.8)
		addBox(box)

	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.White)
	for _, v := range boxes {
		v.Draw(screen)
	}

	for _, v := range joints {
		scale := float64(objects.SCALE)
		ebitenutil.DrawLine(screen, v.GetAnchorA().X*scale, v.GetAnchorA().Y*scale, v.GetAnchorB().X*scale, v.GetAnchorB().Y*scale, color.Black)
	}
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
