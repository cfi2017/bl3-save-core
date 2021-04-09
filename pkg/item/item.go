package item

import (
	"encoding/binary"
	"errors"
	"fmt"
	"hash/crc32"
	"log"
	"strings"

	"github.com/cfi2017/bl3-save-core/internal/item"
	"github.com/cfi2017/bl3-save-core/pkg/assets"
	"github.com/cfi2017/bl3-save-core/pkg/pb"
)

type Item struct {
	Level             int                              `json:"level"`
	Balance           string                           `json:"balance"`
	Manufacturer      string                           `json:"manufacturer"`
	InvData           string                           `json:"inv_data"`
	Parts             []string                         `json:"parts"`
	Generics          []string                         `json:"generics"`
	Overflow          string                           `json:"overflow"`
	Version           uint64                           `json:"version"`
	Wrapper           *pb.OakInventoryItemSaveGameData `json:"wrapper"`
	SkipIntrospection bool                             `json:"skipIntrospection"`
	raw               []byte                           `json:"-"`
	SerialVersion     uint8                            `json:"serialVersion"`
}

func DecryptSerial(data []byte) ([]byte, error) {
	if len(data) < 5 {
		return nil, errors.New("invalid serial length")
	}
	if 0x03 > data[0] || data[0] > 0x04 {
		return nil, errors.New("invalid serial version")
	}
	seed := int32(binary.BigEndian.Uint32(data[1:])) // next four bytes of serial are bogo seed
	decrypted := item.BogoDecrypt(seed, data[5:])
	crc := binary.BigEndian.Uint16(decrypted)                          // first two bytes of decrypted data are crc checksum
	combined := append(append(data[:5], 0xFF, 0xFF), decrypted[2:]...) // combined data with checksum replaced with 0xFF to compute checksum
	computedChecksum := crc32.ChecksumIEEE(combined)
	check := uint16(((computedChecksum) >> 16) ^ ((computedChecksum & 0xFFFF) >> 0))

	if crc != check {
		return nil, errors.New("checksum failure in packed data")
	}

	return decrypted[2:], nil
}

func EncryptSerial(data []byte, seed int32, version uint8) ([]byte, error) {
	prefix := []byte{version}

	// seed to bytes
	seedBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(seedBytes, uint32(seed))

	// prefix + seed + 0xFFFF (checksum blank) + data
	prefix = append(prefix, seedBytes...)
	prefix = append(prefix, 0xFF, 0xFF)
	data = append(prefix, data...)

	// calculate checksum
	crc := crc32.ChecksumIEEE(data)
	checksum := ((crc >> 16) ^ crc) & 0xFFFF
	sumBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(sumBytes, uint16(checksum))

	// set checksum bytes
	data[5], data[6] = sumBytes[0], sumBytes[1] // set crc

	// return prefix + seed + encrypted data
	return append(append([]byte{version}, seedBytes...), item.BogoEncrypt(seed, data[5:])...), nil

}

/*
GetSeedFromSerial returns the seed bytes for the given serial.
Returns an error if the serial does not have the required length.
*/
func GetSeedFromSerial(data []byte) (int32, error) {
	if len(data) < 5 {
		return 0, errors.New("invalid serial length")
	}
	return int32(binary.BigEndian.Uint32(data[1:])), nil
}

/*
Deserialize decrypts and deserializes a serial number into an item.
This requires a valid database to be set.
*/
func Deserialize(data []byte) (i Item, err error) {
	data = makeCopy(data)
	i.raw = make([]byte, len(data))
	copy(i.raw, data)
	i.SerialVersion = data[0]
	data, err = DecryptSerial(data)
	if err != nil {
		return
	}

	r := item.NewReader(data)
	num := item.ReadNBits(r, 8)
	if num != 128 {
		err = errors.New(fmt.Sprintf("value should be %d, is %d", 128, num))
		return
	}

	i.Version = item.ReadNBits(r, 7)

	balanceBits := item.GetBits("InventoryBalanceData", i.Version)
	invDataBits := item.GetBits("InventoryData", i.Version)
	manBits := item.GetBits("ManufacturerData", i.Version)

	i.Balance = item.GetPart("InventoryBalanceData", item.ReadNBits(r,
		balanceBits)-1)
	i.InvData = item.GetPart("InventoryData", item.ReadNBits(r,
		invDataBits)-1)
	i.Manufacturer = item.GetPart("ManufacturerData", item.ReadNBits(r,
		manBits)-1)
	i.Level = int(item.ReadNBits(r, 7))

	if k, e := assets.GetBtik()[strings.ToLower(i.Balance)]; e {
		bits := item.GetBits(k, i.Version)
		partCount := int(item.ReadNBits(r, 6))
		i.Parts = make([]string, partCount)
		for index := 0; index < partCount; index++ {
			i.Parts[index] = item.GetPart(k, item.ReadNBits(r, bits)-1)
		}
		genericCount := item.ReadNBits(r, 4)
		i.Generics = make([]string, genericCount)
		bits = item.GetBits("InventoryGenericPartData", i.Version)
		for index := 0; index < int(genericCount); index++ {
			// looks like the bits are the same
			// for all the parts and generics
			i.Generics[index] = item.GetPart("InventoryGenericPartData", item.ReadNBits(r, bits)-1)
		}
		i.Overflow = r.Overflow()

	} else {
		err = errors.New(fmt.Sprintf("unknown category %s, skipping part introspection", i.Balance))
		i.SkipIntrospection = true
	}

	return
}

func makeCopy(data []byte) []byte {
	tmp := make([]byte, len(data))
	copy(tmp, data)
	return tmp
}

/*
Serialize serializes an item into a serial number with the given seed.
This requires a valid database to be set.
*/
func Serialize(i Item, seed int32) ([]byte, error) {
	// skip introspection if set, don't accidentally remove items
	if i.Wrapper != nil && i.Wrapper.ItemSerialNumber != nil && i.SkipIntrospection {
		return i.Wrapper.ItemSerialNumber, nil
	}
	w := item.NewWriter(i.Overflow)
	var err error

	// how many bits for each generic part?
	bits := item.GetBits("InventoryGenericPartData", i.Version)

	// write each generic, bottom to top
	for index := len(i.Generics) - 1; index >= 0; index-- {
		index, err := item.GetIndexFor("InventoryGenericPartData", i.Generics[index])
		if err != nil {
			return nil, err
		}
		err = w.WriteInt(uint64(index)+1, bits)
		if err != nil {
			log.Printf("tried to fit index %v into %v bits for %s", index, bits, i.Generics[index])
			return nil, err
		}
	}
	// write generic count
	err = w.WriteInt(uint64(len(i.Generics)), 4)
	if err != nil {
		return nil, err
	}
	if k, e := assets.GetBtik()[strings.ToLower(i.Balance)]; e {
		// how many bits per part?
		bits = item.GetBits(k, i.Version)
		// write each part, bottom to top
		for index := len(i.Parts) - 1; index >= 0; index-- {
			partIndex, err := item.GetIndexFor(k, i.Parts[index])
			if err != nil {
				return nil, err
			}
			err = w.WriteInt(uint64(partIndex)+1, bits)
			if err != nil {
				return nil, err
			}
		}
		// write part count
		err = w.WriteInt(uint64(len(i.Parts)), 6)
		if err != nil {
			return nil, err
		}
	}

	err = w.WriteInt(uint64(i.Level), 7)
	if err != nil {
		return nil, err
	}

	manIndex, err := item.GetIndexFor("ManufacturerData", i.Manufacturer)
	if err != nil {
		return nil, err
	}
	manBits := item.GetBits("ManufacturerData", i.Version)
	err = w.WriteInt(uint64(manIndex)+1, manBits)
	if err != nil {
		return nil, err
	}
	invIndex, err := item.GetIndexFor("InventoryData", i.InvData)
	if err != nil {
		return nil, err
	}
	invBits := item.GetBits("InventoryData", i.Version)
	err = w.WriteInt(uint64(invIndex)+1, invBits)
	if err != nil {
		return nil, err
	}
	balanceIndex, err := item.GetIndexFor("InventoryBalanceData", i.Balance)
	if err != nil {
		return nil, err
	}
	balanceBits := item.GetBits("InventoryBalanceData", i.Version)
	err = w.WriteInt(uint64(balanceIndex)+1, balanceBits)
	if err != nil {
		return nil, err
	}

	err = w.WriteInt(i.Version, 7)
	if err != nil {
		return nil, err
	}

	err = w.WriteInt(128, 8)
	if err != nil {
		return nil, err
	}

	return EncryptSerial(w.GetBytes(), seed, i.SerialVersion)

}
