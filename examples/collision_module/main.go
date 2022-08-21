package main

import (
	"fmt"
	b2d "github.com/E4/box2d"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image"
	"image/color"
	"log"
)

// ShapeとかあるけどいったんBoxだけ想定でやる
type Block struct {
	HalfWidth  float64
	HalfHeight float64
	Body       *b2d.B2Body
}

type Point struct {
	X float32
	Y float32
}

func (b *Block) GetRectPath() (Point, Point, Point, Point) {
	// TODO: 回転とかに対応していないから足りない。
	// ほかのドキュメントを読み進めていい感じにできないか模索する

	// 左上
	lt := Point{X: float32(b.Body.GetPosition().X - b.HalfWidth), Y: float32(b.Body.GetPosition().Y - b.HalfHeight)}
	// 右上
	rt := Point{X: float32(b.Body.GetPosition().X + b.HalfWidth), Y: float32(b.Body.GetPosition().Y - b.HalfHeight)}
	// 右下
	rb := Point{X: float32(b.Body.GetPosition().X + b.HalfWidth), Y: float32(b.Body.GetPosition().Y + b.HalfHeight)}
	// 左下
	lb := Point{X: float32(b.Body.GetPosition().X - b.HalfWidth), Y: float32(b.Body.GetPosition().Y + b.HalfHeight)}

	return lt, rt, rb, lb
}

var ground Block
var dBox Block
var World b2d.B2World

/*
参考リンク。
この辺で hello worldして物理エンジン動かしてみる。
https://qiita.com/zenwerk/items/d15ee04335e1d1b8217b
https://github.com/E4/box2d
https://box2d.org/documentation/md__d_1__git_hub_box2d_docs_hello.html
https://github.com/zenwerk/ebiten-example/commit/141fd83be850e7a89b0be12204edc3c13480d31f

これもよさそう
http://vivi.dyndns.org/tech/Qt/Box2D_HelloWorld.html
*/
func init() {
	//////////////////////////////////////////
	// Shapes
	//////////////////////////////////////////

	// Circle Chapes
	circle := b2d.NewB2CircleShape()
	circle.M_p.Set(2.0, 3.0) // ポジションは設定
	circle.M_radius = 0.5    // 半径の設定

	// Polygon Shapes
	// 多角形 shape
	// CCW(逆時計回り）⊿の形になる
	vertices := []b2d.B2Vec2{
		{0.0, 0.0},
		{1.0, 0.0},
		{0.0, 1.0},
	}
	count := 3
	polygon := b2d.NewB2PolygonShape()
	polygon.Set(vertices, count)

	// Edge Shapes (線分の物体） 横幅1mの線分
	v1 := b2d.NewB2Vec2(0.0, 0.0)
	v2 := b2d.NewB2Vec2(1.0, 0.0)

	var edge b2d.B2EdgeShape
	edge.Set(*v1, *v2)

	// Chain Shapes
	// Ghost 何とかは一般的じゃないみたいなのでスキップ
	// 時計回り、反時計回りで判定の向きが異なる？
	vs := []b2d.B2Vec2{
		{1.7, 0.0},
		{1.0, 0.25},
		{0.0, 0.0},
		{-1.7, 0.4},
	}
	chain := b2d.MakeB2ChainShape()
	chain.CreateLoop(vs, 4)

	// なんか世界をつなげるサンプルらしい、よくわからない
	//b2ChainShape::CreateChain(const b2Vec2* vertices, int32 count,
	//	const b2Vec2& prevVertex, const b2Vec2& nextVertex);

	// Unityでもbox2dを改変したものを使用しているらしい

	/////////////////////////////////////
	// Geometric Queries
	/////////////////////////////////////

	// Shape Point Test
	// 点と図形の衝突判定らしい
	transform := b2d.NewB2Transform()
	transform.SetIdentity() // なにこれ？
	v := b2d.NewB2Vec2(5.0, 2.0)

	hit := polygon.TestPoint(*transform, *v)
	fmt.Printf("%v", hit)

	// Shape Ray Cast
	{
		transform := b2d.NewB2Transform()
		transform.SetIdentity()

		input := b2d.NewB2RayCastInput()
		input.P1.Set(0.0, 0.0)
		input.P2.Set(1.0, 0.0)
		input.MaxFraction = 1.0 // なにこれ？
		childIndex := 0

		output := b2d.MakeB2RayCastOutput()
		hit := polygon.RayCast(&output, *input, *transform, childIndex)
		if hit {
			// B2Vecの足し算引き算ってどうするの？
			b2d.B2Vec2Add(input.P2, input.P1)                                                                            // 足し算
			b2d.B2Vec2Sub(input.P2, input.P1)                                                                            // 引き算
			b2d.B2Vec2MulScalar(output.Fraction, b2d.B2Vec2Sub(input.P2, input.P1))                                      // スカラーの掛け算
			hitPoint := b2d.B2Vec2Add(input.P1, b2d.B2Vec2MulScalar(output.Fraction, b2d.B2Vec2Sub(input.P2, input.P1))) // 足し算
			// 線をoutputから得た係数をかけて当たった場所の座標を計算しているようです。
			fmt.Printf("%v", hitPoint)
			//hitPoint := input.P1 + output.Fraction*(input.P2-input.P1)
		}
	}

	////////////////////////////////////////////////
	// Pairwise Functions
	////////////////////////////////////////////////
	// 衝突判定モジュール。 shape同士の判定？
	// overlap
	{
		xfA := b2d.NewB2Transform()
		xfB := b2d.NewB2Transform()
		xfA.SetIdentity()
		xfB.SetIdentity()
		indexA := 0
		indexB := 0
		// なんかdestroyにポインタレシーバが付いているから circle 以外使えないのはバグっぽく感じるが・・・
		// all github で検索してもそもそもこのメソッドはあまり使われていない。気にしなくてもよい？
		overlap := b2d.B2TestOverlapShapes(*circle, indexA, *circle, indexB, *xfA, *xfB)
		fmt.Printf("%v", overlap)
	}

	// Contact Manifolds
	// こっちがshape同士の判定っぽい と思ったけど通常計算することはないらしい。 サンプルも略されているのでいったん省略する。

	// Distance
	// 距離の計算するらしい。 b2DistanceProxy?

	// TimeOfImpact
	// トンネル現象（通り抜けの検知に使うらしい、問題になったら振り返ればよさそう）

	// Dynamic Tree
	// 計算を効率化するために使いそう。通常では使うことはないらしい

	// Broad-phase
	// 通常では使うことはないらしい

	// https://box2d.org/documentation/md__d_1__git_hub_box2d_docs_hello.html
	/*
		なんかよくわからないが、以下の流れで進んでいるチュートリアル
		・Worldの生成
		・床の生成
		・動的BOXの生成
		・シミュレーション
		結果はボックスの落下を座標のログに出したもののよう。
		これだけだと全然イメージできないので、Ebitenで結果を描画できるようにする。
	*/

	//////////////////////////////////////////////
	// Creating a World
	//////////////////////////////////////////////
	// Step.1
	gravity := b2d.NewB2Vec2(0.0, 120.0)
	// step.2
	World = b2d.MakeB2World(*gravity)

	//////////////////////////////////////////////
	// Creating a Ground Box
	//////////////////////////////////////////////
	// Step.1
	groundBodyDef := b2d.NewB2BodyDef()
	groundBodyDef.Position.Set(0.0, 400.0)

	// Step.2
	groundBody := World.CreateBody(groundBodyDef)

	// Step.3
	groundBox := b2d.NewB2PolygonShape()
	groundBox.SetAsBox(300.0, 10.0)

	// Step.4
	f := groundBody.CreateFixture(groundBox, 0.0)
	f.GetBody().GetPosition()

	//////////////////////////////////////////////
	// Creating a Dynamic Body
	//////////////////////////////////////////////
	// Step.1
	bodyDef := b2d.NewB2BodyDef()
	bodyDef.Type = b2d.B2BodyType.B2_dynamicBody
	bodyDef.Position.Set(0.0, 0.0)
	body := World.CreateBody(bodyDef)

	// Step.2
	dynamicBox := b2d.NewB2PolygonShape()
	dynamicBox.SetAsBox(20.0, 20.0)

	// Step.3
	fixtureDef := b2d.B2FixtureDef{
		Shape:    dynamicBox,
		Friction: 0.3,
		Density:  1.0,
	}

	// Step.4
	body.CreateFixtureFromDef(&fixtureDef)

	// いったんグローバル変数に入れて描画できるようにする。
	ground = Block{
		HalfWidth:  300,
		HalfHeight: 10,
		Body:       groundBody,
	}

	dBox = Block{
		HalfWidth:  20,
		HalfHeight: 20,
		Body:       body,
	}
	emptyImage.Fill(color.White)

}

const (
	screenWidth  = 640
	screenHeight = 480
)

type Game struct {
	count int
}

var (
	emptyImage    = ebiten.NewImage(3, 3)
	emptySubImage = emptyImage.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)
)

func (g *Game) Update() error {
	//////////////////////////////////////////////
	// Simulating the World
	//////////////////////////////////////////////
	timeStep := 1.0 / 60.0

	velocityIterations := 8 * 2
	positionIterations := 3 * 2
	World.Step(timeStep, velocityIterations, positionIterations)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	ox := float32(300)
	oy := float32(0)

	// 床の描画
	var path vector.Path
	{
		lt, rt, rb, lb := ground.GetRectPath()
		path.MoveTo(lt.X+ox, lt.Y+oy)
		path.LineTo(rt.X+ox, rt.Y+oy)
		path.LineTo(rb.X+ox, rb.Y+oy)
		path.LineTo(lb.X+ox, lb.Y+oy)
	}
	{
		lt, rt, rb, lb := dBox.GetRectPath()
		path.MoveTo(lt.X+ox, lt.Y+oy)
		path.LineTo(rt.X+ox, rt.Y+oy)
		path.LineTo(rb.X+ox, rb.Y+oy)
		path.LineTo(lb.X+ox, lb.Y+oy)
		ebitenutil.DebugPrint(screen, fmt.Sprintf("%f, %f, %f, %f", lt.X+ox, lt.Y+oy, rb.X+ox, rb.Y+oy))
	}

	op := &ebiten.DrawTrianglesOptions{
		FillRule: ebiten.EvenOdd,
	}

	vs, is := path.AppendVerticesAndIndicesForFilling(nil, nil)
	for i := range vs {
		vs[i].SrcX = 1
		vs[i].SrcY = 1
		vs[i].ColorR = 0xdb / float32(0xff)
		vs[i].ColorG = 0x56 / float32(0xff)
		vs[i].ColorB = 0x20 / float32(0xff)
	}
	screen.DrawTriangles(vs, is, emptySubImage, op)

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {

	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Animation (Ebiten Demo)")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
