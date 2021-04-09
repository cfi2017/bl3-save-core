package item

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log"
	"reflect"
	"testing"
)

var checks = []string{
	// v3 serials
	"A6cRHH+sfCuWGEZz2Lc5FWDbSfcQLmbaOV6SzgYP",
	"AwAAAADFtIC3/mrBkEsaj5NM0xGVIBFDCAAAAAAAMAYA",
	"AwAAAABLZ4A3RhkBkWMalJ8AEtSYWC1gJWYIAQAAAAAAyhgA",
	"AwAAAADuCYA3RhkBkWMalJ8AEtSYWC1gJmYIAQAAAAAAyhgA",
	"AwAAAACiM4A3rVMBk2tkjwhEwkRYO5cMpwkIAQAAAAAAyIAA",
	"AwAAAAA+uIC3syBDllvs4u4gGyP7LLHEEkssMQAA",
	"AwAAAADtqIC31oBBkWMEBcKAJnqQTAdOLWIIAQAAAAAAyIAA",
	"AwAAAACGEoC36JCAkTsKGoSgBASiIgsA",
	"AwAAAAAL94C3t9hAkysShLxMKmMLAA==",
	"AwAAAAByhIC3A/pBkWMGBYLB+IDMbhnFMWYIAQAAAAAAzoAA",
	"AwAAAACr6IC37xABkWsIqFPqeE0YjJJYUxUxhAAAAAAAAGQMAA==",
	"AwAAAAB1WoC3t9hAkysShLxMKkMLAA==",
	"AwAAAAC754A3ElwAmCtYUlWPjAAAAA==",
	"AwAAAADBk4A3ElwAmCtYUlWxjgAAAA==",
	"AwAAAABM+oC33IBBkWMEA0LBZlmkGKfELb4IAQAAAAAAzoAA",
	"AwAAAABDBIC3syBDllvs4u4gDG3MtcQVE0tsRQAA",
	"AwAAAADT2YC3syBDllvs4u6gz2zcdcUVS0ysRAAA",
	"AwAAAAB0xYA3pNNBkXMIKNMJSiplJhFOghEFhAAAAAAAAGUMAA==",
	"AwAAAADh6oA3p+vCkHsiOIpkJRQgNB8QyRCcxAwhAAAAAABAGQMA",
	"AwAAAACr94A3wNBBkWMIJxRMBMrEIJyELGIIAQAAAAAAyoAA",
	"AwAAAAA2EYC3pGNBk2MaDghSkel4SJFSMnQIAQAAAAAAyhgA",
	"AwAAAABvRoA38QgBk0sap9jjQvFwZDFDCAAAAAAAQAYA",
	"AwAAAABRaYC3DUGBkGMGuXk40BtJSotjghAAAAAAAGAMAA==",
	"AwAAAAAHiIC3NmvBkEsaD4dOwwlNchFDCAAAAAAAUAYA",
	"AwAAAAC2soA31ECBkFOGteE+ViSvNkwIAQAAAAAA0oAA",
	"AwAAAAB0hoA3a1IBk3MeMkhEkisIJhZQhqOLGUIAAAAAAIAzIAA=",
	"AwAAAACxi4C3ZINBkXMEA0KBl9EoUowTAtQDhAAAAAAAAGNAAA==",
	"AwAAAABMd4C3y+qAEmCaB3LONQrZ6stiihAAAAAAAGAJAA==",
	"AwAAAABnloC3z5lBk1saN8zHFMEJCKlMzBACAAAAAACQAQA=",
	"AwAAAACIAIC3t9hAkysShLxMKgMLAA==",
	"AwAAAADjeIA3VJMAkSsQUhYFGIMLAA==",
	"AwAAAADEyIA37wgBk1sap5fBcYmAShxYzBACAAAAAACkAQA=",
	// v4 serials
	"BKZ2z3UD2T0rpSE7deWBs4XmRkf8qUZ6rLgIC0adv9qk0yYeYeY=",
	"BOk5kBr6GvVQjoRaByNoBBeOF2BJgiGGd/9ui90SQNIEDVsqleI=",
	"BPlxBUiTWKvpoM5qL5HMO+XOeeNFAxu2P4aYHLgMzJLeorpGeX3wN32fZ8D6eoRE4m/KhJLa4adfWM2xolApJnZUzT1SK3UCpJwj1L1YtmR2Rh9kvM2c8XtnxECSYga0",
	"AwAAAACgzIA+8zyBEmwa54FUxlCkbcXYwEEIAAAAAABgxgAA",
	"AwAAAADP2YC+8ZQAETwQGoRKBQSCogsAAAA=",
	"AwAAAACgzIA+8zyBEmwa54FUxlCkbcXYwEEIAAAAAABgxgAA",
	"AwAAAADvcYA+AS1DFly2svqfzwxE9uqss74aawAAAA==",
	"AwAAAABYdYA+JfUBE4waDwckoWAsMszOyUSCVZoRiRAAAAAAAAANCAA=",
	"AwAAAAC6zIC+91yAEkRY4DpGLJiPJ1sAAAAA",
	"AwAAAADlL4C+8oxCEWQajyZjkfKQXCiKM3IQAgAAAAAAmAEBAA==",
	"AwAAAABSL4A+KWVBE3waD4jjSlAYLB2QiILCxhVCAAAAAAAAMwYAAA==",
	"AwAAAAByXYA+BGVAFlzQ0rqKyvgrdWuut+6aawAAAA==",
	"AwAAAACIEIC+/dzAEzQmgvhLM2BCCwAAAA==",
	"AwAAAABCpIC+9fQBE1waSSRjo/SEWphlnCEEAAAAAABIAwIA",
	"AwAAAAAyO4C+9jQBEVwaDogPYwp9NpFg8CAEAAAAAAAgYwAA",
	"AwAAAABu6IA+KiUBEnxojZoD8jASFpvQKfNBxhlCAAAAAACAYxAAAA==",
	"AwAAAADP2YC+8ZQAETwQGoRKBQSCogsAAAA=",
	"AwAAAAB44IC+BDVDFlzq8nqg/rzk9euvvw4bbAAAAA==",
	"AwAAAAAwO4A+wlwAGBRY0AIAAAA=",
	"AwAAAACMYIA+8txAEzwShPRJX01wggIAAAA=",
	"AwAAAACQJYA+/DQBEWwaDsiDsMCQRMANRRgDCAEAAAAAAMwYAAA=",
	"AwAAAABmO4A+9pTAEDwGGoRMT1F7ggYAAAA=",
	"AwAAAACEPoA+9pTAEDwGGoRMBQR6QgsAAAA=",
	"AwAAAACWzIA+ACVDFlzskvog9jbblOqqq656agAAAA==",
	"AwAAAACso4A+JN3AEzQmgvw/M2BiCwAAAA==",
	"AwAAAABp8IA+91wAGCxY3LKyJQgAAAA=",
	"AwAAAADC3YC+/JSBEHwckMijKpiYQBuP6mWlxBNCAAAAAACANCAAAA==",
}

func TestDecryptSerial(t *testing.T) {
	for _, check := range checks {
		bs, err := base64.StdEncoding.DecodeString(check)
		if err != nil {
			t.Fatal(err)
		}
		item, err := DecryptSerial(bs)
		if err != nil {
			t.Fatal(err)
		}
		log.Println(item)
	}
}

func TestDeserialize(t *testing.T) {
	for _, check := range checks {
		bs, err := base64.StdEncoding.DecodeString(check)
		if err != nil {
			t.Fatal(err)
		}
		item, err := Deserialize(bs)
		if err != nil {
			t.Fatal(err)
		}
		log.Println(item)
	}
}

func TestSerialize(t *testing.T) {
	for ci, check := range checks {
		var result = check
		var history = make([]string, 10)
		var item Item
		var last string
		for i := 0; i < 10; i++ {
			last = result
			bs, err := base64.StdEncoding.DecodeString(result)
			if err != nil {
				t.Fatal(err)
			}
			seed, err := GetSeedFromSerial(bs)
			if err != nil {
				t.Fatal(err)
			}
			item, err = Deserialize(bs)
			if err != nil {
				t.Fatal(err)
			}
			bs2, err := Serialize(item, seed)
			if err != nil {
				t.Fatal(err)
			}
			result = base64.StdEncoding.EncodeToString(bs2)
			history[i] = result
			i2, err := Deserialize(bs2)
			if err != nil {
				fmt.Println(reflect.DeepEqual(bs, bs2))
				t.Fatalf("error in deserialise pass 2 (%d, %d): %v", ci, i, err)
			}
			if item.Level != i2.Level || item.Version != i2.Version {
				t.Fatal("component mismatch in re-serialized item")
			}
			if result != last {
				log.Printf("item mismatch in iteration %d: %s\n", i, result)
				if i > 0 {
					log.Println(err)
					log.Println(check)
					log.Println(result)
					bs1, _ := base64.StdEncoding.DecodeString(check)
					bs2, _ := base64.StdEncoding.DecodeString(result)
					log.Println(hex.EncodeToString(bs1))
					log.Println(hex.EncodeToString(bs2))
					dec1, _ := DecryptSerial(bs1)
					dec2, _ := DecryptSerial(bs2)
					bs1, _ = base64.StdEncoding.DecodeString(check)
					bs2, _ = base64.StdEncoding.DecodeString(result)
					log.Println(hex.EncodeToString(dec1))
					log.Println(hex.EncodeToString(dec2))
					i1, _ := Deserialize(bs1)
					i2, _ := Deserialize(bs2)
					log.Println(i1.Version)
					log.Println(i2.Version)
					t.Fatal("invalid serial")
				}
			}
		}
	}
}

func TestAddPart(t *testing.T) {
	code := "AwAAAADuCYA3RhkBkWMalJ8AEtSYWC1gJmYIAQAAAAAAyhgA"
	part := "/Game/Gear/Weapons/Pistols/Vladof/_Shared/_Design/Parts/Barrels/Barrel_01/Part_PS_VLA_Barrel_01_B.Part_PS_VLA_Barrel_01_B"
	bs, err := base64.StdEncoding.DecodeString(code)
	if err != nil {
		t.Fatal(err)
	}
	seed, err := GetSeedFromSerial(bs)
	if err != nil {
		panic(err)
	}
	item, err := Deserialize(bs)
	if err != nil {
		t.Fatal(err)
	}
	item.Parts = append(item.Parts, part)
	bs, err = Serialize(item, seed)
	if err != nil {
		t.Fatal(err)
	}
	log.Println(base64.StdEncoding.EncodeToString(bs))
	i2, err := Deserialize(bs)
	if err != nil {
		t.Fatal(err)
	}
	if len(i2.Parts) != 13 {
		t.Fatalf("invalid part length %v", len(i2.Parts))
	}

}

func TestAddAnointment(t *testing.T) {
	code := "AwAAAADuCYA3RhkBkWMalJ8AEtSYWC1gJmYIAQAAAAAAyhgA"
	anointment := "/Game/Gear/Weapons/_Shared/_Design/EndGameParts/Character/Operative/CloneSwapDamage/GPart_CloneSwap_WeaponDamage.GPart_CloneSwap_WeaponDamage"
	bs, err := base64.StdEncoding.DecodeString(code)
	if err != nil {
		t.Fatal(err)
	}
	seed, err := GetSeedFromSerial(bs)
	if err != nil {
		t.Fatal(err)
	}
	item, err := Deserialize(bs)
	if err != nil {
		t.Fatal(err)
	}
	item.Generics = append(item.Generics, anointment)
	bs, err = Serialize(item, seed)
	if err != nil {
		t.Fatal(err)
	}
	log.Println(base64.StdEncoding.EncodeToString(bs))
	i2, err := Deserialize(bs)
	if err != nil {
		t.Fatal(err)
	}
	if len(i2.Generics) != 2 {
		t.Fatal("invalid anointment length")
	}

}
