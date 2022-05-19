package raw

import (
	"encoding/binary"
	"errors"
	"math"

	"ekko/raw/conditions"
	"ekko/raw/symbol"
)

// Quote
// -------------|--------------|-------------|-----------
//  Field       | NASDAQ short | NASDAQ long | NYSE long
// -------------|--------------|-------------|-----------
//  timestamp   | uint64       | uint64      | uint64
//  symbol      | [5]byte      | [11]byte    | [11]byte
//  bidPrice    | uint16       | uint64      | uint64
//  bidSize     | uint16       | uint32      | uint32
//  bidExchange | byte         | byte        | byte
//  askPrice    | uint16       | uint64      | uint64
//  askSize     | uint16       | uint32      | uint32
//  askExchange | byte         | byte        | byte
//  condition   | byte         | byte        | byte
//  nbbo        | bool         | bool        | bool
// -------------|--------------|-------------|-----------

type Quote struct {
	Type        byte                       //  1
	Symbol      symbol.LongSymbol          // 12
	AskExchange byte                       // 13
	BidExchange byte                       // 14
	Conditions  conditions.QuoteConditions // 16
	Timestamp   uint64                     // 24
	AskPrice    float64                    // 32
	BidPrice    float64                    // 40
	AskSize     uint32                     // 44
	BidSize     uint32                     // 48
	Nbbo        bool                       // 49
	Tape        byte                       // 50 - not in the original Quote message
	ReceivedAt  uint64                     // 58 - not in the original Quote message
}

func (q *Quote) PutMarshaledRaw(b []byte) {
	nbbo := byte(0)
	if q.Nbbo {
		nbbo = 1
	}

	binary.BigEndian.PutUint64(b[50:], q.ReceivedAt) // go bounds check elimination
	b[0] = 'q'
	copy(b[1:], q.Symbol[:])
	b[12] = q.AskExchange
	b[13] = q.BidExchange
	b[14] = q.Conditions[0]
	b[15] = q.Conditions[1]
	binary.BigEndian.PutUint64(b[16:], q.Timestamp)
	binary.BigEndian.PutUint64(b[24:], math.Float64bits(q.AskPrice))
	binary.BigEndian.PutUint64(b[32:], math.Float64bits(q.BidPrice))
	binary.BigEndian.PutUint32(b[40:], q.AskSize)
	binary.BigEndian.PutUint32(b[44:], q.BidSize)
	b[48] = nbbo
	b[49] = q.Tape
}

func (q *Quote) MarshalRaw() RawQuote {
	b := RawQuote{}
	q.PutMarshaledRaw(b[:])
	return b
}

func (q *Quote) UnmarshalRaw(rq RawQuote) error {
	if rq[0] != 'q' {
		return errors.New("invalid quote")
	}

	q.ReceivedAt = binary.BigEndian.Uint64(rq[50:]) // go bounds check elimination
	q.Symbol = *(*symbol.LongSymbol)(rq[1:])
	q.AskExchange = rq[12]
	q.BidExchange = rq[13]
	q.Conditions[0] = rq[14]
	q.Conditions[1] = rq[15]
	q.Timestamp = binary.BigEndian.Uint64(rq[16:])
	q.AskPrice = math.Float64frombits(binary.BigEndian.Uint64(rq[24:]))
	q.BidPrice = math.Float64frombits(binary.BigEndian.Uint64(rq[32:]))
	q.AskSize = binary.BigEndian.Uint32(rq[40:])
	q.BidSize = binary.BigEndian.Uint32(rq[44:])
	q.Nbbo = rq[48] == 1
	q.Tape = rq[49]

	return nil
}
