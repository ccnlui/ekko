package raw

import (
	"math"
	"testing"
	"time"

	"ekko/raw/conditions"
	"ekko/raw/symbol"
)

func BenchmarkRawQuoteMarshalMsgPack(b *testing.B) {
	now := uint64(time.Now().UnixNano())
	quote := &Quote{
		Symbol:      symbol.NewLongSymbol("12345678901"),
		BidPrice:    math.MaxFloat64,
		AskPrice:    math.MaxFloat64,
		BidSize:     math.MaxUint32,
		AskSize:     math.MaxUint32,
		BidExchange: '!',
		AskExchange: '!',
		Conditions:  conditions.QuoteConditions{'A', 'B'},
		Nbbo:        false,
		Tape:        'C',
		Timestamp:   now,
		ReceivedAt:  now,
	}
	data := quote.MarshalRaw()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = data.MarshalMsgPack()
	}
}
