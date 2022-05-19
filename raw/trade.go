package raw

import (
	"encoding/binary"
	"errors"
	"math"

	"ekko/raw/conditions"
	"ekko/raw/symbol"
)

// Trade
// -----------|--------------|-------------|------------|-----------
//  Field     | NASDAQ short | NASDAQ long | NYSE short | NYSE long
// -----------|--------------|-------------|------------|-----------
//  timestamp | uint64       | uint64      | uint64     | uint64
//  symbol    | [5]byte      | [11]byte    | [5]byte    | [11]byte
//  tradeId   | uint64       | uint64      | int64      | int64
//  price     | uint16       | uint64      | uint16     | uint64
//  volume    | uint16       | uint32      | uint16     | uint32
//  cond      | [4]byte      | [4]byte     | byte       | [4]byte
//  exchange  | byte         | byte        | byte       | byte
// -----------|--------------|-------------|------------|-----------

type Trade struct {
	Type       byte                       //  1
	Symbol     symbol.LongSymbol          // 12
	Volume     uint32                     // 16
	Id         uint64                     // 24
	Timestamp  uint64                     // 32
	Price      float64                    // 40
	Conditions conditions.TradeConditions // 44
	Exchange   byte                       // 45
	Tape       byte                       // 46 - not in the original Trade message
	ReceivedAt uint64                     // 54 - not in the original Trade message
}

func (t *Trade) PutMarshaledRaw(b []byte) {
	binary.BigEndian.PutUint64(b[46:], t.ReceivedAt) // go bounds check elimination
	b[0] = 't'
	copy(b[1:], t.Symbol[:])
	binary.BigEndian.PutUint32(b[12:], t.Volume)
	binary.BigEndian.PutUint64(b[16:], t.Id)
	binary.BigEndian.PutUint64(b[24:], t.Timestamp)
	binary.BigEndian.PutUint64(b[32:], math.Float64bits(t.Price))
	copy(b[40:], t.Conditions[:])
	b[44] = t.Exchange
	b[45] = t.Tape
}

func (t *Trade) MarshalRaw() RawTrade {
	b := RawTrade{}
	t.PutMarshaledRaw(b[:])
	return b
}

func (t *Trade) UnmarshalRaw(b RawTrade) error {
	if b[0] != 't' {
		return errors.New("invalid trade")
	}

	t.ReceivedAt = binary.BigEndian.Uint64(b[46:]) // go bounds check elimination
	t.Symbol = *(*symbol.LongSymbol)(b[1:])
	t.Volume = binary.BigEndian.Uint32(b[12:])
	t.Id = binary.BigEndian.Uint64(b[16:])
	t.Timestamp = binary.BigEndian.Uint64(b[24:])
	t.Price = math.Float64frombits(binary.BigEndian.Uint64(b[32:]))
	t.Conditions = *(*conditions.TradeConditions)(b[40:])
	t.Exchange = b[44]
	t.Tape = b[45]

	return nil
}
