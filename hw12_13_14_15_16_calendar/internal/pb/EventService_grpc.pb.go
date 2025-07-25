// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.12.4
// source: EventService.proto

package pb

import (
	context "context"
	empty "github.com/golang/protobuf/ptypes/empty"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	EventService_CreateEvent_FullMethodName = "/event.EventService/CreateEvent"
	EventService_UpdateEvent_FullMethodName = "/event.EventService/UpdateEvent"
	EventService_DeleteEvent_FullMethodName = "/event.EventService/DeleteEvent"
	EventService_ListDay_FullMethodName     = "/event.EventService/ListDay"
	EventService_ListWeek_FullMethodName    = "/event.EventService/ListWeek"
	EventService_ListMonth_FullMethodName   = "/event.EventService/ListMonth"
)

// EventServiceClient is the client API for EventService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// ==== Service ==============================================================
type EventServiceClient interface {
	CreateEvent(ctx context.Context, in *CreateEventRequest, opts ...grpc.CallOption) (*EventResponse, error)
	UpdateEvent(ctx context.Context, in *UpdateEventRequest, opts ...grpc.CallOption) (*EventResponse, error)
	DeleteEvent(ctx context.Context, in *DeleteEventRequest, opts ...grpc.CallOption) (*empty.Empty, error)
	ListDay(ctx context.Context, in *ListDayRequest, opts ...grpc.CallOption) (*EventsResponse, error)
	ListWeek(ctx context.Context, in *ListWeekRequest, opts ...grpc.CallOption) (*EventsResponse, error)
	ListMonth(ctx context.Context, in *ListMonthRequest, opts ...grpc.CallOption) (*EventsResponse, error)
}

type eventServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewEventServiceClient(cc grpc.ClientConnInterface) EventServiceClient {
	return &eventServiceClient{cc}
}

func (c *eventServiceClient) CreateEvent(ctx context.Context, in *CreateEventRequest, opts ...grpc.CallOption) (*EventResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(EventResponse)
	err := c.cc.Invoke(ctx, EventService_CreateEvent_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *eventServiceClient) UpdateEvent(ctx context.Context, in *UpdateEventRequest, opts ...grpc.CallOption) (*EventResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(EventResponse)
	err := c.cc.Invoke(ctx, EventService_UpdateEvent_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *eventServiceClient) DeleteEvent(ctx context.Context, in *DeleteEventRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, EventService_DeleteEvent_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *eventServiceClient) ListDay(ctx context.Context, in *ListDayRequest, opts ...grpc.CallOption) (*EventsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(EventsResponse)
	err := c.cc.Invoke(ctx, EventService_ListDay_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *eventServiceClient) ListWeek(ctx context.Context, in *ListWeekRequest, opts ...grpc.CallOption) (*EventsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(EventsResponse)
	err := c.cc.Invoke(ctx, EventService_ListWeek_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *eventServiceClient) ListMonth(ctx context.Context, in *ListMonthRequest, opts ...grpc.CallOption) (*EventsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(EventsResponse)
	err := c.cc.Invoke(ctx, EventService_ListMonth_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// EventServiceServer is the server API for EventService service.
// All implementations must embed UnimplementedEventServiceServer
// for forward compatibility.
//
// ==== Service ==============================================================
type EventServiceServer interface {
	CreateEvent(context.Context, *CreateEventRequest) (*EventResponse, error)
	UpdateEvent(context.Context, *UpdateEventRequest) (*EventResponse, error)
	DeleteEvent(context.Context, *DeleteEventRequest) (*empty.Empty, error)
	ListDay(context.Context, *ListDayRequest) (*EventsResponse, error)
	ListWeek(context.Context, *ListWeekRequest) (*EventsResponse, error)
	ListMonth(context.Context, *ListMonthRequest) (*EventsResponse, error)
	mustEmbedUnimplementedEventServiceServer()
}

// UnimplementedEventServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedEventServiceServer struct{}

func (UnimplementedEventServiceServer) CreateEvent(context.Context, *CreateEventRequest) (*EventResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateEvent not implemented")
}
func (UnimplementedEventServiceServer) UpdateEvent(context.Context, *UpdateEventRequest) (*EventResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateEvent not implemented")
}
func (UnimplementedEventServiceServer) DeleteEvent(context.Context, *DeleteEventRequest) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteEvent not implemented")
}
func (UnimplementedEventServiceServer) ListDay(context.Context, *ListDayRequest) (*EventsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListDay not implemented")
}
func (UnimplementedEventServiceServer) ListWeek(context.Context, *ListWeekRequest) (*EventsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListWeek not implemented")
}
func (UnimplementedEventServiceServer) ListMonth(context.Context, *ListMonthRequest) (*EventsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListMonth not implemented")
}
func (UnimplementedEventServiceServer) mustEmbedUnimplementedEventServiceServer() {}
func (UnimplementedEventServiceServer) testEmbeddedByValue()                      {}

// UnsafeEventServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to EventServiceServer will
// result in compilation errors.
type UnsafeEventServiceServer interface {
	mustEmbedUnimplementedEventServiceServer()
}

func RegisterEventServiceServer(s grpc.ServiceRegistrar, srv EventServiceServer) {
	// If the following call pancis, it indicates UnimplementedEventServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&EventService_ServiceDesc, srv)
}

func _EventService_CreateEvent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateEventRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EventServiceServer).CreateEvent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: EventService_CreateEvent_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EventServiceServer).CreateEvent(ctx, req.(*CreateEventRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EventService_UpdateEvent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateEventRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EventServiceServer).UpdateEvent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: EventService_UpdateEvent_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EventServiceServer).UpdateEvent(ctx, req.(*UpdateEventRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EventService_DeleteEvent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteEventRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EventServiceServer).DeleteEvent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: EventService_DeleteEvent_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EventServiceServer).DeleteEvent(ctx, req.(*DeleteEventRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EventService_ListDay_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListDayRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EventServiceServer).ListDay(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: EventService_ListDay_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EventServiceServer).ListDay(ctx, req.(*ListDayRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EventService_ListWeek_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListWeekRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EventServiceServer).ListWeek(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: EventService_ListWeek_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EventServiceServer).ListWeek(ctx, req.(*ListWeekRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EventService_ListMonth_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListMonthRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EventServiceServer).ListMonth(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: EventService_ListMonth_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EventServiceServer).ListMonth(ctx, req.(*ListMonthRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// EventService_ServiceDesc is the grpc.ServiceDesc for EventService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var EventService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "event.EventService",
	HandlerType: (*EventServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateEvent",
			Handler:    _EventService_CreateEvent_Handler,
		},
		{
			MethodName: "UpdateEvent",
			Handler:    _EventService_UpdateEvent_Handler,
		},
		{
			MethodName: "DeleteEvent",
			Handler:    _EventService_DeleteEvent_Handler,
		},
		{
			MethodName: "ListDay",
			Handler:    _EventService_ListDay_Handler,
		},
		{
			MethodName: "ListWeek",
			Handler:    _EventService_ListWeek_Handler,
		},
		{
			MethodName: "ListMonth",
			Handler:    _EventService_ListMonth_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "EventService.proto",
}
