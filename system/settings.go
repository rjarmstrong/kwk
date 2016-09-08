package system

import (
	"fmt"
	"github.com/boltdb/bolt"
	"encoding/json"
	"os"
)

const (
	settingsBucketName = "Settings"
)

type Settings struct {
	Path string
	Db   *bolt.DB
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
	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		panic("DB couldn't be opened :" + err.Error())
	}
	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(settingsBucketName))
		if err != nil {
			panic(err)
		}
		return nil
	})
	return &Settings{Db:db}
}

func (s *Settings) Upsert(key string, value interface{}) {
	s.Db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(settingsBucketName))
		b, err := json.Marshal(value)
		if err != nil {
			return err
		}
		err = bucket.Put([]byte(key), b)
		if err != nil {
			return err
		}
		return nil
	})
}

func (s *Settings) Get(key string, value interface{}) error {
	err := s.Db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(settingsBucketName))
		v := bucket.Get([]byte(key))
		if len(v) == 0 {
			return fmt.Errorf("Key '%s' does not exist.", key)
		}
		err := json.Unmarshal(v, value)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("%s", err)
	}
	return nil
}

func (s *Settings) Delete(key string) {
	s.Db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(settingsBucketName))
		if e := bucket.Delete([]byte(key)); e != nil {
			fmt.Println(e)
		}
		return nil
	})
}

func (s *Settings) Close() {
	s.Db.Close()
}