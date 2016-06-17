package model

// Operator stores details such as ID. This class can be extended in the future to add more attributes.
type Operator struct {
	id uint
}

// Getter methods for all Operator properties that are not exported (read-only):

// The ID.
func (op *Operator) ID() uint {
	return op.id
}

// Constructor to instantiate an Operator. The constructor must be used to create new objects as all properties are read-only.
func NewOperator(id uint) *Operator {
	return &Operator{id: id}
}
