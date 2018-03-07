package common

import (
	"bytes"
	"encoding/binary"
)

const (
	ITEM_NAME_LENGTH = 100
	ITEM_ID_LENGTH   = 10
)

type Item struct {
	Name   string
	Volume uint32
	ID     string
}

func (item *Item) MarshalBinary() ([]byte, error) {
	bs := &bytes.Buffer{}

	bs.Write(FitBytes([]byte(item.Name), ITEM_NAME_LENGTH))
	bs.Write([]byte(item.ID)[:ITEM_ID_LENGTH])
	binary.Write(bs, binary.LittleEndian, item.Volume)
	return bs.Bytes(), nil
}

func (item *Item) UnMarshalBinary(d []byte) error {
	bs := bytes.NewBuffer(d)
	item.Name = string(bs.Next(ITEM_NAME_LENGTH))
	item.ID = string(bs.Next(ITEM_ID_LENGTH))
	binary.Read(bytes.NewBuffer(bs.Next(4)), binary.LittleEndian, &item.Volume)
	return nil
}
