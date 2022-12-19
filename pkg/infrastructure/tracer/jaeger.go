package tracer

import (
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"time"
)

func WithTracer(appName, appEnv string) func() (opentracing.Tracer, error) {
	return func() (opentracing.Tracer, error) {
		return NewTracer(appName, appEnv)
	}
}
func NewTracer(appName, appEnv string) (opentracing.Tracer, error) {
	cfg := &config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Disabled:    appEnv != "production",
		Tags:        []opentracing.Tag{{Key: "environment", Value: appEnv}},
		ServiceName: appName,
		Reporter: &config.ReporterConfig{
			LogSpans:                   true,
			DisableAttemptReconnecting: false,
			BufferFlushInterval:        1 * time.Second,
		},
	}
	cfg, err := cfg.FromEnv()
	if err != nil {
		fmt.Println("error initializing jaeger tracer", err)
		return nil, err
	}

	// jaeger tracer client
	tr, _, err := cfg.NewTracer(
		config.Logger(jaegerlog.DebugLogAdapter(jaegerlog.StdLogger)),
		config.ZipkinSharedRPCSpan(true),
	)
	return tr, err
}
