package consistency

import (
	"bytes"
	"common"
	"encoding/binary"
	"log"
	"math"
)

type Operation struct {
	Optype        byte
	Action        byte
	Timestamp     uint32
	Payload       []byte
	PayloadLength uint32
}

func (op *Operation) SetOPType() {
	switch op.Action {
	case OP_ADDITEM:
		op.Optype = RED
	case OP_ADDCART:
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
	op.PayloadLength = uint32(len(op.Payload))
	binary.Write(bs, binary.LittleEndian, op.PayloadLength)
	switch op.Action {
	case OP_ADDITEM:
		bs.Write(op.Payload)
	case OP_ADDCART:
		bs.Write(op.Payload)
	case OP_REMOVE:

	case OP_CLEAR:

	case OP_SETTLE:

	}
	return bs.Bytes(), nil
}

func (op *Operation) UnMarshalBinary(d []byte) []byte {
	bs := bytes.NewBuffer(d)
	op.Optype = bs.Next(1)[0]
	op.Action = bs.Next(1)[0]
	binary.Read(bytes.NewBuffer(bs.Next(4)), binary.LittleEndian, &op.PayloadLength)
	switch op.Action {
	case OP_ADDITEM:
		op.Payload = bs.Next(int(op.PayloadLength))
	case OP_ADDCART:
		op.Payload = bs.Next(int(op.PayloadLength))
	case OP_REMOVE:

	case OP_CLEAR:

	case OP_SETTLE:

	}

	return bs.Next(math.MaxInt32)
}

type OperationSlice []Operation

func (slice OperationSlice) Len() int {
	return len(slice)
}

func (slice OperationSlice) AddOperation(op *Operation) OperationSlice {
	// for i, o := range slice {
	// 	if o.Timestamp > op.Timestamp {
	// 		slice = append(append(slice[:i], op), slice[i:]...)
	// 		return slice
	// 	}
	// }
	return append(slice, *op)
}

func (slice OperationSlice) ClearOperation() OperationSlice {
	return slice[:0]
}

func (slice *OperationSlice) MarshalBinary() ([]byte, error) {
	bf := &bytes.Buffer{}

	for _, s := range *slice {
		bs, err := s.MarshalBinary()
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
	for len(remaining) >= 2 {
		op := new(Operation)
		rem := op.UnMarshalBinary(remaining)

		(*slice) = append((*slice), *op)
		remaining = rem
	}
	return nil
}

func (slice *OperationSlice) HandleOperations() {
	for _, s := range *slice {
		if !s.generator() {
			log.Println("Generate Fail!")
			break
		}
	}
}

func (op Operation) generator() bool {
	log.Println(op.Action)
	switch op.Action {
	case OP_ADDITEM:
		op.shadow()
	case OP_ADDCART:
		op.shadow()
	case OP_REMOVE:
	case OP_CLEAR:
	case OP_SETTLE:
	}
	return true
}

func (op Operation) shadow() {
	switch op.Action {
	case OP_ADDITEM:
		newItem := common.Item{}
		newItem.UnMarshalBinary(op.Payload)
		AddNewItem(newItem)
	case OP_ADDCART:
		item := common.Item{}
		item.UnMarshalBinary(op.Payload)
		AddItemToCartForClient(item.ID, item.Volume)
	case OP_REMOVE:
	case OP_CLEAR:
	case OP_SETTLE:
	}
}
