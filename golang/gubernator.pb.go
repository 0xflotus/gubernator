// Code generated by protoc-gen-go. DO NOT EDIT.
// source: gubernator.proto

/*
Package gubernator is a generated protocol buffer package.

It is generated from these files:
	gubernator.proto
	peers.proto

It has these top-level messages:
	GetRateLimitsReq
	GetRateLimitsResp
	RateLimitReq
	RateLimitResp
	HealthCheckReq
	HealthCheckResp
	GetPeerRateLimitsReq
	GetPeerRateLimitsResp
	UpdatePeerGlobalsReq
	UpdatePeerGlobalsResp
*/
package gubernator

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "google.golang.org/genproto/googleapis/api/annotations"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Algorithm int32

const (
	// Token bucket algorithm https://en.wikipedia.org/wiki/Token_bucket
	Algorithm_TOKEN_BUCKET Algorithm = 0
	// Leaky bucket algorithm https://en.wikipedia.org/wiki/Leaky_bucket
	Algorithm_LEAKY_BUCKET Algorithm = 1
)

var Algorithm_name = map[int32]string{
	0: "TOKEN_BUCKET",
	1: "LEAKY_BUCKET",
}
var Algorithm_value = map[string]int32{
	"TOKEN_BUCKET": 0,
	"LEAKY_BUCKET": 1,
}

func (x Algorithm) String() string {
	return proto.EnumName(Algorithm_name, int32(x))
}
func (Algorithm) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type Behavior int32

const (
	// BATCHING is the default behavior. This enables batching requests which protects the
	// service from thundering herd. IE: When a service experiences spikes of unexpected high
	// volume requests.
	//
	// Using this option introduces a small amount of latency depending on
	// the `batchWait` setting. Defaults to around 500 Microseconds of additional
	// latency in low throughput situations. For high volume loads, batching can reduce
	// the overall load on the system substantially.
	Behavior_BATCHING Behavior = 0
	// Disables batching. Use this for super low latency rate limit requests when
	// thundering herd is not a concern but latency of requests is of paramount importance.
	Behavior_NO_BATCHING Behavior = 1
	// Enables Global caching of the rate limit. Use this if the rate limit applies globally to
	// all ingress requests. (IE: Throttle hundreds of thousands of requests to an entire
	// datacenter or cluster of http servers)
	//
	// Using this option gubernator will continue to use a single peer as the rate limit coordinator
	// to increment and manage the state of the rate limit, however the result of the rate limit is
	// distributed to each peer and cached locally. A rate limit request received from any peer in the
	// cluster will first check the local cache for a rate limit answer, if it exists the peer will
	// immediately return the answer to the client and asynchronously forward the aggregate hits to
	// the peer coordinator. Because of GLOBALS async nature we lose some accuracy in rate limit
	// reporting, which may result in allowing some requests beyond the chosen rate limit. However we
	// gain massive performance as every request coming into the system does not have to wait for a
	// single peer to decide if the rate limit has been reached.
	Behavior_GLOBAL Behavior = 2
)

var Behavior_name = map[int32]string{
	0: "BATCHING",
	1: "NO_BATCHING",
	2: "GLOBAL",
}
var Behavior_value = map[string]int32{
	"BATCHING":    0,
	"NO_BATCHING": 1,
	"GLOBAL":      2,
}

func (x Behavior) String() string {
	return proto.EnumName(Behavior_name, int32(x))
}
func (Behavior) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

type Status int32

const (
	Status_UNDER_LIMIT Status = 0
	Status_OVER_LIMIT  Status = 1
)

var Status_name = map[int32]string{
	0: "UNDER_LIMIT",
	1: "OVER_LIMIT",
}
var Status_value = map[string]int32{
	"UNDER_LIMIT": 0,
	"OVER_LIMIT":  1,
}

func (x Status) String() string {
	return proto.EnumName(Status_name, int32(x))
}
func (Status) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

// Must specify at least one Request
type GetRateLimitsReq struct {
	Requests []*RateLimitReq `protobuf:"bytes,1,rep,name=requests" json:"requests,omitempty"`
}

func (m *GetRateLimitsReq) Reset()                    { *m = GetRateLimitsReq{} }
func (m *GetRateLimitsReq) String() string            { return proto.CompactTextString(m) }
func (*GetRateLimitsReq) ProtoMessage()               {}
func (*GetRateLimitsReq) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *GetRateLimitsReq) GetRequests() []*RateLimitReq {
	if m != nil {
		return m.Requests
	}
	return nil
}

// RateLimits returned are in the same order as the Requests
type GetRateLimitsResp struct {
	Responses []*RateLimitResp `protobuf:"bytes,1,rep,name=responses" json:"responses,omitempty"`
}

func (m *GetRateLimitsResp) Reset()                    { *m = GetRateLimitsResp{} }
func (m *GetRateLimitsResp) String() string            { return proto.CompactTextString(m) }
func (*GetRateLimitsResp) ProtoMessage()               {}
func (*GetRateLimitsResp) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *GetRateLimitsResp) GetResponses() []*RateLimitResp {
	if m != nil {
		return m.Responses
	}
	return nil
}

type RateLimitReq struct {
	// The name of the rate limit IE: 'requests_per_second', 'gets_per_minute`
	Name string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	// Uniquely identifies this rate limit IE: 'ip:10.2.10.7' or 'account:123445'
	UniqueKey string `protobuf:"bytes,2,opt,name=unique_key,json=uniqueKey" json:"unique_key,omitempty"`
	// Rate limit requests optionally specify the number of hits a request adds to the matched limit. If Hit
	// is zero, the request returns the current limit, but does not increment the hit count.
	Hits int64 `protobuf:"varint,3,opt,name=hits" json:"hits,omitempty"`
	// The number of requests that can occur for the duration of the rate limit
	Limit int64 `protobuf:"varint,4,opt,name=limit" json:"limit,omitempty"`
	// The duration of the rate limit in milliseconds
	// Second = 1000 Milliseconds
	// Minute = 60000 Milliseconds
	// Hour = 3600000 Milliseconds
	Duration int64 `protobuf:"varint,5,opt,name=duration" json:"duration,omitempty"`
	// The algorithm used to calculate the rate limit. The algorithm may change on
	// subsequent requests, when this occurs any previous rate limit hit counts are reset.
	Algorithm Algorithm `protobuf:"varint,6,opt,name=algorithm,enum=pb.gubernator.Algorithm" json:"algorithm,omitempty"`
	// The behavior of the rate limit in gubernator.
	Behavior Behavior `protobuf:"varint,7,opt,name=behavior,enum=pb.gubernator.Behavior" json:"behavior,omitempty"`
}

func (m *RateLimitReq) Reset()                    { *m = RateLimitReq{} }
func (m *RateLimitReq) String() string            { return proto.CompactTextString(m) }
func (*RateLimitReq) ProtoMessage()               {}
func (*RateLimitReq) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *RateLimitReq) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *RateLimitReq) GetUniqueKey() string {
	if m != nil {
		return m.UniqueKey
	}
	return ""
}

func (m *RateLimitReq) GetHits() int64 {
	if m != nil {
		return m.Hits
	}
	return 0
}

func (m *RateLimitReq) GetLimit() int64 {
	if m != nil {
		return m.Limit
	}
	return 0
}

func (m *RateLimitReq) GetDuration() int64 {
	if m != nil {
		return m.Duration
	}
	return 0
}

func (m *RateLimitReq) GetAlgorithm() Algorithm {
	if m != nil {
		return m.Algorithm
	}
	return Algorithm_TOKEN_BUCKET
}

func (m *RateLimitReq) GetBehavior() Behavior {
	if m != nil {
		return m.Behavior
	}
	return Behavior_BATCHING
}

type RateLimitResp struct {
	// The status of the rate limit.
	Status Status `protobuf:"varint,1,opt,name=status,enum=pb.gubernator.Status" json:"status,omitempty"`
	// The currently configured request limit (Identical to RateLimitRequest.rate_limit_config.limit).
	Limit int64 `protobuf:"varint,2,opt,name=limit" json:"limit,omitempty"`
	// This is the number of requests remaining before the limit is hit.
	Remaining int64 `protobuf:"varint,3,opt,name=remaining" json:"remaining,omitempty"`
	// This is the time when the rate limit span will be reset, provided as a unix timestamp in milliseconds.
	ResetTime int64 `protobuf:"varint,4,opt,name=reset_time,json=resetTime" json:"reset_time,omitempty"`
	// Contains the error; If set all other values should be ignored
	Error string `protobuf:"bytes,5,opt,name=error" json:"error,omitempty"`
	// This is additional metadata that a client might find useful. (IE: Additional headers, corrdinator ownership, etc..)
	Metadata map[string]string `protobuf:"bytes,6,rep,name=metadata" json:"metadata,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
}

func (m *RateLimitResp) Reset()                    { *m = RateLimitResp{} }
func (m *RateLimitResp) String() string            { return proto.CompactTextString(m) }
func (*RateLimitResp) ProtoMessage()               {}
func (*RateLimitResp) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *RateLimitResp) GetStatus() Status {
	if m != nil {
		return m.Status
	}
	return Status_UNDER_LIMIT
}

func (m *RateLimitResp) GetLimit() int64 {
	if m != nil {
		return m.Limit
	}
	return 0
}

func (m *RateLimitResp) GetRemaining() int64 {
	if m != nil {
		return m.Remaining
	}
	return 0
}

func (m *RateLimitResp) GetResetTime() int64 {
	if m != nil {
		return m.ResetTime
	}
	return 0
}

func (m *RateLimitResp) GetError() string {
	if m != nil {
		return m.Error
	}
	return ""
}

func (m *RateLimitResp) GetMetadata() map[string]string {
	if m != nil {
		return m.Metadata
	}
	return nil
}

type HealthCheckReq struct {
}

func (m *HealthCheckReq) Reset()                    { *m = HealthCheckReq{} }
func (m *HealthCheckReq) String() string            { return proto.CompactTextString(m) }
func (*HealthCheckReq) ProtoMessage()               {}
func (*HealthCheckReq) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

type HealthCheckResp struct {
	// Valid entries are 'healthy' or 'unhealthy'
	Status string `protobuf:"bytes,1,opt,name=status" json:"status,omitempty"`
	// If 'unhealthy', message indicates the problem
	Message string `protobuf:"bytes,2,opt,name=message" json:"message,omitempty"`
	// The number of peers we know about
	PeerCount int32 `protobuf:"varint,3,opt,name=peer_count,json=peerCount" json:"peer_count,omitempty"`
}

func (m *HealthCheckResp) Reset()                    { *m = HealthCheckResp{} }
func (m *HealthCheckResp) String() string            { return proto.CompactTextString(m) }
func (*HealthCheckResp) ProtoMessage()               {}
func (*HealthCheckResp) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *HealthCheckResp) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

func (m *HealthCheckResp) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

func (m *HealthCheckResp) GetPeerCount() int32 {
	if m != nil {
		return m.PeerCount
	}
	return 0
}

func init() {
	proto.RegisterType((*GetRateLimitsReq)(nil), "pb.gubernator.GetRateLimitsReq")
	proto.RegisterType((*GetRateLimitsResp)(nil), "pb.gubernator.GetRateLimitsResp")
	proto.RegisterType((*RateLimitReq)(nil), "pb.gubernator.RateLimitReq")
	proto.RegisterType((*RateLimitResp)(nil), "pb.gubernator.RateLimitResp")
	proto.RegisterType((*HealthCheckReq)(nil), "pb.gubernator.HealthCheckReq")
	proto.RegisterType((*HealthCheckResp)(nil), "pb.gubernator.HealthCheckResp")
	proto.RegisterEnum("pb.gubernator.Algorithm", Algorithm_name, Algorithm_value)
	proto.RegisterEnum("pb.gubernator.Behavior", Behavior_name, Behavior_value)
	proto.RegisterEnum("pb.gubernator.Status", Status_name, Status_value)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for GubernatorV1 service

type GubernatorV1Client interface {
	// Given a list of rate limit requests, return the rate limits of each.
	GetRateLimits(ctx context.Context, in *GetRateLimitsReq, opts ...grpc.CallOption) (*GetRateLimitsResp, error)
	// This method is for round trip benchmarking and can be used by
	// the client to determine connectivity to the server
	HealthCheck(ctx context.Context, in *HealthCheckReq, opts ...grpc.CallOption) (*HealthCheckResp, error)
}

type gubernatorV1Client struct {
	cc *grpc.ClientConn
}

func NewGubernatorV1Client(cc *grpc.ClientConn) GubernatorV1Client {
	return &gubernatorV1Client{cc}
}

func (c *gubernatorV1Client) GetRateLimits(ctx context.Context, in *GetRateLimitsReq, opts ...grpc.CallOption) (*GetRateLimitsResp, error) {
	out := new(GetRateLimitsResp)
	err := grpc.Invoke(ctx, "/pb.gubernator.GubernatorV1/GetRateLimits", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gubernatorV1Client) HealthCheck(ctx context.Context, in *HealthCheckReq, opts ...grpc.CallOption) (*HealthCheckResp, error) {
	out := new(HealthCheckResp)
	err := grpc.Invoke(ctx, "/pb.gubernator.GubernatorV1/HealthCheck", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for GubernatorV1 service

type GubernatorV1Server interface {
	// Given a list of rate limit requests, return the rate limits of each.
	GetRateLimits(context.Context, *GetRateLimitsReq) (*GetRateLimitsResp, error)
	// This method is for round trip benchmarking and can be used by
	// the client to determine connectivity to the server
	HealthCheck(context.Context, *HealthCheckReq) (*HealthCheckResp, error)
}

func RegisterGubernatorV1Server(s *grpc.Server, srv GubernatorV1Server) {
	s.RegisterService(&_GubernatorV1_serviceDesc, srv)
}

func _GubernatorV1_GetRateLimits_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRateLimitsReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GubernatorV1Server).GetRateLimits(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.gubernator.GubernatorV1/GetRateLimits",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GubernatorV1Server).GetRateLimits(ctx, req.(*GetRateLimitsReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _GubernatorV1_HealthCheck_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HealthCheckReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GubernatorV1Server).HealthCheck(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.gubernator.GubernatorV1/HealthCheck",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GubernatorV1Server).HealthCheck(ctx, req.(*HealthCheckReq))
	}
	return interceptor(ctx, in, info, handler)
}

var _GubernatorV1_serviceDesc = grpc.ServiceDesc{
	ServiceName: "pb.gubernator.GubernatorV1",
	HandlerType: (*GubernatorV1Server)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetRateLimits",
			Handler:    _GubernatorV1_GetRateLimits_Handler,
		},
		{
			MethodName: "HealthCheck",
			Handler:    _GubernatorV1_HealthCheck_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "gubernator.proto",
}

func init() { proto.RegisterFile("gubernator.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 659 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x54, 0x5d, 0x6e, 0xda, 0x4c,
	0x14, 0x8d, 0x21, 0x21, 0xf8, 0x86, 0x1f, 0x67, 0xf4, 0x7d, 0x89, 0x45, 0x49, 0x8b, 0xfc, 0x44,
	0x91, 0x0a, 0x0a, 0x51, 0x7f, 0x94, 0x3e, 0x01, 0xa5, 0x24, 0x82, 0x80, 0x34, 0x25, 0x91, 0xda,
	0x17, 0x34, 0x24, 0x57, 0x60, 0x05, 0xff, 0xe0, 0x19, 0x47, 0xca, 0x5b, 0xd5, 0x2d, 0x74, 0x1b,
	0xdd, 0x4d, 0xb7, 0x50, 0xa9, 0x1b, 0xe8, 0x02, 0xaa, 0x19, 0xc0, 0x60, 0xa4, 0xe6, 0x6d, 0xee,
	0x39, 0xe7, 0x9e, 0x61, 0xce, 0x91, 0x01, 0x63, 0x12, 0x8e, 0x31, 0x70, 0x99, 0xf0, 0x82, 0xaa,
	0x1f, 0x78, 0xc2, 0x23, 0x59, 0x7f, 0x5c, 0x5d, 0x83, 0x85, 0xe2, 0xc4, 0xf3, 0x26, 0x33, 0xac,
	0x31, 0xdf, 0xae, 0x31, 0xd7, 0xf5, 0x04, 0x13, 0xb6, 0xe7, 0xf2, 0x85, 0xd8, 0xea, 0x82, 0xd1,
	0x41, 0x41, 0x99, 0xc0, 0x9e, 0xed, 0xd8, 0x82, 0x53, 0x9c, 0x93, 0xb7, 0x90, 0x0e, 0x70, 0x1e,
	0x22, 0x17, 0xdc, 0xd4, 0x4a, 0xc9, 0xf2, 0x41, 0xfd, 0x59, 0x35, 0xe6, 0x59, 0x8d, 0xf4, 0x14,
	0xe7, 0x34, 0x12, 0x5b, 0x03, 0x38, 0xdc, 0x32, 0xe3, 0x3e, 0x39, 0x07, 0x3d, 0x40, 0xee, 0x7b,
	0x2e, 0xc7, 0x95, 0x5d, 0xf1, 0xdf, 0x76, 0xdc, 0xa7, 0x6b, 0xb9, 0xf5, 0x47, 0x83, 0xcc, 0xe6,
	0x5d, 0x84, 0xc0, 0xae, 0xcb, 0x1c, 0x34, 0xb5, 0x92, 0x56, 0xd6, 0xa9, 0x3a, 0x93, 0x13, 0x80,
	0xd0, 0xb5, 0xe7, 0x21, 0x8e, 0xee, 0xf1, 0xd1, 0x4c, 0x28, 0x46, 0x5f, 0x20, 0x5d, 0x7c, 0x94,
	0x2b, 0x53, 0x5b, 0x70, 0x33, 0x59, 0xd2, 0xca, 0x49, 0xaa, 0xce, 0xe4, 0x3f, 0xd8, 0x9b, 0x49,
	0x4b, 0x73, 0x57, 0x81, 0x8b, 0x81, 0x14, 0x20, 0x7d, 0x17, 0x06, 0x2a, 0x1e, 0x73, 0x4f, 0x11,
	0xd1, 0x4c, 0xde, 0x80, 0xce, 0x66, 0x13, 0x2f, 0xb0, 0xc5, 0xd4, 0x31, 0x53, 0x25, 0xad, 0x9c,
	0xab, 0x9b, 0x5b, 0xaf, 0x68, 0xac, 0x78, 0xba, 0x96, 0x92, 0x33, 0x48, 0x8f, 0x71, 0xca, 0x1e,
	0x6c, 0x2f, 0x30, 0xf7, 0xd5, 0xda, 0xf1, 0xd6, 0x5a, 0x73, 0x49, 0xd3, 0x48, 0x68, 0xfd, 0x48,
	0x40, 0x36, 0x96, 0x09, 0x79, 0x05, 0x29, 0x2e, 0x98, 0x08, 0xb9, 0x7a, 0x79, 0xae, 0xfe, 0xff,
	0x96, 0xc9, 0x27, 0x45, 0xd2, 0xa5, 0x68, 0xfd, 0xbe, 0xc4, 0xe6, 0xfb, 0x8a, 0xb2, 0x09, 0x87,
	0xd9, 0xae, 0xed, 0x4e, 0x96, 0x71, 0xac, 0x01, 0x19, 0x63, 0x80, 0x1c, 0xc5, 0x48, 0xd8, 0x0e,
	0x2e, 0x83, 0xd1, 0x15, 0x32, 0xb4, 0x1d, 0x94, 0x96, 0x18, 0x04, 0x5e, 0xa0, 0x92, 0xd1, 0xe9,
	0x62, 0x20, 0x1f, 0x21, 0xed, 0xa0, 0x60, 0x77, 0x4c, 0x30, 0x33, 0xa5, 0xba, 0xad, 0x3c, 0xd5,
	0x6d, 0xf5, 0x6a, 0x29, 0x6e, 0xbb, 0x22, 0x78, 0xa4, 0xd1, 0x6e, 0xe1, 0x3d, 0x64, 0x63, 0x14,
	0x31, 0x20, 0x29, 0xdb, 0x5c, 0xf4, 0x2c, 0x8f, 0xf2, 0x07, 0x3c, 0xb0, 0x59, 0x88, 0xcb, 0x86,
	0x17, 0xc3, 0x79, 0xe2, 0x9d, 0x66, 0x19, 0x90, 0xbb, 0x40, 0x36, 0x13, 0xd3, 0xd6, 0x14, 0x6f,
	0xef, 0x29, 0xce, 0xad, 0x31, 0xe4, 0x63, 0x08, 0xf7, 0xc9, 0x51, 0x2c, 0x41, 0x3d, 0x8a, 0xca,
	0x84, 0x7d, 0x07, 0x39, 0x67, 0x93, 0x95, 0xf1, 0x6a, 0x94, 0x81, 0xf8, 0x88, 0xc1, 0xe8, 0xd6,
	0x0b, 0x5d, 0xa1, 0xf2, 0xda, 0xa3, 0xba, 0x44, 0x5a, 0x12, 0xa8, 0xd4, 0x40, 0x8f, 0x1a, 0x27,
	0x06, 0x64, 0x86, 0x83, 0x6e, 0xbb, 0x3f, 0x6a, 0x5e, 0xb7, 0xba, 0xed, 0xa1, 0xb1, 0x23, 0x91,
	0x5e, 0xbb, 0xd1, 0xfd, 0xbc, 0x42, 0xb4, 0xca, 0x6b, 0x48, 0xaf, 0xba, 0x26, 0x19, 0x48, 0x37,
	0x1b, 0xc3, 0xd6, 0xc5, 0x65, 0xbf, 0x63, 0xec, 0x90, 0x3c, 0x1c, 0xf4, 0x07, 0xa3, 0x08, 0xd0,
	0x08, 0x40, 0xaa, 0xd3, 0x1b, 0x34, 0x1b, 0x3d, 0x23, 0x51, 0x79, 0x09, 0xa9, 0x45, 0xbb, 0x52,
	0x76, 0xdd, 0xff, 0xd0, 0xa6, 0xa3, 0xde, 0xe5, 0xd5, 0xa5, 0xbc, 0x23, 0x07, 0x30, 0xb8, 0x89,
	0x66, 0xad, 0xfe, 0x5b, 0x83, 0x4c, 0x27, 0x8a, 0xfe, 0xe6, 0x94, 0xf8, 0x90, 0x8d, 0x7d, 0x90,
	0xe4, 0xc5, 0x56, 0x3b, 0xdb, 0xdf, 0x7e, 0xa1, 0xf4, 0xb4, 0x80, 0xfb, 0x56, 0xf1, 0xdb, 0xcf,
	0x5f, 0xdf, 0x13, 0x47, 0xd6, 0x61, 0xed, 0xe1, 0xb4, 0x16, 0xa3, 0xcf, 0xb5, 0x0a, 0x41, 0x38,
	0xd8, 0x48, 0x9e, 0x9c, 0x6c, 0xd9, 0xc5, 0x7b, 0x2a, 0x3c, 0x7f, 0x8a, 0xe6, 0xbe, 0x75, 0xac,
	0xee, 0x3a, 0x24, 0x79, 0x79, 0xd7, 0x06, 0xd9, 0xcc, 0x7f, 0x81, 0xf5, 0xda, 0x57, 0x4d, 0x1b,
	0xa7, 0xd4, 0xdf, 0xd9, 0xd9, 0xdf, 0x00, 0x00, 0x00, 0xff, 0xff, 0x0a, 0x2a, 0x9c, 0x9f, 0x0f,
	0x05, 0x00, 0x00,
}