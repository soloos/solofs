// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package protocol

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type NetBlockPReadResponse struct {
	_tab flatbuffers.Table
}

func GetRootAsNetBlockPReadResponse(buf []byte, offset flatbuffers.UOffsetT) *NetBlockPReadResponse {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &NetBlockPReadResponse{}
	x.Init(buf, n+offset)
	return x
}

func (rcv *NetBlockPReadResponse) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *NetBlockPReadResponse) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *NetBlockPReadResponse) CommonResponse(obj *CommonResponse) *CommonResponse {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		x := rcv._tab.Indirect(o + rcv._tab.Pos)
		if obj == nil {
			obj = new(CommonResponse)
		}
		obj.Init(rcv._tab.Bytes, x)
		return obj
	}
	return nil
}

func (rcv *NetBlockPReadResponse) Length() int32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.GetInt32(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *NetBlockPReadResponse) MutateLength(n int32) bool {
	return rcv._tab.MutateInt32Slot(6, n)
}

func NetBlockPReadResponseStart(builder *flatbuffers.Builder) {
	builder.StartObject(2)
}
func NetBlockPReadResponseAddCommonResponse(builder *flatbuffers.Builder, CommonResponse flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(CommonResponse), 0)
}
func NetBlockPReadResponseAddLength(builder *flatbuffers.Builder, Length int32) {
	builder.PrependInt32Slot(1, Length, 0)
}
func NetBlockPReadResponseEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}