package store

import (
	"bytes"
	"math"
)

type Operation struct {
	Optype     byte
	Action     byte
	Timestamp  uint32
	Item       []byte
	ItemLength uint32
}

func (op *Operation) SetOPType() {
	switch op.Action {
	case OP_REFRESH:
		op.Optype = RED
	case OP_ADD:
		op.Optype = BLUE
	case OP_REMOVE:
		op.Optype = RED
	case OP_CLEAR:
		op.Optype = RED
	case OP_SETTLE:
		op.Optype = RED
	}
}

func NewOperation(action byte) *Operation {
	op := &Operation{Action: action}
	op.SetOPType()
	return op
}

func (op *Operation) MarshalBinary() ([]byte, error) {
	bs := &bytes.Buffer{}
	bs.WriteByte(op.Optype)
	bs.WriteByte(op.Action)
	return bs.Bytes(), nil
}

func (op *Operation) UnMarshalBinary(d []byte) []byte {
	bs := bytes.NewBuffer(d)
	op.Optype = bs.Next(1)[0]
	op.Action = bs.Next(1)[0]

	return bs.Next(math.MaxInt32)
}

type OperationSlice []Operation

func (slice OperationSlice) Len() int {
	return len(slice)
}

func (slice OperationSlice) AddOperation(op Operation) OperationSlice {
	for i, t := range slice {
		if t.Timestamp > op.Timestamp {
			slice = append(append(slice[:i], op), slice[i:]...)
			return slice
		}
	}
	return append(slice, op)
}

func (slice *OperationSlice) MarshalBinary() ([]byte, error) {
	bf := &bytes.Buffer{}

	for _, t := range *slice {
		bs, err := t.MarshalBinary()
		if err != nil {
			return nil, err
		}
		bf.Write(bs)
	}

	return bf.Bytes(), nil
}

func (slice *OperationSlice) UnMarshalBinary(d []byte) error {
	remaining := d

	// TODO: loop to get all ops
	op := new(Operation)
	rem := op.UnMarshalBinary(remaining)

	(*slice) = append((*slice), *op)
	remaining = rem
	return nil
}
