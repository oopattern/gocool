// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package pb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion7

// ObserveClient is the client API for Observe service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ObserveClient interface {
	SayRoute(ctx context.Context, in *RouteReq, opts ...grpc.CallOption) (*RouteResp, error)
}

type observeClient struct {
	cc grpc.ClientConnInterface
}

func NewObserveClient(cc grpc.ClientConnInterface) ObserveClient {
	return &observeClient{cc}
}

func (c *observeClient) SayRoute(ctx context.Context, in *RouteReq, opts ...grpc.CallOption) (*RouteResp, error) {
	out := new(RouteResp)
	err := c.cc.Invoke(ctx, "/pb.Observe/SayRoute", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ObserveServer is the server API for Observe service.
// All implementations must embed UnimplementedObserveServer
// for forward compatibility
type ObserveServer interface {
	SayRoute(context.Context, *RouteReq) (*RouteResp, error)
	mustEmbedUnimplementedObserveServer()
}

// UnimplementedObserveServer must be embedded to have forward compatible implementations.
type UnimplementedObserveServer struct {
}

func (UnimplementedObserveServer) SayRoute(context.Context, *RouteReq) (*RouteResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SayRoute not implemented")
}
func (UnimplementedObserveServer) mustEmbedUnimplementedObserveServer() {}

// UnsafeObserveServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ObserveServer will
// result in compilation errors.
type UnsafeObserveServer interface {
	mustEmbedUnimplementedObserveServer()
}

func RegisterObserveServer(s grpc.ServiceRegistrar, srv ObserveServer) {
	s.RegisterService(&_Observe_serviceDesc, srv)
}

func _Observe_SayRoute_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RouteReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ObserveServer).SayRoute(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.Observe/SayRoute",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ObserveServer).SayRoute(ctx, req.(*RouteReq))
	}
	return interceptor(ctx, in, info, handler)
}

var _Observe_serviceDesc = grpc.ServiceDesc{
	ServiceName: "pb.Observe",
	HandlerType: (*ObserveServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SayRoute",
			Handler:    _Observe_SayRoute_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "observe.proto",
}
