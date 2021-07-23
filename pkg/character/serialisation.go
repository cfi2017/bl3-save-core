package character

import (
	"io"

	"github.com/cfi2017/bl3-save-core/pkg/pb"
	"github.com/cfi2017/bl3-save-core/pkg/shared"
	"google.golang.org/protobuf/proto"
)

var (
	platforms = map[string]shared.PlatformMagic{
		"pc": {
			Prefix: []byte{
				0x71, 0x34, 0x36, 0xB3, 0x56, 0x63, 0x25, 0x5F,
				0xEA, 0xE2, 0x83, 0x73, 0xF4, 0x98, 0xB8, 0x18,
				0x2E, 0xE5, 0x42, 0x2E, 0x50, 0xA2, 0x0F, 0x49,
				0x87, 0x24, 0xE6, 0x65, 0x9A, 0xF0, 0x7C, 0xD7,
			}, Xor: []byte{
				0x7C, 0x07, 0x69, 0x83, 0x31, 0x7E, 0x0C, 0x82,
				0x5F, 0x2E, 0x36, 0x7F, 0x76, 0xB4, 0xA2, 0x71,
				0x38, 0x2B, 0x6E, 0x87, 0x39, 0x05, 0x02, 0xC6,
				0xCD, 0xD8, 0xB1, 0xCC, 0xA1, 0x33, 0xF9, 0xB6,
			},
		},
		"ps4": {
			Prefix: []byte{
				0xd1, 0x7b, 0xbf, 0x75, 0x4c, 0xc1, 0x80, 0x30,
				0x37, 0x92, 0xbd, 0xd0, 0x18, 0x3e, 0x4a, 0x5f,
				0x43, 0xa2, 0x46, 0xa0, 0xed, 0xdb, 0x2d, 0x9f,
				0x56, 0x5f, 0x8b, 0x3d, 0x6e, 0x73, 0xe6, 0xb8,
			}, Xor: []byte{
				0xfb, 0xfd, 0xfd, 0x51, 0x3a, 0x5c, 0xdb, 0x20,
				0xbb, 0x5e, 0xc7, 0xaf, 0x66, 0x6f, 0xb6, 0x9a,
				0x9a, 0x52, 0x67, 0x0f, 0x19, 0x5d, 0xd3, 0x84,
				0x15, 0x19, 0xc9, 0x4a, 0x79, 0x67, 0xda, 0x6d,
			},
		},
	}
)

func Decrypt(reader io.Reader, platform string) (shared.SavFile, []byte) {
	s, data := shared.DeserializeHeader(reader)
	return s, shared.Decrypt(data, platforms[platform].Prefix, platforms[platform].Xor)
}

func Deserialize(reader io.Reader, platform string) (shared.SavFile, pb.Character, error) {
	// deserialise header, decrypt data
	s, data := Decrypt(reader, platform)
	p := pb.Character{}
	if err := proto.Unmarshal(data, &p); err != nil {
		return s, p, err
	}

	return s, p, nil
}

func Serialize(writer io.Writer, s shared.SavFile, p pb.Character, platform string) {
	bs, err := proto.Marshal(&p)
	if err != nil {
		panic(err)
	}
	bs = shared.Encrypt(bs, platforms[platform].Prefix, platforms[platform].Xor)
	shared.SerializeHeader(writer, s, bs)

}
