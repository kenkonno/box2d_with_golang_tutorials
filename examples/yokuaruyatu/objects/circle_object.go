package objects

import (
	b2d "github.com/E4/box2d"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"math"
)

type CircleObject struct {
	Shape   *b2d.B2CircleShape
	Body    *b2d.B2Body
	Fixture *b2d.B2Fixture
	r       float32
	g       float32
	b       float32
}

// TODO: 色々オプションを追加する必要がありそう。いったん静的なポリゴンオブジェクトってことで

func NewCircleObject(px, py, r float64, world *b2d.B2World) CircleObject {

	bodyDef := b2d.NewB2BodyDef()
	bodyDef.Position.Set(px, py)

	body := world.CreateBody(bodyDef)

	circleShape := b2d.NewB2CircleShape()
	circleShape.SetRadius(r)

	fixture := body.CreateFixture(circleShape, 1.0)

	return CircleObject{
		Shape:   circleShape,
		Body:    body,
		Fixture: fixture,
		r:       0xdb,
		g:       0x56,
		b:       0x20,
	}
}
func NewDynamicCircleObject(px, py, r float64, world *b2d.B2World) CircleObject {
	result := NewCircleObject(px, py, r, world)
	result.Body.SetType(b2d.B2BodyType.B2_dynamicBody)
	result.Fixture.SetDensity(1.0)     // 質量
	result.Fixture.SetRestitution(0.4) // 反発係数
	return result
}

// Draw ebiten 向けの画像描画
func (r *CircleObject) Draw(screen *ebiten.Image) {
	// Get the angle in radians.
	radian := r.Body.GetAngle()
	position := r.Body.GetPosition()
	vec := NewVec(position.X, position.Y)

	var path vector.Path

	path.Arc(vec.X*SCALE, vec.Y*SCALE, float32(r.Shape.GetRadius())*SCALE, float32(radian), float32(radian)+2*math.Pi, vector.Clockwise)
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

func (r *CircleObject) SetRGB(_r, g, b float32) {
	r.r = _r
	r.g = g
	r.b = b
}
