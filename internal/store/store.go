package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/uta8a/spzcdd/internal/domain"
	"go.etcd.io/bbolt"
)

var (
	ErrNotFound  = errors.New("not found")
	ErrConflict  = errors.New("conflict")
	ErrInvalid   = errors.New("invalid")
	ErrCorrupted = errors.New("corrupted")
)

const (
	schemaVersion = 1

	bucketMeta  = "meta"
	bucketTasks = "tasks"
	bucketSpecs = "specs"
	bucketExecs = "execs"

	metaKeySchemaVersion = "schema_version"
)

type Store struct {
	db *bbolt.DB
}

type OpenOptions struct {
	FileMode os.FileMode
	Timeout  time.Duration
}

func Open(path string, opt OpenOptions) (*Store, error) {
	if path == "" {
		return nil, fmt.Errorf("%w: path is required", ErrInvalid)
	}
	if opt.FileMode == 0 {
		opt.FileMode = 0o600
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, err
	}

	db, err := bbolt.Open(path, opt.FileMode, &bbolt.Options{Timeout: opt.Timeout})
	if err != nil {
		return nil, err
	}

	s := &Store{db: db}
	if err := s.initBucketsAndSchema(); err != nil {
		_ = db.Close()
		return nil, err
	}
	return s, nil
}

func (s *Store) Close() error {
	if s == nil || s.db == nil {
		return nil
	}
	return s.db.Close()
}

func (s *Store) initBucketsAndSchema() error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		meta, err := tx.CreateBucketIfNotExists([]byte(bucketMeta))
		if err != nil {
			return err
		}
		if _, err := tx.CreateBucketIfNotExists([]byte(bucketTasks)); err != nil {
			return err
		}
		if _, err := tx.CreateBucketIfNotExists([]byte(bucketSpecs)); err != nil {
			return err
		}
		if _, err := tx.CreateBucketIfNotExists([]byte(bucketExecs)); err != nil {
			return err
		}

		v := meta.Get([]byte(metaKeySchemaVersion))
		if v == nil {
			return meta.Put([]byte(metaKeySchemaVersion), []byte(fmt.Sprintf("%d", schemaVersion)))
		}
		if string(v) != fmt.Sprintf("%d", schemaVersion) {
			return fmt.Errorf("%w: unsupported schema version: %s", ErrCorrupted, string(v))
		}
		return nil
	})
}

func marshalJSON(v any) ([]byte, error) {
	return json.Marshal(v)
}

func unmarshalJSON(b []byte, v any) error {
	return json.Unmarshal(b, v)
}

func nowUTC() time.Time {
	return time.Now().UTC()
}

func taskKey(taskID string) []byte {
	return []byte("task/" + taskID)
}

func specKey(taskID string, rev int) []byte {
	return []byte(fmt.Sprintf("spec/%s/%010d", taskID, rev))
}

func execKey(taskID string, rev int) []byte {
	return []byte(fmt.Sprintf("exec/%s/%010d", taskID, rev))
}

func (s *Store) getTask(tx *bbolt.Tx, taskID string) (*domain.Task, error) {
	b := tx.Bucket([]byte(bucketTasks))
	v := b.Get(taskKey(taskID))
	if v == nil {
		return nil, ErrNotFound
	}
	var t domain.Task
	if err := unmarshalJSON(v, &t); err != nil {
		return nil, fmt.Errorf("%w: task json", ErrCorrupted)
	}
	return &t, nil
}
