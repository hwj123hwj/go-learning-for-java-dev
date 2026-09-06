package demo

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DemoServiceClient interface {
	GetUser(ctx context.Context, in *GetUserRequest, opts ...grpc.CallOption) (*GetUserResponse, error)
	ListUsers(ctx context.Context, in *ListUsersRequest, opts ...grpc.CallOption) (DemoService_ListUsersClient, error)
	UploadUsers(ctx context.Context, opts ...grpc.CallOption) (DemoService_UploadUsersClient, error)
	Chat(ctx context.Context, opts ...grpc.CallOption) (DemoService_ChatClient, error)
}

type demoServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewDemoServiceClient(cc grpc.ClientConnInterface) DemoServiceClient {
	return &demoServiceClient{cc}
}

func (c *demoServiceClient) GetUser(ctx context.Context, in *GetUserRequest, opts ...grpc.CallOption) (*GetUserResponse, error) {
	out := new(GetUserResponse)
	err := c.cc.Invoke(ctx, "/demo.DemoService/GetUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *demoServiceClient) ListUsers(ctx context.Context, in *ListUsersRequest, opts ...grpc.CallOption) (DemoService_ListUsersClient, error) {
	stream, err := c.cc.NewStream(ctx, &_DemoService_serviceDesc.Streams[0], "/demo.DemoService/ListUsers", opts...)
	if err != nil {
		return nil, err
	}
	x := &demoServiceListUsersClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type DemoService_ListUsersClient interface {
	Recv() (*ListUsersResponse, error)
	grpc.ClientStream
}

type demoServiceListUsersClient struct {
	grpc.ClientStream
}

func (x *demoServiceListUsersClient) Recv() (*ListUsersResponse, error) {
	m := new(ListUsersResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *demoServiceClient) UploadUsers(ctx context.Context, opts ...grpc.CallOption) (DemoService_UploadUsersClient, error) {
	stream, err := c.cc.NewStream(ctx, &_DemoService_serviceDesc.Streams[1], "/demo.DemoService/UploadUsers", opts...)
	if err != nil {
		return nil, err
	}
	x := &demoServiceUploadUsersClient{stream}
	return x, nil
}

type DemoService_UploadUsersClient interface {
	Send(*UploadUsersRequest) error
	CloseAndRecv() (*UploadUsersResponse, error)
	grpc.ClientStream
}

type demoServiceUploadUsersClient struct {
	grpc.ClientStream
}

func (x *demoServiceUploadUsersClient) Send(m *UploadUsersRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *demoServiceUploadUsersClient) CloseAndRecv() (*UploadUsersResponse, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(UploadUsersResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *demoServiceClient) Chat(ctx context.Context, opts ...grpc.CallOption) (DemoService_ChatClient, error) {
	stream, err := c.cc.NewStream(ctx, &_DemoService_serviceDesc.Streams[2], "/demo.DemoService/Chat", opts...)
	if err != nil {
		return nil, err
	}
	x := &demoServiceChatClient{stream}
	return x, nil
}

type DemoService_ChatClient interface {
	Send(*ChatMessage) error
	Recv() (*ChatMessage, error)
	grpc.ClientStream
}

type demoServiceChatClient struct {
	grpc.ClientStream
}

func (x *demoServiceChatClient) Send(m *ChatMessage) error {
	return x.ClientStream.SendMsg(m)
}

func (x *demoServiceChatClient) Recv() (*ChatMessage, error) {
	m := new(ChatMessage)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

type DemoServiceServer interface {
	GetUser(context.Context, *GetUserRequest) (*GetUserResponse, error)
	ListUsers(*ListUsersRequest, DemoService_ListUsersServer) error
	UploadUsers(DemoService_UploadUsersServer) error
	Chat(DemoService_ChatServer) error
}

type UnimplementedDemoServiceServer struct{}

func (UnimplementedDemoServiceServer) GetUser(context.Context, *GetUserRequest) (*GetUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUser not implemented")
}

func (UnimplementedDemoServiceServer) ListUsers(*ListUsersRequest, DemoService_ListUsersServer) error {
	return status.Errorf(codes.Unimplemented, "method ListUsers not implemented")
}

func (UnimplementedDemoServiceServer) UploadUsers(DemoService_UploadUsersServer) error {
	return status.Errorf(codes.Unimplemented, "method UploadUsers not implemented")
}

func (UnimplementedDemoServiceServer) Chat(DemoService_ChatServer) error {
	return status.Errorf(codes.Unimplemented, "method Chat not implemented")
}

func RegisterDemoServiceServer(s grpc.ServiceRegistrar, srv DemoServiceServer) {
	s.RegisterService(&_DemoService_serviceDesc, srv)
}

func _DemoService_GetUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DemoServiceServer).GetUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/demo.DemoService/GetUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DemoServiceServer).GetUser(ctx, req.(*GetUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DemoService_ListUsers_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ListUsersRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(DemoServiceServer).ListUsers(m, &demoServiceListUsersServer{stream})
}

type DemoService_ListUsersServer interface {
	Send(*ListUsersResponse) error
	grpc.ServerStream
}

type demoServiceListUsersServer struct {
	grpc.ServerStream
}

func (x *demoServiceListUsersServer) Send(m *ListUsersResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _DemoService_UploadUsers_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(DemoServiceServer).UploadUsers(&demoServiceUploadUsersServer{stream})
}

type DemoService_UploadUsersServer interface {
	SendAndClose(*UploadUsersResponse) error
	Recv() (*UploadUsersRequest, error)
	grpc.ServerStream
}

type demoServiceUploadUsersServer struct {
	grpc.ServerStream
}

func (x *demoServiceUploadUsersServer) SendAndClose(m *UploadUsersResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *demoServiceUploadUsersServer) Recv() (*UploadUsersRequest, error) {
	m := new(UploadUsersRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _DemoService_Chat_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(DemoServiceServer).Chat(&demoServiceChatServer{stream})
}

type DemoService_ChatServer interface {
	Send(*ChatMessage) error
	Recv() (*ChatMessage, error)
	grpc.ServerStream
}

type demoServiceChatServer struct {
	grpc.ServerStream
}

func (x *demoServiceChatServer) Send(m *ChatMessage) error {
	return x.ServerStream.SendMsg(m)
}

func (x *demoServiceChatServer) Recv() (*ChatMessage, error) {
	m := new(ChatMessage)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _DemoService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "demo.DemoService",
	HandlerType: (*DemoServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetUser",
			Handler:    _DemoService_GetUser_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "ListUsers",
			Handler:       _DemoService_ListUsers_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "UploadUsers",
			Handler:       _DemoService_UploadUsers_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "Chat",
			Handler:       _DemoService_Chat_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "demo.proto",
}
