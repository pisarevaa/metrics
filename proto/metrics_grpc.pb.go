// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.6.1
// source: proto/metrics.proto

package proto

import (
	context "context"

	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	Metrics_AddMetrics_FullMethodName = "/metrics.Metrics/AddMetrics"
	Metrics_GetMetric_FullMethodName  = "/metrics.Metrics/GetMetric"
	Metrics_GetMetrics_FullMethodName = "/metrics.Metrics/GetMetrics"
)

// MetricsClient is the client API for Metrics service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MetricsClient interface {
	AddMetrics(ctx context.Context, in *AddMetricsRequest, opts ...grpc.CallOption) (*AddMetricsResponse, error)
	GetMetric(ctx context.Context, in *GetMetricRequest, opts ...grpc.CallOption) (*GetMetricResponse, error)
	GetMetrics(ctx context.Context, in *GetMetricsRequest, opts ...grpc.CallOption) (*GetMetricsResponse, error)
}

type metricsClient struct {
	cc grpc.ClientConnInterface
}

func NewMetricsClient(cc grpc.ClientConnInterface) MetricsClient {
	return &metricsClient{cc}
}

func (c *metricsClient) AddMetrics(ctx context.Context, in *AddMetricsRequest, opts ...grpc.CallOption) (*AddMetricsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(AddMetricsResponse)
	err := c.cc.Invoke(ctx, Metrics_AddMetrics_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *metricsClient) GetMetric(ctx context.Context, in *GetMetricRequest, opts ...grpc.CallOption) (*GetMetricResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetMetricResponse)
	err := c.cc.Invoke(ctx, Metrics_GetMetric_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *metricsClient) GetMetrics(ctx context.Context, in *GetMetricsRequest, opts ...grpc.CallOption) (*GetMetricsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetMetricsResponse)
	err := c.cc.Invoke(ctx, Metrics_GetMetrics_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MetricsServer is the server API for Metrics service.
// All implementations must embed UnimplementedMetricsServer
// for forward compatibility.
type MetricsServer interface {
	AddMetrics(context.Context, *AddMetricsRequest) (*AddMetricsResponse, error)
	GetMetric(context.Context, *GetMetricRequest) (*GetMetricResponse, error)
	GetMetrics(context.Context, *GetMetricsRequest) (*GetMetricsResponse, error)
	mustEmbedUnimplementedMetricsServer()
}

// UnimplementedMetricsServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedMetricsServer struct{}

func (UnimplementedMetricsServer) AddMetrics(context.Context, *AddMetricsRequest) (*AddMetricsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddMetrics not implemented")
}
func (UnimplementedMetricsServer) GetMetric(context.Context, *GetMetricRequest) (*GetMetricResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMetric not implemented")
}
func (UnimplementedMetricsServer) GetMetrics(context.Context, *GetMetricsRequest) (*GetMetricsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMetrics not implemented")
}
func (UnimplementedMetricsServer) mustEmbedUnimplementedMetricsServer() {}
func (UnimplementedMetricsServer) testEmbeddedByValue()                 {}

// UnsafeMetricsServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MetricsServer will
// result in compilation errors.
type UnsafeMetricsServer interface {
	mustEmbedUnimplementedMetricsServer()
}

func RegisterMetricsServer(s grpc.ServiceRegistrar, srv MetricsServer) {
	// If the following call pancis, it indicates UnimplementedMetricsServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Metrics_ServiceDesc, srv)
}

func _Metrics_AddMetrics_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddMetricsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetricsServer).AddMetrics(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Metrics_AddMetrics_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetricsServer).AddMetrics(ctx, req.(*AddMetricsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Metrics_GetMetric_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetMetricRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetricsServer).GetMetric(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Metrics_GetMetric_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetricsServer).GetMetric(ctx, req.(*GetMetricRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Metrics_GetMetrics_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetMetricsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetricsServer).GetMetrics(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Metrics_GetMetrics_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetricsServer).GetMetrics(ctx, req.(*GetMetricsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Metrics_ServiceDesc is the grpc.ServiceDesc for Metrics service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Metrics_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "metrics.Metrics",
	HandlerType: (*MetricsServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AddMetrics",
			Handler:    _Metrics_AddMetrics_Handler,
		},
		{
			MethodName: "GetMetric",
			Handler:    _Metrics_GetMetric_Handler,
		},
		{
			MethodName: "GetMetrics",
			Handler:    _Metrics_GetMetrics_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/metrics.proto",
}