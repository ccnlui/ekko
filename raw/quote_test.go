package raw

import (
	"log"
	"math"
	"testing"
	"time"

	"ekko/raw/conditions"

	"ekko/raw/symbol"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQuoteMarshalUnmarshal(t *testing.T) {
	now := time.Now()
	quote := &Quote{
		Symbol:      symbol.NewLongSymbol("AAPL"),
		AskExchange: 'A',
		BidExchange: 'B',
		Conditions:  conditions.QuoteConditions{'@'},
		Timestamp:   uint64(now.UnixNano()),
		AskPrice:    987654.321,
		BidPrice:    123456.789,
		AskSize:     654321,
		BidSize:     123456,
		Nbbo:        true,
		Tape:        'C',
		ReceivedAt:  uint64(now.Add(-1 * time.Second).UnixNano()),
	}

	b := quote.MarshalRaw()

	got := &Quote{}
	err := got.UnmarshalRaw(b)
	require.NoError(t, err)
	assert.EqualValues(t, got, quote)
}

func TestQuoteUnmarshal_error(t *testing.T) {
	err := (&Quote{}).UnmarshalRaw([58]byte{'z'})
	require.Error(t, err)
}

func BenchmarkQuoteMarshalRaw(b *testing.B) {
	quote := &Quote{
		BidPrice:   math.MaxFloat64,
		AskPrice:   math.MaxFloat64,
		BidSize:    math.MaxUint32,
		AskSize:    math.MaxUint32,
		Conditions: conditions.QuoteConditions{'@'},
		Symbol:     symbol.NewLongSymbol("12345678901"),
		Tape:       'C',
	}
	size := 0

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		now := uint64(time.Now().UnixNano())
		exchange := 'A' + byte(i%26)
		quote.Timestamp = now
		quote.BidExchange = exchange
		quote.AskExchange = exchange
		quote.Nbbo = (i % 2) == 0
		quote.ReceivedAt = now

		size += len(quote.MarshalRaw())
	}
	b.ReportMetric(float64(size)/float64(b.N), "B/obj")
}

func BenchmarkQuoteUnmarshalRaw(b *testing.B) {
	now := uint64(time.Now().UnixNano())
	quote := &Quote{
		Timestamp:   now,
		BidPrice:    math.MaxFloat64,
		AskPrice:    math.MaxFloat64,
		BidSize:     math.MaxUint32,
		AskSize:     math.MaxUint32,
		BidExchange: '!',
		AskExchange: '!',
		Conditions:  conditions.QuoteConditions{'@'},
		Nbbo:        false,
		Symbol:      symbol.NewLongSymbol("12345678901"),
		Tape:        'C',
		ReceivedAt:  now,
	}
	data := quote.MarshalRaw()
	result := &Quote{}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		data[13] = 'A' + byte(i)
		if err := result.UnmarshalRaw(data); err != nil || result.BidExchange != 'A'+byte(i) {
			log.Fatal("quote unmarshal failed")
		}
	}
}
