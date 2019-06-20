// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package sdfsprotocol

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type NetINodeNetBlockInfoResponse struct {
	_tab flatbuffers.Table
}

func GetRootAsNetINodeNetBlockInfoResponse(buf []byte, offset flatbuffers.UOffsetT) *NetINodeNetBlockInfoResponse {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &NetINodeNetBlockInfoResponse{}
	x.Init(buf, n+offset)
	return x
}

func (rcv *NetINodeNetBlockInfoResponse) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *NetINodeNetBlockInfoResponse) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *NetINodeNetBlockInfoResponse) CommonResponse(obj *CommonResponse) *CommonResponse {
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

func (rcv *NetINodeNetBlockInfoResponse) Len() int32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.GetInt32(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *NetINodeNetBlockInfoResponse) MutateLen(n int32) bool {
	return rcv._tab.MutateInt32Slot(6, n)
}

func (rcv *NetINodeNetBlockInfoResponse) Cap() int32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.GetInt32(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *NetINodeNetBlockInfoResponse) MutateCap(n int32) bool {
	return rcv._tab.MutateInt32Slot(8, n)
}

func (rcv *NetINodeNetBlockInfoResponse) Backends(j int) []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.ByteVector(a + flatbuffers.UOffsetT(j*4))
	}
	return nil
}

func (rcv *NetINodeNetBlockInfoResponse) BackendsLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func NetINodeNetBlockInfoResponseStart(builder *flatbuffers.Builder) {
	builder.StartObject(4)
}
func NetINodeNetBlockInfoResponseAddCommonResponse(builder *flatbuffers.Builder, CommonResponse flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(CommonResponse), 0)
}
func NetINodeNetBlockInfoResponseAddLen(builder *flatbuffers.Builder, Len int32) {
	builder.PrependInt32Slot(1, Len, 0)
}
func NetINodeNetBlockInfoResponseAddCap(builder *flatbuffers.Builder, Cap int32) {
	builder.PrependInt32Slot(2, Cap, 0)
}
func NetINodeNetBlockInfoResponseAddBackends(builder *flatbuffers.Builder, Backends flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(3, flatbuffers.UOffsetT(Backends), 0)
}
func NetINodeNetBlockInfoResponseStartBackendsVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(4, numElems, 4)
}
func NetINodeNetBlockInfoResponseEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
