package model

type Operator struct {
	id uint
}

func (op *Operator) ID() uint {
	return op.id
}

func NewOperator(id uint) *Operator {
	return &Operator{id: id}
}
