package model

// PONDER: Enforce non-nil values in all reference types for safety and easier error handling in future?

// BaseStation is a taggable structure that stores location, operator and other properties.
type BaseStation struct {
	id       uint
	x, y, ht float64
	ownerOp  *Operator
	//which users are it connected to
	ConnectedUsers []*User //which users are the bs connected to
	Destroyed      uint
}

// Getter methods for all BaseStation properties that are not exported (read-only):

// The ID.
func (bs *BaseStation) ID() uint {
	return bs.id
}

// The X-coordinate.
func (bs *BaseStation) X() float64 {
	return bs.x
}

// The Y-coordinate.
func (bs *BaseStation) Y() float64 {
	return bs.y
}

// The height.
func (bs *BaseStation) Ht() float64 {
	return bs.ht
}

// The owning operator.
func (bs *BaseStation) OwnerOp() *Operator {
	return bs.ownerOp
}

// Constructor to instantiate a BaseStation. The constructor must be used to create new objects as all properties are read-only.
func NewBaseStation(id uint, x, y, ht float64, op *Operator) *BaseStation {
	return &BaseStation{id: id, x: x, y: y, ht: ht, ownerOp: op, Destroyed: 0}
}
