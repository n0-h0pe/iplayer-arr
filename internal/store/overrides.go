package store

import (
	"encoding/json"
	"strings"

	bolt "go.etcd.io/bbolt"
)

func normaliseShowName(name string) string {
	return strings.ToLower(strings.TrimSpace(name))
}

func (s *Store) PutOverride(o *ShowOverride) error {
	key := normaliseShowName(o.ShowName)
	return s.db.Update(func(tx *bolt.Tx) error {
		data, err := json.Marshal(o)
		if err != nil {
			return err
		}
		return tx.Bucket(bucketOverrides).Put([]byte(key), data)
	})
}

func (s *Store) GetOverride(showName string) (*ShowOverride, error) {
	key := normaliseShowName(showName)
	var o *ShowOverride
	err := s.db.View(func(tx *bolt.Tx) error {
		data := tx.Bucket(bucketOverrides).Get([]byte(key))
		if data == nil {
			return nil
		}
		o = &ShowOverride{}
		return json.Unmarshal(data, o)
	})
	return o, err
}

func (s *Store) ListOverrides() ([]*ShowOverride, error) {
	var overrides []*ShowOverride
	err := s.db.View(func(tx *bolt.Tx) error {
		return tx.Bucket(bucketOverrides).ForEach(func(k, v []byte) error {
			var o ShowOverride
			if err := json.Unmarshal(v, &o); err != nil {
				return err
			}
			overrides = append(overrides, &o)
			return nil
		})
	})
	return overrides, err
}

func (s *Store) DeleteOverride(showName string) error {
	key := normaliseShowName(showName)
	return s.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(bucketOverrides).Delete([]byte(key))
	})
}
