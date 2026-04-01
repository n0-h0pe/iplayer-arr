package store

import bolt "go.etcd.io/bbolt"

func (s *Store) GetConfig(key string) (string, error) {
	var val string
	err := s.db.View(func(tx *bolt.Tx) error {
		data := tx.Bucket(bucketConfig).Get([]byte(key))
		if data != nil {
			val = string(data)
		}
		return nil
	})
	return val, err
}

func (s *Store) SetConfig(key, value string) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(bucketConfig).Put([]byte(key), []byte(value))
	})
}
