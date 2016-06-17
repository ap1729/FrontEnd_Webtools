package service

import "math"

// Hexagon models a single cell site and provides useful functionality.
type Hexagon struct {
	CenterX, CenterY, Side float64
	ID                     uint
}

/*
// Getter methods for all Hexagon properties that should be "private set, public read"
func (hex *Hexagon) CenterX() float64 {
	return hex.centerX
}
func (hex *Hexagon) CenterY() float64 {
	return hex.centerY
}
func (hex *Hexagon) Side() float64 {
	return hex.side
}
func (hex *Hexagon) ID() uint {
	return hex.id
}
*/

// Function to check if a point (x, y) lies inside or on the boundary of the Hexagon.
func (hex *Hexagon) Contains(x, y float64) bool {
	relx := x - hex.CenterX
	rely := y - hex.CenterY
	nLen := hex.Side * math.Cos(math.Pi/6)
	if math.Abs(relx*math.Cos(0)-rely*math.Sin(0)) <= nLen && math.Abs(relx*math.Cos(math.Pi/3)-rely*math.Sin(math.Pi/3)) <= nLen && math.Abs(relx*math.Cos(-math.Pi/3)-rely*math.Sin(-math.Pi/3)) <= nLen {
		return true
	} else {
		return false
	}
}
