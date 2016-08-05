package model

// User is a taggable structure that stores location, operator and other properties.
type User struct {
	id        uint
	x, y, ht  float64
	defaultOp *Operator
    bs0 *BaseStation //basestation in level0
    bs1 *BaseStation //basestation in level1
	// A public property for setting the current operator. This property has no internal implications,
	// it is present as for convenience and encapsulation.
	ConnectedBs *BaseStation //basestation it is connected to 
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

//bs level 0
func (ue *User) BS0() *BaseStation {
	return ue.bs0
}

//bs level 1
func (ue *User) BS1() *BaseStation {
	return ue.bs1
}

// Constructor to instantiate a User. The constructor must be used to create new objects as all properties are read-only.
func NewUser(id uint, x, y, ht float64, op *Operator,bs0 *BaseStation) *User {
	return &User{id: id, x: x, y: y, ht: ht, defaultOp: op, CurrOp: op,bs0: bs0}
}
