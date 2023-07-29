package kvdb

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"runtime"
	"strings"
	"time"

	"github.com/dgraph-io/badger/v3"
)

var (
	BadgerDB *badger.DB
	Errdb    error
)

func init() {
	path, _ := GetCurrentPath()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, 0755)
	}
	opts := badger.DefaultOptions(path)
	opts.Dir = path
	opts.ValueDir = path
	opts.SyncWrites = false
	opts.ValueThreshold = 256
	opts.CompactL0OnClose = true
	BadgerDB, Errdb = badger.Open(opts)
	if Errdb != nil {
		fmt.Println(Errdb.Error())
	}
}

func Set(key []byte, value []byte) {
	wb := BadgerDB.NewWriteBatch()
	defer wb.Cancel()
	err := wb.SetEntry(badger.NewEntry(key, value).WithMeta(0))
	if err != nil {
		log.Println("Failed to write data to cache.", "key", string(key), "value", string(value), "err", err)
	}
	err = wb.Flush()
	if err != nil {
		log.Println("Failed to flush data to cache.", "key", string(key), "value", string(value), "err", err)
	}
}

func SetWithTTL(key []byte, value []byte, ttl int64) {
	wb := BadgerDB.NewWriteBatch()
	defer wb.Cancel()
	err := wb.SetEntry(badger.NewEntry(key, value).WithMeta(0).WithTTL(time.Duration(ttl * time.Second.Nanoseconds())))
	if err != nil {
		log.Println("Failed to write data to cache.", "key", string(key), "value", string(value), "err", err)
	}
	err = wb.Flush()
	if err != nil {
		log.Println("Failed to flush data to cache.", "key", string(key), "value", string(value), "err", err)
	}
}

func Get(key []byte) (string, error) {
	var ival []byte
	err := BadgerDB.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err
		}
		ival, err = item.ValueCopy(nil)
		return err
	})
	if err != nil {
		return "", err
	}
	return string(ival), nil
}

func GetToken(key []byte) ([]byte, error) {
	var ival []byte
	err := BadgerDB.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err
		}
		ival, err = item.ValueCopy(nil)
		return err
	})
	if err != nil {
		return []byte(""), err
	}
	return ival, nil
}

func Has(key []byte) (bool, error) {
	var exist bool = false
	err := BadgerDB.View(func(txn *badger.Txn) error {
		_, err := txn.Get(key)
		if err != nil {
			return err
		} else {
			exist = true
		}
		return err
	})
	// align with leveldb, if the key doesn't exist, leveldb returns nil
	if strings.HasSuffix(err.Error(), "not found") {
		err = nil
	}
	return exist, err
}

func Delete(key []byte) error {
	err := BadgerDB.Update(func(txn *badger.Txn) error {
		return txn.Delete(key)
	})
	return err
}

func IteratorKeysAndValues() {

	err := BadgerDB.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			k := item.Key()
			err := item.Value(func(v []byte) error {
				fmt.Printf("key=%s, value=%s\n", k, v)
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		log.Println("Failed to iterator keys and values from the cache.", "error", err)
	}
}

func IteratorKeys() {
	err := BadgerDB.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			k := item.Key()
			fmt.Printf("key=%s\n", k)
		}
		return nil
	})

	if err != nil {
		log.Println("Failed to iterator keys from the cache.", "error", err)
	}
}

func SeekWithPrefix(prefixStr string) {
	err := BadgerDB.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		prefix := []byte(prefixStr)
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			k := item.Key()
			err := item.Value(func(v []byte) error {
				fmt.Printf("key=%s, value=%s\n", k, v)
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		log.Println("Failed to seek prefix from the cache.", "prefix", prefixStr, "error", err)
	}
}

func GetCurrentPath() (string, error) {
	OS := runtime.GOOS
	LinkPathStr := "/"
	var dbPath string
	if OS == "windows" {
		LinkPathStr = "\\"
	}
	u, _ := user.Current()
	homeDir := u.HomeDir

	if OS == "windows" {
		homeDir = homeWindows()
	}
	if OS == "linux" {
		dbPath = strings.Join([]string{homeDir, ".cache", "FleaMarketManager"}, LinkPathStr)
	}
	if OS == "windows" {
		dbPath = strings.Join([]string{homeDir, ".FleaMarketManager"}, LinkPathStr)
	}
	return dbPath, nil
}

func homeWindows() string {
	drive := os.Getenv("HOMEDRIVE")
	path := os.Getenv("HOMEPATH")
	home := drive + path
	if drive == "" || path == "" {
		home = os.Getenv("USERPROFILE")
	}
	if home == "" {
		return ""
	}
	return home
}
