package db

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"

	bolt "go.etcd.io/bbolt"
)

var (
	client *bolt.DB
)

func Init(dbpath string) (err error) {
	dirpath := filepath.Dir(dbpath)
	if _, err = os.Stat(dirpath); errors.Is(err, fs.ErrNotExist) {
		err = os.MkdirAll(dirpath, 0755)
		if err != nil {
			return
		}
	}

	client, err = bolt.Open(dbpath, 0600, nil)
	return
}

func PutPid(pid int, fullname string) (err error) {
	return client.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("pid"))
		if err != nil {
			return err
		}
		return b.Put([]byte(fullname), []byte(fmt.Sprint(pid)))
	})
}

func GetPid(fullname string) (pid int, ok bool) {
	err := client.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("pid"))
		if b == nil {
			return fmt.Errorf("bucket pid not found")
		}
		v := b.Get([]byte(fullname))
		if v == nil {
			return fmt.Errorf("pid not found")
		}

		pid, _ = strconv.Atoi(string(v))
		return nil
	})

	return pid, err == nil
}
