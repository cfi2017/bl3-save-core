package item

import (
	"github.com/cfi2017/bl3-save-core/pkg/assets"
)

var (
	debug bool
)

func BogoEncrypt(seed int32, data []byte) []byte {
	if seed == 0 {
		return data
	}

	steps := int(seed&0x1F) % len(data)
	data = append(data[steps:], data[:steps]...)
	return xor(seed, data)
}

func BogoDecrypt(seed int32, data []byte) []byte {
	if seed == 0 {
		return data
	}

	data = xor(seed, data)
	steps := int(seed&0x1F) % len(data)
	return append(data[len(data)-steps:], data[:len(data)-steps]...)
}

func xor(seed int32, data []byte) []byte {
	x := uint64(seed>>5) & 0xFFFFFFFF
	// target 4248340707
	for i := range data {
		x = (x * 0x10A860C1) % 0xFFFFFFFB
		data[i] = byte((uint64(data[i]) ^ x) & 0xFF)
	}
	return data
}

func GetBits(k string, v uint64) int {
	return assets.GetDB().GetData(k).GetBits(v)
}

func GetIndexFor(k string, v string) int {
	for i, asset := range assets.GetDB().GetData(k).Assets {
		if asset == v {
			return i
		}
	}
	panic("no asset found while serializing")
}

func GetPart(key string, index uint64) string {
	data := assets.GetDB().GetData(key)
	if int(index) >= len(data.Assets) {
		return ""
	}
	return data.GetPart(index)
}

func ReadNBits(r *Reader, n int) uint64 {
	i, err := r.ReadInt(n)
	if err != nil {
		panic(err)
	}
	return i
}
