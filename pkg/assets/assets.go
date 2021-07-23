package assets

import (
	"encoding/json"
	"io/ioutil"
	"sync"
)

var DefaultAssetLoader AssetsLoader

type AssetsLoader interface {
	GetDB() PartsDatabase
	GetBtik() map[string]string
}

func GetDB() PartsDatabase {
	return DefaultAssetLoader.GetDB()
}

func GetBtik() map[string]string {
	return DefaultAssetLoader.GetBtik()
}

type StaticFileAssetLoader struct {
	once sync.Once
	Pwd  string
	btik map[string]string
	db   PartsDatabase
}

func (s *StaticFileAssetLoader) GetDB() PartsDatabase {
	err := s.Load()
	if err != nil {
		panic(err)
	}
	return s.db
}

func (s *StaticFileAssetLoader) Load() error {
	var err error
	s.once.Do(func() {
		s.btik, err = LoadPartMap("balance_to_inv_key.json")
		if err != nil {
			return
		}
		s.db, err = LoadPartsDatabase("inventory_raw.json")
	})
	return err
}

func (s *StaticFileAssetLoader) GetBtik() map[string]string {
	err := s.Load()
	if err != nil {
		panic(err)
	}
	return s.btik
}

func LoadPartMap(file string) (m map[string]string, err error) {
	bs, err := ioutil.ReadFile(file)
	if err != nil {
		return
	}
	err = json.Unmarshal(bs, &m)
	return
}

func LoadPartsDatabase(file string) (db PartsDatabase, err error) {
	bs, err := ioutil.ReadFile(file)
	if err != nil {
		return
	}
	err = json.Unmarshal(bs, &db)
	return
}
