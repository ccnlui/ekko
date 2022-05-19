package symbol

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewShortSymbol(t *testing.T) {
	assert.Equal(t, ShortSymbol{' ', ' ', ' ', ' ', ' '}, NewShortSymbol(""))
	assert.Equal(t, ShortSymbol{'A', ' ', ' ', ' ', ' '}, NewShortSymbol("A"))
	assert.Equal(t, ShortSymbol{'A', 'B', ' ', ' ', ' '}, NewShortSymbol("AB"))
	assert.Equal(t, ShortSymbol{'A', 'B', 'C', ' ', ' '}, NewShortSymbol("ABC"))
	assert.Equal(t, ShortSymbol{'A', 'B', 'C', 'D', ' '}, NewShortSymbol("ABCD"))
	assert.Equal(t, ShortSymbol{'A', 'B', 'C', 'D', 'E'}, NewShortSymbol("ABCDE"))
	assert.Equal(t, ShortSymbol{'A', 'B', 'C', 'D', 'E'}, NewShortSymbol("ABCDEF"))
}

func TestNewLongSymbol(t *testing.T) {
	assert.Equal(t, LongSymbol{' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' '}, NewLongSymbol(""))
	assert.Equal(t, LongSymbol{'A', 'A', 'P', 'L', ' ', ' ', ' ', ' ', ' ', ' ', ' '}, NewLongSymbol("AAPL"))
	assert.Equal(t, LongSymbol{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0', ' '}, NewLongSymbol("1234567890"))
	assert.Equal(t, LongSymbol{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0', '1'}, NewLongSymbol("12345678901"))
}

func TestShortSymbolToLong(t *testing.T) {
	assert.Equal(t, NewLongSymbol(""), NewShortSymbol("").ToLongSymbol())
	assert.Equal(t, NewLongSymbol("A"), NewShortSymbol("A").ToLongSymbol())
	assert.Equal(t, NewLongSymbol("AB"), NewShortSymbol("AB").ToLongSymbol())
	assert.Equal(t, NewLongSymbol("ABC"), NewShortSymbol("ABC").ToLongSymbol())
	assert.Equal(t, NewLongSymbol("ABCD"), NewShortSymbol("ABCD").ToLongSymbol())
	assert.Equal(t, NewLongSymbol("ABCDE"), NewShortSymbol("ABCDE").ToLongSymbol())
}

func TestLongSymbolBytes(t *testing.T) {
	assert.Equal(t, "", NewLongSymbol("").String())
	assert.Equal(t, "AAPL", NewLongSymbol("AAPL").String())
	assert.Equal(t, "12345678901", NewLongSymbol("12345678901").String())
}

func BenchmarkShortTradeToLong(b *testing.B) {
	aapl := NewShortSymbol("AAPL")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = aapl.ToLongSymbol()
	}
}

func BenchmarkLongTradeString(b *testing.B) {
	aapl := NewLongSymbol("AAPL")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = aapl.String()
	}
}
