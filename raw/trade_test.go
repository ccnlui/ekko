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

func TestTradeMarshalUnmarshal(t *testing.T) {
	now := time.Now()
	trade := &Trade{
		Id:         1000,
		Timestamp:  uint64(now.UnixNano()),
		Price:      987654.321,
		Volume:     123456,
		Conditions: conditions.TradeConditions{'@'},
		Symbol:     symbol.NewLongSymbol("AAPL"),
		Exchange:   'V',
		Tape:       'A',
		ReceivedAt: uint64(now.Add(-1 * time.Second).UnixNano()),
	}

	b := trade.MarshalRaw()

	got := &Trade{}
	err := got.UnmarshalRaw(b)
	require.NoError(t, err)
	assert.EqualValues(t, got, trade)
}

func TestTradeUnmarshal_error(t *testing.T) {
	err := (&Trade{}).UnmarshalRaw([54]byte{'z'})
	require.Error(t, err)
}

func BenchmarkTradeMarshalRaw(b *testing.B) {
	trade := &Trade{
		Price:      math.MaxFloat64,
		Volume:     math.MaxUint32,
		Conditions: conditions.TradeConditions{'A', 'B', 'C', 'D'},
		Symbol:     symbol.NewLongSymbol("12345678901"),
		Tape:       'A',
	}
	size := 0

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		now := uint64(time.Now().UnixNano())
		trade.Id = uint64(i)
		trade.Timestamp = now
		trade.Exchange = 'A' + byte(i%26)
		trade.ReceivedAt = now

		size += len(trade.MarshalRaw())
	}
	b.ReportMetric(float64(size)/float64(b.N), "B/obj")
}

func BenchmarkTradeUnmarshalRaw(b *testing.B) {
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
	result := &Trade{}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		data[44] = 'A' + byte(i)
		if err := result.UnmarshalRaw(data); err != nil || result.Exchange != 'A'+byte(i) {
			log.Fatal("trade unmarshal failed")
		}
	}
}
