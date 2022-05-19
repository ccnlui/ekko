package raw

import (
	"math"
	"testing"
	"time"

	"ekko/raw/conditions"
	"ekko/raw/symbol"
)

func BenchmarkRawTradeMarshalMsgPack(b *testing.B) {
	now := uint64(time.Now().UnixNano())
	trade := &Trade{
		Id:         math.MaxUint64,
		Timestamp:  now,
		Price:      math.MaxFloat64,
		Volume:     math.MaxUint32,
		Conditions: conditions.TradeConditions{'A', 'B', 'C', 'D'},
		Symbol:     symbol.NewLongSymbol("12345678901"),
		Exchange:   '!',
		Tape:       'A',
		ReceivedAt: now,
	}
	data := trade.MarshalRaw()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = data.MarshalMsgPack()
	}
}
