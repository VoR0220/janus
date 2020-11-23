package transformer

import (
	"encoding/binary"
	"hash"

	"github.com/qtumproject/janus/pkg/eth"

	"github.com/ethereum/go-ethereum/common"
	"golang.org/x/crypto/sha3"
)

// ETH implementation of the Bloom filter calculation because QTUM RPC currently does not support this

type bytesBacked interface {
	Bytes() []byte
}

type KeccakState interface {
	hash.Hash
	Read([]byte) (int, error)
}

const (
	// BloomByteLength represents the number of bytes used in a header log bloom.
	BloomByteLength = 256

	// BloomBitLength represents the number of bits used in a header log bloom.
	BloomBitLength = 8 * BloomByteLength
)

// Bloom represents a 2048 bit bloom filter.
type Bloom [BloomByteLength]byte

func getLogsBloom(logs []eth.Log) string {
	buf := make([]byte, 6)
	var logsBloom Bloom
	for _, log := range logs {
		logsBloom.add(common.HexToAddress(log.Address).Bytes(), buf)
		for _, b := range log.Topics {
			logsBloom.add(common.HexToHash(b).Bytes(), buf)
		}
	}
	return string(logsBloom[:])
}

func (b *Bloom) add(d []byte, buf []byte) {
	i1, v1, i2, v2, i3, v3 := bloomValues(d, buf)
	b[i1] |= v1
	b[i2] |= v2
	b[i3] |= v3
}

// bloomValues returns the bytes (index-value pairs) to set for the given data
func bloomValues(data []byte, hashbuf []byte) (uint, byte, uint, byte, uint, byte) {
	sha := sha3.NewLegacyKeccak256().(KeccakState)
	sha.Reset()
	sha.Write(data)
	sha.Read(hashbuf)
	// The actual bits to flip
	v1 := byte(1 << (hashbuf[1] & 0x7))
	v2 := byte(1 << (hashbuf[3] & 0x7))
	v3 := byte(1 << (hashbuf[5] & 0x7))
	// The indices for the bytes to OR in
	i1 := BloomByteLength - uint((binary.BigEndian.Uint16(hashbuf)&0x7ff)>>3) - 1
	i2 := BloomByteLength - uint((binary.BigEndian.Uint16(hashbuf[2:])&0x7ff)>>3) - 1
	i3 := BloomByteLength - uint((binary.BigEndian.Uint16(hashbuf[4:])&0x7ff)>>3) - 1

	return i1, v1, i2, v2, i3, v3
}

// Test checks if the given topic is present in the bloom filter
func (b Bloom) Test(topic []byte) bool {
	i1, v1, i2, v2, i3, v3 := bloomValues(topic, make([]byte, 6))
	return v1 == v1&b[i1] &&
		v2 == v2&b[i2] &&
		v3 == v3&b[i3]
}
