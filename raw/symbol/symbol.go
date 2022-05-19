package symbol

type ShortSymbol [5]byte
type LongSymbol [11]byte

func NewShortSymbol(s string) ShortSymbol {
	ss := ShortSymbol{' ', ' ', ' ', ' ', ' '}
	copy(ss[:], s)
	return ss
}

func (ss ShortSymbol) ToLongSymbol() LongSymbol {
	ls := LongSymbol{' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' '}
	copy(ls[:], ss[:])
	return ls
}

func NewLongSymbol(s string) LongSymbol {
	ls := LongSymbol{' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' '}
	copy(ls[:], s)
	return ls
}

func (ls LongSymbol) Len() int {
	for i, c := range ls {
		if c == ' ' {
			return i
		}
	}
	return 11
}

func (ls LongSymbol) Bytes() []byte {
	return ls[:ls.Len()]
}

func (ls LongSymbol) String() string {
	return string(ls.Bytes())
}
