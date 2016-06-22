package model

// User is a taggable structure that stores location, operator and other properties.
type User struct {
	id        uint
	x, y, ht  float64
	defaultOp *Operator

	CurrOp *Operator
}

// Getter methods for all User properties that are not exported (read-only).

// The ID.
func (ue *User) ID() uint {
	return ue.id
}

// The X-coordinate.
func (ue *User) X() float64 {
	return ue.x
}

// The Y-coordinate.
func (ue *User) Y() float64 {
	return ue.y
}

// The height.
func (ue *User) Ht() float64 {
	return ue.ht
}

// The registered operator.
func (ue *User) DefaultOp() *Operator {
	return ue.defaultOp
}

// Constructor to instantiate a User. The constructor must be used to create new objects as all properties are read-only.
func NewUser(id uint, x, y, ht float64, op *Operator) *User {
	return &User{id: id, x: x, y: y, ht: ht, defaultOp: op, CurrOp: op}
}
