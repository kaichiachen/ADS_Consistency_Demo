package consistency

import "bytes"

type Message struct {
	Identifier byte
	Data       []byte
	Reply      chan Message
}

func NewMessage(id byte) *Message {
	return &Message{Identifier: id}
}

func (m *Message) MarshalBinary() ([]byte, error) {
	bs := &bytes.Buffer{}

	bs.WriteByte(m.Identifier)
	return bs.Bytes(), nil
}

func (m *Message) UnMarshalBinary(d []byte) error {
	bs := bytes.NewBuffer(d)

	m.Identifier = bs.Next(1)[0]

	return nil
}
