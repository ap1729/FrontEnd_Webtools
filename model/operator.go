package model

type Operator struct {
	id uint
}

// Getter methods for all Operator properties that should be read-only
func (op *Operator) ID() uint {
	return op.id
}

// "Constructor" for a new operator
func NewOperator(id uint) *Operator {
	return &Operator{id: id}
}
