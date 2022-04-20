package boltdb

import (
	"errors"
	"github.com/aa-trsv/telegram-bot-otrs-builder/pkg/repository"
	"github.com/boltdb/bolt"
	"strconv"
)

type AccessRepository struct {
	db *bolt.DB
}

func NewAccessRepository(db *bolt.DB) *AccessRepository {
	return &AccessRepository{db: db}
}

func (r *AccessRepository) Save(chatID int64, access string, bucket repository.Bucket) error {
	return r.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		return b.Put(intToBytes(chatID), []byte(access))
	})
}

func (r *AccessRepository) Get(chatID int64, bucket repository.Bucket) (string, error) {
	var access string

	err := r.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		data := b.Get(intToBytes(chatID))
		access = string(data)
		return nil
	})
	if err != nil {
		return "", err
	}
	if access == "" {
		// TODO: Вынести ошибку в конфиг
		return "", errors.New("access denied")
	}

	return access, nil
}

func intToBytes(v int64) []byte {
	return []byte(strconv.FormatInt(v, 10))
}
