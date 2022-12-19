package criteria

import "github.com/stretchr/testify/mock"

type MockCriteria struct {
	mock.Mock
}

func (m *MockCriteria) Copy() ICriteriaBuilder {
	return m
}

func (m *MockCriteria) Select(fields ...string) ICriteriaBuilder {
	itr := make([]interface{}, len(fields))
	for i, f := range fields {
		itr[i] = f
	}

	m.Called(itr...)

	return m
}

func (m *MockCriteria) Paginate(page int, perPage int) ICriteriaBuilder {
	m.Called(page, perPage)

	return m
}

func (m *MockCriteria) Order(field string, direction string) ICriteriaBuilder {
	m.Called(field, direction)

	return m
}

func (m *MockCriteria) Where(condition ...ICondition) ICriteriaBuilder {
	itr := make([]interface{}, len(condition))
	for i, f := range condition {
		itr[i] = f
	}

	m.Called(itr...)

	return m
}

func (m *MockCriteria) composeCriteria(other ...ICriteriaBuilder) []interface{} {
	itr := make([]interface{}, len(other))
	for i, f := range other {
		itr[i] = f
	}
	return itr
}

func (m *MockCriteria) And(other ...ICriteriaBuilder) ICriteriaBuilder {
	m.Called(m.composeCriteria(other...))
	return m
}

func (m *MockCriteria) Or(other ...ICriteriaBuilder) ICriteriaBuilder {
	m.Called(m.composeCriteria(other...))
	return m
}

func (m *MockCriteria) ToString() string {
	return m.Called().Get(0).(string)
}
