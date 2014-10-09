// Code generated by protoc-gen-gogo.
// source: message.proto
// DO NOT EDIT!

/*
	Package message is a generated protocol buffer package.

	Generate: protoc --gogo_out=. -I=.:../../dependencies/code.google.com/p/gogoprotobuf/gogoproto:../../dependencies/code.google.com/p/gogoprotobuf/protobuf message.proto

	It is generated from these files:
		message.proto

	It has these top-level messages:
		State
		Event
		Message
		Attribute
*/
package message

import proto "code.google.com/p/gogoprotobuf/proto"
import math "math"

// discarding unused import gogoproto "gogo.pb"

import io "io"
import math1 "math"
import fmt "fmt"
import code_google_com_p_gogoprotobuf_proto "code.google.com/p/gogoprotobuf/proto"

import math2 "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = math.Inf

type State struct {
	Time             *int64   `protobuf:"varint,1,opt,name=time" json:"time,omitempty"`
	State            *string  `protobuf:"bytes,2,opt,name=state" json:"state,omitempty"`
	Service          *string  `protobuf:"bytes,3,opt,name=service" json:"service,omitempty"`
	Host             *string  `protobuf:"bytes,4,opt,name=host" json:"host,omitempty"`
	Description      *string  `protobuf:"bytes,5,opt,name=description" json:"description,omitempty"`
	Once             *bool    `protobuf:"varint,6,opt,name=once" json:"once,omitempty"`
	Tags             []string `protobuf:"bytes,7,rep,name=tags" json:"tags,omitempty"`
	Ttl              *float32 `protobuf:"fixed32,8,opt,name=ttl" json:"ttl,omitempty"`
	XXX_unrecognized []byte   `json:"-"`
}

func (m *State) Reset()         { *m = State{} }
func (m *State) String() string { return proto.CompactTextString(m) }
func (*State) ProtoMessage()    {}

func (m *State) GetTime() int64 {
	if m != nil && m.Time != nil {
		return *m.Time
	}
	return 0
}

func (m *State) GetState() string {
	if m != nil && m.State != nil {
		return *m.State
	}
	return ""
}

func (m *State) GetService() string {
	if m != nil && m.Service != nil {
		return *m.Service
	}
	return ""
}

func (m *State) GetHost() string {
	if m != nil && m.Host != nil {
		return *m.Host
	}
	return ""
}

func (m *State) GetDescription() string {
	if m != nil && m.Description != nil {
		return *m.Description
	}
	return ""
}

func (m *State) GetOnce() bool {
	if m != nil && m.Once != nil {
		return *m.Once
	}
	return false
}

func (m *State) GetTags() []string {
	if m != nil {
		return m.Tags
	}
	return nil
}

func (m *State) GetTtl() float32 {
	if m != nil && m.Ttl != nil {
		return *m.Ttl
	}
	return 0
}

type Event struct {
	Time             *int64       `protobuf:"varint,1,opt,name=time" json:"time,omitempty"`
	State            *string      `protobuf:"bytes,2,opt,name=state" json:"state,omitempty"`
	Service          *string      `protobuf:"bytes,3,opt,name=service" json:"service,omitempty"`
	Host             *string      `protobuf:"bytes,4,opt,name=host" json:"host,omitempty"`
	Description      *string      `protobuf:"bytes,5,opt,name=description" json:"description,omitempty"`
	Tags             []string     `protobuf:"bytes,7,rep,name=tags" json:"tags,omitempty"`
	Ttl              *float32     `protobuf:"fixed32,8,opt,name=ttl" json:"ttl,omitempty"`
	Attributes       []*Attribute `protobuf:"bytes,9,rep,name=attributes" json:"attributes,omitempty"`
	MetricSint64     *int64       `protobuf:"zigzag64,13,opt,name=metric_sint64" json:"metric_sint64,omitempty"`
	MetricD          *float64     `protobuf:"fixed64,14,opt,name=metric_d" json:"metric_d,omitempty"`
	MetricF          *float32     `protobuf:"fixed32,15,opt,name=metric_f" json:"metric_f,omitempty"`
	TsdbService      *string      `protobuf:"bytes,16,opt,name=tsdb_service" json:"tsdb_service,omitempty"`
	TsdbTags         []string     `protobuf:"bytes,17,rep,name=tsdb_tags" json:"tsdb_tags,omitempty"`
	XXX_unrecognized []byte       `json:"-"`
}

func (m *Event) Reset()         { *m = Event{} }
func (m *Event) String() string { return proto.CompactTextString(m) }
func (*Event) ProtoMessage()    {}

func (m *Event) GetTime() int64 {
	if m != nil && m.Time != nil {
		return *m.Time
	}
	return 0
}

func (m *Event) GetState() string {
	if m != nil && m.State != nil {
		return *m.State
	}
	return ""
}

func (m *Event) GetService() string {
	if m != nil && m.Service != nil {
		return *m.Service
	}
	return ""
}

func (m *Event) GetHost() string {
	if m != nil && m.Host != nil {
		return *m.Host
	}
	return ""
}

func (m *Event) GetDescription() string {
	if m != nil && m.Description != nil {
		return *m.Description
	}
	return ""
}

func (m *Event) GetTags() []string {
	if m != nil {
		return m.Tags
	}
	return nil
}

func (m *Event) GetTtl() float32 {
	if m != nil && m.Ttl != nil {
		return *m.Ttl
	}
	return 0
}

func (m *Event) GetAttributes() []*Attribute {
	if m != nil {
		return m.Attributes
	}
	return nil
}

func (m *Event) GetMetricSint64() int64 {
	if m != nil && m.MetricSint64 != nil {
		return *m.MetricSint64
	}
	return 0
}

func (m *Event) GetMetricD() float64 {
	if m != nil && m.MetricD != nil {
		return *m.MetricD
	}
	return 0
}

func (m *Event) GetMetricF() float32 {
	if m != nil && m.MetricF != nil {
		return *m.MetricF
	}
	return 0
}

func (m *Event) GetTsdbService() string {
	if m != nil && m.TsdbService != nil {
		return *m.TsdbService
	}
	return ""
}

func (m *Event) GetTsdbTags() []string {
	if m != nil {
		return m.TsdbTags
	}
	return nil
}

type Message struct {
	Ok               *bool    `protobuf:"varint,2,opt,name=ok" json:"ok,omitempty"`
	Error            *string  `protobuf:"bytes,3,opt,name=error" json:"error,omitempty"`
	States           []*State `protobuf:"bytes,4,rep,name=states" json:"states,omitempty"`
	Events           []*Event `protobuf:"bytes,6,rep,name=events" json:"events,omitempty"`
	Data             *string  `protobuf:"bytes,7,opt,name=data" json:"data,omitempty"`
	DataType         *string  `protobuf:"bytes,8,opt,name=data_type" json:"data_type,omitempty"`
	XXX_unrecognized []byte   `json:"-"`
}

func (m *Message) Reset()         { *m = Message{} }
func (m *Message) String() string { return proto.CompactTextString(m) }
func (*Message) ProtoMessage()    {}

func (m *Message) GetOk() bool {
	if m != nil && m.Ok != nil {
		return *m.Ok
	}
	return false
}

func (m *Message) GetError() string {
	if m != nil && m.Error != nil {
		return *m.Error
	}
	return ""
}

func (m *Message) GetStates() []*State {
	if m != nil {
		return m.States
	}
	return nil
}

func (m *Message) GetEvents() []*Event {
	if m != nil {
		return m.Events
	}
	return nil
}

func (m *Message) GetData() string {
	if m != nil && m.Data != nil {
		return *m.Data
	}
	return ""
}

func (m *Message) GetDataType() string {
	if m != nil && m.DataType != nil {
		return *m.DataType
	}
	return ""
}

type Attribute struct {
	Key              *string `protobuf:"bytes,1,req,name=key" json:"key,omitempty"`
	Value            *string `protobuf:"bytes,2,opt,name=value" json:"value,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *Attribute) Reset()         { *m = Attribute{} }
func (m *Attribute) String() string { return proto.CompactTextString(m) }
func (*Attribute) ProtoMessage()    {}

func (m *Attribute) GetKey() string {
	if m != nil && m.Key != nil {
		return *m.Key
	}
	return ""
}

func (m *Attribute) GetValue() string {
	if m != nil && m.Value != nil {
		return *m.Value
	}
	return ""
}

func init() {
}
func (m *State) Unmarshal(data []byte) error {
	l := len(data)
	index := 0
	for index < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if index >= l {
				return io.ErrUnexpectedEOF
			}
			b := data[index]
			index++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Time", wireType)
			}
			var v int64
			for shift := uint(0); ; shift += 7 {
				if index >= l {
					return io.ErrUnexpectedEOF
				}
				b := data[index]
				index++
				v |= (int64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.Time = &v
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field State", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if index >= l {
					return io.ErrUnexpectedEOF
				}
				b := data[index]
				index++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			postIndex := index + int(stringLen)
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			s := string(data[index:postIndex])
			m.State = &s
			index = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Service", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if index >= l {
					return io.ErrUnexpectedEOF
				}
				b := data[index]
				index++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			postIndex := index + int(stringLen)
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			s := string(data[index:postIndex])
			m.Service = &s
			index = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Host", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if index >= l {
					return io.ErrUnexpectedEOF
				}
				b := data[index]
				index++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			postIndex := index + int(stringLen)
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			s := string(data[index:postIndex])
			m.Host = &s
			index = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Description", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if index >= l {
					return io.ErrUnexpectedEOF
				}
				b := data[index]
				index++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			postIndex := index + int(stringLen)
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			s := string(data[index:postIndex])
			m.Description = &s
			index = postIndex
		case 6:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Once", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if index >= l {
					return io.ErrUnexpectedEOF
				}
				b := data[index]
				index++
				v |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			b := bool(v != 0)
			m.Once = &b
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Tags", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if index >= l {
					return io.ErrUnexpectedEOF
				}
				b := data[index]
				index++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			postIndex := index + int(stringLen)
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Tags = append(m.Tags, string(data[index:postIndex]))
			index = postIndex
		case 8:
			if wireType != 5 {
				return fmt.Errorf("proto: wrong wireType = %d for field Ttl", wireType)
			}
			var v uint32
			i := index + 4
			if i > l {
				return io.ErrUnexpectedEOF
			}
			index = i
			v = uint32(data[i-4])
			v |= uint32(data[i-3]) << 8
			v |= uint32(data[i-2]) << 16
			v |= uint32(data[i-1]) << 24
			v2 := math1.Float32frombits(v)
			m.Ttl = &v2
		default:
			var sizeOfWire int
			for {
				sizeOfWire++
				wire >>= 7
				if wire == 0 {
					break
				}
			}
			index -= sizeOfWire
			skippy, err := code_google_com_p_gogoprotobuf_proto.Skip(data[index:])
			if err != nil {
				return err
			}
			if (index + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, data[index:index+skippy]...)
			index += skippy
		}
	}
	return nil
}
func (m *Event) Unmarshal(data []byte) error {
	l := len(data)
	index := 0
	for index < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if index >= l {
				return io.ErrUnexpectedEOF
			}
			b := data[index]
			index++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Time", wireType)
			}
			var v int64
			for shift := uint(0); ; shift += 7 {
				if index >= l {
					return io.ErrUnexpectedEOF
				}
				b := data[index]
				index++
				v |= (int64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.Time = &v
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field State", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if index >= l {
					return io.ErrUnexpectedEOF
				}
				b := data[index]
				index++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			postIndex := index + int(stringLen)
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			s := string(data[index:postIndex])
			m.State = &s
			index = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Service", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if index >= l {
					return io.ErrUnexpectedEOF
				}
				b := data[index]
				index++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			postIndex := index + int(stringLen)
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			s := string(data[index:postIndex])
			m.Service = &s
			index = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Host", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if index >= l {
					return io.ErrUnexpectedEOF
				}
				b := data[index]
				index++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			postIndex := index + int(stringLen)
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			s := string(data[index:postIndex])
			m.Host = &s
			index = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Description", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if index >= l {
					return io.ErrUnexpectedEOF
				}
				b := data[index]
				index++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			postIndex := index + int(stringLen)
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			s := string(data[index:postIndex])
			m.Description = &s
			index = postIndex
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Tags", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if index >= l {
					return io.ErrUnexpectedEOF
				}
				b := data[index]
				index++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			postIndex := index + int(stringLen)
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Tags = append(m.Tags, string(data[index:postIndex]))
			index = postIndex
		case 8:
			if wireType != 5 {
				return fmt.Errorf("proto: wrong wireType = %d for field Ttl", wireType)
			}
			var v uint32
			i := index + 4
			if i > l {
				return io.ErrUnexpectedEOF
			}
			index = i
			v = uint32(data[i-4])
			v |= uint32(data[i-3]) << 8
			v |= uint32(data[i-2]) << 16
			v |= uint32(data[i-1]) << 24
			v2 := math1.Float32frombits(v)
			m.Ttl = &v2
		case 9:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Attributes", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if index >= l {
					return io.ErrUnexpectedEOF
				}
				b := data[index]
				index++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			postIndex := index + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Attributes = append(m.Attributes, &Attribute{})
			m.Attributes[len(m.Attributes)-1].Unmarshal(data[index:postIndex])
			index = postIndex
		case 13:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field MetricSint64", wireType)
			}
			var v uint64
			for shift := uint(0); ; shift += 7 {
				if index >= l {
					return io.ErrUnexpectedEOF
				}
				b := data[index]
				index++
				v |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			v = (v >> 1) ^ uint64((int64(v&1)<<63)>>63)
			v2 := int64(v)
			m.MetricSint64 = &v2
		case 14:
			if wireType != 1 {
				return fmt.Errorf("proto: wrong wireType = %d for field MetricD", wireType)
			}
			var v uint64
			i := index + 8
			if i > l {
				return io.ErrUnexpectedEOF
			}
			index = i
			v = uint64(data[i-8])
			v |= uint64(data[i-7]) << 8
			v |= uint64(data[i-6]) << 16
			v |= uint64(data[i-5]) << 24
			v |= uint64(data[i-4]) << 32
			v |= uint64(data[i-3]) << 40
			v |= uint64(data[i-2]) << 48
			v |= uint64(data[i-1]) << 56
			v2 := math1.Float64frombits(v)
			m.MetricD = &v2
		case 15:
			if wireType != 5 {
				return fmt.Errorf("proto: wrong wireType = %d for field MetricF", wireType)
			}
			var v uint32
			i := index + 4
			if i > l {
				return io.ErrUnexpectedEOF
			}
			index = i
			v = uint32(data[i-4])
			v |= uint32(data[i-3]) << 8
			v |= uint32(data[i-2]) << 16
			v |= uint32(data[i-1]) << 24
			v2 := math1.Float32frombits(v)
			m.MetricF = &v2
		case 16:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TsdbService", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if index >= l {
					return io.ErrUnexpectedEOF
				}
				b := data[index]
				index++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			postIndex := index + int(stringLen)
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			s := string(data[index:postIndex])
			m.TsdbService = &s
			index = postIndex
		case 17:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TsdbTags", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if index >= l {
					return io.ErrUnexpectedEOF
				}
				b := data[index]
				index++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			postIndex := index + int(stringLen)
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.TsdbTags = append(m.TsdbTags, string(data[index:postIndex]))
			index = postIndex
		default:
			var sizeOfWire int
			for {
				sizeOfWire++
				wire >>= 7
				if wire == 0 {
					break
				}
			}
			index -= sizeOfWire
			skippy, err := code_google_com_p_gogoprotobuf_proto.Skip(data[index:])
			if err != nil {
				return err
			}
			if (index + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, data[index:index+skippy]...)
			index += skippy
		}
	}
	return nil
}
func (m *Message) Unmarshal(data []byte) error {
	l := len(data)
	index := 0
	for index < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if index >= l {
				return io.ErrUnexpectedEOF
			}
			b := data[index]
			index++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		switch fieldNum {
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Ok", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if index >= l {
					return io.ErrUnexpectedEOF
				}
				b := data[index]
				index++
				v |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			b := bool(v != 0)
			m.Ok = &b
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Error", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if index >= l {
					return io.ErrUnexpectedEOF
				}
				b := data[index]
				index++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			postIndex := index + int(stringLen)
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			s := string(data[index:postIndex])
			m.Error = &s
			index = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field States", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if index >= l {
					return io.ErrUnexpectedEOF
				}
				b := data[index]
				index++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			postIndex := index + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.States = append(m.States, &State{})
			m.States[len(m.States)-1].Unmarshal(data[index:postIndex])
			index = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Events", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if index >= l {
					return io.ErrUnexpectedEOF
				}
				b := data[index]
				index++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			postIndex := index + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Events = append(m.Events, &Event{})
			m.Events[len(m.Events)-1].Unmarshal(data[index:postIndex])
			index = postIndex
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Data", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if index >= l {
					return io.ErrUnexpectedEOF
				}
				b := data[index]
				index++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			postIndex := index + int(stringLen)
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			s := string(data[index:postIndex])
			m.Data = &s
			index = postIndex
		case 8:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DataType", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if index >= l {
					return io.ErrUnexpectedEOF
				}
				b := data[index]
				index++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			postIndex := index + int(stringLen)
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			s := string(data[index:postIndex])
			m.DataType = &s
			index = postIndex
		default:
			var sizeOfWire int
			for {
				sizeOfWire++
				wire >>= 7
				if wire == 0 {
					break
				}
			}
			index -= sizeOfWire
			skippy, err := code_google_com_p_gogoprotobuf_proto.Skip(data[index:])
			if err != nil {
				return err
			}
			if (index + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, data[index:index+skippy]...)
			index += skippy
		}
	}
	return nil
}
func (m *Attribute) Unmarshal(data []byte) error {
	l := len(data)
	index := 0
	for index < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if index >= l {
				return io.ErrUnexpectedEOF
			}
			b := data[index]
			index++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Key", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if index >= l {
					return io.ErrUnexpectedEOF
				}
				b := data[index]
				index++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			postIndex := index + int(stringLen)
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			s := string(data[index:postIndex])
			m.Key = &s
			index = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Value", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if index >= l {
					return io.ErrUnexpectedEOF
				}
				b := data[index]
				index++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			postIndex := index + int(stringLen)
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			s := string(data[index:postIndex])
			m.Value = &s
			index = postIndex
		default:
			var sizeOfWire int
			for {
				sizeOfWire++
				wire >>= 7
				if wire == 0 {
					break
				}
			}
			index -= sizeOfWire
			skippy, err := code_google_com_p_gogoprotobuf_proto.Skip(data[index:])
			if err != nil {
				return err
			}
			if (index + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, data[index:index+skippy]...)
			index += skippy
		}
	}
	return nil
}
func (m *State) Size() (n int) {
	var l int
	_ = l
	if m.Time != nil {
		n += 1 + sovMessage(uint64(*m.Time))
	}
	if m.State != nil {
		l = len(*m.State)
		n += 1 + l + sovMessage(uint64(l))
	}
	if m.Service != nil {
		l = len(*m.Service)
		n += 1 + l + sovMessage(uint64(l))
	}
	if m.Host != nil {
		l = len(*m.Host)
		n += 1 + l + sovMessage(uint64(l))
	}
	if m.Description != nil {
		l = len(*m.Description)
		n += 1 + l + sovMessage(uint64(l))
	}
	if m.Once != nil {
		n += 2
	}
	if len(m.Tags) > 0 {
		for _, s := range m.Tags {
			l = len(s)
			n += 1 + l + sovMessage(uint64(l))
		}
	}
	if m.Ttl != nil {
		n += 5
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}
func (m *Event) Size() (n int) {
	var l int
	_ = l
	if m.Time != nil {
		n += 1 + sovMessage(uint64(*m.Time))
	}
	if m.State != nil {
		l = len(*m.State)
		n += 1 + l + sovMessage(uint64(l))
	}
	if m.Service != nil {
		l = len(*m.Service)
		n += 1 + l + sovMessage(uint64(l))
	}
	if m.Host != nil {
		l = len(*m.Host)
		n += 1 + l + sovMessage(uint64(l))
	}
	if m.Description != nil {
		l = len(*m.Description)
		n += 1 + l + sovMessage(uint64(l))
	}
	if len(m.Tags) > 0 {
		for _, s := range m.Tags {
			l = len(s)
			n += 1 + l + sovMessage(uint64(l))
		}
	}
	if m.Ttl != nil {
		n += 5
	}
	if len(m.Attributes) > 0 {
		for _, e := range m.Attributes {
			l = e.Size()
			n += 1 + l + sovMessage(uint64(l))
		}
	}
	if m.MetricSint64 != nil {
		n += 1 + sozMessage(uint64(*m.MetricSint64))
	}
	if m.MetricD != nil {
		n += 9
	}
	if m.MetricF != nil {
		n += 5
	}
	if m.TsdbService != nil {
		l = len(*m.TsdbService)
		n += 2 + l + sovMessage(uint64(l))
	}
	if len(m.TsdbTags) > 0 {
		for _, s := range m.TsdbTags {
			l = len(s)
			n += 2 + l + sovMessage(uint64(l))
		}
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}
func (m *Message) Size() (n int) {
	var l int
	_ = l
	if m.Ok != nil {
		n += 2
	}
	if m.Error != nil {
		l = len(*m.Error)
		n += 1 + l + sovMessage(uint64(l))
	}
	if len(m.States) > 0 {
		for _, e := range m.States {
			l = e.Size()
			n += 1 + l + sovMessage(uint64(l))
		}
	}
	if len(m.Events) > 0 {
		for _, e := range m.Events {
			l = e.Size()
			n += 1 + l + sovMessage(uint64(l))
		}
	}
	if m.Data != nil {
		l = len(*m.Data)
		n += 1 + l + sovMessage(uint64(l))
	}
	if m.DataType != nil {
		l = len(*m.DataType)
		n += 1 + l + sovMessage(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}
func (m *Attribute) Size() (n int) {
	var l int
	_ = l
	if m.Key != nil {
		l = len(*m.Key)
		n += 1 + l + sovMessage(uint64(l))
	}
	if m.Value != nil {
		l = len(*m.Value)
		n += 1 + l + sovMessage(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func sovMessage(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozMessage(x uint64) (n int) {
	return sovMessage(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *State) Marshal() (data []byte, err error) {
	size := m.Size()
	data = make([]byte, size)
	n, err := m.MarshalTo(data)
	if err != nil {
		return nil, err
	}
	return data[:n], nil
}

func (m *State) MarshalTo(data []byte) (n int, err error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.Time != nil {
		data[i] = 0x8
		i++
		i = encodeVarintMessage(data, i, uint64(*m.Time))
	}
	if m.State != nil {
		data[i] = 0x12
		i++
		i = encodeVarintMessage(data, i, uint64(len(*m.State)))
		i += copy(data[i:], *m.State)
	}
	if m.Service != nil {
		data[i] = 0x1a
		i++
		i = encodeVarintMessage(data, i, uint64(len(*m.Service)))
		i += copy(data[i:], *m.Service)
	}
	if m.Host != nil {
		data[i] = 0x22
		i++
		i = encodeVarintMessage(data, i, uint64(len(*m.Host)))
		i += copy(data[i:], *m.Host)
	}
	if m.Description != nil {
		data[i] = 0x2a
		i++
		i = encodeVarintMessage(data, i, uint64(len(*m.Description)))
		i += copy(data[i:], *m.Description)
	}
	if m.Once != nil {
		data[i] = 0x30
		i++
		if *m.Once {
			data[i] = 1
		} else {
			data[i] = 0
		}
		i++
	}
	if len(m.Tags) > 0 {
		for _, s := range m.Tags {
			data[i] = 0x3a
			i++
			l = len(s)
			for l >= 1<<7 {
				data[i] = uint8(uint64(l)&0x7f | 0x80)
				l >>= 7
				i++
			}
			data[i] = uint8(l)
			i++
			i += copy(data[i:], s)
		}
	}
	if m.Ttl != nil {
		data[i] = 0x45
		i++
		i = encodeFixed32Message(data, i, uint32(math2.Float32bits(*m.Ttl)))
	}
	if m.XXX_unrecognized != nil {
		i += copy(data[i:], m.XXX_unrecognized)
	}
	return i, nil
}
func (m *Event) Marshal() (data []byte, err error) {
	size := m.Size()
	data = make([]byte, size)
	n, err := m.MarshalTo(data)
	if err != nil {
		return nil, err
	}
	return data[:n], nil
}

func (m *Event) MarshalTo(data []byte) (n int, err error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.Time != nil {
		data[i] = 0x8
		i++
		i = encodeVarintMessage(data, i, uint64(*m.Time))
	}
	if m.State != nil {
		data[i] = 0x12
		i++
		i = encodeVarintMessage(data, i, uint64(len(*m.State)))
		i += copy(data[i:], *m.State)
	}
	if m.Service != nil {
		data[i] = 0x1a
		i++
		i = encodeVarintMessage(data, i, uint64(len(*m.Service)))
		i += copy(data[i:], *m.Service)
	}
	if m.Host != nil {
		data[i] = 0x22
		i++
		i = encodeVarintMessage(data, i, uint64(len(*m.Host)))
		i += copy(data[i:], *m.Host)
	}
	if m.Description != nil {
		data[i] = 0x2a
		i++
		i = encodeVarintMessage(data, i, uint64(len(*m.Description)))
		i += copy(data[i:], *m.Description)
	}
	if len(m.Tags) > 0 {
		for _, s := range m.Tags {
			data[i] = 0x3a
			i++
			l = len(s)
			for l >= 1<<7 {
				data[i] = uint8(uint64(l)&0x7f | 0x80)
				l >>= 7
				i++
			}
			data[i] = uint8(l)
			i++
			i += copy(data[i:], s)
		}
	}
	if m.Ttl != nil {
		data[i] = 0x45
		i++
		i = encodeFixed32Message(data, i, uint32(math2.Float32bits(*m.Ttl)))
	}
	if len(m.Attributes) > 0 {
		for _, msg := range m.Attributes {
			data[i] = 0x4a
			i++
			i = encodeVarintMessage(data, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(data[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	if m.MetricSint64 != nil {
		data[i] = 0x68
		i++
		i = encodeVarintMessage(data, i, uint64((uint64(*m.MetricSint64)<<1)^uint64((*m.MetricSint64>>63))))
	}
	if m.MetricD != nil {
		data[i] = 0x71
		i++
		i = encodeFixed64Message(data, i, uint64(math2.Float64bits(*m.MetricD)))
	}
	if m.MetricF != nil {
		data[i] = 0x7d
		i++
		i = encodeFixed32Message(data, i, uint32(math2.Float32bits(*m.MetricF)))
	}
	if m.TsdbService != nil {
		data[i] = 0x82
		i++
		data[i] = 0x1
		i++
		i = encodeVarintMessage(data, i, uint64(len(*m.TsdbService)))
		i += copy(data[i:], *m.TsdbService)
	}
	if len(m.TsdbTags) > 0 {
		for _, s := range m.TsdbTags {
			data[i] = 0x8a
			i++
			data[i] = 0x1
			i++
			l = len(s)
			for l >= 1<<7 {
				data[i] = uint8(uint64(l)&0x7f | 0x80)
				l >>= 7
				i++
			}
			data[i] = uint8(l)
			i++
			i += copy(data[i:], s)
		}
	}
	if m.XXX_unrecognized != nil {
		i += copy(data[i:], m.XXX_unrecognized)
	}
	return i, nil
}
func (m *Message) Marshal() (data []byte, err error) {
	size := m.Size()
	data = make([]byte, size)
	n, err := m.MarshalTo(data)
	if err != nil {
		return nil, err
	}
	return data[:n], nil
}

func (m *Message) MarshalTo(data []byte) (n int, err error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.Ok != nil {
		data[i] = 0x10
		i++
		if *m.Ok {
			data[i] = 1
		} else {
			data[i] = 0
		}
		i++
	}
	if m.Error != nil {
		data[i] = 0x1a
		i++
		i = encodeVarintMessage(data, i, uint64(len(*m.Error)))
		i += copy(data[i:], *m.Error)
	}
	if len(m.States) > 0 {
		for _, msg := range m.States {
			data[i] = 0x22
			i++
			i = encodeVarintMessage(data, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(data[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	if len(m.Events) > 0 {
		for _, msg := range m.Events {
			data[i] = 0x32
			i++
			i = encodeVarintMessage(data, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(data[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	if m.Data != nil {
		data[i] = 0x3a
		i++
		i = encodeVarintMessage(data, i, uint64(len(*m.Data)))
		i += copy(data[i:], *m.Data)
	}
	if m.DataType != nil {
		data[i] = 0x42
		i++
		i = encodeVarintMessage(data, i, uint64(len(*m.DataType)))
		i += copy(data[i:], *m.DataType)
	}
	if m.XXX_unrecognized != nil {
		i += copy(data[i:], m.XXX_unrecognized)
	}
	return i, nil
}
func (m *Attribute) Marshal() (data []byte, err error) {
	size := m.Size()
	data = make([]byte, size)
	n, err := m.MarshalTo(data)
	if err != nil {
		return nil, err
	}
	return data[:n], nil
}

func (m *Attribute) MarshalTo(data []byte) (n int, err error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.Key != nil {
		data[i] = 0xa
		i++
		i = encodeVarintMessage(data, i, uint64(len(*m.Key)))
		i += copy(data[i:], *m.Key)
	}
	if m.Value != nil {
		data[i] = 0x12
		i++
		i = encodeVarintMessage(data, i, uint64(len(*m.Value)))
		i += copy(data[i:], *m.Value)
	}
	if m.XXX_unrecognized != nil {
		i += copy(data[i:], m.XXX_unrecognized)
	}
	return i, nil
}
func encodeFixed64Message(data []byte, offset int, v uint64) int {
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
func encodeFixed32Message(data []byte, offset int, v uint32) int {
	data[offset] = uint8(v)
	data[offset+1] = uint8(v >> 8)
	data[offset+2] = uint8(v >> 16)
	data[offset+3] = uint8(v >> 24)
	return offset + 4
}
func encodeVarintMessage(data []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		data[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	data[offset] = uint8(v)
	return offset + 1
}