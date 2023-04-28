package event

import "go.uber.org/fx"

func Load(managerType string) fx.Option {
	var c interface{}
	switch managerType {
	case "serial":
		c = newSerial
	case "async":
		c = newAsyncUnsafe
	default:
		c = NewSimpleManager
	}
	return fx.Options(fx.Provide(
		fx.Annotate(
			c,
			fx.As(new(IManager)),
		),
	),
		fx.Invoke(
			fx.Annotate(
				func(subs []Handler, mgr IManager) {
					for _, sub := range subs {
						for _, evtType := range sub.SubscribeAt() {
							mgr.Subscribe(evtType, sub)
						}
					}
				},
				fx.ParamTags(`group:"event_subscribers"`),
			),
		),
	)
}

var Module = fx.Options(
	fx.Provide(
		fx.Annotate(
			newSerial,
			fx.As(new(IManager)),
		),
	),
	fx.Invoke(
		fx.Annotate(
			func(subs []Handler, mgr IManager) {
				for _, sub := range subs {
					for _, evtType := range sub.SubscribeAt() {
						mgr.Subscribe(evtType, sub)
					}
				}
			},
			fx.ParamTags(`group:"event_subscribers"`),
		),
	),
)
