package main

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"log"
	"time"

	"github.com/boltdb/bolt"
)

//DB :
var DB *bolt.DB

var mappingsBucket = []byte("mappings")

//AddResponseData :
func AddResponseData(response ResponseData) (int, error) {
	var result int
	data, e := serialize(&response)

	if e != nil {
		return result, e
	}

	err := DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(mappingsBucket)
		sequance, err := put(b, data)

		if err == nil {
			result = int(sequance)
		}
		return err
	})

	return result, err
}

//GetResponseData :
func GetResponseData(id int) (ResponseData, error) {
	result := ResponseData{}

	err := DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(mappingsBucket)
		data := b.Get(itob(id))
		return deserialize(data, &result)
	})

	return result, err
}

func put(b *bolt.Bucket, data []byte) (int, error) {
	sequence, err := b.NextSequence()

	if err != nil {
		return 0, err
	}

	s := int(sequence)
	return s, b.Put(itob(s), data)
}

func serialize(r *ResponseData) ([]byte, error) {
	b := bytes.Buffer{}
	e := gob.NewEncoder(&b)

	err := e.Encode(r)
	return b.Bytes(), err
}

func deserialize(data []byte, response *ResponseData) error {
	b := bytes.Buffer{}
	b.Write(data)
	d := gob.NewDecoder(&b)
	return d.Decode(response)
}

func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

func init() {
	if db, err := bolt.Open("rest-mock.db", 0600, &bolt.Options{
		Timeout: 1 * time.Second,
	}); err == nil {
		DB = db
	} else {
		log.Fatalf("Failed to open connection to db\n%v", err)
	}

	err := DB.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists(mappingsBucket); err != nil {
			log.Fatalf("Failed to create mappings bucket\n%v", err)
		}

		return nil
	})

	if err != nil {
		log.Fatalf("Failed to create mappings bucket\n%v", err)
	}
}
