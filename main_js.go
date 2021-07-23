package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"syscall/js"

	"github.com/cfi2017/bl3-save-core/pkg/assets"
	"github.com/cfi2017/bl3-save-core/pkg/character"
	"github.com/cfi2017/bl3-save-core/pkg/item"
	"github.com/cfi2017/bl3-save-core/pkg/pb"
	"github.com/cfi2017/bl3-save-core/pkg/profile"
	shared2 "github.com/cfi2017/bl3-save-core/pkg/shared"
)

var c chan bool

func init() {
	c = make(chan bool)
}

func main() {
	js.Global().Set("shutdown", js.FuncOf(quit))
	js.Global().Set("decodeCharacter", js.FuncOf(decodeCharacter))
	js.Global().Set("decodeProfile", js.FuncOf(decodeProfile))
	js.Global().Set("encodeCharacter", js.FuncOf(encodeCharacter))
	js.Global().Set("encodeProfile", js.FuncOf(encodeProfile))
	js.Global().Set("deserialiseItem", js.FuncOf(deserialiseItem))
	js.Global().Set("deserialiseItemBase64", js.FuncOf(deserialiseItemBase64))
	js.Global().Set("serialiseItem", js.FuncOf(serialiseItem))
	js.Global().Set("serialiseItemBase64", js.FuncOf(serialiseItemBase64))
	js.Global().Set("getSeedFromSerial", js.FuncOf(getSeedFromSerial))
	js.Global().Set("setAssetDB", js.FuncOf(setAssetDB))
	<-c
}

func quit(_ js.Value, _ []js.Value) interface{} {
	c <- true
	return nil
}

type InMemoryAssetLoader struct {
	DB   assets.PartsDatabase
	Btik map[string]string
}

type ItemRequest struct {
	Items    []item.Item                         `json:"items"`
	Equipped []*pb.EquippedInventorySaveGameData `json:"equipped"`
	Active   []int32                             `json:"active"`
}

func (i *InMemoryAssetLoader) GetDB() assets.PartsDatabase {
	return i.DB
}

func (i *InMemoryAssetLoader) GetBtik() map[string]string {
	return i.Btik
}

func setAssetDB(_ js.Value, args []js.Value) interface{} {
	var db = make(assets.PartsDatabase)
	var btik = make(map[string]string)
	err := json.Unmarshal([]byte(args[0].String()), &db)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal([]byte(args[1].String()), &btik)
	if err != nil {
		panic(err)
	}

	assets.DefaultAssetLoader = &InMemoryAssetLoader{
		DB:   db,
		Btik: btik,
	}
	return nil
}

func deserialiseItemBase64(_ js.Value, args []js.Value) interface{} {
	bs, err := base64.StdEncoding.DecodeString(args[0].String())
	if err != nil {
		fmt.Println(err)
		return ""
	}
	i, err := item.Deserialize(bs)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	bs, err = json.Marshal(i)
	return string(bs)
}

func serialiseItemBase64(_ js.Value, args []js.Value) interface{} {
	data := item.Item{}
	err := json.Unmarshal([]byte(args[0].String()), &data)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	seed, err := item.GetSeedFromSerial(data.Wrapper.ItemSerialNumber)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	serial, err := item.Serialize(data, seed)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return base64.StdEncoding.EncodeToString(serial)
}

func deserialiseItem(_ js.Value, args []js.Value) interface{} {
	bs := make([]byte, args[0].Length())
	js.CopyBytesToGo(bs, args[0])
	i, err := item.Deserialize(bs)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	bs, err = json.Marshal(i)
	return string(bs)
}

func serialiseItem(_ js.Value, args []js.Value) interface{} {
	data := item.Item{}
	err := json.Unmarshal([]byte(args[0].String()), &data)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	seed, err := item.GetSeedFromSerial(data.Wrapper.ItemSerialNumber)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	serial, err := item.Serialize(data, seed)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	dst := args[1].Invoke(len(serial))
	js.CopyBytesToJS(dst, serial)
	return dst
}

func getSeedFromSerial(_ js.Value, args []js.Value) interface{} {
	bs := make([]byte, args[0].Length())
	js.CopyBytesToGo(bs, args[0])
	num, err := item.GetSeedFromSerial(bs)
	if err != nil {
		panic(err)
	}
	return num
}

func decodeCharacter(_ js.Value, args []js.Value) interface{} {
	bs := make([]byte, args[0].Length())
	js.CopyBytesToGo(bs, args[0])
	r := bytes.NewReader(bs)
	s, c := character.Deserialize(r, args[1].String())
	// workaround for invalid json parsing values
	for _, d := range c.GbxZoneMapFodSaveGameData.LevelData {
		if d.DiscoveryPercentage > math.MaxFloat32 {
			d.DiscoveryPercentage = -1
		}
	}

	items := pbArrayToItems(c.InventoryItems)

	bs, err := json.Marshal(struct {
		Save      shared2.SavFile `json:"save"`
		Character pb.Character    `json:"character"`
		Items     ItemRequest     `json:"items"`
	}{s, c, ItemRequest{
		Items:    items,
		Equipped: c.EquippedInventoryList,
		Active:   c.ActiveWeaponList,
	}})
	if err != nil {
		panic(err)
	}
	return string(bs)
}

func decodeProfile(_ js.Value, args []js.Value) interface{} {
	bs := make([]byte, args[0].Length())
	js.CopyBytesToGo(bs, args[0])
	r := bytes.NewReader(bs)
	s, p := profile.Deserialize(r, args[1].String())

	items := bankToItems(p.BankInventoryList)

	bs, err := json.Marshal(struct {
		Save    shared2.SavFile `json:"save"`
		Profile pb.Profile      `json:"profile"`
		Items   []item.Item     `json:"items"`
	}{s, p, items})
	if err != nil {
		panic(err)
	}
	return string(bs)

}

func bankToItems(list [][]byte) []item.Item {
	items := make([]item.Item, 0)
	for _, data := range list {
		d := make([]byte, len(data))
		copy(d, data)
		i, err := item.Deserialize(d)
		if err != nil {
			log.Println(err)
			log.Println(base64.StdEncoding.EncodeToString(data))
			// c.AbortWithStatus(500)
			// return
		}
		i.Wrapper = &pb.OakInventoryItemSaveGameData{
			ItemSerialNumber: data,
		}
		items = append(items, i)
	}
	return items
}

func encodeCharacter(_ js.Value, args []js.Value) interface{} {
	var data struct {
		Save      shared2.SavFile `json:"save"`
		Character pb.Character    `json:"character"`
		Items     ItemRequest     `json:"items"`
	}
	err := json.Unmarshal([]byte(args[0].String()), &data)
	if err != nil {
		return nil
	}

	data.Character.InventoryItems, err = itemsToPBArray(data.Items.Items)
	if err != nil {
		panic(err)
	}

	data.Character.EquippedInventoryList = data.Items.Equipped
	data.Character.ActiveWeaponList = data.Items.Active

	buf := new(bytes.Buffer)
	character.Serialize(buf, data.Save, data.Character, args[2].String())
	bs, _ := ioutil.ReadAll(buf)
	dst := args[1].Invoke(len(bs))
	js.CopyBytesToJS(dst, bs)
	return dst
}

func encodeProfile(_ js.Value, args []js.Value) interface{} {
	var data struct {
		Save    shared2.SavFile `json:"save"`
		Profile pb.Profile      `json:"profile"`
		Items   []item.Item     `json:"items"`
	}
	err := json.Unmarshal([]byte(args[0].String()), &data)
	if err != nil {
		return nil
	}
	buf := new(bytes.Buffer)

	pba, err := itemsToPBArray(data.Items)
	data.Profile.BankInventoryList = make([][]byte, len(pba))
	for i := range pba {
		data.Profile.BankInventoryList[i] = pba[i].ItemSerialNumber
	}
	profile.Serialize(buf, data.Save, data.Profile, args[2].String())
	bs, _ := ioutil.ReadAll(buf)
	dst := args[1].Invoke(len(bs))
	js.CopyBytesToJS(dst, bs)
	return dst
}

func pbArrayToItems(inventoryItems []*pb.OakInventoryItemSaveGameData) []item.Item {
	items := make([]item.Item, 0)
	for _, data := range inventoryItems {
		d := make([]byte, len(data.ItemSerialNumber))
		copy(d, data.ItemSerialNumber)
		i, err := item.Deserialize(d)
		if err != nil {
			log.Println(err)
			log.Println(base64.StdEncoding.EncodeToString(data.ItemSerialNumber))
		}
		i.Wrapper = data
		items = append(items, i)
	}
	return items
}

func itemsToPBArray(items []item.Item) ([]*pb.OakInventoryItemSaveGameData, error) {
	result := make([]*pb.OakInventoryItemSaveGameData, len(items))
	for index, i := range items {
		result[index] = i.Wrapper
		seed, err := item.GetSeedFromSerial(i.Wrapper.ItemSerialNumber)
		if err != nil {
			// set seed to be 0
			seed = 0
		}
		if i.Balance == "" {
			// sanity check, if the balance is empty, just write the original item back
			continue
		}
		result[index].ItemSerialNumber, err = item.Serialize(i, seed)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}
