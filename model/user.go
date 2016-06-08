package model

type User struct {
	id        uint
	x, y, ht  float64
	defaultOp *Operator
}

// Getter methods for all User properties that should be read-only
func (ue *User) ID() uint {
	return ue.id
}
func (ue *User) X() float64 {
	return ue.x
}
func (ue *User) Y() float64 {
	return ue.y
}
func (ue *User) Ht() float64 {
	return ue.ht
}
func (ue *User) DefaultOp() *Operator {
	return ue.defaultOp
}

// "Constructor" for User object
func NewUser(id uint, x, y, ht float64, op *Operator) *User {
	return &User{id: id, x: x, y: y, ht: ht, defaultOp: op}
}
