package criteria

import "fmt"

const (
	Eq Operator = iota
	Not
	Like
	NotLike
	In
	NotIn
	Gt
	Gte
	Lt
	Lte
	Between
)

type ICondition interface {
	Field() string
	Operator() Operator
	Value() interface{}
	ToString() string
}
type baseCondition struct {
	field    string
	operator Operator
	value    interface{}
}

func NewCondition(field string, operator Operator, value interface{}) ICondition {
	if value == nil {
		return nil
	}
	return &baseCondition{field: field, operator: operator, value: value}
}

func (b *baseCondition) Field() string {
	return b.field
}

func (b *baseCondition) Operator() Operator {
	return b.operator
}

func (b *baseCondition) Value() interface{} {
	return b.value
}

func (b *baseCondition) ToString() string {
	return fmt.Sprintf("%s-%d=%v", b.field, b.operator, b.value)
}
