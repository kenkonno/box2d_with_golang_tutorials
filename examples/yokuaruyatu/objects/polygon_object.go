package objects

import (
	b2d "github.com/E4/box2d"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type PolygonObject struct {
	Shape   *b2d.B2PolygonShape
	Body    *b2d.B2Body
	Fixture *b2d.B2Fixture
	r       float32
	g       float32
	b       float32
}

// TODO: 色々オプションを追加する必要がありそう。いったん静的なポリゴンオブジェクトってことで

func NewPolygonObject(px, py float64, vertices []b2d.B2Vec2, angle, density float64, world *b2d.B2World) PolygonObject {

	bodyDef := b2d.NewB2BodyDef()
	bodyDef.Position.Set(px, py)
	bodyDef.Angle = angle

	body := world.CreateBody(bodyDef)
	body.SetUserData(NewObjectBase())

	polygon := b2d.NewB2PolygonShape()
	polygon.Set(vertices, len(vertices))

	fixture := body.CreateFixture(polygon, density)

	return PolygonObject{
		Shape:   polygon,
		Body:    body,
		Fixture: fixture,
		r:       0xdb,
		g:       0x56,
		b:       0x20,
	}
}

// NewPolygonBox 静的な箱を作る
func NewPolygonBox(px, py, halfWidth, halfHeight float64, world *b2d.B2World, angle float64) PolygonObject {
	vertices := createBox(halfWidth, halfHeight)
	return NewPolygonObject(px, py, vertices, angle, 0, world)
}

func NewDynamicBox(px, py, halfWidth, halfHeight, density float64, world *b2d.B2World) PolygonObject {
	vertices := createBox(halfWidth, halfHeight)
	result := NewPolygonObject(px, py, vertices, 0, density, world)
	result.Body.SetType(b2d.B2BodyType.B2_dynamicBody)
	result.Fixture.SetRestitution(0.4) // 反発係数
	return result
}

func createBox(halfWidth, halfHeight float64) []b2d.B2Vec2 {
	return []b2d.B2Vec2{
		{-halfWidth, -halfHeight},
		{-halfWidth, +halfHeight},
		{+halfWidth, +halfHeight},
		{+halfWidth, -halfHeight},
	}
}

// Draw ebiten 向けの画像描画
func (r *PolygonObject) Draw(screen *ebiten.Image) {
	// Get the angle in radians.
	rad := r.Body.GetAngle()
	var vec []Vec
	for i := 0; i < r.Shape.M_count; i++ {
		v := r.Shape.M_vertices[i]
		vec = append(vec, RotateVec(v, rad))
	}
	var path vector.Path
	centerX := r.Body.GetPosition().X * SCALE
	centerY := r.Body.GetPosition().Y * SCALE
	for i, v := range vec {
		if i == 0 {
			path.MoveTo(float32(centerX)+v.X*SCALE32, float32(centerY)+v.Y*SCALE32)
		} else {
			path.LineTo(float32(centerX)+v.X*SCALE32, float32(centerY)+v.Y*SCALE32)
		}
	}
	// TODO: こっから下はいまいち思い出せていない。 path を三角形にして描画設定をしているように記憶している。
	op := &ebiten.DrawTrianglesOptions{
		FillRule: ebiten.EvenOdd,
	}
	vs, is := path.AppendVerticesAndIndicesForFilling(nil, nil)
	for i := range vs {
		vs[i].SrcX = 1
		vs[i].SrcY = 1
		vs[i].ColorR = r.r / float32(0xff)
		vs[i].ColorG = r.g / float32(0xff)
		vs[i].ColorB = r.b / float32(0xff)
	}
	screen.DrawTriangles(vs, is, emptySubImage, op)
}

func (r *PolygonObject) SetRGB(_r, g, b float32) {
	r.r = _r
	r.g = g
	r.b = b
}
