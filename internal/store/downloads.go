package store

import (
	"encoding/json"
	"fmt"

	bolt "go.etcd.io/bbolt"
)

func (s *Store) PutDownload(dl *Download) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		data, err := json.Marshal(dl)
		if err != nil {
			return fmt.Errorf("marshal download: %w", err)
		}
		return tx.Bucket(bucketDownloads).Put([]byte(dl.ID), data)
	})
}

func (s *Store) GetDownload(id string) (*Download, error) {
	var dl *Download
	err := s.db.View(func(tx *bolt.Tx) error {
		data := tx.Bucket(bucketDownloads).Get([]byte(id))
		if data == nil {
			return nil
		}
		dl = &Download{}
		return json.Unmarshal(data, dl)
	})
	return dl, err
}

func (s *Store) ListDownloads() ([]*Download, error) {
	var downloads []*Download
	err := s.db.View(func(tx *bolt.Tx) error {
		return tx.Bucket(bucketDownloads).ForEach(func(k, v []byte) error {
			var dl Download
			if err := json.Unmarshal(v, &dl); err != nil {
				return err
			}
			downloads = append(downloads, &dl)
			return nil
		})
	})
	return downloads, err
}

func (s *Store) DeleteDownload(id string) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(bucketDownloads).Delete([]byte(id))
	})
}

func (s *Store) FindDownloadByPIDQuality(pid, quality string) (*Download, error) {
	var found *Download
	err := s.db.View(func(tx *bolt.Tx) error {
		return tx.Bucket(bucketDownloads).ForEach(func(k, v []byte) error {
			var dl Download
			if err := json.Unmarshal(v, &dl); err != nil {
				return err
			}
			if dl.PID == pid && dl.Quality == quality {
				found = &dl
			}
			return nil
		})
	})
	return found, err
}

func (s *Store) MoveToHistory(id string) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		dlBucket := tx.Bucket(bucketDownloads)
		hBucket := tx.Bucket(bucketHistory)

		data := dlBucket.Get([]byte(id))
		if data == nil {
			return fmt.Errorf("download %s not found", id)
		}

		if err := hBucket.Put([]byte(id), data); err != nil {
			return err
		}
		return dlBucket.Delete([]byte(id))
	})
}
