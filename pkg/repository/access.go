package repository

type Bucket string

const AccessList Bucket = "access_list"

type AccessRepository interface {
	Save(chatID int64, access string, bucket Bucket) error
	Get(chatID int64, bucket Bucket) (string, error)
}
