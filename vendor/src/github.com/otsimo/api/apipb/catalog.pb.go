// Code generated by protoc-gen-gogo.
// source: catalog.proto
// DO NOT EDIT!

package apipb

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

import io "io"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type CatalogCategory int32

const (
	CatalogCategory_NONE             CatalogCategory = 0
	CatalogCategory_FEATURED         CatalogCategory = 1
	CatalogCategory_NEW              CatalogCategory = 2
	CatalogCategory_RECOMMENDATION   CatalogCategory = 3
	CatalogCategory_POPULAR          CatalogCategory = 4
	CatalogCategory_RECENTLY_UPDATED CatalogCategory = 5
)

var CatalogCategory_name = map[int32]string{
	0: "NONE",
	1: "FEATURED",
	2: "NEW",
	3: "RECOMMENDATION",
	4: "POPULAR",
	5: "RECENTLY_UPDATED",
}
var CatalogCategory_value = map[string]int32{
	"NONE":             0,
	"FEATURED":         1,
	"NEW":              2,
	"RECOMMENDATION":   3,
	"POPULAR":          4,
	"RECENTLY_UPDATED": 5,
}

func (x CatalogCategory) String() string {
	return proto.EnumName(CatalogCategory_name, int32(x))
}

type CatalogItem struct {
	GameId   string          `protobuf:"bytes,1,opt,name=game_id,proto3" json:"game_id,omitempty"`
	Category CatalogCategory `protobuf:"varint,2,opt,name=category,proto3,enum=apipb.CatalogCategory" json:"category,omitempty"`
	Index    int32           `protobuf:"varint,3,opt,name=index,proto3" json:"index,omitempty"`
}

func (m *CatalogItem) Reset()         { *m = CatalogItem{} }
func (m *CatalogItem) String() string { return proto.CompactTextString(m) }
func (*CatalogItem) ProtoMessage()    {}

type CatalogPullRequest struct {
	ProfileId     string `protobuf:"bytes,1,opt,name=profile_id,proto3" json:"profile_id,omitempty"`
	ClientVersion string `protobuf:"bytes,2,opt,name=client_version,proto3" json:"client_version,omitempty"`
}

func (m *CatalogPullRequest) Reset()         { *m = CatalogPullRequest{} }
func (m *CatalogPullRequest) String() string { return proto.CompactTextString(m) }
func (*CatalogPullRequest) ProtoMessage()    {}

type Catalog struct {
	Title     string         `protobuf:"bytes,1,opt,name=title,proto3" json:"title,omitempty"`
	CreatedAt int64          `protobuf:"varint,2,opt,name=created_at,proto3" json:"created_at,omitempty"`
	VisibleAt int64          `protobuf:"varint,3,opt,name=visible_at,proto3" json:"visible_at,omitempty"`
	ExpiresAt int64          `protobuf:"varint,4,opt,name=expires_at,proto3" json:"expires_at,omitempty"`
	Items     []*CatalogItem `protobuf:"bytes,5,rep,name=items" json:"items,omitempty"`
}

func (m *Catalog) Reset()         { *m = Catalog{} }
func (m *Catalog) String() string { return proto.CompactTextString(m) }
func (*Catalog) ProtoMessage()    {}

func (m *Catalog) GetItems() []*CatalogItem {
	if m != nil {
		return m.Items
	}
	return nil
}

func init() {
	proto.RegisterEnum("apipb.CatalogCategory", CatalogCategory_name, CatalogCategory_value)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// Client API for CatalogService service

type CatalogServiceClient interface {
	Pull(ctx context.Context, in *CatalogPullRequest, opts ...grpc.CallOption) (*Catalog, error)
	Push(ctx context.Context, in *Catalog, opts ...grpc.CallOption) (*Response, error)
}

type catalogServiceClient struct {
	cc *grpc.ClientConn
}

func NewCatalogServiceClient(cc *grpc.ClientConn) CatalogServiceClient {
	return &catalogServiceClient{cc}
}

func (c *catalogServiceClient) Pull(ctx context.Context, in *CatalogPullRequest, opts ...grpc.CallOption) (*Catalog, error) {
	out := new(Catalog)
	err := grpc.Invoke(ctx, "/apipb.CatalogService/Pull", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *catalogServiceClient) Push(ctx context.Context, in *Catalog, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := grpc.Invoke(ctx, "/apipb.CatalogService/Push", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for CatalogService service

type CatalogServiceServer interface {
	Pull(context.Context, *CatalogPullRequest) (*Catalog, error)
	Push(context.Context, *Catalog) (*Response, error)
}

func RegisterCatalogServiceServer(s *grpc.Server, srv CatalogServiceServer) {
	s.RegisterService(&_CatalogService_serviceDesc, srv)
}

func _CatalogService_Pull_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error) (interface{}, error) {
	in := new(CatalogPullRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	out, err := srv.(CatalogServiceServer).Pull(ctx, in)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func _CatalogService_Push_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error) (interface{}, error) {
	in := new(Catalog)
	if err := dec(in); err != nil {
		return nil, err
	}
	out, err := srv.(CatalogServiceServer).Push(ctx, in)
	if err != nil {
		return nil, err
	}
	return out, nil
}

var _CatalogService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "apipb.CatalogService",
	HandlerType: (*CatalogServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Pull",
			Handler:    _CatalogService_Pull_Handler,
		},
		{
			MethodName: "Push",
			Handler:    _CatalogService_Push_Handler,
		},
	},
	Streams: []grpc.StreamDesc{},
}

func (m *CatalogItem) Marshal() (data []byte, err error) {
	size := m.Size()
	data = make([]byte, size)
	n, err := m.MarshalTo(data)
	if err != nil {
		return nil, err
	}
	return data[:n], nil
}

func (m *CatalogItem) MarshalTo(data []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.GameId) > 0 {
		data[i] = 0xa
		i++
		i = encodeVarintCatalog(data, i, uint64(len(m.GameId)))
		i += copy(data[i:], m.GameId)
	}
	if m.Category != 0 {
		data[i] = 0x10
		i++
		i = encodeVarintCatalog(data, i, uint64(m.Category))
	}
	if m.Index != 0 {
		data[i] = 0x18
		i++
		i = encodeVarintCatalog(data, i, uint64(m.Index))
	}
	return i, nil
}

func (m *CatalogPullRequest) Marshal() (data []byte, err error) {
	size := m.Size()
	data = make([]byte, size)
	n, err := m.MarshalTo(data)
	if err != nil {
		return nil, err
	}
	return data[:n], nil
}

func (m *CatalogPullRequest) MarshalTo(data []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.ProfileId) > 0 {
		data[i] = 0xa
		i++
		i = encodeVarintCatalog(data, i, uint64(len(m.ProfileId)))
		i += copy(data[i:], m.ProfileId)
	}
	if len(m.ClientVersion) > 0 {
		data[i] = 0x12
		i++
		i = encodeVarintCatalog(data, i, uint64(len(m.ClientVersion)))
		i += copy(data[i:], m.ClientVersion)
	}
	return i, nil
}

func (m *Catalog) Marshal() (data []byte, err error) {
	size := m.Size()
	data = make([]byte, size)
	n, err := m.MarshalTo(data)
	if err != nil {
		return nil, err
	}
	return data[:n], nil
}

func (m *Catalog) MarshalTo(data []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Title) > 0 {
		data[i] = 0xa
		i++
		i = encodeVarintCatalog(data, i, uint64(len(m.Title)))
		i += copy(data[i:], m.Title)
	}
	if m.CreatedAt != 0 {
		data[i] = 0x10
		i++
		i = encodeVarintCatalog(data, i, uint64(m.CreatedAt))
	}
	if m.VisibleAt != 0 {
		data[i] = 0x18
		i++
		i = encodeVarintCatalog(data, i, uint64(m.VisibleAt))
	}
	if m.ExpiresAt != 0 {
		data[i] = 0x20
		i++
		i = encodeVarintCatalog(data, i, uint64(m.ExpiresAt))
	}
	if len(m.Items) > 0 {
		for _, msg := range m.Items {
			data[i] = 0x2a
			i++
			i = encodeVarintCatalog(data, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(data[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	return i, nil
}

func encodeFixed64Catalog(data []byte, offset int, v uint64) int {
	data[offset] = uint8(v)
	data[offset+1] = uint8(v >> 8)
	data[offset+2] = uint8(v >> 16)
	data[offset+3] = uint8(v >> 24)
	data[offset+4] = uint8(v >> 32)
	data[offset+5] = uint8(v >> 40)
	data[offset+6] = uint8(v >> 48)
	data[offset+7] = uint8(v >> 56)
	return offset + 8
}
func encodeFixed32Catalog(data []byte, offset int, v uint32) int {
	data[offset] = uint8(v)
	data[offset+1] = uint8(v >> 8)
	data[offset+2] = uint8(v >> 16)
	data[offset+3] = uint8(v >> 24)
	return offset + 4
}
func encodeVarintCatalog(data []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		data[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	data[offset] = uint8(v)
	return offset + 1
}
func (m *CatalogItem) Size() (n int) {
	var l int
	_ = l
	l = len(m.GameId)
	if l > 0 {
		n += 1 + l + sovCatalog(uint64(l))
	}
	if m.Category != 0 {
		n += 1 + sovCatalog(uint64(m.Category))
	}
	if m.Index != 0 {
		n += 1 + sovCatalog(uint64(m.Index))
	}
	return n
}

func (m *CatalogPullRequest) Size() (n int) {
	var l int
	_ = l
	l = len(m.ProfileId)
	if l > 0 {
		n += 1 + l + sovCatalog(uint64(l))
	}
	l = len(m.ClientVersion)
	if l > 0 {
		n += 1 + l + sovCatalog(uint64(l))
	}
	return n
}

func (m *Catalog) Size() (n int) {
	var l int
	_ = l
	l = len(m.Title)
	if l > 0 {
		n += 1 + l + sovCatalog(uint64(l))
	}
	if m.CreatedAt != 0 {
		n += 1 + sovCatalog(uint64(m.CreatedAt))
	}
	if m.VisibleAt != 0 {
		n += 1 + sovCatalog(uint64(m.VisibleAt))
	}
	if m.ExpiresAt != 0 {
		n += 1 + sovCatalog(uint64(m.ExpiresAt))
	}
	if len(m.Items) > 0 {
		for _, e := range m.Items {
			l = e.Size()
			n += 1 + l + sovCatalog(uint64(l))
		}
	}
	return n
}

func sovCatalog(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozCatalog(x uint64) (n int) {
	return sovCatalog(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *CatalogItem) Unmarshal(data []byte) error {
	l := len(data)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowCatalog
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := data[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: CatalogItem: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: CatalogItem: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field GameId", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCatalog
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := data[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthCatalog
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.GameId = string(data[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Category", wireType)
			}
			m.Category = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCatalog
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := data[iNdEx]
				iNdEx++
				m.Category |= (CatalogCategory(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Index", wireType)
			}
			m.Index = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCatalog
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := data[iNdEx]
				iNdEx++
				m.Index |= (int32(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipCatalog(data[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthCatalog
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *CatalogPullRequest) Unmarshal(data []byte) error {
	l := len(data)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowCatalog
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := data[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: CatalogPullRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: CatalogPullRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ProfileId", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCatalog
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := data[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthCatalog
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ProfileId = string(data[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ClientVersion", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCatalog
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := data[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthCatalog
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ClientVersion = string(data[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipCatalog(data[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthCatalog
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Catalog) Unmarshal(data []byte) error {
	l := len(data)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowCatalog
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := data[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Catalog: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Catalog: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Title", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCatalog
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := data[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthCatalog
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Title = string(data[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field CreatedAt", wireType)
			}
			m.CreatedAt = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCatalog
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := data[iNdEx]
				iNdEx++
				m.CreatedAt |= (int64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field VisibleAt", wireType)
			}
			m.VisibleAt = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCatalog
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := data[iNdEx]
				iNdEx++
				m.VisibleAt |= (int64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ExpiresAt", wireType)
			}
			m.ExpiresAt = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCatalog
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := data[iNdEx]
				iNdEx++
				m.ExpiresAt |= (int64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Items", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCatalog
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := data[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthCatalog
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Items = append(m.Items, &CatalogItem{})
			if err := m.Items[len(m.Items)-1].Unmarshal(data[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipCatalog(data[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthCatalog
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipCatalog(data []byte) (n int, err error) {
	l := len(data)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowCatalog
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := data[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowCatalog
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if data[iNdEx-1] < 0x80 {
					break
				}
			}
			return iNdEx, nil
		case 1:
			iNdEx += 8
			return iNdEx, nil
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowCatalog
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := data[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			iNdEx += length
			if length < 0 {
				return 0, ErrInvalidLengthCatalog
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowCatalog
					}
					if iNdEx >= l {
						return 0, io.ErrUnexpectedEOF
					}
					b := data[iNdEx]
					iNdEx++
					innerWire |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				innerWireType := int(innerWire & 0x7)
				if innerWireType == 4 {
					break
				}
				next, err := skipCatalog(data[start:])
				if err != nil {
					return 0, err
				}
				iNdEx = start + next
			}
			return iNdEx, nil
		case 4:
			return iNdEx, nil
		case 5:
			iNdEx += 4
			return iNdEx, nil
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
	}
	panic("unreachable")
}

var (
	ErrInvalidLengthCatalog = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowCatalog   = fmt.Errorf("proto: integer overflow")
)
