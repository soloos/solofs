// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package protocol

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type INodeNetBlockInfoResponse struct {
	_tab flatbuffers.Table
}

func GetRootAsINodeNetBlockInfoResponse(buf []byte, offset flatbuffers.UOffsetT) *INodeNetBlockInfoResponse {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &INodeNetBlockInfoResponse{}
	x.Init(buf, n+offset)
	return x
}

func (rcv *INodeNetBlockInfoResponse) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *INodeNetBlockInfoResponse) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *INodeNetBlockInfoResponse) CommonResponse(obj *CommonResponse) *CommonResponse {
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

func (rcv *INodeNetBlockInfoResponse) NetBlockID() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *INodeNetBlockInfoResponse) Len() int32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.GetInt32(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *INodeNetBlockInfoResponse) MutateLen(n int32) bool {
	return rcv._tab.MutateInt32Slot(8, n)
}

func (rcv *INodeNetBlockInfoResponse) Cap() int32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		return rcv._tab.GetInt32(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *INodeNetBlockInfoResponse) MutateCap(n int32) bool {
	return rcv._tab.MutateInt32Slot(10, n)
}

func (rcv *INodeNetBlockInfoResponse) Backends(obj *NetBlockBackend, j int) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(12))
	if o != 0 {
		x := rcv._tab.Vector(o)
		x += flatbuffers.UOffsetT(j) * 4
		x = rcv._tab.Indirect(x)
		obj.Init(rcv._tab.Bytes, x)
		return true
	}
	return false
}

func (rcv *INodeNetBlockInfoResponse) BackendsLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(12))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func INodeNetBlockInfoResponseStart(builder *flatbuffers.Builder) {
	builder.StartObject(5)
}
func INodeNetBlockInfoResponseAddCommonResponse(builder *flatbuffers.Builder, CommonResponse flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(CommonResponse), 0)
}
func INodeNetBlockInfoResponseAddNetBlockID(builder *flatbuffers.Builder, NetBlockID flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(1, flatbuffers.UOffsetT(NetBlockID), 0)
}
func INodeNetBlockInfoResponseAddLen(builder *flatbuffers.Builder, Len int32) {
	builder.PrependInt32Slot(2, Len, 0)
}
func INodeNetBlockInfoResponseAddCap(builder *flatbuffers.Builder, Cap int32) {
	builder.PrependInt32Slot(3, Cap, 0)
}
func INodeNetBlockInfoResponseAddBackends(builder *flatbuffers.Builder, Backends flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(4, flatbuffers.UOffsetT(Backends), 0)
}
func INodeNetBlockInfoResponseStartBackendsVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(4, numElems, 4)
}
func INodeNetBlockInfoResponseEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
