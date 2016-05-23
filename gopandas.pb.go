// Code generated by protoc-gen-go.
// source: gopandas.proto
// DO NOT EDIT!

/*
Package gopandas is a generated protocol buffer package.

It is generated from these files:
	gopandas.proto

It has these top-level messages:
	FlyTimePanel
*/
package gopandas

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
const _ = proto.ProtoPackageIsVersion1

type FlyTimePanel struct {
	Data             []float64 `protobuf:"fixed64,1,rep,packed,name=data" json:"data,omitempty"`
	Dates            []uint64  `protobuf:"varint,2,rep,packed,name=dates" json:"dates,omitempty"`
	Secondary        []string  `protobuf:"bytes,3,rep,name=secondary" json:"secondary,omitempty"`
	Thirdly          []string  `protobuf:"bytes,4,rep,name=thirdly" json:"thirdly,omitempty"`
	XXX_unrecognized []byte    `json:"-"`
}

func (m *FlyTimePanel) Reset()                    { *m = FlyTimePanel{} }
func (m *FlyTimePanel) String() string            { return proto.CompactTextString(m) }
func (*FlyTimePanel) ProtoMessage()               {}
func (*FlyTimePanel) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *FlyTimePanel) GetData() []float64 {
	if m != nil {
		return m.Data
	}
	return nil
}

func (m *FlyTimePanel) GetDates() []uint64 {
	if m != nil {
		return m.Dates
	}
	return nil
}

func (m *FlyTimePanel) GetSecondary() []string {
	if m != nil {
		return m.Secondary
	}
	return nil
}

func (m *FlyTimePanel) GetThirdly() []string {
	if m != nil {
		return m.Thirdly
	}
	return nil
}

func init() {
	proto.RegisterType((*FlyTimePanel)(nil), "gopandas.FlyTimePanel")
}

var fileDescriptor0 = []byte{
	// 121 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe2, 0xe2, 0x4b, 0xcf, 0x2f, 0x48,
	0xcc, 0x4b, 0x49, 0x2c, 0xd6, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x80, 0xf1, 0x95, 0xc2,
	0xb9, 0x78, 0xdc, 0x72, 0x2a, 0x43, 0x32, 0x73, 0x53, 0x03, 0x12, 0xf3, 0x52, 0x73, 0x84, 0x04,
	0xb8, 0x58, 0x52, 0x12, 0x4b, 0x12, 0x25, 0x18, 0x15, 0x98, 0x35, 0x18, 0x9d, 0x98, 0x04, 0x18,
	0x85, 0x04, 0xb9, 0x58, 0x81, 0x22, 0xa9, 0xc5, 0x12, 0x4c, 0x40, 0x21, 0x16, 0xa8, 0x10, 0x67,
	0x71, 0x6a, 0x72, 0x3e, 0xd0, 0x80, 0xa2, 0x4a, 0x09, 0x66, 0xa0, 0x30, 0xa7, 0x10, 0x3f, 0x17,
	0x7b, 0x49, 0x46, 0x66, 0x51, 0x4a, 0x4e, 0xa5, 0x04, 0x0b, 0x48, 0x00, 0x10, 0x00, 0x00, 0xff,
	0xff, 0x69, 0x93, 0x5b, 0x41, 0x73, 0x00, 0x00, 0x00,
}
