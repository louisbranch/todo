package gob

import (
	"encoding/gob"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/luizbranco/todo"
)

type entry struct {
	Text    string
	Done    bool
	Entries []*entry
}

type Storage struct {
	file    string
	Entries []*entry
}

func Open(file string) (*Storage, error) {
	s := &Storage{file: file}
	flags := os.O_RDONLY
	f, err := os.OpenFile(file, flags, 0644)
	if os.IsNotExist(err) {
		return s, nil
	}
	defer f.Close()
	if err != nil {
		return s, err
	}
	dec := gob.NewDecoder(f)
	err = dec.Decode(s)
	if err != nil && err != io.EOF {
		return s, err
	}
	return s, nil
}

func (s *Storage) All() []todo.Task {
	return cpTasks(nil, s.Entries)
}

func cpTasks(ids []string, entries []*entry) []todo.Task {
	tasks := make([]todo.Task, len(entries))

	for i, e := range entries {
		ids := append([]string{}, ids...)
		id := strconv.Itoa(i + 1)
		ids = append(ids, id)
		t := todo.Task{
			ID:   strings.Join(ids, "."),
			Text: e.Text,
			Done: e.Done,
		}
		if len(e.Entries) > 0 {
			t.Tasks = cpTasks(ids, e.Entries)
		}
		tasks[i] = t
	}

	return tasks
}

func (s *Storage) Create(id string, text string) error {
	indexes, err := parseIDs(id)
	if err != nil {
		return err
	}

	if len(indexes) == 0 {
		s.Entries = append(s.Entries, &entry{Text: text})
		return s.save()
	}

	var e *entry
	entries := s.Entries

	for _, index := range indexes {
		if len(entries) <= index {
			return fmt.Errorf("parent task not found %s", id)
		}
		e = entries[index]
		entries = e.Entries
	}

	e.Entries = append(e.Entries, &entry{Text: text})
	return s.save()
}

func (s *Storage) Update(id string, text string) error {
	e, err := s.find(id)
	if err != nil {
		return err
	}

	e.Text = text
	return s.save()
}

func (s *Storage) Finish(id string) error {
	e, err := s.find(id)
	if err != nil {
		return err
	}

	entries := []*entry{e}

	for len(entries) > 0 {
		e, entries = entries[0], entries[1:]
		e.Done = true
		entries = append(entries, e.Entries...)
	}

	return s.save()
}

func (s *Storage) Delete(id string) error {
	indexes, err := parseIDs(id)
	if err != nil {
		return err
	}

	err = fmt.Errorf("task not found %s", id)

	if len(indexes) == 1 {
		i := indexes[0]
		if len(s.Entries) > i {
			copy(s.Entries[i:], s.Entries[i+1:])
			s.Entries[len(s.Entries)-1] = nil
			s.Entries = s.Entries[:len(s.Entries)-1]
			return s.save()
		}
		return err
	}

	var e *entry
	entries := s.Entries

	for _, index := range indexes[:len(indexes)-1] {
		if len(entries) <= index {
			return err
		}
		e = entries[index]
		entries = e.Entries
	}

	if e == nil {
		return err
	}

	i := indexes[len(indexes)-1]
	copy(e.Entries[i:], e.Entries[i+1:])
	e.Entries[len(e.Entries)-1] = nil
	e.Entries = e.Entries[:len(e.Entries)-1]

	if len(e.Entries) == 0 {
		e.Entries = nil
	}

	return s.save()
}

// Clean remove tasks that are done.
func (s *Storage) Clean() error {

	s.Entries = clean(s.Entries)

	return s.save()
}

func clean(entries []*entry) []*entry {
	var result []*entry

	for _, e := range entries {
		if e.Done {
			continue
		}
		e.Entries = clean(e.Entries)
		result = append(result, e)
	}

	return result
}

func (s *Storage) save() error {
	flags := os.O_CREATE | os.O_TRUNC | os.O_WRONLY
	f, err := os.OpenFile(s.file, flags, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := gob.NewEncoder(f)
	return enc.Encode(s)
}

func (s *Storage) find(id string) (*entry, error) {
	indexes, err := parseIDs(id)
	if err != nil {
		return nil, err
	}

	err = fmt.Errorf("task not found %s", id)

	var e *entry
	entries := s.Entries

	for _, index := range indexes {
		if len(entries) <= index {
			return nil, err
		}
		e = entries[index]
		entries = e.Entries
	}

	if e == nil {
		return nil, err
	}

	return e, nil
}

func parseIDs(token string) ([]int, error) {
	if token == "" {
		return nil, nil
	}

	var nums []int

	tokens := strings.Split(token, ".")

	for _, t := range tokens {
		n, err := strconv.Atoi(t)
		if err != nil {
			return nil, fmt.Errorf("invalid id %s", token)
		}
		if n < 1 {
			return nil, fmt.Errorf("invalid id %s", token)
		}
		nums = append(nums, n-1)
	}

	return nums, nil
}
