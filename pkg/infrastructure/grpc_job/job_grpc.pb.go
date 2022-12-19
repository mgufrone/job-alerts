// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package grpc_job

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// JobServiceClient is the client API for JobService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type JobServiceClient interface {
	// read
	GetAll(ctx context.Context, in *JobFilter, opts ...grpc.CallOption) (*JobData, error)
	Count(ctx context.Context, in *JobFilter, opts ...grpc.CallOption) (*JobCountResult, error)
	GetAndCount(ctx context.Context, in *JobFilter, opts ...grpc.CallOption) (*JobCount, error)
	// write
	Create(ctx context.Context, in *Job, opts ...grpc.CallOption) (*Job, error)
	Update(ctx context.Context, in *Job, opts ...grpc.CallOption) (*Job, error)
	Delete(ctx context.Context, in *Job, opts ...grpc.CallOption) (*JobResult, error)
}

type jobServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewJobServiceClient(cc grpc.ClientConnInterface) JobServiceClient {
	return &jobServiceClient{cc}
}

func (c *jobServiceClient) GetAll(ctx context.Context, in *JobFilter, opts ...grpc.CallOption) (*JobData, error) {
	out := new(JobData)
	err := c.cc.Invoke(ctx, "/service.JobService/GetAll", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *jobServiceClient) Count(ctx context.Context, in *JobFilter, opts ...grpc.CallOption) (*JobCountResult, error) {
	out := new(JobCountResult)
	err := c.cc.Invoke(ctx, "/service.JobService/Count", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *jobServiceClient) GetAndCount(ctx context.Context, in *JobFilter, opts ...grpc.CallOption) (*JobCount, error) {
	out := new(JobCount)
	err := c.cc.Invoke(ctx, "/service.JobService/GetAndCount", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *jobServiceClient) Create(ctx context.Context, in *Job, opts ...grpc.CallOption) (*Job, error) {
	out := new(Job)
	err := c.cc.Invoke(ctx, "/service.JobService/Create", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *jobServiceClient) Update(ctx context.Context, in *Job, opts ...grpc.CallOption) (*Job, error) {
	out := new(Job)
	err := c.cc.Invoke(ctx, "/service.JobService/Update", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *jobServiceClient) Delete(ctx context.Context, in *Job, opts ...grpc.CallOption) (*JobResult, error) {
	out := new(JobResult)
	err := c.cc.Invoke(ctx, "/service.JobService/Delete", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// JobServiceServer is the server API for JobService service.
// All implementations must embed UnimplementedJobServiceServer
// for forward compatibility
type JobServiceServer interface {
	// read
	GetAll(context.Context, *JobFilter) (*JobData, error)
	Count(context.Context, *JobFilter) (*JobCountResult, error)
	GetAndCount(context.Context, *JobFilter) (*JobCount, error)
	// write
	Create(context.Context, *Job) (*Job, error)
	Update(context.Context, *Job) (*Job, error)
	Delete(context.Context, *Job) (*JobResult, error)
	mustEmbedUnimplementedJobServiceServer()
}

// UnimplementedJobServiceServer must be embedded to have forward compatible implementations.
type UnimplementedJobServiceServer struct {
}

func (UnimplementedJobServiceServer) GetAll(context.Context, *JobFilter) (*JobData, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAll not implemented")
}
func (UnimplementedJobServiceServer) Count(context.Context, *JobFilter) (*JobCountResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Count not implemented")
}
func (UnimplementedJobServiceServer) GetAndCount(context.Context, *JobFilter) (*JobCount, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAndCount not implemented")
}
func (UnimplementedJobServiceServer) Create(context.Context, *Job) (*Job, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (UnimplementedJobServiceServer) Update(context.Context, *Job) (*Job, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}
func (UnimplementedJobServiceServer) Delete(context.Context, *Job) (*JobResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
func (UnimplementedJobServiceServer) mustEmbedUnimplementedJobServiceServer() {}

// UnsafeJobServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to JobServiceServer will
// result in compilation errors.
type UnsafeJobServiceServer interface {
	mustEmbedUnimplementedJobServiceServer()
}

func RegisterJobServiceServer(s grpc.ServiceRegistrar, srv JobServiceServer) {
	s.RegisterService(&JobService_ServiceDesc, srv)
}

func _JobService_GetAll_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(JobFilter)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(JobServiceServer).GetAll(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/service.JobService/GetAll",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(JobServiceServer).GetAll(ctx, req.(*JobFilter))
	}
	return interceptor(ctx, in, info, handler)
}

func _JobService_Count_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(JobFilter)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(JobServiceServer).Count(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/service.JobService/Count",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(JobServiceServer).Count(ctx, req.(*JobFilter))
	}
	return interceptor(ctx, in, info, handler)
}

func _JobService_GetAndCount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(JobFilter)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(JobServiceServer).GetAndCount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/service.JobService/GetAndCount",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(JobServiceServer).GetAndCount(ctx, req.(*JobFilter))
	}
	return interceptor(ctx, in, info, handler)
}

func _JobService_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Job)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(JobServiceServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/service.JobService/Create",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(JobServiceServer).Create(ctx, req.(*Job))
	}
	return interceptor(ctx, in, info, handler)
}

func _JobService_Update_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Job)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(JobServiceServer).Update(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/service.JobService/Update",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(JobServiceServer).Update(ctx, req.(*Job))
	}
	return interceptor(ctx, in, info, handler)
}

func _JobService_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Job)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(JobServiceServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/service.JobService/Delete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(JobServiceServer).Delete(ctx, req.(*Job))
	}
	return interceptor(ctx, in, info, handler)
}

// JobService_ServiceDesc is the grpc.ServiceDesc for JobService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var JobService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "service.JobService",
	HandlerType: (*JobServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetAll",
			Handler:    _JobService_GetAll_Handler,
		},
		{
			MethodName: "Count",
			Handler:    _JobService_Count_Handler,
		},
		{
			MethodName: "GetAndCount",
			Handler:    _JobService_GetAndCount_Handler,
		},
		{
			MethodName: "Create",
			Handler:    _JobService_Create_Handler,
		},
		{
			MethodName: "Update",
			Handler:    _JobService_Update_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _JobService_Delete_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "job.proto",
}
