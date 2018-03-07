package common

func FitBytes(bs []byte, l int) []byte {
	if len(bs) < l {
		for len(bs) < l {
			bs = append([]byte{0}, bs...)
		}
	}
	return bs[:l]
}
