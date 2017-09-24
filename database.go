package main

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"log"
	"time"

	"github.com/boltdb/bolt"
)

//DB :
var DB *bolt.DB

var mappingsBucket = []byte("mappings")

//AddResponseData :
func AddResponseData(item BoltDbItem) (int, error) {
	log.Println("Adding new item")
	var result int
	data, err := serialize(&item)

	if err != nil {
		return result, err
	}

	err = DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(mappingsBucket)
		sequance, err := put(b, data)
		log.Printf("Added data with sequence %v\n", sequance)

		if err == nil {
			result = int(sequance)
		}
		return err
	})

	return result, err
}

//GetBoltDbItem :
func GetBoltDbItem(id int) (*BoltDbItem, error) {
	log.Printf("Getting data with id %v\n", id)
	result := BoltDbItem{}

	err := DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(mappingsBucket)
		data := b.Get(itob(id))

		if data == nil {
			return fmt.Errorf("BoltDbItem with id %v not found", id)
		}
		return deserialize(data, &result)
	})

	return &result, err
}

//GetAll :
func GetAll() (map[int]*BoltDbItem, error) {
	log.Println("Getting all endpoints from DB")
	result := make(map[int]*BoltDbItem, 0)

	err := DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(mappingsBucket)
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			item := BoltDbItem{}
			if err := deserialize(v, &item); err != nil {
				return err
			}

			result[int(binary.BigEndian.Uint64(k))] = &item
		}

		return nil
	})

	return result, err
}

func put(b *bolt.Bucket, data []byte) (int, error) {
	sequence, err := b.NextSequence()

	if err != nil {
		return 0, err
	}

	key := int(sequence)
	return key, b.Put(itob(key), data)
}

func serialize(item *BoltDbItem) ([]byte, error) {
	b := bytes.Buffer{}
	e := gob.NewEncoder(&b)

	err := e.Encode(item)
	return b.Bytes(), err
}

func deserialize(data []byte, item *BoltDbItem) error {
	b := bytes.Buffer{}
	b.Write(data)
	d := gob.NewDecoder(&b)
	return d.Decode(item)
}

func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

func init() {
	if db, err := bolt.Open(fmt.Sprintf("rest-mock:%v.db", VERSION), 0600, &bolt.Options{
		Timeout: 1 * time.Second,
	}); err == nil {
		DB = db
	} else {
		log.Fatalf("Failed to open connection to db\n%v\n", err)
	}

	err := DB.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists(mappingsBucket); err != nil {
			log.Fatalf("Failed to create mappings bucket\n%v\n", err)
		}

		return nil
	})

	if err != nil {
		log.Fatalf("Failed to create mappings bucket\n%v\n", err)
	}
}
