package conditions

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTradeConditionsString(t *testing.T) {
	assert.Equal(t, "", TradeConditions{}.String())
	assert.Equal(t, "A", TradeConditions{'A'}.String())
	assert.Equal(t, "AB", TradeConditions{'A', 'B'}.String())
	assert.Equal(t, "ABC", TradeConditions{'A', 'B', 'C'}.String())
	assert.Equal(t, "ABCD", TradeConditions{'A', 'B', 'C', 'D'}.String())

	assert.Equal(t, "", TradeConditions{0, 'B', 'C', 'D'}.String())
	assert.Equal(t, "A", TradeConditions{'A', 0, 'C', 'D'}.String())
	assert.Equal(t, "AB", TradeConditions{'A', 'B', 0, 'D'}.String())
	assert.Equal(t, "ABC", TradeConditions{'A', 'B', 'C', 0}.String())
}

func TestTradeConditionsRemoveSpaces(t *testing.T) {
	assert.Equal(t, "@", TradeConditions{' ', '@'}.RemoveSpaces().String())
	assert.Equal(t, "@", TradeConditions{'@', ' '}.RemoveSpaces().String())
	assert.Equal(t, "A", TradeConditions{'A', ' ', ' ', ' '}.RemoveSpaces().String())
	assert.Equal(t, "A", TradeConditions{' ', 'A', ' ', ' '}.RemoveSpaces().String())
	assert.Equal(t, "A", TradeConditions{' ', ' ', 'A', ' '}.RemoveSpaces().String())
	assert.Equal(t, "A", TradeConditions{' ', ' ', ' ', 'A'}.RemoveSpaces().String())
}

func TestTradeConditionsRemoveNonFirstSpaces(t *testing.T) {
	assert.Equal(t, " @", TradeConditions{' ', '@'}.RemoveNonFirstSpaces().String())
	assert.Equal(t, "@", TradeConditions{'@', ' '}.RemoveNonFirstSpaces().String())
	assert.Equal(t, "A", TradeConditions{'A', ' ', ' ', ' '}.RemoveNonFirstSpaces().String())
	assert.Equal(t, " A", TradeConditions{' ', 'A', ' ', ' '}.RemoveNonFirstSpaces().String())
	assert.Equal(t, " A", TradeConditions{' ', ' ', 'A', ' '}.RemoveNonFirstSpaces().String())
	assert.Equal(t, " A", TradeConditions{' ', ' ', ' ', 'A'}.RemoveNonFirstSpaces().String())
}

func TestQuoteConditionsString(t *testing.T) {
	assert.Equal(t, "", QuoteConditions{}.String())
	assert.Equal(t, "A", QuoteConditions{'A'}.String())
	assert.Equal(t, "AB", QuoteConditions{'A', 'B'}.String())

	assert.Equal(t, "", QuoteConditions{0, 'B'}.String())
	assert.Equal(t, "A", QuoteConditions{'A', 0}.String())
}
