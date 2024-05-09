package database

import (
	"encoding/json"
	"errors"
	"path/filepath"

	"sub2clash/model"

	"go.etcd.io/bbolt"
)

var DB *bbolt.DB

func ConnectDB() error {
	// 确保数据目录存在
	path := filepath.Join("data", "sub2clash.db")
	// 打开或创建数据库文件
	db, err := bbolt.Open(path, 0600, nil)
	if err != nil {
		return err
	}
	DB = db

	// 确保存储桶存在
	return db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("ShortLinks"))
		return err
	})
}

func FindShortLinkByHash(hash string) (*model.ShortLink, error) {
	var shortLink model.ShortLink
	err := DB.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("ShortLinks"))
		v := b.Get([]byte(hash))
		if v == nil {
			return errors.New("ShortLink not found")
		}
		return json.Unmarshal(v, &shortLink)
	})
	if err != nil {
		return nil, err
	}
	return &shortLink, nil
}

func SaveShortLink(shortLink *model.ShortLink) error {
	return DB.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("ShortLinks"))
		encoded, err := json.Marshal(shortLink)
		if err != nil {
			return err
		}
		return b.Put([]byte(shortLink.Hash), encoded)
	})
}

func CheckShortLinkHashExists(hash string) (bool, error) {
	exists := false
	err := DB.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("ShortLinks"))
		v := b.Get([]byte(hash))
		exists = v != nil
		return nil
	})
	if err != nil {
		return false, err
	}
	return exists, nil
}
