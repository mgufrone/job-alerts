package criteria

type Operator int

type ICriteriaBuilder interface {
	Copy() ICriteriaBuilder
	Select(fields ...string) ICriteriaBuilder
	Paginate(page int, perPage int) ICriteriaBuilder
	Order(field string, direction string) ICriteriaBuilder
	// by default, it will run ands
	Where(condition ...ICondition) ICriteriaBuilder
	And(other ...ICriteriaBuilder) ICriteriaBuilder
	Or(other ...ICriteriaBuilder) ICriteriaBuilder
	ToString() string
}
