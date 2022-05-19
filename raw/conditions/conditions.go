package conditions

type TradeConditions [4]byte
type QuoteConditions [2]byte

func (tc TradeConditions) Len() int {
	for i, c := range tc {
		if c == 0 {
			return i
		}
	}
	return 4
}

func (tc TradeConditions) Bytes() []byte {
	return tc[:tc.Len()]
}

func (tc TradeConditions) String() string {
	return string(tc.Bytes())
}

func (tc TradeConditions) RemoveSpaces() TradeConditions {
	conditions := TradeConditions{}
	for i, j := 0, 0; i < 4; i++ {
		if tc[i] != ' ' {
			conditions[j] = tc[i]
			j++
		}
	}
	return conditions
}

func (tc TradeConditions) RemoveNonFirstSpaces() TradeConditions {
	conditions := TradeConditions{}
	for i, j := 0, 0; i < 4; i++ {
		if tc[i] != ' ' || i == 0 {
			conditions[j] = tc[i]
			j++
		}
	}
	return conditions
}

func (qc QuoteConditions) Len() int {
	switch {
	case qc[0] == 0:
		return 0
	case qc[1] == 0:
		return 1
	}
	return 2
}

func (qc QuoteConditions) Bytes() []byte {
	return qc[:qc.Len()]
}

func (qc QuoteConditions) String() string {
	return string(qc.Bytes())
}
