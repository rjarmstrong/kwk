package system

import (
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
	"encoding/json"
	"os"
	"errors"
)

const (
	settingsBucketName = "Settings"
)

type Settings struct {
	Path string
	Db   *leveldb.DB
}

func GetCachePath() string {
	path := fmt.Sprintf("/Users/%s/Library/Caches/kwk", os.Getenv("USER"))
	if err := os.Mkdir(path, os.ModeDir); err != nil {
		if os.IsExist(err) {
			return path
		}
		panic(err)
	}
	return path
}

func NewSettings(dbName string) *Settings {
	path := fmt.Sprintf("%s/%s", GetCachePath(), dbName)
	return &Settings{Path:path}
}

func (s *Settings) open(){
	db, err := leveldb.OpenFile(s.Path, nil)
	if err != nil {
		panic("DB couldn't be opened :" + err.Error())
	}
	s.Db = db
}

func (s *Settings) Upsert(key string, value interface{}) {
	s.open()
	defer s.Close()
	str, _ := json.Marshal(value)
	if err := s.Db.Put([]byte(settingsBucketName + key), []byte(str), nil); err != nil {
		panic(err)
	}
}

func (s *Settings) Get(key string, value interface{}) error {
	s.open()
	defer s.Close()
	if v, err := s.Db.Get([]byte(settingsBucketName + key), nil); err != nil {
		if err.Error() == "leveldb: not found" {
			return errors.New("Not found.")
		}
		return err
	} else {
		err := json.Unmarshal(v, value)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Settings) Delete(key string) error {
	s.open()
	defer s.Close()
	return s.Db.Delete([]byte(settingsBucketName + key), nil);
}

func (s *Settings) Close() {
	s.Db.Close()
}