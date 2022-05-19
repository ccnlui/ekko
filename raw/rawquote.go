package raw

import (
	"encoding/binary"

	"ekko/raw/conditions"
	"ekko/raw/symbol"
)

type RawQuote [58]byte

func (rt *RawQuote) ReceivedAt() uint64 {
	return binary.BigEndian.Uint64(rt[50:58])
}

func (rq *RawQuote) MarshalMsgPack() []byte {
	symbol := *(*symbol.LongSymbol)(rq[1:])
	symbolLen := symbol.Len()

	conds := *(*conditions.QuoteConditions)(rq[14:])
	condsLen := conds.Len()

	// 91 for timestamp64, +8 timestamp96
	b := make([]byte, 0, 91+symbolLen+2*condsLen)
	b = append(b, 0x8c, 0xa1, 'T', 0xa1, rq[0])    // type
	b = append(b, 0xa1, 'S', 0xa0+byte(symbolLen)) // symbol
	b = append(b, symbol.Bytes()...)               //
	b = append(b, 0xa2, 'b', 'x', 0xa1, rq[13])    // bid exchange
	b = append(b, 0xa2, 'b', 'p', 0xcb)            // bid price
	b = append(b, rq[32:40]...)                    //
	b = append(b, 0xa2, 'b', 's', 0xce)            // bid size
	b = append(b, rq[44:48]...)                    //
	b = append(b, 0xa2, 'a', 'x', 0xa1, rq[12])    // ask exchange
	b = append(b, 0xa2, 'a', 'p', 0xcb)            // ask price
	b = append(b, rq[24:32]...)                    //
	b = append(b, 0xa2, 'a', 's', 0xce)            // ask size
	b = append(b, rq[40:44]...)                    //
	b = append(b, 0xa1, 'c', 0x90+byte(condsLen))  // conditions
	for _, c := range conds.Bytes() {
		b = append(b, 0xa1, c)
	}
	b = append(b, 0xa1, 'z', 0xa1, rq[49])     // tape
	b = append(b, 0xa1, 't')                   // timestamp
	b = appendMsgpackTimestamp64(b, rq[16:24]) //
	b = append(b, 0xa1, 'r')                   // receivedAt
	b = appendMsgpackTimestamp64(b, rq[50:58]) //

	return b
}

func appendMsgpackTimestamp64(dst, src []byte) []byte {
	timestamp := binary.BigEndian.Uint64(src)
	dst = append(dst, 0xd7, 0xff)
	dst = append(dst, 0, 0, 0, 0, 0, 0, 0, 0)
	binary.BigEndian.PutUint64(dst[len(dst)-8:], timestamp%1e9<<34+timestamp/1e9)
	return dst
}

// func appendMsgpackTimestamp96(dst, src []byte) []byte {
// 	timestamp := binary.BigEndian.Uint64(src)
// 	dst = append(dst, 0xc7, 0x0c, 0xff)
// 	i := len(dst)
// 	dst = append(dst, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0)
// 	binary.BigEndian.PutUint32(dst[i:], uint32(timestamp%1e9))
// 	binary.BigEndian.PutUint64(dst[i+4:], timestamp/1e9)
// 	return dst
// }
