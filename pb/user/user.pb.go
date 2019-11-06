// Code generated by protoc-gen-go. DO NOT EDIT.
// source: user.proto

package pb

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type GetUserRequest struct {
	Uid                  string   `protobuf:"bytes,1,opt,name=uid,proto3" json:"uid,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetUserRequest) Reset()         { *m = GetUserRequest{} }
func (m *GetUserRequest) String() string { return proto.CompactTextString(m) }
func (*GetUserRequest) ProtoMessage()    {}
func (*GetUserRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_116e343673f7ffaf, []int{0}
}

func (m *GetUserRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetUserRequest.Unmarshal(m, b)
}
func (m *GetUserRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetUserRequest.Marshal(b, m, deterministic)
}
func (m *GetUserRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetUserRequest.Merge(m, src)
}
func (m *GetUserRequest) XXX_Size() int {
	return xxx_messageInfo_GetUserRequest.Size(m)
}
func (m *GetUserRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetUserRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetUserRequest proto.InternalMessageInfo

func (m *GetUserRequest) GetUid() string {
	if m != nil {
		return m.Uid
	}
	return ""
}

type GetUserReply struct {
	Uid                  string   `protobuf:"bytes,1,opt,name=uid,proto3" json:"uid,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetUserReply) Reset()         { *m = GetUserReply{} }
func (m *GetUserReply) String() string { return proto.CompactTextString(m) }
func (*GetUserReply) ProtoMessage()    {}
func (*GetUserReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_116e343673f7ffaf, []int{1}
}

func (m *GetUserReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetUserReply.Unmarshal(m, b)
}
func (m *GetUserReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetUserReply.Marshal(b, m, deterministic)
}
func (m *GetUserReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetUserReply.Merge(m, src)
}
func (m *GetUserReply) XXX_Size() int {
	return xxx_messageInfo_GetUserReply.Size(m)
}
func (m *GetUserReply) XXX_DiscardUnknown() {
	xxx_messageInfo_GetUserReply.DiscardUnknown(m)
}

var xxx_messageInfo_GetUserReply proto.InternalMessageInfo

func (m *GetUserReply) GetUid() string {
	if m != nil {
		return m.Uid
	}
	return ""
}

type LoginRequest struct {
	Username             string   `protobuf:"bytes,1,opt,name=username,proto3" json:"username,omitempty"`
	Password             string   `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *LoginRequest) Reset()         { *m = LoginRequest{} }
func (m *LoginRequest) String() string { return proto.CompactTextString(m) }
func (*LoginRequest) ProtoMessage()    {}
func (*LoginRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_116e343673f7ffaf, []int{2}
}

func (m *LoginRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LoginRequest.Unmarshal(m, b)
}
func (m *LoginRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LoginRequest.Marshal(b, m, deterministic)
}
func (m *LoginRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LoginRequest.Merge(m, src)
}
func (m *LoginRequest) XXX_Size() int {
	return xxx_messageInfo_LoginRequest.Size(m)
}
func (m *LoginRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_LoginRequest.DiscardUnknown(m)
}

var xxx_messageInfo_LoginRequest proto.InternalMessageInfo

func (m *LoginRequest) GetUsername() string {
	if m != nil {
		return m.Username
	}
	return ""
}

func (m *LoginRequest) GetPassword() string {
	if m != nil {
		return m.Password
	}
	return ""
}

type LoginReply struct {
	Id                   int32    `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Username             string   `protobuf:"bytes,2,opt,name=username,proto3" json:"username,omitempty"`
	Avatar               string   `protobuf:"bytes,3,opt,name=avatar,proto3" json:"avatar,omitempty"`
	RoleID               int32    `protobuf:"varint,4,opt,name=roleID,proto3" json:"roleID,omitempty"`
	RecentTime           string   `protobuf:"bytes,5,opt,name=recentTime,proto3" json:"recentTime,omitempty"`
	CreatedTime          string   `protobuf:"bytes,6,opt,name=createdTime,proto3" json:"createdTime,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *LoginReply) Reset()         { *m = LoginReply{} }
func (m *LoginReply) String() string { return proto.CompactTextString(m) }
func (*LoginReply) ProtoMessage()    {}
func (*LoginReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_116e343673f7ffaf, []int{3}
}

func (m *LoginReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LoginReply.Unmarshal(m, b)
}
func (m *LoginReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LoginReply.Marshal(b, m, deterministic)
}
func (m *LoginReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LoginReply.Merge(m, src)
}
func (m *LoginReply) XXX_Size() int {
	return xxx_messageInfo_LoginReply.Size(m)
}
func (m *LoginReply) XXX_DiscardUnknown() {
	xxx_messageInfo_LoginReply.DiscardUnknown(m)
}

var xxx_messageInfo_LoginReply proto.InternalMessageInfo

func (m *LoginReply) GetId() int32 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *LoginReply) GetUsername() string {
	if m != nil {
		return m.Username
	}
	return ""
}

func (m *LoginReply) GetAvatar() string {
	if m != nil {
		return m.Avatar
	}
	return ""
}

func (m *LoginReply) GetRoleID() int32 {
	if m != nil {
		return m.RoleID
	}
	return 0
}

func (m *LoginReply) GetRecentTime() string {
	if m != nil {
		return m.RecentTime
	}
	return ""
}

func (m *LoginReply) GetCreatedTime() string {
	if m != nil {
		return m.CreatedTime
	}
	return ""
}

type SendCodeRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SendCodeRequest) Reset()         { *m = SendCodeRequest{} }
func (m *SendCodeRequest) String() string { return proto.CompactTextString(m) }
func (*SendCodeRequest) ProtoMessage()    {}
func (*SendCodeRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_116e343673f7ffaf, []int{4}
}

func (m *SendCodeRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SendCodeRequest.Unmarshal(m, b)
}
func (m *SendCodeRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SendCodeRequest.Marshal(b, m, deterministic)
}
func (m *SendCodeRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SendCodeRequest.Merge(m, src)
}
func (m *SendCodeRequest) XXX_Size() int {
	return xxx_messageInfo_SendCodeRequest.Size(m)
}
func (m *SendCodeRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_SendCodeRequest.DiscardUnknown(m)
}

var xxx_messageInfo_SendCodeRequest proto.InternalMessageInfo

type SendCodeReply struct {
	CodeID               string   `protobuf:"bytes,1,opt,name=codeID,proto3" json:"codeID,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SendCodeReply) Reset()         { *m = SendCodeReply{} }
func (m *SendCodeReply) String() string { return proto.CompactTextString(m) }
func (*SendCodeReply) ProtoMessage()    {}
func (*SendCodeReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_116e343673f7ffaf, []int{5}
}

func (m *SendCodeReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SendCodeReply.Unmarshal(m, b)
}
func (m *SendCodeReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SendCodeReply.Marshal(b, m, deterministic)
}
func (m *SendCodeReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SendCodeReply.Merge(m, src)
}
func (m *SendCodeReply) XXX_Size() int {
	return xxx_messageInfo_SendCodeReply.Size(m)
}
func (m *SendCodeReply) XXX_DiscardUnknown() {
	xxx_messageInfo_SendCodeReply.DiscardUnknown(m)
}

var xxx_messageInfo_SendCodeReply proto.InternalMessageInfo

func (m *SendCodeReply) GetCodeID() string {
	if m != nil {
		return m.CodeID
	}
	return ""
}

func init() {
	proto.RegisterType((*GetUserRequest)(nil), "pb.GetUserRequest")
	proto.RegisterType((*GetUserReply)(nil), "pb.GetUserReply")
	proto.RegisterType((*LoginRequest)(nil), "pb.LoginRequest")
	proto.RegisterType((*LoginReply)(nil), "pb.LoginReply")
	proto.RegisterType((*SendCodeRequest)(nil), "pb.SendCodeRequest")
	proto.RegisterType((*SendCodeReply)(nil), "pb.SendCodeReply")
}

func init() { proto.RegisterFile("user.proto", fileDescriptor_116e343673f7ffaf) }

var fileDescriptor_116e343673f7ffaf = []byte{
	// 307 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x52, 0xdd, 0x4a, 0xc3, 0x30,
	0x14, 0xb6, 0x9d, 0xeb, 0xe6, 0x71, 0xce, 0x2d, 0xc2, 0x28, 0xbd, 0x90, 0x91, 0x1b, 0x05, 0x61,
	0xe0, 0xcf, 0x1b, 0x28, 0x8a, 0xe0, 0x55, 0xd5, 0x07, 0x48, 0x9b, 0x83, 0x14, 0xba, 0x26, 0x26,
	0xe9, 0x64, 0x4f, 0xe3, 0xbd, 0x4f, 0x29, 0x49, 0xd3, 0xd9, 0x0e, 0xef, 0xfa, 0xfd, 0xe4, 0x3b,
	0x39, 0x5f, 0x03, 0x50, 0x6b, 0x54, 0x2b, 0xa9, 0x84, 0x11, 0x24, 0x94, 0x19, 0xa5, 0x30, 0x7d,
	0x42, 0xf3, 0xae, 0x51, 0xa5, 0xf8, 0x59, 0xa3, 0x36, 0x64, 0x06, 0x83, 0xba, 0xe0, 0x71, 0xb0,
	0x0c, 0x2e, 0x8f, 0x52, 0xfb, 0x49, 0x97, 0x30, 0xd9, 0x79, 0x64, 0xb9, 0xfd, 0xc7, 0xf1, 0x08,
	0x93, 0x17, 0xf1, 0x51, 0x54, 0x6d, 0x46, 0x02, 0x63, 0x3b, 0xa7, 0x62, 0x6b, 0xf4, 0xb6, 0x1d,
	0xb6, 0x9a, 0x64, 0x5a, 0x7f, 0x09, 0xc5, 0xe3, 0xb0, 0xd1, 0x5a, 0x4c, 0x7f, 0x02, 0x00, 0x1f,
	0x64, 0x07, 0x4d, 0x21, 0xf4, 0x73, 0x86, 0x69, 0x58, 0xf0, 0x5e, 0x6c, 0xb8, 0x17, 0xbb, 0x80,
	0x88, 0x6d, 0x98, 0x61, 0x2a, 0x1e, 0x38, 0xc5, 0x23, 0xcb, 0x2b, 0x51, 0xe2, 0xf3, 0x43, 0x7c,
	0xe8, 0x72, 0x3c, 0x22, 0xe7, 0x00, 0x0a, 0x73, 0xac, 0xcc, 0x5b, 0xb1, 0xc6, 0x78, 0xe8, 0xce,
	0x74, 0x18, 0xb2, 0x84, 0xe3, 0x5c, 0x21, 0x33, 0xc8, 0x9d, 0x21, 0x72, 0x86, 0x2e, 0x45, 0xe7,
	0x70, 0xfa, 0x8a, 0x15, 0xbf, 0x17, 0x1c, 0xfd, 0xde, 0xf4, 0x02, 0x4e, 0xfe, 0x28, 0xbb, 0xc1,
	0x02, 0xa2, 0x5c, 0x70, 0x3b, 0xbd, 0xa9, 0xc1, 0xa3, 0x9b, 0xef, 0x00, 0x46, 0xb6, 0x50, 0xbd,
	0xc9, 0xc9, 0x35, 0x8c, 0x7c, 0xbd, 0x84, 0xac, 0x64, 0xb6, 0xea, 0xff, 0x8f, 0x64, 0xd6, 0xe3,
	0x64, 0xb9, 0xa5, 0x07, 0xe4, 0x0a, 0x86, 0xae, 0x26, 0xe2, 0xc4, 0x6e, 0xf5, 0xc9, 0xb4, 0xc3,
	0x34, 0xe6, 0x3b, 0x18, 0xb7, 0x97, 0x22, 0x67, 0x56, 0xdd, 0xbb, 0x75, 0x32, 0xef, 0x93, 0xee,
	0x54, 0x16, 0xb9, 0x37, 0x72, 0xfb, 0x1b, 0x00, 0x00, 0xff, 0xff, 0xb8, 0x84, 0x09, 0xd9, 0x31,
	0x02, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// UsersvcClient is the client API for Usersvc service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type UsersvcClient interface {
	GetUser(ctx context.Context, in *GetUserRequest, opts ...grpc.CallOption) (*GetUserReply, error)
	Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginReply, error)
	SendCode(ctx context.Context, in *SendCodeRequest, opts ...grpc.CallOption) (*SendCodeReply, error)
}

type usersvcClient struct {
	cc *grpc.ClientConn
}

func NewUsersvcClient(cc *grpc.ClientConn) UsersvcClient {
	return &usersvcClient{cc}
}

func (c *usersvcClient) GetUser(ctx context.Context, in *GetUserRequest, opts ...grpc.CallOption) (*GetUserReply, error) {
	out := new(GetUserReply)
	err := c.cc.Invoke(ctx, "/pb.Usersvc/GetUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *usersvcClient) Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginReply, error) {
	out := new(LoginReply)
	err := c.cc.Invoke(ctx, "/pb.Usersvc/Login", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *usersvcClient) SendCode(ctx context.Context, in *SendCodeRequest, opts ...grpc.CallOption) (*SendCodeReply, error) {
	out := new(SendCodeReply)
	err := c.cc.Invoke(ctx, "/pb.Usersvc/SendCode", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UsersvcServer is the server API for Usersvc service.
type UsersvcServer interface {
	GetUser(context.Context, *GetUserRequest) (*GetUserReply, error)
	Login(context.Context, *LoginRequest) (*LoginReply, error)
	SendCode(context.Context, *SendCodeRequest) (*SendCodeReply, error)
}

// UnimplementedUsersvcServer can be embedded to have forward compatible implementations.
type UnimplementedUsersvcServer struct {
}

func (*UnimplementedUsersvcServer) GetUser(ctx context.Context, req *GetUserRequest) (*GetUserReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUser not implemented")
}
func (*UnimplementedUsersvcServer) Login(ctx context.Context, req *LoginRequest) (*LoginReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Login not implemented")
}
func (*UnimplementedUsersvcServer) SendCode(ctx context.Context, req *SendCodeRequest) (*SendCodeReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendCode not implemented")
}

func RegisterUsersvcServer(s *grpc.Server, srv UsersvcServer) {
	s.RegisterService(&_Usersvc_serviceDesc, srv)
}

func _Usersvc_GetUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UsersvcServer).GetUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.Usersvc/GetUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UsersvcServer).GetUser(ctx, req.(*GetUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Usersvc_Login_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UsersvcServer).Login(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.Usersvc/Login",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UsersvcServer).Login(ctx, req.(*LoginRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Usersvc_SendCode_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendCodeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UsersvcServer).SendCode(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.Usersvc/SendCode",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UsersvcServer).SendCode(ctx, req.(*SendCodeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Usersvc_serviceDesc = grpc.ServiceDesc{
	ServiceName: "pb.Usersvc",
	HandlerType: (*UsersvcServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetUser",
			Handler:    _Usersvc_GetUser_Handler,
		},
		{
			MethodName: "Login",
			Handler:    _Usersvc_Login_Handler,
		},
		{
			MethodName: "SendCode",
			Handler:    _Usersvc_SendCode_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "user.proto",
}
