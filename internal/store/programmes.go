package store

import (
	"encoding/json"
	"time"

	bolt "go.etcd.io/bbolt"
)

func (s *Store) PutProgramme(p *Programme) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		data, err := json.Marshal(p)
		if err != nil {
			return err
		}
		return tx.Bucket(bucketProgrammes).Put([]byte(p.PID), data)
	})
}

func (s *Store) GetProgramme(pid string) (*Programme, error) {
	var p *Programme
	err := s.db.View(func(tx *bolt.Tx) error {
		data := tx.Bucket(bucketProgrammes).Get([]byte(pid))
		if data == nil {
			return nil
		}
		p = &Programme{}
		return json.Unmarshal(data, p)
	})
	return p, err
}

func (s *Store) DeleteProgramme(pid string) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(bucketProgrammes).Delete([]byte(pid))
	})
}

func (s *Store) PurgeStaleProgrammes(maxAge time.Duration) error {
	cutoff := time.Now().Add(-maxAge)
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketProgrammes)
		var stale [][]byte
		b.ForEach(func(k, v []byte) error {
			var p Programme
			if json.Unmarshal(v, &p) == nil && p.CachedAt.Before(cutoff) {
				stale = append(stale, k)
			}
			return nil
		})
		for _, k := range stale {
			b.Delete(k)
		}
		return nil
	})
}
