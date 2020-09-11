package profile

import (
	"io"

	"github.com/cfi2017/bl3-save-core/pkg/pb"
	shared2 "github.com/cfi2017/bl3-save-core/pkg/shared"
	"google.golang.org/protobuf/proto"
)

type Magic struct {
	Prefix []byte
	Xor    []byte
}

var (
	PCMagic = Magic{
		Prefix: []byte{
			0xD8, 0x04, 0xB9, 0x08, 0x5C, 0x4E, 0x2B, 0xC0,
			0x61, 0x9F, 0x7C, 0x8D, 0x5D, 0x34, 0x00, 0x56,
			0xE7, 0x7B, 0x4E, 0xC0, 0xA4, 0xD6, 0xA7, 0x01,
			0x14, 0x15, 0xA9, 0x93, 0x1F, 0x27, 0x2C, 0x8F,
		},
		Xor: []byte{
			0xE8, 0xDC, 0x3A, 0x66, 0xF7, 0xEF, 0x85, 0xE0,
			0xBD, 0x4A, 0xA9, 0x73, 0x57, 0x99, 0x30, 0x8C,
			0x94, 0x63, 0x59, 0xA8, 0xC9, 0xAE, 0xD9, 0x58,
			0x7D, 0x51, 0xB0, 0x1E, 0xBE, 0xD0, 0x77, 0x43,
		},
	}

	PS4Magic = Magic{
		Prefix: []byte{
			0xd1, 0x7b, 0xbf, 0x75, 0x4c, 0xc1, 0x80, 0x30,
			0x37, 0x92, 0xbd, 0xd0, 0x18, 0x3e, 0x4a, 0x5f,
			0x43, 0xa2, 0x46, 0xa0, 0xed, 0xdb, 0x2d, 0x9f,
			0x56, 0x5f, 0x8b, 0x3d, 0x6e, 0x73, 0xe6, 0xb8,
		},
		Xor: []byte{
			0xfb, 0xfd, 0xfd, 0x51, 0x3a, 0x5c, 0xdb, 0x20,
			0xbb, 0x5e, 0xc7, 0xaf, 0x66, 0x6f, 0xb6, 0x9a,
			0x9a, 0x52, 0x67, 0x0f, 0x19, 0x5d, 0xd3, 0x84,
			0x15, 0x19, 0xc9, 0x4a, 0x79, 0x67, 0xda, 0x6d,
		},
	}
)

func Decrypt(reader io.Reader, magic Magic) (shared2.SavFile, []byte) {
	s, data := shared2.DeserializeHeader(reader)
	return s, shared2.Decrypt(data, PCMagic.Prefix, PCMagic.Xor)
}

func Deserialize(reader io.Reader, magic Magic) (shared2.SavFile, pb.Profile) {
	// deserialise header, decrypt data
	s, data := Decrypt(reader, magic)

	p := pb.Profile{}
	if err := proto.Unmarshal(data, &p); err != nil {
		panic("couldn't unmarshal protobuf data")
	}

	return s, p
}

func Serialize(writer io.Writer, s shared2.SavFile, p pb.Profile) {
	bs, err := proto.Marshal(&p)
	if err != nil {
		panic(err)
	}
	bs = shared2.Encrypt(bs, PCMagic.Prefix, PCMagic.Xor)
	shared2.SerializeHeader(writer, s, bs)

}
