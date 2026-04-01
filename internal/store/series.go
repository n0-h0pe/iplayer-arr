package store

import (
	"encoding/json"

	bolt "go.etcd.io/bbolt"
)

func (s *Store) PutSeriesMapping(m *SeriesMapping) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		data, err := json.Marshal(m)
		if err != nil {
			return err
		}
		return tx.Bucket(bucketSeries).Put([]byte(m.TVDBId), data)
	})
}

func (s *Store) GetSeriesMapping(tvdbId string) (*SeriesMapping, error) {
	var m *SeriesMapping
	err := s.db.View(func(tx *bolt.Tx) error {
		data := tx.Bucket(bucketSeries).Get([]byte(tvdbId))
		if data == nil {
			return nil
		}
		m = &SeriesMapping{}
		return json.Unmarshal(data, m)
	})
	return m, err
}
