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

type StartupMsg struct {
	RedNum        uint32
}

var typeMap = map[int]int{
	OP_ADDITEM:  RED,
	OP_ADDCART:  BLUE,
	OP_REMOVE:   BLUE,
	OP_CLEAR:    BLUE,
	OP_CHECKOUT: RED,
}

func (msg *StartupMsg) UnMarshalInt(d []byte) {
       bs := bytes.NewBuffer(d)
       binary.Read(bytes.NewBuffer(bs.Next(4)), binary.LittleEndian, &msg.RedNum)
}

func (op *Operation) SetOPType() {
	op.Optype = byte(typeMap[int(op.Action)])
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
		bs.Write(op.Payload)
	case OP_CLEAR:
		bs.Write(op.Payload)
	case OP_CHECKOUT:
		bs.Write(op.Payload)
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
		op.Payload = bs.Next(int(op.PayloadLength))
	case OP_CLEAR:
		op.Payload = bs.Next(int(op.PayloadLength))
	case OP_CHECKOUT:
		op.Payload = bs.Next(int(op.PayloadLength))
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
		if RedNum < s.Timestamp {
			RedNum = s.Timestamp
		}
		if !s.generator() {
			log.Println("Generate Fail!")
			break
		}
	}
}

func (slice *OperationSlice) Count() (red int, blue int) {
	red = 0
	blue = 0
	for _, s := range *slice {
		if s.Optype == RED {
			red += 1
		} else {
			blue += 1
		}
	}
	return red, blue
}

func (op Operation) generator() OP_RESULT {
	OpResult := OPERATION_FAIL
	switch op.Action {
	case OP_ADDITEM:
		OpResult = op.shadow()
	case OP_ADDCART:
		OpResult = op.shadow()
	case OP_REMOVE:
		OpResult = op.shadow()
	case OP_CLEAR:
		OpResult = op.shadow()
	case OP_CHECKOUT:
		OpResult = op.shadow()
	}
	return OpResult
}

func (op Operation) shadow() OP_RESULT {
	OpResult := OPERATION_FAIL
	switch op.Action {
	case OP_ADDITEM:
		newItem := common.Item{}
		newItem.UnMarshalBinary(op.Payload)
		OpResult = AddNewItem(newItem)
	case OP_ADDCART:
		item := common.Item{}
		item.UnMarshalBinary(op.Payload)
		OpResult = AddItemToCartForClient(item.ID, item.Volume)
	case OP_REMOVE:
		item := common.Item{}
		item.UnMarshalBinary(op.Payload)
		OpResult = RemoveItemFromCartForClient(item.ID, item.Volume)
	case OP_CLEAR:
		OpResult = ClearCartForServer()
	case OP_CHECKOUT:
		OpResult = CheckoutForServer(op)
	}

	return OpResult
}
