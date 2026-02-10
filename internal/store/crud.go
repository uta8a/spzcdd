package store

import (
	"fmt"
	"sort"

	"github.com/uta8a/spzcdd/internal/domain"
	"go.etcd.io/bbolt"
)

// CreateTask creates a new task and assigns a monotonically increasing TaskNumber.
func (s *Store) CreateTask(title, draftBody string) (*domain.Task, error) {
	if s == nil || s.db == nil {
		return nil, fmt.Errorf("%w: store is closed", ErrInvalid)
	}
	now := nowUTC()
	returnTask := &domain.Task{}

	err := s.db.Update(func(tx *bbolt.Tx) error {
		meta := tx.Bucket([]byte(bucketMeta))
		if meta == nil {
			return fmt.Errorf("%w: meta bucket missing", ErrCorrupted)
		}
		seq, err := meta.NextSequence()
		if err != nil {
			return err
		}

		t, err := domain.NewTask(now, title, draftBody)
		if err != nil {
			return fmt.Errorf("%w: %v", ErrInvalid, err)
		}
		t.Number = domain.TaskNumber(seq)

		b := tx.Bucket([]byte(bucketTasks))
		if b == nil {
			return fmt.Errorf("%w: tasks bucket missing", ErrCorrupted)
		}
		k := taskKey(t.ID)
		if b.Get(k) != nil {
			return fmt.Errorf("%w: task id already exists", ErrConflict)
		}
		buf, err := marshalJSON(t)
		if err != nil {
			return err
		}
		if err := b.Put(k, buf); err != nil {
			return err
		}
		*returnTask = *t
		return nil
	})
	if err != nil {
		return nil, err
	}
	return returnTask, nil
}

// UpdateDraft updates draft body only when task is in DRAFT.
func (s *Store) UpdateDraft(taskID, draftBody string) (*domain.Task, error) {
	if s == nil || s.db == nil {
		return nil, fmt.Errorf("%w: store is closed", ErrInvalid)
	}
	now := nowUTC()
	var out *domain.Task

	err := s.db.Update(func(tx *bbolt.Tx) error {
		t, err := s.getTask(tx, taskID)
		if err != nil {
			return err
		}
		if err := t.SetDraftBody(now, draftBody); err != nil {
			return fmt.Errorf("%w: %v", ErrConflict, err)
		}
		b := tx.Bucket([]byte(bucketTasks))
		buf, err := marshalJSON(t)
		if err != nil {
			return err
		}
		if err := b.Put(taskKey(taskID), buf); err != nil {
			return err
		}
		out = t
		return nil
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SetState sets the task state.
func (s *Store) SetState(taskID string, state domain.TaskState) (*domain.Task, error) {
	if s == nil || s.db == nil {
		return nil, fmt.Errorf("%w: store is closed", ErrInvalid)
	}
	if !state.IsValid() {
		return nil, fmt.Errorf("%w: invalid state: %q", ErrInvalid, state)
	}
	now := nowUTC()
	var out *domain.Task

	err := s.db.Update(func(tx *bbolt.Tx) error {
		t, err := s.getTask(tx, taskID)
		if err != nil {
			return err
		}
		t.State = state
		t.Touch(now)
		b := tx.Bucket([]byte(bucketTasks))
		buf, err := marshalJSON(t)
		if err != nil {
			return err
		}
		if err := b.Put(taskKey(taskID), buf); err != nil {
			return err
		}
		out = t
		return nil
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SetError sets task state to ERROR and stores the error message.
func (s *Store) SetError(taskID, errMsg string) (*domain.Task, error) {
	if s == nil || s.db == nil {
		return nil, fmt.Errorf("%w: store is closed", ErrInvalid)
	}
	now := nowUTC()
	var out *domain.Task

	err := s.db.Update(func(tx *bbolt.Tx) error {
		t, err := s.getTask(tx, taskID)
		if err != nil {
			return err
		}
		t.State = domain.StateError
		t.LastError = errMsg
		t.Touch(now)
		b := tx.Bucket([]byte(bucketTasks))
		buf, err := marshalJSON(t)
		if err != nil {
			return err
		}
		if err := b.Put(taskKey(taskID), buf); err != nil {
			return err
		}
		out = t
		return nil
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (s *Store) GetTask(taskID string) (*domain.Task, error) {
	if s == nil || s.db == nil {
		return nil, fmt.Errorf("%w: store is closed", ErrInvalid)
	}
	var out *domain.Task
	err := s.db.View(func(tx *bbolt.Tx) error {
		t, err := s.getTask(tx, taskID)
		if err != nil {
			return err
		}
		out = t
		return nil
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (s *Store) ListTasks() ([]domain.Task, error) {
	if s == nil || s.db == nil {
		return nil, fmt.Errorf("%w: store is closed", ErrInvalid)
	}
	var tasks []domain.Task
	err := s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(bucketTasks))
		if b == nil {
			return fmt.Errorf("%w: tasks bucket missing", ErrCorrupted)
		}
		return b.ForEach(func(k, v []byte) error {
			var t domain.Task
			if err := unmarshalJSON(v, &t); err != nil {
				return fmt.Errorf("%w: task json", ErrCorrupted)
			}
			tasks = append(tasks, t)
			return nil
		})
	})
	if err != nil {
		return nil, err
	}
	sort.Slice(tasks, func(i, j int) bool { return tasks[i].Number < tasks[j].Number })
	return tasks, nil
}

// PutSpec stores a spec (immutable per (task_id, rev)).
func (s *Store) PutSpec(spec *domain.Spec) error {
	if s == nil || s.db == nil {
		return fmt.Errorf("%w: store is closed", ErrInvalid)
	}
	if spec == nil {
		return fmt.Errorf("%w: spec is required", ErrInvalid)
	}
	if spec.TaskID == "" || spec.Rev <= 0 {
		return fmt.Errorf("%w: invalid spec", ErrInvalid)
	}
	if !spec.Status.IsValid() {
		return fmt.Errorf("%w: invalid spec status: %q", ErrInvalid, spec.Status)
	}

	return s.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(bucketSpecs))
		if b == nil {
			return fmt.Errorf("%w: specs bucket missing", ErrCorrupted)
		}
		k := specKey(spec.TaskID, spec.Rev)
		if b.Get(k) != nil {
			return fmt.Errorf("%w: spec already exists", ErrConflict)
		}
		buf, err := marshalJSON(spec)
		if err != nil {
			return err
		}
		return b.Put(k, buf)
	})
}

func (s *Store) GetSpec(taskID string, rev int) (*domain.Spec, error) {
	if s == nil || s.db == nil {
		return nil, fmt.Errorf("%w: store is closed", ErrInvalid)
	}
	if taskID == "" || rev <= 0 {
		return nil, fmt.Errorf("%w: invalid spec key", ErrInvalid)
	}
	var out *domain.Spec
	err := s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(bucketSpecs))
		if b == nil {
			return fmt.Errorf("%w: specs bucket missing", ErrCorrupted)
		}
		v := b.Get(specKey(taskID, rev))
		if v == nil {
			return ErrNotFound
		}
		var spec domain.Spec
		if err := unmarshalJSON(v, &spec); err != nil {
			return fmt.Errorf("%w: spec json", ErrCorrupted)
		}
		out = &spec
		return nil
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (s *Store) GetCurrentSpec(taskID string) (*domain.Spec, error) {
	t, err := s.GetTask(taskID)
	if err != nil {
		return nil, err
	}
	if t.CurrentSpecRev <= 0 {
		return nil, ErrNotFound
	}
	return s.GetSpec(taskID, t.CurrentSpecRev)
}

// PutExecution stores an execution (immutable per (task_id, rev)).
func (s *Store) PutExecution(exec *domain.Execution) error {
	if s == nil || s.db == nil {
		return fmt.Errorf("%w: store is closed", ErrInvalid)
	}
	if exec == nil {
		return fmt.Errorf("%w: execution is required", ErrInvalid)
	}
	if exec.TaskID == "" || exec.Rev <= 0 {
		return fmt.Errorf("%w: invalid execution", ErrInvalid)
	}
	if !exec.Status.IsValid() {
		return fmt.Errorf("%w: invalid execution status: %q", ErrInvalid, exec.Status)
	}

	return s.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(bucketExecs))
		if b == nil {
			return fmt.Errorf("%w: execs bucket missing", ErrCorrupted)
		}
		k := execKey(exec.TaskID, exec.Rev)
		if b.Get(k) != nil {
			return fmt.Errorf("%w: execution already exists", ErrConflict)
		}
		buf, err := marshalJSON(exec)
		if err != nil {
			return err
		}
		return b.Put(k, buf)
	})
}

func (s *Store) GetExecution(taskID string, rev int) (*domain.Execution, error) {
	if s == nil || s.db == nil {
		return nil, fmt.Errorf("%w: store is closed", ErrInvalid)
	}
	if taskID == "" || rev <= 0 {
		return nil, fmt.Errorf("%w: invalid execution key", ErrInvalid)
	}
	var out *domain.Execution
	err := s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(bucketExecs))
		if b == nil {
			return fmt.Errorf("%w: execs bucket missing", ErrCorrupted)
		}
		v := b.Get(execKey(taskID, rev))
		if v == nil {
			return ErrNotFound
		}
		var exec domain.Execution
		if err := unmarshalJSON(v, &exec); err != nil {
			return fmt.Errorf("%w: execution json", ErrCorrupted)
		}
		out = &exec
		return nil
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (s *Store) GetCurrentExecution(taskID string) (*domain.Execution, error) {
	t, err := s.GetTask(taskID)
	if err != nil {
		return nil, err
	}
	if t.CurrentExecRev <= 0 {
		return nil, ErrNotFound
	}
	return s.GetExecution(taskID, t.CurrentExecRev)
}
