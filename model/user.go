package model

type User struct {
	id        uint
	x, y, ht  float64
	defaultOp *Operator
}

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

func NewUser(id uint, x, y, ht float64, op *Operator) *User {
	return &User{id: id, x: x, y: y, ht: ht, defaultOp: op}
}
