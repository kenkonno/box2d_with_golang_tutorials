# インクレディブルマシーンのようなゲームを作る

# 手版

## 物体の作り方
1. bodyDefを作る以上 NewB2BodyDef
2. BoxShapeを作る NewB2PolygonShape
3. BodyにShapeをくっつける body.CreateFixture

## メートル to Pixel
1メートルを100px として描画する

## 描画クラス
 - Polygon
   - NewB2PolygonShape の場合は Shape.M_vertices がたぶん頂点の集合

 - Circle
   - Circle の場合は 半径と角度だけだと思う

