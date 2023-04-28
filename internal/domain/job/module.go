package job

import "go.uber.org/fx"

var MockModule = fx.Provide(
	func() *MockQueryRepository {
		return &MockQueryRepository{}
	},
	func(source *MockQueryRepository) QueryResolver {
		return func() IQueryRepository {
			return source
		}
	},
	func() *MockCommandRepository {
		return &MockCommandRepository{}
	},
	func(source *MockCommandRepository) CommandResolver {
		return func() ICommandRepository {
			return source
		}
	},
)
