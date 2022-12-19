package grpc

import (
	"context"
	"google.golang.org/grpc/health/grpc_health_v1"
)

type healthService struct {
}

func (h healthService) Check(ctx context.Context, request *grpc_health_v1.HealthCheckRequest) (res *grpc_health_v1.HealthCheckResponse, err error) {
	res = &grpc_health_v1.HealthCheckResponse{Status: grpc_health_v1.HealthCheckResponse_SERVING}
	return
}

func (h healthService) Watch(request *grpc_health_v1.HealthCheckRequest, server grpc_health_v1.Health_WatchServer) error {
	res := &grpc_health_v1.HealthCheckResponse{Status: grpc_health_v1.HealthCheckResponse_SERVING}
	return server.Send(res)
}

func NewHealthcheck() grpc_health_v1.HealthServer {
	return &healthService{}
}
