package gob

import (
	"io/ioutil"
	"log"
	"reflect"
	"testing"

	"github.com/luizbranco/todo"
)

type tcase struct {
	parent, text string
}

func TestAll(t *testing.T) {
	s := new(t, []tcase{
		{"", "1"},
		{"1", "1-1"},
		{"1", "1-2"},
		{"1", "1-3"},
		{"1.2", "1-2-1"},
		{"1.2", "1-2-2"},
		{"", "2"},
		{"2", "2-1"},
		{"2.1", "2-1-1"},
	})

	want := []todo.Task{
		todo.Task{
			ID:   "1",
			Text: "1",
			Tasks: []todo.Task{
				todo.Task{ID: "1.1", Text: "1-1"},
				todo.Task{
					ID:   "1.2",
					Text: "1-2",
					Tasks: []todo.Task{
						todo.Task{ID: "1.2.1", Text: "1-2-1"},
						todo.Task{ID: "1.2.2", Text: "1-2-2"},
					},
				},
				todo.Task{ID: "1.3", Text: "1-3"},
			},
		},
		todo.Task{
			ID:   "2",
			Text: "2",
			Tasks: []todo.Task{
				todo.Task{
					ID:   "2.1",
					Text: "2-1",
					Tasks: []todo.Task{
						todo.Task{ID: "2.1.1", Text: "2-1-1"},
					},
				},
			},
		},
	}

	got := s.All()

	if !reflect.DeepEqual(want, got) {
		t.Errorf("want %v, got %v", want, got)
	}
}

func TestCreate(t *testing.T) {
	s := new(t, []tcase{
		{"", "1"},
		{"1", "1.1"},
		{"1", "1.2"},
		{"1", "1.3"},
		{"1.2", "1.2.1"},
		{"1.2", "1.2.2"},
		{"", "2"},
		{"2", "2.1"},
		{"2.1", "2.1.1"},
	})

	want := []*entry{
		&entry{
			Text: "1",
			Entries: []*entry{
				&entry{Text: "1.1"},
				&entry{
					Text: "1.2",
					Entries: []*entry{
						&entry{Text: "1.2.1"},
						&entry{Text: "1.2.2"},
					},
				},
				&entry{Text: "1.3"},
			},
		},
		&entry{
			Text: "2",
			Entries: []*entry{
				&entry{
					Text: "2.1",
					Entries: []*entry{
						&entry{Text: "2.1.1"},
					},
				},
			},
		},
	}

	errors := []string{
		"11",      // does not exist
		"1.1.1.1", // does not exist
		"1a",      // invalid id
	}

	for _, e := range errors {
		err := s.Create(e, "")
		if err == nil {
			t.Error("want error, got nothing")
		}
	}

	got := s.Entries

	if !reflect.DeepEqual(want, got) {
		t.Errorf("want %v, got %v", want, got)
	}
}

func TestUpdate(t *testing.T) {
	s := new(t, []tcase{
		{"", "1"},
		{"1", "1.1"},
	})

	cases := []struct {
		id   string
		text string
		err  bool
	}{
		{"1", "First", false},
		{"1.1", "Nested", false},
		{"2", "", true},
	}

	want := []*entry{
		&entry{
			Text: "First",
			Entries: []*entry{
				&entry{Text: "Nested"},
			},
		},
	}

	for _, c := range cases {
		err := s.Update(c.id, c.text)
		if c.err && err == nil {
			t.Error("want error, got nothing")
		}
	}

	got := s.Entries

	if !reflect.DeepEqual(want, got) {
		t.Errorf("want %v, got %v", want, got)
	}
}

func TestFinish(t *testing.T) {
	s := new(t, []tcase{
		{"", "1"},
		{"1", "1.1"},
		{"1.1", "1.1.1"},
		{"", "2"},
		{"2", "2.1"},
	})

	cases := []struct {
		id  string
		err bool
	}{
		{"1", false},
		{"2.1", false},
		{"3", true},
	}

	want := []*entry{
		&entry{
			Text: "1",
			Done: true,
			Entries: []*entry{
				&entry{
					Text: "1.1",
					Done: true,
					Entries: []*entry{
						&entry{
							Text: "1.1.1",
							Done: true,
						},
					},
				},
			},
		},
		&entry{
			Text: "2",
			Done: false,
			Entries: []*entry{
				&entry{
					Text: "2.1",
					Done: true,
				},
			},
		},
	}

	for _, c := range cases {
		err := s.Finish(c.id)
		if c.err && err == nil {
			t.Error("want error, got nothing")
		}
	}

	got := s.Entries

	if !reflect.DeepEqual(want, got) {
		t.Errorf("want %v, got %v", want, got)
	}
}

func TestDelete(t *testing.T) {
	s := new(t, []tcase{
		{"", "1"},
		{"1", "1.1"},
		{"1.1", "1.1.1"},
		{"", "2"},
		{"2", "2.1"},
	})

	cases := []struct {
		id  string
		err bool
	}{
		{"2.1", false},
		{"1", false},
		{"3", true},
	}

	want := []*entry{
		&entry{
			Text: "2",
		},
	}

	for _, c := range cases {
		err := s.Delete(c.id)
		if c.err && err == nil {
			t.Error("want error, got nothing")
		}
	}

	got := s.Entries

	if !reflect.DeepEqual(want, got) {
		t.Errorf("want %v, got %v", want, got)
	}
}

func TestClean(t *testing.T) {
	s := new(t, []tcase{
		{"", "1"},
		{"1", "1-1"},
		{"1", "1-2"},
		{"1", "1-3"},
		{"1.2", "1-2-1"},
		{"1.2", "1-2-2"},
		{"", "2"},
		{"2", "2-1"},
		{"2.1", "2-1-1"},
	})

	want := []todo.Task{
		todo.Task{
			ID:   "1",
			Text: "1",
			Tasks: []todo.Task{
				todo.Task{ID: "1.1", Text: "1-1"},
				todo.Task{
					ID:   "1.2",
					Text: "1-2",
					Tasks: []todo.Task{
						todo.Task{ID: "1.2.1", Text: "1-2-2"},
					},
				},
				todo.Task{ID: "1.3", Text: "1-3"},
			},
		},
		todo.Task{
			ID:   "2",
			Text: "2",
		},
	}

	s.Finish("1.2.1")
	s.Finish("2.1")

	err := s.Clean()
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}

	got := s.All()

	if !reflect.DeepEqual(want, got) {
		t.Errorf("want\n%v got\n%v", want, got)
	}
}

func new(t *testing.T, cases []tcase) *Storage {
	f, err := ioutil.TempFile("", "todo")
	if err != nil {
		log.Fatalf("failed to create file %s", f.Name())
	}

	file := f.Name()
	s := &Storage{file: file}

	for _, c := range cases {
		err := s.Create(c.parent, c.text)
		if err != nil {
			t.Errorf("got error %s", err)
		}
	}

	return s
}
