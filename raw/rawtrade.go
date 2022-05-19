package raw

import (
	"encoding/binary"

	"ekko/raw/conditions"
	"ekko/raw/symbol"
)

type RawTrade [54]byte

func (rt *RawTrade) ReceivedAt() uint64 {
	return binary.BigEndian.Uint64(rt[46:54])
}

func (rt *RawTrade) MarshalMsgPack() []byte {
	symbol := *(*symbol.LongSymbol)(rt[1:])
	symbolLen := symbol.Len()

	conds := *(*conditions.TradeConditions)(rt[40:])
	condsLen := conds.Len()

	// 74 for timestamp64, +8 for timestamp96
	b := make([]byte, 0, 74+symbolLen+2*condsLen)
	b = append(b, 0x8a, 0xa1, 'T', 0xa1, rt[0])    // type
	b = append(b, 0xa1, 'i', 0xd3)                 // id
	b = append(b, rt[16:24]...)                    //
	b = append(b, 0xa1, 'S', 0xa0+byte(symbolLen)) // symbol
	b = append(b, symbol.Bytes()...)               //
	b = append(b, 0xa1, 'x', 0xa1, rt[44])         // exchange
	b = append(b, 0xa1, 'p', 0xcb)                 // price
	b = append(b, rt[32:40]...)                    //
	b = append(b, 0xa1, 's', 0xce)                 // size
	b = append(b, rt[12:16]...)                    //
	b = append(b, 0xa1, 'c', 0x90+byte(condsLen))  // conditions
	for _, c := range conds.Bytes() {
		b = append(b, 0xa1, c)
	}
	b = append(b, 0xa1, 'z', 0xa1, rt[45])     // tape
	b = append(b, 0xa1, 't')                   // timestamp
	b = appendMsgpackTimestamp64(b, rt[24:32]) //
	b = append(b, 0xa1, 'r')                   // receivedAt
	b = appendMsgpackTimestamp64(b, rt[46:54]) //
	return b
}
