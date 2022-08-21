package objects

import (
	b2d "github.com/E4/box2d"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"image/color"
	"math"
)

var (
	emptyImage    = ebiten.NewImage(3, 3)
	emptySubImage = emptyImage.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)
)

func init() {
	emptyImage.Fill(color.White)
}

type Vec struct {
	X float32
	Y float32
}

func NewVec(x, y float64) Vec {
	return Vec{float32(x), float32(y)}
}

// RotateVec 回転させる
func RotateVec(v b2d.B2Vec2, radian float64) Vec {
	//x2 = x1 * cos(α) - y1 * sin(α)
	//y2 = x1 * sin(α) + y1 * cos(α)
	return NewVec(
		v.X*math.Cos(radian)-v.Y*math.Sin(radian),
		v.X*math.Sin(radian)+v.Y*math.Cos(radian),
	)
}

// 1メートル 100px にスケールさせる
const SCALE float32 = 100

type ObjectBase struct {
	ID int
}

var id = 0

func NewObjectBase() ObjectBase {
	r := ObjectBase{
		ID: id,
	}
	id += 1
	return r
}
