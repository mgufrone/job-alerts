package notification

import "go.uber.org/fx"

var MockModule = fx.Provide(
	func() *MockQueryRepository {
		return &MockQueryRepository{}
	},
	func() *MockCommandRepository {
		return &MockCommandRepository{}
	},
	func(source *MockQueryRepository) QueryResolver {
		return func() IQueryRepository {
			return source
		}
	},
	func(source *MockCommandRepository) CommandResolver {
		return func() ICommandRepository {
			return source
		}
	},
)
