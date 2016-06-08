package model

type BaseStation struct {
	id       uint
	x, y, ht float64
	ownerOp  *Operator
}

// Getter methods for all BaseStation properties that should be read-only
func (bs *BaseStation) ID() uint {
	return bs.id
}
func (bs *BaseStation) X() float64 {
	return bs.x
}
func (bs *BaseStation) Y() float64 {
	return bs.y
}
func (bs *BaseStation) Ht() float64 {
	return bs.ht
}
func (bs *BaseStation) OwnerOp() *Operator {
	return bs.ownerOp
}

// "Constructor" for BaseStation object
func NewBaseStation(id uint, x, y, ht float64, op *Operator) *BaseStation {
	return &BaseStation{id: id, x: x, y: y, ht: ht, ownerOp: op}
}
