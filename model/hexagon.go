package model

import "math"

type Hexagon struct {
	centerX, centerY, side float64
	id                     uint
}

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

// Function to check if a point lies inside the Hexagon
func (hex *Hexagon) Contains(x, y float64) bool {
	relx := x - hex.centerX
	rely := y - hex.centerY
	if math.Abs(relx*math.Cos(0)-rely*math.Sin(0)) < 500.00 && math.Abs(relx*math.Cos(math.Pi/3)-rely*math.Sin(math.Pi/3)) < 500.00 && math.Abs(relx*math.Cos(-math.Pi/3)-rely*math.Sin(-math.Pi/3)) < 500.00 {
		return true
	} else {
		return false
	}
}
