package character

import (
	"io"

	"github.com/cfi2017/bl3-save-core/pkg/pb"
	"github.com/cfi2017/bl3-save-core/pkg/shared"
	"google.golang.org/protobuf/proto"
)

var (
	PCMagic = shared.Magic{
		Prefix: []byte{
			0x71, 0x34, 0x36, 0xB3, 0x56, 0x63, 0x25, 0x5F,
			0xEA, 0xE2, 0x83, 0x73, 0xF4, 0x98, 0xB8, 0x18,
			0x2E, 0xE5, 0x42, 0x2E, 0x50, 0xA2, 0x0F, 0x49,
			0x87, 0x24, 0xE6, 0x65, 0x9A, 0xF0, 0x7C, 0xD7,
		},

		Xor: []byte{
			0x7C, 0x07, 0x69, 0x83, 0x31, 0x7E, 0x0C, 0x82,
			0x5F, 0x2E, 0x36, 0x7F, 0x76, 0xB4, 0xA2, 0x71,
			0x38, 0x2B, 0x6E, 0x87, 0x39, 0x05, 0x02, 0xC6,
			0xCD, 0xD8, 0xB1, 0xCC, 0xA1, 0x33, 0xF9, 0xB6,
		},
	}
)

func Decrypt(reader io.Reader, magic shared.Magic) (shared.SavFile, []byte) {
	s, data := shared.DeserializeHeader(reader)
	return s, shared.Decrypt(data, magic.Prefix, magic.Xor)
}

func Deserialize(reader io.Reader, magic shared.Magic) (shared.SavFile, pb.Character) {
	// deserialise header, decrypt data
	s, data := Decrypt(reader, magic)
	p := pb.Character{}
	if err := proto.Unmarshal(data, &p); err != nil {
		panic("couldn't unmarshal protobuf data")
	}

	return s, p
}

func Serialize(writer io.Writer, s shared.SavFile, p pb.Character, magic shared.Magic) {
	bs, err := proto.Marshal(&p)
	if err != nil {
		panic(err)
	}
	bs = shared.Encrypt(bs, magic.Prefix, magic.Xor)
	shared.SerializeHeader(writer, s, bs)

}
